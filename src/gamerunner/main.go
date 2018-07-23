package main

import (
	"fmt"
	"gamerunner/entity"
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
	tilesImage *ebiten.Image
	walktick   <-chan time.Time
	walkleft   bool
)

func init() {
	lwf, err := os.Open("resources/images/link_walking.png")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(lwf)
	if err != nil {
		panic(err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	walktick = time.Tick(500 * time.Millisecond)
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	select {
	case <-walktick:
		walkleft = !walkleft
	default:
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %0.2f", ebiten.CurrentFPS()))
	imageMax := tilesImage.Bounds().Max
	tilewidth := imageMax.X / 2
	tileheight := imageMax.Y / 3
	op := &ebiten.DrawImageOptions{}
	if walkleft {
		r := image.Rect(0, 0, tilewidth, tileheight)
		op.SourceRect = &r
	} else {
		r := image.Rect(tilewidth, 0, imageMax.X, tileheight)
		op.SourceRect = &r
	}
	op.GeoM.Scale(512.0/float64(tilewidth), 512.0/float64(tileheight))
	err := screen.DrawImage(tilesImage, op)
	return err
}

func main() {
	//router := eventrouter.NewRouter("test_router.lua")
	space, err := eventrouter.NewRouter("space.lua")
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
	err = ebiten.Run(update, 512, 512, 2, "Tiles Test")
	if err != nil {
		panic(err)
	}
}
