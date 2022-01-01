package mu8_test

import (
	"fmt"
	"math/rand"

	"github.com/soypat/mu8"
	"github.com/soypat/mu8/genes"
	"github.com/soypat/mu8/genetic"
)

// This simple program evolves
// a genome to maximize it's ConstrainedFloat
// genome so that it reaches the max value.
func ExamplePopulation() {
	src := rand.NewSource(1)
	const (
		Nprints      = 10
		genomelen    = 5
		Nindividuals = 100
		Ngenerations = 10000
		mutationRate = 0.5
		polygamy     = 1
	)
	individuals := make([]*mygenome, Nindividuals)
	for i := 0; i < Nindividuals; i++ {
		genome := newGenome(genomelen)
		mu8.Mutate(genome, src, .01)
		individuals[i] = genome
	}

	pop := genetic.NewPopulation(individuals, src, func() *mygenome {
		return newGenome(genomelen)
	})
	for i := 0; i < Ngenerations; i++ {
		err := pop.Advance()
		if err != nil {
			panic(err.Error())
		}
		err = pop.Selection(mutationRate, polygamy)
		if err != nil {
			panic(err.Error())
		}
		champFitness := pop.ChampionFitness()
		if i%(Ngenerations/Nprints) == 0 {
			fmt.Printf("champ fitness=%.3f\n", champFitness)
		}
	}
	// Output:
	// champ fitness=0.180
	// champ fitness=0.942
	// champ fitness=0.942
	// champ fitness=0.949
	// champ fitness=0.949
	// champ fitness=0.949
	// champ fitness=0.949
	// champ fitness=0.949
	// champ fitness=0.949
	// champ fitness=0.953
}

type mygenome struct {
	genoma []genes.ConstrainedFloat
}

func newGenome(n int) *mygenome {
	return &mygenome{genoma: make([]genes.ConstrainedFloat, n)}
}

func (g *mygenome) GetGene(i int) mu8.Gene { return &g.genoma[i] }
func (g *mygenome) Len() int               { return len(g.genoma) }

// Simulate simply adds the genes. We'd expect the genes to reach the max values of the constraint.
func (g *mygenome) Simulate() (fitness float64) {
	for i := range g.genoma {
		fitness += g.genoma[i].Value()
	}
	return fitness / float64(g.Len())
}

func ExampleIslands() {
	src := rand.NewSource(1)
	const (
		Ncrossovers      = 10
		genomelen        = 5
		Nislands         = 5
		Nconcurrent      = Nislands // Must be <= number of islands.
		Nindividuals     = 50
		NgenPerCrossover = 2000
		mutationRate     = 0.2
		polygamy         = 2
	)
	individuals := make([]*mygenome, Nindividuals)
	for i := 0; i < Nindividuals; i++ {
		genome := newGenome(genomelen)
		mu8.Mutate(genome, src, .80)
		individuals[i] = genome
	}

	isls := genetic.NewIslands(Nislands, individuals, src, func() *mygenome {
		return newGenome(genomelen)
	})
	for i := 0; i < Ncrossovers; i++ {
		err := isls.Advance(mutationRate, polygamy, NgenPerCrossover, Nconcurrent)
		if err != nil {
			panic(err.Error())
		}
		isls.Crossover()
		champFitness := isls.ChampionFitness()
		fmt.Printf("champ fitness=%.3f\n", champFitness)
	}
	// Output:
	// champ fitness=0.849
	// champ fitness=0.849
	// champ fitness=0.849
	// champ fitness=0.941
	// champ fitness=0.941
	// champ fitness=0.941
	// champ fitness=0.941
	// champ fitness=0.941
	// champ fitness=0.941
	// champ fitness=0.941
}
