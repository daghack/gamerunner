entity = {
	tileset = "resources/images/link_walking.png",
	velocity = 3.0,
	animations = {
		walking_left = {
			frames = {4, 5},
			interval = 0.15
		},
		walking_right = {
			frames = {6, 7},
			interval = 0.15
		},
		walking_up = {
			frames = {2, 3},
			interval = 0.15
		},
		walking_down = {
			frames = {0, 1},
			interval = 0.15
		},
		idle_left = {
			frames = {4},
			interval = 1
		},
		idle_right = {
			frames = {6},
			interval = 1
		},
		idle_up = {
			frames = {2},
			interval = 1
		},
		idle_down = {
			frames = {0},
			interval = 1
		}
	}
}

state = {
	animation = "walking_down",
	last_updated = 0,
	animation_index = 0,
	location = "test_area",
	x_pos = 4,
	y_pos = 4
}

function set_idle()
	local animp = ""
	if state.animation == "walking_down" then
		animp = "idle_down"
	elseif state.animation == "walking_up" then
		animp = "idle_up"
	elseif state.animation == "walking_left" then
		animp = "idle_left"
	elseif state.animation == "walking_right" then
		animp = "idle_right"
	else
		return
	end
	state.animation = animp
	state.animation_index = 0
end

events_emptied = false
events_present = false

command_handler = {
	move_up = "walking_up",
	move_down = "walking_down",
	move_left = "walking_left",
	move_right = "walking_right"
}

function update_state(milliseconds, updates_per_second)
	if updates_per_second == 0 then
		updates_per_second = 60
	end
	events_emptied = false
	events_present = false
	while not events_emptied do
		channel.select(
			{"|<-", controller, function(ok, event)
				events_present = true
				state.animation = command_handler[event]
			end},
			{"default", function()
				events_emptied = true
			end}
		)
	end
	if not events_present then
		set_idle()
	end
	move(updates_per_second)
	local animation = entity.animations[state.animation]
	if milliseconds - state.last_updated > (1000 * animation.interval) then
		state.animation_index = state.animation_index + 1
		state.animation_index = state.animation_index % #animation.frames
		state.last_updated = milliseconds
	end
end

function move(updates_per_second)
	if state.animation == "walking_down" then
		state.y_pos = state.y_pos + (entity.velocity / updates_per_second)
	elseif state.animation == "walking_up" then
		state.y_pos = state.y_pos - (entity.velocity / updates_per_second)
	elseif state.animation == "walking_left" then
		state.x_pos = state.x_pos - (entity.velocity / updates_per_second)
	elseif state.animation == "walking_right" then
		state.x_pos = state.x_pos + (entity.velocity / updates_per_second)
	end
end

function active_frame()
	local frame = entity.animations[state.animation].frames[state.animation_index + 1]
	return frame
end
