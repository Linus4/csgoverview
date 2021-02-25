package main

import (
	"fmt"
	"log"
	"math"
	"sort"

	common "github.com/linus4/csgoverview/common"
	demoinfo "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	radiusPlayer      int32   = 10
	radiusPlayerFloat float64 = float64(radiusPlayer)
	radiusSmoke       float64 = 25
	killfeedHeight    int32   = 15
	shotLength        float64 = 1000
	headshotRune      rune    = '\u205c'
)

var (
	colorTerror                = sdl.Color{252, 176, 12, 255}
	colorCounter               = sdl.Color{89, 206, 200, 255}
	colorMoney                 = sdl.Color{45, 135, 45, 255}
	colorBomb                  = sdl.Color{255, 0, 0, 255}
	colorEqDecoy               = sdl.Color{102, 34, 0, 255}
	colorEqMolotov             = sdl.Color{255, 153, 0, 255}
	colorEqIncendiary          = sdl.Color{255, 153, 0, 255}
	colorInferno               = sdl.Color{255, 153, 0, 100}
	colorEqFlash               = sdl.Color{128, 170, 255, 255}
	colorEqSmoke               = sdl.Color{153, 153, 153, 255}
	colorEqSmokeOutlineTerror  = sdl.Color{252, 176, 12, 120}
	colorEqSmokeOutlineCounter = sdl.Color{89, 206, 200, 120}
	colorSmoke                 = sdl.Color{153, 153, 153, 100}
	colorEqHE                  = sdl.Color{85, 150, 0, 255}
	colorDarkWhite             = sdl.Color{200, 200, 200, 255}
	colorDarkGrey              = sdl.Color{125, 125, 125, 255}
	colorFlashEffect           = sdl.Color{200, 200, 200, 180}
	colorAwpShot               = sdl.Color{255, 50, 0, 255}
)

