package genetic

import (
	"math/rand"
)

type Breeder interface {
	// Breeds two parent Population and returns two children
	Breed(a, b Genome) (ca, cb Genome)
	// String name of breeder
	String() string
}

//Combines genomes by selecting 2 points to exchange
// Example: Parent 1 = 111111, Parent 2 = 000000, Child = 111001
type GA2PointBreeder struct{}

func (breeder *GA2PointBreeder) Breed(a, b Genome) (ca, cb Genome) {
	if a.Len() != b.Len() {
		panic("Length mismatch in pmx")
	}
	p1 := rand.Intn(a.Len())
	p2 := rand.Intn(b.Len())
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	ca, cb = a.Crossover(b, p1, p2)
	return
}

func (b *GA2PointBreeder) String() string { return "GA2PointBreeder" }