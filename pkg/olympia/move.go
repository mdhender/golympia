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
	"strings"
)

const (
	FORCED_RIDE  = 2
	FORCED_MARCH = 1
)

var (
	ocean_chars []int
)

func ifndef(t bool, fn func()) {
	if !t {
		fn()
	}
}

func departure_message(who int, v *exit_view) {
	var to string
	var with string
	var desc string
	var comma string

	assert(valid_box(who))

	if !is_loc_or_ship(v.destination) ||
		!is_loc_or_ship(v.orig) {
		return
	}

	if char_really_hidden(who) {
		return
	}

	if loc_depth(v.orig) == LOC_province &&
		(weather_here(v.orig, sub_fog) != FALSE || weather_here(v.orig, sub_fog) != FALSE) &&
		is_priest(who) != sk_domingo {
		return
	}

	if v.dest_hidden != FALSE {
		return
	}

	desc = liner_desc(who) /* consumes mucho souts */

	if subloc(v.destination) == v.orig {
		to = sout(" entered %s", box_name(v.destination))
	} else if subloc(v.orig) == v.destination {
		to = sout(" exited %s", box_name(v.orig))
	} else if viewloc(v.orig) != viewloc(v.destination) {
		if v.direction >= 1 && v.direction <= 4 {
			to = sout(" went %s", full_dir_s[v.direction])
		} else {
			to = sout(" left for %s", box_name(v.destination))
		}
	} else {
		return
	}

	with = display_with(who)

	if strings.IndexByte(desc, ',') != -1 {
		comma = ","
	}

	if with != "" {
		with = "."
	}

	if viewloc(v.orig) != viewloc(v.destination) {
		garr := garrison_here(v.orig)

		if garr != 0 && garrison_notices(garr, who) != FALSE {
			show_to_garrison = true
		}
	}

	wout(v.orig, "%s%s%s%s", desc, comma, to, with)
	show_chars_below(v.orig, who)

	show_to_garrison = false
}

func arrival_message(who int, v *exit_view) {
	var from string
	var with string
	var desc string
	var comma string

	if char_really_hidden(who) {
		return
	}
	/*
	 *  Show nothing if foggy or misty here.
	 *
	 */
	if loc_depth(v.destination) == LOC_province &&
		(weather_here(v.destination, sub_fog) != FALSE ||
			weather_here(v.destination, sub_mist) != FALSE) {
		return
	}

	desc = liner_desc(who) /* consumes mucho souts */

	if v.orig_hidden == FALSE {
		if v.direction >= 1 && v.direction <= 4 {
			from = sout(" from the %s",
				full_dir_s[exit_opposite[v.direction]])
		} else {
			from = sout(" from %s", box_name(v.orig))
		}
	}

	with = display_with(who)

	if strings.IndexByte(desc, ',') != -1 {
		comma = ","
	}

	if with != "" {
		with = "."
	}

	if viewloc(v.orig) != viewloc(v.destination) {
		garr := garrison_here(v.destination)

		if garr != 0 {
			if garrison_notices(garr, who) != FALSE {
				show_to_garrison = true
			}

			if garrison_spot_check(garr, who) != FALSE {
				indent += 3
				wout(garr, "%s%s", desc, with)
				show_chars_below(garr, who)
				indent -= 3
			}
		}
	}

	wout(v.destination, "%s%s arrived%s%s", desc, comma, from, with)
	show_chars_below(v.destination, who)

	show_to_garrison = false
}

/*
 *  Mark that we know both ends of a hidden road we're about to go through
 */

func discover_road(who, where int, v *exit_view) {
	l := exits_from_loc(who, v.destination)
	for i := 0; i < len(l); i++ {
		if l[i].road != FALSE && l[i].destination == where {
			set_known(who, l[i].road)
			set_known(who, v.road)

			for _, j := range loop_char_here(who) {
				set_known(j, l[i].road)
				set_known(j, v.road)
			}

		}
	}
}

func parse_exit_dir(c *command, where int, zero_arg string) *exit_view {
	l := exits_from_loc(c.who, where)

	if valid_box(c.a) && subkind(c.a) != sub_region {
		if where == c.a {
			if zero_arg != "" {
				wout(c.who, "Already in %s.", box_name(where))
			}
			return nil
		}
		/*
		 *  Give priority to passable routes.  A secret passable route may
		 *  parallel a visible impassable route.
		 */
		{
			var ret, impass_ret *exit_view
			for i := 0; i < len(l); i++ {
				if l[i].destination == c.a && (l[i].hidden == FALSE || see_all(c.who) != FALSE) {
					if l[i].impassable != FALSE {
						impass_ret = l[i]
					} else {
						ret = l[i]
					}
				}
			}
			if ret != nil {
				return ret
			} else if impass_ret != nil {
				return impass_ret
			}
		}

		//#if 0
		//        if (zero_arg)
		//          wout(c.who, "No visible route from %s to %s.",
		//           box_name(where),
		//           c.parse[1]);
		//#endif

		return nil
	}

	dir := lookup_sb(full_dir_s, c.parse[1])
	if dir < 0 {
		dir = lookup_sb(short_dir_s, c.parse[1])
	}

	if dir < 0 {
		//#if 0
		//        if (zero_arg)
		//          wout(c.who, "Unknown direction or destination '%s'.",
		//           c.parse[1]);
		//#endif
		return nil
	}

	for i := 0; i < len(l); i++ {
		if l[i].direction == dir && (l[i].hidden == FALSE || see_all(c.who) != FALSE) {
			if dir == DIR_IN && zero_arg != "" {
				wout(c.who, "(assuming '%s %s')", zero_arg, box_code_less(l[i].destination))
			}
			return l[i]
		}
	}

	/*
	 *  Wed Nov 18 18:28:23 1998 -- Scott Turner
	 *
	 *  Convert "move out" to "move up" if out isn't available and
	 *  up is.
	 *
	 */
	if dir == DIR_OUT {
		for i := 0; i < len(l); i++ {
			if l[i].direction == DIR_UP && (l[i].hidden == FALSE || see_all(c.who) != FALSE) {
				wout(c.who, "(assuming 'move up')")
				return l[i]
			}
		}
	}

	if dir == DIR_IN {
		for i := 0; i < len(l); i++ {
			if l[i].direction == DIR_DOWN && (l[i].hidden == FALSE || see_all(c.who) != FALSE) {
				wout(c.who, "(assuming 'move down')")
				return l[i]
			}
		}
	}

	//#if 0
	//    if (zero_arg)
	//      wout(c.who, "No visible %s route from %s.",
	//       full_dir_s[dir], box_name(where));
	//#endif

	return nil
}

func kill_random_mount(who int) {
	for _, e := range loop_inventory(who) {
		if item_ride_cap(e.item) >= 100 {
			wout(who, "The forced ride kills one %s.", just_name(e.item))
			sub_item(who, e.item, 1)
			return
		}
	}
	panic("!reached")
}

// move_exit_land returns the number of days of travel needed?
func move_exit_land(c *command, v *exit_view, show bool) int {
	delay := v.distance
	if delay == 0 {
		return 0
	}

	/* destination terrain */
	terr := subkind(province(v.destination))
	/* traveling in a swamp? */
	swamp := terr == sub_swamp || subkind(v.destination) == sub_bog || subkind(v.destination) == sub_pits
	/* traveling in the mountains? */
	mountains := terr == sub_mountain

	var w weights
	determine_stack_weights(c.who, &w, mountains)

	/*
	 *  Are you riding?
	 *
	 */
	if w.ride_cap >= w.ride_weight && !swamp {
		/*
		 *  You're riding.
		 *
		 */
		delay -= delay / 2
		/*
		 *  Are you also forcing the pace?
		 *
		 *  Tue May 25 06:45:58 1999 -- Scott Turner
		 *
		 *  Do not ignore ninjas, etc!
		 *
		 */
		if get_effect(c.who, ef_forced_march, 0, 0) != FALSE && char_really_alone(c.who) {
			delay -= delay / 2
			v.forced_march = FORCED_RIDE
		}
		/*
		 *  Or perhaps you're doing a forced march?
		 *
		 */
	} else if get_effect(c.who, ef_forced_march, 0, 0) != FALSE && char_really_alone(c.who) {
		delay -= delay / 2
		v.forced_march = FORCED_MARCH
	} else {
		/*
		 *  You're walking.
		 *
		 */
		if w.land_weight > w.land_cap*2 {
			wout(c.who, "%s is too overloaded to travel.", box_name(c.who))
			wout(c.who, "You have a total of %d weight, and your maximum capacity is %d.", w.land_weight, w.land_cap*2)
			return -1
		}
		/*
		 *  In a swamp with animals?
		 *
		 */
		if swamp && w.animals != FALSE {
			if show {
				wout(c.who, "Difficult terrain slows the animals.  Travel will take an extra day.")
			}
			delay += 1
		}

		/*
		 *  Overloaded?
		 *
		 */
		if w.land_weight > w.land_cap {
			ratio := (w.land_weight - w.land_cap) * 100 / w.land_cap
			additional := delay * ratio / 100
			if show {
				if additional == 1 {
					wout(c.who, "Excess inventory slows movement.  Travel will take an extra day.")
				} else if additional > 1 {
					wout(c.who, "Excess inventory slows movement.  Travel will take an extra %s days.", nice_num(additional))
				}
			}
			delay += additional
		}
	}

	/*
	 *  The fast move effect speeds you by 2x if you're moving into your holy terrain.
	 *
	 */
	val := get_all_effects(c.who, ef_fast_move, 0, 0)
	if val != FALSE && is_holy_terrain(c.who, v.destination) {
		additional := delay / 2
		if show {
			if additional == 1 {
				wout(c.who, "Holy guidance speeds your travel by one day.")
			} else {
				wout(c.who, "Holy guidance speeds your travel by %d days.", additional)
			}
		}
		delay -= additional
		for i := 0; i < val; i++ {
			delete_effect(c.who, ef_fast_move, 0)
		}
	}

	/*
	 *  The slow move effect slows you by 2x.
	 *
	 *  Mon Oct  7 12:19:17 1996 -- Scott Turner
	 *
	 *  Shouldn't effect other priests...
	 *
	 */
	if get_effect(v.destination, ef_slow_move, 0, 0) != FALSE {
		/*
		 *  A priest?
		 *
		 */
		if is_holy_terrain(c.who, v.destination) {
			if show {
				wout(c.who, "You sense the power of %s obscuring the entrances to %s.", god_name(is_priest(c.who)), box_name(v.destination))
			}
		} else {
			if show {
				if delay == 1 {
					wout(c.who, "A magical effect slows your movement by one day.")
				} else {
					wout(c.who, "A magical effect slows your movement by %d days.", delay)
				}
			}
			delay += delay
		}
	}

	/*
	 *  Tue Oct  6 18:10:55 1998 -- Scott Turner
	 *
	 *  Perchance you have an artifact.
	 *
	 */
	if art := best_artifact(c.who, ART_FAST_TERR, v.destination, 0); art != FALSE {
		delay -= rp_item_artifact(art).param2
		wout(c.who, "Your passage is magically sped by %s day%s.", nice_num(rp_item_artifact(art).param2), or_string(rp_item_artifact(art).param2 > 1, "s", ""))
	}

	if delay < 0 {
		delay = 0
	}

	return delay
}

