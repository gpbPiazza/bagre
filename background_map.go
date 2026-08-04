package main

import "github.com/hajimehoshi/ebiten/v2"

type BackgroundMap struct {
}

type Stages struct {
}

type StageType int

// const (
// 	eletri
// )

type StageController interface {
	Stage() StageType
}

// stages -> one map one final boss
// inside of each stages we will have a many events/StageEvents
// Where they will behave like a stack, each stack will have some different condition to happend
// StageController struct
// StageEvents Type
// Stages:
//
//
// background will depend on StageController where he needs to know which stage we are and when a event start
// the BackgrounEffects will be a struct responsiable to
// 1 -> make the background
// 2 -> when a stage start do something cool to user know
//
//
//
//// stage checks what ever should happend given the stage that we are
// currently we have stage 1 - eletricJellies,
//
// map 1
// stages of map 1
// Make something COOL when a stage start, like
// the entire game slow down and we get something awsome in the screen shwoing
// stage 1 started, with miniMap of stages.
//
// StageName -> Wes Vs the holdHands
// StageName -> Wes Vs the unbreakable
// turtle shield --- shark bite
// miniMap icons like:
// thunder --- octopus hand --- octopus EYE
// all this are StageEvents
//
//	0 - nothing passive jelly
//	1 - eletric jelly
//	2 - octopus hold
//	3 - octopus boss

func NewBackgroundMap(_ StageController) *BackgroundMap {
	return &BackgroundMap{}
}

func (b *BackgroundMap) Update() {

}
func (b *BackgroundMap) Draw(_ *ebiten.Image) {

}
