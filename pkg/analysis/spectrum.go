package analysis

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

type Spectrum struct {
	Sub      float64 // 20-50Hz (Rumble)
	Kick     float64 // 50-100Hz (Thump)
	Bass     float64 // 100-200Hz (Basslines)
	LowMid   float64 // 200-500Hz (Mud/Warmth)
	Snare    float64 // 500-1000Hz (Snare Crack)
	Mids     float64 // 1000-2000Hz (Leads)
	HighMids float64 // 2000-4000Hz (Edge)
	Presence float64 // 4000-6000Hz (Clarity)
	Highs    float64 // 6000-10000Hz (Hats)
	Air      float64 // 10000-20000Hz (Sparkle)
}

type Analyzer struct {
	windowSize int
	leftBuf    []float64
	rightBuf   []float64
}

func NewAnalyzer(windowSize int) *Analyzer {
	return &Analyzer{
		windowSize: windowSize,
		leftBuf:    make([]float64, 0, windowSize),
		rightBuf:   make([]float64, 0, windowSize),
	}
}

func (a *Analyzer) Process(samples []float32) Spectrum {
	// 1. De-interleave and Append to Sliding Buffers
	for i := 0; i < len(samples); i += 2 {
		a.leftBuf = append(a.leftBuf, float64(samples[i]))
		a.rightBuf = append(a.rightBuf, float64(samples[i+1]))
	}

	// 2. Check if we have enough data for a full window
	if len(a.leftBuf) < a.windowSize {
		return Spectrum{}
	}

	// 3. Sliding Window: Keep only the last 'windowSize' samples
	if len(a.leftBuf) > a.windowSize {
		a.leftBuf = a.leftBuf[len(a.leftBuf)-a.windowSize:]
		a.rightBuf = a.rightBuf[len(a.rightBuf)-a.windowSize:]
	}

	// 4. Prepare inputs with Window Function
	lInput := make([]float64, a.windowSize)
	rInput := make([]float64, a.windowSize)

	for i := 0; i < a.windowSize; i++ {
		// Hanning Window
		window := 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(i)/float64(a.windowSize-1)))
		lInput[i] = a.leftBuf[i] * window
		rInput[i] = a.rightBuf[i] * window
	}

	// 5. Perform Stereo FFT
	lCoeffs := fft.FFTReal(lInput)
	rCoeffs := fft.FFTReal(rInput)

	// 6. Analyze
	var s Spectrum
	resolution := 44100.0 / float64(a.windowSize)

	for i := 0; i < len(lCoeffs)/2; i++ {
		// Magnitude is average of L and R
		lMag := cmplx.Abs(lCoeffs[i])
		rMag := cmplx.Abs(rCoeffs[i])

		magnitude := (lMag + rMag) / 2.0
		magnitude = math.Log10(magnitude + 1.0)

		if magnitude < 0.1 {
			magnitude = 0
		}

		frequency := float64(i) * resolution

		switch {
		case frequency >= 20 && frequency < 50:
			s.Sub = math.Max(s.Sub, magnitude)
		case frequency >= 50 && frequency < 100:
			s.Kick = math.Max(s.Kick, magnitude)
		case frequency >= 100 && frequency < 200:
			s.Bass = math.Max(s.Bass, magnitude)
		case frequency >= 200 && frequency < 500:
			s.LowMid = math.Max(s.LowMid, magnitude)
		case frequency >= 500 && frequency < 1000:
			s.Snare = math.Max(s.Snare, magnitude)
		case frequency >= 1000 && frequency < 2000:
			s.Mids = math.Max(s.Mids, magnitude)
		case frequency >= 2000 && frequency < 4000:
			s.HighMids = math.Max(s.HighMids, magnitude)
		case frequency >= 4000 && frequency < 6000:
			s.Presence = math.Max(s.Presence, magnitude)
		case frequency >= 6000 && frequency < 10000:
			s.Highs = math.Max(s.Highs, magnitude)
		case frequency >= 10000 && frequency < 20000:
			s.Air = math.Max(s.Air, magnitude)
		}
	}

	// Scale up (Peak Detection needs higher scale)
	scale := 50.0
	s.Sub *= scale
	s.Kick *= scale
	s.Bass *= scale
	s.LowMid *= scale
	s.Snare *= scale
	s.Mids *= scale
	s.HighMids *= scale
	s.Presence *= scale
	s.Highs *= scale
	s.Air *= scale

	return s
}
