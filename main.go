package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	minPets         = 10
	maxPets         = 20
	catFoodPerMonth = 7
	catMinWeight    = 1.
	catMaxWeight    = 7.
	dogFoodPerMonth = 10 / 5
	dogMinWeight    = 1.
	dogMaxWeight    = 20.
	cowFoodPerMonth = 25
	cowMinWeight    = 10.
	cowMaxWeight    = 300.
)

var (
	// App`s settings
	locale       = "en"     // default output language
	grows, gcols = 1, 80    // row delimitor's vars
	sleepInt     = 50       // row delimitor's average timeout in milliseconds
	randSource   *rand.Rand // random seed

	// Animals list
	pets = []Pet{Dog{}, Cat{}, Cow{}}

	// Locales list
	locales = map[string]map[string]string{
		"en": {
			"pet_info":        "A %v, weighting %.2f kg, needs %v kg of food per month.\n",
			"choose_language": "Please choose output language or choose exit\n",
			"language":        "English",
			"exit":            "Exit",
			"gen_farm":        "Generating new farm\n",
			"ffood_info":      "Summary %v kg food per month needed for %v animals.\n",
			"calc_farm_info":  "Calculating all food needed\n",
			"cat":             "cat",
			"dog":             "dog",
			"cow":             "cow",
		},
		"ua": {
			"pet_info":        "%v, важить %.2f кілограм, потребує %v кілограмів кормів на місяць.\n",
			"choose_language": "Будь ласка, оберіть мову або оберіть вихід \n",
			"language":        "Українська",
			"exit":            "Вихід",
			"gen_farm":        "Генеруємо нову ферму\n",
			"ffood_info":      "Загалом потрібно %v кілограмів кормів на місяць для %v тварин.\n",
			"calc_farm_info":  "Рахуємо загальну вагу кормів\n",
			"cat":             "кішка",
			"dog":             "пес",
			"cow":             "корова",
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

type PetCreator interface {
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
	PetCreator
}

// Custom function of rounding food weight, int + adding 1 extra spare kg
func foodRound(weight float64) int {
	return int(weight) + 1
}

// Return pet_info string
func getPetInfo(t string, p Pet) string {
	output := fmt.Sprintf(locales[locale]["pet_info"], locales[locale][t], p.Weight(), p.foodNeded())
	output_arr := strings.SplitAfterN(output, " ", 2)

	if len(output_arr) < 2 {
		return output
	}
	return fmt.Sprint(cases.Title(language.Und, cases.NoLower).String(output_arr[0]), output_arr[1])
}

// Return weight within ranges
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

	return Dog{weight: genWeight(weight, dogMinWeight, dogMaxWeight)}

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

	// Animals output loop
	if f.Pets != nil {

		for i, pet := range f.Pets {

			fmt.Printf("%v) ", i+1)
			pet.showInfo()
		}
	}

	prettyBarsProcessOutput(1, gcols, func(i, j int) {})

	fmt.Printf(locales[locale]["calc_farm_info"])

	// Calculating summary foods needed
	var foodSum, ind int

	prettyBarsProcessOutput(grows, gcols, func(i, j int) {

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

	// Return summary
	fmt.Printf(locales[locale]["ffood_info"], foodSum, len(f.Pets))
}

// Generate random farm
func (f *Farm) genPets(max, min int) {

	f.Pets = make([]Pet, 0, max)
	numPets := randSource.Intn(max-min) + min

	gsteps := grows * gcols

	fmt.Printf(locales[locale]["gen_farm"])

	// Farm generation process with indication
	prettyBarsProcessOutput(grows, gcols, func(i int, j int) {

		time.Sleep(time.Millisecond * time.Duration(randSource.Intn(sleepInt)))

		if (len(f.Pets)*100)/numPets < ((i+1)*(j+1)*100)/gsteps {

			f.Pets = append((f.Pets), pets[randSource.Intn(len(pets))].giveBirth(0))
		}

	})

}

// Function for output bars during processes
func prettyBarsProcessOutput(grows int, gcols int, fn func(int, int)) {
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

// Main logic
func main() {

	// Initializing
	randS := rand.NewSource(time.Now().UnixNano())
	randSource = rand.New(randS)

	// Choose language
	languages := []string{}

	for loc, _ := range locales {
		languages = append(languages, loc)

	}

	sort.Strings(languages)

	for _, loc := range languages {
		prettyBarsProcessOutput(1, gcols, func(i, j int) {})
		fmt.Printf(locales[loc]["choose_language"])
	}

	prettyBarsProcessOutput(1, gcols, func(i, j int) {})

	// Output exits
	fmt.Print("0) ")

	for i, loc := range languages {
		if i != 0 {
			fmt.Print("/")
		}
		fmt.Print(locales[loc]["exit"])
	}

	fmt.Println()

	// Output languages
	for i, loc := range languages {
		fmt.Printf("%v) %v\n", i+1, locales[loc]["language"])
	}

	prettyBarsProcessOutput(1, gcols, func(i, j int) {})

	// Console manipulations
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

LANG_INPUT_LOOP:

	for {

		b := make([]byte, 1)
		_, err = os.Stdin.Read(b)
		if err != nil {
			fmt.Println(err)
			return
		}

		l, err := strconv.Atoi(string(b))

		if err == nil {
			switch l {

			case 0:
				return
			default:

				if l <= len(languages) {

					locale = languages[l-1]
					break LANG_INPUT_LOOP
				}
			}
		}

	}

	term.Restore(int(os.Stdin.Fd()), oldState)

	// Generate and show farm
	var farm Farm
	farm.genPets(maxPets, minPets)
	farm.showInfo()

}
