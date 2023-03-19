package main

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"math"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type IDialogueBox interface {
	Draw(*ebiten.Image)
}

type DialogueBox struct {
	IDialogueBox
	game *Game
	name string
	font *Font
}

func (dialogueBox *DialogueBox) Draw(screen *ebiten.Image) {
	game := dialogueBox.game
	game.DrawHerbGUIFrame(screen, 0, 0, 128, 128)

	fontRenderer := game.fontRenderer
	fontRenderer.PushState()

	fontRenderer.SetFont(dialogueBox.font)
	fontRenderer.DrawTextAt(screen, dialogueBox.name, Vec2f{16.0, 16.0})

	fontRenderer.PopState()
}

func InitDialogueBox(dialogueBox *DialogueBox, game *Game) {
	dialogueBox.game = game
}

func NewDialogueBox(game *Game) *DialogueBox {
	dialogueBox := new(DialogueBox)
	InitDialogueBox(dialogueBox, game)

	return dialogueBox
}

type IVNSpriteMenu interface {
	Draw(*ebiten.Image)
	SetCharacter(*VNCharacter)
	GetCharacter() *VNCharacter
}

type VNSpriteMenu struct {
	IVNSpriteMenu
	game        *Game
	vnCharacter *VNCharacter
	font        *Font
}

var editgroup = 0
var vnEditTileSize float64 = 48

func (menu *VNSpriteMenu) Draw(screen *ebiten.Image) {
	menu.game.DrawHerbGUIFrame(screen, screenWidth-vnEditTileSize*5, 0, vnEditTileSize*5, screenHeight)

	var i int

	e := menu.vnCharacter.sprites[editgroup].sprites

	fontRenderer := menu.game.fontRenderer

	fontRenderer.PushState()

	fontRenderer.SetFont(menu.font)
	fontRenderer.SetScale(2.0)

	assetGroupTitle := menu.vnCharacter.sprites[editgroup].name

	dim := fontRenderer.GetStringDimensions(assetGroupTitle)
	fontRenderer.DrawTextAt(screen, assetGroupTitle, Vec2f{screenWidth - vnEditTileSize*5, dim.Y})

	screenTranslate := ebiten.GeoM{}
	screenTranslate.Translate(screenWidth-vnEditTileSize*5, dim.Y*3)

	fontRenderer.PopState()

	for _, spriteVariant := range e {
		op := &ebiten.DrawImageOptions{}

		imgWidth, imgHeight := spriteVariant.image.Size()
		op.GeoM.Translate(-(float64(imgWidth) / 2), -(float64(imgHeight) / 2))

		scale := float64(vnEditTileSize) / math.Max(float64(imgWidth), float64(imgHeight))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(vnEditTileSize/2, vnEditTileSize/2)

		sx := float64(i%5) * float64(vnEditTileSize)
		sy := float64(i/5) * float64(vnEditTileSize)

		op.GeoM.Translate(sx, sy)

		op.GeoM.Concat(screenTranslate)

		if menu.vnCharacter.sprites[editgroup].current == spriteVariant {
			x, y := op.GeoM.Apply(1, 1)
			ebitenutil.DrawRect(screen, x, y, vnEditTileSize, vnEditTileSize, color.RGBA{0, 255, 0, 128})
		}

		screen.DrawImage(spriteVariant.image, op)

		i++
	}
}

func (menu *VNSpriteMenu) GetCharacter() *VNCharacter {
	return menu.vnCharacter
}

func (menu *VNSpriteMenu) SetCharacter(char *VNCharacter) {
	menu.vnCharacter = char
}

func InitVNSpriteMenu(menu *VNSpriteMenu, g *Game) {
	rm := ResourceManager_GetInstance()

	menu.game = g
	menu.font = rm.LoadFontJSON("font/font_fantasy.json")
}

func NewVNSpriteMenu(g *Game) *VNSpriteMenu {
	menu := new(VNSpriteMenu)
	InitVNSpriteMenu(menu, g)
	return menu
}

type ShopScreen struct {
	Screen
	backgroundImage *ebiten.Image
	tableImage      *ebiten.Image
	endou           *VNCharacter
	morshu          *VNCharacter
	spriteMenu      IVNSpriteMenu
	dialogueBox     IDialogueBox
}

var le = 0

