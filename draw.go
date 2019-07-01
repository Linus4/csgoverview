package main

import (
	"fmt"
	"log"
	"sort"

	ocom "github.com/linus4/csgoverview/common"
	"github.com/linus4/csgoverview/match"
	common "github.com/markus-wa/demoinfocs-golang/common"
	meta "github.com/markus-wa/demoinfocs-golang/metadata"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	terrorR      uint8 = 252
	terrorG      uint8 = 176
	terrorB      uint8 = 12
	counterR     uint8 = 89
	counterG     uint8 = 206
	counterB     uint8 = 200
	radiusPlayer int32 = 10
)

var (
	colorTerror sdl.Color = sdl.Color{
		R: 252,
		G: 176,
		B: 12,
		A: 255,
	}
	colorCounter sdl.Color = sdl.Color{
		R: 89,
		G: 206,
		B: 200,
		A: 255,
	}
	colorMoney sdl.Color = sdl.Color{
		R: 45,
		G: 135,
		B: 45,
		A: 255,
	}
	colorBomb sdl.Color = sdl.Color{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}
	colorEqDecoy sdl.Color = sdl.Color{
		R: 102,
		G: 34,
		B: 0,
		A: 255,
	}
	colorEqMolotov sdl.Color = sdl.Color{
		R: 255,
		G: 153,
		B: 0,
		A: 255,
	}
	colorEqIncendiary sdl.Color = sdl.Color{
		R: 255,
		G: 153,
		B: 0,
		A: 255,
	}
	colorEqFlash sdl.Color = sdl.Color{
		R: 128,
		G: 170,
		B: 255,
		A: 255,
	}
	colorEqSmoke sdl.Color = sdl.Color{
		R: 153,
		G: 153,
		B: 153,
		A: 255,
	}
	colorEqHE sdl.Color = sdl.Color{
		R: 85,
		G: 150,
		B: 0,
		A: 255,
	}
)

func DrawPlayer(renderer *sdl.Renderer, player *ocom.OverviewPlayer, font *ttf.Font, match *match.Match) {
	pos := player.LastAlivePosition

	scaledX, scaledY := meta.MapNameToMap[match.MapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	var color sdl.Color

	if player.Team == common.TeamTerrorists {
		color = colorTerror
	} else {
		color = colorCounter
	}

	if player.Hp > 0 {
		gfx.CircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer, color)

		DrawString(renderer, player.Name, color, scaledXInt+10, scaledYInt+10, font)
		//gfx.StringRGBA(renderer, scaledXInt+15, scaledYInt+15, player.Name, colorR, colorG, colorB, 255)

		viewAngle := -int32(player.ViewDirectionX) // negated because of sdl
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer+1, viewAngle-20, viewAngle+20, 200, 200, 200, 255)
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer+2, viewAngle-10, viewAngle+10, 200, 200, 200, 255)
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer+3, viewAngle-5, viewAngle+5, 200, 200, 200, 255)

		// FlashDuration is not the time remaining but always the total amount of time flashed from a single flashbang
		if player.FlashDuration > 0.8 {
			gfx.FilledCircleRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer-5, 200, 200, 200, 200)
		}

		for _, w := range player.Weapons {
			if w.Weapon == common.EqBomb {
				gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer-1, 255, 0, 0, 255)
				gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer-2, 255, 0, 0, 255)
			}
		}

		if player.IsDefusing {
			color.A = 200
			gfx.CharacterColor(renderer, scaledXInt-radiusPlayer/4, scaledYInt-radiusPlayer/4, 'D', color)
			color.A = 255
		}
	} else {
		color.A = 150
		gfx.CharacterColor(renderer, scaledXInt, scaledYInt, 'X', color)
		color.A = 255
	}
}

