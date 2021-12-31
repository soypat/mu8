package genes

import (
	"fmt"
	"math/rand"

	"github.com/soypat/mu8"
)

type NormalDistribution struct {
	gene         float64
	stdDevMinus1 float64
}

func NewNormalDistribution(val, stdDeviation float64) *NormalDistribution {
	if stdDeviation <= 0 {
		panic("standard deviation must be non-zero and positive")
	}
	return &NormalDistribution{gene: val, stdDevMinus1: stdDeviation - 1}
}

func (n *NormalDistribution) Mutate(rng *rand.Rand) {
	// Normal mutation distribution around current gene position.
	n.gene = rng.NormFloat64()*n.StdDev() + n.gene
}

func (n *NormalDistribution) CloneFrom(g mu8.Gene) {
	co := castGene[*NormalDistribution](g)
	n.gene = co.gene
}

func (n *NormalDistribution) Splice(rng *rand.Rand, g mu8.Gene) {
	co := castGene[*NormalDistribution](g)
	n.gene = rng.NormFloat64()*n.StdDev() + (co.gene+n.gene)/2
}

func (n *NormalDistribution) Copy() *NormalDistribution {
	clone := *n
	return &clone
}

func (n *NormalDistribution) StdDev() float64 {
	return n.stdDevMinus1 + 1
}

func (n *NormalDistribution) SetValue(f float64) {
	n.gene = f
}

func (n *NormalDistribution) Value() float64 {
	return n.gene
}

func (n *NormalDistribution) String() string {
	return fmt.Sprintf("%f", n.gene)
}
