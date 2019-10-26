# csgoverview - 2D Demoviewer

A 2D demo replay tool for Counter Strike: Global Offensive.

Package match povides a high-level parser you can use for your own demoviewer.

Current version is `0.5.0`. Master branch is currently used for development.

[![GoDoc](https://godoc.org/github.com/Linus4/csgoverview?status.svg)](https://godoc.org/github.com/Linus4/csgoverview) [![Go Report Card](https://goreportcard.com/badge/github.com/linus4/csgoverview)](https://goreportcard.com/report/github.com/linus4/csgoverview)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Linus4/csgoverview/blob/master/LICENSE)

Check out the [Roadmap](https://github.com/Linus4/csgoverview/projects/1) where
I keep track of ideas and todos.

## Tools

* [golang](https://golang.org/)
* [SDL2](https://wiki.libsdl.org/Introduction)
* [go-sdl2](https://github.com/veandco/go-sdl2)
* [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang)
* [csgo-overviews](https://github.com/zoidbergwill/csgo-overviews)

## Hardware Requirements

* Display with a resolution of **1920x1080** or higher
* about 1.5 GB of memory for a 32 tick demo
* about 5 GB of memory for a 128 tick demo

## Installation

### Windows

I did not sign the application - Windows *will* prevent the app from running.

1. Download the latest compiled version (`.zip`) from the [releases
  page](https://github.com/Linus4/csgoverview/releases).
1. Create a folder and extract the `.zip` file it.
1. Download the overview images from [this
   repository](https://github.com/zoidbergwill/csgo-overviews) and put them
   into the same folder.
1. Put the demo you want to watch into the same folder.
1. Right click the demo; Open with; Choose another app; check 'Always use this
   app to open .dem files'; More apps; Look for another app on this PC; select
   `csgoverview.exe` (weirdly, this does not open the demo yet).
1. Double click any demo in the folder to open it with csgoverview.
1. *alternatively*, you can launch the app from the command-line.

### Dependencies

#### Fedora

```sh
dnf install git golang SDL2{,_gfx,_image,_ttf}-devel dejavu-sans-fonts
```

#### Ubuntu

```sh
sudo apt install git golang libsdl2{,-gfx,-image,-ttf}-dev fonts-dejavu
```

### Build

This project uses go modules. Make sure you have go version `1.12` or higher
installed. Run `go version` to check.

```sh
git clone https://github.com/Linus4/csgoverview.git
cd csgoverview
go build
```

### Cross-compiling for Windows

Using a Fedora 30 machine:

#### Dependencies

```sh
sudo dnf install make git golang SDL2{,_gfx,_image,_ttf,_mixer} mingw64-SDL2{,_image,_ttf}
```

#### Installing SDL2_gfx library

```sh
wget http://www.ferzkopp.net/Software/SDL2_gfx/SDL2_gfx-1.0.4.tar.gz
tar xf SDL2_gfx-1.0.4.tar.gz
cd SDL2_gfx-1.0.4
mingw64-configure
mingw64-make
sudo mingw64-make install
cd ..
```

#### Cloning the repository

```sh
git clone https://github.com/Linus4/csgoverview.git
cd csgoverview
```

#### Switch to new go-sdl2 version

Edit line in `go.mod`:

```sh
-       github.com/veandco/go-sdl2 v0.3.3
+       github.com/veandco/go-sdl2 v0.4.0-rc.0
```

#### Build

```sh
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -tags static -ldflags "-s -w"
```

#### Required files

Put the font file (`.ttf`) from the repository, the required map overviews
(next section) and the executable (`.exe`) into a directory. You can launch the
application from the command-line or by double-clicking on a demo in the same
directory after you have set `csgoverview.exe` to be your default app to use
for `.dem` files.

### Get overviews

Use [this repository](https://github.com/zoidbergwill/csgo-overviews)
(overviews directory), create a directory with `mkdir
$HOME/.local/share/csgoverview`  and copy the overviews that you need to
`$HOME/.local/share/csgoverview`.

You can use other overviews as long as they are `.jpg` files and they match the
naming pattern (e.g. `de_nuke.jpg`). Ideally, their size should be 1024x1024
pixels or larger.

More overviews are available here:

* [alternative generated from game
  files](https://github.com/CSGO-Analysis/csgo-maps-overviews)
* [Simple Radar](www.simpleradar.com)

On Linux, you can convert images with `convert image.png image.jpg` if you
have `ImageMagick` installed.

### Executable

You can move or symlink the executable into a directory in your `$PATH` to make
the program available everywhere on your system.

Example:

```sh
sudo ln -s /usr/bin/csgoverview <path to cloned repository>/csgoverview
```

### Desktop file (Linux)

In order to add csgoverview to your applications menu, create a `.desktop`
file (use the path to the executable on your computer in Exec):

```sh
echo "[Desktop Entry]
Name=CSGOverview
Exec=/usr/bin/csgoverview %F
Type=Application
Terminal=false
Categories=Games;" > $HOME/.local/share/applications/csgoverview.desktop
```

## Usage

```sh
./csgoverview
    -fontpath string
    	Path to font file (.ttf) (default "/usr/share/fonts/dejavu/DejaVuSans.ttf")
    -framerate float
    	Fallback GOTV Framerate (default -1)
    -tickrate float
    	Fallback Gameserver Tickrate (default -1)

  [path to demo]
```

If you're using GTK, you can also double-click the executable and select a
demo with the GUI.

After you've created the `.desktop` file, you can right click on a demo and
select csgoverview when you select 'Open With'. This way, you can just
double-click on a demo to watch it with csgoverview.

Looking for font file in the following directories: the one you supply with the
`-fontpath` flag, `/usr/share/fonts/dejavu/DejaVuSans.ttf`,
`./DejaVuSans.ttf`.

Looking for overview file in the following directories:
`$HOME/.local/share/csgoverview` and in the current directory.


### Keybinds

* a -> 5 s backwards
* d -> 5 s forwards
* A -> 10 s backwards
* D -> 10 s forwards
* w -> hold to speed up 5 x
* s -> hold to slow down to 0.5 x
* q -> round backwards
* e -> round forwards
* Q -> to start of previous half
* E -> to start of next half
* space -> toggle pause

![Screenshot 1 de_mirage](https://i.imgur.com/BKTTBfW.png)

![Screenshot 2 de_dust2](https://i.imgur.com/2kfkpvP.png)

![Screenshot 3 de_inferno](https://i.imgur.com/sNYT4eH.png)

## Credits

Thank you for helping me or contributing to the project!

* [markus-wa](https://github.com/markus-wa)
  ([demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang))
* [veeableful](https://github.com/veeableful)
  ([go-sdl2](https://github.com/veandco/go-sdl2/))
