package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	catFoodPerMonth = 7
	catMinWeight    = 1
	catMaxWeight    = 7
	dogFoodPerMonth = 10 / 5
	dogMinWeight    = 1
	dogMaxWight     = 20
	cowFoodPerMonth = 25
	cowMinWeight    = 10
	cowMaxWeight    = 300
	maxPets         = 20
)

var (
	locale     = "ua"
	randSource *rand.Rand
	petNames   = map[string]map[string]string{
		"en": {
			"cat": "cat",
			"dog": "dog",
			"cow": "cow",
		},
		"ua": {
			"cat": "кішка",
			"dog": "пес",
			"cow": "корова",
		},
	}
	phrases = map[string]map[string]string{
		"en": {
			"whoami": "I'm a %v. I'm weighting %v kg and I need %v kg of food per month\n",
		},
		"ua": {
			"whoami": "Я %v. Я важу %v кілограм і мені потрібно %v кілограм їжи в місяць\n",
		},
	}
)

// Common types declarations
type Eater interface {
	foodNeded() int
}

type WhoAmI interface {
	whoAmI()
}

type ProCreator interface {
	giveBirth(weight float64) Pet
}
type Animal struct {
	weight float64
}

type HaveWeight interface {
	Weight() float64
}

type Pet interface {
	HaveWeight
	Eater
	WhoAmI
	ProCreator
}

// Return animal weight

// Custom foonction of rounding food weight, int + adding 1 extra spare kg
func foodRound(weight float64) int {
	return int(weight) + 1
}

// Return whoami string
func returnWhoAmI(t string, p Pet) string {
	return fmt.Sprintf(phrases[locale]["whoami"], petNames[locale][t], p.Weight(), p.foodNeded())

}

// Cat declarations
type Cat Animal

func (cat Cat) foodNeded() int {
	return foodRound(cat.weight * catFoodPerMonth)
}

func (cat Cat) whoAmI() {
	fmt.Println(returnWhoAmI("cat", cat))
}

func (cat Cat) giveBirth(weight float64) Pet {
	return Cat{weight: weight}
}

func (cat Cat) Weight() float64 {
	return cat.weight
}

// Dog declarations
type Dog Animal

func (dog Dog) foodNeded() int {
	return foodRound(dog.weight * dogFoodPerMonth)
}

func (dog Dog) whoAmI() {
	fmt.Println(returnWhoAmI("dog", dog))

}

func (dog Dog) giveBirth(weight float64) Pet {

	return Dog{weight: weight}
}

func (dog Dog) Weight() float64 {
	return dog.weight
}

// Cow declarations
type Cow Animal

func (cow Cow) foodNeded() int {
	return foodRound(cow.weight * cowFoodPerMonth)
}

func (cow Cow) whoAmI() {
	fmt.Println(returnWhoAmI("cow", cow))
}

func (cow Cow) giveBirth(weight float64) Pet {
	return Cow{weight: weight}
}

func (cow Cow) Weight() float64 {
	return cow.weight
}

// Farm declaration
type Farm []*Pet

func main() {
	randS := rand.NewSource(time.Now().UnixNano())
	randSource = rand.New(randS)
	pets := []Pet{Dog{}, Cat{}, Cow{}}

	for _, p := range pets {
		p.giveBirth(10).whoAmI()

	}

}
