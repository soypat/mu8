package genetic

import "errors"

var (
	ErrNegativeFitness      = errors.New("fitness cannot yield negative values. Use zero instead.")
	ErrZeroFitnessSum       = errors.New("zero fitness sum: cannot make decisions")
	ErrBadPolygamy          = errors.New("bad polygamy: must be in range [0, Nindividuals)")
	ErrChampFitnessDecrease = errors.New("champion fitness should be monotonically increasing. check for preserved references in newIndividual function.")
	// This should never trigger.
	errGenomeCast = errors.New("theoretically unreachable. Bad Genome->G cast. make sure your Genomes return same type")
)
