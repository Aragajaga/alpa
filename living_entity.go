package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ILivingEntity interface {
	SetGame(*Game)
	StartWalk(LookDirection)
	EndWalk()
	GetTilePos() (int, error)
	Update()
	Draw(*ebiten.Image)
	DrawDebugInfo(*ebiten.Image)
	GetHealth() float64
	GetWorldPos() Vec2f
	GetLivingEntity() *LivingEntity
	ProcessWalk()
}

type LivingEntity struct {
	etype         interface{}
	entityClass   string
	look          LookDirection
	walking       bool
	worldPos      Vec2f
	anchorPos     Vec2f
	baseSpeed     float64
	health        float64
	speedModifier float64
	sprite        *ebiten.Image
	game          *Game
	id            int
	spells        []ISpell
	prevTilePos   int
}

func (e *LivingEntity) SetGame(g *Game) {
	e.game = g
}

func (e *LivingEntity) StartWalk(look LookDirection) {
	e.walking = true
	e.look = look
}

func (e *LivingEntity) EndWalk() {
	e.walking = false
}

func (e *LivingEntity) ProcessWalk() {
	if e.walking {
		speed := e.baseSpeed * e.speedModifier

		switch look := e.look; look {
		case LooksRight:
			if !e.game.IsTileSolidAt(e.worldPos.X+speed, e.worldPos.Y) {
				e.worldPos.X += speed
			}

		case LooksLeft:
			if !e.game.IsTileSolidAt(e.worldPos.X-speed, e.worldPos.Y) {
				e.worldPos.X -= speed
			}

		case LooksUp:
			if !e.game.IsTileSolidAt(e.worldPos.X, e.worldPos.Y-speed) {
				e.worldPos.Y -= speed
			}

		case LooksDown:
			if !e.game.IsTileSolidAt(e.worldPos.X, e.worldPos.Y+speed) {
				e.worldPos.Y += speed
			}
		}
	}
}

func (e *LivingEntity) GetTilePos() (int, error) {
	return e.game.WorldPosToTilePos(e.worldPos.X, e.worldPos.Y)
}

func (e *LivingEntity) Update() {
	for _, spell := range e.spells {
		spell.Update()
	}
}

func (e *LivingEntity) SetSpeedModifier(speed float64) {
	e.speedModifier = speed
}

const CHAR_TILE_SIZE int = 16

func (e *LivingEntity) DrawHealthBar(screen *ebiten.Image) {
	cameraZoom := e.game.camera.GetZoom()
	// screenPos := e.game.camera.WorldToScreen(e.worldPos)

	borderColor := color.RGBA{0, 0, 0, 255}
	fillColor := HSL{math.Pow(e.health/100.0, 2) * 0.33, 1.0, 0.5}.ToRGB()

	// health bar offset from character origin in world coordinates
	healthBarOffset := Vec2f{-(float64(CHAR_TILE_SIZE) / 2), -(float64(CHAR_TILE_SIZE) + 4)}

	pos := Vec2f{screenWidth / 2, screenHeight / 2}
	pos = pos.Translate(
		e.game.camera.GetCurrentPosition().Translate(
			e.GetWorldPos().Translate(healthBarOffset).Negate()).Scale(
			e.game.camera.GetZoom()).Negate())

	ebitenutil.DrawRect(
		screen,
		pos.X,
		pos.Y,
		16*cameraZoom,
		3*cameraZoom, borderColor)

	ebitenutil.DrawRect(
		screen,
		pos.X+cameraZoom,
		pos.Y+cameraZoom,
		(e.health/100.0)*(16-2)*cameraZoom, cameraZoom, fillColor)

	pos = pos.Translate(Vec2f{0, -11})

	fontRenderer := e.game.fontRenderer

	fontRenderer.PushState()

	fontRenderer.Reset()
	fontRenderer.SetTextColor(color.White)

	if e.game.camera.targetEntity.GetLivingEntity() == e {
		fontRenderer.SetTextColor(color.RGBA{255, 0, 0, 255})
	}

	fontRenderer.DrawTextAt(screen, e.entityClass, pos)

	fontRenderer.PopState()
}

type ISpell interface {
	Update()
	Draw(*ebiten.Image)
	GetSpell() *Spell
}

type Spell struct {
	creationTime int
	caster       ILivingEntity
}

func (s *Spell) GetSpell() *Spell {
	return s
}

type MonobearExplosion struct {
	Spell
}

func (s *MonobearExplosion) Update() {
	/*
		elapsed := tickCounter - s.creationTime

		if elapsed >= 256 {
			for i, spell := range s.caster.GetLivingEntity().spells {
				if s.GetSpell() == spell.GetSpell() {
					s.caster.GetLivingEntity().spells[i] = nil
					s.caster.GetLivingEntity().spells = append(s.caster.GetLivingEntity().spells[:i], s.caster.GetLivingEntity().spells[i+1:]...)
					break
				}
			}
		}
	*/
}

func CreateMonobearExplosion(caster ILivingEntity) *MonobearExplosion {
	s := new(MonobearExplosion)
	s.caster = caster
	s.creationTime = tickCounter

	return s
}

