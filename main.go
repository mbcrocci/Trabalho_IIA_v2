package main

import (
	"math/rand"
	"os"
	"time"
	"fmt"
	"bufio"
	"math"
	"strconv"
)
var (
	distanceTable DistanceTable
 	size int
)

func readFile(filename string) DistanceTable {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erro ao abrir o ficheiro: ", err)
		os.Exit(1)
	}
	defer file.Close()

	localTable := DistanceTable{}
	verts := []int{}
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

		localTable.Add(node1, node2, cost)
	}
	size = len(verts)
	return localTable
}

func Distance(s []int) float64 {
	verts := []int{}
	for i := 1; i <= len(s); i++ {
		if s[i-1] == 1 {
			verts = append(verts, i)
		}
	}

	var total float64
	var count int
	for i := 0; i < len(verts); i++{
		for j := i+1; j < len(verts); j++ {
			dist, _ := distanceTable.Search(verts[i], verts[j])
			total += dist
			count++
		}
	}

	return  total / float64(count)
}

func Neighbour(s []int) []int {
	n := make([]int, len(s))
	copy(n, s)
	i := rand.Intn(len(s))

	if n[i] == 0 {
		n[i] = 1
	} else {
		n[i] = 0
	}

	return n
}

type Edge struct {
	From, To int
	Dist float64
}

type Edge2 string

//type DistanceTable []Edge

type DistanceTable map[string]float64

func (d *DistanceTable) Add(n1, n2 int, dist float64) {
	edge := strconv.Itoa(n1) + "-" + strconv.Itoa(n2)
	//*d = append(*d, Edge{n1,n2, dist})
	(*d)[edge] = dist
}

// Search devolve o custo de ir de n1 para n2
// Devolve um false caso essa ligacao nao exista
func (d DistanceTable) Search(n1, n2 int) (float64, bool) {
	/*for _, edge := range d {
		if (edge.From == n1 && edge.To == n2) || (edge.From == n2 && edge.To == n1) {
			return edge.Dist, true
		}
	}*/
	edge := strconv.Itoa(n1) + "-" + strconv.Itoa(n2)
	if val, ok := d[edge]; ok {
		return val, ok
	}
	return 0.0, false
}

