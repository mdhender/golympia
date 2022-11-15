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
	"math"
)

type strength_function_t func(int) int

func controlled_humans_here(where int) int {
	var i int
	for _, i = range loop_all_here(where) {
		if kind(i) == T_char &&
			subkind(i) == 0 &&
			loyal_kind(i) != LOY_unsworn {
			return i
		}
	}
	return FALSE
}

func get_exit_dir(l []*exit_view, dir int) *exit_view {
	var i int

	for i = 0; i < len(l); i++ {
		if l[i].direction == dir { /* && l[i].hidden == FALSE? */
			return l[i]
		}
	}

	return nil
}

/*
 *  Tue Aug 11 06:05:19 1998 -- Scott Turner
 *
 *  Modified to permit monsters to avoid garrisons, cities, or
 *  civilized areas.
 *
 */

func choose_npc_direction(who, where, dir int, avoid_garrisons, avoid_cities, avoid_civ bool) *exit_view {
	var e *exit_view
	var l exit_view_l
	var i int

	l = exits_from_loc_nsew_select(who, where, LAND, RAND != FALSE)

	/*
	 *  Don't try to go someplace we can't enter.
	 *
	 */
	for i = 0; i < len(l); {
		if peaceful_enter(who, where, l[i].destination) {
			l = l.delete(i)
		} else {
			i++
		}
	}

	/*
	 *  Delete exits based on the conditions; don't advance
	 *  the list pointer if we delete.
	 *
	 */
	if avoid_garrisons {
		for i = 0; i < len(l); {
			if garrison_here(l[i].destination) != FALSE {
				l = l.delete(i)
			} else {
				i++
			}
		}
	}

	if avoid_cities {
		for i = 0; i < len(l); {
			if city_here(l[i].destination) != FALSE {
				l = l.delete(i)
			} else {
				i++
			}
		}
	}

	if avoid_civ {
		for i = 0; i < len(l); {
			if has_item(l[i].destination, item_peasant) > 100 {
				l = l.delete(i)
			} else {
				i++
			}
		}
	}

	if len(l) == 0 {
		return nil
	}

	/*
	 *  There is a 90% chance an NPC will keep going in the same
	 *  direction, if it can.
	 */

	if dir != FALSE {
		e = get_exit_dir(l, dir)
		if e != nil && e.impassable == FALSE && e.hidden == FALSE && e.magic_barrier == FALSE && rnd(1, 10) < 10 {
			return e
		}
	}

	for i = 0; i < len(l); i++ {
		if l[i].impassable == FALSE && l[i].hidden == FALSE && l[i].magic_barrier == FALSE {
			return l[i]
		}
	}

	return nil /* order of l has already been randomized */
}

func npc_move(who int) {
	var e *exit_view
	where := subloc(who)
	avoid_garrisons, avoid_cities, avoid_civ := false, true, false

	if loc_depth(where) != LOC_province {
		/*
		 *  Sat Nov 27 12:23:25 1999 -- Scott Turner
		 *
		 *  If you're on an Island you're screwed; if you're not the top
		 *  dog might as well suicide.
		 *
		 */
		if subkind(where) == sub_island && place_here(where, who) != FALSE {
			wout(who, "Island suicide.")
			queue(who, "die")
			return
		}
		queue(who, "move out")
		return
	}

	/*
	 *  Avoid some areas based on your "program"
	 *
	 */
	switch npc_program(who) {
	/*
	 *  Bandits stay away from the law but will enter
	 *  cities.
	 *
	 */
	case PROG_bandit:
		avoid_cities = false
		avoid_garrisons = true
		break
	case PROG_balrog:
		avoid_cities = false
		break
		/*
		 *  Dumb monsters avoid people.
		 *
		 */
	case PROG_dumb_monster:
		avoid_civ = true
		break
		/*
		 *  Smart monsters avoid cities and garrisons.
		 *
		 */
	case PROG_smart_monster:
		avoid_garrisons = true
		break
		/*
		 *  Orcs hate everything.
		 *
		 */
	case PROG_orc:
		avoid_cities = false
		break
	default:
		break
	}

	e = choose_npc_direction(who, where, npc_last_dir(who),
		avoid_garrisons, avoid_cities,
		avoid_civ)
	if e != nil {
		p_misc(who).npc_dir = e.direction
		queue(who, "move %s", full_dir_s[e.direction])
	}
}

func auto_unsworn(who int) {
	var n int
	where := subloc(who)

	if loc_depth(where) == LOC_build {
		return
	}

	if rnd(1, 2) == 1 {
		if n = city_here(where); n != FALSE && rnd(1, 2) == 1 {
			queue(who, "move %s", box_code_less(n))
		} else {
			npc_move(who)
		}
	}
}

func auto_mob(who int) {
	var p *entity_misc

	p = rp_misc(who)

	if p == nil {
		log.Printf("warning: mob's rp_misc is nil, who=%d\n", who)
		return
	}

	/*
	 *  Disperse if unstacked and not at home.
	 *  50% chance of dispersing each turn after five turns guarding
	 *
	 *  Since auto npc orders are only queued at the beginning of a turn,
	 *  a mob unstacked will appear in the end of turn location report.
	 *  Someone may try to rally the mob, so give them a chance before
	 *  dispersing the mob.
	 */

	if (subloc(who) != p.npc_home) ||
		(sysclock.turn-p.npc_created >= 5 && rnd(1, 2) == 1) {
		queue(who, "wait time %d", rnd(10, 20))
		queue(who, "reclaim \"disperses.\"")
		return
	}
}

func create_bandit(where int) int {
	var name string
	var item int
	var foo string

	switch rnd(1, 4) {
	case 1:
		item = item_bandit
		name = "Band of thieves"
		foo = "come"
		break

	case 2:
		item = item_wolf
		name = "Pack of wolves"
		foo = "come"
		break

	case 3:
		item = item_rat
		name = "Horde of rats"
		foo = "comes"
		break

	case 4:
		item = item_savage
		name = "Wild men"
		foo = "come"
		break

	default:
		panic("!reached")
	}
	// todo: create_bandit: name never used

	newt := create_monster_stack(item, rnd(4, 9), where)

	if newt < 0 {
		return -1
	}
	set_name(newt, name)

	wout(where, "%s %s out of hiding.", box_name(newt), foo)

	gen_item(newt, item_gold, rnd(25, 250))

	return newt
}

