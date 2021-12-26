package mu8

import "math/rand"

func Mutate(g Genome, src rand.Source, mutationRate float64) {
	if mutationRate == 0 {
		panic("can't mutate with zero mutation rate")
	}
	for i := 0; i < g.Len(); i++ {
		r := randfloat(src)
		if r < mutationRate {
			g.GetGene(i).Mutate(r / mutationRate)
		}
	}
}

func randfloat(r rand.Source) float64 {
again:
	f := float64(r.Int63()) / (1 << 63)
	if f == 1 {
		goto again // resample; this branch is taken O(never)
	}
	return f
}
