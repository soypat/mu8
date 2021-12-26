package mu8_test

import (
	"math/rand"
	"testing"

	"github.com/soypat/mu8/genes"
	"github.com/soypat/mu8/genetic"
)

func TestGenome(t *testing.T) {
	src := rand.NewSource(1)
	const (
		genomelen    = 2
		Nindividuals = 4
		Ngenerations = 10
		mutationRate = 0.01
		polygamy     = 0
	)
	individuals := make([]genomeimpl, Nindividuals)
	for i := 0; i < Nindividuals; i++ {
		individuals[i] = newgenome(genomelen)
	}
	pop := genetic.NewPopulation[geneimpl](individuals, src)
	for i := 0; i < Ngenerations; i++ {
		pop.Advance()
		pop.Selection(mutationRate, polygamy)
	}
	pop.Champion()
}

// remove type aliases once API is well defined
type geneimpl = *genes.ConstrainedFloat
type genomeimpl = *genome

type genome struct {
	genoma []geneimpl
}

func newgenome(n int) genomeimpl {
	g := &genome{
		genoma: make([]*genes.ConstrainedFloat, n),
	}
	for i := 0; i < n; i++ {
		g.genoma[i] = genes.NewConstrainedFloat(0, -1, 1)
	}
	return g
}

func (g *genome) GetGene(i int) geneimpl { return g.genoma[i] }
func (g *genome) Len() int               { return len(g.genoma) }

// Simulate simply adds the genes. We'd expect the genes to reach the max values of the constraint.
func (g *genome) Simulate() (fitness float64) {
	for i := range g.genoma {
		fitness += g.genoma[i].Value()
	}
	return fitness / float64(g.Len())
}

func (g *genome) Clone() genomeimpl {
	clone := &genome{
		genoma: make([]*genes.ConstrainedFloat, g.Len()),
	}
	for i := range clone.genoma {
		clone.genoma[i] = g.genoma[i].Copy()
	}
	return clone
}