func (app *app) drawPlayer(player *common.Player, index int) {
	m := app.match
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

		scaledX, scaledY := m.TranslateScale(pos.X, pos.Y)
		var scaledXInt int32 = int32(scaledX) + mapXOffset
		var scaledYInt int32 = int32(scaledY) + mapYOffset

		if common.MapHasAlternateVersion(m.MapName) && (!player.IsOnNormalElevation && app.isOnNormalElevation) ||
			(player.IsOnNormalElevation && !app.isOnNormalElevation) {
			color.A = 100
			colorLOS.A = 100
			colorC4.A = 100
		}

		var name string
		number := index + 1
		if player.Team == demoinfo.TeamTerrorists {
			number = (number + 5) % 10
		}
		if !app.hidePlayerNames {
			name = fmt.Sprintf("%v %v", number, player.Name)
		} else {
			name = fmt.Sprintf("%v", number)
		}

		switch player.ActiveWeapon {
		case demoinfo.EqBomb:
			gfx.BoxColor(app.renderer, scaledXInt+radiusPlayer+4, scaledYInt-2, scaledXInt+radiusPlayer+10, scaledYInt+2, colorC4)

		case demoinfo.EqHE, demoinfo.EqFlash, demoinfo.EqSmoke, demoinfo.EqIncendiary, demoinfo.EqMolotov, demoinfo.EqDecoy:
			var colorGrenade sdl.Color
			switch player.ActiveWeapon {
			case demoinfo.EqDecoy:
				colorGrenade = colorEqDecoy
			case demoinfo.EqMolotov:
				colorGrenade = colorEqMolotov
			case demoinfo.EqIncendiary:
				colorGrenade = colorEqIncendiary
			case demoinfo.EqFlash:
				colorGrenade = colorEqFlash
			case demoinfo.EqSmoke:
				colorGrenade = colorEqSmoke
			case demoinfo.EqHE:
				colorGrenade = colorEqHE
			}
			if common.MapHasAlternateVersion(m.MapName) && (!player.IsOnNormalElevation && app.isOnNormalElevation) ||
				(player.IsOnNormalElevation && !app.isOnNormalElevation) {
				colorGrenade.A = 100
			}
			gfx.BoxColor(app.renderer, scaledXInt+radiusPlayer+5, scaledYInt-3, scaledXInt+radiusPlayer+9, scaledYInt+3, colorGrenade)
		}

		app.drawString(cropStringToN(name, 12), color, scaledXInt+10, scaledYInt+10)

		if player.Health == 100 {
			gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer, color)
		} else {
			// start == 0 is facing right
			// health left
			var healthArc int32 = int32(player.Health) * 360 / 100
			start := 90 - (healthArc / 2)
			end := 90 + (healthArc / 2)
			gfx.ArcColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer, start, end, color)
			// health lost
			color.R = uint8(float32(color.R) * 0.6)
			color.G = uint8(float32(color.G) * 0.6)
			color.B = uint8(float32(color.B) * 0.6)
			start = -90 - ((360 - healthArc) / 2)
			end = -90 + ((360 - healthArc) / 2)
			gfx.ArcColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer, start, end, color)
		}

		viewAngle := -int32(player.ViewDirectionX) // negated because of sdl
		if player.HasAwp() {
			colorLOS = colorAwpShot
		}
		gfx.ArcColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer+1, viewAngle-20, viewAngle+20, colorLOS)
		gfx.ArcColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer+2, viewAngle-10, viewAngle+10, colorLOS)
		gfx.ArcColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer+3, viewAngle-5, viewAngle+5, colorLOS)

		colorFlash := colorFlashEffect
		if player.FlashDuration.Seconds() > 0.5 {
			remaining := player.FlashTimeRemaining
			if remaining.Seconds() >= 3.1 {
				colorFlash.A = 255
			} else {
				colorFlash.A = uint8((remaining.Seconds() * 255) / 3.1)
			}
			if common.MapHasAlternateVersion(m.MapName) && (!player.IsOnNormalElevation && app.isOnNormalElevation) ||
				(player.IsOnNormalElevation && !app.isOnNormalElevation) {
				colorFlash.A /= 2
			}
			gfx.FilledCircleColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer-5, colorFlash)
		}

		if player.HasBomb {
			gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer-1, colorC4)
			gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, radiusPlayer-2, colorC4)
		}

		if player.IsDefusing {
			color.A -= 55
			gfx.CharacterColor(app.renderer, scaledXInt-radiusPlayer/4, scaledYInt-radiusPlayer/4, 'D', color)
			color.A += 55
		}
	} else {
		pos := player.LastAlivePosition

		scaledX, scaledY := m.TranslateScale(pos.X, pos.Y)
		var scaledXInt int32 = int32(scaledX) + mapXOffset
		var scaledYInt int32 = int32(scaledY) + mapYOffset

		color.A -= 105
		gfx.CharacterColor(app.renderer, scaledXInt, scaledYInt, 'X', color)
	}
}

func (app *app) drawGrenade(grenade *common.GrenadeProjectile) {
	pos := grenade.Position
	m := app.match

	scaledX, scaledY := m.TranslateScale(pos.X, pos.Y)
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
	if common.MapHasAlternateVersion(m.MapName) && (!grenade.IsOnNormalElevation && app.isOnNormalElevation) ||
		(grenade.IsOnNormalElevation && !app.isOnNormalElevation) {
		color.A = 100
	}

	gfx.BoxColor(app.renderer, scaledXInt-2, scaledYInt-3, scaledXInt+2, scaledYInt+3, color)
}

func (app *app) drawEffects(effect *common.Effect) {
	m := app.match
	pos := effect.Position

	scaledX, scaledY := m.TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	var alphaModifier uint8
	if common.MapHasAlternateVersion(m.MapName) && (!effect.IsOnNormalElevation && app.isOnNormalElevation) ||
		(effect.IsOnNormalElevation && !app.isOnNormalElevation) {
		alphaModifier = 2
	} else {
		alphaModifier = 1
	}

	switch effect.Type {
	case demoinfo.EqFlash:
		color := colorEqFlash
		color.A /= alphaModifier
		gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, effect.Lifetime, color)
	case demoinfo.EqHE:
		color := colorEqHE
		color.A /= alphaModifier
		gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, effect.Lifetime, color)
	case demoinfo.EqSmoke:
		color := colorSmoke
		color.A /= alphaModifier
		colorCircles := colorDarkWhite
		colorCircles.A /= alphaModifier
		// 4.9 is the reference on Inferno for the value for radiusSmoke
		scaledRadiusSmoke := int32(radiusSmoke * 4.9 / float64(m.MapScale))
		gfx.FilledCircleColor(app.renderer, scaledXInt, scaledYInt, scaledRadiusSmoke, color)
		// only draw the outline if the smoke is not fading
		if effect.Lifetime < 15*m.SmokeEffectLifetime/18 {
			var smokeOutlineColor = colorEqSmokeOutlineTerror
			if effect.Team == demoinfo.TeamCounterTerrorists {
				smokeOutlineColor = colorEqSmokeOutlineCounter
			}
			gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, scaledRadiusSmoke, smokeOutlineColor)
		}
		gfx.ArcColor(app.renderer, scaledXInt, scaledYInt, 10, 270+effect.Lifetime*360/m.SmokeEffectLifetime, 630, colorCircles)
	case demoinfo.EqDefuseKit:
		gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, effect.Lifetime, colorMoney)
	case demoinfo.EqBomb:
		gfx.AACircleColor(app.renderer, scaledXInt, scaledYInt, effect.Lifetime, colorBomb)
	}
}

