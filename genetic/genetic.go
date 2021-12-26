package genetic

import (
	"math/rand"

	"github.com/soypat/mu8"
)

// Population provides a generic implementation
// of Genetic Algorithm.
type Population[T any] struct {
	individuals []mu8.Genome[T]
	champ       mu8.Genome[T]
	fitness     []float64
	fitnessSum  float64
	gen         int
	rng         rand.Rand
}

func NewPopulation[T any](individuals []mu8.Genome[T], src rand.Source) Population[T] {
	return Population[T]{
		individuals: individuals,
		rng:         *rand.New(src),
		fitness:     make([]float64, len(individuals)),
	}
}

func (pop *Population[T]) Advance() {
	var maxFitness float64
	champIdx := -1
	for i := range pop.individuals {
		fitness := pop.individuals[i].Simulate()
		pop.fitness[i] = fitness
		if fitness > maxFitness {
			maxFitness = fitness
			champIdx = i
		}
		pop.fitnessSum += fitness
	}
	// Clone the champion so that his legacy may live on, untarnished by interbreeding and mutations.
	pop.champ = pop.individuals[champIdx].Clone()
}

func (pop *Population[T]) Selection(mutationRate float64, polygamy int) {
	newGeneration := make([]mu8.Genome[T], len(pop.individuals))
	// Skip first index, reserved for our champion.
	for i := 1; i < len(pop.individuals); i++ {
		p := pop.individuals[i]
		// Find the meanest, greenest individuals
		parents := pop.selectFittest(polygamy)
		child := Breed(p, parents...)
		mu8.Mutate(child, &pop.rng, mutationRate)
		newGeneration[i] = child
	}
	// Looking out for our one and only, Champ.
	newGeneration[0] = pop.champ
	pop.individuals = newGeneration
	pop.gen++
}

func (pop *Population[T]) selectFittest(sample int) (fittest []mu8.Genome[T]) {
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
			if runningSum > threshold {
				fittest = append(fittest, pop.individuals[i])
			}
		}
	}
	return fittest
}

func (pop *Population[T]) Champion() mu8.Genome[T] {
	return pop.champ
}

func slicemap(n int, f func(int) float64) []float64 {
	result := make([]float64, n)
	for i := range result {
		result[i] = f(i)
	}
	return result
}

// Breed breeds receiver Genome with other genomes by splicing.
// An argument of no genomes returns a non-referential copy of the receiver,
// which could be described as a cloning procedure.
func Breed[T any](firstParent mu8.Genome[T], conjugates ...mu8.Genome[T]) mu8.Genome[T] {
	child := firstParent.Clone()
	if len(conjugates) == 0 {
		return child
	}
	for i := 0; i < child.Len(); i++ {
		gene := child.GetGene(i)
		for _, c := range conjugates {
			gene.Splice(c.GetGene(i).Instance())
		}
	}
	return child
}
