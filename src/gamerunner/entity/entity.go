package entity

import (
	"gamerunner/controllers"
	"github.com/hajimehoshi/ebiten"
	"github.com/yuin/gopher-lua"
	"image"
	_ "image/png"
	"os"
	"time"
)

type Entity struct {
	state      *lua.LState
	tileset    *ebiten.Image
	controller controllers.Controller
}

func NewEntity(entityfile string, controller controllers.Controller) (*Entity, error) {
	l := lua.NewState()
	toret := &Entity{
		state:      l,
		controller: controller,
	}
	err := l.DoFile(entityfile)
	if err != nil {
		return nil, err
	}
	l.SetGlobal("controller", lua.LChannel(controller.CommandChannel()))
	entityTable := l.GetGlobal("entity").(*lua.LTable)
	tilesetFile := string(entityTable.RawGetString("tileset").(lua.LString))
	lwf, err := os.Open(tilesetFile)
	if err != nil {
		return nil, err
	}
	defer lwf.Close()
	img, _, err := image.Decode(lwf)
	if err != nil {
		return nil, err
	}
	toret.tileset, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	return toret, nil
}

func (e *Entity) Draw(screen *ebiten.Image, tilesize float64) error {
	timestamp := lua.LNumber(float64(time.Now().UnixNano() / 1000000))
	ups := lua.LNumber(ebiten.CurrentTPS())
	err := e.state.CallByParam(lua.P{
		Fn:      e.state.GetGlobal("update_state"),
		Protect: true,
	}, timestamp, ups)
	if err != nil {
		return err
	}
	err = e.state.CallByParam(lua.P{
		Fn:      e.state.GetGlobal("active_frame"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return err
	}
	frame := e.state.ToInt(-1)
	e.state.Pop(1)
	imageMax := e.tileset.Bounds().Max
	tilewidth := imageMax.X / 2
	tileheight := imageMax.Y / 4
	x := frame % 2
	y := frame / 2
	r := image.Rect(x*tilewidth, y*tileheight, (x+1)*tilewidth, (y+1)*tileheight)
	op := &ebiten.DrawImageOptions{
		SourceRect: &r,
	}
	state := e.state.GetGlobal("state").(*lua.LTable)
	currentX := float64(state.RawGetString("x_pos").(lua.LNumber)) * tilesize
	currentY := float64(state.RawGetString("y_pos").(lua.LNumber)) * tilesize
	op.GeoM.Translate(currentX, currentY)
	screen.DrawImage(e.tileset, op)

	return nil
}

func (e *Entity) Controller() controllers.Controller {
	return e.controller
}
