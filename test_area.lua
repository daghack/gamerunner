map_config = {
	tileset = "resources/images/zelda_tiles.png",
	tile_size = 16,
	map_width = 8,
	map_height = 8,
	tilemap = {
		209, 209, 209, 209, 209, 209, 209, 209,
		209, 209, 209, 209, 209, 209, 209, 209,
		209, 209, 005, 006, 006, 007, 209, 209,
		209, 209, 008, 112, 113, 009, 209, 209,
		209, 209, 008, 128, 129, 009, 209, 209,
		209, 209, 010, 011, 011, 012, 209, 209,
		209, 209, 209, 209, 209, 209, 209, 209,
		209, 209, 209, 209, 209, 209, 209, 209}
}

function update_tiles()
	for k, v in pairs(map_config.tilemap) do
		if v == 209 then
			map_config.tilemap[k] = 210
		elseif v == 210 then
			map_config.tilemap[k] = 209
		end
	end
	reload_map()
end