func move_exit_fly(c *command, v *exit_view, show bool) int {
	if subkind(v.destination) == sub_under {
		wout(c.who, "Cannot fly underground.")
		return -1
	}

	delay := v.distance
	if delay > 3 {
		delay = 3
	}

	var w weights
	determine_stack_weights(c.who, &w, false)

	if w.fly_cap < w.fly_weight {
		wout(c.who, "%s is too overloaded to fly.", box_name(c.who))
		wout(c.who, "You have a total of %d weight, and your maximum flying capacity is %d.", w.fly_weight, w.fly_cap)
		return -1
	}

	return delay
}

func save_v_array(c *command, v *exit_view) {

	//#if 0
	//    c.v.direction = v.direction;
	//    c.v.destination = v.destination;
	//    c.v.road = v.road;
	//    c.v.dest_hidden = v.dest_hidden;
	//    c.v.distance = v.distance;
	//    c.v.orig = v.orig;
	//    c.v.orig_hidden = v.orig_hidden;
	//#endif
	c.b = v.direction
	c.c = v.destination
	c.d = v.road
	c.e = v.dest_hidden
	c.f = v.distance
	c.g = v.orig
	c.h = v.orig_hidden
	c.i = v.seize
	c.j = v.forced_march
}

func restore_v_array(c *command, v *exit_view) {
	*v = exit_view{} // bzero(v, sizeof(*v));

	//#if 0
	//    v.direction = c.v.direction;
	//    v.destination = c.v.destination;
	//    v.road = c.v.road;
	//    v.dest_hidden = c.v.dest_hidden;
	//    v.distance = c.v.distance;
	//    v.orig = c.v.orig;
	//    v.orig_hidden = c.v.orig_hidden;
	//#endif

	v.direction = c.b
	v.destination = c.c
	v.road = c.d
	v.dest_hidden = c.e
	v.distance = c.f
	v.orig = c.g
	v.orig_hidden = c.h
	v.seize = c.i
	v.forced_march = c.j
}

func suspend_stack_actions(who int) {
	for _, i := range loop_stack(who) {
		p_char(i).moving = sysclock.days_since_epoch
	}
}

func restore_stack_actions(who int) {
	for _, i := range loop_stack(who) {
		p_char(i).moving = 0
	}
}

func clear_guard_flag(who int) {
	if kind(who) == T_char {
		p_char(who).guard = FALSE
	}
	for _, i := range loop_char_here(who) {
		p_char(i).guard = FALSE
	}
}

func land_check(c *command, v *exit_view, show bool) bool {
	if v.water != FALSE {
		if show {
			wout(c.who, "A sea-worthy ship is required for travel across water.")
		}
		return false
	}

	if v.impassable != FALSE {
		if show {
			wout(c.who, "That route is impassable.")
		}
		return false
	}

	if v.in_transit != FALSE {
		if show {
			wout(c.who, "%s is underway.  Boarding is not possible.", box_name(v.destination))
		}
		return false
	}

	//#if 0
	//var owner int        /* owner of loc we're moving into, if any */
	//    if (loc_depth(v.destination) == LOC_build &&
	//        (owner = building_owner(v.destination)) &&
	//        FALSE == will_admit(owner, c.who, v.destination) &&
	//        v.direction != DIR_OUT) {
	//      if (show) {
	//        wout(c.who, "%s refused to let us enter.",
	//         box_name(owner));
	//        wout(owner, "Refused to let %s enter.",
	//         box_name(c.who));
	//      }
	//      return FALSE;
	//    }
	//#endif

	return true
}

func can_move_here(where int, c *command) bool {
	v := parse_exit_dir(c, where, "")
	if v != nil && v.direction != DIR_IN && land_check(c, v, false) && move_exit_land(c, v, false) >= 0 {
		return true
	}
	return false
}

func can_move_at_outer_level(where int, c *command) bool {
	outer := subloc(where)
	for loc_depth(outer) > LOC_region {
		if can_move_here(outer, c) {
			// todo: kinda guessing at this
			return loc_depth(outer) == loc_depth(where)
		}
		outer = subloc(outer)
	}
	return false
}

func is_smuggling(who, ef int) bool {
	return get_effect(who, ef, 0, 0) != FALSE
}

func check_smuggling(who, ef, sk int) bool {
	// if they don't have the skill, they automatically fail
	if rp_skill_ent(who, sk) == nil {
		return false
	}
	// how experienced are they?
	exp := 50 + skill_exp(who, sk)
	// do they succeed?
	return (rnd(1, 100) < exp)
}

func smuggle_savings(who, cost, sk int) int {
	reduce := (cost * rnd(50, 100)) / 100
	if reduce < 1 {
		reduce = 1
	}
	if reduce > cost {
		reduce = cost
	}
	/* Only increase 1x month? */
	add_skill_experience(who, sk)
	return reduce
}

func smuggle_fine(cost int) int {
	fine := (cost * rnd(5, 13)) / 8
	if fine < 1 {
		fine = 1
	}
	return fine
}

func pay_fine(c *command, fine, ruler int) bool {
	total := stack_has_item(c.who, item_gold)

	if total < fine {
		autocharge(c.who, total)
		gen_item(ruler, item_gold, total)
		wout(VECT, "Since you could not pay the entire fine, you are questioned for 3 days and then released.")
		c.second_wait += 3
		return false
	}
	autocharge(c.who, fine)
	gen_item(ruler, item_gold, fine)
	wout(VECT, "You pay the fine and %s is detained for questioning for 3 days.", box_name(c.who))
	prepend_order(player(c.who), c.who, "wait time 3")
	return true
}

func calc_entrance_fee(control *loc_control_ent, c *command, ruler int) int {
	vector_stack(c.who, true)

	w_cost := 0
	if control.weight != 0 {
		var w weights
		determine_stack_weights(c.who, &w, false)
		w_cost = (control.weight * w.land_weight) / 1000
		if is_smuggling(c.who, ef_smuggle_goods) {
			if check_smuggling(c.who, ef_smuggle_goods, sk_smuggle_goods) {
				// success
				reduce := smuggle_savings(c.who, w_cost, sk_smuggle_goods)
				if reduce != 0 {
					indent += 3
					wout(c.who, "You smuggle %s worth of goods unseen through the gates.", gold_s(reduce))
					indent -= 3
				}
				w_cost -= reduce
			} else if rnd(1, 4) == 1 {
				// caught smuggling
				fine := smuggle_fine(w_cost)
				indent += 3
				wout(VECT, "You are caught attempting to smuggle and fined %s!", gold_s(fine))
				indent -= 3
				if !pay_fine(c, fine, ruler) {
					return -1
				}
			}
		}

		if w_cost != 0 {
			wout(c.who, "Entrance fee is %s on %d weight.", gold_s(w_cost), w.land_weight)
		}
	}

	m_cost := 0
	if control.men != 0 {
		m_cost = (control.men * count_stack_any(c.who)) / 100
		if is_smuggling(c.who, ef_smuggle_men) {
			if check_smuggling(c.who, ef_smuggle_men, sk_smuggle_men) {
				// success
				reduce := smuggle_savings(c.who, m_cost, sk_smuggle_men)
				if reduce != 0 {
					indent += 3
					wout(c.who, "You smuggle %s worth of men unseen through the gates.", gold_s(reduce))
					indent -= 3
				}
				m_cost -= reduce
			} else if rnd(1, 4) == 1 {
				// caught smuggling
				fine := smuggle_fine(m_cost) // mdhender: changed from w_cost to m_cost
				indent += 3
				wout(VECT, "You are caught attempting to smuggle and fined %s!", gold_s(fine))
				indent -= 3
				if !pay_fine(c, fine, ruler) {
					return -1
				}
			}
		}

		if m_cost != 0 {
			wout(c.who, "Entrance fee is %s on %d men/beasts.", gold_s(m_cost), count_stack_any(c.who))
		}
	}

	n_cost := 0
	if control.nobles != 0 {
		n_cost = control.nobles * count_stack_units(c.who)
		if is_smuggling(c.who, ef_smuggle_men) {
			if check_smuggling(c.who, ef_smuggle_men, sk_smuggle_men) {
				// success
				reduce := smuggle_savings(c.who, n_cost, sk_smuggle_men)
				if reduce != 0 {
					indent += 3
					wout(c.who, "You smuggle %s worth of nobles unseen through the gates.", gold_s(reduce))
					indent -= 3
				}
				n_cost -= reduce
			} else if rnd(1, 4) == 1 {
				// caught smuggling
				fine := smuggle_fine(n_cost) // mdhender: changed from w_cost to n_cost
				indent += 3
				wout(VECT, "You are caught attempting to smuggle and fined %s!", gold_s(fine))
				indent -= 3
				if !pay_fine(c, fine, ruler) {
					return -1
				}
			}
		}

		if n_cost != 0 {
			wout(c.who, "Entrance fee is %s on %d nobles.", gold_s(n_cost), count_stack_units(c.who))
		}
	}

	return (w_cost + m_cost + n_cost)
}

