package audio

// AudioStream manages the audio capture context and device
type AudioStream struct {
	// TODO: Add malgo context and device fields here
}

// NewAudioStream initializes the audio system
func NewAudioStream() *AudioStream {
	// TODO: Initialize malgo context
	return &AudioStream{}
}

// StartCapture begins the loopback audio capture
func (as *AudioStream) StartCapture() error {
	// TODO: Find default playback device
	// TODO: Init device with Capture config using Playback ID
	// TODO: Start device
	return nil
}

// Close cleans up resources
func (as *AudioStream) Close() {
	// TODO: Free context and device
}
