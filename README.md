# csgoverview - 2D Demoviewer

A 2D demo replay tool for Counter Strike: Global Offensive.

Package match povides a high-level parser you can use for your own demoviewer.

Current version is `0.6.0`.

[![GoDoc](https://godoc.org/github.com/Linus4/csgoverview?status.svg)](https://godoc.org/github.com/Linus4/csgoverview) [![Go Report Card](https://goreportcard.com/badge/github.com/linus4/csgoverview)](https://goreportcard.com/report/github.com/linus4/csgoverview)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Linus4/csgoverview/blob/master/LICENSE) [![Paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donate_SM.gif)](https://www.paypal.me/linuswbr)

Check out the [Roadmap](https://github.com/Linus4/csgoverview/projects/1) where
I keep track of ideas and todos.

## Table of Contents

* [Hardware Requirements](#hardware-requirements)
* [Windows Installation](#windows-installation)
* [Linux Installation / Build Instructions](#linux-installation)
* [Get overviews](#get-overviews)
* [Keybinds](#keybinds)
* [Tool recommendations](#tool-recommendations)
* [Cross-compiling](#cross-compiling)
* [Credits](#credits)

## Hardware Requirements

* Display with a resolution of **1920x1080** or higher
* about 1.5 GB of memory for a 32 tick demo
* about 5 GB of memory for a 128 tick demo

## Windows Installation

I did not sign the application - Windows *will* prevent the app from running.

1. Download latest version from the [releases
   page](https://github.com/Linus4/csgoverview/releases).
1. Create a folder and extract `csgoverview.exe` into it.
1. Create a folder called 'csgoverview' in your user directory. (e.g.
   `C:\Users\Username\csgoverview`)
1. Move the `.ttf` file from the .zip archive into the csgoverview folder.
1. Download the overview images from [this
   repository](https://github.com/zoidbergwill/csgo-overviews) and put them
   into the csgoverview folder.
1. Right click a demo and select 'Open with' to open it with csgoverview.

## Linux Installation

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

## Get overviews

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

## Usage

```sh
./csgoverview
    -fontpath string
    	Path to font file (.ttf) (default "/usr/share/fonts/dejavu/DejaVuSans.ttf")
    -framerate float
    	Fallback GOTV Framerate (default -1)
    -overviewdir string
        Path to overview directory (default "$HOME/.local/share/csgoverview")
    -tickrate float
    	Fallback Gameserver Tickrate (default -1)

  [path to demo]
```

If you're using GTK, you can also double-click the executable and select a
demo with the GUI.

After you've created the `.desktop` file, you can right click on a demo and
select csgoverview when you select 'Open With'. This way, you can just
double-click on a demo to watch it with csgoverview.

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

## Tool recommendations

* [gInk](https://github.com/geovens/gInk): draw on the screen (windows, free
  software)
* [mixer.com](https://mixer.com/): streaming service without delay (FTL)
* [join.me](https://www.join.me/): streaming service without delay and drawing on the
  screen
* [Applications](https://askubuntu.com/questions/4428/how-can-i-record-my-screen)
  to record your screen (free software, linux)
* [Draw on your Screen GNOME Shell
  Extension](https://extensions.gnome.org/extension/1683/draw-on-you-screen/):
  draw on the screen (linux, GNOME, free software)
* [Gfycat.com](https://gfycat.com): share videos/gifs

## Cross-compiling

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

#### Build

```sh
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -tags static -ldflags "-s -w"
```

#### Required files

Put the font file (`.ttf`) from the repository and the required map overviews
(next section) into a folder called 'csgoverview' in the user-directory
(`C:\Users\Name\csgoverview`). You can launch the application from the
command-line or right-clicking a demo file and selecting 'Open with' to open
the file with csgoverview.


![Screenshot 1 de_mirage](https://i.imgur.com/BKTTBfW.png)

![Screenshot 2 de_dust2](https://i.imgur.com/2kfkpvP.png)

![Screenshot 3 de_inferno](https://i.imgur.com/sNYT4eH.png)

## Credits

Thank you for helping me or contributing to the project!

* [markus-wa](https://github.com/markus-wa)
  ([demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang))
* [veeableful](https://github.com/veeableful)
  ([go-sdl2](https://github.com/veandco/go-sdl2/))