func charge_entrance_fees(who, ruler, cost int) bool {
	if !autocharge(who, cost) {
		wout(VECT, "Can't afford %s in fees to enter, so you are turned away.", gold_s(cost))
		return false
	}
	gen_item(ruler, item_gold, cost)
	wout(ruler, "Received %s in entrance fees from %s.", gold_s(cost), box_name(who))
	wout(VECT, "%s took %s in entrance fees from us.", box_name(ruler), gold_s(cost))
	gold_fees += cost
	return true
}

/*
 *  Thu Mar 20 11:44:30 1997 -- Scott Turner
 *
 *  Who is immediately in control of this location?
 *
 */
func controls_loc(where int) int {
	if !valid_box(where) || (kind(where) != T_loc && kind(where) != T_ship) {
		return 0
	} else if kind(where) == T_loc && loc_depth(where) == LOC_province {
		return garrison_here(where)
	} else if subkind(where) == sub_city {
		return garrison_here(province(where))
	}
	return first_character(where)
}

func player_controls_loc(where int) int {
	if !valid_box(where) ||
		(kind(where) != T_loc && kind(where) != T_ship) {
		return 0
	} else if kind(where) == T_loc && loc_depth(where) == LOC_province {
		return player(province_admin(where))
	} else if subkind(where) == sub_city {
		return player(province_admin(province(where)))
	}
	return player(first_character(where))
}

/*
 *  Tue Dec 29 11:30:19 1998 -- Scott Turner
 *
 *  Check to see if <who> can join <g>.
 *
 */
func can_join_guild(who, g int) bool {
	/*
	 *  Can't already belong to a guild.
	 *
	 */
	if guild_member(who) != FALSE {
		wout(who, "%s is already a guild member.", box_name(who))
		return false
	}

	/*
	 *  Priests can only join their strength guild.
	 *
	 */
	if is_priest(who) != FALSE && g != rp_relig_skill(is_priest(who)).strength {
		wout(who, "You may only join the %s Guild.", box_name(rp_relig_skill(is_priest(who)).strength))
		return false
	}
	//#if 0
	//    /*
	//     *  Obvious magicians can't join any guild.
	//     *
	//     */
	//    if (is_magician(who) && !char_hide_mage(who)) {
	//      wout(who,"Magicians are not welcome here.");
	//      return FALSE;
	//    }
	//#endif

	/*
	 *  You need to already know all the learnable and researchable
	 *  skills in this school.
	 *
	 */
	p := rp_skill(g)
	if p == nil {
		return false
	}

	for i := 0; i < len(p.offered); i++ {
		if FALSE == has_skill(who, p.offered[i]) {
			wout(who, "You must master all available skills before entering a guild.")
			return false
		}
	}
	for i := 0; i < len(p.research); i++ {
		if FALSE == has_skill(who, p.research[i]) {
			wout(who, "You must master all available skills before entering a guild.")
			return false
		}
	}

	return true
}

/*
 *  Tue Dec 29 11:35:49 1998 -- Scott Turner
 *
 *  Join a guild.
 *
 */
func join_guild(who, g int) bool {
	if kind(g) != T_skill || skill_school(g) != g {
		return false
	}
	if !can_join_guild(who, g) {
		return false
	}
	wout(who, "%s joins the %s Guild.", just_name(who), just_name(g))
	rp_char(who).guild = g
	return true
}

func check_guilds(c *command, v *exit_view) bool {
	/*
	 *  Temples
	 *
	 *  Only a priest of the appropriate religion can move into
	 *  a temple (but this doesn't prevent you from ATTACKing it or
	 *  RAZEing it).
	 *
	 */
	if is_temple(v.destination) != FALSE && is_temple(v.destination) != is_priest(c.who) {
		wout(c.who, "Only a priest of %s may enter %s.", box_name(is_temple(v.destination)), box_name(v.destination))
		return false
	}

	/*
	 *  Guilds
	 *
	 *  Temples are really just guilds, so we do the same thing here.
	 *
	 */
	if is_guild(v.destination) != FALSE && is_guild(v.destination) != guild_member(c.who) {
		return join_guild(c.who, is_guild(v.destination))
	}

	/*
	 *  Otherwise, you may enter...
	 *
	 */
	return true
}

func do_actual_move(c *command, v *exit_view, delay int) {
	attack := (strcasecmp_bs(c.parse[0], "attack") == 0)

	v.distance = delay
	c.wait = delay

	save_v_array(c, v)
	leave_stack(c.who)

	// todo: did i translate this correctly?
	if delay-or_int(attack, 1, 0) > 0 { // was if ((delay - attack) > 0)
		wout(VECT, "Travel to %s will take %s day%s.", box_name(v.destination), nice_num(delay), or_string(delay == 1, "", "s"))
	}

	suspend_stack_actions(c.who)
	clear_guard_flag(c.who)

	// todo: did i translate this correctly?
	if delay-or_int(attack, 1, 0) > 1 { // was if ((delay - attack) > 1)
		prisoner_movement_escape_check(c.who)
	}

	departure_message(c.who, v)

}

/*
 *  Sun Mar 30 18:40:12 1997 -- Scott Turner
 *
 *  Do an attack after failing a move.
 *
 */
func do_move_attack(c *command, v *exit_view) {
	prepend_order(player(c.who), c.who, fmt.Sprintf("attack %d", v.destination))
}

/*
 *  Sun Mar 30 20:21:19 1997 -- Scott Turner
 *
 *  Move permitted?
 *
 *  Mon May 24 12:52:20 1999 -- Scott Turner
 *
 *  Moves into sublocs now have to check their size.
 */
func move_permitted(c *command, v *exit_view) bool {
	if v.hades_cost != 0 {
		n := count_stack_any(c.who)
		for _, j := range loop_stack(c.who) {
			if has_artifact(j, ART_PROT_HADES, 0, 0, 0) != FALSE {
				wout(j, "Enter, daemon lord!")
				n--
			}
		}

		if n < 0 {
			n = 0
		}
		cost := n * v.hades_cost

		if cost != 0 {
			if !autocharge(c.who, cost) {
				wout(c.who, "Can't afford %s to enter Hades.", gold_s(cost))
				return false
			}
			wout(c.who, "The Gatekeeper Spirit of Hades took %s from us.", gold_s(cost))
		}

		log_output(LOG_SPECIAL, "%s (%s) tries to enter Hades", box_name(player(c.who)), box_name(c.who))
	}

	/*
	 *  Check and/or join guilds.
	 *
	 */
	if !check_guilds(c, v) {
		return false
	}

	/*
	 *  If the destination is a subloc with an entrance size,
	 *  then we'd better be that size.
	 *
	 */
	if entrance_size(v.destination) != 0 {
		/*
		 *  Use the special count function to also count ninjas
		 *  and angels.
		 *
		 */
		n := count_stack_any_real(c.who, false, false)
		if n > entrance_size(v.destination) {
			wout(c.who, "The entrance will pass only %s at a time.", nice_num(entrance_size(v.destination)))
			return false
		}
	}

	return true
}

/*
 *  Tue Sep  7 14:17:06 1999 -- Scott Turner
 *
 *  You can sneak into a province, but not anything smaller.
 *
 *  Mon Jan 15 10:25:18 2001 -- Scott Turner
 *
 *  Whoops, you'd better also be hidden!
 */
func can_sneak(who, where int) bool {
	return (char_hidden(who) != FALSE && province(where) == where && char_alone_stealth(who))
}

/*
 *  Mon Jan 15 10:43:50 2001 -- Scott Turner
 *
 *  This checks to see if "who" can enter "dest" peacefully from "from".
 *
 */
