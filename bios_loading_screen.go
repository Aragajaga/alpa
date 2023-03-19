package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type BIOSLoadingScreen struct {
	Screen
	startTime     int
	parentManager *ScreenManager
	strings       []string
}

func (bs *BIOSLoadingScreen) ProcessKeyEvents() bool {
	return true
}

func (bs *BIOSLoadingScreen) Update() {
	if bs.game.appTicker-bs.startTime >= 120 {
		bs.parentManager.SetScreen(NewXPBootScreen(bs.game, bs.parentManager))
	}
}

func (bs *BIOSLoadingScreen) Draw(screen *ebiten.Image) {

	duration := bs.game.appTicker - bs.startTime

	pos := Vec2f{16.0, 16.0}

	fontRenderer := bs.game.fontRenderer
	fontRenderer.PushState()

	fontRenderer.Reset()
	fontRenderer.SetScale(2.0)
	fontRenderer.SetTextColor(color.White)

	for i := 0; i < int(float64(len(bs.strings))*math.Min(float64(duration)/60.0, 1.0)); i++ {
		bs.game.fontRenderer.DrawTextAt(screen, bs.strings[i], pos)
		pos.Translate(Vec2f{0, float64(i) * fontRenderer.GetGlyphSize().Y})
	}

	fontRenderer.PopState()
}

func NewBIOSLoadingScreen(g *Game, parentManager *ScreenManager) *BIOSLoadingScreen {
	screen := new(BIOSLoadingScreen)
	screen.game = g
	screen.parentManager = parentManager
	screen.startTime = g.appTicker

	screen.strings = []string{
		"Adward Modular BIOS v05312",
		"(c) Adward Inc. 1998-2001",
		"",
		"Main Processor: Intel Core i9-12900K",
		"RDRAM Clock: 5000MHz",
		"Memory Test: 2625400K OK",
		"",
		"PNP Init Completed",
		"",
		"IDE Slot 0: STA241243622124",
		"IDE Slot 1: Optiarc MB62134",
	}

	return screen
}
