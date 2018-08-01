function send_event(id, event)
	return true
end

function pre_join(id, chan)
	print(id, "joining!")
	return not listener_exists(id)
end

function post_join(id, chan)
	print(id, "joined!")
	self.send_event(id, "space_join")
	return true
end

function pre_leave(id, chan)
	print(id, "leaving!")
	return true
end

function post_leave(id, chan)
	print(id, "left!")
	self.send_event(id, "space_leave")
	return true
end

areas = {
	test_area = "test_area.lua"
}
