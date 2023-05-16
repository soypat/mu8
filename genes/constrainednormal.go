package genes

import (
	"math"
	"math/rand"

	"github.com/soypat/mu8"
)

// ConstrainedFloat is a float64 gene that is constrained to a range and whose
// mutations are drawn from a uniform distribution.
type ConstrainedNormalDistr struct {
	NormalDistribution
	minPlus3sd, maxMinus3sd float64
}

// NewConstrainedNormalDistr returns a new ConstrainedNormalDistr.
// The arguments are:
//   - val: the initial value of the gene.
//   - stdDeviation: the standard deviation of mutations.
//   - min: the minimum value the gene can take.
//   - max: the maximum value the gene can take.
func NewConstrainedNormalDistr(val, stdDeviation, min, max float64) *ConstrainedNormalDistr {
	if min > max {
		panic(errBadConstraints)
	}
	return &ConstrainedNormalDistr{
		NormalDistribution: *NewNormalDistribution(val, stdDeviation),
		minPlus3sd:         min + 3*stdDeviation,
		maxMinus3sd:        max - 3*stdDeviation,
	}
}

// Mutate mutates the gene according to a random normal distribution.
// It implements the [mu8.Gene] interface.
func (cn *ConstrainedNormalDistr) Mutate(rng *rand.Rand) {
	cn.NormalDistribution.Mutate(rng)
	cn.clamp()
}

// CloneFrom copies the argument gene into the receiver. CloneFrom implements the
// [mu8.Gene] interface. If g is not of type *ConstrainedNormalDistr, CloneFrom panics.
func (cn *ConstrainedNormalDistr) CloneFrom(g mu8.Gene) {
	co := castGene[*ConstrainedNormalDistr](g)
	cn.gene = co.gene
}

// Splice performs a crossover between the argument and the receiver genes
// and stores the result in the receiver. It implements the [mu8.Gene] interface.
// If g is not of type *ConstrainedNormalDistr, Splice panics.
func (cn *ConstrainedNormalDistr) Splice(rng *rand.Rand, g mu8.Gene) {
	co := castGene[*ConstrainedNormalDistr](g)
	cn.NormalDistribution.Splice(rng, &co.NormalDistribution)
	cn.clamp()
}

// Copy returns a copy of the gene.
func (cn *ConstrainedNormalDistr) Copy() *ConstrainedNormalDistr {
	clone := *cn
	return &clone
}

// clamp clamps the gene value to the constraints.
func (cn *ConstrainedNormalDistr) clamp() {
	sd3 := cn.StdDev() * 3
	min := cn.minPlus3sd - sd3
	max := cn.maxMinus3sd + sd3
	cn.gene = math.Max(min, math.Min(max, cn.gene))
}

// SetValue sets the value of the gene.
func (cn *ConstrainedNormalDistr) SetValue(f float64) {
	cn.gene = f
	cn.clamp()
}
