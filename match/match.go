// Package match contains a high-level parser for demos.
package match

import (
	"errors"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	common "github.com/linus4/csgoverview/common"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	demoinfo "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	event "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	meta "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/metadata"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	flashEffectLifetime  int32 = 10
	heEffectLifetime     int32 = 10
	defuseEffectLifetime int32 = 45
	bombEffectLifetime   int32 = 60
	killfeedLifetime     int   = 10
	c4timer              int   = 40
)

// Match contains general information about the demo and all relevant, parsed
// data from every tick of the demo that will be displayed.
type Match struct {
	MapName              string
	MapPZero             common.Point
	MapScale             float32
	HalfStarts           []int
	RoundStarts          []int
	Effects              map[int][]common.Effect
	FrameRate            float64
	TickRate             float64
	FrameRateRounded     int
	States               []common.OverviewState
	SmokeEffectLifetime  int32
	Killfeed             map[int][]common.Kill
	Shots                map[int][]common.Shot
	currentPhase         common.Phase
	latestTimerEventTime time.Duration
	takeNthFrame         int
	// used when handling events
	currentFrame int
}

// NewMatch parses the demo at the specified path in the argument and returns a
// match.Match containing all relevant data from the demo.
// fallbackFrameRate and fallbackTickRate are used in case the values cannot be
// parsed from the demo. If they are not set, they must be -1.
func NewMatch(demoFileName string, fallbackFrameRate, fallbackTickRate float64) (*Match, error) {
	demo, err := os.Open(demoFileName)
	if err != nil {
		return nil, err
	}
	defer demo.Close()

	parser := dem.NewParser(demo)
	defer parser.Close()
	header, err := parser.ParseHeader()
	if err != nil {
		return nil, err
	}

	match := &Match{
		HalfStarts:   make([]int, 0),
		RoundStarts:  make([]int, 0),
		Effects:      make(map[int][]common.Effect),
		Killfeed:     make(map[int][]common.Kill),
		Shots:        make(map[int][]common.Shot),
		currentPhase: common.PhaseRegular,
		takeNthFrame: 1,
	}

	match.FrameRate = header.FrameRate()
	if math.IsNaN(match.FrameRate) || match.FrameRate == 0 {
		if fallbackFrameRate == -1 {
			messageBoxButtonData := []sdl.MessageBoxButtonData{
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 0,
					Text:     "128",
				},
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 1,
					Text:     "64",
				},
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 2,
					Text:     "32",
				},
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 3,
					Text:     "24",
				},
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 4,
					Text:     "Cancel",
				},
			}
			messageBoxData := sdl.MessageBoxData{
				Flags:  sdl.MESSAGEBOX_ERROR,
				Window: nil,
				Title:  "Error",
				Message: "Could not parse GOTV framerate from demo." +
					" Please choose framerate from options below.",
				Buttons:     messageBoxButtonData,
				ColorScheme: nil,
			}
			buttonid, _ := sdl.ShowMessageBox(&messageBoxData)
			switch buttonid {
			case 0:
				match.FrameRate = 128
			case 1:
				match.FrameRate = 64
			case 2:
				match.FrameRate = 32
			case 3:
				match.FrameRate = 24
			case 4:
				err := errors.New("could not parse Framerate from demo." +
					" Please provide a fallback value (command-line option -framerate)")
				return nil, err
			}
		} else {
			match.FrameRate = fallbackFrameRate
		}
	}
	match.TickRate = parser.TickRate()
	if math.IsNaN(match.TickRate) || match.TickRate == 0 {
		if fallbackTickRate == -1 {
			messageBoxButtonData := []sdl.MessageBoxButtonData{
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 0,
					Text:     "128",
				},
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 1,
					Text:     "64",
				},
				sdl.MessageBoxButtonData{
					Flags:    0,
					ButtonID: 2,
					Text:     "Cancel",
				},
			}
			messageBoxData := sdl.MessageBoxData{
				Flags:  sdl.MESSAGEBOX_ERROR,
				Window: nil,
				Title:  "Error",
				Message: "Could not parse server tickrate from demo." +
					" Please choose tickrate from options below.",
				Buttons:     messageBoxButtonData,
				ColorScheme: nil,
			}
			buttonid, _ := sdl.ShowMessageBox(&messageBoxData)
			switch buttonid {
			case 0:
				match.TickRate = 128
			case 1:
				match.TickRate = 64
			case 2:
				err := errors.New("could not parse Tickrate from demo." +
					"Please provide a fallback value (command-line option -tickrate)")
				return nil, err
			}
		} else {
			match.TickRate = fallbackTickRate
		}
	}
	match.FrameRateRounded = int(math.Round(match.FrameRate))
	if match.FrameRateRounded == 128 {
		match.takeNthFrame = 4
		match.FrameRateRounded = 32
	}
	if match.FrameRateRounded == 64 {
		match.takeNthFrame = 2
		match.FrameRateRounded = 32
	}
	match.MapName = header.MapName
	match.MapPZero = common.Point{
		X: float32(meta.MapNameToMap[match.MapName].PZero.X),
		Y: float32(meta.MapNameToMap[match.MapName].PZero.Y),
	}
	match.MapScale = float32(meta.MapNameToMap[match.MapName].Scale)
	match.SmokeEffectLifetime = int32(18 * match.FrameRateRounded)

	registerEventHandlers(parser, match)
	match.States = parseGameStates(parser, match)

	return match, nil
}

