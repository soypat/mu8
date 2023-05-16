[![go.dev reference](https://pkg.go.dev/badge/github.com/soypat/mu8)](https://pkg.go.dev/github.com/soypat/mu8)
[![Go Report Card](https://goreportcard.com/badge/github.com/soypat/mu8)](https://goreportcard.com/report/github.com/soypat/mu8)
[![codecov](https://codecov.io/gh/soypat/mu8/branch/main/graph/badge.svg)](https://codecov.io/gh/soypat/mu8/branch/main)
[![License](https://img.shields.io/badge/License-BSD_2--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)

# μ8

<img align="right" width="190px" src="https://user-images.githubusercontent.com/26156425/147430929-bd9adebd-9c00-4ee2-a5bd-bc8642ee9a82.png">

Simple unsupervised machine learning package using Go 1.18 generics.

## User information
μ8 (mu8) uses a simple genetic algorithm implementation to optimize a objective function. It allows optimizing floating point numbers, integers and anything else that can implement the 3 method [`Gene`](./mu8.go) interface

The genetic algorithm implementation is currently ~150 lines long and is contained in [`population.go`](./genetic/population.go). It consists of the following steps:

1. Natural selection. Best individual conserved (population champion)
2. Mate.
3. Mutate babies.
4. Rinse and repeat.

The file [`mu8.go`](./mu8.go) contains `Genome` and `Gene` interface definitions. Users should implement `Genome` interface and use `Gene` implementations from [`genes`](./genes) package.

There is an Islands Model Genetic Algorithm (IMGA) implementation in [`islands.go`](./genetic/islands.go) using the `Islands` type that makes use of a parallel optimization algorithm to make use of multi-core machines.

## μ8 examples

### Basic usage example
Everything starts with the `mu8.Genome` type on the user side. We define a type that implements it
using a helper type `genes.ContrainedFloat` from the `genes` package. All this `genes` type does
is save us the trouble of writing our own `mu8.Gene` implementation.

```go
type mygenome struct {
	genoma []genes.ConstrainedFloat
}

func (g *mygenome) GetGene(i int) mu8.Gene { return &g.genoma[i] }
func (g *mygenome) Len() int               { return len(g.genoma) }

// Simulate simply adds the genes. We'd expect the genes to reach the max values of the constraint.
func (g *mygenome) Simulate() (fitness float64) {
	for i := range g.genoma {
		fitness += g.genoma[i].Value()
	}
    // fitness must ALWAYS be greater than zero for succesful simulation.
	return math.Max(0, fitness/float64(g.Len()))
}
```
We're almost ready to optimize our implementation to maximize it's fitness, which would simply be the addition of all it's genes.

Let's write the function that initializes a blank-slate `mygenome`

```go
func newGenome(n int) *mygenome {
	return &mygenome{genoma: make([]genes.ConstrainedFloat, n)}
}
```
The function above may be confusing... what is the constraint on the number? By default
`genes.ConstrainedFloat` uses the range [0, 1]. 

```go
const Nindividuals = 100
individuals := make([]*mygenome, Nindividuals)
for i := 0; i < Nindividuals; i++ {
	genome := newGenome(genomelen)
	// This spices up the initial population so fitnesses are not all zero.
	mu8.Mutate(genome, src, .1)
	individuals[i] = genome
}

pop := genetic.NewPopulation(individuals, rand.NewSource(1), func() *mygenome {
		return newGenome(3)
})

const Ngeneration = 100
ctx := context.Background()
for i := 0; i < Ngenerations; i++ {
		err := pop.Advance(ctx)
		if err != nil {
			panic(err.Error())
		}
		err = pop.Selection(0.5, 1)
		if err != nil {
			panic(err.Error())
		}
}
fmt.Printf("champ fitness=%.3f\n", pop.ChampionFitness())
```
The final fitness should be close to 1.0 if the algorithm did it's job. For the code see 
[`mu8_test.go`](./mu8_test.go)

### Rocket stage optimization example

See [`rocket`](./examples/rocket/main.go) for a demonstration on rocket stage optimization. 
Below is the output of said program
```
champHeight:117.967km
champHeight:136.748km
champHeight:140.633km
champHeight:141.873km
champHeight:141.873km
champHeight:141.873km
champHeight:142.883km
champHeight:143.292km
champHeight:143.292km
champHeight:143.292km
champHeight:143.292km
our champion: 
Stage 0: coast=281.2s, propMass=0.0kg, Δm=99.35kg/s, totalMass=200.0
Stage 1: coast=0.0s, propMass=1.6kg, Δm=0.01kg/s, totalMass=21.6
```

### Gradient "ascent" example
```go
src := rand.NewSource(1)
const (
	genomelen      = 6
	gradMultiplier = 10.0
	epochs         = 6
)
// Create new individual and mutate it randomly.
individual := newGenome(genomelen)
rng := rand.New(src)
for i := 0; i < genomelen; i++ {
	individual.GetGene(i).Mutate(rng)
}
// Prepare for gradient descent.
grads := make([]float64, genomelen)
ctx := context.Background()
// Champion will harbor our best individual.
champion := newGenome(genomelen)
for epoch := 0; epoch < epochs; epoch++ {
	// We calculate the gradients of the individual passing a nil
	// newIndividual callback since the GenomeGrad type we implemented
	// does not require blank-slate initialization.
	err := mu8.Gradient(ctx, grads, individual, nil)
	if err != nil {
		panic(err)
	}
	// Apply gradients.
	for i := 0; i < individual.Len(); i++ {
		gene := individual.GetGeneGrad(i)
		grad := grads[i]
		gene.SetValue(gene.Value() + grad*gradMultiplier)
	}
	mu8.CloneGrad(champion, individual)
	fmt.Printf("fitness=%f with grads=%f\n", individual.Simulate(ctx), grads)
}
```

## Contributing
Contributions very welcome! I myself have no idea what I'm doing so I welcome
issues on any matter :)

Pull requests also welcome but please submit an issue first on what you'd like to change.
I promise I'll answer as fast as I can.

Please take a look at the TODO's in the project: <kbd>Ctrl</kbd>+<kbd>F</kbd> `TODO`

## References
Inspired by [CodeBullets amazing video](https://www.youtube.com/watch?v=BOZfhUcNiqk) on the subject.

## Logo work
Gopher rendition by [Juliette Whittingslow](https://www.instagram.com/artewitty/).  
Gopher design authored by [Renée French](https://www.instagram.com/reneefrench)
is licensed by the Creative Commons Attribution 3.0 licensed.
