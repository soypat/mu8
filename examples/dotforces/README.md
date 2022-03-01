# CodeBullet's example

This example was built following along [CodeBullet's proposed problem](https://www.youtube.com/watch?v=BOZfhUcNiqk). It attempts to minimize the time it takes for points to move towards a goal.


The following command will run an example of what genetic algorithm is possible of doing (GA built from scratch not using a pre-defined library).

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