func grenadeEventHandler(lifetime int32, e event.GrenadeEvent, match *Match) {
	effectLifetime := int(lifetime)
	var isOnNormalElevation bool
	if common.MapHasAlternateVersion(match.MapName) {
		if e.Position.Z > common.MapGetHeightThreshold(match.MapName) {
			isOnNormalElevation = true
		} else {
			isOnNormalElevation = false
		}
	}
	for i := 0; i < effectLifetime; i++ {
		effect := common.Effect{
			Position: common.Point{
				X: float32(e.Position.X),
				Y: float32(e.Position.Y),
			},
			Type:                e.GrenadeType,
			Lifetime:            int32(i),
			IsOnNormalElevation: isOnNormalElevation,
		}
		effects, ok := match.Effects[match.currentFrame+i]
		if ok {
			match.Effects[match.currentFrame+i] = append(effects, effect)
		} else {
			match.Effects[match.currentFrame+i] = []common.Effect{effect}
		}
	}
}

func bombEventHandler(lifetime int32, eqType demoinfo.EquipmentType, bomb *demoinfo.Bomb, match *Match) {
	effectLifetime := int(lifetime)
	for i := 0; i < effectLifetime; i++ {
		effect := common.Effect{
			Position: common.Point{
				X: float32(bomb.Position().X),
				Y: float32(bomb.Position().Y),
			},
			Type:                eqType,
			Lifetime:            int32(i),
			IsOnNormalElevation: true,
		}
		effects, ok := match.Effects[match.currentFrame+i]
		if ok {
			match.Effects[match.currentFrame+i] = append(effects, effect)
		} else {
			match.Effects[match.currentFrame+i] = []common.Effect{effect}
		}
	}
}

func weaponFireEventHandler(e event.WeaponFire, match *Match) {
	if e.Shooter == nil {
		return
	}
	if e.Weapon.Class() == demoinfo.EqClassEquipment ||
		e.Weapon.Class() == demoinfo.EqClassGrenade ||
		e.Weapon.Class() == demoinfo.EqClassUnknown {
		return
	}
	isAwpShot := e.Weapon.Type == demoinfo.EqAWP
	var isOnNormalElevation bool
	if e.Shooter.Position().Z > common.MapGetHeightThreshold(match.MapName) {
		isOnNormalElevation = true
	} else {
		isOnNormalElevation = false
	}
	shot := common.Shot{
		Position: common.Point{
			X: float32(e.Shooter.Position().X),
			Y: float32(e.Shooter.Position().Y),
		},
		ViewDirectionX:      e.Shooter.ViewDirectionX(),
		IsAwpShot:           isAwpShot,
		IsOnNormalElevation: isOnNormalElevation,
	}
	var lifetime int

	if isAwpShot {
		lifetime = 4
	} else {
		lifetime = 1
	}
	for i := 0; i < lifetime; i++ {
		shots, ok := match.Shots[match.currentFrame+i]
		if ok {
			match.Shots[match.currentFrame+i] = append(shots, shot)
		} else {
			match.Shots[match.currentFrame+i] = []common.Shot{shot}
		}
	}
}