func peaceful_enter(who, from, where int) bool {
	ruler := controls_loc(where)
	pl := player_controls_loc(where)

	/*
	 *  Mon Mar 17 11:52:50 1997 -- Scott Turner
	 *
	 *  If someone rules here and is hostile to you, you won't
	 *  be admitted.
	 *
	 */
	if ruler != 0 && is_hostile(ruler, who) != FALSE && !can_sneak(who, where) {
		return false
	}

	/*
	 *  Thu Mar 20 12:33:41 1997 -- Scott Turner
	 *
	 *  If the border is closed and you aren't on the admit list,
	 *  then you are refused admission.  However, there needs to
	 *  actually be someone ruling the location.
	 *
	 *  Mon Mar  1 10:04:12 1999 -- Scott Turner
	 *
	 *  There's a problem here, in that this is checking "where", which
	 *  is probably not the right thing!
	 *
	 *  Mon May 17 08:10:56 1999 -- Scott Turner
	 *
	 *  Don't prevent people from leaving a mine if the border to the
	 *  enclosing province is closed.
	 *
	 *  Tue Sep  7 14:19:17 1999 -- Scott Turner
	 *
	 *  Let concealed people sneak into provinces.
	 *
	 */
	if pl != 0 && ruler != 0 &&
		((rp_loc(where) != nil && rp_loc(where).control.closed) || (rp_subloc(where) != nil && rp_subloc(where).control.closed)) &&
		FALSE == somewhere_inside(where, from) && subkind(from) != sub_mine && subkind(from) != sub_mine_notdone &&
		FALSE == will_admit(pl, who, ruler) && !can_sneak(who, where) {
		return false
	}
	return true
}

/*
 * Sun Mar 30 20:00:15 1997 -- Scott Turner
 *
 *  Check to see if a move can be accomplished w/o a fight.
 *
 */
func check_peaceful_move(c *command, v *exit_view) bool {
	cost := 0

	if !peaceful_enter(c.who, v.orig, v.destination) {
		wout(c.who, "You are refused admission to %s.", box_name(v.destination))
		return false
	}

	var pl int    // todo: pl is never updated; we said no bug fixes
	var ruler int // todo: ruler is never updated; we said no bug fixes
	//#if 0
	//    /*
	//     *  Mon Mar 17 11:52:50 1997 -- Scott Turner
	//     *
	//     *  If someone rules here and is hostile to you, you won't
	//     *  be admitted.
	//     *
	//     */
	//    ruler = controls_loc(v.destination);
	//    if (ruler && is_hostile(ruler,c.who) &&
	//        !can_sneak(c.who, v.destination)) {
	//      wout(c.who,"%s refused you admission to %s.",
	//       box_name(ruler),box_name(v.destination));
	//      return FALSE;
	//    }
	//
	//    /*
	//     *  Thu Mar 20 12:33:41 1997 -- Scott Turner
	//     *
	//     *  If the border is closed and you aren't on the admit list,
	//     *  then you are refused admission.  However, there needs to
	//     *  actually be someone ruling the location.
	//     *
	//     *  Mon Mar  1 10:04:12 1999 -- Scott Turner
	//     *
	//     *  There's a problem here, in that this is checking "where", which
	//     *  is probably not the right thing!
	//     *
	//     *  Mon May 17 08:10:56 1999 -- Scott Turner
	//     *
	//     *  Don't prevent people from leaving a mine if the border to the
	//     *  enclosing province is closed.
	//     *
	//     *  Tue Sep  7 14:19:17 1999 -- Scott Turner
	//     *
	//     *  Let concealed people sneak into provinces.
	//     *
	//     */
	//    pl = player_controls_loc(v.destination);
	//    if (pl &&
	//        ruler &&
	//        rp_loc(v.destination) &&
	//        rp_loc(v.destination).control.closed &&
	//        FALSE == somewhere_inside(v.destination, v.orig) &&
	//        subkind(v.orig) != sub_mine &&
	//        subkind(v.orig) != sub_mine_notdone &&
	//        FALSE == will_admit(pl, c.who, ruler) &&
	//        !can_sneak(c.who, v.destination)) {
	//      wout(c.who,"%s refused you admission to %s.",
	//       box_name(ruler),box_name(v.destination));
	//      return FALSE;
	//    } else if (pl &&
	//           ruler &&
	//           rp_subloc(v.destination) &&
	//           rp_subloc(v.destination).control.closed &&
	//           FALSE == somewhere_inside(v.destination, v.orig) &&
	//           subkind(v.orig) != sub_mine &&
	//           subkind(v.orig) != sub_mine_notdone &&
	//           FALSE == will_admit(pl, c.who, ruler) &&
	//           !can_sneak(c.who, v.destination)) {
	//      wout(c.who,"%s refused you admission to %s.",
	//       box_name(ruler),box_name(v.destination));
	//      return FALSE;
	//    }
	//#endif

	/*
	 *  Mon Mar 17 11:56:23 1997 -- Scott Turner
	 *
	 *  Entrance fees to locations.
	 *
	 *  Subloc fee if:
	 *
	 *   1 -- you're not coming from inside the subloc
	 *   2 -- fee is set.
	 *   3 -- subloc has an owner
	 *   4 -- owner is a different player (will_admit)
	 *   4 -- owner is not admitting you free.
	 *
	 */
	// todo: pl is always zero because of the code commented out above
	var control *loc_control_ent
	if pl != 0 && rp_subloc(v.destination) != nil &&
		FALSE == somewhere_inside(v.destination, v.orig) &&
		(rp_subloc(v.destination).control.men != 0 || rp_subloc(v.destination).control.weight != 0 || rp_subloc(v.destination).control.nobles != 0) &&
		FALSE == will_admit(pl, c.who, v.destination) && !can_sneak(c.who, v.destination) {
		control = &rp_subloc(v.destination).control
		/*
		 *  Loc fee if:
		 *
		 *   1 -- fee is set.
		 *   2 -- not coming from inside
		 *   2 -- loc has a garrison
		 *   3 -- castle is occupied.
		 *   4 -- you don't rule here.
		 *   5 -- coming from a loc and not a subloc
		 *   6 -- orig is not administered by same as this loc
		 *   7 -- you're not admitted free.
		 *
		 */
	} else if pl != 0 /* this also checks for a controlled garrison */ &&
		rp_loc(v.destination) != nil &&
		(rp_loc(v.destination).control.men != 0 || rp_loc(v.destination).control.weight != 0 || rp_loc(v.destination).control.nobles != 0) &&
		FALSE == will_admit(pl, c.who, v.destination) &&
		FALSE == may_rule_here(c.who, v.destination) &&
		kind(v.orig) == T_loc &&
		loc_depth(v.orig) == LOC_province &&
		province_admin(v.destination) != province_admin(v.orig) &&
		!can_sneak(c.who, v.destination) {
		control = &rp_loc(v.destination).control
	}

	// todo: ruler is always zero because of the code commented out above
	if ruler != 0 && control != nil && (control.men != 0 || control.nobles != 0 || control.weight != 0) {
		if FALSE == will_pay(c.who) {
			wout(VECT, "%s is charging for entrance and you are not entering any location charging for entrance.", box_name(v.destination))
			return false
		}
		cost = calc_entrance_fee(control, c, ruler)
		if cost < 0 {
			return false
		} else if will_pay(c.who) < cost {
			wout(VECT, "Refused to pay %s to enter %s, so you are turned away.",
				gold_s(cost), box_name(v.destination))
			return false
		} else if !charge_entrance_fees(c.who, ruler, cost) {
			return false
		}
	}

	return true
}

//#if 0
///*
// *  Sun Mar 30 18:37:03 1997 -- Scott Turner
// *
// *  Check for an "attack" flag and prepend the appropriate
// *  order if turned away...
// *
// */
//int
//v_move(struct command *c)
//{
//    var v *exit_view
//    var delay int
//    where := subloc(c.who);
//    check_outer := TRUE;
//    prepend := 0;
//
//    if (numargs(c) < 1)
//    {
//        wout(c.who, "Specify direction or destination to MOVE.");
//        return FALSE;
//    }
//
//    while (numargs(c) > 0)
//    {
//      v = parse_exit_dir(c, where, "move");
//
//      if (v) {
//        check_outer = FALSE;
//            if (land_check(c, v, TRUE)) break;
//        v = nil;
//      }
//
//      cmd_shift(c);
//    }
//
//    vector_stack(c.who, TRUE);
//
//    if (v == nil && check_outer && can_move_at_outer_level(where, c))
//    {
//        c.a = subloc(where);
//        v = parse_exit_dir(c, where, nil);
//        assert(v);
//
//        if (v)
//        {
//            assert(move_exit_land(c, v, FALSE) >= 0);
//            wout(c.who, "(assuming 'move out' first)");
//            prepend = 1;
//        }
//    }
//
//    if (v == nil)
//        return FALSE;
//
//    delay = move_exit_land(c, v, TRUE);
//
//    if (delay < 0)
//        return FALSE;
//
//    if (!move_permitted(c,v) ||
//        !check_peaceful_move(c,v)) {
//      return FALSE;
//    }
//
//    /*
//     *  Actual movement; encapsulated so it can also be
//     *  called from v_attack.
//     *
//     */
//    do_actual_move(c,v,delay);
//
//    if (prepend) prepend_order(player(c.who), c.who, c.line);
//    return TRUE;
//}
//#endif

func touch_loc_after_move(who, where int) {
	if kind(who) == T_char {
		touch_loc(who)
	}
	for _, i := range loop_char_here(who) {
		if !is_prisoner(i) {
			touch_loc(i)
		}
	}
}

func match_trades(i int) { panic("!implemented") }

