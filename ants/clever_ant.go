package clever_ant

import (
	"rand"
	"antwar"
)

type pos struct {
	x,y int;
}

var ENV_POS [5]pos = [5]pos{pos{0, 0}, pos{0, -1}, pos{1, 0}, pos{0, 1}, pos{-1, 0}}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (this *pos) Equals(other *pos) bool {
	return this.x == other.x && this.y == other.y
}

func (self* pos) distanceTo(other *pos) int {
	return abs(self.x - other.x) + abs(self.y - other.y)
}

func (self* pos) length() int {
	return abs(self.x) + abs(self.y)
}


type cleverAnt struct {
	direction antwar.Action
	home, food pos
	turn int
}

func flipCoin() bool {
	if rand.Intn(2) == 0 {
		return true;
	}
	return false;
}

func (p pos) direction() (d antwar.Action) {
	horizontal := antwar.HERE
	vertical := antwar.HERE
	if p.x > 0 {
		horizontal = antwar.WEST
	}
	if p.x < 0 {
		horizontal = antwar.EAST
	}
	if p.y > 0 {
		vertical = antwar.NORTH
	}
	if p.y < 0 {
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

func (me *cleverAnt) directionOut() (d antwar.Action) {
	var horizontal,	vertical antwar.Action
	home := me.home
	if home.x < 0 {
		horizontal = antwar.WEST
	} else if home.x > 0 {
		horizontal = antwar.EAST
	} else if flipCoin() {
		horizontal = antwar.WEST
	} else {
		horizontal = antwar.EAST
	}

	if home.y < 0 {
		vertical = antwar.NORTH
	} else if home.y > 0 {
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

func (a *cleverAnt) update(decision antwar.Action) {
	a.home.update(decision)
	a.food.update(decision)
	a.turn++
}

func (p *pos) update(decision antwar.Action) {
	switch decision {
	case antwar.NORTH: p.y--
	case antwar.SOUTH: p.y++
	case antwar.EAST: p.x++
	case antwar.WEST: p.x--
	}
}

func (from *pos) inDirectionOf(decision antwar.Action) (p pos) {
	p = *from
	p.update(decision);
	return
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

func (a *cleverAnt) LookForFood(e *antwar.Environment) {
	var foodPos *pos
	var shortestSoFar int
	for i := 0; i < len(e); i++ {
		tile := e[i]
		envPos := ENV_POS[i];
		candidate := &envPos
		candidateLen := envPos.distanceTo(&a.home);
		// TODO: Obv. check if ants belong to other team
		if tile.FoodCount() > tile.AntCount() && (foodPos == nil || candidateLen < shortestSoFar) {
			foodPos = candidate
			shortestSoFar = candidateLen
		}
	}
	if foodPos != nil {
		a.food.Equals(foodPos) 
	}
}

func (a *cleverAnt) Decide(env *antwar.Environment) (decision antwar.Action, bringFood bool) {
	a.LookForFood(env);
	
	if env[0].FoodCount() > 0 {
		decision = a.home.direction();
	} else if (a.food.Equals(&a.home)) {
		decision = a.directionOut();
	} else {
		decision = a.food.direction();
	}
	
	a.update(decision)
	bringFood = env[0].FoodCount() > 0
	return
}

func Spawn () antwar.AntBrain { return new(cleverAnt) }

