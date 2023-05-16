package genes

import "github.com/soypat/mu8"

// Package level defined default step size.
const defaultStep = 5e-7

// Step returns the step size that should be performed during gradient descent.
// It implements the [mu8.GeneGrad] interface by returning a package-level defined default step size.
func (*ConstrainedFloat) Step() float64 { return defaultStep }

// Step returns the step size that should be performed during gradient descent.
// It implements the [mu8.GeneGrad] interface by returning a package-level defined default step size.
func (*ConstrainedNormalDistr) Step() float64 { return defaultStep }

// Step returns the step size that should be performed during gradient descent.
// It implements the [mu8.GeneGrad] interface by returning a package-level defined default step size.
func (*NormalDistribution) Step() float64 { return defaultStep }

// compile time check that these types implement the GeneGrad interface.
var (
	_ mu8.GeneGrad = (*ConstrainedFloatGrad)(nil)
	_ mu8.GeneGrad = (*ConstrainedNormalDistrGrad)(nil)
	_ mu8.GeneGrad = (*ConstrainedFloat)(nil)
	_ mu8.GeneGrad = (*ConstrainedNormalDistr)(nil)
	_ mu8.GeneGrad = (*NormalDistribution)(nil)
)

// ConstrainedFloatGrad is a ConstrainedFloat that implements the GeneGrad interface
// with a programmable step size.
// It implements the [mu8.GeneGrad] interface.
type ConstrainedFloatGrad struct {
	ConstrainedFloat
	stepMinusDefaultStep float64
}

// Step returns the step size that should be performed during gradient descent.
// It implements the [mu8.GeneGrad] interface.
func (cf *ConstrainedFloatGrad) Step() float64 {
	return cf.stepMinusDefaultStep + defaultStep
}

// NewConstrainedFloatGrad returns a new ConstrainedFloatGrad.
// The arguments are:
//   - start: the initial value of the gene.
//   - min: the minimum value the gene can take.
//   - max: the maximum value the gene can take.
//   - step: the step size that should be used gradient descent.
func NewConstrainedFloatGrad(start, min, max, step float64) *ConstrainedFloatGrad {
	return &ConstrainedFloatGrad{
		ConstrainedFloat:     *NewConstrainedFloat(start, min, max),
		stepMinusDefaultStep: step - defaultStep,
	}
}

// ConstrainedNormalDistrGrad is a ConstrainedNormalDistr that implements the GeneGrad interface
// with a programmable step size.
// It implements the [mu8.GeneGrad] interface.
type ConstrainedNormalDistrGrad struct {
	ConstrainedNormalDistr
	stepMinusDefaultStep float64
}

// Step returns the step size that should be performed during gradient descent.
// It implements the [mu8.GeneGrad] interface.
func (cf *ConstrainedNormalDistrGrad) Step() float64 {
	return cf.stepMinusDefaultStep + defaultStep
}

// NewConstrainedNormalGrad returns a new ConstrainedNormalDistrGrad
// where arguments are:
//   - start: the initial value of the gene.
//   - stddev: the standard deviation of mutations.
//   - min: the minimum value the gene can take.
//   - max: the maximum value the gene can take.
//   - step: the step size that should be used during gradient descent.
func NewConstrainedNormalGrad(start, stddev, min, max, step float64) *ConstrainedNormalDistrGrad {
	return &ConstrainedNormalDistrGrad{
		ConstrainedNormalDistr: *NewConstrainedNormalDistr(start, stddev, min, max),
		stepMinusDefaultStep:   step - defaultStep,
	}
}
