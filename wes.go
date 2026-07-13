package main

import (
	"bytes"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	attackWesImg *ebiten.Image
	walkWesImg   *ebiten.Image
)

// conforme o wes for comendo os peixes ele fica maior
// seu tamanho determina se  o inimigo do wes te medo dele ou não
// inimigo do wes quer matar ele. Wes quer comer maximo de peixer possíveis
// Ai nivel de dificuldade pro player
// - o cardume vai ter sprint com CD por tempo
// - ⁠se o predador ficar envolto pot muito jelly fish, eles vai soltar choque
// - ⁠no meio do jogo vai vir o inimigo do wes, um polvo vermelho puto, querendo matar o wes
// - ⁠o wes toda vez que come um  peixe ele cresce
// - ⁠o wes ficar mt grande faz o inimigo do wes ter medo dele, logo, inimigo do wes nao mata ele
// - ⁠wes ganha o jogo de comer todos os peixes
// - ⁠wes perde de ele morrer
// - ⁠wes perde se acabar o tempo
// - ⁠antes do inimigo do wes aparecer, vai ter tartarugas, elas nao dao dano no wes mas elas comem os peixes dele, logo, ele fica menor. Ele compete com elas e elas empurram ele tomando stun

const (
	walk = iota
	attack
)

type Wes struct {
	position Vector2D
	velocity Vector2D
	id       int

	state int
}

func NewWes(id int) *Wes {
	w := &Wes{
		position: Vector2D{x: screenWidth / 2, y: screenHeight / 2},
		velocity: Vector2D{x: 1.0, y: 1.0},
		id:       id,
		state:    walk,
	}

	return w
}

func (w *Wes) Draw() (img *ebiten.Image, tickCountPerPose int, frameCount int) {
	switch w.state {
	case attack:
		return attackWesImg, 3, 6
	case walk:
		return walkWesImg, 5, 4
	default:
		return walkWesImg, 5, 4
	}
}

func (w *Wes) Position() (float64, float64) {
	return w.position.x, w.position.y
}

func (w *Wes) move() {
	const speed = 2.0

	newPosition := w.position
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newPosition.y = w.position.SubY(speed).y
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newPosition.y = w.position.AddY(speed).y
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newPosition.x = w.position.SubX(speed).x
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newPosition.x = w.position.AddX(speed).x
	}

	// 10 px visually good
	if newPosition.x > screenWidth-10 {
		newPosition.x = screenWidth - 10
	}
	if newPosition.x < 10 {
		newPosition.x = 10
	}
	if newPosition.y > screenHeight-10 {
		newPosition.y = screenHeight - 10
	}
	if newPosition.y < 10 {
		newPosition.y = 10
	}

	log.Println("WES position!", newPosition)
	w.position = newPosition

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		w.state = walk
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		w.state = attack
	}
}

func loadWesImg() error {
	walk, err := os.ReadFile("./assets/wes/Walk.png")
	if err != nil {
		return err
	}

	walkImg, _, err := image.Decode(bytes.NewReader(walk))
	if err != nil {
		return err
	}

	walkWesImg = ebiten.NewImageFromImage(walkImg)

	attack, err := os.ReadFile("./assets/wes/Attack.png")
	if err != nil {
		return err
	}

	attackImg, _, err := image.Decode(bytes.NewReader(attack))
	if err != nil {
		return err
	}

	attackWesImg = ebiten.NewImageFromImage(attackImg)

	return nil
}

func DrawUnit(screen *ebiten.Image, unit Unit, tick int) {
	img, tickCountPerPose, frameCount := unit.Draw()

	frameWidth, frameHeight := calcFrame(img, frameCount)

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2) // center pin
	op.GeoM.Translate(wes.position.x, wes.position.y)                  // move

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
