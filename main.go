package main

import (
	"fmt"
	"visualizer/logger"
	"visualizer/pkg/audio"
)

func main() {
	stream := audio.NewAudioStream()
	defer stream.Close()

	if err := stream.StartCapture(); err != nil {
		logger.Error("Failed to start capture: %v", err)
		return
	}

	for {
		samples := <-stream.Samples
		fmt.Println(samples)
	}
}
