package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var xpCaption *ebiten.Image
var xpLeftFrame *ebiten.Image
var xpRightFrame *ebiten.Image
var xpBottomFrame *ebiten.Image
var xpCloseButton *ebiten.Image
var xpCloseGlyph *ebiten.Image
var xpIconError *ebiten.Image
var xpButton *ebiten.Image

type WinXPScreen struct {
	Screen
	windows      []IXPWindow
	activeWindow IXPWindow
	cwDefX       int
	cwDefY       int
	hwndCounter  int

	wallpaperImage   *ebiten.Image
	taskbandImage    *ebiten.Image
	startButtonImage *ebiten.Image
	font             *Font

	nextScreenBuilder NextScreenBuilder
}

func (s *WinXPScreen) OnAttach() {
	s.Screen.OnAttach()

	s.windows = append(s.windows, NewWindowListWindow(s))
}

var grabX, grabY int
var windowGrab IXPWindow

const cwUseDefault = -2147483648

func (s *WinXPScreen) SetActiveWindow(wnd IXPWindow) {
	if windowGrab == nil {
		s.activeWindow = wnd
	}
}

func (s *WinXPScreen) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.game.SetScreen(s.nextScreenBuilder(s.game))
	}

	curX, curY := ebiten.CursorPosition()

	if windowGrab != nil {
		s.activeWindow = windowGrab
		windowGrab.SetPosition(float64(curX-grabX), float64(curY-grabY))
		windowGrab.OnMove()
	}

	for i := len(s.windows) - 1; i >= 0; i-- {
		rc := s.windows[i].GetRect()

		if curX >= rc.X1 && curY >= rc.Y1 &&
			curX < rc.X2 && curY < rc.Y2 {
			s.windows[i].OnMouseMove()
			break
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for i := len(s.windows) - 1; i >= 0; i-- {
			rc := s.windows[i].GetRect()

			if curX >= rc.X1 && curY >= rc.Y1 &&
				curX < rc.X2 && curY < rc.Y2 {
				temp := s.windows[i]
				s.windows[i] = nil
				s.windows = append(s.windows, temp)
				s.windows = append(s.windows[:i], s.windows[i+1:]...)

				winPosX, winPosY := temp.GetPosition()
				grabX = curX - int(winPosX)
				grabY = curY - int(winPosY)
				s.activeWindow = temp
				temp.OnLButtonDown()
				break
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		windowGrab = nil

		for i := len(s.windows) - 1; i >= 0; i-- {
			rc := s.windows[i].GetRect()

			if curX >= rc.X1 && curY >= rc.Y1 &&
				curX < rc.X2 && curY < rc.Y2 {
				s.windows[i].OnLButtonUp()
				break
			}
		}
	}

	/*

		if s.game.appTicker%30 == 0 {
			s.windows = append(s.windows, NewXPMessageBox(s, "Файл kernel32.dll не найден.", "Ошибка", 0))
		}
	*/
}

func (s *WinXPScreen) LoadResources() {
	rm := ResourceManager_GetInstance()
	s.wallpaperImage = rm.LoadImage("assets/computer/wallpaper.png")
	s.taskbandImage = rm.LoadImage("assets/computer/taskband.png")
	s.startButtonImage = rm.LoadImage("assets/computer/appmenu.png")
	s.font = rm.LoadFontJSON("font/font_fantasy.json")
}

func (s *WinXPScreen) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screenWidth)/800, float64(screenHeight)/600)
	screen.DrawImage(s.wallpaperImage, op)

	for _, window := range s.windows {
		window.Draw(screen)
	}

	s.game.DrawNineGrid(screen, s.taskbandImage, 1.0, NineGridInfo{Top: 15, Bottom: 11, Left: 1, Right: 1}, 0, screenHeight-30, screenWidth, 30)
	s.game.DrawStatedNineGrid(screen, s.startButtonImage, 0, 3, 1.0, NineGridInfo{Left: 6, Top: 13, Right: 52, Bottom: 14}, 0, screenHeight-30, 100, 30)

	fontRenderer := s.game.fontRenderer

	fontRenderer.PushState()
	fontRenderer.SetTextColor(color.White)

	fontRenderer.PushState()
	fontRenderer.SetScale(2.0)
	fontRenderer.EnableShadow(true)

	startString := "Пуск"
	strDim := fontRenderer.GetStringDimensions(startString)
	fontRenderer.DrawTextAt(screen, startString, Vec2f{(100 - strDim.X) / 2, screenHeight - (30+strDim.Y)/2})

	fontRenderer.PopState()

	hour, min, _ := time.Now().Clock()
	time := fmt.Sprintf("%d:%d", hour, min)
	strDim2 := fontRenderer.GetStringDimensions(time)
	fontRenderer.DrawTextAt(screen, time, Vec2f{float64(screenWidth - 64), float64(screenHeight - (30+strDim2.Y)/2)})

	fontRenderer.PopState()
}

func (s *WinXPScreen) ProcessKeyEvents() bool {
	return true
}

func NewWinXPScreen(g *Game, nsb NextScreenBuilder) *WinXPScreen {
	s := new(WinXPScreen)
	s.IScreen = s
	s.game = g
	s.cwDefX = 10
	s.cwDefY = 10
	s.nextScreenBuilder = nsb

	s.windows = append(s.windows, NewXPTaskbarWindow(s))
	return s
}

func (g *Game) DrawWinXPCaption(screen *ebiten.Image, x float64, y float64, width float64, height float64) {
	g.DrawStatedNineGrid(screen, xpCaption, 0, 2, 1.0, NineGridInfo{Right: 35, Top: 9, Left: 28, Bottom: 17}, x, y, width, height)
}

func (g *Game) DrawWinXPLeftFrame(screen *ebiten.Image, x float64, y float64, width float64, height float64) {
	g.DrawStatedNineGrid(screen, xpLeftFrame, 0, 2, 1.0, NineGridInfo{Right: 2, Top: 0, Left: 2, Bottom: 0}, x, y, width, height)
}

func (g *Game) DrawWinXPRightFrame(screen *ebiten.Image, x float64, y float64, width float64, height float64) {
	g.DrawStatedNineGrid(screen, xpRightFrame, 0, 2, 1.0, NineGridInfo{Right: 2, Top: 0, Left: 2, Bottom: 0}, x, y, width, height)
}

func (g *Game) DrawWinXPBottomFrame(screen *ebiten.Image, x float64, y float64, width float64, height float64) {
	g.DrawStatedNineGrid(screen, xpBottomFrame, 0, 2, 1.0, NineGridInfo{Right: 5, Top: 2, Left: 5, Bottom: 2}, x, y, width, height)
}

func (g *Game) DrawWinXPCloseButton(screen *ebiten.Image, inactive bool, state int, x float64, y float64, width float64, height float64) {

	if inactive {
		state += 4
	}

	g.DrawStatedNineGrid(screen, xpCloseButton, state, 8, 1.0, NineGridInfo{Right: 5, Top: 5, Left: 5, Bottom: 5}, x, y, width, height)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x+4, y+4)
	screen.DrawImage(xpCloseGlyph.SubImage(image.Rect(0, 0, 13, 13)).(*ebiten.Image), op)
}
