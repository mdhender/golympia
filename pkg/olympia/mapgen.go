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
	"encoding/json"
	"fmt"
	"github.com/mdhender/golympia/pkg/io"
	"log"
	"math/rand"
	"os"
	"strings"
)

// Olympia map generator:
//   Abandon all hope, ye who enter here.

// mapgen.h

// mapgen.c

// todo: remove generation of special cities, skills for cities
// todo: gate rings might be better made in the engine
// todo: don't forget to implement Summerbridge and Uldim pass in the engine
// todo: make hades bigger to accommodate more graveyards

var SEED_FILE = "randseed.json"

const (
	MAX_BOX = 100_000
	MAX_ROW = 100
	MAX_COL = 100
)

const (
	GATES_CONTINENTAL_TOUR = TRUE  // on or off
	GATES_OTHER            = FALSE // the rest of the gates
	GATES_STONE_CIRCLES    = FALSE // on or off
	GATE_TIMES             = 0     // VLN: number of gates?  was 25
)

const (
	DIR_N   = 1
	DIR_NE  = 2
	DIR_E   = 3
	DIR_SE  = 4
	DIR_S   = 5
	DIR_SW  = 6
	DIR_W   = 7
	DIR_NW  = 8
	MAX_DIR = DIR_NW + 1
)

// G2 Entity coding system:
//
//	range___________   extent   use_______________________________
//	     1 -     999      999   items
//	 1,000 -   8,999    8,000   chars
//	 9,000 -   9,999    1,000   skills
//	10,000 -  19,999   10,000   provinces        (CCNN: AA00-DV99)
//	20,000 -  49,999   20,000   more provinces   (CCNN: DW00-ZZ99)
//	50,000 -  56,759    6,760   player entities  (CCN)
//	56,760 -  58,759    2,000   lucky locs       (CNN)
//	58,760 -  58,999      240   regions          (NNNNN)
//	59,000 -  78,999   20,000   sublocs, misc    (CNNN: A000-Z999)
//	79,000 - 100,000   21,000   storms           (NNNNN)
//
//	Note: restricted alphabet, no vowels (except a) or l:
//	  "abcdfghjkmnpqrstvwxz"
const (
	CITY_LOW    = 56_760
	CITY_HIGH   = 58_759
	REGION_OFF  = 58_760 /* where to start numbering regions */
	SUBLOC_LOW  = 59_000
	SUBLOC_HIGH = 78_999
)

var alloc_flag [MAX_BOX]int
var bridge_dir_s = []string{"-invalid-", "  n-s", "  e-w", "ne-sw", "nw-se"}
var dir_vector [MAX_DIR]int

type guild_name_t struct {
	skill  int
	weight int
	name   string
}

var guild_names = []guild_name_t{
	{terr_stone_cir, 1, ""},
	{terr_grove, 1, ""},
	{terr_bog, 1, ""},
	{terr_cave, 1, ""},
	{terr_grave, 20, ""},
	{terr_grave, 1, "Barrows"},
	{terr_grave, 1, "Barrow Downs"},
	{terr_grave, 1, "Barrow Hills"},
	{terr_grave, 1, "Cairn Hills"},
	{terr_grave, 1, "Catacombs"},
	{terr_grave, 1, "Grave Mounds"},
	{terr_grave, 1, "Place of the Dead"},
	{terr_grave, 1, "Cemetery Hill"},
	{terr_grave, 1, "Fields of Death"},
	{terr_ruins, 1, ""},
	{terr_battlefield, 3, "Old battlefield"},
	{terr_battlefield, 1, "Ancient battlefield"},
	{terr_battlefield, 1, ""},
	{terr_ench_for, 1, ""},
	{terr_rocky_hill, 1, ""},
	{terr_tree_cir, 1, ""},
	{terr_pits, 1, "Cursed Pits"},
	{terr_pasture, 1, ""},
	{terr_pasture, 1, "Grassy field"},
	{terr_sac_grove, 1, ""},
	{terr_oasis, 1, ""},
	{terr_pop_field, 1, ""},
	{terr_sand_pit, 1, ""},
	{terr_yew_grove, 1, ""},
	{terr_temple, 1, ""},
	{terr_lair, 1, ""},
}

type loc_table_t struct {
	terr   int // terrain appropriate
	kind   int // what to make there
	weight int // weight given to selection
	hidden int // 0=no, 1=yes, 2=rnd(0,1)
}

var loc_table = []loc_table_t{
	{terr_desert, terr_cave, 10, 1},
	{terr_desert, terr_oasis, 10, 2},
	{terr_desert, terr_sand_pit, 10, 2},
	{terr_mountain, terr_ruins, 10, 1},
	{terr_mountain, terr_cave, 10, 1},
	{terr_mountain, terr_yew_grove, 10, 2},
	{terr_mountain, terr_lair, 10, 2},
	{terr_mountain, terr_battlefield, 6, 2},
	{terr_swamp, terr_bog, 10, 2},
	{terr_swamp, terr_pits, 10, 2},
	{terr_swamp, terr_battlefield, 6, 2},
	{terr_swamp, terr_lair, 5, 2},
	{terr_forest, terr_ruins, 10, 1},
	{terr_forest, terr_tree_cir, 10, 1},
	{terr_forest, terr_ench_for, 10, 1},
	{terr_forest, terr_yew_grove, 10, 2},
	{terr_forest, terr_cave, 10, 1},
	{terr_forest, terr_grove, 9, 1},
	{terr_forest, terr_battlefield, 6, 2},
	{terr_swamp, terr_lair, 3, 1},
	{terr_plain, terr_ruins, 10, 1},
	{terr_plain, terr_pasture, 10, 0},
	{terr_plain, terr_rocky_hill, 10, 0},
	{terr_plain, terr_sac_grove, 10, 2},
	{terr_plain, terr_pop_field, 10, 0},
	{terr_plain, terr_cave, 10, 1},
	{terr_plain, terr_battlefield, 6, 2},
}

const (
	terr_land        = 1
	terr_ocean       = 2
	terr_forest      = 3
	terr_swamp       = 4
	terr_mountain    = 5
	terr_plain       = 6
	terr_desert      = 7
	terr_water       = 8
	terr_island      = 9
	terr_stone_cir   = 10 /* circle of stones */
	terr_grove       = 11 /* mallorn grove */
	terr_bog         = 12
	terr_cave        = 13
	terr_city        = 14
	terr_guild       = 15
	terr_grave       = 16
	terr_ruins       = 17
	terr_battlefield = 18
	terr_ench_for    = 19 /* enchanted forest */
	terr_rocky_hill  = 20
	terr_tree_cir    = 21
	terr_pits        = 22
	terr_pasture     = 23
	terr_oasis       = 24
	terr_yew_grove   = 25
	terr_sand_pit    = 26
	terr_sac_grove   = 27 /* sacred grove */
	terr_pop_field   = 28 /* poppy field */
	terr_temple      = 29
	terr_lair        = 30 /* dragon lair */
)

var terr_s = []string{
	"<null>",
	"land",
	"ocean",
	"forest",
	"swamp",
	"mountain",
	"plain",
	"desert",
	"water",
	"island",
	"ring of stones",
	"mallorn grove",
	"bog",
	"cave",
	"city",
	"guild",
	"graveyard",
	"ruins",
	"field",
	"enchanted forest",
	"rocky hill",
	"circle of trees",
	"pits",
	"pasture",
	"oasis",
	"yew grove",
	"sand pit",
	"sacred grove",
	"poppy field",
	"temple",
	"lair",
}

const MAX_INSIDE = 500 /* max continents/regions */

var inside_names [MAX_INSIDE]string

var inside_list [MAX_INSIDE][]*tile   // list of province tiles in each region
var inside_gates_to [MAX_INSIDE]int   // for info gathering only
var inside_gates_from [MAX_INSIDE]int // for info gathering only
var inside_num_cities [MAX_INSIDE]int // for info gathering only
var inside_top = 0

var max_col = 0

var max_row = 0

var water_count = 0 /* count of water provinces */
var land_count = 0  /* count of land provinces */
var num_islands = 0

var map_ [MAX_ROW][MAX_COL]*tile

const MAX_SUBLOC = 20000

var subloc [MAX_SUBLOC]*tile

var top_subloc = 0

var loc_fp *io.FILE
var gate_fp *io.FILE
var road_fp *io.FILE

func alloc_inside() int {
	inside_top++
	if !(inside_top < MAX_INSIDE) {
		panic("assert(inside_top < MAX_INSIDE)")
	}
	return inside_top
}

func add_road(from *tile, to_loc int, hidden int, name string) {
	from.roads = append(from.roads, &road{
		ent_num: rnd_alloc_num(SUBLOC_LOW, SUBLOC_HIGH),
		to_loc:  to_loc,
		hidden:  hidden,
		name:    name,
	})
}

func adjacent_tile_terr(row, col int) *tile {
	randomize_dir_vector()

	p := adjacent_tile_sup(row, col, dir_vector[1])
	if p == nil || p.terrain == terr_land || p.terrain == terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[2])
	}
	if p == nil || p.terrain == terr_land || p.terrain == terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[3])
	}
	if p == nil || p.terrain == terr_land || p.terrain == terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[4])
	}
	if p == nil || p.terrain == terr_land || p.terrain == terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[5])
	}
	if p == nil || p.terrain == terr_land || p.terrain == terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[6])
	}
	if p == nil || p.terrain == terr_land || p.terrain == terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[7])
	}
	if p == nil || p.terrain == terr_land || p.terrain == terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[8])
	}

	return p
}

