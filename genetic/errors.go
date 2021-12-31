package genetic

import (
	"errors"
	"math"
	"math/rand"
	"strconv"

	"github.com/soypat/mu8"
)

var (
	ErrNegativeFitness          = errors.New("fitness cannot yield negative values. Use zero instead.")
	ErrZeroFitnessSum           = errors.New("zero fitness sum: cannot make decisions")
	ErrBadPolygamy              = errors.New("bad polygamy: must be in range [0, Nindividuals)")
	ErrCodependencyChampFitness = errors.New("codependency found: champion fitness should be monotonically increasing. check for preserved references in newIndividual function. See FindCodependency")
	// This should never trigger.
	errGenomeCast = errors.New("theoretically unreachable. Bad Genome->G cast. make sure your Genomes return same type")
)

// FindCodependecy returns error if inconsistency detected in newIndividual function
// for use with mu8.Genome genetic algorithm implementations.
func FindCodependecy[G mu8.Genome](src rand.Source, newIndividual func() G) error {
	starter1 := newIndividual()
	starter2 := newIndividual()
	fit1 := starter1.Simulate()
	fit2 := starter2.Simulate()
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
		fit1 := parent1.Simulate()
		g := parent1.GetGene(i)
		// This line should have no effect on parent2's simulation (should be "initial" fitness).
		g.Mutate(rng)
		fit2 := parent2.Simulate()
		if math.IsNaN(fit1) || math.IsNaN(fit2) {
			return errors.New("NaN fitness")
		} else if fit1 != fit2 {
			return errors.New(strconv.Itoa(i) + "th gene codependent")
		}
	}
	return nil
}
