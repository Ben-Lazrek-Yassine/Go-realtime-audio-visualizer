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
			maxVal := 200.0
			width := 15
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

		fmt.Printf("\033[2K\r BASS: %s  %.0f\n", drawBar(spectrum.Bass, logger.Red), spectrum.Bass)
		fmt.Printf("\033[2K\r MIDS: %s  %.0f\n", drawBar(spectrum.Mids, logger.Green), spectrum.Mids)
		fmt.Printf("\033[2K\r HIGH: %s  %.0f\n", drawBar(spectrum.Highs, logger.Cyan), spectrum.Highs)
		fmt.Printf("\033[3A")
	}

}
