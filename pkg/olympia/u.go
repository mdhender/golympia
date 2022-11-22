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
	"log"
	"strings"
	"time"
)

const (
	MATES        = (-1)
	MATES_SILENT = (-2)
	TAKE_ALL     = 1
	TAKE_SOME    = 2
	TAKE_NI      = 3 // noble item: wrapper adds one
)

var (
	dot_count = 0
	num_s     = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
)

type weights struct {
	total_weight int // total weight of unit or stack
	animals      int
	land_cap     int // carryable weight on land
	land_weight  int
	ride_cap     int // carryable weight on horseback
	ride_weight  int
	fly_cap      int // carryable weight on flying carpet?
	fly_weight   int
}

func ifHERO(who, sk int) bool {
	if !HERO {
		return false
	}
	return rnd(1, 100) > min((2*skill_exp(who, sk)), 80)
}

func ifnNEW_TRADE(who, item, old int) {
	if NEW_TRADE {
		return
	}
	investigate_possible_trade(who, item, old)
}

func investigate_possible_trade(who, item, old int) { panic("!implemented") }

/*
 *  u.c -- the useful function junkyard
 */

/*
 *  Mon Jun 16 11:45:05 1997 -- Scott Turner
 *
 *  True if you have to cross an ocean between the two
 *  locs; cheats right now by simply seeing if they're in
 *  the same region.  This works if you define the map correctly :-)
 *
 */
func crosses_ocean(a, b int) bool {
	return region(a) != region(b)
}

func kill_stack_ocean(who int) {
	var l []int
	for _, i := range loop_stack(who) {
		l = append(l, i)
	}

	for i := len(l) - 1; i >= 0; i-- {
		// todo: would be neat to have the ocean inherit, then wash stuff up on the beach
		kill_char(l[i], 0, S_body)
		if kind(l[i]) == T_char { /* not dead yet! */
			extract_stacked_unit(l[i])
			where := find_nearest_land(province(l[i]))
			out(l[i], "%s washed ashore at %s.", box_name(l[i]), box_name(where))
			log_output(LOG_SPECIAL, "kill_stack_ocean, swam ashore, who=%s", box_code_less(l[i]))
			move_stack(l[i], where)
		}
	}
}

func survive_fatal(who int) bool {
	if forget_skill(who, sk_survive_fatal) || forget_skill(who, sk_survive_fatal2) {
		wout(who, "%s would have died, but survived a fatal wound!", box_name(who))
		wout(who, "Forgot %s.", box_code(sk_survive_fatal))
		wout(who, "Health is now 100.")

		p_char(who).health = 100
		p_char(who).sick = FALSE

		if options.survive_np && skill_np_req(sk_survive_fatal) != FALSE {
			wout(who, "Received back %d noble points.", skill_np_req(sk_survive_fatal))
			add_np(player(who), skill_np_req(sk_survive_fatal))
		}

		return true
	}

	return false
}

func char_reclaim(who int) {
	p_char(who).melt_me = TRUE
	// QUIT shouldn't give items to stackmates
	// kill_char(who, MATES);
	kill_char(who, 0, S_body)
}

func v_reclaim(c *command) int {
	var what string
	if numargs(c) < 1 || len(c.parse[1]) == 0 {
		what = "disperses."
	} else {
		what = string(c.parse[1])
	}
	wout(subloc(c.who), "%s %s", box_name(c.who), what)
	char_reclaim(c.who)
	return TRUE
}

func new_char(sk, ni, where, health, pl, loy_kind, loy_lev int, name string) int {
	newt := new_ent(T_char, sk)
	if newt < 0 {
		return -1
	}
	if name != "" {
		set_name(newt, name)
	}
	p := p_char(newt)
	p.health = health
	p.unit_item = ni
	p.break_point = 50
	if ni != FALSE && (item_attack(ni) != FALSE || item_defense(ni) != FALSE) {
		p_char(newt).attack = item_attack(ni)
		rp_char(newt).defense = item_defense(ni)
	} else if sk == sub_garrison {
		p_char(newt).attack = 0
		rp_char(newt).defense = 0
	} else {
		p_char(newt).attack = 60
		rp_char(newt).defense = 60
	}

	// set NPC program?
	if ni != FALSE && item_prog(ni) != FALSE {
		rp_char(newt).npc_prog = int(item_prog(ni))
	}

	if is_loc_or_ship(where) {
		set_where(newt, where)
	} else {
		set_where(newt, subloc(where))
	}

	set_lord(newt, pl, loy_kind, loy_lev)

	if kind(where) == T_char {
		join_stack(newt, where)
	}

	if beast_capturable(newt) || is_npc(newt) {
		p.break_point = 0
	}

	return newt
}

