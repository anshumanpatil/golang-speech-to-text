package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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

	v, err := verifyPythonVersion("3")
	if err != nil {
		log.Fatalf("Without python3 library won't work!!!")

	}
	if !strings.Contains(*v, "3") {
		log.Fatalf("Without python3 library won't work!!!")
	}

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

func verifyPythonVersion(v string) (*string, error) {
	_, err := exec.LookPath("python" + v)
	if err != nil {
		// log.Fatalf("No python version located")
		return nil, err
	}
	out, err := exec.Command("python"+v, "--version").CombinedOutput()
	// log.Print(out)
	if err != nil {
		// log.Fatalf("Error checking Python version with the 'python' command: %v", err)
		return nil, err
	}
	fields := strings.Fields(string(out))
	if len(fields) > 0 {
		version := fields[1]
		log.Print(version)
		return &version, nil
	}
	return nil, fmt.Errorf("No version found!!")
}