func (s *MonobearExplosion) Draw(screen *ebiten.Image) {
	e := s.caster.GetLivingEntity()

	{
		op := &ebiten.DrawImageOptions{}

		pos := e.game.camera.WorldToScreen2(e.worldPos)

		op.GeoM.Translate(-16, -16)
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(skillMonobearExplosion, op)
	}

	{
		screenCenter := Vec2f{screenWidth / 2, screenHeight / 2}
		cameraWorldPos := e.game.camera.GetCurrentPosition()
		entityWorldPos := e.GetWorldPos()
		cameraZoom := e.game.camera.GetZoom()

		var pos Vec2f
		pos = cameraWorldPos.Subtract(entityWorldPos)
		pos = pos.Scale(cameraZoom)
		pos = screenCenter.Subtract(pos)

		e.game.fontRenderer.DrawTextAt(screen, "Monobear Explosion", pos)
	}

	{
		op := &ebiten.DrawImageOptions{}

		screenCenter := Vec2f{screenWidth / 2, screenHeight / 2}
		cameraWorldPos := e.game.camera.GetCurrentPosition()
		entityWorldPos := e.GetWorldPos()
		cameraZoom := e.game.camera.GetZoom()

		var pos Vec2f
		pos = cameraWorldPos.Subtract(entityWorldPos)
		pos = pos.Scale(cameraZoom)
		pos = screenCenter.Subtract(pos)

		i := tickCounter / 4 % 12

		sx := (i % 4) * 16
		sy := (i / 4) * 16

		op.GeoM.Translate(-(float64(tileSize) / 2), -(float64(tileSize) / 2))
		op.GeoM.Scale(cameraZoom, cameraZoom)
		op.GeoM.Translate(pos.X, pos.Y)

		tile := explosionSprite.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image)
		screen.DrawImage(tile, op)
	}
}

func (e *LivingEntity) Draw(screen *ebiten.Image) {
	var walkAnim int
	var spriteLine int
	var sx int
	var sy int
	var sprite *ebiten.Image

	op := &ebiten.DrawImageOptions{}
	walkAnim = 0

	if e.walking {
		speed := e.baseSpeed * e.speedModifier
		animTicker := int(float64(float64(tickCounter) * speed))
		walkAnim += animTicker / 4 % 3
	}
	spriteLine = int(e.look) * 4

	sx = ((spriteLine + walkAnim) % 4) * tileSize
	sy = ((spriteLine + walkAnim) / 4) * tileSize

	cameraZoom := e.game.camera.GetZoom()

	pos := e.game.camera.WorldToScreen2(e.worldPos)

	op.GeoM.Translate(-e.anchorPos.X, -e.anchorPos.Y)
	op.GeoM.Scale(cameraZoom, cameraZoom)
	op.GeoM.Translate(pos.X, pos.Y)

	sprite = e.sprite.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image)
	screen.DrawImage(sprite, op)

	e.DrawHealthBar(screen)
}

func (e *LivingEntity) DrawDebugInfo(screen *ebiten.Image) {
	// health bar offset from character origin in world coordinates
	healthBarOffset := Vec2f{float64(CHAR_TILE_SIZE) / 2, -(float64(CHAR_TILE_SIZE))}

	pos := Vec2f{screenWidth / 2, screenHeight / 2}
	pos = pos.Translate(
		e.game.camera.GetCurrentPosition().Translate(
			e.GetWorldPos().Translate(healthBarOffset).Negate()).Scale(
			e.game.camera.GetZoom()).Negate())

	infoPrinters := []func(*LivingEntity) string{
		func(e *LivingEntity) string {
			return fmt.Sprintf("x: %f, y: %f", e.GetWorldPos().X, e.GetWorldPos().Y)
		},
		func(e *LivingEntity) string {
			return fmt.Sprintf("walking: %t", e.walking)
		},
		func(e *LivingEntity) string {
			return fmt.Sprintf("health: %f", e.health)
		},
		func(e *LivingEntity) string {
			tilePos, _ := e.GetTilePos()
			return fmt.Sprintf("CurTilePos: %d", tilePos)
		},
		func(e *LivingEntity) string {
			return fmt.Sprintf("PrevTilePos: %d", e.GetLivingEntity().prevTilePos)
		},
	}

	fontRenderer := e.game.fontRenderer

	fontRenderer.PushState()
	fontRenderer.Reset()

	for i, printer := range infoPrinters {
		fontRenderer.DrawTextAt(screen, printer(e), pos)
		pos = pos.Translate(Vec2f{0, float64(fontRenderer.GetGlyphSize().X+1) * float64(i)})
	}

	fontRenderer.PopState()
}

var eidCounter int

func (e *LivingEntity) _ConstructLivingEntity(g *Game) {
	e.worldPos = Vec2f{CHARACTER_SPAWN_X, CHARACTER_SPAWN_Y}
	e.anchorPos = Vec2f{CHARACTER_ANCHOR_X, CHARACTER_ANCHOR_Y}
	e.baseSpeed = 1.0
	e.speedModifier = 1.0
	e.health = 100.0
	e.id = eidCounter
	eidCounter++
	e.SetGame(g)
}

func (e *LivingEntity) GetUnderlyingTiles() ([]Tile, error) {
	var tiles []Tile

	tilePos, err := e.GetTilePos()
	if err != nil {
		return nil, err
	}

	for _, layer := range e.game.level.tileLayers {
		if layer[tilePos] != Tile(tileIDEmpty) {
			tiles = append(tiles, layer[tilePos])
		}
	}

	return tiles, err
}

func (e *LivingEntity) GetHealth() float64 {
	return e.health
}

func (e *LivingEntity) GetWorldPos() Vec2f {
	return e.worldPos
}

func (e *LivingEntity) GetLivingEntity() *LivingEntity {
	return e
}