func move_stack(who, where int) {

	assert(kind(who) == T_char)

	if !in_faery(subloc(who)) && in_faery(where) {
		log_output(LOG_SPECIAL, "%s enters Faery at %s.",
			box_name(who), box_name(where))
	}

	set_where(who, where)
	mark_loc_stack_known(who, where)
	touch_loc_after_move(who, where)
	update_weather_view_locs(who, where)
	clear_contacts(who)

	//#ifndef NEW_TRADE
	ifndef(NEW_TRADE, func() {
		if subkind(where) == sub_city || (subkind(where) == sub_guild && rp_subloc(where).guild == sk_trading) {
			for _, i := range loop_stack(who) {
				match_trades(i)
			}
		}
	})
	//#endif

	if subkind(where) == sub_ocean {
		p := p_char(who)
		if p != nil && p.time_flying == 0 {
			p.time_flying++
			ocean_chars = append(ocean_chars, who)
		}
	}

	if subkind(where) != sub_ocean {
		p := rp_char(who)
		if p != nil && p.time_flying != 0 {
			p.time_flying = 0
			ocean_chars = rem_value(ocean_chars, who)
		}
	}

	//#if 0
	//    /*
	//     *  Don't need new bandits "popping" up.
	//     *
	//     */
	//    if (loc_depth(where) == LOC_province &&
	//        subkind(where) != sub_ocean &&
	//        has_item(where, item_peasant) < 100 &&
	//        !in_hades(where) &&
	//        !in_faery(where) &&
	//        !in_clouds(where))
	//        wilderness_attack_check(who, where);
	//#endif
}

/*
 *  Traps
 *
 *  Encodes some common traps...
 *
 */
var traps = []trap_struct{
	{sk_eres,
		sk_eres,
		10,
		50,
		"water elemental",
		"A water elemental eyes you warily as you arrive.",
		"You see an odd spot of water in the ocean below.",
		"A water elemental attacks as you enter the province!"},

	{sk_dol,
		sk_dol,
		10,
		50,
		"snake",
		"A nest of snakes rattles menacingly as you pass by.",
		"You spot a nest of snakes in the province below.",
		"A nest of snakes attacks as you enter the province!"},

	{sk_anteus,
		sk_anteus,
		10,
		50,
		"boulder",
		"A precarious pile of boulders shifts slightly as you pass by.",
		"You spot a strange pile of boulders in the province below.",
		"A boulder trap collapses on you as you enter the province!"},

	{sk_timeid,
		sk_timeid,
		10,
		50,
		"deadfall",
		"A dead tree shifts slightly as you pass by.",
		"You spot an odd-looking dead tree in the province below.",
		"A deadfall drops on you as you enter the province!"},

	{sk_kireus,
		sk_kireus,
		10,
		50,
		"quicksand",
		"A pool of quicksand gurgles as you pass by.",
		"You spot a patch of odd-colored sand in the province below.",
		"You stumble into a patch of quicksand as you enter the province!"},

	{}}

/*
 *  Thu Sep 26 13:13:24 1996 -- Scott Turner
 *
 *  Check to see if there are any effects hanging on this province that will effect the arriving stack.
 *
 */
func check_arrival_effects(who, where int, flying bool) {
	vector_stack(who, true)
	/*
	 *  Generic Trap Check
	 *
	 *  If there is a generic trap on this location (ef_religion_trap), then
	 *  we use the subtype as an index into traps[] to figure out exactly what
	 *  should happen.
	 *
	 */
	for _, e := range loop_effects(where) {
		/*
		 *  Look for religion traps.
		 *
		 */
		if e.type_ == ef_religion_trap {
			/*
			 *  Look up the subtype in traps[]
			 *
			 */
			i := 0
			for ; traps[i].type_ != 0 && traps[i].type_ != e.subtype; i++ {
				//
			}
			/*
			 *  Find something?
			 *
			 */
			if traps[i].type_ != 0 {
				/*
				 *  Implement this trap.
				 *
				 */
				if priest_in_stack(who, traps[i].religion) {
					wout(VECT, traps[i].ignored)
					continue
				}
				/*
				 *  Ignored if you're flying.
				 *
				 */
				if flying {
					if rnd(1, 4) > 2 {
						wout(VECT, traps[i].flying)
					}
					continue
				}
				// #ifdef HERO
				if HERO {
					/*
					 *  Wed Nov 25 13:05:44 1998 -- Scott Turner
					 *
					 *  You might be a hero with "acute senses"
					 *
					 */
					if rnd(1, 100) < min(80, 15+2*skill_exp(who, sk_acute_senses)) {
						wout(VECT, traps[i].flying)
						continue
					}
				}
				//#endif // HERO

				/*
				 *  Otherwise you get attacked.  Note that the ship
				 *  is the "who" being passed in.
				 *
				 */
				wout(VECT, traps[i].attack)
				do_trap_attack(who, traps[i].num_attacks, traps[i].attack_chance)
				delete_effect(where, e.type_, e.subtype)
			}
		}
	}
}

//#if 0
//int
//d_move(struct command *c)
//{
//    var vv exit_view
//    v := &exit_view{}
//
//    restore_v_array(c, v);
//
//    if (v.road)
//        discover_road(c.who, subloc(c.who), v);
//
//    vector_stack(c.who, TRUE);
//    wout(VECT, "Arrival at %s.", box_name(v.destination));
//    if (loc_depth(v.destination) == LOC_province &&
//        viewloc(subloc(c.who)) != viewloc(v.destination) &&
//        weather_here(v.destination, sub_fog))
//    {
//        wout(VECT, "The province is blanketed in fog.");
//    }
//
//    if (loc_depth(v.destination) == LOC_province &&
//        viewloc(subloc(c.who)) != viewloc(v.destination) &&
//        weather_here(v.destination, sub_mist))
//    {
//        wout(VECT, "The province is covered with a dank mist.");
//    }
//
//    check_arrival_effects(c.who, v.destination, 0);
//
//    restore_stack_actions(c.who);
//
//#if 0
///*
// *  Stackmates who are executing commands have gotten a free day.
// *  Force a one evening delay into their command completion loop.
// */
//
//    if (v.distance > 0)
//    {
//        var i int
//
//        for _, j = range loop_char_here(c.who, i)
//        {
//            struct command *nc = rp_command(i);
//
//            if (nc.wait != 0)
//                nc.move_skip = TRUE;
//        }
//
//    }
//#endif
//
//    move_stack(c.who, v.destination);
//
//    if (viewloc(v.orig) != viewloc(v.destination))
//        arrival_message(c.who, v);
//
//    return TRUE;
//}
//#endif

func init_ocean_chars() {
	for _, i := range loop_char() {
		if where := subloc(i); subkind(where) == sub_ocean {
			ocean_chars = append(ocean_chars, i)
		}
	}
}

func check_ocean_chars() {
	l := ilist_copy(ocean_chars)
	for i := 0; i < len(l); i++ {
		who := l[i]
		where := subloc(who)
		p := p_char(who)
		if !alive(who) || subkind(where) != sub_ocean {
			p.time_flying = 0
			ocean_chars = rem_value(ocean_chars, who)
			continue
		}

		p.time_flying++
		if p.time_flying <= 15 {
			continue
		} else if stack_parent(who) != 0 {
			continue
		}

		// flying stack plunges into the sea
		vector_stack(who, true)
		wout(VECT, "Flight can no longer be maintained.  %s plunges into the sea.", box_name(who))

		kill_stack_ocean(who)
	}
}

// todo: should this do something?
func fly_check(c *command, v *exit_view) bool {
	return true
}

func can_fly_here(where int, c *command) bool {
	v := parse_exit_dir(c, where, "")
	return v != nil && v.direction != DIR_IN && move_exit_fly(c, v, false) >= 0
}

func can_fly_at_outer_level(where int, c *command) bool {
	for outer := subloc(where); loc_depth(outer) > LOC_region; outer = subloc(outer) {
		if can_fly_here(outer, c) {
			return loc_depth(outer) == loc_depth(where)
		}
	}
	return false
}

/*
 *  Wed Mar 26 08:21:51 1997 -- Scott Turner
 *
 *  Note that flying avoids all control -- closed borders and/or fees.
 *
 */
func v_fly(c *command) int {
	var delay int
	where := subloc(c.who)
	check_outer := true
	prepend := false

	if numargs(c) < 1 {
		wout(c.who, "Specify direction or destination to FLY.")
		return FALSE
	}

	var v *exit_view
	for numargs(c) > 0 {
		v = parse_exit_dir(c, where, "fly")
		if v != nil {
			check_outer = false
			if fly_check(c, v) {
				break
			}
			v = nil
		}
		cmd_shift(c)
	}

	if v == nil && check_outer && can_fly_at_outer_level(where, c) {
		c.a = subloc(where)
		v = parse_exit_dir(c, where, "")
		assert(v != nil)
		if v != nil {
			assert(move_exit_fly(c, v, false) >= 0)
			wout(c.who, "(assuming 'fly out' first)")
			prepend = true
		}
	}

	if v == nil {
		return FALSE
	}

	delay = move_exit_fly(c, v, true)
	if delay < 0 {
		return FALSE
	}

	/*
	 *  Check and/or join guilds.
	 *
	 */
	if !check_guilds(c, v) {
		return FALSE
	}

	v.distance = delay
	c.wait = delay

	save_v_array(c, v)
	leave_stack(c.who)

	if delay > 0 {
		vector_stack(c.who, true)
		wout(VECT, "Flying to %s will take %s day%s.", box_name(v.destination), nice_num(delay), or_string(delay == 1, "", "s"))
	}

	suspend_stack_actions(c.who)
	clear_guard_flag(c.who)

	departure_message(c.who, v)

	if prepend {
		prepend_order(player(c.who), c.who, c.line)
	}
	return TRUE
}