// Return the region immediately adjacent to <location> in direction <dir>
// Returns 0 if there is no adjacent location in the given direction.
func adjacent_tile_sup(row, col int, dir int) *tile {
	switch dir {
	case DIR_N:
		row--
	case DIR_NE:
		row--
		col++
	case DIR_E:
		col++
	case DIR_SE:
		row++
		col++
	case DIR_S:
		row++
	case DIR_SW:
		row++
		col--
	case DIR_W:
		col--
	case DIR_NW:
		row--
		col--
	default:
		panic(fmt.Sprintf("assert(dir != %d)", dir))
	}
	// TODO: maybe this should use the same wrapping logic as prov_dest()?
	// TODO: should this use max_row, MAX_ROW?
	if row < 0 || row > 99 || col < 0 || col > 99 { // off the map
		// fixme: vln: might need to fix wrapping here, too
		return nil
	}

	return map_[row][col]
}

func adjacent_tile_water(row, col int) *tile {
	randomize_dir_vector()

	p := adjacent_tile_sup(row, col, dir_vector[1])
	if p == nil || p.terrain != terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[2])
	}
	if p == nil || p.terrain != terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[3])
	}
	if p == nil || p.terrain != terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[4])
	}
	if p == nil || p.terrain != terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[5])
	}
	if p == nil || p.terrain != terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[6])
	}
	if p == nil || p.terrain != terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[7])
	}
	if p == nil || p.terrain != terr_ocean {
		p = adjacent_tile_sup(row, col, dir_vector[8])
	}

	return p
}

func bridge_caddy_corners() {
	// todo: the max_row vs MAX_ROW
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].terrain != terr_ocean && rnd(1, 35) == 35 {
				bridge_corner_sup(row, col)
			}
		}
	}
}

var _static_bridge_corner_sup = struct {
	road_name_cnt int
}{
	road_name_cnt: 0,
}

func bridge_corner_sup(row, col int) int {
	// find all squares neighboring the hole
	n := adjacent_tile_sup(row, col, DIR_N)
	if n != nil && n.mark != 0 {
		return FALSE
	}
	s := adjacent_tile_sup(row, col, DIR_S)
	if s != nil && s.mark != 0 {
		return FALSE
	}
	e := adjacent_tile_sup(row, col, DIR_E)
	if s != nil && e.mark != 0 {
		return FALSE
	}
	w := adjacent_tile_sup(row, col, DIR_W)
	if w != nil && w.mark != 0 {
		return FALSE
	}
	nw := adjacent_tile_sup(row, col, DIR_NW)
	if nw != nil && nw.mark != 0 {
		return FALSE
	}
	sw := adjacent_tile_sup(row, col, DIR_SW)
	if sw != nil && sw.mark != 0 {
		return FALSE
	}
	ne := adjacent_tile_sup(row, col, DIR_NE)
	if ne != nil && ne.mark != 0 {
		return FALSE
	}
	se := adjacent_tile_sup(row, col, DIR_SE)
	if se != nil && se.mark != 0 {
		return FALSE
	}

	var name string
	switch _static_bridge_corner_sup.road_name_cnt % 3 {
	case 0:
		name = "Secret pass"
	case 1:
		name = "Secret route"
	case 2:
		name = "Old road"
	default:
		panic("!reached")
	}
	_static_bridge_corner_sup.road_name_cnt++

	var l []int
	if nw != nil && nw.terrain != terr_ocean {
		l = append(l, 1)
	}
	if ne != nil && ne.terrain != terr_ocean {
		l = append(l, 2)
	}
	if se != nil && se.terrain != terr_ocean {
		l = append(l, 3)
	}
	if sw != nil && sw.terrain != terr_ocean {
		l = append(l, 4)
	}

	if len(l) == 0 {
		return FALSE
	}

	// the horror, the horror
	if n != nil {
		n.mark += rnd(0, 1)
	}
	if e != nil {
		e.mark += rnd(0, 1)
	}
	if w != nil {
		w.mark += rnd(0, 1)
	}
	if s != nil {
		s.mark += rnd(0, 1)
	}
	if nw != nil {
		nw.mark += rnd(0, 1)
	}
	if ne != nil {
		ne.mark += rnd(0, 1)
	}
	if sw != nil {
		sw.mark += rnd(0, 1)
	}
	if se != nil {
		se.mark += rnd(0, 1)
	}

	i := rnd(0, len(l)-1)
	switch l[i] {
	case 1:
		link_roads(map_[row][col], nw, TRUE, name)
		map_[row][col].mark = 1
		nw.mark = 1
	case 2:
		link_roads(map_[row][col], ne, TRUE, name)
		map_[row][col].mark = 1
		ne.mark = 1
	case 3:
		link_roads(map_[row][col], se, TRUE, name)
		map_[row][col].mark = 1
		se.mark = 1
	case 4:
		link_roads(map_[row][col], sw, TRUE, name)
		map_[row][col].mark = 1
		sw.mark = 1
	default:
		panic("!reached")
	}

	return l[i]
}

var _static_bridge_map_hole_sup = struct {
	road_name_cnt int
}{
	road_name_cnt: 0,
}

func bridge_map_hole_sup(row, col int) int {
	// find all squares neighboring the hole
	n := adjacent_tile_sup(row, col, DIR_N)
	if n != nil && n.mark != 0 {
		return FALSE
	}
	s := adjacent_tile_sup(row, col, DIR_S)
	if s != nil && s.mark != 0 {
		return FALSE
	}
	e := adjacent_tile_sup(row, col, DIR_E)
	if s != nil && e.mark != 0 {
		return FALSE
	}
	w := adjacent_tile_sup(row, col, DIR_W)
	if w != nil && w.mark != 0 {
		return FALSE
	}
	nw := adjacent_tile_sup(row, col, DIR_NW)
	if nw != nil && nw.mark != 0 {
		return FALSE
	}
	sw := adjacent_tile_sup(row, col, DIR_SW)
	if sw != nil && sw.mark != 0 {
		return FALSE
	}
	ne := adjacent_tile_sup(row, col, DIR_NE)
	if ne != nil && ne.mark != 0 {
		return FALSE
	}
	se := adjacent_tile_sup(row, col, DIR_SE)
	if se != nil && se.mark != 0 {
		return FALSE
	}

	// for every path across the hole, add it to the list of possibilities if it's land-to-land
	// and we haven't already used one of the destination tiles for another hole-crossing.
	var l []int
	if n != nil && s != nil && n.terrain != terr_ocean && s.terrain != terr_ocean && map_[n.row][n.col].mark+map_[s.row][s.col].mark == 0 {
		l = append(l, 1)
	}
	if e != nil && w != nil && e.terrain != terr_ocean && w.terrain != terr_ocean && map_[e.row][e.col].mark+map_[w.row][w.col].mark == 0 {
		l = append(l, 2)
	}
	if ne != nil && sw != nil && ne.terrain != terr_ocean && sw.terrain != terr_ocean && map_[ne.row][ne.col].mark+map_[sw.row][sw.col].mark == 0 {
		l = append(l, 3)
	}
	if se != nil && nw != nil && se.terrain != terr_ocean && nw.terrain != terr_ocean && map_[se.row][se.col].mark+map_[nw.row][nw.col].mark == 0 {
		l = append(l, 4)
	}

	if len(l) == 0 {
		return FALSE
	}

	var name string
	switch _static_bridge_map_hole_sup.road_name_cnt % 3 {
	case 0:
		name = "Secret pass"
	case 1:
		name = "Secret route"
	case 2:
		name = "Old road"
	default:
		panic("!reached")
	}
	_static_bridge_map_hole_sup.road_name_cnt++

	// the horror, the horror
	if n != nil {
		n.mark += rnd(0, 1)
	}
	if e != nil {
		e.mark += rnd(0, 1)
	}
	if w != nil {
		w.mark += rnd(0, 1)
	}
	if s != nil {
		s.mark += rnd(0, 1)
	}
	if nw != nil {
		nw.mark += rnd(0, 1)
	}
	if ne != nil {
		ne.mark += rnd(0, 1)
	}
	if sw != nil {
		sw.mark += rnd(0, 1)
	}
	if se != nil {
		se.mark += rnd(0, 1)
	}

	i := rnd(0, len(l)-1)
	switch l[i] {
	case 1:
		link_roads(n, s, TRUE, name)
		n.mark = 1
		s.mark = 1
	case 2:
		link_roads(e, w, TRUE, name)
		e.mark = 1
		w.mark = 1
	case 3:
		link_roads(ne, sw, TRUE, name)
		ne.mark = 1
		sw.mark = 1
	case 4:
		link_roads(se, nw, TRUE, name)
		se.mark = 1
		nw.mark = 1
	default:
		panic("!reached")
	}

	return l[i]
}

// bridge a # map hole with a secret road.
// do not put two roads in the same square.
func bridge_map_holes() {
	for row := 0; row < max_row; row++ {
		for col := 0; col < max_col; col++ {
			if map_[row][col] == nil {
				if n := bridge_map_hole_sup(row, col); n != 0 {
					log.Printf("%s map hole bridge at (%d,%d)\n", bridge_dir_s[n], row, col)
				}
			}
		}
	}
	log.Println("")
}

func bridge_mountain_ports() {
	// todo: max_row vs MAX_ROW
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].terrain == terr_mountain && is_port_city(row, col) != FALSE && rnd(1, 7) == 7 {
				bridge_mountain_sup(row, col)
			}
		}
	}
}

