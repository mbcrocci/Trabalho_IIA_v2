package genetic

import "math/rand"

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