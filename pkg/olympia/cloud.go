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

import "log"

const (
	CLOUD_SZ = 4
)

var (
	cloud_region = 0
)

func cloud_disconnect(here, there int) int {
	p1 := p_loc(here)
	if p1 == nil {
		panic("assert(p1)")
	}
	for i := 0; i < len(p1.prov_dest); i++ {
		if p1.prov_dest[i] == there {
			log.Printf("Disconnecting %s from %s.\n", box_name(here), box_name(there))
			p1.prov_dest[i] = 0
			return TRUE
		}
	}
	return FALSE
}

// %%%%
// %*%%	*=Nimbus
// %^%*	*=Aerovia	^=Mt. Olympus link
// *%%%	*=Stratos
func create_cloudlands() {
	var cmap [CLOUD_SZ + 1][CLOUD_SZ + 1]int

	log.Printf("INIT: creating cloudlands\n")

	/*
	 *  Create region wrapper
	 */

	cloud_region = new_ent(T_loc, sub_region)
	set_name(cloud_region, "Cloudlands")

	/*
	 *  Fill cmap[row,col] with locations.
	 */

	for r := 0; r <= CLOUD_SZ; r++ {
		for c := 0; c <= CLOUD_SZ; c++ {
			n := new_ent(T_loc, sub_cloud)
			cmap[r][c] = n
			set_name(n, "Cloud")
			set_where(n, cloud_region)
		}
	}

	// set the NESW exit routes for every map location
	for r := 0; r <= CLOUD_SZ; r++ {
		for c := 0; c <= CLOUD_SZ; c++ {
			p := p_loc(cmap[r][c])
			north, east, south, west := 0, 0, 0, 0
			if r != 0 {
				north = cmap[r-1][c]
			}
			if c != CLOUD_SZ {
				east = cmap[r][c+1]
			}
			if r != CLOUD_SZ {
				south = cmap[r+1][c]
			}
			if c != 0 {
				west = cmap[r][c-1]
			}

			// order is important, we're assuming they are N E S W
			p.prov_dest = append(p.prov_dest, north)
			p.prov_dest = append(p.prov_dest, east)
			p.prov_dest = append(p.prov_dest, south)
			p.prov_dest = append(p.prov_dest, west)
		}
	}

	nimbus := new_ent(T_loc, sub_city)
	set_where(nimbus, cmap[1][1])
	set_name(nimbus, "Nimbus")
	seed_city(nimbus)

	aerovia := new_ent(T_loc, sub_city)
	set_where(aerovia, cmap[2][3])
	set_name(aerovia, "Aerovia")
	seed_city(aerovia)

	stratos := new_ent(T_loc, sub_city)
	set_where(stratos, cmap[3][0])
	set_name(stratos, "Stratos")
	seed_city(stratos)

	// create gates to rings of stones at the four corners of the Cloudlands.
	// may not have rings of stones.
	var l []int
	for _, i := range loop_loc() {
		if subkind(i) == sub_stone_cir {
			l = append(l, i)
		}
	}

	if len(l) >= 4 {
		l = shuffle_ints(l)

		gate1 := new_ent(T_gate, 0)
		set_where(gate1, cmap[0][0])
		p_gate(gate1).to_loc = l[0]
		rp_gate(gate1).seal_key = rnd(111, 999)

		gate2 := new_ent(T_gate, 0)
		set_where(gate2, cmap[CLOUD_SZ][0])
		p_gate(gate2).to_loc = l[1]
		rp_gate(gate2).seal_key = rnd(111, 999)

		gate3 := new_ent(T_gate, 0)
		set_where(gate3, cmap[0][CLOUD_SZ])
		p_gate(gate3).to_loc = l[2]
		rp_gate(gate3).seal_key = rnd(111, 999)

		gate4 := new_ent(T_gate, 0)
		set_where(gate4, cmap[CLOUD_SZ][CLOUD_SZ])
		p_gate(gate4).to_loc = l[3]
		rp_gate(gate4).seal_key = rnd(111, 999)
	}

	log.Printf("Aerovia is in %s\n", box_name(cmap[2][3]))

	// link a cloud to a Mt. Olympus below
	olympus := 0
	for _, i := range loop_mountain() {
		if name(i) == "Mt. Olympus" || olympus == 0 {
			olympus = i
			break
		}
	}
	if olympus == 0 {
		log.Printf("ERROR: Can't find mountain 'Mt. Olympus'\n")
		return
	}
	p := p_loc(cmap[2][1])
	for len(p.prov_dest) < DIR_DOWN {
		p.prov_dest = append(p.prov_dest, 0)
	}
	p.prov_dest[DIR_DOWN-1] = olympus
	p = p_loc(olympus)
	for len(p.prov_dest) < DIR_UP {
		p.prov_dest = append(p.prov_dest, 0)
	}
	p.prov_dest[DIR_UP-1] = cmap[2][1]
}

