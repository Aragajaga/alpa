package main

type Michael struct {
	WanderingNPC
}

func CreateMichael(g *Game) *Michael {
	e := new(Michael)
	e.etype = e
	e._ConstructLivingEntity(g)
	e.entityClass = langData["entity_michael"]
	e.sprite = michaelSprite
	return e
}
