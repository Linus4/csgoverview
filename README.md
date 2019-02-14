A 2D demo replay tool for Counter Strike: Global Offensive.

# Tools

* [golang](https://golang.org/)
* [SDL2](https://wiki.libsdl.org/Introduction)
* [go-sdl2](https://github.com/veandco/go-sdl2)
* [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang)
* [csgo-overviews](https://github.com/zoidbergwill/csgo-overviews)

# Roadmap

## Milestone v1.0.0

### Demo playback

* Map
* Playerpositions
* Playernames
* Playernumbers
* Player-line of vision
* Shots
* Grenades during flight
* NadeTails
* Effects for grenades
* Fade-out effect for smokes
* Timer for mollys and smokes
* `x` at places where players died
* Indicator for flash-effect / -duration

### Keybinds

* a -> 10 s backwards
* d -> 10 s forwards
* A -> 30 s backwards
* D -> 30 s forwards
* w -> hold to speed up 10 x
* s -> hold to slow down to 0.5 x
* W -> toggle 2 x speed
* S -> toggle 0.33 x speed
* q -> round backwards
* e -> round forwards
* Q -> to start of last half
* E -> to start of next half
* space -> toggle pause
* [num]g -> go to round [num]
* o -> open demo
* p -> take screenshot
* i -> export gif

### Additional information about round and players

* #round / #total
* Score
* Teamnames
* Warmup-indicator
* Freezetime timer
* Time remaining in current round
* Bombplant Indicator
* Bomb timer
* Defuse timer
* Killfeed
* Player details (left/right)
    - Name
    - Armor
    - Helmet
    - Primary
    - Secondary
    - Grenades
    - Defkit
    - (Taser)
    - Money
    - Kills in current round
* Results of previous rounds (survivors?)

### Misc

* Splashscreen with keybinds
* Scaling
* Multi-platform
* Additional information about players optional - instead letters for weapons
  instead of playernumbers?
* tickrate <-> framerate <-> real time?

## Milestone v2.0.0

### Screenshot export

* jpg of current view

### Gif-Export of a single round

* from
* until

### Command Line Interface

* take Screenshot
* export Gif

## Milestone v3.0.0

### Configurable keybinds and scaling

* hardcoded defaults

# Non-Features

### GUI / Buttons

* Keybinds

### Networking

* join.me
* screenshare
* stream
* gif export
* cli

### Drawing

* twiddla
* paint

### Analysis

* CSGO Demos Manager
* other 3rd-party tools
