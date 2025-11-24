package ui

import (
	"image/color"
	"visualizer/logger"
	"visualizer/pkg/analysis"
	"visualizer/pkg/audio"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Visualizer struct {
	Stream   *audio.AudioStream
	Analyzer *analysis.Analyzer

	Current []float64
	Peaks   []float64
}

func NewVisualizer(stream *audio.AudioStream, analyzer *analysis.Analyzer) *Visualizer {
	return &Visualizer{
		Stream:   stream,
		Analyzer: analyzer,
		Current:  make([]float64, 10),
		Peaks:    make([]float64, 10),
	}
}

func (v *Visualizer) Update() error {
	var raw analysis.Spectrum
	var hasNewData bool

loop:
	for {
		select {
		case samples := <-v.Stream.Samples:
			raw = v.Analyzer.Process(samples)
			hasNewData = true
		default:
			break loop
		}
	}

	if !hasNewData {
		return nil
	}

	rawValues := []float64{
		raw.Sub, raw.Kick, raw.Bass, raw.LowMid, raw.Snare,
		raw.Mids, raw.HighMids, raw.Presence, raw.Highs, raw.Air,
	}

	decay := 2.0

	for i := 0; i < 10; i++ {
		target := rawValues[i]

		if target > v.Current[i] {
			v.Current[i] += (target - v.Current[i]) * 0.5
		} else {
			v.Current[i] -= decay
		}
		if v.Current[i] < 0 {
			v.Current[i] = 0
		}

		if v.Current[i] > v.Peaks[i] {
			v.Peaks[i] = v.Current[i]
		} else {
			v.Peaks[i] -= decay * 0.5
		}
		if v.Peaks[i] < 0 {
			v.Peaks[i] = 0
		}
	}

	return nil
}

func (v *Visualizer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 25, 255}) // Dark Navy Background

	barWidth := float32(60)
	spacing := float32(80)
	startX := float32(112)
	bottomY := float32(500)

	colors := []color.RGBA{
		{255, 50, 50, 255},   // Sub
		{255, 100, 50, 255},  // Kick
		{255, 150, 50, 255},  // Bass
		{255, 200, 50, 255},  // LowMid
		{255, 255, 50, 255},  // Snare
		{150, 255, 100, 255}, // Mids
		{50, 255, 200, 255},  // HighMids
		{50, 200, 255, 255},  // Pres
		{50, 100, 255, 255},  // Highs
		{150, 50, 255, 255},  // Air
	}

	for i := 0; i < 10; i++ {
		x := startX + (float32(i) * spacing)
		height := float32(v.Current[i]) * 5
		peakY := bottomY - (float32(v.Peaks[i]) * 5)

		vector.DrawFilledRect(screen, x+5, bottomY-height+5, barWidth, height, color.RGBA{0, 0, 0, 100}, true)
		vector.DrawFilledRect(screen, x, bottomY-height, barWidth, height, colors[i], true)
		vector.DrawFilledRect(screen, x, peakY-5, barWidth, 4, color.White, true)

		vector.DrawFilledCircle(screen, x+barWidth/2, bottomY+20, 3, color.White, true)
	}
}

func (v *Visualizer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func Run() {
	stream := audio.NewAudioStream()
	defer stream.Close()

	if err := stream.StartCapture(); err != nil {
		logger.Error("Failed to start capture: %v", err)
		return
	}

	analyzer := analysis.NewAnalyzer(4096)
	visualizer := NewVisualizer(stream, analyzer)

	ebiten.SetWindowSize(1024, 600)
	ebiten.SetWindowTitle("Go Audio Visualizer - Pro Mode")
	err := ebiten.RunGame(visualizer)
	if err != nil {
		logger.Error("Failed to run game: %v", err)
	}
}
