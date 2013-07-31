package antwar

import (
	"image/color"
	"fmt"
)

type Tile struct {
	ants                           AntSet
	food                           int
	team                           *Team
	board                          *Board
	pos                            Pos
	antHill                        *AntHill
	here, north, east, south, west *Tile
}

func (t *Tile) removeAnt(theAnt *Ant) {
	t.ants.Remove(theAnt)
}

func (t *Tile) putAnt(a *Ant) {
	if t.team != a.team {
		t.killAnts()
		t.destroyAntHill()
		t.team = a.team
	}

	t.ants.Put(a)
	a.tile = t
}

func (t *Tile) putFood(amount int) {
	t.food += amount
}

func (t *Tile) removeFood(amount int) {
	t.food -= amount
}

func (t *Tile) createAntHill(team *Team) *AntHill {
	anthill := new(AntHill)
	anthill.team = team
	anthill.tile = t
	t.team = team
	t.antHill = anthill
	fmt.Printf("Created anthill at %v for team %v\n", anthill.tile.pos, team.name)
	return anthill
}

func (t *Tile) destroyAntHill() {
	if t.antHill == nil {
		return
	}
	fmt.Printf("Destroyed anthill for %v", t.antHill.team.name)
	t.board.removeAntHill(t.antHill)
	t.antHill = nil

}

func (t *Tile) killAnts() {
	if 0 == t.ants.Len() {
		return
	}
	t.ants.Do(func(a *Ant) {
		t.board.ants.Remove(a)
		t.ants.Remove(a)
		a.team.ants.Remove(a)
	})
}

func (t *Tile) color() color.Color {
	if t.ants.Len() > 0 {
		return t.team.color
	}
	if t.team != nil {
		r, g, b, _ := t.team.color.RGBA()
		return color.RGBA{uint8(int(r) << 7), uint8(int(g) << 7), uint8(int(b) << 7), 10}
	}
	return color.RGBA{255, 255, 255, uint8(t.food)}
}

func (t *Tile) update() {
	t.board.Update(t)
}

func (t *Tile) AntCount() int {
	return t.ants.Len()
}

func (t *Tile) FoodCount() int {
	return t.food
}

func (t *Tile) Here() *Tile {
	return t
}

func (t *Tile) North() *Tile {
	return t.north
}

func (t *Tile) East() *Tile {
	return t.east
}

func (t *Tile) South() *Tile {
	return t.south
}

func (t *Tile) West() *Tile {
	return t.west
}