func bridge_mountain_sup(row, col int) {
	from := map_[row][col]
	if from == nil {
		panic("assert(from != nil)")
	}
	to := adjacent_tile_water(row, col)
	if to == nil {
		panic("assert(to != nil)")
	} else if !(to.terrain == terr_ocean) {
		panic("assert(to.terrain == terr_ocean)")
	}

	var name string
	switch rnd(1, 3) {
	case 1:
		name = "Narrow channel"
	case 2:
		name = "Rocky channel"
	case 3:
		name = "Secret sea route"
	default:
		panic("!reached")
	}

	add_road(from, to.region, TRUE, name)
	add_road(to, from.region, TRUE, name)

	log.Printf("secret sea route at (%2d,%2d)\n", from.row, from.col)
}

func choose_random_stone_circle(l []*tile, avoid1, avoid2 *tile) *tile {
	for {
		i := rnd(0, len(l)-1)
		if l[i] == avoid1 || l[i] == avoid2 {
			continue
		}
		return l[i]
	}
}

func clear_alloc_flag() {
	for i := range alloc_flag {
		alloc_flag[i] = 0
	}
}

func clear_province_marks() {
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil {
				map_[row][col].mark = 0
			}
		}
	}
}

func clear_subloc_marks() {
	for i := 1; i <= top_subloc; i++ {
		subloc[i].mark = 0
	}
}

// count_cities updates the global inside_num_sides array.
func count_cities() {
	for i := 1; i <= top_subloc; i++ {
		if subloc[i].terrain == terr_city {
			inside_num_cities[map_[subloc[i].row][subloc[i].col].inside]++
		}
	}
}

func count_continents() {
	log.Println("")
	log.Printf("Land regions:")
	log.Printf("%-25s  %8s  %6s  %7s  %s\n", "name", "coord", "nprovs", "ncities", "gates (out/in)")
	log.Printf("%-25s  %8s  %6s  %7s  %s\n", "-------------------------", "-----", "------", "-------", "--------------")
	for i := 1; i <= inside_top; i++ {
		if p := inside_list[i][0]; p != nil && p.terrain != terr_ocean {
			print_continent(i)
		}
	}
	log.Println("")
	log.Println("")
	log.Printf("Oceans:")
	log.Printf("%-25s  %8s  %6s  %7s  %s\n", "name", "coord", "nprovs", "ncities", "gates (out/in)")
	log.Printf("%-25s  %8s  %6s  %7s  %s\n", "-------------------------", "-----", "------", "-------", "--------------")
	for i := 1; i <= inside_top; i++ {
		if p := inside_list[i][0]; p != nil && p.terrain == terr_ocean {
			print_continent(i)
		}
	}
	log.Println("")
	log.Println("")
	log.Printf("  %8d continents\n", inside_top)
	log.Printf("  %8d land  locs\n", land_count)
	log.Printf("  %8d water locs\n", water_count)
}

func count_subloc_coverage() {
	clear_province_marks()

	for i := 1; i <= top_subloc; i++ {
		if subloc[i].depth == 3 {
			map_[subloc[i].row][subloc[i].col].mark++
			if map_[subloc[i].row][subloc[i].col].mark >= 5 {
				log.Printf("(%d,%d) has %d sublocs (region %d)\n", subloc[i].row, subloc[i].col, map_[subloc[i].row][subloc[i].col].mark, map_[subloc[i].row][subloc[i].col].region)
			}
		}
	}

	log.Println("")
	log.Printf("subloc coverage:")

	var count [100]int
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].terrain != terr_ocean {
				count[map_[row][col].mark]++
			}
		}
	}
	for i := 99; i >= 0 && count[i] == 0; i-- {
		count[i] = -1
	}
	for i := 0; i < 100 && count[i] != -1; i++ {
		var locHas string
		if count[i] == 1 {
			locHas = "loc has"
		} else {
			locHas = "locs have"
		}
		var s string
		if i == 1 {
			s = " "
		} else {
			s = "s"
		}
		log.Printf("%6d %s %d subloc%s (%d%%)\n", count[i], locHas, i, s, count[i]*100/land_count)
	}
}

func count_sublocs() {
	log.Println("")
	log.Println("subloc counts:")

	clear_province_marks()

	for i := 1; i <= top_subloc; i++ {
		if subloc[i].terrain == terr_island {
			row, col := subloc[i].row, subloc[i].col
			map_[row][col].mark++
		}
	}

	var count [100]int
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].terrain == terr_ocean {
				count[map_[row][col].mark]++
			}
		}
	}

	for i := 99; i >= 0; i-- {
		if count[i] != 0 {
			break
		}
		count[i] = -1
	}

	for i := 0; i < 100; i++ {
		if count[i] == -1 {
			break
		}

		var locHas string
		if count[i] == 1 {
			locHas = "loc has"
		} else {
			locHas = "locs have"
		}
		var s string
		if i == 1 {
			s = " "
		} else {
			s = "s"
		}
		log.Printf("%6d %s %d island%s (%d%%)\n", count[i], locHas, i, s, count[i]*100/water_count)
	}
}

func count_tiles() {
	var count [1000]int
	for r := 0; r < MAX_ROW; r++ {
		for c := 0; c < MAX_COL; c++ {
			if map_[r][c] != nil {
				count[map_[r][c].terrain]++
			}
		}
	}
	for i := 1; i <= top_subloc; i++ {
		count[subloc[i].terrain]++
	}
	for i := 1; i < len(terr_s[i]); i++ {
		fprintf(stderr, "%-30s %d\n", terr_s[i], count[i])
	}
}

func create_a_building(sl int, hidden int, kind int) int {
	top_subloc++
	if !(top_subloc < MAX_SUBLOC) {
		panic("assert(top_subloc < MAX_SUBLOC)")
	}

	subloc[top_subloc] = &tile{}
	subloc[top_subloc].region = rnd_alloc_num(SUBLOC_LOW, SUBLOC_HIGH)
	subloc[top_subloc].inside = subloc[sl].region

	subloc[top_subloc].row = subloc[sl].row
	subloc[top_subloc].col = subloc[sl].col

	subloc[top_subloc].hidden = hidden
	subloc[top_subloc].terrain = kind
	subloc[top_subloc].depth = 4

	subloc[sl].subs = append(subloc[sl].subs, subloc[top_subloc].region)

	return top_subloc
}

func create_a_city(row, col int, name string, major int) int {
	if name == "" {
		name = random_city_name()
	}
	n := create_a_subloc(row, col, 0, terr_city)
	subloc[n].name = name
	subloc[n].major_city = major
	return n
}

func create_a_graveyard(row, col int) {
	n := create_a_subloc(row, col, rnd(0, 1), terr_grave)
	s := name_guild(terr_grave)
	if s != "" {
		subloc[n].name = s
	}
}

func create_a_subloc(row, col int, hidden int, kind int) int {
	top_subloc++
	if !(top_subloc < MAX_SUBLOC) {
		panic("assert(top_subloc < MAX_SUBLOC)")
	}

	subloc[top_subloc] = &tile{}
	if kind == terr_city {
		subloc[top_subloc].region = rnd_alloc_num(CITY_LOW, CITY_HIGH)
	} else {
		subloc[top_subloc].region = rnd_alloc_num(SUBLOC_LOW, SUBLOC_HIGH)
	}
	subloc[top_subloc].inside = map_[row][col].region
	subloc[top_subloc].row = row
	subloc[top_subloc].col = col
	subloc[top_subloc].hidden = hidden
	subloc[top_subloc].terrain = kind
	subloc[top_subloc].depth = 3

	if kind == terr_city {
		map_[row][col].city = 2
	}

	map_[row][col].subs = append(map_[row][col].subs, subloc[top_subloc].region)

	return top_subloc
}

func dir_assert() {
	row, col := 1, 1
	reg := rc_to_region(row, col)
	if !(reg == 10101) {
		panic("assert(rc_to_region(row, col) == 10101)")
	}
	if !(row == region_row(reg)) {
		panic("assert(region_row(reg) == row)")
	}
	if !(col == region_col(reg)) {
		panic("assert(region_col(reg) == col)")
	}

	row, col = 99, 99
	reg = rc_to_region(row, col)
	if !(reg == 19999) {
		panic("assert(rc_to_region(row, col) == 19999)")
	}
	if !(row == region_row(reg)) {
		panic("assert(region_row(reg) == row)")
	}
	if !(col == region_col(reg)) {
		panic("assert(region_col(reg) == col)")
	}

	row, col = 57, 63
	reg = rc_to_region(row, col)
	if !(reg == 15763) {
		panic("assert(rc_to_region(row, col) == 15763)")
	}
	if !(row == region_row(reg)) {
		panic("assert(region_row(reg) == row)")
	}
	if !(col == region_col(reg)) {
		panic("assert(region_col(reg) == col)")
	}
}

func dump_continents(name string) {
	type continent struct {
		Name string `json:"na,omitempty"`
	}
	var data struct {
		Continents []continent `json:"loc-region"`
	}
	for i := 1; i <= inside_top; i++ {
		c := continent{Name: inside_names[i]}
		data.Continents = append(data.Continents, c)
		// todo: print_inside_locs needs to be hoisted into here
	}
	if buf, err := json.MarshalIndent(data, "", "  "); err != nil {
		panic(err)
	} else if err = os.WriteFile(name, buf, 0666); err != nil {
		panic(err)
	}

	for i := 1; i <= inside_top; i++ {
		fprintf(loc_fp, "%d loc region\n", REGION_OFF+i)
		if inside_names[i] != "" {
			fprintf(loc_fp, "na %s\n", inside_names[i])
		}
		print_inside_locs(i)
		fprintf(loc_fp, "\n")
	}
}

