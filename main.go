package main

import (
	"hello/views"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	window := app.NewWindow("Estacionamiento")
	window.CenterOnScreen()
	views.NewMainView(window)
	window.ShowAndRun()
}