func DrawGrenade(renderer *sdl.Renderer, grenade *common.GrenadeProjectile, match *match.Match) {
	pos := grenade.Position

	scaledX, scaledY := meta.MapNameToMap[match.MapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	var color sdl.Color

	switch grenade.Weapon {
	case common.EqDecoy:
		color = colorEqDecoy
	case common.EqMolotov:
		color = colorEqMolotov
	case common.EqIncendiary:
		color = colorEqIncendiary
	case common.EqFlash:
		color = colorEqFlash
	case common.EqSmoke:
		color = colorEqSmoke
	case common.EqHE:
		color = colorEqHE
	}

	gfx.BoxColor(renderer, scaledXInt-2, scaledYInt-3, scaledXInt+2, scaledYInt+3, color)
}

func DrawGrenadeEffect(renderer *sdl.Renderer, effect *ocom.GrenadeEffect, match *match.Match) {
	pos := effect.GrenadeEvent.Position

	scaledX, scaledY := meta.MapNameToMap[match.MapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	var colorR, colorG, colorB uint8

	switch effect.GrenadeEvent.GrenadeType {
	case common.EqFlash:
		colorR = 128
		colorG = 170
		colorB = 255
	case common.EqSmoke:
		colorR = 153
		colorG = 153
		colorB = 153
	case common.EqHE:
		colorR = 85
		colorG = 150
		colorB = 0
	}

	switch effect.GrenadeEvent.GrenadeType {
	case common.EqFlash:
		gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, int32(effect.Lifetime), colorR, colorG, colorB, 255)
	case common.EqHE:
		gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, int32(effect.Lifetime), colorR, colorG, colorB, 255)
	case common.EqSmoke:
		gfx.FilledCircleRGBA(renderer, scaledXInt, scaledYInt, 25, colorR, colorG, colorB, 100)
		// only draw the outline if the smoke is not fading
		if effect.Lifetime < 15*match.SmokeEffectLifetime/18 {
			gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, 25, colorR, colorG, colorB, 255)
		}
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, 10, int32(270+effect.Lifetime*360/match.SmokeEffectLifetime), 630, colorR, colorG, colorB, 255)
	}
}

func DrawInferno(renderer *sdl.Renderer, inferno *common.Inferno, match *match.Match) {
	hull := inferno.ConvexHull2D()
	var colorR, colorG, colorB uint8 = 255, 153, 0
	xCoordinates := make([]int16, 0)
	yCoordinates := make([]int16, 0)

	for _, v := range hull {
		scaledX, scaledY := meta.MapNameToMap[match.MapName].TranslateScale(v.X, v.Y)
		scaledXInt := int16(scaledX) + int16(mapXOffset)
		scaledYInt := int16(scaledY) + int16(mapYOffset)
		xCoordinates = append(xCoordinates, scaledXInt)
		yCoordinates = append(yCoordinates, scaledYInt)
	}

	gfx.FilledPolygonRGBA(renderer, xCoordinates, yCoordinates, colorR, colorG, colorB, 100)
	gfx.PolygonRGBA(renderer, xCoordinates, yCoordinates, colorR, colorG, colorB, 100)
}

func DrawBomb(renderer *sdl.Renderer, bomb *common.Bomb, match *match.Match) {
	pos := bomb.Position()
	if bomb.Carrier != nil {
		return
	}

	scaledX, scaledY := meta.MapNameToMap[match.MapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset

	gfx.BoxColor(renderer, scaledXInt-3, scaledYInt-2, scaledXInt+3, scaledYInt+2, colorBomb)
}

func DrawString(renderer *sdl.Renderer, text string, color sdl.Color, x, y int32, font *ttf.Font) {
	textSurface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		log.Fatal(err)
	}
	defer textSurface.Free()
	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		log.Fatal(err)
	}
	defer textTexture.Destroy()
	textRect := &sdl.Rect{
		X: x,
		Y: y,
		W: textSurface.W,
		H: textSurface.H,
	}
	err = renderer.Copy(textTexture, nil, textRect)
	if err != nil {
		log.Fatal(err)
	}
}

