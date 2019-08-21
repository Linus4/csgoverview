// Package common contains types that are used throughout this project.
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

// GrenadeEffect extends the GrenadeEvent type from the parser by the Lifetime
// variable that is used to draw the effect.
type GrenadeEffect struct {
	event.GrenadeEvent
	Lifetime int
}

// Kill contains all information that is displayed on the killfeed.
type Kill struct {
	KillerName string
	KillerTeam common.Team
	VictimName string
	VictimTeam common.Team
	Weapon     string
}
