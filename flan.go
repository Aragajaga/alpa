package main

type Flan struct {
	WanderingNPC
}

func CreateFlan(g *Game) *Flan {
	e := new(Flan)
	e.etype = e
	e._ConstructLivingEntity(g)
	e.entityClass = langData["entity_flan"]
	e.sprite = flanSprite
	return e
}
