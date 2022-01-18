package genes

import (
	"fmt"
	"math/rand"

	"github.com/soypat/mu8"
)

// NewConstrainedInt returns a mu8.Gene implementation for a number
// that should be kept within bounds [min,max] during mutation.
func NewConstrainedInt(start, min, max int) *ConstrainedInt {
	if min >= max {
		panic(errBadConstraints)
	}
	if start > max || start < min {
		panic(errStartOutOfBounds)
	}
	return &ConstrainedInt{
		gene:        start,
		min:         min,
		rangeMinus1: max - min - 1,
	}
}

// ConstrainedInt implements Gene interface.
// Is automatically initialized to the domain [0,1]
type ConstrainedInt struct {
	// functional unit of heredity.
	gene int
	min  int
	// rangeMinus1 is the max value constraint can reach, minus 1+min.
	// This allows initialization without a call to NewConstrainedInt.
	rangeMinus1 int
}

// Value returns actual value of constrained float.
func (c *ConstrainedInt) Value() int { return c.gene }

// SetValue sets the gene's actual value. This method may be useful
// for setting best gene value for a single individual in the
// population by hand between runs.
func (c *ConstrainedInt) SetValue(f int) {
	if f < c.min || f > c.rangeMinus1+1 {
		panic("value not within constraints")
	}
	c.gene = f
}

func (c *ConstrainedInt) Mutate(rng *rand.Rand) {
	// Uniform mutation distribution.
	c.gene = c.min + rng.Intn(c.rangeMinus1+1)
}

func (c *ConstrainedInt) CloneFrom(g mu8.Gene) {
	co := castGene[*ConstrainedInt](g)
	c.gene = co.gene
}

func (c *ConstrainedInt) Copy() *ConstrainedInt {
	clone := *c
	return &clone
}

func (c *ConstrainedInt) Splice(rng *rand.Rand, g mu8.Gene) {
	co := castGene[*ConstrainedInt](g)
	diff := c.gene - co.gene
	if diff <= 0 {
		if diff == 0 {
			return // no work to do if genes are equal, also avoids a panic.
		}
		diff = -diff
	}
	random := rng.Intn(diff)
	minGene := min(c.gene, co.gene)
	// Pick a random uniformly distributed gene between two values.
	c.gene = minGene + random
}

func (c *ConstrainedInt) String() string {
	return fmt.Sprintf("%d", c.gene)
}
