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
)

var (
	// App`s settings
	locale       = "en"     // default output language
	grows, gcols = 1, 80    // row delimitor's vars
	sleepInt     = 50       // row delimitor's average timeout in milliseconds
	randSource   *rand.Rand // random seed

	// Animals list
	pets = []pet{dog{}, cat{}, cow{}}

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
type eater interface {
	foodNeded() int
}

type showInfo interface {
	showInfo()
}

type petCreator interface {
	giveBirth(weight float64) pet
}
type animal struct {
	weight float64
}

type haveWeight interface {
	getWeight() float64
}

type pet interface {
	haveWeight
	eater
	showInfo
	petCreator
}

// Custom function of rounding food weight, int + adding 1 extra spare kg
func foodRound(weight float64) int {
	return int(weight) + 1
}

// Return pet_info string. Make title from first word
func getPetInfo(typeName string, p pet) string {
	output := fmt.Sprintf(locales[locale]["pet_info"], locales[locale][typeName], p.getWeight(), p.foodNeded())
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
type cat animal

func (c cat) foodNeded() int {
	return foodRound(c.weight * catFoodPerMonthPerKg)
}

func (c cat) showInfo() {
	fmt.Println(getPetInfo("cat", c))
}

func (c cat) giveBirth(weight float64) pet {
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

func (d dog) showInfo() {
	fmt.Println(getPetInfo("dog", d))

}

func (d dog) giveBirth(weight float64) pet {

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

func (c cow) showInfo() {
	fmt.Println(getPetInfo("cow", c))
}

func (c cow) giveBirth(weight float64) pet {
	return cow{weight: genWeight(weight, cowMinWeight, cowMaxWeight)}
}

func (c cow) getWeight() float64 {
	return c.weight
}

// Farm declaration
type farm struct {
	pets []pet
}

// Show farm details
func (f *farm) detailedInfo() int {

	gsteps := grows * gcols

	// Animals output loop
	if f.pets != nil {

		for i, pet := range f.pets {

			fmt.Printf("%v) ", i+1)
			pet.showInfo()
		}
	}

	prettyBarsProcessOutput(1, gcols, func(i, j int) {})

	fmt.Printf(locales[locale]["calc_farm_info"])

	// Calculating summary foods needed
	var foodSum, ind int

	prettyBarsProcessOutput(grows, gcols, func(i, j int) {

		if f.pets != nil && len(f.pets) > 0 {

			time.Sleep(time.Millisecond * time.Duration(randSource.Intn(sleepInt)))

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

// Generate random farm
func (f *farm) genPets(max, min int) {

	f.pets = make([]pet, 0, max)
	numPets := randSource.Intn(max-min) + min

	gsteps := grows * gcols

	fmt.Printf(locales[locale]["gen_farm"])

	// Farm generation process with indication
	prettyBarsProcessOutput(grows, gcols, func(i int, j int) {

		time.Sleep(time.Millisecond * time.Duration(randSource.Intn(sleepInt)))

		if (len(f.pets)*100)/numPets < ((i+1)*(j+1)*100)/gsteps {

			f.pets = append((f.pets), pets[randSource.Intn(len(pets))].giveBirth(0))
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

func renderMenu() {

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
	defer term.Restore(int(os.Stdin.Fd()), oldState)

LANG_SELECT_LOOP:

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
					break LANG_SELECT_LOOP
				}
			}
		}

	}

	term.Restore(int(os.Stdin.Fd()), oldState)
}

// Main logic
func main() {

	// Initializing
	randS := rand.NewSource(time.Now().UnixNano())
	randSource = rand.New(randS)

	// renderMenu
	renderMenu()

	// Generate and show farm
	var f farm
	f.genPets(maxPets, minPets)
	foodSum := f.detailedInfo()

	// Return summary
	fmt.Printf(locales[locale]["ffood_info"], foodSum, len(f.pets))
}
