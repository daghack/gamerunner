package controllers

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/yuin/gopher-lua"
)

var keysListenedFor = []ebiten.Key{
	ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD,
	ebiten.KeyUp, ebiten.KeyLeft, ebiten.KeyDown, ebiten.KeyRight,
}

const (
	CmdMoveUp    = "move_up"
	CmdMoveDown  = "move_down"
	CmdMoveLeft  = "move_left"
	CmdMoveRight = "move_right"
)

type KeyboardController struct {
	inverted       bool
	commandChannel chan lua.LValue
}

func NewKeyboardController(inverted bool) *KeyboardController {
	return &KeyboardController{
		commandChannel: make(chan lua.LValue, 16),
		inverted:       inverted,
	}
}

func (kc *KeyboardController) Frame() {
	for _, key := range keysListenedFor {
		command_str := ""
		switch key {
		case ebiten.KeyW, ebiten.KeyUp:
			command_str = CmdMoveUp
		case ebiten.KeyA, ebiten.KeyLeft:
			command_str = CmdMoveLeft
		case ebiten.KeyS, ebiten.KeyDown:
			command_str = CmdMoveDown
		case ebiten.KeyD, ebiten.KeyRight:
			command_str = CmdMoveRight
		}
		if kc.inverted {
			switch command_str {
			case CmdMoveUp:
				command_str = CmdMoveDown
			case CmdMoveDown:
				command_str = CmdMoveUp
			case CmdMoveLeft:
				command_str = CmdMoveRight
			case CmdMoveRight:
				command_str = CmdMoveLeft
			}
		}
		if ebiten.IsKeyPressed(key) && command_str != "" {
			kc.commandChannel <- lua.LString(command_str)
		}
	}
}

func (kc *KeyboardController) CommandChannel() chan lua.LValue {
	return kc.commandChannel
}
