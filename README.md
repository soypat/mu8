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

### μ8 examples
See [`rocket`](./examples/rocket/main.go) for a demonstration on rocket stage optimization. 
Below is the output of said program
```
champHeight:101.619km
champHeight:138.558km
champHeight:141.474km
champHeight:141.474km
champHeight:141.474km
champHeight:141.474km
champHeight:141.474km
champHeight:141.478km
champHeight:141.478km
champHeight:141.538km
our champion: 
Stage 0: coast=135.2s, propMass=195.4kg, Δm=99.77kg/s, totalMass=395.4
Stage 1: coast=145.9s, propMass=1.2kg, Δm=0.14kg/s, totalMass=21.2
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
