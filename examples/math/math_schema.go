package main

type API interface {
	// Add two numbers
	Add(x, y int) int
	// Sub Subtract two numbers
	Sub(x, y int) int
	// Mul Multiply two numbers
	Mul(x, y int) int
	// Div Divide two numbers
	Div(x, y int) int
	// Neg Negate a number
	Neg(x int) int
	// ID Identity function
	ID(x int) int
	// Unknown request
	Unknown(text string) int
}
