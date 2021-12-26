package mu8

import "math/rand"

func Mutate[T any](g Genome[T], src rand.Source) {
	random := rand.New(src)
	for i := 0; i < g.Len(); i++ {
		g.GetGene(i).Mutate(random.Float64())
	}
}
