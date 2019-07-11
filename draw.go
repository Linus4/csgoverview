package main

import (
	"fmt"
	"log"
	"sort"
	"unicode/utf8"

	ocom "github.com/linus4/csgoverview/common"
	"github.com/linus4/csgoverview/match"
	common "github.com/markus-wa/demoinfocs-golang/common"
	event "github.com/markus-wa/demoinfocs-golang/events"
	meta "github.com/markus-wa/demoinfocs-golang/metadata"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	radiusPlayer   int32 = 10
	killfeedHeight int32 = 15
)

var (
	colorTerror       sdl.Color = sdl.Color{252, 176, 12, 255}
	colorCounter      sdl.Color = sdl.Color{89, 206, 200, 255}
	colorMoney        sdl.Color = sdl.Color{45, 135, 45, 255}
	colorBomb         sdl.Color = sdl.Color{255, 0, 0, 255}
	colorEqDecoy      sdl.Color = sdl.Color{102, 34, 0, 255}
	colorEqMolotov    sdl.Color = sdl.Color{255, 153, 0, 255}
	colorEqIncendiary sdl.Color = sdl.Color{255, 153, 0, 255}
	colorInferno      sdl.Color = sdl.Color{255, 153, 0, 100}
	colorEqFlash      sdl.Color = sdl.Color{128, 170, 255, 255}
	colorEqSmoke      sdl.Color = sdl.Color{153, 153, 153, 255}
	colorSmoke        sdl.Color = sdl.Color{153, 153, 153, 100}
	colorEqHE         sdl.Color = sdl.Color{85, 150, 0, 255}
	colorDarkWhite    sdl.Color = sdl.Color{200, 200, 200, 255}
	colorFlashEffect  sdl.Color = sdl.Color{200, 200, 200, 180}
)

func DrawPlayer(renderer *sdl.Renderer, player *common.Player, font *ttf.Font, match *match.Match) {
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

		viewAngle := -int32(player.ViewDirectionX) // negated because of sdl
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, radiusPlayer+1, viewAngle-20, viewAngle+20, colorDarkWhite)
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, radiusPlayer+2, viewAngle-10, viewAngle+10, colorDarkWhite)
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, radiusPlayer+3, viewAngle-5, viewAngle+5, colorDarkWhite)

		// FlashDuration is not the time remaining but always the total amount of time flashed from a single flashbang
		if player.FlashDuration > 0.8 {
			gfx.FilledCircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer-5, colorFlashEffect)
		}

		for _, w := range player.Weapons() {
			if w.Weapon == common.EqBomb {
				gfx.CircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer-1, colorBomb)
				gfx.CircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer-2, colorBomb)
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

	switch effect.GrenadeEvent.GrenadeType {
	case common.EqFlash:
		gfx.CircleColor(renderer, scaledXInt, scaledYInt, int32(effect.Lifetime), colorEqFlash)
	case common.EqHE:
		gfx.CircleColor(renderer, scaledXInt, scaledYInt, int32(effect.Lifetime), colorEqHE)
	case common.EqSmoke:
		gfx.FilledCircleColor(renderer, scaledXInt, scaledYInt, 25, colorSmoke)
		// only draw the outline if the smoke is not fading
		if effect.Lifetime < 15*match.SmokeEffectLifetime/18 {
			gfx.CircleColor(renderer, scaledXInt, scaledYInt, 25, colorDarkWhite)
		}
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, 10, int32(270+effect.Lifetime*360/match.SmokeEffectLifetime), 630, colorDarkWhite)
	}
}

func DrawInferno(renderer *sdl.Renderer, inferno *common.Inferno, match *match.Match) {
	hull := inferno.ConvexHull2D()
	xCoordinates := make([]int16, 0)
	yCoordinates := make([]int16, 0)

	for _, v := range hull {
		scaledX, scaledY := meta.MapNameToMap[match.MapName].TranslateScale(v.X, v.Y)
		scaledXInt := int16(scaledX) + int16(mapXOffset)
		scaledYInt := int16(scaledY) + int16(mapYOffset)
		xCoordinates = append(xCoordinates, scaledXInt)
		yCoordinates = append(yCoordinates, scaledYInt)
	}

	gfx.FilledPolygonColor(renderer, xCoordinates, yCoordinates, colorInferno)
	gfx.PolygonColor(renderer, xCoordinates, yCoordinates, colorInferno)
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
	var cts, ts []common.Player
	for _, player := range match.States[curFrame].Players {
		if player.Team == common.TeamCounterTerrorists {
			cts = append(cts, player)

		} else {
			ts = append(ts, player)
		}
	}
	sort.Slice(cts, func(i, j int) bool { return cts[i].SteamID < cts[j].SteamID })
	sort.Slice(ts, func(i, j int) bool { return ts[i].SteamID < ts[j].SteamID })
	DrawInfobar(renderer, cts, 0, mapYOffset, colorCounter, font)
	DrawInfobar(renderer, ts, mapXOffset+mapOverviewWidth, mapYOffset, colorTerror, font)
	DrawKillfeed(renderer, match.Killfeed[curFrame], mapXOffset+mapOverviewWidth, mapYOffset+600, font)
}

func DrawInfobar(renderer *sdl.Renderer, players []common.Player, x, y int32, color sdl.Color, font *ttf.Font) {
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
		weapons := player.Weapons()
		sort.Slice(weapons, func(i, j int) bool { return weapons[i].Weapon < weapons[j].Weapon })
		for _, w := range weapons {
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
					// there seems to be only one flashbang in player.Weapons() even if he has two
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

func DrawKillfeed(renderer *sdl.Renderer, killfeed []event.Kill, x, y int32, font *ttf.Font) {
	var yOffset int32 = 0
	for _, kill := range killfeed {
		var colorKiller, colorVictim sdl.Color
		if kill.Killer.Team == common.TeamCounterTerrorists {
			colorKiller = colorCounter
		} else {
			colorKiller = colorTerror
		}
		if kill.Victim.Team == common.TeamCounterTerrorists {
			colorVictim = colorCounter
		} else {
			colorVictim = colorTerror
		}
		killerName := kill.Killer.Name
		if utf8.RuneCountInString(kill.Killer.Name) > 15 {
			killerRunes := []rune(kill.Killer.Name)
			killerName = string(killerRunes[:15])
		}
		victimName := kill.Victim.Name
		if utf8.RuneCountInString(kill.Victim.Name) > 15 {
			victimRunes := []rune(kill.Victim.Name)
			victimName = string(victimRunes[:15])
		}
		weaponName := kill.Weapon.Weapon.String()
		if len(kill.Weapon.Weapon.String()) > 10 {
			weaponName = kill.Weapon.Weapon.String()[:10]
		}
		DrawString(renderer, killerName, colorKiller, x+5, y+yOffset, font)
		DrawString(renderer, weaponName, colorDarkWhite, x+110, y+yOffset, font)
		DrawString(renderer, victimName, colorVictim, x+185, y+yOffset, font)
		yOffset += killfeedHeight
	}
}
