package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EditMode struct {
	cursor         TileCursor
	game           *Game
	brushTile      Tile
	swapSampleView bool
	selLayer       int
}

func (m *EditMode) Draw(screen *ebiten.Image) {
	g := m.game

	DrawTileCursor := func() {
		cameraZoom := g.camera.GetZoom()
		pos := g.camera.WorldToScreen2(m.cursor.GetWorldPos())

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-2, -2)
		op.GeoM.Scale(cameraZoom, cameraZoom)
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(tileCursor, op)

		if !m.swapSampleView {
			tile := GetTileSprite(tilesImage, tileXNum, tileSize, m.brushTile)
			op.GeoM.Translate(2*cameraZoom, 2*cameraZoom)
			op.ColorM.Scale(1.0, 1.0, 1.0, 0.8+math.Sin(float64(tickCounter)/2.0)*0.2)
			screen.DrawImage(tile, op)
		} else {
			tilePos := m.cursor.y*g.level.width + m.cursor.x
			selTile := g.level.tileLayers[m.selLayer][tilePos]

			tile := GetTileSprite(tilesImage, tileXNum, tileSize, 63)

			op.GeoM.Translate(2*cameraZoom, 2*cameraZoom)
			screen.DrawImage(tile, op)

			tile = GetTileSprite(tilesImage, tileXNum, tileSize, selTile)
			screen.DrawImage(tile, op)

			fontRenderer := g.fontRenderer

			fontRenderer.PushState()

			fontRenderer.SetScale(g.view.guiScale)
			fontRenderer.SetTextColor(color.White)
			fontRenderer.EnableShadow(true)

			labelPos := pos.Translate(Vec2f{0.0, -16.0})
			g.fontRenderer.DrawTextAt(screen,
				fmt.Sprintf("%s: %d", I18n("string_layer", "Layer"), m.selLayer), labelPos)

			fontRenderer.PopState()
		}
	}

	DrawTooltip := func(previewTile Tile, title string, lineProviders []func() (string, TextFormat)) {
		g.DrawHerbGUIFrame(screen, 0, float64(screenHeight-96), 256, 96)

		currentOp := &ebiten.DrawImageOptions{}
		previewTileScale := 2.0

		tile := GetTileSprite(tilesImage, tileXNum, tileSize, previewTile)

		// 4 pixels in GUI cordinates
		margin := float64(8)

		previewTileGUISize := Vec2f{tileSize, tileSize}
		previewTileGUISize = previewTileGUISize.Scale(previewTileScale)

		guiPos := Vec2f{margin, screenHeight/g.view.guiScale - previewTileGUISize.Y - margin}
		screenPos := guiPos.Scale(g.view.guiScale)

		currentOp.GeoM.Scale(previewTileScale*g.view.guiScale, previewTileScale*g.view.guiScale)
		currentOp.GeoM.Translate(screenPos.X, screenPos.Y)
		screen.DrawImage(tile, currentOp)

		guiPos = guiPos.Translate(previewTileGUISize.ScaleVec2f(Vec2f{1.0, 0.0}))
		guiPos = guiPos.Translate(Vec2f{margin, margin / 2})
		screenPos = guiPos.Scale(g.view.guiScale)

		fontRenderer := g.fontRenderer

		fontRenderer.PushState()
		fontRenderer.Reset()
		fontRenderer.SetScale(g.view.guiScale)

		g.fontRenderer.DrawTextAt(screen, title, screenPos)

		strDim := g.fontRenderer.GetStringDimensions(title)
		guiPos = guiPos.Translate(Vec2f{0, strDim.Y})
		guiPos = guiPos.Translate(Vec2f{0, 4.0})
		screenPos = guiPos.Scale(g.view.guiScale)

		fontRenderer.PopState()

		for _, lineProvider := range lineProviders {
			line, format := lineProvider()

			screenScaleFormat := format
			screenScaleFormat.scale = math.Max(1.0, g.view.guiScale*format.scale)

			fontRenderer := g.fontRenderer

			fontRenderer.PushState()
			fontRenderer.Reset()
			fontRenderer.SetScale(math.Max(1.0, g.view.guiScale+format.scale))

			fontRenderer.DrawTextAt(screen, line, screenPos)

			// strDim := fontRenderer.GetStringDimensions(line, math.Max(0.5, format.scale))
			strDim := fontRenderer.GetStringDimensions(line)
			guiPos = guiPos.Translate(Vec2f{0, strDim.Y})
			guiPos = guiPos.Translate(Vec2f{0, 4.0 * math.Max(1, format.scale)})
			screenPos = guiPos.Scale(g.view.guiScale)

			fontRenderer.PopState()
		}
	}

	DrawTileInfoTooltip := func(tile Tile, title string) {
		stringProviders := []func() (string, TextFormat){
			func() (string, TextFormat) {
				c := color.RGBA{0, 0, 0, 255}

				var tileName string
				_, has := tileDescStorage[tile]
				if has {
					tileName = I18n(tileDescStorage[tile].JSONName, tileDescStorage[tile].JSONName)
				} else {
					tileName = I18n("string_unknown", "Unknown")
					c = color.RGBA{255, 0, 0, 255}
				}

				return tileName, TextFormat{1.0, c, true}
			},
		}

		_, has := tileNameMap[tile]
		if has {
			stringProviders = append(stringProviders, func() (string, TextFormat) {
				walkableString := I18n("string_walkable", "Walkable")
				c := color.RGBA{0, 0, 255, 255}

				if !tileDescStorage[tile].Walkable {
					walkableString = I18n("string_solid", "Solid")
					c = color.RGBA{255, 0, 0, 255}
				}

				return walkableString, TextFormat{scale: 0.5, textColor: c, shadow: false}
			})
		}

		stringProviders = append(stringProviders, func() (string, TextFormat) {
			return fmt.Sprintf("ID: %d", tile), TextFormat{scale: 0.5, textColor: color.Black, shadow: false}
		})

		DrawTooltip(tile, title, stringProviders)
	}

	var modeTitle string = I18n("string_edit_mode", "Edit Mode")

	if m.swapSampleView {
		modeTitle += " (" + I18n("string_edit_mode_xray", "x-ray") + ")"
	} else {
		modeTitle += " (" + I18n("string_edit_mode_paint", "Paint") + ")"
	}

	g.DrawModeTitle(screen, modeTitle)

	DrawTileCursor()

	// Draw current layer tile
	if !m.swapSampleView {
		tilePos := m.cursor.y*g.level.width + m.cursor.x
		selTile := g.level.tileLayers[m.selLayer][tilePos]

		DrawTileInfoTooltip(selTile, fmt.Sprintf("%s: %d", I18n("string_layer", "Layer"), m.selLayer))
	} else {
		DrawTileInfoTooltip(m.brushTile, I18n("string_brush", "Brush"))
	}
}

