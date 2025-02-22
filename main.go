package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	goruntime "runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create main menu
	appMenu := menu.NewMenu()

	// macOS: Create app menu (appears next to Apple logo)
	if goruntime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())           // Add default App menu items (About, Preferences, etc.)
		appSubmenu := appMenu.AddSubmenu("File") // Empty label becomes app name menu

		// Save
		appSubmenu.AddText("Save", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(app.ctx, "saveRequested")
		})

		// Open
		appSubmenu.AddText("Open", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(app.ctx, "openRequested")
		})
	} else {
		// Windows/Linux: Use traditional File menu
		fileMenu := appMenu.AddSubmenu("File")
		fileMenu.AddText("Save", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(app.ctx, "saveRequested")
		})
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "nyxpad",
		Menu:   appMenu,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
