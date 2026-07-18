package main

import (
	"fmt"
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
	ScreenWidth   int
	ScreenHeight  int
	tick          int // grows by 1 every Update (~60/sec); our clock for animation
	id            int
	units         gameUnits
	evenetManager *EventManager
}

func NewGame(id int, gUnits gameUnits, evenEventManager *EventManager) *Game {
	return &Game{
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		units:         gUnits,
		evenetManager: evenEventManager,
		id:            id,
	}
}

func (g *Game) ID() int {
	return g.id
}

// Update archtecture move, we always calculate first what will happend
// then we write.
//
// KNOWN ISSUES / TODO (deferred, none block running):
//
//  1. Dead mutex leftover: rwLocker (jelly.go:34) + the sync import have zero
//     usages now that everything runs on this single goroutine. It only compiles
//     because Go tolerates unused package-level vars. Delete both to finish the
//     no-locking migration.
//
//  2. Dead jellies are never removed from g.units.smack. On Die() a jelly is
//     deleted from the units map but stays in the smack slice, so it keeps
//     running nextMove/writeMove every tick. It's invisible (not in units, so
//     not drawn or gridded), just wasted work that grows as more die. TODO:
//     prune eaten jellies out of the smack slice (and until then, skip them with
//     `if j.state == unitStateDead { continue }` in both loops below). Same gap
//     that blocks the death-animation TODO.
//
//  3. Attack ordering (not a bug): Attack() runs AFTER the jellies commit this
//     tick's move (writeMove) but reads the grid from LAST tick's rebuildGrid, so
//     wes eats based on where jellies were ~1 tick ago (~1.5px off from where
//     they're drawn). Negligible at 60fps and safe (no nil, since rebuildGrid
//     runs after Die). Side effect: Die()'s grid-clear (jelly.go:120) is now
//     redundant because rebuildGrid always runs before any reader sees the grid.
func (g *Game) Update() error {
	g.tick++

	g.units.wes.move()

	for _, j := range g.units.smack {
		j.checkState(g.tick)
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
		u.Die(g.tick)
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

func (g *Game) Subscribe(et EventType, payload any) {
	switch et {
	case removeUnit:
		u, ok := payload.(Unit)
		if !ok {
			fmt.Println("Não é uma unit no remove unit event")
			return
		}

		delete(units, u.ID())
		position := u.VecPosition()
		unitsByPositions[int(position.x)][int(position.y)] = -1

		lastIndex := len(g.units.smack) - 1
		smack := g.units.smack
		var newSmack []*JellyFish
		for i, jelly := range g.units.smack {
			if jelly.ID() == u.ID() {
				smack[i] = smack[lastIndex]
				smack[lastIndex] = smack[i]
				newSmack = smack[:lastIndex]
				break
			}
		}
		g.units.smack = newSmack

	default:
		return
	}
}

type gameUnits struct {
	wes   *Wes
	smack []*JellyFish
}

func NewUnits(eventManager *EventManager, logger *slog.Logger) gameUnits {
	wes := NewWes(jellysCount+2, logger)

	for i, row := range unitsByPositions {
		for j := range row {
			unitsByPositions[i][j] = -1
		}
	}

	var smack []*JellyFish
	for i := range jellysCount {
		jelly := newJellyFish(i, logger, eventManager)

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
	Draw() (img *ebiten.Image, ticksWhenDead int, tickCountPerPose int, frameCount int)

	Scale() (float64, float64)

	Position() (float64, float64)

	VecVelocity() Vector2D

	VecPosition() Vector2D

	ID() int

	Die(tick int)

	IsPlayer() bool
}

func drawUnit(screen *ebiten.Image, unit Unit, tick int) {
	frameOX := 0 // frames start at the left edge
	frameOY := 0 // ...and at the top edge (single row)

	img, tickesWhenDead, tickCountPerPose, frameCount := unit.Draw()

	frameWidth, frameHeight := calcFrame(img, frameCount)

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2) // center pin
	op.GeoM.Scale(unit.Scale())                                        // shrink or grow the sprite around that pin
	op.GeoM.Translate(unit.Position())                                 // move

	// Pick which frame to show, based on the clock.
	// tick/5  -> hold each pose for 5 ticks (~12fps instead of 60)
	// % frameCount -> loop back to frame 0 after the last frame (0..7)
	i := (tick / tickCountPerPose) % frameCount
	if tickesWhenDead != 0 {
		elapsed := tick - tickesWhenDead
		i = min(elapsed/tickCountPerPose, frameCount-1)
	}

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
