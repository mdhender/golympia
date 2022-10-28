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
 *
 */

package olympia

import (
	"log"
)

func GenerateMap(options ...Option) error {
	for _, option := range options {
		if err := option(); err != nil {
			return err
		}
	}

	log.Printf("%-18s == %q\n", "map data", mapDataFilename)
	log.Printf("%-18s == %q\n", "continent data", continentDataFilename)
	log.Printf("%-18s == %q\n", "gate data", gateDataFilename)
	log.Printf("%-18s == %q\n", "land data", landDataFilename)
	log.Printf("%-18s == %q\n", "location data", locationDataFilename)
	log.Printf("%-18s == %q\n", "region data", regionDataFilename)
	log.Printf("%-18s == %q\n", "road data", roadDataFilename)
	log.Printf("%-18s == %q\n", "seed data", seedDataFilename)

	clear_alloc_flag()
	dir_assert()
	open_fps()
	load_seed(seedDataFilename)
	map_init()
	read_map(mapDataFilename)
	fix_terrain_land()
	set_regions(regionDataFilename)
	set_province_clumps(landDataFilename)
	unnamed_province_clumps()
	make_islands()
	make_graveyards()
	place_sublocations()
	make_gates()
	make_roads()
	print_map(stdout)
	print_sublocs()
	dump_continents(continentDataFilename)
	count_cities()
	count_continents()
	count_sublocs()
	count_subloc_coverage()
	dump_roads(roadDataFilename)
	dump_gates(gateDataFilename)

	fclose(loc_fp)
	fclose(gate_fp)
	fclose(road_fp)

	count_tiles()

	log.Println("")
	log.Println("")
	log.Printf("max row, col     = (%2d,%2d)\n", max_row, max_col)
	log.Printf("subloc_low       = %8d\n", SUBLOC_LOW)
	log.Printf("highest province = %8d\n", map_[max_row][max_col].region)

	// if the province allocation spilled into the subloc range, we have to increase SUBLOC_MAX
	if !(SUBLOC_LOW > map_[max_row][max_col].region) {
		panic("assert(SUBLOC_LOW > map_[max_row][max_col].region)")
	}

	return nil
}
