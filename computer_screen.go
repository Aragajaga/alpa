package main

import (
	alpacolor "github.com/aragajaga/alpa/util/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ComputerScreen struct {
	Screen
	childManager           ScreenManager
	nextScreenBuilder      NextScreenBuilder
	controlModifierPressed bool
	startTime              int
}

func (cs *ComputerScreen) ProcessKeyEvents() bool {
	if inpututil.IsKeyJustReleased(ebiten.KeyControlRight) {
		cs.controlModifierPressed = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyControlRight) {
		cs.controlModifierPressed = true
	}

	if cs.controlModifierPressed && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		cs.game.SetScreen(cs.nextScreenBuilder(cs.game))
		return false
	}

	return true
}

func (cs *ComputerScreen) Update() {
	cs.ProcessKeyEvents()

	cs.childManager.currentScreen.Update()
}

func (cs *ComputerScreen) Draw(screen *ebiten.Image) {
	cs.childManager.currentScreen.Draw(screen)

	fontRenderer := cs.game.fontRenderer

	fontRenderer.PushState()
	fontRenderer.SetTextColor(alpacolor.Red)
	fontRenderer.SetScale(2.0)

	var ticksElapsed int = cs.game.appTicker - cs.startTime
	var text string = I18n("string_pc_release_tip", "Press Right Ctrl + R to release")
	textDim := cs.game.fontRenderer.GetStringDimensions(text)
	var y float64 = textDim.Y

	if ticksElapsed <= 30 {
		y = -textDim.Y + (textDim.Y*2)*(float64(ticksElapsed)/30)
	} else if ticksElapsed >= 120 && ticksElapsed <= 150 {
		y = -textDim.Y + (textDim.Y*2)*(1.0-float64(ticksElapsed-120)/30)
	}

	if ticksElapsed <= 150 {
		cs.game.DrawHerbGUIFrame(screen, (float64(screenWidth)-textDim.X)/2, y, textDim.X, textDim.Y)
		cs.game.fontRenderer.DrawTextAt(screen, text, Vec2f{(float64(screenWidth) - textDim.X) / 2, y})
	}

	fontRenderer.PopState()
}

func NewComputerScreen(g *Game, nsb NextScreenBuilder) *ComputerScreen {
	screen := new(ComputerScreen)
	screen.IScreen = screen
	screen.game = g
	screen.nextScreenBuilder = nsb
	screen.startTime = g.appTicker

	screen.childManager.SetScreen(NewBIOSLoadingScreen(g, &screen.childManager))

	g.audioManager.PlayBackgroundMusic("bgm/computer")

	return screen
}
