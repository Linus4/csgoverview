A 2D demo replay tool for Counter Strike: Global Offensive.

# Tools

* [golang](https://golang.org/)
* [SDL2](https://wiki.libsdl.org/Introduction)
* [go-sdl2](https://github.com/veandco/go-sdl2)
* [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang)
* [csgo-overviews](https://github.com/zoidbergwill/csgo-overviews)

# Roadmap

## Milestone v0.1.0

### Demo playback

* [X] Map
* [X] Playerpositions
* [X] Playernames
* [X] Player-line of vision
* [ ] ~~Shots~~
* [X] Grenades during flight
* [ ] ~~NadeTails~~
* [X] Effects for flashbangs and hes
* [X] Effects for smokes
* [X] Effects for mollys
* [X] Fade-out effect for smokes
* [X] Timer for smokes
* [X] `x` at places where players died
* [X] Indicator for flash-effect / -duration
* [X] Indicator for defusing player
* [ ] C4 on player, on the ground, planted and defused
* [ ] Smoke radius scaling (map metadata)

### Keybinds

* [X] a -> 10 s backwards
* [X] d -> 10 s forwards
* [X] A -> 30 s backwards
* [X] D -> 30 s forwards
* [X] w -> hold to speed up 10 x
* [X] s -> hold to slow down to 0.5 x
* [X] q -> round backwards
* [X] e -> round forwards
* [ ] Q -> to start of last half
* [ ] E -> to start of next half
* [X] space -> toggle pause
* [ ] [num]g -> go to round [num]
* [ ] p -> take screenshot
* [ ] i -> export gif/mp4

### Screenshot export

* jpg of current frame

### Gif/ mp4 export of a single round

* from
* until

## Milestone v0.2.0

### Misc

* [ ] Scaling
* [ ] Multi-platform
 

### Additional information about round and players

* [ ] #round / #total
* [ ] Score
* [ ] Teamnames
* [ ] Warmup-indicator
* [ ] Freezetime timer
* [ ] Time remaining in current round
* [ ] Bombplant Indicator
* [ ] Bomb timer
* [ ] Defuse timer
* [ ] Killfeed
* [ ] Player details (left/right)
    - Name
    - Hp
    - Armor
    - Helmet
    - Primary
    - Secondary
    - Grenades
    - Defkit
    - (Taser)
    - Money
    - Kills in current round
* [ ] Results of previous rounds (survivors?)

### Command Line Interface

* take Screenshot
* export Gif

## Milestone v0.3.0

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
* live overlay ?

### Analysis

* CSGO Demos Manager
* other 3rd-party tools
