map_config = {
	tileset = "resources/images/zelda_tiles.png",
	tile_size = 16,
	map_width = 8,
	map_height = 8,
	tilemap = {
		undrawn_tiles = {0, 1, 2},
		collision_tiles = {0, 1, 2},
		layers = {
			{209, 209, 209, 209, 209, 209, 209, 209,
			209, 209, 209, 209, 209, 209, 209, 209,
			209, 209, 0, 0, 0, 0, 209, 209,
			209, 209, 0, 0, 0, 0, 209, 209,
			209, 209, 0, 0, 0, 0, 209, 209,
			209, 209, 0, 0, 0, 0, 209, 209,
			209, 209, 209, 209, 209, 209, 209, 209,
			209, 209, 209, 209, 209, 209, 209, 209},
			{0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 005, 006, 006, 007, 0, 0,
			0, 0, 008, 0, 0, 009, 0, 0,
			0, 0, 008, 0, 0, 009, 0, 0,
			0, 0, 010, 011, 011, 012, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 112, 113, 0, 0, 0,
			0, 0, 0, 128, 129, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0}
		}
	}
}

function update_tiles()
	for k, layer in pairs(map_config.tilemap.layers) do
		for i, v in pairs(layer) do
			if v == 209 then
				layer[i] = 210
			elseif v == 210 then
				layer[i] = 209
			end
		end
	end
	reload_map()
end

last_updated = 0

function update_area(milliseconds)
	if milliseconds - last_updated > 500 then
		update_tiles()
		last_updated = milliseconds
	end
end
