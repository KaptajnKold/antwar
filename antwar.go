package antwar

type Action int;

type Base struct {
	Team string
	Pos
}

type AntBrain interface {
	Decide(env *Environment) (Action, bool)
}

type Ant struct {
	Brain AntBrain
	Team string
	Pos
}

