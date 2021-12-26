package main

import (
	"math"
	"math/rand"

	"github.com/soypat/mu8"
	"github.com/soypat/mu8/genes"
)

const (
	// [m/s]
	gravity = 9.8
)

// Having constant length arrays lets us more dynamically program other aspects
// of this Genome. The process of cloning also becomes dead simple.
const Nstages = 2

var baseRocket = &rocket{
	payloadMass: 10,
	CD:          0.4,
	stages: [Nstages]stage{
		{
			isp:       300,
			massStruc: 20,
			massProp:  genes.NewConstrainedFloat(50, 30, 100),
			deltaMass: genes.NewConstrainedFloat(2, 1, 3),
			coastTime: genes.NewConstrainedFloat(10, 0, 300),
		},
		{
			isp:       300,
			massStruc: 200,
			massProp:  genes.NewConstrainedFloat(1500, 800, 3000),
			deltaMass: genes.NewConstrainedFloat(13, 5, 30),
			coastTime: genes.NewConstrainedFloat(10, 0, 300),
		},
	},
}

func main() {
	src := rand.NewSource(1)
	individuals := make([]*rocket, 100)
	for i := range individuals {
		clone := baseRocket.Clone()
		mu8.Mutate[*rocket](clone, src) // fak... me...
		individuals[i] = clone
	}
}

// atmosphere thermodynamic property calculation, done horribly wrong!
func atmos(height float64) (Temp, Press, Density float64) {
	const (
		baseTemp, spaceTemp = 300, 7
		baseRho, spaceRho   = 1.2, 1e-6
		baseP, spaceP       = 101325., 1e-6
	)
	// Normalize height so 0km = -2, 60km=+2 => 30km = 0. Domain ratio 60e3:4
	normalized := (height + 30e3) / (60e3 / 4)
	cmpErf := (1 + math.Erfc(normalized)) / 2

	Density = spaceRho + (baseRho-spaceRho)*cmpErf
	Temp = spaceTemp + (baseTemp-spaceTemp)*cmpErf
	Press = spaceP + (baseP-spaceP)*cmpErf
	return Temp, Press, Density
}