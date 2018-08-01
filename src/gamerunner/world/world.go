package world

import (
	"gamerunner/area"
	"gamerunner/eventrouter"
	"github.com/hajimehoshi/ebiten"
	"github.com/yuin/gopher-lua"
	"image"
	_ "image/png"
	"os"
	"time"
)

type World struct {
	state               *lua.LState
	manager             *eventrouter.Router
	areaMap             map[string]*area.Area
	linkTS              *ebiten.Image
	keyboardControllers []chan lua.LValue
	keyboardstate       map[ebiten.Key]bool
}

func LoadWorld(worldfile string) (*World, error) {
	l := lua.NewState()
	err := l.DoFile(worldfile)
	if err != nil {
		return nil, err
	}
	router, err := eventrouter.NewRouterFromState(l)
	if err != nil {
		return nil, err
	}
	toret := &World{
		state:         l,
		manager:       router,
		areaMap:       map[string]*area.Area{},
		keyboardstate: map[ebiten.Key]bool{},
	}
	err = toret.loadAreas()
	if err != nil {
		return nil, err
	}
	return toret, nil
}

func (w *World) loadArea(key, areafile string) error {
	area, err := area.NewArea(key, areafile)
	if err != nil {
		return err
	}
	w.areaMap[key] = area
	return nil
}

func (w *World) loadAreas() error {
	var err error
	lvtable := w.state.GetGlobal("areas").(*lua.LTable)
	lvtable.ForEach(func(idval lua.LValue, fileval lua.LValue) {
		id := idval.(lua.LString)
		file := fileval.(lua.LString)
		if err != nil {
			return
		}
		err = w.loadArea(string(id), string(file))
	})
	return err
}

func (w *World) LoadEntity(id, entityfile string, controller chan lua.LValue) error {
	// #TODO THIS IS NOT HOW I WANT TO DO THIS: FOR INITIAL TESTING ONLY
	err := w.state.DoFile(entityfile)
	if err != nil {
		return err
	}
	w.state.SetGlobal("controller", lua.LChannel(controller))
	lwf, err := os.Open("resources/images/link_walking.png")
	if err != nil {
		return err
	}
	defer lwf.Close()
	img, _, err := image.Decode(lwf)
	if err != nil {
		return err
	}
	w.linkTS, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	return err
}

func (w *World) drawEntity(screen *ebiten.Image) error {
	// #TODO THIS IS NOT HOW I WANT TO DO THIS: FOR INITIAL TESTING ONLY
	timestamp := lua.LNumber(float64(time.Now().UnixNano() / 1000000))
	err := w.state.CallByParam(lua.P{
		Fn:      w.state.GetGlobal("update_state"),
		Protect: true,
	}, timestamp)
	if err != nil {
		return err
	}
	err = w.state.CallByParam(lua.P{
		Fn:      w.state.GetGlobal("active_frame"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return err
	}
	frame := w.state.ToInt(-1)
	w.state.Pop(1)
	imageMax := w.linkTS.Bounds().Max
	tilewidth := imageMax.X / 2
	tileheight := imageMax.Y / 4
	x := frame % 2
	y := frame / 2
	r := image.Rect(x*tilewidth, y*tileheight, (x+1)*tilewidth, (y+1)*tileheight)
	op := &ebiten.DrawImageOptions{
		SourceRect: &r,
	}
	op.GeoM.Translate(64, 64)
	screen.DrawImage(w.linkTS, op)

	return nil
}

func (w *World) Run(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	for i := ebiten.Key0; i <= ebiten.KeyMax; i++ {
		p := w.keyboardstate[i]
		keypressed := ebiten.IsKeyPressed(i)
		if p != keypressed {
			w.keyboardstate[i] = keypressed
			table := w.state.NewTable()
			table.RawSetString("key", lua.LNumber(i))
			table.RawSetString("pressed", lua.LBool(keypressed))
			for _, ch := range w.keyboardControllers {
				select {
				case ch <- table:
				default:
				}
			}
		}
	}
	w.areaMap["test_area"].Draw(screen)
	return w.drawEntity(screen)
}

func (w *World) KeyboardController() chan lua.LValue {
	newChannel := make(chan lua.LValue, 16)
	w.keyboardControllers = append(w.keyboardControllers, newChannel)
	return newChannel
}
