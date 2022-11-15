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

func scry_show_where(who, target int) {

	if crosses_ocean(who, target) {
		out(who, "%s is in %s.", just_name(region(target)))
		return
	}

	out(who, "%s is in %s.", box_name(province(target)))
}

func cast_where(who int) int {
	var where int

	where = char_proj_cast(who)

	if is_loc_or_ship(where) {
		return where
	}

	return subloc(who)
}

func reset_cast_where(who int) int {
	var where int

	where = char_proj_cast(who)

	if is_loc_or_ship(where) {
		p_magic(who).project_cast = 0
		return where
	}

	return subloc(who)
}

func cast_check_char_here(who, target int) int {
	var where int
	var basic int
	var pl int

	basic = char_proj_cast(who)

	if is_loc_or_ship(basic) {
		where = basic
	} else {
		where = subloc(who)
	}

	p_magic(who).project_cast = 0

	if crosses_ocean(where, who) {
		wout(who, "Something seems to block your magic.")
		return FALSE
	}

	if kind(target) != T_char || where != subloc(target) {
		wout(who, "%s is not a character in range of this cast.",
			box_code(target))
		return FALSE
	}

	if char_really_hidden(target) {
		pl = player(who)
		if pl == player(target) {
			return TRUE
		}

		if contacted(target, who) {
			return TRUE
		}

		return FALSE
	}

	//#if 0
	//    if (basic == where)
	//        p_magic(who).project_cast = 0;
	//#endif

	return TRUE
}

