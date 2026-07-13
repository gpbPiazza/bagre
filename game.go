package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth, screenHeight = 640, 360

	frameOX        = 0 // frames start at the left edge
	frameOY        = 0 // ...and at the top edge (single row)
	walkFrameCount = 4 // Walk.png has 4 poses
)

var (
	darkGrey    = color.RGBA{R: 40, G: 45, B: 60, A: 255}
	jellyRunner *ebiten.Image
)

type Game struct {
	ScreenWidth  int
	ScreenHeight int
	tick         int // grows by 1 every Update (~60/sec); our clock for animation
}

func NewGame() *Game {
	return &Game{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
}

func (g *Game) Update() error {
	g.tick++

	wes.move()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(darkGrey)
	for _, jelly := range smack {
		DrawJellyWalk(screen, jellyRunner, g.tick, walkFrameCount, jelly)
	}

	DrawUnit(screen, wes, g.tick)
}

type Unit interface {
	// Draw return every property needed to propertly draw a unit
	// Draw itself dont draw the unit. just return data
	Draw() (img *ebiten.Image, tickCountPerPose int, frameCount int)
	Position() (float64, float64)
}

func calcFrame(img image.Image, frameCount int) (width, height int) {
	rec := img.Bounds()
	imgWidht := rec.Dx()
	imgHeight := rec.Dy()

	height = imgHeight
	width = imgWidht / frameCount
	return
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}
