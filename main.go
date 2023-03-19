package main

import (
	"encoding/json"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	alpacolor "github.com/aragajaga/alpa/util/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TextGrid struct {
	rows [][]rune
}

const (
	screenWidth  = 640
	screenHeight = 480
)

const (
	tileSize = 16
	tileXNum = 8
)

const (
	GUIFRAMETILE_NW     uint8 = 0
	GUIFRAMETILE_TOP    uint8 = 1
	GUIFRAMETILE_NE     uint8 = 2
	GUIFRAMETILE_LEFT   uint8 = 3
	GUIFRAMETILE_RIGHT  uint8 = 4
	GUIFRAMETILE_SW     uint8 = 5
	GUIFRAMETILE_BOTTOM uint8 = 6
	GUIFRAMETILE_SE     uint8 = 7
)

const (
	GUIFRAMETILE_SIZE int = 5
)

const (
	FONTTILEMAP_ROW_LATIN    = 0
	FONTTILEMAP_ROW_CYRILLIC = 6
)

const (
	UNICODE_LATIN_FIRST    = '\u0020'
	UNICODE_LATIN_LAST     = '\u007F'
	UNICODE_CYRILLIC_FIRST = '\u0400'
	UNICODE_CYRILLIC_LAST  = '\u045F'
)

const (
	GLYPH_HORIZONTAL_SPACING = 1
	GLYPH_VERTICAL_SPACING   = 2
)

var (
	tilesImage             *ebiten.Image
	charSprite             *ebiten.Image
	morgenSprite           *ebiten.Image
	flanSprite             *ebiten.Image
	monobearSprite         *ebiten.Image
	explosionSprite        *ebiten.Image
	guiFrameHerb           *ebiten.Image
	guiFrameTest           *ebiten.Image
	guiButton              *ebiten.Image
	tileCursor             *ebiten.Image
	skillMonobearExplosion *ebiten.Image
	seeYaTileSet           *ebiten.Image
	michaelSprite          *ebiten.Image
	worldBorderImage       *ebiten.Image
	tickCounter            int
)

type KeyBind uint8
type KeyBindMap map[KeyBind]ebiten.Key

const (
	kbPlayerMoveRight           KeyBind = 1
	kbPlayerMoveLeft            KeyBind = 2
	kbPlayerMoveUp              KeyBind = 3
	kbPlayerMoveDown            KeyBind = 4
	kbToggleEditMode            KeyBind = 5
	kbToggleEntityFocusRotation KeyBind = 6
	kbEditorNextLayer           KeyBind = 7
	kbEditorPrevLayer           KeyBind = 8
	kbEditorNextBrush           KeyBind = 9
	kbEditorPrevBrush           KeyBind = 10
	kbEditorPlace               KeyBind = 11
	kbEditorDelete              KeyBind = 12
	kbEditorSwitchMode          KeyBind = 13
	kbShowDebugInfo             KeyBind = 14
	kbWorldZoomIn               KeyBind = 15
	kbWorldZoomOut              KeyBind = 16
)

var keyBinds KeyBindMap

var r *rand.Rand

type LookDirection int64

const (
	LooksRight LookDirection = 0
	LooksLeft  LookDirection = 1
	LooksUp    LookDirection = 2
	LooksDown  LookDirection = 3
)

type View struct {
	guiScale float64
}

type Tile uint8

type TileLayer []Tile

type Level struct {
	tileLayers []TileLayer
	width      int
	height     int
	fileName   string
}

func (lv *Level) ReplaceAll(from, to Tile) {
	for i, layer := range lv.tileLayers {
		for j, tile := range layer {
			if tile == from {
				lv.tileLayers[i][j] = to
			}
		}
	}
}

func I18n(stringID, fallbackText string) string {
	translationString, has := langData[stringID]

	if !has {
		return fallbackText
	}

	return translationString
}

type IResource interface{}

type Rect struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type XPRunWindow struct {
	XPWindow
}

func NewXPRunWindow(xps *WinXPScreen) *XPRunWindow {
	s := new(XPRunWindow)
	s.hWnd = xps.hwndCounter
	xps.hwndCounter++
	s.xpScreen = xps
	s.game = xps.game

	return s
}

type XPStartButtonWidget struct {
	XPWindow
}

func (wnd *XPStartButtonWidget) Draw(screen *ebiten.Image) {
	wnd.game.DrawStatedNineGrid(screen, wnd.xpScreen.startButtonImage, 0, 3, 1.0, NineGridInfo{Left: 6, Top: 13, Right: 52, Bottom: 14}, wnd.posX, wnd.posY, wnd.width, wnd.height)
}

func NewXPStartButtonWidget(xps *WinXPScreen) *XPStartButtonWidget {
	s := new(XPStartButtonWidget)
	InitXPWindow(&s.XPWindow, xps, "_CiceronStart", 0, 0, 32, 32)
	return s
}

type ScreenManager struct {
	currentScreen IScreen
}

func (bs *ScreenManager) SetScreen(screen IScreen) {
	bs.currentScreen = screen
}

func CreateTextGrid() *TextGrid {
	tg := new(TextGrid)
	return tg
}

func (tg *TextGrid) Row(i int) []rune {
	return tg.rows[i]
}

func (tg *TextGrid) SetRow(i int, row []rune) {
	tg.rows[i] = row
}

type TileCursor struct {
	x int
	y int
}

func (c *TileCursor) GetWorldPos() Vec2f {
	return Vec2f{float64(c.x), float64(c.y)}.Scale(tileSize)
}

type TileDesc struct {
	JSONName string
	Walkable bool
}

type TileDescStorage map[Tile]TileDesc

func (tds *TileDescStorage) RegisterTile(tileID Tile, jsonName string, walkable bool) {
	// loadingLog := lazyAppend(loadingLog, fmt.Sprintf("Registering tile \"%s\"", jsonName))
	(*tds)[tileID] = TileDesc{
		JSONName: jsonName,
		Walkable: walkable,
	}
}

var tileDescStorage TileDescStorage

