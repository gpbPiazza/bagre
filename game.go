package main

import (
	"image"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth, screenHeight = 950, 550
)

var (
	darkGrey = color.RGBA{R: 40, G: 45, B: 60, A: 255}
)

type Game struct {
	ScreenWidth  int
	ScreenHeight int
	tick         int // grows by 1 every Update (~60/sec); our clock for animation
	units        gameUnits
}

func NewGame(gUnits gameUnits) *Game {
	return &Game{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		units:        gUnits,
	}
}

// Update archtecture move, we always calculate first what will happend
// then we write.
func (g *Game) Update() error {
	g.tick++

	g.units.wes.move()

	for _, j := range g.units.smack {
		j.nextMove()
	}
	for _, j := range g.units.smack {
		j.writeMove()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		g.units.wes.state = unitStateWalk
	}
	var unitsEaten []Unit
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		unitsEaten = g.units.wes.Attack()
	}
	for _, u := range unitsEaten {
		u.Die()
	}

	rebuildGrid()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(darkGrey)

	for _, u := range units {
		drawUnit(screen, u, g.tick)
	}
}

type gameUnits struct {
	wes   *Wes
	smack []*JellyFish
}

func NewUnits(logger *slog.Logger) gameUnits {
	wes := NewWes(jellysCount+2, logger)

	for i, row := range unitsByPositions {
		for j := range row {
			unitsByPositions[i][j] = -1
		}
	}

	var smack []*JellyFish
	for i := range jellysCount {
		jelly := newJellyFish(i, logger)

		units[jelly.id] = jelly
		unitsByPositions[int(jelly.position.x)][int(jelly.position.y)] = jelly.ID()
		smack = append(smack, jelly)
	}

	unitsByPositions[int(wes.position.x)][int(wes.position.y)] = wes.id
	units[wes.id] = wes

	return gameUnits{
		wes, smack,
	}
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

	Die()
	IsPlayer() bool
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

func rebuildGrid() {
	for i := range unitsByPositions {
		for k := range unitsByPositions[i] {
			unitsByPositions[i][k] = -1
		}
	}

	for id, u := range units {
		p := u.VecPosition()
		unitsByPositions[int(p.x)][int(p.y)] = id
	}
}
