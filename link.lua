entity = {
	texture = "resources/images/link.png",
	animations = {
		walking_left = {
			frames = {4, 5},
			interval = 0.4
		},
		walking_right = {
			frames = {6, 7},
			interval = 0.4
		},
		walking_up = {
			frames = {2, 3},
			interval = 0.4
		},
		walking_down = {
			frames = {0, 1},
			interval = 0.4
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
	last_updated = 0.0,
	animation_index = 0,
	location = "test_area"
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

function update_state(milliseconds)
	channel.select(
		{"|<-", controller, function(ok, event)
			if event.pressed then
				if event.key == 32 then
					state.animation = "walking_up"
				elseif event.key == 28 then
					state.animation = "walking_down"
				elseif event.key == 10 then
					state.animation = "walking_left"
				elseif event.key == 13 then
					state.animation = "walking_right"
				end
			else
				set_idle()
			end
		end},
		{"default", function()
		end}
	)
	local animation = entity.animations[state.animation]
	if milliseconds - state.last_updated > (1000 * animation.interval) then
		state.animation_index = state.animation_index + 1
		state.animation_index = state.animation_index % #animation.frames
		state.last_updated = milliseconds
	end
end

function active_frame()
	local frame = entity.animations[state.animation].frames[state.animation_index + 1]
	return frame
end
