package naive_ant

import (
	"rand"
	"antwar"
	"fmt"
)

type pos struct {
	x,y int;
}

type naiveAnt struct {
	direction antwar.Action
	pos
	turn int
}

func (a *naiveAnt) directionHome() antwar.Action {
	p := (*a).pos;
	if p.x > 0 {
		return antwar.WEST
	}
	if p.y > 0 {
		return antwar.SOUTH
	}
	if p.x < 0 {
		return antwar.EAST
	}
	if p.y < 0 {
		return antwar.NORTH
	}
	return antwar.HERE
}

func (a *naiveAnt) update(decision antwar.Action) {
	if decision == antwar.WEST {
		a.x--
	}
	if decision == antwar.EAST {
		a.x++
	}
	if decision == antwar.SOUTH {
		a.y++
	}
	if decision == antwar.NORTH {
		a.y--
	}
	a.turn++
}

func oppositeDirectionOf(d antwar.Action) (opposite antwar.Action) {
	switch d {
	case antwar.NORTH: opposite = antwar.SOUTH
	case antwar.SOUTH: opposite = antwar.NORTH
	case antwar.EAST: opposite = antwar.WEST
	case antwar.WEST: opposite = antwar.EAST
	default: opposite = (antwar.Action)(rand.Intn(4) + 1)
	}
	return
}

func (a *naiveAnt) Decide(env antwar.Environment) (decision antwar.Action, bringFood bool) {
	if env[0].Food > 0 {
		decision = a.directionHome();
	} else if env[1].Food > 0 {
		fmt.Printf("Food seen to the NORTH")
		decision = antwar.NORTH
	} else if env[2].Food > 0 {
		fmt.Printf("Food seen to the EAST")
		decision = antwar.EAST
	} else if env[3].Food > 0 {
		fmt.Printf("Food seen to the SOUTH")
		decision = antwar.SOUTH
	} else if env[4].Food > 0 {
		fmt.Printf("Food seen to the WEST")
		decision = antwar.WEST
	} else {
		decision = oppositeDirectionOf(a.directionHome());
	}
	a.update(decision)
	bringFood = env[0].Food > 0
	return
}

func Spawn () antwar.AntBrain { return new(naiveAnt) }
