package engine

type Player struct {
	*Inventory
}

func NewPlayer() *Player {
	return &Player{
		Inventory: NewInventory(nil),
	}
}
