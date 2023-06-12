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

import "fmt"

type exit_view struct {
	direction     int // which direction does the exit go
	destination   int // where the exit goes
	orig          int // loc we're coming from
	distance      int // how far, in days
	impassable    int // set if not possible to go there
	dest_hidden   int // set if destination hidden
	orig_hidden   int // set if origination or road is hidden
	hidden        int // set if hidden destination unknown to us
	inside        int // different region destinion is in
	road          int // road entity number, if this is a road
	water         int // is a water link
	in_transit    int // is link to a ship that is moving?
	magic_barrier int // a magical barrier prevents travel
	hades_cost    int // Gate Spirit of Hades fee to enter
	seize         int // Whether we're seizing in an attack.
	forced_march  int // A forced ride?
}

const (
	RAND  = 1
	LAND  = 1
	WATER = 2
)

var (
	max_map_row  = 0
	max_map_col  = 0
	max_map_init = false
)

func DIR_NSEW(a int) bool { return ((a) >= MG_DIR_N && (a) <= MG_DIR_W) }

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

func rc_to_region(row, col int) int {
	assert(0 <= row && row < MAX_ROW)
	assert(0 <= col && col < MAX_COL)
	return 10_000 + (row * 100) + col
}

func region_col(where int) int {
	return where % 100
}

func region_row(where int) int {
	return (where / 100) % 100
}

func region_row_col(where int) (int, int) {
	return (where / 100) % 100, where % 100
}

func determine_map_edges() {
	var i int
	var row, col int

	for _, i = range loop_province() {
		if region(i) == faery_region || region(i) == hades_region || region(i) == cloud_region {
			continue
		}

		row = region_row(i)
		col = region_col(i)

		if row > max_map_row {
			max_map_row = row
		}

		if col > max_map_col {
			max_map_col = col
		}
	}

	max_map_init = true
}

