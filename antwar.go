package antwar

import "image"

type Action int;

type AntBrain interface {
	Decide(env *Environment) (Action, bool)
}

type Ant struct {
	Brain AntBrain
	Team *Team
	Pos
}

type AntSpawner (func() AntBrain);

type AntSet map[*Ant]*Ant

func (s AntSet) Put(a *Ant) {
	s[a] = a
}

func (s AntSet) Remove(a *Ant) {
	s[a] = a, false
}

func (s AntSet) Do(f func(a *Ant)) {
	for _, ant := range(s) {
		f(ant)
	}
}

func (s AntSet) Len() int {
	return len(s)
}

func NewAntSet(capacity int) AntSet {
	s := make(AntSet, capacity)
	return s
}

type AntHill struct {
	Team *Team
	Pos
	Ants AntSet
}

func NewAntHill(team *Team, pos Pos) *AntHill {
	return &AntHill{team, pos, NewAntSet(5000)}
}

var colorIndex = 0

type Team struct {
	Name string
	Ants AntSet
	Spawn AntSpawner
	Color image.Color
}

func NewTeam(name string, spawn AntSpawner) *Team {
	teamColors := [...]image.RGBAColor{
		image.RGBAColor{255,128,128,255},
		image.RGBAColor{128,255,128,255},
		image.RGBAColor{128,128,255,255},
		image.RGBAColor{127,128,127,255},
		image.RGBAColor{127,127,128,255},
		image.RGBAColor{128,127,127,255},
	}
	color := teamColors[colorIndex]
	colorIndex++
	return &Team{name, NewAntSet(5000), spawn, color}
}
