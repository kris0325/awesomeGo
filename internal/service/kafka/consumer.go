package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// 初始化 Kafka 消费者
func InitConsumer(broker, topic string) *kafka.Consumer {
	cConfig := kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           "test-group-kris-222",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	cons, err := kafka.NewConsumer(&cConfig)
	if err != nil {
		fmt.Printf("Failed to create Kafka consumer : %s\n", err)
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	go func() {
		for {
			err := consumeMessage(cons, topic)
			if err != nil {
				log.Printf("Consumer: Error consuming message: %v", err)
				fmt.Printf("Error consuming message : %s\n", err)
			}
		}
	}()

	return cons
}

// 消费消息从 Kafka
func consumeMessage(consumer *kafka.Consumer, topic string) error {
	msg, err := consumer.ReadMessage(-1)
	if err != nil {
		return err
	}
	log.Printf("Consumer received message: %s\n", string(msg.Value))

	// 成功处理消息后手动提交偏移量
	_, err = consumer.CommitMessage(msg)
	if err != nil {
		log.Printf("Error committing offset: %v", err)
		return err
	}

	return nil
}
