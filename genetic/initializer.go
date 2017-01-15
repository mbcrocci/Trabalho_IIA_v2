package genetic

type Initializer interface {
	// Initializes popsize length []GAGenome from i
	InitPop(i Genome, popsize int) []Genome
	// String name of initializers
	String() string
}

type RandomInitializer struct{}

func (i *RandomInitializer) InitPop(first Genome, popsize int) (pop []Genome) {
	pop = make([]Genome, popsize)
	for x := 0; x < popsize; x++ {
		pop[x] = first.Copy()
		pop[x].Randomize()
	}
	return pop
}

func (i *RandomInitializer) String() string { return "RandomInitializer" }
