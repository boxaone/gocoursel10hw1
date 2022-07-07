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
	gcols = 80
	grows = 1
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
		},
	}

	return &conf
}

// Common types declarations
type eater interface {
	foodNeded() int
}

type animal struct {
	weight float64
}

type haveWeight interface {
	getWeight() float64
}
type speciable interface {
	getSpecie() string
}

type pet interface {
	haveWeight
	eater
	speciable
}

// Custom function of rounding food weight, int + adding 1 extra spare kg
func foodRound(weight float64) int {
	return int(weight) + 1
}

// Return pet_info string. Make title from first word
func getPetInfo(pet_info string, specie string, weight float64, food int) string {
	output := fmt.Sprintf(pet_info, specie, weight, food)
	output_arr := strings.SplitAfterN(output, " ", 2)

	if len(output_arr) < 2 {
		return output
	}
	return fmt.Sprint(cases.Title(language.Und, cases.NoLower).String(output_arr[0]), output_arr[1])
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
type cat animal

func (c cat) foodNeded() int {
	return foodRound(c.weight * catFoodPerMonthPerKg)
}

func (c cat) getSpecie() string {
	return "cat"
}
func (c cat) giveBirth(weight float64) cat {
	return cat{weight: genWeight(weight, catMinWeight, catMaxWeight)}
}

func (c cat) getWeight() float64 {
	return c.weight
}

// Dog declarations
type dog animal

func (d dog) foodNeded() int {
	return foodRound(d.weight * dogFoodPerMonthPerKg)
}

func (d dog) getSpecie() string {
	return "dog"
}
func (d dog) giveBirth(weight float64) dog {

	return dog{weight: genWeight(weight, dogMinWeight, dogMaxWeight)}

}

func (d dog) getWeight() float64 {
	return d.weight
}

// Cow declarations
type cow animal

func (c cow) foodNeded() int {
	return foodRound(c.weight * cowFoodPerMonthPerKg)
}

func (c cow) giveBirth(weight float64) cow {
	return cow{weight: genWeight(weight, cowMinWeight, cowMaxWeight)}
}
func (c cow) getSpecie() string {
	return "cow"
}

func (c cow) getWeight() float64 {
	return c.weight
}

// Farm declaration
type farm struct {
	pets []pet
}

// Show farm details
func (f *farm) detailedInfo(conf *config) int {

	gsteps := grows * gcols

	// Animals output loop
	if f.pets != nil {

		for i, pet := range f.pets {

			fmt.Printf("%v) %v", i+1, getPetInfo(conf.locales[conf.locale]["pet_info"], conf.locales[conf.locale][pet.getSpecie()], pet.getWeight(), pet.foodNeded()))
		}
	}

	prettyBarsProcessOutput(1, gcols, func(i, j int) {})

	fmt.Printf(conf.locales[conf.locale]["calc_farm_info"])

	// Calculating summary foods needed
	var foodSum, ind int

	prettyBarsProcessOutput(grows, gcols, func(i, j int) {

		if f.pets != nil && len(f.pets) > 0 {

			time.Sleep(time.Millisecond * time.Duration(randSource()().Intn(sleepInt)))

			if (ind*100)/(len(f.pets)) <= ((i+1)*(j+1)*100)/gsteps {

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
func getRandomPet() pet {
	switch randSource()().Intn(3) {
	case 1:
		return dog{}.giveBirth(0)

	case 2:
		return cat{}.giveBirth(0)
	default:
		return cow{}.giveBirth(0)

	}
}

// Generate random farm
func (f *farm) genPets(max, min int, conf *config) {
	rands := randSource()()
	f.pets = make([]pet, 0, max)
	numPets := rands.Intn(max-min) + min

	gsteps := grows * gcols

	fmt.Printf(conf.locales[conf.locale]["gen_farm"])

	// Farm generation process with indication
	prettyBarsProcessOutput(grows, gcols, func(i int, j int) {

		time.Sleep(time.Millisecond * time.Duration(rands.Intn(sleepInt)))

		if (len(f.pets)*100)/numPets < ((i+1)*(j+1)*100)/gsteps {
			f.pets = append((f.pets), getRandomPet())
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

	for loc, _ := range conf.locales {
		languages = append(languages, loc)

	}

	sort.Strings(languages)

	for _, loc := range languages {
		prettyBarsProcessOutput(1, gcols, func(i, j int) {})
		fmt.Printf(conf.locales[loc]["choose_language"])
	}

	// Show delimiter
	prettyBarsProcessOutput(1, gcols, func(i, j int) {})

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

	prettyBarsProcessOutput(1, gcols, func(i, j int) {})

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
	return ""
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
	foodSum := f.detailedInfo(conf)

	// Return summary
	fmt.Printf(conf.locales[conf.locale]["ffood_info"], foodSum, len(f.pets))
}