func wilderness_attack_check(who, where int) {
	if rnd(1, 16) > 1 ||
		char_really_hidden(who) ||
		is_npc(who) ||
		kind(who) != T_char ||
		count_stack_units(who) > 1 ||
		count_stack_figures(who) > 10 ||
		char_really_hidden(who) || /* todo: this is here twice */
		char_gone(who) != FALSE {
		return
	}
	/*
	 *  No new beasts to attack Priests of Dol
	 *
	 */
	if priest_in_stack(who, sk_dol) {
		return
	}

	newt := create_bandit(where)

	if newt < 0 {
		return
	}

	queue(newt, "wait time 0")
	//#if 0
	//    init_load_sup(new);   /* make ready to execute commands immediately */
	//#endif

	/*
	 *  ART_SAFETY
	 *
	 */
	if FALSE == has_artifact(who, ART_SAFETY, noble_item(newt), 0, 0) &&
		rnd(1, 2) == 1 {
		queue(newt, "attack %s", box_code_less(who))
	}
}

/*
 *   Wed Apr 22 11:33:44 1998 -- Scott Turner
 *
 *   A little bit of overkill.
 *
 *  Tue Jul  6 13:45:47 1999 -- Scott Turner
 *
 *  Don't execute LOY_unsworn/indep_player nobles (e.g., hostages)
 *
 */
func execute_prisoners(who int, reason string) {
	var pris int

	for _, pris = range loop_here(who) {
		if is_prisoner(pris) && (loyal_rate(pris) != LOY_unsworn ||
			player(pris) != indep_player) {
			wout(viewloc(subloc(who)), reason, box_name(who), box_name(pris))
			queue(who, "execute %s", box_code_less(pris))
		}
	}

}

func release_prisoners(who int, reason string) {
	var pris int
	found := false

	for _, pris = range loop_here(who) {
		if is_prisoner(pris) {
			wout(viewloc(subloc(who)), reason, box_name(who), box_name(pris))
			queue(who, "unstack %s", box_code_less(pris))
			found = true
		}
	}

	/*
	 *  Move away, so we don't automatically re-attack.
	 *
	 */
	if found {
		npc_move(who)
	}
}

/*
 *  Tue Apr 27 18:30:15 1999 -- Scott Turner
 *
 *  Okay to breed when you're the top dog in your province
 *  within 25% of your MM.
 *
 *  Tue Oct 12 09:52:30 1999 -- Scott Turner
 *
 *  You also have to be the only stack of your sort in the
 *  province.
 */
func breed_okay(who int) bool {
	mm := MM(who)
	where := subloc(who)
	var i, item int
	//first := true

	if mm < 1 || mm > MAX_MM {
		return false
	}

	/*
	 *  Wed Nov 24 17:37:45 1999 -- Scott Turner
	 *
	 *  Don't breed if there's another similar pack here.
	 *
	 */
	for _, i = range loop_here(where) {
		if i == who {
			break
		}
		if noble_item(i) == noble_item(who) {
			return false
		}
	}

	/*
	 *  Check the units in this location to see if I'm
	 *  the top dog in my ecological "niche"
	 *
	 *  Wed Nov 24 07:48:35 1999 -- Scott Turner
	 *
	 *  Widening the niche.
	 *
	 */
	if first_character(where) != stack_leader(who) {
		return false
	}

	/*
	 *  I can't breed if I'm more than 2x my split limit.
	 *
	 */
	item = noble_item(who)
	if has_item(who, item) > item_split(item)*2 {
		return false
	}

	/*
	 *  I can breed.
	 *
	 */
	return true
}

/*
 *  Wed Jun 30 10:37:43 1999 -- Scott Turner
 *
 *  Drop useless stuff.
 *
 */
func auto_drop(who int) bool {
	var e *item_ent
	found_item := false

	for _, e = range loop_inventory(who) {
		if e.item != item_gold &&
			subkind(e.item) != sub_trade_good && /* todo: was "!subkind(...) == ... */
			FALSE == item_unique(e.item) &&
			FALSE == item_attack(e.item) {
			queue(who, "drop %s %d", box_code_less(e.item), e.qty)
			found_item = true
		}
	}

	return found_item
}

/*
 *  Thu Oct  3 11:49:02 1996 -- Scott Turner
 *
 *  Capturable beasts don't attack priests of Dol
 *
 */
func auto_bandit(who int) {
	where := subloc(who)
	var i int
	victim := 0
	var spend int

	/*
	 *  If we have stuff in our inventory, get rid of it.  Spend
	 *  any money we find.  Unstack nobles.
	 */
	release_prisoners(who, "%s releases %s.")
	if auto_drop(who) {
		return
	}

	if has_item(who, item_gold) > 0 {
		/* Spend some... 1-3 per man. */
		spend = rnd(has_item(who, noble_item(who))+1,
			3*has_item(who, noble_item(who))+3)
		if spend > has_item(who, item_gold) {
			spend = has_item(who, item_gold)
		}
		consume_item(who, item_gold, spend)
		queue(who, "wait time 7")
		return
	}

	for _, i = range loop_here(where) {
		if kind(i) == T_char && !is_real_npc(i) &&
			count_stack_units(i) == 1 &&
			count_stack_figures(i) <= has_item(who, noble_item(who)) &&
			FALSE == char_gone(i) &&
			!char_really_hidden(i) {
			if beast_capturable(who) && priest_in_stack(i, sk_dol) {
				continue
			}

			victim = i
			break
		}
	}

	if victim != FALSE && FALSE == has_artifact(victim, ART_SAFETY, noble_item(who), 0, 0) {
		queue(who, "attack %s", box_code_less(victim))
	} else if rnd(1, 2) == 1 && breed_okay(who) {
		queue(who, "breed")
	} else {
		npc_move(who)
	}
}

/*
 *  Wed Nov 24 06:10:00 1999 -- Scott Turner
 *
 *  Grab a subloc if you're in one; look for one if you're not.
 *
 */
func grab_subloc(who int) bool {
	where := subloc(who)
	var i int

	if npc_program(who) == PROG_subloc_monster {
		return false
	}

	/*
	 *  If we're in a subloc, we can establish ourselves as a subloc monster.
	 *
	 */
	if loc_depth(where) > LOC_province &&
		first_character(where) == stack_leader(who) {
		rp_char(who).npc_prog = PROG_subloc_monster
		wout(who, "Established as subloc monster in %s.", box_name(where))
		return true
	}

	/*
	 *  Perhaps there's an empty subloc here to move into?
	 *
	 */
	for _, i = range loop_here(where) {
		if kind(i) == T_loc &&
			loc_depth(i) > LOC_province &&
			subkind(i) != sub_city &&
			FALSE == first_character(i) {
			wout(who, "Moving into %s to establish residency.",
				box_name(i))
			queue(who, "move %s", box_code_less(i))
			rp_char(who).npc_prog = PROG_subloc_monster
			return true
		}
	}

	return false
}

