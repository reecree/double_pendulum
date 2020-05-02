// Copyright 2020 Chris Ruehs

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"image/color"
	"syscall/js"
	"time"

	"github.com/reecree/double_pendulum/models"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
)

type gameState struct {
	dp *models.DoublePendulum
}

var (
	done chan struct{}

	cvs    *canvas.Canvas2d
	width  float64
	height float64
)

const ballSize = 2
const stepSize = 0.00250

var gs = gameState{}

// This specifies how long a delay between calls to 'render'.     To get Frame Rate,   1s / renderDelay
var renderDelay time.Duration = 20 * time.Millisecond

func main() {

	FrameRate := time.Second / renderDelay
	println("FPS:", FrameRate)

	cvs, _ = canvas.NewCanvas2d(false)
	// Make Canvas 90% of window size.
	cvs.Create(
		int(js.Global().Get("innerWidth").Float()*0.9),
		int(js.Global().Get("innerHeight").Float()*0.9),
	)

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(60, Render)
	<-done
}

// Render is called from the 'requestAnnimationFrame' function.
// It may also be called seperatly from a 'doEvery' function, if the user
// prefers drawing to be seperate from the annimationFrame callback
func Render(gc *draw2dimg.GraphicContext) bool {
	if gs.dp == nil {
		gs.dp = models.NewDP(1, 0)
	}

	// if gs.laserX+gs.directionX > width-gs.laserSize || gs.laserX+gs.directionX < gs.laserSize {
	// 	gs.directionX = -gs.directionX
	// }
	// if gs.laserY+gs.directionY > height-gs.laserSize || gs.laserY+gs.directionY < gs.laserSize {
	// 	gs.directionY = -gs.directionY
	// }

	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.Clear()
	// move red laser

	gs.dp.Step(stepSize)
	gs.dp.Modify()

	// gs.laserX += gs.directionX
	// gs.laserY += gs.directionY

	// draws red 🔴 laser
	gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})

	gc.BeginPath()
	x1, y1 := gs.dp.Pendulum1.Ball.GetLocation()
	x2, y2 := gs.dp.Pendulum2.Ball.GetLocation()
	x1, y1 = shift(x1, -y1, 6)
	x2, y2 = shift(x2, -y2, 6)

	z1, z2 := shift(0, 0, 6)

	draw2dkit.Circle(gc, x1, y1, ballSize)
	draw2dkit.Circle(gc, x2, y2, ballSize)
	draw2dkit.Circle(gc, z1, z2, 5)

	gc.FillStroke()
	gc.Close()

	return true
}

// func shift(x, y, max float64) (float64, float64) {
// 	return ballSize + (width-ballSize*2)/2 + x*(width-ballSize*2)/(max*2),
// 		ballSize + (height-ballSize*2)/2 + y*(height-ballSize*2)/(max*2)
// }
func shift(x, y, max float64) (float64, float64) {
	return width/2 + x*width/max, height/2 + y*height/max
}
