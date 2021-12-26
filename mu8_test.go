package mu8_test

import (
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
func ExampleGenome() {
	src := rand.NewSource(1)
	const (
		Nprints      = 10
		genomelen    = 5
		Nindividuals = 100
		Ngenerations = 10000
		mutationRate = 0.1
		polygamy     = 1
	)
	individuals := make([]*genome, Nindividuals)
	for i := 0; i < Nindividuals; i++ {
		genome := newgenome(genomelen)
		mu8.Mutate(genome, src, .01)
		individuals[i] = genome
	}

	pop := genetic.NewPopulation(individuals, func() *genome { return newgenome(genomelen) }, src)
	prevFitness := 0.0
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
		if champFitness < prevFitness {
			panic("fitness not monotonically increasing")
		}
		prevFitness = champFitness
	}
	// Output:
	// champ fitness=0.154
	// champ fitness=0.859
	// champ fitness=0.859
	// champ fitness=0.898
	// champ fitness=0.898
	// champ fitness=0.911
	// champ fitness=0.911
	// champ fitness=0.925
	// champ fitness=0.928
	// champ fitness=0.928
}

type genome struct {
	genoma []*genes.ConstrainedFloat
}

func newgenome(n int) *genome {
	g := &genome{
		genoma: make([]*genes.ConstrainedFloat, n),
	}
	for i := 0; i < n; i++ {
		g.genoma[i] = genes.NewConstrainedFloat(0, 0, 1)
	}
	return g
}

func (g *genome) GetGene(i int) mu8.Gene { return g.genoma[i] }
func (g *genome) Len() int               { return len(g.genoma) }

// Simulate simply adds the genes. We'd expect the genes to reach the max values of the constraint.
func (g *genome) Simulate() (fitness float64) {
	for i := range g.genoma {
		fitness += g.genoma[i].Value()
	}
	return math.Max(0, fitness/float64(g.Len()))
}
