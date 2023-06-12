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

/*
*  Wed Oct  9 14:20:52 1996 -- Scott Turner
*
*  Doesn't seem to exist in the game?
*
 */
func v_bind_storm(c *command) int {
	storm := c.a
	ship := subloc(c.who)

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "%s doesn't control any storm %s.",
			box_name(c.who),
			box_code(storm))
		return FALSE
	}

	if !is_ship(ship) {
		wout(c.who, "%s must be on a ship to bind the storm to.",
			box_name(c.who))
		return FALSE
	}

	if !has_piety(c.who, 3) {
		return FALSE
	}

	return TRUE
}

func d_bind_storm(c *command) int {
	storm := c.a
	ship := subloc(c.who)

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "%s doesn't control storm %s anymore.",
			box_name(c.who),
			box_code(storm))
		return FALSE
	}

	if !is_ship(ship) {
		wout(c.who, "%s is no longer on a ship.", box_name(c.who))
		return FALSE
	}

	if !use_piety(c.who, 3) {
		return FALSE
	}

	old := storm_bind(storm)
	if old != FALSE {
		p := rp_subloc(old)
		if p != nil {
			p.bound_storms = rem_value(p.bound_storms, storm)
		}
	}

	p_misc(storm).bind_storm = storm
	p := p_subloc(ship)
	p.bound_storms = append(p.bound_storms, storm)

	wout(c.who, "Bound %s to %s.", box_name(storm), box_name(ship))
	return TRUE
}

func move_storm(storm, dest int) {
	orig := subloc(storm)
	sk := subkind(storm)
	var before, i int
	var owner int

	before = weather_here(dest, int(sk))

	set_where(storm, dest)

	owner = npc_summoner(storm)

	if owner != FALSE {
		touch_loc_pl(player(owner), dest)
	}

	show_to_garrison = true

	if weather_here(orig, int(sk)) == 0 {
		switch sk {
		case sub_rain:
			wout(orig, "It has stopped raining.")
			break

		case sub_wind:
			wout(orig, "It is no longer windy.")
			break

		case sub_fog:
			wout(orig, "The fog has cleared.")
			break

		case sub_mist:
			wout(orig, "The mist dissipates.")
			break

		default:
			panic("!reached")
		}
	}

	if FALSE == before {
		switch sk {
		case sub_rain:
			wout(dest, "It has begun to rain.")
			break

		case sub_wind:
			wout(dest, "It has become quite windy.")
			break

		case sub_fog:
			wout(dest, "It has become quite foggy.")
			break

		case sub_mist:
			wout(dest, "A dank mist has risen from the ground.")
			break
		}
	}

	show_to_garrison = false

	/*
	 *  Tue Jan  7 12:25:57 1997 -- Scott Turner
	 *
	 *  A storm has potential to do damage to a ship only if it is a rain
	 *  storm and there previously was no rain; this avoids hitting a ship
	 *  multiply with the existing storms every time a new storm wanders
	 *  by.
	 *
	 */
	if before != sub_rain {
		for _, i = range loop_here(dest) {
			if is_ship(i) {
				ship_storm_check(i)
			}
		}
	}
}

//#if 0
//void
//move_bound_storms(int ship, int where)
//{
//    struct entity_subloc *p;
//    var i int
//    int storm;
//
//    p = rp_subloc(ship);
//    if (p == nil)
//        return;
//
//    for i =  0; i < ilist_len(p.bound_storms); i++
//    {
//        storm = p.bound_storms[i];
//        if (kind(storm) != T_storm)
//        {
//            ilist_delete(&p.bound_storms, storm);
//            i--;
//            continue;
//        }
//
//        move_storm(storm, where);
//    }
//}
//#endif

