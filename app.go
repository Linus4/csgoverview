package main

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cheggaaa/pb/v3"
	"github.com/linus4/csgoverview/common"
	"github.com/linus4/csgoverview/match"
	demoinfo "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	winWidth             int32  = 1624
	winHeight            int32  = 1024
	nameMapFontSize      int    = 14
	mapOverviewWidth     int32  = 1024
	mapOverviewHeight    int32  = 1024
	mapXOffset           int32  = 300
	mapYOffset           int32  = 0
	infobarElementHeight int32  = 100
	appVersion           string = "v1.2.0"
	hotkeysString        string = `
* a -> 3 s backwards
* d -> 3 s forwards
* A -> 10 s backwards
* D -> 10 s forwards
* w -> increase playback speed
* s -> decrease playback speed
* r -> reset playback speed to x 1
* W -> hold to speed up 5 x
* S -> hold to slow down to 0.5 x
* q -> round backwards
* e -> round forwards
* Q -> to start of previous half
* E -> to start of next half
* space -> toggle pause
* mouse wheel -> scroll 1 second forwards/backwards
* c -> switch to alternate overview image (normal / lower)
* h -> hide player names on map
* 0-9 -> copy players position and view angle to clipboard
* k -> show hotkeys
`
)

// Config contains information the application requires in order to run
type Config struct {
	// Path to font file (.ttf)
	FontPath string

	// Path to overview directory
	OverviewDir string

	// Whether to just print the version number
	PrintVersion bool
}

// DefaultConfig contains standard parameters for the application.
var DefaultConfig = Config{}

// App contains the state of the application.
type app struct {
	window              *sdl.Window
	renderer            *sdl.Renderer
	mapTexture          *sdl.Texture
	alternateMapTexture *sdl.Texture
	mapRect             *sdl.Rect
	font                *ttf.Font
	match               *match.Match
	config              *Config

	lastDrawnAt                 time.Time
	isPaused                    bool
	curFrame                    int
	isOnNormalElevation         bool
	playbackSpeedModifier       float64
	staticPlaybackSpeedModifier float64
	hidePlayerNames             bool
}

