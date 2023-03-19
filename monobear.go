package main

type Monobear struct {
	WanderingNPC
}

func CreateMonobear(g *Game) *Monobear {
	e := new(Monobear)
	e.etype = e
	e._ConstructLivingEntity(g)
	e.entityClass = langData["entity_monobear"]
	e.speedModifier = 0.75
	e.sprite = monobearSprite
	// e.spells = append(e.spells, CreateMonobearExplosion(e))
	return e
}
