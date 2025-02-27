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

import "math"

/*
 *  Default is that species is compatible with itself,
 *  unless an explicit {self, self, 0} is given.
 */

type breed struct {
	i1, i2, result int
}

const BREED_ACCIDENT = 5

var breed_tbl = []breed{
	{item_peasant, item_ox, item_minotaur},
	{item_peasant, item_wild_horse, item_centaur},
	{item_wild_horse, item_wild_horse, item_wild_horse},
	{item_lion, item_lizard, item_chimera},
	{item_peasant, item_lion, item_harpie},
	{item_lizard, item_bird, item_dragon},
	{item_wild_horse, item_bird, item_pegasus},
	{item_peasant, item_lizard, item_gorgon},
	{item_rat, item_spider, item_ratspider},
	{item_pegasus, item_dragon, item_nazgul},
}

func breed_translate(item int) int {

	switch item {
	case item_riding_horse:
		return item_wild_horse
	case item_warmount:
		return item_wild_horse
	}

	return item
}

func breed_match(which, i1, i2 int) bool {
	a := []int{i1, i2}
	b := []int{breed_tbl[which].i1, breed_tbl[which].i2}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if a[i] == b[j] {
				a[i], b[j] = 0, 0
			}
		}
	}

	for i := 0; i < 2; i++ {
		if a[i] != 0 || b[i] != 0 {
			return false
		}
	}

	return false
}

func d_bird_spy(c *command) int {
	targ := c.d

	if !has_holy_symbol(c.who) {
		wout(c.who, "A holy symbol is required to bird spy.")
		return FALSE
	}

	if !use_piety(c.who, skill_piety(c.use_skill)) {
		wout(c.who, "You don't have the piety required to use that prayer.")
		return FALSE
	}

	if !is_loc_or_ship(targ) {
		wout(c.who, "%s is not a location.", box_code(targ))
		return FALSE
	}

	wout(c.who, "The bird returns with a report:")
	out(c.who, "")
	show_loc(c.who, targ)

	return TRUE
}

func d_breed(c *command) int {
	i1, i2 := c.a, c.b
	breed_accident := BREED_ACCIDENT

	if is_real_npc(c.who) {
		return d_npc_breed(c)
	}

	if kind(i1) != T_item {
		wout(c.who, "%s is not an item.", c.parse[1])
		return FALSE
	}

	if kind(i2) != T_item {
		wout(c.who, "%s is not an item.", c.parse[2])
		return FALSE
	}

	if has_item(c.who, i1) < 1 {
		wout(c.who, "Don't have any %s.", box_code(i1))
		return FALSE
	}

	if has_item(c.who, i2) < 1 {
		wout(c.who, "Don't have any %s.", box_code(i2))
		return FALSE
	}

	if i1 == i2 && has_item(c.who, i1) < 2 {
		wout(c.who, "Don't have two %s.", box_code(i1))
		return FALSE
	}

	/*
	 *  A normal union just succeeds.
	 *
	 */
	if normal_union(i1, i2) {
		offspring := find_breed(i1, i2)
		wout(c.who, "Produced %s.", box_name_qty(offspring, 1))
		gen_item(c.who, offspring, 1)
		add_skill_experience(c.who, sk_breed_beasts)
		p_skill(sk_breed_beasts).use_count++
		return TRUE
	}

	/*
	 *  A non-normal union is more problematic.
	 *
	 */
	if !has_holy_symbol(c.who) {
		wout(c.who, "A holy symbol is required for that breeding.")
		return FALSE
	}

	/*
	 *  Wed Feb 23 12:01:17 2000 -- Scott Turner
	 *
	 *  Have to directly encode the piety required here so it won't
	 *  be charged automatically in use.c
	 *
	 */
	if !use_piety(c.who, 3) {
		wout(c.who, "You don't have the piety required to use that prayer.")
		return FALSE
	}

	p_skill(sk_breed_beasts).use_count++

	/*
	 *  The following isn't quite right -- there is no chance of
	 *  killing both the breeders if they are of the same type.
	 */
	killed, offspring := false, find_breed(i1, i2)
	if offspring == item_dragon {
		breed_accident = 13
	}
	if i1 == i2 { // why?
		breed_accident *= 2
	}
	if i1 != 0 && rnd(1, 100) <= breed_accident {
		wout(c.who, "%s was killed in the breeding attempt.", cap_(box_name_qty(i1, 1)))
		consume_item(c.who, i1, 1)
		killed = true
	}
	if i2 != 0 && rnd(1, 100) <= breed_accident && i1 != i2 {
		wout(c.who, "%s was killed in the breeding attempt.", cap_(box_name_qty(i2, 1)))
		consume_item(c.who, i2, 1)
		killed = true
	}
	if killed || offspring == 0 || rnd(1, 4) == 1 {
		wout(c.who, "No offspring was produced.")
		return FALSE
	}

	wout(c.who, "Produced %s.", box_name_qty(offspring, 1))

	gen_item(c.who, offspring, 1)
	add_skill_experience(c.who, sk_breed_beasts)

	return TRUE
}

