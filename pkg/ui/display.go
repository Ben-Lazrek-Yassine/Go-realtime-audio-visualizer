package ui

// Visualizer implements the Ebiten game interface
type Visualizer struct {
	// TODO: Add state for the visualizer (current frequency data)
}

// Update is called every tick (60Hz)
func (v *Visualizer) Update() error {
	// TODO: Get latest audio data
	return nil
}

// Draw renders the screen
func (v *Visualizer) Draw(screen interface{}) {
	// Note: screen will be *ebiten.Image
	// TODO: Draw bars/shapes based on frequency data
}

// Layout defines the screen dimensions
func (v *Visualizer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Run starts the Ebiten window
func Run() error {
	// TODO: ebiten.RunGame(&Visualizer{})
	return nil
}
