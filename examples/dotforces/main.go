package main

import (
	"fmt"
)

// inspired by https://www.youtube.com/watch?v=BOZfhUcNiqk
func main() {
	const dt = 1.
	goal := Vec{X: 30, Y: 40}
	pop := NewPop(200, goal)
	lastElite := 0.0
	for pop.gen < 300 {
		for !pop.extinct() {
			pop.update(dt)
		}
		pop.naturalSelection()
		sum := pop.fitnessSum()
		currentElite := pop.eliteFitness()
		if currentElite < lastElite {
			panic("bad elite")
		}
		pop.mutate()
		if pop.gen%10 == 0 {
			fmt.Printf("gen %d: totalfitness=%.2f, elitescore=%.5f\n", pop.gen, sum, pop.eliteFitness())
		}
	}
}
