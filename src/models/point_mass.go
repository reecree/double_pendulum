package models

type PointMass struct {
	locationX    float64
	locationY    float64
	velocityX    float64
	velocityY    float64
	acceleration float64
	mass         float64
	moment       float64
}

func (p *PointMass) SetPosition(x, y float64) {
	p.locationX = x
	p.locationY = y
}

func (p *PointMass) SetVelocity(x, y float64) {
	p.velocityX = x
	p.velocityY = y
}

func (p *PointMass) GetKineticEnergy() float64 {
	return p.GetRotationalEnergy()+p.GetTranslationalEnergy()
}

func (p *PointMass) GetTranslationalEnergy() float64 {
	return 0.5 * p.mass * p.v this.velocity_.lengthSquared();
}

func (p *PointMass) GetRotationalEnergy() float64 {
	return 0.5*p.GetMomentAboutCM()*this.angular_velocity_*this.angular_velocity_;
}

func (p *PointMass) GetMomentAboutCM() float64 {
	return  p.mass * p.moment
}
