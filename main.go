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
	}
	ants := new(vector.Vector);
	bases := new(vector.Vector);
	board := antwar.NewBoard();
	gui := antwar.NewGUI(board);
	defer gui.Close()
	
	for name, _ := range teams {
		base := antwar.Base{name, antwar.RandomPos()}
		bases.Push(base)
		board.At(base.Pos).Base = true;
		board.At(base.Pos).Team = name;
	}
	
	for i := 0; i < 10; i++ {
		for j := 0; j < bases.Len(); j++ {
			base, _ := bases.At(j).(antwar.Base);
			ant := &antwar.Ant{teams[base.Team](), base.Team, base.Pos}
			ants.Push(ant)
			board.At(base.Pos).Ants++
		}
	}
	
	board.CreateFood(400)
	
	for i := 0; i < 100000; i++ {
		bases.Do(func (b interface{}) {
			base, _ := b.(antwar.Base)
			tile := board.At(base.Pos)
			for ; 0 < tile.Food; {
				ant := &antwar.Ant{teams[base.Team](), base.Team, base.Pos}
				ants.Push(ant)
				tile := board.At(base.Pos)
				tile.Ants++
				tile.Food--
			}
		})
		
		board.CreateFood(1);
		
		ants.Do(func (a interface{}){
			
			ant, _ := a.(*antwar.Ant)
			pos := ant.Pos
			timers["move"].begin()
			env := board.Environment(pos)
			timers["move"].end()
			
			timers["decide"].begin()
			decision, bringFood := ant.Brain.Decide(env)
			timers["decide"].end()
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
		})
	}
	for name, stats := range(timers) {
		fmt.Println(name + ": ")
		stats.String()
	}
	fmt.Println("Game over")
}