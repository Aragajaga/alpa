package main

import (
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EditBox struct {
	Widget
}

func (eb *EditBox) Draw(screen *ebiten.Image) {
	pos := Vec2f{eb.PosX, eb.PosY}

	eb.game.DrawGUIFrame(screen, pos.X, pos.Y, eb.Width, eb.Height)

	if eb.selected {
		ebitenutil.DrawRect(screen, pos.X, pos.Y, eb.Width, eb.Height, color.RGBA{0, 255, 0, 64})
	}

	fontRenderer := eb.game.fontRenderer

	fontRenderer.PushState()

	fontRenderer.Reset()
	fontRenderer.SetScale(2.0)

	strDim := fontRenderer.GetStringDimensions(eb.text)
	pos = pos.Add(Vec2f{eb.Width, eb.Height}.Subtract(strDim).Scale(0.5))

	eb.game.fontRenderer.DrawTextAt(screen, eb.text, pos)

	fontRenderer.PopState()
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func (eb *EditBox) ProcessKeyEvents() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		eb.text = trimLastChar(eb.text)
		return false
	} else {
		runes := ebiten.InputChars()
		eb.text = eb.text + string(runes)
	}

	return true
}

func CreateCommonEditBox(g *Game, callback func(*Game)) *EditBox {
	editBox := new(EditBox)
	InitializeCommonWidget(&editBox.Widget, g)

	editBox.text = ""
	editBox.callback = callback
	return editBox
}