func loc_depth(n int) int {
	switch subkind(n) {
	case sub_region:
		return LOC_region

	case sub_ocean,
		sub_forest,
		sub_plain,
		sub_mountain,
		sub_desert,
		sub_swamp,
		sub_under,
		sub_cloud,
		sub_mine_shaft:
		return LOC_province

	case sub_island,
		sub_stone_cir,
		sub_mallorn_grove,
		sub_bog,
		sub_cave,
		sub_city,
		sub_city_notdone,
		sub_lair,
		sub_graveyard,
		sub_ruins,
		sub_battlefield,
		sub_ench_forest,
		sub_rocky_hill,
		sub_tree_circle,
		sub_pits,
		sub_pasture,
		sub_oasis,
		sub_yew_grove,
		sub_sand_pit,
		sub_sacred_grove,
		sub_poppy_field,
		sub_faery_hill,
		sub_hades_pit,
		sub_mine,
		sub_mine_notdone,
		sub_mine_collapsed:
		return LOC_subloc

	case sub_guild,
		sub_temple,
		sub_galley,
		sub_roundship,
		sub_ship,
		sub_castle,
		sub_galley_notdone,
		sub_roundship_notdone,
		sub_ship_notdone,
		sub_ghost_ship,
		sub_temple_notdone,
		sub_inn,
		sub_inn_notdone,
		sub_castle_notdone,
		sub_tower,
		sub_tower_notdone,
		sub_mine_shaft_notdone,
		sub_orc_stronghold,
		sub_orc_stronghold_notdone:
		return LOC_build
	}

	panic(fmt.Sprintf("assert(subkind != %d)", subkind(n)))
}

// first try to give items to someone below, then to someone above.
func stackmate_inheritor(who int) int {
	for _, i := range loop_here(who) {
		if kind(i) == T_char && !is_prisoner(i) {
			if i == 0 {
				return stack_parent(who)
			}
			return i
		}
	}
	return stack_parent(who)
}

func take_unit_items(from, inherit, how_many int) {
	first := true

	var silent bool
	var to int
	switch inherit {
	case 0:
		to = 0
		silent = true
		break

	case MATES:
		to = stackmate_inheritor(from)
		silent = false
		break

	case MATES_SILENT:
		to = stackmate_inheritor(from)
		silent = true
		break

	default:
		to = inherit
		silent = false
	}

	if how_many == TAKE_NI {
		gen_item(from, noble_item(from), 1)
	}

	for _, e := range loop_inventory(from) {
		/*
		 *  Thu Mar 29 12:30:23 2001 -- Scott Turner
		 *
		 *  An auraculum should stay with its maker if he dies.
		 *
		 */
		if how_many == TAKE_SOME && e.item == char_auraculum(from) {
			continue
		}

		/*
		 *  Ungiveable items can't be taken this way.
		 *
		 *  Mon May  3 08:49:42 1999 -- Scott Turner
		 *
		 *  We can't just continue here; we need to get rid of the ungiveables.
		 *
		 */
		if rp_item(e.item) != nil && rp_item(e.item).ungiveable != FALSE {
			move_item(from, 0, e.item, e.qty)
			continue
		}

		/*
		 *  Wed Mar 12 11:56:30 1997 -- Scott Turner
		 *
		 *  No beast transfers this way?
		 *
		 *  Mon May  3 08:51:30 1999 -- Scott Turner
		 *
		 *  Turn these into monster stacks.
		 *
		 */
		if item_capturable(e.item) != FALSE && e.item != noble_item(from) {
			create_monster_stack(e.item, e.qty, subloc(from))
			move_item(from, 0, e.item, e.qty)
			continue
		}

		/*
		 *  Wed Mar 12 11:58:54 1997 -- Scott Turner
		 *
		 *  Only some trained men will come over.  att+def = 10
		 *  is the 50/50 point, so the chance to come over is:
		 *
		 *               5/(att+def)
		 *
		 */
		qty := 0
		if item_attack(e.item) != FALSE || item_defense(e.item) != FALSE {
			for i := 0; i < e.qty; i++ {
				if rnd(1, 100) < (500 / (item_attack(e.item) + item_defense(e.item))) {
					qty++
				}
			}
		} else {
			qty = e.qty
		}

		/*
		 *  Don't let beasts grab men this way.
		 *
		 */
		if man_item(e.item) != FALSE && subkind(to) == sub_ni {
			qty = 0
		}

		if how_many == TAKE_SOME && rnd(1, 2) == 1 {
			qty = rnd(qty/2, qty)
		}

		/*
		 *  Don't let unique items get dropped this way
		 */
		if qty == 0 && item_unique(e.item) != FALSE {
			qty = 1
		}

		if qty > 0 && !silent && valid_box(to) {
			if first {
				first = true
				wout(to, "Taken from %s:", box_name(from))
				indent += 3
			}

			wout(to, "%s", box_name_qty(e.item, qty))
		}

		move_item(from, to, e.item, qty)

		if e.item == item_gold && player(from) != player(to) && player(to) > 1000 {
			if player(from) < 1000 {
				gold_combat_indep += qty
			} else {
				gold_combat += qty
			}
		}

		if qty != e.qty {
			move_item(from, 0, e.item, e.qty-qty)
		}
	}

	/*
	 *  tranfer prisoners, too
	 */

	for _, i := range loop_here(from) {
		if !(kind(i) == T_char && is_prisoner(i)) {
			continue
		}
		if to > 0 {
			if first && !silent {
				wout(to, "Taken from %s:", box_name(from))
				indent += 3
				first = false
			}
			move_prisoner(from, to, i)
			if !silent {
				wout(to, "%s", liner_desc(i))
			}
			if player(i) == player(to) {
				p_char(i).prisoner = FALSE
			}
		} else {
			p_magic(i).swear_on_release = FALSE
			drop_stack(from, i)
		}
	}

	if !first && !silent {
		indent -= 3
	}
}