func new_storm(newStorm, sk, aura, where int) int {
	var before, i int

	assert(sk == sub_rain || sk == sub_wind || sk == sub_fog || sk == sub_mist)
	assert(loc_depth(where) == LOC_province)

	before = weather_here(where, int(sk))

	if newStorm == 0 {
		newStorm = new_ent(T_storm, sk)

		if newStorm <= 0 {
			return -1
		}
	}

	p_misc(newStorm).storm_str = aura
	set_where(newStorm, where)

	show_to_garrison = true

	if FALSE == before {
		switch sk {
		case sub_rain:
			wout(where, "It has begun to rain.")
			break

		case sub_wind:
			wout(where, "It has become quite windy.")
			break

		case sub_fog:
			wout(where, "It has become quite foggy.")
			break

		case sub_mist:
			wout(where, "The province is shrouded in a dank mist.")
			break
		}
	}

	show_to_garrison = false

	/*
	 *  Tue Jan  7 12:25:57 1997 -- Scott Turner
	 *
	 *  A storm has potential to do damage to a ship only if it is a rain
	 *  storm and there previously was no rain; this avoids hitting a ship
	 *  multiply with the existing storms every time a new storm wanders
	 *  by.
	 *
	 */
	if before != sub_rain {
		for _, i = range loop_here(where) {
			if is_ship(i) {
				ship_storm_check(i)
			}
		}
	}
	return TRUE
}

func storm_report(pl int) {
	first := TRUE
	var owner int
	var where int
	var i int

	for _, i = range loop_storm() {
		owner = npc_summoner(i)

		if owner == 0 || player(owner) != pl {
			continue
		}

		where = province(i)

		/*
		 *  Tue Jun 10 11:54:48 1997 -- Scott Turner
		 *
		 *  Check for crossing water...
		 *
		 */
		if crosses_ocean(owner, i) {
			continue
		}

		if first != FALSE {
			if options.output_tags != FALSE {
				out(pl, "<tag type=storm_report pl=%d>", pl)
			}
			out(pl, "")
			out(pl, "%5s  %4s  %5s  %4s  %s",
				"storm", "kind", "owner", "loc", "strength")
			out(pl, "%5s  %4s  %5s  %4s  %s",
				"-----", "----", "-----", "----", "--------")

			first = FALSE
		}

		if options.output_tags != FALSE {
			out(pl, "<tag type=storm storm=%d kind=%s owner=%d loc=%d strength=%d>",
				i, subkind_s[subkind(i)], owner, where,
				storm_strength(i))
		}

		out(pl, "%5s  %4s  %5s  %4s     %s",
			box_code_less(i),
			subkind_s[subkind(i)],
			box_code_less(owner),
			box_code_less(where),
			comma_num(storm_strength(i)))

		tagout(pl, "</tag type=storm>")
	}

	if FALSE == first {
		tagout(pl, "</tag type=storm_report>")
	}
}

func dissipate_storm(storm int, show bool) {
	var owner int
	var p *entity_misc
	where := subloc(storm)

	assert(kind(storm) == T_storm)

	owner = npc_summoner(storm)

	/*
	 *  Tue Jun 10 11:56:04 1997 -- Scott Turner
	 *
	 *  Same region...
	 *
	 */
	if owner != FALSE && kind(owner) == T_char &&
		region(owner) == region(storm) {
		wout(owner, "%s has dissipated.", box_name(storm))
	}

	if show {
		sk := subkind(storm)

		if weather_here(where, int(sk)) == 0 {
			switch sk {
			case sub_rain:
				wout(where, "It has stopped raining.")
				break

			case sub_wind:
				wout(where, "It is no longer windy.")
				break

			case sub_fog:
				wout(where, "The fog has cleared.")
				break

			case sub_mist:
				wout(where, "The mist has dissipated.")
				break

			default:
				panic("!reached")
			}
		}
	}

	set_where(storm, 0)

	p = p_misc(storm)

	if p.npc_home != FALSE && p.npc_cookie != FALSE {
		gen_item(p.npc_home, p.npc_cookie, 1)
	}

	//#if 0
	//    if (ship = storm_bind(storm)) {
	//        p := rp_subloc(ship);
	//        if (p) {
	//            ilist_rem_value(&p.bound_storms, storm);
	//        }
	//        rp_misc(storm).bind_storm = 0;
	//    }
	//#endif

	delete_box(storm)
}