/*
 *  Fri Nov 26 06:29:08 1999 -- Scott Turner
 *
 *  Merge with another stack of the same type if it won't
 *  put you up to the split size.
 *
 *
 */
func merge_stacks(who int) bool {
	where := subloc(who)
	var i int

	/*
	 *  Perhaps there's an empty subloc here to move into?
	 *
	 */
	for _, i = range loop_char_here(where) {
		if i != who &&
			is_real_npc(i) &&
			noble_item(who) == noble_item(i) &&
			has_item(i, noble_item(i))+has_item(who, noble_item(who)) <
				(item_split(noble_item(i))-2) {
			wout(who, "Merging stacks with %s.", box_name(i))
			queue(who, "stack %s", box_code_less(i))
			queue(who, "give %s %s %d 0", box_code_less(i),
				box_code_less(noble_item(who)),
				has_item(who, noble_item(i)))
			queue(who, "die")
			return true
		}
	}

	return false
}

/*
 *  Fri Nov 26 08:13:01 1999 -- Scott Turner
 *
 *  Estimate strength of a stack -- dumb version.
 *
 */
func dumb_estimate(i int) int {
	return count_stack_any(i)
}

func estimate_strength(who, i int, strength_function strength_function_t) int {
	/*
	 *  If there's a subloc here, we'll assume the first character
	 *  there will defend it alone.  Probably suicidal, but hey.
	 *
	 */
	if kind(i) == T_loc && first_character(i) != FALSE {
		return estimate_strength(who, first_character(i), strength_function)
	}
	/*
	 *  Otherwise, make sure he's suitable before evaluating.
	 *
	 */
	if kind(i) != T_char || is_real_npc(i) {
		return 0
	}
	if beast_capturable(who) && priest_in_stack(i, sk_dol) {
		return 0
	}
	if has_artifact(i, ART_SAFETY, noble_item(who), 0, 0) != FALSE {
		return 0
	}
	if char_really_hidden(i) || char_gone(i) != FALSE {
		return 0
	}

	return strength_function(i)

}

/*
 *  Fri Nov 26 08:20:43 1999 -- Scott Turner
 *
 *  Find and attack the weakest stack here, if you think
 *  you can "take" them.
 *
 *  Fri Aug 18 12:03:54 2000 -- Scott Turner
 *
 *  Don't attack anyone in Hades with "Transcend Death" if they're alone.
 *
 */
func find_weakest(who, where int, strength_function strength_function_t) int {
	var strength, victim_strength int
	victim := 0
	var i int

	/*
	 *  Attack the weakest person hanging around here.
	 *
	 */
	for _, i = range loop_here(where) {

		/*
		 *  Don't consider wizards in Hades.
		 *
		 */
		if in_hades(who) &&
			has_skill(i, sk_transcend_death) != FALSE &&
			char_alone(i) {
			continue
		}

		strength = estimate_strength(who, i, strength_function)
		if strength == FALSE {
			continue
		}

		if victim == FALSE || strength < victim_strength {
			victim = i
			victim_strength = strength
		}
	}

	/*
	 *  If we're stronger than the weakest guy, attack him.
	 *
	 */
	if victim != FALSE && strength_function(who) > int(math.Ceil(1.25*float64(victim_strength))) {
		return victim
	}

	return FALSE
}

/*
 *  Fri Jan 19 10:02:48 2001 -- Scott Turner
 *
 *  Go through all our exits, and if there's an exit that we cannot
 *  enter peacefully, assess whether we might want to jump on it.
 *
 */
func find_weak_garrison(who int, strength_function strength_function_t) bool {
	var l []*exit_view
	where := subloc(who)
	var i int
	my_strength := strength_function(who)

	/*
	 *  Note: L is randomized.
	 *
	 */
	l = exits_from_loc_nsew_select(who, where, LAND, RAND != FALSE)

	/*
	 *  Any garrisons to attack?
	 *
	 */
	for i = 0; i < len(l); i++ {
		if peaceful_enter(who, where, l[i].destination) &&
			strength_function(controls_loc(l[i].destination)) <= my_strength {
			queue(who, "move %s", box_code_less(l[i].destination))
			queue(who, "attack %s", box_code_less(l[i].destination))
			return true
		}
	}
	return false
}

/*
 *  Mon Mar  3 12:50:31 1997 -- Scott Turner
 *
 *  Dumb beasts wander around until they find an empty location and then
 *  they sit there.
 *
 */
func auto_dumb(who int) {
	where := subloc(who)
	var first, victim int

	/*
	 *  Fri Nov 26 08:21:36 1999 -- Scott Turner
	 *
	 *  Attack anyone who looks weak.
	 *
	 */
	if victim = find_weakest(who, subloc(who), dumb_estimate); victim != FALSE {
		queue(who, "attack %s", box_code_less(victim))
		return
	}

	/*
	 *  Fri Jan 19 09:59:06 2001 -- Scott Turner
	 *
	 *  Maybe there's a garrison blocking our path into an adjacent
	 *  province.  If so, we ought to beat it up if we can and head
	 *  over there -- there's probably some tasty peasants, etc.
	 *
	 */
	if find_weak_garrison(who, dumb_estimate) {
		return
	}

	/*
	 *  Are there peasants here to pillage?
	 *
	 */
	if count_fighters(who, item_attack(item_peasant)) >= 10 &&
		has_item(where, item_peasant) > 100 && rnd(1, 100) < 50 {
		queue(who, "pillage 1") /* Attack the guards if necessary. */
		return
	}

	/*
	 *  Might be a subloc to grab.
	 *
	 */
	if grab_subloc(who) {
		return
	}

	/*
	 *  Maybe we should merge with someone?
	 *
	 */
	if merge_stacks(who) {
		return
	}

	/*
	 *  Next, if we're the top dog we can just stay here.
	 *
	 */
	if first_character(where) == stack_leader(who) {
		if breed_okay(who) && rnd(1, 100) < 75 {
			queue(who, "breed")
		} else {
			queue(who, "wait time 7")
		}
		return
	}

	/*
	 *  If we're not top dog, but we're bigger than the top dog,
	 *  maybe we'll just attack him and grab his spot.
	 *
	 */
	first = first_character(where)
	if rnd(1, 10) == 1 &&
		first != stack_leader(who) &&
		is_real_npc(first) &&
		noble_item(first) != noble_item(who) &&
		dumb_estimate(who) > (2*dumb_estimate(first)) {
		queue(who, "attack %s", box_code_less(first))
		return
	}

	/*
	 *  If we have stuff in our inventory, get rid of it.  We just
	 *  drop the stuff, since we're dumb.  If we have a person in our
	 *  inventory, we "execute" (i.e., eat) him.  Others we keep as
	 *  fighters in our cause.
	 *
	 */
	if rnd(1, 10) == 1 {
		execute_prisoners(who, "%s feasts on the body of %s.")
	} else {
		release_prisoners(who, "%s releases %s.")
	}
	if auto_drop(who) {
		return
	}

	/*
	 *  Otherwise we'll be moving along.
	 *
	 */
	npc_move(who)

	return
}

