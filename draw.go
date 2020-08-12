package main

import (
	"fmt"
	"log"
	"math"
	"sort"

	common "github.com/linus4/csgoverview/common"
	"github.com/linus4/csgoverview/match"
	demoinfo "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	radiusPlayer      int32   = 10
	radiusPlayerFloat float64 = float64(radiusPlayer)
	radiusSmoke       float64 = 25
	killfeedHeight    int32   = 15
	shotLength        float64 = 1000
)

var (
	colorTerror       = sdl.Color{252, 176, 12, 255}
	colorCounter      = sdl.Color{89, 206, 200, 255}
	colorMoney        = sdl.Color{45, 135, 45, 255}
	colorBomb         = sdl.Color{255, 0, 0, 255}
	colorEqDecoy      = sdl.Color{102, 34, 0, 255}
	colorEqMolotov    = sdl.Color{255, 153, 0, 255}
	colorEqIncendiary = sdl.Color{255, 153, 0, 255}
	colorInferno      = sdl.Color{255, 153, 0, 100}
	colorEqFlash      = sdl.Color{128, 170, 255, 255}
	colorEqSmoke      = sdl.Color{153, 153, 153, 255}
	colorSmoke        = sdl.Color{153, 153, 153, 100}
	colorEqHE         = sdl.Color{85, 150, 0, 255}
	colorDarkWhite    = sdl.Color{200, 200, 200, 255}
	colorFlashEffect  = sdl.Color{200, 200, 200, 180}
	colorAwpShot      = sdl.Color{255, 50, 0, 255}
	hidePlayerNames   bool
)

func drawPlayer(renderer *sdl.Renderer, player *common.Player, font *ttf.Font, index int, match *match.Match) {
	var color sdl.Color
	if player.Team == demoinfo.TeamTerrorists {
		color = colorTerror
	} else {
		color = colorCounter
	}
	colorLOS := colorDarkWhite
	colorC4 := colorBomb

	if player.IsAlive {
		pos := player.Position

		scaledX, scaledY := match.TranslateScale(pos.X, pos.Y)
		var scaledXInt int32 = int32(scaledX) + mapXOffset
		var scaledYInt int32 = int32(scaledY) + mapYOffset

		if common.MapHasAlternateVersion(match.MapName) && (!player.IsOnNormalElevation && isOnNormalElevation) ||
			(player.IsOnNormalElevation && !isOnNormalElevation) {
			color.A = 100
			colorLOS.A = 100
			colorC4.A = 100
		}
		gfx.AACircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer, color)

		var name string
		if !hidePlayerNames {
			name = fmt.Sprintf("%v %v", index+1, player.Name)
		} else {
			name = fmt.Sprintf("%v", index+1)
		}
		drawString(renderer, cropStringToN(name, 12), color, scaledXInt+10, scaledYInt+10, font)

		viewAngle := -int32(player.ViewDirectionX) // negated because of sdl
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, radiusPlayer+1, viewAngle-20, viewAngle+20, colorLOS)
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, radiusPlayer+2, viewAngle-10, viewAngle+10, colorLOS)
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, radiusPlayer+3, viewAngle-5, viewAngle+5, colorLOS)

		colorFlash := colorFlashEffect
		if player.FlashDuration.Seconds() > 0.5 {
			remaining := player.FlashTimeRemaining
			if remaining.Seconds() >= 3.1 {
				colorFlash.A = 255
			} else {
				colorFlash.A = uint8((remaining.Seconds() * 255) / 3.1)
			}
			if common.MapHasAlternateVersion(match.MapName) && (!player.IsOnNormalElevation && isOnNormalElevation) ||
				(player.IsOnNormalElevation && !isOnNormalElevation) {
				colorFlash.A /= 2
			}
			gfx.FilledCircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer-5, colorFlash)
		}

		if player.HasBomb {
			gfx.AACircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer-1, colorC4)
			gfx.AACircleColor(renderer, scaledXInt, scaledYInt, radiusPlayer-2, colorC4)
		}

		if player.IsDefusing {
			color.A -= 55
			gfx.CharacterColor(renderer, scaledXInt-radiusPlayer/4, scaledYInt-radiusPlayer/4, 'D', color)
			color.A += 55
		}
	} else {
		pos := player.LastAlivePosition

		scaledX, scaledY := match.TranslateScale(pos.X, pos.Y)
		var scaledXInt int32 = int32(scaledX) + mapXOffset
		var scaledYInt int32 = int32(scaledY) + mapYOffset

		color.A -= 105
		gfx.CharacterColor(renderer, scaledXInt, scaledYInt, 'X', color)
	}
}

