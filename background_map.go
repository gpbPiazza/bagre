package main

import (
	"log/slog"

	"github.com/gpbPiazza/bagre/pkg/log"
	"github.com/hajimehoshi/ebiten/v2"
)

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

// type Stages struct {
// }
//
// type StageType int
//
// // const (
// // 	eletri
// // )
//
// type StageController interface {
// 	Stage() StageType
// /}

type BackgroundMap struct {
	menuImg, flatGround *ebiten.Image
	menuImgOpt          *ebiten.DrawImageOptions

	flatGroundOpts []*ebiten.DrawImageOptions
}

func NewBackgroundMap(l *slog.Logger) *BackgroundMap {
	menuImg, err := loadImage("./assets/platform/background/Background.png")
	if err != nil {
		l.Error("failed to load menuBackground", log.Err(err))
		panic(err)
	}

	flatGround, err := loadImage("./assets/platform/tileset/Tile_38.png")
	if err != nil {
		l.Error("failed to load flatGround", log.Err(err))
		panic(err)
	}

	opt := &ebiten.DrawImageOptions{}
	dx, dy := menuImg.Bounds().Dx(), menuImg.Bounds().Dy()
	opt.GeoM.Scale(screenWidth/float64(dx), screenHeight/float64(dy))

	return &BackgroundMap{
		menuImg:        menuImg,
		menuImgOpt:     opt,
		flatGround:     flatGround,
		flatGroundOpts: newGroundOpts(flatGround),
	}
}

func (b *BackgroundMap) Update() {
}

func (b *BackgroundMap) Draw(screen *ebiten.Image) {
	screen.DrawImage(b.menuImg, b.menuImgOpt)

	for _, opt := range b.flatGroundOpts {
		screen.DrawImage(b.flatGround, opt)
	}
}

func newGroundOpts(ground *ebiten.Image) []*ebiten.DrawImageOptions {
	dy := ground.Bounds().Dy()
	dx := ground.Bounds().Dx()

	groundY := float64(screenHeight - dy)

	space := screenWidth
	var opts []*ebiten.DrawImageOptions

	for space >= 0 {
		opt := &ebiten.DrawImageOptions{}

		// math.Abs(x float64)

		opt.GeoM.Translate(float64(space-dx), groundY)

		opts = append(opts, opt)

		space -= dx
	}

	return opts
}