func weather_here(where, sk int) int {
	var i int
	sum := 0

	if loc_depth(where) == LOC_build {
		return 0
	}

	where = province(where)

	if FALSE == where {
		return 0
	}

	for _, i = range loop_here(where) {
		if kind(i) == T_storm && subkind(i) == schar(sk) {
			sum += storm_strength(i)
		}
	}

	return sum
}

func v_summon_rain(c *command) int {
	aura := c.a
	var where int

	where = province(cast_where(c.who))
	c.d = where

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if aura < 3 {
		c.a = 3
		aura = 3
	}

	if !may_cookie_npc(c.who, where, item_rain_cookie) {
		return FALSE
	}

	if !has_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	return TRUE
}

func d_summon_rain(c *command) int {
	aura := c.a
	where := c.d
	name := c.parse[2]

	if !may_cookie_npc(c.who, where, item_rain_cookie) {
		return FALSE
	}

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if !use_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	newCookie := do_cookie_npc(c.who, where, item_rain_cookie, where)

	if newCookie <= 0 {
		wout(c.who, "Failed to summon a storm.")
		return FALSE
	}

	reset_cast_where(c.who)

	if len(name) != 0 {
		set_name(newCookie, string(name))
	}

	new_storm(newCookie, sub_rain, aura*2, where)

	wout(c.who, "Summoned %s.", box_name_kind(newCookie))

	return TRUE
}

func v_summon_wind(c *command) int {
	aura := c.a
	var where int

	where = province(cast_where(c.who))
	c.d = where

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if aura < 3 {
		c.a = 3
		aura = 3
	}

	if !may_cookie_npc(c.who, where, item_wind_cookie) {
		return FALSE
	}

	if !has_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	return TRUE
}

func d_summon_wind(c *command) int {
	aura := c.a
	name := c.parse[2]
	where := c.d

	if !may_cookie_npc(c.who, where, item_wind_cookie) {
		return FALSE
	}

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if !use_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	newCookie := do_cookie_npc(c.who, where, item_wind_cookie, where)

	if newCookie <= 0 {
		wout(c.who, "Failed to summon a storm.")
		return FALSE
	}

	reset_cast_where(c.who)

	if len(name) != 0 {
		set_name(newCookie, string(name))
	}

	new_storm(newCookie, sub_wind, aura*2, where)

	wout(c.who, "Summoned %s.", box_name_kind(newCookie))

	return TRUE
}

func v_summon_fog(c *command) int {
	aura := c.a
	var where int

	where = province(cast_where(c.who))
	c.d = where

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if aura < 3 {
		c.a = 3
		aura = 3
	}

	if !may_cookie_npc(c.who, where, item_fog_cookie) {
		return FALSE
	}

	if !has_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	return TRUE
}

func d_summon_fog(c *command) int {
	aura := c.a
	name := c.parse[2]
	where := c.d

	if !may_cookie_npc(c.who, subloc(c.who), item_fog_cookie) {
		return FALSE
	}

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if !use_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	newCookie := do_cookie_npc(c.who, where, item_fog_cookie, where)

	if newCookie <= 0 {
		wout(c.who, "Failed to summon a storm.")
		return FALSE
	}

	reset_cast_where(c.who)

	if len(name) != 0 {
		set_name(newCookie, string(name))
	}

	new_storm(newCookie, sub_fog, aura*2, where)

	wout(c.who, "Summoned %s.", box_name_kind(newCookie))

	return TRUE
}

