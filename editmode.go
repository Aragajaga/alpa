package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameplayModeEdit struct {
	GameplayMode
	editMode *EditMode
}

func (mode *GameplayModeEdit) ProcessKeyEvents() bool {
	mode.editMode.ProcessKeyEvents()

	if inpututil.IsKeyJustPressed(keyBinds[kbToggleEditMode]) {
		mode.gameplayScreen.SetGameplayMode(NewGameplayModeDefault(mode.gameplayScreen))
		return false
	}

	return true
}

func (mode *GameplayModeEdit) Draw(screen *ebiten.Image) {
	mode.editMode.Draw(screen)
}

func (mode *GameplayModeEdit) Update() {
	//mode.editMode.Update()
}

func NewGameplayModeEdit(s *GameplayScreen) *GameplayModeEdit {
	mode := new(GameplayModeEdit)
	mode.gameplayScreen = s
	mode.editMode = NewEditMode(s.game)

	return mode
}
