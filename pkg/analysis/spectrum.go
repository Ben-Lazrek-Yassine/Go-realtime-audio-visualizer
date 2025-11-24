package analysis

import (
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

type Spectrum struct {
	Bass  float64
	Mids  float64
	Highs float64
}

func Process(samples []float32) (s Spectrum) {
	if len(samples) == 0 {
		return Spectrum{}
	}

	input := make([]float64, len(samples))
	for i, v := range samples {
		input[i] = float64(v)
	}

	coeffs := fft.FFTReal(input)

	for i, c := range coeffs {
		if i > len(coeffs)/2 {
			break
		}

		magnitude := cmplx.Abs(c)

		switch {
		case i > 0 && i < 5:
			s.Bass += magnitude
		case i >= 5 && i < 50:
			s.Mids += magnitude
		case i >= 50:
			s.Highs += magnitude
		}
	}

	return s
}