var tileNameMap map[Tile]string
var brandingImage *ebiten.Image

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	tileNameMap = make(map[Tile]string)

	tileDescStorage = make(TileDescStorage)

	tileDescStorage.RegisterTile(tileIDEmpty, "tile_id_empty", true)
	tileDescStorage.RegisterTile(tileIDGrass, "tile_id_grass", true)
	tileDescStorage.RegisterTile(tileIDSand, "tile_id_sand", true)
	tileDescStorage.RegisterTile(tileIDWater, "tile_id_water", true)
	tileDescStorage.RegisterTile(tileIDVoid, "tile_id_void", true)
	tileDescStorage.RegisterTile(tileIDHouseWall, "tile_id_house_wall", false)
	tileDescStorage.RegisterTile(tileIDCarvedStone, "tile_id_carved_stone", false)
	tileDescStorage.RegisterTile(tileIDBricks, "tile_id_bricks", false)
	tileDescStorage.RegisterTile(tileIDRock, "tile_id_rock", false)
	tileDescStorage.RegisterTile(tileIDDoor, "tile_id_door", false)
	tileDescStorage.RegisterTile(tileIDWallCornerL, "tile_id_wall_corner_l", true)
	tileDescStorage.RegisterTile(tileIDWallCornerR, "tile_id_wall_corner_r", true)
	tileDescStorage.RegisterTile(tileIDWindow, "tile_id_window", true)
	tileDescStorage.RegisterTile(tileIDBarrier, "tile_id_barrier", false)
	tileDescStorage.RegisterTile(tileIDSwitch, "tile_id_switch", true)
	tileDescStorage.RegisterTile(tileIDSwitchActive, "tile_id_switch_active", true)
	tileDescStorage.RegisterTile(tileIDRoofNW, "tile_id_roof_nw", false)
	tileDescStorage.RegisterTile(tileIDRoofTop, "tile_id_roof_top", false)
	tileDescStorage.RegisterTile(tileIDRoofNE, "tile_id_roof_ne", false)
	tileDescStorage.RegisterTile(tileIDRoofLeft, "tile_id_roof_left", false)
	tileDescStorage.RegisterTile(tileIDRoof, "tile_id_roof", false)
	tileDescStorage.RegisterTile(tileIDRoofRight, "tile_id_roof_right", false)
	tileDescStorage.RegisterTile(tileIDRoofSW, "tile_id_roof_sw", false)
	tileDescStorage.RegisterTile(tileIDRoofBottom, "tile_id_roof_bottom", false)
	tileDescStorage.RegisterTile(tileIDRoofSE, "tile_id_roof_se", false)
	tileDescStorage.RegisterTile(tileIDThorns, "tile_id_thorns", true)
	tileDescStorage.RegisterTile(tileIDThornsActive, "tile_id_thorns_active", true)
	tileDescStorage.RegisterTile(tileIDButton, "tile_id_button", true)
	tileDescStorage.RegisterTile(tileIDButtonPushed, "tile_id_button_pushed", true)
	tileDescStorage.RegisterTile(tileIDBlenderNW, "tile_id_blender_nw", true)
	tileDescStorage.RegisterTile(tileIDBlenderTop, "tile_id_blender_top", true)
	tileDescStorage.RegisterTile(tileIDBlenderNE, "tile_id_blender_ne", true)
	tileDescStorage.RegisterTile(tileIDBlenderLeft, "tile_id_blender_left", true)
	tileDescStorage.RegisterTile(tileIDBlenderRight, "tile_id_blender_right", true)
	tileDescStorage.RegisterTile(tileIDBlenderSW, "tile_id_blender_sw", true)
	tileDescStorage.RegisterTile(tileIDBlenderBottom, "tile_id_blender_bottom", true)
	tileDescStorage.RegisterTile(tileIDBlenderSE, "tile_id_blender_se", true)
	tileDescStorage.RegisterTile(tileIDTableBottom, "tile_id_table_bottom", false)
	tileDescStorage.RegisterTile(tileIDTableTop, "tile_id_table_top", false)
	tileDescStorage.RegisterTile(tileIDBedTop, "tile_id_bed_top", false)
	tileDescStorage.RegisterTile(tileIDBedBottom, "tile_id_bed_bottom", false)
	tileDescStorage.RegisterTile(tileIDLaptop, "tile_id_laptop", true)
	tileDescStorage.RegisterTile(tileIDHouseFloor, "tile_id_house_floor", true)
	tileDescStorage.RegisterTile(tileIDHouseWall, "tile_id_house_inner_wall", false)
	tileDescStorage.RegisterTile(tileIDMine, "tile_id_mine", true)
	tileDescStorage.RegisterTile(tileIDMineFlag, "tile_id_mine_flag", false)
	tileDescStorage.RegisterTile(tileIDBloodDigit1, "tile_id_blood_digit_1", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit2, "tile_id_blood_digit_2", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit3, "tile_id_blood_digit_3", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit4, "tile_id_blood_digit_4", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit5, "tile_id_blood_digit_5", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit6, "tile_id_blood_digit_6", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit7, "tile_id_blood_digit_7", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit8, "tile_id_blood_digit_8", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit9, "tile_id_blood_digit_9", true)
	tileDescStorage.RegisterTile(tileIDBloodDigit0, "tile_id_blood_digit_0", true)
	tileDescStorage.RegisterTile(tileIDSokobanBox, "tile_id_sokoban_box", false)
	tileDescStorage.RegisterTile(tileIDPoppingBarrierPushed, "tile_id_popping_barrier_pushed", true)
	tileDescStorage.RegisterTile(tileIDPoppingBarrier, "tile_id_popping_barrier", true)
	tileDescStorage.RegisterTile(tileIDPoppingBarrierActive, "tile_id_popping_barrier_active", false)
	tileDescStorage.RegisterTile(tileIDWeigth, "tile_id_weight", false)
	tileDescStorage.RegisterTile(tileIDArrowRight, "tile_id_arrow_right", true)
	tileDescStorage.RegisterTile(tileIDBerryBush, "tile_id_berry_bush", false)
	tileDescStorage.RegisterTile(tileIDBall, "tile_id_ball", false)
	tileDescStorage.RegisterTile(tileIDHoneyBall, "tile_id_honey_ball", true)
	tileDescStorage.RegisterTile(tileIDElevator, "tile_id_elevator", true)
	tileDescStorage.RegisterTile(tileIDTarget, "tile_id_target", true)
	tileDescStorage.RegisterTile(tileIDFenceSinleTop, "tile_id_fence_single_top", false)
	tileDescStorage.RegisterTile(tileIDPavedRoad, "tile_id_paved_road", true)
	tileDescStorage.RegisterTile(tileIDPit, "tile_id_pit", true)
	tileDescStorage.RegisterTile(tileIDTree1, "tile_id_tree_1", true)
	tileDescStorage.RegisterTile(tileIDBush, "tile_id_bush", false)
	tileDescStorage.RegisterTile(tileIDWell, "tile_id_well", false)
	tileDescStorage.RegisterTile(tileIDFenceNW, "tile_id_fence_nw", false)
	tileDescStorage.RegisterTile(tileIDPavedRoad2, "tile_id_paved_road_2", true)
	tileDescStorage.RegisterTile(tileIDCastleTower, "tile_id_castle_tower", false)
	tileDescStorage.RegisterTile(tileIDTree2, "tile_id_tree_2", true)
	tileDescStorage.RegisterTile(tileIDTree3, "tile_id_tree_3", true)
	tileDescStorage.RegisterTile(tileIDFlower1, "tile_id_flower_1", true)
	tileDescStorage.RegisterTile(tileIDSlab, "tile_id_slab", false)
	tileDescStorage.RegisterTile(tileIDArrowLeft, "tile_id_arrow_left", true)
	tileDescStorage.RegisterTile(tileIDAid, "tile_id_aid", true)
	tileDescStorage.RegisterTile(tileIDBlueRose, "tile_id_blue_rose", true)
	tileDescStorage.RegisterTile(tileIDColorDigit1, "tile_id_color_digit_1", true)
	tileDescStorage.RegisterTile(tileIDColorDigit2, "tile_id_color_digit_2", true)
	tileDescStorage.RegisterTile(tileIDColorDigit3, "tile_id_color_digit_3", true)
	tileDescStorage.RegisterTile(tileIDPoop, "tile_id_poop", true)

	tileNameMap = map[Tile]string{
		tileIDEmpty:                "Empty",
		tileIDGrass:                "Grass",
		tileIDSand:                 "Sand",
		tileIDWater:                "Water",
		tileIDVoid:                 "Void",
		tileIDHouseWall:            "House Wall",
		tileIDCarvedStone:          "Carved Stone",
		tileIDBricks:               "Bricks",
		tileIDRock:                 "Rock",
		tileIDDoor:                 "Door",
		tileIDWallCornerL:          "Corner Bricks (Left)",
		tileIDWallCornerR:          "Corner Bricks (Right)",
		tileIDWindow:               "Window",
		tileIDBarrier:              "Barrier",
		tileIDSwitch:               "Switch",
		tileIDSwitchActive:         "Pushed Switch",
		tileIDRoofNW:               "Roof (Top-Left)",
		tileIDRoofTop:              "Roof (Top)",
		tileIDRoofNE:               "Roof (Top-Right)",
		tileIDRoofLeft:             "Roof (Left)",
		tileIDRoof:                 "Roof (Center)",
		tileIDRoofRight:            "Roof (Right)",
		tileIDRoofSW:               "Roof (Bottom-Left)",
		tileIDRoofBottom:           "Roof (Bottom)",
		tileIDRoofSE:               "Roof (Bottom-Right)",
		tileIDThorns:               "Thorns",
		tileIDThornsActive:         "Activated Thorns",
		tileIDButton:               "Button",
		tileIDButtonPushed:         "Pushed Button",
		tileIDBlenderNW:            "Blender (Top-Left)",
		tileIDBlenderTop:           "Blender (Top)",
		tileIDBlenderNE:            "Blender (Top-Right)",
		tileIDBlenderLeft:          "Blender (Left)",
		tileIDBlenderRight:         "Blender (Right)",
		tileIDBlenderSW:            "Blender (Bottom-Left)",
		tileIDBlenderBottom:        "Blender (Bottom)",
		tileIDBlenderSE:            "Blender (Bottom-Right)",
		tileIDTableBottom:          "Table (Bottom)",
		tileIDTableTop:             "Table (Top)",
		tileIDBedTop:               "Bed (Top)",
		tileIDBedBottom:            "Bed (Bottom)",
		tileIDLaptop:               "Laptop",
		tileIDHouseFloor:           "Floor",
		tileIDHouseInnerWall:       "House Wall (Inner)",
		tileIDMine:                 "Mine",
		tileIDMineFlag:             "Mine Flag",
		tileIDBloodDigit1:          "Blood Digit 1",
		tileIDBloodDigit2:          "Blood Digit 2",
		tileIDBloodDigit3:          "Blood Digit 3",
		tileIDBloodDigit4:          "Blood Digit 4",
		tileIDBloodDigit5:          "Blood Ditit 5",
		tileIDBloodDigit6:          "Blood Digit 6",
		tileIDBloodDigit7:          "Blood Digit 7",
		tileIDBloodDigit8:          "Blood Digit 8",
		tileIDBloodDigit9:          "Blood Digit 9",
		tileIDBloodDigit0:          "Blood Digit 0",
		tileIDSokobanBox:           "Sokoban Box",
		tileIDPoppingBarrierPushed: "Popping Barrier (Pushed)",
		tileIDPoppingBarrier:       "Popping Barrier",
		tileIDPoppingBarrierActive: "Popping Barrier (Active)",
		tileIDWeigth:               "Weight",
		tileIDArrowRight:           "Arrow Right",
		tileIDBerryBush:            "Berried Bush",
		tileIDBall:                 "Ball",
		tileIDHoneyBall:            "Honey Ball",
		tileIDElevator:             "Elevator",
		tileIDTarget:               "Target",
		tileIDFenceSinleTop:        "Fence Single (Top)",
		tileIDPavedRoad:            "Paved Road",
		tileIDPit:                  "Pit",
		tileIDTree1:                "Tree 1",
		tileIDBush:                 "Bush",
		tileIDWell:                 "Well",
		tileIDFenceNW:              "Fence (NW)",
		tileIDPavedRoad2:           "Paved Road 2",
		tileIDCastleTower:          "Castle Tower",
		tileIDTree2:                "Tree 2",
		tileIDTree3:                "Tree 3",
		tileIDFlower1:              "Red Flower",
		tileIDSlab:                 "Slab",
		tileIDArrowLeft:            "Arrow Left",
		tileIDAid:                  "Aid",
		tileIDBlueRose:             "Blue Rose",
		tileIDColorDigit1:          "Color Digit 1",
		tileIDColorDigit2:          "Color Digit 2",
		tileIDColorDigit3:          "Color Digit 3",
		tileIDPoop:                 "Poop",
	}
}

