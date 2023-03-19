package main

import "github.com/hajimehoshi/ebiten/v2"

type XPTaskbarWindow struct {
	XPWindow
	windows []IXPWindow
}

/*
func (wnd *XPTaskbarWindow) Draw(screen *ebiten.Image) {

} */

func (wnd *XPTaskbarWindow) DrawTaskBand(screen *ebiten.Image, x, y, width, height float64) {
	g := wnd.game

	g.DrawNineGrid(screen, wnd.xpScreen.taskbandImage, 1.0, NineGridInfo{Top: 15, Bottom: 11, Left: 1, Right: 1}, x, y, width, height)

	// Кнопка ПУСК:
	// s.game.DrawStatedNineGrid(screen, xpStartButton, 0, 3, 1.0, NineGridInfo{Left: 6, Top: 13, Right: 52, Bottom: 14}, 0, screenHeight-30, 100, 30)
}

func (wnd *XPTaskbarWindow) Draw(screen *ebiten.Image) {
	x := wnd.posX
	y := wnd.posY
	width := wnd.width
	height := wnd.height
	// g := wnd.game

	for _, child := range wnd.windows {
		child.Draw(screen)
	}

	wnd.DrawTaskBand(screen, x, y, width, height)
}

func (wnd *XPTaskbarWindow) OnMove() {
	x, y := wnd.GetPosition()

	for _, child := range wnd.windows {
		child.SetPosition(x, y)
	}
}

func NewXPTaskbarWindow(xps *WinXPScreen) *XPTaskbarWindow {
	s := new(XPTaskbarWindow)
	InitXPWindow(&s.XPWindow, xps, "_CiceronUI", 0, 0, 100, 100)

	s.windows = append(s.windows, NewXPStartButtonWidget(xps))
	return s
}
