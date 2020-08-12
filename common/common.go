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

// Effect contains information about graphical effects from grenades, bombs, defuses
type Effect struct {
	Position            Point
	Type                demoinfo.EquipmentType
	Lifetime            int32
	IsOnNormalElevation bool
}

// GrenadeProjectile conains all information that is used to draw a grenade
// mid air on the map.
type GrenadeProjectile struct {
	Position            Point
	Type                demoinfo.EquipmentType
	IsOnNormalElevation bool
}

// Kill contains all information that is displayed on the killfeed.
type Kill struct {
	KillerName string
	KillerTeam demoinfo.Team
	VictimName string
	VictimTeam demoinfo.Team
	Weapon     demoinfo.EquipmentType
}

// Timer contains the time remaining in the current phase of the round.
type Timer struct {
	TimeRemaining time.Duration
	Phase         Phase
}

// Shot contains information about a shot from a weapon.
type Shot struct {
	Position            Point
	ViewDirectionX      float32
	IsAwpShot           bool
	IsOnNormalElevation bool
}

// Inferno contains the hull points of the surface area of a molotov or
// incendiary grenade.
type Inferno struct {
	ConvexHull2D        []Point
	IsOnNormalElevation bool
}

// Bomb contains all relevant information about the C4.
type Bomb struct {
	Position            Point
	IsBeingCarried      bool
	IsOnNormalElevation bool
}

// Player contains all relevant information about a player in the match.
type Player struct {
	Name                string
	ID                  int
	Team                demoinfo.Team
	Position            Point
	LastAlivePosition   Point
	ViewDirectionX      float32
	FlashDuration       time.Duration
	FlashTimeRemaining  time.Duration
	Inventory           []demoinfo.EquipmentType
	Health              int16
	Armor               int16
	Money               int16
	Kills               int16
	Deaths              int16
	Assists             int16
	IsAlive             bool
	IsDefusing          bool
	IsOnNormalElevation bool
	HasHelmet           bool
	HasDefuseKit        bool
	HasBomb             bool
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

// MapInfo contains information about maps in regards to alternate versions
// of the overview image (normal/lower/upper).
type MapInfo struct {
	AlternateOverview string
	// Height threshold to determine if a player is on the normal or on the
	// alternate version of the overview map
	HeightThreshold float64
}

// Golang Maps return default values for keys that are not in a map.
var mapInfos = map[string]MapInfo{
	"de_vertigo": MapInfo{"de_vertigo_lower.jpg", 11598},
	"de_nuke":    MapInfo{"de_nuke_lower.jpg", -550},
}

// MapHasAlternateVersion returns whether a map has an alternative overview image.
func MapHasAlternateVersion(mapName string) bool {
	if mapInfos[mapName].AlternateOverview != "" {
		return true
	} else {
		return false
	}
}

// MapGetAlternateVersion returns the filename for the alternate overview file.
func MapGetAlternateVersion(mapName string) string {
	return mapInfos[mapName].AlternateOverview
}

// MapGetHeightThreshold returns the corresponding field of the specified map.
func MapGetHeightThreshold(mapName string) float64 {
	return mapInfos[mapName].HeightThreshold
}