func (app *app) drawInferno(inferno *common.Inferno) {
	m := app.match
	hull := inferno.ConvexHull2D
	color := colorInferno
	xCoordinates := make([]int16, 0)
	yCoordinates := make([]int16, 0)

	for _, v := range hull {
		scaledX, scaledY := m.TranslateScale(v.X, v.Y)
		scaledXInt := int16(scaledX) + int16(mapXOffset)
		scaledYInt := int16(scaledY) + int16(mapYOffset)
		xCoordinates = append(xCoordinates, scaledXInt)
		yCoordinates = append(yCoordinates, scaledYInt)
	}
	if common.MapHasAlternateVersion(m.MapName) && (!inferno.IsOnNormalElevation && app.isOnNormalElevation) ||
		(inferno.IsOnNormalElevation && !app.isOnNormalElevation) {
		color.A /= 2
	}

	gfx.FilledPolygonColor(app.renderer, xCoordinates, yCoordinates, color)
	gfx.AAPolygonColor(app.renderer, xCoordinates, yCoordinates, color)
}

func (app *app) drawBomb(bomb *common.Bomb) {
	m := app.match
	pos := bomb.Position
	if bomb.IsBeingCarried {
		return
	}

	scaledX, scaledY := m.TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset
	color := colorBomb
	if common.MapHasAlternateVersion(m.MapName) && (!bomb.IsOnNormalElevation && app.isOnNormalElevation) ||
		(bomb.IsOnNormalElevation && !app.isOnNormalElevation) {
		color.A = 100
	}

	gfx.BoxColor(app.renderer, scaledXInt-3, scaledYInt-2, scaledXInt+3, scaledYInt+2, color)
}

func (app *app) drawString(text string, color sdl.Color, x, y int32) {
	textSurface, err := app.font.RenderUTF8Blended(text, color)
	if err != nil {
		log.Fatal(err)
	}
	defer textSurface.Free()
	textTexture, err := app.renderer.CreateTextureFromSurface(textSurface)
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
	err = app.renderer.Copy(textTexture, nil, textRect)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *app) drawInfobars() {
	m := app.match
	var cts, ts []common.Player
	for _, player := range m.States[app.curFrame].Players {
		if player.Team == demoinfo.TeamCounterTerrorists {
			cts = append(cts, player)

		} else {
			ts = append(ts, player)
		}
	}
	sort.Slice(cts, func(i, j int) bool { return cts[i].ID < cts[j].ID })
	sort.Slice(ts, func(i, j int) bool { return ts[i].ID < ts[j].ID })
	app.drawInfobar(cts, 0, mapYOffset, colorCounter)
	app.drawInfobar(ts, mapXOffset+mapOverviewWidth, mapYOffset, colorTerror)
	app.drawKillfeed(m.Killfeed[app.curFrame], mapXOffset+mapOverviewWidth, mapYOffset+600)
	app.drawTimer(m.States[app.curFrame].Timer, 0, mapYOffset+600)
	app.drawPlaybackSpeedModifier(5, mapYOffset+630)
}