func los_province_distance(a, b int) int {
	var ra, rb, r1, r2 int
	var ca, cb, c1, c2 int
	var d1, d2 int

	a = province(a)
	b = province(b)

	ra = region_row(a) - 1
	rb = region_row(b) - 1
	ca = region_col(a) - 1
	cb = region_col(b) - 1

	r1 = min(ra, rb)
	r2 = max(ra, rb)
	c1 = min(ca, cb)
	c2 = max(ca, cb)

	/* no wrap N-S */
	/*
	 *  Mon Apr 10 10:06:53 2000 -- Scott Turner
	 *
	 *  Well, we do have N-S wrapping now.  How does that
	 *  change this?
	 *
	 *  Hmm, looks like I already put this fix in :-)
	 *
	 */

	d1 = r2 - r1
	if d1 > (r1+max_map_row)-r2 {
		d1 = (r1 + max_map_row) - r2
	}

	d2 = c2 - c1
	if d2 > (c1+max_map_col)-c2 {
		d2 = (c1 + max_map_col) - c2
	}

	//#if 1
	return d1 + d2 /* since there is no diagonal movement */
	//#else
	//    d1 *= d1;
	//    d2 *= d2;
	//
	//    return my_sqrt(d1 + d2);
	//#endif
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

/*
 *  Return the region immediately adjacent to <location>
 *  in direction <dir>
 *
 *  Returns 0 if there is no adjacent location in the given
 *  direction.
 *
 *  Cache version.
 */

func location_direction(where, dir int) int {
	var p *entity_loc

	dir--

	p = rp_loc(where)
	if p == nil || dir >= len(p.prov_dest) {
		return 0
	}

	return p.prov_dest[dir]
}

func exit_distance(loc1, loc2 int) int {
	var dist int
	var w_d int /* where depth */
	var d_d int /* dest depth */

	if subkind(loc1) == sub_hades_pit || subkind(loc2) == sub_hades_pit {
		return 28 // todo: why 28?
	}

	if subkind(loc1) == sub_mine_shaft || subkind(loc1) == sub_mine_shaft_notdone ||
		subkind(loc2) == sub_mine_shaft || subkind(loc2) == sub_mine_shaft_notdone {
		return 0
	}

	if loc_depth(loc1) > loc_depth(loc2) {
		var tmp int

		tmp = loc1
		loc1 = loc2
		loc2 = tmp
	}

	w_d = loc_depth(loc1)
	d_d = loc_depth(loc2)

	if d_d == LOC_build {
		return 0
	}
	if d_d == LOC_subloc {
		return 1
	}

	/*
	 *  water-land links are distance=1
	 */

	if subkind(loc1) == sub_ocean && subkind(loc2) != sub_ocean {
		return 2
	}

	if subkind(loc1) != sub_ocean && subkind(loc2) == sub_ocean {
		return 2
	}

	/*
	 *  Linked sublocs between regions
	 */

	if province(loc1) != province(loc2) {
		loc1 = province(loc1)
		loc2 = province(loc2)
	}

	switch subkind(loc2) {
	case sub_ocean:
		if loc_sea_lane(loc1) != 0 && loc_sea_lane(loc2) != 0 {
			dist = 5
		} else {
			dist = 6
		}
		break

	case sub_mountain:
		dist = 10
		break
	case sub_forest:
		dist = 8
		break
	case sub_swamp:
		dist = 14
		break
	case sub_desert:
		dist = 8
		break
	case sub_plain:
		dist = 7
		break
	case sub_under:
		dist = 7
		break
	case sub_cloud:
		dist = 7
		break

	default:
		panic(fmt.Sprintf("exit_distance: subkind=%s, loc1=%d, loc2=%d, w_d=%d, d_d=%d", subkind_s[subkind(loc2)], loc1, loc2, w_d, d_d))
	}

	return dist
}

func is_port_city(where int) bool {
	var p int
	var n, s, e, w int

	if subkind(where) != sub_city {
		return false
	}

	assert(loc_depth(where) == LOC_subloc)

	p = province(where)

	if subkind(p) == sub_mountain {
		return false
	}

	n = location_direction(p, DIR_N)
	s = location_direction(p, DIR_S)
	e = location_direction(p, DIR_E)
	w = location_direction(p, DIR_W)

	if (n != 0 && subkind(n) == sub_ocean) ||
		(s != 0 && subkind(s) == sub_ocean) ||
		(e != 0 && subkind(e) == sub_ocean) ||
		(w != 0 && subkind(w) == sub_ocean) {
		return true
	}

	return false
}

func province_has_port_city(where int) int {
	var i int
	ret := 0

	assert(loc_depth(where) == LOC_province)

	for _, i = range loop_here(where) {
		if subkind(i) == sub_city && is_port_city(i) {
			ret = i
			break
		}
	}

	return ret
}

//#if 0
//static int
//summer_uldim_open_now()
//{
//    extern int month_done;
//
//    if (oly_month(sysclock) >= 3 && oly_month(sysclock) <= 6)
//        return TRUE;
//
//    if (oly_month(sysclock) == 2 && month_done)
//        return TRUE;
//
//    return FALSE;
//}
//#endif

func add_province_exit(who, where, dest, dir int, l []*exit_view) []*exit_view {
	var n int

	assert(valid_box(dest))

	v := &exit_view{}

	if (is_ship_either(where) && ship_gone(where) != FALSE) || (is_ship_either(dest) && ship_gone(dest) != FALSE) {
		v.in_transit = TRUE
	}

	if is_ship_either(where) && subkind(dest) == sub_ocean {
		v.impassable = TRUE
	}

	if subkind(where) == sub_ocean && !is_ship_either(dest) {
		v.water = TRUE
	}

	if subkind(dest) == sub_ocean {
		v.water = TRUE
	}

	/*
	 *  if land.water && land has a city, then impassable
	 */

	if loc_depth(where) == LOC_province &&
		subkind(dest) == sub_ocean &&
		province_has_port_city(where) != FALSE {
		v.impassable = TRUE
	}

	/*
	 *  Can't go into collapsed mines
	 */

	if subkind(dest) == sub_mine_collapsed {
		v.impassable = TRUE
	}

	/*
	 *  if water.land && land has a city, then dest = the city
	 */

	if subkind(where) == sub_ocean && /* from ocean */
		subkind(dest) != sub_ocean && /* to land */
		subkind(dest) != sub_mountain && /* no mountain ports */
		loc_depth(dest) == LOC_province /* and not islands */ {
		if n = province_has_port_city(dest); n != 0 {
			//#if 0
			//   dest = n
			//#else
			v.impassable = TRUE
			add_province_exit(who, where, n, dir, l)
			//#endif
		}
	}

	/*
	 *  if water-mountain, then impassable
	 */

	if (subkind(where) == sub_mountain && subkind(dest) == sub_ocean) ||
		(subkind(where) == sub_ocean && subkind(dest) == sub_mountain) {
		v.impassable = TRUE
	}

	/*
	 *  if surface-cloud, then impassable (except by FLYing)
	 */

	if (dir == DIR_UP || dir == DIR_DOWN) &&
		(subkind(where) == sub_cloud || subkind(dest) == sub_cloud) {
		v.impassable = TRUE
	}

	//#if 0
	//    /*
	//     *  If Uldim mountains, then impassable
	//     */
	//
	//        if ((dir == DIR_N && uldim(where) == 1) ||
	//            (dir == DIR_S && uldim(where) == 2))
	//        {
	//            v.impassable = TRUE;
	//        }
	//
	//    /*
	//     *  Uldim pass and Summerbridge are passable part of the year
	//     */
	//
	//        if ((dir == DIR_N && (uldim(where) == 4 || summerbridge(where) == 1)) ||
	//            (dir == DIR_S && (uldim(where) == 3 || summerbridge(where) == 2)))
	//        {
	//            if (!summer_uldim_open_now())
	//                v.impassable = TRUE;
	//        }
	//#endif

	v.orig = where
	v.destination = dest
	v.direction = dir
	v.distance = exit_distance(where, dest)

	if loc_hidden(where) {
		v.orig_hidden = TRUE
	}

	if loc_hidden(dest) {
		v.dest_hidden = TRUE
	}

	/*
	 *  Don't make Out routes be hidden.  The character may have poofed
	 *  into a building, and it's unreasonable not to know how to leave.
	 */

	if loc_hidden(dest) && !test_known(who, dest) && dir != DIR_OUT {
		v.hidden = TRUE
	}

	if region(where) != region(dest) {
		v.inside = region(dest)

		if !in_hades(where) && in_hades(dest) {
			v.hades_cost = 100
		}
	}

	/*
	 *  If the destination location is protected by a magical barrier,
	 *  then don't allow travel.
	 */

	if loc_barrier(dest) != FALSE && dir != DIR_OUT {
		v.impassable = TRUE
		v.magic_barrier = TRUE
	}

	l = append(l, v)

	return l
}

func extra_routes(who int, where int, l []*exit_view) []*exit_view {
	var i int
	var dest int
	var v *exit_view

	for _, i = range loop_here(where) {
		if kind(i) == T_road {
			dest = road_dest(i)
			assert(valid_box(dest))

			v = &exit_view{}
			v.orig = where
			v.destination = dest
			v.distance = exit_distance(where, dest)
			v.road = i

			/*
			 *  Surface-cloud links are impassable (except by flying)
			 */

			if (subkind(where) == sub_mountain && subkind(dest) == sub_cloud) ||
				(subkind(where) == sub_cloud && subkind(dest) == sub_mountain) {
				v.impassable = TRUE
			}

			if road_hidden(i) != FALSE {
				v.orig_hidden = TRUE
				v.dest_hidden = TRUE
			}

			if road_hidden(i) != FALSE && !test_known(who, i) {
				v.hidden = TRUE
			}

			if region(where) != region(dest) {
				v.inside = region(dest)
			}

			if subkind(where) == sub_ocean ||
				subkind(dest) == sub_ocean {
				v.water = TRUE
			}

			l = append(l, v)
		}
	}

	return l
}

/*
 *  Exits from a province
 *
 *	cycle through the possible directions, checking each
 *	add all inner locations
 */

func province_exits(who int, where int, l []*exit_view) []*exit_view {
	var dir int
	var n int

	for dir = 1; dir <= DIR_DOWN; dir++ {
		n = location_direction(where, dir)
		if n != 0 {
			l = add_province_exit(who, where, n, dir, l)
		}
	}
	return l
}

func province_sub_exits(who int, where int, l []*exit_view) []*exit_view {
	var i int
	var p *entity_subloc

	for _, i = range loop_here(where) {
		if is_loc_or_ship(i) {
			l = add_province_exit(who, where, i, DIR_IN, l)
		}
	}

	/*
	 *  Mon Dec  9 15:22:12 1996 -- Scott Turner
	 *
	 *  All links are "open" now.
	 *
	 */
	p = rp_subloc(where)

	if p != nil {
		for i = 0; i < len(p.link_from); i++ {
			if p.link_from[i] != 0 {
				add_province_exit(who, where, p.link_from[i], DIR_IN, l)
			}
		}
	}

	return l
}

func subloc_exits(who int, where int, l []*exit_view) []*exit_view {
	var dir int
	var n int
	var i int
	var p *entity_subloc

	if is_port_city(where) {
		p := province(where)

		for dir = 1; dir <= 4; dir++ {
			n = location_direction(p, dir)
			if n != 0 && subkind(n) == sub_ocean {
				l = add_province_exit(who, where, n, dir, l)
			}
		}
	}

	for _, i = range loop_here(where) {
		if is_loc_or_ship(i) {
			add_province_exit(who, where, i, DIR_IN, l)
		}
	}

	/*
	 *  Mon Dec  9 14:42:04 1996 -- Scott Turner
	 *
	 *  Faery hills may be "floating" with no outer location.
	 *
	 *  Wed Jun 24 21:06:47 1998 -- Scott Turner
	 *
	 *  Huh? So you can't exit out into nowhere...
	 *
	 */
	if loc(where) != 0 && subkind(loc(where)) != sub_region {
		add_province_exit(who, where, loc(where), DIR_OUT, l)
	}

	/*
	 *  Mon Dec  9 15:22:12 1996 -- Scott Turner
	 *
	 *  All links are "open" now.
	 *
	 */
	p = rp_subloc(where)

	if p != nil { /*  && loc_link_open(where)) */
		for i = 0; i < len(p.link_to); i++ {
			add_province_exit(who, where, p.link_to[i], 0, l)
		}
	}

	return l
}

/*
 *  If we're on a ship at sea and there are other ships in
 *  the same location, add boarding routes between the ships
 */

func ship_exits(who int, ship int, l []*exit_view) []*exit_view {
	var i int
	var outer_loc int

	assert(is_ship_either(ship))
	outer_loc = loc(ship)

	/*
	 *  We want to maintain the link to the other ships, even after
	 *  they have left.  Attacking another ship requires the link to
	 *  be in place, and we want to allow combat against a ship that
	 *  may have left.  Therefore, we add the ship link, even if
	 *  ship_gone(ship) is true.
	 */

	/*
	 *  Exit from ship to location it is in
	 */

	l = add_province_exit(who, ship, outer_loc, DIR_OUT, l)

	//#if 0
	//    if (subkind(outer_loc) != sub_ocean)
	//        return;
	//#endif

	for _, i = range loop_here(outer_loc) {
		if i != ship && is_ship_either(i) {
			add_province_exit(who, ship, i, 0, l)
		}
	}

	return l
}

func building_exits(who int, where int, l []*exit_view) []*exit_view {
	var i int

	if is_ship_either(where) {
		l = ship_exits(who, where, l)
	} else {
		l = add_province_exit(who, where, loc(where), DIR_OUT, l)
	}

	for _, i = range loop_here(where) {
		if is_loc_or_ship(i) {
			add_province_exit(who, where, i, DIR_IN, l)
		}
	}

	return l
}

func exits_from_loc(who, where int) []*exit_view {
	var l []*exit_view

	switch loc_depth(where) {
	case LOC_province:
		l = province_exits(who, where, l)
		l = province_sub_exits(who, where, l)
	case LOC_subloc:
		l = subloc_exits(who, where, l)
	case LOC_build:
		l = building_exits(who, where, l)
	default:
		panic(fmt.Sprintf("where=%d, depth=%d", where, loc_depth(where)))
	}

	l = extra_routes(who, where, l) /* add secret hidden roads */

	return l
}

func exits_from_loc_nsew(who, where int) []*exit_view {
	if loc_depth(where) != LOC_province {
		return nil
	}
	return province_exits(who, where, nil)
}

func exits_from_loc_nsew_select(who int, where int, land int, do_scramble bool) []*exit_view {
	var l []*exit_view
	var ret []*exit_view
	var i int

	if loc_depth(where) != LOC_province {
		return nil
	}

	ret = nil
	l = exits_from_loc_nsew(who, where)

	for i = 0; i < len(l); i++ {
		if ((land&LAND) != 0 && l[i].water == 0) || ((land&WATER) != 0 && l[i].water != 0) {
			ret = append(ret, l[i])
		}
	}

	if do_scramble {
		ret = shuffle_exits(ret)
	}

	return ret
}

/*
 *  Returns:
 *
 *	0	no ocean access
 *	1	ocean access, but impassable
 *	2	passable ocean access (but may be hidden)
 */

func has_ocean_access(where int) int {
	var l []*exit_view
	var i int
	ret := 0

	l = exits_from_loc(0, where)

	for i = 0; i < len(l); i++ {
		if l[i].water != FALSE {
			if l[i].impassable != FALSE {
				if ret == 0 {
					ret = 1
				}
			} else {
				ret = 2
			}
		}
	}

	return ret
}

func list_exit_extras(who int, v *exit_view) {

	if v.magic_barrier != FALSE {
		indent += 3

		//#if 0
		//        wout(who, "A magical barrier surrounds %s.", box_name(v.destination));
		//        wout(who, "Entry to %s is prevented by a magical barrier.", box_name(v.destination));
		//#endif

		tagout(who, "<tag type=loc_barrier id=%d>", v.destination)
		wout(who, "A magical barrier prevents entry.")
		tagout(who, "</tag type=loc_barrier>")
		indent -= 3
	}

	if v.hades_cost != FALSE {
		indent += 3
		wiout(who, 1, "\"Notice to mortals, from the Gatekeeper Spirit of Hades: 100 gold/head is removed from any stack taking this road.\"")
		indent -= 3
	}

	if rp_loc(v.destination) != nil && garrison_here(v.destination) != FALSE {
		indent += 3

		tagout(who, "<tag type=loc_garr id=%d garr=%d>",
			v.destination, garrison_here(v.destination))
		wout(who, "%s%s",
			liner_desc(garrison_here(v.destination)),
			display_with(garrison_here(v.destination)))
		tagout(who, "</tag type=loc_garr>")
		indent -= 3
	}

	if rp_loc(v.destination) != nil &&
		rp_loc(v.destination).control.closed &&
		province_admin(v.destination) != FALSE {
		indent += 6
		tagout(who, "<tag type=border_closed id=%d>", v.destination)
		wout(who, "Border closed.")
		tagout(who, "</tag type=border_closed>")
		indent -= 6
	}

	if rp_loc(v.destination) != nil {
		tagout(who, "<tag type=fees id=%d nobles=%d weight=%d men=%d>",
			v.destination, rp_loc(v.destination).control.nobles,
			rp_loc(v.destination).control.weight,
			rp_loc(v.destination).control.men)
	}

	if rp_loc(v.destination) != nil &&
		rp_loc(v.destination).control.nobles != FALSE &&
		province_admin(v.destination) != FALSE {
		indent += 6
		wout(who, "Fee of %s per noble to enter.",
			gold_s(rp_loc(v.destination).control.nobles))
		indent -= 6
	}

	if rp_loc(v.destination) != nil &&
		rp_loc(v.destination).control.men != FALSE &&
		province_admin(v.destination) != FALSE {
		indent += 6
		wout(who, "Fee of %s per 100 men to enter.",
			gold_s(rp_loc(v.destination).control.men))
		indent -= 6
	}

	if rp_loc(v.destination) != nil &&
		rp_loc(v.destination).control.weight != FALSE &&
		province_admin(v.destination) != FALSE {
		indent += 6
		wout(who, "Fee of %s per 1000 weight to enter.",
			gold_s(rp_loc(v.destination).control.weight))
		indent -= 6
	}

	if rp_loc(v.destination) != nil {
		tagout(who, "</tag type=fees>")
	}
}

/*
 *	East, swamp, to Athens [aa59], 15 days
 */

/*
 *  We don't list routes to ships that are in_transit
 *
 *  We create them anyway so the player can get a useful error
 *  message if he attempts to traverse them.  Something like "That
 *  ship has left."
 */

func list_exits_sup(who int, where int, v *exit_view, first []byte) {
	ret := ""
	var s string

	//#if 0
	//    if (v.in_transit)
	//        return;
	//#endif

	if v.hidden != FALSE && see_all(who) == FALSE {
		return
	}

	if len(first) != 0 {
		out(who, "%s", first)
		indent += 3
		first = nil
	}

	if v.direction > 0 {
		ret = comma_append(ret, full_dir_s[v.direction])
	}

	s = name(v.destination)

	if len(s) != 0 && !is_ship_either(v.destination) {
		ret = comma_append(ret, subkind_s[subkind(v.destination)])
	}

	ret = comma_append(ret, sout("to %s", box_name(v.destination)))

	if v.inside != FALSE {
		s = name(v.inside)
		if len(s) != 0 {
			ret = comma_append(ret, s)
		}
	}

	if v.dest_hidden != FALSE {
		ret = comma_append(ret, "hidden")
	}

	if v.impassable != FALSE {
		ret = comma_append(ret, "impassable")
	} else {
		ret = comma_append(ret,
			sout("%d~day%s", v.distance, add_s(v.distance)))
	}

	tagout(who, "<tag type=exit id=%d dir=%d dest=%d inside=%d hidden=%d impassable=%d distance=%d>",
		where,
		v.direction,
		v.destination,
		v.inside,
		v.dest_hidden,
		v.impassable,
		v.distance)

	wout(who, "%s", cap_(ret))

	list_exit_extras(who, v)

	tagout(who, "</tag type=exit>")
}

func list_road_sup(who int, where int, v *exit_view, first []byte) {
	var dist string
	hid := ""

	//#if 0
	//    if (v.in_transit)
	//        return;
	//#endif

	if v.hidden != FALSE && see_all(who) == FALSE {
		return
	}

	if v.dest_hidden != FALSE {
		hid = "hidden, "
	}

	if len(first) != 0 {
		out(who, "")
		out(who, "%s", first)
		first = nil
		indent += 3
	}

	if v.impassable != FALSE || v.in_transit != FALSE {
		dist = "impassable"
	} else {
		dist = sout("%d~day%s", add_ds(v.distance))
	}

	tagout(who, "<tag type=exit id=%d dir=%d dest=%d inside=%d hidden=%d impassable=%d distance=%d>",
		where,
		v.direction,
		v.destination,
		v.inside,
		v.dest_hidden,
		v.impassable,
		v.distance)

	out(who, "%s, to %s, %s%s",
		just_name(v.road),
		box_name(v.destination),
		hid,
		dist)

	list_exit_extras(who, v)
	tagout(who, "</tag type=exit>")
}

func list_exits(who, where int) {
	var l []*exit_view
	var i int

	tagout(who, "<tag type=exits loc=%d>", where)

	l = exits_from_loc(who, where)

	/*
	 *  direction may be zero for roads and secret passages
	 */

	first := fmt.Sprintf("Routes leaving %s: ", just_name(where))

	for i = 0; i < len(l); i++ {
		if l[i].road == 0 && (l[i].direction != DIR_IN || see_all(who) == 2) {
			list_exits_sup(who, where, l[i], []byte(first))
		}
	}

	for i = 0; i < len(l); i++ {
		if l[i].road != FALSE {
			list_road_sup(who, where, l[i], []byte(first))
		}
	}

	if len(first) != 0 {
		//#if 0
		//        wout(who, "No known routes leaving %s", box_name(where));
		//#else
		if is_ship_either(where) {
			wout(who, "No current exits from %s", box_name(where))
		} else {
			wout(who, "No known routes leaving %s", box_name(where))
		}
		//#endif
	} else {
		indent -= 3
	}

	tagout(who, "</tag type=exits>")

}

func list_sailable_routes(who, ship int) {
	var outer_loc int
	var l []*exit_view
	var i int

	if !is_ship_either(ship) {
		return
	}

	outer_loc = loc(ship)
	l = exits_from_loc(who, outer_loc)

	first := fmt.Sprintf("Ocean routes:")

	tagout(who, "<tag type=sail_routes id=%d>", outer_loc)

	for i = 0; i < len(l); i++ {
		if l[i].direction > 0 &&
			(l[i].direction != DIR_IN || see_all(who) == 2) &&
			l[i].water != FALSE {
			list_exits_sup(who, outer_loc, l[i], []byte(first))
		}
	}

	for i = 0; i < len(l); i++ {
		if l[i].road != FALSE && l[i].water != FALSE {
			list_road_sup(who, outer_loc, l[i], []byte(first))
		}
	}

	if len(first) != 0 {
		out(who, "No visible sailable routes")
	} else {
		indent -= 3
	}

	out(who, "")
	tagout(who, "</tag type=sail_routes>")
}

func count_hidden_exits(l []*exit_view) int {
	sum := 0
	var i int

	for i = 0; i < len(l); i++ {
		if l[i].hidden != FALSE {
			sum++
		}
	}

	return sum
}

func hidden_count_to_index(which int, l []*exit_view) int {
	var i int

	for i = 0; i < len(l); i++ {
		if l[i].hidden != FALSE {
			which--
		}

		if which <= 0 {
			assert(l[i].hidden != FALSE)
			return i
		}
	}

	panic("!reached")
}

func find_hidden_exit(who int, l []*exit_view, which int) {
	where := subloc(who)

	if is_ship(where) {
		where = subloc(where)
	}

	assert(valid_box(who))
	assert(which < len(l))
	assert(l[which].hidden != FALSE)

	if l[which].road != FALSE {
		wout(who, "A hidden route has been found in %s!", box_name(where))
		out(who, "")

		set_known(who, l[which].road)
		l[which].hidden = FALSE

		indent += 3
		list_road_sup(who, subloc(who), l[which], nil)
		indent -= 3
	} else if l[which].direction == DIR_IN {
		wout(who, "A hidden inner location has been found in %s!", box_name(where))
		out(who, "")

		set_known(who, l[which].destination)
		l[which].hidden = FALSE

		indent += 3
		wout(who, "%s", liner_desc(l[which].destination))
		indent -= 3
	} else {
		wout(who, "A hidden route has been found in %s!", box_name(where))
		out(who, "")

		set_known(who, l[which].destination)
		l[which].hidden = FALSE

		indent += 3
		list_exits_sup(who, subloc(who), l[which], nil)
		indent -= 3
	}

	out(who, "")
}
