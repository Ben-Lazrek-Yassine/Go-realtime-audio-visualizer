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
	Spectrum analysis.Spectrum
}

func NewVisualizer(stream *audio.AudioStream, analyzer *analysis.Analyzer) *Visualizer {
	return &Visualizer{
		Stream:   stream,
		Analyzer: analyzer,
	}
}

func (v *Visualizer) Update() error {
	select {
	case samples := <-v.Stream.Samples:
		v.Spectrum = v.Analyzer.Process(samples)
	default:
		// No new audio yet, keep showing the last frame
	}
	return nil
}

func (v *Visualizer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 20, 255})
	centerY := float32(300)
	startX := float32(100)
	spacing := float32(90)

	drawBeat := func(index int, val float64, clr color.RGBA) {
		x := startX + (float32(index) * spacing)

		radius := float32(val) * 0.8
		if radius < 10 {
			radius = 10
		}

		vector.DrawFilledCircle(screen, x, centerY, radius, clr, true)
	}

	drawBeat(0, v.Spectrum.Sub, color.RGBA{255, 0, 0, 255})
	drawBeat(1, v.Spectrum.Kick, color.RGBA{255, 50, 0, 255})
	drawBeat(2, v.Spectrum.Bass, color.RGBA{255, 100, 0, 255})

	drawBeat(3, v.Spectrum.LowMid, color.RGBA{255, 200, 0, 255})
	drawBeat(4, v.Spectrum.Snare, color.RGBA{0, 255, 0, 255})
	drawBeat(5, v.Spectrum.Mids, color.RGBA{0, 255, 128, 255})

	drawBeat(6, v.Spectrum.HighMids, color.RGBA{0, 255, 255, 255})
	drawBeat(7, v.Spectrum.Presence, color.RGBA{0, 128, 255, 255})
	drawBeat(8, v.Spectrum.Highs, color.RGBA{0, 0, 255, 255})
	drawBeat(9, v.Spectrum.Air, color.RGBA{128, 0, 255, 255})
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
	ebiten.SetWindowTitle("Go Audio Visualizer")
	err := ebiten.RunGame(visualizer)
	if err != nil {
		logger.Error("Failed to run game: %v", err)
	}
}
