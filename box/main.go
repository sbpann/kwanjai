package main

type geometry interface {
	area() float64
	perim() float64
}

type rectangle struct {
	width  int
	height int
}
