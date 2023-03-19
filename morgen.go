package main

type Morgen struct {
	WanderingNPC
}

func CreateMorgen(g *Game) *Morgen {
	e := new(Morgen)
	e.etype = e
	e._ConstructLivingEntity(g)
	e.entityClass = langData["entity_morgen"]
	e.speedModifier = 0.75
	e.sprite = morgenSprite
	return e
}
