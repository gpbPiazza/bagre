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
	green       = color.RGBA{R: 10, G: 255, B: 50}
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(darkGrey)

	frameWidth, frameHeight := calcFrame(jellyRunner, walkFrameCount)

	op := new(ebiten.DrawImageOptions)
	// Here we are manipulating the top left corner of the image position on the screen
	// When we perform the (-16, -16) we are setting the image pin element to his self center.
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)

	// now  here we are setting the image position to the screen
	op.GeoM.Translate(screenWidth/2, screenHeight/2)

	// Pick which frame to show, based on the clock.
	// tick/5  -> hold each pose for 5 ticks (~12fps instead of 60)
	// % frameCount -> loop back to frame 0 after the last frame (0..7)
	i := (g.tick / 5) % walkFrameCount

	// Slide the crop rectangle to frame i: only X moves, the row (Y) is fixed.
	sx, sy := frameOX+i*frameWidth, frameOY

	// SubImage returns a cropped VIEW into the sheet — the exact pixel rectangle
	// (sx,sy)..(sx+48,sy+48), i.e. one 48x48 frame. No pixels are copied; it's a window.
	//
	// The sheet is one row: [frame0, frame1, frame2, frame3]
	// Each Draw call crops just ONE frame (the one `i` points to right now).
	// As g.tick advances over successive Draw calls, `i` steps 0→1→2→3→0…,
	// so the sequence of stills played over time reads as animation.
	//
	// Cut frame i (a 32x32 window) out of the sheet, then draw it with our transform.
	cropRect := image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)
	frame := jellyRunner.SubImage(cropRect).(*ebiten.Image)
	screen.DrawImage(frame, op)

	for _, jelly := range smack {
		screen.Set(int(jelly.position.x+1), int(jelly.position.y+1), green)
		screen.Set(int(jelly.position.x-1), int(jelly.position.y-1), green)
		screen.Set(int(jelly.position.x), int(jelly.position.y+1), green)
		screen.Set(int(jelly.position.x), int(jelly.position.y-1), green)
	}
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
