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
	genoma []genes.ConstrainedNormalDistrGrad
}

func newGenome(n int) *mygenome {
	return &mygenome{genoma: make([]genes.ConstrainedNormalDistrGrad, n)}
}

func (g *mygenome) GetGene(i int) mu8.Gene         { return &g.genoma[i].ConstrainedNormalDistr }
func (g *mygenome) GetGeneGrad(i int) mu8.GeneGrad { return &g.genoma[i] }
func (g *mygenome) Len() int                       { return len(g.genoma) }
func (g *mygenome) LenGrad() int                   { return g.Len() }

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

func ExampleGradient() {
	src := rand.NewSource(1)
	const (
		genomelen      = 6
		gradMultiplier = 10.0
		epochs         = 6
	)
	// Create new individual and mutate it randomly.
	individual := newGenome(genomelen)
	rng := rand.New(src)
	for i := 0; i < genomelen; i++ {
		individual.GetGene(i).Mutate(rng)
	}
	// Prepare for gradient descent.
	grads := make([]float64, genomelen)
	ctx := context.Background()
	// Champion will harbor our best individual.
	champion := newGenome(genomelen)
	for epoch := 0; epoch < epochs; epoch++ {
		// We calculate the gradients of the individual passing a nil
		// newIndividual callback since the GenomeGrad type we implemented
		// does not require blank-slate initialization.
		err := mu8.Gradient(ctx, grads, individual, nil)
		if err != nil {
			panic(err)
		}
		// Apply gradients.
		for i := 0; i < individual.Len(); i++ {
			gene := individual.GetGeneGrad(i)
			grad := grads[i]
			gene.SetValue(gene.Value() + grad*gradMultiplier)
		}
		mu8.CloneGrad(champion, individual)
		fmt.Printf("fitness=%f with grads=%f\n", individual.Simulate(ctx), grads)
	}

	// Output:
	// fitness=0.467390 with grads=[-0.055556 -0.055556 -0.055556 0.055556 0.055556 0.055556]
	// fitness=0.630529 with grads=[-0.055556 -0.055556 -0.055556 0.055556 0.055556 0.055556]
	// fitness=0.784850 with grads=[-0.055556 -0.055556 -0.055556 0.000000 0.055556 0.055556]
	// fitness=0.913839 with grads=[-0.055556 -0.055556 -0.055556 0.000000 0.055556 0.055556]
	// fitness=0.994674 with grads=[-0.055556 -0.055556 -0.055556 0.000000 0.055556 0.055556]
	// fitness=1.000000 with grads=[-0.055556 -0.055556 -0.055556 0.000000 0.000000 0.000000]
}
