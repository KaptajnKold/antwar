package antwar

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

const (
	nTurnsMax                = 50000
	nStartingAnts            = 1
	nTilesStartingWithFood   = 50
	nFoodOnStartingTile      = 50
	nTilesToReceiveExtraFood = 3
	nFoodToPutOnTiles        = 10
	boardWidth               = 1000
	boardHeight              = 1000
)

func printStats(teams []*Team, turn int) {
	fmt.Printf("\nTurn: %v\n", turn)
	for _, team := range teams {
		fmt.Printf("%v: %v\n", team.Name, team.Ants.Len())
	}
	fmt.Printf("\x1b[%vA", len(teams)+2)
}

func printTeams(teams Teams) {
	for i, team := range teams {
		fmt.Printf("%v. %v\n", i+1, team.Name)
	}
}

func NewGame(teams []*Team) {
	runtime.GOMAXPROCS(2)
	rand.Seed(time.Now().UTC().UnixNano())

	board := NewBoard(boardWidth, boardHeight)
	NewGUI(board)

	board.Teams = teams

	board.createStartingAntHills()
	board.CreateStartingAnts(nStartingAnts)
	board.SpawnFoodRandomly(nTilesStartingWithFood, nFoodOnStartingTile)

	for i := 0; i < nTurnsMax; i++ {
		board.takeTurn()
		printStats(board.Teams, i)
		if board.CheckForWin() {
			break
		}

	}
	printTeams(board.TeamsByRank())
}
