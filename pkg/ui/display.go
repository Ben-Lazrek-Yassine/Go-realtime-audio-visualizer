package ui

import (
	"visualizer/logger"
	"visualizer/pkg/analysis"
	"visualizer/pkg/audio"

	"github.com/hajimehoshi/ebiten/v2"
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
	// TODO: Draw the bars here!
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

	err := ebiten.RunGame(visualizer)
	if err != nil {
		logger.Error("Failed to run game: %v", err)
	}
}
