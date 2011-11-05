package main

import (
	"antwar"
	"ants/random_ant"
	"ants/naive_ant"
	"container/vector"
	"fmt"
	"os"
)

type timeKeeper struct {
	start, nanos, count int64
}

func (t *timeKeeper) begin() {
	secs, nsecs, _ := os.Time()
	t.start = 1e9*secs+nsecs
}

func (t *timeKeeper) end() {
	secs, nsecs, _ := os.Time()
	stopped :=  1e9*secs+nsecs - t.start
	t.count++
	t.nanos = t.nanos + stopped
}

func (t *timeKeeper) String() {
	fmt.Printf("count:%v avg:%v total:%v \n", t.count, t.nanos/t.count, t.nanos)
}

type AntSpawner (func() antwar.AntBrain);

func main() {
	timers := make(map[string]*timeKeeper)
	timers["decide"] = new(timeKeeper)
	timers["move"] = new(timeKeeper)
	
	teams := map[string]AntSpawner{
		"randomAnt": random_ant.Spawn,
		"naiveAnt": naive_ant.Spawn,
		"cleverAnt": naive_ant.Spawn,
	}
	ants := new(vector.Vector);
	bases := new(vector.Vector);
	board := antwar.NewBoard();
	gui := antwar.NewGUI(board);
	defer gui.Close()
	
	for name, _ := range teams {
		base := antwar.Base{name, antwar.RandomPos()}
		bases.Push(base)
		board.At(base.Pos).CreateBase(name);
	}
	
	for i := 0; i < 10; i++ {
		for j := 0; j < bases.Len(); j++ {
			base, _ := bases.At(j).(antwar.Base);
			ant := &antwar.Ant{teams[base.Team](), base.Team, base.Pos}
			ants.Push(ant)
			board.At(base.Pos).PutAnt(ant)
		}
	}
	
	board.CreateFood(100)
	
	for i := 0; i < 100000; i++ {
		bases.Do(func (b interface{}) {
			base, _ := b.(antwar.Base)
			tile := board.At(base.Pos)
			for ; 0 < tile.FoodCount(); {
				ant := &antwar.Ant{teams[base.Team](), base.Team, base.Pos}
				ants.Push(ant)
				tile := board.At(base.Pos)
				tile.PutAnt(ant)
				tile.RemoveFood(1)
			}
		})
		
		board.CreateFood(1);
		
		ants.Do(func (a interface{}){
			
			ant, _ := a.(*antwar.Ant)
			origin := ant.Pos
			env := board.Environment(origin)
			destination := origin
			
			decision, bringFood := ant.Brain.Decide(env)
			switch decision {
			case antwar.NORTH: destination = origin.North()
			case antwar.SOUTH: destination = origin.South()
			case antwar.EAST: destination = origin.East()
			case antwar.WEST: destination = origin.West()
			case antwar.HERE: return
			}
			
			fromTile := board.At(ant.Pos)
			fromTile.RemoveAnt(ant)
			
			ant.Pos = destination;
			toTile := board.At(ant.Pos)

			toTile.PutAnt(ant)

			if bringFood && fromTile.FoodCount() > 0{
				fromTile.RemoveFood(1)
				toTile.PutFood(1)
			}
			
			board.Update(origin)
			board.Update(destination)
		})
		fmt.Printf("# Ants: %v \n", ants.Len());
	}
	
	for name, stats := range(timers) {
		fmt.Println(name + ": ")
		stats.String()
	}
	fmt.Println("Game over")
}