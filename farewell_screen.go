package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// A goodbye screen which appears when player clicked on an exit button
//
// This will play animation and save/closes the game in background
type FarewellScreen struct {
	Screen
	startTime int
}

func (s *FarewellScreen) Update() {
	if s.game.appTicker-s.startTime > ebiten.MaxTPS()*3 {
		os.Exit(0)
	}
}

func (s *FarewellScreen) Draw(screen *ebiten.Image) {
	text := I18n("string_seeya", "See you again.")

	fontRenderer := s.game.fontRenderer
	fontRenderer.PushState()
	fontRenderer.SetScale(2.0)

	textDim := fontRenderer.GetStringDimensions(text)
	pos := Vec2f{screenWidth, screenHeight}
	pos = pos.Subtract(textDim)
	pos = pos.Scale(0.5)

	s.game.fontRenderer.DrawTextAt(screen, text, pos)

	fontRenderer.PopState()

	op := &ebiten.DrawImageOptions{}

	pos = Vec2f{screenWidth, screenHeight}.Scale(0.5)
	pos = pos.Add(Vec2f{0, textDim.Y})
	tile := GetTileSprite(seeYaTileSet, 4, 16, Tile((s.game.appTicker/8)%3))
	op.GeoM.Translate(-(float64(tile.Bounds().Dx()) / 2), -(float64(tile.Bounds().Dy()) / 2))
	op.GeoM.Scale(4.0, 4.0)
	op.GeoM.Translate(pos.X, pos.Y+100)

	screen.DrawImage(tile, op)
}

func (*FarewellScreen) ProcessKeyEvents() bool {
	return true
}

func NewFarewellScreen(g *Game) *FarewellScreen {
	s := new(FarewellScreen)
	s.IScreen = s
	s.game = g
	s.startTime = g.appTicker

	return s
}
