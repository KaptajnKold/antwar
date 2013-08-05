package antwar

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
	"flag"
)

const (
	nTurnsMax                = 50000
	nStartingAnts            = 50
	nTilesStartingWithFood   = 500
	nFoodOnStartingTile      = 50
	nTilesToReceiveExtraFood = 2
	nFoodToPutOnTiles        = 70
)

var (
	width, height int
)

func printStats(teams []*Team, turn int) {
	fmt.Printf("\nTurn: %v\n", turn)
	for _, team := range teams {
		fmt.Printf("%v: %v\n", team.name, team.ants.Len())
	}
	fmt.Printf("\x1b[%vA", len(teams)+2)
}

func printTeams(teams Teams) {
	for i, team := range teams {
		fmt.Printf("%v. %v\n", i+1, team.name)
	}
}

func NewGame(teams []*Team) {
	runtime.GOMAXPROCS(2)
	rand.Seed(time.Now().UTC().UnixNano())

	flag.IntVar(&width, "width", 500, "Width of game board")
	flag.IntVar(&height, "height", 500, "Height of game board")
	flag.Parse()

	board := NewBoard(width, height)
	NewGUI(board)

	board.teams = teams

	board.createStartingAntHills()
	board.CreateStartingAnts(nStartingAnts)
	board.SpawnFoodRandomly(nTilesStartingWithFood, nFoodOnStartingTile)

	for i := 0; i < nTurnsMax; i++ {
		board.takeTurn()
		printStats(board.teams, i)
		if board.CheckForWin() {
			break
		}

	}
	printTeams(board.TeamsByRank())
}