func d_fly(c *command) int {
	v := &exit_view{}

	restore_v_array(c, v)
	restore_stack_actions(c.who)

	// can he still fly
	var w weights
	determine_stack_weights(c.who, &w, false)
	if w.fly_cap < w.fly_weight {
		wout(c.who, "%s is too overloaded to fly.", box_name(c.who))
		wout(c.who, "You have a total of %d weight, and your maximum flying capacity is %d.", w.fly_weight, w.fly_cap)
		return FALSE
	}

	if c.wait != 0 {
		return TRUE
	}

	if v.road != 0 {
		discover_road(c.who, subloc(c.who), v)
	}

	vector_stack(c.who, true)
	wout(VECT, "Arrival at %s.", box_name(v.destination))

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(c.who)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_fog) != FALSE {
		wout(VECT, "The province is blanketed in fog.")
	}

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(c.who)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_mist) != FALSE {
		wout(VECT, "The province is covered with a dank mist.")
	}

	check_arrival_effects(c.who, v.destination, true)

	//#if 0
	//    /*
	//     *  Stackmates who are executing commands have gotten a free day.
	//     *  Force a one evening delay into their command completion loop.
	//     */
	//
	//        if (v.distance > 0)
	//        {
	//            var i int
	//
	//            for _, j = range loop_char_here(c.who, i)
	//            {
	//                struct command *nc = rp_command(i);
	//
	//                if (nc.wait != 0)
	//                    nc.move_skip = TRUE;
	//            }
	//
	//        }
	//#endif

	move_stack(c.who, v.destination)
	arrival_message(c.who, v)

	return TRUE
}

/*
 *  Synonym for 'move out'
 */

func v_exit(c *command) int {
	ret := oly_parse_s(c, "move out")
	assert(ret)
	return v_move_attack(c)
}

func v_enter(c *command) int {
	if numargs(c) < 1 {
		ret := oly_parse_s(c, "move in")
		assert(ret)
		return v_move_attack(c)
	}

	ret := oly_parse_s(c, sout("move %s", c.parse[1]))
	assert(ret)

	return v_move_attack(c)
}

func v_north(c *command) int {
	ret := oly_parse_s(c, "move north")
	assert(ret)

	return v_move_attack(c)
}

func v_south(c *command) int {
	ret := oly_parse_s(c, "move south")
	assert(ret)

	return v_move_attack(c)
}

func v_east(c *command) int {
	ret := oly_parse_s(c, "move east")
	assert(ret)

	return v_move_attack(c)
}

func v_west(c *command) int {
	ret := oly_parse_s(c, "move west")
	assert(ret)

	return v_move_attack(c)
}

func check_captain_loses_sailors(qty, target, inform int) {
	panic("this is broken")
	//    static int cmd_sail = -1;
	//    where := subloc(target);
	//    struct command *c;
	//    hands_short := 0;
	//    int before, now;
	//    int should_have;
	//    int penalty;
	//
	///* This is broken */
	//    return;
	//
	//    if (cmd_sail < 0) {
	//        cmd_sail = find_command("sail");
	//        assert(cmd_sail > 0);
	//    }
	//
	//    if (!is_ship(subloc(target))) {
	//        return;
	//    }
	//
	//    c = rp_command(target);
	//
	//    if (c == nil || c.state != RUN || c.cmd != cmd_sail) {
	//        return;
	//    }
	//
	//    now = has_item(target, item_sailor) + has_item(target, item_pirate);
	//    before = now + qty;
	//
	//    switch (subkind(where)) {
	//        case sub_galley:
	//            should_have = 14;
	//            break;
	//
	//        case sub_roundship:
	//            should_have = 8;
	//            break;
	//
	//        default:
	//            fprintf(stderr, "kind is %d\n", subkind(where));
	//            assert(FALSE);
	//    }
	//
	//    if (now >= should_have) {
	//        return;
	//    }        /* still have enough sailors */
	//
	//    if (before > should_have){
	//before = should_have;
	//}
	//    penalty = before - now;
	//
	//    assert(penalty > 0);
	//
	//    vector_clear();
	//    vector_add(target);
	//    if (inform && target != inform)
	//        vector_add(inform);
	//
	//    if (penalty == 1){
	//wout(VECT, "Loss of crew will cause travel to take an extra day.");
	//}else{
	//wout(VECT, "Loss of crew will cause travel to take %s extra days.", nice_num(penalty));}
	//
	//    assert(c.wait > 0);
	//    c.wait += penalty;
	//
	//    log_output(LOG_SPECIAL, "Loss of sailors incurs penalty for %s.",
	//               box_code(player(target)));
}

func move_exit_water(c *command, v *exit_view, ship int, show bool) int {
	delay := v.distance
	hands_short := 0 /* how many hands we are short */
	var n int
	var s string
	sail_time := v.distance
	row_time := v.distance
	wind_bonus := 0
	where := subloc(ship)
	sp := rp_ship(ship)
	var ports, sails int

	switch subkind(ship) {

	case sub_roundship:
		n = has_item(c.who, item_sailor) + has_item(c.who, item_pirate)
		if n < 8 {
			hands_short = 8 - n

			if hands_short == 1 {
				s = "day"
			} else {
				s = sout("%s days", nice_num(hands_short))
			}

			if show {
				wout(c.who, "The crew of a roundship is eight sailors, but you have %s.  Travel will take an extra %s.",
					or_string(n == 0, "none", nice_num(n)), s)
			}

			if weather_here(where, sub_wind) != FALSE && delay > 1 {
				wind_bonus = 1
				if show {
					wout(c.who, "Favorable winds speed our progress.")
				}
			}
			delay = delay + hands_short - wind_bonus
		}
		break

	case sub_galley:
		n = has_item(c.who, item_sailor) + has_item(c.who, item_pirate)
		/* n += has_item(c.who, item_slave); */

		if n < 14 {
			hands_short = 14 - n

			if hands_short == 1 {
				s = "day"
			} else {
				s = sout("%s days", nice_num(hands_short))
			}

			if show {
				wout(c.who, "The crew of a galley is fourteen pirates or sailors, but you have %s.  Travel will take an extra %s.",
					or_string(n == 0, "none", nice_num(n)), s)
			}
		}
		delay = delay + hands_short - wind_bonus
		break

		/*
		 *  Fri Jan  3 11:54:06 1997 -- Scott Turner
		 *
		 *  For the general-purpose ship, we calculate a row time and a
		 *  sail time and then use the lesser of the two.
		 *
		 */
	case sub_ship:
		if sp == nil {
			return sail_time
		}
		/*
		 *  Sailors available.
		 *
		 */
		n = has_item(c.who, item_sailor) + has_item(c.who, item_pirate) - max((sp.hulls/20)*2, 2)
		/*
		 *  Can we row?
		 *
		 */

		/*
		 *  How many manned ports do we have?
		 *
		 */
		ports = min(sp.ports, n/ROWERS_PER_PORT)

		/*
		 *  Do we have more or fewer ports than we need?
		 *
		 */
		if (ports * ROWERS_PER_PORT) < sp.hulls {
			/*
			 *  Need at least one oarsmen per unit hull.
			 *
			 */
			row_time = 1000
		} else if ports > sp.hulls {
			/*
			 *  Save one day per each multiple of sp.hulls over 1
			 *
			 */
			row_time -= (ports / sp.hulls) - 1
		} else if ports < sp.hulls {
			row_time += sp.hulls - ports
		}

		if row_time < 1 {
			row_time = 1
		}
		/*
		 *  Can we sail?
		 *
		 */

		/*
		 *  How many sails available?
		 *
		 */
		sails = min(sp.sails, n)

		/*
		 *  Do we have more or fewer sails than we need?
		 *
		 */
		if sails < sp.hulls/2 {
			/*
			 *  Need at least 1 sail per 2 hull to move at all.
			 *
			 */
			sail_time = 1000
		} else if sails > 2*sp.hulls {
			/*
			 *  Save one day for each full sail per unit hull
			 *
			 */
			sail_time -= ((int)((sails - 2*sp.hulls) / sp.hulls)) * 2
		} else if sails < 2*sp.hulls {
			sail_time += 2
		}

		/*
		 *  Wind bonus
		 *
		 */
		if weather_here(where, sub_wind) != FALSE && sail_time > 1 && sail_time != 1000 {
			sail_time--
			if show {
				wout(c.who, "Favorable winds speed our progress.")
			}
		}
		if sail_time < 1 {
			sail_time = 1
		}

		/*
		 *  Which is it?
		 *
		 */
		if sail_time == 1000 && row_time == 1000 {
			if show {
				wout(c.who, "The ship is unable to be rowed or sailed at this time.")
			}
			return -1
		} else if sail_time < row_time {
			if show {
				wout(c.who, "The sails are unfurled.")
			}
			return sail_time
		}
		if show {
			wout(c.who, "The rowers dip their oars to the cadence of the drums.")
		}
		return row_time

	default:
		panic(fmt.Sprintf("assert(subkind != %d)", subkind(ship)))
	}

	return delay
}

func sail_depart_message(ship int, v *exit_view) {
	var to string
	var desc string
	var comma string

	desc = liner_desc(ship)

	if v.dest_hidden == FALSE {
		to = sout(" for %s.", box_name(v.destination))
	}

	if strings.IndexByte(desc, ',') != -1 {
		comma = ","
	}

	wout(v.orig, "%s%s departed%s", desc, comma, to)
}

func sail_arrive_message(ship int, v *exit_view) {
	var from string
	var desc string
	var comma string
	var with string

	desc = liner_desc(ship)

	if v.orig_hidden == FALSE {
		from = sout(" from %s", box_name(v.orig))
	}

	if strings.IndexByte(desc, ',') != -1 {
		comma = ","
	}

	with = display_owner(ship)
	if with != "" {
		with = "."
	}

	show_to_garrison = true

	wout(v.destination, "%s%s arrived%s%s", desc, comma, from, with)
	show_owner_stack(v.destination, ship)

	show_to_garrison = false
}

