package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var loadingLog []string

type LoadingScreen struct {
	Screen
}

// Process drawing of loading screen in drawing goroutine
func (s *LoadingScreen) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(splash.Bounds().Dx())/2, -float64(splash.Bounds().Dy())/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(splash, op)

	s.game.systemFontRenderer.DrawTextFormattedAt(screen, "Loading...", TextFormat{textColor: color.White, scale: 2.0, shadow: true}, 10, 10)

	for i, logString := range loadingLog {
		s.game.systemFontRenderer.DrawTextFormattedAt(screen, logString, TextFormat{textColor: color.White, scale: 1.0, shadow: true}, 10, 32+float64(i*16))
	}
}

func (s *LoadingScreen) ProcessKeyEvents() bool {
	return true
}

func (s *LoadingScreen) Update() {
	if s.game.ready {
		s.game.SetScreen(CreateMainMenu(s.game))
	}
}

func NewLoadingScreen(g *Game) *LoadingScreen {
	screen := new(LoadingScreen)
	screen.game = g

	return screen
}