func DrawInfobars(renderer *sdl.Renderer, match *match.Match, font *ttf.Font) {
	var cts, ts []ocom.OverviewPlayer
	for _, player := range match.States[curFrame].Players {
		if player.Team == common.TeamCounterTerrorists {
			cts = append(cts, player)
			sort.Slice(cts, func(i, j int) bool { return cts[i].SteamID < cts[j].SteamID })

		} else {
			ts = append(ts, player)
			sort.Slice(ts, func(i, j int) bool { return ts[i].SteamID < ts[j].SteamID })
		}
	}
	DrawInfobar(renderer, cts, 0, mapYOffset, colorCounter, font)
	DrawInfobar(renderer, ts, mapXOffset+mapOverviewWidth, mapYOffset, colorTerror, font)
}

func DrawInfobar(renderer *sdl.Renderer, players []ocom.OverviewPlayer, x, y int32, color sdl.Color, font *ttf.Font) {
	var yOffset int32 = 0
	for _, player := range players {
		if player.Hp > 0 {
			gfx.BoxColor(renderer, x+int32(player.Hp)*(mapXOffset/infobarElementHeight), yOffset, x, yOffset+5, color)
		}
		DrawString(renderer, player.Name, color, x+80, yOffset+10, font)
		DrawString(renderer, fmt.Sprintf("%v", player.Hp), color, x+5, yOffset+10, font)
		if player.Armor > 0 && player.HasHelmet {
			DrawString(renderer, "H", color, x+30, yOffset+10, font)
		} else if player.Armor > 0 {
			DrawString(renderer, "A", color, x+30, yOffset+10, font)
		}
		if player.HasDefuseKit {
			DrawString(renderer, "D", color, x+45, yOffset+10, font)
		}
		DrawString(renderer, fmt.Sprintf("%v $", player.Money), colorMoney, x+5, yOffset+25, font)
		var nadeCounter int32 = 0
		sort.Slice(player.Weapons, func(i, j int) bool { return player.Weapons[i].Weapon < player.Weapons[j].Weapon })
		for _, w := range player.Weapons {
			if w.Class() == common.EqClassSMG || w.Class() == common.EqClassHeavy || w.Class() == common.EqClassRifle {
				DrawString(renderer, w.Weapon.String(), color, x+150, yOffset+25, font)
			}
			if w.Class() == common.EqClassPistols {
				DrawString(renderer, w.Weapon.String(), color, x+150, yOffset+40, font)
			}
			if w.Class() == common.EqClassGrenade {
				var nadeColor sdl.Color
				switch w.Weapon {
				case common.EqDecoy:
					nadeColor = colorEqDecoy
				case common.EqMolotov:
					nadeColor = colorEqMolotov
				case common.EqIncendiary:
					nadeColor = colorEqIncendiary
				case common.EqFlash:
					// there seems to be only one flashbang in player.Weapons even if he has two
					nadeColor = colorEqFlash
				case common.EqSmoke:
					nadeColor = colorEqSmoke
				case common.EqHE:
					nadeColor = colorEqHE
				}

				gfx.BoxColor(renderer, x+150+nadeCounter*12, yOffset+60, x+150+nadeCounter*12+6, yOffset+60+9, nadeColor)
				nadeCounter++
			}
			if w.Class() == common.EqClassEquipment {
				if w.Weapon == common.EqBomb {
					gfx.BoxColor(renderer, x+45, yOffset+12, x+45+12, yOffset+12+9, colorBomb)
				}
			}
		}
		addInfo := player.AdditionalPlayerInformation
		kdaInfo := fmt.Sprintf("%v / %v / %v", addInfo.Kills, addInfo.Assists, addInfo.Deaths)
		DrawString(renderer, kdaInfo, color, x+5, yOffset+40, font)

		yOffset += infobarElementHeight
	}
}
