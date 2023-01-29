package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/pkg/term"
	wave "github.com/zenwerk/go-wave"
)

func getch() []byte {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)
	numRead, err := t.Read(bytes)
	t.Restore()
	t.Close()
	if err != nil {
		return nil
	}
	return bytes[0:numRead]
}

func errCheck(err error, n ...string) {
	if err != nil {
		fmt.Println("errr - ", n)
		panic(err)
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s audiofilename.wav\n", os.Args[0])
		os.Exit(0)
	}

	audioFileName := os.Args[1]

	fmt.Println("Recording. Press ESC to quit.")

	if !strings.HasSuffix(audioFileName, ".wav") {
		audioFileName += ".wav"
	}
	waveFile, err := os.Create(audioFileName)
	errCheck(err, "1")

	// www.people.csail.mit.edu/hubert/pyaudio/  - under the Record tab
	inputChannels := 2
	outputChannels := 0
	sampleRate := 44100
	framesPerBuffer := make([]byte, 64)

	// init PortAudio

	portaudio.Initialize()
	//defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err, "2")
	//defer stream.Close()

	// setup Wave file writer

	param := wave.WriterParam{
		Out:           waveFile,
		Channel:       inputChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 8, // if 16, change to WriteSample16()
	}

	waveWriter, err := wave.NewWriter(param)
	errCheck(err, "3")

	//defer waveWriter.Close()

	go func() {
		key := getch()
		fmt.Println()
		fmt.Println("Cleaning up ...")
		fmt.Println(key, len(key))
		// if len(key) == 27 {
		// better to control
		// how we close then relying on defer
		waveWriter.Close()
		stream.Close()
		portaudio.Terminate()
		fmt.Println("Play", audioFileName, "with a audio player to hear the result.")
		os.Exit(0)

		// }

	}()

	// recording in progress ticker. From good old DOS days.
	// ticker := []string{
	// 	"-",
	// 	"\\",
	// 	"/",
	// 	"|",
	// }
	rand.Seed(time.Now().UnixNano())

	// start reading from microphone
	errCheck(stream.Start(), "4")
	for {
		errCheck(stream.Read(), "5")

		// fmt.Printf("\rRecording is live now. Say something to your microphone! [%v]", ticker[rand.Intn(len(ticker)-1)])

		// write to wave file
		_, err := waveWriter.Write([]byte(framesPerBuffer)) // WriteSample16 for 16 bits
		errCheck(err, "6")
	}
	errCheck(stream.Stop(), "7")
}
