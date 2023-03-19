package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ISliderWidget interface {
	IWidget
	GetValue() float64
	SetValue(float64)
}

type SliderWidget struct {
	Widget
	value float64
}

func (widget *SliderWidget) ProcessKeyEvents() bool {

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		widget.value = math.Max(0.0, widget.value-0.1)
		return false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		widget.value = math.Min(1.0, widget.value+0.1)
		return false
	}

	return true
}

func (widget *SliderWidget) Update() {
	widget.ProcessKeyEvents()
}

func (widget *SliderWidget) Draw(screen *ebiten.Image) {
	widget.game.DrawGUIFrame(screen, widget.PosX, widget.PosY, widget.Width, widget.Height)

	widget.game.DrawGUIButton(screen, buttonStateNormal, widget.PosX+8+(widget.Width-48)*widget.value,
		widget.PosY+8, 32, widget.Height-16)

	if widget.selected {
		widget.game.DrawGUIButton(screen, buttonStateHover, widget.PosX+8+(widget.Width-48)*widget.value,
			widget.PosY+8, 32, widget.Height-16)
	}
}

func CreateCommonSliderWidget(g *Game, callback func(*Game)) *SliderWidget {
	slider := new(SliderWidget)
	InitializeCommonWidget(&slider.Widget, g)

	slider.value = 0.5
	slider.text = ""
	slider.callback = callback
	return slider
}
