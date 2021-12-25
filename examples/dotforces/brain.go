package main

import (
	"math"
)

type Brain struct {
	directions []Vec
	step       int
}

func NewBrain(n int) Brain {
	b := Brain{
		directions: make([]Vec, n),
	}
	b.Randomize()
	return b
}

func (b *Brain) Randomize() {
	for i := range b.directions {
		randomAngle := random(2 * math.Pi)
		b.directions[i] = dir(randomAngle, 0)
	}
}

func (b Brain) Clone() (clone Brain) {
	clone.directions = make([]Vec, len(b.directions))
	copy(clone.directions, b.directions)
	return clone
}

// mutation rate between 0 and 1.
func (b *Brain) Mutate(mutationRate float64) {
	for i := range b.directions {
		if random(1) < mutationRate {
			randomAngle := random(2 * math.Pi)
			b.directions[i] = dir(randomAngle, 0)
		}
	}
}
