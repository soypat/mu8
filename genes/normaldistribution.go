package genes

import (
	"fmt"
	"math/rand"

	"github.com/soypat/mu8"
)

// NormalDistribution is a gene that mutates using a normal distribution.
// It implements the [mu8.Gene] interface.
type NormalDistribution struct {
	gene         float64
	stdDevMinus1 float64
}

// NewNormalDistribution returns a new NormalDistribution gene.
// The arguments are:
//   - val: the initial value of the gene.
//   - stdDeviation: the standard deviation of mutations.
func NewNormalDistribution(val, stdDeviation float64) *NormalDistribution {
	if stdDeviation <= 0 {
		panic("standard deviation must be non-zero and positive")
	}
	return &NormalDistribution{gene: val, stdDevMinus1: stdDeviation - 1}
}

// Mutate mutates the gene randomly using rng. It implements the [mu8.Gene] interface.
func (n *NormalDistribution) Mutate(rng *rand.Rand) {
	// Normal mutation distribution around current gene position.
	n.gene = rng.NormFloat64()*n.StdDev() + n.gene
}

// CloneFrom copies the argument gene into the receiver. It requires g to be of type
// *NormalDistribution. CloneFrom implements the [mu8.Gene] interface.
func (n *NormalDistribution) CloneFrom(g mu8.Gene) {
	co := castGene[*NormalDistribution](g)
	n.gene = co.gene
}

// Splice performs a crossover between the argument and the receiver genes
// and stores the result in the receiver. It implements the [mu8.Gene] interface.
func (n *NormalDistribution) Splice(rng *rand.Rand, g mu8.Gene) {
	co := castGene[*NormalDistribution](g)
	n.gene = rng.NormFloat64()*n.StdDev() + (co.gene+n.gene)/2
}

// Copy returns a copy of the gene.
func (n *NormalDistribution) Copy() *NormalDistribution {
	clone := *n
	return &clone
}

// StdDev returns the standard deviation of mutations on the gene.
func (n *NormalDistribution) StdDev() float64 {
	return n.stdDevMinus1 + 1
}

// SetValue sets the value of the gene.
func (n *NormalDistribution) SetValue(f float64) {
	n.gene = f
}

// Value returns the current value of the gene.
func (n *NormalDistribution) Value() float64 {
	return n.gene
}

// String returns a string representation of the NormalDistribution.
func (n *NormalDistribution) String() string {
	return fmt.Sprintf("%f", n.gene)
}
