package mu8

// Genome represents a candidate for genetic algorithm selection.
// It is parametrized with the backing Gene type.
type Genome[T any] interface {
	// Simulate runs the backing simulation which the genetic
	// algorithm seeks to optimize. It returns a number quantifying
	// how well the Genome did in the simulation. This is then
	// used to compare between other Genomes during the Selection phase.
	Simulate() (fitness float64)
	// GetGene gets ith gene in the
	GetGene(i int) Gene[T]
	// Clone produces a new copy of Genome with no past information of simulation.
	// It should ideally hold no references to receiver to prevent data leaks.
	Clone() Genome[T]
	// Number of genes in genome.
	Len() int
}

// Gene is the basic physical and functional unit of heredity.
type Gene[T any] interface {
	// Instance returns the backing gene data.
	Instance() T
	// Splice modifies the receiver with the attributes of the argument. It should NOT
	// modify the argument.
	Splice(T)
	// Copy returns a copy of the gene so that modifying the receiver is not reflected in the returned parameter.
	Copy() T
	// Mutate performs a mutation on the receiver. rand is a random number between [0, 1)
	// to aid the user with randomness. The distribution of rand is expected to be normal.
	Mutate(rand float64)
}
