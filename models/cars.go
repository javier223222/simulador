package models

import (
    "fmt"
    "image/color"
    "math/rand"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
)

type Vehicle struct {
    ID           int
    shape        *canvas.Image
    timerText    *canvas.Text
    remainingTime int
    quitSignal   chan bool
}

const (
    minWaitTime = 25
    maxWaitTime = 30
)

var (
    exitingVehicles []*Vehicle
)


const carImagePath = "resources/car.png"

const emptySlotImagePath = "resources/empty_slot.png" // Ruta de la imagen del cuadro vacío

func NewEmptyVehicle() *Vehicle {
    // Cargar la imagen del cuadro vacío
    shape := canvas.NewImageFromFile(emptySlotImagePath)
    shape.FillMode = canvas.ImageFillContain // Asegurar que la imagen se ajuste sin distorsión
    shape.SetMinSize(fyne.NewSize(40, 20))   // Tamaño adecuado para el espacio vacío

    timerText := canvas.NewText(fmt.Sprintf("%d", 0), color.White)
    timerText.Hide()

    vehicle := &Vehicle{
        ID:           -1,
        shape:        shape, // Imagen del cuadro vacío
        remainingTime: 0,
        timerText:    timerText,
    }

    return vehicle
}

func NewVehicle(id int, quitSignal chan bool) *Vehicle {
    waitTime := rand.Intn(maxWaitTime-minWaitTime) + minWaitTime

    shape := canvas.NewImageFromFile(carImagePath) 
    shape.FillMode = canvas.ImageFillContain       
    shape.SetMinSize(fyne.NewSize(40, 20))         

    timerText := canvas.NewText(fmt.Sprintf("%d", waitTime), color.White)
    timerText.Hide()

    vehicle := &Vehicle{
        ID:           id,
        shape:        shape,
        remainingTime: waitTime,
        timerText:    timerText,
        quitSignal:   quitSignal,
    }

    return vehicle
}

func (v *Vehicle) StartCountdown(id int) {
    for {
        select {
        case <-v.quitSignal:
            fmt.Printf("StartCountdown Close")
            return
        default:
            if v.remainingTime <= 0 {
                v.ID = id
                exitingVehicles = append(exitingVehicles, v)
                return
            }
            v.remainingTime--
            v.timerText.Text = fmt.Sprintf("%d", v.remainingTime)
            time.Sleep(1 * time.Second)
        }
    }
}

func (v *Vehicle) GetShape() *canvas.Image {
    return v.shape
}

func (v *Vehicle) UpdateData(vehicle *Vehicle) {
    v.ID = vehicle.ID
    v.remainingTime = vehicle.remainingTime
    v.timerText.Text = vehicle.timerText.Text
    v.timerText.Color = vehicle.timerText.Color

    if v.ID == -1 {
        // Si el espacio está vacío, muestra la imagen del cuadro vacío
        v.shape.File = emptySlotImagePath
    } else {
        // Si el espacio está ocupado, muestra la imagen del carro
        v.shape.File = carImagePath
    }
    v.shape.Refresh() // Refresca la imagen para mostrar el cambio
}

func (v *Vehicle) GetTimerText() *canvas.Text {
    return v.timerText
}

func (v *Vehicle) GetRemainingTime() int {
    return v.remainingTime
}

func (v *Vehicle) GetID() int {
    return v.ID
}

func GetExitingVehicles() []*Vehicle {
    return exitingVehicles
}

func PopNextExitingVehicle() *Vehicle {
    vehicle := exitingVehicles[0]
    if !IsExitingQueueEmpty() {
        exitingVehicles = exitingVehicles[1:]
    }
    return vehicle
}

func IsExitingQueueEmpty() bool {
    return len(exitingVehicles) == 0
}