func add_char_damage(who, amount, inherit int) {
	if amount <= 0 {
		return
	}

	p := p_char(who)
	if p.health == -1 {
		if amount >= 50 {
			kill_char(who, inherit, S_body)
		}
		return
	}

	if p.health > 0 {
		if amount > p.health {
			amount = p.health
		}

		p.health -= amount
		assert(p.health >= 0)

		wout(who, "%s is wounded.  Health is now %d.",
			box_name(who), p.health)
	}

	if p.health <= 0 {
		kill_char(who, inherit, S_body)
	} else if FALSE == p.sick && rnd(1, 100) > p.health && ifHERO(who, sk_avoid_illness) && FALSE == has_artifact(who, ART_SICKNESS, 0, 0, 0) {
		/*
		 *  Wed Nov 25 12:51:14 1998 -- Scott Turner
		 *
		 *  Hero skill "Avoid Illness" can help negate this.
		 *
		 */

		p.sick = TRUE
		wout(who, "%s has fallen ill.", box_name(who))
	}
}

func put_back_cookie(who int) {
	p := rp_misc(who)
	if p == nil || p.npc_home == 0 {
		return
	}
	gen_item(p.npc_home, p.npc_cookie, 1)
}

/*
 *  Has a contacted b || has b found a?
 */
func contacted(a, b int) bool {
	p := p_char(a)

	if ilist_lookup(p.contact, b) >= 0 {
		return true
	}

	if ilist_lookup(p.contact, player(b)) >= 0 {
		return true
	}

	return false
}

func char_here(who, target int) bool {
	if where := subloc(who); where != subloc(target) {
		return false
	} else if char_really_hidden(target) {
		if pl := player(who); pl == player(target) {
			return true
		} else if contacted(target, who) {
			return true
		}
		return false
	}
	return true
}

func check_char_here(who, target int) bool {
	if target == garrison_magic {
		wout(who, "There is no garrison here.")
		return false
	} else if kind(target) != T_char || !char_here(who, target) {
		wout(who, "%s is not here.", box_code(target))
		return false
	}
	return true
}

func check_char_gone(who, target int) bool {
	if target == garrison_magic {
		wout(who, "There is no garrison here.")
		return false
	} else if kind(target) != T_char {
		wout(who, "%s is not a character.", box_code(target))
		return false
	} else if !char_here(who, target) {
		wout(who, "%s can not be seen here.", box_code(target))
		return false
	} else if char_gone(target) != FALSE {
		wout(who, "%s has left.", box_name(target))
		return false
	}
	return true
}

func check_still_here(who, target int) bool {
	if target == garrison_magic {
		wout(who, "There is no garrison here.")
		return false
	} else if kind(target) != T_char {
		wout(who, "%s is not a character.", box_code(target))
		return false
	} else if !char_here(who, target) {
		wout(who, "%s can no longer be seen here.", box_name(target))
		return false
	}
	return true
}

func check_skill(who, skill int) bool {
	if has_skill(who, skill) < 1 {
		wout(who, "Requires %s.", box_name(skill))
		return false
	}
	return true
}

func sink_ship(ship int) {
	log_output(LOG_SPECIAL, "%s has sunk in %s.", box_name(ship), box_name(subloc(ship)))

	wout(ship, "%s has sunk!", box_name(ship))
	wout(subloc(ship), "%s has sunk!", box_name(ship))

	where := subloc(ship)
	if subkind(where) == sub_ocean {
		for _, who := range loop_here(ship) {
			if kind(who) == T_char {
				kill_stack_ocean(who)
			}
		}
	} else {
		for _, who := range loop_here(ship) {
			if kind(who) == T_char {
				move_stack(who, where)
			} else {
				set_where(who, where)
			}
		}

	}

	//#if 0
	//    /*
	//     *  Unbind any storms bound to this ship
	//     */
	//
	//        p = rp_subloc(ship);
	//        if (p)
	//        {
	//            for i = 0; i < len(p.bound_storms); i++
	//            {
	//                storm = p.bound_storms[i];
	//                if (kind(storm) == T_storm)
	//                    p_misc(storm).bind_storm = 0;
	//            }
	//
	//            ilist_clear(&p.bound_storms);
	//        }
	//#endif

	set_where(ship, 0)
	delete_box(ship)
}

func get_rid_of_collapsed_mine(fort int) {
	assert(subkind(fort) == sub_mine_collapsed)

	/*
	 *  Move anything inside, out, just in case
	 */

	where := subloc(fort)
	for _, who := range loop_here(fort) {
		if kind(who) == T_char {
			move_stack(who, where)
		} else {
			set_where(who, where)
		}
	}

	set_where(fort, 0)
	delete_box(fort)
}

func building_collapses(fort int) {
	where := subloc(fort)
	log_output(LOG_SPECIAL, "%s collapsed in %s.", box_name(fort), box_name(where))

	vector_char_here(fort)
	vector_add(where)
	wout(VECT, "%s collapses!", box_name(fort))

	for _, who := range loop_here(fort) {
		if kind(who) == T_char {
			move_stack(who, where)
		} else {
			set_where(who, where)
		}
	}

	if subkind(fort) == sub_mine {
		change_box_subkind(fort, sub_mine_collapsed)
		p_misc(fort).mine_delay = 8
		return
	} else if subkind(fort) == sub_castle {
		for _, i := range loop_garrison() {
			if garrison_castle(i) == fort {
				p_misc(i).garr_castle = 0
			}
		}
	}

	set_where(fort, 0)
	delete_box(fort)
}

func add_structure_damage(fort, damage int) bool {
	assert(damage >= 0)

	p := p_subloc(fort)
	if p.damage+damage > 100 {
		p.damage = 100
	} else {
		p.damage += damage
	}
	if p.damage < 100 { // only partially destroyed
		return false
	}

	// completely destroyed, so sink or collapse it
	if is_ship(fort) {
		sink_ship(fort)
	} else {
		building_collapses(fort)
	}

	return true
}

