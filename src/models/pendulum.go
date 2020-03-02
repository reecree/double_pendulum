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

func (dp *DoublePendulum) evaluate() {
	var (
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

	num := -g * (2*m1 + m2) * math.Sin(th1)
	num = num - g*m2*math.Sin(th1-2*th2)
	num = num - 2*m2*dth2*dth2*L2*math.Sin(th1-th2)
	num = num - m2*dth1*dth1*L1*math.Sin(2*(th1-th2))
	num = num / (L1 * (2*m1 + m2 - m2*math.Cos(2*(th1-th2))))
	dp.Pendulum1.Ball.acceleration = num

	num = (m1 + m2) * dth1 * dth1 * L1
	num = num + g*(m1+m2)*math.Cos(th1)
	num = num + m2*dth2*dth2*L2*math.Cos(th1-th2)
	num = num * 2 * math.Sin(th1-th2)
	num = num / (L2 * (2*m1 + m2 - m2*math.Cos(2*(th1-th2))))
	dp.Pendulum2.Ball.acceleration = num
}

func (dp *DoublePendulum) Modify() {
	// limit the pendulum angle to +/- Pi
	theta1 := limitAngle(dp.Pendulum1.theta)
	if theta1 != dp.Pendulum1.theta {
		dp.Pendulum1.theta = theta1
		//this.getVarsList().setValue(0, theta1, /*continuous=*/false);
		//vars[0] = theta1;
	}
	theta2 := limitAngle(dp.Pendulum2.theta)
	if theta2 != dp.Pendulum2.theta {
		dp.Pendulum2.theta = theta2
		//this.getVarsList().setValue(0, theta2, /*continuous=*/false);
		//vars[0] = theta2;
	}
	// update the variables that track energy
	//   0        1       2        3        4      5      6   7   8    9
	// theta1, theta1', theta2, theta2', accel1, accel2, KE, PE, TE, time
	dp.move()
	dp.evaluate()

	// t := .01
	// dp.Pendulum1.theta += dp.Pendulum1.Ball.acceleration * t
	// dp.Pendulum2.theta += dp.Pendulum2.Ball.acceleration * t
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
// func (dp *DoublePendulum) Step(stepSize float64) {
// 		// var error, i;
// 		// var va = this.ode_.getVarsList();
// 		// var vars = va.getValues();
// 		// var N = vars.length;
// 		// if (this.inp_.length < N) {
// 		//   this.inp_ = /** @type {!Array<number>}*/(new Array(N));
// 		//   this.k1_ = /** @type {!Array<number>}*/(new Array(N));
// 		//   this.k2_ = /** @type {!Array<number>}*/(new Array(N));
// 		//   this.k3_ = /** @type {!Array<number>}*/(new Array(N));
// 		//   this.k4_ = /** @type {!Array<number>}*/(new Array(N));
// 		// }
// 		// var inp = this.inp_;
// 		// var k1 = this.k1_;
// 		// var k2 = this.k2_;
// 		// var k3 = this.k3_;
// 		// var k4 = this.k4_;
// 		// evaluate at time t
// 		// for (i=0; i<N; i++) {
// 		//   inp[i] = vars[i];
// 		// }
// 		// Util.zeroArray(k1);

// 		error = this.ode_.evaluate(inp, k1, 0);
// 		if (error !== null) {
// 		  return error;
// 		}
// 		// evaluate at time t+stepSize/2
// 		for (i=0; i<N; i++) {
// 		  inp[i] = vars[i]+k1[i]*stepSize/2;
// 		}
// 		Util.zeroArray(k2);
// 		error = this.ode_.evaluate(inp, k2, stepSize/2);
// 		if (error !== null) {
// 		  return error;
// 		}
// 		// evaluate at time t+stepSize/2
// 		for (i=0; i<N; i++) {
// 		  inp[i] = vars[i]+k2[i]*stepSize/2;
// 		}
// 		Util.zeroArray(k3);
// 		error = this.ode_.evaluate(inp, k3, stepSize/2);
// 		if (error !== null) {
// 		  return error;
// 		}
// 		// evaluate at time t+stepSize
// 		for (i=0; i<N; i++) {
// 		  inp[i] = vars[i]+k3[i]*stepSize;
// 		}
// 		Util.zeroArray(k4);
// 		error = this.ode_.evaluate(inp, k4, stepSize);
// 		if (error !== null) {
// 		  return error;
// 		}
// 		for (i=0; i<N; i++) {
// 			vars[i] += (k1[i] + 2*k2[i] + 2*k3[i] + k4[i])*stepSize/6;
// 		}
// 		va.setValues(vars, /*continuous=*/true);
// 		return null;
// }