func (g *Game) ToggleDebugInfoShow() {
	g.showDebugInfo = !g.showDebugInfo
}

func (g *Game) Update() error {
	if g.currentScreen != nil {
		g.currentScreen.Update()
	}

	g.appTicker++
	return nil
}

func GetTileSprite(tileSet *ebiten.Image, tileSetWidth int, tileSize int, tile Tile) *ebiten.Image {
	sx := (int(tile) % tileSetWidth) * tileSize
	sy := (int(tile) / tileSetWidth) * tileSize

	return tileSet.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image)
}

type IRunningMode interface {
	ProcessKeyEvents()
	Draw(*ebiten.Image)
}

type TextFormat struct {
	scale     float64
	textColor color.Color
	shadow    bool
}

type IGameplayMode interface {
	Draw(*ebiten.Image)
	ProcessKeyEvents() bool
	Update()
}

type GameplayMode struct {
	gameplayScreen *GameplayScreen
}

func (*GameplayMode) Draw(screen *ebiten.Image) {

}

func (*GameplayMode) Update() {

}

func (*GameplayMode) ProcessKeyEvents() bool {
	return true
}

type GameplayModeDefault struct {
	GameplayMode
}

func (mode *GameplayModeDefault) ProcessKeyEvents() bool {

	player := mode.gameplayScreen.game.char

	if inpututil.IsKeyJustReleased(keyBinds[kbPlayerMoveRight]) ||
		inpututil.IsKeyJustReleased(keyBinds[kbPlayerMoveLeft]) ||
		inpututil.IsKeyJustReleased(keyBinds[kbPlayerMoveUp]) ||
		inpututil.IsKeyJustReleased(keyBinds[kbPlayerMoveDown]) {
		player.EndWalk()
		return false
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbPlayerMoveRight]) {
		player.StartWalk(LooksRight)
		return false
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbPlayerMoveLeft]) {
		player.StartWalk(LooksLeft)
		return false
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbPlayerMoveUp]) {
		player.StartWalk(LooksUp)
		return false
	}

	if inpututil.IsKeyJustPressed(keyBinds[kbPlayerMoveDown]) {
		player.StartWalk(LooksDown)
		return false
	}

	return true
}

func NewGameplayModeDefault(s *GameplayScreen) *GameplayModeDefault {
	mode := new(GameplayModeDefault)
	mode.gameplayScreen = s

	game := mode.gameplayScreen.game

	game.camera.TargetEntity(game.char)

	return mode
}

type GameplayModeEntityFocusRotation struct {
	GameplayMode
}

func (mode *GameplayModeEntityFocusRotation) Draw(screen *ebiten.Image) {
	game := mode.gameplayScreen.game

	game.DrawModeTitle(screen, I18n("string_entity_focus_rotation", "Entity Focus Rotation"))
	game.DrawEntityInfo(screen, game.camera.targetEntity)
}

func (mode *GameplayModeEntityFocusRotation) Update() {
	if tickCounter%16 == 0 {
		game := mode.gameplayScreen.game

		if len(game.entities) == 1 {
			game.camera.TargetEntity(mode.gameplayScreen.game.entities[0])
		} else if len(game.entities) > 0 {
			game.camera.TargetEntity(game.entities[r.Intn(len(game.entities)-1)])
		}
	}
}
func (mode *GameplayModeEntityFocusRotation) ProcessKeyEvents() bool {
	if inpututil.IsKeyJustPressed(keyBinds[kbToggleEntityFocusRotation]) {
		mode.gameplayScreen.SetGameplayMode(NewGameplayModeDefault(mode.gameplayScreen))
		return false
	}

	return true
}

func NewGameplayModeEntityFocusRotation(s *GameplayScreen) *GameplayModeEntityFocusRotation {
	mode := new(GameplayModeEntityFocusRotation)
	mode.gameplayScreen = s

	return mode
}

// Delay the time before next loading log entry will append
func lazyAppend(target []string, el string) []string {
	target = append(target, el)

	// time.Sleep(time.Second)
	return target
}

func LoadImage(path string) *ebiten.Image {
	loadingLog = lazyAppend(loadingLog, "Loading image \""+path+"\"")

	image, _, _ := ebitenutil.NewImageFromFile(path)
	return image
}

var splash *ebiten.Image

// The main class that maintains global game state
//
// Implements ebiten.Game
type Game struct {
	level              *Level
	char               ICharacter
	entities           []ILivingEntity
	entityListMutex    sync.RWMutex
	view               View
	gameOver           bool
	fontRenderer       *FontRenderer
	systemFontRenderer *FontRenderer
	debugScreen        *DebugScreen
	showDebugInfo      bool
	camera             Camera
	ready              bool
	appTicker          int
	currentScreen      IScreen
	audioManager       *AudioManager
	volumeMusic        float64
}