/*
 *  Mon Mar  3 14:06:58 1997 -- Scott Turner
 *
 *  Estimate a stack's strength, assuming limited control of men
 *  and beasts, etc.
 *
 *  Fri Nov 26 08:28:32 1999 -- Scott Turner
 *
 *  Added in a factor for the sheer size of the stack.
 *
 */
func smart_estimate(who int) int {
	var i int
	var e *item_ent
	combat_total := 0
	var num_beasts, num_men, num int

	for _, i = range loop_stack(who) {
		if !is_npc(i) {
			num_beasts = 15
			num_men = DEFAULT_CONTROLLED
		} else if subkind(who) == sub_garrison {
			num_beasts = 0
			num_men = GARRISON_CONTROLLED
		} else {
			num_men = 10000
			num_beasts = 10000
		}
		for _, e = range loop_sorted_inv(i) {
			if is_fighter(e.item) != FALSE {
				if item_animal(e.item) != FALSE {
					num = min(num_beasts, e.qty)
					num_beasts -= num
				} else {
					num = min(num_men, e.qty)
					num_men -= num
				}
				combat_total += num * (item_attack(e.item) + item_defense(e.item))
			}
		}

		/*
		 *  Plus the noble himself.
		 *
		 */
		combat_total += char_attack(i) + char_defense(i)
	}
	return combat_total * count_stack_figures(who)
}

/*
 *  Fri Apr 10 06:39:54 1998 -- Scott Turner
 *
 *  Learning skills for npcs.
 *
 */
func study_skill(who, sk int) {
	if has_skill(who, sk) != FALSE {
		return
	}

	/*
	 *  Might need to learn the parent skill.
	 *
	 */
	if sk != skill_school(sk) &&
		FALSE == has_skill(who, skill_school(sk)) {
		study_skill(who, skill_school(sk))
		return
	}

	/*
	 *  Spend the time and learn the skill.
	 *
	 */
	queue(who, "study %s", box_code_less(sk))
}

func use_skill(who, skill int) {
	var buf string

	if FALSE == has_skill(who, skill) {
		study_skill(who, skill)
		return
	}
	buf = fmt.Sprintf("use %d", skill)
	queue(who, buf)
}

func need_to_study(who, skill int) bool {
	var p *skill_ent

	p = rp_skill_ent(who, skill)

	if p != nil && p.know < SKILL_know {
		return true
	}

	return false
}

/*
 *  Fri Nov 26 19:04:52 1999 -- Scott Turner
 *
 *  Raze any structures you find.
 *
 */
func raze_structures(who int) bool {
	where := subloc(who)
	var i int

	if savage_hates(where) && building_owner(where) == who {
		queue(who, "raze")
		return true
	}

	for _, i = range loop_here(where) {
		if savage_hates(i) &&
			FALSE == first_character(i) {
			queue(who, "raze %s", box_code_less(i))
			return true
		}
	}

	return false
}

/*
 *  Mon Mar  3 13:19:14 1997 -- Scott Turner
 *
 *  Smart beasts are just a little cleverer about what they do.
 *
 */
func auto_smart(who int) {
	where := subloc(who)
	victim := 0

	/*
	 *  Fri Nov 26 08:24:54 1999 -- Scott Turner
	 *
	 *  Possibly attack the weakest stack here.
	 *
	 */
	if victim = find_weakest(who, subloc(who), smart_estimate); victim != FALSE {
		queue(who, "attack %s", box_code_less(victim))
		return
	}

	/*
	 *  Fri Jan 19 09:59:06 2001 -- Scott Turner
	 *
	 *  Maybe there's a garrison blocking our path into an adjacent
	 *  province.  If so, we ought to beat it up if we can and head
	 *  over there -- there's probably some tasty peasants, etc.
	 *
	 */
	if find_weak_garrison(who, smart_estimate) {
		return
	}

	/*
	 *  Are there peasants here to pillage?
	 *
	 */
	if count_fighters(who, item_attack(item_peasant)) >= 10 &&
		has_item(where, item_peasant) > 100 && rnd(1, 100) < 50 {
		queue(who, "pillage 1") /* Attack the guards if necessary. */
		return
	}

	/*
	 *  Are there structures to destroy?
	 *
	 */
	if raze_structures(who) {
		return
	}

	/*
	 *  Might be a subloc to grab.
	 *
	 */
	if grab_subloc(who) {
		return
	}

	/*
	 *  Maybe we should merge with someone?
	 *
	 */
	if merge_stacks(who) {
		return
	}

	/*
	 *  Next, if we're the top dog we can just stay here.  We'll
	 *  try to breed, but failing that we can study.
	 *
	 */
	if first_character(where) == stack_leader(who) {
		if breed_okay(who) && rnd(1, 10) < 6 {
			queue(who, "breed")
		} else {
			/*
			 *  First continue learning anything we don't know yet.
			 *
			 */
			if need_to_study(who, sk_attack_tactics) {
				use_skill(who, sk_attack_tactics)
				return
			}
			if need_to_study(who, sk_defense_tactics) {
				use_skill(who, sk_defense_tactics)
				return
			}
			/*
			 *   Otherwise learn or study a new one.
			 *
			 */
			if rnd(1, 2) == 1 {
				use_skill(who, sk_attack_tactics)
				return
			} else {
				use_skill(who, sk_defense_tactics)
				return
			}
			queue(who, "wait time 7") // todo: unreachable code
			return
		}
		return
	}

	/*
	 *  If we have stuff in our inventory, get rid of it.  We just drop
	 *  the stuff, since we're dumb.  If we have a person in our
	 *  inventory, we "execute" him, so that he doesn't come back to
	 *  bother us again.  Others we keep as fighters in our cause.
	 *
	 *  We like money, though :-)
	 */
	if rnd(1, 10) == 5 {
		execute_prisoners(who, "%s feasts on the body of %s.")
	} else {
		release_prisoners(who, "%s releases %s.")
	}
	if auto_drop(who) {
		return
	}

	/*
	 *  Otherwise we'll be moving along.
	 *
	 */
	npc_move(who)
	return
}

