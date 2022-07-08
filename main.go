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
	minPets              = 10
	maxPets              = 20
	catFoodPerMonthPerKg = 7
	catMinWeight         = 1.
	catMaxWeight         = 7.
	dogFoodPerMonthPerKg = 10 / 5
	dogMinWeight         = 1.
	dogMaxWeight         = 20.
	cowFoodPerMonthPerKg = 25
	cowMinWeight         = 10.
	cowMaxWeight         = 300.
	// Row delimitors'
	gCols = 80
	gRows = 1
	// Sleep between actions in milliseconds
	sleepInt = 50
)

type config struct {
	locale  string // locale
	locales map[string]map[string]string
}

// Return initialized rand
func randSource() func() *rand.Rand {

	randS := rand.NewSource(time.Now().UnixNano())

	return func() *rand.Rand {
		return rand.New(randS)
	}
}

func returnInitialConfig() *config {

	conf := config{
		locale: "en",

		locales: map[string]map[string]string{
			"en": {
				"pet_info":        "%v, weighting %.2f kg, needs %v kg of food per month.\n",
				"choose_language": "Please choose output language or choose exit\n",
				"language":        "English",
				"exit":            "Exit",
				"gen_farm":        "Generating new farm\n",
				"ffood_info":      "Summary %v kg food per month needed for %v animals.\n",
				"calc_farm_info":  "Calculating all food needed\n",
				"cat":             "cat",
				"dog":             "dog",
				"cow":             "cow",
				"usual_name":      "a %v",
				"special_name":    "a %v named %v",
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
				"usual_name":      "%v",
				"special_name":    "%v на ім'я %v",
			},
		},
	}

	return &conf
}

// Common types declarations
type eater interface {
	foodNeded() int
}

type animal struct {
	weight         float64
	name           string
	specie         string
	foodPerMonthKg float64
}

func (a animal) foodNeded() int {
	return foodRound(a.weight * a.foodPerMonthKg)
}

func (a animal) getWeight() float64 {
	return a.weight
}

func (a animal) getSpecie() string {
	return a.specie
}

func (a animal) getName() string {

	if a.name != "" {

		return a.name
	}
	return a.specie
}

type haveWeight interface {
	getWeight() float64
}
type namable interface {
	getSpecie() string
	getName() string
}

type pet interface {
	haveWeight
	eater
	namable
}

// Custom function of rounding food weight, int + adding 1 extra spare kg
func foodRound(weight float64) int {
	return int(weight) + 1
}

// Return pet_info string. Make title from first word
func getPetInfo(petInfo string, specie string, weight float64, food int) string {
	output := fmt.Sprintf(petInfo, specie, weight, food)
	outputArr := strings.SplitAfterN(output, " ", 2)

	if len(outputArr) < 2 {
		return output
	}
	return fmt.Sprint(cases.Title(language.Und, cases.NoLower).String(outputArr[0]), outputArr[1])
}

