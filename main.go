package main

import (
	"fmt"
	"os"

	"github.com/gen2brain/malgo"
)

func main() {
	context, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		fmt.Println("Error initializing context:", err)
		os.Exit(1)
	}

	defer func() {
		_ = context.Uninit()
		context.Free()
	}()

	devices, err := context.Devices(malgo.Playback)
	if err != nil {
		fmt.Println("Error getting default device:", err)
		os.Exit(1)
	}
	fmt.Println("Default device:", devices)

}
