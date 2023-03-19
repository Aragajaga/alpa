package main

import "github.com/hajimehoshi/ebiten/v2"

type XPMessageBox struct {
	XPWindow
	message string
}

func XPDrawButton(g *Game, screen *ebiten.Image, text string, state int, x float64, y float64, width float64, height float64) {
	g.DrawStatedNineGrid(screen, xpButton, state, 5, 1.0, NineGridInfo{Top: 9, Left: 8, Right: 8, Bottom: 9},
		x, y, width, height)

	fontRenderer := g.fontRenderer

	fontRenderer.PushState()
	fontRenderer.Reset()

	textDim := g.fontRenderer.GetStringDimensions(text)
	g.fontRenderer.DrawTextAt(screen, text, Vec2f{x + (width-textDim.X)/2, y + (height-textDim.Y)/2})

	fontRenderer.PopState()
}

func (wnd *XPMessageBox) Draw(screen *ebiten.Image) {
	wnd.XPWindow.Draw(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(wnd.posX+4+12, wnd.posY+4+12+30)
	screen.DrawImage(xpIconError, op)

	XPDrawButton(wnd.game, screen, "Кратко", 4, wnd.posX+4+11, wnd.posY+90, 75, 23)
	XPDrawButton(wnd.game, screen, "ДА", 0, wnd.posX+96, wnd.posY+90, 75, 23)
	XPDrawButton(wnd.game, screen, "НЕТ", 0, wnd.posX+177, wnd.posY+90, 75, 23)

	fontRenderer := wnd.game.fontRenderer

	fontRenderer.PushState()
	fontRenderer.Reset()

	textDim := wnd.game.fontRenderer.GetStringDimensions(wnd.message)
	wnd.game.fontRenderer.DrawTextAt(screen, wnd.message,
		Vec2f{wnd.posX + 4 + 12 + 20 + 32, wnd.posY + 46 + (32-textDim.Y)/2})

	fontRenderer.PopState()

}

func NewXPMessageBox(xps *WinXPScreen, message string, caption string, icon int) *XPMessageBox {
	s := new(XPMessageBox)
	InitXPWindow(&s.XPWindow, xps, caption, cwUseDefault, cwUseDefault, 266, 126)
	s.message = message

	rm := ResourceManager_GetInstance()
	rm.LoadSound("sound/computer.ogg").Play()

	return s
}
