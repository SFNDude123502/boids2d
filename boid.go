package main

import (
	"math"
)

type Boid struct {
	Loc Vector // position
	V   Vector // velocity
}

func (b1 *Boid) Dist(b2 *Boid) float64 {
	return math.Sqrt(math.Pow(b2.Loc.X-b1.Loc.X, 2) + math.Pow(b2.Loc.Y-b1.Loc.Y, 2))
}

func (b *Boid) Cohesion(flock []*Boid) Vector { // moving towards nearby flock members
	if flockSize == 1 {
		return Vector{0, 0}
	}
	var totalCoords *Vector = &Vector{0, 0}
	var neighbors = 0

	for _, b1 := range flock {
		if detectRadius >= b.Dist(b1) && b.Dist(b1) > 0 {
			totalCoords.add(b1.Loc)
			neighbors++
		}
	}
	if neighbors == 0 {
		return Vector{0, 0}
	}
	var localCenter = &Vector{
		X: totalCoords.X / float64(neighbors),
		Y: totalCoords.Y / float64(neighbors),
	}
	localCenter.add(Vector{X: -b.Loc.X, Y: -b.Loc.Y})
	localCenter.scale(cohesionFactor)
	return *localCenter
}

func (b *Boid) Seperation(flock []*Boid) Vector { // personal space from VERY nearby flock members
	var out *Vector = &Vector{0, 0}
	for _, bi := range flock {
		if b.Dist(bi) < personalSpace {
			out.add(Vector{X: b.Loc.X - bi.Loc.X, Y: b.Loc.Y - bi.Loc.Y})
		}
	}
	out.scale(avoidanceFactor)
	return *out
}

func (b *Boid) Allignment(flock []*Boid) Vector { // matching speed and direction of nearby flock members
	if flockSize == 1 {
		return Vector{0, 0}
	}
	var totalVelo *Vector = &Vector{0, 0}
	var neighbors = 0

	for _, b1 := range flock {
		if detectRadius >= b.Dist(b1) && b.Dist(b1) > 0 {
			totalVelo.add(b1.V)
			neighbors++
		}
	}
	if neighbors == 0 {
		return Vector{0, 0}
	}

	var avgVelo = &Vector{
		X: totalVelo.X / float64(neighbors),
		Y: totalVelo.Y / float64(neighbors),
	}
	avgVelo.add(Vector{X: -b.V.X, Y: -b.V.Y})
	avgVelo.scale(allignmentFactor) // 5% of the total diference between speed and avrg speed of flock
	return *avgVelo
}

func (b *Boid) Bounds() Vector { // not running into/through walls
	var out Vector = Vector{0, 0}

	if b.Loc.X < margin {
		out.X += turnSpeed
	} else if b.Loc.X > 1920-margin {
		out.X -= turnSpeed
	}

	if b.Loc.Y < margin {
		out.Y += turnSpeed
	} else if b.Loc.Y > 1080-margin {
		out.Y -= turnSpeed
	}

	return out
}

func (b *Boid) EnforceSpeedLimit() { // self-explanitory
	var speed = math.Sqrt(b.V.X*b.V.X + b.V.Y*b.V.Y)
	if speed > speedLimit {
		b.V.X *= speedLimit / speed
		b.V.Y *= speedLimit / speed
	}
}
