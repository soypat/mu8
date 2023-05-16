package genes

import (
	"math"
	"math/rand"

	"github.com/soypat/mu8"
)

type ConstrainedNormalDistr struct {
	NormalDistribution
	minPlus3sd, maxMinus3sd float64
}

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

func (cn *ConstrainedNormalDistr) Mutate(rng *rand.Rand) {
	cn.NormalDistribution.Mutate(rng)
	cn.clamp()
}

func (cn *ConstrainedNormalDistr) CloneFrom(g mu8.Gene) {
	co := castGene[*ConstrainedNormalDistr](g)
	cn.gene = co.gene
}

func (cn *ConstrainedNormalDistr) Splice(rng *rand.Rand, g mu8.Gene) {
	co := castGene[*ConstrainedNormalDistr](g)
	cn.NormalDistribution.Splice(rng, &co.NormalDistribution)
	cn.clamp()
}

func (cn *ConstrainedNormalDistr) Copy() *ConstrainedNormalDistr {
	clone := *cn
	return &clone
}

func (cn *ConstrainedNormalDistr) clamp() {
	sd3 := cn.StdDev() * 3
	min := cn.minPlus3sd - sd3
	max := cn.maxMinus3sd + sd3
	cn.gene = math.Max(min, math.Min(max, cn.gene))
}

func (cn *ConstrainedNormalDistr) SetValue(f float64) {
	cn.gene = f
	cn.clamp()
}
