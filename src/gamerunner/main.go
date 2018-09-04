package main

import (
	"gamerunner/controllers"
	"gamerunner/world"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	ebiten.SetFullscreen(false)
	defer ebiten.SetFullscreen(false)
	w, err := world.LoadWorld("world_manager.lua")
	if err != nil {
		panic(err)
	}
	err = w.LoadEntity("link", "link.lua", controllers.NewKeyboardController(false))
	if err != nil {
		panic(err)
	}
	err = ebiten.Run(w.Run, 256, 256, 1, "World Manager")
	if err != nil {
		panic(err)
	}
}
