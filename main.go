//go:build gui

package main

import (
	"embed"

	"muitoolunlock/ui"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/lang"
)

//go:embed translations/*
var translations embed.FS

func main() {
	a := app.New()
	lang.AddTranslationsFS(translations, "translations")

	// Create init screen
	initScreen := ui.NewInitScreen(a)

	// Start with init screen - it will handle navigation internally
	initScreen.Show()

	// Run the app - must be called from main goroutine
	a.Run()
}
