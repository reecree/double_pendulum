package models

import (
	"fmt"
	"math"
	"testing"
)

func TestPendulum(t *testing.T) {
	dp := NewDP(math.Pi/2, 0)

	for i := 0; i < 10; i++ {
		//x, y := dp.Pendulum1.Ball.GetLocation()
		fmt.Println(dp.Pendulum1.String())
		fmt.Println(dp.Pendulum2.String())
		//fmt.Println(x, y)
		dp.Modify()
		// t := .001
		// dp.Pendulum1.theta += dp.Pendulum1.Ball.acceleration * t
		// dp.Pendulum2.theta += dp.Pendulum2.Ball.acceleration * t
	}
	t.Error("dd")
}
