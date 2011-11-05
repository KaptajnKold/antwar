package main

import (
	"antwar"
	"ants/random_ant"
	"ants/naive_ant"
	"ants/clever_ant"
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

func printStats(teams map[string]antwar.Team) {
	fmt.Printf("\x1b[2J");
	for name, team := range(teams) {
		fmt.Printf("%v: %v\n", name, team.Ants.Len())
	}
}

func main() {
	fmt.Println("Start game…")
	timers := make(map[string]*timeKeeper)
	timers["decide"] = new(timeKeeper)
	timers["move"] = new(timeKeeper)
	
	teams := map[string]antwar.Team{
		"randomAnt": antwar.Team{"randomAnt", antwar.NewAntSet(5000), random_ant.Spawn},
		"naiveAnt": antwar.Team{"randomAnt", antwar.NewAntSet(5000), naive_ant.Spawn},
		"cleverAnt": antwar.Team{"randomAnt", antwar.NewAntSet(5000), clever_ant.Spawn},
	}
	antHills := new(vector.Vector);
	board := antwar.NewBoard();
	gui := antwar.NewGUI(board);
	defer gui.Close()
	
	// Create starting antHill for each team
	fmt.Println("Creating starting ant hills…")
	for name, _ := range teams {
		antHill := antwar.NewAntHill(name, antwar.RandomPos())
		antHills.Push(antHill)
		board.At(antHill.Pos).CreateAntHill(name);
	}
	
	// TODO: Make starting number of ants a command line parameter
	fmt.Println("Creating starting ants…")
	for i := 0; i < 10; i++ {		
		for j := 0; j < antHills.Len(); j++ {
			fmt.Println("Spawn ant…")
			antHill, _ := antHills.At(j).(*antwar.AntHill)
			fmt.Println("%v", antHill)
			ant := &antwar.Ant{teams[antHill.Team].Spawn(), antHill.Team, antHill.Pos}
			board.Ants.Put(ant)
			fmt.Println("Putting ant on hill…")
			teams[antHill.Team].Ants.Put(ant)
			fmt.Println("Done putting ant on hill!")
			board.At(antHill.Pos).PutAnt(ant)
		}
	}
	
	board.CreateFood(100)

	fmt.Println("Starting main loop…")
	for i := 0; i < 100000; i++ {
		antHills.Do(func (b interface{}) {
			antHill, _ := b.(*antwar.AntHill)
			tile := board.At(antHill.Pos)
			for ; 0 < tile.FoodCount(); {
				ant := &antwar.Ant{teams[antHill.Team].Spawn(), antHill.Team, antHill.Pos}
				board.Ants.Put(ant)
				tile := board.At(antHill.Pos)
				teams[antHill.Team].Ants.Put(ant)
				tile.PutAnt(ant)
				tile.RemoveFood(1)
			}
		})
		
		board.CreateFood(1);
		
		board.Ants.Do(func (ant *antwar.Ant){			
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

			if 0 < toTile.AntCount() && toTile.Team != ant.Team {
				toTile.Ants.Do(func (anAnt *antwar.Ant) {
					board.Ants.Remove(anAnt)
					toTile.Ants.Remove(anAnt)
					teams[anAnt.Team].Ants.Remove(anAnt)
				})
			}

			toTile.PutAnt(ant)

			if bringFood && fromTile.FoodCount() > 0{
				fromTile.RemoveFood(1)
				toTile.PutFood(1)
			}
			
			board.Update(origin)
			board.Update(destination)
		})
		printStats(teams);
	}
	
	for name, stats := range(timers) {
		fmt.Println(name + ": ")
		stats.String()
	}
	fmt.Println("Game over")
}