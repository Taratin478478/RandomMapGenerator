package main

import (
	"fmt"
	"image"
	"log"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/valyala/fastrand"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
	//balancingIterations = 10
	//g.pixelSize = 10
)

type Game struct {
	i          int
	pixSlice   [screenHeight][screenWidth]uint32
	pixSlices  [50][screenHeight][screenWidth]uint32
	noiseImage *image.RGBA
	level      int
	pixelSize  int
	x          int
	y          int
}

func (g *Game) init() error {
	g.level = 10
	g.pixelSize = 5
	return nil
}

func (g *Game) UpdatePixels(j int) error {
	//filling screen pixels
	t := time.Now()
	var rc, gc, bc uint8
	for i1 := 0; i1 < screenHeight; i1++ {
		for i2 := 0; i2 < screenWidth; i2++ {
			n := uint8(g.pixSlices[j][i1/g.pixelSize+g.y][i2/g.pixelSize+g.x])
			if n > 132 { //mountains
				rc, gc, bc = 20, 15, 11
			} else if n > 127 { //hills
				rc, gc, bc = 159, 130, 0
			} else if n > 122 { //upper meadows
				rc, gc, bc = 50, 153, 0
			} else if n > 117 { //lower meadows
				rc, gc, bc = 86, 222, 71
			} else if n > 114 { //shallow water
				rc, gc, bc = 0, 127, 255
			} else { //deep water
				rc, gc, bc = 51, 51, 255
			}
			g.noiseImage.Pix[4*((i1)*screenHeight+(i2))] = rc
			g.noiseImage.Pix[4*((i1)*screenHeight+(i2))+1] = gc
			g.noiseImage.Pix[4*((i1)*screenHeight+(i2))+2] = bc
			g.noiseImage.Pix[4*((i1)*screenHeight+(i2))+3] = 0xff
			/*
				for i3 := 0; i3 < g.pixelSize; i3++ {
					for i4 := 0; i4 < g.pixelSize; i4++ {
						g.noiseImage.Pix[4*((i1*g.pixelSize+i3)*screenHeight+(i2*g.pixelSize+i4))] = rc
						g.noiseImage.Pix[4*((i1*g.pixelSize+i3)*screenHeight+(i2*g.pixelSize+i4))+1] = gc
						g.noiseImage.Pix[4*((i1*g.pixelSize+i3)*screenHeight+(i2*g.pixelSize+i4))+2] = bc
						g.noiseImage.Pix[4*((i1*g.pixelSize+i3)*screenHeight+(i2*g.pixelSize+i4))+3] = 0xff
					}
				}

			*/
		}
	}
	fmt.Println((time.Now().Nanosecond() - t.Nanosecond()) / 1000000)
	return nil
}

func (g *Game) Update() error {
	/*
		if ebiten.IsKeyPressed(ebiten.KeyUp) && g.pixelSize < 100 {
			g.pixelSize++
			g.x = g.x + float64((screenWidth / g.pixelSize - screenWidth / (g.pixelSize + 1)) / 2)
			g.y = g.y + float64((screenHeight / g.pixelSize - screenHeight / (g.pixelSize + 1)) / 2)
			fmt.Println(g.x, g.y)
			err := g.UpdatePixels(g.level)
			if err != nil {
				return err
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.pixelSize > 1 {
			g.pixelSize--
			g.x = g.x + float64((screenWidth / g.pixelSize - screenWidth / (g.pixelSize + 1)) / 2)
			g.y = g.y + float64((screenHeight / g.pixelSize - screenHeight / (g.pixelSize + 1)) / 2)
			fmt.Println(g.x, g.y)
			err := g.UpdatePixels(g.level)
			if err != nil {
				return err
			}
		}
	*/
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.y > 0 {
		g.y--
		err := g.UpdatePixels(g.level)
		if err != nil {
			return err
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.x > 0 {
		g.x--
		err := g.UpdatePixels(g.level)
		if err != nil {
			return err
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.y < screenHeight {
		g.y++
		err := g.UpdatePixels(g.level)
		if err != nil {
			return err
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.x < screenWidth {
		g.x++
		err := g.UpdatePixels(g.level)
		if err != nil {
			return err
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//fill matrix with random numbers
		g.level = 10
		for i1 := 0; i1 < screenHeight; i1++ {
			for i2 := 0; i2 < screenWidth; i2++ {
				x := uint32(uint8(fastrand.Uint32()))
				g.pixSlices[0][i1][i2] = x
			}
		}
		//balancing {balancingIterations} times
		for j := 1; j < 50; j++ {
			for i1 := 0; i1 < screenHeight; i1++ {
				for i2 := 0; i2 < screenWidth; i2++ {
					if i1 == 0 {
						if i2 == 0 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1] + g.pixSlices[j-1][i1+1][i2] + g.pixSlices[j-1][i1+1][i2+1]) / 4
						} else if i2 == screenWidth-1 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1+1][i2-1] + g.pixSlices[j-1][i1+1][i2]) / 4
						} else {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1] + g.pixSlices[j-1][i1+1][i2-1] + g.pixSlices[j-1][i1+1][i2] + g.pixSlices[j-1][i1+1][i2+1]) / 6
						}
					} else if i1 == screenHeight-1 {
						if i2 == 0 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1-1][i2+1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1]) / 4
						} else if i2 == screenWidth-1 {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2-1] + g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2]) / 4
						} else {
							g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2-1] + g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1-1][i2+1] + g.pixSlices[j-1][i1][i2-1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1]) / 6
						}

					} else if i2 == 0 {
						g.pixSlices[j][i1][i2] = (g.pixSlices[j-1][i1-1][i2] + g.pixSlices[j-1][i1-1][i2+1] + g.pixSlices[j-1][i1][i2] + g.pixSlices[j-1][i1][i2+1] + g.pixSlices[j-1][i1+1][i2] + g.pixSlices[j-1][i1+1][i2+1]) / 6
					} else if i2 == screenWidth-1 {
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
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.level < 49 {
			g.level++
			err := g.UpdatePixels(g.level)
			if err != nil {
				return err
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
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
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%d %d\n%s\nLevel: %d\nSpace - new map\nArrows - change level", ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.x, g.y, txt, g.level)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RandomMaps")
	ebiten.SetMaxTPS(60)
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
