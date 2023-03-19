package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type IGenericWidgetContainerScreen interface {
	IScreen
	SetTitle(string)
	GetTitle() string
}

type GenericWidgetContainerScreen struct {
	Screen
	widgets       []IWidget
	focusedWidget IWidget
	title         string
}

func (s *GenericWidgetContainerScreen) SetInitialFocus() {
	if len(s.widgets) > 0 {
		s.focusedWidget = s.widgets[0]
		s.focusedWidget.SetSelection(true)
	}
}

func (s *GenericWidgetContainerScreen) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), color.RGBA{0, 0, 0, 128 + 64})

	var containerWidth float64
	var containerHeight float64
	for _, widget := range s.widgets {
		containerWidth = math.Max(float64(containerWidth), widget.GetSize().X)
		containerHeight += widget.GetSize().Y + 4
	}

	x := (float64(screenWidth) - containerWidth) / 2
	y := (float64(screenHeight) - containerHeight) / 2

	s.game.DrawHerbGUIFrame(screen, x-32, y-32, containerWidth+64, containerHeight+64+64)

	fontRenderer := s.game.fontRenderer

	fontRenderer.PushState()
	fontRenderer.SetScale(2.0)

	titleDim := s.game.fontRenderer.GetStringDimensions(s.title)
	s.game.fontRenderer.DrawTextAt(screen, s.title, Vec2f{(screenWidth - titleDim.X) / 2, y})

	for _, widget := range s.widgets {
		widget.Draw(screen)
	}

	fontRenderer.PopState()
}

func (gm *GenericWidgetContainerScreen) GetTitle() string {
	return gm.title
}

func (gm *GenericWidgetContainerScreen) SetTitle(title string) {
	gm.title = title
}

func (gm *GenericWidgetContainerScreen) TabStopPrev() {
	gm.focusedWidget.SetSelection(false)

	for i, widget := range gm.widgets {
		if gm.focusedWidget == widget {
			if i == 0 {
				gm.focusedWidget = gm.widgets[len(gm.widgets)-1]
			} else {
				gm.focusedWidget = gm.widgets[i-1]
			}
			break
		}
	}

	gm.focusedWidget.SetSelection(true)
}

func (gm *GenericWidgetContainerScreen) TabStopNext() {
	gm.focusedWidget.SetSelection(false)

	for i, widget := range gm.widgets {
		if gm.focusedWidget == widget {
			if i == len(gm.widgets)-1 {
				gm.focusedWidget = gm.widgets[0]
			} else {
				gm.focusedWidget = gm.widgets[i+1]
			}
			break
		}
	}

	gm.focusedWidget.SetSelection(true)
}

func (gm *GenericWidgetContainerScreen) ProcessKeyEvents() bool {

	if gm.focusedWidget.ProcessKeyEvents() {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
			inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			gm.TabStopNext()
			return false
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
			inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			gm.TabStopPrev()
			return false
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			gm.focusedWidget.Click()
			return false
		}
	}

	return true
}

func (gm *GenericWidgetContainerScreen) Update() {
	var containerWidth float64
	var containerHeight float64
	for _, widget := range gm.widgets {
		containerWidth = math.Max(float64(containerWidth), widget.GetSize().X)
		containerHeight += widget.GetSize().Y + 4
	}

	x := (float64(screenWidth) - containerWidth) / 2
	y := 64 + (float64(screenHeight)-containerHeight)/2

	pos := Vec2f{x, y}

	for _, widget := range gm.widgets {
		widget.SetPosition(pos)
		pos = pos.Add(Vec2f{0, widget.GetSize().Y + 4})
	}
}

func InititalizeGenericWidgetContainerScreen(s *GenericWidgetContainerScreen) {
}