func sail_check(c *command, v *exit_view, show bool) bool {

	if v.water == FALSE {
		if show {
			wout(c.who, "There is no water route in that direction.")
		}
		return false
	}

	if v.impassable != FALSE {
		if show {
			wout(c.who, "That route is impassable.")
		}
		return false
	}

	return true
}

func can_sail_here(where int, c *command, ship int) bool {
	v := parse_exit_dir(c, where, "")
	if v != nil && v.direction != DIR_IN &&
		sail_check(c, v, false) &&
		move_exit_water(c, v, ship, false) >= 0 {
		return true
	}

	return false
}

func can_sail_at_outer_level(ship, where int, c *command) bool {
	if ship_cap(ship) != 0 {
		loaded := ship_weight(ship) * 100 / ship_cap(ship)
		if loaded > 100 {
			wout(c.who, "%s is too overloaded to sail.",
				box_name(ship))
			return false
		}
	}

	outer := subloc(where)
	for loc_depth(outer) > LOC_region {
		if can_sail_here(outer, c, ship) {
			return loc_depth(outer) == loc_depth(where)
		}
		outer = subloc(outer)
	}

	return false
}

func v_sail(c *command) int {
	var v *exit_view
	var delay int
	ship := subloc(c.who)
	var outer_loc, rocky_coast int
	check_outer := true
	var result int

	if !is_ship(ship) {
		if is_ship_notdone(ship) {
			wout(c.who, "%s is not yet completed.",
				box_name(ship))
		} else {
			wout(c.who, "Must be on a sea-worthy ship to sail.")
		}
		return FALSE
	}

	if building_owner(ship) != c.who {
		wout(c.who, "Only the captain of a ship may sail.")
		return FALSE
	}

	if has_skill(c.who, sk_pilot_ship) <= 0 {
		wout(c.who, "Knowledge of %s is required to sail.",
			box_name(sk_pilot_ship))
		return FALSE
	}

	if numargs(c) < 1 {
		wout(c.who, "Specify direction or destination to sail.")
		return FALSE
	}

	outer_loc = subloc(ship)
	/*
	 *  Wed Jan  8 12:07:09 1997 -- Scott Turner
	 *
	 *  We have to do this *before* the parse_exit_dir below, because
	 *  they both use a static local in exits_from_loc, and since v is
	 *  saved here, bad things happen...
	 *
	 */
	rocky_coast = near_rocky_coast(outer_loc)
	for numargs(c) > 0 {
		v = parse_exit_dir(c, outer_loc, "sail")
		if v != nil {
			check_outer = false

			if sail_check(c, v, true) {
				break
			}
			v = nil
		}

		cmd_shift(c)
	}

	if v == nil && check_outer && can_sail_at_outer_level(ship, outer_loc, c) {
		c.a = subloc(outer_loc)
		v = parse_exit_dir(c, outer_loc, "")
		assert(v != nil)

		if v != nil {
			assert(move_exit_water(c, v, ship, false) >= 0)
			wout(c.who, "(assuming 'sail out' first)")
			prepend_order(player(c.who), c.who, c.line)
		}
	}

	if v == nil {
		wout(c.who, "No valid arguments to sail.")
		return FALSE
	}

	if ship_cap(ship) != 0 {
		loaded := ship_weight(ship) * 100 / ship_cap(ship)
		if loaded > 100 {
			wout(c.who, "%s is too overloaded to sail.", box_name(ship))
			return FALSE
		}
	}

	assert(v.in_transit == FALSE)

	delay = move_exit_water(c, v, ship, true)
	if delay < 0 {
		return FALSE
	}

	/*
	 *  Tue Jan  7 15:28:02 1997 -- Scott Turner
	 *
	 *  Possibility of problems if you're sailing out of a deep
	 *  water province and you don't have Deep Sea Navigation.
	 *
	 *  Tue Dec 15 08:29:33 1998 -- Scott Turner
	 *
	 *  The Pen don't get lost.
	 *
	 */
	if rocky_coast == 3 &&
		has_skill(c.who, sk_deep_sea) <= 0 &&
		!(nation(first_character(ship)) != FALSE &&
			strncmp(rp_nation(nation(first_character(ship))).name, "Pen", 3) == 0) &&
		FALSE == has_artifact(c.who, ART_SAFE_SEA, 0, 0, 0) &&
		rnd(1, 3) == 1 {
		wout(c.who, "The lack of landmarks to sail by confuses you.")
		result = rnd(1, 100) - or_int(priest_in_stack(c.who, sk_eres), 10, 0)
		if result <= 40 {
			tmp := rnd(3, 10)
			wout(c.who, "It takes %s days to reorient yourself.",
				nice_num(tmp))
			delay += tmp
		} else if result <= 80 {
			/*
			 *  Lost; so move like an NPC (randomly).
			 *
			 */
			choices := exits_from_loc_nsew_select(c.who, outer_loc, WATER, true)

			if len(choices) != 0 {
				if choices[0].destination != v.destination {
					wout(c.who, "Somehow you sail off in an unintended direction.")
					v = choices[0]
					delay = move_exit_water(c, v, ship, true)
				}
			}
		} else if result <= 95 {
			wout(c.who, "You sail into rough seas that damage the ship.")
			add_structure_damage(ship, rnd(1, 20))
		} else {
			wout(c.who, "In your confusion you somehow swamp the ship!")
			add_structure_damage(ship, 100)
		}
	}

	/*
	 *  Wed Jan  8 12:08:48 1997 -- Scott Turner
	 *
	 *  Ship might have been destroyed!
	 *
	 */
	if !valid_box(ship) {
		return FALSE
	}

	c.wait = delay
	v.distance = delay

	save_v_array(c, v)

	if delay > 0 {
		vector_char_here(c.who)
		vector_add(c.who)

		wout(VECT, "Sailing to %s will take %s day%s.",
			box_name(v.destination),
			nice_num(delay),
			or_string(delay == 1, "", "s"))
	}

	sail_depart_message(ship, v)

	/*
	 *  Mark the in_transit field with the daystamp of the beginning
	 *  of the voyage
	 */

	p_subloc(ship).moving = sysclock.days_since_epoch

	if ferry_horn(ship) { /* clear ferry horn signal */
		p_magic(ship).ferry_flag = false
	}

	return TRUE
}

/*
 *  Tue Oct  6 18:01:38 1998 -- Scott Turner
 *
 *  This now must be polled for deep sea effects!
 *
 */
func d_sail(c *command) int {
	ship := subloc(c.who)
	v := &exit_view{}

	assert(is_ship(ship))                 /* still a ship */
	assert(building_owner(ship) == c.who) /* still captain */

	restore_v_array(c, v)

	if v.road != FALSE {
		discover_road(c.who, subloc(ship), v)
	}

	vector_char_here(ship)
	wout(VECT, "Arrival at %s.", box_name(v.destination))

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(ship)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_fog) != FALSE {
		wout(VECT, "The province is blanketed in fog.")
	}

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(ship)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_mist) != FALSE {
		wout(VECT, "The province is covered in a dank mist.")
	}

	p_subloc(ship).moving = 0 /* no longer moving */
	set_where(ship, v.destination)
	mark_loc_stack_known(ship, v.destination)
	//#if 0
	//    move_bound_storms(ship, v.destination);
	//#endif
	if ferry_horn(ship) { /* clear ferry horn signal */
		p_magic(ship).ferry_flag = false
	}

	touch_loc_after_move(ship, v.destination)
	check_arrival_effects(ship, v.destination, false)
	sail_arrive_message(ship, v)

	if c.use_skill == 0 {
		add_skill_experience(c.who, sk_pilot_ship)
	}

	ship_storm_check(ship) /* Might destroy ship? */

	return TRUE
}

/*
 *  If sailing is interrupted, we must zero subloc.moving
 *  to indicate that the ship is no longer in transit.
 */

func i_sail(c *command) int {
	ship := subloc(c.who)

	assert(is_ship(ship))

	p_subloc(ship).moving = 0

	if ferry_horn(ship) { /* clear ferry horn signal */
		p_magic(ship).ferry_flag = false
	}

	return TRUE
}

/*
 *  Thu Apr  3 19:42:14 1997 -- Scott Turner
 *
 *  This version doesn't look for locations.
 *
 */
func select_target_local(c *command) int {
	target := c.a
	where := subloc(c.who)

	if kind(target) == T_char {
		if !check_char_here(c.who, target) {
			return 0
		}

		if loc_depth(where) == LOC_province &&
			weather_here(where, sub_fog) != FALSE &&
			!contacted(target, c.who) {
			wout(c.who, "That target is not visible in the fog.")
			return 0
		}

		if is_prisoner(target) {
			wout(c.who, "Cannot attack prisoners.")
			return 0
		}

		if c.who == target {
			wout(c.who, "Can't attack oneself.")
			return 0
		}

		if stack_leader(c.who) == stack_leader(target) {
			wout(c.who, "Can't attack a member of the same stack.")
			return 0
		}

		return stack_leader(target)
	}

	return 0
}