func parse_storm_dir(c *command, storm int) *exit_view {
	where := subloc(storm)
	var l []*exit_view
	var i int
	var dir int

	l = exits_from_loc_nsew(c.who, where)

	if valid_box(c.a) {
		if where == c.a {
			wout(c.who, "%s is already in %s.",
				box_name(storm), box_name(where))
			return nil
		}

		var ret *exit_view

		for i = 0; i < len(l); i++ {
			if l[i].destination == c.a {
				ret = l[i]
			}
		}

		if ret != nil {
			return ret
		}

		wout(c.who, "No route from %s to %s.",
			box_name(where),
			c.parse[1])

		return nil
	}

	dir = lookup_sb(full_dir_s, (c.parse[1]))
	if dir < 0 {
		dir = lookup_sb(short_dir_s, (c.parse[1]))
	}

	if dir < 0 {
		wout(c.who, "Unknown direction or destination '%s'.",
			c.parse[1])
		return nil
	}

	if !DIR_NSEW(dir) {
		wout(c.who, "Direction must be N, S, E or W.")
		return nil
	}

	for i = 0; i < len(l); i++ {
		if l[i].direction == dir {
			return l[i]
		}
	}

	wout(c.who, "No %s route from %s.",
		full_dir_s[dir], box_name(where))
	return nil
}