func (s *ShopScreen) ProcessKeyEvents() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		facial := s.spriteMenu.GetCharacter().sprites[editgroup]

		var j int
		for i, e := range facial.sprites {
			if facial.current == e {
				j = i
			}
		}

		if j-1 < 0 {
			j = len(facial.sprites) - 1
		} else {
			j--
		}

		facial.current = facial.sprites[j]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		facial := s.spriteMenu.GetCharacter().sprites[editgroup]

		var j int
		for i, e := range facial.sprites {
			if facial.current == e {
				j = i
			}
		}

		if j+1 > len(facial.sprites)-1 {
			j = 0
		} else {
			j++
		}

		facial.current = facial.sprites[j]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		if editgroup-1 < 0 {
			editgroup = len(s.endou.sprites) - 1
		} else {
			editgroup--
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		if editgroup+1 > len(s.endou.sprites)-1 {
			editgroup = 0
		} else {
			editgroup++
		}
	}

	return true
}

func (s *ShopScreen) Update() {
	s.endou.matTransform.Reset()
	s.endou.matTransform.Translate(-(float64(s.endou.width) / 2), -(float64(s.endou.height) / 2))

	scale := 0.975 + math.Sin(float64(s.game.appTicker)/72.0)*0.025
	jump := math.Sin(float64(s.game.appTicker) / 16.0)
	s.endou.matTransform.Translate(0, jump*10)
	s.endou.matTransform.Scale(scale, scale)
	s.endou.matTransform.Rotate(math.Sin(float64(s.game.appTicker)/48.0) * 0.05)
	s.endou.matTransform.Translate(float64(s.endou.width)/2, float64(s.endou.height)/2)
	s.IScreen.ProcessKeyEvents()
}

func (s *ShopScreen) LoadResources() {
	rm := ResourceManager_GetInstance()

	s.backgroundImage = rm.LoadImage("assets/shop/background.png")
	s.tableImage = rm.LoadImage("assets/shop/table.png")

	s.endou = NewVNCharacterFromJSON("assets/shop/test/endou/endou.json", "Endou", g_endouMapping)
	s.morshu = NewVNCharacterFromJSON("assets/shop/morshu/morshu.json", "Morshu", g_morshuMapping)
	s.spriteMenu.SetCharacter(s.morshu)
}

func DrawStretchedImage(screen *ebiten.Image, img *ebiten.Image, rc Rect) {
	var sx, sy, sx2, sy2 float64

	sx = float64(img.Bounds().Min.X)
	sy = float64(img.Bounds().Min.Y)

	sx2 = float64(img.Bounds().Max.X)
	sy2 = float64(img.Bounds().Max.Y)

	vs := []ebiten.Vertex{
		{
			DstX:   float32(rc.X1),
			DstY:   float32(rc.Y1),
			SrcX:   float32(sx),
			SrcY:   float32(sy),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(rc.X2),
			DstY:   float32(rc.Y1),
			SrcX:   float32(sx2),
			SrcY:   float32(sy),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(rc.X1),
			DstY:   float32(rc.Y2),
			SrcX:   float32(sx),
			SrcY:   float32(sy2),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(rc.X2),
			DstY:   float32(rc.Y2),
			SrcX:   float32(sx2),
			SrcY:   float32(sy2),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
	}

	triOp := &ebiten.DrawTrianglesOptions{}
	triOp.Address = ebiten.AddressClampToZero
	screen.DrawTriangles(vs, []uint16{0, 1, 2, 1, 2, 3}, img, triOp)
}

func (s *ShopScreen) Draw(screen *ebiten.Image) {

	DrawStretchedImage(screen, s.backgroundImage, Rect{0, 0, screenWidth, screenHeight})

	op := &ebiten.DrawImageOptions{}
	imgWidth, imgHeight := s.tableImage.Size()
	op.GeoM.Translate(-float64(imgWidth/2), -float64(imgHeight/2))
	op.GeoM.Translate(screenWidth/2, float64(screenHeight-imgHeight/2))

	s.morshu.Draw(screen)

	screen.DrawImage(s.tableImage, op)

	s.spriteMenu.Draw(screen)

	s.dialogueBox.Draw(screen)
}

func InitShopScreen(s *ShopScreen, g *Game) {
	s.game = g
	s.spriteMenu = NewVNSpriteMenu(g)
	s.dialogueBox = NewDialogueBox(g)
}

func NewShopScreen(g *Game) *ShopScreen {
	s := new(ShopScreen)
	InitShopScreen(s, g)
	s.IScreen = s
	return s
}

type JSONLayeredImage struct {
	Width  uint                    `json:"width"`
	Height uint                    `json:"height"`
	Layers []JSONLayeredImageEntry `json:"layers"`
}

type JSONLayeredImageEntry struct {
	Name    string `json:"name"`
	Index   uint   `json:"index"`
	OffsetX uint   `json:"offset_x"`
	OffsetY uint   `json:"offset_y"`
	Width   uint   `json:"width"`
	Height  uint   `json:"height"`
	image   *ebiten.Image
}

type IVNCharacter interface {
	Draw(*ebiten.Image)
	SetExpression(string)
	SetBody(string)
}

type VNCharacterSprite struct {
	pos   Vec2f
	image *ebiten.Image
}

func (sprite *VNCharacterSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sprite.pos.X, sprite.pos.Y)

	screen.DrawImage(sprite.image, op)
}

