// Package common contains types that are used throughout this project.
package common

import (
	"time"

	demoinfo "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
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
	Players               []Player
	Grenades              []GrenadeProjectile
	Infernos              []Inferno
	Bomb                  Bomb
	TeamCounterTerrorists TeamState
	TeamTerrorists        TeamState
	Timer                 Timer
}

// GrenadeEffect extends the GrenadeEvent type from the parser by the Lifetime
// variable that is used to draw the effect.
type GrenadeEffect struct {
	Position    Point
	GrenadeType demoinfo.EquipmentType
	Lifetime    int
}

// GrenadeProjectile conains all information that is used to draw a grenade
// mid air on the map.
type GrenadeProjectile struct {
	Position Point
	Type     demoinfo.EquipmentType
}

// Kill contains all information that is displayed on the killfeed.
type Kill struct {
	KillerName string
	KillerTeam demoinfo.Team
	VictimName string
	VictimTeam demoinfo.Team
	Weapon     string
}

// Timer contains the time remaining in the current phase of the round.
type Timer struct {
	TimeRemaining time.Duration
	Phase         Phase
}

// Shot contains information about a shot from a weapon.
type Shot struct {
	Position       Point
	ViewDirectionX float32
	IsAwpShot      bool
}

// Inferno contains the hull points of the surface area of a molotov or
// incendiary grenade.
type Inferno struct {
	ConvexHull2D []Point
}

// Bomb contains all relevant information about the C4.
type Bomb struct {
	Position       Point
	IsBeingCarried bool
}

// Player contains all relevant information about a player in the match.
type Player struct {
	Name               string
	SteamID64          uint64
	Team               demoinfo.Team
	Position           Point
	LastAlivePosition  Point
	ViewDirectionX     float32
	FlashDuration      time.Duration
	FlashTimeRemaining time.Duration
	Inventory          []demoinfo.EquipmentType
	Health             int16
	Armor              int16
	Money              int16
	Kills              int16
	Deaths             int16
	Assists            int16
	IsAlive            bool
	IsDefusing         bool
	HasHelmet          bool
	HasDefuseKit       bool
	HasBomb            bool
}

// TeamState contains information about a team in the match.
type TeamState struct {
	ClanName string
	Score    byte
}

// Point contains the coordinates for a point on the map.
type Point struct {
	X float32
	Y float32
}
