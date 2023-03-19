package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SaveLevelScreen struct {
	GenericWidgetContainerScreen
	gameplayScreen *GameplayScreen
}

func (s *SaveLevelScreen) ProcessKeyEvents() bool {
	if s.GenericWidgetContainerScreen.ProcessKeyEvents() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.gameplayScreen.overlayStack.Pop()
		}
	}

	return false
}

func CreateSaveLevelScreen(gameplayScreen *GameplayScreen) *SaveLevelScreen {
	s := new(SaveLevelScreen)
	InititalizeGenericWidgetContainerScreen(&s.GenericWidgetContainerScreen)

	s.gameplayScreen = gameplayScreen
	s.game = gameplayScreen.game
	g := gameplayScreen.game
	s.title = langData["string_save_level"]

	editBox := CreateCommonEditBox(g, func(g *Game) {})
	editBox.SetText("level0.lvl")
	s.widgets = append(s.widgets, editBox)

	s.widgets = append(s.widgets, CreateCommonButton(g, I18n("string_verb_save", "Save"), func(g *Game) {
		g.SaveLevel(s.widgets[0].GetText())
		gameplayScreen.overlayStack.Pop()
	}))

	s.SetInitialFocus()
	return s
}
