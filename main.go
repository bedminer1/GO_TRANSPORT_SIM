package main

import (
	"fmt"
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

func main() {
	fmt.Println("Hello World")
}