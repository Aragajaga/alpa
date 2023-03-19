package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type WindowListWindow struct {
	XPWindow
}

func (wnd *WindowListWindow) Draw(screen *ebiten.Image) {
	wnd.XPWindow.Draw(screen)

	fontRenderer := wnd.game.fontRenderer
	fontRenderer.PushState()

	fontRenderer.SetFont(wnd.font)
	fontRenderer.SetTextColor(color.Black)

	for i, window := range wnd.xpScreen.windows {
		fontRenderer.DrawTextAt(screen, window.GetText(), Vec2f{wnd.posX, wnd.posY + float64(i*wnd.font.glyphHeight)})
	}
	fontRenderer.PopState()
}

func NewWindowListWindow(xps *WinXPScreen) *WindowListWindow {
	s := new(WindowListWindow)
	InitXPWindow(&s.XPWindow, xps, "Top-level windows", 10, 10, 240, 320)
	return s
}
