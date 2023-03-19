package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

/*
	Use this somewhere:

	mapScreenWidth := float64(g.level.width*tileSize) * g.camera.GetZoom()
	mapScreenHeight := float64(g.level.height*tileSize) * g.camera.GetZoom()
	if mapScreenWidth < screenWidth {
		g.view.offsetX = (screenWidth - mapScreenWidth) / 2
	} else {
		g.view.offsetX = math.Max(-(mapScreenWidth - screenWidth), math.Min(0, screenWidth/2-cameraFocusEntiy.worldPosX*g.camera.GetZoom()))
	}

	if mapScreenHeight < screenHeight {
		g.view.offsetY = (screenHeight - mapScreenHeight) / 2
	} else {
		g.view.offsetY = math.Max(-(mapScreenHeight - screenHeight), math.Min(0, screenHeight/2-cameraFocusEntiy.worldPosY*g.camera.GetZoom()))
	}
*/

type ICamera interface {
	TargetPosition(Vec2f)
	TargetEntity(*ILivingEntity)
	TargetCursor(*TileCursor)
	Update(Vec2f)
	GetZoom() float64
	SetZoom(float64)
	GetCurrentPosition() Vec2f
	WorldToScreen2(Vec2f) Vec2f
}

type CameraTargetType uint8

const (
	CAMERA_TARGET_POSITION CameraTargetType = 1
	CAMERA_TARGET_ENTITY   CameraTargetType = 2
	CAMERA_TARGET_CURSOR   CameraTargetType = 3
)

type Camera struct {
	zoom            float64
	currentWorldPos Vec2f
	targetWorldPos  Vec2f
	targetEntity    ILivingEntity
	targetCursor    *TileCursor
	targetType      CameraTargetType
	tickAnimStarted int
}

func (c *Camera) TargetPosition(pos Vec2f) {
	c.targetWorldPos = pos
	c.targetType = CAMERA_TARGET_POSITION
	c.tickAnimStarted = tickCounter
}

func (c *Camera) TargetEntity(e ILivingEntity) {
	c.targetEntity = e
	c.targetType = CAMERA_TARGET_ENTITY
	c.tickAnimStarted = tickCounter
}

func (c *Camera) TargetCursor(cur *TileCursor) {
	c.targetCursor = cur
	c.targetType = CAMERA_TARGET_CURSOR
	c.tickAnimStarted = tickCounter
}

func EaseOutQuad(factor float64) float64 {
	return 1 - (1-factor)*(1-factor)
}

func (c *Camera) Update() {
	var duration int = ebiten.MaxTPS() * 1

	if c.targetType == CAMERA_TARGET_POSITION {
		c.currentWorldPos = c.targetWorldPos
	} else if c.targetType == CAMERA_TARGET_ENTITY && c.targetEntity != nil {
		animIteration := tickCounter - c.tickAnimStarted
		entityPos := c.targetEntity.GetWorldPos()

		if animIteration <= duration {
			animFactor := float64(animIteration) / float64(duration)
			delta := entityPos.Translate(c.currentWorldPos.Negate()).Scale(animFactor)

			c.currentWorldPos = c.currentWorldPos.Translate(delta)

		} else {
			c.currentWorldPos = entityPos
		}
	} else if c.targetType == CAMERA_TARGET_CURSOR {
		animIteration := tickCounter - c.tickAnimStarted

		cursorPos := Vec2f{float64(c.targetCursor.x), float64(c.targetCursor.y)}.Scale(tileSize).Translate(Vec2f{tileSize / 2, tileSize / 2})

		if animIteration <= duration {
			animFactor := float64(animIteration) / float64(duration)
			delta := cursorPos.Translate(c.currentWorldPos.Negate()).Scale(animFactor)

			c.currentWorldPos = c.currentWorldPos.Translate(delta)

		} else {
			c.currentWorldPos = cursorPos
		}
	}
}

func (c *Camera) GetCurrentPosition() Vec2f {
	return c.currentWorldPos
}

func (c *Camera) WorldToScreen2(pos Vec2f) Vec2f {
	screenCenter := Vec2f{screenWidth / 2, screenHeight / 2}
	cameraWorldPos := c.currentWorldPos
	entityWorldPos := pos
	cameraZoom := c.zoom

	var screenPos Vec2f
	screenPos = cameraWorldPos.Subtract(entityWorldPos)
	screenPos = screenPos.Scale(cameraZoom)
	screenPos = screenCenter.Subtract(screenPos)

	return screenPos
}

func (c *Camera) SetZoom(zoom float64) {
	c.zoom = math.Max(1.0, zoom)
}

func (c *Camera) GetZoom() float64 {
	return c.zoom
}