type VNCharacter struct {
	IVNCharacter
	sprites []*VNSpriteListCategory

	width        uint
	height       uint
	name         string
	matTransform ebiten.GeoM
}

func (char *VNCharacter) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	for i := 0; i < len(char.sprites); i++ {
		sprite := char.sprites[i].current

		op.GeoM.Reset()
		op.GeoM.Translate(sprite.pos.X, sprite.pos.Y)
		op.GeoM.Concat(char.matTransform)
		screen.DrawImage(sprite.image, op)
	}
}

type VNAssetGroupMapping struct {
	name       string
	assetNames []string
}

var g_morshuMapping = []VNAssetGroupMapping{
	{
		name: "body",
		assetNames: []string{
			"morshu_090.png",
		},
	},
	{
		name: "facial",
		assetNames: []string{
			"morshu_001.png",
			"morshu_002.png",
		},
	},
}

var g_endouMapping = []VNAssetGroupMapping{
	{
		name: "body:",
		assetNames: []string{
			"endou_044.png",
			"endou_047.png",
			"endou_050.png",
			"endou_053.png",
		},
	},
	{
		name: "facial",
		assetNames: []string{
			"endou_001.png",
			"endou_002.png",
			"endou_003.png",
			"endou_004.png",
			"endou_005.png",
			"endou_006.png",
			"endou_007.png",
			"endou_008.png",
			"endou_009.png",
			"endou_010.png",
			"endou_011.png",
			"endou_012.png",
			"endou_013.png",
			"endou_014.png",
			"endou_015.png",
			"endou_016.png",
			"endou_017.png",
			"endou_018.png",
			"endou_019.png",
			"endou_020.png",
			"endou_021.png",
			"endou_022.png",
			"endou_023.png",
			"endou_024.png",
			"endou_025.png",
			"endou_026.png",
			"endou_027.png",
			"endou_028.png",
			"endou_029.png",
			"endou_030.png",
			"endou_031.png",
			"endou_032.png",
			"endou_033.png",
			"endou_034.png",
			"endou_035.png",
			"endou_036.png",
			"endou_037.png",
			"endou_038.png",
			"endou_039.png",
			"endou_040.png",
		},
	},
}

func LoadGroupAssets(char *VNCharacter, j JSONLayeredImage, coherence []VNAssetGroupMapping, dir string) {

	rm := ResourceManager_GetInstance()

	for _, coh := range coherence {
		list := new(VNSpriteListCategory)
		for _, e := range j.Layers {
			for _, z := range coh.assetNames {
				if e.Name == z {
					sp := new(VNCharacterSprite)
					sp.image = rm.LoadImage(dir + e.Name)
					sp.pos.X = float64(e.OffsetX)
					sp.pos.Y = float64(e.OffsetY)

					list.sprites = append(list.sprites, sp)
					break
				}
			}
		}
		list.name = coh.name
		list.current = list.sprites[0]
		char.sprites = append(char.sprites, list)
	}

	char.width = j.Width
	char.height = j.Height
}

func InitVNCharacter(char *VNCharacter, name string) {
	char.name = name
}

func NewVNCharacter(name string) *VNCharacter {
	char := new(VNCharacter)
	InitVNCharacter(char, name)
	return char
}

func NewVNCharacterFromJSON(path string, name string, assetGroupMapping []VNAssetGroupMapping) *VNCharacter {
	char := new(VNCharacter)
	InitVNCharacter(char, name)

	assetPath, _ := filepath.Split(path)

	var layeredImage JSONLayeredImage
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, &layeredImage)
	if err != nil {
		return nil
	}

	LoadGroupAssets(char, layeredImage, assetGroupMapping, assetPath)

	return char
}

type VNSpriteListCategory struct {
	name    string
	sprites []*VNCharacterSprite
	current *VNCharacterSprite
}
