package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"deflater/internal/elevate"
	"deflater/internal/logging"
	"deflater/internal/maintain"
)

//go:embed all:frontend/dist
var assets embed.FS

const appVersion = "0.1.2"

func main() {
	// One line at the very first moment of every launch, so an elevated
	// relaunch is visible in the log even if it exits before the window.
	logging.Logf("launch: args=%v elevated=%v", os.Args[1:], elevate.IsElevated())

	// Headless maintenance pass, run by the scheduled task. No window.
	for _, arg := range os.Args[1:] {
		if arg == "--maintenance" {
			maintain.Run()
			return
		}
	}

	app := NewApp()
	// Note: deliberately NO Wails SingleInstanceLock. It rejects the app's
	// own elevated self-relaunch as a "second instance", breaking the
	// apply-with-admin flow. Concurrent-write safety is already handled by
	// the config file lock and the one-shot pending token.
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
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
	})
	if err != nil {
		logging.Logf("fatal: %v", err)
	}
}