/*
 *  Sun Mar  9 21:05:39 1997 -- Scott Turner
 *
 *  Orcs build (or occupy) little orc strongholds!
 *
 *  Wed Mar 12 10:51:00 1997 -- Scott Turner
 *
 *  Orcs holed up in strongholds should look into the province
 *  for possible victims, move out, whack them and move back.
 *
 */
func auto_orc(who int) {
	found := false
	var victim, i int
	//int real_orc_loc_okay(who, where int);
	where := subloc(who)

	/*
	 *  If we have stuff in our inventory, get rid of it.  We just drop
	 *  the stuff, since we're dumb.  If we have a person in our
	 *  inventory, we "execute" him, so that he doesn't come back to
	 *  bother us again.  Others we keep as fighters in our cause.
	 *
	 *  We like money, though :-)
	 */
	execute_prisoners(who, "%s chops %s into pieces and defiles the body.")
	if auto_drop(who) {
		return
	}

	/*
	 *  If there's an empty orc stronghold here, let's
	 *  occupy it.  If there's an unfinished, we'll occupy it
	 *  and finish it.
	 *
	 */
	for _, i = range loop_here(where) {
		if subkind(i) == sub_orc_stronghold ||
			subkind(i) == sub_orc_stronghold_notdone {
			found = true
		}
		if subkind(i) == sub_orc_stronghold &&
			FALSE == building_owner(i) {
			queue(who, "enter %s", box_code_less(i))
			queue(who, "border %s close", box_code_less(i))
			return
		} else if subkind(i) == sub_orc_stronghold_notdone &&
			FALSE == building_owner(i) {
			queue(who, "enter %s", box_code_less(i))
			queue(who, "border %s close", box_code_less(i))
			queue(who, "build stronghold")
			return
		}
	}

	/*
	 *  If there's not a stronghold here, and it's a suitable
	 *  location, then build one.
	 *
	 *
	 */
	if !found && real_orc_loc_okay(who, where) {
		queue(who, "build stronghold")
		return
	}

	/*
	 *  If we're in a stronghold, then look out to the outer
	 *  level for possible victims.
	 *
	 */
	if subkind(where) == sub_orc_stronghold {
		if victim = find_weakest(who, subloc(where), smart_estimate); victim != FALSE {
			queue(who, "move %s", box_code_less(subloc(where)))
			queue(who, "attack %s", box_code_less(victim))
			queue(who, "move %s", box_code_less(where))
			queue(who, "border %s close", box_code_less(where))
			return
		}
	}

	/*
	 *  Otherwise, be a smart monster.
	 *
	 */
	auto_smart(who)

}

/*
 *  Mon Dec  7 18:23:59 1998 -- Scott Turner
 *
 *  Returns TRUE if it scheduled an attack on someone; that
 *  may require prepending an order!
 *
 */
func elf_attack(who, where int) bool {
	var i, j int
	p := p_misc(faery_player)
	assert(p != nil)

	/*
	 *  If a (sub-)location then check everything there.
	 *
	 */
	if kind(where) == T_loc {
		wout(who, "Checking out subloc %s.", box_name(where))
		for _, i = range loop_here(where) {
			if elf_attack(who, i) {
				wout(who, "Faery hunt prepending move to %s.",
					box_name(i))
				prepend_order(player(who), who, sout("attack %s", box_code_less(i)))
				return true
			}
		}
		return false
	}

	if kind(where) != T_char || subkind(player(where)) != sub_pl_regular {
		return false
	}

	/*
	 *  Go through every noble in the stack.  We don't warn (nor attack)
	 *  nobles carrying elfstones or Priests of Timeid.
	 *
	 */
	for _, j = range loop_stack(where) {
		if has_use_key(j, use_faery_stone) != FALSE ||
			has_artifact(j, ART_PROT_FAERY, 0, 0, 1) != FALSE {
			queue(who, "message 1 %s", box_code_less(j))
			queue(who, "Well met, elf friend.")
			continue
		}
		if is_priest(j) == sk_timeid {
			queue(who, "message 1 %s", box_code_less(j))
			queue(who, "The peace of the Lady to you, priest.")
			continue
		}
		if test_bit(p.npc_memory, j) {
			queue(who, "message 1 %s", box_code_less(j))
			queue(who, "You are not welcome in Faery.  Leave or die.")
			wout(who, "Faery hunt warned %s.", box_name(j))
			p.npc_memory = set_bit(p.npc_memory, j)
			add_effect(j, ef_faery_warning, 0, 30, 1)
		} else if FALSE == get_effect(j, ef_faery_warning, 0, 0) {
			log_output(LOG_SPECIAL, "Faery hunt attacking %s.", box_name(j))
			queue(who, "attack %s", box_code_less(j))
			return true
		}
	}
	return true
}

/*
 *  Tue Oct  6 07:19:43 1998 -- Scott Turner
 *
 *  Elves are a special case.  They won't attack anyone who is an
 *  elf friend, but they're brutal to everyone else.
 *
 */
