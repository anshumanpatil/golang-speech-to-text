package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// {   'alternative': [{'confidence': 0.88687539, 'transcript': 'hello'}],
//
//	'final': True}

type Possibility struct {
	confidence float64
	transcript string
}

type Verb struct {
	alternative []interface{}
	final       interface{}
}

func main() {

	m := kafka.ConfigMap{}
	m["bootstrap.servers"] = "localhost:9092"
	m["group.id"] = "kafka-go-getting-started"
	m["auto.offset.reset"] = "earliest"

	c, err := kafka.NewConsumer(&m)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "speech"
	err = c.SubscribeTopics([]string{topic}, nil)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			speech := string(ev.Value)
			fmt.Println("string speech", speech)
		}
	}

	c.Close()

}
