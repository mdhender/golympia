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
	"os"
)

const (
	SZ = 6 /* SZ x SZ is size of Hades */
)

var (
	hades_pit    = 0 /* Pit of Hades */
	hades_player = 0
	hades_region = 0
)

func auto_hades() {
	for n_hades := len(loop_units(hades_player)); n_hades < 25; n_hades++ {
		create_hades_nasty()
	}
}

func create_hades() {
	var hmap [SZ + 1][SZ + 1]int

	log.Printf("INIT: creating hades\n")

	// create region wrapper for Hades
	hades_region = new_ent(T_loc, sub_region)
	set_name(hades_region, "Hades")

	// create the King of Hades player
	if hades_player != 0 {
		panic("assert(hades_player == 0)")
	}

	hades_player = 205
	alloc_box(hades_player, T_player, sub_pl_npc)
	set_name(hades_player, "King of Hades")
	p_player(hades_player).password = DEFAULT_PASSWORD

	// fill hmap[row,col] with locations
	for r := 0; r <= SZ; r++ {
		for c := 0; c <= SZ; c++ {
			n := new_ent(T_loc, sub_under)
			hmap[r][c] = n
			set_name(n, "Hades")
			set_where(n, hades_region)
			p_loc(n).hidden = TRUE
			set_known(hades_player, n)
		}
	}

	// set the NSEW exit routes for every map location
	for r := 0; r <= SZ; r++ {
		for c := 0; c <= SZ; c++ {
			var north, south, east, west int

			p := p_loc(hmap[r][c])
			if r == 0 {
				north = 0
			} else {
				north = hmap[r-1][c]
			}
			if r == SZ {
				south = 0
			} else {
				south = hmap[r+1][c]
			}
			if c == SZ {
				east = 0
			} else {
				east = hmap[r][c+1]
			}
			if c == 0 {
				west = 0
			} else {
				west = hmap[r][c-1]
			}
			p.prov_dest = append(p.prov_dest, north)
			p.prov_dest = append(p.prov_dest, east)
			p.prov_dest = append(p.prov_dest, south)
			p.prov_dest = append(p.prov_dest, west)
		}
	}

	// place a city in the center of the map, with the Pit of Hades inside the city.
	// this is the only place where Necromancy is taught.
	city := new_ent(T_loc, sub_city)
	set_where(city, hmap[SZ/2][SZ/2])
	set_name(city, "City of the Dead")
	set_known(hades_player, city)

	p := p_subloc(city)
	p.teaches = append(p.teaches, sk_necromancy)

	pit := new_ent(T_loc, sub_hades_pit)
	set_where(pit, city)
	set_name(pit, "Pit of Hades")
	set_known(hades_player, pit)

	hades_pit = pit

	// dual-link every graveyard from the world into one of the Hades locations except the center one containing the pit
	var l []int
	for _, i := range loop_loc() {
		if subkind(i) == sub_graveyard {
			l = append(l, i)
			set_known(hades_player, i)
		}

	}
	l = shuffle_ints(l)

	for i := 0; i < len(l); {
		for r := 0; r <= SZ && i < len(l); r++ {
			for c := 0; c <= SZ && i < len(l); c++ {
				if r == SZ/2 && c == SZ/2 {
					continue
				}

				p := p_subloc(l[i])
				p.link_to = append(p.link_to, hmap[r][c])

				// p.link_when = -1
				// p.link_open = -1

				p = p_subloc(hmap[r][c])
				p.link_from = append(p.link_from, l[i])

				i++
			}
		}
	}

	log.Printf("hades loc is %q\n", box_name(hmap[1][1]))
}

func create_hades_nasty() {
	p := rp_loc_info(hades_region)
	if p == nil {
		panic("assert(p)")
	}
	where := p.here_list[rnd(0, ilist_len(p.here_list)-1)]

	var nasty int
	switch rnd(1, 5) {
	case 1:
		nasty = new_char(sub_ni, item_spirit, where, 100, hades_player, LOY_npc, 0, "Tortured spirits")
		if nasty < 0 {
			return
		}
		p_char(nasty).npc_prog = PROG_balrog
		gen_item(nasty, item_spirit, rnd(25, 75))
	case 2:
		nasty = new_char(sub_ni, item_spirit, where, 100, hades_player, LOY_npc, 0, "Ghostly presence")
		if nasty < 0 {
			return
		}
		p_char(nasty).npc_prog = PROG_balrog
		p_char(nasty).attack = 100
		rp_char(nasty).defense = 100
	case 3:
		nasty = new_char(sub_ni, item_lesser_demon, where, 100, hades_player, LOY_npc, 0, "Lesser Demon")
		if nasty < 0 {
			return
		}
		p_char(nasty).npc_prog = PROG_balrog
		p_char(nasty).attack = 250
		rp_char(nasty).defense = 250
		gen_item(nasty, item_spirit, rnd(50, 150))
	case 4:
		nasty = new_char(sub_ni, item_greater_demon, where, 100, hades_player, LOY_npc, 0, "Greater Demon")
		if nasty < 0 {
			return
		}
		p_char(nasty).npc_prog = PROG_balrog
		p_char(nasty).attack = 500
		rp_char(nasty).defense = 500
		gen_item(nasty, item_spirit, rnd(100, 250))
	case 5:
		nasty = create_new_beasts(where, sub_undead)
		p_char(nasty).npc_prog = PROG_balrog
		fprintf(os.Stderr, "Created undead in Hades: %s.\n", box_name(nasty))
	default:
		panic("!reached")
	}

	queue(nasty, "wait time 0")
	// init_load_sup(new);   // make ready to execute commands immediately
}
