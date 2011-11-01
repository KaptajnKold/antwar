package antwar

import (
	"rand"
)

const (
	WIDTH = 800
	HEIGHT = 600
)

type Board struct {
	Tiles [WIDTH][HEIGHT]*Tile
	OnUpdate func(p Pos)
}

func (b *Board) At(p Pos) *Tile {
	return (*b).Tiles[p.X][p.Y];
}

func (b *Board) Environment(pos Pos) *Environment {
	env := new(Environment)
	env[0] = *b.At(pos)
	env[1] = *b.At(pos.North())
	env[2] = *b.At(pos.East())
	env[3] = *b.At(pos.South())
	env[4] = *b.At(pos.West())
	return env
}

func (b *Board) CreateFood(n int) {
	for i := 0; i < n; i++ {
		pos := RandomPos()
		tile := b.At(pos)
		tile.Food += rand.Intn(10)
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