package genetic

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/soypat/mu8"
)

// Errors that may be encountered during the execution of the genetic algorithm.
var (
	ErrZeroFitnessSum      = errors.New("zero fitness sum: cannot make decisions")
	ErrInfFitnessSum       = errors.New("infinite fitness sum: fitnesses returned by individuals are too large")
	errCodependency        = fmt.Errorf("%w: check for preserved references in newIndividual function. See mu8.FindCodependency", mu8.ErrCodependency)
	errNegativeFitness     = fmt.Errorf("%w: use zero instead. See pop.DubiousIndividual to recover problematic Genome information", mu8.ErrNegativeFitness)
	errInvalidFitness      = fmt.Errorf("%w:  See pop.DubiousIndividual to recover problematic Genome information", mu8.ErrInvalidFitness)
	errChampionZeroFitness = errors.New("zero fitness champion: consider initializing Population with a non-zero fitness individuals or you may never get results")
	errBadPolygamy         = errors.New("bad polygamy: must be in range [0, Nindividuals)")
	errBadMutationRate     = errors.New("bad mutation rate: must be in range (0, 1]")
)

// Population provides a generic implementation
// of Genetic Algorithm.
type Population[G mu8.Genome] struct {
	individuals []G
	generator   func() G
	champ       G
	// dubiousIndividual shall be set with the problematic
	// Genome when encountering an unrecoverable error during algorithm execution
	dubiousIndividual G
	dubious           float64
	champFitness      float64
	fitness           []float64
	fitnessSum        float64
	gen               int
	rng               rand.Rand
}

// NewPopulation should be called when instantiating a new
// optimization instance of Genetic Algorithm.
//
// individuals are the initial population and should be independent of one another. These
// implement the mu8.Genome interface.
//
// newIndividual should instantiate a blank-slate Genome that is ready for cloning via
// mu8.Clone. Changes to newIndividual result should not be reflected in other individuals in
// the algorithm.
//
// src provides randomness necessary for the Genetic Algorithm to work. Use math.NewSource(1)
// for consistent results. It is not recommended for one to use a crypto/rand source directly.
// If a true random run is required then it is strongly suggested rand.NewSource(trueRandomSeed)
// is used and the seed saved in between runs to be able to replicate bugs.
//
// Example:
//
//	 pop := NewPopulation([]*ind{a, b, c}, rand.NewSource(1), func() *ind {
//			return newIndividual() // return a blank slate individual
//	 })
func NewPopulation[G mu8.Genome](individuals []G, src rand.Source, newIndividual func() G) Population[G] {
	if len(individuals) == 0 {
		panic("length of individuals must be of length greater than 0")
	}
	if individuals[0].Len() == 0 {
		panic("individuals must have at least one gene for algorithm to work")
	}
	return Population[G]{
		individuals: individuals,
		rng:         *rand.New(src),
		fitness:     make([]float64, len(individuals)),
		generator:   newIndividual,
		champ:       newIndividual(),
	}
}

// Individuals returns a reference the pool of individuals participating in
// the simulation. Calling Selection will update the value returned by
// Individuals if not cloned before calling Selection.
func (pop *Population[G]) Individuals() []G { return pop.individuals }

// Advance simulates current population and saves fitness scores. Multiple
// calls to Advance without calling Selection may have undesired effects.
func (pop *Population[G]) Advance(ctx context.Context) error {
	pop.fitnessSum = 0
	maxFitness := math.Inf(-1)
	champIdx := -1
	fitnessSum := 0.0
	for i := 0; i < len(pop.individuals) && ctx.Err() == nil; i++ {
		fitness := pop.individuals[i].Simulate(ctx)
		// We now check for errors that impede the continuation of the algorithm.
		if fitness < 0 {
			pop.dubious = fitness
			pop.dubiousIndividual = pop.individuals[i]
			return errNegativeFitness
		} else if math.IsInf(fitness, 0) || math.IsNaN(fitness) {
			pop.dubious = fitness
			pop.dubiousIndividual = pop.individuals[i]
			return errInvalidFitness
		}
		fitnessSum += fitness
		pop.fitness[i] = fitness
		if fitness > maxFitness {
			maxFitness = fitness
			champIdx = i
			if fitness > pop.champFitness {
				// we perform a greedy save of the new possible champion in case context is cancelled before Advance finishes.
				pop.champ = pop.individuals[i]
				pop.champFitness = fitness
			}
		}
	}
	switch {
	case champIdx < 0 || fitnessSum == 0:
		return ErrZeroFitnessSum // No decision can be taken and no progress can be made.
	case pop.fitness[champIdx] < pop.champFitness:
		// This is a big error. It means new instances of individuals are
		// affected by previous instances Simulation call or calls to gene's Mutate.
		// If this panic triggers consider all champion data has been compromised
		// and may not accurately represent "optimal" Genome.
		panic(errCodependency)
	case math.IsInf(pop.fitnessSum, 0):
		return ErrInfFitnessSum
	}
	newChamp := pop.generator()
	// Clone the champion so that his legacy may live on, untarnished by interbreeding and mutations.
	err := mu8.Clone(newChamp, pop.individuals[champIdx])
	if err != nil {
		return err // TODO: Should this panic?
	}
	// Looks like all went down OK. We can now update the population's champion.
	pop.champ = newChamp
	pop.fitnessSum = fitnessSum
	pop.champ = pop.individuals[champIdx]
	pop.champFitness = pop.fitness[champIdx]
	return nil
}