func drawGrenade(renderer *sdl.Renderer, grenade *common.GrenadeProjectile, match *match.Match) {
	pos := grenade.Position

	scaledX, scaledY := match.TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	var color sdl.Color

	switch grenade.Type {
	case demoinfo.EqDecoy:
		color = colorEqDecoy
	case demoinfo.EqMolotov:
		color = colorEqMolotov
	case demoinfo.EqIncendiary:
		color = colorEqIncendiary
	case demoinfo.EqFlash:
		color = colorEqFlash
	case demoinfo.EqSmoke:
		color = colorEqSmoke
	case demoinfo.EqHE:
		color = colorEqHE
	}
	if common.MapHasAlternateVersion(match.MapName) && (!grenade.IsOnNormalElevation && isOnNormalElevation) ||
		(grenade.IsOnNormalElevation && !isOnNormalElevation) {
		color.A = 100
	}

	gfx.BoxColor(renderer, scaledXInt-2, scaledYInt-3, scaledXInt+2, scaledYInt+3, color)
}

func drawEffects(renderer *sdl.Renderer, effect *common.Effect, match *match.Match) {
	pos := effect.Position

	scaledX, scaledY := match.TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	var alphaModifier uint8
	if common.MapHasAlternateVersion(match.MapName) && (!effect.IsOnNormalElevation && isOnNormalElevation) ||
		(effect.IsOnNormalElevation && !isOnNormalElevation) {
		alphaModifier = 2
	} else {
		alphaModifier = 1
	}

	switch effect.Type {
	case demoinfo.EqFlash:
		color := colorEqFlash
		color.A /= alphaModifier
		gfx.AACircleColor(renderer, scaledXInt, scaledYInt, effect.Lifetime, color)
	case demoinfo.EqHE:
		color := colorEqHE
		color.A /= alphaModifier
		gfx.AACircleColor(renderer, scaledXInt, scaledYInt, effect.Lifetime, color)
	case demoinfo.EqSmoke:
		color := colorSmoke
		color.A /= alphaModifier
		colorCircles := colorDarkWhite
		colorCircles.A /= alphaModifier
		// 4.9 is the reference on Inferno for the value for radiusSmoke
		scaledRadiusSmoke := int32(radiusSmoke * 4.9 / float64(match.MapScale))
		gfx.FilledCircleColor(renderer, scaledXInt, scaledYInt, scaledRadiusSmoke, color)
		// only draw the outline if the smoke is not fading
		if effect.Lifetime < 15*match.SmokeEffectLifetime/18 {
			gfx.AACircleColor(renderer, scaledXInt, scaledYInt, scaledRadiusSmoke, colorCircles)
		}
		gfx.ArcColor(renderer, scaledXInt, scaledYInt, 10, 270+effect.Lifetime*360/match.SmokeEffectLifetime, 630, colorCircles)
	case demoinfo.EqDefuseKit:
		gfx.AACircleColor(renderer, scaledXInt, scaledYInt, effect.Lifetime, colorMoney)
	case demoinfo.EqBomb:
		gfx.AACircleColor(renderer, scaledXInt, scaledYInt, effect.Lifetime, colorBomb)
	}
}

func drawInferno(renderer *sdl.Renderer, inferno *common.Inferno, match *match.Match) {
	hull := inferno.ConvexHull2D
	color := colorInferno
	xCoordinates := make([]int16, 0)
	yCoordinates := make([]int16, 0)

	for _, v := range hull {
		scaledX, scaledY := match.TranslateScale(v.X, v.Y)
		scaledXInt := int16(scaledX) + int16(mapXOffset)
		scaledYInt := int16(scaledY) + int16(mapYOffset)
		xCoordinates = append(xCoordinates, scaledXInt)
		yCoordinates = append(yCoordinates, scaledYInt)
	}
	if common.MapHasAlternateVersion(match.MapName) && (!inferno.IsOnNormalElevation && isOnNormalElevation) ||
		(inferno.IsOnNormalElevation && !isOnNormalElevation) {
		color.A /= 2
	}

	gfx.FilledPolygonColor(renderer, xCoordinates, yCoordinates, color)
	gfx.AAPolygonColor(renderer, xCoordinates, yCoordinates, color)
}

