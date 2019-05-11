package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

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
	mapName          string
	halfStarts       []int
	roundStarts      []int
	grenadeEffects   map[int][]GrenadeEffect
	curFrame         int = 0
	frameRate        float64
	frameRateRounded int
)

type OverviewState struct {
	IngameTick int
	Players    []common.Player
	Grenades   []common.GrenadeProjectile
}

type GrenadeEffect struct {
	GrenadeEvent event.GrenadeEvent
	Lifetime     int32
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./csgoverview [demoname]")
		return
	}

	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
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

	demo, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}
	defer demo.Close()

	// MatchStart + GameHalfEnd
	halfStarts = make([]int, 0)
	roundStarts = make([]int, 0)
	roundStarts = append(roundStarts, 0)
	grenadeEffects = make(map[int][]GrenadeEffect)

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
	parser.RegisterEventHandler(func(e event.FlashExplode) {
		lifetime := 10
		frame := parser.CurrentFrame()
		for i := 0; i < lifetime; i++ {
			effect := GrenadeEffect{
				GrenadeEvent: e.GrenadeEvent,
				Lifetime:     int32(i),
			}
			effects, ok := grenadeEffects[frame+i]
			if ok {
				grenadeEffects[frame+i] = append(effects, effect)
			} else {
				grenadeEffects[frame+i] = []GrenadeEffect{effect}
			}
		}
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

	frameRate = parser.Header().FrameRate()
	frameRateRounded = int(math.Round(frameRate))
	mapName = parser.Header().MapName

	surface, err := img.Load(fmt.Sprintf("%v.jpg", mapName))
	if err != nil {
		log.Println(err)
		return
	}

	mapTexture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Println(err)
		return
	}
	defer mapTexture.Destroy()

	// err
	renderer.Clear()
	// nil, nil stretches texture to fill the screen
	// err
	renderer.Copy(mapTexture, nil, nil)
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

		grenades := make([]common.GrenadeProjectile, 0)

		for _, grenade := range parser.GameState().GrenadeProjectiles() {
			grenades = append(grenades, *grenade)
		}

		state := OverviewState{
			IngameTick: parser.GameState().IngameTick(),
			Players:    players,
			Grenades:   grenades,
		}

		states = append(states, state)
	}
	fmt.Printf("Got %v frames\n", len(states))

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
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch eventT := event.(type) {
			case *sdl.QuitEvent:
				return

			case *sdl.KeyboardEvent:
				if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_SPACE {
					paused = !paused
				}

				if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_a {
					if eventT.Keysym.Mod == sdl.KMOD_LSHIFT || eventT.Keysym.Mod == sdl.KMOD_RSHIFT {
						if curFrame < frameRateRounded*30 {
							curFrame = 0
						} else {
							curFrame -= frameRateRounded * 30
						}
					} else {
						if curFrame < frameRateRounded*10 {
							curFrame = 0
						} else {
							curFrame -= frameRateRounded * 10
						}
					}
				}

				if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_d {
					if eventT.Keysym.Mod == sdl.KMOD_LSHIFT || eventT.Keysym.Mod == sdl.KMOD_RSHIFT {
						if curFrame+frameRateRounded*30 > len(states)-1 {
							curFrame = len(states) - 1
						} else {
							curFrame += frameRateRounded * 30
						}
					} else {
						if curFrame+frameRateRounded*10 > len(states)-1 {
							curFrame = len(states) - 1
						} else {
							curFrame += frameRateRounded * 10
						}
					}
				}

				if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_q {
					set := false
					for i, frame := range roundStarts {
						if curFrame < frame {
							if i > 1 && curFrame < roundStarts[i-1]+frameRateRounded/2 {
								curFrame = roundStarts[i-2]
								set = true
								break
							}
							curFrame = roundStarts[i-1]
							set = true
							break
						}
					}
					// not set -> last round of match
					if !set {
						if len(roundStarts) > 1 && curFrame < roundStarts[len(roundStarts)-1]+frameRateRounded/2 {
							curFrame = roundStarts[len(roundStarts)-2]
						} else {
							curFrame = roundStarts[len(roundStarts)-1]
						}
					}
				}

				if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_e {
					for _, frame := range roundStarts {
						if curFrame < frame {
							curFrame = frame
							break
						}
					}
				}
			}
		}

		if paused {
			sdl.Delay(32)
			// draw?
			continue
		}

		renderer.Clear()
		renderer.Copy(mapTexture, nil, nil)

		players := states[curFrame].Players

		for _, player := range players {
			DrawPlayer(renderer, &player)
		}

		effects, ok := grenadeEffects[curFrame]

		if ok {
			for _, effect := range effects {
				DrawGrenadeEffect(renderer, &effect)
			}
		}

		grenades := states[curFrame].Grenades

		for _, grenade := range grenades {
			DrawGrenade(renderer, &grenade)
		}

		fmt.Printf("Ingame Tick %v\n", states[curFrame].IngameTick)
		renderer.Present()

		// frameDuration is in ms
		frameDuration := float64(time.Since(frameStart) / 1000000)
		sdl.Delay(uint32(frameRate - frameDuration))
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

		// FlashDuration is not the time remaining but always the total amount of time flashed from a single flashbang
		if player.FlashDuration > 0.8 {
			gfx.FilledCircleRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer-5, 200, 200, 200, 200)
		}

		if player.IsDefusing {
			gfx.CharacterRGBA(renderer, scaledXInt-radiusPlayer/4, scaledYInt-radiusPlayer/4, 'D', counterR, counterG, counterB, 200)
		}
	} else {
		//gfx.SetFont(fontdata, 10, 10)
		gfx.CharacterRGBA(renderer, scaledXInt, scaledYInt, 'X', colorR, colorG, colorB, 150)
	}
}

