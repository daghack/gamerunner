package eventrouter

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

type Router struct {
	state       *lua.LState
	listenerMap map[string]chan lua.LValue
}

func NewRouter(routingFile string) (*Router, error) {
	l := lua.NewState()
	l.OpenLibs()
	err := l.DoFile(routingFile)
	if err != nil {
		l.Close()
		return nil, err
	}
	toret := &Router{
		state:       l,
		listenerMap: map[string]chan lua.LValue{},
	}
	l.SetGlobal("listener_exists", l.NewFunction(toret.listenerExists))
	toret.AddToLua("self", l)
	return toret, nil
}

func (r *Router) Close() {
	r.state.Close()
}

func (r *Router) idChanHook(funcname, id string, echan chan lua.LValue) (bool, error) {
	err := r.state.CallByParam(lua.P{
		Fn:      r.state.GetGlobal(funcname),
		NRet:    1,
		Protect: true,
	}, lua.LString(id), lua.LChannel(echan))
	if err != nil {
		return false, err
	}
	ret := r.state.ToBool(-1)
	r.state.Pop(1)
	return ret, nil
}

func (r *Router) listenerExists(l *lua.LState) int {
	id := l.ToString(-1)
	l.Pop(1)
	_, ok := r.listenerMap[id]
	l.Push(lua.LBool(ok))
	return 1
}

func (r *Router) Join(l *lua.LState) int {
	id := l.ToString(-2)
	echan := l.ToChannel(-1)
	l.Pop(2)
	cont, err := r.idChanHook("pre_join", id, echan)
	if err != nil {
		fmt.Println(err)
	}
	if cont {
		r.listenerMap[id] = echan
		cont, err = r.idChanHook("post_join", id, echan)
		if err != nil {
			fmt.Println(err)
		}
		if cont {
			l.Push(lua.LBool(true))
		} else {
			delete(r.listenerMap, id)
			l.Push(lua.LBool(false))
		}
	} else {
		l.Push(lua.LBool(false))
	}
	return 1
}

func (r *Router) Leave(l *lua.LState) int {
	id := l.ToString(-1)
	l.Pop(1)
	if echan, ok := r.listenerMap[id]; !ok {
		l.Push(lua.LBool(true))
	} else {
		cont, err := r.idChanHook("pre_leave", id, echan)
		if err != nil {
			fmt.Println(err)
		}
		if cont {
			r.listenerMap[id] = echan
			cont, err = r.idChanHook("post_leave", id, echan)
			if err != nil {
				fmt.Println(err)
			}
			if cont {
				l.Push(lua.LBool(true))
			} else {
				delete(r.listenerMap, id)
				l.Push(lua.LBool(false))
			}
		} else {
			l.Push(lua.LBool(false))
		}
	}
	return 1
}

func (r *Router) SendEvent(l *lua.LState) int {
	id := l.ToString(-2)
	event := l.Get(-1)
	l.Pop(2)

	err := r.state.CallByParam(lua.P{
		Fn:      r.state.GetGlobal("send_event"),
		NRet:    1,
		Protect: true,
	}, lua.LString(id), event)

	if err != nil {
		l.Push(lua.LBool(false))
		return 1
	}
	ret := r.state.ToBool(-1)
	r.state.Pop(1)

	if echan, ok := r.listenerMap[id]; ok && ret {
		echan <- event
		l.Push(lua.LBool(true))
	} else {
		l.Push(lua.LBool(false))
	}
	return 1
}

func (r *Router) AddToLua(name string, l *lua.LState) {
	l.SetGlobal(name, l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"send_event": r.SendEvent,
		"join":       r.Join,
		"leave":      r.Leave,
	}))
}
