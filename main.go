package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/linus4/csgoverview/draw"
	"github.com/linus4/csgoverview/match"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	winWidth             int32 = 1624
	winHeight            int32 = 1024
	mapOverviewWidth     int32 = 1024
	mapOverviewHeight    int32 = 1024
	mapXOffset           int32 = 300
	mapYOffset           int32 = 0
	infobarElementHeight int32 = 100
	nameMapFontSize      int   = 14
)

var (
	paused   bool
	curFrame int
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./csgoverview [demoname]")
		return
	}
	demoFileName := os.Args[1]

	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		log.Println("trying to initialize SDL:", err)
		return
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		log.Println("trying to initialize the TTF lib:", err)
		return
	}
	defer ttf.Quit()

	font, err := ttf.OpenFont("liberationserif-regular.ttf", nameMapFontSize)
	if err != nil {
		log.Println("trying to open the font:", err)
		return
	}
	defer font.Close()

	window, err := sdl.CreateWindow("csgoverview", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		log.Println("trying to create SDL window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println("trying to create SDL renderer:", err)
		return
	}
	defer renderer.Destroy()
	renderer.SetLogicalSize(mapOverviewWidth+2*mapXOffset, mapOverviewHeight+mapYOffset)

	match, err := match.NewMatch(demoFileName)
	if err != nil {
		log.Println(err)
		return
	}

	mapSurface, err := img.Load(fmt.Sprintf("%v.jpg", match.MapName))
	if err != nil {
		log.Println("trying to load map overview image:", err)
		return
	}
	defer mapSurface.Free()

	mapTexture, err := renderer.CreateTextureFromSurface(mapSurface)
	if err != nil {
		log.Println("trying to create mapTexture from Surface", err)
		return
	}
	defer mapTexture.Destroy()

	mapRect := &sdl.Rect{mapXOffset, mapYOffset, mapOverviewWidth, mapOverviewHeight}

	// MAIN GAME LOOP
	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch eventT := event.(type) {
			case *sdl.QuitEvent:
				return

			case *sdl.KeyboardEvent:
				handleKeyboardEvents(eventT, window, match)
			}
		}

		if paused {
			sdl.Delay(32)
			// draw?
			continue
		}

		renderer.SetDrawColor(10, 10, 10, 255)
		renderer.Clear()

		draw.DrawInfobars(renderer, match, font, curFrame)
		renderer.Copy(mapTexture, nil, mapRect)

		infernos := match.States[curFrame].Infernos
		for _, inferno := range infernos {
			draw.DrawInferno(renderer, &inferno, match)
		}

		effects := match.GrenadeEffects[curFrame]
		for _, effect := range effects {
			draw.DrawGrenadeEffect(renderer, &effect, match)
		}

		grenades := match.States[curFrame].Grenades
		for _, grenade := range grenades {
			draw.DrawGrenade(renderer, &grenade, match)
		}

		bomb := match.States[curFrame].Bomb
		draw.DrawBomb(renderer, &bomb, match)

		players := match.States[curFrame].Players
		for _, player := range players {
			draw.DrawPlayer(renderer, &player, font, match)
		}

		renderer.Present()

		updateWindowTitle(window, match)

		var playbackSpeed float64 = 1

		// frameDuration is in ms
		frameDuration := float64(time.Since(frameStart) / 1000000)
		keyboardState := sdl.GetKeyboardState()
		if keyboardState[sdl.GetScancodeFromKey(sdl.K_w)] != 0 {
			playbackSpeed = 5
		}
		if keyboardState[sdl.GetScancodeFromKey(sdl.K_s)] != 0 {
			playbackSpeed = 0.5
		}
		delay := (1/playbackSpeed)*(1000/match.FrameRate) - frameDuration
		if delay < 0 {
			delay = 0
		}
		sdl.Delay(uint32(delay))
		if curFrame < len(match.States)-1 {
			curFrame++
		}
	}

}

