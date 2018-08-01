package main

import (
	"gamerunner/world"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	w, err := world.LoadWorld("world_manager.lua")
	if err != nil {
		panic(err)
	}
	err = w.LoadEntity("link", "link.lua", w.KeyboardController())
	if err != nil {
		panic(err)
	}
	err = ebiten.Run(w.Run, 256, 256, 2, "World Manager")
	if err != nil {
		panic(err)
	}
}