func auto_elf(who int) {
	where := subloc(who)

	/*
	 *  Whenever we think, the first thing we do is look around and note
	 *  if there are any intruders around.  We warn them and add them to
	 *  our memory.  This is a little tricky, since we need to go
	 *  through and put out the bounty on everyone we see in the stack.
	 *
	 */
	if elf_attack(who, where) {
		return
	}

	/*
	 *  If we have stuff in our inventory, get rid of it.  We just drop
	 *  the stuff, since we're dumb.  If we have a person in our
	 *  inventory, we "execute" him, so that he doesn't come back to
	 *  bother us again.  Others we keep as fighters in our cause.
	 *
	 *  We like money, though :-)
	 */
	if rnd(1, 10) == 5 {
		execute_prisoners(who, "%s burns %s at a pyre to the Elven Gods.")
	} else {
		release_prisoners(who, "%s releases %s.")
	}
	if auto_drop(who) {
		return
	}

	/*
	 *  If you're less than 25, then just breed, breed, breed.
	 *
	 */
	if breed_okay(who) && rnd(1, 10) < 6 {
		queue(who, "breed")
		return
	}

	/*
	 *  Next, if we're the top dog we can just stay here.
	 *
	 */
	if first_character(where) == stack_leader(who) {
		/*
		 *  First continue learning anything we don't know yet.
		 *
		 */
		if need_to_study(who, sk_attack_tactics) {
			use_skill(who, sk_attack_tactics)
			return
		}
		if need_to_study(who, sk_defense_tactics) {
			use_skill(who, sk_defense_tactics)
			return
		}
		/*
		 *   Otherwise learn or study a new one.
		 *
		 */
		if rnd(1, 10) < 6 {
			use_skill(who, sk_attack_tactics)
			return
		} else {
			use_skill(who, sk_defense_tactics)
			return
		}
		queue(who, "wait time 7")
		return
	}

	/*
	 *  Otherwise we'll be moving along.
	 *
	 */
	npc_move(who)
	return

}

/*
 *  Tue Feb  4 12:35:13 1997 -- Scott Turner
 *
 *  Kill and attack everyone...
 *
 */
func auto_balrog(who int) {
	where := subloc(who)
	var i int

	/*
	 *  Execute prisoners.
	 *
	 */
	execute_prisoners(who, "%s burns %s and consumes the body.")

	/*
	 *  If we have stuff in our inventory, get rid of it.  We just drop
	 *  the stuff, since we're dumb.  If we have a person in our
	 *  inventory, we "execute" him, so that he doesn't come back to
	 *  bother us again.  Others we keep as fighters in our cause.
	 *
	 *  We like money, though :-)
	 */
	if auto_drop(who) {
		return
	}

	for _, i = range loop_here(where) {
		/*
		 *  Don't consider wizards in Hades.
		 *
		 */
		if in_hades(who) &&
			has_skill(i, sk_transcend_death) != FALSE &&
			char_alone(i) {
			continue
		}

		if kind(i) == T_char && !is_npc(i) &&
			FALSE == char_gone(i) &&
			!char_really_hidden(i) &&
			FALSE == has_artifact(i, ART_SAFETY, noble_item(who), 0, 0) {
			queue(who, "attack %s", box_code_less(i))
			return
		}
	}

	npc_move(who)
}

/*
 *  Tue Oct  6 18:51:25 1998 -- Scott Turner
 *
 *  Programs for daemons.
 *
 */
func auto_daemon(who int) {
	where := subloc(who)
	var i int

	/*
	 *  Execute prisoners.
	 *
	 */
	execute_prisoners(who, "%s defiles %s and consumes the body.")

	/*
	 *  If we have stuff in our inventory, get rid of it.  We just drop
	 *  the stuff, since we're dumb.  If we have a person in our
	 *  inventory, we "execute" him, so that he doesn't come back to
	 *  bother us again.  Others we keep as fighters in our cause.
	 *
	 *  We like money, though :-)
	 */
	if auto_drop(who) {
		return
	}

	for _, i = range loop_here(where) {
		if kind(i) == T_char && !is_npc(i) &&
			FALSE == char_gone(i) &&
			!char_really_hidden(i) &&
			!(has_skill(i, sk_transcend_death) != FALSE && char_alone(i)) &&
			!(in_hades(who) && has_artifact(i, ART_PROT_HADES, 0, 0, 0) != FALSE) &&
			FALSE == has_artifact(i, ART_SAFETY, noble_item(who), 0, 0) {
			queue(who, "attack %s", box_code_less(i))
			return
		}
	}

	npc_move(who)

}

/*
 *  Wed Jun 30 07:44:43 1999 -- Scott Turner
 *
 *  Programs for subloc monsters.  Basically they're going to sit
 *  and breed until they're full; then if they're smart they'll start
 *  studying things.
 *
 */

func count_sub_artifact(item int) int {
	if !is_sub_artifact(item) {
		return 0
	}
	return 1
}

func is_sub_artifact(item int) bool {
	if valid_box(item) && subkind(item) == sub_magic_artifact {
		return true
	}
	return false
}

//int count_generic(int who, int stack, int (*func)(int));

func auto_subloc(who int) {
	where := subloc(who)
	var target int

	/*
	 *  It's unlikely, but if there's anyone here to attack...
	 *
	 */
	if target = find_weakest(who, subloc(who), smart_estimate); target != FALSE {
		queue(who, "Subloc attack.")
		queue(who, "attack %s", box_code_less(target))
		return
	}

	/*
	 *  If we're not the top dog in our subloc then we should
	 *  modify our program and leave, don't you think?
	 *
	 */
	if first_character(where) != stack_leader(who) {
		wout(who, "Abandoning my subloc %s.", box_name(where))
		/* Temporary -- don't flood world with subloc monsters */
		/*FIXME: move out some percentage of time, but not every time.  Need to
		  check to see if monsters are randomly generated in the wild and
		  if so it may not be necessary to move monsters out*/

		/*queue(who,"die");    */
		queue(who, "move out")
		rp_char(who).npc_prog = int(item_prog(noble_item(who)))
		return
	}

	/*
	 *  If our subloc isn't closed, then try to close it.
	 *
	 */
	if where != province(where) && !p_subloc(where).control.closed {
		queue(who, "border %s close", box_code_less(where))
	}

	/*
	 *  If we have stuff in our inventory, get rid of it.
	 *
	 *  We like money, though :-)
	 *
	 *  Sun Nov 28 12:47:30 1999 -- Scott Turner
	 *
	 *  Subloc monsters don't slay prisoners; they're treasure!
	 *
	 */
	if auto_drop(who) {
		return
	}

	/*
	 *  If we're under our limit, we breed.
	 *
	 */
	//item := noble_item(who)
	if breed_okay(who) && rnd(1, 10) < 6 {
		queue(who, "breed")
		return
	}

	/*
	 *  Perhaps we can generate a little treasure?
	 *
	 */
	if rnd(1, 200) == 100 && count_generic(who, false, count_sub_artifact) < 2 {
		wout(who, "Found a big treasure!")
		generate_one_treasure(who)
	} else {
		gen_item(who, item_gold, rnd(1, 5))
	}

	/*
	 *  If we're dumb, that's it.  We can just wait.
	 *
	 */
	if item_wild(noble_item(who)) != PROG_smart_monster {
		queue(who, "wait time 7")
		return
	}

	/*
	 *  Otherwise, we're smart and we can do some learning.  First
	 *  continue learning anything we don't know yet.
	 *
	 */
	if need_to_study(who, sk_attack_tactics) {
		use_skill(who, sk_attack_tactics)
		return
	}
	if need_to_study(who, sk_defense_tactics) {
		use_skill(who, sk_defense_tactics)
		return
	}
	/*
	 *   Otherwise learn or study a new one.
	 *
	 */
	if rnd(1, 10) < 6 {
		use_skill(who, sk_attack_tactics)
		return
	} else {
		use_skill(who, sk_defense_tactics)
		return
	}
	queue(who, "wait time 7")
	return
}

