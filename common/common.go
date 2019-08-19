package common

import (
	"github.com/markus-wa/demoinfocs-golang/common"
	event "github.com/markus-wa/demoinfocs-golang/events"
)

// OverviewState contains all information that will be displayed for a single tick.
type OverviewState struct {
	IngameTick            int
	Players               []common.Player
	Grenades              []common.GrenadeProjectile
	Infernos              []common.Inferno
	Bomb                  common.Bomb
	TeamCounterTerrorists common.TeamState
	TeamTerrorists        common.TeamState
}

// GrenadeEffect extends by the Lifetime variable that is used to draw the effect.
type GrenadeEffect struct {
	event.GrenadeEvent
	Lifetime int
}
