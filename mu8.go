package mu8

// Genome contains one or more Gene instances which
type Genome interface {
	// Breed breeds receiver Genome with other genomes by splicing.
	// An argument of no genomes returns a non-referential copy of the receiver,
	// which could be described as a cloning procedure.
	Breed(gens ...Genome) Genome
	Mutate(mutationRate float64)
	// Simulate runs the backing simulation which the genetic
	// algorithm seeks to optimize. It returns a number quantifying
	// how well the Genome did in the simulation. This is then
	// used to compare between other Genomes during the Selection phase.
	Simulate() (fitness float64)

	Len() int
	GetGene(i int) Gene
}

// Gene is the basic physical and functional unit of heredity.
type Gene interface {
	// Splice modifies the receiver with the attributes of the argument.
	Splice(Gene)
	// Copy returns a copy of the gene so that modifying the receiver is not reflected in the returned parameter.
	Copy() Gene
}