func registerEventHandlers(parser dem.Parser, match *Match) {
	parser.RegisterEventHandler(func(event.RoundStart) {
		match.RoundStarts = append(match.RoundStarts, match.currentFrame+1)
	})
	parser.RegisterEventHandler(func(event.MatchStart) {
		match.HalfStarts = append(match.HalfStarts, match.currentFrame+1)
	})
	parser.RegisterEventHandler(func(event.GameHalfEnded) {
		match.HalfStarts = append(match.HalfStarts, match.currentFrame+1)
	})
	parser.RegisterEventHandler(func(event.AnnouncementWinPanelMatch) {
		match.HalfStarts = append(match.HalfStarts, match.currentFrame+1)
	})
	parser.RegisterEventHandler(func(e event.WeaponFire) {
		weaponFireEventHandler(e, match)
	})
	parser.RegisterEventHandler(func(e event.FlashExplode) {
		grenadeEventHandler(flashEffectLifetime, e.GrenadeEvent, match)
	})
	parser.RegisterEventHandler(func(e event.HeExplode) {
		grenadeEventHandler(heEffectLifetime, e.GrenadeEvent, match)
	})
	parser.RegisterEventHandler(func(e event.SmokeStart) {
		grenadeEventHandler(match.SmokeEffectLifetime, e.GrenadeEvent, match)
	})
	parser.RegisterEventHandler(func(e event.BombDefused) {
		bombEventHandler(defuseEffectLifetime, demoinfo.EqDefuseKit, parser.GameState().Bomb(), match)
	})
	parser.RegisterEventHandler(func(e event.BombExplode) {
		bombEventHandler(bombEffectLifetime, demoinfo.EqBomb, parser.GameState().Bomb(), match)
	})
	parser.RegisterEventHandler(func(e event.Kill) {
		var killerName, victimName string
		var killerTeam, victimTeam demoinfo.Team
		if e.Killer == nil {
			killerName = "World"
			killerTeam = demoinfo.TeamUnassigned
		} else {
			killerName = e.Killer.Name
			killerTeam = e.Killer.Team
		}
		if e.Victim == nil {
			victimName = "World"
			victimTeam = demoinfo.TeamUnassigned
		} else {
			victimName = e.Victim.Name
			victimTeam = e.Victim.Team
		}
		kill := common.Kill{
			KillerName: killerName,
			KillerTeam: killerTeam,
			VictimName: victimName,
			VictimTeam: victimTeam,
			Weapon:     e.Weapon.Type,
		}

		for i := 0; i < match.FrameRateRounded*killfeedLifetime; i++ {
			kills, ok := match.Killfeed[match.currentFrame+i]
			if ok {
				if len(kills) > 5 {
					match.Killfeed[match.currentFrame+i] = match.Killfeed[match.currentFrame+i][1:]
				}
				match.Killfeed[match.currentFrame+i] = append(kills, kill)
			} else {
				match.Killfeed[match.currentFrame+i] = []common.Kill{kill}
			}
		}
	})
	parser.RegisterEventHandler(func(e event.RoundStart) {
		match.currentPhase = common.PhaseFreezetime
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.RoundFreezetimeEnd) {
		match.currentPhase = common.PhaseRegular
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.BombPlanted) {
		match.currentPhase = common.PhasePlanted
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.RoundEnd) {
		match.currentPhase = common.PhaseRestart
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.GameHalfEnded) {
		match.currentPhase = common.PhaseHalftime
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(event.RoundStart) {
		for i := 1; i < int(match.SmokeEffectLifetime); i++ {
			match.Effects[match.currentFrame+i] = make([]common.Effect, 0)
		}
	})
}

// parse demo and save GameStates in slice
func parseGameStates(parser dem.Parser, match *Match) []common.OverviewState {
	playbackFrames := parser.Header().PlaybackFrames
	states := make([]common.OverviewState, 0, playbackFrames)

	for ok, err := parser.ParseNextFrame(); ok; ok, err = parser.ParseNextFrame() {
		if err != nil {
			log.Println(err)
			// return here or not?
			continue
		}
		if parser.CurrentFrame()%match.takeNthFrame != 0 {
			continue
		}
		match.currentFrame = len(states)

		gameState := parser.GameState()

		var isOnNormalElevation bool
		players := make([]common.Player, 0, 10)

		for _, p := range gameState.Participants().Playing() {
			var hasBomb bool
			inventory := make([]demoinfo.EquipmentType, 0)
			for _, w := range p.Weapons() {
				if w.Type == demoinfo.EqBomb {
					hasBomb = true
				}
				if isWeaponOrGrenade(w.Type) {
					if w.Type == demoinfo.EqFlash && w.AmmoReserve() > 0 {
						inventory = append(inventory, w.Type)
					}
					inventory = append(inventory, w.Type)
				}
			}
			sort.Slice(inventory, func(i, j int) bool { return inventory[i] < inventory[j] })
			if common.MapHasAlternateVersion(match.MapName) {
				if p.Position().Z > common.MapGetHeightThreshold(match.MapName) {
					isOnNormalElevation = true
				} else {
					isOnNormalElevation = false
				}
			}
			player := common.Player{
				Name: p.Name,
				ID:   p.UserID,
				Team: p.Team,
				Position: common.Point{
					X: float32(p.Position().X),
					Y: float32(p.Position().Y),
				},
				LastAlivePosition: common.Point{
					X: float32(p.LastAlivePosition.X),
					Y: float32(p.LastAlivePosition.Y),
				},
				ViewDirectionX:      p.ViewDirectionX(),
				FlashDuration:       p.FlashDurationTime(),
				FlashTimeRemaining:  p.FlashDurationTimeRemaining(),
				Inventory:           inventory,
				Health:              int16(p.Health()),
				Armor:               int16(p.Armor()),
				Money:               int16(p.Money()),
				Kills:               int16(p.Kills()),
				Deaths:              int16(p.Deaths()),
				Assists:             int16(p.Assists()),
				IsAlive:             p.IsAlive(),
				IsDefusing:          p.IsDefusing,
				IsOnNormalElevation: isOnNormalElevation,
				HasHelmet:           p.HasHelmet(),
				HasDefuseKit:        p.HasDefuseKit(),
				HasBomb:             hasBomb,
			}
			players = append(players, player)
		}
		sort.Slice(players, func(i, j int) bool { return players[i].ID < players[j].ID })

		grenades := make([]common.GrenadeProjectile, 0)

		for _, grenade := range gameState.GrenadeProjectiles() {
			if common.MapHasAlternateVersion(match.MapName) {
				if grenade.Position().Z > common.MapGetHeightThreshold(match.MapName) {
					isOnNormalElevation = true
				} else {
					isOnNormalElevation = false
				}
			}
			g := common.GrenadeProjectile{
				Position: common.Point{
					X: float32(grenade.Position().X),
					Y: float32(grenade.Position().Y),
				},
				Type:                grenade.WeaponInstance.Type,
				IsOnNormalElevation: isOnNormalElevation,
			}
			grenades = append(grenades, g)
		}

		infernos := make([]common.Inferno, 0)
		for _, inferno := range gameState.Infernos() {
			r2Points := inferno.Fires().Active().ConvexHull2D()
			commonPoints := make([]common.Point, 0)
			for _, point := range r2Points {
				commonPoint := common.Point{
					X: float32(point.X),
					Y: float32(point.Y),
				}
				commonPoints = append(commonPoints, commonPoint)
			}
			if common.MapHasAlternateVersion(match.MapName) {
				if inferno.Fires().Active().ConvexHull3D().Vertices[0].Z > common.MapGetHeightThreshold(match.MapName) {
					isOnNormalElevation = true
				} else {
					isOnNormalElevation = false
				}
			}
			i := common.Inferno{
				ConvexHull2D:        commonPoints,
				IsOnNormalElevation: isOnNormalElevation,
			}
			infernos = append(infernos, i)
		}

		var isBeingCarried bool
		if gameState.Bomb().Carrier != nil {
			isBeingCarried = true
		} else {
			isBeingCarried = false
		}
		if common.MapHasAlternateVersion(match.MapName) {
			if gameState.Bomb().Position().Z > common.MapGetHeightThreshold(match.MapName) {
				isOnNormalElevation = true
			} else {
				isOnNormalElevation = false
			}
		}
		bomb := common.Bomb{
			Position: common.Point{
				X: float32(gameState.Bomb().Position().X),
				Y: float32(gameState.Bomb().Position().Y),
			},
			IsBeingCarried:      isBeingCarried,
			IsOnNormalElevation: isOnNormalElevation,
		}

		cts := common.TeamState{
			ClanName: gameState.TeamCounterTerrorists().ClanName(),
			Score:    byte(gameState.TeamCounterTerrorists().Score()),
		}
		ts := common.TeamState{
			ClanName: gameState.TeamTerrorists().ClanName(),
			Score:    byte(gameState.TeamTerrorists().Score()),
		}

		var timer common.Timer

		if gameState.IsWarmupPeriod() {
			timer = common.Timer{
				TimeRemaining: 0,
				Phase:         common.PhaseWarmup,
			}
		} else {
			switch match.currentPhase {
			case common.PhaseFreezetime:
				freezetime, _ := strconv.Atoi(gameState.ConVars()["mp_freezetime"])
				remaining := time.Duration(freezetime)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = common.Timer{
					TimeRemaining: remaining,
					Phase:         common.PhaseFreezetime,
				}
			case common.PhaseRegular:
				roundtime, _ := strconv.ParseFloat(gameState.ConVars()["mp_roundtime_defuse"], 64)
				remaining := time.Duration(roundtime*60)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = common.Timer{
					TimeRemaining: remaining,
					Phase:         common.PhaseRegular,
				}
			case common.PhasePlanted:
				// mp_c4timer is not set in testdemo
				//bombtime, _ := strconv.Atoi(gameState.ConVars()["mp_c4timer"])
				bombtime := c4timer
				remaining := time.Duration(bombtime)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = common.Timer{
					TimeRemaining: remaining,
					Phase:         common.PhasePlanted,
				}
			case common.PhaseRestart:
				restartDelay, _ := strconv.Atoi(gameState.ConVars()["mp_round_restart_delay"])
				remaining := time.Duration(restartDelay)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = common.Timer{
					TimeRemaining: remaining,
					Phase:         common.PhaseRestart,
				}
			case common.PhaseHalftime:
				halftimeDuration, _ := strconv.Atoi(gameState.ConVars()["mp_halftime_duration"])
				remaining := time.Duration(halftimeDuration)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = common.Timer{
					TimeRemaining: remaining,
					Phase:         common.PhaseRestart,
				}
			}
		}

		state := common.OverviewState{
			IngameTick:            parser.GameState().IngameTick(),
			Players:               players,
			Grenades:              grenades,
			Infernos:              infernos,
			Bomb:                  bomb,
			TeamCounterTerrorists: cts,
			TeamTerrorists:        ts,
			Timer:                 timer,
		}

		states = append(states, state)
	}

	return states
}

func isWeaponOrGrenade(e demoinfo.EquipmentType) bool {
	return e.Class() == demoinfo.EqClassSMG ||
		e.Class() == demoinfo.EqClassHeavy ||
		e.Class() == demoinfo.EqClassRifle ||
		e.Class() == demoinfo.EqClassPistols ||
		e.Class() == demoinfo.EqClassGrenade

}

// Translate translates in-game world-relative coordinates to (0, 0) relative coordinates.
func (m Match) Translate(x, y float32) (float32, float32) {
	return x - m.MapPZero.X, m.MapPZero.Y - y
}

// TranslateScale translates and scales in-game world-relative coordinates to (0, 0) relative coordinates.
func (m Match) TranslateScale(x, y float32) (float32, float32) {
	x, y = m.Translate(x, y)
	return x / m.MapScale, y / m.MapScale
}
