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
	"fmt"
	"log"
	"path/filepath"
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

	if _, err := SysDataLoad(filepath.Join(libdir, "sysdata.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}

	clear_alloc_flag()
	dir_assert()
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

	count_cities()
	count_continents()
	count_sublocs()
	count_subloc_coverage()
	count_tiles()

	// add in the items that the engine expects to always have
	for _, monster := range monster_tbl {
		BoxAlloc(monster.Id, strKind[monster.Kind], strSubKind[monster.SubKind])
		bx[monster.Id].x_item = monster.toBox()
	}

	log.Printf("max row, col     = (%2d,%2d)\n", max_row, max_col)
	log.Printf("subloc_low       = %8d\n", SUBLOC_LOW)
	log.Printf("highest province = %8d\n", map_[max_row][max_col].region)

	// if the province allocation spilled into the subloc range, we have to increase SUBLOC_MAX
	if !(SUBLOC_LOW > map_[max_row][max_col].region) {
		panic("assert(SUBLOC_LOW > map_[max_row][max_col].region)")
	}

	log.Println("")
	log.Println("")

	/* check database integrity */
	if err := check_db(); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	// and save
	if err := save_db(); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}

	if err := print_map(filepath.Join(libdir, "map-data.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := print_sublocs(filepath.Join(libdir, "subloc-data.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := dump_continents(filepath.Join(libdir, continentDataFilename)); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := RoadDataSave(filepath.Join(libdir, roadDataFilename)); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := GateDataSave(filepath.Join(libdir, gateDataFilename)); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := CharacterDataSave(filepath.Join(libdir, "characters")); err != nil {
		log.Println(fmt.Errorf("GenerateMap: %w", err))
		//return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := LocationDataSave(filepath.Join(libdir, "locations.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := EntityItemDataSave(filepath.Join(libdir, "items.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := MiscDataSave(filepath.Join(libdir, "misc.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := NationDataSave(filepath.Join(libdir, "nations.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := ShipDataSave(filepath.Join(libdir, "ships.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := SkillDataSave(filepath.Join(libdir, "skills.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := SysDataSave(filepath.Join(libdir, "sysdata.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}
	if err := UnformDataSave(filepath.Join(libdir, "unform.json")); err != nil {
		return fmt.Errorf("GenerateMap: %w", err)
	}

	return nil
}
