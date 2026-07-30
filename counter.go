package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Counter struct {
	unitsEaten          int
	totalJellyFishCount int
	eventManager        *EventManager
	startEletricJelly   int

	alreadyStartEletricJelly bool
}

func NewCounter(jellyFishCount int, eventManager *EventManager) *Counter {
	return &Counter{
		unitsEaten:               0,
		totalJellyFishCount:      jellyFishCount,
		eventManager:             eventManager,
		startEletricJelly:        startEletricJellyAt,
		alreadyStartEletricJelly: false,
	}
}

func (c *Counter) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf("%d/%d", c.unitsEaten, c.totalJellyFishCount)

	face := &text.GoTextFace{
		Source: textFont,
		Size:   counterFontSize,
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(counterPosX, counterPosY)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, msg, face, op)
}

func (c *Counter) Add(val int) {
	c.unitsEaten += val

	if c.unitsEaten > c.startEletricJelly && !c.alreadyStartEletricJelly {
		c.eventManager.Publish(startEletricJellyFishs, nil)
		c.alreadyStartEletricJelly = true
	}
}
