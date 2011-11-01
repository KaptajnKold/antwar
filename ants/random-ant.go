package random_ant

import (
	"rand"
	"antwar"
)

type randomAnt struct {
	direction antwar.Action
}

func (a *randomAnt) Decide(env *antwar.Environment) (antwar.Action, bool) {
	return (antwar.Action)(rand.Intn(4) + 1), true
}

func Spawn () antwar.AntBrain { return &randomAnt{(antwar.Action)(rand.Intn(4) + 1)}}
