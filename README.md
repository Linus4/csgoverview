A 2D demo replay tool for Counter Strike: Global Offensive.

# Tools

* [golang](https://golang.org/)
* [SDL2](https://wiki.libsdl.org/Introduction)
* [go-sdl2](https://github.com/veandco/go-sdl2)
* [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang)
* [csgo-overviews](https://github.com/zoidbergwill/csgo-overviews)

# Installation

This project uses go modules, so make sure you have go version `1.11` or higher
installed. Run `go version` to check.

## Dependencies

### Fedora

```sh
dnf install git golang SDL2-devel SDL2_gfx-devel SDL2_image-devel SDL2_ttf-devel
```

## Build

```sh
git clone https://github.com/Linus4/csgoverview.git
cd csgoverview
go build
```

## Get overviews

Use [this repository](https://github.com/zoidbergwill/csgo-overviews)
(overviews directory) and copy the overviews that you need into the directory
you cloned.

# Usage

* the demo you want to watch must be in the directory you cloned
* you must be in the directory you cloned

```sh
./csgoverview [demoname]
```

## Keybinds

* a -> 10 s backwards
* d -> 10 s forwards
* A -> 30 s backwards
* D -> 30 s forwards
* w -> hold to speed up 5 x
* s -> hold to slow down to 0.5 x
* q -> round backwards
* e -> round forwards
* Q -> to start of previous half
* E -> to start of next half
* space -> toggle pause

![Screenshot 1 de_dust2](https://i.imgur.com/FpPy5WV.png)

![Screenshot 2 de_inferno](https://i.imgur.com/VrWOKzJ.png)

