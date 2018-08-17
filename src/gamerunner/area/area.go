package area

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/yuin/gopher-lua"
	"image"
	_ "image/png"
	"os"
	"time"
)

type Area struct {
	state       *lua.LState
	renderedMap *ebiten.Image
}

func NewArea(areaId, luafile string) (*Area, error) {
	l := lua.NewState()
	err := l.DoFile(luafile)
	if err != nil {
		return nil, err
	}
	toret := &Area{
		state: l,
	}
	err = toret.loadArea()
	if err != nil {
		return nil, err
	}
	l.SetGlobal("reload_map", l.NewFunction(toret.ReloadMap))
	return toret, nil
}

func (area *Area) renderMap(tileset *ebiten.Image, tilesize, w, h int, layers *lua.LTable, undrawn_tiles_table *lua.LTable) error {
	var err error
	area.renderedMap, err = ebiten.NewImage(tilesize*w, tilesize*h, ebiten.FilterDefault)
	if err != nil {
		return err
	}
	undrawn_tiles := map[int]bool{}
	undrawn_tiles_table.ForEach(func(key, val lua.LValue) {
		undrawn_tiles[int(val.(lua.LNumber))] = true
	})
	layers.ForEach(func(key, val lua.LValue) {
		layer_tilemap := map[int]int{}
		layer := val.(*lua.LTable)
		layer.ForEach(func(key, value lua.LValue) {
			layer_tilemap[int(key.(lua.LNumber))-1] = int(value.(lua.LNumber))
		})
		area.renderLayer(tileset, tilesize, w, h, layer_tilemap, undrawn_tiles)
	})
	return nil
}

func (area *Area) renderLayer(tileset *ebiten.Image, tilesize, w, h int, layer map[int]int, undrawn_tiles map[int]bool) error {
	tileXNum := tileset.Bounds().Max.X / tilesize
	k := len(layer)
	for i := 0; i < k; i += 1 {
		t := layer[i]
		if undrawn_tiles[t] {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%w)*tilesize), float64((i/w)*tilesize))
		sx := (t % tileXNum) * tilesize
		sy := (t / tileXNum) * tilesize
		r := image.Rect(sx, sy, sx+tilesize, sy+tilesize)
		op.SourceRect = &r
		area.renderedMap.DrawImage(tileset, op)
	}
	return nil
}

func (area *Area) ReloadMap(state *lua.LState) int {
	area.loadArea()
	return 0
}

func (area *Area) loadArea() error {
	config := area.state.GetGlobal("map_config").(*lua.LTable)
	tilesetfile := config.RawGetString("tileset").(lua.LString)
	tilesize := config.RawGetString("tile_size").(lua.LNumber)
	mapwidth := config.RawGetString("map_width").(lua.LNumber)
	mapheight := config.RawGetString("map_height").(lua.LNumber)
	tilemapTable := config.RawGetString("tilemap").(*lua.LTable)
	layers := tilemapTable.RawGetString("layers").(*lua.LTable)
	undrawn_tiles := tilemapTable.RawGetString("undrawn_tiles").(*lua.LTable)
	tilesetR, err := os.Open(string(tilesetfile))
	if err != nil {
		return err
	}
	defer tilesetR.Close()
	img, _, err := image.Decode(tilesetR)
	if err != nil {
		return err
	}
	tileset, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		return err
	}
	return area.renderMap(tileset, int(tilesize), int(mapwidth), int(mapheight), layers, undrawn_tiles)
}

func (area *Area) TileSize() float64 {
	config := area.state.GetGlobal("map_config").(*lua.LTable)
	return float64(config.RawGetString("tile_size").(lua.LNumber))
}

func (area *Area) updateArea() error {
	timestamp := lua.LNumber(float64(time.Now().UnixNano() / 1000000))
	return area.state.CallByParam(lua.P{
		Fn:      area.state.GetGlobal("update_area"),
		Protect: true,
	}, timestamp)
}

func (area *Area) Draw(screen *ebiten.Image) error {
	area.updateArea()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(256.0/128.0, 256.0/128.0)
	return screen.DrawImage(area.renderedMap, op)
}

func (area *Area) UpdateTiles() error {
	return area.state.CallByParam(lua.P{Fn: area.state.GetGlobal("update_tiles")})
}
