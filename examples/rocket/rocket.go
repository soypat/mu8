package main

import (
	"math"

	"github.com/soypat/mu8"
)

// Rocket is the full genome.
type rocket struct {
	// drag coefficient. Usually around 0.3 - 0.5
	CD           float64
	payloadMass  float64
	expendedFuel [Nstages]float64
	burnoutTime  [Nstages]float64
	stages       [Nstages]stage
}

func (r *rocket) Simulate() (fitness float64) {
	currentStage := 0
	// Integration variables.
	position := 0.0
	velocity := 0.0
	// simulation time step
	dt := 0.1
	maxPosition := 0.0
	for t := 0.0; t < 100; t += dt {
		// Calculate burnout time
		propRemain := math.Max(0, r.stages[currentStage].massProp.Value()-r.expendedFuel[currentStage])
		if propRemain == 0 && r.burnoutTime[currentStage] == 0 {
			r.burnoutTime[currentStage] = t
		}
		burnout := r.burnoutTime[currentStage]
		coastTime := t - burnout

		var thrust, dm float64
		if burnout == 0 {
			// we still in business providing thrust.
			dm = r.stages[currentStage].deltaMass.Value()
			isp := r.stages[currentStage].isp
			thrust = gravity * dm * isp
		}

		mass := r.mass(currentStage)
		_, _, Rho := atmos(position)
		force := thrust + math.Copysign(0.5*Rho*r.CD*velocity*velocity, -velocity)
		accel := force/mass - gravity
		// Integrate variables
		r.expendedFuel[currentStage] -= dm * dt
		velocity = accel*dt + velocity
		position = velocity*dt + position
		if currentStage < Nstages-1 && coastTime >= r.stages[currentStage].coastTime.Value() {
			// Conditions for stage separation is that total coast time of stage be greater|equal to expected stage coast time
			currentStage++ //equivalent to a stage separation
		}
		maxPosition = math.Max(maxPosition, position)
	}
	return maxPosition / 1e3 // Fitness in kilometers
}

// Calculate current rocket mass.
func (r *rocket) mass(currentStage int) (mass float64) {
	for i := currentStage; i < Nstages; i++ {
		mass += r.stages[i].massStruc + r.stages[i].massProp.Value() - r.expendedFuel[i]
	}
	return mass + r.payloadMass
}

func (r *rocket) GetGene(i int) mu8.Gene {
	NGenesPerStage := r.stages[0].Len()
	stageIdx := i / NGenesPerStage
	geneIdx := i % NGenesPerStage
	return r.stages[stageIdx].GetGene(geneIdx)
}

func (r *rocket) Len() int {
	NGenesPerStage := r.stages[0].Len()
	return len(r.stages) * NGenesPerStage
}

func (r *rocket) Clone() mu8.Genome {
	clone := &rocket{
		CD:          r.CD,
		payloadMass: r.payloadMass,
	}
	for k := range r.stages {
		clone.stages[k] = *r.stages[k].Clone()
	}
	return clone
}
