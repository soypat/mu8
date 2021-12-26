package mu8

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
	Splice(Gene)
	// CloneFrom copies the Gene argument into the receiver, replacing all genetic information
	// in receiving Gene.
	CloneFrom(Gene)
	// Mutate performs a random mutation on the receiver. rand is a random number between [0, 1)
	// which is usually calculated beforehand to determine if Gene is to be mutated.
	// The distribution of rand is expected to be normal.
	Mutate(rand float64)
}