func float_cloudlands() {
	var cmap [SZ + 1][SZ + 1]int

	tags_off()

	log.Printf("Floating cloudlands.\n")

	// to fill in "map" we take advantage of how the locations were added to the Cloudlands region.
	row, col := 0, 0
	for _, i := range loop_here(cloud_region) {
		if row > CLOUD_SZ {
			panic("assert(!(row > CLOUD_SZ))")
		}
		cmap[row][col] = i
		col++
		if col > CLOUD_SZ {
			col = 0
			row++
		}
	}

	// now let's break all connections between the Cloudlands and any non-cloudlands regions.
	for row := 0; row <= CLOUD_SZ; row++ {
		for col := 0; col <= CLOUD_SZ; col++ {
			for i := 0; i < ilist_len(rp_loc(cmap[row][col]).prov_dest); i++ {
				tmp := rp_loc(cmap[row][col]).prov_dest[i]
				if tmp != 0 && region(tmp) != cloud_region {
					// disconnect...
					cloud_disconnect(cmap[row][col], tmp)
					cloud_disconnect(tmp, cmap[row][col])
				}
			}
		}
	}

	// now we need to determine where over the map the cloudlands is.
	// this is actually stored in the Cloudlands region in that regions "where" box.
	// although a region isn't actually anywhere, that's a convenient fiction for us.
	// it holds the Oly location of map[0][0].
	y, x := region_row(subloc(cloud_region)), region_col(subloc(cloud_region))

	// this causes a problem if there's a hole in the map right here (when we try to save the DB).
	// so let's loop and make sure we don't use an offset that causes a hole.
	for {
		// now calculate the offsets.
		y_off, x_off := rnd(1, 7)-3, rnd(1, 7)-3

		// check for "wrap".
		// Scott Turner: this really needs to be the size of the world!
		x = (x + x_off + xsize) % xsize
		y = (y + y_off + ysize) % ysize

		// so now change the "where" of the cloud_region.
		// we don't use "set_where" since it's not really there :-).
		p_loc_info(cloud_region).where = rc_to_region(y, x)
		if !valid_box(p_loc_info(cloud_region).where) {
			continue
		}

		log.Printf("Cloudlands floating (%d, %d).\n", x_off, y_off)
		break
	}

	// now go through and recalculate where each of the cloudlands is over, and make connections to mountains as appropriate.
	for row := 0; row <= CLOUD_SZ; row++ {
		for col := 0; col <= CLOUD_SZ; col++ {
			// check for "wrap".
			nx, ny := (x+col+xsize)%xsize, (y+row+ysize)%ysize

			// possibly a hole in the map?
			where := rc_to_region(ny, nx)
			if !valid_box(where) {
				continue
			}

			if subkind(where) == sub_mountain {
				// create the connection -- both ways.
				log.Printf("Creating connection between %s and %s.\n", box_name(cmap[row][col]), box_name(where))
				connect_locations(cmap[row][col], DIR_DOWN, where, DIR_UP)
			} else {
				out(where, "Strangely solid clouds cover the sky.")
			}
		}
	}

	tags_on()
}
