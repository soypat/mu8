package genes

import (
	"math"

	"github.com/soypat/mu8"
)

var _ mu8.Gene = (*ConstrainedFloat)(nil)

// NewConstrainedFloat returns a mu8.Gene implementation for a number
// that should be kept within bounds [min,max] during mutation.
func NewConstrainedFloat(start, min, max float64) *ConstrainedFloat {
	if min > max {
		panic("got min less than max in NewConstrainedFloat")
	}
	return &ConstrainedFloat{
		gene:      start,
		min:       min,
		maxMinus1: max - 1,
	}
}

// ConstrainedFloat implements Gene interface.
// Is automatically initialized to the domain [0,1]
type ConstrainedFloat struct {
	// functional unit of heredity.
	gene float64
	min  float64
	// maxMinus1 is the max value constraint can reach, minus 1.
	// This allows initialization without a call to NewConstrainedFloat.
	maxMinus1 float64
}

// Value returns actual value of constrained float.
func (c *ConstrainedFloat) Value() float64 { return c.gene }

// SetValue sets the gene's actual value. This method may be useful
// for setting best gene value for a single individual in the
// population by hand between runs.
func (c *ConstrainedFloat) SetValue(f float64) {
	if f < c.min || f > c.maxMinus1+1 {
		panic("value not within constraints")
	}
	c.gene = f
}

func (c *ConstrainedFloat) Mutate(random float64) {
	// Uniform mutation distribution.
	random = c.min + random*c.rangeLength()
	c.gene = c.clamp(random)
}

func (c *ConstrainedFloat) CloneFrom(g mu8.Gene) {
	co := castGene[*ConstrainedFloat](g)
	c.gene = co.gene
}

func (c *ConstrainedFloat) Copy() *ConstrainedFloat {
	clone := *c
	return &clone
}

func (c *ConstrainedFloat) Splice(random float64, g mu8.Gene) {
	co := castGene[*ConstrainedFloat](g)
	randi := 1 - random
	// Pick a random uniformly distributed gene between two values.
	c.gene = c.clamp((c.gene*randi + co.gene*random) / 2)
}

func (c *ConstrainedFloat) clamp(f float64) float64 {
	return math.Max(math.Min(f, c.maxMinus1+1), c.min)
}

func (c *ConstrainedFloat) rangeLength() float64 {
	return c.maxMinus1 + 1 - c.min
}