func dump_gates(name string) {
	type gate_to struct {
		ToLoc int `json:"tl"`
		Key   int `json:"sk,omitempty"`
	}
	type link_to struct {
		WhereTo int `json:"wh"`
	}
	type gate struct {
		EntNum  int       `json:"ent-num"`
		LinksTo []link_to `json:"LI"`
		GatesTo []gate_to `json:"GA"`
	}
	var data struct {
		Gates []gate `json:"gates"`
	}

	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] == nil {
				continue
			}
			for j := 0; j < len(map_[row][col].gates_dest); j++ {
				data.Gates = append(data.Gates, gate{
					EntNum: map_[row][col].gates_num[j],
					LinksTo: []link_to{
						link_to{WhereTo: map_[row][col].region},
					},
					GatesTo: []gate_to{
						gate_to{
							ToLoc: map_[row][col].gates_dest[j],
							Key:   map_[row][col].gates_key[j],
						},
					},
				})
			}
		}
	}

	for i := 1; i <= top_subloc; i++ {
		for j := 0; j < len(subloc[i].gates_num); j++ {
			fprintf(gate_fp, "%d gate 0\n", subloc[i].gates_num[j])
			fprintf(gate_fp, "LI\n")
			fprintf(gate_fp, " wh %d\n", subloc[i].region)
			fprintf(gate_fp, "GA\n")
			fprintf(gate_fp, " tl %d\n", subloc[i].gates_dest[j])
			if subloc[i].gates_key[j] != 0 {
				fprintf(gate_fp, " sk %d\n", subloc[i].gates_key[j])
			}
			fprintf(gate_fp, "\n")
		}
	}

	if buf, err := json.MarshalIndent(data, "", "  "); err != nil {
		panic(err)
	} else if err = os.WriteFile(name, buf, 0666); err != nil {
		panic(err)
	}

	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil {
				for j := 0; j < len(map_[row][col].gates_dest); j++ {
					fprintf(gate_fp, "%d gate 0\n", map_[row][col].gates_num[j])
					fprintf(gate_fp, "LI\n")
					fprintf(gate_fp, " wh %d\n", map_[row][col].region)
					fprintf(gate_fp, "GA\n")
					fprintf(gate_fp, " tl %d\n", map_[row][col].gates_dest[j])
					if map_[row][col].gates_key[j] != 0 {
						fprintf(gate_fp, " sk %d\n", map_[row][col].gates_key[j])
					}
					fprintf(gate_fp, "\n")
				}
			}
		}
	}

	for i := 1; i <= top_subloc; i++ {
		for j := 0; j < len(subloc[i].gates_num); j++ {
			fprintf(gate_fp, "%d gate 0\n", subloc[i].gates_num[j])
			fprintf(gate_fp, "LI\n")
			fprintf(gate_fp, " wh %d\n", subloc[i].region)
			fprintf(gate_fp, "GA\n")
			fprintf(gate_fp, " tl %d\n", subloc[i].gates_dest[j])
			if subloc[i].gates_key[j] != 0 {
				fprintf(gate_fp, " sk %d\n", subloc[i].gates_key[j])
			}
			fprintf(gate_fp, "\n")
		}
	}
}

func dump_roads(name string) {
	type gate_to struct {
		ToLoc  int `json:"tl"`
		Hidden int `json:"rh,omitempty"`
	}
	type link_to struct {
		WhereTo int `json:"wh"`
	}
	type road struct {
		EntNum  int       `json:"ent-num"`
		Name    string    `json:"na,omitempty"`
		Region  int       `json:"region"`
		LinksTo []link_to `json:"LI"`
		GatesTo []gate_to `json:"GA"`
	}
	var data struct {
		Roads []road `json:"roads"`
	}

	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] == nil {
				continue
			}
			for j := 0; j < len(map_[row][col].roads); j++ {
				data.Roads = append(data.Roads, road{
					EntNum: map_[row][col].roads[j].ent_num,
					Name:   map_[row][col].roads[j].name,
					LinksTo: []link_to{
						link_to{WhereTo: map_[row][col].region},
					},
					GatesTo: []gate_to{
						gate_to{
							ToLoc:  map_[row][col].roads[j].to_loc,
							Hidden: map_[row][col].roads[j].hidden,
						},
					},
				})
			}
		}
	}

	for i := 1; i <= top_subloc; i++ {
		for j := 0; j < len(subloc[i].roads); j++ {
			data.Roads = append(data.Roads, road{
				EntNum: subloc[i].roads[j].ent_num,
				Name:   subloc[i].roads[j].name,
				LinksTo: []link_to{
					link_to{WhereTo: subloc[i].region},
				},
				GatesTo: []gate_to{
					gate_to{
						ToLoc:  subloc[i].roads[j].to_loc,
						Hidden: subloc[i].roads[j].hidden,
					},
				},
			})
		}
	}

	if buf, err := json.MarshalIndent(data, "", "  "); err != nil {
		panic(err)
	} else if err = os.WriteFile(name, buf, 0666); err != nil {
		panic(err)
	}

	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil {
				for j := 0; j < len(map_[row][col].roads); j++ {
					fprintf(road_fp, "%d road 0\n", map_[row][col].roads[j].ent_num)
					if map_[row][col].roads[j].name != "" {
						fprintf(road_fp, "na %s\n", map_[row][col].roads[j].name)
					}
					fprintf(road_fp, "LI\n")
					fprintf(road_fp, " wh %d\n", map_[row][col].region)
					fprintf(road_fp, "GA\n")
					fprintf(road_fp, " tl %d\n", map_[row][col].roads[j].to_loc)
					if map_[row][col].roads[j].hidden != FALSE {
						fprintf(road_fp, " rh %d\n", map_[row][col].roads[j].hidden)
					}
					fprintf(road_fp, "\n")
				}
			}
		}
	}

	for i := 1; i <= top_subloc; i++ {
		for j := 0; j < len(subloc[i].roads); j++ {
			fprintf(road_fp, "%d road 0\n", subloc[i].roads[j].ent_num)
			if subloc[i].roads[j].name != "" {
				fprintf(road_fp, "na %s\n", subloc[i].roads[j].name)
			}
			fprintf(road_fp, "LI\n")
			fprintf(road_fp, " wh %d\n", subloc[i].region)
			fprintf(road_fp, "GA\n")
			fprintf(road_fp, " tl %d\n", subloc[i].roads[j].to_loc)
			if subloc[i].roads[j].hidden != FALSE {
				fprintf(road_fp, " rh %d\n", subloc[i].roads[j].hidden)
			}
			fprintf(road_fp, "\n")
		}
	}
}

func fix_terrain_land() {
	// todo: max_row vs MAX_ROW
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].terrain == terr_land {
				if p := adjacent_tile_terr(row, col); p != nil && p.terrain != terr_land && p.terrain != terr_ocean {
					map_[row][col].terrain = p.terrain
					map_[row][col].color = p.color
				} else {
					log.Printf("fix_terrain: could not infer type of (%d,%d), assuming 'forest'\n", row, col)
					map_[row][col].terrain = terr_forest
				}
			}
		}
	}
}

func flood_land_clumps(row, col int, name string) int {
	map_[row][col].name = name

	count := 0
	for dir := 1; dir < MAX_DIR; dir++ {
		p := adjacent_tile_sup(row, col, dir)
		if p == nil || p.terrain == terr_ocean || p.color == -1 || p.color != map_[row][col].color {
			continue
		} else if p.name == name { // already been here
			continue
		} else if p.name != "" {
			panic(fmt.Sprintf("flood_land_clumps(%d,%d,%q) error)", row, col, name))
		}

		count += flood_land_clumps(p.row, p.col, name)
	}

	return count
}

func flood_land_inside(row, col int, ins int) int {
	count := 0

	map_[row][col].inside = ins
	if map_[row][col].region_boundary != FALSE {
		return count
	}

	for dir := 1; dir < MAX_DIR; dir++ {
		p := adjacent_tile_sup(row, col, dir)
		if p == nil || p.terrain == terr_ocean {
			continue
		} else if p.inside == ins { // already been here
			continue
		} else if p.inside != FALSE {
			panic(fmt.Sprintf("error: flood_land_inside(%d,%d,%q)\n", row, col, inside_names[ins]))
		}
		count += flood_land_inside(p.row, p.col, ins)
	}

	return count
}

func flood_water_inside(row, col int, ins int) int {
	count := 0
	map_[row][col].inside = ins

	for dir := 1; dir < MAX_DIR; dir++ {
		p := adjacent_tile_sup(row, col, dir)
		if p == nil || p.color == -1 || p.color != map_[row][col].color {
			continue
		} else if p.inside == ins { // already been here
			continue
		} else if p.inside != FALSE {
			panic(fmt.Sprintf("error: flood_water_inside(%d,%d,%q)\n", row, col, inside_names[ins]))
		}
		count += flood_water_inside(p.row, p.col, ins)
	}

	return count
}

func gate_continental_tour() {
	log.Println("")
	log.Println("Continental gate tour:")

	l := random_tile_from_each_region()
	m := shift_tour_endpoints(l)

	if !(len(l) == len(m)) {
		panic("assert(len(l) == len(m))")
	}

	i := 0
	for ; i < len(l)-1; i++ {
		log.Printf("\t(%2d,%2d) . (%2d,%2d)\n", l[i].row, l[i].col, m[i+1].row, m[i+1].col)
		new_gate(l[i], m[i+1], 0)
	}
	log.Printf("\t(%2d,%2d) . (%2d,%2d)\n\n", l[i].row, l[i].col, m[0].row, m[0].col)

	new_gate(l[i], m[0], rnd(111, 333))
}