func v_direct_storm(c *command) int {
	storm := c.a
	var v *exit_view
	var dest int

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	cmd_shift(c)

	v = parse_storm_dir(c, storm)

	if v == nil {
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if loc_depth(v.destination) != LOC_province {
		wout(c.who, "Can't direct storm to %s.",
			box_code(v.destination))
		return FALSE
	}

	dest = v.destination
	p_misc(storm).storm_move = dest
	p_misc(storm).npc_dir = v.direction

	wout(c.who, "%s will move to %s at month end.",
		box_name(storm), box_name(dest))

	return TRUE
}

func v_dissipate(c *command) int {
	storm := c.a
	var where int
	var here_s string

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	where = province(cast_where(c.who))

	if where == province(subloc(c.who)) {
		here_s = ("here")
	} else {
		here_s = fmt.Sprintf("in %s", box_name(where))
	}

	c.d = where

	if subloc(storm) != where {
		wout(c.who, "%s is not %s.", box_name(storm), here_s)
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	return TRUE
}

func d_dissipate(c *command) int {
	storm := c.a
	where := c.d
	var here_s string
	var p *entity_misc

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	if where == province(subloc(c.who)) {
		here_s = ("here")
	} else {
		here_s = fmt.Sprintf("in %s", box_name(where))
	}

	if subloc(storm) != where {
		wout(c.who, "%s is not %s.", box_name(storm), here_s)
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	p = p_misc(storm)

	/* add_aura(c.who,p.storm_str/2); */
	add_piety(c.who, p.storm_str/2, true)
	p.storm_str = 0

	dissipate_storm(storm, true)
	out(c.who, "Current piety is now %s.", nice_num(rp_char(c.who).religion.piety))

	return TRUE
}

func v_renew_storm(c *command) int {
	storm := c.a
	aura := c.b
	var where int
	var here_s string

	if kind(storm) != T_storm {
		wout(c.who, "%s is not a storm.", box_code(storm))
		return FALSE
	}

	if aura < 1 {
		c.b = 1
		aura = 1
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if !has_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	where = province(cast_where(c.who))

	if where == province(subloc(c.who)) {
		here_s = ("here")
	} else {
		here_s = fmt.Sprintf("in %s", box_name(where))
	}

	c.d = where

	if subloc(storm) != where {
		wout(c.who, "%s is not %s.", box_name(storm), here_s)
		return FALSE
	}

	return TRUE
}

func d_renew_storm(c *command) int {
	storm := c.a
	aura := c.b
	where := c.d
	var here_s string
	var p *entity_misc

	if kind(storm) != T_storm {
		wout(c.who, "%s is not a storm.", box_code(storm))
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if where == province(subloc(c.who)) {
		here_s = ("here")
	} else {
		here_s = fmt.Sprintf("in %s", box_name(where))
	}

	if subloc(storm) != where {
		wout(c.who, "%s is not %s.", box_name(storm), here_s)
		return FALSE
	}

	if !use_piety(c.who, aura) {
		wout(c.who, "You haven't the aura for that prayer.")
		return FALSE
	}

	p = p_misc(storm)
	p.storm_str += aura * 2

	out(c.who, "%s is now strength %s.",
		box_name(storm), comma_num(p.storm_str))

	return TRUE
}

func v_lightning(c *command) int {
	storm := c.a
	target := c.b
	var where int

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	if subkind(storm) != sub_rain {
		wout(c.who, "%s is not a rain storm.", box_name(storm))
		return FALSE
	}

	where = subloc(storm)

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if kind(target) != T_char && !is_loc_or_ship(target) {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if is_loc_or_ship(target) && loc_depth(target) != LOC_build {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if subloc(target) != where {
		wout(c.who, "Target %s isn't in the same place as the storm.",
			box_code(target))
		return FALSE
	}

	if in_safe_now(target) != FALSE {
		wout(c.who, "Not allowed in a safe haven.")
		return FALSE
	}

	return TRUE
}

func d_lightning(c *command) int {
	storm := c.a
	target := c.b
	aura := c.c
	var where int
	var p *entity_misc

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	if subkind(storm) != sub_rain {
		wout(c.who, "%s is not a rain storm.", box_name(storm))
		return FALSE
	}

	where = subloc(storm)

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if kind(target) != T_char && !is_loc_or_ship(target) {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if is_loc_or_ship(target) && loc_depth(target) != LOC_build {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if subloc(target) != where {
		wout(c.who, "Target %s isn't in the same place as the storm.",
			box_code(target))
		return FALSE
	}

	if in_safe_now(target) != FALSE {
		wout(c.who, "Not allowed in a safe haven.")
		return FALSE
	}

	p = p_misc(storm)

	if aura == 0 {
		aura = p.storm_str
	}

	if aura > p.storm_str {
		aura = p.storm_str
	}

	p.storm_str -= aura

	wout(c.who, "%s strikes %s with a lightning bolt!",
		box_name(storm), box_name(target))

	vector_clear()
	vector_add(where)
	vector_add(target)

	if is_loc_or_ship(target) {
		/*
		 *  Protection provided by priests of Eres
		 *
		 */
		if priest_in_stack(target, sk_eres) {
			wout(VECT, "Lightning narrowly misses %s!", box_name(target))
		} else {
			wout(VECT, "%s was struck by lightning!",
				box_name(target))
			add_structure_damage(target, aura)
		}
	} else {
		if is_priest(target) == sk_eres {
			wout(VECT, "Lightning strikes %s but does no harm!",
				box_name(target))
		} else {
			wout(VECT, "%s was struck by lightning!",
				box_name(target))
			add_char_damage(target, aura, MATES)
		}
	}

	if p.storm_str <= 0 {
		dissipate_storm(storm, true)
	}

	return TRUE
}

//#if 0
//int
//v_list_storms(c *command)
//{
//
//  if (!has_piety(c.who, 1)) {
//    wout(c.who, "You haven't the aura for that prayer.");
//    return FALSE;
//  };
//
//    return TRUE;
//}
//
//
//int
//d_list_storms(c *command)
//{
//    var where int
//    var i int
//    first := TRUE;
//    var here_s string;
//
//    if (!use_piety(c.who, 1)) {
//      wout(c.who, "You haven't the aura for that prayer.");
//      return FALSE;
//    };
//
//    where = province(reset_cast_where(c.who));
//
//    if (crosses_ocean(where,c.who)) {
//      wout(c.who,"Something seems to block your magic.");
//      return FALSE;
//    };
//
//    if (where == province(subloc(c.who)))
//        here_s = ( "here");
//    else
//        here_s = fmt.Sprintf( "in %s", box_name(where));
//
//    for _, i = range loop_here(where, i)
//    {
//        if (kind(i) != T_storm)
//            continue;
//
//        if (first)
//        {
//            wout(c.who, "Storms %s:", here_s);
//            indent += 3;
//            first = FALSE;
//        }
//
//        wout(c.who, "%s", liner_desc(i));
//    }
//
//
//    if (first)
//        wout(c.who, "There are no storms %s.", here_s);
//    else
//        indent -= 3;
//
//    return TRUE;
//}
//#endif

func v_seize_storm(c *command) int {
	storm := c.a
	var where int
	var here_s string

	if kind(storm) != T_storm {
		wout(c.who, "%s isn't a storm.", box_code(storm))
		return FALSE
	}

	if npc_summoner(storm) == c.who {
		wout(c.who, "You already control %s.", box_name(storm))
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	where = province(cast_where(c.who))

	if where == province(subloc(c.who)) {
		here_s = ("here")
	} else {
		here_s = fmt.Sprintf("in %s", box_name(where))
	}

	c.d = where

	if subloc(storm) != where {
		wout(c.who, "%s is not %s.", box_name(storm), here_s)
		return FALSE
	}

	return TRUE
}

func d_seize_storm(c *command) int {
	storm := c.a
	where := c.d
	var here_s string
	var owner int

	if kind(storm) != T_storm {
		wout(c.who, "%s isn't a storm.", box_code(storm))
		return FALSE
	}

	owner = npc_summoner(storm)

	if owner != FALSE && owner == c.who {
		wout(c.who, "You already control %s.", box_name(storm))
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	if where == province(subloc(c.who)) {
		here_s = ("here")
	} else {
		here_s = fmt.Sprintf("in %s", box_name(where))
	}

	if subloc(storm) != where {
		wout(c.who, "%s is not %s.", box_name(storm), here_s)
		return FALSE
	}

	vector_clear()
	vector_add(c.who)
	if owner != FALSE {
		vector_add(owner)
	}

	wout(VECT, "%s seized control of %s!",
		box_name(c.who), box_name(storm))

	p_misc(storm).summoned_by = c.who

	return TRUE
}

func v_death_fog(c *command) int {
	storm := c.a
	target := c.b
	var where int

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	if subkind(storm) != sub_fog {
		wout(c.who, "%s is not a fog.", box_name(storm))
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	where = subloc(storm)

	if kind(target) != T_char {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if in_safe_now(target) != FALSE {
		wout(c.who, "Not allowed in a safe haven.")
		return FALSE
	}

	if subloc(target) != where {
		wout(c.who, "Target %s isn't in the same place as the fog.",
			box_code(target))
		return FALSE
	}

	return TRUE
}

func fog_excuse() string {

	switch rnd(1, 3) {
	case 1:
		return "wandered off in the fog and were lost."
	case 2:
		return "choked to death in the poisonous fog."
	case 3:
		return "disappeared in the fog."
	}
	panic("!reached")
}

func d_death_fog(c *command) int {
	storm := c.a
	target := c.b
	aura := c.c
	var where int
	var p *entity_misc
	//var save_aura int
	var aura_used int
	var e *item_ent

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	if subkind(storm) != sub_fog {
		wout(c.who, "%s is not a fog.", box_name(storm))
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	where = subloc(storm)

	if kind(target) != T_char {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if subloc(target) != where {
		wout(c.who, "Target %s isn't in the same place as the fog.",
			box_code(target))
		return FALSE
	}

	p = p_misc(storm)

	aura = min(aura, p.storm_str)

	for _, e = range loop_inventory(target) {
		/*
		 *   Quit if we've run out of aura.
		 *
		 */
		if aura < 1 {
			break
		}

		if item_attack(e.item) > 0 &&
			item_attack(e.item)+item_defense(e.item)+
				item_missile(e.item) < 25 {
			/*
			 *  Kill some or all of these...
			 *
			 */
			amount := min(has_item(target, e.item), aura)
			consume_item(target, e.item, amount)
			wout(target, "%s %s %s.", cap_(nice_num(amount)),
				plural_item_name(e.item, amount),
				fog_excuse())
			aura_used += amount
		}
	}

	if aura_used == 0 {
		wout(c.who, "%s has no vulnerable men.",
			box_name(target))
		return FALSE
	}

	wout(c.who, "Killed %s men of %s.",
		nice_num(aura_used),
		or_string(aura_used == 1, "man", "men"),
		box_name(target))

	p.storm_str -= aura_used

	if p.storm_str <= 0 {
		dissipate_storm(storm, true)
	}

	return TRUE
}

func v_banish_corpses(c *command) int {
	target := c.a
	sum := 0

	sum = has_item(target, item_corpse)

	if sum == 0 {
		wout(c.who, "There are no %s here.",
			plural_item_name(item_corpse, 2))
		return FALSE
	}

	return TRUE
}

func d_banish_corpses(c *command) int {
	target := c.a
	max_aura := c.b
	var i int
	sum := 0
	var n int
	var max int

	sum = has_item(target, item_corpse)

	if sum == 0 {
		wout(c.who, "There are no %s here.",
			plural_item_name(item_corpse, 2))
		return FALSE
	}

	if is_priest(c.who) != FALSE {
		max = rp_char(c.who).religion.piety
	} else {
		max = char_cur_aura(c.who)
	}

	if max == 0 {
		wout(c.who, "You lack the strength to challenge the undead.")
		return FALSE
	}

	if max < sum {
		sum = max
	}

	if max_aura != FALSE && (sum > max_aura) {
		sum = max_aura
	}

	if is_priest(c.who) != FALSE {
		use_piety(c.who, sum)
	} else {
		charge_aura(c.who, sum)
	}

	wout(c.who, "Banished %s %s.", comma_num(sum),
		plural_item_name(item_corpse, sum))
	wout(province(c.who), "%s banished %s %s!",
		box_name(c.who), comma_num(sum),
		plural_item_name(item_corpse, sum))

	consume_item(target, item_corpse, sum)
	// todo: i and n are uninitialized here
	wout(i, "%s banished our %s!", box_name(c.who),
		plural_item_name(item_corpse, n))

	return TRUE
}

func v_fierce_wind(c *command) int {
	storm := c.a
	target := c.b
	var where int

	if !has_holy_symbol(c.who) {
		wout(c.who, "You must have a holy symbol to use fierce wind.")
		return FALSE
	}

	if !has_holy_plant(c.who) {
		wout(c.who, "You must have a holy plant to use fierce wind.")
		return FALSE
	}

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	if subkind(storm) != sub_wind {
		wout(c.who, "%s is not a rain storm.", box_name(storm))
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	where = subloc(storm)

	if !is_loc_or_ship(target) || loc_depth(target) != LOC_build {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if subloc(target) != where {
		wout(c.who, "Target %s isn't in the same place as the storm.",
			box_code(target))
		return FALSE
	}

	return TRUE
}

func d_fierce_wind(c *command) int {
	storm := c.a
	target := c.b
	aura := c.c
	var where int
	var p *entity_misc

	if !has_holy_symbol(c.who) {
		wout(c.who, "You must have a holy symbol to use fierce wind.")
		return FALSE
	}

	if !has_holy_plant(c.who) {
		wout(c.who, "You must have a holy plant to use fierce wind.")
		return FALSE
	}

	if kind(storm) != T_storm || npc_summoner(storm) != c.who {
		wout(c.who, "You don't control any storm %s.",
			box_code(storm))
		return FALSE
	}

	if subkind(storm) != sub_rain {
		wout(c.who, "%s is not a rain storm.", box_name(storm))
		return FALSE
	}

	if crosses_ocean(storm, c.who) {
		wout(c.who, "Something seems to block your prayer.")
		return FALSE
	}

	where = subloc(storm)

	if !is_loc_or_ship(target) || loc_depth(target) != LOC_build {
		wout(c.who, "%s is not a valid target.", box_code(target))
		return FALSE
	}

	if subloc(target) != where {
		wout(c.who, "Target %s isn't in the same place as the storm.",
			box_code(target))
		return FALSE
	}

	if !use_piety(c.who, skill_piety(c.use_skill)) {
		wout(c.who, "You don't have the piety required to use that prayer.")
		return FALSE
	}

	/*
	 *  Use up 1 holy plant.
	 *
	 */
	move_item(c.who, 0, holy_plant(c.who), 1)

	p = p_misc(storm)

	if aura == 0 {
		aura = p.storm_str
	}

	if p.storm_str > aura {
		aura = p.storm_str
	}

	p.storm_str -= aura

	vector_clear()
	vector_add(where)
	vector_add(target)
	vector_add(c.who)
	/*
	 *  Protection provided by priests of Eres
	 *
	 */
	if priest_in_stack(target, sk_eres) {
		wout(VECT, "A fierce wind blows around %s but does no harm!",
			box_name(target))
	} else {
		wout(VECT, "%s is buffeted by a fierce wind!", box_name(target))
		add_structure_damage(target, aura)
	}

	if p.storm_str <= 0 {
		dissipate_storm(storm, true)
	}

	return TRUE
}

func create_some_storms(num, kind int) {
	var l []int
	var i int

	l = nil

	for _, i = range loop_province() {
		if in_clouds(i) || in_hades(i) || in_faery(i) {
			continue
		}

		if subkind(i) == sub_mine_shaft {
			continue
		}

		if weather_here(i, kind) != FALSE {
			continue
		}

		l = append(l, i)
	}

	l = shuffle_ints(l)

	for i = 0; i < len(l) && i < num; i++ {
		new_storm(0, kind, rnd(2, 3), l[i])
	}

	l = nil
}

func natural_weather() {
	nprov := nprovinces()
	var n int

	/*
	 *  One natural storm per 4 (formerly 16) provinces.
	 *  Half of storms made each month.
	 *  Called four times per month.
	 */

	n = nprov / 4 / 2 / 4

	switch oly_month(&sysclock) {
	case 0: /* Fierce winds */
		create_some_storms(n, sub_fog)
		create_some_storms(n, sub_wind)
		break

	case 1: /* Snowmelt */
		create_some_storms(n, sub_fog)
		create_some_storms(n, sub_rain)
		break

	case 2: /* Blossom bloom */
		break

	case 3: /* Sunsear */
		create_some_storms(n, sub_rain)
		break

	case 4: /* Thunder and rain */
		create_some_storms(n*2, sub_rain)
		break

	case 5: /* Harvest */
		break

	case 6: /* Waning days */
		create_some_storms(n, sub_rain)
		create_some_storms(n, sub_fog)
		create_some_storms(n, sub_rain)
		break

	case 7: /* Dark night */
		create_some_storms(n, sub_wind)
		break

	default:
		panic("!reached")
	}

}

func update_weather_view_loc_sup(who, where int) {
	var pl int

	pl = player(who)
	assert(valid_box(pl))

	p_player(pl).weatherSeen = set_bit(p_player(pl).weatherSeen, where)
}

func update_weather_view_locs(stack, where int) {
	var i int

	where = province(where)

	if kind(stack) == T_char && weather_mage(stack) != FALSE {
		update_weather_view_loc_sup(stack, where)
	}

	for _, i = range loop_char_here(stack) {
		if !is_prisoner(i) && weather_mage(i) != FALSE {
			update_weather_view_loc_sup(i, where)
		}
	}

}

func init_weather_views() {
	var who int

	for _, who = range loop_char() {
		if weather_mage(who) != FALSE {
			update_weather_view_loc_sup(who, province(who))
		}
	}

}

func can_see_weather_here(who, where int) bool {
	pl := player(who)

	assert(valid_box(pl))

	where = province(where)

	p := rp_player(pl)
	if p == nil {
		return false
	}

	return test_bit(p.weatherSeen, where)
}
