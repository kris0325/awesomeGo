package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// 初始化 Kafka 生产者
func InitProducer(broker, topic string) *kafka.Producer {
	pConfig := kafka.ConfigMap{
		"bootstrap.servers": broker,
	}

	prod, err := kafka.NewProducer(&pConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	go func() {
		for {
			value := fmt.Sprintf("producer.go start produce Message from  ********* at %v", time.Now().Format(time.StampMilli))
			err := produceMessage(prod, topic, value)
			if err != nil {
				log.Printf("Producer: Error producing message: %v", err)
			}
			fmt.Printf("producer.go finished produce : %s: %s\n, at %v\"", broker, "producer.go", time.Now().Format(time.StampMilli))
			time.Sleep(time.Second) // 生产者间隔时间
		}
	}()

	return prod
}

// 生产消息到 Kafka
func produceMessage(producer *kafka.Producer, topic, value string) error {
	deliveryChan := make(chan kafka.Event)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}
