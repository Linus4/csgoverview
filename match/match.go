// Package match contains a high-level parser for demos.
package match

import (
	"errors"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	ocom "github.com/linus4/csgoverview/common"
	dem "github.com/markus-wa/demoinfocs-golang"
	common "github.com/markus-wa/demoinfocs-golang/common"
	event "github.com/markus-wa/demoinfocs-golang/events"
)

const (
	flashEffectLifetime int = 10
	heEffectLifetime    int = 10
	killfeedLifetime    int = 10
	c4timer             int = 40
)

// Match contains general information about the demo and all relevant, parsed
// data from every tick of the demo that will be displayed.
type Match struct {
	MapName              string
	HalfStarts           []int
	RoundStarts          []int
	GrenadeEffects       map[int][]ocom.GrenadeEffect
	FrameRate            float64
	TickRate             float64
	FrameRateRounded     int
	States               []ocom.OverviewState
	SmokeEffectLifetime  int
	Killfeed             map[int][]ocom.Kill
	Shots                map[int][]ocom.Shot
	currentPhase         ocom.Phase
	latestTimerEventTime time.Duration
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
		HalfStarts:     make([]int, 0),
		RoundStarts:    make([]int, 0),
		GrenadeEffects: make(map[int][]ocom.GrenadeEffect),
		Killfeed:       make(map[int][]ocom.Kill),
		Shots:          make(map[int][]ocom.Shot),
	}

	match.FrameRate = header.FrameRate()
	if math.IsNaN(match.FrameRate) || match.FrameRate == 0 {
		if fallbackFrameRate == -1 {
			err := errors.New("could not parse Framerate from demo." +
				"Please provide a fallback value (command-line option -framerate)")
			return nil, err
		}
		match.FrameRate = fallbackFrameRate
	}
	match.TickRate = parser.TickRate()
	if math.IsNaN(match.TickRate) || match.TickRate == 0 {
		if fallbackTickRate == -1 {
			err := errors.New("could not parse Tickrate from demo." +
				"Please provide a fallback value (command-line option -tickrate)")
			return nil, err
		}
		match.TickRate = fallbackTickRate
	}
	match.FrameRateRounded = int(math.Round(match.FrameRate))
	match.MapName = header.MapName
	match.SmokeEffectLifetime = int(18 * match.FrameRate)

	registerEventHandlers(parser, match)
	match.States = parseGameStates(parser, match)

	return match, nil
}

func grenadeEventHandler(lifetime int, frame int, e event.GrenadeEvent, match *Match) {
	for i := 0; i < lifetime; i++ {
		effect := ocom.GrenadeEffect{
			GrenadeEvent: e,
			Lifetime:     i,
		}
		effects, ok := match.GrenadeEffects[frame+i]
		if ok {
			match.GrenadeEffects[frame+i] = append(effects, effect)
		} else {
			match.GrenadeEffects[frame+i] = []ocom.GrenadeEffect{effect}
		}
	}
}

func weaponFireEventHandler(frame int, e event.WeaponFire, match *Match) {
	if e.Shooter == nil {
		return
	}
	if e.Weapon.Class() == common.EqClassEquipment ||
		e.Weapon.Class() == common.EqClassGrenade ||
		e.Weapon.Class() == common.EqClassUnknown {
		return
	}
	isAwpShot := e.Weapon.Weapon == common.EqAWP
	shot := ocom.Shot{
		Position:       e.Shooter.Position,
		ViewDirectionX: e.Shooter.ViewDirectionX,
		IsAwpShot:      isAwpShot,
	}

	lifetime := int((match.FrameRate + 1) / 32)
	if lifetime == 0 {
		lifetime = 1
	}
	if isAwpShot {
		lifetime = int((match.FrameRate + 1) / 8)
	}
	for i := 0; i < lifetime; i++ {
		shots, ok := match.Shots[frame+i]
		if ok {
			match.Shots[frame+i] = append(shots, shot)
		} else {
			match.Shots[frame+i] = []ocom.Shot{shot}
		}
	}
}

