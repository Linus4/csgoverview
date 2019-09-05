# csgoverview - 2D Demoviewer

A 2D demo replay tool for Counter Strike: Global Offensive.

Package match povides a high-level parser you can use for your own demoviewer.

Current version is `0.3.0`. Master branch is currently used for development.

[![GoDoc](https://godoc.org/github.com/Linus4/csgoverview?status.svg)](https://godoc.org/github.com/Linus4/csgoverview) [![Go Report Card](https://goreportcard.com/badge/github.com/linus4/csgoverview)](https://goreportcard.com/report/github.com/linus4/csgoverview)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Linus4/csgoverview/blob/master/LICENSE)

## Tools

* [golang](https://golang.org/)
* [SDL2](https://wiki.libsdl.org/Introduction)
* [go-sdl2](https://github.com/veandco/go-sdl2)
* [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang)
* [csgo-overviews](https://github.com/zoidbergwill/csgo-overviews)

## Installation

### Dependencies

#### Fedora

```sh
dnf install git golang SDL2-devel SDL2_gfx-devel SDL2_image-devel SDL2_ttf-devel
```

#### Ubuntu

```sh
sudo apt install git golang libsdl2-dev libsdl2-gfx-dev libsdl2-image-dev libsdl2-ttf-dev
```

### Build

This project uses go modules. Make sure you have go version `1.12` or higher
installed. Run `go version` to check.

```sh
git clone https://github.com/Linus4/csgoverview.git
cd csgoverview
go build
```

### Get overviews

Use [this repository](https://github.com/zoidbergwill/csgo-overviews)
(overviews directory) and copy the overviews that you need into the directory
you cloned.

You can use other overviews as long as they are `.jpg` files and they match the
naming pattern (e.g. `de_nuke.jpg`). Ideally, their size should be 1024x1024
pixels or larger.

More overviews are available here:

* [alternative generated from game
  files](https://github.com/CSGO-Analysis/csgo-maps-overviews)
* [Simple Radar](www.simpleradar.com)

On Linux, you can convert images with `convert image.png image.jpg` if you
have `ImageMagick` installed.

## Usage

1. you must be in the directory you cloned that contains the executable
1. the font (`liberationserif-regular.ttf`) must be in the same directory as the
  executable
1. the overview (`e.g. de_nuke.jpg`) must be in the same directory as the
  executable

```sh
./csgoverview [path to demo]
```

### Keybinds

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

![Screenshot 2 de_dust2](https://i.imgur.com/j3BDQhz.png)

![Screenshot 3 de_inferno](https://i.imgur.com/VrWOKzJ.png)

