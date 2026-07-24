package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const counterFontSize = 24

func NewCounter(jellyFishCount int) *Counter {
	return &Counter{
		unitsEaten:          0,
		totalJellyFishCount: jellyFishCount,
	}
}

type Counter struct {
	unitsEaten          int
	totalJellyFishCount int
}

func (c *Counter) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf("%d/%d", c.unitsEaten, c.totalJellyFishCount)

	face := &text.GoTextFace{
		Source: textFont,
		Size:   counterFontSize,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(150, 10)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, msg, face, op)
}

func (c *Counter) Add(val int) {
	c.unitsEaten += val
}