func (g *Game) WorldPosToTilePos(worldX float64, worldY float64) (int, error) {
	x := int(worldX)
	y := int(worldY)

	return y/tileSize*g.level.width + x/tileSize, nil
}

func (g *Game) GetUnderlyingTilesAt(worldX float64, worldY float64) ([]Tile, error) {
	var tiles []Tile

	tilePos, err := g.WorldPosToTilePos(worldX, worldY)
	if err != nil {
		return nil, err
	}

	for _, layer := range g.level.tileLayers {
		if layer[tilePos] != tileIDEmpty {
			tiles = append(tiles, layer[tilePos])
		}
	}

	return tiles, err
}

func (g *Game) IsTileSolidAt(worldX float64, worldY float64) bool {
	tiles, _ := g.GetUnderlyingTilesAt(worldX, worldY)

	for _, underlyingTile := range tiles {
		if !tileDescStorage[underlyingTile].Walkable {
			return true
		}
	}

	return false
}

func (g *Game) GameOver() {
	g.gameOver = true
}

func (g *Game) Reset() {
	g.char = CreateCharacter(g)
	g.gameOver = false
}

func (g *Game) ProcessTileEffects(e ILivingEntity, tile Tile) {

	switch id := tile; id {

	case tileIDVoid:
		e.GetLivingEntity().health -= 0.1

	case tileIDWater:
		e.GetLivingEntity().SetSpeedModifier(0.25)

	case tileIDSwitch:
		tilePos, _ := e.GetTilePos()
		for _, layer := range g.level.tileLayers {
			if layer[tilePos] == tileIDSwitch {
				layer[tilePos] = tileIDSwitchActive

				for j, layer := range g.level.tileLayers {
					for k, tile := range layer {
						if tile == tileIDThornsActive {
							g.level.tileLayers[j][k] = tileIDThorns
						}
					}
				}
			}
		}

	case tileIDButton:
		tilePos, _ := e.GetTilePos()
		for _, layer := range g.level.tileLayers {
			if layer[tilePos] == tileIDButton {
				layer[tilePos] = tileIDButtonPushed
			}
		}

	case tileIDThornsActive:
		e.GetLivingEntity().health -= 10.0

	case tileIDLaptop:
		switch e.GetLivingEntity().etype.(type) {
		case Character:
			g.SetScreen(NewComputerScreen(g, func(g *Game) IScreen { return NewGameplayScreen(g) }))
		}
	}

}

var langData map[string]string

func (g *Game) Load() {
	loadingLog = append(loadingLog, "Loading Keybinds")

	keyBinds = KeyBindMap{
		kbPlayerMoveRight:           ebiten.KeyArrowRight,
		kbPlayerMoveLeft:            ebiten.KeyArrowLeft,
		kbPlayerMoveUp:              ebiten.KeyArrowUp,
		kbPlayerMoveDown:            ebiten.KeyArrowDown,
		kbToggleEntityFocusRotation: ebiten.KeyF7,
		kbToggleEditMode:            ebiten.KeyF8,
		kbEditorNextLayer:           ebiten.KeyX,
		kbEditorPrevLayer:           ebiten.KeyZ,
		kbEditorNextBrush:           ebiten.KeyNumpadAdd,
		kbEditorPrevBrush:           ebiten.KeyNumpadSubtract,
		kbEditorPlace:               ebiten.KeyEnter,
		kbEditorDelete:              ebiten.KeyDelete,
		kbEditorSwitchMode:          ebiten.KeyS,
		kbShowDebugInfo:             ebiten.KeyF3,
		kbWorldZoomOut:              ebiten.KeyO,
		kbWorldZoomIn:               ebiten.KeyP,
	}

	langFile, _ := ioutil.ReadFile("lang/ru_ru.json")

	langData = make(map[string]string)
	_ = json.Unmarshal([]byte(langFile), &langData)

	g.audioManager = NewAudioManager(g)

	g.volumeMusic = 0.5
	g.audioManager.Load("bgm/stage_prepare", "sound/prepare.ogg")
	g.audioManager.Load("bgm/main_menu", "sound/bgm_main_menu.ogg")
	g.audioManager.Load("bgm/level0", "sound/bgm_level0.ogg")
	g.audioManager.Load("bgm/computer", "sound/computer.ogg")
	g.audioManager.Load("winxp/critical_stop", "sound/critical_stop.ogg")

	am := AssetManager_GetInstance()
	am.Load("game/tile_map", "assets/tilemap2.png")
	am.Load("game/character/player", "assets/char2.png")
	am.Load("game/character/morgen", "assets/morgen.png")
	am.Load("game/character/michael", "assets/michael.png")
	am.Load("game/character/flan", "assets/flan.png")
	am.Load("game/character/monobear", "assets/monobear.png")
	am.Load("gui/frame", "assets/gui_frame.png")
	am.Load("gui/frame_herb", "assets/gui_frame_herb.png")
	am.Load("gui/button", "assets/gui_button.png")
	am.Load("winxp/boot_logo", "assets/boot.png")

	tilesImage = LoadImage("assets/tilemap2.png")
	charSprite = LoadImage("assets/char2.png")
	morgenSprite = LoadImage("assets/morgen.png")
	michaelSprite = LoadImage("assets/michael.png")
	flanSprite = LoadImage("assets/flan.png")
	monobearSprite = LoadImage("assets/monobear.png")
	guiFrameTest = LoadImage("assets/gui_frame_test.png")
	guiFrameHerb = LoadImage("assets/gui_frame_herb.png")
	guiButton = LoadImage("assets/gui_button.png")
	tileCursor = LoadImage("assets/tile_selector.png")
	explosionSprite = LoadImage("assets/explosion.png")
	skillMonobearExplosion = LoadImage("assets/spell_monobear_explosion.png")
	seeYaTileSet = LoadImage("assets/seeya.png")
	worldBorderImage = LoadImage("assets/world_border.png")
	xpCaption = LoadImage("assets/computer/frame_caption.png")
	xpLeftFrame = LoadImage("assets/computer/frame_left.png")
	xpRightFrame = LoadImage("assets/computer/frame_right.png")
	xpBottomFrame = LoadImage("assets/computer/frame_bottom.png")
	xpCloseButton = LoadImage("assets/computer/close_button.png")
	xpCloseGlyph = LoadImage("assets/computer/close_glyph.png")
	xpIconError = LoadImage("assets/computer/icon_error.png")
	xpButton = LoadImage("assets/computer/button.png")

	loadingLog = lazyAppend(loadingLog, "Loading font_fantasy")
	g.fontRenderer = NewFontRenderer()

	loadingLog = lazyAppend(loadingLog, "Loading level")
	g.LoadLevel("level/level0.lvl")

	loadingLog = lazyAppend(loadingLog, "Setting camera zoom")
	g.camera.SetZoom(4.0)

	loadingLog = lazyAppend(loadingLog, "Setting GUI scale")
	g.view.guiScale = 2.0

	loadingLog = lazyAppend(loadingLog, "Creating character")
	g.char = CreateCharacter(g)

	loadingLog = lazyAppend(loadingLog, "Appending character to entity list")
	g.entities = append(g.entities, g.char)

	for i := 0; i < 8; i++ {
		g.entities = append(g.entities, CreateMichael(g))
	}

	loadingLog = lazyAppend(loadingLog, "Creating Debug Screen")
	g.debugScreen = CreateDebugScreen(g)

	loadingLog = lazyAppend(loadingLog, "Ready")
	g.ready = true
}

func (g *Game) SetScreen(screen IScreen) {
	if g.currentScreen != nil {
		g.currentScreen.OnDetach()
	}

	if screen != nil {
		screen.OnAttach()
	}

	g.currentScreen = screen
}

const (
	buttonStateNormal = iota
	buttonStateHover
	buttonStatePressed
)

