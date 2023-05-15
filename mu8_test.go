package mu8_test

import (
	"context"
	"fmt"
	"math"
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
		genomelen    = 8
		Nindividuals = 100
		Ngenerations = 1000
		mutationRate = 0.05
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
		err := pop.Advance(context.Background())
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
	// champ fitness=0.081
	// champ fitness=0.832
	// champ fitness=0.860
	// champ fitness=0.887
	// champ fitness=0.887
	// champ fitness=0.926
	// champ fitness=0.926
	// champ fitness=0.926
	// champ fitness=0.926
	// champ fitness=0.953
}

type mygenome struct {
	genoma []genes.ConstrainedNormalDistr
}

func newGenome(n int) *mygenome {
	return &mygenome{genoma: make([]genes.ConstrainedNormalDistr, n)}
}

func (g *mygenome) GetGene(i int) mu8.Gene { return &g.genoma[i] }
func (g *mygenome) Len() int               { return len(g.genoma) }

// Simulate simply adds the genes. We'd expect the genes to reach the max values of the constraint.
func (g *mygenome) Simulate(context.Context) (fitness float64) {
	for i := range g.genoma {
		fitness += math.Abs(g.genoma[i].Value())
	}
	return fitness / float64(g.Len()) / 3
}

func ExampleIslands() {
	src := rand.NewSource(1)
	const (
		Ncrossovers      = 10
		genomelen        = 8
		Nindividuals     = 100
		Nislands         = 5
		Nconcurrent      = Nislands // Must be <= number of islands.
		NgenPerCrossover = 10
		mutationRate     = 0.1
		polygamy         = 1
	)
	individuals := make([]*mygenome, Nindividuals)
	for i := 0; i < Nindividuals; i++ {
		genome := newGenome(genomelen)
		mu8.Mutate(genome, src, .05)
		individuals[i] = genome
	}

	isls := genetic.NewIslands(Nislands, individuals, src, func() *mygenome {
		return newGenome(genomelen)
	})
	for i := 0; i < Ncrossovers; i++ {
		err := isls.Advance(context.Background(), mutationRate, polygamy, NgenPerCrossover, Nconcurrent)
		if err != nil {
			panic(err.Error())
		}
		isls.Crossover()
		champFitness := isls.ChampionFitness()
		fmt.Printf("champ fitness=%.3f\n", champFitness)
	}
	// Output:
	// champ fitness=0.882
	// champ fitness=0.897
	// champ fitness=0.923
	// champ fitness=0.946
	// champ fitness=0.946
	// champ fitness=0.946
	// champ fitness=0.946
	// champ fitness=0.946
	// champ fitness=0.956
	// champ fitness=0.956
}
