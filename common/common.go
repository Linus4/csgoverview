// Package common contains types that are used throughout this project.
package common

import (
	"time"

	"github.com/golang/geo/r3"
	"github.com/markus-wa/demoinfocs-golang/common"
	event "github.com/markus-wa/demoinfocs-golang/events"
)

// Phase corresponds to a phase of a round.
type Phase int

// Possible values for Phase type.
const (
	PhaseFreezetime Phase = iota
	PhaseRegular
	PhasePlanted
	PhaseRestart
	PhaseWarmup
	PhaseHalftime
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
	Timer                 Timer
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

// Timer contains the time remaining in the current phase of the round.
type Timer struct {
	TimeRemaining time.Duration
	Phase         Phase
}

// Shot contains information about a shot from a weapon.
type Shot struct {
	Position       r3.Vector
	ViewDirectionX float32
	IsAwpShot      bool
}
