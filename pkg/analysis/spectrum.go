package analysis

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

type Spectrum struct {
	SubBass    float64 // 20 - 60 Hz   (Rumble)
	Bass       float64 // 60 - 250 Hz  (Thump)
	LowMids    float64 // 250 - 500 Hz (Punch)
	Mids       float64 // 500 - 2k Hz  (Voice)
	Presence   float64 // 4k-6k Hz       (Air/Hats)
	Brilliance float64 // 6-20kHz
}

func Process(samples []float32) (s Spectrum) {
	resolution := 44100.0 / (float64(len(samples)) / 2.0)

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
		frequency := float64(i) * resolution

		switch {
		case frequency >= 20 && frequency < 60:
			s.SubBass += magnitude
		case frequency >= 60 && frequency < 250:
			s.Bass += magnitude
		case frequency >= 250 && frequency < 500:
			s.LowMids += magnitude
		case frequency >= 500 && frequency < 2000:
			s.Mids += magnitude
		case frequency >= 2000 && frequency < 6000:
			s.Presence += magnitude
		case frequency >= 6000 && frequency < 20000:
			s.Brilliance += magnitude
		default:
			// ignore
		}
	}

	scale := 30.0
	s.SubBass *= scale
	s.Bass *= scale
	s.LowMids *= scale
	s.Mids *= scale
	s.Presence *= scale
	s.Brilliance *= scale

	return s
}
