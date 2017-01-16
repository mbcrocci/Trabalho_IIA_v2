package main

import (
	"fmt"
	"math/rand"
	"sort"
	"math"
)

type Parameter struct {
	//Chance of breeding
	PBreed float64
	//Chance of mutation
	PMutate float64

	// Initializer, Selector, Mutator, Breeder Objects this GA will use
	Initializer Initializer
	Selector    Selector
	Mutator     Mutator
	Breeder     Breeder
}

type GA struct {
	pop     Population
	popsize int

	Parameter Parameter
}

func NewGA(parameter Parameter) *GA {
	ga := new(GA)
	ga.Parameter = parameter
	return ga
}

func (ga *GA) String() string {
	return fmt.Sprintf("Initializer = %s, Selector = %s, Mutator = %s Breeder = %s",
		ga.Parameter.Initializer,
		ga.Parameter.Selector,
		ga.Parameter.Mutator,
		ga.Parameter.Breeder)
}

func (ga *GA) Init(popsize int, i Genome) {
	ga.pop = ga.Parameter.Initializer.InitPop(i, popsize)
	ga.popsize = popsize
}

func (ga *GA) Optimize(gen int) {
	for i := 0; i < gen; i++ {
		l := len(ga.pop)
		for p := 0; p < l; p++ {
			if ga.Parameter.PBreed > rand.Float64() {
				children := make(Population, 2)
				children[0], children[1] = ga.Parameter.Breeder.Breed(
					ga.Parameter.Selector.Select(ga.pop),
					ga.Parameter.Selector.Select(ga.pop))

				ga.pop = AppendGenomes(ga.pop, children)
			}
			//Mutate
			if ga.Parameter.PMutate > rand.Float64() {
				children := make(Population, 1)
				children[0] = ga.Parameter.Mutator.Mutate(ga.pop[p])
				ga.pop = AppendGenomes(ga.pop, children)
			}
		}
		//cleanup remove some from pop
		// this should probably use a type of selector
		sort.Sort(Population(ga.pop))
		ga.pop = ga.pop[0:ga.popsize]
	}
}

func (ga *GA) OptimizeUntil(stop func(best Genome) bool) {
	for !stop(ga.Best()) {
		ga.Optimize(1)
	}
}

func (ga *GA) Best() Genome {
	sort.Sort(ga.pop)
	best := ga.pop[0]
	for _, pop := range ga.pop {
		if !math.IsNaN(pop.Fitness()) && pop.Valid() {
			best = pop
			break
		}
	}
	return best
}

func (ga *GA) PrintTop(n int) {
	sort.Sort(ga.pop)
	if len(ga.pop) < n {
		for i := 0; i < len(ga.pop); i++ {
			fmt.Printf("%2d: %s Score = %f\n", i, ga.pop[i], ga.pop[i].Fitness())
		}
		return
	}
	for i := 0; i < n; i++ {
		fmt.Printf("%2d: %s Score = %f\n", i, ga.pop[i], ga.pop[i].Fitness())
	}
}

func (ga *GA) PrintPop() {
	fmt.Printf("Current Population:\n")
	for i := 0; i < len(ga.pop); i++ {
		fmt.Printf("%2d: %s Score = %f\n", i, ga.pop[i], ga.pop[i].Fitness())
	}
}
