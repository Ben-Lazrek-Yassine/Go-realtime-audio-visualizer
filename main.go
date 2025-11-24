package main

import (
	"fmt"
	"visualizer/logger"
	"visualizer/pkg/analysis"
	"visualizer/pkg/audio"
)

func main() {
	stream := audio.NewAudioStream()
	defer stream.Close()

	if err := stream.StartCapture(); err != nil {
		logger.Error("Failed to start capture: %v", err)
		return
	}

	var frameCount int
	for {
		samples := <-stream.Samples
		frameCount++
		spectrum := analysis.Process(samples)
		drawBar := func(val float64, color string) string {
			maxVal := 100.0 // Adjusted for Log scale
			width := 10     // Slightly smaller to fit 5 bars
			percent := val / maxVal
			if percent > 1.0 {
				percent = 1.0
			}
			if percent < 0 {
				percent = 0
			}

			filled := int(float64(width) * percent)
			bar := ""
			for i := 0; i < width; i++ {
				if i < filled {
					bar += "|"
				} else {
					bar += " "
				}
			}
			return fmt.Sprintf("%s[%s]%s", color, bar, logger.Reset)
		}

		fmt.Printf("\033[2K\r SUB:   %s  %.0f\n", drawBar(spectrum.Sub, logger.Red), spectrum.Sub)
		fmt.Printf("\033[2K\r KICK:  %s  %.0f\n", drawBar(spectrum.Kick, logger.Red), spectrum.Kick)
		fmt.Printf("\033[2K\r SNARE: %s  %.0f\n", drawBar(spectrum.Snare, logger.Yellow), spectrum.Snare)
		fmt.Printf("\033[2K\r VOCAL: %s  %.0f\n", drawBar(spectrum.Vocals, logger.Green), spectrum.Vocals)
		fmt.Printf("\033[2K\r HIGH:  %s  %.0f\n", drawBar(spectrum.Highs, logger.Cyan), spectrum.Highs)
		fmt.Printf("\033[5A") // Move up 5 lines
	}

}