func DrawGrenade(renderer *sdl.Renderer, grenade *common.GrenadeProjectile) {
	pos := grenade.Position

	scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX)
	var scaledYInt int32 = int32(scaledY)
	var colorR, colorG, colorB uint8

	switch grenade.Weapon {
	case common.EqDecoy:
		colorR = 102
		colorG = 34
		colorB = 0
	case common.EqMolotov:
		colorR = 255
		colorG = 153
		colorB = 0
	case common.EqIncendiary:
		colorR = 255
		colorG = 153
		colorB = 0
	case common.EqFlash:
		colorR = 128
		colorG = 170
		colorB = 255
	case common.EqSmoke:
		colorR = 153
		colorG = 153
		colorB = 153
	case common.EqHE:
		colorR = 85
		colorG = 150
		colorB = 0
	}

	gfx.BoxRGBA(renderer, scaledXInt-2, scaledYInt-3, scaledXInt+2, scaledYInt+3, colorR, colorG, colorB, 255)

	// GenadeEvent SmokeExpired SmokeStart InfernoStart InfernoExpired HeExplode GrenadeProjectileDestroy FlashExplode

	// trajectory
	// not drawing trajectories at the moment

}

func DrawGrenadeEffect(renderer *sdl.Renderer, effect *GrenadeEffect) {
	pos := effect.GrenadeEvent.Position

	scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX)
	var scaledYInt int32 = int32(scaledY)
	var colorR, colorG, colorB uint8

	switch effect.GrenadeEvent.GrenadeType {
	case common.EqDecoy:
		colorR = 102
		colorG = 34
		colorB = 0
	case common.EqMolotov:
		colorR = 255
		colorG = 153
		colorB = 0
	case common.EqIncendiary:
		colorR = 255
		colorG = 153
		colorB = 0
	case common.EqFlash:
		colorR = 128
		colorG = 170
		colorB = 255
	case common.EqSmoke:
		colorR = 153
		colorG = 153
		colorB = 153
	case common.EqHE:
		colorR = 85
		colorG = 150
		colorB = 0
	}

	gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, effect.Lifetime, colorR, colorG, colorB, 255)
}
