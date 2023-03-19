package main

import "github.com/hajimehoshi/ebiten/v2"

type IWidget interface {
	Click()
	Draw(screen *ebiten.Image)
	GetPosition() Vec2f
	GetSize() Vec2f
	GetText() string
	ProcessKeyEvents() bool
	SetPosition(Vec2f)
	SetSelection(bool)
	SetSize(Vec2f)
	SetText(string)
}

type Widget struct {
	PosX, PosY, Width, Height float64
	text                      string
	game                      *Game
	selected                  bool
	callback                  func(*Game)
}

func (b *Widget) ProcessKeyEvents() bool {
	return true
}

func (b *Widget) GetText() string {
	return b.text
}

func (b *Widget) SetText(text string) {
	b.text = text
}

func (b *Widget) Click() {
	b.callback(b.game)
}

func (b *Widget) GetPosition() Vec2f {
	return Vec2f{b.PosX, b.PosY}
}

func (b *Widget) GetSize() Vec2f {
	return Vec2f{b.Width, b.Height}
}

func (b *Widget) SetPosition(pos Vec2f) {
	b.PosX = pos.X
	b.PosY = pos.Y
}

func (b *Widget) SetSize(size Vec2f) {
	b.Width = size.X
	b.Height = size.Y
}

func (b *Widget) SetSelection(selected bool) {
	b.selected = selected
}

func InitializeCommonWidget(widget *Widget, game *Game) {
	widget.Width = 200
	widget.Height = 40
	widget.game = game
}