func run(c *Config) error {
	// Progress Bar template
	pbTmpl := `{{string . "step" | green}} {{ bar . "[" "#" (cycle . "⬒" "⬔" "⬓" "⬕" ) "." "]"}} {{percent .}}`
	pbInstance := pb.ProgressBarTemplate(pbTmpl).Start(8)
	pbInstance.Set("step", "Starting...")
	pbInstance.Start()

	var demoFileName string
	if c.PrintVersion {
		fmt.Println(appVersion)
		return nil
	}
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

	pbInstance.Set("step", "Initializing SDL")
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		errorString := fmt.Sprintf("trying to initialize SDL:\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, nil)
		return err
	}
	defer sdl.Quit()
	pbInstance.Increment()

	pbInstance.Set("step", "Initializing TTF lib")
	err = ttf.Init()
	if err != nil {
		errorString := fmt.Sprintf("trying to initialize the TTF lib:\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, nil)
		return err
	}
	defer ttf.Quit()
	pbInstance.Increment()

	pbInstance.Set("step", "Opening font file")
	font, err := ttf.OpenFont(c.FontPath, nameMapFontSize)
	if err != nil {
		errorString := fmt.Sprintf("trying to open font file (system):\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, nil)
		return err
	}
	defer font.Close()
	pbInstance.Increment()

	pbInstance.Set("step", "Creating SDL window")
	window, err := sdl.CreateWindow("csgoverview", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		errorString := fmt.Sprintf("trying to create SDL window:\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, nil)
		return err
	}
	defer window.Destroy()
	pbInstance.Increment()

	pbInstance.Set("step", "Creating SDL renderer")
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		errorString := fmt.Sprintf("trying to create SDL renderer:\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, window)
		return err
	}
	defer renderer.Destroy()
	renderer.SetLogicalSize(mapOverviewWidth+2*mapXOffset, mapOverviewHeight+mapYOffset)
	pbInstance.Increment()

	pbInstance.Set("step", "Parsing demo file")
	match, err := match.NewMatch(demoFileName, pbInstance)
	if err != nil {
		errorString := fmt.Sprintf("trying to parse demo file:\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, window)
		return err
	}
	pbInstance.Increment()

	pbInstance.Set("step", "Loading map overview surface")
	mapSurface, err := img.Load(filepath.Join(c.OverviewDir, fmt.Sprintf("%v.jpg", match.MapName)))
	if err != nil {
		errorString := fmt.Sprintf("trying to load map overview image from %v: \n"+
			"%v \nFollow the instructions on https://github.com/linus4/csgoverview "+
			"to place the overview images in this directory.", c.OverviewDir, err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, window)
		return err
	}
	defer mapSurface.Free()
	pbInstance.Increment()

	pbInstance.Set("step", "Creating texture from surface")
	mapTexture, err := renderer.CreateTextureFromSurface(mapSurface)
	if err != nil {
		errorString := fmt.Sprintf("trying to create mapTexture from Surface:\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, window)
		return err
	}
	defer mapTexture.Destroy()

	var alternateMapTexture *sdl.Texture
	if common.MapHasAlternateVersion(match.MapName) {
		alternateFileName := common.MapGetAlternateVersion(match.MapName)
		alternateMapSurface, err := img.Load(filepath.Join(c.OverviewDir, fmt.Sprintf("%v", alternateFileName)))
		if err != nil {
			errorString := fmt.Sprintf("trying to load map overview image from %v: \n"+
				"%v \nFollow the instructions on https://github.com/linus4/csgoverview "+
				"to place the overview images in this directory.", c.OverviewDir, err)
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, window)
			return err
		}
		defer alternateMapSurface.Free()

		alternateMapTexture, err = renderer.CreateTextureFromSurface(alternateMapSurface)
		if err != nil {
			errorString := fmt.Sprintf("trying to create alternateMapTexture from Surface:\n%v", err)
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, window)
			return err
		}
		defer alternateMapTexture.Destroy()
	}
	pbInstance.Increment()

	mapRect := &sdl.Rect{mapXOffset, mapYOffset, mapOverviewWidth, mapOverviewHeight}

	app := app{
		window:                      window,
		renderer:                    renderer,
		mapTexture:                  mapTexture,
		alternateMapTexture:         alternateMapTexture,
		mapRect:                     mapRect,
		font:                        font,
		match:                       match,
		config:                      c,
		lastDrawnAt:                 time.Now().UTC(),
		isPaused:                    false,
		curFrame:                    0,
		isOnNormalElevation:         true,
		playbackSpeedModifier:       1,
		staticPlaybackSpeedModifier: 1,
		hidePlayerNames:             false,
	}

	pbInstance.Set("step", "Starting app")
	pbInstance.Finish()

	err = app.run()
	if err != nil {
		errorString := fmt.Sprintf("while running the app:\n%v", err)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error", errorString, nil)
		return err
	}

	return nil
}

func (app *app) run() error {
	m := app.match
	// MAIN GAME LOOP
	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch eventT := event.(type) {
			case *sdl.QuitEvent:
				return nil

			case *sdl.KeyboardEvent:
				if eventT.Type != sdl.KEYDOWN {
					break
				}
				app.handleKeyboardEvents(eventT)
				if eventT.Keysym.Sym == sdl.K_c {
					if common.MapHasAlternateVersion(m.MapName) {
						tmp := app.mapTexture
						app.mapTexture = app.alternateMapTexture
						app.alternateMapTexture = tmp
						app.isOnNormalElevation = !app.isOnNormalElevation
					}
				}

			case *sdl.MouseWheelEvent:
				// back
				if eventT.Type == sdl.MOUSEWHEEL {
					if eventT.Y > 0 {
						if app.curFrame < m.FrameRate*1 {
							app.curFrame = 0
						} else {
							app.curFrame -= m.FrameRate * 1
						}
					}
					if eventT.Y < 0 {
						// forward
						if app.curFrame+m.FrameRate*1 > len(m.States)-1 {
							app.curFrame = len(m.States) - 1
						} else {
							app.curFrame += m.FrameRate * 1
						}
					}
				}
			}

		}

		if app.isPaused {
			sdl.Delay(32)
			app.updateGraphics()
			app.updateWindowTitle()
			continue
		}

		app.updateGraphics()
		app.updateWindowTitle()
		app.playbackSpeedModifier = 1

		// frameDuration is in ms
		frameDuration := float64(time.Since(frameStart) / 1000000)
		keyboardState := sdl.GetKeyboardState()
		if keyboardState[sdl.GetScancodeFromKey(sdl.K_w)] != 0 &&
			keyboardState[sdl.GetScancodeFromKey(sdl.K_LSHIFT)] != 0 {
			app.playbackSpeedModifier = 5
		}
		if keyboardState[sdl.GetScancodeFromKey(sdl.K_s)] != 0 &&
			keyboardState[sdl.GetScancodeFromKey(sdl.K_LSHIFT)] != 0 {
			app.playbackSpeedModifier = 0.5
		}
		delay := (1/(app.playbackSpeedModifier*app.staticPlaybackSpeedModifier))*(1000/float64(m.FrameRate)) - frameDuration
		if delay < 0 {
			delay = 0
		}
		sdl.Delay(uint32(delay))
		if app.curFrame < len(m.States)-1 {
			app.curFrame = app.curFrame + 1
		}
	}

}

func (app *app) handleKeyboardEvents(eventT *sdl.KeyboardEvent) {
	m := app.match
	switch eventT.Keysym.Sym {
	case sdl.K_SPACE:
		app.isPaused = !app.isPaused

	case sdl.K_a:
		if isShiftPressed(eventT) {
			if app.curFrame < m.FrameRate*10 {
				app.curFrame = 0
			} else {
				app.curFrame -= m.FrameRate * 10
			}
		} else {
			if app.curFrame < m.FrameRate*3 {
				app.curFrame = 0
			} else {
				app.curFrame -= m.FrameRate * 3
			}
		}

	case sdl.K_d:
		if isShiftPressed(eventT) {
			if app.curFrame+m.FrameRate*10 > len(m.States)-1 {
				app.curFrame = len(m.States) - 1
			} else {
				app.curFrame += m.FrameRate * 10
			}
		} else {
			if app.curFrame+m.FrameRate*3 > len(m.States)-1 {
				app.curFrame = len(m.States) - 1
			} else {
				app.curFrame += m.FrameRate * 3
			}
		}

	case sdl.K_q:
		if isShiftPressed(eventT) {
			set := false
			for i, frame := range m.HalfStarts {
				if app.curFrame < frame {
					if i > 1 && app.curFrame < m.HalfStarts[i-1]+m.FrameRate/2 {
						app.curFrame = m.HalfStarts[i-2]
						set = true
						break
					}
					if i-1 < 0 {
						app.curFrame = 0
						set = true
						break
					}
					app.curFrame = m.HalfStarts[i-1]
					set = true
					break
				}
			}
			// not set -> last round of match
			if !set {
				if len(m.HalfStarts) > 1 && app.curFrame < m.HalfStarts[len(m.HalfStarts)-1]+m.FrameRate/2 {
					app.curFrame = m.HalfStarts[len(m.HalfStarts)-2]
				} else {
					app.curFrame = m.HalfStarts[len(m.HalfStarts)-1]
				}
			}
		} else {
			set := false
			for i, frame := range m.RoundStarts {
				if app.curFrame < frame {
					if i > 1 && app.curFrame < m.RoundStarts[i-1]+m.FrameRate/2 {
						app.curFrame = m.RoundStarts[i-2]
						set = true
						break
					}
					if i-1 < 0 {
						app.curFrame = 0
						set = true
						break
					}
					app.curFrame = m.RoundStarts[i-1]
					set = true
					break
				}
			}
			// not set -> last round of match
			if !set {
				if len(m.RoundStarts) > 1 && app.curFrame < m.RoundStarts[len(m.RoundStarts)-1]+m.FrameRate/2 {
					app.curFrame = m.RoundStarts[len(m.RoundStarts)-2]
				} else {
					app.curFrame = m.RoundStarts[len(m.RoundStarts)-1]
				}
			}
		}

	case sdl.K_e:
		if isShiftPressed(eventT) {
			for _, frame := range m.HalfStarts {
				if app.curFrame < frame {
					app.curFrame = frame
					break
				}
			}
		} else {
			for _, frame := range m.RoundStarts {
				if app.curFrame < frame {
					app.curFrame = frame
					break
				}
			}
		}

	case sdl.K_h:
		app.hidePlayerNames = !app.hidePlayerNames
	case sdl.K_0:
		app.copyPositionToClipboard(9)
	case sdl.K_1:
		app.copyPositionToClipboard(0)
	case sdl.K_2:
		app.copyPositionToClipboard(1)
	case sdl.K_3:
		app.copyPositionToClipboard(2)
	case sdl.K_4:
		app.copyPositionToClipboard(3)
	case sdl.K_5:
		app.copyPositionToClipboard(4)
	case sdl.K_6:
		app.copyPositionToClipboard(5)
	case sdl.K_7:
		app.copyPositionToClipboard(6)
	case sdl.K_8:
		app.copyPositionToClipboard(7)
	case sdl.K_9:
		app.copyPositionToClipboard(8)

	case sdl.K_w:
		if isShiftPressed(eventT) {
			break
		}
		switch app.staticPlaybackSpeedModifier {
		case 0.25:
			app.staticPlaybackSpeedModifier = 0.5
		case 0.5:
			app.staticPlaybackSpeedModifier = 1
		case 1:
			app.staticPlaybackSpeedModifier = 1.25
		case 1.25:
			app.staticPlaybackSpeedModifier = 1.5
		case 1.5:
			app.staticPlaybackSpeedModifier = 2
		}

	case sdl.K_s:
		if isShiftPressed(eventT) {
			break
		}
		switch app.staticPlaybackSpeedModifier {
		case 0.5:
			app.staticPlaybackSpeedModifier = 0.25
		case 1:
			app.staticPlaybackSpeedModifier = 0.5
		case 1.25:
			app.staticPlaybackSpeedModifier = 1
		case 1.5:
			app.staticPlaybackSpeedModifier = 1.25
		case 2:
			app.staticPlaybackSpeedModifier = 1.5
		}

	case sdl.K_r:
		app.staticPlaybackSpeedModifier = 1

	case sdl.K_k:
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "Hotkeys", hotkeysString, nil)
	}
}

