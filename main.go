package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
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
	log.Println("starting server ...")

	config := getConfigFromEnv()

	_ = &kafka.Writer{
		Addr:     kafka.TCP(config.KafkaAddress),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	msgCh := make(chan []byte)
	go readLoop(strings.Split(config.KafkaAddress, ","), config.Topic, func(m kafka.Message) {
		// go func() {
		log.Printf("Sleep 10 seconds to simulate delay...")
		time.Sleep(10 * time.Second)
		msgCh <- m.Value
		log.Printf("Wake up...")
		// }()
	})
	for {
		select {
		case msg := <-msgCh:
			fmt.Printf("msg: %s\n", string(msg))
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
	mechanism, err := scram.Mechanism(scram.SHA512, "username", "password")
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "consumer-group-id",
		Topic:   topic,
		Dialer:  dialer,
	})

	for {
		m, err := r.ReadMessage(context.Background())

		if err != nil {
			fmt.Println(err)
			break
		}
		// TODO: process message
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		consumer(m)
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