// Selection performs natural selection of individuals in the population.
// It first breeds individuals (fittest are most likely to be bred) and then
// mutates the babies obtained from the breeding procedure. The Individuals
// are updated once this function terminates.
func (pop *Population[G]) Selection(mutationRate float64, polygamy int) error {
	switch {
	case pop.champFitness == 0 && pop.fitnessSum == 0:
		// If you are getting this error try redesigning your fitness function
		// so that it yields a non-zero fitness value for at least one individual.
		return errChampionZeroFitness
	case pop.fitnessSum == 0:
		// If you are getting this error you either
		//  - Had a bad run, try running Advance again.
		//  - Population was not initialized properly via the NewPopulation function.
		return ErrZeroFitnessSum
	case mutationRate <= 0 || mutationRate > 1:
		return errBadMutationRate
	case polygamy < 0 || polygamy > len(pop.individuals):
		return errBadPolygamy
	}

	newGeneration := make([]G, len(pop.individuals))
	// Skip first index, reserved for our champion.
	for i := 1; i < len(pop.individuals); i++ {
		// Find the meanest, greenest individuals
		parents := pop.selectFittest(polygamy + 1)
		child, err := pop.breed(parents[0], parents...)
		if err != nil {
			return err
		}
		mu8.Mutate(child, &pop.rng, mutationRate)
		newGeneration[i] = child
	}
	// Looking out for our one and only, Champ.
	newGeneration[0] = pop.champ
	pop.individuals = newGeneration
	pop.gen++
	return nil
}

// Champion returns the best candidate of the population, this
// individual possessing the highest fitness score from last call to Advance().
func (pop *Population[G]) Champion() G {
	return pop.champ
}

// ChampionFitness returns the highest fitness score of the population found
// during the last call to Advance().
func (pop *Population[G]) ChampionFitness() float64 {
	return pop.champFitness
}

// DubiousIndividual returns the last individual that caused a NaN or Inf
// result during simulation to aid with debugging. It returns the
// Genome's zero value and 0 if no dubious individual was encountered.
// The presence of a dubious individual can mean one of two things:
//
// 1. The algorithm found an individual that maximizes its fitness
// beyond the precision of float64 type.
//
// 2. Presence of an edge case with the fitness function that breaks
// the real, positive number assumption.
//
// Both of mentioned cases means the fitness function is ill defined and should be corrected.
func (pop *Population[G]) DubiousIndividual() (dubious G, problematicFitness float64) {
	return pop.dubiousIndividual, pop.dubious
}

// selectFittest selects `sample` individuals from the population and returns a slice
// containing them. The most fittest are the most likely to be selected.
func (pop *Population[G]) selectFittest(sample int) (fittest []G) {
	// Quick return for clone case.
	if sample == 0 {
		return nil
	}
	// The lucky few selected will statistically be more likely to be fitter, proportional to their fitness.
	luckOfTheFit := slicemap(sample, func(int) float64 { return pop.fitnessSum * pop.rng.Float64() })
	runningSum := 0.0
	for i := 0; len(fittest) < sample; i++ {
		runningSum += pop.fitness[i]
		for _, threshold := range luckOfTheFit {
			// TODO(soypat): This code has to be overhauled so same parent is not selected
			// more than once.
			if runningSum > threshold {
				fittest = append(fittest, pop.individuals[i])
			}
		}
	}
	return fittest
}

// breed breeds receiver Genome with other genomes by splicing.
// An argument of no genomes returns a non-referential copy of the receiver,
// which could be described as a cloning procedure.
func (pop *Population[G]) breed(firstParent G, conjugates ...G) (G, error) {
	child := pop.generator()
	err := mu8.Clone(child, firstParent)
	if err != nil {
		return child, err
	}
	if len(conjugates) == 0 {
		return child, nil
	}

	for i := 0; i < child.Len(); i++ {
		gene := child.GetGene(i)
		for _, c := range conjugates {
			gene.Splice(&pop.rng, c.GetGene(i))
		}
	}
	return child, nil
}

func slicemap(n int, f func(int) float64) []float64 {
	result := make([]float64, n)
	for i := range result {
		result[i] = f(i)
	}
	return result
}
