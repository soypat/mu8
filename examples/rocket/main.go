package main

import (
	"fmt"
	"math/rand"

	"github.com/soypat/mu8"
	"github.com/soypat/mu8/genes"
	"github.com/soypat/mu8/genetic"
)

const (
	// [m/s]
	gravity = 9.8
)

// Having constant length arrays lets us more dynamically program other aspects
// of this Genome. The process of cloning also becomes dead simple.
const Nstages = 2

var baseRocket = &rocket{
	payloadMass: 10,
	CD:          0.4,
	stages: [Nstages]stage{
		{
			isp:       300,
			massStruc: 200,
			// The way this works is the first value is the starting value,
			// other two are the permissible range the value may take during optimization.
			massProp:  genes.NewConstrainedFloat(0, 0, 3000),
			deltaMass: genes.NewConstrainedFloat(0, 0, 100),
			coastTime: genes.NewConstrainedFloat(0, 0, 300),
		},
		{
			isp:       300,
			massStruc: 20,
			massProp:  genes.NewConstrainedFloat(0, 0, 100),
			deltaMass: genes.NewConstrainedFloat(0, 0, 3),
			coastTime: genes.NewConstrainedFloat(0, 0, 300),
		},
	},
}

func main() {
	const (
		// How many iterations to print, spaced evenly between generations
		// Do not set to less than 2.
		Nprints = 10
		// Number of generations to simulate.
		Ngen = 4000
		// Number of individuals in populations
		Nindividuals = 8
		// Polygamy (how many partners a rocket has)
		polygamy = 2
		// Mutation rate
		mutrate = 0.1
	)
	type genoma = *rocket
	src := rand.NewSource(2)
	individuals := make([]*rocket, Nindividuals)
	for i := range individuals {
		clone := baseRocket.Clone()
		mu8.Mutate(clone, src, 0.3)
		individuals[i] = clone
	}

	pop := genetic.NewPopulation(individuals, src, func() *rocket { return baseRocket.Clone() })
	for i := 0; i < Ngen; i++ {
		err := pop.Advance()
		if err != nil {
			panic(err)
		}
		err = pop.Selection(mutrate, polygamy)
		if err != nil {
			panic(err)
		}
		if i%(Ngen/(Nprints-1)) == 0 || i == Ngen-1 {
			bestFitness := pop.ChampionFitness()
			fmt.Printf("champHeight:%.3fkm\n", bestFitness)
		}
	}
	best := pop.Champion()
	fmt.Println("our champion:", best)
}