func (g *Game) DrawStatedNineGrid(screen *ebiten.Image,
	img *ebiten.Image,
	state int,
	nStates int,
	scale float64,
	gridInfo NineGridInfo,
	x float64,
	y float64,
	width float64,
	height float64) {

	atlasWidth, atlasHeight := img.Size()

	stateWidth := atlasWidth
	stateHeight := atlasHeight / nStates
	sx := 0
	sy := stateHeight * state

	stateImage := img.SubImage(image.Rect(sx, sy, sx+stateWidth, sy+stateHeight)).(*ebiten.Image)
	g.DrawNineGrid(screen, stateImage, scale, gridInfo, x, y, width, height)
}

func (g *Game) DrawGUIButton(screen *ebiten.Image, state int, x float64, y float64, width float64, height float64) {

	am := AssetManager_GetInstance()
	atlasWidth, atlasHeight := am.Get("gui/button").Size()

	buttonWidth := atlasWidth
	buttonHeight := atlasHeight / 3
	// sx := 0
	sy := buttonHeight * state

	buttonImage := am.Get("gui/button").SubImage(image.Rect(0, sy, buttonWidth, sy+buttonHeight)).(*ebiten.Image)

	g.DrawNineGrid(screen, buttonImage, g.view.guiScale, NineGridInfo{Right: 5, Top: 5, Left: 5, Bottom: 5}, x, y, width, height)
}

type NineGridInfo struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

const (
	nineGridNW     = 0b00001001
	nineGridTop    = 0b00001010
	nineGridNE     = 0b00001100
	nineGridLeft   = 0b00010001
	nineGridCenter = 0b00010010
	nineGridRight  = 0b00010100
	nineGridSW     = 0b00100001
	nineGridBottom = 0b00100010
	nineGridSE     = 0b00100100
)

const (
	nineGridColumnLeft   = 1
	nineGridColumnCenter = 2
	nineGridColumnRight  = 4
	nineGridRowTop       = 8
	nineGridRowCenter    = 16
	nineGridRowBottom    = 32
)

func NineGridGetTile_0(img *ebiten.Image, gridTile int, gridInfo NineGridInfo) *ebiten.Image {
	var x, y, width, height int

	atlasWidth, atlasHeight := img.Size()

	bounds := Recti{
		Vec2i{gridInfo.Left, gridInfo.Top},
		Vec2i{gridInfo.Right, gridInfo.Bottom},
	}

	switch gridTile {
	case nineGridNW:
		x = 0
		y = 0
		width = bounds.p1.X
		height = bounds.p1.Y

	case nineGridTop:
		x = bounds.p1.X
		y = 0
		width = atlasWidth - (bounds.p1.X + bounds.p2.X)
		height = bounds.p1.Y

	case nineGridNE:
		x = atlasWidth - bounds.p2.X
		y = 0
		width = bounds.p2.X
		height = bounds.p1.Y

	case nineGridLeft:
		x = 0
		y = bounds.p1.Y
		width = bounds.p1.X
		height = atlasHeight - (bounds.p1.Y + bounds.p2.Y)

	case nineGridCenter:
		x = bounds.p1.X
		y = bounds.p1.Y
		width = atlasWidth - (bounds.p1.X + bounds.p2.X)
		height = atlasHeight - (bounds.p1.Y + bounds.p2.Y)

	case nineGridRight:
		x = atlasWidth - bounds.p2.X
		y = bounds.p1.Y
		width = bounds.p2.X
		height = atlasHeight - (bounds.p1.Y + bounds.p2.Y)

	case nineGridSW:
		x = 0
		y = atlasHeight - bounds.p2.Y
		width = bounds.p1.X
		height = bounds.p2.Y

	case nineGridBottom:
		x = bounds.p1.X
		y = atlasHeight - bounds.p2.Y
		width = atlasWidth - (bounds.p1.X + bounds.p2.X)
		height = bounds.p2.Y

	case nineGridSE:
		x = atlasWidth - bounds.p2.X
		y = atlasHeight - bounds.p2.Y
		width = bounds.p2.X
		height = bounds.p2.Y
	}

	min := img.Bounds().Min
	return img.SubImage(image.Rect(
		min.X+x,
		min.Y+y,
		min.X+x+width,
		min.Y+y+height)).(*ebiten.Image)
}

func NineGridGetTile(img *ebiten.Image, gridTile int, gridInfo NineGridInfo) *ebiten.Image {
	var x, y, width, height int
	atlasWidth, atlasHeight := img.Size()

	bounds := Recti{
		Vec2i{gridInfo.Left, gridInfo.Top},
		Vec2i{gridInfo.Right, gridInfo.Bottom},
	}

	if gridTile&nineGridColumnLeft != 0 {
		width = bounds.p1.X
	} else if gridTile&nineGridColumnCenter != 0 {
		x = bounds.p1.X
		width = atlasWidth - (bounds.p1.X + bounds.p2.X)
	} else if gridTile&nineGridColumnRight != 0 {
		x = atlasWidth - bounds.p2.X
		width = bounds.p2.X
	}

	if gridTile&nineGridRowTop != 0 {
		height = bounds.p1.Y
	} else if gridTile&nineGridRowCenter != 0 {
		y = bounds.p1.Y
		height = atlasHeight - (bounds.p1.Y + bounds.p2.Y)
	} else if gridTile&nineGridRowBottom != 0 {
		y = atlasHeight - bounds.p2.Y
		height = bounds.p2.Y
	}

	min := img.Bounds().Min
	return img.SubImage(image.Rect(
		min.X+x,
		min.Y+y,
		min.X+x+width,
		min.Y+y+height)).(*ebiten.Image)
}

func NineGridGetTile2(gridTile int, gridInfo *NineGridInfo2) *ebiten.Image {
	var x, y, width, height int

	img := gridInfo.sprite
	atlasWidth, atlasHeight := img.Size()

	bounds := gridInfo.renderBounds

	if gridTile&nineGridColumnLeft != 0 {
		width = bounds.p1.X
	} else if gridTile&nineGridColumnCenter != 0 {
		x = bounds.p1.X
		width = atlasWidth - (bounds.p1.X + bounds.p2.X)
	} else if gridTile&nineGridColumnRight != 0 {
		x = atlasWidth - bounds.p2.X
		width = bounds.p2.X
	}

	if gridTile&nineGridRowTop != 0 {
		height = bounds.p1.Y
	} else if gridTile&nineGridRowCenter != 0 {
		y = bounds.p1.Y
		height = atlasHeight - (bounds.p1.Y + bounds.p2.Y)
	} else if gridTile&nineGridRowBottom != 0 {
		y = atlasHeight - bounds.p2.Y
		height = bounds.p2.Y
	}

	min := img.Bounds().Min
	return img.SubImage(image.Rect(
		min.X+x,
		min.Y+y,
		min.X+x+width,
		min.Y+y+height)).(*ebiten.Image)
}

const (
	nineGridModeScale = iota
	nineGridModeRepeat
)

type NineGridInfo2 struct {
	actualBounds Recti
	renderBounds Recti
	mode         int
	sprite       *ebiten.Image
	tileCache    [9]*ebiten.Image
}

func InitializeNineGrid2(ng *NineGridInfo2) {
	ng.tileCache = [9]*ebiten.Image{
		NineGridGetTile2(nineGridNW, ng),
		NineGridGetTile2(nineGridTop, ng),
		NineGridGetTile2(nineGridNE, ng),
		NineGridGetTile2(nineGridLeft, ng),
		NineGridGetTile2(nineGridCenter, ng),
		NineGridGetTile2(nineGridRight, ng),
		NineGridGetTile2(nineGridSW, ng),
		NineGridGetTile2(nineGridBottom, ng),
		NineGridGetTile2(nineGridSE, ng),
	}
}

