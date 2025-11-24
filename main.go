package main

import (
	"visualizer/logger"
	"visualizer/pkg/audio"
)

func main() {
	stream := audio.NewAudioStream()
	defer stream.Close()
	logger.Info("Audio stream initialized")
	stream.StartCapture()
}
