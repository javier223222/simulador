package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MainView struct {
	window fyne.Window
}

func NewMainView(window fyne.Window) *MainView {
	MainView := &MainView{
		window: window,
	}
	MainView.InitApp()
	return MainView
}

func (m *MainView) InitApp() {
	m.DrawSceneMenu()
}

func (m *MainView) DrawSceneMenu() {
	   // Crear el fondo con el color especificado #E2F1E7
	   background := canvas.NewRectangle(color.RGBA{R: 226, G: 241, B: 231, A: 255}) // #E2F1E7 en RGBA
	   background.Resize(fyne.NewSize(1000, 500)) // Ajusta el tamaño del fondo al tamaño de la ventana
   
	   // Crear el título
	   title := canvas.NewText("Simulador de Estacionamiento", color.RGBA{R: 34, G: 54, B: 66, A: 255})
	   title.TextSize = 24 // Ajustar el tamaño de fuente
	   titleContainer := container.NewCenter(title)
   
	   // Crear los botones
	   start := widget.NewButton("Empezar simulacion", m.StartParkingSimulation)
	   exit := widget.NewButton("Salir", m.ExitGame)
   
	   // Contenedor principal
	   container_center := container.NewVBox(
		   titleContainer,
		   layout.NewSpacer(),
		   start,
		   exit,
		   layout.NewSpacer(),
	   )
   
	   // Colocar el fondo detrás del contenedor principal
	   content := container.NewMax(background, container.NewCenter(container_center))
   
	   // Configurar la ventana
	   m.window.SetContent(content)
	   m.window.Resize(fyne.NewSize(1000, 500))
	   m.window.SetFixedSize(true)
}

func (m *MainView) ExitGame() {
	m.window.Close()
}



func (m *MainView) StartParkingSimulation() {
	NewParkingView(m.window)
}