func NewEditMode(g *Game) *EditMode {
	editMode := new(EditMode)
	editMode.game = g

	g.camera.TargetCursor(&editMode.cursor)
	return editMode
}

func (m *EditMode) ProcessKeyEvents() {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		m.cursor.x = int(math.Min(float64(m.game.level.width-1), float64(m.cursor.x+1)))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		m.cursor.x = int(math.Max(0, float64(m.cursor.x-1)))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.cursor.y = int(math.Max(0, float64(m.cursor.y-1)))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.cursor.y = int(math.Min(float64(m.game.level.height-1), float64(m.cursor.y+1)))
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbEditorPlace]) {
		tilePos := m.cursor.y*m.game.level.width + m.cursor.x

		m.game.level.tileLayers[m.selLayer][tilePos] = m.brushTile
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbEditorDelete]) {
		tilePos := m.cursor.y*m.game.level.width + m.cursor.x

		m.game.level.tileLayers[m.selLayer][tilePos] = tileIDEmpty
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbEditorPrevBrush]) {
		m.brushTile--
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbEditorNextBrush]) {
		m.brushTile++
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbEditorSwitchMode]) {
		m.swapSampleView = !m.swapSampleView
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbEditorPrevLayer]) {
		m.selLayer = int(math.Max(0, float64(m.selLayer-1)))
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbEditorNextLayer]) {
		m.selLayer = int(math.Min(float64(len(m.game.level.tileLayers)-1), float64(m.selLayer+1)))
	}
}
