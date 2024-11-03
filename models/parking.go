package models

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	DarkGray         = color.RGBA{R: 24, G: 30, B: 30, A: 255}
	BlackColor       = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	mutexParkingLock sync.Mutex
)

const (
	arrivalRate         = 2.0
	MaxWaitingVehicles int = 30
	MaxParkingSpots    int = 20
)

type ParkingLot struct {
	waitingVehicles          []*Vehicle
	parkingSpaces            [MaxParkingSpots]*Vehicle
	entryStation             *Vehicle
	exitStation              *Vehicle
	outStation               *Vehicle
	quitChannel              chan bool
	newVehicleWaitingChannel chan bool
}

func NewParkingLot(vehicleWaitChannel chan bool, quitChannel chan bool) *ParkingLot {
	parkingLot := &ParkingLot{
		newVehicleWaitingChannel: vehicleWaitChannel,
		quitChannel:              quitChannel,
	}
	return parkingLot
}

func (p *ParkingLot) InitializeParking() {
	for i := range p.parkingSpaces {
		spaceVehicle := NewEmptyVehicle()
		p.parkingSpaces[i] = spaceVehicle
	}
}

func (p *ParkingLot) CreateOutStation() *Vehicle {
	p.outStation = NewEmptyVehicle()
	return p.outStation
}

func (p *ParkingLot) CreateExitStation() *Vehicle {
	p.exitStation = NewEmptyVehicle()
	return p.exitStation
}

func (p *ParkingLot) CreateEntryStation() *Vehicle {
	p.entryStation = NewEmptyVehicle()
	return p.entryStation
}

func (p *ParkingLot) GenerateVehicleQueue() {
	vehicleID := 20
	for {
		select {
		case <-p.quitChannel:
			fmt.Printf("GenerateVehicleQueue Close")
			return
		default:
			interarrivalTime := -math.Log(1-rand.Float64()) / arrivalRate // Poisson distribution
			time.Sleep(time.Duration(interarrivalTime * float64(time.Second)))
			if len(p.waitingVehicles) < MaxWaitingVehicles {
				vehicle := NewVehicle(vehicleID, p.quitChannel)
				vehicleID++
				p.waitingVehicles = append(p.waitingVehicles, vehicle)
				p.newVehicleWaitingChannel <- true
			}
		}
	}
}

func (p *ParkingLot) MonitorParkingSpaces() {
	for {
		select {
		case <-p.quitChannel:
			fmt.Printf("MonitorParkingSpaces Close")
			return
		default:
			availableIndex := p.FindAvailableSpace()
			if availableIndex != -1 && !p.IsWaitingQueueEmpty() {
				mutexParkingLock.Lock()
				p.TransferToEntry()
				p.TransferToParking(availableIndex)
				mutexParkingLock.Unlock()
			}
		}
	}
}

func (p *ParkingLot) TransferToEntry() {
	vehicle := p.RemoveFirstWaitingVehicle()
	p.entryStation.UpdateData(vehicle)
	time.Sleep(1 * time.Second)
}

func (p *ParkingLot) TransferToParking(index int) {
	p.parkingSpaces[index].UpdateData(p.entryStation)
	p.parkingSpaces[index].timerText.Show()

	p.entryStation.UpdateData(NewEmptyVehicle())
	go p.parkingSpaces[index].StartCountdown(index)
	time.Sleep(1 * time.Second)
}

func (p *ParkingLot) MoveVehicleToExit() {
	for {
		select {
		case <-p.quitChannel:
			fmt.Printf("MoveVehicleToExit Close")
			return
		default:
			if !IsExitingQueueEmpty() {
				mutexParkingLock.Lock()
				vehicle := PopNextExitingVehicle()

				p.TransferToExit(vehicle.ID)
				p.TransferToOut()
				mutexParkingLock.Unlock()

				time.Sleep(1 * time.Second)
				p.outStation.UpdateData(NewEmptyVehicle())
			}
		}
	}
}

func (p *ParkingLot) TransferToExit(index int) {
	if index < MaxParkingSpots {
		p.exitStation.UpdateData(p.parkingSpaces[index])
		p.parkingSpaces[index].timerText.Hide()
		p.parkingSpaces[index].UpdateData(NewEmptyVehicle())
		time.Sleep(1 * time.Second)
	}
}

func (p *ParkingLot) TransferToOut() {
	p.outStation.UpdateData(p.exitStation)
	p.exitStation.UpdateData(NewEmptyVehicle())
	time.Sleep(1 * time.Second)
}

func (p *ParkingLot) FindAvailableSpace() int {
	for i := range p.parkingSpaces {
		if p.parkingSpaces[i].GetID() == -1 {
			return i
		}
	}
	return -1
}

func (p *ParkingLot) RemoveFirstWaitingVehicle() *Vehicle {
	vehicle := p.waitingVehicles[0]
	if !p.IsWaitingQueueEmpty() {
		p.waitingVehicles = p.waitingVehicles[1:]
	}
	return vehicle
}

func (p *ParkingLot) IsWaitingQueueEmpty() bool {
	return len(p.waitingVehicles) == 0
}

func (p *ParkingLot) GetWaitingVehicles() []*Vehicle {
	return p.waitingVehicles
}

func (p *ParkingLot) GetEntryStationVehicle() *Vehicle {
	return p.entryStation
}

func (p *ParkingLot) GetExitStationVehicle() *Vehicle {
	return p.exitStation
}

func (p *ParkingLot) GetParkingSpaces() [MaxParkingSpots]*Vehicle {
	return p.parkingSpaces
}

func (p *ParkingLot) ClearParkingSpaces() {
	for i := range p.parkingSpaces {
		p.parkingSpaces[i] = nil
	}
}
