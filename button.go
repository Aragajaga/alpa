package main

import (
	alpacolor "github.com/aragajaga/alpa/util/color"
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	Widget
}

func (button *Button) Draw(screen *ebiten.Image) {
	pos := Vec2f{button.PosX, button.PosY}

	/*
		var buttonState int

		buttonState = buttonStateNormal

		if button.selected {
			buttonState = buttonStateHover
		}
	*/

	// button.game.DrawGUIButton(screen, buttonState, pos.X, pos.Y, button.Width, button.Height)

	fontRenderer := button.game.fontRenderer

	fontRenderer.PushState()

	fontRenderer.Reset()
	fontRenderer.SetScale(2.0)
	fontRenderer.SetTextColor(alpacolor.Red)

	strDim := fontRenderer.GetStringDimensions(button.text)
	pos = pos.Add(Vec2f{button.Width, button.Height}.Subtract(strDim).Scale(0.5))

	button.game.fontRenderer.DrawTextAt(screen, button.text, pos)

	fontRenderer.PopState()
}

func CreateCommonButton(g *Game, text string, callback func(*Game)) *Button {
	button := new(Button)
	InitializeCommonWidget(&button.Widget, g)

	button.text = text
	button.callback = callback
	return button
}
