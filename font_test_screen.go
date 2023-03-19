package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type FontTestScreen struct {
	Screen
	font *Font
}

func (s *FontTestScreen) LoadResources() {
	rm := ResourceManager_GetInstance()
	s.font = rm.LoadFontJSON("font/font_runic.json")
}

var pangrams = []string{
	"Съешь же ещё этих мягких французских булок, да выпей чаю.",
	"В чащах юга жил бы цитрус? Да, но фальшивый экземпляр!",
	"Эй, жлоб! Где туз? Прячь юных съёмщиц в шкаф.",
	"БУКВОПЕЧАТАЮЩЕЙ СВЯЗИ НУЖНЫ ХОРОШИЕ Э/МАГНИТНЫЕ РЕЛЕ. ДАТЬ ЦИФРЫ (123456789+=.?-)",
}

func (s *FontTestScreen) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.White)

	fontRenderer := s.game.fontRenderer
	fontRenderer.PushState()

	fontRenderer.SetFont(s.font)
	//fontRenderer.SetTextColor(alpacolor.Red)
	//fontRenderer.EnableShadow(true)

	var y float64

	for i := 0.0; i < 5.0; i += 0.5 {
		fontRenderer.SetScale(i + 1.0)
		fontRenderer.DrawTextAt(screen, pangrams[3], Vec2f{0, y})

		y += fontRenderer.GetGlyphSize().Y + 16
	}

	fontRenderer.PopState()
}

func InitFontTestScreen(s *FontTestScreen, g *Game) {
	s.game = g
}

func NewFontTestScreen(g *Game) *FontTestScreen {
	s := new(FontTestScreen)
	s.IScreen = s
	InitFontTestScreen(s, g)

	return s
}