func DrawSimpleRepeatedTexture(screen *ebiten.Image, img *ebiten.Image, scale, x, y, w, h float64) {

	var sx, sy, sw, sh float64

	sx = float64(img.Bounds().Min.X)
	sy = float64(img.Bounds().Min.Y)
	sw = w / scale
	sh = h / scale

	vs := []ebiten.Vertex{
		{
			DstX:   float32(x),
			DstY:   float32(y),
			SrcX:   float32(sx),
			SrcY:   float32(sy),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(x + w),
			DstY:   float32(y),
			SrcX:   float32(sx + sw),
			SrcY:   float32(sy),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(x),
			DstY:   float32(y + h),
			SrcX:   float32(sx),
			SrcY:   float32(sy + sh),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(x + w),
			DstY:   float32(y + h),
			SrcX:   float32(sx + sw),
			SrcY:   float32(sy + sh),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
	}

	triOp := &ebiten.DrawTrianglesOptions{}
	triOp.Address = ebiten.AddressRepeat
	screen.DrawTriangles(vs, []uint16{0, 1, 2, 1, 2, 3}, img, triOp)
}

func (g *Game) DrawNineGridRepeat2(screen *ebiten.Image, ng *NineGridInfo2, scale float64, offx float64, offy float64, width float64, height float64) {
	var x, y, w, h float64

	for i := 0; i < 9; i++ {

		x = offx
		y = offy

		switch i % 3 {
		case 0:
			w = float64(ng.renderBounds.p1.X) * scale
		case 1:
			x += float64(ng.renderBounds.p1.X) * scale
			w = float64(width) - float64(ng.renderBounds.p1.X+ng.renderBounds.p2.X)*scale
		case 2:
			x += float64(width) - float64(ng.renderBounds.p2.X)*scale
			w = float64(ng.renderBounds.p2.X) * scale
		}

		switch i / 3 {
		case 0:
			h = float64(ng.renderBounds.p1.Y) * scale
		case 1:
			y += float64(ng.renderBounds.p1.Y) * scale
			h = float64(height) - float64(ng.renderBounds.p1.Y+ng.renderBounds.p2.Y)*scale
		case 2:
			y += float64(height) - float64(ng.renderBounds.p2.Y)*scale
			h = float64(ng.renderBounds.p2.Y) * scale
		}

		DrawSimpleRepeatedTexture(
			screen,
			ng.tileCache[i],
			g.view.guiScale,
			x,
			y,
			w,
			h)
	}
}

func (g *Game) DrawNineGrid2Adjustable(screen *ebiten.Image, ng *NineGridInfo2, scale float64, offx float64, offy float64, width float64, height float64) {
	var x, y, w, h float64

	for i := 0; i < 9; i++ {

		x = offx
		y = offy

		switch i % 3 {
		case 0:
			x -= float64(ng.actualBounds.p1.X) * scale
			w = float64(ng.renderBounds.p1.X) * scale
		case 1:
			x += float64(ng.renderBounds.p1.X-ng.actualBounds.p1.X) * scale
			w = float64(width) - float64((ng.renderBounds.p1.X-ng.actualBounds.p1.X)+(ng.renderBounds.p2.X-ng.actualBounds.p2.X))*scale
		case 2:
			x += float64(width)
			w = float64(ng.renderBounds.p2.X) * scale
		}

		switch i / 3 {
		case 0:
			y -= float64(ng.actualBounds.p1.Y) * scale
			h = float64(ng.renderBounds.p1.Y) * scale
		case 1:
			y += float64(ng.renderBounds.p1.Y-ng.actualBounds.p1.Y) * scale
			h = float64(height) - float64((ng.renderBounds.p1.Y-ng.actualBounds.p1.Y)+(ng.renderBounds.p2.Y-ng.actualBounds.p2.Y))*scale
		case 2:
			y += float64(height)
			h = float64(ng.renderBounds.p2.Y) * scale
		}

		DrawSimpleRepeatedTexture(
			screen,
			ng.tileCache[i],
			g.view.guiScale,
			x,
			y,
			w,
			h)
	}
}

func (g *Game) DrawNineGrid2RepeatOuter(screen *ebiten.Image, ng *NineGridInfo2, scale float64, offx float64, offy float64, width float64, height float64) {
	var x, y, w, h float64

	for i := 0; i < 9; i++ {

		x = offx
		y = offy

		switch i % 3 {
		case 0:
			x -= float64(ng.renderBounds.p1.X) * scale
			w = float64(ng.renderBounds.p1.X) * scale
		case 1:
			w = float64(width)
		case 2:
			x += float64(width)
			w = float64(ng.renderBounds.p2.X) * scale
		}

		switch i / 3 {
		case 0:
			y -= float64(ng.renderBounds.p1.Y) * scale
			h = float64(ng.renderBounds.p1.Y) * scale
		case 1:
			h = float64(height)
		case 2:
			y += float64(height)
			h = float64(ng.renderBounds.p2.Y) * scale
		}

		DrawSimpleRepeatedTexture(
			screen,
			ng.tileCache[i],
			g.view.guiScale,
			x,
			y,
			w,
			h)
	}
}

func (g *Game) DrawNineGrid2(screen *ebiten.Image, ng *NineGridInfo2, scale float64, offx float64, offy float64, width float64, height float64) {
	op := &ebiten.DrawImageOptions{}

	atlasWidth, atlasHeight := ng.sprite.Size()
	renderBounds := ng.renderBounds
	cm := renderBounds.p1.Add(renderBounds.p2).ToVec2f()
	am := ng.actualBounds.p1.Add(ng.actualBounds.p2).ToVec2f()
	stretchW := (float64(width) - (cm.X-am.X)*scale) / (float64(atlasWidth) - cm.X)
	stretchH := (float64(height) - (cm.Y-am.Y)*scale) / (float64(atlasHeight) - cm.Y)

	for i := 0; i < 9; i++ {
		x := float64(offx) - float64(ng.actualBounds.p1.X)*scale
		y := float64(offy) - float64(ng.actualBounds.p1.Y)*scale
		sw := scale
		sh := scale

		switch i % 3 {
		case 1:
			x += float64(renderBounds.p1.X) * scale
			sw = stretchW
		case 2:
			x += (float64(width) + float64(ng.actualBounds.p1.X+ng.actualBounds.p2.X)*scale) - float64(renderBounds.p2.X)*scale
		}

		switch i / 3 {
		case 1:
			y += float64(renderBounds.p1.Y) * scale
			sh = stretchH
		case 2:
			y += (float64(height) + float64(ng.actualBounds.p1.X+ng.actualBounds.p2.Y)*scale) - float64(renderBounds.p2.Y)*scale
		}

		op.GeoM.Reset()
		op.GeoM.Scale(sw, sh)
		op.GeoM.Translate(x, y)
		screen.DrawImage(ng.tileCache[i], op)
	}
}

func (g *Game) DrawNineGrid(screen *ebiten.Image, img *ebiten.Image, scale float64, inf NineGridInfo, x float64, y float64, width float64, height float64) {
	cornerNW := NineGridGetTile(img, nineGridNW, inf)
	cornerTop := NineGridGetTile(img, nineGridTop, inf)
	cornerNE := NineGridGetTile(img, nineGridNE, inf)
	cornerLeft := NineGridGetTile(img, nineGridLeft, inf)
	cornerCenter := NineGridGetTile(img, nineGridCenter, inf)
	cornerRight := NineGridGetTile(img, nineGridRight, inf)
	cornerSW := NineGridGetTile(img, nineGridSW, inf)
	cornerBottom := NineGridGetTile(img, nineGridBottom, inf)
	cornerSE := NineGridGetTile(img, nineGridSE, inf)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	screen.DrawImage(cornerNW, op)

	op.GeoM.Reset()
	cornerTopWidth, _ := cornerTop.Size()
	topTileWidth := (width - float64(inf.Left+inf.Right)*scale) / float64(cornerTopWidth)
	op.GeoM.Scale(topTileWidth, scale)
	op.GeoM.Translate(x+float64(inf.Left)*scale, y)
	screen.DrawImage(cornerTop, op)

	op.GeoM.Reset()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(width-float64(inf.Right)*scale), y)
	screen.DrawImage(cornerNE, op)

	_, cornerLeftHeight := cornerLeft.Size()
	op.GeoM.Reset()
	leftTileHeight := (height - float64(inf.Top+inf.Bottom)*scale) / float64(cornerLeftHeight)
	op.GeoM.Scale(scale, leftTileHeight)
	op.GeoM.Translate(x, y+float64(inf.Top)*scale)
	screen.DrawImage(cornerLeft, op)

	op.GeoM.Reset()
	op.GeoM.Scale(topTileWidth, leftTileHeight)
	op.GeoM.Translate(x+float64(inf.Left)*scale, y+float64(inf.Top)*scale)
	screen.DrawImage(cornerCenter, op)

	op.GeoM.Reset()
	op.GeoM.Scale(scale, leftTileHeight)
	op.GeoM.Translate(x+(width-float64(inf.Right)*scale), y+float64(inf.Top)*scale)
	screen.DrawImage(cornerRight, op)

	op.GeoM.Reset()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y+(height-float64(inf.Bottom)*scale))
	screen.DrawImage(cornerSW, op)

	op.GeoM.Reset()
	op.GeoM.Scale(topTileWidth, scale)
	op.GeoM.Translate(x+float64(inf.Left)*scale, y+(height-float64(inf.Bottom)*scale))
	screen.DrawImage(cornerBottom, op)

	op.GeoM.Reset()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(width-float64(inf.Right)*scale), y+(height-float64(inf.Bottom)*scale))
	screen.DrawImage(cornerSE, op)
}

