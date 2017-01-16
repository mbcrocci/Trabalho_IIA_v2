package main

import (
	"math/rand"
)

type Mutator interface {
	Mutate(a Genome) Genome
	String() string
}

type RandomBitFlipMutator struct{}

func (m RandomBitFlipMutator) Mutate(genome Genome) Genome {
	n := genome.Copy()
	n.Flip(rand.Intn(genome.Len()))
	return n
}

func (m RandomBitFlipMutator) String() string { return "RandomBitFlipMutator" }
