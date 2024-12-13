package main

import (
	"fmt"       // 用于格式化输入输出
	"io/ioutil" // 用于读取 HTTP 响应体
	"net/http"  // 用于发起 HTTP 请求
	"sync"      // 用于同步等待多个 Goroutine 完成
)

func main() {
	// 创建一个 WaitGroup 来等待所有 Goroutine 完成
	var wg sync.WaitGroup

	// 定义我们要发起的并发请求数
	requestCount := 1000

	// 循环发起 1000 个并发请求
	for i := 0; i < requestCount; i++ {
		wg.Add(1) // 每次启动 Goroutine 前，增加 WaitGroup 的计数

		// 启动 Goroutine 并发执行 HTTP 请求
		go func(i int) {
			defer wg.Done() // 当此 Goroutine 完成时，减少 WaitGroup 计数

			// 构建 API 请求的 URL
			url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", i)

			// 发起 HTTP GET 请求
			resp, err := http.Get(url)
			if err != nil {
				// 如果请求出错，输出错误信息并返回
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close() // 确保响应体在处理完后关闭

			// 读取响应体的内容
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}

			// 打印当前请求的结果（Goroutine 编号和响应内容）
			fmt.Printf("Task %d Response: %s\n", i, body)
		}(i) // 传递 i 作为 Goroutine 的参数，以避免闭包问题
	}

	// 等待所有 Goroutine 完成
	wg.Wait()
}
