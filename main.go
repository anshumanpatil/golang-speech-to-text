package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Possibility struct {
	Confidence float64 `json:"confidence"`
	Transcript string  `json:"transcript"`
}

type Verb struct {
	Alternative []Possibility `json:"alternative"`
	Final       bool          `json:"final"`
}

func main() {
	mode := flag.String("mode", "", "")
	flag.Parse()

	if mode != nil {
		fmt.Println("mode = all: running speech recogniser!")
		cmd := exec.Command("python3", "speech.py")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
		go copyOutput(stdout)
		go copyOutput(stderr)

		if *mode == "recogniser" {
			return
		}

	}

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
			spchVerb := Verb{}
			errJson := json.Unmarshal([]byte(speech), &spchVerb)
			if errJson != nil {
				fmt.Println("errJson Value - ", speech)
				continue
			}
			fmt.Println("Value - ", spchVerb)

		}
	}

	c.Close()

}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
