package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// forks represented simply by their index at the table
var philosophers = []Philosopher{
	{name: "    Plato", leftFork: 0, rightFork: 1},
	{name: " Socrates", leftFork: 1, rightFork: 2},
	{name: "Aristotle", leftFork: 2, rightFork: 0},
	// {name: "   Pascal", leftFork: 3, rightFork: 4},
	// {name: "Locke", leftFork: 4, rightFork: 0},
}

// Define a few variables.
var hunger = 3                  // how many times a philosopher eats
var eatTime = 1 * time.Second   // how long it takes to eat
var thinkTime = 1 * time.Second // how long it takes to think

var historyLeaveMu sync.Mutex
var historyLeave = []int{}

var historyEatMu sync.Mutex
var historyEat = []int{}

func main() {
	fmt.Println()
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("---------------------------")
	fmt.Println("= The table is empty =")

	dine()

	// print out finished message
	fmt.Println("= The table is empty =")

	fmt.Println()
	fmt.Println("--------- EATING TURNS ----------")
	fmt.Println(historyEat)
	fmt.Println("--------- ------------ ----------")
	fmt.Println()
	fmt.Println("--------- LEFT THE TABLE ---------")
	fmt.Println(historyLeave)
	fmt.Println("--------- ------------ -----------")
}

func dine() {

	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	// We want everyone to be seated before they start eating
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// map each fork to a mutex
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go diningProblem(philosophers[i], wg, forks, seated, i)
	}

	wg.Wait()
}

func diningProblem(
	philosopher Philosopher,
	wg *sync.WaitGroup,
	forks map[int]*sync.Mutex,
	seated *sync.WaitGroup,
	philIdx int,
) {
	defer wg.Done()

	fmt.Printf("  %s is seated.\n", philosopher.name)
	seated.Done()
	// Wait until others are also seated
	seated.Wait()

	for i := hunger; i > 0; i-- {
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			color.Yellow(fmt.Sprintf("\t%s takes fork #%v.\n", philosopher.name, philosopher.rightFork))
			forks[philosopher.leftFork].Lock()
			color.Yellow(fmt.Sprintf("\t%s takes fork #%v.\n", philosopher.name, philosopher.leftFork))
		} else {
			forks[philosopher.leftFork].Lock()
			color.Yellow(fmt.Sprintf("\t%s takes fork #%v.\n", philosopher.name, philosopher.leftFork))
			forks[philosopher.rightFork].Lock()
			color.Yellow(fmt.Sprintf("\t%s takes fork #%v.\n", philosopher.name, philosopher.rightFork))
		}

		historyEatMu.Lock()
		color.Green(fmt.Sprintf("\t\t\t%s is eating...\n", philosopher.name))
		historyEat = append(historyEat, philIdx)
		historyEatMu.Unlock()
		time.Sleep(eatTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		color.Cyan(fmt.Sprintf("\t%s puts down the forks.\n", philosopher.name))
		color.Blue(fmt.Sprintf("\t\t\t%s is thinking...\n", philosopher.name))
		time.Sleep(thinkTime)
	}

	// we only get here if philosopher has satisfied their hunger
	historyLeaveMu.Lock()
	historyLeave = append(historyLeave, philIdx)
	fmt.Println(philosopher.name, "left.")
	historyLeaveMu.Unlock()
}
