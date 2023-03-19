package main

import (
	"encoding/json"
	"errors"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type IFont interface {
	GetGlyphSize() Vec2i
	GetGlyphHeight() int
	GetGlyphWidth() int
	GetName() string
}

type Font struct {
	IFont
	name        string
	charmap     *ebiten.Image
	glyphWidth  int
	glyphHeight int
}

func (font *Font) GetGlyphSize() Vec2i {
	return Vec2i{font.glyphWidth, font.glyphHeight}
}

func (font *Font) GetGlyphWidth() int {
	return font.glyphWidth
}

func (font *Font) GetGlyphHeight() int {
	return font.glyphHeight
}

func (font *Font) GetName() string {
	return font.name
}

type IFontRenderer interface {
	GetTextColor() color.Color
	SetTextColor(color.Color)

	GetShadowColor() color.Color
	SetShadowColor(color.Color)

	IsShadowEnabled() bool
	EnableShadow(bool)

	GetScale() float64
	SetScale(float64)

	GetFont() IFont
	SetFont(IFont)

	GetStringDimensions() Vec2f
	GetGlyphSize() Vec2f

	PushState()
	PopState()

	Reset()
}

type ShadowOffsetMode uint8

const (
	shadowOffsetModeLocal ShadowOffsetMode = iota
	shadowOffsetModeScreen
)

type FormatOptionsLIFO struct {
	format *FontRendererFormatOptions
	next   *FormatOptionsLIFO
}

type FontRendererFormatOptions struct {
	font             *Font
	scale            float64
	textColor        color.Color
	shadowColor      color.Color
	shadow           bool
	shadowOffsetX    float64
	shadowOffsetY    float64
	shadowOffsetMode ShadowOffsetMode
}

type FontRenderer struct {
	IFontRenderer
	optHead *FormatOptionsLIFO
	op      ebiten.DrawImageOptions
}

func (fontRenderer *FontRenderer) GetTextColor() color.Color {
	return fontRenderer.optHead.format.textColor
}

func (fontRenderer *FontRenderer) SetTextColor(clr color.Color) {
	fontRenderer.optHead.format.textColor = clr
}

func (fontRenderer *FontRenderer) GetShadowColor() color.Color {
	return fontRenderer.optHead.format.shadowColor
}

func (fontRenderer *FontRenderer) SetShadowColor(clr color.Color) {
	fontRenderer.optHead.format.shadowColor = clr
}

func (fontRenderer *FontRenderer) IsShadowEnabled() bool {
	return fontRenderer.optHead.format.shadow
}

func (fontRenderer *FontRenderer) EnableShadow(shadow bool) {
	fontRenderer.optHead.format.shadow = shadow
}

func (fontRenderer *FontRenderer) GetScale() float64 {
	return fontRenderer.optHead.format.scale
}

func (fontRenderer *FontRenderer) SetScale(scale float64) {
	fontRenderer.optHead.format.scale = scale
}

func (fontRenderer *FontRenderer) SetFont(font IFont) {
	fontRenderer.optHead.format.font = font.(*Font)
}

func (fontRenderer *FontRenderer) GetFont() IFont {
	return fontRenderer.optHead.format.font
}

func (fontRenderer *FontRenderer) GetStringDimensions(text string) Vec2f {
	dim := fontRenderer.GetFont().GetGlyphSize().ToVec2f()
	dim = dim.Translate(Vec2f{1.0, 0})
	dim = dim.ScaleVec2f(Vec2f{float64(len([]rune(text))), 1.0})
	dim = dim.Scale(fontRenderer.GetScale())
	return dim
}

func (fontRenderer *FontRenderer) GetGlyphSize() Vec2f {
	dim := fontRenderer.GetFont().GetGlyphSize().ToVec2f()
	dim = dim.Translate(Vec2f{1.0, 0})
	dim = dim.Scale(fontRenderer.GetScale())
	return dim
}

func (fontRenderer *FontRenderer) PushState() {
	format := new(FontRendererFormatOptions)
	*format = *fontRenderer.optHead.format

	node := new(FormatOptionsLIFO)
	node.format = format
	node.next = fontRenderer.optHead

	fontRenderer.optHead = node
}

func (fontRenderer *FontRenderer) PopState() {
	if fontRenderer.optHead.next == nil {
		log.Println("[FontRenderer] Warning: Tried to pop the last state. It is an abnorbal behaviour, check your code.")
		return
	}

	fontRenderer.optHead = fontRenderer.optHead.next
}

var g_fontRendererDefaultFormatOptions = FontRendererFormatOptions{
	font:             nil,
	textColor:        color.Black,
	shadowColor:      color.RGBA{0, 0, 0, 128},
	scale:            1.0,
	shadow:           false,
	shadowOffsetX:    1,
	shadowOffsetY:    1,
	shadowOffsetMode: shadowOffsetModeLocal,
}

func (fontRenderer *FontRenderer) Reset() {
	rm := ResourceManager_GetInstance()
	font := rm.LoadFontJSON("font/font_system.json")

	g_fontRendererDefaultFormatOptions.font = font
	*fontRenderer.optHead.format = g_fontRendererDefaultFormatOptions
}

func (fontRenderer *FontRenderer) _DrawGlyph(screen *ebiten.Image, glyph *ebiten.Image) {
	format := fontRenderer.optHead.format

	if format.shadow {
		geomSave := fontRenderer.op.GeoM

		/*
			translateMat := ebiten.GeoM{}
			translateMat.Reset()

			translateMat.Translate(format.shadowOffsetX, format.shadowOffsetY)

			x, y := 1.0, 1.0
			x, y := translateMat.Apply(x, y)

			fontRenderer.op.GeoM.Translate(x, y)
		*/
		fontRenderer.op.GeoM.Translate(format.shadowOffsetX, format.shadowOffsetY)

		fontRenderer.op.ColorM.Reset()
		ColorM_Colorize(&fontRenderer.op.ColorM, format.shadowColor)
		screen.DrawImage(glyph, &fontRenderer.op)

		fontRenderer.op.GeoM = geomSave
	}

	fontRenderer.op.ColorM.Reset()
	ColorM_Colorize(&fontRenderer.op.ColorM, format.textColor)
	screen.DrawImage(glyph, &fontRenderer.op)
}

/*

func (fr *FontRenderer) _DrawTextRect_DrawGlyph(screen *ebiten.Image, ch rune, scale float64, cx int, cy int, tableOffsetY int, chOffset int, i int, op *ebiten.DrawImageOptions, x float64, y float64, width int, height int) {
	chTablePos := int(ch) - chOffset

	sx := (chTablePos % 16) * fr.glyphWidth
	sy := (chTablePos/16 + tableOffsetY) * fr.glyphHeight

	glyph := fr.tilemap.SubImage(image.Rect(sx, sy, sx+fr.glyphWidth, sy+fr.glyphHeight)).(*ebiten.Image)

	glyphOuterWidth := fr.glyphWidth + GLYPH_HORIZONTAL_SPACING
	glyphOuterHeight := fr.glyphHeight + GLYPH_VERTICAL_SPACING

	op.GeoM.Reset()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)

	op.GeoM.Translate(float64(cx*glyphOuterWidth)*scale, float64(cy*glyphOuterHeight)*scale)

	fr._DrawGlyph(screen, glyph, color.RGBA{0, 0, 0, 255}, op)
}

func (fr *FontRenderer) DrawTextRect(screen *ebiten.Image, text string, scale float64, x float64, y float64, width int, height int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)

	glyphOuterWidth := fr.glyphWidth + GLYPH_HORIZONTAL_SPACING
	// glyphOuterHeight := fr.glyphHeight + GLYPH_VERTICAL_SPACING

	screenBufferWidth := int(float64(width) / scale / float64(glyphOuterWidth))
	// screenBufferHeight := height / glyphOuterHeight

	var cx, cy int = 0, 0

	for i, ch := range []rune(text) {
		if ch == '\n' {
			cx = 0
			cy++
		} else {
			if ch >= UNICODE_LATIN_FIRST && ch <= UNICODE_LATIN_LAST {
				fr._DrawTextRect_DrawGlyph(screen, ch, scale, cx, cy, FONTTILEMAP_ROW_LATIN, UNICODE_LATIN_FIRST, i, op, x, y, width, height)
			} else if ch >= UNICODE_CYRILLIC_FIRST && ch <= UNICODE_CYRILLIC_LAST {
				fr._DrawTextRect_DrawGlyph(screen, ch, scale, cx, cy, FONTTILEMAP_ROW_CYRILLIC, UNICODE_CYRILLIC_FIRST, i, op, x, y, width, height)
			}

			cx++
		}

		if cx != 0 && cx%screenBufferWidth == 0 {
			cx = 0
			cy++
		}

	}
}
*/

const (
	CHARMAP_LATIN_OFFSET    = 0
	CHARMAP_CYRILLIC_OFFSET = 96
)

func (fontRenderer *FontRenderer) _DrawTextAt_DrawGlyph(screen *ebiten.Image, ch rune) {
	font := fontRenderer.optHead.format.font

	atlasWidth, _ := font.charmap.Size()

	glyphSize := font.GetGlyphSize()
	colCount := atlasWidth / glyphSize.X

	chOffset := CHARMAP_LATIN_OFFSET + int(rune('?')) - UNICODE_LATIN_FIRST

	if ch >= UNICODE_LATIN_FIRST && ch <= UNICODE_LATIN_LAST {
		chOffset = CHARMAP_LATIN_OFFSET + int(ch) - UNICODE_LATIN_FIRST
	} else if ch >= UNICODE_CYRILLIC_FIRST && ch <= UNICODE_CYRILLIC_LAST {
		chOffset = CHARMAP_CYRILLIC_OFFSET + int(ch) - UNICODE_CYRILLIC_FIRST
	}

	sx := (chOffset % colCount) * glyphSize.X
	sy := (chOffset / colCount) * glyphSize.Y

	glyph := font.charmap.SubImage(image.Rect(sx, sy, sx+glyphSize.X, sy+glyphSize.Y)).(*ebiten.Image)
	fontRenderer._DrawGlyph(screen, glyph)
}

func (fontRenderer *FontRenderer) DrawTextAt(screen *ebiten.Image, text string, pos Vec2f) {
	format := fontRenderer.optHead.format

	fontRenderer.op.GeoM.Reset()
	fontRenderer.op.ColorM.Reset()

	localTranslation := ebiten.GeoM{}
	localTranslation.Reset()
	screenScale := ebiten.GeoM{}
	screenScale.Reset()
	screenTranslation := ebiten.GeoM{}
	screenTranslation.Reset()

	screenScale.Scale(format.scale, format.scale)
	screenTranslation.Translate(pos.X, pos.Y)

	// localTranslation * screenScale * screenTranslation

	for _, ch := range text {
		fontRenderer.op.GeoM.Reset()
		fontRenderer.op.GeoM.Concat(localTranslation)
		fontRenderer.op.GeoM.Concat(screenScale)
		fontRenderer.op.GeoM.Concat(screenTranslation)

		fontRenderer._DrawTextAt_DrawGlyph(screen, rune(ch))

		localTranslation.Translate(float64(format.font.GetGlyphWidth()+1), 0)
	}
	time.Sleep(5000)
}

func (fontRenderer *FontRenderer) DrawTextFormattedAt(screen *ebiten.Image, text string, format TextFormat, x float64, y float64) {
	fontRenderer.PushState()

	fontRenderer.EnableShadow(format.shadow)
	fontRenderer.SetScale(format.scale)
	fontRenderer.SetTextColor(format.textColor)

	fontRenderer.DrawTextAt(screen, text, Vec2f{x, y})

	fontRenderer.PopState()
}

func IsWordChar(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= 'А' && ch <= 'Я') || (ch >= 'а' && ch <= 'я')
}

func InitFontRenderer(fontRenderer *FontRenderer) {

	format := new(FontRendererFormatOptions)
	node := new(FormatOptionsLIFO)
	node.format = format
	node.next = nil

	fontRenderer.optHead = node
	fontRenderer.Reset()
}

func NewFontRenderer() *FontRenderer {
	s := new(FontRenderer)
	InitFontRenderer(s)
	return s
}

type JSONFontData struct {
	Name        string `json:"name"`
	Asset       string `json:"asset"`
	GlyphWidth  int    `json:"glyph_width"`
	GlyphHeight int    `json:"glyph_height"`
	LineHeight  int    `json:"line_height"`
}

func LoadFontFromJSON(path string) (*Font, error) {
	font := new(Font)
	if font == nil {
		return nil, errors.New("Font object allocation failed")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var jsonFont JSONFontData
	err = json.Unmarshal(data, &jsonFont)
	if err != nil {
		return nil, err
	}

	font.charmap, _, err = ebitenutil.NewImageFromFile("font/" + jsonFont.Asset)
	if err != nil {
		return nil, err
	}

	font.name = jsonFont.Name
	font.glyphWidth = jsonFont.GlyphWidth
	font.glyphHeight = jsonFont.GlyphHeight

	return font, nil
}
