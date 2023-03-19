package main

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameplayScreen struct {
	Screen
	gameplayMode IGameplayMode
	overlayStack ScreenStack
}

func (s *GameplayScreen) SetGameplayMode(mode IGameplayMode) {
	s.gameplayMode = mode
}

func (s *GameplayScreen) Draw(screen *ebiten.Image) {
	game := s.game

	game.DrawWorld(screen)

	var entities []ILivingEntity

	entityListLocker := game.entityListMutex.RLocker()
	entityListLocker.Lock()
	entities = append(entities, game.entities...)
	entityListLocker.Unlock()

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].GetWorldPos().Y < entities[j].GetWorldPos().Y
	})

	for _, entity := range entities {
		entity.Draw(screen)
		for _, spell := range entity.GetLivingEntity().spells {
			spell.Draw(screen)
		}
	}

	if game.showDebugInfo {
		for _, entity := range entities {
			entity.DrawDebugInfo(screen)
		}
	}

	s.gameplayMode.Draw(screen)

	s.overlayStack.Draw(screen)

	if game.showDebugInfo {
		game.debugScreen.Draw(screen)
	}
}

func (s *GameplayScreen) ProcessKeyEvents() bool {
	game := s.game

	if s.gameplayMode.ProcessKeyEvents() {

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.overlayStack.Push(CreateGameMenu(s))
		}

		if inpututil.IsKeyJustPressed(keyBinds[kbToggleEditMode]) {
			s.SetGameplayMode(NewGameplayModeEdit(s))
		}

		if inpututil.IsKeyJustPressed(keyBinds[kbToggleEntityFocusRotation]) {
			s.SetGameplayMode(NewGameplayModeEntityFocusRotation(s))
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
			game.ToggleDebugInfoShow()
		}
	}

	return true
}

func (s *GameplayScreen) Update() {

	keyOpaque := true
	for i := s.overlayStack.head; i != nil; i = i.next {
		keyOpaque = keyOpaque && i.screen.ProcessKeyEvents()
	}

	if keyOpaque {
		s.ProcessKeyEvents()
	}

	for i := s.overlayStack.head; i != nil; i = i.next {
		i.screen.Update()
	}

	for _, entity := range s.game.entities {
		entity.Update()
	}

	a := &s.game.entities
	for i := len(*a) - 1; i >= 0; i-- {
		if (*a)[i].GetHealth() <= 0 {
			s.game.entityListMutex.Lock()

			(*a)[i] = (*a)[len(*a)-1]
			(*a)[len(*a)-1] = nil
			*a = (*a)[:len(*a)-1]
			s.game.entityListMutex.Unlock()
		}
	}

	s.game.camera.Update()

	tickCounter++
}

func NewGameplayScreen(g *Game) *GameplayScreen {
	s := new(GameplayScreen)
	s.IScreen = s
	s.game = g

	s.SetGameplayMode(NewGameplayModeDefault(s))

	/*
		g.audioManager.PlayBackgroundMusic("bgm/level0")
	*/

	return s
}
