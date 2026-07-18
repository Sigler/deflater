package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"deflater/internal/logging"
	"deflater/internal/maintain"
)

//go:embed all:frontend/dist
var assets embed.FS

const appVersion = "0.1.0"

func main() {
	// Headless maintenance pass, run by the scheduled task. No window.
	for _, arg := range os.Args[1:] {
		if arg == "--maintenance" {
			maintain.Run()
			return
		}
	}

	app := NewApp()
	err := wails.Run(&options.App{
		Title:     "Deflater",
		Width:     1080,
		Height:    780,
		MinWidth:  920,
		MinHeight: 620,
		// Warm near-black behind the webview so resizes never flash white.
		BackgroundColour: &options.RGBA{R: 24, G: 22, B: 21, A: 255},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		Bind:             []any{app},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
	})
	if err != nil {
		logging.Logf("fatal: %v", err)
	}
}