/*
 *  Wed Dec  1 18:36:53 1999 -- Scott Turner
 *
 *  Counting function.
 *  Count the items in a unit/stack that meet the acceptance function.
 *
 */
func count_generic(who int, stack bool, fn func(int) int) int {
	sum := 0
	if stack {
		// todo: what does this do?
		// it looks like it loops for every char in the stack and then adds the same qty.
		// why not just sum += len(loop_char_here(who)) * count_generic(who, false, fn)?
		for _ = range loop_char_here(who) {
			sum += count_generic(who, false, fn)
		}
	} else {
		sum += fn(who)
		for _, e := range loop_inventory(who) {
			if fn(e.item) != 0 {
				sum += e.qty
			}
		}

	}
	return sum
}

func is_man_item(item int) bool {
	return man_item(item) != FALSE
}

func count_man_items(who int) int {
	// todo: why check here?
	// mdhender: commented out - seems to have no use
	//sum := 0;
	//if (subkind(who) == sub_garrison) {
	//    sum = 0;
	//}

	return count_generic(who, false, man_item)
}

// return the number of units in the stack, including who
func count_stack_units(who int) int {
	return len(loop_char_here(who)) + 1
}

func count_stack_figures(who int) int {
	sum := 0
	for _, i := range loop_stack(who) {
		sum += count_man_items(i)
	}
	return sum
}

func count_fighters_2(who, attack_min int) int {
	beasts, men := 0, 0
	for _, e := range loop_inventory(who) {
		if item_attack(e.item) >= attack_min {
			if item_animal(e.item) != FALSE {
				beasts += e.qty
			} else {
				men += e.qty
			}
		}
	}

	if beast_limit := calc_beast_limit(who, FALSE); beasts > beast_limit {
		beasts = beast_limit
	}
	if man_limit := calc_man_limit(who, FALSE); men > man_limit {
		men = man_limit
	}

	return beasts + men
}

func count_fighters(who, attack_min int) int {
	sum := 0
	for _, i := range loop_stack(who) {
		sum += count_fighters_2(i, attack_min)
	}
	return sum
}

/*
 *  Mon Oct 28 16:16:54 1996 -- Scott Turner
 *
 *  Ninjas don't "count" :-).  Neither do angels.
 *
 *  Mon Jan  4 08:16:19 1999 -- Scott Turner
 *
 *  And golems, which are treated as angels?
 *
 */
func count_any_real(who int, ignore_ninjas, ignore_angels bool) int {
	ignore_golems := ignore_angels
	sum := 1 // why?
	if subkind(who) == sub_garrison {
		sum = 0 // why?
	}
	for _, e := range loop_inventory(who) {
		if (ignore_ninjas && e.item == item_ninja) ||
			(ignore_angels && e.item == item_angel) ||
			(ignore_golems && (e.item == item_dirt_golem || e.item == item_flesh_golem || e.item == item_iron_golem)) {
			continue
		}

		if man_item(e.item) != FALSE || is_fighter(e.item) != FALSE {
			sum += e.qty
		}
	}

	return sum
}

func count_any(who, j, k int) int {
	return count_any_real(who, true, true)
}

func count_stack_any_real(who int, ignore_ninjas, ignore_angels bool) int {
	sum := 0
	for _, i := range loop_stack(who) {
		sum += count_any_real(i, ignore_ninjas, ignore_angels)
	}
	return sum
}

func count_stack_any(who int) int {
	return count_stack_any_real(who, true, true)
}

func count_loc_char_item(where, item int) int {
	sum := 0
	for _, i := range loop_char_here(where) {
		sum += has_item(i, item)
	}
	return sum
}

func clear_temps(kind int) {
	for _, i := range loop_kind(kind) {
		bx[i].temp = 0
	}
}

func olytime_increment(p *olytime) {
	p.days_since_epoch++
	p.day++
}

/*
 *  Ready counter for next turn
 *  Must be followed by an olytime_increment
 */