const PROV_OR_CITY = -1

type cookie_monster_tbl struct {
	cookie              int
	kind, sk, ni        int
	terrain             int
	man_kind, low, high int
	not_here            string
	no_cookies          string
}

var cookie_monster = []cookie_monster_tbl{
	{
		item_mob_cookie,
		T_char, sub_ni, item_angry_peasant,
		PROV_OR_CITY,
		item_angry_peasant, 12, 36,
		"Mobs can only be raised in provinces and cities.",
		"A mob has already been raised from this place.",
	},
	{
		item_undead_cookie,
		T_char, sub_demon_lord, 0,
		sub_graveyard,
		item_corpse, 15, 25,
		"Demon lords may only be summoned in graveyards.",
		"A demon lord has already been summoned from this graveyard.",
	},
	{
		item_rain_cookie,
		T_storm, sub_rain, 0,
		0,
		0, 0, 0,
		"Rain may not be summoned here.",
		"There is no storm to be summoned from this province.",
	},
	{
		item_wind_cookie,
		T_storm, sub_wind, 0,
		0,
		0, 0, 0,
		"Wind may not be summoned here.",
		"There is no storm to be summoned from this province.",
	},
	{
		item_fog_cookie,
		T_storm, sub_fog, 0,
		0,
		0, 0, 0,
		"Fog may not be summoned here.",
		"There is no storm to be summoned from this province.",
	},
	{}}

func find_cookie(k int) *cookie_monster_tbl {
	var i int

	assert(kind(k) == T_item)

	for i = 0; cookie_monster[i].cookie != FALSE; i++ {
		if cookie_monster[i].cookie == k {
			return &cookie_monster[i]
		}
	}

	return nil
}

func may_cookie_npc(who, where, cookie int) bool {
	var t *cookie_monster_tbl
	bad_place := false

	t = find_cookie(cookie)
	assert(t != nil)

	if t.terrain > 0 && subkind(where) != schar(t.terrain) {
		bad_place = true
	}

	if t.terrain == PROV_OR_CITY &&
		subkind(where) != sub_city && loc_depth(where) != LOC_province {
		bad_place = true
	}

	if bad_place {
		if who != FALSE {
			wout(who, "%s", t.not_here)
		}
		return false
	}

	if has_item(where, cookie) == 0 {
		if who != FALSE {
			wout(who, "%s", t.no_cookies)
		}
		return false
	}

	return true
}

func do_cookie_npc(who, where, cookie, place int) int {
	var t *cookie_monster_tbl
	var p *entity_misc
	var newt int

	if !may_cookie_npc(who, where, cookie) {
		return 0
	}

	t = find_cookie(cookie)
	assert(t != nil)

	if t.kind == T_char {
		newt = new_char(t.sk, t.ni, place, 100, indep_player,
			LOY_npc, 0, "")
	} else {
		newt = new_ent(t.kind, t.sk)

		if newt > 0 {
			set_where(newt, place)
		}
	}

	if newt <= 0 {
		return 0
	}

	if t.sk == sub_ni {
		p_char(newt).health = -1
	}
	p = p_misc(newt)
	p.npc_home = where
	p.npc_cookie = cookie
	p.summoned_by = who
	p.npc_created = sysclock.turn

	if t.man_kind != FALSE {
		gen_item(newt, t.man_kind, rnd(t.low, t.high))
	}

	consume_item(where, cookie, 1)

	return newt
}

func create_peasant_mob(where int) int {
	newt := do_cookie_npc(0, where, item_mob_cookie, where)
	if newt <= 0 {
		return 0
	}

	set_name(newt, or_string(rnd(1, 2) == 1, "Mob", "Crowd"))

	queue(newt, "guard 1")
	init_load_sup(newt) /* make ready to execute commands immediately */

	return newt
}

func do_npc_orders(who, flush, queue_orders int) {
	if loyal_kind(who) == LOY_summon {
		return
	}

	if is_prisoner(who) {
		return
	}

	if queue_orders == FALSE && rp_command(who) != nil && rp_command(who).state != DONE {
		return
	} /* running an order */

	if top_order(indep_player, who) != nil && flush != FALSE {
		flush_unit_orders(player(who), who)
	}

	switch npc_program(who) {
	case 0:
		switch subkind(who) {
		case 0:
			auto_unsworn(who)
			break

		case sub_lost_soul:
			npc_move(who)
			break

		case sub_demon_lord:
			auto_undead(who)
			break

		case sub_ni:
			switch noble_item(who) {
			case item_savage:
				auto_savage(who)
				break

			case item_peasant, item_angry_peasant:
				auto_mob(who)
				break
			}
			break
		}
		break

	case PROG_bandit:
		auto_bandit(who)
		break

	case PROG_balrog:
		auto_balrog(who)
		break

	case PROG_subloc_monster:
		auto_subloc(who)
		break

	case PROG_npc_token:
		npc_move(who)
		break

	case PROG_dumb_monster:
		auto_dumb(who)
		break

	case PROG_smart_monster:
		auto_smart(who)
		break

	case PROG_orc:
		auto_orc(who)
		break

	case PROG_elf:
		auto_elf(who)
		break

	case PROG_daemon:
		auto_daemon(who)
		break

	default:
		panic("!reached")
	}

	/*
	 *  And then tell them to think again.
	 *
	 */
	if valid_box(who) {
		queue(who, "think")
	}
}

func queue_npc_orders() {
	var who int

	stage("queue_npc_orders()")

	init_savage_attacks()
	auto_hades()

	for _, who = range loop_units(indep_player) {
		do_npc_orders(who, 1, 0)
	}
}

