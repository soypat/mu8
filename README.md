# μ8
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
* `genetic` directory contains genetic algorithm implementation
* `genes` contains useful `Gene` interface implementations.

### Example
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
