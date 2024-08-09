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
	terminal := 0

	for i := range busStops {
		busStops[i] = BusStop{id: i, queue: []Passenger{}}
	}

	buses := make ([]Bus, 5)
	for i := range buses {
		buses[i] = Bus{id: i, route: generateRoute(terminal, 7, len(busStops)), capacity: 20, passengers: []Passenger{}, location: 0}
	}

	return &System{
		busStops: busStops,
		buses: buses,
		terminal: terminal,
		totalWaitingTime: 0,
		totalPassengers: 0,
	}
}

func generateRoute(terminal int, numStops int, totalStops int) []int {
    // Start the route at the terminal
    route := []int{terminal}
    
    // Create a list of all possible stops, excluding the terminal
    allStops := make([]int, 0, totalStops-1)
    for i := 0; i < totalStops; i++ {
        if i != terminal {
            allStops = append(allStops, i)
        }
    }

    // Shuffle the list of stops and select a subset
    rand.Shuffle(len(allStops), func(i, j int) {
        allStops[i], allStops[j] = allStops[j], allStops[i]
    })
    
    // Add up to numStops - 2 stops (because the terminal is already added at the beginning)
    route = append(route, allStops[:numStops-2]...)
    
    // End the route at the terminal
    route = append(route, terminal)
    
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

			fmt.Printf("\nAverage Waiting Time: %v\nStop number: %v\nPeople waiting at this stop: %v\n", s.averageWaitingTime(), bus.location, len(s.busStops[bus.location].queue)) // stabilizes at 9.772898928s now
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

	for i := 0; i < len(stop.queue); {
		passenger := stop.queue[i]

		if isDirectRoute(bus.route, passenger.source, passenger.destination) {
            // Calculate waiting time for this passenger
            waitingTime := currentTime.Sub(passenger.arrivalTime)
            s.totalWaitingTime += waitingTime
            s.totalPassengers++
            
            // Add the passenger to the bus
            bus.passengers = append(bus.passengers, passenger)
			if len(bus.passengers) == bus.capacity {
				break
			}
            
            // Remove the passenger from the stop's queue
            stop.queue = append(stop.queue[:i], stop.queue[i+1:]...)
        } else {
            i++  // Move to the next passenger
        }
	}

}

func isDirectRoute(route []int, source int, destination int) bool {
	sourceFound := false

	for _, stop := range route {
		if stop == source {
			sourceFound = true
		}

		if sourceFound && stop == destination {
			return true
		}
	}

	return false
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
	bus.route = s.generateRouteToStop(longestQueueStop)
}

func (s *System) generateRouteToStop(stop int) []int {
	route := generateRoute(s.terminal, 7, len(s.busStops))
	 // Ensure the stop is included if not already
	 found := false
	 for _, r := range route {
		 if r == stop {
			 found = true
			 break
		 }
	 }
	 if !found {
		 route = append(route[:len(route) - 1], stop, s.terminal)
	 }

	return route
}

func main() {
	rand.Seed(time.Now().UnixNano())
	system := initializeSystem()
	go system.simulatePassengerArrivals()
	system.moveBuses()
}