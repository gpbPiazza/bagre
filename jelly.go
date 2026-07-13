package main

import (
	"bytes"
	"image"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// How far a jelly senses neighbors. Keep it LOCAL: too big and every jelly
	// sees the whole population, so cohesion's "average position" becomes the
	// screen center and all of them collapse into one static blob. Small radius
	// = many sub-groups that drift, split, and re-form (a living school).
	// Lower bound: must exceed the drawn sprite size (48*0.5 = 24px) so a jelly
	// senses a neighbor BEFORE their images overlap.
	jellyViewRadius      = 20
	adjustRateAligment   = 0.15
	adjustRateCohesion   = 0.09
	adjustRateSeparation = 0.30
	jellysCount          = 250
)

var (
	smackMapPositions = [screenWidth + 1][screenHeight + 1]int{}
	// smack is the collective noun for jellyFish
	smack    = make(map[int]*JellyFish, 0)
	wes      *Wes
	rwLocker = sync.RWMutex{}
)

type JellyFish struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func NewJellyFish(id int) *JellyFish {
	b := &JellyFish{
		position: Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeight},
		velocity: Vector2D{x: (rand.Float64() * 2) - 1.0, y: (rand.Float64() * 2) - 1.0},
		id:       id,
	}

	go b.swim()

	return b
}

func (b *JellyFish) swim() {
	for {
		b.move()
		time.Sleep(5 * time.Millisecond)
	}
}

// The behavior you actually want
//
// "Come together, keep moving, separate, re-form" = a living school. That requires:
//
// - Constant motion — fish always cruise; forces change direction, not whether they move.
// - Local awareness — each fish reacts to a few near neighbors, not the whole tank.
// - Tension, not equilibrium — separation slightly wins up close (no overlap), cohesion pulls loosely at medium range, alignment keeps sub-groups coherent — but they never perfectly cancel.
func (b *JellyFish) calcAcceleration() Vector2D {
	upperView := b.position.AddVal(jellyViewRadius)
	lowerView := b.position.AddVal(-jellyViewRadius)
	// all variables with prefix all here mean all elements inside of viewBox, inside of JellyFish View Radius
	allJellyFishsVelocity := Vector2D{x: 0, y: 0}
	allJellyFishsPosition := Vector2D{x: 0, y: 0}
	allJellyFishsSeparation := Vector2D{x: 0, y: 0}
	jellysCount := 0.0

	rwLocker.RLock()
	for i := math.Max(lowerView.x, 0); i <= math.Min(upperView.x, screenWidth); i++ {
		for j := math.Max(lowerView.y, 0); j <= math.Min(upperView.y, screenHeight); j++ {
			seenJellyFishID := smackMapPositions[int(i)][int(j)]
			if seenJellyFishID != -1 && b.id != seenJellyFishID {
				seenJellyFish := smack[seenJellyFishID]
				dist := seenJellyFish.position.Distance(b.position)
				if dist < jellyViewRadius {
					jellysCount++
					allJellyFishsVelocity = allJellyFishsVelocity.Add(seenJellyFish.velocity)   // Direction
					allJellyFishsPosition = allJellyFishsPosition.Add(seenJellyFish.position)   // move to the center
					separation := b.position.Subtract(seenJellyFish.position).DivisionVal(dist) // push aways too close
					allJellyFishsSeparation = allJellyFishsSeparation.Add(separation)           // push aways too close
				}
			}
		}
	}
	rwLocker.RUnlock()

	borderBounceX := b.borderBounce(b.position.x, screenWidth)
	borderBouncey := b.borderBounce(b.position.y, screenHeight)
	accel := Vector2D{x: borderBounceX, y: borderBouncey}

	if jellysCount > 0 {
		avgVelocity := allJellyFishsVelocity.DivisionVal(jellysCount)
		avgPosition := allJellyFishsPosition.DivisionVal(jellysCount)
		accelAligment := avgVelocity.Subtract(b.velocity).MultiplyVal(adjustRateAligment)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyVal(adjustRateCohesion)
		accelSepartion := allJellyFishsSeparation.MultiplyVal(adjustRateSeparation)
		accel = accel.Add(accelAligment).Add(accelCohesion).Add(accelSepartion)
	}

	return accel
}

// Quanto mais próximo da borda mais rápido será o bounce
func (b *JellyFish) borderBounce(pos, border float64) float64 {

	// Está próximo da bater na borda, passou do limite de vistualização
	// ou seja o passarinho viu a parede e irá mudar de direção
	if pos < jellyViewRadius {
		return 1 / pos
	}
	// Is the same thing but in the other side of the screenView
	// o primeiro If é para o jelly que está próximo da parede em que X é muito pequeno
	// Aqui o x é grande, o mesmo para Y.
	if pos > border-jellyViewRadius {
		return 1 / (pos - border)
	}

	return 0
}

func (b *JellyFish) move() {
	minSpeed := 0.5
	maxSpeed := 1.5

	accel := b.calcAcceleration()

	rwLocker.Lock()
	b.velocity = b.velocity.Add(accel)

	// Cruising speed: only change the LENGTH of velocity, never its direction.
	// Floor (minSpeed) = never stall into a frozen blob.
	// Ceiling (maxSpeed) = never blast off across the screen.
	velocityMag := b.velocityMagnitude()
	if velocityMag < minSpeed {
		b.velocity = b.velocity.ScaleToLength(minSpeed)
	}
	if velocityMag > maxSpeed {
		b.velocity = b.velocity.ScaleToLength(maxSpeed)
	}

	//set the current position to -1, empty space
	smackMapPositions[int(b.position.x)][int(b.position.y)] = -1
	// move
	b.position = b.position.Add(b.velocity)
	// fill the new position into the map
	smackMapPositions[int(b.position.x)][int(b.position.y)] = b.id

	rwLocker.Unlock()
}

func NewSmack() {

	for i := 0; i < jellysCount; i++ {
		jelly := NewJellyFish(i)
		smack[jelly.id] = jelly
	}
}

func setEmptySmacksMapsPosition() {
	for i, row := range smackMapPositions {
		for j := range row {
			smackMapPositions[i][j] = -1
		}
	}
}

func setEachJellyFishPositionIntoSmackMap() {
	for _, jelly := range smack {
		smackMapPositions[int(jelly.position.x)][int(jelly.position.y)] = jelly.id
	}
}

func loadJellyImg() error {
	jellyBytesPng, err := os.ReadFile("./assets/jellyfish/Walk.png")
	if err != nil {
		return err
	}

	jellyImg, imgName, err := image.Decode(bytes.NewReader(jellyBytesPng))
	log.Default().Printf("open img %s", imgName)
	if err != nil {
		return err
	}

	jellyRunner = ebiten.NewImageFromImage(jellyImg)

	return nil
}

func DrawJellyWalk(screen, img *ebiten.Image, tick, frameCount int, jelly *JellyFish) {
	frameWidth, frameHeight := calcFrame(img, frameCount)

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2) // center pin
	op.GeoM.Scale(0.5, 0.5)                                            // shrink the sprite around that pin
	op.GeoM.Translate(jelly.position.x, jelly.position.y)              // move

	// Pick which frame to show, based on the clock.
	// tick/5  -> hold each pose for 5 ticks (~12fps instead of 60)
	// % frameCount -> loop back to frame 0 after the last frame (0..7)
	i := (tick / 5) % frameCount

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

func (b *JellyFish) velocityMagnitude() float64 {
	return b.velocity.Pythagoras()
}
