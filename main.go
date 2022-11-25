package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

const (
	DefaultPort         = "3333"
	DefaultKafkaAddress = "localhost:9092"
	DefaultTopic        = "my-topic"
)

type config struct {
	KafkaAddress string
	Topic        string
}

func main() {
	log.Println("starting server ... 1")

	config := getConfigFromEnv()

	msgCh := make(chan []byte)
	pod := os.Getenv("POD_NAME")
	go readLoop(strings.Split(config.KafkaAddress, ","), config.Topic, func(m kafka.Message) {
		log.Printf("[%s] Sleep 10 seconds to simulate delay...", pod)
		time.Sleep(10 * time.Second)
		msg := fmt.Sprintf("message offset:'%d' key:'%s' value:'%s'", m.Offset, string(m.Key), string(m.Value))
		msgCh <- []byte(msg)
		//log.Printf("Wake up...")
	})
	for {
		select {
		case msg := <-msgCh:
			log.Printf("[%s] %s\n", pod, string(msg))
		}
	}
}

func getConfigFromEnv() config {
	config := config{}
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

func readLoop(brokers []string, topic string, consumer func(m kafka.Message)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "consumer-group-id",
		Topic:   topic,
	})

	ctx := context.Background()
	for {
		m, err := r.ReadMessage(ctx)

		if err != nil {
			log.Printf("%s", err)
			break
		}

		// process message
		consumer(m)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