func gate_land_ring(rings int) {
	clear_province_marks()
	mark_bad_locs()

	for j := 1; j <= rings; j++ {
		num := rnd(5, 12)

		var r_first, c_first int
		random_province(&r_first, &c_first, 0)

		r_n, c_n := r_first, c_first

		for i := 1; i < num; i++ {
			var r_next, c_next int
			random_province(&r_next, &c_next, 0)
			new_gate(map_[r_n][c_n], map_[r_next][c_next], 0)
			r_n, c_n = r_next, c_next
		}

		new_gate(map_[r_n][c_n], map_[r_first][c_first], 0)
	}
}

func gate_link_islands(rings int) {
	clear_subloc_marks()

	for j := 1; j <= rings; j++ {
		num := rnd(6, 12)
		first := random_island()
		n := first
		for i := 1; i < num; i++ {
			next := random_island()
			new_gate(subloc[n], subloc[next], 0)
			n = next
		}

		new_gate(subloc[n], subloc[first], 0)
	}
}

func gate_province_islands(times int) {
	clear_province_marks()
	mark_bad_locs()
	clear_subloc_marks()

	for j := 1; j <= times; j++ {
		var r1, c1 int
		random_province(&r1, &c1, 0)
		isle := random_island()

		var r2, c2 int
		random_province(&r2, &c2, 0)

		new_gate(map_[r1][c1], subloc[isle], 0)
		new_gate(subloc[isle], map_[r2][c2], 0)
	}
}

// every region gets a hidden stone circle.
// each stone circle has two gates to other stone circles,
// chosen at random, and five gates to random provinces
func gate_stone_circles() {
	log.Println("")
	log.Println("Ring of stones:")

	var circs []*tile

	l := random_tile_from_each_region()
	for i := 0; i < len(l); i++ {
		n := create_a_subloc(l[i].row, l[i].col, 1, terr_stone_cir)
		circs = append(circs, subloc[n])
		log.Printf("	(%2d,%2d) in %s\n", l[i].row, l[i].col, inside_names[l[i].inside])
	}

	for i := 0; i < len(circs); i++ {
		first := choose_random_stone_circle(circs, circs[i], nil)
		second := choose_random_stone_circle(circs, circs[i], first)

		new_gate(circs[i], first, rnd(111, 333))
		new_gate(circs[i], second, rnd(111, 333))
	}

	clear_province_marks()
	mark_bad_locs()

	for i := 0; i < len(circs); i++ {
		for j := 1; j <= 5; j++ {
			var row, col int
			random_province(&row, &col, 0)
			if rnd(0, 1) != 0 {
				new_gate(circs[i], map_[row][col], 0)
			} else {
				new_gate(circs[i], map_[row][col], rnd(111, 333))
			}
		}
	}
}

func is_port_city(row, col int) int {
	for _, dir := range []int{DIR_N, DIR_S, DIR_E, DIR_W} {
		t := adjacent_tile_sup(row, col, dir)
		if t != nil && t.terrain == terr_ocean {
			return TRUE
		}
	}
	return FALSE
}

func island_allowed(row, col int) bool {
	inside := map_[row][col].inside
	if inside == 0 {
		return true
	}
	for p := 0; p < len(inside_names) && inside_names[p] != ""; p++ {
		if strings.HasPrefix(inside_names[p], "Deep") {
			return false
		}
	}
	return true
}

// If there is a sublocation at an endpoint of the secret road,
// move the road to come from the sublocation instead of the province.
//
// Since only 1/3 of the locations have sublocs, this doesn't happen all the time.
// A very few locations will have the route between two sublocs.
func link_roads(from, to *tile, hidden int, name string) {
	for i := 1; i <= top_subloc; i++ {
		if subloc[i].inside == from.region && subloc[i].terrain != terr_city {
			from = subloc[i]
			break
		}
	}

	for i := 1; i <= top_subloc; i++ {
		if subloc[i].inside == to.region && subloc[i].terrain != terr_city {
			to = subloc[i]
			break
		}
	}

	add_road(from, to.region, hidden, name)
	add_road(to, from.region, hidden, name)
}

func make_appropriate_subloc(row, col int, unused int) {
	sum := 0
	terr := map_[row][col].terrain

	for _, loc := range loc_table {
		if loc.terr == terr {
			sum += loc.weight
		}
	}

	if sum <= 0 {
		log.Printf("no subloc appropriate for (%d,%d)\n", row, col)
		return
	}

	n := rnd(1, sum)
	for _, loc := range loc_table {
		if loc.terr == terr {
			n -= loc.weight
			if n <= 0 {
				if loc.kind < 0 {
					break
				}
				var hid int
				if loc.hidden == 2 {
					hid = rnd(0, 1)
				} else {
					hid = loc.hidden
				}

				n = create_a_subloc(row, col, hid, loc.kind)
				s := name_guild(loc.kind)
				if s != "" {
					subloc[n].name = s
				}
				break
			}
		}
	}
}

// Gate laying plan
//
//	province => island   => province  (gate_province_islands)
//	province             => province  (random_province_gates)
//	continental tour                  (gate_continental_tour)
//	each region gets a stone ring     (gate_stone_circles)
//	   each with links to 2 other
//	   rings, and 5 random provinces
//	 7 rings of 5-12 provinces        (gate_land_ring)
//	12 rings of 6-12 islands          (gate_link_islands)
func make_gates() {
	if GATES_OTHER != FALSE {
		gate_province_islands(GATE_TIMES)
		random_province_gates(GATE_TIMES)
	}
	if GATES_CONTINENTAL_TOUR != FALSE {
		gate_continental_tour()
	}
	if GATES_STONE_CIRCLES != FALSE {
		gate_stone_circles()
	}
	if GATES_OTHER != FALSE {
		// todo: 5 random provinces in comments above?
		gate_land_ring(2)    // VLN: was 5
		gate_link_islands(1) // VLN: disjoint
		// gate_link_islands(9)  // disjoint
		// gate_link_islands(3)  // might cross
	}
	show_gate_coverage()
}

func make_graveyards() {
	for i := 1; i <= inside_top; i++ {
		if inside_list[i][0].terrain == terr_ocean {
			continue
		}
		n := len(inside_list[i])
		if n < 10 {
			continue
		}
		l := shuffle_tiles(inside_list[i])
		for j := 0; j < n/10; j++ {
			create_a_graveyard(l[j].row, l[j].col)
		}
	}
}

func make_islands() {
	// gather all allowable provinces
	var available []*tile
	for row := 0; row < max_row; row++ {
		for col := 0; col < max_col; col++ {
			if map_[row][col] != nil {
				if map_[row][col].terrain == terr_ocean {
					if island_allowed(row, col) {
						available = append(available, map_[row][col])
					}
				}
			}
		}
	}
	rand.Shuffle(len(available), func(i, j int) {
		available[i], available[j] = available[j], available[i]
	})

	for i := 3; i < 100; i++ {
		num_islands = water_count / i
		if num_islands < len(available) {
			break
		}
	}
	log.Printf("make_islands: available provinces %d, water_count %d, num_islands %d\n", len(available), water_count, num_islands)

	if !(num_islands < len(available)) { // too many islands, not enough room
		panic("assert(num_islands < len(available)")
	}

	for i := 1; i <= num_islands; i++ {
		row, col := available[i].row, available[i].col
		create_a_subloc(row, col, rnd(0, 1), terr_island)
	}
}

func make_roads() {
	clear_province_marks()
	bridge_map_holes()
	bridge_caddy_corners()
	bridge_mountain_ports()
}

func map_init() {
	for i := range map_ {
		for j := range map_[i] {
			map_[i][j] = nil
		}
	}
}

func mark_bad_locs() {
	for i := 1; i <= inside_top; i++ {
		if inside_names[i] == "Impassable Mountains" {
			for j := 0; j < len(inside_list[i]); j++ {
				inside_list[i][j].mark = 1
			}
		}
	}
}

func name_guild(skill int) string {
	sum := 0
	for _, guild := range guild_names {
		if guild.skill == skill {
			sum += guild.weight
		}
	}
	if !(sum > 0) {
		panic("assert(sum > 0)")
	}

	n := rnd(1, sum)
	for _, guild := range guild_names {
		if guild.skill == skill {
			n -= guild.weight
			if n <= 0 {
				return guild.name
			}
		}
	}

	panic("!reached")
}

func new_gate(from *tile, to *tile, key int) {
	gate_num := rnd_alloc_num(SUBLOC_LOW, SUBLOC_HIGH)

	from.gates_num = append(from.gates_num, gate_num)
	from.gates_dest = append(from.gates_dest, to.region)
	from.gates_key = append(from.gates_key, key)

	// gather statistics
	inside_gates_from[map_[from.row][from.col].inside]++
	inside_gates_to[map_[to.row][to.col].inside]++
}

func not_place_random_subloc(kind int, hidden int) int {
	var row, col int
	not_random_province(&row, &col)
	return create_a_subloc(row, col, hidden, kind)
}

