package world

import (
	"fmt"
	"gamerunner/area"
	"gamerunner/controllers"
	"gamerunner/entity"
	"gamerunner/eventrouter"
	"github.com/hajimehoshi/ebiten"
	"github.com/yuin/gopher-lua"
)

type World struct {
	state     *lua.LState
	manager   *eventrouter.Router
	areaMap   map[string]*area.Area
	entityMap map[string]*entity.Entity
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
		state:     l,
		manager:   router,
		areaMap:   map[string]*area.Area{},
		entityMap: map[string]*entity.Entity{},
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

func (w *World) LoadEntity(id, entityfile string, controller controllers.Controller) error {
	var err error
	w.entityMap[id], err = entity.NewEntity(entityfile, controller)
	return err
}

func (w *World) Run(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("Esc Pressed")
	}
	for _, entity := range w.entityMap {
		entity.Controller().Frame()
	}
	w.areaMap["test_area"].Draw(screen)
	for _, entity := range w.entityMap {
		err := entity.Draw(screen, w.areaMap["test_area"].TileSize())
		if err != nil {
			return err
		}
	}
	return nil
}
