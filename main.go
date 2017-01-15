package main

import (
	"math/rand"
	"os"
	"time"
	"fmt"
	"bufio"
	"runtime"
	"math"
	"github.com/mbcrocci/Trabalho_IIA_v2/genetic"
)
var (
	distanceTable DistanceTable
 	size int
)

type Solution []int
func (s Solution) Distance() float64 {
	verts := []int{}
	for i := 1; i < len(s); i++ {
		if s[i-1] == 1 {
			verts = append(verts, i)
		}
	}

	var total float64
	var count float64
	for i := 0; i < len(verts); i++{
		for j := i+1; j < len(verts); j++ {
			dist, _ := distanceTable.Search(verts[i], verts[j])
			total += dist
			count++
		}
	}

	return  total / count
}

func (s Solution) Neighbour() Solution {
	n := s
	i := rand.Intn(len(s))
	j := rand.Intn(len(s))

	n[i], n[j] = n[j], n[i]

	return n
}

type Edge struct {
	From, To int
	Dist float64
}

type DistanceTable []Edge

func (d *DistanceTable) Add(n1, n2 int, dist float64) {
	if _, found := (*d).Search(n1, n2); !found {
		*d = append(*d, Edge{n1,n2, dist})
	}
}

// Search devolve o custo de ir de n1 para n2
// Devolve um false caso essa ligacao nao exista
func (d DistanceTable) Search(n1, n2 int) (float64, bool) {
	for _, edge := range d {
		if (edge.From == n1 && edge.To == n2) || (edge.From == n2 && edge.To == n1) {
			return edge.Dist, true
		}
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

func graphAlg(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erro ao abrir o ficheiro: ", err)
		os.Exit(1)
	}
	defer file.Close()

	distanceTable = DistanceTable{}
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

		distanceTable.Add(node1, node2, cost)
	}
	file.Close()

	//distanceTable.Print()
	size = len(verts)

	solution := make([]int, size)

	for i, _ := range solution {
		solution[i] = rand.Intn(2)
	}

	annealing := func (solution Solution) (Solution, float64) {
		oldDist := solution.Distance()
		T := 1.0
		TMin := 0.0
		alpha := 0.00001
		for T > TMin {
			for i := 1; i <= 100; i++ {
				newSolution := solution.Neighbour()
				newDist := newSolution.Distance()

				ap := math.Exp((oldDist - newDist) / T)
				if ap > rand.Float64() {
					solution = newSolution
					oldDist = newDist
				}
			}
			T = T * alpha
		}
		return solution, oldDist
	}

	solution, fitness := annealing(solution)

	fmt.Println("Solution: ", solution)
	fmt.Println("Fitness: ", fitness)
}

func fitness(g genetic.MyGenome) float64 {
	verts := []int{}
	for i, b := range g.Gene {
		if b == 1 {
			verts = append(verts, i+1)
		}
	}

	var total float64
	var count float64
	for i := 1; i < len(verts); i++{
		for j := i+1; j < len(verts); j++ {
			if dist, found := distanceTable.Search(verts[i - 1], verts[j]); found {
				total += dist
				count++
			}
		}
	}

	return float64(size) - ( total / count)
}


func GeneticAlg(filename string) {
	distanceTable = DistanceTable{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erro ao abrir o ficheiro: ", err)
		os.Exit(1)
	}

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

		distanceTable.Add(node1, node2, cost)
	}
	file.Close()

	size = len(verts)

	param := genetic.Parameter{
		Initializer: new(genetic.RandomInitializer),
		Selector: new(genetic.RandomSelector),
		Breeder: new(genetic.GA2PointBreeder),
		Mutator: new(genetic.RandomBitFlipMutator),
		PMutate: 0.1,
		PBreed: 0.7,
	}

	ga := genetic.NewGA(param)

	genome := genetic.NewMyGenome(size, fitness)

	ga.Init(100, genome)
	ga.Optimize(10)

	ga.PrintTop(10)
}

func hybrid(filename string) {}

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
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())


	op := menu()
	switch op {
	case 1:
		graphAlg(os.Args[1])
	case 2:
		GeneticAlg(os.Args[1])
	case 3:
		hybrid(os.Args[1])
	case 4:
		test_all_graph()
	case 5:
		test_all_genetic()
	case 6:
		test_all_hybrid()
	}
}
func test_all_graph() {}

func test_all_genetic() {
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
		GeneticAlg(filename)
	}
}

func test_all_hybrid() {}

