package antwar

import (
	"rand"
	"image"
)

type Action int;

const (
	HERE Action = iota
	NORTH
	EAST
	SOUTH
	WEST
)

type Pos struct {
	X, Y int
}

type Base struct {
	Team string
	Pos
}

type Tile struct {
	Ants, Food int
	Team string
	Base bool
}

type Environment [5](*Tile)

type AntBrain interface {
	Decide(env *Environment) Action
}

type Ant struct {
	Brain AntBrain
	Pos
}

type Board struct {
	Tiles [800][600]Tile
}

func RandomPos() Pos {
	return Pos{rand.Intn(800), rand.Intn(600)}
}

func (t *Tile) Color() image.Color {
	if t.Ants > 0 {
		return image.RGBAColor{255,255,255,100}
	}
	if t.Team != "" {
		return image.RGBAColor{100,100,100,100}
	}
	return image.Black;
}