func d_capture_beasts(c *command) int {
	target := c.a

	if !has_holy_symbol(c.who) {
		wout(c.who, "You must have a holy symbol to capture wild beasts.")
		return FALSE
	}

	if !has_holy_plant(c.who) {
		wout(c.who, "Capturing wild beasts requires a holy plant.")
		return FALSE
	}

	/*
	 *  Target should be a character (oddly enough)
	 *
	 */
	if kind(target) != T_char {
		wout(c.who, "You cannot capture beasts from %s.", box_name(target))
		return FALSE
	}

	/*
	 *  In same location.
	 *
	 */

	if !check_char_here(c.who, target) {
		wout(c.who, "%s is not here.", box_name(c.a))
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Cannot capture beasts from prisoners.")
		return FALSE
	}

	if c.who == target {
		wout(c.who, "Can't capture beasts from oneself.")
		return FALSE
	}

	if stack_leader(c.who) == stack_leader(target) {
		wout(c.who, "Can't capture beasts from a member of the same stack.")
		return FALSE
	}

	/*
	 *  Now select a random beast from the target.  That can be either the
	 *  target itself, or one of the beasts in the target's inventory.  In
	 *  any case, the beast's statistics determine the piety cost of the
	 *  spell.
	 *
	 *  Note that who the beast is "from" is also returned.
	 *
	 */
	type_, from := get_random_beast(target)

	if type_ == 0 || item_animal(type_) == 0 {
		wout(c.who, "%s has no beasts that you can capture.", box_name(target))
		return FALSE
	}

	piety := int(math.Ceil(float64(item_attack(type_)+item_defense(type_))/50.0 + 1.5))

	/*
	 *  Perhaps he hasn't the piety.
	 *
	 */
	if !has_piety(c.who, piety) {
		wout(c.who, "You lack the piety to lure a beast from %s.", box_name(target))
		return FALSE
	}

	/*
	 *  Lure the beast away.
	 *
	 */
	if !move_item(from, c.who, type_, 1) {
		/*
		 *  Possibly it is "from" himself who we are capturing.
		 *
		 */
		if subkind(from) == sub_ni && item_animal(noble_item(from)) == 0 {
			wout(c.who, "You capture %s!", box_name(target))
			use_piety(c.who, piety)
			take_prisoner(c.who, from)
			/*
			 *  Fri Mar 24 06:56:46 2000 -- Scott Turner
			 *
			 *  Take prisoner doesn't actually transfer the animal into your
			 *  inventory, so we'll have to "create" one.
			 *
			 */
			gen_item(c.who, type_, 1)
			return TRUE
		}
		wout(c.who, "For some reason, you capture no beast.")
		return FALSE
	}
	wout(c.who, "You capture a %s from %s!", box_name(type_), box_name(target))
	use_piety(c.who, piety)
	return TRUE

}

func find_breed(i1, i2 int) int {
	i1 = breed_translate(i1)
	i2 = breed_translate(i2)

	for i := 0; i < len(breed_tbl) && breed_tbl[i].i1 != 0; i++ {
		if breed_match(i, i1, i2) {
			return breed_tbl[i].result
		}
	}

	if item_animal(i1) != 0 && i1 == i2 {
		return i1
	}

	return 0
}

/*
 *  Fri Oct 11 13:29:43 1996 -- Scott Turner
 *
 *  Select a random beast (type) to steal from the target.
 *
 */
func get_random_beast(target int) (choice, from int) {
	/*
	 *  Select (in one pass!) a random beast.
	 *
	 */
	sum := 0
	for _, i := range loop_stack(target) {
		if subkind(i) == sub_ni && item_animal(noble_item(i)) != 0 {
			sum++
			if rnd(1, sum) == 1 {
				choice = noble_item(i)
				from = i
			}
		}
		for _, e := range inventory_loop(i) {
			if item_animal(e.item) != 0 {
				sum += e.qty
				if rnd(1, sum) <= e.qty {
					choice = e.item
					from = i
				}
			}
		}
	}

	return choice, from
}

func normal_union(b1, b2 int) bool {
	/* Thu Oct 10 12:14:35 1996 -- Scott Turner     */
	/* The following animals are considered normal: */

	/* 	[52] riding horse */
	/* 	[53] warmount     */
	/* 	[76] ox           */

	return ((b1 == 52 || b1 == 53 || b1 == 76) && b1 == b2)
}

