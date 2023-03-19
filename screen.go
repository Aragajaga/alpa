package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type IScreen interface {
	Draw(*ebiten.Image)
	Update()
	ProcessKeyEvents() bool
	LoadResources()
	UnloadResources()
	OnAttach()
	OnDetach()
}

type Screen struct {
	IScreen
	game *Game
}

func (s *Screen) LoadResources() {

}

func (s *Screen) UnloadResources() {

}

func (s *Screen) OnAttach() {
	log.Printf("screen %T attached", s.IScreen)
	s.IScreen.LoadResources()
}

func (s *Screen) OnDetach() {
	log.Printf("screen %T detached", s.IScreen)
	s.IScreen.UnloadResources()
}

func (s *Screen) Draw(*ebiten.Image) {

}

func (s *Screen) Update() {

}
