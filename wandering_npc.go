package main

type WanderingNPC struct {
	LivingEntity

	walkerState int
}

func (e *WanderingNPC) Update() {
	e.LivingEntity.Update()

	if tickCounter%64 == 0 {

		switch e.walkerState {
		case 0:
			e.StartWalk(LookDirection(r.Int() % 4))
		default:
			e.EndWalk()
		}

		e.walkerState = (e.walkerState + 1) % 2
	}

	e.ProcessWalk()

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
