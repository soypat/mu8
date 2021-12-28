# μ8

<img align="right" width="220px" src="https://user-images.githubusercontent.com/26156425/147430929-bd9adebd-9c00-4ee2-a5bd-bc8642ee9a82.png">

Genetic algorithm for machine learning. 
Inspired by [CodeBullets amazing video](https://www.youtube.com/watch?v=BOZfhUcNiqk) on the subject.
---
_This is a work in progress_

Steps
1. Natural selection.
2. Mate.
3. Mutate babies.
4. Rinse and repeat.


### Info
* [`mu8.go`](./mu8.go) `Genome` and `Gene` interface definitions. Users should implement `Genome` interface and use `Gene` implementations from `genes` package.
* `genetic` directory contains genetic algorithm implementation
* `genes` contains useful `Gene` interface implementations.

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
for i := 0; i < Ngenerations; i++ {
		err := pop.Advance()
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

### CodeBullet's example
The following command will run an example of what genetic algorithm is possible of doing.
It is solving [CodeBullet's proposed problem](https://www.youtube.com/watch?v=BOZfhUcNiqk) of moving points 
towards a goal point.
```shell
go run ./examples/dotforces/
```
`elitescore` is the fitness or "score" of the best child in the generation. As you can see it gets larger.

```
gen 10: totalfitness=1032.37, elitescore=10.24671
gen 20: totalfitness=1325.90, elitescore=11.98153
gen 30: totalfitness=1374.32, elitescore=11.98153
... ten seconds later...
gen 300: totalfitness=2169.45, elitescore=16.26937
```
The score went from 10 to 16 with help of a genetic algorithm.

### Logo work
Gopher rendition by [Juliette Whittingslow](https://www.instagram.com/artewitty/).  
Gopher design authored by [Renee French](https://www.instagram.com/reneefrench)
is licensed by the Creative Commons Attribution 3.0 licensed.
