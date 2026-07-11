package main

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	boidsCount     = 600
	boidViewRadius = 13
	adjustRate     = 0.15
)

var (
	flockMapPositions = [screenWidth + 1][screenHeight + 1]int{}
	flock             = make(map[int]*Boid, 0)
	rwLocker          = sync.RWMutex{}
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func NewBoid(id int) *Boid {
	b := &Boid{
		position: Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeight},
		velocity: Vector2D{x: (rand.Float64() * 2) - 1.0, y: (rand.Float64() * 2) - 1.0},
		id:       id,
	}

	go b.fly()

	return b
}

func (b *Boid) fly() {
	for {
		b.move()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) calcAcceleration() Vector2D {
	upperView := b.position.AddVal(boidViewRadius)
	lowerView := b.position.AddVal(-boidViewRadius)
	// all variables with prefix all here mean all elements inside of viewBox, inside of Boid View Radius
	allBoidsVelocity := Vector2D{x: 0, y: 0}
	allBoidsPosition := Vector2D{x: 0, y: 0}
	allBoidsSeparation := Vector2D{x: 0, y: 0}
	boidsCount := 0.0

	rwLocker.RLock()
	for i := math.Max(lowerView.x, 0); i <= math.Min(upperView.x, screenWidth); i++ {
		for j := math.Max(lowerView.y, 0); j <= math.Min(upperView.y, screenHeight); j++ {
			otherBoidId := flockMapPositions[int(i)][int(j)]
			if otherBoidId != -1 && b.id != otherBoidId {
				otherBoid := flock[otherBoidId]
				dist := otherBoid.position.Distance(b.position)
				if dist < boidViewRadius {
					boidsCount++
					allBoidsVelocity = allBoidsVelocity.Add(otherBoid.velocity)
					allBoidsPosition = allBoidsPosition.Add(otherBoid.position)
					separation := b.position.Subtract(otherBoid.position).DivisionVal(dist)
					allBoidsSeparation = allBoidsSeparation.Add(separation)
				}
			}
		}
	}
	rwLocker.RUnlock()

	borderBounceX := b.borderBounce(b.position.x, screenWidth)
	borderBouncey := b.borderBounce(b.position.y, screenHeight)
	accel := Vector2D{x: borderBounceX, y: borderBouncey}

	if boidsCount > 0 {
		avgVelocity := allBoidsVelocity.DivisionVal(boidsCount)
		avgPosition := allBoidsPosition.DivisionVal(boidsCount)
		accelAligment := avgVelocity.Subtract(b.velocity).MultiplyVal(adjustRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyVal(adjustRate)
		accelSepartion := allBoidsSeparation.MultiplyVal(adjustRate)
		accel = accel.Add(accelAligment).Add(accelCohesion).Add(accelSepartion)
	}

	return accel
}

// Quanto mais próximo da borda mais rápido será o bounce
func (b *Boid) borderBounce(pos, border float64) float64 {

	// Está próximo da bater na borda, passou do limite de vistualização
	// ou seja o passarinho viu a parede e irá mudar de direção
	if pos < boidViewRadius {
		return 1 / pos
	}
	// Is the same thing but in the other side of the screenView
	// o primeiro If é para o boid que está próximo da parede em que X é muito pequeno
	// Aqui o x é grande, o mesmo para Y.
	if pos > border-boidViewRadius {
		return 1 / (pos - border)
	}

	return 0
}

func (b *Boid) move() {
	accel := b.calcAcceleration()

	rwLocker.Lock()
	// the limit method its to ensure the boid will not run faster than 1 px per cycle
	b.velocity = b.velocity.Add(accel).LimitVal(-1, 1)
	//set the current position to -1, empty space
	flockMapPositions[int(b.position.x)][int(b.position.y)] = -1
	// move
	b.position = b.position.Add(b.velocity)
	// fill the new position into the map
	flockMapPositions[int(b.position.x)][int(b.position.y)] = b.id

	rwLocker.Unlock()
}

func NewFlock() {
	for i := 0; i < boidsCount; i++ {
		boid := NewBoid(i)
		flock[boid.id] = boid
	}
}

func setEmptyFlocksMapsPosition() {
	for i, row := range flockMapPositions {
		for j := range row {
			flockMapPositions[i][j] = -1
		}
	}
}

func setEachBoidPositionIntoFlockMap() {
	for _, boid := range flock {
		flockMapPositions[int(boid.position.x)][int(boid.position.y)] = boid.id
	}
}
