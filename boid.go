package main

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	jellysCount     = 600
	jellyViewRadius = 13
	adjustRate      = 0.15
)

var (
	smackMapPositions = [screenWidth + 1][screenHeight + 1]int{}
	// smack is the collective noun for jellyFish
	smack    = make(map[int]*JellyFish, 0)
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

	go b.fly()

	return b
}

func (b *JellyFish) fly() {
	for {
		b.move()
		time.Sleep(5 * time.Millisecond)
	}
}

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
			otherJellyFishId := smackMapPositions[int(i)][int(j)]
			if otherJellyFishId != -1 && b.id != otherJellyFishId {
				otherJellyFish := smack[otherJellyFishId]
				dist := otherJellyFish.position.Distance(b.position)
				if dist < jellyViewRadius {
					jellysCount++
					allJellyFishsVelocity = allJellyFishsVelocity.Add(otherJellyFish.velocity)
					allJellyFishsPosition = allJellyFishsPosition.Add(otherJellyFish.position)
					separation := b.position.Subtract(otherJellyFish.position).DivisionVal(dist)
					allJellyFishsSeparation = allJellyFishsSeparation.Add(separation)
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
		accelAligment := avgVelocity.Subtract(b.velocity).MultiplyVal(adjustRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyVal(adjustRate)
		accelSepartion := allJellyFishsSeparation.MultiplyVal(adjustRate)
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
	accel := b.calcAcceleration()

	rwLocker.Lock()
	// the limit method its to ensure the jelly will not run faster than 1 px per cycle
	b.velocity = b.velocity.Add(accel).LimitVal(-1, 1)
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
