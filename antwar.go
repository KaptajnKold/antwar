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

func mod(n, d int) int {
	return ((n % d) + d) % d;
}

type Pos struct {
	X, Y int
}

func (p *Pos) North() Pos {
	pos := *p
	return Pos{pos.X, mod(pos.Y - 1, 600)}
}
func (p *Pos) South() Pos {
	pos := *p
	return Pos{pos.X, mod(pos.Y + 1, 600)}
}
func (p *Pos) East() Pos {
	pos := *p
	return Pos{mod(pos.X + 1, 800), pos.Y}
}
func (p *Pos) West() Pos {
	pos := *p
	return Pos{mod(pos.X - 1, 800), pos.Y}
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

type Environment [5](Tile)

type AntBrain interface {
	Decide(env *Environment) (Action, bool)
}

type Ant struct {
	Brain AntBrain
	Team string
	Pos
}


func RandomPos() Pos {
	return Pos{rand.Intn(800), rand.Intn(600)}
}

func (t *Tile) Color() image.Color {
	if t.Ants > 0 {
		return image.RGBAColor{255,255,255,100}
	}
	if t.Food > 0 {
		return image.RGBAColor{255,255,0,100}
	}
	if t.Team != "" {
		return image.RGBAColor{100,100,100,100}
	}
	return image.Black;
}