package main

import "github.com/hajimehoshi/ebiten/v2"

func NewUI(wes *Wes, counter *Counter) *UI {
	return &UI{}
}

type UI struct {
}

func (u *UI) Draw(screen *ebiten.Image) {
}
