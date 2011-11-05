package naive_ant

import (
	"rand"
	"antwar"
)

type pos struct {
	x,y int;
}

type naiveAnt struct {
	direction antwar.Action
	pos
	turn int
}

func flipCoin() bool {
	if rand.Intn(2) == 0 {
		return true;
	}
	return false;
}

func (me *naiveAnt) directionHome() (d antwar.Action) {
	horizontal := antwar.HERE
	vertical := antwar.HERE
	if me.x > 0 {
		horizontal = antwar.WEST
	}
	if me.x < 0 {
		horizontal = antwar.EAST
	}
	if me.y > 0 {
		vertical = antwar.NORTH
	}
	if me.y < 0 {
		vertical = antwar.SOUTH
	}
	if horizontal == antwar.HERE && vertical == antwar.HERE {
		d = antwar.HERE
	} else if horizontal == antwar.HERE {
		d = vertical
	} else if vertical == antwar.HERE {
		d = horizontal
	} else if flipCoin() {
		d = vertical
	} else {
		d = horizontal
	}
	return;
}

func (me *naiveAnt) directionOut() (d antwar.Action) {
	var horizontal,	vertical antwar.Action
	
	if me.x < 0 {
		horizontal = antwar.WEST
	} else if me.x > 0 {
		horizontal = antwar.EAST
	} else if flipCoin() {
		horizontal = antwar.WEST
	} else {
		horizontal = antwar.EAST
	}

	if me.y < 0 {
		vertical = antwar.NORTH
	} else if me.y > 0 {
		vertical = antwar.SOUTH
	} else if flipCoin() {
		vertical = antwar.NORTH
	} else {
		vertical = antwar.SOUTH
	}
	
	if flipCoin() {
		d = horizontal
	} else {
		d = vertical
	}
	return
}

func (a *naiveAnt) update(decision antwar.Action) {
	switch decision {
	case antwar.NORTH: a.y--
	case antwar.SOUTH: a.y++
	case antwar.EAST: a.x++
	case antwar.WEST: a.x--
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

func (a *naiveAnt) Decide(env *antwar.Environment) (decision antwar.Action, bringFood bool) {
	if env[0].FoodCount() > 0 {
		decision = a.directionHome();
	} else if env[1].FoodCount() > 0 {
		decision = antwar.NORTH
	} else if env[2].FoodCount() > 0 {
		decision = antwar.EAST
	} else if env[3].FoodCount() > 0 {
		decision = antwar.SOUTH
	} else if env[4].FoodCount() > 0 {
		decision = antwar.WEST
	} else {
		decision = a.directionOut();
	}
	a.update(decision)
	bringFood = env[0].FoodCount() > 0
	return
}

func Spawn () antwar.AntBrain { return new(naiveAnt) }
