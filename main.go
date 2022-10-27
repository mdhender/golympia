/*
 * golympia - a turn based game
 * Copyright (c) 2022 Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import "log"

func main() {
	clear_alloc_flag()
	dir_assert()
	open_fps()
	load_seed(SEED_FILE)
	map_init()
	read_map("map-data.txt")
	fix_terrain_land()
	set_regions("regions.json")
	set_province_clumps("lands.json")
	unnamed_province_clumps()
	make_islands()
	make_graveyards()
	place_sublocations()
	make_gates()
	make_roads()
	print_map(stdout)
	print_sublocs()
	dump_continents("continents.json")
	count_cities()
	count_continents()
	count_sublocs()
	count_subloc_coverage()
	dump_roads("roads.json")
	dump_gates("gates.json")

	fclose(loc_fp)
	fclose(gate_fp)
	fclose(road_fp)

	count_tiles()

	log.Printf("highest province = %d\n\n", map_[max_row][max_col].region)

	// if the province allocation spilled into the subloc range, we have to increase SUBLOC_MAX
	assert(SUBLOC_LOW > map_[max_row][max_col].region)
}
