package antwar

import (
    "image/color"
)

type Action int

type AntBrain interface {
	Decide(env *Tile, brains []AntBrain) (Action, bool)
}

type Ant struct {
	brain AntBrain
	team  *Team
	tile *Tile
}

type AntSpawner (func() AntBrain)

type AntSet map[*Ant]bool

func (s AntSet) Put(a *Ant) {
	s[a] = true
}

func (s AntSet) Remove(a *Ant) {
	delete(s, a)
}

func (s AntSet) Do(f func(a *Ant)) {
	for ant, isPresent := range s {
		if isPresent {
			f(ant)
		}
	}
}

func (s AntSet) Len() int {
	return len(s)
}

func (set AntSet) brainsExcept(exception *Ant) []AntBrain {
	brains := make([]AntBrain, 0, set.Len()-1)
	set.Do(func(ant *Ant) {
		if ant != exception {
			brains = append(brains, ant.brain)
		}
	})
	return brains
}

func NewAntSet(capacity int) AntSet {
	s := make(AntSet, capacity)
	return s
}

type AntHill struct {
	team *Team
	ants AntSet
	tile *Tile
}

func (h *AntHill) spawnAnt() *Ant {
	ant := &Ant{h.team.spawn(), h.team, h.tile}
	h.tile.putAnt(ant)
    h.tile.board.ants.Put(ant)
	ant.team.ants.Put(ant)
    return ant
}

func (hill *AntHill) spawnAnts() {
	for 0 < hill.tile.food {
    	hill.spawnAnt()
    	hill.tile.removeFood(1)
	}
    hill.tile.update()
}

var colorIndex = 0

type Team struct {
	name  string
	ants  AntSet
	spawn AntSpawner
	color color.Color
}

type Teams []*Team

func (self *Team) RanksLowerThan(other *Team) bool {
	return self.ants.Len() < other.ants.Len()
}

type ByRank struct {
	Teams
}

func (s Teams) Len() int { return len(s) }

func (s Teams) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByRank) Less(i, j int) bool { return s.Teams[i].RanksLowerThan(s.Teams[j]) }

func NewTeam(name string, spawn AntSpawner) *Team {
	teamColors := [...]color.RGBA{
		color.RGBA{255, 128, 128, 255},
		color.RGBA{128, 255, 128, 255},
		color.RGBA{128, 128, 255, 255},
		color.RGBA{127, 128, 127, 255},
		color.RGBA{127, 127, 128, 255},
		color.RGBA{128, 127, 127, 255},
	}
	color := teamColors[colorIndex]
	colorIndex++
	return &Team{name, NewAntSet(5000), spawn, color}
}
