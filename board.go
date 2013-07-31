package antwar

import (
	"fmt"
	"math/rand"
	"sort"
)

const (
	HERE Action = iota
	NORTH
	EAST
	SOUTH
	WEST
)

func mod(n, d int) int {
	return ((n % d) + d) % d
}

type Pos struct {
	X, Y int
}

type column []*Tile

func (board *Board) NorthFrom(pos Pos) Pos {
	return Pos{pos.X, mod(pos.Y-1, board.height)}
}
func (board *Board) SouthFrom(pos Pos) Pos {
	return Pos{pos.X, mod(pos.Y+1, board.height)}
}
func (board *Board) EastFrom(pos Pos) Pos {
	return Pos{mod(pos.X+1, board.width), pos.Y}
}
func (board *Board) WestFrom(pos Pos) Pos {
	return Pos{mod(pos.X-1, board.width), pos.Y}
}

type Board struct {
	width, height int
	columns       []column
	onUpdate      func(p Pos)
	ants          AntSet
	antHills      []*AntHill
	teams         []*Team
}

func (b *Board) randomPos() Pos {
	return Pos{rand.Intn(b.width), rand.Intn(b.height)}
}

func (b *Board) randomTile() *Tile {
	return b.At(b.randomPos())
}

func (b *Board) At(p Pos) *Tile {
	return (*b).columns[p.X][p.Y]
}

func (board *Board) createStartingAntHills() {
	numberOfTeams := len(board.teams)
	board.antHills = make([](*AntHill), numberOfTeams)
	for i, team := range board.teams {
		board.antHills[i] = board.randomTile().createAntHill(team)
	}
	fmt.Printf("Created %v starting ant hills\n", len(board.antHills))
}

func (board *Board) CreateStartingAnts(startingNumberOfAnts int) {
	// Add 1 ant at a time to each team to avoid one team's ants all get to move before the other teams' ants
	for _, antHill := range board.antHills {
		fmt.Printf("Initially %v starting ants for team %v.\n", antHill.team.ants.Len(), antHill.team.name)
	}
	for i := 0; i < startingNumberOfAnts; i++ {
		for _, antHill := range board.antHills {
			antHill.spawnAnt()
		}
	}
	for _, antHill := range board.antHills {
		fmt.Printf("Created %v starting ants for team %v.\n", antHill.team.ants.Len(), antHill.team.name)
	}
	fmt.Printf("Created %v starting ants.\n", board.ants.Len())
}

func (b *Board) SpawnFoodRandomly(ncolumns, nFood int) {
	for i, j := 0, rand.Intn(ncolumns); i < j; i++ {
		t := b.randomTile()
		t.putFood(rand.Intn(nFood))
		t.update()
	}
}

func (board *Board) SpawnAnts() {
	for _, antHill := range board.antHills {
		antHill.spawnAnts()
	}
}

func moveFood(fromTile, toTile *Tile) {
	if fromTile.food > 0 {
		fromTile.removeFood(1)
		toTile.putFood(1)
	}
}

func (b *Board) removeAntHill(h *AntHill) {
	for index, candidate := range(b.antHills) {
		if candidate == h {
			b.antHills = append(b.antHills[:index], b.antHills[index + 1:]...)
			return
		}
	}
}

func findDestination(origin *Tile, decision Action) (destination *Tile) {
	switch decision {
	case NORTH:
		destination = origin.North()
	case SOUTH:
		destination = origin.South()
	case EAST:
		destination = origin.East()
	case WEST:
		destination = origin.West()
	case HERE:
		destination = origin
	}
	return
}

func (board *Board) MoveAnts() {
	board.ants.Do(func(ant *Ant) {
		origin := ant.tile
		brains := origin.ants.brainsExcept(ant)

		decision, bringFood := ant.brain.Decide(origin, brains)
		destination := findDestination(origin, decision)

		if origin == destination {
			return
		}

		origin.removeAnt(ant)
		destination.putAnt(ant)

		if bringFood {
			moveFood(origin, destination)
		}

		origin.update()
		destination.update()
	})
}

func (b *Board) CheckForWin() bool {
	nTeamsWithAnts := 0
	for _, team := range b.teams {
		if team.ants.Len() > 0 {
			nTeamsWithAnts++
			if nTeamsWithAnts > 1 {
				return false
			}
		}
	}
	return true
}

func (b *Board) TeamsByRank() []*Team {
	teams := make(Teams, len(b.teams))
	copy(teams, b.teams)
	sort.Sort(ByRank{teams})
	return teams
}

func (b *Board) Update(t *Tile) {
	if b.onUpdate != nil {
		b.onUpdate(t.pos)
	}
}

func (b *Board) Width() int {
	return b.width
}

func (b *Board) Height() int {
	return b.height
}

func (board *Board) createGrid() {
	board.columns = make([]column, board.width)
	for x := 0; x < board.width; x++ {
		board.columns[x] = make([]*Tile, board.height)
	}
}

func (b *Board) createTiles() {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			tile := new(Tile)
			tile.board = b
			b.columns[x][y] = tile
			tile.ants = NewAntSet(20)
			tile.pos = Pos{x, y}
		}
	}
}

func (board *Board) linkTiles() {
	for x := 0; x < board.width; x++ {
		for y := 0; y < board.height; y++ {
			pos := Pos{x, y}
			here := board.At(pos)
			here.north = board.At(board.NorthFrom(pos))
			here.east = board.At(board.EastFrom(pos))
			here.south = board.At(board.SouthFrom(pos))
			here.west = board.At(board.WestFrom(pos))
		}
	}
}

func (board *Board) takeTurn() {
	board.SpawnAnts()
	board.SpawnFoodRandomly(nTilesToReceiveExtraFood, nFoodToPutOnTiles)
	board.MoveAnts()
}

func NewBoard(width, height int) (board *Board) {
	board = new(Board)
	board.width = width
	board.height = height
	board.ants = NewAntSet(10000)
	board.createGrid()
	board.createTiles()
	board.linkTiles()
	return
}