func attack_okay(c *command, target int) bool {
	if stack_leader(c.who) != c.who {
		wout(c.who, "Only the stack leader may initiate combat.")
		return false
	}

	if target <= 0 {
		wout(c.who, "You must specify a target to attack.")
		return false
	}

	attacker := select_attacker(c.who, target)
	if attacker <= 0 {
		wout(c.who, "You must specify a target to attack.")
		return false
	}

	var targ_who int
	if is_loc_or_ship(target) {
		targ_who = loc_target(target)
	} else {
		targ_who = target
	}

	/*
	 *  Target should be a character?
	 *
	 */
	if kind(targ_who) != T_char {
		wout(c.who, "Nothing there to attack.")
		return false
	}

	if !is_real_npc(c.who) && player(c.who) == player(targ_who) {
		wout(c.who, "Units in the same faction may not engage in combat.")
		return false
	}

	/*
	 *  Could be gone?
	 *
	 *  Mon Oct  5 18:58:37 1998 -- Scott Turner
	 *
	 *  Note difference from char_gone!
	 *
	 */
	//#if 0
	//    if (!is_loc_or_ship(target) && char_gone(target)) {
	//        wout(c.who, "%s has already left.",box_name(target));
	//        return FALSE;
	//      }
	//
	//    if (is_ship(target) && ship_gone(target)) {
	//        wout(c.who, "The ship slips away before you can attack it!");
	//        return FALSE;
	//      }
	//
	//    if (subloc(c.who) != subloc(target)) {
	//        wout(c.who, "%s is not here.",box_name(target));
	//        return FALSE;
	//    }
	//#endif

	if char_really_hidden(target) && !contacted(target, c.who) {
		wout(c.who, "%s is not here.", box_name(target))
		return false
	}

	if is_prisoner(targ_who) {
		wout(c.who, "Cannot attack prisoners.")
		return false
	}

	if c.who == targ_who {
		wout(c.who, "Can't attack oneself.")
		return false
	}

	if stack_leader(c.who) == stack_leader(targ_who) {
		wout(c.who, "Can't attack a member of the same stack.")
		return false
	}

	/*
	 *  Can't attack because we need an item?
	 *
	 */
	n := only_defeatable(targ_who)
	if n != 0 && FALSE == has_item(c.who, n) {
		wout(c.who, "To defeat %s you need %s.", box_name(targ_who),
			box_name(n))
		return false
	}

	/*
	 *  If the destination is a subloc with an entrance size,
	 *  then we'd better be that size.
	 *
	 */
	if entrance_size(target) != 0 {
		/*
		 *  Use the special count function to also count ninjas
		 *  and angels.
		 *
		 */
		n := count_stack_any_real(c.who, false, false)
		if n > entrance_size(target) {
			wout(c.who, "The entrance will pass only %s at a time.",
				nice_num(entrance_size(target)))
			return false
		}
	}

	return true
}

/*
 *  Tue Apr  1 12:22:42 1997 -- Scott Turner
 *
 *  So far doesn't have
 *     - "no seize" flag for attack use
 *     - non-location targets
 *
 *  Fri Dec  3 06:10:20 1999 -- Scott Turner
 *
 *  Need to add in restrictions for entrance size.  We need to add
 *  a check in move_permitted (which will stop peaceful moves through
 *  too small an opening) and in attack_okay (which will stop attacks
 *  through too small an opening).
 *
 *  Wed Nov  1 13:48:40 2000 -- Scott Turner
 *
 *  Fold in flying.
 *
 */
func v_move_attack(c *command) int {
	var delay int
	where := subloc(c.who)
	//check_outer := true;
	//prepend := 0;

	/*
	 *  Check for the flag that says not to seize the slot.
	 *
	 */
	seize := true
	if numargs(c) != 0 && atoi_b(c.parse[numargs(c)]) == 1 {
		seize = false
	}

	if numargs(c) < 1 {
		wout(c.who, "Specify direction or destination to %s.",
			c.parse[0])
		return FALSE
	}

	attack := strcasecmp_bs(c.parse[0], "attack") == 0

	fly := strcasecmp_bs(c.parse[0], "fly") == 0

	var v *exit_view
	for numargs(c) > 0 {
		v := parse_exit_dir(c, where, string(c.parse[0]))
		if v != nil {
			//check_outer = false;
			/*
			 *  Accept this move if:
			 *
			 *  (1) we can get there overland.
			 *  (2) we're attacking OR
			 *      we're peaceful and this move is permitted.
			 *
			 *  In d_move, we'll avoid making an unpermitted move after
			 *  an attack.  So we can attack into somewhere we can't move,
			 *  we just can't follow up by occupying that location.
			 *
			 *  Wed Nov  1 13:51:53 2000 -- Scott Turner
			 *
			 *  Updated for flying.
			 */
			if (fly && fly_check(c, v)) ||
				(!fly && land_check(c, v, true)) &&
					(attack || (move_permitted(c, v) && check_peaceful_move(c, v))) {
				break
			}
			v = nil
		}

		/*
		 *  Maybe he's trying to attack someone here?
		 *
		 */
		target := select_target_local(c)
		if target != 0 {
			/*
			 *  V is only needed till the end of the routine.
			 *
			 */
			v = &exit_view{}
			v.destination = target
			v.distance = 0
		}

		cmd_shift(c)
	}

	if v == nil {
		if attack {
			wout(c.who, "No target found!", c.parse[0])
		} else {
			wout(c.who, "No valid arguments to %s.", c.parse[0])
		}
		return FALSE
	}

	v.seize = or_int(seize, TRUE, FALSE) /* Hacky place to keep it :-) */

	vector_stack(c.who, true)

	/*
	 *  If v.distance = 0, this returns 0, so it works for the phony
	 *  attack exit_view.
	 *
	 */
	delay = move_exit_land(c, v, true)

	if delay < 0 {
		return FALSE
	}

	/*
	 *  An attack costs 1 day to prepare.
	 *
	 */
	if attack {
		delay++
	}

	/*
	 *  Attack might not be okay?
	 *
	 */
	if attack && !attack_okay(c, v.destination) {
		/*
		 *  In that case, let's attempt this as a peaceful move.
		 *
		 *  Tue Sep 21 12:54:48 1999 -- Scott Turner
		 *
		 *  No, let's fail.  The player can always get this functionality
		 *  with an attack/move comibination.
		 */
		return FALSE
	}

	/*
	 *  Actual movement; encapsulated so it can also be
	 *  called from v_attack.
	 *
	 */
	do_actual_move(c, v, delay)

	return TRUE
}

func d_move_attack(c *command) int {
	v := &exit_view{}
	attacker := 0
	result := 0
	delay := 0

	attack := strcasecmp_bs(c.parse[0], "attack") == 0

	restore_v_array(c, v)
	restore_stack_actions(c.who)

	if !valid_box(v.destination) {
		wout(c.who, "Your destination no longer exists!")
		return FALSE
	}

	/*
	 *  Maybe you are now overloaded?
	 *
	 */
	delay = move_exit_land(c, v, false)

	if delay < 0 {
		return FALSE
	}

	/*
	 *  Check for a forced march.
	 *
	 */
	if v.forced_march == FORCED_MARCH {
		wout(c.who, "Forced march costs %s health point%s.",
			nice_num(c.f),
			or_string(c.f == 1, "", "s"))
		add_char_damage(c.who, c.f, MATES)
	} else if v.forced_march == FORCED_RIDE {
		if rnd(1, 100) < (c.f * 20) {
			kill_random_mount(c.who)
		}
	}

	/*
	 *  Remove the effect if it should have been used.
	 *
	 */
	if c.f > 0 {
		delete_effect(c.who, ef_forced_march, 0)
	}

	if attack {
		/*
		 *  Attack might not be okay?
		 *
		 *  Tue Sep 21 12:45:31 1999 -- Scott Turner
		 *
		 *  If the attack is "not okay", then you should just automatically
		 *  win the battle, since what this means is that while you were moving
		 *  the enemy did something to make your attack irrelevant -- like
		 *  ungarrisoning the province.
		 *
		 */
		target := v.destination
		if target == 0 {
			wout(c.who, "You didn't specify anyone to attack.")
			return FALSE
		}

		/*
		 *  If the attack is okay, then go ahead and run the battle.  If
		 *  it isn't okay, then you "win".
		 *
		 */
		if attack_okay(c, target) {
			attacker = select_attacker(c.who, target)
			result = regular_combat(attacker, target, v.seize, c.who)
			if result == 0 {
				return FALSE
			}
		}

		/*
		 *  Seize the slot?
		 */
		if v.seize == FALSE {
			return TRUE
		}
	}

	if !is_loc_or_ship(v.destination) {
		return TRUE
	}

	if !move_permitted(c, v) {
		wout(c.who, "You cannot move into that location.")
		return FALSE
	}

	if v.road != FALSE {
		discover_road(c.who, subloc(c.who), v)
	}

	vector_stack(c.who, true)
	wout(VECT, "Arrival at %s.", box_name(v.destination))
	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(c.who)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_fog) != FALSE {
		wout(VECT, "The province is blanketed in fog.")
	}

	if loc_depth(v.destination) == LOC_province &&
		viewloc(subloc(c.who)) != viewloc(v.destination) &&
		weather_here(v.destination, sub_mist) != FALSE {
		wout(VECT, "The province is covered with a dank mist.")
	}

	check_arrival_effects(c.who, v.destination, true)

	move_stack(c.who, v.destination)

	if attack && result == A_WON {
		promote(c.who, 0)
	}

	if viewloc(v.orig) != viewloc(v.destination) {
		arrival_message(c.who, v)
	}

	return TRUE
}

func v_maxpay(c *command) int {
	m := c.a
	if rp_char(c.who) == nil {
		wout(c.who, "Oddly enough, you cannot use the maxpay command.")
		return FALSE
	}
	rp_char(c.who).pay = m
	wout(c.who, "Maximum amount paid to enter a location set to %s.", nice_num(m))

	return TRUE
}
