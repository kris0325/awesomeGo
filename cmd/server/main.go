package main

import (
	"awesomeGo/internal/config"
	"awesomeGo/internal/routes"
	"log"
)

func main() {
	// Initialize DB connection
	config.InitDB()

	// Create a Gin router
	r := routes.InitRoutes()

	// Start the server
	log.Println("************** Server successfully  running on port 8080 **********")
	r.Run(":8080")
}

//func main() {
//// Kafka broker地址，根据您的本地环境配置
////broker := "localhost:9092"
//broker := "127.0.0.1:9092"
//
//// Kafka topic名称
//topic := "test-topic"
//
//// Kafka 消费者组ID
//consumerGroupID := "test-consumer-group"
//
//// 创建Kafka消费者配置
//consumerConfig := kafka.ConfigMap{
//	"bootstrap.servers": broker,
//	"group.id":          consumerGroupID,
//	"auto.offset.reset": "earliest",
//}
//
//// 创建Kafka消费者
//consumer, err := kafka.NewConsumer(&consumerConfig)
//if err != nil {
//	fmt.Printf("Failed to create consumer: %s\n", err)
//	os.Exit(1)
//}
//
//// 订阅主题
//err = consumer.SubscribeTopics([]string{topic}, nil)
//if err != nil {
//	fmt.Printf("Error subscribing to topic: %v\n", err)
//	os.Exit(1)
//}
//
//// 捕获中断信号，优雅地关闭消费者
//sigchan := make(chan os.Signal, 1)
//signal.Notify(sigchan, os.Interrupt)
//
//run := true
//for run {
//	select {
//	case sig := <-sigchan:
//		fmt.Printf("Caught signal %v: terminating\n", sig)
//		run = false
//	default:
//		// 从Kafka消费消息
//		msg, err := consumer.ReadMessage(-1)
//		if err == nil {
//			fmt.Printf("main.go consumer **** Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
//		} else {
//			fmt.Printf("Error reading message: %v\n", err)
//		}
//	}
//}
//
//// 关闭消费者
//err = consumer.Close()
//if err != nil {
//	return
//}
//}
