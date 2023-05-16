package mu8

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
)

var (
	ErrNegativeFitness = errors.New("negative fitness")
	ErrInvalidFitness  = errors.New("got infinite or NaN fitness")
	ErrCodependency    = errors.New("codependency between individuals")
)

// FindCodependecy returns error if inconsistency detected in newIndividual function
// for use with mu8.Genome genetic algorithm implementations.
//
// Users can check if a codependency is found by checking the error:
//
//	if errors.Is(err, ErrCodependency) {
//		// Handle codependency case.
//	}
//
// The error will have a descriptive text of how the Genome is codependent.
func FindCodependecy[G Genome](src rand.Source, newIndividual func() G) error {
	ctx := context.Background()
	starter1 := newIndividual()
	starter2 := newIndividual()
	// These two fitnesses should be equal if there is not codependency.
	fit1 := starter1.Simulate(ctx)
	fit2 := starter2.Simulate(ctx)
	switch {
	case fit1 != fit2:
		return fmt.Errorf("%w: during subsequent calls to newIndividual which should return identical fitnesses. check for closure variable capture modification or preserved slice reference?", ErrCodependency)
	case fit1 == 0:
		return errors.New("cannot reliably determine codependency with zero fitness simulation results")
	case fit1 < 0 || fit2 < 0:
		return ErrNegativeFitness
	case math.IsNaN(fit1) || math.IsNaN(fit2) || math.IsInf(fit1, 0) || math.IsInf(fit2, 0):
		return ErrInvalidFitness
	}

	rng := rand.New(src)
	var codependents []int
	for i := 0; i < starter1.Len(); i++ {
		parent1 := newIndividual()
		parent2 := newIndividual()
		fit1 := parent1.Simulate(ctx)
		g := parent1.GetGene(i)
		// This line should have no effect on parent2's simulation (should be "initial" fitness).
		g.Mutate(rng)
		fit2 := parent2.Simulate(ctx)
		if math.IsNaN(fit1) || math.IsNaN(fit2) || math.IsInf(fit1, 0) || math.IsInf(fit2, 0) {
			if codependents != nil {
				return fmt.Errorf("invalid fitness and %w: detected in genes indices: %v", ErrCodependency, codependents)
			}
			return ErrInvalidFitness
		} else if fit1 != fit2 {
			codependents = append(codependents, i)
		}
	}
	if codependents == nil {
		return nil // No codependency detected.
	}
	return fmt.Errorf("%w: genes indices: %v", ErrCodependency, codependents)
}
