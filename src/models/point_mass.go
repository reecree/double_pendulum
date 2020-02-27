package models

import "math"

type PointMass struct {
	locationX    float64
	locationY    float64
	velocityX    float64
	velocityY    float64
	acceleration float64
	mass         float64
	moment       float64
	direction    float64
}

func (p *PointMass) SetPosition(x, y float64) {
	p.locationX = x
	p.locationY = y
}

func (p *PointMass) SetVelocity(x, y float64) {
	p.velocityX = x
	p.velocityY = y
}

func (p *PointMass) GetLength() float64 {
	return math.Sqrt(p.locationX*p.locationX + p.locationY*p.locationY)
}

func (p *PointMass) GetVelocity() float64 {
	return p.direction * math.Sqrt(p.velocityX*p.velocityX+p.velocityY*p.velocityY)
}

func (p *PointMass) GetKineticEnergy() float64 {
	return p.GetRotationalEnergy() + p.GetTranslationalEnergy()
}

func (p *PointMass) GetTranslationalEnergy() float64 {
	return 0.5 * p.mass * ((p.velocityX * p.velocityX) + (p.velocityY * p.velocityY))
}

func (p *PointMass) GetRotationalEnergy() float64 {
	angularVelocity := p.GetVelocity() / p.GetLength()
	return 0.5 * p.GetMomentAboutCM() * angularVelocity * angularVelocity
}

func (p *PointMass) GetMomentAboutCM() float64 {
	return p.mass * p.moment
}
