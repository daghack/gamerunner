package space

import (
	"github.com/yuin/gopher-lua"
)

const luaCode = `
Router:join("johnny", ch)
local exit = false
function handle(ok, v)
	print("handling", v)
	if not ok then
		print("channel close")
		exit = true
	else
		print(v)
		if v == "space_join" then
			Router:send_event("johnny", "space_leave")
		elseif v == "space_leave" then
		else
			Router:send_event("johnny", "space_join")
		end
	end
end
while not exit do
	channel.select({"|<-", ch, handle})
end
`

type Space struct {
	id           string
	state        *lua.LState
	eventChannel chan lua.LValue
	entityMap    map[string]chan<- lua.LValue
}

type SpaceManager struct {
	spaceMap map[string]*Space
}

func (s *Space) processEvents() {
	s.state.SetGlobal("ch", lua.LChannel(s.eventChannel))
	err := s.state.DoString(luaCode)
	if err != nil {
		panic(err)
	}
}

func LoadSpace(base *lua.LState) (*Space, error) {
	toret := &Space{
		state:        base,
		eventChannel: make(chan lua.LValue, 16),
		entityMap:    map[string]chan<- lua.LValue{},
	}
	base.SetGlobal("entity_join", base.NewFunction(toret.EntityJoin))
	base.SetGlobal("entity_leave", base.NewFunction(toret.EntityLeave))
	go toret.processEvents()
	return toret, nil
}

func (s *Space) SendEvent(event string) {
	s.eventChannel <- lua.LString(event)
}

func (s *Space) EntityJoin(l *lua.LState) int {
	id := l.ToString(1)
	echan := l.ToChannel(2)
	l.Pop(2)
	if _, ok := s.entityMap[id]; !ok {
		s.entityMap[id] = echan
		echan <- lua.LString("space_join")
	}
	return 0
}

func (s *Space) EntityLeave(l *lua.LState) int {
	id := l.ToString(1)
	l.Pop(1)
	if echan, ok := s.entityMap[id]; ok {
		delete(s.entityMap, id)
		echan <- lua.LString("space_leave")
	}
	return 0
}
