package main

type ICharacter interface {
	ILivingEntity
}

type Character struct {
	LivingEntity
}

const (
	CHARACTER_SPAWN_X  = 40
	CHARACTER_SPAWN_Y  = 200
	CHARACTER_ANCHOR_X = 7
	CHARACTER_ANCHOR_Y = 15
)

func (e *Character) Update() {

	if !e.game.gameOver {
		e.ProcessWalk()
	}

	curTilePos, _ := e.GetTilePos()

	if e.prevTilePos != curTilePos {
		e.game.ProcessTileLeaving(e, e.prevTilePos)

		e.prevTilePos = curTilePos
	}

	underlyingTiles, _ := e.game.GetUnderlyingTilesAt(e.worldPos.X, e.worldPos.Y)
	for _, tile := range underlyingTiles {
		e.game.ProcessTileEffects(e, tile)
	}
}

func CreateCharacter(g *Game) *Character {
	e := new(Character)
	e.etype = e
	e._ConstructLivingEntity(g)
	e.entityClass = langData["entity_player"]
	e.sprite = charSprite
	return e
}
