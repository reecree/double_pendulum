package models

import "math"

type Pendulum struct {
	length     float64
	theta      float64
	thetaPrime float64
	ball       PointMass
}

func (p *Pendulum) SetPosition(x, y float64) {
	p.ball.SetPosition(x, y)
}

func (p *Pendulum) SetVelocity(x, y float64) {
	p.ball.SetVelocity(x, y)
}

type DoublePendulum struct {
	pendulum1       *Pendulum
	pendulum2       *Pendulum
	potentialOffset float64

	gravity float64
	originX float64
	originY float64
}

func (dp *DoublePendulum) getEnergy() (float64, float64) {
	var (
		L1 = dp.pendulum1.length
		L2 = dp.pendulum2.length
		ke = dp.pendulum1.ball.GetKineticEnergy() + dp.pendulum2.ball.GetKineticEnergy()
		// lowest point that bob1 can be is -L1, define that as zero potential energy
		// lowest point that bob2 can be is -L1 -L2
		y1 = dp.pendulum1.ball.locationY
		y2 = dp.pendulum2.ball.locationY

		pe = dp.gravity*dp.pendulum1.ball.mass*(y1+L1) +
			dp.gravity*dp.pendulum2.ball.mass*(y2+L1+L2)
	)
	return pe + dp.potentialOffset, ke
}

func (dp *DoublePendulum) SetPotentialEnergy(val float64) {
	dp.potentialOffset = 0
	pe, _ := dp.getEnergy()
	dp.potentialOffset = val - pe
}

func (dp *DoublePendulum) evaluate() {
	var (
		th1  = dp.pendulum1.theta
		dth1 = dp.pendulum1.thetaPrime
		th2  = dp.pendulum2.theta
		dth2 = dp.pendulum2.thetaPrime
		m1   = dp.pendulum1.ball.mass
		m2   = dp.pendulum2.ball.mass
		L1   = dp.pendulum1.length
		L2   = dp.pendulum2.length
		g    = dp.gravity
	)

	num := -g * (2*m1 + m2) * math.Sin(th1)
	num = num - g*m2*math.Sin(th1-2*th2)
	num = num - 2*m2*dth2*dth2*L2*math.Sin(th1-th2)
	num = num - m2*dth1*dth1*L1*math.Sin(2*(th1-th2))
	num = num / (L1 * (2*m1 + m2 - m2*math.Cos(2*(th1-th2))))
	dp.pendulum1.ball.acceleration = num

	num = (m1 + m2) * dth1 * dth1 * L1
	num = num + g*(m1+m2)*math.Cos(th1)
	num = num + m2*dth2*dth2*L2*math.Cos(th1-th2)
	num = num * 2 * math.Sin(th1-th2)
	num = num / (L2 * (2*m1 + m2 - m2*math.Cos(2*(th1-th2))))
	dp.pendulum2.ball.acceleration = num
}

func (dp *DoublePendulum) modify() {
	// limit the pendulum angle to +/- Pi
	theta1 := limitAngle(dp.pendulum1.theta)
	if theta1 != dp.pendulum1.theta {
		dp.pendulum1.theta = theta1
		//this.getVarsList().setValue(0, theta1, /*continuous=*/false);
		//vars[0] = theta1;
	}
	theta2 := limitAngle(dp.pendulum2.theta)
	if theta2 != dp.pendulum2.theta {
		dp.pendulum2.theta = theta2
		//this.getVarsList().setValue(0, theta2, /*continuous=*/false);
		//vars[0] = theta2;
	}
	// update the variables that track energy
	//   0        1       2        3        4      5      6   7   8    9
	// theta1, theta1', theta2, theta2', accel1, accel2, KE, PE, TE, time
	dp.move()
	dp.evaluate()
	//dp.getEnergy()
	//var ei = this.getEnergyInfo_(vars);
	//   dp.
	//   vars[6] = ei.getTranslational();
	//   vars[7] = ei.getPotential();
	//   vars[8] = ei.getTotalEnergy();
	//   va.setValues(vars, /*continuous=*/true);
}

func (dp *DoublePendulum) move() {
	var (
		sinTheta1 = math.Sin(dp.pendulum1.theta)
		cosTheta1 = math.Cos(dp.pendulum1.theta)
		sinTheta2 = math.Sin(dp.pendulum2.theta)
		cosTheta2 = math.Cos(dp.pendulum2.theta)
		L1        = dp.pendulum1.length
		L2        = dp.pendulum2.length
		x1        = L1 * sinTheta1
		y1        = -L1 * cosTheta1
		x2        = x1 + L2*sinTheta2
		y2        = y1 - L2*cosTheta2
	)

	dp.pendulum1.SetPosition(x1, y1)
	dp.pendulum2.SetPosition(x2, y2)

	var (
		v1x = dp.pendulum1.thetaPrime * L1 * cosTheta1
		v1y = dp.pendulum1.thetaPrime * L1 * sinTheta1
		v2x = v1x + dp.pendulum2.thetaPrime*L2*cosTheta2
		v2y = v1y + dp.pendulum2.thetaPrime*L2*sinTheta2
	)
	dp.pendulum1.SetVelocity(v1x, v1y)
	dp.pendulum2.SetVelocity(v2x, v2y)
}

func NewDP(x1, y1 float64) *DoublePendulum {
	dp := DoublePendulum{
		pendulum1: &Pendulum{
			length: 1,
			theta:  math.Pi / 4,
			ball:   PointMass{},
		},
		pendulum2: &Pendulum{
			length: 1,
			theta:  math.Pi / 8,
			ball:   PointMass{},
		},
		gravity: 9.8,
	}
	dp.pendulum1.SetPosition(x1, y1)
	dp.pendulum2.SetPosition(x1, y1-1)

	return nil
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
