package mu8

import (
	"context"
	"errors"
	"math"
	"math/rand"
)

// Genome represents a candidate for genetic algorithm selection.
// It is parametrized with the backing Gene type.
type Genome interface {
	// Simulate runs the backing simulation which the genetic
	// algorithm seeks to optimize. It returns a number quantifying
	// how well the Genome did in the simulation. This is then
	// used to compare between other Genomes during the Selection phase.
	//
	// The input context is cancelled when the optimization is terminated early
	// by user or by an error encountered if using IMGA. The context should be used
	// for long running simulations. It should not be used to pass
	// values into the simulation
	Simulate(context.Context) (fitness float64)

	// GetGene gets ith gene in the Genome. It is expected the ith Genes
	// of two Genomes in a Genetic Algorithm instance have matching types.
	GetGene(i int) Gene

	// Number of Genes in Genome.
	Len() int
}

// Gene is the basic physical and functional unit of heredity.
type Gene interface {
	// Splice modifies the receiver with the attributes of the argument. It should NOT
	// modify the argument. Splice is called during breeding of multiple Genomes.
	// It is expected Splice receives an argument matching the type of the receiver.
	// The rng argument intends to aid with randomness and Splice implementation process.
	Splice(rng *rand.Rand, g Gene)

	// CloneFrom copies the Gene argument into the receiver, replacing all genetic information
	// in receiving Gene.
	CloneFrom(Gene)

	// Mutate performs a random mutation on the receiver with the aid of rng.
	Mutate(rng *rand.Rand)
}

// Mutate mutates the Genes in the Genome g, modifying g in place.
// The probability of a Gene being mutated is mutationRate/1.
func Mutate(g Genome, src rand.Source, mutationRate float64) {
	switch {
	case mutationRate == 0:
		panic("can't mutate with zero mutation rate")
	case mutationRate < 0 || mutationRate > 1:
		panic("mutation rate outside valid bounds 0..1")
	}

	rng := rand.New(src)
	for i := 0; i < g.Len(); i++ {
		r := rng.Float64()
		if r < mutationRate {
			g.GetGene(i).Mutate(rng)
		}
	}
}

// Clone clones the Genes of src to dst. It does not
// modify src. dst should be initialized beforehand.
func Clone(dst, src Genome) error {
	if dst == nil {
		return errors.New("got nil destination for Clone")
	} else if src == nil {
		return errors.New("got nil source to Clone")
	} else if dst.Len() != src.Len() {
		return errors.New("destination and source mismatch")
	}

	for i := 0; i < dst.Len(); i++ {
		dst.GetGene(i).CloneFrom(src.GetGene(i))
	}
	return nil
}

// GenomeGrad is a Genome that can be used with gradient descent.
type GenomeGrad interface {
	Simulate(context.Context) (fitness float64)
	GetGeneGrad(i int) GeneGrad
	Len() int
}

// GeneGrad is a Gene that can be used with gradient descent.
type GeneGrad interface {
	SetValue(float64)
	Value() float64
	Step() float64
}

// GradientDescent computes the Gradient of the GenomeGrad g using finite differences.
// It stores the result of the calculation to grad. The length of grad must match
// the number of Genes in g. The startIndividual argument is used to seed the
// individual on every run of the simulation if it is not possible to reuse
// startIndividual between simulations. If newIndividual is nil then the same individual is
// used for all runs.
func Gradient[T GenomeGrad](ctx context.Context, grad []float64, startIndividual T, newIndividual func() T) error {
	if startIndividual.Len() != len(grad) {
		panic("scratch length mismatch")
	}
	startFitness := startIndividual.Simulate(ctx)
	for i := 0; i < startIndividual.Len() && ctx.Err() == nil; i++ {
		if newIndividual != nil {
			blankSlate := newIndividual()
			CloneGrad(blankSlate, startIndividual)
			startIndividual = blankSlate
		}
		gene := startIndividual.GetGeneGrad(i)
		start := gene.Value()
		step := gene.Step()
		if step == 0 {
			return errors.New("zero step size")
		}
		gene.SetValue(start + step)
		newFitness := startIndividual.Simulate(ctx)
		if newFitness < 0 {
			return errors.New("negative fitness")
		} else if math.IsNaN(newFitness) || math.IsInf(newFitness, 0) {
			return errors.New("invalid fitness (NaN or Inf))")
		}
		grad[i] = (newFitness - startFitness) / step
		gene.SetValue(start) // Return gene to original value.
	}
	return nil
}

// CloneGrad clones all the genes of src to dst. It does not modify src.
func CloneGrad(dst, src GenomeGrad) error {
	if dst == nil {
		return errors.New("got nil destination for Clone")
	} else if src == nil {
		return errors.New("got nil source to Clone")
	} else if dst.Len() != src.Len() {
		return errors.New("destination and source mismatch")
	}

	for i := 0; i < dst.Len(); i++ {
		dst.GetGeneGrad(i).SetValue(src.GetGeneGrad(i).Value())
	}
	return nil
}