func drawBomb(renderer *sdl.Renderer, bomb *common.Bomb, match *match.Match) {
	pos := bomb.Position
	if bomb.IsBeingCarried {
		return
	}

	scaledX, scaledY := match.TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	color := colorBomb
	if common.MapHasAlternateVersion(match.MapName) && (!bomb.IsOnNormalElevation && isOnNormalElevation) ||
		(bomb.IsOnNormalElevation && !isOnNormalElevation) {
		color.A = 100
	}

	gfx.BoxColor(renderer, scaledXInt-3, scaledYInt-2, scaledXInt+3, scaledYInt+2, color)
}

func drawString(renderer *sdl.Renderer, text string, color sdl.Color, x, y int32, font *ttf.Font) {
	textSurface, err := font.RenderUTF8Blended(text, color)
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

func drawInfobars(renderer *sdl.Renderer, match *match.Match, font *ttf.Font) {
	var cts, ts []common.Player
	for _, player := range match.States[curFrame].Players {
		if player.Team == demoinfo.TeamCounterTerrorists {
			cts = append(cts, player)

		} else {
			ts = append(ts, player)
		}
	}
	sort.Slice(cts, func(i, j int) bool { return cts[i].ID < cts[j].ID })
	sort.Slice(ts, func(i, j int) bool { return ts[i].ID < ts[j].ID })
	drawInfobar(renderer, cts, 0, mapYOffset, colorCounter, font)
	drawInfobar(renderer, ts, mapXOffset+mapOverviewWidth, mapYOffset, colorTerror, font)
	drawKillfeed(renderer, match.Killfeed[curFrame], mapXOffset+mapOverviewWidth, mapYOffset+600, font)
	drawTimer(renderer, match.States[curFrame].Timer, 0, mapYOffset+600, font)
}

func drawInfobar(renderer *sdl.Renderer, players []common.Player, x, y int32, color sdl.Color, font *ttf.Font) {
	var yOffset int32
	for i, player := range players {
		if player.IsAlive {
			gfx.BoxColor(renderer, x+int32(player.Health)*(mapXOffset/infobarElementHeight), yOffset, x, yOffset+5, color)
		}
		if !player.IsAlive {
			color.A = 150
		}
		name := fmt.Sprintf("%v %v", i+1, player.Name)
		drawString(renderer, cropStringToN(name, 20), color, x+85, yOffset+10, font)
		color.A = 255
		drawString(renderer, fmt.Sprintf("%v", player.Health), color, x+5, yOffset+10, font)
		if player.Armor > 0 && player.HasHelmet {
			drawString(renderer, "H", color, x+35, yOffset+10, font)
		} else if player.Armor > 0 {
			drawString(renderer, "A", color, x+35, yOffset+10, font)
		}
		if player.HasDefuseKit {
			drawString(renderer, "D", color, x+50, yOffset+10, font)
		}
		drawString(renderer, fmt.Sprintf("%v $", player.Money), colorMoney, x+5, yOffset+25, font)
		var nadeCounter int32
		inventory := player.Inventory
		for _, w := range inventory {
			if w.Class() == demoinfo.EqClassSMG || w.Class() == demoinfo.EqClassHeavy || w.Class() == demoinfo.EqClassRifle {
				drawString(renderer, w.String(), color, x+150, yOffset+25, font)
			}
			if w.Class() == demoinfo.EqClassPistols {
				drawString(renderer, w.String(), color, x+150, yOffset+40, font)
			}
			if w.Class() == demoinfo.EqClassGrenade {
				var nadeColor sdl.Color
				switch w {
				case demoinfo.EqDecoy:
					nadeColor = colorEqDecoy
				case demoinfo.EqMolotov:
					nadeColor = colorEqMolotov
				case demoinfo.EqIncendiary:
					nadeColor = colorEqIncendiary
				case demoinfo.EqFlash:
					nadeColor = colorEqFlash
				case demoinfo.EqSmoke:
					nadeColor = colorEqSmoke
				case demoinfo.EqHE:
					nadeColor = colorEqHE
				}

				gfx.BoxColor(renderer, x+150+nadeCounter*12, yOffset+60, x+150+nadeCounter*12+6, yOffset+60+9, nadeColor)
				nadeCounter++
			}
			if player.HasBomb {
				gfx.BoxColor(renderer, x+50, yOffset+12, x+45+12, yOffset+12+9, colorBomb)
			}
		}
		kdaInfo := fmt.Sprintf("%v / %v / %v", player.Kills, player.Assists, player.Deaths)
		drawString(renderer, kdaInfo, color, x+5, yOffset+40, font)

		yOffset += infobarElementHeight
	}
}

func drawKillfeed(renderer *sdl.Renderer, killfeed []common.Kill, x, y int32, font *ttf.Font) {
	var yOffset int32
	for _, kill := range killfeed {
		var colorKiller, colorVictim sdl.Color
		if kill.KillerTeam == demoinfo.TeamCounterTerrorists {
			colorKiller = colorCounter
		} else if kill.KillerTeam == demoinfo.TeamTerrorists {
			colorKiller = colorTerror
		} else {
			colorKiller = colorDarkWhite
		}
		if kill.VictimTeam == demoinfo.TeamCounterTerrorists {
			colorVictim = colorCounter
		} else {
			colorVictim = colorTerror
		}
		killerName := cropStringToN(kill.KillerName, 10)
		victimName := cropStringToN(kill.VictimName, 10)
		weaponName := cropStringToN(kill.Weapon.String(), 10)
		drawString(renderer, killerName, colorKiller, x+5, y+yOffset, font)
		drawString(renderer, weaponName, colorDarkWhite, x+110, y+yOffset, font)
		drawString(renderer, victimName, colorVictim, x+200, y+yOffset, font)
		yOffset += killfeedHeight
	}
}

func drawTimer(renderer *sdl.Renderer, timer common.Timer, x, y int32, font *ttf.Font) {
	if timer.Phase == common.PhaseWarmup {
		drawString(renderer, "Warmup", colorDarkWhite, x+5, y, font)
	} else {
		minutes := int(timer.TimeRemaining.Minutes())
		seconds := int(timer.TimeRemaining.Seconds()) - 60*minutes
		timeString := fmt.Sprintf("%d:%2d", minutes, seconds)
		/* ESEA demos have no RoundFreezetimeEnd events so in what
		should be PhaseRegular the minutes and seconds are negative.
		if minutes < 0 || seconds < 0 {
			timeString = "Paused / Timeout"
		}
		*/
		var color sdl.Color
		if timer.Phase == common.PhasePlanted {
			color = colorBomb
		} else if timer.Phase == common.PhaseRestart {
			color = colorEqHE
		} else {
			color = colorDarkWhite
		}
		drawString(renderer, timeString, color, x+5, y, font)
	}
}

func drawShot(renderer *sdl.Renderer, shot *common.Shot, match *match.Match) {
	pos := shot.Position
	viewAngleDegrees := -shot.ViewDirectionX // negated because of sdl
	viewAngleRadian := float64(viewAngleDegrees * math.Pi / 180)
	color := colorDarkWhite
	if shot.IsAwpShot {
		color = colorAwpShot
	}
	if common.MapHasAlternateVersion(match.MapName) && (!shot.IsOnNormalElevation && isOnNormalElevation) ||
		(shot.IsOnNormalElevation && !isOnNormalElevation) {
		color.A = 100
	}

	scaledX, scaledY := match.TranslateScale(pos.X, pos.Y)
	scaledX += float32(math.Cos(viewAngleRadian) * radiusPlayerFloat)
	scaledY += float32(math.Sin(viewAngleRadian) * radiusPlayerFloat)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset

	targetX := int32(scaledXInt) + int32(math.Cos(viewAngleRadian)*shotLength/float64(match.MapScale))
	targetY := int32(scaledYInt) + int32(math.Sin(viewAngleRadian)*shotLength/float64(match.MapScale))

	gfx.AALineColor(renderer, scaledXInt, scaledYInt, targetX, targetY, color)
}

func cropStringToN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	return s
}
