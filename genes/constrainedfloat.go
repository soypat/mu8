package genes

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/soypat/mu8"
)

// NewConstrainedFloat returns a mu8.Gene implementation for a number
// that should be kept within bounds [min,max] during mutation.
// start is the initial value of the gene.
func NewConstrainedFloat(start, min, max float64) *ConstrainedFloat {
	if min > max {
		panic(errBadConstraints)
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
	c.gene = c.clamp(f)
}

// Mutate changes the gene's value by a random amount within constraints.
// Mutate implements the [mu8.Gene] interface.
func (c *ConstrainedFloat) Mutate(rng *rand.Rand) {
	// Uniform mutation distribution.
	random := rng.Float64()
	random = c.min + random*c.rangeLength()
	c.gene = c.clamp(random)
}

// CloneFrom copies the argument gene into the receiver. CloneFrom implements
// the [mu8.Gene] interface. If g is not of type *ConstrainedFloat, CloneFrom panics.
func (c *ConstrainedFloat) CloneFrom(g mu8.Gene) {
	co := castGene[*ConstrainedFloat](g)
	c.gene = co.gene
}

// Copy returns a new ConstrainedFloat with the same value as the receiver.
func (c *ConstrainedFloat) Copy() *ConstrainedFloat {
	clone := *c
	return &clone
}

// Splice performs a crossover operation between the receiver and g.
// Splice implements the [mu8.Gene] interface. If g is not of type *ConstrainedFloat, Splice panics.
func (c *ConstrainedFloat) Splice(rng *rand.Rand, g mu8.Gene) {
	co := castGene[*ConstrainedFloat](g)
	random := rng.Float64()
	randi := 1 - random
	// Pick a random uniformly distributed gene between two values.
	c.gene = c.clamp((c.gene*randi + co.gene*random) / 2)
}

func (c *ConstrainedFloat) Format(state fmt.State, verb rune) {
	var val string
	prec, okp := state.Precision()
	width, okw := state.Width()
	switch verb {
	case 's', 'f':
		if !okw && !okp {
			val = fmt.Sprintf("%f", c.gene)
		} else {
			val = fmt.Sprintf("%[1]*.[2]*[3]f", width, prec, c.gene)
		}
	case 'v':
		val = fmt.Sprintf("{gene:%g, min:%g, max:%g}", c.gene, c.min, c.maxMinus1+1)
	case 'g':
		val = fmt.Sprintf("%[1]*.[2]*[3]g", width, prec, c.gene)
	default:
		val = fmt.Sprintf("!ERR(%%%v)", verb)
	}
	fmt.Fprint(state, val)
}

func (c *ConstrainedFloat) clamp(f float64) float64 {
	return math.Max(math.Min(f, c.maxMinus1+1), c.min)
}

func (c *ConstrainedFloat) rangeLength() float64 {
	return c.maxMinus1 + 1 - c.min
}

// String implements fmt.Stringer interface.
func (c *ConstrainedFloat) String() string {
	return fmt.Sprintf("%f", c.gene)
}
