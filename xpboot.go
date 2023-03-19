package main

import "github.com/hajimehoshi/ebiten/v2"

type XPBootScreen struct {
	Screen
	game          *Game
	startTime     int
	parentManager *ScreenManager
}

func (s *XPBootScreen) ProcessKeyEvents() bool {
	return true
}

func (s *XPBootScreen) Update() {
	if s.game.appTicker-s.startTime >= 120 {
		s.parentManager.SetScreen(NewWinXPScreen(s.game, nil))
	}
}

func (s *XPBootScreen) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	img := AssetManager_GetInstance().Get("winxp/boot_logo")
	imgWidth, imgHeight := img.Size()

	op.GeoM.Translate(-(float64(imgWidth) / 2), -(float64(imgHeight) / 2))
	op.GeoM.Translate(float64(screenWidth)/2, float64(screenHeight)/2)

	screen.DrawImage(img, op)
}

func NewXPBootScreen(g *Game, parentManager *ScreenManager) *XPBootScreen {
	screen := new(XPBootScreen)
	screen.game = g
	screen.startTime = g.appTicker
	screen.parentManager = parentManager

	return screen
}
