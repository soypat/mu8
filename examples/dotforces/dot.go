package main

// distance from center which dots survive.
const (
	dotSafeLimits = 100.
	nerves        = 1000
)

type Dot struct {
	pos, vel, accel Vec
	brain           Brain
	dead            bool
}

func NewDot(v Vec) Dot {
	return Dot{
		pos:   v,
		brain: NewBrain(nerves),
	}
}

func (d *Dot) step(dt float64) {
	if len(d.brain.directions) > d.brain.step {
		d.accel = d.brain.directions[d.brain.step]
		d.brain.step++
	} else {
		d.dead = true
	}
	d.vel = Add(d.vel, Scale(dt, d.accel))
	d.vel = clamp(2, d.vel)
	d.pos = Add(d.pos, Scale(dt, d.vel))
}

func (d *Dot) update(dt float64) {
	if !d.dead {
		d.step(dt)
		if Norm(d.pos) > dotSafeLimits {
			d.dead = true
		}
	}
}

type fitnessModel int

const (
	crappyModel fitnessModel = iota
	betterModel
)

func (d *Dot) fitness(m fitnessModel, goal Vec) (fitness float64) {
	fstep := float64(d.brain.step + 1) // set to at least 1 to avoid divide by zero, just in case
	distance := Norm(Sub(d.pos, goal))
	switch m {
	case crappyModel:
		fitness = 1 / distance
	case betterModel:
		fitness = fstep / (10/fstep + distance)
	}
	return fitness
}
