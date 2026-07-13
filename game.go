package main

import (
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
	green       = color.RGBA{R: 10, G: 255, B: 50}
	darkGrey    = color.RGBA{R: 40, G: 45, B: 60, A: 255}
	jellyRunner *ebiten.Image
	wesRunner   *ebiten.Image
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

	wes.swim()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(darkGrey)
	for _, jelly := range smack {
		DrawJellyWalk(screen, jellyRunner, g.tick, walkFrameCount, jelly)
	}
	DrawWesWalk(screen, wesRunner, g.tick, walkFrameCount, wes)
}

func calcFrame(img *ebiten.Image, frameCount int) (width, height int) {
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