// the 'not' refers to not desert and not swamp (and ocean, too?).
// we don't want to make any cities in deserts or swamps (or oceans,either?).
// (well, except for the lost city and the city of the ancients.)
func not_random_province(row, col *int) { // oh, hack upon hack
	sum := 0 // total number of eligible provinces
	for r := 0; r <= max_row; r++ {
		for c := 0; c < max_col; c++ {
			if map_[r][c] != nil {
				if map_[r][c].terrain != terr_ocean && map_[r][c].terrain != terr_swamp && map_[r][c].terrain != terr_desert && map_[r][c].mark == 0 {
					sum++
				}
			}
		}
	}
	if sum == 0 { // there are no provinces to pick
		panic("assert(sum != 0)")
	}

	n := rnd(1, sum) // randomly pick one of those provinces

	for r := 0; r <= max_row; r++ {
		for c := 0; c < max_col; c++ {
			if map_[r][c] != nil {
				if map_[r][c].terrain != terr_ocean && map_[r][c].terrain != terr_swamp && map_[r][c].terrain != terr_desert && map_[r][c].mark == 0 {
					n--
					if n <= 0 { // this picks that nth province
						*row, *col = r, c
						map_[r][c].mark = TRUE
						return
					}
				}
			}
		}
	}

	// should never fall through to here
	panic("!reached")
}

func open_fps() {
	var err error
	if loc_fp, err = fopen("loc", "w"); err != nil {
		panic(err)
	} else if gate_fp, err = fopen("gate", "w"); err != nil {
		panic(err)
	} else if road_fp, err = fopen("road", "w"); err != nil {
		panic(err)
	}
}

func place_random_subloc(kind int, hidden int, terr int) int {
	var row, col int
	random_province(&row, &col, terr)
	return create_a_subloc(row, col, hidden, kind)
}

func place_sublocations() {
	var l []int
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].terrain != terr_ocean {
				l = append(l, row*1000+col)
			}
		}
	}
	l = shuffle_ints(l)

	for i := 0; i < len(l); i++ {
		row, col := l[i]/1000, l[i]%1000

		// put a city everywhere there is a * or every 1 in 15 locs, randomly.
		// don't put one where there already is a city (city != 2).
		if map_[row][col].city == 1 || (rnd(1, 15) == 1 && map_[row][col].city != 2) {
			create_a_city(row, col, "", FALSE)
		}
		if rnd(1, 100) <= 35 {
			make_appropriate_subloc(row, col, 0)
		}
		if rnd(1, 100) <= 35 {
			make_appropriate_subloc(row, col, 0)
		}
		if rnd(1, 100) <= 35 {
			make_appropriate_subloc(row, col, 0)
		}
	}
}

func print_continent(i int) {
	p := inside_list[i][0]
	name := inside_names[i]
	if name == "" {
		name = fmt.Sprintf("?? (%2d,%2d)", p.row, p.col)
	}
	coord := fmt.Sprintf("(%2d,%2d)", p.row, p.col)
	nprovs := fmt.Sprintf("%d", len(inside_list[i]))
	ncities := fmt.Sprintf("%d", inside_num_cities[i])
	gates := fmt.Sprintf("%d/%d", inside_gates_from[i], inside_gates_to[i])
	log.Printf("%-25s  %8s  %6s  %7s  %s\n", name, coord, nprovs, ncities, gates)
}

func print_inside_locs(n int) {
	count := 0
	for i := 0; i < len(inside_list[n]); i++ {
		count++
		if count == 1 {
			fprintf(loc_fp, "LI\n")
			fprintf(loc_fp, " hl ")
		}
		if count%11 == 10 { // continuation line
			fprintf(loc_fp, "\\\n\t")
		}
		fprintf(loc_fp, "%d ", inside_list[n][i].region)
	}
	if count != 0 {
		fprintf(loc_fp, "\n")
	}
}

func print_map(fp *io.FILE) {
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] == nil { // hole in map?
				continue
			}
			flag, sl := TRUE, FALSE

			fprintf(loc_fp, "%d loc %s\n", map_[row][col].region,
				terr_s[map_[row][col].terrain])

			if map_[row][col].name != "" && map_[row][col].name != "Unnamed" {
				fprintf(loc_fp, "na %s\n", map_[row][col].name)
			}

			if map_[row][col].uldim_flag != 0 {
				if sl == FALSE {
					fprintf(loc_fp, "SL\n")
					sl = TRUE
				}
				fprintf(loc_fp, " uf %d\n", map_[row][col].uldim_flag)
			}

			if map_[row][col].summerbridge_flag != 0 {
				if sl == FALSE {
					fprintf(loc_fp, "SL\n")
					sl = TRUE
				}
				fprintf(loc_fp, " sf %d\n", map_[row][col].summerbridge_flag)
			}

			if map_[row][col].safe_haven != FALSE {
				if sl == FALSE {
					fprintf(loc_fp, "SL\n")
					sl = TRUE
				}
				fprintf(loc_fp, " sh 1\n")
			}

			if map_[row][col].inside != 0 {
				fprintf(loc_fp, "LI\n")
				flag = FALSE

				if map_[row][col].inside != 0 {
					fprintf(loc_fp, " wh %d\n", map_[row][col].inside+REGION_OFF)
				}
			}

			print_inside_sublocs(flag, row, col)

			fprintf(loc_fp, "LO\n")
			fprintf(loc_fp, " pd %d %d %d %d\n", prov_dest(map_[row][col], DIR_N), prov_dest(map_[row][col], DIR_E), prov_dest(map_[row][col], DIR_S), prov_dest(map_[row][col], DIR_W))

			if map_[row][col].hidden != FALSE {
				fprintf(loc_fp, " hi %d\n", map_[row][col].hidden)
			}
			/* untested */
			if map_[row][col].sea_lane != FALSE {
				fprintf(loc_fp, " sl 1\n")
			}

			fprintf(loc_fp, "\n")
		}
	}
}

func print_inside_sublocs(flag int, row, col int) {
	count := 0
	for i := 0; i < len(map_[row][col].roads); i++ {
		count++
		if count == 1 {
			if flag != 0 {
				fprintf(loc_fp, "LI\n")
			}
			fprintf(loc_fp, " hl ")
		}

		if count%11 == 10 { // continuation line
			fprintf(loc_fp, "\\\n\t")
		}

		fprintf(loc_fp, "%d ", map_[row][col].roads[i].ent_num)
	}

	for i := 0; i < len(map_[row][col].gates_num); i++ {
		count++
		if count == 1 {
			if flag != 0 {
				fprintf(loc_fp, "LI\n")
			}
			fprintf(loc_fp, " hl ")
		}

		if count%11 == 10 { // continuation line
			fprintf(loc_fp, "\\\n\t")
		}

		fprintf(loc_fp, "%d ", map_[row][col].gates_num[i])
	}

	for i := 0; i < len(map_[row][col].subs); i++ {
		count++
		if count == 1 {
			if flag != 0 {
				fprintf(loc_fp, "LI\n")
			}
			fprintf(loc_fp, " hl ")
		}

		if count%11 == 10 { // continuation line
			fprintf(loc_fp, "\\\n\t")
		}

		fprintf(loc_fp, "%d ", map_[row][col].subs[i])
	}

	if count != 0 {
		fprintf(loc_fp, "\n")
	}
}

func print_subloc_gates(n int) { // and inside buildings...
	count := 0

	for i := 0; i < len(subloc[n].roads); i++ {
		count++
		if count == 1 {
			fprintf(loc_fp, " hl ")
		}
		if count%11 == 10 { // continuation line
			fprintf(loc_fp, "\\\n\t")
		}
		fprintf(loc_fp, "%d ", subloc[n].roads[i].ent_num)
	}

	for i := 0; i < len(subloc[n].gates_num); i++ {
		count++
		if count == 1 {
			fprintf(loc_fp, " hl ")
		}
		if count%11 == 10 { // continuation line
			fprintf(loc_fp, "\\\n\t")
		}
		fprintf(loc_fp, "%d ", subloc[n].gates_num[i])
	}

	for i := 0; i < len(subloc[n].subs); i++ {
		count++
		if count == 1 {
			fprintf(loc_fp, " hl ")
		}
		if count%11 == 10 { // continuation line
			fprintf(loc_fp, "\\\n\t")
		}
		fprintf(loc_fp, "%d ", subloc[n].subs[i])
	}

	if count != 0 {
		fprintf(loc_fp, "\n")
	}
}

func print_sublocs() {
	for i := 1; i <= top_subloc; i++ {
		sl := FALSE

		fprintf(loc_fp, "%d loc %s\n", subloc[i].region, terr_s[subloc[i].terrain])

		if subloc[i].name != "" && subloc[i].name != "Unnamed" {
			fprintf(loc_fp, "na %s\n", subloc[i].name)
		}

		if !(subloc[i].inside != 0) {
			panic("assert(subloc[i].inside != 0)")
		}

		fprintf(loc_fp, "LI\n")
		fprintf(loc_fp, " wh %d\n", subloc[i].inside)
		print_subloc_gates(i)

		fprintf(loc_fp, "LO\n")

		if subloc[i].hidden != FALSE {
			fprintf(loc_fp, " hi %d\n", subloc[i].hidden)
		}

		if subloc[i].major_city != FALSE {
			if sl == FALSE {
				fprintf(loc_fp, "SL\n")
				sl = TRUE
			}
			fprintf(loc_fp, " mc %d\n", subloc[i].major_city)
		}

		if subloc[i].safe_haven != FALSE {
			if sl == FALSE {
				fprintf(loc_fp, "SL\n")
				sl = TRUE
			}
			fprintf(loc_fp, " sh 1\n")
		}

		fprintf(loc_fp, "\n")
	}
}

