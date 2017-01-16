package main

import (
	"math/rand"
	"sort"
)

type Selector interface {
	// Select one from pop
	Select(pop []Genome) Genome

	// String name of selector
	String() string
}

type RandomSelector struct{}

func (s RandomSelector) Select(pop []Genome) Genome {
	i := rand.Intn(len(pop))
	return pop[i]
}

func (s *RandomSelector) String() string {
	return "RandomSelector"
}

type TournamentSelector struct {}

func (s TournamentSelector) Select(pop []Genome) Genome {
	g := make(Population, 5)
	l := len(pop)

	for i := 0; i < 5; i++ {
		g[i] = pop[rand.Intn(l)]
	}

	sort.Sort(g)

	for i := 0; i < 5; i++ {
		if g[i].Valid() {
			return g[i]
		}
	}
	return g[0]
}

func (s TournamentSelector) String() string {
	return "TournamentSelector"
}