func (app *app) copyPositionToClipboard(player int) {
	m := app.match
	players := m.States[app.curFrame].Players
	if player >= len(players) {
		return
	}

	sort.Slice(players, func(i, j int) bool {
		if players[i].Team > players[j].Team {
			return true
		}
		if players[i].Team < players[j].Team {
			return false
		}
		return players[i].ID < players[j].ID
	})

	clipboard.WriteAll("setpos " +
		strconv.FormatFloat(float64(m.States[app.curFrame].Players[player].Position.X), 'f', 2, 32) + " " +
		strconv.FormatFloat(float64(m.States[app.curFrame].Players[player].Position.Y), 'f', 2, 32) + " " +
		strconv.FormatFloat(float64(m.States[app.curFrame].Players[player].Position.Z), 'f', 2, 32) +
		";setang " +
		strconv.FormatFloat(float64(m.States[app.curFrame].Players[player].ViewDirectionY), 'f', 2, 32) + " " +
		strconv.FormatFloat(float64(m.States[app.curFrame].Players[player].ViewDirectionX), 'f', 2, 32))
}

func (app *app) updateWindowTitle() {
	m := app.match
	cts := m.States[app.curFrame].TeamCounterTerrorists
	ts := m.States[app.curFrame].TeamTerrorists
	clanNameCTs := cts.ClanName
	if clanNameCTs == "" {
		clanNameCTs = "Counter Terrorists"
	}
	clanNameTs := ts.ClanName
	if clanNameTs == "" {
		clanNameTs = "Terrorists"
	}
	windowTitle := fmt.Sprintf("%s  [%d:%d]  %s - Round %d", clanNameCTs, cts.Score, ts.Score, clanNameTs, cts.Score+ts.Score+1)
	// expensive?
	app.window.SetTitle(windowTitle)
}

