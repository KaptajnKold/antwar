package antwar

import "image/color"

type Action int

type AntBrain interface {
	Decide(env *Environment, brains []AntBrain) (Action, bool)
}

type Ant struct {
	Brain AntBrain
	Team  *Team
	Pos
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
			brains = append(brains, ant.Brain)
		}
	})
	return brains
}

func NewAntSet(capacity int) AntSet {
	s := make(AntSet, capacity)
	return s
}

type AntHill struct {
	Team *Team
	Ants AntSet
	tile *Tile
}

func (hill *AntHill) spawnAnt() *Ant {
	ant := &Ant{hill.Team.Spawn(), hill.Team, hill.tile.pos}
	hill.tile.PutAnt(ant)
	return ant
}

var colorIndex = 0

type Team struct {
	Name  string
	Ants  AntSet
	Spawn AntSpawner
	Color color.Color
}

type Teams []*Team

func (self *Team) RanksLowerThan(other *Team) bool {
	return self.Ants.Len() < other.Ants.Len()
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
