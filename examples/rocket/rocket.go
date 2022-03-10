package main

import (
	"context"
	"fmt"
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

func (r *rocket) String() (output string) {
	for k := range r.stages {
		stage := r.stages[k]
		propmass := stage.massProp.Value()
		output += fmt.Sprintf("\nStage %d: coast=%.1fs, propMass=%.1fkg, Δm=%.2fkg/s, totalMass=%.1f",
			k, stage.coastTime.Value(), propmass, stage.deltaMass.Value(), stage.massStruc+propmass)
	}
	return output
}

func (r *rocket) Simulate(context.Context) (fitness float64) {
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

func (r *rocket) Clone() *rocket {
	clone := &rocket{
		CD:          r.CD,
		payloadMass: r.payloadMass,
	}
	for k := range r.stages {
		clone.stages[k] = *r.stages[k].Clone()
	}
	return clone
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