func v_scry_region(c *command) int {
	targ_loc := c.a
	aura := c.b

	if !is_loc_or_ship(targ_loc) {
		wout(c.who, "%s is not a location.", box_code(targ_loc))
		return FALSE
	}

	if crosses_ocean(targ_loc, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	if c.b < 1 {
		c.b = 1
	}
	aura = c.b

	if !check_aura(c.who, aura) {
		return FALSE
	}

	return TRUE
}

func alert_scry_attempt(who, where int, t string) {
	var has_detect int
	var source string

	for _, n := range loop_char_here(where) {
		has_detect = has_skill(n, sk_detect_scry)

		if has_detect > exp_novice {
			source = box_name(who)
		} else {
			source = "Someone"
		}

		if has_detect != FALSE {
			wout(n, "%s%s cast %s on this location.",
				source, t, box_name(sk_scry_region))
		}

		if has_detect >= exp_master {
			wout(n, "%s is in %s.", box_name(who),
				char_rep_location(who))
		}
		/*
		 *  If another exp gradient is wanted, use box_name for the loc,
		 *  then graduate to char_rep_location, since the latter gives
		 *  more info
		 */

	}

}

func alert_palantir_scry(who, where int) {
	var has_detect int

	for _, n := range loop_char_here(where) {
		has_detect = has_skill(n, sk_detect_scry)

		if has_detect < exp_master {
			continue
		}

		wout(n, "%s used a palantir to scry this location.",
			box_name(who))

		if has_detect > exp_master {
			wout(n, "%s is in %s.", box_name(who),
				char_rep_location(who))
		}
	}

}

func alert_scry_generic(who, where int) {
	var has_detect int

	for _, n := range loop_char_here(where) {
		has_detect = has_skill(n, sk_detect_scry)

		if has_detect != FALSE {
			wout(n, "%s scried %s from %s.",
				box_name(who),
				box_name(where),
				char_rep_location(who))
		}
	}

}

func check_shield_artifact(where int) int {

	for _, n := range loop_all_here(where) {
		if has_artifact(n, ART_SHIELD_PROV, 0, 0, 0) != FALSE {
			return TRUE
		}
	}

	if has_artifact(where, ART_SHIELD_PROV, 0, 0, 0) != FALSE {
		return TRUE
	}

	return FALSE
}

func d_scry_region(c *command) int {
	targ_loc := c.a
	aura := c.b

	if !is_loc_or_ship(targ_loc) {
		wout(c.who, "%s is no longer a valid location.",
			box_code(targ_loc))
		return FALSE
	}

	if crosses_ocean(targ_loc, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	if !charge_aura(c.who, aura) {
		return FALSE
	}

	if aura <= loc_shroud(province(targ_loc)) {
		wout(c.who, "%s is shrouded from your scry.",
			box_code(targ_loc))

		alert_scry_attempt(c.who, targ_loc, " unsuccessfully")

		return FALSE
	}

	/*
	 *  Sun Oct 11 18:26:34 1998 -- Scott Turner
	 *
	 *  Might be an artifact somewhere in the province.
	 *
	 */
	if check_shield_artifact(targ_loc) != FALSE {
		wout(c.who, "A grey mist blocks your vision.")
		alert_scry_attempt(c.who, targ_loc, " unsuccessfully")
		return FALSE
	}

	wout(c.who, "A vision of %s appears:", box_name(targ_loc))
	out(c.who, "")
	show_loc(c.who, targ_loc)

	alert_scry_attempt(c.who, targ_loc, "")

	return TRUE
}

func v_shroud_region(c *command) int {
	var aura int
	where := province(cast_where(c.who))

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	if c.a < 1 {
		c.a = 1
	}
	aura = c.a

	if !has_piety(c.who, aura) {
		wout(c.who, "You do not have that much piety.")
		return FALSE
	}

	wout(c.who, "Attempt to create a magical shroud to conceal %s from scry attempts.", box_code(where))

	reset_cast_where(c.who)
	c.b = where

	return TRUE
}

func d_shroud_region(c *command) int {
	aura := c.a
	var p *entity_loc
	where := c.b

	if !use_piety(c.who, aura) {
		wout(c.who, "You do not have that much piety.")
		return FALSE
	}

	p = p_loc(where)
	p.shroud += aura * 2

	wout(c.who, "%s is now cloaked with a strength %s location shroud.",
		box_name(where), nice_num(p.shroud))

	for _, n := range loop_char_here(where) {
		//var exp int

		if n == c.who {
			continue
		}

		if has_skill(n, sk_shroud_region) != FALSE {
			wout(n, "%s cast %s here.  %s is now cloaked with a strength %s location shroud.",
				box_name(c.who),
				box_name(sk_shroud_region),
				box_code(where),
				nice_num(p.shroud))
		}
	}

	return TRUE
}

func v_detect_scry(c *command) int {

	if !check_aura(c.who, 1) {
		return FALSE
	}

	wout(c.who, "Will practice location scry detection.")
	return TRUE
}

func d_detect_scry(c *command) int {

	if !charge_aura(c.who, 1) {
		return FALSE
	}

	return TRUE
}

func notify_loc_shroud(where int) {
	var p *entity_loc

	p = rp_loc(where)

	if p == nil {
		return
	}

	for _, who := range loop_char_here(where) {
		//var exp int

		if has_skill(who, sk_shroud_region) != FALSE {
			if p.shroud > 0 {
				wout(who, "The magical shroud over %s has diminished to %s aura.",
					box_name(where),
					nice_num(p.shroud))
			} else {
				wout(who, "The magical shroud over %s has dissipated.", box_name(where))
			}
		}
	}

}

func v_dispel_region(c *command) int {
	targ_loc := province(c.a)

	if !is_loc_or_ship(targ_loc) {
		wout(c.who, "%s is not a location.", box_code(targ_loc))
		return FALSE
	}

	if !check_aura(c.who, 3) {
		return FALSE
	}

	wout(c.who, "Attempt to dispel any magical shroud over %s.",
		box_name(targ_loc))

	return TRUE
}

func d_dispel_region(c *command) int {
	targ_loc := province(c.a)
	var p *entity_loc

	if !is_loc_or_ship(targ_loc) {
		wout(c.who, "%s is no longer a valid location.",
			box_code(targ_loc))
		return FALSE
	}

	p = rp_loc(targ_loc)

	if p != nil && p.shroud > 0 {
		if !charge_aura(c.who, 3) {
			return FALSE
		}

		wout(c.who, "Removed an aura %s magical shroud from %s.",
			nice_num(p.shroud),
			box_name(targ_loc))
		p.shroud = 0
		notify_loc_shroud(targ_loc)
	} else {
		wout(c.who, "%s was not magically shrouded.",
			box_name(targ_loc))
	}

	return TRUE
}

func show_item_where(who, target int) {
	var owner int
	var prov int

	assert(kind(target) == T_item)

	owner = item_unique(target)
	assert(owner != 0)

	prov = province(owner)

	if prov == owner {
		wout(who, "%s is in %s.", box_name(target), box_name(prov))
		return
	}

	if subkind(owner) == sub_graveyard {
		wout(who, "%s is buried in %s, in %s.",
			box_name(target),
			box_name(owner),
			box_name(prov))
		return
	}

	wout(who, "%s is held by %s, in %s.",
		box_name(target), box_name(owner), box_name(prov))
}

func v_locate_char(c *command) int {
	target := c.a
	aura := c.b

	if c.b < 1 {
		c.b = 1
	}
	aura = c.b

	if !check_aura(c.who, aura) {
		return FALSE
	}

	wout(c.who, "Attempt to locate %s.", box_code(target))

	return TRUE
}

func d_locate_char(c *command) int {
	target := c.a
	aura := c.b
	var chance int

	if kind(target) != T_char && subkind(target) != sub_dead_body {
		wout(c.who, "%s is not a character.", box_code(target))
		charge_aura(c.who, 1)
		return FALSE
	}

	if crosses_ocean(target, c.who) {
		wout(c.who, "Something seems to block your magic.")
		charge_aura(c.who, 1)
		return FALSE
	}

	if !charge_aura(c.who, aura) {
		return FALSE
	}

	switch aura {
	case 1:
		chance = 50
		break

	case 2:
		chance = 75
		break

	case 3:
		chance = 90
		break

	default:
		assert(false)
	}

	if rnd(1, 100) > chance {
		wout(c.who, "Character location failed.")
		return FALSE
	}

	if subkind(target) == sub_dead_body {
		show_item_where(c.who, target)
	} else {
		wout(c.who, "%s is in %s.", box_name(target),
			char_rep_location(target))
	}

	return TRUE
}

func v_bar_loc(c *command) int {
	var aura int
	var where int

	where = cast_where(c.who)

	if kind(where) != T_loc {
		wout(c.who, "%s is not a location.", box_code(where))
		return FALSE
	}

	if in_safe_now(where) != FALSE {
		wout(c.who, "Can't put a barrier around a safe haven.")
		return FALSE
	}

	if loc_depth(where) < LOC_subloc {
		wout(c.who, "Can't put a barrier around %s.", box_code(where))
		return FALSE
	}

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	if c.a < 1 {
		c.a = 1
	}
	if c.a > 4 {
		c.a = 4
	}
	aura = c.a * c.a

	if !check_aura(c.who, aura) {
		return FALSE
	}

	if loc_barrier(where) > 0 {
		wout(c.who, "%s already has a barrier.",
			box_name(where))
		return FALSE
	}

	if no_barrier(where) > 0 {
		wout(c.who, "You cannot cast another barrier on %s.",
			box_name(where))
		return FALSE
	}

	c.d = where
	reset_cast_where(c.who)

	wout(c.who, "Create a magical barrier over %s.", box_name(where))
	return TRUE
}

func d_bar_loc(c *command) int {
	aura := c.a * c.a
	where := c.d
	//var p *entity_loc
	//var old_val int

	if kind(where) != T_loc {
		wout(c.who, "%s is not a location.", box_code(where))
		return FALSE
	}

	if in_safe_now(where) != FALSE {
		wout(c.who, "Can't put a barrier around a safe haven.")
		return FALSE
	}

	if loc_depth(where) < LOC_subloc {
		wout(c.who, "Can't put a barrier around %s.", box_code(where))
		return FALSE
	}

	if loc_barrier(where) > 0 {
		wout(c.who, "%s already has a barrier.",
			box_name(where))
		return FALSE
	}

	if no_barrier(where) > 0 {
		wout(c.who, "You cannot cast another barrier on %s.",
			box_name(where))
		return FALSE
	}

	if !charge_aura(c.who, aura) {
		return FALSE
	}

	add_effect(where, ef_magic_barrier, 0, c.a*30, 1)

	wout(c.who, "Cast a %s aura barrier over %s.",
		nice_num(c.a),
		box_name(where))

	return TRUE
}

func v_unbar_loc(c *command) int {
	var aura int
	var where int
	var v *exit_view

	if c.parse[1][0] == '0' {
		where = cast_where(c.who)
	} else {
		v = parse_exit_dir(c, cast_where(c.who),
			sout("use %d", sk_unbar_loc))

		if v == nil {
			return FALSE
		}

		where = v.destination
	}

	if crosses_ocean(where, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	if c.b < 1 {
		c.b = 1
	}
	if c.b > 4 {
		c.b = 4
	}
	aura = c.b

	if !check_aura(c.who, aura) {
		return FALSE
	}

	if loc_barrier(where) == 0 {
		wout(c.who, "There is no barrier over %s.", box_name(where))
		return FALSE
	}

	c.d = where
	reset_cast_where(c.who)

	return TRUE
}

func d_unbar_loc(c *command) int {
	aura := c.b
	where := c.d
	//var p *entity_loc
	//var old_val int
	var chance int

	if kind(where) != T_loc {
		wout(c.who, "%s is not a location.", box_code(where))
		return FALSE
	}

	if loc_barrier(where) == 0 {
		wout(c.who, "There is no barrier over %s.", box_name(where))
		return FALSE
	}

	if !charge_aura(c.who, aura) {
		return FALSE
	}

	switch aura {
	case 1:
		chance = 10
		break

	case 2:
		chance = 25
		break

	case 3:
		chance = 50
		break

	case 4:
		chance = 75
		break

	default:
		assert(false)
	}

	if rnd(1, 100) > chance {
		wout(c.who, "Attempt to remove barrier fails.")
		return FALSE
	}

	wout(c.who, "The barrier over %s has been removed.", box_name(where))
	wout(where, "The barrier over %s has dissipated.", box_name(where))
	delete_effect(where, ef_magic_barrier, 0)
	add_effect(where, ef_inhibit_barrier, 0, 4, 1)

	return TRUE
}

func v_proj_cast(c *command) int {
	var to_where, curr_where int
	var aura int

	if c.a == 0 {
		c.a = subloc(c.who)
	}
	to_where = c.a

	if !is_loc_or_ship(to_where) {
		wout(c.who, "%s is not a location.", box_code(to_where))
		return FALSE
	}

	if in_safe_now(to_where) != FALSE {
		wout(c.who, "Magic may not be projected to safe havens.")
		return FALSE
	}

	curr_where = cast_where(c.who)

	if crosses_ocean(curr_where, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	aura = los_province_distance(curr_where, to_where) + 1
	c.d = aura
	assert(aura >= 0 && aura < 100)

	if crosses_ocean(to_where, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	/*
	 *  Don't needlessly give away the exact distance with
	 *  check_aura.
	 */

	if char_cur_aura(c.who) < aura {
		wout(c.who, "Not enough current aura.")
		return FALSE
	}

	wout(c.who, "Attempt to project next cast to %s.",
		box_name(to_where))

	reset_cast_where(c.who)

	return TRUE
}

func d_proj_cast(c *command) int {
	to_where := c.a
	aura := c.d
	var p *char_magic

	if !is_loc_or_ship(to_where) {
		wout(c.who, "%s is not a valid location.",
			box_code(to_where))
		return FALSE
	}

	if crosses_ocean(to_where, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	if !charge_aura(c.who, aura) {
		return FALSE
	}

	if subloc(c.who) != to_where && loc_shroud(province(to_where)) != 0 {
		wout(c.who, "%s is protected with a magical shroud.",
			box_name(to_where))
		wout(c.who, "Spell fails.")
		return FALSE
	}

	p = p_magic(c.who)
	p.project_cast = to_where

	wout(c.who, "Next cast will be based from %s.", box_name(to_where))

	return TRUE
}

func v_save_proj(c *command) int {

	if !valid_box(char_proj_cast(c.who)) {
		wout(c.who, "No projected cast state is active.")
		return FALSE
	}

	if !check_aura(c.who, 3) {
		return FALSE
	}

	wout(c.who, "Attempt to save projected cast state.")
	return TRUE
}

func d_save_proj(c *command) int {
	var p *char_magic
	var im *item_magic

	if !charge_aura(c.who, 3) {
		return FALSE
	}

	newPotion := new_potion(c.who)

	p = p_magic(c.who)
	im = p_item_magic(newPotion)

	im.use_key = use_proj_cast
	im.project_cast = p.project_cast

	p.project_cast = 0

	return TRUE
}

func v_use_proj_cast(c *command) int {
	item := c.a
	var im *item_magic

	assert(kind(item) == T_item)

	wout(c.who, "%s drinks the potion...", just_name(c.who))

	im = rp_item_magic(item)

	if im == nil ||
		is_loc_or_ship(im.project_cast) ||
		!is_magician(c.who) {
		destroy_unique_item(c.who, item)
		wout(c.who, "Nothing happens.")
		return FALSE
	}

	p_magic(c.who).project_cast = im.project_cast

	wout(c.who, "Project next cast to %s.",
		char_rep_location(im.project_cast))
	destroy_unique_item(c.who, item)

	return TRUE
}
