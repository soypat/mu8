package main

const desiredFitness = 5.
const desiredModel = betterModel

// Pop represents a dot Population
type Pop struct {
	dots         []Dot
	champ        Dot
	champFitness float64
	fitness      []float64
	goal         Vec
	gen          int
}

func NewPop(size int, goal Vec) Pop {
	p := Pop{
		dots:    make([]Dot, size),
		fitness: make([]float64, size),
		goal:    goal,
	}
	for i := range p.dots {
		p.dots[i] = NewDot(Vec{})
	}
	return p
}

func (p *Pop) update(dt float64) {
	for i := range p.dots {
		dot := &p.dots[i]
		dot.update(dt)
		if dot.fitness(desiredModel, p.goal) > desiredFitness {
			// kill dots if they accomplished their life goal.
			dot.dead = true
		}
	}
}

// true when all dots are dead.
func (p *Pop) extinct() bool {
	for i := range p.dots {
		if !p.dots[i].dead {
			return false
		}
	}
	return true
}

func (p *Pop) naturalSelection() {
	p.setFitness()

	newGeneration := make([]Dot, len(p.dots))
	// First element reserved for elite/champion.
	for i := 1; i < len(p.dots); i++ {
		parent := p.selectParent()
		// ideally mating between two parents would happen here
		// and we'd splice the results. Not necessary for small simple programs
		child := NewDot(Vec{})
		child.brain = parent.brain.Clone()
		newGeneration[i] = child
	}
	p.dots = newGeneration
	// Set champion back to original position.
	p.gen++
}

func (p *Pop) mutate() {
	for i := 1; i < len(p.dots); i++ {
		p.dots[i].brain.Mutate(0.02)
	}
	p.dots[0] = p.champ
}

func (p *Pop) setFitness() {
	eliteIdx := 0
	maxFitness := 0.0
	for i := range p.dots {
		fit := p.dots[i].fitness(desiredModel, p.goal)
		if fit > maxFitness {
			maxFitness = fit
			eliteIdx = i
		}
		p.fitness[i] = fit
	}
	p.champFitness = maxFitness
	p.champ = NewDot(Vec{})
	p.champ.brain = p.dots[eliteIdx].brain.Clone()
}

func (p *Pop) fitnessSum() (fitsum float64) {
	for i := range p.dots {
		fitsum += p.fitness[i]
	}
	return fitsum
}

func (p *Pop) selectParent() Dot {
	luckOfTheFit := random(p.fitnessSum())
	runningSum := 0.0
	for i := range p.dots {
		runningSum += p.fitness[i]
		if runningSum > luckOfTheFit {
			return p.dots[i]
		}
	}
	panic("unreachable")
}

func (p *Pop) elite() Dot {
	return p.champ
}
func (p *Pop) eliteFitness() float64 {
	return p.champFitness
}
