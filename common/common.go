package common

import (
	"github.com/markus-wa/demoinfocs-golang/common"
	event "github.com/markus-wa/demoinfocs-golang/events"
)

type OverviewState struct {
	IngameTick            int
	Players               []OverviewPlayer
	Grenades              []common.GrenadeProjectile
	Infernos              []common.Inferno
	Bomb                  common.Bomb
	TeamCounterTerrorists common.TeamState
	TeamTerrorists        common.TeamState
}

type GrenadeEffect struct {
	event.GrenadeEvent
	Lifetime int
}

// Do not use Weapons(), but do use Weapons instead
type OverviewPlayer struct {
	common.Player
	Weapons []common.Equipment
}
