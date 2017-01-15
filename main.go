package main

import (
	"math/rand"
	"os"
	"time"
	"fmt"
	"bufio"
	"github.com/MaxHalford/gago"
	"strconv"
	"runtime"
	"math"
)
var (
	distanceTable DistanceTable
 	size int
)

type Solution []int
func (s Solution) Distance() float64 {
	verts := []int{}
	for i, b := range s {
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

	return count - ( total / count)
}

func (s Solution) Neighbour() Solution {
	n := s
	i := rand.Intn(len(s))
	if n[i] == 0 {
		n[i] = 1
	} else {
		n[i] = 0
	}
	return n
}

type Edge string
func makeEdge(n1, n2 int) Edge {
	return Edge(strconv.Itoa(n1) + "-" + strconv.Itoa(n2))
}

type DistanceTable map[Edge]float64

func (d *DistanceTable) Add(n1, n2 int, dist float64) {
	if _, found := (*d).Search(n1, n2); !found {
		e := makeEdge(n1, n2)
		(*d)[e] = dist
	}
}

// Search devolve o custo de ir de n1 para n2
// Devolve um false caso essa ligacao nao exista
func (d DistanceTable) Search(n1, n2 int) (float64, bool) {
	e := makeEdge(n1, n2)
	if dist, ok := d[e]; !ok {
		return 0.0, false
	} else {
		return dist, true
	}
}

func (d DistanceTable) Print() {
	fmt.Println("Distance Table")
	for edge, dist := range d {
		fmt.Println(edge, " = ", dist)
	}
}

type Digits []int
func (X Digits) Evaluate() float64 {
	verts := []int{}
	for i, b := range X {
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

	return (total / count)
}

// Mutate a slice of digits by permuting it's values.
func (X Digits) Mutate(rng *rand.Rand) {
	gago.MutPermuteInt(X, 3, rng)
}

// Crossover a slice of digits with another by applying 2-point crossover.
func (X Digits) Crossover(Y gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossGNXInt(X, Y.(Digits), 2, rng)
	return Digits(o1), Digits(o2)
}

// MakeDigits creates random slices of digits by randomly assigning them 1s and
// 0s.
func MakeDigits(rng *rand.Rand) gago.Genome {
	var digits = make(Digits, size)
	for i := range digits {
		if rng.Float64() < 0.5 {
			digits[i] = 1
		}
	}
	return gago.Genome(digits)
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

	distanceTable = make(map[Edge]float64)
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
		TMin := 0.00001
		alpha := 0.001
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


func genetic(filename string) {
	distanceTable = make(map[Edge]float64)
	// TESTING GAGO
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

	//distanceTable.Print()
	size = len(verts)

	var ga = gago.Generational(MakeDigits)
	ga.Initialize()

	for i := 1; i < 1000; i++ {
		ga.Enhance()
	}
	fmt.Printf("Best fitness -> %f\n", ga.Best.Fitness)
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
		genetic(os.Args[1])
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
		genetic(filename)
	}
}

func test_all_hybrid() {}

