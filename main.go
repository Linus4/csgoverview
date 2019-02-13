package main

import (
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winHeight int32 = 926
	winWidth  int32 = 926
)

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

	surface, err := img.Load("de_cache.jpg")
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

	// nil, nil stretches texture to fill the screen
	// err
	renderer.Clear()
	// err
	renderer.Copy(texture, nil, nil)
	renderer.Present()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		sdl.Delay(32)
	}

	/*
		demo, err := os.Open("test.dem")
		if err != nil {
			log.Println("Could not open test.dem")
			log.Println(err)
			return
		}
		defer demo.Close()

		parser := dem.NewParser(demo)
		fmt.Println("Created parser")

		states := make([]dem.IGameState, 0)
		b := true
		for b {
			// err!
			b, _ = parser.ParseNextFrame()
			states = append(states, parser.GameState())
		}
		for i, gs := range states {
			fmt.Printf("Frame %v: Grenades: %v\n", i, gs.TotalRoundsPlayed())
		}
		fmt.Println("Finished parsing the demo")
	*/

	// get information from header (map, tickrate)

	// find matchstart, round strats, half-time, overtime half starts

	// parse and present demo frame by frame (positions)

	// translate coordinates

	// draw entities
}
