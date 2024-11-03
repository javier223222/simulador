package views

import (
	"fmt"
	"image/color"
	"hello/models"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var semRenderNewCarWait chan bool
var semQuit chan bool

type ParkingView struct {
	window               fyne.Window
	waitRectangleStation [models.MaxWaitingVehicles]*canvas.Image // Cambiado a *canvas.Image para usar la imagen del carro
}

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

var parking *models.ParkingLot

func NewParkingView(window fyne.Window) *ParkingView {
	parkingView := &ParkingView{window: window}

	semQuit = make(chan bool)
	semRenderNewCarWait = make(chan bool)

	parking = models.NewParkingLot(semRenderNewCarWait, semQuit)
	parkingView.MakeScene()
	parkingView.StartSimulation()

	return parkingView
}

func (p *ParkingView) MakeScene() {
	background := canvas.NewRectangle(color.RGBA{R: 56, G: 149, B: 132, A: 255}) // Cambia el color aquí
    background.Resize(p.window.Canvas().Size())
	containerParkingView := container.New(layout.NewVBoxLayout())
	containerParkingOut := container.New(layout.NewHBoxLayout())
	containerButtons := container.New(layout.NewHBoxLayout())

	restart := widget.NewButton("Reinciar la simulacion", func() {
		dialog.ShowConfirm("Salir", "¿Desea Reiniciar la aplicación?", func(response bool) {
			if response {
				p.RestartSimulation()
			} else {
				fmt.Println("No")
			}
		}, p.window)
	})

	exit := widget.NewButton("Salir",
		func() {
			dialog.ShowConfirm("Salir", "¿Desea salir de la aplicación?", func(response bool) {
				if response {
					p.BackToMenu()
				} else {
					fmt.Println("No")
				}
			}, p.window)
		},
	)

	containerButtons.Add(restart)
	containerButtons.Add(exit)

	containerParkingOut.Add(p.MakeWaitStation())
	containerParkingOut.Add(layout.NewSpacer())
	containerParkingOut.Add(p.MakeExitStation())
	containerParkingOut.Add(layout.NewSpacer())

	containerParkingView.Add(containerParkingOut)
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeParkingLotEntrance())
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeEnterAndExitStation())
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeParking())
	containerParkingView.Add(layout.NewSpacer())

	containerParkingView.Add(container.NewCenter(containerButtons))
	content := container.NewMax(background, containerParkingView)
	p.window.SetContent(content)
	
	p.window.Resize(fyne.NewSize(1000, 750))
	p.window.CenterOnScreen()
}

func (p *ParkingView) MakeParking() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parking.InitializeParking()

	parkingArray := parking.GetParkingSpaces()
	for i := 0; i < len(parkingArray); i++ {
		if i == 10 {
			addSpace(parkingContainer)
		}
		parkingContainer.Add(container.NewCenter(parkingArray[i].GetShape(), parkingArray[i].GetTimerText()))
	}
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeWaitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	for i := len(p.waitRectangleStation) - 1; i >= 0; i-- {
		car := canvas.NewImageFromFile("resources/car.png") 
		car.FillMode = canvas.ImageFillContain               
		car.SetMinSize(fyne.NewSize(40, 20))                 

		p.waitRectangleStation[i] = car
		p.waitRectangleStation[i].Hide()
		parkingContainer.Add(p.waitRectangleStation[i])
	}
	return parkingContainer
}

func (p *ParkingView) MakeExitStation() *fyne.Container {
	out := parking.CreateOutStation()
	return container.NewCenter(out.GetShape())
}

func (p *ParkingView) MakeEnterAndExitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parkingContainer.Add(layout.NewSpacer())
	entrace := parking.CreateEntryStation()
	parkingContainer.Add(entrace.GetShape())
	parkingContainer.Add(layout.NewSpacer())
	exit := parking.CreateExitStation()
	parkingContainer.Add(exit.GetShape())
	parkingContainer.Add(layout.NewSpacer())
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeParkingLotEntrance() *fyne.Container {
	EntraceContainer := container.New(layout.NewGridLayout(3))
	EntraceContainer.Add(makeBorder())
	EntraceContainer.Add(layout.NewSpacer())
	EntraceContainer.Add(makeBorder())
	return EntraceContainer
}

func (p *ParkingView) RenderNewCarWaitStation() {
	for {
		select {
		case <-semQuit:
			fmt.Printf("RenderNewCarWaitStation Close")
			return
		case <-semRenderNewCarWait:
			waitCars := parking.GetWaitingVehicles()
			for i := len(waitCars) - 1; i >= 0; i-- {
				if waitCars[i].ID != -1 {
					p.waitRectangleStation[i].Show()
					p.waitRectangleStation[i].Refresh()
				}
			}
			p.window.Content().Refresh()
		}
	}
}

func (p *ParkingView) RenderUpdate() {
	for {
		select {
		case <-semQuit:
			fmt.Printf("RenderUpdate Close")
			return
		default:
			p.window.Content().Refresh()
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *ParkingView) StartSimulation() {
	go parking.GenerateVehicleQueue()
	go parking.MoveVehicleToExit()
	go parking.MonitorParkingSpaces()
	go p.RenderNewCarWaitStation()
	go p.RenderUpdate()

}

func (p *ParkingView) BackToMenu() {
	close(semQuit)
	NewMainView(p.window)
}

func (p *ParkingView) RestartSimulation() {
	close(semQuit)
	NewParkingView(p.window)
}

func addSpace(parkingContainer *fyne.Container) {
	for j := 0; j < 5; j++ {
		parkingContainer.Add(layout.NewSpacer())
	}
}

func makeBorder() *canvas.Rectangle {
	square := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 0})
	square.SetMinSize(fyne.NewSquareSize(float32(30)))
	square.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	square.StrokeWidth = float32(1)
	return square
}

