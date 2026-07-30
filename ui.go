package main

import "github.com/hajimehoshi/ebiten/v2"

type UI struct {
}

func NewUI(_ *Wes, _ *Counter) *UI {
	return &UI{}
}

func (u *UI) Draw(_ *ebiten.Image) {
}
