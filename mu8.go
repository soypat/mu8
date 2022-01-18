package mu8

import (
	"errors"
	"math/rand"
)

// Genome represents a candidate for genetic algorithm selection.
// It is parametrized with the backing Gene type.
type Genome interface {
	// Simulate runs the backing simulation which the genetic
	// algorithm seeks to optimize. It returns a number quantifying
	// how well the Genome did in the simulation. This is then
	// used to compare between other Genomes during the Selection phase.
	Simulate() (fitness float64)

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

	// Mutate performs a random mutation on the receiver. rand is a random number between [0, 1)
	// which is usually calculated beforehand to determine if Gene is to be mutated.
	// The rng argument intends to aid with randomness and Mutate implementation process.
	Mutate(rng *rand.Rand)
}

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
