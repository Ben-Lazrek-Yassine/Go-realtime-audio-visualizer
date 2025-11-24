package analysis

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

type Spectrum struct {
	Sub    float64 // 20 - 60 Hz   (Rumble)
	Kick   float64 // 60 - 250 Hz  (Thump)
	Snare  float64 // 250 - 500 Hz (Punch)
	Vocals float64 // 500 - 2k Hz  (Voice)
	Highs  float64 // 2k+ Hz       (Air/Hats)
}

func Process(samples []float32) (s Spectrum) {
	if len(samples) == 0 {
		return Spectrum{}
	}

	monoLength := len(samples) / 2
	input := make([]float64, monoLength)
	for i := 0; i < monoLength; i++ {
		val := (float64(samples[i*2]) + float64(samples[i*2+1])) / 2.0
		// hanning window to reduce spectral leakage
		window := 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(i)/float64(monoLength-1)))
		input[i] = val * window
	}

	// FFT
	coeffs := fft.FFTReal(input)

	for i, c := range coeffs {
		if i > len(coeffs)/2 {
			break
		}

		magnitude := cmplx.Abs(c)
		magnitude = math.Log10(magnitude + 1.0)

		// ignore quiet sounds
		if magnitude < 0.1 {
			magnitude = 0
		}

		switch {
		case i == 1: // ~43Hz
			s.Sub += magnitude
		case i >= 2 && i < 4: // ~86Hz - 130Hz (Tight Kick)
			s.Kick += magnitude
		case i >= 4 && i < 9: // ~170Hz - 350Hz (Snare Body)
			s.Snare += magnitude
		case i >= 9 && i < 40: // ~380Hz - 1.7kHz (Vocals)
			s.Vocals += magnitude
		case i >= 40: // 1.7kHz+ (Highs)
			s.Highs += magnitude
		}
	}

	// Scale up a bit since Log10 squashed the numbers
	scale := 30.0
	s.Sub *= scale
	s.Kick *= scale
	s.Snare *= scale
	s.Vocals *= scale
	s.Highs *= scale

	return s
}
