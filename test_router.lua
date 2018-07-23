function send_event(id, event)
	print(id, event)
	return event == "space_join"
end

function join(id, chan)
	print(id .. " joining.")
	return true
end

function leave(id, chan)
end
