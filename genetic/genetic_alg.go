package genetic

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


func calcFitness(g MyGenome) float64 {
	verts := []int{}
	for i, b := range g.Gene {
		if b == 1 {
			verts = append(verts, i)
		}
	}

	var total float64
	for i := 0; i < len(verts); i++ {
		for j := i+1; j < len(verts); j++ {
			//dist, _ := distanceTable.Search(i-1, i)
			//total += dist
		}
	}

	return total / float64(len(verts))
}



func GeneticAlgorithm(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erro ao abrir o ficheiro: ", err)
		os.Exit(1)
	}

	verts := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		node1, node2, cost := 0, 0, 0.0
		fmt.Sscanf(line, "%d %d %f", &node1, &node2, &cost)

		if !contains(verts, node1) {
			verts = append(verts, node1)
		}
		if !contains(verts, node2) {
			verts = append(verts, node2)
		}

		//distanceTable.Add(node1, node2, cost)
	}
	file.Close()

	//distanceTable.Print()

	param := Parameter{
		Initializer: new(RandomInitializer),
		Selector:    new(RandomSelector),
		Mutator:     new(RandomBitFlipMutator),
		Breeder:     new(GA2PointBreeder),
		PMutate: 0.1,
		PBreed: 0.7,
	}

	simulation := NewGA(param)

	fmt.Println("VERTS: ", verts)
	genome := NewMyGenome(len(verts), calcFitness)

	simulation.Init(100, genome)
	fmt.Println("\n\n")
	simulation.PrintTop(10)
	simulation.Optimize(100)
	simulation.PrintTop(10)
}

type Simulation struct {
	initializer Initializer
	selector    Selector
	mutator     Mutator
	breeder     Breeder

	// Parametros
	breedProb      float64
	mutationProb   float64
	populationSize int

	population Population
}

func (sim *Simulation) Init(genome Genome, popSize int) {
	sim.population = sim.initializer.InitPop(genome, popSize)
	sim.populationSize = popSize
}

func (sim *Simulation) Simulate(gen int) {
	for i := 0; i < gen; i++ {
		l := len(sim.population)
		for p := 0; p < l; p++ {
			if sim.breedProb > rand.Float64() {
				children := make(Population, 2)
				children[0], children[1] = sim.breeder.Breed(
					sim.selector.Select(sim.population),
					sim.selector.Select(sim.population))

				sim.population = append(sim.population, children[0])
				sim.population = append(sim.population, children[1])

			}

			if sim.mutationProb > rand.Float64() {
				children := make(Population, 1)
				children[0] = sim.mutator.Mutate(sim.population[p])
				sim.population = append(children)
			}
		}

		// oderna a populacao de acordo com o seu fitness e elimina os 3 piores
		sort.Sort(Population(sim.population))
		sim.population = sim.population[0:sim.populationSize]
	}
}

func (sim Simulation) PrintTop(n int) {
	sort.Sort(sim.population)
	if len(sim.population) < n {
		for i := 0; i < len(sim.population); i++ {
			fmt.Printf("%2d: %s Score = %f", i, sim.population[i], sim.population[i].Fitness())
		}
		return
	}
	for i := 0; i < n; i++ {
		fmt.Printf("%2d: %s Score = %f", i, sim.population[i], sim.population[i].Fitness())
	}
}