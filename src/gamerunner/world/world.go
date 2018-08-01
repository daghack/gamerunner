package world

import (
	"gamerunner/area"
	"gamerunner/entity"
	"gamerunner/eventrouter"
	"github.com/hajimehoshi/ebiten"
	"github.com/yuin/gopher-lua"
)

type World struct {
	state               *lua.LState
	manager             *eventrouter.Router
	areaMap             map[string]*area.Area
	entityMap           map[string]*entity.Entity
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
		entityMap:     map[string]*entity.Entity{},
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
	var err error
	w.entityMap[id], err = entity.NewEntity(entityfile, controller)
	return err
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
	for _, entity := range w.entityMap {
		err := entity.Draw(screen)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) KeyboardController() chan lua.LValue {
	newChannel := make(chan lua.LValue, 16)
	w.keyboardControllers = append(w.keyboardControllers, newChannel)
	return newChannel
}
