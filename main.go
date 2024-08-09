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
	totalWaitingTime time.Duration
	totalPassengers int
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
		totalWaitingTime: 0,
		totalPassengers: 0,
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

			fmt.Printf("Average Waiting Time: %v\n", s.averageWaitingTime()) // stabilizes at 9.772898928s now
			nextStop := (bus.location + 1) % len(bus.route)
			bus.location = bus.route[nextStop]
			s.pickupPassengers(bus)
		}
		time.Sleep(time.Second)
	}
}

func (s *System) pickupPassengers(bus *Bus) {
	stop := s.busStops[bus.location]
	currentTime := time.Now()

	for len(stop.queue) > 0 && len(bus.passengers) < bus.capacity {
		passenger := stop.queue[0]
		stop.queue = stop.queue[1:]

		// calculater waiting time
		waitingTime := currentTime.Sub(passenger.arrivalTime)
		s.totalWaitingTime += waitingTime
		s.totalPassengers++

		// add passenger to bus
		bus.passengers = append(bus.passengers, passenger)
	}

}

func (s *System) averageWaitingTime() time.Duration {
	if s.totalPassengers == 0 {
		return 0
	}
	
	return s.totalWaitingTime / time.Duration(s.totalPassengers)
}

func dispatchBus(s *System, bus *Bus) {
	// Current Strategy: select the bus stop with the longest queue
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