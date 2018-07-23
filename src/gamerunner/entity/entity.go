package entity

import (
	"gamerunner/eventrouter"
	"github.com/yuin/gopher-lua"
)

const luaCode = `
Router:join("johnny", echan)
Router:send_event("johnny", "ping")
function handle(ok, event)
	print(event)
	if event == "ping" then
		--os.execute("say ping")
		Router:send_event("johnny", "pong")
	elseif event == "pong" then
		--os.execute("say pong")
		Router:leave("johnny")
	end
end
while true do
	channel.select({"|<-", echan, handle})
end
`

type Entity struct {
	state     *lua.LState
	eventChan chan lua.LValue
}

func NewEntity(r *eventrouter.Router) *Entity {
	l := lua.NewState()
	toret := &Entity{
		state:     l,
		eventChan: make(chan lua.LValue, 16),
	}

	r.AddToLua("Router", l)
	l.SetGlobal("echan", lua.LChannel(toret.eventChan))

	return toret
}

func (e *Entity) EventLoop() error {
	return e.state.DoString(luaCode)
}
