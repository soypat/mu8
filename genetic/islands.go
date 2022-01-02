package genetic

import (
	"context"
	"math/rand"
	"sync"

	"github.com/soypat/mu8"
)

// Not implemented yet.
// Islands Model Genetic Algorithm (IMGA) is a multi-population based GA.
// IMGA aimed to avoid local optimum by maintaining population (island) diversity using migration.
type Islands[G mu8.Genome] struct {
	islands []island[G]
	rng     rand.Rand
	// Migration Window, a buffer to keep best individual from each island.
	mw []migrant[G]
}

type migrant[G mu8.Genome] struct {
	ind    G
	origin int
}

// NewIslands WARNING: EXPERIMENTAL. Not stable or ready for use.
func NewIslands[G mu8.Genome](Nislands int, individuals []G, src rand.Source, newIndividual func() G) Islands[G] {
	if Nislands <= 1 {
		panic("need at least 2 islands")
	} else if len(individuals) < Nislands {
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
					populations[j] = append(populations[j], individuals[i])
					i++
					break
				}
			}
		}
	}

	islands := make([]island[G], Nislands)
	for i := range islands {
		islands[i] = newIsland(populations[i], rand.NewSource(src.Int63()), newIndividual)

	}
	return Islands[G]{
		islands: islands,
		mw:      make([]migrant[G], len(islands)),
		rng:     *rand.New(src),
	}
}

func (is Islands[G]) Islands() []island[G] {
	return is.islands
}

func newIsland[G mu8.Genome](individuals []G, src rand.Source, newIndividual func() G) island[G] {
	return island[G]{
		prevFitness: make([]float64, len(individuals)),
		Population:  NewPopulation(individuals, src, newIndividual),
	}
}

type island[G mu8.Genome] struct {
	Population[G]
	prevFitness []float64
	attr        float64
}

// receiveMigrant replaces individual with zero or minimum
// fitness with the migrant
func (is *island[G]) receiveMigrant(migrant G) {
	minidx := -1
	minFitness := is.fitness[0] + 1
	for i := 0; i < len(is.fitness); i++ {
		fitness := is.fitness[i]
		if fitness == 0 {
			minidx = i
			break
		} else if fitness < minFitness {
			minidx = i
			minFitness = fitness
		}
	}
	mu8.Clone(is.individuals[minidx], migrant)
}

func (is *island[G]) Individuals() []G {
	return is.individuals
}

type errmsg struct {
	err error
	i   int
}

func (is *Islands[G]) Advance(ctx context.Context, mutationRate float64, polygamy, Ngen, Nconcurrent int) (err error) {
	I := len(is.islands)

	if Nconcurrent <= 0 {
		panic("concurrency must be greater than 0")
	} else if Ngen <= 0 {
		panic("number of generations must be greater or equal to 1")
	} else if Nconcurrent > I {
		panic("cannot have more goroutines than number of islands.")
	} else if Ngen <= 1 {
		panic("number of generations between crossovers should be positive and it is HIGHLY recommended it is above 1")
	}

	defer func() {
		a := recover()
		if a != nil {
			if perr, ok := a.(error); ok {
				err = perr
			} else {
				panic(a)
			}
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Concurrency limiting mechanism ensures only Nconcurrent
	// goroutines are running population.Advance at a time.
	checkin := make(chan struct{}, Nconcurrent)

	// When Ngen is implemented there shall be exactly I goroutines
	// each entrusted with it's own population to prevent data-races.
	var wg sync.WaitGroup

	errChan := make(chan errmsg, I)
	go func() {
		select {
		case <-ctx.Done():
			for i := 0; i < I; i++ {
				is.islands[i].exit <- struct{}{} // Signal to end all calls to Advance immediately.
			}
		}
	}()
	for i := 0; i < I; i++ {
		i := i // Loop variable escape for closures.
		wg.Add(1)
		go func() {
			defer wg.Done()
			for g := 0; g < Ngen; g++ {
				checkin <- struct{}{}
				err := is.islands[i].Advance()
				if err == nil && ctx.Err() != nil {
					errChan <- errmsg{ctx.Err(), i}
					return
				}
				if err != nil {
					errChan <- errmsg{err, i}
					return
				}
				err = is.islands[i].Selection(mutationRate, polygamy)
				if err == nil && ctx.Err() != nil {
					errChan <- errmsg{ctx.Err(), i}
					return
				}
				if err != nil {
					errChan <- errmsg{err, i}
					return
				}
				<-checkin
			}
		}()
	}
	wg.Wait()
	close(checkin)
	// Population error handling.
	popErrs := make([]*errmsg, I)
	for len(errChan) > 0 {
		gotErr := <-errChan
		popErrs[gotErr.i] = &gotErr
	}

	// is.updateAttractiveness()
	for i := 0; i < I; i++ {
		if popErrs[i] != nil {
			err = popErrs[i].err
			continue
		}
		mig := is.islands[i].generator()
		errclone := mu8.Clone(mig, is.islands[i].Champion())
		if errclone != nil {
			return err
		}
		is.mw[i] = migrant[G]{mig, i}
	}

	return nil
}

func (is *Islands[G]) Crossover() {
	I := len(is.islands)
	// Perform crossover.
	for i := 0; i < I; i++ {
		migrant := is.mw[i]
		for {
			j := is.rng.Intn(I)
			if migrant.origin != j {
				is.islands[j].receiveMigrant(migrant.ind)
				break
			}
		}
	}
}

func (is *Islands[G]) champIdx() int {
	maxFitness := 0.
	maxidx := -1
	for i := range is.islands {
		champFitness := is.islands[i].ChampionFitness()
		if champFitness > maxFitness {
			maxFitness = champFitness
			maxidx = i
		}
	}
	if maxidx == -1 {
		panic("all fitnesses zero. can't select champion before Advance completed")
	}
	return maxidx
}

func (is *Islands[G]) Champion() G {
	champi := is.champIdx()
	return is.islands[champi].Champion()
}

func (is *Islands[G]) ChampionFitness() float64 {
	champi := is.champIdx()
	return is.islands[champi].ChampionFitness()
}

func (is *Islands[G]) updateAttractiveness() {
	I := len(is.islands)
	for i := 0; i < I; i++ {
		sum := 0.0
		isle := &is.islands[i]
		for k := range isle.prevFitness {
			// iterate over individuals
			sum += isle.prevFitness[k] - isle.fitness[k]
		}
		Sp := float64(len(isle.prevFitness)) // TODO: should this include only "live" (non-zero) fitnesses?
		isle.attr = sum / Sp
	}
}
