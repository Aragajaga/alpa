package main

import "github.com/hajimehoshi/ebiten/v2"

type RunWindow struct {
	XPWindow

	runIcon16 *ebiten.Image
	runIcon32 *ebiten.Image
}

func (wnd *RunWindow) Draw(screen *ebiten.Image) {
	x := wnd.posX + 14
	y := wnd.posY + 40

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(wnd.runIcon32, op)
}

func InitRunWindow(wnd *RunWindow, xps *WinXPScreen) {
	InitXPWindow(&wnd.XPWindow, xps, "Run", 0, 0, 320, 240)

	rm := ResourceManager_GetInstance()
	wnd.runIcon16 = rm.LoadImage("assets/computer/run16.png")
	wnd.runIcon32 = rm.LoadImage("assets/computer/run32.png")
}

func NewRunWindow(xps *WinXPScreen) *RunWindow {
	s := new(RunWindow)
	InitRunWindow(s, xps)
	return s
}