func v_bird_spy(c *command) int {
	targ := c.a
	where := subloc(c.who)
	var v *exit_view

	if !has_holy_symbol(c.who) {
		wout(c.who, "A holy symbol is required to bird spy.")
		return FALSE
	}

	if !has_piety(c.who, skill_piety(c.use_skill)) {
		wout(c.who, "You don't have the piety required to use that prayer.")
		return FALSE
	}

	if is_ship(where) {
		where = loc(where)
	}

	if numargs(c) < 1 {
		wout(c.who, "Specify what location the bird should spy on.")
		return FALSE
	}

	if !is_loc_or_ship(c.a) {
		v = parse_exit_dir(c, where, sout("use %d", sk_bird_spy))

		if v == nil {
			return FALSE
		}

		targ = v.destination
	}

	if province(targ) != province(c.who) {
		okay := false
		l := exits_from_loc(c.who, where)
		for i := 0; i < len(l); i++ {
			if l[i].destination == targ {
				okay = true
			}
		}
		if !okay {
			wout(c.who, "The location to be spied upon must be a sublocation in the same province or a neighboring location.")
			return FALSE
		}
	}

	c.d = targ

	return TRUE
}

/*
 *  Wed Mar  5 12:14:55 1997 -- Scott Turner
 *
 *  Added hooks for npc_breed, which is used by the NPC chars.
 *
 */
func v_breed(c *command) int {
	i1, i2 := c.a, c.b

	if is_real_npc(c.who) {
		c.wait += 7
		return TRUE
	}

	if has_skill(c.who, sk_breed_beasts) == FALSE {
		wout(c.who, "Requires %s.", box_name(sk_breed_beasts))
		return FALSE
	}
	c.use_skill = sk_breed_beasts

	if numargs(c) < 2 {
		wout(c.who, "Usage: breed <item> <item>")
		return FALSE
	}

	if kind(i1) != T_item {
		wout(c.who, "%s is not an item.", c.parse[1])
		return FALSE
	}

	if kind(i2) != T_item {
		wout(c.who, "%s is not an item.", c.parse[2])
		return FALSE
	}

	if has_item(c.who, i1) < 1 {
		wout(c.who, "Don't have any %s.", box_code(i1))
		return FALSE
	}

	if has_item(c.who, i2) < 1 {
		wout(c.who, "Don't have any %s.", box_code(i2))
		return FALSE
	}

	if i1 == i2 && has_item(c.who, i1) < 2 {
		wout(c.who, "Don't have two %s.", box_code(i1))
		return FALSE
	}

	/*
	 *  Thu Oct 10 12:15:09 1996 -- Scott Turner
	 *
	 *  May need a holy symbol & piety.
	 *
	 */
	if !normal_union(i1, i2) {
		if !has_holy_symbol(c.who) {
			wout(c.who, "A holy symbol is required for that breeding.")
			return FALSE
		}

		if !has_piety(c.who, skill_piety(c.use_skill)) {
			wout(c.who, "You don't have the piety required to use that prayer.")
			return FALSE
		}
	}

	/*
	 *  Hack to fold experience_use_speedup into this skill
	 *  if they use BREED instead of USE xxxx
	 */

	c.wait = 7
	exp := max(has_skill(c.who, sk_breed_beasts)-1, 0)
	if exp != 0 {
		c.wait--
	}

	return TRUE
}

/*
 *  Thu Oct 10 14:50:13 1996 -- Scott Turner
 *
 *  Capture Beasts is now an explicit skill
 *
 */
func v_capture_beasts(c *command) int {
	target := c.a

	if !has_holy_symbol(c.who) {
		wout(c.who, "You must have a holy symbol to capture wild beasts.")
		return FALSE
	}
	if !has_holy_plant(c.who) {
		wout(c.who, "Capturing wild beasts requires a holy plant.")
		return FALSE
	}
	/*
	 *  Target should be a character (oddly enough)
	 *
	 */
	if kind(target) != T_char {
		wout(c.who, "You cannot capture beasts from %s.", box_name(c.a))
		return FALSE
	}
	/*
	 *  In same location.
	 *
	 */
	if !check_char_here(c.who, target) {
		wout(c.who, "%s is not here.", box_name(c.a))
		return FALSE
	}
	if is_prisoner(target) {
		wout(c.who, "Cannot capture beasts from prisoners.")
		return FALSE
	}
	if c.who == target {
		wout(c.who, "Can't capture beasts from oneself.")
		return FALSE
	}
	if stack_leader(c.who) == stack_leader(target) {
		wout(c.who, "Can't capture beasts from a member of the same stack.")
		return FALSE
	}

	return TRUE
}
