package genetic

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/soypat/mu8"
	"github.com/soypat/mu8/genes"
)

// Run this test with -race argument to detect
// race conditions
func TestIslandsRace(t *testing.T) {
	src := rand.NewSource(1)
	const (
		Ncrossovers      = 10
		genomelen        = 8
		Nindividuals     = 1000
		Nislands         = 5
		Nconcurrent      = Nislands // Must be <= number of islands.
		NgenPerCrossover = 10
		mutationRate     = 0.1
		polygamy         = 1
	)
	individuals := make([]*cfgenome, Nindividuals)
	for i := 0; i < Nindividuals; i++ {
		genome := newGenome(genomelen)
		mu8.Mutate(genome, src, .05)
		individuals[i] = genome
	}

	isls := NewIslands(Nislands, individuals, src, func() *cfgenome {
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
}

type cfgenome struct {
	genoma []genes.ConstrainedFloat
	// modify this during Simulate and read from it to catch race conditions
	racey int
}

func newGenome(n int) *cfgenome {
	return &cfgenome{genoma: make([]genes.ConstrainedFloat, n)}
}

func (g *cfgenome) GetGene(i int) mu8.Gene { return &g.genoma[i] }
func (g *cfgenome) Len() int               { return len(g.genoma) }

// Simulate simply adds the genes. We'd expect the genes to reach the max values of the constraint.
func (g *cfgenome) Simulate(context.Context) (fitness float64) {
	for i := range g.genoma {
		g.racey = i
		fitness += math.Abs(g.genoma[i].Value())
		g.racey++
		i = g.racey
	}
	return fitness / float64(g.Len()) / 3
}
