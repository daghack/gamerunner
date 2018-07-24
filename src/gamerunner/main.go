package main

import (
	"fmt"
	"gamerunner/entity"
	"gamerunner/entity/area"
	"gamerunner/eventrouter"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/yuin/gopher-lua"
	"image"
	_ "image/png"
	"os"
	"time"
)

var (
	index      int
	tilesImage *ebiten.Image
	walktick   <-chan time.Time
	roundtick  <-chan time.Time
	tiletick   <-chan time.Time
	walkleft   bool
	areatest   *area.Area
)

func init() {
	lwf, err := os.Open("resources/images/link_walking.png")
	if err != nil {
		panic(err)
	}
	defer lwf.Close()
	img, _, err := image.Decode(lwf)
	if err != nil {
		panic(err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	walktick = time.Tick(400 * time.Millisecond)
	roundtick = time.Tick(1600 * time.Millisecond)
	tiletick = time.Tick(350 * time.Millisecond)
	areatest, err = area.NewArea("test_area", "test_area.lua")
	if err != nil {
		panic(err)
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	select {
	case <-walktick:
		walkleft = !walkleft
	case <-tiletick:
		err := areatest.UpdateTiles()
		if err != nil {
			panic(err)
		}
	case <-roundtick:
		index = (index + 1) % 4
	default:
	}
	imageMax := tilesImage.Bounds().Max
	tilewidth := imageMax.X / 2
	tileheight := imageMax.Y / 4
	op := &ebiten.DrawImageOptions{}
	if walkleft {
		r := image.Rect(0, index*tileheight, tilewidth, (index+1)*tileheight)
		op.SourceRect = &r
	} else {
		r := image.Rect(tilewidth, index*tileheight, imageMax.X, (index+1)*tileheight)
		op.SourceRect = &r
	}
	areatest.Draw(screen)
	op.GeoM.Translate(64, 64)
	err := screen.DrawImage(tilesImage, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %0.2f", ebiten.CurrentFPS()))
	return err
}

func main() {
	//router := eventrouter.NewRouter("test_router.lua")
	space, err := eventrouter.NewRouter("world_manager.lua")
	if err != nil {
		panic(err)
	}
	e := entity.NewEntity(space)
	l := lua.NewState()
	defer l.Close()
	l.OpenLibs()
	space.AddToLua("space", l)
	//router.AddToLua("Router", l)
	go e.EventLoop()
	err = ebiten.Run(update, 256, 256, 2, "Tiles Test")
	if err != nil {
		panic(err)
	}
}
