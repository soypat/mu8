package main

import "github.com/soypat/mu8/genes"

// stage is a sub-genome.
type stage struct {
	massStruc float64
	isp       float64
	// Optmized parameter: total rocket fuel [kg]
	massProp *genes.ConstrainedFloat
	// Optmized parameter: mass expulsion rate [kg/s]
	deltaMass *genes.ConstrainedFloat
	// Optimized parameter: coast time [s]
	coastTime *genes.ConstrainedFloat
}

func (s *stage) Len() int {
	return 3
}

func (s *stage) GetGene(i int) (c *genes.ConstrainedFloat) {
	switch i {
	case 0:
		c = s.massProp
	case 1:
		c = s.deltaMass
	case 2:
		c = s.coastTime
	default:
		panic("unreachable gene index")
	}
	return c
}

func (s *stage) Clone() *stage {
	return &stage{
		massStruc: s.massStruc,
		isp:       s.isp,
		massProp:  s.massProp.CopyT(),
		deltaMass: s.deltaMass.CopyT(),
		coastTime: s.coastTime.CopyT(),
	}
}
