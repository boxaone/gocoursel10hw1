package main

import "fmt"

const (
	catFoodPerMonth = 7
	dogFoodPerMonth = 10 / 5
	cowFoodPerMonth = 25
)

// Common types declarations
type Eater interface {
	foodNeded() int
}

type WhoAmI interface {
	whoAmI() string
}

type Pet struct {
	weight int
}

// Cat declarations
type Cat Pet

func (cat Cat) foodNeded() int {
	return cat.weight * catFoodPerMonth
}
func (cat Cat) whoAmI() {
	fmt.Printf("I'm a cat. I'm weighting %v kg and I need %v kg of food per month\n", cat.weight, cat.foodNeded())
}

// Dog declarations
type Dog Pet

func (dog Dog) foodNeded() int {
	return dog.weight * dogFoodPerMonth
}
func (dog Dog) whoAmI() {
	fmt.Printf("I'm a dog. I'm weighting %v kg and I need %v kg of food per month\n", dog.weight, dog.foodNeded())
}

// Cow declarations
type Cow Pet

func (cow Cow) foodNeded() int {
	return cow.weight * cowFoodPerMonth
}
func (cow Cow) whoAmI() {
	fmt.Printf("I'm a cow. I'm weighting %v kg and I need %v kg of food per month\n", cow.weight, cow.foodNeded())
}
