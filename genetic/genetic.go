package genetic

import (
	"math"
	"math/rand"

	"github.com/soypat/mu8"
)

// Population provides a generic implementation
// of Genetic Algorithm.
type Population[G mu8.Genome] struct {
	individuals  []G
	generator    func() G
	champ        G
	champFitness float64
	fitness      []float64
	fitnessSum   float64
	gen          int
	rng          rand.Rand
}

func NewPopulation[G mu8.Genome](individuals []G, newIndividual func() G, src rand.Source) Population[G] {
	return Population[G]{
		individuals: individuals,
		rng:         *rand.New(src),
		fitness:     make([]float64, len(individuals)),
		generator:   newIndividual,
		champ:       newIndividual(),
	}
}

func (pop *Population[G]) Advance() error {
	pop.fitnessSum = 0
	maxFitness := math.Inf(-1)
	champIdx := -1
	for i := range pop.individuals {
		fitness := pop.individuals[i].Simulate()
		if fitness < 0 {
			return ErrNegativeFitness
		}
		pop.fitnessSum += fitness
		pop.fitness[i] = fitness
		if fitness > maxFitness {
			maxFitness = fitness
			champIdx = i
		}
	}
	pop.champ = pop.generator()
	// Clone the champion so that his legacy may live on, untarnished by interbreeding and mutations.
	mu8.Clone(pop.champ, pop.individuals[champIdx])
	bestFitness := pop.fitness[champIdx]
	if pop.fitnessSum == 0 {
		return ErrZeroFitnessSum
	} else if bestFitness < pop.champFitness {
		return ErrChampFitnessDecrease
	}
	pop.champFitness = bestFitness
	return nil
}

func (pop *Population[G]) Selection(mutationRate float64, polygamy int) error {
	if polygamy < 0 || polygamy > len(pop.individuals) {
		return ErrBadPolygamy
	}
	newGeneration := make([]G, len(pop.individuals))
	// Skip first index, reserved for our champion.
	for i := 1; i < len(pop.individuals); i++ {
		// Find the meanest, greenest individuals
		parents := pop.selectFittest(polygamy + 1)
		child := pop.breed(parents[0], parents...)
		mu8.Mutate(child, &pop.rng, mutationRate)
		newGeneration[i] = child
	}
	// Looking out for our one and only, Champ.
	newGeneration[0] = pop.champ
	pop.individuals = newGeneration
	pop.gen++
	return nil
}

func (pop *Population[G]) selectFittest(sample int) (fittest []G) {
	// Quick return for clone case.
	if sample == 0 {
		return nil // make([]mu8.Genome, 0) // empty slice
	}
	// The lucky few selected will statistically be more likely to be fitter, proportional to their fitness.
	luckOfTheFit := slicemap(sample, func(int) float64 { return pop.fitnessSum * pop.rng.Float64() })
	runningSum := 0.0
	for i := 0; len(fittest) < sample; i++ {
		runningSum += pop.fitness[i]
		for _, threshold := range luckOfTheFit {
			// TODO: This code has to be overhauled so same parent is not selected
			// more than once.
			if runningSum > threshold {
				fittest = append(fittest, pop.individuals[i])
			}
		}
	}
	return fittest
}

func (pop *Population[G]) Champion() G {
	return pop.champ
}
func (pop *Population[G]) ChampionFitness() float64 {
	return pop.champFitness
}

func slicemap(n int, f func(int) float64) []float64 {
	result := make([]float64, n)
	for i := range result {
		result[i] = f(i)
	}
	return result
}

// breed breeds receiver Genome with other genomes by splicing.
// An argument of no genomes returns a non-referential copy of the receiver,
// which could be described as a cloning procedure.
func (pop *Population[G]) breed(firstParent G, conjugates ...G) G {
	child := pop.generator()
	mu8.Clone(child, firstParent)
	if len(conjugates) == 0 {
		return child
	}
	for i := 0; i < child.Len(); i++ {
		gene := child.GetGene(i)
		for _, c := range conjugates {
			gene.Splice(c.GetGene(i))
		}
	}
	return child
}
