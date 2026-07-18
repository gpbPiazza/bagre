package main

import (
	"fmt"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type unitState int

func (js unitState) String() string {
	switch js {
	case unitStateWalk:
		return "walk"
	case unitStateDead:
		return "dead"
	default:
		return fmt.Sprintf("not mapped state - %d", js)
	}
}

const (
	unitStateWalk unitState = iota
	unitStateAttack
	unitStateDead
)

type Unit interface {
	// Draw return every property needed to propertly draw a unit
	// Draw itself dont draw the unit. just return data
	Draw() (img *ebiten.Image, ticksWhenDead int, tickCountPerPose int, frameCount int)

	Scale() (float64, float64)

	Position() (float64, float64)

	VecVelocity() Vector2D

	VecPosition() Vector2D

	ID() int

	State() unitState

	Die(tick int)

	IsPlayer() bool
}

type Units struct {
	wes        *Wes
	smack      []*JellyFish
	unitsEaten []Unit
}

func NewUnits(eventManager *EventManager, logger *slog.Logger) Units {
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

	return Units{
		wes:        wes,
		smack:      smack,
		unitsEaten: nil,
	}
}

func checkState(u Unit, tick int, events *EventManager) {
	switch u.State() {
	case unitStateDead:
		_, ticksWhenDead, tickCountPerFrame, frameCount := u.Draw()
		howLongItLast := tickCountPerFrame * frameCount
		// 10 * 6 -> 60 -> animataçõa demora 60 ticks para terminar
		elapsed := tick - ticksWhenDead
		// elapsed for menor ou igual ao fim da animação, significa que ja acabou

		if elapsed >= howLongItLast {
			events.Publish(removeUnit, u)
		}
	default:
		return
	}
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
