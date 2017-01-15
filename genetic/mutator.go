package genetic

import (
	"math/rand"
	"fmt"
)

type Mutator interface {
	Mutate(a Genome) Genome
	String() string
}

type RandomBitFlipMutator struct{}

func (m RandomBitFlipMutator) Mutate(genome Genome) Genome {
	n := genome.Copy()
	fmt.Println("Copied this gene: ", n)
	n.Flip(rand.Intn(genome.Len()))
	return n

}

func (m RandomBitFlipMutator) String() string { return "RandomBitFlipMutator" }
