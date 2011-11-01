package main

import (
	"antwar"
	"ants/random_ant"
	"ants/naive_ant"
	"container/vector"
	"fmt"
)

type AntSpawner (func() antwar.AntBrain);

func Mod(n, d int) int {
	return ((n % d) + d) % d;
}


func main() {
	teams := map[string]AntSpawner{
		"randomAnt": random_ant.Spawn,
		"naiveAnt": naive_ant.Spawn,
	}
	ants := new(vector.Vector);
	bases := new(vector.Vector);
	board := antwar.NewBoard();
	gui := antwar.NewGUI(board);
	
	for name, _ := range teams {
		base := antwar.Base{name, antwar.RandomPos()}
		bases.Push(base)
		board.Tiles[base.X][base.Y].Base = true;
		board.Tiles[base.X][base.Y].Team = name;
	}
	
	for i := 0; i < 10; i++ {
		for j := 0; j < bases.Len(); j++ {
			base, _ := bases.At(j).(antwar.Base);
			ant := &antwar.Ant{teams[base.Team](), base.Team, base.Pos}
			ants.Push(ant)
			board.At(base.Pos).Ants++
		}
	}
	
	board.CreateFood(4000)
	
	for i := 0; i < 1000; i++ {
		// Spawn new ants
		bases.Do(func (b interface{}) {
			base, _ := b.(antwar.Base)
			tile := board.At(base.Pos)
			for ; 0 < tile.Food; {
				ant := &antwar.Ant{teams[base.Team](), base.Team, base.Pos}
				ants.Push(ant)
				tile := board.At(base.Pos)
				tile.Ants++
				tile.Food--
				fmt.Printf("+");
			}
		})
		
		board.CreateFood(5);
		
		for j := 0; j < ants.Len(); j++ {
			ant, bringFood := ants.At(j).(*antwar.Ant)
			pos := ant.Pos
			env := board.Environment(pos)
			decision, _ := ant.Brain.Decide(env)
			switch decision {
			case antwar.NORTH: pos = pos.North()
			case antwar.SOUTH: pos = pos.South()
			case antwar.EAST: pos = pos.East()
			case antwar.WEST: pos = pos.West()
			}
			fromTile := board.At(ant.Pos)
			fromTile.Ants--
			bringFood = bringFood && fromTile.Food > 0
			if bringFood {
				fromTile.Food--
			}
			board.Update(ant.Pos)
			ant.Pos = pos;
			toTile := board.At(ant.Pos)
			toTile.Ants++
			toTile.Team = ant.Team
			if bringFood {
				toTile.Food++
			}
			board.Update(ant.Pos)
		}
	}
	fmt.Println("Game over")
	gui.Close();
}