// Return the region immediately adjacent to <location> in direction <dir>.
// Returns 0 if there is no adjacent location in the given direction.
func prov_dest(t *tile, dir int) int {
	row, col := t.row, t.col

	switch dir {
	case DIR_N:
		row--
	case DIR_E:
		col++
	case DIR_S:
		row++
	case DIR_W:
		col--
	default:
		panic(fmt.Sprintf("assert(dir != %d)", dir))
	}

	// this way wraps just E-W
	//  TODO: should this use max_row, MAX_ROW?
	//  if row < 0 || row > 99 {
	//	  return 0
	//  }
	//  if col < 0 {
	//	  col = 99
	//  } else if col > 99 {
	//	  col = 0
	//  }

	// this way wraps both N-S and E-W
	if row < 0 {
		row = max_row
	} else if row > max_row {
		row = 0
	}
	if col < 0 {
		col = max_col
	} else if col > max_col {
		col = 0
	}

	if map_[row][col] == nil {
		return 0
	}
	return map_[row][col].region
}

var _static_random_city_name struct {
	cities []string
}

func random_city_name() string {
	if _static_random_city_name.cities == nil {
		var data struct {
			Cities []string `json:"cities"`
		}
		buf, err := os.ReadFile("city-names.json")
		if err != nil {
			log.Printf("random_city_name: %+v\n", err)
		} else if err = json.Unmarshal(buf, &data); err != nil {
			log.Printf("random_city_name: %+v\n", err)
		} else {
			_static_random_city_name.cities = data.Cities
		}
		log.Printf("%q: loaded %d city names\n", "city-names.json", len(_static_random_city_name.cities))
		if _static_random_city_name.cities == nil {
			_static_random_city_name.cities = []string{"T'othville"}
		}
	}

	var name string
	if len(_static_random_city_name.cities) != 0 {
		name = _static_random_city_name.cities[0]
		_static_random_city_name.cities = _static_random_city_name.cities[1:]
	} else {
		as := []string{"a", "a", "a", "ai", "au", "a'i", "a'u", "e", "e", "e", "i", "i", "i", "i'i", "o", "o", "o", "o'a", "u", "u", "u'tu"}
		cs := []string{"b", "bh", "c", "ch", "d", "f", "g", "gh", "h", "je", "k", "ka'", "ke'", "l", "ll", "m", "n", "p", "r", "s", "t", "w"}
		if rnd(0, 10) == 0 {
			name = as[rnd(0, len(as)-1)]
		}
		for i := rnd(1, 3) + rnd(1, 3) + rnd(1, 3); i >= 0; i-- {
			name = name + cs[rnd(0, len(cs)-1)] + as[rnd(0, len(as)-1)]
		}
		name = strings.ToUpper(string(name[0])) + name[1:]
	}

	return name
}

func random_island() int {
	i := 0
	for { // todo: understand what this is doing
		n := rnd(1, num_islands)
		for i = 1; i <= top_subloc; i++ {
			if subloc[i].terrain == terr_island {
				n--
				if n <= 0 {
					break
				}
			}
		}
		if !(n <= top_subloc) {
			panic("assert(n <= top_subloc);")
		}
		if !(subloc[i].mark != FALSE) {
			break
		}
	}
	subloc[i].mark = TRUE
	return i
}

func random_province(row, col *int, terr int) {
	sum := 0 // number of eligible provinces

	if terr == 0 {
		for r := 0; r <= max_row; r++ {
			for c := 0; c < max_col; c++ {
				if map_[r][c] != nil && map_[r][c].terrain != terr_ocean && map_[r][c].mark == FALSE {
					sum++
				}
			}
		}
	} else {
		for r := 0; r <= max_row; r++ {
			for c := 0; c < max_col; c++ {
				if map_[r][c] != nil && map_[r][c].terrain == terr && map_[r][c].mark == FALSE {
					sum++
				}
			}
		}
	}
	if sum == 0 { // nothing available
		panic("assert(sum != 0)")
	}

	n := rnd(1, sum) // pick one of those provinces at random
	if terr == 0 {
		for r := 0; r <= max_row; r++ {
			for c := 0; c < max_col; c++ {
				if map_[r][c] != nil {
					if map_[r][c].terrain == terr_ocean {
						if map_[r][c].mark == FALSE {
							n--
							if n <= 0 {
								*row, *col = r, c
								map_[r][c].mark = TRUE
								return
							}
						}
					}
				}
			}
		}
	} else {
		for r := 0; r <= max_row; r++ {
			for c := 0; c < max_col; c++ {
				if map_[r][c] != nil {
					if map_[r][c].terrain == terr {
						if map_[r][c].mark == FALSE {
							n--
							if n <= 0 {
								*row, *col = r, c
								map_[r][c].mark = TRUE
								return
							}
						}
					}
				}
			}
		}
	}

	// should never reach here
	panic("!reached")
}

func random_province_gates(n int) {
	clear_province_marks()
	mark_bad_locs()

	for i := 0; i < n; i++ {
		var r1, c1 int
		random_province(&r1, &c1, 0)

		var r2, c2 int
		random_province(&r2, &c2, 0)

		// todo: should this be r1,c1 and r2,c2?
		new_gate(map_[r1][c1], map_[r1][c2], 0)
	}
}

func random_tile_from_each_region() []*tile {
	var l []*tile

	for i := 1; i <= inside_top; i++ {
		if inside_list[i][0].terrain == terr_ocean {
			continue
		} else if inside_names[i] == "Impassable Mountains" {
			continue
		}

		j := rnd(0, len(inside_list[i])-1)

		l = append(l, inside_list[i][j])
	}

	return shuffle_tiles(l)
}

func randomize_dir_vector() {
	snap := []int{DIR_N, DIR_NE, DIR_E, DIR_SE, DIR_S, DIR_SW, DIR_W, DIR_NW}
	rand.Shuffle(len(snap), func(i, j int) {
		snap[i], snap[j] = snap[j], snap[i]
	})
	dir_vector[0] = 0
	for i := 1; i < MAX_DIR; i++ {
		dir_vector[i] = snap[i-1]
	}
	return
}

func rc_to_region(row, col int) int {
	assert(0 <= row && row < MAX_ROW)
	assert(0 <= col && col < MAX_COL)
	return 10_000 + (row * 100) + col
}

