package mu8

import "math/rand"

// Mutate mutates the Genes in the Genome g, modifying g in place.
// The probability of a Gene being mutated is mutationRate/1.
func Mutate(g Genome, src rand.Source, mutationRate float64) {
	if mutationRate == 0 {
		panic("can't mutate with zero mutation rate")
	}
	rng := rand.New(src)
	for i := 0; i < g.Len(); i++ {
		r := rng.Float64()
		if r < mutationRate {
			g.GetGene(i).Mutate(rng)
		}
	}
}