var nineGridGUIFrame *NineGridInfo2

func (g *Game) DrawGUIFrame(screen *ebiten.Image, x float64, y float64, width float64, height float64) {
	if nineGridGUIFrame == nil {
		nineGridGUIFrame = &NineGridInfo2{
			actualBounds: Recti{p1: Vec2i{X: 5, Y: 5}, p2: Vec2i{X: 5, Y: 5}},
			renderBounds: Recti{p1: Vec2i{X: 5, Y: 5}, p2: Vec2i{X: 5, Y: 5}},
			sprite:       guiFrameTest,
			mode:         nineGridModeRepeat,
		}

		InitializeNineGrid2(nineGridGUIFrame)
	}

	g.DrawNineGrid2RepeatOuter(screen, nineGridGUIFrame, g.view.guiScale, x, y, width, height)
}

var guiFrameHerbGridInfo *NineGridInfo2

func (g *Game) DrawHerbGUIFrame(screen *ebiten.Image, x float64, y float64, width float64, height float64) {
	if guiFrameHerbGridInfo == nil {
		rm := ResourceManager_GetInstance()

		guiFrameHerbGridInfo = &NineGridInfo2{
			actualBounds: Recti{p1: Vec2i{X: 9, Y: 19}, p2: Vec2i{X: 5, Y: 11}},
			renderBounds: Recti{p1: Vec2i{X: 26, Y: 24}, p2: Vec2i{X: 5, Y: 11}},
			mode:         nineGridModeRepeat,
			sprite:       rm.LoadImage("assets/gui/gui_frame_herb.png"),
			tileCache:    [9]*ebiten.Image{},
		}

		InitializeNineGrid2(guiFrameHerbGridInfo)
	}

	g.DrawNineGrid2Adjustable(
		screen,
		guiFrameHerbGridInfo,
		g.view.guiScale,
		x,
		y,
		width,
		height)
}

/*
func (g *Game) DrawHerbGUIFrame(screen *ebiten.Image, x float64, y float64, width float64, height float64) {
	if guiFrameHerbGridInfo == nil {
		guiFrameHerbGridInfo = &NineGridInfo2{
			actualBounds: Recti{p1: Vec2i{X: 16, Y: 16}, p2: Vec2i{X: 16, Y: 16}},
			renderBounds: Recti{p1: Vec2i{X: 8, Y: 8}, p2: Vec2i{X: 8, Y: 8}},
			mode:         nineGridModeRepeat,
			sprite:       guiFrameHerb,
			tileCache:    [9]*ebiten.Image{},
		}

		InitializeNineGrid2(guiFrameHerbGridInfo)
	}

	g.DrawNineGrid2Adjustable(
		screen,
		guiFrameHerbGridInfo,
		g.view.guiScale,
		x,
		y,
		width,
		height)
}
*/

