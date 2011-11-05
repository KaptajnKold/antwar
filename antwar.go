package antwar

type Action int;

type AntBrain interface {
	Decide(env *Environment) (Action, bool)
}

type Ant struct {
	Brain AntBrain
	Team string
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
	Team string
	Pos
	Ants AntSet
}

func NewAntHill(team string, pos Pos) *AntHill {
	return &AntHill{team, pos, NewAntSet(5000)}
}

type Team struct {
	Name string
	Ants AntSet
	Spawn AntSpawner
}
