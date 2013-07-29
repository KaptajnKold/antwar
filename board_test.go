package antwar

import (
	"testing"
)

type dummyAnt struct {
}

var (
	decideWasCalled bool
)

func (me *dummyAnt) Decide(env *Environment, brains []AntBrain) (Action, bool) {
	decideWasCalled = true
	return HERE, false
}

func dummySpawn() AntBrain { return new(dummyAnt) }

func createBoard() *Board {
	board := NewBoard(100, 100)
	board.Teams = []*Team{
		NewTeam("dummy", dummySpawn),
		NewTeam("dummy2", dummySpawn),
	}
	return board
}

func TestCreateStartingAntHills(t *testing.T) {
	board := createBoard()
	decideWasCalled = false
	board.CreateStartingAntHills()
	if l := len(board.AntHills); l != 2 {
		t.Errorf("len(board.AntHills) was %v. Exptected %v", l, 2)
	}
	board.CreateStartingAnts(3)
	if totalAntCount := board.Ants.Len(); totalAntCount != 6 {
		t.Errorf("board.Ants.Len() was %v. Exptected %v", totalAntCount, 6)
	}
	board.MoveAnts()
	if !decideWasCalled {
		t.Error("Exptected Decide to be called.")
	}
}

func TestCreateAntHill(t *testing.T) {
	board := createBoard()
	pos := Pos{0, 0}
	hill := board.At(pos).CreateAntHill(board.Teams[0])
	if hill.tile != board.At(pos) {
		t.Error("AntHill.tile not the same as the tile that created it")
	}
}

func TestSpawnAnts(t *testing.T) {
	board := createBoard()
	pos := Pos{0, 0}
	antHill := board.At(pos).CreateAntHill(board.Teams[0])
	board.AntHills = []*AntHill{antHill}
	board.At(pos).PutFood(1)
	board.SpawnAnts()
	if numberOfAnts := board.At(pos).Ants.Len(); numberOfAnts != 1 {
		t.Errorf("Expected one ant on tile. Was %v.", numberOfAnts)
	}
	if foodCount := board.At(pos).Food; foodCount > 0 {
		t.Errorf("Expected all food to be consumed. Was %v.", foodCount)
	}
}
