package random_ant

import (
	"rand"
	"antwar"
)

type randomAnt struct {
	direction antwar.Action
}

func (a *randomAnt) Decide(env *antwar.Environment) antwar.Action {
	return (antwar.Action)(rand.Intn(4) + 1)
}

func Spawn () antwar.AntBrain { return &randomAnt{(antwar.Action)(rand.Intn(4) + 1)}}
