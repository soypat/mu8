package main

import (
	"math"
	"math/rand"
)

func V(x, y, z float64) Vec {
	return Vec{X: x, Y: y, Z: z}
}

func randV(mag float64) Vec {
	return V(random(mag), random(mag), random(mag))
}

func dir(bearing, inclination float64) Vec {
	dir := Vec{}
	dir.X = math.Cos(bearing) * math.Cos(inclination)
	dir.Y = math.Sin(-bearing) * math.Cos(inclination)
	dir.Z = -math.Sin(inclination)
	dir = Scale(1/Norm(dir), dir)
	return dir
}

func clamp(abs float64, v Vec) Vec {
	return Vec{
		X: math.Copysign(math.Min(abs, math.Abs(v.X)), v.X),
		Y: math.Copysign(math.Min(abs, math.Abs(v.Y)), v.Y),
		Z: math.Copysign(math.Min(abs, math.Abs(v.Z)), v.Z),
	}
}

func random(mag float64) float64 {
	return rand.Float64() * mag
}

type Vec struct {
	X, Y, Z float64
}

// Add returns the vector sum of p and q.
func Add(p, q Vec) Vec {
	return Vec{
		X: p.X + q.X,
		Y: p.Y + q.Y,
		Z: p.Z + q.Z,
	}
}

// Sub returns the vector sum of p and -q.
func Sub(p, q Vec) Vec {
	return Vec{
		X: p.X - q.X,
		Y: p.Y - q.Y,
		Z: p.Z - q.Z,
	}
}

func Norm(p Vec) float64 {
	return math.Hypot(p.X, math.Hypot(p.Y, p.Z))
}

func Scale(f float64, p Vec) Vec {
	return Vec{
		X: f * p.X,
		Y: f * p.Y,
		Z: f * p.Z,
	}
}