/*
func (g *Game) DrawTextPopup(screen *ebiten.Image, text string, x int, y int, width int, height int) {
	g.DrawHerbGUIFrame(screen, float64(x), float64(y), float64(width), float64(height))
	g.fontRenderer.DrawTextRect(screen, text, g.view.guiScale, float64(x)+7*g.view.guiScale, float64(y)+7*g.view.guiScale, width-14, height-14)
}
*/

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScreen != nil {
		g.currentScreen.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) LoadLevel(path string) error {
	fd, err := os.Open(path)
	if fd != nil {
		defer fd.Close()
	}

	if err != nil {
		return err
	}

	level := new(Level)
	level.fileName = path
	level.width = 15
	level.height = 15

	for {
		byteLayer := make([]byte, 225)

		_, err := fd.Read(byteLayer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		tileLayer := make([]Tile, 225)

		for i := 0; i < len(byteLayer); i++ {
			tileLayer[i] = Tile(byteLayer[i])
		}

		level.tileLayers = append(level.tileLayers, tileLayer)
	}

	g.level = level

	return nil
}

func (g *Game) Exit() {
	g.SetScreen(NewFarewellScreen(g))
}

func (g *Game) SaveLevel(path string) {
	log.Println("Writing level file " + path)

	level := g.level

	buffer := make([]byte, len(level.tileLayers)*len(level.tileLayers[0]))
	for i, layer := range level.tileLayers {
		for j, tile := range layer {
			buffer[i*len(level.tileLayers[0])+j] = byte(tile)
		}
	}

	_ = os.WriteFile(path, buffer, 0644)

	/*
		fd, _ := os.Open(path)
		defer fd.Close()

		fd.Write(buffer)
	*/
}
func (g *Game) GetUnderlyingTilesAtTilePos(tilePos int) []*Tile {
	var tiles []*Tile

	for _, layer := range g.level.tileLayers {
		tiles = append(tiles, &layer[tilePos])
	}

	return tiles
}

func (g *Game) ProcessTileLeaving(e ILivingEntity, tilePos int) {
	tiles := g.GetUnderlyingTilesAtTilePos(tilePos)

	for _, tile := range tiles {
		switch *tile {
		case tileIDThorns:
			*tile = tileIDThornsActive

			g.level.ReplaceAll(tileIDSwitchActive, tileIDSwitch)
		}
	}
}

func (g *Game) DrawModeTitle(screen *ebiten.Image, text string) {
	fontRenderer := g.fontRenderer

	fontRenderer.PushState()
	fontRenderer.SetScale(g.view.guiScale)
	fontRenderer.SetTextColor(alpacolor.Red)

	textDim := g.fontRenderer.GetStringDimensions(text)

	glyphSize := fontRenderer.GetGlyphSize()
	frameSize := textDim.Add(glyphSize.Scale(2.0))

	g.DrawHerbGUIFrame(screen, 0, 0, frameSize.X, frameSize.Y)
	g.fontRenderer.DrawTextAt(screen, text, glyphSize)

	fontRenderer.PopState()
}

func (g *Game) DrawEntityInfo(screen *ebiten.Image, e ILivingEntity) {
	g.DrawHerbGUIFrame(screen,
		float64((screenWidth-(screenWidth-128))/2),
		float64(screenHeight-128),
		screenWidth-128,
		128)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(4.0, 4.0)
	op.GeoM.Translate(float64((screenWidth-(screenWidth-128))/2)+16, float64(screenHeight-128)+16)

	sprite := e.GetLivingEntity().sprite.SubImage(image.Rect(0, tileSize*3, tileSize, tileSize*4)).(*ebiten.Image)
	screen.DrawImage(sprite, op)

	fontRenderer := g.fontRenderer

	fontRenderer.PushState()
	fontRenderer.SetScale(g.view.guiScale)

	g.fontRenderer.DrawTextAt(screen,
		g.camera.targetEntity.GetLivingEntity().entityClass,
		Vec2f{float64((screenWidth-(screenWidth-128))/2) + tileSize*4.0 + 32,
			float64(screenHeight-128) + 32})

	fontRenderer.PopState()

	healthFactor := g.camera.targetEntity.GetLivingEntity().GetHealth() / 100.0

	if healthFactor > 0 {
		ebitenutil.DrawRect(screen,
			float64((screenWidth-(screenWidth-128))/2)+tileSize*4.0+32,
			float64(screenHeight-128)+56,
			128,
			16,
			color.RGBA{0, 0, 0, 255})

		ebitenutil.DrawRect(screen,
			float64((screenWidth-(screenWidth-128))/2)+tileSize*4.0+36,
			float64(screenHeight-128)+60,
			healthFactor*(128-8),
			8,
			HSL{math.Pow(healthFactor, 2) * 0.33, 1.0, 0.5}.ToRGB())
	} else {
		fontRenderer.PushState()
		fontRenderer.SetScale(g.view.guiScale)
		fontRenderer.SetTextColor(alpacolor.Red)

		g.fontRenderer.DrawTextAt(screen,
			"Thank you for watching!",
			Vec2f{float64((screenWidth-(screenWidth-128))/2) + tileSize*4.0 + 32,
				float64(screenHeight-128) + 56})

		fontRenderer.PopState()
	}
}

func (g *Game) DrawWorld(screen *ebiten.Image) {
	DrawSimpleRepeatedTexture(screen, worldBorderImage, 1.0, 0, 0, float64(screenWidth), float64(screenHeight))

	for _, l := range g.level.tileLayers {
		for i, t := range l {
			if t == tileIDEmpty {
				continue
			}

			op := &ebiten.DrawImageOptions{}

			cameraZoom := g.camera.GetZoom()

			tileWorldPos := Vec2f{
				float64(i % g.level.width * tileSize),
				float64(i / g.level.width * tileSize),
			}
			pos := g.camera.WorldToScreen2(tileWorldPos)

			tile := GetTileSprite(tilesImage, tileXNum, tileSize, t)

			op.GeoM.Scale(cameraZoom, cameraZoom)
			op.GeoM.Translate(pos.X, pos.Y)
			screen.DrawImage(tile, op)
		}
	}
}

type NextScreenBuilder func(*Game) IScreen

func CreateGame() (*Game, error) {
	g := new(Game)

	g.fontRenderer = NewFontRenderer()
	g.view.guiScale = 2.0

	if directScreenSet {
		screenBuilder, has := screenNames[directScreenName]
		if has {
			g.SetScreen(screenBuilder(g))
		}
	} else {
		g.systemFontRenderer = NewFontRenderer()
		brandingImage = LoadImage("assets/aragajaga.png")
		splash = LoadImage("assets/splash.png")

		go g.Load()

		g.SetScreen(NewBrandingScreen(g, func(g *Game) IScreen {
			if g.ready {
				return CreateMainMenu(g)
			} else {
				return NewLoadingScreen(g)
			}
		}))
	}

	return g, nil
}

var g_audioContext *audio.Context

func GetAudioContext() *audio.Context {
	if g_audioContext == nil {
		g_audioContext = audio.NewContext(44100)
	}

	return g_audioContext
}

var directScreenSet bool
var directScreenName string

type ScreenBuilder func(*Game) IScreen

var screenNames map[string]ScreenBuilder

var g_currentProcessCmdLine *CommandLine

type SwitchMap map[string]string

type CommandLine struct {
	switches   SwitchMap
	argv       []string
	begin_args uint
}

func (cmdLine *CommandLine) Init() bool {
	if g_currentProcessCmdLine == nil {
		return false
	}

	g_currentProcessCmdLine = new(CommandLine)
	cmdLine.FromOsArgs(os.Args)
	return true
}

type AnyMap map[interface{}]interface{}

func (cmdLine *CommandLine) FromOsArgs(argv []string) {
	cmdLine.argv = []string{""}

	// Clear the switch map
	for k := range cmdLine.switches {
		delete(cmdLine.switches, k)
	}

	cmdLine.begin_args = 1

	if len(argv) == 0 {
		cmdLine.SetProgram("")
	} else {
		cmdLine.SetProgram(argv[0])
	}
}

func (cmdLine *CommandLine) GetSwitches() *SwitchMap {
	return &cmdLine.switches
}

func (cmdLine *CommandLine) SetProgram(program string) {
	cmdLine.argv[0] = strings.TrimSpace(program)
}

var switchPrefixes []string = []string{"--", "-", "/"}

func GetSwitchPrefixLen(arg string) int {
	for _, prefix := range switchPrefixes {
		if substr(arg, 0, strlen(prefix)) == "prefix" {
			return strlen(prefix)
		}
	}

	return 0
}

var switchValueSeparator string = "="

/*
 *  Fill in |switch_string| and |switch_value| if |string| is a switch.
 *	This will preserve the input switch prefix in the output |switch_string|.
 */
func IsSwitch(arg string) (bool, string, string) {
	prefixLen := GetSwitchPrefixLen(arg)
	if prefixLen == 0 || prefixLen == strlen(arg) {
		return false, "", ""
	}

	pair := strings.SplitN(arg, switchValueSeparator, 2)
	return true, pair[0], pair[1]
}

func IsSwitchWithKey(arg string, key string) bool {
	prefixLen := GetSwitchPrefixLen(arg)

	if prefixLen == 0 || prefixLen == strlen(arg) {
		return false
	}

	equalsPos := strings.Index(arg, switchValueSeparator)
	return substr(arg, prefixLen, equalsPos-prefixLen) == key
}

func RegisterScreenName(name string, builder func(g *Game) IScreen) {
	screenNames[name] = builder
}

func main() {
	screenNames = make(map[string]ScreenBuilder)

	RegisterScreenName("main-menu", func(g *Game) IScreen {
		return CreateMainMenu(g)
	})

	RegisterScreenName("branding", func(g *Game) IScreen {
		return NewBrandingScreen(g, nil)
	})

	RegisterScreenName("cherryos-desktop", func(g *Game) IScreen {
		return NewWinXPScreen(g, nil)
	})

	RegisterScreenName("font-test", func(g *Game) IScreen {
		return NewFontTestScreen(g)
	})

	RegisterScreenName("shop", func(g *Game) IScreen {
		return NewShopScreen(g)
	})

	RegisterScreenName("game", func(g *Game) IScreen {
		return NewGameplayScreen(g)
	})

	argv := os.Args[1:]
	for _, argument := range argv {
		if strings.HasPrefix(argument, "--screen=") {
			param := strings.SplitN(argument, "=", 2)[1]

			_, has := screenNames[param]
			if has {
				directScreenSet = true
				directScreenName = param
				log.Println("Setting screen " + param)
				break
			}
		}
	}

	game, err := CreateGame()
	if err != nil {
		log.Fatal(err)
		return
	}

	var icon image.Image
	_, icon, _ = ebitenutil.NewImageFromFile("assets/icon.png")
	var iconImages [](image.Image)
	iconImages = append(iconImages, icon)

	ebiten.SetWindowIcon(iconImages)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	var cmdLine string
	var title string
	for _, arg := range os.Args[1:] {
		cmdLine += arg + " "
	}

	title = "Ronery"

	if len(os.Args[1:]) > 0 {
		title += " (" + cmdLine + ")"
	}

	ebiten.SetWindowTitle(title)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
