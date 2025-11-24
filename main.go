package main

import (
	"flag"
	"visualizer/pkg/visualizer/cli"
	"visualizer/pkg/visualizer/ui"
)

func main() {
	useCli := flag.Bool("cli", false, "Run in terminal mode")
	useUi := flag.Bool("ui", false, "Run in UI mode")
	flag.Parse()

	if *useCli {
		cli.Cli_Visualizer()
	} else if *useUi {
		ui.Run()
	}
}
