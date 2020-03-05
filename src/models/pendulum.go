package models

import (
	"fmt"
	"math"
)

type Pendulum struct {
	length     float64
	theta      float64
	thetaPrime float64
	Ball       PointMass
}

func (p *Pendulum) SetPosition(x, y float64) {
	p.Ball.SetPosition(x, y)
}

func (p *Pendulum) SetVelocity(x, y float64) {
	p.Ball.SetVelocity(x, y)
}

type DoublePendulum struct {
	Pendulum1       *Pendulum
	Pendulum2       *Pendulum
	potentialOffset float64

	gravity float64
	originX float64
	originY float64
}

func (p *Pendulum) String() string {
	return fmt.Sprintf("Length: %f, Theta: %f, Ball: %s, ", p.length, p.theta, p.Ball.String())
}

func (dp *DoublePendulum) getEnergy() (float64, float64) {
	var (
		L1 = dp.Pendulum1.length
		L2 = dp.Pendulum2.length
		ke = dp.Pendulum1.Ball.GetKineticEnergy() + dp.Pendulum2.Ball.GetKineticEnergy()
		// lowest point that bob1 can be is -L1, define that as zero potential energy
		// lowest point that bob2 can be is -L1 -L2
		y1 = dp.Pendulum1.Ball.locationY
		y2 = dp.Pendulum2.Ball.locationY

		pe = dp.gravity*dp.Pendulum1.Ball.mass*(y1+L1) +
			dp.gravity*dp.Pendulum2.Ball.mass*(y2+L1+L2)
	)
	return pe + dp.potentialOffset, ke
}

func (dp *DoublePendulum) SetPotentialEnergy(val float64) {
	dp.potentialOffset = 0
	pe, _ := dp.getEnergy()
	dp.potentialOffset = val - pe
}

func (dp *DoublePendulum) evaluate() []float64 {
	var (
		vals = make([]float64, 4)
		th1  = dp.Pendulum1.theta
		dth1 = dp.Pendulum1.thetaPrime
		th2  = dp.Pendulum2.theta
		dth2 = dp.Pendulum2.thetaPrime
		m1   = dp.Pendulum1.Ball.mass
		m2   = dp.Pendulum2.Ball.mass
		L1   = dp.Pendulum1.length
		L2   = dp.Pendulum2.length
		g    = dp.gravity
	)

	vals[0] = dth1
	num := -g * (2*m1 + m2) * math.Sin(th1)
	num = num - g*m2*math.Sin(th1-2*th2)
	num = num - 2*m2*dth2*dth2*L2*math.Sin(th1-th2)
	num = num - m2*dth1*dth1*L1*math.Sin(2*(th1-th2))
	num = num / (L1 * (2*m1 + m2 - m2*math.Cos(2*(th1-th2))))
	vals[1] = num

	vals[2] = dth2
	num = (m1 + m2) * dth1 * dth1 * L1
	num = num + g*(m1+m2)*math.Cos(th1)
	num = num + m2*dth2*dth2*L2*math.Cos(th1-th2)
	num = num * 2 * math.Sin(th1-th2)
	num = num / (L2 * (2*m1 + m2 - m2*math.Cos(2*(th1-th2))))
	vals[3] = num
	return vals
}

func (dp *DoublePendulum) Modify() {
	// limit the pendulum angle to +/- Pi
	theta1 := limitAngle(dp.Pendulum1.theta)
	if theta1 != dp.Pendulum1.theta {
		dp.Pendulum1.theta = theta1
	}
	theta2 := limitAngle(dp.Pendulum2.theta)
	if theta2 != dp.Pendulum2.theta {
		dp.Pendulum2.theta = theta2
	}

	dp.move()
	res := dp.evaluate()
	dp.Pendulum1.Ball.acceleration = res[1]
	dp.Pendulum2.Ball.acceleration = res[3]
}

