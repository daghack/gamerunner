package controllers

import (
	"github.com/yuin/gopher-lua"
)

type Controller interface {
	CommandChannel() chan lua.LValue
	Frame()
}
