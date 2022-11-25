package main

import (
	"log"
	"os"

	kafka "github.com/segmentio/kafka-go"
)

const (
	DefaultPort         = "3333"
	DefaultKafkaAddress = "localhost:9092"
	DefaultTopic        = "my-topic"
)

type config struct {
	Port         string
	KafkaAddress string
	Topic        string
}

func main() {
	log.Println("starting server ...")

	config := getConfigFromEnv()

	_ = &kafka.Writer{
		Addr:     kafka.TCP(config.KafkaAddress),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

}

func getConfigFromEnv() config {
	config := config{}
	config.Port = os.Getenv("PORT")
	if len(config.Port) <= 0 {
		config.Port = DefaultPort
	}
	config.KafkaAddress = os.Getenv("KAFKA")
	if len(config.KafkaAddress) <= 0 {
		config.KafkaAddress = DefaultKafkaAddress
	}
	config.Topic = os.Getenv("TOPIC")
	if len(config.Topic) <= 0 {
		config.Topic = DefaultTopic
	}
	return config
}
