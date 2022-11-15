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
	MAX_SAVAGES = 200 /* max # of savage units allowed */
)

var (
	num_savages = 0 /* total # of savages in the world */
)

func auto_savage(who int) {
	where := subloc(who)

	/*
	 *  If stacked under someone, do nothing.
	 */
	if stack_parent(who) != FALSE {
		return
	}

	/*
	 *  If in a structure, issue RAZE
	 */
	if savage_hates(where) && building_owner(where) == who {
		queue(who, "raze")
		return
	}

	/*
	 *  If there is an inn/castle/tower/ship/temple here,
	 *	if empty,
	 *		enter,
	 *		raze.
	 *	if human-occupied,
	 *		beat drums to attract other savages,
	 *		wait,
	 *		attack.
	 */
	if target := savage_hates_here(where); target != 0 {
		if building_owner(target) == 0 {
			queue(who, "enter %s", box_code_less(target))
			queue(who, "raze")
			return
		} else if controlled_humans_here(target) != FALSE {
			if has_item(who, item_drum) < 1 {
				gen_item(who, item_drum, 1)
			}
			queue(who, "use %d 1", item_drum)
			queue(who, "wait time %d", rnd(35, 50))
			queue(who, "attack %s", box_code_less(target))
			return
		}
	}

	/*
	 *  Unstack any savages under us
	 */
	for _, i := range loop_here(who) {
		if kind(i) == T_char {
			queue(who, "unstack %d", i)
		}
	}

	if controlled_humans_here(where) == FALSE {
		queue(who, "die")
		return
	}

	npc_move(who)
}

func call_savage(where, to_where, who, why int) bool {
	if controlled_humans_here(where) != FALSE {
		return false
	}

	savage := create_savage(where)
	queue(savage, "move %s", box_code_less(to_where))
	switch why {
	case 0: /* battle challenge */
		queue(savage, "attack %s", box_code_less(who))
	case 1: /* call to arms */
		set_loyal(savage, LOY_summon, 3)
		queue(savage, "stack %s", box_code_less(who))
	case 2: /* move and attack structure */
		queue(savage, "use 98 1")
		queue(savage, "wait time %d", rnd(35, 50))
		queue(savage, "attack %s", box_code_less(who))
	}

	// init_load_sup(savage);   /* make ready to execute commands immediately */

	return true
}

func create_savage(where int) int {
	savage := create_monster_stack(item_savage, rnd(4, 26), where)
	if savage < 0 {
		return -1
	}
	gen_item(savage, item_drum, 1)
	return savage
}

func d_keep_savage(c *command) int {
	if !keep_savage_check(c) {
		return FALSE
	}

	target := c.a
	set_loyal(target, LOY_summon, max(loyal_rate(target)+2, 4))
	wout(c.who, "%s will remain for %d months.", box_code(target), loyal_rate(target))

	return TRUE
}

// 5% chance that a savage will be created to attack a structure in an un-garrisoned province each turn.
func init_savage_attacks() {
	num_savages = 0 // yep, this uses the global
	for _, i := range loop_kind(T_char) {
		if noble_item(i) == item_savage {
			num_savages++
		}
	}
	if num_savages >= MAX_SAVAGES {
		return
	}

	for _, fort := range loop_loc() {
		if loc_depth(fort) != LOC_build {
			continue
		} else if garrison_here(province(fort)) != FALSE {
			continue
		} else if rnd(1, 100) < 95 {
			continue
		}

		where := subloc(fort)
		if loc_depth(where) != LOC_province {
			continue
		}

		l := exits_from_loc_nsew_select(0, where, LAND, RAND)
		if len(l) == 0 {
			log.Printf("init_savage_attacks: no exits?\n")
			continue /* probably shouldn't happen */
		}

		for _, ex := range l {
			if call_savage(ex.destination, where, fort, 2) {
				break
			}
		}
	}
}

func keep_savage_check(c *command) bool {
	target := c.a
	if kind(target) != T_char || noble_item(target) != item_savage {
		wout(c.who, "%s is not a group of savages.", box_code(target))
		return false
	} else if subloc(target) != subloc(c.who) {
		wout(c.who, "%s is not here.", box_code(target))
		return false
	} else if loyal_kind(target) != LOY_summon {
		wout(c.who, "%s is no longer bonded.", box_code(target))
		return false
	}
	return true
}

func savage_hates(where int) bool {
	switch subkind(where) {
	case sub_castle, sub_castle_notdone,
		sub_galley, sub_galley_notdone,
		sub_inn, sub_inn_notdone,
		sub_roundship, sub_roundship_notdone,
		sub_temple, sub_temple_notdone,
		sub_tower, sub_tower_notdone:
		return true
	default:
		return false
	}
}

func savage_hates_here(where int) int {
	for _, i := range loop_here(where) {
		if is_loc_or_ship(i) && loc_depth(i) == LOC_build && savage_hates(i) {
			return i
		}
	}
	return 0
}

func v_keep_savage(c *command) int {
	if !keep_savage_check(c) {
		return FALSE
	}
	return TRUE
}

func v_summon_savage(c *command) int {
	if has_item(c.who, item_drum) < 1 {
		wout(c.who, "Must first make a drum with MAKE 98 1.")
		return FALSE
	}

	// todo: Scott Turner -- a limit on how many wild men you can summon, since they are orthogonal to "control men in battle"?
	num_savages := 0 // yeah, this obscures the global
	for _, i := range loop_stack(c.who) {
		if kind(i) == T_char && subkind(i) == sub_ni && noble_item(i) == item_savage {
			num_savages++
		}
	}

	if num_savages > 2 {
		wout(c.who, "You may only summon 3 savage stacks at a time.")
		return FALSE
	}

	c.a = 1 /* speed = summon */

	return v_use_drum(c)
}

func v_use_drum(c *command) int {
	where := subloc(c.who)
	speed := c.a
	var s, speed_s string

	switch speed {
	case 1:
		speed_s = "slow "
	case 2:
		speed_s = "fast "
	}

	wout(c.who, "%s sounds a %sdrumbeat.", box_name(c.who), speed_s)
	wout(where, "%s sounds a %sdrumbeat.", box_name(c.who), speed_s)

	if loc_depth(where) != LOC_province {
		s = sout("%sbeating drums may be heard coming from %s.", speed_s, box_name(where))
		wout(subloc(where), "%s", cap_(s))
		return TRUE
	}

	n := false
	for _, ex := range exits_from_loc_nsew_select(c.who, province(where), LAND, RAND) {
		dir := exit_opposite[ex.direction]
		s = sout("%sbeating drums may be heard to the %s.", speed_s, full_dir_s[dir])
		wout(ex.destination, "%s", cap_(s))
		if num_savages < MAX_SAVAGES && subkind(ex.destination) != sub_ocean && (speed == 0 || speed == 1) && !n {
			n = call_savage(ex.destination, where, c.who, speed)
		}
	}

	switch speed {
	case 0:
		s = "battle challenge"
	case 1:
		s = "call to arms"
	default:
		s = "call"
	}

	if !n && (speed == 0 || speed == 1) {
		wout(c.who, "No savages are responding to the %s.", s)
	} else if n && (speed == 0 || speed == 1) {
		wout(c.who, "Savages will surely respond to the %s.", s)
	}

	return TRUE
}