func handleKeyboardEvents(eventT *sdl.KeyboardEvent, window *sdl.Window, match *match.Match) {
	if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_SPACE {
		paused = !paused
	}

	if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_a {
		if eventT.Keysym.Mod == sdl.KMOD_LSHIFT || eventT.Keysym.Mod == sdl.KMOD_RSHIFT {
			if curFrame < match.FrameRateRounded*30 {
				curFrame = 0
			} else {
				curFrame -= match.FrameRateRounded * 30
			}
		} else {
			if curFrame < match.FrameRateRounded*10 {
				curFrame = 0
			} else {
				curFrame -= match.FrameRateRounded * 10
			}
		}
	}

	if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_d {
		if eventT.Keysym.Mod == sdl.KMOD_LSHIFT || eventT.Keysym.Mod == sdl.KMOD_RSHIFT {
			if curFrame+match.FrameRateRounded*30 > len(match.States)-1 {
				curFrame = len(match.States) - 1
			} else {
				curFrame += match.FrameRateRounded * 30
			}
		} else {
			if curFrame+match.FrameRateRounded*10 > len(match.States)-1 {
				curFrame = len(match.States) - 1
			} else {
				curFrame += match.FrameRateRounded * 10
			}
		}
	}

	if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_q {
		if eventT.Keysym.Mod == sdl.KMOD_LSHIFT || eventT.Keysym.Mod == sdl.KMOD_RSHIFT {
			set := false
			for i, frame := range match.HalfStarts {
				if curFrame < frame {
					if i > 1 && curFrame < match.HalfStarts[i-1]+match.FrameRateRounded/2 {
						curFrame = match.HalfStarts[i-2]
						set = true
						break
					}
					curFrame = match.HalfStarts[i-1]
					set = true
					break
				}
			}
			// not set -> last round of match
			if !set {
				if len(match.HalfStarts) > 1 && curFrame < match.HalfStarts[len(match.HalfStarts)-1]+match.FrameRateRounded/2 {
					curFrame = match.HalfStarts[len(match.HalfStarts)-2]
				} else {
					curFrame = match.HalfStarts[len(match.HalfStarts)-1]
				}
			}
		} else {
			set := false
			for i, frame := range match.RoundStarts {
				if curFrame < frame {
					if i > 1 && curFrame < match.RoundStarts[i-1]+match.FrameRateRounded/2 {
						curFrame = match.RoundStarts[i-2]
						set = true
						break
					}
					curFrame = match.RoundStarts[i-1]
					set = true
					break
				}
			}
			// not set -> last round of match
			if !set {
				if len(match.RoundStarts) > 1 && curFrame < match.RoundStarts[len(match.RoundStarts)-1]+match.FrameRateRounded/2 {
					curFrame = match.RoundStarts[len(match.RoundStarts)-2]
				} else {
					curFrame = match.RoundStarts[len(match.RoundStarts)-1]
				}
			}
		}
	}

	if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_e {
		if eventT.Keysym.Mod == sdl.KMOD_LSHIFT || eventT.Keysym.Mod == sdl.KMOD_RSHIFT {
			for _, frame := range match.HalfStarts {
				if curFrame < frame {
					curFrame = frame
					break
				}
			}
		} else {
			for _, frame := range match.RoundStarts {
				if curFrame < frame {
					curFrame = frame
					break
				}
			}
		}
	}

	/*
		if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_p {
			fmt.Println("take screenshot")
			fileName := fmt.Sprintf("screenshot_"+demoFileName+"_%v", curFrame)
			// using a renderer so window does not have a surface
			screenshotSurface, err := window.GetSurface()
			if err != nil {
				log.Println(err)
				return
			}
			err = img.SavePNG(screenshotSurface, fileName)
			if err != nil {
				log.Println(err)
				return
			}
		}
	*/
}

func updateWindowTitle(window *sdl.Window, match *match.Match) {
	cts := match.States[curFrame].TeamCounterTerrorists
	ts := match.States[curFrame].TeamTerrorists
	clanNameCTs := cts.ClanName
	if clanNameCTs == "" {
		clanNameCTs = "Counter Terrorists"
	}
	clanNameTs := ts.ClanName
	if clanNameTs == "" {
		clanNameTs = "Terrorists"
	}
	windowTitle := fmt.Sprintf("%s  [%d:%d]  %s", clanNameCTs, cts.Score, ts.Score, clanNameTs)
	// expensive?
	window.SetTitle(windowTitle)
}
