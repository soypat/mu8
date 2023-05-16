package genes

import "github.com/soypat/mu8"

const defaultStep = 5e-7

var (
	_ mu8.GeneGrad = (*ConstrainedFloatGrad)(nil)
)

type ConstrainedFloatGrad struct {
	ConstrainedFloat
	stepMinusDefaultStep float64
}

func (cf *ConstrainedFloatGrad) Step() float64 {
	return cf.stepMinusDefaultStep + defaultStep
}

func NewConstrainedFloatGrad(start, min, max, step float64) *ConstrainedFloatGrad {
	return &ConstrainedFloatGrad{
		ConstrainedFloat:     *NewConstrainedFloat(start, min, max),
		stepMinusDefaultStep: step - defaultStep,
	}
}

type ConstrainedNormalDistrGrad struct {
	ConstrainedNormalDistr
	stepMinusDefaultStep float64
}

func (cf *ConstrainedNormalDistrGrad) Step() float64 {
	return cf.stepMinusDefaultStep + defaultStep
}
func NewConstrainedNormalGrad(start, stddev, min, max, step float64) *ConstrainedNormalDistrGrad {
	return &ConstrainedNormalDistrGrad{
		ConstrainedNormalDistr: *NewConstrainedNormalDistr(start, stddev, min, max),
		stepMinusDefaultStep:   step - defaultStep,
	}
}
