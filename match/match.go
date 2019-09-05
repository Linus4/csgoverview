// Package match contains a high-level parser for demos.
package match

import (
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
	FrameRateRounded     int
	States               []ocom.OverviewState
	SmokeEffectLifetime  int
	Killfeed             map[int][]ocom.Kill
	currentPhase         ocom.Phase
	LatestTimerEventTime time.Duration
}

// NewMatch parses the demo at the specified path in the argument and returns a
// match.Match containing all relevant data from the demo.
func NewMatch(demoFileName string) (*Match, error) {
	demo, err := os.Open(demoFileName)
	if err != nil {
		return nil, err
	}
	defer demo.Close()

	parser := dem.NewParser(demo)
	header, err := parser.ParseHeader()
	if err != nil {
		return nil, err
	}

	match := &Match{
		HalfStarts:     make([]int, 0),
		RoundStarts:    make([]int, 0),
		GrenadeEffects: make(map[int][]ocom.GrenadeEffect),
		Killfeed:       make(map[int][]ocom.Kill),
	}

	match.FrameRate = header.FrameRate()
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

func registerEventHandlers(parser *dem.Parser, match *Match) {
	h1 := parser.RegisterEventHandler(func(event.RoundStart) {
		match.RoundStarts = append(match.RoundStarts, parser.CurrentFrame())
	})
	h2 := parser.RegisterEventHandler(func(event.MatchStart) {
		match.HalfStarts = append(match.HalfStarts, parser.CurrentFrame())
	})
	h3 := parser.RegisterEventHandler(func(event.GameHalfEnded) {
		match.HalfStarts = append(match.HalfStarts, parser.CurrentFrame())
	})
	h4 := parser.RegisterEventHandler(func(event.TeamSideSwitch) {
		match.HalfStarts = append(match.HalfStarts, parser.CurrentFrame())
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
		kill := ocom.Kill{
			KillerName: e.Killer.Name,
			KillerTeam: e.Killer.Team,
			VictimName: e.Victim.Name,
			VictimTeam: e.Victim.Team,
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
		match.LatestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.RoundFreezetimeEnd) {
		match.currentPhase = ocom.PhaseRegular
		match.LatestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.BombPlanted) {
		match.currentPhase = ocom.PhasePlanted
		match.LatestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.RoundEnd) {
		match.currentPhase = ocom.PhaseRestart
		match.LatestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(e event.GameHalfEnded) {
		match.currentPhase = ocom.PhaseHalftime
		match.LatestTimerEventTime = parser.CurrentTime()
	})
	parser.RegisterEventHandler(func(event.AnnouncementWinPanelMatch) {
		parser.UnregisterEventHandler(h1)
		parser.UnregisterEventHandler(h2)
		parser.UnregisterEventHandler(h3)
		parser.UnregisterEventHandler(h4)
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
				remaining := time.Duration(freezetime)*time.Second - (parser.CurrentTime() - match.LatestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhaseFreezetime,
				}
			case ocom.PhaseRegular:
				roundtime, _ := strconv.ParseFloat(gameState.ConVars()["mp_roundtime_defuse"], 64)
				remaining := time.Duration(roundtime*60)*time.Second - (parser.CurrentTime() - match.LatestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhaseRegular,
				}
			case ocom.PhasePlanted:
				// mp_c4timer is not set in testdemo
				//bombtime, _ := strconv.Atoi(gameState.ConVars()["mp_c4timer"])
				bombtime := c4timer
				remaining := time.Duration(bombtime)*time.Second - (parser.CurrentTime() - match.LatestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhasePlanted,
				}
			case ocom.PhaseRestart:
				restartDelay, _ := strconv.Atoi(gameState.ConVars()["mp_round_restart_delay"])
				remaining := time.Duration(restartDelay)*time.Second - (parser.CurrentTime() - match.LatestTimerEventTime)
				timer = ocom.Timer{
					TimeRemaining: remaining,
					Phase:         ocom.PhaseRestart,
				}
			case ocom.PhaseHalftime:
				halftimeDuration, _ := strconv.Atoi(gameState.ConVars()["mp_halftime_duration"])
				remaining := time.Duration(halftimeDuration)*time.Second - (parser.CurrentTime() - match.LatestTimerEventTime)
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
