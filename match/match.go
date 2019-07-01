package match

import (
	"log"
	"math"
	"os"

	ocom "github.com/linus4/csgoverview/common"
	dem "github.com/markus-wa/demoinfocs-golang"
	common "github.com/markus-wa/demoinfocs-golang/common"
	event "github.com/markus-wa/demoinfocs-golang/events"
)

const (
	flashEffectLifetime int = 10
	heEffectLifetime    int = 10
)

type Match struct {
	MapName             string
	HalfStarts          []int
	RoundStarts         []int
	GrenadeEffects      map[int][]ocom.GrenadeEffect
	FrameRate           float64
	FrameRateRounded    int
	States              []ocom.OverviewState
	SmokeEffectLifetime int
}

func NewMatch(demoFileName string) (*Match, error) {
	demo, err := os.Open(demoFileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer demo.Close()

	parser := dem.NewParser(demo)
	header, err := parser.ParseHeader()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	match := &Match{
		HalfStarts:     make([]int, 0),
		RoundStarts:    make([]int, 0),
		GrenadeEffects: make(map[int][]ocom.GrenadeEffect),
		States:         make([]ocom.OverviewState, 0),
	}

	match.FrameRate = header.FrameRate()
	match.FrameRateRounded = int(math.Round(match.FrameRate))
	match.MapName = header.MapName
	match.SmokeEffectLifetime = int(18 * match.FrameRate)

	// First pass of the parser for events
	registerEventHandlers(parser, match)

	err = parser.ParseToEnd()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Second pass of the parser for gamestates
	_, err = demo.Seek(0, 0)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	parser = dem.NewParser(demo)
	match.States = parseGameStates(parser)

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
	parser.RegisterEventHandler(func(event.AnnouncementWinPanelMatch) {
		parser.UnregisterEventHandler(h1)
		parser.UnregisterEventHandler(h2)
		parser.UnregisterEventHandler(h3)
		parser.UnregisterEventHandler(h4)
	})
}

// parse demo and save GameStates in slice
func parseGameStates(parser *dem.Parser) []ocom.OverviewState {
	states := make([]ocom.OverviewState, 0)

	for ok, err := parser.ParseNextFrame(); ok; ok, err = parser.ParseNextFrame() {
		if err != nil {
			log.Println(err)
			// return here or not?
			continue
		}

		gameState := parser.GameState()

		players := make([]ocom.OverviewPlayer, 0)

		for _, p := range gameState.Participants().Playing() {
			// common.RawWeapons is map[int]*common.Equipment, but it is unclear what the key means
			equipment := make([]common.Equipment, 0)
			for _, eq := range p.Weapons() {
				equipment = append(equipment, *eq)
			}
			player := ocom.OverviewPlayer{
				Player:  *p,
				Weapons: equipment,
			}
			// I will save all information like this.......... FIXME
			additionalPlayerInformation := *p.AdditionalPlayerInformation
			player.AdditionalPlayerInformation = &additionalPlayerInformation
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

		state := ocom.OverviewState{
			IngameTick:            parser.GameState().IngameTick(),
			Players:               players,
			Grenades:              grenades,
			Infernos:              infernos,
			Bomb:                  bomb,
			TeamCounterTerrorists: cts,
			TeamTerrorists:        ts,
		}

		states = append(states, state)
	}

	return states
}
