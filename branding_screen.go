package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Branding screen at the game startup
type BrandingScreen struct {
	Screen
	nextScreenBuilder NextScreenBuilder
	skip              bool
	aragajagaImage    *ebiten.Image
}

func (s *BrandingScreen) LoadResources() {
	rm := ResourceManager_GetInstance()

	s.aragajagaImage = rm.LoadImage("assets/aragajaga.png")
}

// Sets the this.skip to true, which then would be processed in Update method to switch to another screen
func (s *BrandingScreen) Skip() {
	s.skip = true
}

func (s *BrandingScreen) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.White)
	op := &ebiten.DrawImageOptions{}

	fadeDuration := 1.0
	fadeTPSDuration := float64(ebiten.MaxTPS()) * fadeDuration
	fadeFactor := 1.0

	if float64(s.game.appTicker) < fadeTPSDuration {
		fadeFactor = SinFade(float64(s.game.appTicker) / fadeTPSDuration)
	}

	iWidth, iHeight := s.aragajagaImage.Size()

	op.GeoM.Translate(-(float64(iWidth / 2)), -(float64(iHeight / 2)))
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	op.ColorM.Scale(1.0, 1.0, 1.0, fadeFactor)

	screen.DrawImage(s.aragajagaImage, op)
}

func (s *BrandingScreen) ProcessKeyEvents() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.Skip()
		return false
	}

	return true
}

func (s *BrandingScreen) Update() {
	s.ProcessKeyEvents()

	if s.nextScreenBuilder != nil {
		endDuration := 2.0
		endTPSDuration := float64(ebiten.MaxTPS()) * endDuration

		if float64(s.game.appTicker) >= endTPSDuration || s.skip {
			s.game.SetScreen(s.nextScreenBuilder(s.game))
		}
	}
}

func NewBrandingScreen(g *Game, nextScreenBuilder func(*Game) IScreen) *BrandingScreen {
	screen := new(BrandingScreen)
	screen.IScreen = screen
	screen.game = g
	screen.nextScreenBuilder = nextScreenBuilder

	return screen
}
