package mu8

import "math/rand"

func Mutate[T any](g Genome[T], src rand.Source, mutationRate float64) {
	if mutationRate == 0 {
		panic("can't mutate with zero mutation rate")
	}
	random := rand.New(src)
	for i := 0; i < g.Len(); i++ {
		r := random.Float64()
		if r < mutationRate {
			g.GetGene(i).Mutate(r / mutationRate)
		}
	}
}