func olytime_turn_change(p *olytime) {
	p.day = 0
	p.turn++
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
 *  Olympian weight system
 *
 *	Each item has three fields related to weight and carrying
 *	capacity:
 *
 *		weight			fetch with item_weight(item)
 *		land capacity		fetch with land_cap(item)
 *		ride capacity		fetch with ride_cap(item)
 *
 *	Weight is the complete weight of the item, such as 100 for
 *	men, or 1,000 for oxen.
 *
 *	land capacity is how much the item can carry walking,
 *	not counting its own weight.
 *
 *	ride capacity is how much the item can carry on horseback,
 *	not counting its own weight.
 *
 *	if the item can carry itself riding or walking, but can not
 *	carry any extra goods, set the capacity to -1.  This is because
 *	0 represents "not set" instead of a value.
 *
 *	For example, a wild horse can walk and ride, but cannot be laden
 *	with rider or inventory.  Therefore, its land_cap is -1, and its
 *	ride_cap is -1.
 *
 *	An ox can carry great loads: land_cap 1500.  Perhaps it can trot
 *	alongside horses, but can carry no inventory if doing so.
 *	ride_cap -1.
 *
 */

/*
 *  Fri Oct  9 18:45:42 1998 -- Scott Turner
 *
 *  item_ride_cap and item_fly_cap can no longer be macros because
 *  we need to factor in artifacts.  They are also being changed to
 *  take T_char as well as T_item.
 *
 */
func item_ride_cap(who int) int {
	var a int
	base := who
	capacity := 0

	if kind(who) == T_char {
		base = noble_item(who)
	}
	if base == 0 {
		base = item_peasant
	}
	capacity = or_int(rp_item(base) != nil, rp_item(base).ride_cap, 0)

	if a = best_artifact(who, ART_RIDING, 0, 0); a != 0 {
		capacity += rp_item_artifact(a).param1
	}
	return capacity
}

func item_fly_cap(who int) int {
	var a int
	base := who
	capacity := 0

	if kind(who) == T_char {
		base = noble_item(who)
	}
	if base == 0 {
		base = item_peasant
	}
	capacity = or_int(rp_item(base) != nil, rp_item(base).fly_cap, 0)

	if a = best_artifact(who, ART_FLYING, 0, 0); a != 0 {
		capacity += rp_item_artifact(a).param1
	}
	return capacity
}

func item_land_cap(who int) int {
	var a int
	base := who
	capacity := 0

	if kind(who) == T_char {
		base = noble_item(who)
	}
	if base == 0 {
		base = item_peasant
	}
	capacity = or_int(rp_item(base) != nil, rp_item(base).land_cap, 0)

	if a = best_artifact(who, ART_CARRY, 0, 0); a != 0 {
		capacity += rp_item_artifact(a).param1
	}
	return capacity
}

func item_weight(who int) int {
	base := who
	capacity := 0

	if kind(who) == T_char {
		base = noble_item(who)
	}
	if base == 0 {
		base = item_peasant
	}
	capacity = or_int(rp_item(base) != nil, rp_item(base).weight, 0)

	/*
	 *  Sun Feb 16 22:18:31 1997 -- Scott Turner
	 *
	 *  Check for a potion of weightlessness.
	 *
	 */
	if get_effect(who, ef_weightlessness, 0, 0) != 0 {
		capacity -= 500
	}

	/*
	 *  And an artifact of weightlessness.
	 *
	 */
	if is_artifact(who) != nil && rp_item_artifact(who).type_ == ART_WEIGHTLESS {
		capacity -= rp_item_artifact(who).param1
	}

	return capacity
}

func add_item_weight(item, qty int, w *weights, mountains bool) {
	wt := item_weight(item) * qty
	lc := item_land_cap(item)
	rc := item_ride_cap(item)
	fc := item_fly_cap(item)

	if lc != 0 {
		w.land_cap += max(lc, 0) * qty
	} else {
		w.land_weight += wt
	}

	/*
	 *  Tue Dec  8 17:56:36 1998 -- Scott Turner
	 *
	 *  New wagons don't "ride" in the mountains.
	 *
	 */
	if rc != 0 && (item == item_new_wagon || item == item_war_wagon) && mountains {
		rc = 0
	}
	if rc != 0 {
		w.ride_cap += max(rc, 0) * qty
	} else {
		w.ride_weight += wt
	}

	if fc != 0 {
		w.fly_cap += max(fc, 0) * qty
	} else {
		w.fly_weight += wt
	}

	w.total_weight += wt

	if item_animal(item) != FALSE {
		w.animals += qty
	}
}

func determine_unit_weights(who int, w *weights, mountains bool) {
	assert(kind(who) == T_char)

	// zero out weights
	*w = weights{}

	add_item_weight(who, 1, w, mountains)
	for _, e := range loop_inventory(who) {
		add_item_weight(e.item, e.qty, w, mountains)
	}
}

func determine_stack_weights(who int, w *weights, mountains bool) {
	determine_unit_weights(who, w, mountains)

	var v weights
	for _, i := range loop_all_here(who) {
		determine_unit_weights(i, &v, mountains)
		w.total_weight += v.total_weight
		w.land_cap += v.land_cap
		w.ride_cap += v.ride_cap
		w.fly_cap += v.fly_cap
		w.land_weight += v.land_weight
		w.ride_weight += v.ride_weight
		w.fly_weight += v.fly_weight
		w.animals += v.animals
	}
}

func ship_weight(ship int) int {
	assert(kind(ship) == T_ship)

	var w weights
	sum := 0
	for _, i := range loop_char_here(ship) {
		determine_unit_weights(i, &w, false)
		sum += w.total_weight
	}
	return sum
}

func lookup(table [][]byte, s []byte) int {
	if len(s) == 0 {
		return -1
	}
	return lookup_bb(table, s)
}

func lookup_bb(table [][]byte, s []byte) int {
	if len(s) == 0 {
		return -1
	}
	for i := range table {
		if table[i] != nil && i_strcmp(s, table[i]) == 0 {
			return i
		}
	}
	return -1
}

func lookup_sb(table []string, s []byte) int {
	if len(s) == 0 {
		return -1
	}
	var tbl [][]byte
	for _, v := range table {
		tbl = append(tbl, []byte(v))
	}
	return lookup_bb(tbl, s)
}

func lookup_ss(table []string, s string) int {
	if len(s) == 0 {
		return -1
	}
	var tbl [][]byte
	for _, v := range table {
		tbl = append(tbl, []byte(v))
	}
	return lookup_bb(tbl, []byte(s))
}

func loyal_s(who int) string {
	switch loyal_kind(who) {
	case 0:
		return fmt.Sprintf("unsworn-%d", loyal_rate(who))
	case LOY_contract:
		return fmt.Sprintf("contract-%d", loyal_rate(who))
	case LOY_oath:
		return fmt.Sprintf("oath-%d", loyal_rate(who))
	case LOY_fear:
		return fmt.Sprintf("fear-%d", loyal_rate(who))
	case LOY_npc:
		return fmt.Sprintf("npc-%d", loyal_rate(who))
	case LOY_summon:
		return fmt.Sprintf("summon-%d", loyal_rate(who))
	}
	panic("!reached")
}

func gold_s(n int) string {
	return fmt.Sprintf("%s~gold", comma_num(n))
}

func weeks(n int) string {
	if n == 0 {
		return "0~days"
	} else if n%7 == 0 {
		n = n / 7
		return fmt.Sprintf("%s~week%s", nice_num(n), add_s(n))
	}
	return fmt.Sprintf("%s~day%s", nice_num(n), add_s(n))
}

func more_weeks(n int) string {
	if n == 0 {
		return "0~more days"
	} else if n%7 == 0 {
		n = n / 7
		return fmt.Sprintf("%d~more week%s", n, add_s(n))
	}
	return fmt.Sprintf("%d~more day%s", n, add_s(n))
}

func comma_num(n int) string {
	further := n / 1000000000
	n = n % 1000000000

	millions := n / 1000000
	n = n % 1000000

	thousands := n / 1000
	ones := n % 1000

	if further == 0 && millions == 0 && thousands == 0 {
		return fmt.Sprintf("%d", ones)
	} else if further == 0 && millions == 0 {
		return fmt.Sprintf("%d,%03d", thousands, ones)
	} else if further == 0 {
		return fmt.Sprintf("%d,%03d,%03d", millions, thousands, ones)
	}
	return fmt.Sprintf("%d,%03d,%03d,%03d", further, millions, thousands, ones)
}

func nice_num(n int) string {
	if 0 < n && n < 10 {
		return num_s[n]
	}
	return comma_num(n)
}

func knum(n int, nozero bool) string {
	if n == 0 && nozero {
		return ""
	} else if n < 9999 {
		return fmt.Sprintf("%d", n)
	} else if n < 1000000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return fmt.Sprintf("%dM", n/1000000)
}

func ordinal(n int) string {
	if 10 < n && n <= 19 {
		return fmt.Sprintf("%sth", comma_num(n))
	}
	switch n % 10 {
	case 1:
		return fmt.Sprintf("%sst", comma_num(n))
	case 2:
		return fmt.Sprintf("%snd", comma_num(n))
	case 3:
		return fmt.Sprintf("%srd", comma_num(n))
	default:
		return fmt.Sprintf("%sth", comma_num(n))
	}
}

func mylog_output(base, num int) int {
	assert(base > 10)

	power := 1
	num = num * 10
	for value := base; value < num; {
		power++
		value = value * base / 10
	}

	return power
}

func my_sqrt(n int) int {
	var power int
	for power = 1; power*power <= n; power++ {
		//
	}
	return power - 1
}

/* return a capitalized copy of s */
func cap_(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// deduct_np deducts `num` noble points from the player.
// it returns true if the player had enough points.
// if there weren't, it returns false and does not charge the player.
func deduct_np(pl, num int) bool {
	assert(kind(pl) == T_player)

	p := p_player(pl)
	if p == nil || p.noble_points < num {
		return false
	}
	p.noble_points -= num
	p.np_spent += num

	return true
}

func add_np(pl, num int) {
	assert(kind(pl) == T_player)

	p := p_player(pl)
	if p != nil {
		p.noble_points += num
		p.np_gained += num
	}
}

func deduct_aura(who, amount int) bool {
	p := rp_magic(who)
	if p == nil || p.cur_aura < amount {
		return false
	}
	p.cur_aura -= amount
	return true
}

func charge_aura(who, amount int) bool {
	if !deduct_aura(who, amount) {
		wout(who, "%s aura or piety required, current level is %s.", cap_(nice_num(amount)), nice_num(char_cur_aura(who)))
		return false
	}
	return true
}

func check_aura(who, amount int) bool {
	if char_cur_aura(who) < amount {
		wout(who, "%s aura or piety required, current level is %s.", cap_(nice_num(amount)), nice_num(char_cur_aura(who)))
		return false
	}

	return true
}

// has_item returns the quantity of the item in the character's possession.
func has_item(who, item int) int {
	assert(valid_box(who))
	//#if 0
	//    if (!valid_box(item)) {
	//        fprintf(stderr, "has_item(who=%s, item=%s) failure\n", box_code_less(who), box_code_less(item));
	//        fprintf(stderr, "player(who) = %s\n", box_code_less(player(who)));
	//        fprintf(stderr, "c.line '%s'\n", bx[who].cmd.line);
	//        assert(FALSE);
	//    }
	//#endif
	assert(valid_box(item))

	for i := 0; i < len(bx[who].items); i++ {
		if bx[who].items[i].item == item {
			return bx[who].items[i].qty
		}
	}

	return 0
}

func add_item(who, item, qty int) {
	assert(valid_box(who))
	assert(valid_box(item))
	assert(qty >= 0)

	lore := item_lore(item)
	if lore != 0 && kind(who) == T_char && !test_known(who, item) {
		queue_lore(who, item, false)
	}

	for i := 0; i < len(bx[who].items); i++ {
		if bx[who].items[i].item == item {
			old := bx[who].items[i].qty
			bx[who].items[i].qty += qty
			ifnNEW_TRADE(who, item, old)
			return
		}
	}

	newt := &item_ent{item: item, qty: qty}
	bx[who].items = append(bx[who].items, newt)
	ifnNEW_TRADE(who, item, 0)
}

func sub_item(who, item, qty int) bool {
	assert(valid_box(who))
	assert(valid_box(item))
	assert(qty >= 0)

	for i := 0; i < len(bx[who].items); i++ {
		if bx[who].items[i].item == item {
			if bx[who].items[i].qty < qty {
				return false
			}
			bx[who].items[i].qty -= qty
			return true
		}
	}

	return false
}

func gen_item(who, item, qty int) {
	assert(item_unique(item) == FALSE)
	add_item(who, item, qty)
}

func consume_item(who, item, qty int) bool {
	if item_unique(item) != FALSE {
		wout(gm_player, "Destroying unique item %s via consume_item.", box_name(item))
		destroy_unique_item(who, item)
		return true
	}
	assert(item_unique(item) == FALSE)
	return sub_item(who, item, qty)
}

/*
 *  Move item from one unit to another
 *  Destination=0 means discard the items
 */
func move_item(from, to, item, qty int) bool {
	if qty <= 0 {
		return true
	} else if to == 0 {
		return drop_item(from, item, qty)
	} else if sub_item(from, item, qty) {
		add_item(to, item, qty)
		if item_unique(item) != FALSE {
			assert(qty == 1)
			p_item(item).who_has = to

			//#if 0
			//if (subkind(item) == sub_npc_token)
			//   move_token(item, from, to);
			//#endif
		}
		return true
	}
	return false
}

func hack_unique_item(item, owner int) {
	p_item(item).who_has = owner
	add_item(owner, item, 1)
}

func create_unique_item(who, sk int) int {
	newt := new_ent(T_item, sk)
	if newt < 0 {
		return -1
	}
	if who != 0 {
		p_item(newt).who_has = who
		add_item(who, newt, 1)
	}
	return newt
}

func destroy_unique_item(who, item int) {
	assert(kind(item) == T_item)
	assert(item_unique(item) != FALSE)

	ret := sub_item(who, item, 1)
	assert(ret)

	delete_box(item)
}

func find_nearest_land(where int) int {
	orig_where := where
	var dir, check int

	for try_two := 100; try_two > 0; try_two-- {
		dir = rnd(1, 4)

		for try_one := 1000; try_one > 0; try_one-- {
			if subkind(where) != sub_ocean {
				return where
			}

			for _, i := range loop_here(where) {
				if subkind(i) == sub_island {
					assert(kind(i) == T_loc)
					if i != 0 {
						return i
					}
					break
				}
			}

			// todo: this does something, i'm sure...
			where = location_direction(where, dir)
			for where == 0 {
				where = orig_where
				dir = (dir % 4) + 1
				check++ // todo: should this be reset each loop? guess not.
				assert(check <= 4)
				where = location_direction(where, dir)
			}
		}

		if try_two == 100 {
			log_output(LOG_CODE, "find_nearest_land: Plan B")
		}
	}

	log_output(LOG_CODE, "find_nearest_land: Plan C")
	var l []int
	for _, i := range loop_loc() {
		if region(i) != region(orig_where) {
			continue
		} else if loc_depth(i) != LOC_province {
			continue
		} else if subkind(i) == sub_ocean {
			continue
		}
		l = append(l, i)
	}
	if len(l) == 0 {
		return 0
	}

	return l[rnd(0, len(l)-1)]
}

/*
 *  Simply throw away non-unique items
 *  Put unique items into the province to be found with EXPLORE
 *  If we're at sea, look for a nearby island or shore to move
 *  the item to.
 *
 *  Mon Apr 20 10:00:56 1998 -- Scott Turner
 *
 *  Province should collect peasants & gold.
 */
func drop_item(who, item, qty int) bool {
	if item_unique(item) == FALSE && item != item_peasant && item != item_gold {
		return consume_item(who, item, qty)
	}

	who_gets := province(who)
	if subkind(item) == sub_dead_body {
		who_gets = nearby_grave(who_gets)
		if who_gets == 0 {
			destroy_unique_item(who, item)
			return true
		}
	}

	if subkind(who_gets) == sub_ocean {
		who_gets = find_nearest_land(who_gets)
	}

	if who_gets == 0 {
		who_gets = subloc(who) /* oh well */
	}

	if item != item_gold && item != item_peasant {
		log_output(LOG_CODE, "drop_item: %s from %s to %s", box_name(item), box_name(subloc(who)), box_name(who_gets))
	}

	return move_item(who, who_gets, item, qty)
}

func can_pay(who, amount int) bool {
	return has_item(who, item_gold) >= amount
}

func charge(who, amount int) bool {
	return sub_item(who, item_gold, amount)
}

func stack_has_item(who, item int) (sum int) {
	head := stack_leader(who)
	for _, i := range loop_stack(head) {
		if player(i) != player(who) { /* friendly with us */
			continue
		}

		sum += has_item(i, item)
	}
	return sum
}

func has_use_key(who, key int) int {
	for _, e := range loop_inventory(who) {
		if p := rp_item_magic(e.item); p != nil && p.use_key == key && e.item != 0 {
			return e.item
		}
	}
	return 0
}

func stack_has_use_key(who, key int) int {
	head := stack_leader(who)
	for _, i := range loop_stack(head) {
		if player(i) != player(who) { /* friendly with us */
			continue
		}
		if ret := has_use_key(i, key); ret != 0 {
			return ret
		}
	}
	return 0
}

/*
 *  Subtract qty of item from a stack
 *  Take it from who first, then take it from
 *  anyone else in the stack who has it, starting from
 *  the stack leader and working down.
 *
 *  Return FALSE if the stack doesn't have qty of item.
 */
func stack_sub_item(who, item, qty int) bool {
	if stack_has_item(who, item) < qty {
		return false
	}
	if n := min(has_item(who, item), qty); n > 0 {
		qty -= n
		sub_item(who, item, n)
	}
	assert(qty >= 0)
	if qty == 0 {
		return true
	}

	// try to borrow what we need from friendly stackmates
	head := stack_leader(who)
	for _, i := range loop_stack(head) {
		if qty <= 0 {
			break
		} else if player(i) != player(who) { /* friendly with us */
			continue
		}

		if n := min(has_item(i, item), qty); n > 0 {
			qty -= n
			sub_item(i, item, n)

			//#if 0
			// if (show_day) {
			//    wout(who, "Borrowed %s from %s.",
			//    box_item_desc(item, n),
			//    box_name(i));
			//    wout(i, "%s borrowed %s.", box_name(who), box_item_desc(item, n));
			// }
			//#endif
		}
	}

	assert(qty == 0) /* or else stack_has_item above lied */

	return true
}

func autocharge(who, amount int) bool {
	return stack_sub_item(who, item_gold, amount)
}

func test_bit(kr sparse, i int) bool {
	return ilist_lookup(kr, i) != -1
}

func set_bit(kr sparse, i int) sparse {
	if ilist_lookup(kr, i) == -1 {
		return append(kr, i)
	}
	return kr
}

func clear_know_rec(kr sparse) sparse {
	panic("replace me with an assignment to nil")
	return nil
}

func test_known(who, i int) bool {
	if who == 0 {
		return false
	}

	assert(valid_box(who))
	assert(valid_box(i))

	ep := rp_player(player(who))
	if ep != nil && test_bit(ep.known, i) {
		return true
	}

	return false
}

func set_known(who, i int) {
	assert(valid_box(who))
	assert(valid_box(i))

	pl := player(who)
	assert(valid_box(pl))

	p_player(pl).known = set_bit(p_player(pl).known, i)
}

func print_dot(ch byte) {
	if dot_count == 0 {
		log.Printf("   ")
		dot_count++
	}
	dot_count++ // because we put two spaces in the buf

	if dot_count%60 == 0 {
		log.Printf("\n   ")
	}
	log.Printf("%c", ch)
}

func first_character(where int) int {
	for _, i := range loop_here(where) {
		if kind(i) == T_char {
			return i
		}
	}
	return 0
}

func entab(pl int) string {
	if player_notab(pl) {
		return ""
	}
	return "entab | "
}

func loc_hidden(n int) bool {
	//#if 0
	//    if (loc_depth(n) > LOC_province && weather_here(n, sub_fog))
	//        return TRUE;
	//#endif
	return rp_loc(n).hidden != FALSE
}

func loc_contains_hidden(n int) int {
	for _, enclosed := range loop_here(n) {
		if loc_hidden(enclosed) {
			return enclosed
		}
	}
	return 0
}

func rest_name(c *command, a int) string {
	if numargs(c) < a {
		return ""
	}
	s := string(c.parse[a])
	for i := a + 1; i <= numargs(c); i++ {
		s = fmt.Sprintf("%s %s", s, string(c.parse[i]))
	}
	return s
}

var nprov int

func nprovinces() int {
	if nprov == 0 {
		nprov = len(loop_province())
	}
	return nprov
}

func my_prisoner(who, pris int) bool {
	if kind(pris) != T_char {
		return false
	} else if !is_prisoner(pris) {
		return false
	} else if loc(pris) != who {
		return false
	}
	return true
}

func beast_capturable(who int) bool {
	if subkind(who) != sub_ni {
		return false
	} else if ni := noble_item(who); item_capturable(ni) != FALSE {
		return true
	}
	return false
}

func beast_wild(who int) bool {
	if subkind(who) != sub_ni {
		return false
	} else if ni := noble_item(who); item_wild(ni) != FALSE {
		return true
	}
	return false
}

var (
	_stage_old   time.Time
	_stage_first time.Time
)

func stage(s string) {
	if !time_self {
		if len(s) != 0 {
			log.Printf("%s\n", s)
		}
		return
	}

	t := time.Now()
	if _stage_old.IsZero() {
		_stage_first = t
	} else {
		log.Printf("\t%v\n", t.Sub(_stage_old))
	}
	_stage_old = t

	if s != "" {
		log.Printf("%s", s)
	} else {
		log.Printf("%v\n", t.Sub(_stage_first))
	}
}

/*
 *  Thu Jan  2 13:38:03 1997 -- Scott Turner
 *
 *  Ship capacity is now calculated from the ship's hulls, etc.
 *  For testing's sake, we also do the default values for capacity
 *  for galleys and roundships.
 *
 */
func ship_cap(ship int) int {
	s := rp_ship(ship)
	sc, dam := 0, loc_damage(ship)

	if s != nil {
		sc = (s.hulls * HULL_CAPACITY) - (s.forts * FORT_WEIGHT) - (s.sails * SAIL_WEIGHT) - (s.keels * KEEL_WEIGHT)
		if sc < 0 {
			sc = 0
		}
	} else if subkind(ship) == sub_galley {
		sc = 5000
	} else if subkind(ship) == sub_roundship {
		sc = 25000
	} else {
		sc = 0
	}

	if loc_hp(ship) != 0 {
		sc -= sc * dam / loc_hp(ship)
	}

	return sc
}

// prevent multiple TAGS running in same path
func lock_tag() {
	log.Printf("lock_tag: not implemented\n")
}