func (d DistanceTable) Print() {
	fmt.Println("Distance Table")
	for edge, dist := range d {
		fmt.Println(edge, " = ", dist)
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GraphAlg(filename string, alp float64, it int) {
	distanceTable = readFile(filename)

	solution := make([]int, size)
	for i, _ := range solution {
		solution[i] = rand.Intn(2)
	}
	annealing := func (solution []int) ([]int, float64, []int, float64) {
		var bestSolution []int
		bestDistance := 0.0
		oldDist := Distance(solution)
		T := 1.0
		TMin := 0.0
		alpha := alp
		for T > TMin {
			for i := 1; i <= it; i++ {
				newSolution := Neighbour(solution)
				newDist := Distance(newSolution)

				ap := math.Exp((oldDist - newDist) / T)
				if ap > rand.Float64() {
					solution = newSolution
					oldDist = newDist
				}

				if newDist > bestDistance {
					bestSolution = newSolution
					bestDistance = newDist
				}
			}
			T = T - alpha
		}
		return solution, oldDist, bestSolution, bestDistance
	}

	solution, oldfitness, best, bestfit := annealing(solution)

	fmt.Println("Solution: ", solution)
	fmt.Println("Fitness: ", oldfitness)

	fmt.Println("Best: ", best)
	fmt.Println("BESTFIT: ", bestfit)
}

func fitness(g MyGenome) float64 {
	verts := []int{}
	for i := 1; i <= len(g.Gene); i++ {
		if g.Gene[i-1] == 1 {
			verts = append(verts, i)
		}
	}

	var total float64
	var count int
	for i := 0; i < len(verts); i++{
		for j := i+1; j < len(verts); j++ {
			dist, _ := distanceTable.Search(verts[i], verts[j])
			total += dist
			count++
		}
	}

	return  total / float64(count)
}


func GeneticAlg(filename string, gens int, pmut, pbreed float64) {
	distanceTable = readFile(filename)

	param := Parameter{
		Initializer: new(RandomInitializer),
		Selector: new(TournamentSelector),
		Breeder: new(GA2PointBreeder),
		Mutator: new(RandomBitFlipMutator),
		PMutate: pmut,
		PBreed: pbreed,
	}

	ga := NewGA(param)

	genome := NewMyGenome(size, fitness)

	ga.Init(100, genome)
	ga.Optimize(gens)

	//ga.PrintTop(10)
	best := ga.Best()
	fmt.Println("\nBest: ", best, " -> Fit: ", best.Fitness())
}

func HybridAlg(filename string, gens int, pmut, pbreed float64, alp float64, it int) {
	distanceTable = readFile(filename)


	param := Parameter{
		Initializer: new(RandomInitializer),
		Selector: new(RandomSelector),
		Breeder: new(GA2PointBreeder),
		Mutator: new(RandomBitFlipMutator),
		PMutate: pmut,
		PBreed: pbreed,
	}

	ga := NewGA(param)

	genome := NewMyGenome(size, fitness)

	ga.Init(100, genome)
	ga.Optimize(gens)

	//ga.PrintTop(10)
	bestga := ga.Best().(MyGenome)
	solution := bestga.Gene

	annealing := func (solution []int) ([]int, float64, []int, float64) {
		var bestSolution []int
		bestDistance := 0.0
		oldDist := Distance(solution)
		T := 1.0
		TMin := 0.0
		alpha := alp
		for T > TMin {
			for i := 1; i <= it; i++ {
				newSolution := Neighbour(solution)
				newDist := Distance(newSolution)

				ap := math.Exp((oldDist - newDist) / T)
				if ap > rand.Float64() {
					solution = newSolution
					oldDist = newDist
				}

				if newDist > bestDistance {
					bestSolution = newSolution
					bestDistance = newDist
				}
			}
			T = T - alpha
		}
		return solution, oldDist, bestSolution, bestDistance
	}

	solution, oldfitness, best, bestfit := annealing(solution)

	fmt.Println("Solution: ", solution)
	fmt.Println("Fitness: ", oldfitness)

	fmt.Println("Best: ", best)
	fmt.Println("BESTFIT: ", bestfit)
}

func menu() int {
	fmt.Println("1 - Graph Alg")
	fmt.Println("2 - Genetic Alg")
	fmt.Println("3 - Hybrid Alg")
	fmt.Println("4 - All GraphsAlg")
	fmt.Println("5 - All Genetic Alg")
	fmt.Println("6 - All Hybrid Alg")


	var op int
	for op < 1 || 6 < op {
		fmt.Print("Option: ")
		fmt.Scanf("%d", &op)
	}
	return op


}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	op := menu()
	switch op {
	case 1:
		GraphAlg(os.Args[1], 0.00001, 100)
	case 2:
		GeneticAlg(os.Args[1], 100, 0.1, 0.7)
	case 3:
		HybridAlg(os.Args[1], 10, 0.1, 0.7, 0.00001, 100)
	case 4:
		test_all_graph()
	case 5:
		test_all_genetic()
	case 6:
		test_all_hybrid()
	}
}
func test_all_graph() {
	filenames := []string{
		"instancias/MDPI1_20.txt",
		"instancias/MDPI1_150.txt",
		"instancias/MDPI1_500.txt",
		"instancias/MDPI2_25.txt",
		"instancias/MDPI2_30.txt",
		"instancias/MDPI2_150.txt",
		"instancias/MDPII1_20.txt",
		"instancias/MDPII1_150.txt",
		"instancias/MDPII2_25.txt",
		"instancias/MDPII2_30.txt",
		"instancias/MDPII2_150.txt",
		"instancias/MDPII2_500.txt",
	}

	for _, filename := range filenames {
		fmt.Print(filename, ": ")
		GraphAlg(filename, 0.00001, 100)
	}
}

func test_all_genetic() {
	filenames := []string{
		"instancias/MDPI1_20.txt",
		"instancias/MDPI2_25.txt",
		"instancias/MDPI2_30.txt",
		"instancias/MDPII1_20.txt",
		"instancias/MDPII2_25.txt",
		"instancias/MDPII2_30.txt",
		"instancias/MDPI1_150.txt",
		"instancias/MDPI1_500.txt",
		"instancias/MDPII1_150.txt",
		"instancias/MDPI2_150.txt",
		"instancias/MDPII2_150.txt",
		"instancias/MDPII2_500.txt",
	}

	for _, filename := range filenames {
		fmt.Print(filename, ": ")
		GeneticAlg(os.Args[1], 10, 0.1, 0.7)
	}
}

func test_all_hybrid() {
	filenames := []string{
		"instancias/MDPI1_20.txt",
		"instancias/MDPI2_25.txt",
		"instancias/MDPI2_30.txt",
		"instancias/MDPII1_20.txt",
		"instancias/MDPII2_25.txt",
		"instancias/MDPII2_30.txt",
		"instancias/MDPI1_150.txt",
		"instancias/MDPI1_500.txt",
		"instancias/MDPII1_150.txt",
		"instancias/MDPI2_150.txt",
		"instancias/MDPII2_150.txt",
		"instancias/MDPII2_500.txt",
	}

	for _, filename := range filenames {
		fmt.Print(filename, ": ")
		HybridAlg(os.Args[1], 10, 0.1, 0.7, 0.00001, 100)
	}
}

