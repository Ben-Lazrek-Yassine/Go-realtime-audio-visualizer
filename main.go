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

	infos, err := context.Devices(malgo.Playback)
	if err != nil {
		fmt.Println("Error getting devices:", err)
		os.Exit(1)
	}
	fmt.Println("Devices:", infos)

	for i, info := range infos {
		e := "ok"
		full, err := context.DeviceInfo(malgo.Playback, info.ID, malgo.Shared)
		if err != nil {
			e = err.Error()
		}
		fmt.Printf("    %d: %v, %s, [%s], formats: %+v\n",
			i, info.Name, info.ID, e, full.Formats)
	}

}
