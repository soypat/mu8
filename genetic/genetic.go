package genetic

import (
	"math/rand"

	"github.com/soypat/mu8"
)

// Population provides a generic implementation
// of Genetic Algorithm.
type Population struct {
	individuals []mu8.Genome
	champ       mu8.Genome
	fitness     []float64
	fitnessSum  float64
	gen         int
	rng         rand.Rand
}

func NewPopulation(individuals []mu8.Genome, src rand.Source) Population {
	return Population{
		individuals: individuals,
		rng:         *rand.New(src),
		fitness:     make([]float64, len(individuals)),
	}
}

func (pop *Population) Advance() {
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
	pop.champ = pop.individuals[champIdx].Breed()
}

func (pop *Population) Selection(mutationRate float64, polygamy int) {
	newGeneration := make([]mu8.Genome, len(pop.individuals))
	// Skip first index, reserved for our champion.
	for i := 1; i < len(pop.individuals); i++ {
		p := pop.individuals[i]
		// Find the meanest, greenest individuals
		parents := pop.selectFittest(polygamy)
		child := p.Breed(parents...)
		child.Mutate(mutationRate)
		newGeneration[i] = child
	}
	// Looking out for our one and only, Champ.
	newGeneration[0] = pop.champ
	pop.individuals = newGeneration
	pop.gen++
}

func (pop *Population) selectFittest(sample int) (fittest []mu8.Genome) {
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

func (pop *Population) Champion() mu8.Genome {
	return pop.champ
}

func slicemap(n int, f func(int) float64) []float64 {
	result := make([]float64, n)
	for i := range result {
		result[i] = f(i)
	}
	return result
}