func (dp *DoublePendulum) move() {
	var (
		sinTheta1 = math.Sin(dp.Pendulum1.theta)
		cosTheta1 = math.Cos(dp.Pendulum1.theta)
		sinTheta2 = math.Sin(dp.Pendulum2.theta)
		cosTheta2 = math.Cos(dp.Pendulum2.theta)
		L1        = dp.Pendulum1.length
		L2        = dp.Pendulum2.length
		x1        = L1 * sinTheta1
		y1        = -L1 * cosTheta1
		x2        = x1 + L2*sinTheta2
		y2        = y1 - L2*cosTheta2
	)

	dp.Pendulum1.SetPosition(x1, y1)
	dp.Pendulum2.SetPosition(x2, y2)

	var (
		v1x = dp.Pendulum1.thetaPrime * L1 * cosTheta1
		v1y = dp.Pendulum1.thetaPrime * L1 * sinTheta1
		v2x = v1x + dp.Pendulum2.thetaPrime*L2*cosTheta2
		v2y = v1y + dp.Pendulum2.thetaPrime*L2*sinTheta2
	)
	dp.Pendulum1.SetVelocity(v1x, v1y)
	dp.Pendulum2.SetVelocity(v2x, v2y)
}

func NewDP(theta1, theta2 float64) *DoublePendulum {
	// x2 := x
	// y2 := y - 1
	dp := DoublePendulum{
		Pendulum1: &Pendulum{
			length: 1,
			theta:  theta1,
			Ball:   PointMass{mass: 2},
		},
		Pendulum2: &Pendulum{
			length: 1,
			theta:  theta2,
			Ball:   PointMass{mass: 2},
		},
		gravity: 9.8,
	}
	// dp.Pendulum1.SetPosition(x, y)
	// dp.Pendulum2.SetPosition(x2, y2)
	return &dp
}

func limitAngle(angle float64) float64 {
	if angle > math.Pi {
		n := math.Floor((angle - -math.Pi) / (2 * math.Pi))
		return angle - 2*math.Pi*n
	} else if angle < -math.Pi {
		n := math.Floor(-(angle - math.Pi) / (2 * math.Pi))
		return angle + 2*math.Pi*n
	}
	return angle
}

// // RungeKutta
func (dp *DoublePendulum) Step(stepSize float64) {
	temp := *dp

	k1 := temp.evaluate()
	temp.Pendulum1.theta = dp.Pendulum1.theta + k1[0]*stepSize/2
	temp.Pendulum1.thetaPrime = dp.Pendulum1.thetaPrime + k1[1]*stepSize/2
	temp.Pendulum2.theta = dp.Pendulum2.theta + k1[2]*stepSize/2
	temp.Pendulum2.thetaPrime = dp.Pendulum2.thetaPrime + k1[3]*stepSize/2
	k2 := temp.evaluate()
	temp.Pendulum1.theta = dp.Pendulum1.theta + k2[0]*stepSize/2
	temp.Pendulum1.thetaPrime = dp.Pendulum1.thetaPrime + k2[1]*stepSize/2
	temp.Pendulum2.theta = dp.Pendulum2.theta + k2[2]*stepSize/2
	temp.Pendulum2.thetaPrime = dp.Pendulum2.thetaPrime + k2[3]*stepSize/2
	k3 := temp.evaluate()
	temp.Pendulum1.theta = dp.Pendulum1.theta + k3[0]*stepSize
	temp.Pendulum1.thetaPrime = dp.Pendulum1.thetaPrime + k3[1]*stepSize
	temp.Pendulum2.theta = dp.Pendulum2.theta + k3[2]*stepSize
	temp.Pendulum2.thetaPrime = dp.Pendulum2.thetaPrime + k3[3]*stepSize
	k4 := temp.evaluate()
	dp.Pendulum1.theta += (k1[0] + 2*k2[0] + 2*k3[0] + k4[0]) * stepSize / 6
	dp.Pendulum1.thetaPrime += (k1[1] + 2*k2[1] + 2*k3[1] + k4[1]) * stepSize / 6
	dp.Pendulum2.theta += (k1[2] + 2*k2[2] + 2*k3[2] + k4[2]) * stepSize / 6
	dp.Pendulum2.thetaPrime += (k1[3] + 2*k2[3] + 2*k3[3] + k4[3]) * stepSize / 6
}
