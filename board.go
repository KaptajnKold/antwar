package antwar

import (
	"rand"
	"container/vector"
	"image"
)

const (
	WIDTH = 800
	HEIGHT = 600
)


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

func RandomPos() Pos {
	return Pos{rand.Intn(800), rand.Intn(600)}
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

type Tile struct {
	ants vector.Vector
	food int
	team string
	base bool
}

func (t *Tile) AntCount() int {
	return t.ants.Len()
}

func (t *Tile) FoodCount() int {
	return t.food
}

func (t *Tile) RemoveAnt(theAnt *Ant) {
	for i := 0; i < t.ants.Len(); i++ {
		if anAnt, _ := t.ants.At(i).(*Ant); theAnt == anAnt {
			t.ants.Delete(i);
			return;
		}
	}
}

func (t *Tile) PutAnt(theAnt *Ant) {
	for i := 0; i < t.ants.Len(); i++ {
		if anAnt, _ := t.ants.At(i).(*Ant); theAnt == anAnt {
			return;
		}
	}
	t.ants.Push(theAnt)
	t.team = theAnt.Team
}

func (t *Tile) PutFood(amount int) {
	t.food += amount
}

func (t *Tile) RemoveFood(amount int) {
	t.food -= amount
}

func (t *Tile) CreateBase(team string) {
	t.base = true
	t.team = team
}

func (t *Tile) Color() image.Color {
	if t.AntCount() > 0 {
		return image.RGBAColor{255,255,255,100}
	}
	if t.FoodCount() > 0 {
		return image.RGBAColor{255,255,0,100}
	}
	if t.team != "" {
		return image.RGBAColor{100,100,100,100}
	}
	return image.Black;
}

type Board struct {
	Tiles [WIDTH][HEIGHT]*Tile
	OnUpdate func(p Pos)
}

func (b *Board) At(p Pos) *Tile {
	return (*b).Tiles[p.X][p.Y];
}

func (b *Board) Environment(pos Pos) *Environment {
	env := new(Environment)
	env[0] = b.At(pos)
	env[1] = b.At(pos.North())
	env[2] = b.At(pos.East())
	env[3] = b.At(pos.South())
	env[4] = b.At(pos.West())
	return env
}

func (b *Board) CreateFood(n int) {
	for i := 0; i < n; i++ {
		pos := RandomPos()
		b.At(pos).PutFood(rand.Intn(10))
		b.Update(pos)
	}
}

func (b *Board) Update(pos Pos) {
	if b.OnUpdate != nil {
		b.OnUpdate(pos)
	}
}

func NewBoard() *Board {
	board := new(Board);
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			board.Tiles[x][y] = new(Tile)
		}
	}
	return board
}

type Environment [5](*Tile)
