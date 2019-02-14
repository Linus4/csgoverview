package main

import (
	"fmt"
	"log"
	"os"
	// "time"

	dem "github.com/markus-wa/demoinfocs-golang"
	event "github.com/markus-wa/demoinfocs-golang/events"
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

	// First pass to get round starts, half starts and header info
	demo, err := os.Open("test.dem")
	if err != nil {
		log.Println("Could not open test.dem")
		log.Println(err)
		return
	}
	defer demo.Close()

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

	frameTime := parser.Header().FrameTime()
	mapName := parser.Header().MapName

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

	// parser = dem.NewParser(demo)

	// states := make([]dem.IGameState, 0)

	// parse and present demo frame by frame (positions)

	// translate coordinates

	// draw the things

	fmt.Println("Time per frame: %v", frameTime)
	fmt.Println("Round starts:")
	for i, tick := range roundStarts {
		fmt.Printf("Round %v:\t%v\n", i, tick)
	}
	fmt.Println("Half starts:")
	for i, tick := range halfStarts {
		fmt.Printf("Half %v:\t%v\n", i, tick)
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		sdl.Delay(32)
	}

}