func (app *app) drawInfobar(players []common.Player, x, y int32, color sdl.Color) {
	var yOffset int32
	for i, player := range players {
		if player.IsAlive {
			gfx.BoxColor(app.renderer, x+int32(player.Health)*(mapXOffset/infobarElementHeight), yOffset, x, yOffset+5, color)
		}
		if !player.IsAlive {
			color.A = 150
		}
		number := i + 1
		if player.Team == demoinfo.TeamTerrorists {
			number = (number + 5) % 10
		}
		name := fmt.Sprintf("%v %v", number, player.Name)
		app.drawString(cropStringToN(name, 20), color, x+85, yOffset+10)
		color.A = 255
		app.drawString(fmt.Sprintf("%v", player.Health), color, x+5, yOffset+10)
		if player.Armor > 0 && player.HasHelmet {
			app.drawString("H", color, x+35, yOffset+10)
		} else if player.Armor > 0 {
			app.drawString("A", color, x+35, yOffset+10)
		}
		if player.HasDefuseKit {
			app.drawString("D", color, x+50, yOffset+10)
		}
		app.drawString(fmt.Sprintf("%v $", player.Money), colorMoney, x+5, yOffset+25)
		var nadeCounter int32
		inventory := player.Inventory
		for _, w := range inventory {
			if w.Class() == demoinfo.EqClassSMG || w.Class() == demoinfo.EqClassHeavy || w.Class() == demoinfo.EqClassRifle {
				app.drawString(w.String(), color, x+150, yOffset+25)
			}
			if w.Class() == demoinfo.EqClassPistols {
				app.drawString(w.String(), color, x+150, yOffset+40)
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

				gfx.BoxColor(app.renderer, x+150+nadeCounter*12, yOffset+60, x+150+nadeCounter*12+6, yOffset+60+9, nadeColor)
				nadeCounter++
			}
			if player.HasBomb {
				gfx.BoxColor(app.renderer, x+50, yOffset+12, x+45+12, yOffset+12+9, colorBomb)
			}
		}
		kdaInfo := fmt.Sprintf("%v / %v / %v", player.Kills, player.Assists, player.Deaths)
		app.drawString(kdaInfo, color, x+5, yOffset+40)

		yOffset += infobarElementHeight
	}
}

func (app *app) drawKillfeed(killfeed []common.Kill, x, y int32) {
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
		app.drawString(killerName, colorKiller, x+5, y+yOffset)
		if kill.Headshot {
			weaponName = weaponName + " " + string(headshotRune)
		}
		app.drawString(weaponName, colorDarkWhite, x+112, y+yOffset)
		app.drawString(victimName, colorVictim, x+222, y+yOffset)
		yOffset += killfeedHeight
	}
}

func (app *app) drawTimer(timer common.Timer, x, y int32) {
	if timer.Phase == common.PhaseWarmup {
		app.drawString("Warmup", colorDarkWhite, x+5, y)
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
		app.drawString(timeString, color, x+5, y)
	}
}

func (app *app) drawShot(shot *common.Shot) {
	m := app.match
	pos := shot.Position
	viewAngleDegrees := -shot.ViewDirectionX // negated because of sdl
	viewAngleRadian := float64(viewAngleDegrees * math.Pi / 180)
	color := colorDarkWhite
	if shot.IsAwpShot {
		color = colorAwpShot
	}
	if common.MapHasAlternateVersion(m.MapName) && (!shot.IsOnNormalElevation && app.isOnNormalElevation) ||
		(shot.IsOnNormalElevation && !app.isOnNormalElevation) {
		color.A = 100
	}

	scaledX, scaledY := m.TranslateScale(pos.X, pos.Y)
	scaledX += float32(math.Cos(viewAngleRadian) * radiusPlayerFloat)
	scaledY += float32(math.Sin(viewAngleRadian) * radiusPlayerFloat)
	var scaledXInt int32 = int32(scaledX) + mapXOffset
	var scaledYInt int32 = int32(scaledY) + mapYOffset

	targetX := int32(scaledXInt) + int32(math.Cos(viewAngleRadian)*shotLength/float64(m.MapScale))
	targetY := int32(scaledYInt) + int32(math.Sin(viewAngleRadian)*shotLength/float64(m.MapScale))

	gfx.AALineColor(app.renderer, scaledXInt, scaledYInt, targetX, targetY, color)
}

func (app *app) drawPlaybackSpeedModifier(x, y int32) {
	str := fmt.Sprintf("Playback-Speed: x %v", app.staticPlaybackSpeedModifier*app.playbackSpeedModifier)
	app.drawString(str, colorDarkWhite, x, y)
}

func cropStringToN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	return s
}
