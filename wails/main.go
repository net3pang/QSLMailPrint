package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	if err := wails.Run(&options.App{
		Title:  "QSL Mail",
		Width:  1440,
		Height: 960,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 243, G: 246, B: 246, A: 1},
		OnStartup: app.startup,
		Bind:      []interface{}{app},
	}); err != nil {
		println("Error:", err.Error())
	}
}
