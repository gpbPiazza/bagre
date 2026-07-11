package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	jellyBytesPng, err := os.ReadFile("./assets/jellyfish/Walk.png")
	if err != nil {
		log.Fatal(err)
	}

	jellyImg, imgName, err := image.Decode(bytes.NewReader(jellyBytesPng))
	log.Default().Printf("open img %s", imgName)
	if err != nil {
		log.Fatal(err)
	}
	jellyRunner = ebiten.NewImageFromImage(jellyImg)

	// setEmptySmacksMapsPosition()
	// NewSmack()
	// setEachJellyFishPositionIntoSmackMap()

	s := NewGame()

	ebiten.SetWindowSize(s.ScreenWidth, s.ScreenWidth)

	ebiten.SetWindowTitle("Wes bagre")

	if err := ebiten.RunGame(s); err != nil {
		log.Fatal(err)
	}
}
