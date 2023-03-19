package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type IXPWindow interface {
	Draw(*ebiten.Image)
	GetRect() Rect
	SetPosition(float64, float64)
	GetPosition() (float64, float64)
	SetActive(bool)
	GetHWND() int
	OnLButtonDown()
	OnLButtonUp()
	OnClose()
	OnMove()
	OnMouseMove()
	GetText() string
}

type XPWindow struct {
	xpScreen                  *WinXPScreen
	game                      *Game
	posX, posY, width, height float64
	title                     string
	active                    bool
	closeButtonHover          bool
	closeButtonDown           bool
	hWnd                      int

	captionFrameImage *ebiten.Image
	leftFrameImage    *ebiten.Image
	rightFrameImage   *ebiten.Image
	bottomFrameImage  *ebiten.Image
	//closeButtonImage  *ebiten.Image
	//closeGlyphImage   *ebiten.Image

	font *Font
}

func (wnd *XPWindow) GetText() string {
	return wnd.title
}

func (wnd *XPWindow) Draw(screen *ebiten.Image) {
	x := wnd.posX
	y := wnd.posY
	width := wnd.width
	height := wnd.height
	g := wnd.game

	captionBarState := 1
	if wnd.xpScreen.activeWindow == wnd {
		captionBarState = 0
	}

	g.DrawStatedNineGrid(screen, wnd.captionFrameImage, captionBarState, 2, 1.0, NineGridInfo{Right: 35, Top: 9, Left: 28, Bottom: 17}, x, y, width, height)
	g.DrawStatedNineGrid(screen, wnd.leftFrameImage, captionBarState, 2, 1.0, NineGridInfo{Right: 2, Top: 0, Left: 2, Bottom: 0}, x, y+30, 4, height-30-5)
	g.DrawStatedNineGrid(screen, wnd.rightFrameImage, captionBarState, 2, 1.0, NineGridInfo{Right: 2, Top: 0, Left: 2, Bottom: 0}, x+width-4, y+30, 4, height-30-5)
	g.DrawStatedNineGrid(screen, wnd.bottomFrameImage, captionBarState, 2, 1.0, NineGridInfo{Right: 5, Top: 2, Left: 5, Bottom: 2}, x, y+height-5, width, 5)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x+6, y+8)
	// screen.DrawImage(xpNotepadIcon, op)

	/*
		closeBtnState := 0
		if wnd.closeButtonDown {
			closeBtnState = 2
		} else if wnd.closeButtonHover {
			closeBtnState = 1
		}
	*/

	//g.DrawWinXPCloseButton(screen, wnd.xpScreen.activeWindow != wnd, closeBtnState, x+width-25, y+5, 21, 21)
	fontRenderer := g.fontRenderer

	fontRenderer.PushState()

	fontRenderer.Reset()
	fontRenderer.SetTextColor(color.White)
	fontRenderer.SetScale(2.0)
	fontRenderer.SetFont(wnd.font)

	textDim := fontRenderer.GetStringDimensions(wnd.title)
	fontRenderer.DrawTextAt(screen, wnd.title, Vec2f{x + 28, y + (30-textDim.Y)/2})

	fontRenderer.PopState()

	ebitenutil.DrawRect(screen, x+4, y+30, width-4-4, height-30-5, color.RGBA{236, 233, 216, 255})
}

func (wnd *XPWindow) GetHWND() int {
	return wnd.hWnd
}

func (wnd *XPWindow) GetRect() Rect {
	return Rect{int(wnd.posX), int(wnd.posY), int(wnd.posX + wnd.width), int(wnd.posY + wnd.height)}
}

func (wnd *XPWindow) SetPosition(x float64, y float64) {
	wnd.posX = x
	wnd.posY = y
}

func (wnd *XPWindow) GetPosition() (x float64, y float64) {
	return wnd.posX, wnd.posY
}

func (wnd *XPWindow) OnMove() {
}

func (wnd *XPWindow) SetActive(status bool) {
	wnd.active = status
}

func (wnd *XPWindow) OnMouseMove() {
	curX, curY := ebiten.CursorPosition()
	ncPosX := curX - int(wnd.posX)
	ncPosY := curY - int(wnd.posY)

	if ncPosX >= int(wnd.width)-5-21 && ncPosY >= 5 &&
		ncPosX < int(wnd.width)-5 && ncPosY < 30-5 {
		wnd.closeButtonHover = true
		return
	} else {
		wnd.closeButtonHover = false
	}
}

func (wnd *XPWindow) OnLButtonDown() {
	curX, curY := ebiten.CursorPosition()
	ncPosX := curX - int(wnd.posX)
	ncPosY := curY - int(wnd.posY)

	if ncPosX >= int(wnd.width)-5-21 && ncPosY >= 5 &&
		ncPosX < int(wnd.width)-5 && ncPosY < 30-5 {
		wnd.closeButtonDown = true

		return
	}

	if ncPosX >= 0 && ncPosY >= 0 &&
		ncPosX <= int(wnd.width) && ncPosY <= 30 {
		windowGrab = wnd
		wnd.OnMove()
	}
}

func (wnd *XPWindow) OnLButtonUp() {
	if wnd.closeButtonDown {
		wnd.OnClose()
		wnd.closeButtonDown = false
	}
}

func (wnd *XPWindow) OnClose() {
	wnd.DestroyWindow()
}

func (wnd *XPWindow) DestroyWindow() {
	wnd.OnDestroy()

	a := &wnd.xpScreen.windows

	for i := len(*a) - 1; i >= 0; i-- {
		if (*a)[i].GetHWND() == wnd.hWnd {
			*a = append((*a)[:i], (*a)[i+1:]...)
			break
		}
	}
}

func (wnd *XPWindow) OnDestroy() {

}

func InitXPWindow(wnd *XPWindow, xps *WinXPScreen, title string, x float64, y float64, width float64, height float64) {
	resMan := ResourceManager_GetInstance()

	wnd.captionFrameImage = resMan.LoadImage("assets/computer/frame_caption.png")
	wnd.leftFrameImage = resMan.LoadImage("assets/computer/frame_left.png")
	wnd.rightFrameImage = resMan.LoadImage("assets/computer/frame_right.png")
	wnd.bottomFrameImage = resMan.LoadImage("assets/computer/frame_bottom.png")

	wnd.font = resMan.LoadFontJSON("font/font_fantasy.json")

	wnd.hWnd = xps.hwndCounter
	xps.hwndCounter++
	wnd.xpScreen = xps
	wnd.game = xps.game
	wnd.posX = x
	wnd.posY = y
	if x == cwUseDefault {
		wnd.posX = float64(xps.cwDefX)
		wnd.posY = float64(xps.cwDefY)

		if xps.cwDefY+int(height) < screenHeight {
			xps.cwDefY += 10
		}

		xps.cwDefX += 10

		if xps.cwDefX >= screenWidth {
			xps.cwDefX = 0
			xps.cwDefY = 0
		}

	} else {
		wnd.posX = x
		wnd.posY = y
	}
	wnd.width = width
	wnd.height = height
	wnd.title = title

	xps.SetActiveWindow(wnd)
}

func NewXPWindow(xps *WinXPScreen, title string, x float64, y float64, width float64, height float64) *XPWindow {
	s := new(XPWindow)
	InitXPWindow(s, xps, title, x, y, width, height)
	return s
}
