package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SettingsScreen struct {
	GenericWidgetContainerScreen
	gameplayScreen *GameplayScreen
}

func (s *SettingsScreen) ProcessKeyEvents() bool {
	if s.GenericWidgetContainerScreen.ProcessKeyEvents() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.gameplayScreen.overlayStack.Pop()
		}
	}

	return false
}

func CreateSettingsScreen(gs *GameplayScreen) *SettingsScreen {
	s := new(SettingsScreen)
	InititalizeGenericWidgetContainerScreen(&s.GenericWidgetContainerScreen)
	s.gameplayScreen = gs
	s.game = gs.game
	s.title = "Settings"

	s.widgets = append(s.widgets, CreateCommonButton(s.game, fmt.Sprintf("%s (Language)", I18n("string_language", "Language")), func(g *Game) {
		gs.overlayStack.Push(NewKeybindSettingsScreen(gs))
	}))

	s.widgets = append(s.widgets, CreateCommonButton(s.game, I18n("string_keybinds", "Keybinds"), func(g *Game) {
		gs.overlayStack.Push(NewKeybindSettingsScreen(gs))
	}))

	s.widgets = append(s.widgets, CreateCommonButton(s.game, I18n("string_music", "Music"), func(g *Game) {

	}))

	s.widgets = append(s.widgets, CreateCommonSliderWidget(s.game, func(g *Game) {

	}))

	s.SetInitialFocus()

	return s
}