func read_map(name string) {
	lines, err := io.ReadLines(name)
	if err != nil {
		panic(err)
	}

	for row, line := range lines {
		for col := 0; col < len(line) && line[col] != '\n'; col++ {
			if line[col] == '#' { // hole in map
				// todo: should this be after the check for max row and max col?
				continue
			}

			if row > max_row {
				max_row = row
			}
			if col > max_col {
				max_col = col
			}

			// todo: shouldn't we check for MAX_ROW and MAX_COL?
			map_[row][col] = &tile{
				row:    row,
				col:    col,
				region: rc_to_region(row, col),
				depth:  2,
			}

			color, terrain := 0, 0

			switch line[col] {
			case ';':
				map_[row][col].sea_lane = TRUE
				terrain = terr_ocean
				color = 1
			case ',':
				terrain = terr_ocean
				color = 1
			case ':':
				map_[row][col].sea_lane = TRUE
				terrain = terr_ocean
				color = 2
			case '.':
				terrain = terr_ocean
				color = 2
			case '~':
				map_[row][col].sea_lane = TRUE
				terrain = terr_ocean
				color = 3
			case ' ':
				terrain = terr_ocean
				color = 3
			case '"':
				map_[row][col].sea_lane = TRUE
				terrain = terr_ocean
				color = 4
			case '\'':
				terrain = terr_ocean
				color = 4
			case 'p':
				color = 5
				terrain = terr_plain
			case 'P':
				color = 6
				terrain = terr_plain
			case 'd':
				color = 7
				terrain = terr_desert
			case 'D':
				color = 8
				terrain = terr_desert
			case 'm':
				color = 9
				terrain = terr_mountain
			case 'M':
				color = 10
				terrain = terr_mountain
			case 's':
				color = 11
				terrain = terr_swamp
			case 'S':
				color = 12
				terrain = terr_swamp
			case 'f':
				color = 13
				terrain = terr_forest
			case 'F':
				color = 14
				terrain = terr_forest
			case 'o':
				color = -1
				switch rnd(1, 10) {
				case 1, 2, 3:
					terrain = terr_forest
				case 4, 5, 6:
					terrain = terr_plain
				case 7, 8:
					terrain = terr_mountain
				case 9:
					terrain = terr_swamp
				case 10:
					terrain = terr_desert
				}
			//case '?':
			//	map_[row][col].hidden = TRUE;
			//	terrain = terr_land;

			// special stuff

			case '^':
				color = 9 /* was 15, unique */
				terrain = terr_mountain
				map_[row][col].uldim_flag = 1
				map_[row][col].region_boundary = TRUE
			case 'v':
				color = 9 /* was 15, unique */
				terrain = terr_mountain
				map_[row][col].uldim_flag = 2
				map_[row][col].region_boundary = TRUE
			case '{':
				color = 16
				terrain = terr_mountain
				map_[row][col].uldim_flag = 3
				map_[row][col].name = "Uldim pass"
				map_[row][col].region_boundary = TRUE
			case '}':
				color = 16
				terrain = terr_mountain
				map_[row][col].uldim_flag = 4
				map_[row][col].name = "Uldim pass"
				map_[row][col].region_boundary = TRUE
			case ']':
				terrain = terr_swamp
				map_[row][col].summerbridge_flag = 1
				map_[row][col].name = "Summerbridge"
				map_[row][col].region_boundary = TRUE
			case '[':
				terrain = terr_swamp
				map_[row][col].summerbridge_flag = 2
				map_[row][col].name = "Summerbridge"
				map_[row][col].region_boundary = TRUE
			case 'O':
				terrain = terr_mountain
				color = -1
				map_[row][col].name = "Mt. Olympus"
			case '1':
				terrain = terr_forest
				color = 19
				map_[row][col].safe_haven = TRUE
				n := create_a_city(row, col, "Drassa", TRUE)
				subloc[n].safe_haven = TRUE
				log.Printf("Start city #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '2':
				terrain = terr_forest
				color = 19
				map_[row][col].safe_haven = TRUE
				n := create_a_city(row, col, "Rimmon", TRUE)
				subloc[n].safe_haven = TRUE
				log.Printf("Start city #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '3':
				terrain = terr_forest
				color = 19
				map_[row][col].safe_haven = TRUE
				n := create_a_city(row, col, "Harn", TRUE)
				subloc[n].safe_haven = TRUE
				log.Printf("Start city #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '4':
				terrain = terr_forest
				color = 19
				map_[row][col].safe_haven = TRUE
				n := create_a_city(row, col, "Imperial City", TRUE)
				subloc[n].safe_haven = TRUE
				log.Printf("Imperical City #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '5':
				terrain = terr_forest
				color = 19
				map_[row][col].safe_haven = TRUE
				n := create_a_city(row, col, "Port Aurnos", TRUE)
				subloc[n].safe_haven = TRUE
				log.Printf("Start city #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '6':
				terrain = terr_forest
				color = 19
				map_[row][col].safe_haven = TRUE
				n := create_a_city(row, col, "Greyfell", TRUE)
				subloc[n].safe_haven = TRUE
				log.Printf("Start city #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '7':
				terrain = terr_forest
				color = 19
				map_[row][col].safe_haven = TRUE
				n := create_a_city(row, col, "Yellowleaf", TRUE)
				subloc[n].safe_haven = TRUE
				log.Printf("Start city #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '8':
				terrain = terr_forest
				color = 19
				n := create_a_city(row, col, "Golden City", TRUE)
				log.Printf("Golden City #%c %s at (%d,%d)\n", line[col], subloc[n].name, row, col)
			case '*':
				terrain = terr_land
				create_a_city(row, col, "", TRUE)
			case '%':
				terrain = terr_land
				create_a_city(row, col, "", FALSE)
			default:
				if isdigit(line[col]) {
					panic(fmt.Sprintf("terrain %q should not fall through", line[col]))
				}
				panic(fmt.Sprintf("unknown terrain %q", line[col]))
			}

			map_[row][col].save_char = line[col]
			map_[row][col].terrain = terrain
			map_[row][col].color = color

			if terrain == terr_water || terrain == terr_ocean {
				water_count++
			} else {
				land_count++
			}
		}
	}

	log.Printf("map_init: rows %d, cols %d\n", max_row, max_col)
}

// The entity number of a region determines where it is on the map.
// Here is how:
//
//   (r,c)
// 	+-------------------+
// 	|(1,1)        (1,99)|
// 	|                   |
// 	|                   |
// 	|(n,1)        (n,99)|
// 	+-------------------+
//
// Entity [10101] corresponds to (1, 1).
// Entity [10199] corresponds to (1,99).
//
// Note that for player convenience an alternate method of representing
// location entity numbers may be used, i.e. 'aa'. 101, 'ab' . 102,
// so [aa01] . [10101], [ab53] . [10253].

func region_col(where int) int {
	return where % 100
}

func region_row(where int) int {
	return (where / 100) % 100
}

func region_row_col(where int) (int, int) {
	return (where / 100) % 100, where % 100
}

// rnd_alloc_num allocates a number in the range low...high.
// it panics if it can't find an available number in that range.
func rnd_alloc_num(low, high int) int {
	n := rnd(low, high)

	for i := n; i <= high; i++ {
		if alloc_flag[i] == 0 {
			alloc_flag[i] = 1
			return i
		}
	}

	for i := low; i < n; i++ {
		if alloc_flag[i] == 0 {
			alloc_flag[i] = 1
			return i
		}
	}

	panic(fmt.Sprintf("rnd_alloc_num(%d, %d) failed", low, high))
}

// name groups of provinces
func set_province_clumps(name string) {
	var lands []struct {
		Row  int    `json:"row"`
		Col  int    `json:"col"`
		Kind string `json:"kind"`
		Name string `json:"name"`
	}
	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &lands)
	if err != nil {
		panic(err)
	}
	log.Printf("set_province_clumps: loaded %d areas\n", len(lands))

	count := 0
	for _, land := range lands {
		if land.Kind == "" {
			land.Kind = " "
		}
		if map_[land.Row][land.Col] == nil {
			log.Printf("set_province_clumps: error: map_[%d][%d] == nil\n", land.Row, land.Col)
		} else if map_[land.Row][land.Col].name != "" {
			log.Printf("set_province_clumps: error: clump collision between %q and %q at (%d,%d)\n", land.Name, map_[land.Row][land.Col].name, land.Row, land.Col)
		} else if map_[land.Row][land.Col].save_char != byte(land.Kind[0]) {
			log.Printf("set_province_clumps: error: land %q expects '%c' at (%d,%d), got '%c'\n", land.Name, land.Kind, land.Row, land.Col, map_[land.Row][land.Col].save_char)
		} else {
			flood_land_clumps(land.Row, land.Col, land.Name)
			count++
		}
	}

	log.Printf("set_province_clumps: named %d areas\n", count)
}

func set_regions(name string) {
	var regions []struct {
		Row  int    `json:"row"`
		Col  int    `json:"col"`
		Name string `json:"name"`
	}
	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &regions)
	if err != nil {
		panic(err)
	}
	log.Printf("set_regions_clumps: loaded %d regions\n", len(regions))

	sr_land_count, sr_water_count := 0, 0

	for _, region := range regions {
		if map_[region.Row][region.Col] == nil {
			panic(fmt.Sprintf("assert(map_[%d][%d] != nil)", region.Row, region.Col))
		}
		if map_[region.Row][region.Col].inside != FALSE {
			log.Printf("set_regions: collision between %q and %q at (%d,%d)\n", region.Name, inside_names[map_[region.Row][region.Col].inside], region.Row, region.Col)
			continue
		}

		ins := alloc_inside()
		inside_names[ins] = region.Name
		if map_[region.Row][region.Col].terrain == terr_ocean {
			sr_water_count++
			flood_water_inside(region.Row, region.Col, ins)
		} else {
			sr_land_count++
			flood_land_inside(region.Row, region.Col, ins)
		}
	}

	log.Printf("set_regions: named %d land regions, %d water regions\n", sr_land_count, sr_water_count)

	// locate unnamed regions
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].inside == FALSE {
				ins := alloc_inside()
				if map_[row][col].terrain == terr_ocean {
					n := flood_water_inside(row, col, ins)
					log.Printf("\tunnamed sea at  %d,%d (%d locs)\n", row, col, n)
				} else {
					n := flood_land_inside(row, col, ins)
					log.Printf("\tunnamed land at %d,%d (%d locs)\n", row, col, n)
				}
			}
		}
	}

	// collect the list of provinces in each region
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] != nil && map_[row][col].inside != FALSE {
				inside_list[map_[row][col].inside] = append(inside_list[map_[row][col].inside], map_[row][col])
			}
		}
	}
}

func shift_tour_endpoints(l []*tile) []*tile {
	var other []*tile
	for i := 0; i < len(l); i++ {
		p := adjacent_tile_terr(l[i].row, l[i].col)
		if p == nil {
			p = l[i]
		}

		q := adjacent_tile_terr(p.row, p.col)
		if q == l[i] { // doubled back, retry
			q = adjacent_tile_terr(p.row, p.col)
		}

		if q == nil {
			log.Printf("shift_tour_endpoints: couldn't shift tour (%d,%d): no such q\n", l[i].row, l[i].col)
			other = append(other, l[i])
		} else if q.terrain == terr_ocean {
			log.Printf("shift_tour_endpoints: couldn't shift tour (%d,%d): ocean\n", l[i].row, l[i].col)
			other = append(other, l[i])
		} else {
			other = append(other, q)
		}
	}

	return other
}

func show_gate_coverage() {
	log.Println("")
	log.Println("Gate coverage:  (in/out)")

	for i := 1; i <= inside_top; i++ {
		if inside_list[i][0].terrain == terr_ocean {
			continue
		}
		log.Printf("\t%d/%d\t%s\n", inside_gates_to[i], inside_gates_from[i], inside_names[i])
	}
}

func shuffle_ints(i []int) (l []int) {
	l = append(l, i...)
	rand.Shuffle(len(l), func(i, j int) {
		l[i], l[j] = l[j], l[i]
	})
	return l
}

func shuffle_tiles(t []*tile) (l []*tile) {
	l = append(l, t...)
	rand.Shuffle(len(l), func(i, j int) {
		l[i], l[j] = l[j], l[i]
	})
	return l
}

func unnamed_province_clumps() {
	log.Println("")
	log.Println("")
	log.Println("Unnamed areas")
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] == nil {
				continue
			} else if map_[row][col].name != "" {
				continue
			} else if map_[row][col].terrain == terr_ocean {
				continue
			}
			n := flood_land_clumps(row, col, "Unnamed")
			if map_[row][col].save_char != 'o' {
				log.Printf("\t%2d,%2d\t%c\t%d unnamed\n", row, col, map_[row][col].save_char, n)
			}
		}
	}
	log.Println("")
}
