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

func add_char_here(who int, l []int) []int {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}

	l = append(l, who)

	p := rp_loc_info(who)
	if p == nil {
		panic("assert(p != nil)")
	}

	for i := 0; i < len(p.here_list); i++ {
		if kind(p.here_list[i]) == T_char {
			l = add_char_here(p.here_list[i], l)
		}
	}

	return l
}

func add_here(who int, l []int) []int {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	l = append(l, who)
	p := rp_loc_info(who)
	if p == nil {
		panic("assert(p != nil)")
	}

	for i := 0; i < len(p.here_list); i++ {
		l = add_here(p.here_list[i], l)
	}

	return l
}

func add_to_here_list(loc, who int) {
	if in_here_list(loc, who) {
		panic("assert(!in_here_list(loc, who))")
	}
	if p_loc_info(loc) == nil {
		panic(fmt.Sprintf("assert(p_loc_info(%d) != nil)", loc))
	}
	p_loc_info(loc).here_list = append(p_loc_info(loc).here_list, who)
	if !in_here_list(loc, who) {
		panic(fmt.Sprintf("assert(in_here_list(%d, %d))", loc, who))
	}
}

func all_char_here(who int, l []int) []int {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	assert(valid_box(who))

	l = nil

	p := rp_loc_info(who)
	if p == nil {
		return l
	}

	for i := 0; i < len(p.here_list); i++ {
		if kind(p.here_list[i]) == T_char {
			l = add_char_here(p.here_list[i], l)
		}
	}

	return l
}

func all_here(who int, l []int) []int {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	l = nil

	p := rp_loc_info(who)
	if p == nil {
		return l
	}

	for i := 0; i < len(p.here_list); i++ {
		l = add_here(p.here_list[i], l)
	}

	return l
}

func all_stack(who int, l []int) []int {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	assert(valid_box(who))

	l = nil
	l = append(l, who)

	p := rp_loc_info(who)

	if p == nil {
		return l
	}

	for i := 0; i < len(p.here_list); i++ {
		if kind(p.here_list[i]) == T_char {
			l = add_char_here(p.here_list[i], l)
		}
	}

	return l
}

func building_owner(where int) int {
	if loc_depth(where) != LOC_build {
		panic("assert(loc_depth(where) == LOC_build)")
	}
	return first_character(where)
}

func city_here(a int) int { return subloc_here((a), sub_city) }

func count_loc_structures(where, a, b int) int {
	sum := 0
	for _, i := range loop_here(where) {
		if kind(i) == T_loc && (int(subkind(i)) == a || int(subkind(i)) == b) {
			sum++
		}
	}

	return sum
}

func in_here_list(loc, who int) bool {
	p := rp_loc_info(loc)
	if p == nil {
		return false
	}
	return ilist_lookup(p.here_list, who) != -1
}

func in_safe_now(who int) int {
	for {
		if safe_haven(who) {
			return TRUE
		}
		who = loc(who)
		if who > 0 {
			continue
		}
		break
	}
	return FALSE
}

func loc_owner(where int) int {
	panic("!implemented")
}

/*
 *  Mark that each member of a stack (or a ship) has visited a location
 */
func mark_loc_stack_known(stack, where int) {
	if kind(stack) == T_char {
		set_known(stack, where)
	}

	for _, i := range loop_char_here(stack) {
		if !is_prisoner(i) {
			set_known(i, where)
		}
	}
}

/*
 *  Return the ultimate province a character is in
 */
func province(who int) int {
	if item_unique(who) != FALSE {
		who = item_unique(who)
	}
	for who > 0 && (kind(who) != T_loc || loc_depth(who) != LOC_province) {
		who = loc(who)
	}
	return who
}

// todo: this was commented out?
func province_owner(where int) int {
	prov := province(where)
	city := city_here(prov)
	castle := subloc_here(prov, sub_castle)
	if castle == 0 && city != 0 {
		castle = subloc_here(city, sub_castle)
	}
	if castle != 0 {
		return first_character(castle)
	}
	return 0
}

/*
 *  Return the ultimate region a character is in
 */
func region(who int) int {
	for who > 0 && (kind(who) != T_loc || loc_depth(who) != LOC_region) {
		who = loc(who)
	}
	return who
}

func remove_from_here_list(loc, who int) {
	if !in_here_list(loc, who) {
		panic("assert(in_here_list(loc, who))")
	}
	var l []int
	for _, el := range rp_loc_info(loc).here_list {
		if el == who {
			continue
		}
		l = append(l, el)
	}
	rp_loc_info(loc).here_list = l

	/*
	 *  Mon Apr 20 18:08:44 1998 -- Scott Turner
	 *
	 *  Thanks to Rich's nice encapsulation, this is the only
	 *  place we have to worry about someone being the last person
	 *  out of a subloc...so here's where we reset the fees.
	 *
	 *  Anyway, that's plan :-).
	 *
	 */
	if len(rp_loc_info(loc).here_list) == 0 && rp_subloc(loc) != nil {
		rp_subloc(loc).control.weight = 0
		rp_subloc(loc).control.nobles = 0
		rp_subloc(loc).control.men = 0

		rp_subloc(loc).control2.weight = 0
		rp_subloc(loc).control2.nobles = 0
		rp_subloc(loc).control2.men = 0
	}

	if in_here_list(loc, who) {
		panic("assert(!in_here_list(loc, who))")
	}
}

/*
 *  This check could be expanded to make sure that new_loc is
 *  not anywhere inside of who, by walking up new_loc to the top
 *  and making sure we don't go through who
 */
func set_where(who, new_loc int) {
	if who == new_loc {
		panic("assert(who != new_loc)")
	}
	old_loc := loc(who)
	if old_loc > 0 {
		remove_from_here_list(old_loc, who)
	}
	if new_loc > 0 {
		add_to_here_list(new_loc, who)
	}
	p_loc_info(who).where = new_loc
	//if is_loc_or_ship(loc(who)) {
	//	if is_prisoner(who) {
	//		panic("assert(!is_prisoner(who))")
	//	}
	//}
}

/*
 *  is b somewhere inside of a?
 */
func somewhere_inside(a, b int) int {
	if a == b {
		return FALSE
	}
	for b > 0 {
		b = loc(b)
		if a == b {
			return TRUE
		}
	}
	return FALSE
}

/*
 *  Return the immediate location a character is in, irrespective
 *  of who or what we are stacked with.  This may be a sublocation
 *  such as a city, or a province.
 */
func subloc(who int) int {
	for {
		who = loc(who)
		if who > 0 && kind(who) != T_loc && kind(who) != T_ship {
			continue
		}
		break
	}
	return who
}

func subloc_here(where, sub_kind int) int {
	for _, i := range loop_here(where) {
		if kind(i) == T_loc && int(subkind(i)) == sub_kind {
			return i
		}
	}
	return 0
}

/*
 *  Step out from a location until we get to the appropriate viewing
 *  level.  Currently, provinces see into everything except cities.
 *
 *  I have added graveyards and faery hills to the viewloc hierarchy
 *  so that players will be able to see the hidden links into Faery
 *  and Hades when they are inside these locations.  Otherwise, they
 *  must carefully note the message given when they find the hidden
 *  route.
 */
func viewloc(who int) int {
	for who > 0 && loc_depth(who) != LOC_province && subkind(who) != sub_guild && subkind(who) != sub_city && subkind(who) != sub_graveyard && subkind(who) != sub_faery_hill {
		who = loc(who)
	}
	return who
}
