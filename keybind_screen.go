package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeybindSettingsScreen struct {
	GenericWidgetContainerScreen
	gameplayScreen *GameplayScreen
}

func (s *KeybindSettingsScreen) ProcessKeyEvents() bool {
	if s.GenericWidgetContainerScreen.ProcessKeyEvents() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.gameplayScreen.overlayStack.Pop()
		}
	}

	return false
}

func NewKeybindSettingsScreen(gs *GameplayScreen) *KeybindSettingsScreen {
	s := new(KeybindSettingsScreen)
	s.gameplayScreen = gs
	s.game = gs.game

	s.widgets = append(s.widgets, CreateCommonButton(s.game, "Move Up: W", func(g *Game) {

	}))

	s.widgets = append(s.widgets, CreateCommonButton(s.game, "Move Left: A", func(g *Game) {

	}))

	s.widgets = append(s.widgets, CreateCommonButton(s.game, "Move Down: S", func(g *Game) {

	}))

	s.widgets = append(s.widgets, CreateCommonButton(s.game, "Move Right: D", func(g *Game) {

	}))

	s.widgets = append(s.widgets, CreateCommonButton(s.game, "Toggle edit mode: F8", func(g *Game) {

	}))

	s.widgets = append(s.widgets, CreateCommonButton(s.game, "Toggle debug screen: F3", func(g *Game) {

	}))

	s.SetInitialFocus()

	return s
}