/*
 *  Fri Feb 28 11:48:19 1997 -- Scott Turner
 *
 *  Create a new beast stack.  Type of stack is (inversely) porportional
 *  to the beast's combat values.  Size is random, exponential.
 *
 *  Tue Jun 10 11:37:10 1997 -- Scott Turner
 *
 *  Possibly restricted as to subkind -- useful in Hades.
 *
 */
func create_monster_stack(selection, total, where int) int {
	var name string
	var newt int

	if total < 1 {
		return 0
	}

	name = fmt.Sprintf("Stack of %s", plural_item_name(selection, total))
	newt = new_char(sub_ni, selection, where, -1, indep_player,
		LOY_npc, 0, name)
	p_char(newt).break_point = 0
	rp_char(newt).npc_prog = int(item_prog(selection))
	gen_item(newt, selection, total-1)
	do_npc_orders(newt, 0, 0)
	return newt
}

func create_new_beasts(where, sk int) int {
	var total, selection int

	/*
	 *  First go through all the beasts and select
	 *  which one we'd like to create.
	 *
	 *  Tue Oct 13 13:34:44 1998 -- Scott Turner
	 *
	 *  Combined this with random_beast in artifacts.c
	 */
	selection = random_beast(sk)
	assert(selection != FALSE)
	/*
	 *  Now figure out how many...
	 *
	 */
	total = rnd(1, 6)
	for total < 100 && rnd(1, 6) == 1 {
		total += rnd(1, 6)
	}
	wout(gm_player,
		"Selected %d %s for new beast stack in %s.",
		total, box_name(selection), box_name(where))
	/*
	 *  Now create it...
	 *
	 */
	newt := create_monster_stack(selection, total, where)
	/*
	 *  And (possibly) give it some treasure.
	 *
	 */
	generate_treasure(newt, 1)

	return newt
}

/*
 *  Wed Mar  5 12:16:26 1997 -- Scott Turner
 *
 *  NPC breed command.
 *
 */
func d_npc_breed(c *command) int {
	var item, num, breed_chance, season int

	if FALSE == noble_item(c.who) {
		return FALSE
	}

	item = noble_item(c.who)

	/*
	 *  Wed Nov 24 17:29:09 1999 -- Scott Turner
	 *
	 *  Most animals breed in the Spring.  A few breed at other times
	 *  of the year.  Overall, we should reject 75% of the breeding ops.
	 *
	 *     .25 * .75 = 19%
	 */
	season = oly_month(&sysclock) / 2
	if item%3 == 0 {
		season = (season + item) % 4
	}
	switch season {
	case 0: /* Spring */
		if rnd(1, 100) <= 40 {
			return TRUE
		}
		break
	case 1: /* Summer */
		if rnd(1, 100) <= 80 {
			return TRUE
		}
		break
	case 2: /* Fall */
		if rnd(1, 100) <= 99 {
			return TRUE
		}
		break
	default: /* Winter */
		if rnd(1, 100) <= 99 {
			return TRUE
		}
		break
	}

	/*
	 *  Shall we breed?
	 *
	 *  Tue Apr 27 18:26:30 1999 -- Scott Turner
	 *
	 *  Should be based on MM -- lesser creatures breed more easily than
	 *  stronger -- maximum 50%, minimum 1%
	 */
	breed_chance = ((25 * (MAX_MM - MM(item))) / MAX_MM) + 1

	/*
	 *  Always a chance.
	 *
	 */
	if breed_chance < 1 {
		breed_chance = 1
	}
	if breed_chance > 25 {
		breed_chance = 25
	}
	/*
	 *  Fri Mar 14 10:17:03 1997 -- Scott Turner
	 *
	 *  Orcs in an orc stronghold breed faster.
	 *
	 */
	if noble_item(c.who) == item_orc &&
		subkind(subloc(c.who)) == sub_orc_stronghold {
		breed_chance = 3 * breed_chance / 2 // mdhender: was *= 1.5
	}

	if rnd(1, 100) <= breed_chance {
		/*
		 *  Add 1-10% to this group.
		 *
		 *  Tue Aug 11 05:37:20 1998 -- Scott Turner
		 *
		 *  Put a cap on the number that can be added to prevent large
		 *  stacks from growing too quickly.
		 *
		 */
		num = has_item(c.who, item) / (rnd(10, 100))
		if num < 1 {
			num = 1
		}
		if num > 10 {
			num = 10
		}
		gen_item(c.who, item, num)
		wout(c.who, "Successfully bred %d.", num)
		/*
		 *  Some chance to generate new treasure.
		 *
		 */
		if num >= rnd(1, 50) {
			generate_one_treasure(c.who)
		}
	}

	/*
	 *  Shall we split up?
	 */
	if (item_split(item) != FALSE && has_item(c.who, item) > item_split(item)) &&
		rnd(1, 3) == 1 {
		/*
		 *  Tue Sep 14 12:44:27 1999 -- Scott Turner
		 *
		 *  Split the inventory between the stacks.
		 *
		 */
		var e *item_ent
		//flag := 0;
		num = has_item(c.who, item)
		newt := new_char(sub_ni, item, loc(c.who), -1, indep_player,
			LOY_npc, 0, name(c.who))
		/*
		 *
		 *  Mon May  3 09:23:57 1999 -- Scott Turner
		 *
		 *  Need to empty his inventory of similar items first.
		 *
		 */
		p_char(newt).break_point = rp_char(c.who).break_point
		rp_char(newt).npc_prog = rp_char(c.who).npc_prog
		gen_item(newt, item, num/2)
		consume_item(c.who, item, num/2)
		for _, e = range loop_inventory(c.who) {
			if FALSE == item_attack(e.item) && e.qty > 1 {
				consume_item(c.who, e.item, e.qty/2)
				gen_item(newt, e.item, e.qty/2)
			}
		}

		wout(gm_player, "Splitting into %s and %s.", box_name(c.who), box_name(newt))
	}

	return TRUE

}

/*
 *  Wed Mar  5 12:26:09 1997 -- Scott Turner
 *
 *  A "dummy" command that merely sets the NPC unit to picking new commands.
 *
 */
func v_think(c *command) int {
	var item int
	if !is_real_npc(c.who) {
		return FALSE
	}

	item = noble_item(c.who)
	if valid_box(item) &&
		has_item(c.who, item) > item_split(item)*2 {
		wout(gm_player, "Overstuffed stack %s dropping %d.",
			box_name(c.who),
			has_item(c.who, item)-item_split(item))
		drop_item(c.who, item, has_item(c.who, item)-item_split(item))

	}

	do_npc_orders(c.who, 0, 1)

	return TRUE
}
