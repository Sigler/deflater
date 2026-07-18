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
		// Matches the mascot art's navy so resizes never flash white.
		BackgroundColour: &options.RGBA{R: 11, G: 17, B: 27, A: 255},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind:             []any{app},
		// One window at a time: a second launch focuses the first, so two
		// instances can never race on config or double-apply.
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "deflater-single-instance",
		},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
	})
	if err != nil {
		logging.Logf("fatal: %v", err)
	}
}
