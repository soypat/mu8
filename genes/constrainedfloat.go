package genes

import (
	"math"

	"github.com/soypat/mu8"
)

// ConstrainedFloat implements Gene interface.
type ConstrainedFloat struct {
	// functional unit of heredity.
	gene     float64
	max, min float64
}

func (c *ConstrainedFloat) Copy() mu8.Gene {
	clone := *c
	return &clone
}

func (c *ConstrainedFloat) Splice(g mu8.Gene) {
	cother, ok := g.(*ConstrainedFloat)
	if !ok {
		panic("expected same type in Splice")
	}
	// Naive average.
	c.gene = c.clamp((c.gene + cother.gene) / 2)
}

func (c *ConstrainedFloat) clamp(f float64) float64 {
	return math.Max(math.Min(f, c.max), c.min)
}
