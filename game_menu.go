package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameMenu struct {
	GenericWidgetContainerScreen
	gameplayScreen *GameplayScreen
}

func (s *GameMenu) ProcessKeyEvents() bool {
	if s.GenericWidgetContainerScreen.ProcessKeyEvents() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.gameplayScreen.overlayStack.Pop()
		}
	}

	return false
}
func CreateGameMenu(s *GameplayScreen) *GameMenu {
	gm := new(GameMenu)
	gm.gameplayScreen = s

	g := s.game

	gm.game = g
	gm.title = "Paused"

	gm.widgets = append(gm.widgets, CreateCommonButton(g, I18n("string_resume_game", "Resume Game"), func(g *Game) {
		s.overlayStack.Pop()
	}))

	gm.widgets = append(gm.widgets, CreateCommonButton(g, I18n("string_save_level", "Save Level"), func(g *Game) {
		s.overlayStack.Push(CreateSaveLevelScreen(s))
	}))

	gm.widgets = append(gm.widgets, CreateCommonButton(g, I18n("string_settings", "Settings"), func(g *Game) {
		s.overlayStack.Push(CreateSettingsScreen(s))
	}))

	gm.widgets = append(gm.widgets, CreateCommonButton(g, I18n("string_main_menu", "Main Menu"), func(g *Game) {
		s.game.SetScreen(CreateMainMenu(s.game))
	}))

	gm.widgets[0].SetSelection(true)
	gm.focusedWidget = gm.widgets[0]
	return gm
}