// Return weight within ranges
func genWeight(weight float64, min float64, max float64) float64 {

	if weight == 0 {
		return (max-min)*randSource()().Float64() + min
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
type cat struct {
	animal
}

func (c cat) giveBirth(weight float64, petName string) cat {
	return cat{animal: animal{weight: genWeight(weight, catMinWeight, catMaxWeight), foodPerMonthKg: catFoodPerMonthPerKg, specie: "cat", name: petName}}
}

// Dog declarations
type dog struct {
	animal
}

func (d dog) giveBirth(weight float64, petName string) dog {

	return dog{animal: animal{weight: genWeight(weight, dogMinWeight, dogMaxWeight), foodPerMonthKg: dogFoodPerMonthPerKg, specie: "dog", name: petName}}

}

// Cow declarations
type cow struct {
	animal
}

func (c cow) giveBirth(weight float64, petName string) cow {
	return cow{animal: animal{weight: genWeight(weight, cowMinWeight, cowMaxWeight), foodPerMonthKg: cowFoodPerMonthPerKg, specie: "cow", name: petName}}

}

// Farm declaration
type farm struct {
	pets []pet
}

// Show farm details
func (f *farm) monthlyFarmFoodWeightDetailed(conf *config) int {

	gSteps := gRows * gCols
	rands := randSource()()

	// Animals output loop
	if f.pets != nil {

		for i, pet := range f.pets {
			name := pet.getName()
			specie := conf.locales[conf.locale][pet.getSpecie()]

			nameString := fmt.Sprintf(conf.locales[conf.locale]["usual_name"], specie)

			if name != specie && name != pet.getSpecie() {

				nameString = fmt.Sprintf(conf.locales[conf.locale]["special_name"], specie, name)

			}

			fmt.Printf("%v) %v", i+1, getPetInfo(conf.locales[conf.locale]["pet_info"], nameString, pet.getWeight(), pet.foodNeded()))
		}
	}

	prettyBarsProcessOutput(1, gCols, func(i, j int) {})

	fmt.Printf(conf.locales[conf.locale]["calc_farm_info"])

	// Calculating summary foods needed
	var foodSum, ind int

	prettyBarsProcessOutput(gRows, gCols, func(i, j int) {

		if f.pets != nil && len(f.pets) > 0 {

			time.Sleep(time.Millisecond * time.Duration(rands.Intn(sleepInt)))

			if (ind*100)/(len(f.pets)) <= ((i+1)*(j+1)*100)/gSteps {

				if ind < len(f.pets) {
					foodSum += (f.pets)[ind].foodNeded()
				}
				ind++
			}
		}

	})

	return foodSum
}

// Get random pet
func getRandomPet(name string) pet {

	switch randSource()().Intn(3) {
	case 1:
		return dog{}.giveBirth(0, name)

	case 2:
		return cat{}.giveBirth(0, name)
	default:
		return cow{}.giveBirth(0, name)

	}
}

// Generate random farm
func (f *farm) genPets(max, min int, conf *config) {

	// Initialization
	rands := randSource()()
	f.pets = make([]pet, 0, max)
	numPets := rands.Intn(max-min) + min
	gSteps := gRows * gCols

	// Farm generation process with indication
	fmt.Printf(conf.locales[conf.locale]["gen_farm"])

	prettyBarsProcessOutput(gRows, gCols, func(i int, j int) {

		// Sleep for more visual effects
		time.Sleep(time.Millisecond * time.Duration(rands.Intn(sleepInt)))

		// Calculating processed pets and displayed steps
		petsProcessedPercent := (len(f.pets) * 100) / numPets
		stepsProcessed := (i + 1) * (j + 1)
		stepsProcessedPercent := (stepsProcessed * 100) / gSteps

		// Process pet if pets percent less then steps percent
		if petsProcessedPercent < stepsProcessedPercent {

			f.pets = append((f.pets), getRandomPet("Зірочка"))
		}

	})

}

// Manage cursor movements
func manageCursor(commands ...string) {

	for _, command := range commands {
		switch command {
		case "saveCursorPosition":
			fmt.Print("\033[s")
		case "eraseToEndOfLine":
			fmt.Print("\033[K")
		case "restoreCursorPosition":
			fmt.Print("\033[u")
		}
	}
}

// Function for output bars during processes
func prettyBarsProcessOutput(grows int, gcols int, fn func(int, int)) {

	manageCursor("saveCursorPosition")

	for i := 0; i < grows; i++ {

		manageCursor("restoreCursorPosition", "eraseToEndOfLine")

		for j := 0; j < gcols; j++ {

			fmt.Print("-")
			fn(i, j)

		}
	}
	fmt.Println()

}

func renderMenu(conf *config) {

	// Filter languages
	languages := []string{}

	for loc := range conf.locales {
		languages = append(languages, loc)

	}

	sort.Strings(languages)

	for _, loc := range languages {
		prettyBarsProcessOutput(1, gCols, func(i, j int) {})
		fmt.Printf(conf.locales[loc]["choose_language"])
	}

	// Show delimiter
	prettyBarsProcessOutput(1, gCols, func(i, j int) {})

	// Output exits
	fmt.Print("0) ")

	for i, loc := range languages {
		if i != 0 {
			fmt.Print("/")
		}
		fmt.Print(conf.locales[loc]["exit"])
	}
	fmt.Println()

	// Output languages
	for i, loc := range languages {
		fmt.Printf("%v) %v\n", i+1, conf.locales[loc]["language"])
	}

	prettyBarsProcessOutput(1, gCols, func(i, j int) {})

	locale := processInput(&languages)

	if len(locale) > 0 {
		conf.locale = locale
	}

}
func processInput(languages *[]string) string {

	// Console manipulations
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {

		fmt.Println(err)
		panic("Error changing term to raw")
	}

	restore := func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}

	defer restore()

	for {
		b := make([]byte, 1)
		_, err = os.Stdin.Read(b)

		if err != nil {

			fmt.Println(err)
			panic("Error retrieving key")

		}

		l, err := strconv.Atoi(string(b))

		if err == nil {
			switch l {

			case 0:
				restore()
				os.Exit(0)

			default:

				if l <= len(*languages) {

					return (*languages)[l-1]
				}
			}
		}

	}
}

// Main logic
func main() {

	// Initializing
	conf := returnInitialConfig()

	// renderMenu
	renderMenu(conf)

	// Generate and show farm
	var f farm
	f.genPets(maxPets, minPets, conf)
	foodSum := f.monthlyFarmFoodWeightDetailed(conf)

	// Return summary
	fmt.Printf(conf.locales[conf.locale]["ffood_info"], foodSum, len(f.pets))
}
