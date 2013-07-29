package antwar

import (
	"image/color"
)

type Tile struct {
	Ants            AntSet
	Food            int
	Team            *Team
	Environment     *Environment
	environmentTile *TileInfo
	board           *Board
	pos             Pos
}

func (t *Tile) RemoveAnt(theAnt *Ant) {
	t.Ants.Remove(theAnt)
}

func (t *Tile) PutAnt(theAnt *Ant) {
	t.Ants.Put(theAnt)
	t.Team = theAnt.Team
}

func (t *Tile) PutFood(amount int) {
	t.Food += amount
}

func (t *Tile) RemoveFood(amount int) {
	t.Food -= amount
}

func (t *Tile) CreateAntHill(team *Team) *AntHill {
	anthill := new(AntHill)
	anthill.Team = team
	anthill.tile = t
	return anthill
}

func (t *Tile) Color() color.Color {
	if t.Ants.Len() > 0 {
		return t.Team.Color
	}
	if t.Team != nil {
		r, g, b, _ := t.Team.Color.RGBA()
		return color.RGBA{uint8(int(r) << 7), uint8(int(g) << 7), uint8(int(b) << 7), 10}
	}
	return color.RGBA{255, 255, 255, uint8(t.Food)}
}

type TileInfo struct {
	tile *Tile
}

func (e *TileInfo) AntCount() int {
	return e.tile.Ants.Len()
}

func (e *TileInfo) FoodCount() int {
	return e.tile.Food
}

func (e *TileInfo) Team() string {
	if e.tile.Ants.Len() > 0 {
		return e.tile.Team.Name
	}
	return ""
}
