package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//r.Use(logPingRequest) // 将中间件函数添加到所有路由
	r.GET("/ping", func(c *gin.Context) {
		//fmt.Println("call ping")
		logPingRequest()
		c.JSON(200, gin.H{
			"message": "hello world Go!",
		})
	})

	err := r.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	} // listen and serve on 0.0.0.0:8080
}

//func logPingRequest(c *gin.Context) {
//	// 在处理请求之前打印日志
//	fmt.Println("call ping222")
//	c.Next() // 继续处理请求
//}
