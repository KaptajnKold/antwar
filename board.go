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

type Environment struct {
	Here, North, East, South, West *TileInfo
}

type Board struct {
	width, height int
	columns       []column
	OnUpdate      func(p Pos)
	Ants          AntSet
	AntHills      []*AntHill
	Teams         []*Team
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
	numberOfTeams := len(board.Teams)
	board.AntHills = make([](*AntHill), numberOfTeams)
	for i, team := range board.Teams {
		board.AntHills[i] = board.randomTile().CreateAntHill(team)
		fmt.Printf("Created anthill at %v\n", board.AntHills[i].tile.pos)
	}
}

func (board *Board) destroyAntHill(antHill AntHill) {

}

func (board *Board) CreateStartingAnts(startingNumberOfAnts int) {
	// Add 1 ant at a time to each team to avoid one team's ants all get to move before the other teams' ants
	for i := 0; i < startingNumberOfAnts; i++ {
		for _, antHill := range board.AntHills {
			ant := antHill.spawnAnt()
			board.Ants.Put(ant)
			ant.Team.Ants.Put(ant)
			fmt.Printf("Created %v starting ants for team %v.\n", ant.Team.Ants.Len(), ant.Team.Name)
		}
	}
	fmt.Printf("Created %v starting ants.\n", board.Ants.Len())
}

func (b *Board) SpawnFoodRandomly(ncolumns, nFood int) {
	for i, j := 0, rand.Intn(ncolumns); i < j; i++ {
		pos := b.randomPos()
		b.At(pos).PutFood(rand.Intn(nFood))
		b.Update(pos)
	}
}

func (board *Board) SpawnAnts() {
	for _, antHill := range board.AntHills {
		for 0 < antHill.tile.Food {
			ant := antHill.spawnAnt()
			board.Ants.Put(ant)
			ant.Team.Ants.Put(ant)
			antHill.tile.RemoveFood(1)
		}
	}
}

func moveFood(fromTile, toTile *Tile) {
	if fromTile.Food > 0 {
		fromTile.RemoveFood(1)
		toTile.PutFood(1)
	}
}

func (board *Board) MoveAnts() {
	board.Ants.Do(func(ant *Ant) {
		origin := ant.Pos
		env := board.At(origin).Environment

		// TODO: Enable killing of bases
		brains := make([]AntBrain, 0) //board.At(origin).Ants.brainsExcept(ant)

		var destination Pos

		decision, bringFood := ant.Brain.Decide(env, brains)
		switch decision {
		case NORTH:
			destination = board.NorthFrom(origin)
		case SOUTH:
			destination = board.SouthFrom(origin)
		case EAST:
			destination = board.EastFrom(origin)
		case WEST:
			destination = board.WestFrom(origin)
		case HERE:
			return
		}

		fromTile := board.At(ant.Pos)
		fromTile.RemoveAnt(ant)

		ant.Pos = destination
		toTile := board.At(ant.Pos)

		if 0 < toTile.Ants.Len() && toTile.Team != ant.Team {
			toTile.Ants.Do(func(anAnt *Ant) {
				board.Ants.Remove(anAnt)
				toTile.Ants.Remove(anAnt)
				anAnt.Team.Ants.Remove(anAnt)
			})
		}

		toTile.PutAnt(ant)

		if bringFood {
			moveFood(fromTile, toTile)
		}

		board.Update(origin)
		board.Update(destination)
	})
}

func (b *Board) CheckForWin() bool {
	nTeamsWithAnts := 0
	for _, team := range b.Teams {
		if team.Ants.Len() > 0 {
			nTeamsWithAnts++
			if nTeamsWithAnts > 1 {
				return false
			}
		}
	}
	return true
}

func (b *Board) TeamsByRank() []*Team {
	teams := make(Teams, len(b.Teams))
	copy(teams, b.Teams)
	sort.Sort(ByRank{teams})
	return teams
}

func (b *Board) Update(pos Pos) {
	if b.OnUpdate != nil {
		b.OnUpdate(pos)
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
			tile.Ants = NewAntSet(20)
			tile.environmentTile = &TileInfo{tile}
			tile.Environment = new(Environment)
			tile.pos = Pos{x, y}
		}
	}
}

func (board *Board) linkEnvironmentTiles() {
	for x := 0; x < board.width; x++ {
		for y := 0; y < board.height; y++ {
			pos := Pos{x, y}
			env := board.At(pos).Environment
			env.Here = board.At(pos).environmentTile
			env.North = board.At(board.NorthFrom(pos)).environmentTile
			env.East = board.At(board.EastFrom(pos)).environmentTile
			env.South = board.At(board.SouthFrom(pos)).environmentTile
			env.West = board.At(board.WestFrom(pos)).environmentTile
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
	board.Ants = NewAntSet(10000)
	board.createGrid()
	board.createTiles()
	board.linkEnvironmentTiles()
	return
}
