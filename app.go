package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/linus4/csgoverview/match"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	winWidth             int32 = 1624
	winHeight            int32 = 1024
	nameMapFontSize      int   = 14
	mapOverviewWidth     int32 = 1024
	mapOverviewHeight    int32 = 1024
	mapXOffset           int32 = 300
	mapYOffset           int32 = 0
	infobarElementHeight int32 = 100
)

var (
	paused   bool
	curFrame int
)

type Config struct {
	// Path to font file (.ttf)
	FontPath string

	// Path to overview directory
	OverviewDir string

	// Fallback GOTV Framerate
	FrameRate float64

	// Fallback Gameserver Tickrate
	TickRate float64
}

var DefaultConfig = Config{
	FrameRate: -1,
	TickRate:  -1,
}

func run(c *Config) error {
	var demoFileName string
	if len(flag.Args()) < 1 {
		demoFileNameB, err := exec.Command("zenity", "--file-selection").Output()
		if err != nil {
			fmt.Println("Usage: ./csgoverview [path to demo]")
			return err
		}
		demoFileName = string(demoFileNameB)[:len(demoFileNameB)-1]
	} else {
		demoFileName = flag.Args()[0]
	}

	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		log.Println("trying to initialize SDL:", err)
		return err
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		log.Println("trying to initialize the TTF lib:", err)
		return err
	}
	defer ttf.Quit()

	font, err := ttf.OpenFont(c.FontPath, nameMapFontSize)
	if err != nil {
		log.Println("trying to open the font file (system or flag):", err)
		font, err = ttf.OpenFont("DejaVuSans.ttf", nameMapFontSize)
		if err != nil {
			log.Println("trying to open the font file in the curent directory:", err)
			return err
		}
	}
	defer font.Close()
	font.SetStyle(ttf.STYLE_BOLD)

	window, err := sdl.CreateWindow("csgoverview", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		log.Println("trying to create SDL window:", err)
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Println("trying to create SDL renderer:", err)
		return err
	}
	defer renderer.Destroy()
	renderer.SetLogicalSize(mapOverviewWidth+2*mapXOffset, mapOverviewHeight+mapYOffset)

	match, err := match.NewMatch(demoFileName, c.FrameRate, c.TickRate)
	if err != nil {
		log.Println("trying to parse demo file:", err)
		return err
	}

	mapSurface, err := img.Load(filepath.Join(c.OverviewDir, fmt.Sprintf("%v.jpg", match.MapName)))
	if err != nil {
		log.Println("trying to load map overview image from %v:", c.OverviewDir, err)
		mapSurface, err = img.Load(fmt.Sprintf("%v.jpg", match.MapName))
		if err != nil {
			log.Println("trying to load map overview image from current directory:", err)
			return err
		}
	}
	defer mapSurface.Free()

	mapTexture, err := renderer.CreateTextureFromSurface(mapSurface)
	if err != nil {
		log.Println("trying to create mapTexture from Surface", err)
		return err
	}
	defer mapTexture.Destroy()

	mapRect := &sdl.Rect{mapXOffset, mapYOffset, mapOverviewWidth, mapOverviewHeight}

	// MAIN GAME LOOP
	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch eventT := event.(type) {
			case *sdl.QuitEvent:
				return err

			case *sdl.KeyboardEvent:
				handleKeyboardEvents(eventT, window, match)
			}
		}

		if paused {
			sdl.Delay(32)
			updateGraphics(renderer, match, font, mapTexture, mapRect)
			updateWindowTitle(window, match)
			continue
		}

		updateGraphics(renderer, match, font, mapTexture, mapRect)
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
		if isShiftPressed(eventT) {
			if curFrame < match.FrameRateRounded*10 {
				curFrame = 0
			} else {
				curFrame -= match.FrameRateRounded * 10
			}
		} else {
			if curFrame < match.FrameRateRounded*5 {
				curFrame = 0
			} else {
				curFrame -= match.FrameRateRounded * 5
			}
		}
	}

	if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_d {
		if isShiftPressed(eventT) {
			if curFrame+match.FrameRateRounded*10 > len(match.States)-1 {
				curFrame = len(match.States) - 1
			} else {
				curFrame += match.FrameRateRounded * 10
			}
		} else {
			if curFrame+match.FrameRateRounded*5 > len(match.States)-1 {
				curFrame = len(match.States) - 1
			} else {
				curFrame += match.FrameRateRounded * 5
			}
		}
	}

	if eventT.Type == sdl.KEYDOWN && eventT.Keysym.Sym == sdl.K_q {
		if isShiftPressed(eventT) {
			set := false
			for i, frame := range match.HalfStarts {
				if curFrame < frame {
					if i > 1 && curFrame < match.HalfStarts[i-1]+match.FrameRateRounded/2 {
						curFrame = match.HalfStarts[i-2]
						set = true
						break
					}
					if i-1 < 0 {
						curFrame = 0
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
					if i-1 < 0 {
						curFrame = 0
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
		if isShiftPressed(eventT) {
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
				return err
			}
			err = img.SavePNG(screenshotSurface, fileName)
			if err != nil {
				log.Println(err)
				return err
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

func updateGraphics(renderer *sdl.Renderer, match *match.Match, font *ttf.Font, mapTexture *sdl.Texture, mapRect *sdl.Rect) {
	renderer.SetDrawColor(10, 10, 10, 255)
	renderer.Clear()

	drawInfobars(renderer, match, font)
	renderer.Copy(mapTexture, nil, mapRect)

	infernos := match.States[curFrame].Infernos
	for _, inferno := range infernos {
		drawInferno(renderer, &inferno, match)
	}

	effects := match.GrenadeEffects[curFrame]
	for _, effect := range effects {
		drawGrenadeEffect(renderer, &effect, match)
	}

	grenades := match.States[curFrame].Grenades
	for _, grenade := range grenades {
		drawGrenade(renderer, &grenade, match)
	}

	bomb := match.States[curFrame].Bomb
	drawBomb(renderer, &bomb, match)

	players := match.States[curFrame].Players
	for _, player := range players {
		drawPlayer(renderer, &player, font, match)
	}

	renderer.Present()
}

func isShiftPressed(event *sdl.KeyboardEvent) bool {
	pressed := event.Keysym.Mod & sdl.KMOD_SHIFT

	if pressed > 0 {
		return true
	} else {
		return false
	}
}
