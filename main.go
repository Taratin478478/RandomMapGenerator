package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/valyala/fastrand"
	"image"
	"log"
	"strconv"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
	//balancingIterations = 10
	pixelSize = 10
)

type Game struct {
	i          int
	pixSlice   [screenHeight / pixelSize][screenWidth / pixelSize]uint32
	pixSlices  [50][screenHeight / pixelSize][screenWidth / pixelSize]uint32
	noiseImage *image.RGBA
	level      int
}

func (g *Game) init() error {
	g.level = 10
	return nil
}

func (g *Game) UpdatePixels(j int) error {
	//filling screen pixels
	var rc, gc, bc uint8
	for i1 := 0; i1 < screenHeight/pixelSize; i1++ {
		for i2 := 0; i2 < screenWidth/pixelSize; i2++ {
			n := uint8(g.pixSlices[j][i1][i2])
			if n > 137 { //mountains
				rc, gc, bc = 20, 15, 11
			} else if n > 132 { //hills
				rc, gc, bc = 159, 130, 0
			} else if n > 127 { //upper meadows
				rc, gc, bc = 50, 153, 0
			} else if n > 122 { //lower meadows
				rc, gc, bc = 86, 222, 71
			} else if n > 117 { //shallow water
				rc, gc, bc = 0, 127, 255
			} else { //deep water
				rc, gc, bc = 51, 51, 255
			}
			for i3 := 0; i3 < pixelSize; i3++ {
				for i4 := 0; i4 < pixelSize; i4++ {
					g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))] = rc
					g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))+1] = gc
					g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))+2] = bc
					g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))+3] = 0xff
				}
			}
		}
	}
	return nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//fill matrix with random numbers
		g.level = 10
		for i1 := 0; i1 < screenHeight/pixelSize; i1++ {
			for i2 := 0; i2 < screenWidth/pixelSize; i2++ {
				x := uint32(uint8(fastrand.Uint32()))
				g.pixSlices[0][i1][i2] = x
			}
		}
		//balancing {balancingIterations} times
		for j := 1; j < 50; j++ {
			for i1 := 0; i1 < screenHeight/pixelSize; i1++ {
				for i2 := 0; i2 < screenWidth/pixelSize; i2++ {
					if i1 == 0 {
						if i2 == 0 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1] + g.pixSlices[j-1][i1+1][i2] + g.pixSlices[j-1][i1+1][i2+1]) / 4
						} else if i2 == screenWidth/pixelSize-1 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1+1][i2-1] + g.pixSlices[j-1][i1+1][i2]) / 4
						} else {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1] + g.pixSlices[j-1][i1+1][i2-1] + g.pixSlices[j-1][i1+1][i2] + g.pixSlices[j-1][i1+1][i2+1]) / 6
						}
					} else if i1 == screenHeight/pixelSize-1 {
						if i2 == 0 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1-1][i2+1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1]) / 4
						} else if i2 == screenWidth/pixelSize-1 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2-1] + g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2]) / 4
						} else {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2-1] + g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1-1][i2+1] + g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1]) / 6
						}

					} else if i2 == 0 {
						g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1-1][i2+1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1] + g.pixSlices[j-1][i1+1][i2] + g.pixSlices[j-1][i1+1][i2+1]) / 6
					} else if i2 == screenWidth/pixelSize-1 {
						g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2-1] + g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1+1][i2-1] + g.pixSlices[j-1][i1+1][i2]) / 6
					} else {
						g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2-1] + g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1-1][i2+1] + g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1] + g.pixSlices[j-1][i1+1][i2-1] + g.pixSlices[j-1][i1+1][i2] + g.pixSlices[j-1][i1+1][i2+1]) / 9
					}
				}
			}
			err := g.UpdatePixels(g.level)
			if err != nil {
				return err
			}

		}
		//filling screen pixels
		g.i++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.level < 49 {
			g.level++
			err := g.UpdatePixels(g.level)
			if err != nil {
				return err
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.level > 0 {
			g.level--
			err := g.UpdatePixels(g.level)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.noiseImage.Pix)
	txt := "Iteration " + strconv.Itoa(g.i)
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s\nLevel: %d\nSpace - new map\nArrows - change level", ebiten.CurrentTPS(), ebiten.CurrentFPS(), txt, g.level)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RandomMaps")
	ebiten.SetMaxTPS(10)
	g := &Game{
		noiseImage: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	err := g.init()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
