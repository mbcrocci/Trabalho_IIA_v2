package genetic

import (
	"fmt"
	"math/rand"
)

// Genome interface, Not final.
type Genome interface {
	//Randomize.Genens
	Randomize()
	//Copy a genome;
	Copy() Genome
	//Calculate score
	Fitness() float64
	//Crossover for this genome
	Crossover(bi Genome, p1, p2 int) (ca Genome, cb Genome)
	Flip(index int)

	//Check if genome is valid
	Valid() bool

	String() string
	Len() int
}

type Population []Genome

func (g Population) Len() int           { return len(g) }
func (g Population) Less(i, j int) bool { return g[i].Fitness() < g[j].Fitness() }
func (g Population) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }

func AppendGenomes(slice, data Population) Population {
	l := len(slice)
	if l+len(data) > cap(slice) {
		newSlice := make(Population, (l+len(data))*2)
		for i, c := range slice {
			newSlice[i] = c
		}
		slice = newSlice
	}
	slice = slice[0 : l+len(data)]
	for i, c := range data {
		slice[l+i] = c
	}
	return slice
}

type MyGenome struct {
	Gene    []int
	fitness float64

	fitnessFunc func(genome MyGenome) float64
}

func NewMyGenome(size int, fitnessFunc func(genome MyGenome) float64) MyGenome {
	g := MyGenome{}
	g.Gene = make([]int, size)

	fmt.Println("MADE THE GENES: ", g.Gene)
	g.Randomize()
	g.fitness = 0
	g.fitnessFunc = fitnessFunc

	return g
}

func (g MyGenome) Copy() Genome {
	n := MyGenome{}
	n.Gene = make([]int, len(g.Gene))
	copy(n.Gene, g.Gene)
	n.fitnessFunc = g.fitnessFunc
	n.fitness = g.fitness

	return n
}

func (g MyGenome) Crossover(bi Genome, p1, p2 int) (Genome, Genome) {
	ca := g.Copy().(MyGenome)
	b := bi.(MyGenome)
	cb := b.Copy().(MyGenome)

	copy(ca.Gene[p1:p2 + 1], b.Gene[p1:p2 + 1])
	copy(cb.Gene[p1:p2 + 1], g.Gene[p1:p2 + 1])

	return ca, cb
}

func (g MyGenome) Randomize() {
	for i := 0; i < len(g.Gene); i++ {
		g.Gene[i] = rand.Intn(2)
	}
}

func (g MyGenome) Fitness() float64 {
	g.fitness = g.fitnessFunc(g)
	return g.fitness
}

func (g MyGenome) Flip(index int) {
	if g.Gene[index] == 0 {
		g.Gene[index] = 1
	} else {
		g.Gene[index] = 0
	}
}

func (g MyGenome) Valid() bool {
	verts := []int{}
	for i, b := range g.Gene {
		if b == 1 {
			verts = append(verts, i)
		}
	}
	for i := 1; i < len(verts); i++ {
		//if _, found := distanceTable.Search(i-1, i); !found {
		//	return false
		//}

	}
	return true
}
func (g MyGenome) Len() int { return len(g.Gene) }
func (g MyGenome) String() string { return fmt.Sprintf("%v", g.Gene) }