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
* [X] Player line of vision
* [X] Grenades during flight
* [X] Effects for flashbangs and hes
* [X] Effects for smokes
* [X] Effects for mollys
* [X] Fade-out effect for smokes
* [X] Timer for smokes
* [X] `x` at places where players died
* [X] Indicator for flash-effect / -duration
* [X] Indicator for defusing player
* [X] C4 on the ground
* [X] C4 on player
* [ ] ~~NadeTails~~
* [ ] ~~Shots~~ (maybe fake-shots?)

### Keybinds

* [X] a -> 10 s backwards
* [X] d -> 10 s forwards
* [X] A -> 30 s backwards
* [X] D -> 30 s forwards
* [X] w -> hold to speed up 5 x
* [X] s -> hold to slow down to 0.5 x
* [X] q -> round backwards
* [X] e -> round forwards
* [X] Q -> to start of previous half
* [X] E -> to start of next half
* [X] space -> toggle pause

### Misc

* [X] font support for playernames (ttf)
* [X] Smoke radius scaling (map metadata)

## Milestone v0.2.0

### Additional information about round and players

* [X] Score
* [X] Teamnames
* [ ] Freezetime timer
* [ ] Time remaining in current round
* [ ] Bomb timer
* [X] Killfeed

### Player details (left/right)

* [X] Name
* [X] Hp
* [X] Armor
* [X] Helmet
* [X] Primary
* [X] Secondary
* [X] Grenades
* [X] Defkit
* [X] Money
* [X] C4
* [X] Kills / Assists / Deaths

## Milestone v0.3.0

### Misc

* [ ] Scaling (Resolution Independent Rendering)
* [ ] Build for windows

## Future

* Results of previous rounds + survivors?
* Fake shots?
* Bundle with overviews (map jpegs)
* Highlight active weapon
* Flatpak
* Icons for weapons


# Non-Features

### GUI / Buttons

* Keybinds

### Analysis

* CSGO Demos Manager
* other 3rd-party tools

### Networking?

* join.me
* screenshare
* stream
* gif export
* cli

### Drawing?

* twiddla
* paint
* live overlay ?

### Export

* Screenshot tools
* SimpleScreenRecorder
