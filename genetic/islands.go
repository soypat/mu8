package genetic

import (
	"math/rand"

	"github.com/soypat/mu8"
)

// Not implemented yet.
type Islands[G mu8.Genome] struct {
	islands []Population[G]
}

func NewIslands[G mu8.Genome](Nislands int, individuals []G, src rand.Source, newIndividual func() G) Islands[G] {
	if len(individuals) < Nislands {
		panic("must be more individuals than islands for scheme to work")
	}
	rng := rand.New(src)
	populations := make([][]G, Nislands)

	// Set a max individual count per island so as to evenly distribute individuals.
	maxIndividuals := 1 + len(individuals)/Nislands
	for i := 0; i < len(individuals); {
		// Distribute individuals randomly across islands
		finalDest := rng.Intn(Nislands)
		if len(populations[finalDest]) < maxIndividuals {
			populations[finalDest] = append(populations[finalDest], individuals[i])
			i++
		} else if len(individuals)/(i+1) <= 2 {
			// If random append unsuccesful, append to first available island.
			for j := range populations {
				if len(populations[j]) < maxIndividuals {
					populations[j] = append(populations[j], individuals[j])
					i++
					break
				}
			}
		}
	}

	islands := make([]Population[G], Nislands)
	for i := range islands {
		islands[i] = NewPopulation(populations[i], rand.NewSource(src.Int63()), newIndividual)
	}
	return Islands[G]{
		islands: islands,
	}
}

func (is Islands[G]) Islands() []Population[G] {
	return is.islands
}
