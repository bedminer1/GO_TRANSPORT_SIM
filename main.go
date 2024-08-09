package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Passenger struct {
	source int
	destination int
	arrivalTime time.Time
}

type Bus struct {
	id int
	route []int
	capacity int
	passengers []Passenger
	location int
}

type BusStop struct {
	id int
	queue []Passenger
}

type System struct {
	busStops []BusStop
	buses []Bus
	terminal int
}

// init System struct
func initializeSystem() *System {
	busStops := make([]BusStop, 10)

	for i := range busStops {
		busStops[i] = BusStop{id: i, queue: []Passenger{}}
	}

	buses := make ([]Bus, 5)
	for i := range buses {
		buses[i] = Bus{id: i, route: generateRoute(), capacity: 20, passengers: []Passenger{}, location: 0}
	}

	return &System{
		busStops: busStops,
		buses: buses,
		terminal: 0,
	}
}

func generateRoute() []int {
	route := rand.Perm(10)
	route = append(route, 0)

	return route
}

// simulate passenger arrivals
func (s *System) simulatePassengerArrivals() {
	for {
		source := rand.Intn(10)
		destination := rand.Intn(10)
		if source != destination {
			passenger := Passenger{
				source: source,
				destination: destination,
				arrivalTime: time.Now(),
			}

			s.busStops[source].queue = append(s.busStops[source].queue, passenger)
		}

		time.Sleep(time.Second * time.Duration(rand.Intn(5)))

	}
}

// Bus movement and Scheduling
func (s *System) moveBuses() {
	for {
		for i := range s.buses {
			bus := &s.buses[i]

			if bus.location == s.terminal {

				dispatchBus(s, bus)
			}

			nextStop := (bus.location + 1) % len(bus.route)
			bus.location = bus.route[nextStop]
			s.pickupPassengers(bus)
		}
		time.Sleep(time.Second)
	}
}

func (s *System) pickupPassengers(bus *Bus) {
	stop := s.busStops[bus.location]

	for len(stop.queue) > 0 && len(bus.passengers) < bus.capacity {
		passenger := stop.queue[0]
		stop.queue = stop.queue[1:]
		bus.passengers = append(bus.passengers, passenger)
	}

}

func dispatchBus(s *System, bus *Bus) {
	longestQueueStop := 0
	maxQueueLength := 0

	for i, stop := range s.busStops {
		if len(stop.queue) > maxQueueLength {
			longestQueueStop = i
			maxQueueLength = len(stop.queue)
		}
	}
	bus.route = generateRouteToStop(longestQueueStop)
}

func generateRouteToStop(stop int) []int {
	route := rand.Perm(10)
	for i, s := range route {
		if s == stop {
			route = append(route[:i], append([]int{stop}, route[i:]...)...)
			break
		}
	}

	return route
}

func main() {
	rand.Seed(time.Now().UnixNano())
	system := initializeSystem()
	go system.simulatePassengerArrivals()
	system.moveBuses()
}