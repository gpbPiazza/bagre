package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth, screenHeight = 1150, 550
)

var (
	darkGrey = color.RGBA{R: 40, G: 45, B: 60, A: 255}
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

	for _, jelly := range units {
		drawUnit(screen, jelly, g.tick)
	}

	drawUnit(screen, wes, g.tick)
}

type Unit interface {
	// Draw return every property needed to propertly draw a unit
	// Draw itself dont draw the unit. just return data
	Draw() (img *ebiten.Image, tickCountPerPose int, frameCount int)

	Scale() (float64, float64)

	Position() (float64, float64)

	VecVelocity() Vector2D
	VecPosition() Vector2D
	ID() int
}

func drawUnit(screen *ebiten.Image, unit Unit, tick int) {
	frameOX := 0 // frames start at the left edge
	frameOY := 0 // ...and at the top edge (single row)

	img, tickCountPerPose, frameCount := unit.Draw()

	frameWidth, frameHeight := calcFrame(img, frameCount)

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2) // center pin
	op.GeoM.Scale(unit.Scale())                                        // shrink or grow the sprite around that pin
	op.GeoM.Translate(unit.Position())                                 // move

	// Pick which frame to show, based on the clock.
	// tick/5  -> hold each pose for 5 ticks (~12fps instead of 60)
	// % frameCount -> loop back to frame 0 after the last frame (0..7)
	i := (tick / tickCountPerPose) % frameCount

	sx, sy := frameOX+i*frameWidth, frameOY

	// SubImage returns a cropped VIEW into the sheet — the exact pixel rectangle
	// (sx,sy)..(sx+48,sy+48), i.e. one 48x48 frame. No pixels are copied; it's a window.
	//
	// The sheet is one row: [frame0, frame1, frame2, frame3]
	// Each Draw call crops just ONE frame (the one `i` points to right now).
	// As g.tick advances over successive Draw calls, `i` steps 0→1→2→3→0…,
	// so the sequence of stills played over time reads as animation.
	cropRect := image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)
	walkFrame := img.SubImage(cropRect).(*ebiten.Image)

	screen.DrawImage(walkFrame, op)
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
