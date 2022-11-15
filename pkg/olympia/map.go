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
	"io"
	"log"
)

// some cute routines to print out a view of the map
const (
	MAX_X = 100
	MAX_Y = 100
)

var (
	cmap [MAX_X][MAX_Y]byte
)

// fill the map with characters corresponding to the provinces.
func load_cmap() int {
	if xsize >= MAX_X || ysize >= MAX_Y {
		return 0
	}

	// default everything to impassible
	for x := 0; x < xsize; x++ {
		for y := 0; y < ysize; y++ {
			cmap[x][y] = '#'
		}
	}

	// now update base on the neighbors
	for _, i := range loop_kind(T_loc) {
		if loc_depth(i) != LOC_province {
			continue
		} else if region(i) == faery_region || region(i) == hades_region || region(i) == cloud_region {
			continue
		}
		x, y := region_col(i), region_row(i)
		switch subkind(i) {
		case sub_ocean:
			cmap[x][y] = ' '
		case sub_forest:
			cmap[x][y] = '%'
		case sub_plain, sub_island:
			cmap[x][y] = '.'
		case sub_mountain:
			cmap[x][y] = '^'
		case sub_mine_shaft:
			cmap[x][y] = '0'
		case sub_desert:
			cmap[x][y] = '-'
		case sub_swamp:
			cmap[x][y] = 's'
		default:
			log.Printf("Unknown subtype: %d.\n", subkind(i))
			cmap[x][y] = '?'
		}
	}
	return 1
}

// indicate where players are...
func load_cmap_players() int {
	if xsize >= MAX_X || ysize >= MAX_Y {
		return 0
	}

	for _, i := range loop_kind(T_loc) {
		if loc_depth(i) != LOC_province {
			continue
		} else if region(i) == faery_region || region(i) == hades_region || region(i) == cloud_region {
			continue
		}
		x, y := region_col(i), region_row(i)

		count := 0
		for _, j := range loop_all_here(i) {
			if kind(j) == T_char && !is_real_npc(j) {
				count++
			}
		}
		if count == 0 {
			continue
		} else if count < 10 {
			cmap[x][y] = byte('0' + byte(count))
		} else {
			cmap[x][y] = '*'
			log.Printf("%q at (%2d,%2d) has %d nobles.\n", box_name(i), x, y, count)
		}
	}

	return 1
}

func print_cmap(w io.Writer) {
	if xsize >= MAX_X {
		return
	} else if ysize >= MAX_Y {
		return
	}

	for y := 0; y < ysize; y++ {
		for x := 0; x < xsize; x++ {
			_, _ = w.Write([]byte{cmap[x][y]})
		}
		_, _ = w.Write([]byte{'\n'})
	}
}
