package genes

import (
	"math"

	"github.com/soypat/mu8"
)

var _ mu8.Gene[*ConstrainedFloat] = (*ConstrainedFloat)(nil)

// NewConstrainedFloat returns a mu8.Gene implementation for a number
// that should be kept within bounds [min,max] during mutation.
func NewConstrainedFloat(start, min, max float64) *ConstrainedFloat {
	return &ConstrainedFloat{
		gene: start,
		min:  min,
		max:  max,
	}
}

// ConstrainedFloat implements Gene interface.
type ConstrainedFloat struct {
	// functional unit of heredity.
	gene     float64
	max, min float64
}

// Value returns actual value of constrained float.
func (c *ConstrainedFloat) Value() float64 { return c.gene }

func (c *ConstrainedFloat) Mutate(rand float64) {
	// Uniform mutation distribution.
	rand = c.min + rand*(c.max-c.min)
	c.gene = c.clamp(rand)
}

func (c *ConstrainedFloat) Copy() *ConstrainedFloat {
	clone := *c
	return &clone
}

func (c *ConstrainedFloat) Instance() *ConstrainedFloat { return c }

func (c *ConstrainedFloat) Splice(g *ConstrainedFloat) {
	// Naive average.
	c.gene = c.clamp((c.gene + g.gene) / 2)
}

func (c *ConstrainedFloat) clamp(f float64) float64 {
	return math.Max(math.Min(f, c.max), c.min)
}
