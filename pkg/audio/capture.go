package audio

import (
	"runtime"
	"unsafe"
	"visualizer/logger"

	"github.com/gen2brain/malgo"
)

type AudioStream struct {
	Context *malgo.AllocatedContext
	Device  *malgo.Device
	Samples chan []float32
}

type AudioConfig struct {
	Format           malgo.FormatType
	Channels         int
	SampleRate       int
	PlaybackDeviceID malgo.DeviceID
	CaptureDeviceID  malgo.DeviceID
}

func NewAudioStream() *AudioStream {
	context, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil
	}
	return &AudioStream{
		Context: context,
		Device:  nil,
		Samples: make(chan []float32, 1024),
	}
}

func (as *AudioStream) StartCapture() error {
	logger.Info("Starting capture")
	devices, err := as.Context.Devices(malgo.Playback)
	if err != nil {
		logger.Info("error getting default device")
		return err
	}

	speakerID := devices[0].ID
	idPtr := new(malgo.DeviceID)
	*idPtr = speakerID

	var pinner runtime.Pinner
	pinner.Pin(idPtr)
	defer pinner.Unpin()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Loopback)
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = 2
	deviceConfig.SampleRate = 44100

	deviceConfig.Capture.DeviceID = unsafe.Pointer(idPtr)

	onRecvFrame := func(pOutputSample, pInputSamples []byte, frameCount uint32) {
		sampleCount := frameCount * 2
		logger.Info("Audio frame received")
		pSamples := (*float32)(unsafe.Pointer(&pInputSamples[0]))
		rawSamples := unsafe.Slice(pSamples, sampleCount)

		newSlice := make([]float32, sampleCount)
		copy(newSlice, rawSamples)

		select {
		case as.Samples <- newSlice:
			logger.Info("Audio frame received")
		default:
			logger.Warning("Dropping audio frame! Channel full.")
		}
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrame,
	}

	device, err := malgo.InitDevice(as.Context.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return err
	}

	as.Device = device

	err = device.Start()
	if err != nil {
		return err
	}

	return nil
}
func (as *AudioStream) Close() {
	as.Context.Uninit()
	as.Context.Free()
}
