package cli

import (
	"fmt"
	"visualizer/logger"
	"visualizer/pkg/analysis"
	"visualizer/pkg/audio"
)

func Cli_Visualizer() {
	stream := audio.NewAudioStream()
	defer stream.Close()

	if err := stream.StartCapture(); err != nil {
		logger.Error("Failed to start capture: %v", err)
		return
	}

	analyzer := analysis.NewAnalyzer(4096)

	var frameCount int
	for {
		samples := <-stream.Samples
		frameCount++
		spectrum := analyzer.Process(samples)
		drawBar := func(val float64, color string) string {
			maxVal := 100.0 // Adjusted for Log scale
			width := 50
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
		fmt.Printf("\033[2K\r BASS:  %s  %.0f\n", drawBar(spectrum.Bass, logger.Yellow), spectrum.Bass)
		fmt.Printf("\033[2K\r L-MID: %s  %.0f\n", drawBar(spectrum.LowMid, logger.Yellow), spectrum.LowMid)
		fmt.Printf("\033[2K\r SNARE: %s  %.0f\n", drawBar(spectrum.Snare, logger.Green), spectrum.Snare)
		fmt.Printf("\033[2K\r LEAD:  %s  %.0f\n", drawBar(spectrum.Mids, logger.Green), spectrum.Mids)
		fmt.Printf("\033[2K\r EDGE:  %s  %.0f\n", drawBar(spectrum.HighMids, logger.Cyan), spectrum.HighMids)
		fmt.Printf("\033[2K\r PRES:  %s  %.0f\n", drawBar(spectrum.Presence, logger.Cyan), spectrum.Presence)
		fmt.Printf("\033[2K\r HATS:  %s  %.0f\n", drawBar(spectrum.Highs, logger.Blue), spectrum.Highs)
		fmt.Printf("\033[2K\r AIR:   %s  %.0f\n", drawBar(spectrum.Air, logger.Blue), spectrum.Air)
		fmt.Printf("\033[10A")
	}
}