func registerEventHandlers(parser *dem.Parser, match *Match) {
	parser.RegisterEventHandler(func(event.RoundStart) {
		match.RoundStarts = append(match.RoundStarts, parser.CurrentFrame())
	})
	parser.RegisterEventHandler(func(event.MatchStart) {
		match.HalfStarts = append(match.HalfStarts, parser.CurrentFrame())
	})
	parser.RegisterEventHandler(func(event.GameHalfEnded) {
		match.HalfStarts = append(match.HalfStarts, parser.CurrentFrame())
	})
	parser.RegisterEventHandler(func(e event.WeaponFire) {
		frame := parser.CurrentFrame()
		weaponFireEventHandler(frame, e, match)
	})
	parser.RegisterEventHandler(func(e event.FlashExplode) {
		frame := parser.CurrentFrame()
		grenadeEventHandler(flashEffectLifetime, frame, e.GrenadeEvent, match)
	})
	parser.RegisterEventHandler(func(e event.HeExplode) {
		frame := parser.CurrentFrame()
		grenadeEventHandler(heEffectLifetime, frame, e.GrenadeEvent, match)
	})
	parser.RegisterEventHandler(func(e event.SmokeStart) {
		frame := parser.CurrentFrame()
		grenadeEventHandler(match.SmokeEffectLifetime, frame, e.GrenadeEvent, match)
	})
	parser.RegisterEventHandler(func(e event.Kill) {
		frame := parser.CurrentFrame()
		var killerName, victimName string
		var killerTeam, victimTeam common.Team
		if e.Killer == nil {
			killerName = "World"
			killerTeam = common.TeamUnassigned
		} else {
			killerName = e.Killer.Name
			killerTeam = e.Killer.Team
		}
		if e.Victim == nil {
			victimName = "World"
			victimTeam = common.TeamUnassigned
		} else {
			victimName = e.Victim.Name
			victimTeam = e.Victim.Team
		}
		kill := ocom.Kill{
			KillerName: killerName,
			KillerTeam: killerTeam,
			VictimName: victimName,
			VictimTeam: victimTeam,
			Weapon:     e.Weapon.Weapon.String(),
		}

		for i := 0; i < match.FrameRateRounded*killfeedLifetime; i++ {
			kills, ok := match.Killfeed[frame+i]
			if ok {
				if len(kills) > 5 {
					match.Killfeed[frame+i] = match.Killfeed[frame+i][1:]
				}
				match.Killfeed[frame+i] = append(kills, kill)
			} else {
				match.Killfeed[frame+i] = []ocom.Kill{kill}
			}
		}
	})
	parser.RegisterEventHandler(func(e event.RoundStart) {
		match.currentPhase = ocom.PhaseFreezetime
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.RoundFreezetimeEnd) {
		match.currentPhase = ocom.PhaseRegular
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.BombPlanted) {
		match.currentPhase = ocom.PhasePlanted
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.RoundEnd) {
		match.currentPhase = ocom.PhaseRestart
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.GameHalfEnded) {
		match.currentPhase = ocom.PhaseHalftime
		match.latestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(event.AnnouncementWinPanelMatch) {
		match.HalfStarts = append(match.HalfStarts, parser.CurrentFrame())
	})
}

// parse demo and save GameStates in slice
func parseGameStates(parser *dem.Parser, match *Match) []ocom.OverviewState {
	playbackFrames := parser.Header().PlaybackFrames
	states := make([]ocom.OverviewState, 0, playbackFrames)

	for ok, err := parser.ParseNextFrame(); ok; ok, err = parser.ParseNextFrame() {
		if err != nil {
			log.Println(err)
			// return here or not?
			continue
		}

		gameState := parser.GameState()

		players := make([]common.Player, 0, 10)

		for _, p := range gameState.Participants().Playing() {
			equipment := make(map[int]*common.Equipment)
			for k := range p.RawWeapons {
				eq := *p.RawWeapons[k]
				equipment[k] = &eq
			}
			player := *p
			additionalPlayerInformation := *p.AdditionalPlayerInformation
			player.AdditionalPlayerInformation = &additionalPlayerInformation
			player.RawWeapons = equipment
			players = append(players, player)
		}

		grenades := make([]common.GrenadeProjectile, 0)

		for _, grenade := range gameState.GrenadeProjectiles() {
			grenades = append(grenades, *grenade)
		}

		infernos := make([]common.Inferno, 0)

		for _, inferno := range gameState.Infernos() {
			infernos = append(infernos, *inferno)
		}

		bomb := *gameState.Bomb()

		cts := *gameState.TeamCounterTerrorists()
		ts := *gameState.TeamTerrorists()

		var timer ocom.Timer

		if gameState.IsWarmupPeriod() {
			timer = ocom.Timer{
				TimeRemaining: 0,
				Phase:         ocom.PhaseWarmup,
			}
		} else {
			switch match.currentPhase {
			case ocom.PhaseFreezetime:
				freezetime, _ := strconv.Atoi(gameState.ConVars()["mp_freezetime"])
				remaining := time.Duration(freezetime)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhaseFreezetime,
				}
			case ocom.PhaseRegular:
				roundtime, _ := strconv.ParseFloat(gameState.ConVars()["mp_roundtime_defuse"], 64)
				remaining := time.Duration(roundtime*60)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhaseRegular,
				}
			case ocom.PhasePlanted:
				// mp_c4timer is not set in testdemo
				//bombtime, _ := strconv.Atoi(gameState.ConVars()["mp_c4timer"])
				bombtime := c4timer
				remaining := time.Duration(bombtime)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhasePlanted,
				}
			case ocom.PhaseRestart:
				restartDelay, _ := strconv.Atoi(gameState.ConVars()["mp_round_restart_delay"])
				remaining := time.Duration(restartDelay)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhaseRestart,
				}
			case ocom.PhaseHalftime:
				halftimeDuration, _ := strconv.Atoi(gameState.ConVars()["mp_halftime_duration"])
				remaining := time.Duration(halftimeDuration)*time.Second - (parser.CurrentTime() - match.latestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhaseRestart,
				}
			}
		}

		state := ocom.OverviewState{
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
