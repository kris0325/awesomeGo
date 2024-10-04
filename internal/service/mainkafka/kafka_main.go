package main

import (
	kafka2 "awesomeGo/internal/service/kafka"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// Kafka broker地址和topic名称
	//broker := "localhost:9092"
	broker := "127.0.0.1:9092"
	topic := "test-topic"

	// 初始化 Kafka 生产者和消费者
	producer := kafka2.InitProducer(broker, topic)
	consumer := kafka2.InitConsumer(broker, topic)

	fmt.Printf("kafka_main.go start success: %s: %s\n", broker, "kafka_main.go")

	// 等待信号以优雅关闭程序
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// 关闭 Kafka 生产者和消费者
	producer.Close()
	consumer.Close()
}
