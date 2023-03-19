package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type DebugStringBuilder func(*Game) string

type DebugScreen struct {
	stringBuilders []DebugStringBuilder
	game           *Game
}

func (ds *DebugScreen) Draw(screen *ebiten.Image) {
	fontRenderer := ds.game.fontRenderer

	fontRenderer.PushState()

	fontRenderer.Reset()
	fontRenderer.SetScale(ds.game.view.guiScale)
	fontRenderer.SetTextColor(color.White)
	fontRenderer.EnableShadow(true)

	for i, builder := range ds.stringBuilders {
		ds.game.fontRenderer.DrawTextAt(screen, builder(ds.game), Vec2f{0.0, float64(24 * i)})
	}

	fontRenderer.PopState()
}

func _DSB_CharWorldPos(game *Game) string {
	ch := game.char

	return fmt.Sprintf("World pos x: %f, y: %f", ch.GetWorldPos().X, ch.GetWorldPos().Y)
}

func _DSB_CharTilePos(game *Game) string {
	ch := game.char

	tilePos, _ := game.WorldPosToTilePos(ch.GetWorldPos().X, ch.GetWorldPos().Y)
	return fmt.Sprintf("Tile pos: %d", tilePos)
}

func _DBS_CurrentTPS(game *Game) string {
	return fmt.Sprintf("TPS: %f, TPS Max: %d", ebiten.CurrentTPS(), ebiten.MaxTPS())
}

func _DBS_CurrentFPS(game *Game) string {
	return fmt.Sprintf("FPS: %f", ebiten.CurrentFPS())
}

func CreateDebugScreen(g *Game) *DebugScreen {
	ds := new(DebugScreen)
	ds.game = g

	ds.stringBuilders = []DebugStringBuilder{
		_DBS_CurrentTPS,
		_DBS_CurrentFPS,
		_DSB_CharWorldPos,
		_DSB_CharTilePos,
	}

	return ds
}
