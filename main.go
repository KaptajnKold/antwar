package main

import (
	"antwar"
	"ants/random_ant"
	"container/vector"
	"exp/gui/x11" 
	"fmt"
)

type AntSpawner (func() antwar.AntBrain);

func Mod(n, d int) int {
	return ((n % d) + d) % d;
}


func main() {
	teams := map[string]AntSpawner{
		"randomAnt": random_ant.Spawn,
	}
	ants := new(vector.Vector);
	bases := new(vector.Vector);
	board := new(antwar.Board);
	win, err := x11.NewWindow()
	if (err != nil) {
		println(err);
	}
	// create board
	// create bases
	
	// turn:
	//   - create initial food
	//   - spawn new ants
	//   - move ants
	//   - draw board
	//   - print statistics
	
	for name, _ := range teams {
		base := antwar.Base{name, antwar.RandomPos()}
		bases.Push(base)
		board.Tiles[base.X][base.Y].Base = true;
		board.Tiles[base.X][base.Y].Team = name;
	}
	
	for i := 0; i < 30; i++ {
		for j := 0; j < bases.Len(); j++ {
			base, _ := bases.At(j).(antwar.Base);
			ant := &antwar.Ant{teams[base.Team](), base.Pos}
			ants.Push(ant)
			board.Tiles[base.X][base.Y].Ants++
		}
	}
	
	for i := 0; i < 100000; i++ {
		for j := 0; j < ants.Len(); j++ {
			ant, _ := ants.At(j).(*antwar.Ant)
			decision := ant.Brain.Decide(nil)
			pos := ant.Pos
			if (decision == antwar.NORTH) {
				pos.Y = Mod(pos.Y + 1, 600)
			}
			if (decision == antwar.SOUTH) {
				pos.Y = Mod(pos.Y - 1, 600)
			}
			if (decision == antwar.EAST) {
				pos.X = Mod(pos.X + 1, 800)
			}
			if (decision == antwar.WEST) {
				pos.X = Mod(pos.X - 1,800)
			}
			
			board.Tiles[ant.X][ant.Y].Ants--
			win.Screen().Set(ant.X, ant.Y, board.Tiles[ant.X][ant.Y].Color());
			ant.Pos = pos;
			board.Tiles[ant.X][ant.Y].Ants++
			win.Screen().Set(ant.X, ant.Y, board.Tiles[ant.X][ant.Y].Color());
		}
		win.FlushImage()
		fmt.Println("turn over")
	}
	
	fmt.Println("Game over")
	
	<-win.EventChan();
	win.Close();
	println("Done");
}