func (app *app) updateGraphics() {
	m := app.match
	app.renderer.SetDrawColor(10, 10, 10, 255)
	app.renderer.Clear()
	app.lastDrawnAt = time.Now().UTC()

	app.drawInfobars()
	app.renderer.Copy(app.mapTexture, nil, app.mapRect)

	shots := m.Shots[app.curFrame]
	for _, shot := range shots {
		app.drawShot(&shot)
	}

	infernos := m.States[app.curFrame].Infernos
	for _, inferno := range infernos {
		app.drawInferno(&inferno)
	}

	effects := m.Effects[app.curFrame]
	for _, effect := range effects {
		app.drawEffects(&effect)
	}

	grenades := m.States[app.curFrame].Grenades
	for _, grenade := range grenades {
		app.drawGrenade(&grenade)
	}

	bomb := m.States[app.curFrame].Bomb
	app.drawBomb(&bomb)

	players := m.States[app.curFrame].Players
	var indexT, indexCT int
	for _, player := range players {
		if player.Team == demoinfo.TeamTerrorists {
			app.drawPlayer(&player, indexT)
			indexT++
		} else {
			app.drawPlayer(&player, indexCT)
			indexCT++
		}
	}

	app.drawString("k shows hotkeys", colorDarkGrey, mapXOffset+mapOverviewWidth+150, mapYOffset+mapOverviewHeight-40)

	app.renderer.Present()
}

func isShiftPressed(event *sdl.KeyboardEvent) bool {
	pressed := event.Keysym.Mod & sdl.KMOD_SHIFT

	if pressed > 0 {
		return true
	}
	return false
}
