package main

import (
	"fmt"
	"log"
	"os"
	// "time"

	dem "github.com/markus-wa/demoinfocs-golang"
	common "github.com/markus-wa/demoinfocs-golang/common"
	event "github.com/markus-wa/demoinfocs-golang/events"
	meta "github.com/markus-wa/demoinfocs-golang/metadata"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winHeight    int32 = 1024
	winWidth     int32 = 1024
	terrorR      uint8 = 252
	terrorG      uint8 = 176
	terrorB      uint8 = 12
	counterR     uint8 = 89
	counterG     uint8 = 206
	counterB     uint8 = 200
	radiusPlayer int32 = 10
)

var (
	mapName string
)

type OverviewState struct {
	IngameTick int
	Players    []common.Player
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Println(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("csgoverview", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winHeight, winWidth, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println(err)
		return
	}
	defer renderer.Destroy()

	// First pass to get round starts, half starts and header info
	demo, err := os.Open("test_cache.dem")
	if err != nil {
		log.Println(err)
		return
	}
	defer demo.Close()

	curFrame := 0

	// MatchStart + GameHalfEnd
	halfStarts := make([]int, 0)
	roundStarts := make([]int, 0)

	// find round starts and half starts
	parser := dem.NewParser(demo)
	h1 := parser.RegisterEventHandler(func(event.MatchStart) {
		halfStarts = append(halfStarts, parser.CurrentFrame())
	})
	h2 := parser.RegisterEventHandler(func(event.RoundStart) {
		roundStarts = append(roundStarts, parser.CurrentFrame())
	})
	h3 := parser.RegisterEventHandler(func(event.GameHalfEnded) {
		halfStarts = append(halfStarts, parser.CurrentFrame())
	})
	parser.RegisterEventHandler(func(event.AnnouncementWinPanelMatch) {
		parser.UnregisterEventHandler(h1)
		parser.UnregisterEventHandler(h2)
		parser.UnregisterEventHandler(h3)

	})
	// RoundEndOfficial / reason

	err = parser.ParseToEnd()
	if err != nil {
		log.Println(err)
		return
	}

	// frametime or frames per second?
	frameTime := parser.Header().FrameTime()
	mapName = parser.Header().MapName

	surface, err := img.Load(fmt.Sprintf("%v.jpg", mapName))
	if err != nil {
		log.Println(err)
		return
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Println(err)
		return
	}
	defer texture.Destroy()

	// err
	renderer.Clear()
	// nil, nil stretches texture to fill the screen
	// err
	renderer.Copy(texture, nil, nil)
	renderer.Present()

	_, err = demo.Seek(0, 0)
	if err != nil {
		log.Println(err)
		return
	}

	parser = dem.NewParser(demo)

	states := make([]OverviewState, 0)

	// parse demo and save GameStates in slice

	for ok, err := parser.ParseNextFrame(); ok; ok, err = parser.ParseNextFrame() {
		if err != nil {
			log.Println(err)
			// return here or not?
		}

		players := make([]common.Player, 0)

		for _, player := range parser.GameState().Participants().Playing() {
			players = append(players, *player)
		}

		state := OverviewState{
			IngameTick: parser.GameState().IngameTick(),
			Players:    players,
		}

		states = append(states, state)
	}
	fmt.Printf("Got %v frames\n", len(states))

	fmt.Println("Time per frame: %v", frameTime)
	fmt.Println("Round starts:")
	for i, tick := range roundStarts {
		fmt.Printf("Round %v:\t%v\n", i, tick)
	}
	fmt.Println("Half starts:")
	for i, tick := range halfStarts {
		fmt.Printf("Half %v:\t%v\n", i, tick)
	}

	paused := false

	// MAIN GAME LOOP

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch eventT := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_SPACE {
					paused = !paused
				}
			}
		}

		if paused {
			sdl.Delay(32)
			continue
		}

		renderer.Clear()
		renderer.Copy(texture, nil, nil)

		players := states[curFrame].Players

		for _, player := range players {
			DrawPlayer(renderer, &player)
		}

		// translate coordinates

		// draw the things

		fmt.Printf("Frame %v\n", curFrame)
		fmt.Printf("Ingame Tick %v\n", states[curFrame].IngameTick)
		renderer.Present()

		sdl.Delay(32)
		if curFrame < len(states)-1 {
			curFrame++
		}
	}

}

func DrawPlayer(renderer *sdl.Renderer, player *common.Player) {
	pos := player.LastAlivePosition

	scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX)
	var scaledYInt int32 = int32(scaledY)
	var colorR, colorG, colorB uint8

	if player.Team == common.TeamTerrorists {
		colorR = terrorR
		colorG = terrorG
		colorB = terrorB
	} else { // if player.Team == common.TeamCounterTerrorists {
		colorR = counterR
		colorG = counterG
		colorB = counterB
	}

	if player.Hp > 0 {
		gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer, colorR, colorG, colorB, 255)
		gfx.StringRGBA(renderer, scaledXInt+15, scaledYInt+15, player.Name, colorR, colorG, colorB, 255)

		viewAngle := -int32(player.ViewDirectionX) // negated because of sdl
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer+1, viewAngle-20, viewAngle+20, 200, 200, 200, 255)
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer+2, viewAngle-10, viewAngle+10, 200, 200, 200, 255)
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer+3, viewAngle-5, viewAngle+5, 200, 200, 200, 255)
	} else {
		//gfx.SetFont(fontdata, 10, 10)
		gfx.CharacterRGBA(renderer, scaledXInt, scaledYInt, 'X', colorR, colorG, colorB, 150)
	}

}
