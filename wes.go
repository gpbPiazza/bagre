package main

import (
	"bytes"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

var (
	keyboard chan (ebiten.Key)
)

type Wes struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func NewWes(id int) *Wes {
	w := &Wes{
		position: Vector2D{x: screenWidth / 2, y: screenHeight / 2},
		velocity: Vector2D{x: 1.0, y: 1.0},
		id:       id,
	}

	return w
}

func (w *Wes) swim() {
	const speed = 2.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		wes.position = wes.position.SubY(speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		wes.position = wes.position.AddY(speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		wes.position = wes.position.SubX(speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		wes.position = wes.position.AddX(speed)
	}

	// discrete action — still JustPressed, fires once per tap
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		log.Println("WES attaque agora!! ROAR!")
	}
}

func loadWesImg() error {
	wesBytesPng, err := os.ReadFile("./assets/wes/Walk.png")
	if err != nil {
		return err
	}

	wesImg, imgName, err := image.Decode(bytes.NewReader(wesBytesPng))
	log.Default().Printf("open img %s", imgName)
	if err != nil {
		return err
	}

	wesRunner = ebiten.NewImageFromImage(wesImg)

	return nil
}

func DrawWesWalk(screen, img *ebiten.Image, tick, frameCount int, wes *Wes) {
	frameWidth, frameHeight := calcFrame(img, frameCount)

	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2) // center pin
	op.GeoM.Translate(wes.position.x, wes.position.y)                  // move

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
