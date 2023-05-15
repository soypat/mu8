package genetic

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"strconv"

	"github.com/soypat/mu8"
)

var (
	ErrNegativeFitness          = errors.New("fitness cannot yield negative values. Use zero instead. See pop.DubiousIndividual to recover problematic Genome information")
	ErrInvalidFitness           = errors.New("got infinite or NaN fitness. See pop.DubiousIndividual to recover problematic Genome information")
	ErrZeroFitnessSum           = errors.New("zero fitness sum: cannot make decisions")
	errChampionZeroFitness      = errors.New("zero fitness champion: consider initializing Population with a non-zero fitness individuals or you may never get results")
	ErrInfFitnessSum            = errors.New("infinite fitness sum: fitnesses returned by individuals are too large")
	ErrBadPolygamy              = errors.New("bad polygamy: must be in range [0, Nindividuals)")
	ErrBadMutationRate          = errors.New("bad mutation rate: must be in range (0, 1]")
	ErrCodependencyChampFitness = errors.New("codependency found: champion fitness should be monotonically increasing. check for preserved references in newIndividual function. See FindCodependency")

	errExitRequested = errors.New("call to Advance termination requested")
	// This should never trigger.
	errGenomeCast = errors.New("theoretically unreachable bad Genome->G cast. Make sure your Genomes return same type when calling GetGene(i) for a given ith index")
)

// FindCodependecy returns error if inconsistency detected in newIndividual function
// for use with mu8.Genome genetic algorithm implementations.
func FindCodependecy[G mu8.Genome](src rand.Source, newIndividual func() G) error {
	ctx := context.Background()

	starter1 := newIndividual()
	starter2 := newIndividual()
	fit1 := starter1.Simulate(ctx)
	fit2 := starter2.Simulate(ctx)
	if math.IsNaN(fit1) || math.IsNaN(fit2) {
		return errors.New("NaN fitness")
	} else if fit1 != fit2 {
		return errors.New("codependency between simulation results of subsequent calls to newIndividual. check for closure variable capture modification or preserved slice reference?")
	} else if fit1 == 0 {
		return errors.New("cannot reliably determine codependency with zero fitness simulation results")
	}

	rng := rand.New(src)
	for i := 0; i < starter1.Len(); i++ {
		parent1 := newIndividual()
		parent2 := newIndividual()
		fit1 := parent1.Simulate(ctx)
		g := parent1.GetGene(i)
		// This line should have no effect on parent2's simulation (should be "initial" fitness).
		g.Mutate(rng)
		fit2 := parent2.Simulate(ctx)
		if math.IsNaN(fit1) || math.IsNaN(fit2) {
			return errors.New("NaN fitness")
		} else if fit1 != fit2 {
			return errors.New(strconv.Itoa(i) + "th gene codependent")
		}
	}
	return nil
}
