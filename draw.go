package main

import (
	common "github.com/markus-wa/demoinfocs-golang/common"
	meta "github.com/markus-wa/demoinfocs-golang/metadata"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
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

func DrawPlayer(renderer *sdl.Renderer, player *OverviewPlayer, mapName string) {
	pos := player.LastAlivePosition

	scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX)
	var scaledYInt int32 = int32(scaledY)
	var colorR, colorG, colorB uint8

	if player.Team == common.TeamTerrorists {
		colorR = terrorR
		colorG = terrorG
		colorB = terrorB
	} else { // if player.Team == common.TeamCounterTerrorists {
		colorR = counterR
		colorG = counterG
		colorB = counterB
	}

	if player.Hp > 0 {
		gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, radiusPlayer, colorR, colorG, colorB, 255)
		gfx.StringRGBA(renderer, scaledXInt+15, scaledYInt+15, player.Name, colorR, colorG, colorB, 255)

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
			gfx.CharacterRGBA(renderer, scaledXInt-radiusPlayer/4, scaledYInt-radiusPlayer/4, 'D', counterR, counterG, counterB, 200)
		}
	} else {
		//gfx.SetFont(fontdata, 10, 10)
		gfx.CharacterRGBA(renderer, scaledXInt, scaledYInt, 'X', colorR, colorG, colorB, 150)
	}
}

func DrawGrenade(renderer *sdl.Renderer, grenade *common.GrenadeProjectile, mapName string) {
	pos := grenade.Position

	scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX)
	var scaledYInt int32 = int32(scaledY)
	var colorR, colorG, colorB uint8

	switch grenade.Weapon {
	case common.EqDecoy:
		colorR = 102
		colorG = 34
		colorB = 0
	case common.EqMolotov:
		colorR = 255
		colorG = 153
		colorB = 0
	case common.EqIncendiary:
		colorR = 255
		colorG = 153
		colorB = 0
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

	gfx.BoxRGBA(renderer, scaledXInt-2, scaledYInt-3, scaledXInt+2, scaledYInt+3, colorR, colorG, colorB, 255)

	// SmokeStart InfernoStart InfernoExpired
}

func DrawGrenadeEffect(renderer *sdl.Renderer, effect *GrenadeEffect, mapName string) {
	pos := effect.GrenadeEvent.Position

	scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX)
	var scaledYInt int32 = int32(scaledY)
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
		if effect.Lifetime < 15*smokeEffectLifetime/18 {
			gfx.CircleRGBA(renderer, scaledXInt, scaledYInt, 25, colorR, colorG, colorB, 255)
		}
		gfx.ArcRGBA(renderer, scaledXInt, scaledYInt, 10, int32(270+effect.Lifetime*360/smokeEffectLifetime), 630, colorR, colorG, colorB, 255)
	}
}

func DrawInferno(renderer *sdl.Renderer, inferno *common.Inferno, mapName string) {
	hull := inferno.ConvexHull2D()
	var colorR, colorG, colorB uint8 = 255, 153, 0
	xCoordinates := make([]int16, 0)
	yCoordinates := make([]int16, 0)

	for _, v := range hull {
		scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(v.X, v.Y)
		scaledXInt := int16(scaledX)
		scaledYInt := int16(scaledY)
		xCoordinates = append(xCoordinates, scaledXInt)
		yCoordinates = append(yCoordinates, scaledYInt)
	}

	gfx.FilledPolygonRGBA(renderer, xCoordinates, yCoordinates, colorR, colorG, colorB, 100)
	gfx.PolygonRGBA(renderer, xCoordinates, yCoordinates, colorR, colorG, colorB, 100)
}

func DrawBomb(renderer *sdl.Renderer, bomb *common.Bomb, mapName string) {
	pos := bomb.Position()
	if bomb.Carrier != nil {
		return
	}

	scaledX, scaledY := meta.MapNameToMap[mapName].TranslateScale(pos.X, pos.Y)
	var scaledXInt int32 = int32(scaledX)
	var scaledYInt int32 = int32(scaledY)
	var colorR, colorG, colorB uint8

	colorR = 255
	colorG = 0
	colorB = 0

	gfx.BoxRGBA(renderer, scaledXInt-3, scaledYInt-2, scaledXInt+3, scaledYInt+2, colorR, colorG, colorB, 255)
}
