package genes

import (
	"math"

	"github.com/soypat/mu8"
)

var _ mu8.Gene[*ConstrainedFloat] = (*ConstrainedFloat)(nil)

// ConstrainedFloat implements Gene interface.
type ConstrainedFloat struct {
	// functional unit of heredity.
	gene     float64
	max, min float64
}

func (c *ConstrainedFloat) Mutate(rand float64) {
	// Uniform mutation distribution.
	rand = c.min + rand*(c.max-c.min)
	c.gene = rand
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
