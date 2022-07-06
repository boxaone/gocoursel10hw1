package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	catFoodPerMonth = 7
	catMinWeight    = 1
	catMaxWeight    = 7
	dogFoodPerMonth = 10 / 5
	dogMinWeight    = 1.
	dogMaxWight     = 20.
	cowFoodPerMonth = 25.
	cowMinWeight    = 10.
	cowMaxWeight    = 300.
	minPets         = 10
	maxPets         = 20
)

var (
	// Apps settings
	locale       = "en"
	grows, gcols = 1, 80
	sleepInt     = 50
	randSource   *rand.Rand

	// Animals settings
	pets = []Pet{Dog{}, Cat{}, Cow{}}

	petNames = map[string]map[string]string{
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
			"pet_info_string":        "I'm a %v. I'm weighting %.2f kg and I need %v kg of food per month\n",
			"choose_language_string": "Please choose output language or exit\n",
			"languages_string":       "1) English    2) Ukrainian  3) Exit\n",
			"gen_farm_string":        "Generating new farm\n",
			"ffood_info_string":      "Summary %v kg food per month needed \n",
			"calc_farm_info_string":  "Calculating all food needed\n",
		},
		"ua": {
			"pet_info_string":        "Я %v. Я важу %.2f кілограм і мені потрібно %v кілограм кормів в місяць\n",
			"choose_language_string": "Будь ласка, оберіть мову або вийти \n",
			"languages_string":       "1) Англійська 2) Українська 3) Вийти\n",
			"gen_farm_string":        "Генеруємо нову ферму\n",
			"ffood_info_string":      "Загалом треба %v кілограмів кормів в місяць",
			"calc_farm_info_string":  "Рахуємо загальну вагу кормів\n",
		},
	}
)

// Common types declarations
type Eater interface {
	foodNeded() int
}

type ShowInfo interface {
	showInfo()
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
	ShowInfo
	ProCreator
}

// Return animal weight

// Custom foonction of rounding food weight, int + adding 1 extra spare kg
func foodRound(weight float64) int {
	return int(weight) + 1
}

// Return pet_info_string string
func getPetInfo(t string, p Pet) string {
	return fmt.Sprintf(phrases[locale]["pet_info_string"], petNames[locale][t], p.Weight(), p.foodNeded())

}

func genWeight(weight float64, min float64, max float64) float64 {
	if weight == 0 {
		return (max-min)*randSource.Float64() + min
	}
	if weight < min {
		return min
	}
	if weight > max {
		return max
	}
	return weight

}

// Cat declarations
type Cat Animal

func (cat Cat) foodNeded() int {
	return foodRound(cat.weight * catFoodPerMonth)
}

func (cat Cat) showInfo() {
	fmt.Println(getPetInfo("cat", cat))
}

func (cat Cat) giveBirth(weight float64) Pet {
	return Cat{weight: genWeight(weight, catMinWeight, catMaxWeight)}
}

func (cat Cat) Weight() float64 {
	return cat.weight
}

// Dog declarations
type Dog Animal

func (dog Dog) foodNeded() int {
	return foodRound(dog.weight * dogFoodPerMonth)
}

func (dog Dog) showInfo() {
	fmt.Println(getPetInfo("dog", dog))

}

func (dog Dog) giveBirth(weight float64) Pet {

	return Dog{weight: genWeight(weight, dogMinWeight, dogMaxWight)}

}

func (dog Dog) Weight() float64 {
	return dog.weight
}

// Cow declarations
type Cow Animal

func (cow Cow) foodNeded() int {
	return foodRound(cow.weight * cowFoodPerMonth)
}

func (cow Cow) showInfo() {
	fmt.Println(getPetInfo("cow", cow))
}

func (cow Cow) giveBirth(weight float64) Pet {
	return Cow{weight: genWeight(weight, cowMinWeight, cowMaxWeight)}
}

func (cow Cow) Weight() float64 {
	return cow.weight
}

// Farm declaration
type Farm struct {
	Pets []Pet
}

// Show farm details
func (f *Farm) showInfo() {
	gsteps := grows * gcols
	if f.Pets != nil {
		for _, pet := range f.Pets {
			pet.showInfo()
		}
	}
	genNiceOutput(1, gcols, func(i, j int) {})
	fmt.Printf(phrases[locale]["calc_farm_info_string"])
	var foodSum, ind int

	genNiceOutput(grows, gcols, func(i, j int) {
		if f.Pets != nil && len(f.Pets) > 0 {
			time.Sleep(time.Millisecond * time.Duration(randSource.Intn(sleepInt)))
			if (ind*100)/(len(f.Pets)) <= ((i+1)*(j+1)*100)/gsteps {
				if ind < len(f.Pets) {
					foodSum += (f.Pets)[ind].foodNeded()
				}
				ind++
			}
		}

	})

	fmt.Printf(phrases[locale]["ffood_info_string"], foodSum)
}

// Generate random farm
func (f *Farm) genPets(max, min int) {

	f.Pets = make([]Pet, 0, max)
	numPets := randSource.Intn(max-min) + min

	gsteps := grows * gcols

	genNiceOutput(grows, gcols, func(i1, i2 int) {})
	fmt.Printf(phrases[locale]["gen_farm_string"])

	// Farm generation process with indication

	genNiceOutput(grows, gcols, func(i int, j int) {

		time.Sleep(time.Millisecond * time.Duration(randSource.Intn(sleepInt)))

		if (len(f.Pets)*100)/numPets < ((i+1)*(j+1)*100)/gsteps {
			f.Pets = append((f.Pets), pets[randSource.Intn(len(pets))].giveBirth(0))
		}

	})

}

func genNiceOutput(grows int, gcols int, fn func(int, int)) {
	fmt.Print("\033[s")

	for i := 0; i < grows; i++ {

		fmt.Print("\033[u\033[K")

		for j := 0; j < gcols; j++ {

			fmt.Print("-")
			fn(i, j)

		}
	}
	fmt.Println()

}

// Output farm detailed info
func main() {
	// Initializing
	randS := rand.NewSource(time.Now().UnixNano())
	randSource = rand.New(randS)
	// Choose language
	for _, loc := range []string{"en", "ua"} {
		fmt.Printf(phrases[loc]["choose_language_string"])
		fmt.Printf(phrases[loc]["languages_string"])
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

L:
	for {

		b := make([]byte, 1)
		_, err = os.Stdin.Read(b)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch string(b) {

		case "1":
			locale = "en"
			break L
		case "2":
			locale = "ua"
			break L
		case "3":
			return
		}

	}
	term.Restore(int(os.Stdin.Fd()), oldState)

	//
	var farm Farm
	farm.genPets(maxPets, minPets)
	farm.showInfo()

}
