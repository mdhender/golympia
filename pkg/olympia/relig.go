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
 *  Wed Oct 30 11:54:25 1996 -- Scott Turner
 *
 *  Holy Terrain -- because of Eres and Domingo, this can't be simple.
 *
 *  Wed Mar 22 12:45:19 2000 -- Scott Turner
 *
 *  Hmm, this is sometimes taken to mean a province w/ holy terrain in
 *  it e.g., for Ham, it has a city there :-(
 *
 */
func is_holy_terrain(n, where int) bool {
	relig := is_priest(n)
	ht := rp_relig_skill(relig).terrain

	if relig == 0 {
		return false
	}

	if relig == sk_domingo {
		if loc_hidden(where) {
			return true
		}
		return false
	}

	if ht == sub_ocean {
		if sea_dist(province(where)) == 1 {
			return true
		}
		return false
	}

	return ht == int(subkind(province(where))) || ht == int(subkind(loc(where))) || ht == int(subkind(where))
}

/*
 *  Fri Aug 23 12:22:07 1996 -- Scott Turner
 *
 *  Determine if the priest has a holy symbol of the right
 *  sort in his possession -- and no others!
 *
 */
func has_holy_symbol(who int) bool {
	if is_priest(who) == 0 {
		return false
	}

	for _, e := range loop_inventory(who) {
		if subkind(e.item) == sub_holy_symbol {
			if rp_item_magic(e.item).religion == is_priest(who) {
				return true
			}
		}
	}

	return false
}

/*
 *  Fri Aug 23 14:37:49 1996 -- Scott Turner
 *
 *  Has a holy plant?
 *
 */
func has_holy_plant(who int) bool {
	if is_priest(who) == 0 {
		return false
	}

	assert(rp_relig_skill(is_priest(who)) != nil)

	for _, e := range loop_inventory(who) {
		if e.item == holy_plant(who) {
			return true
		}
	}

	return false
}

/*
 *  Is_Member
 *  Thu Sep  5 12:07:58 1996 -- Scott Turner
 *
 *  Determine if something is a member of a given religion.
 *
 */
func is_member(who, relig int) bool {
	if rp_char(who) == nil {
		return false
	}

	if is_priest(who) == relig {
		return true
	}

	if rp_char(who).religion.priest != 0 &&
		rp_char(rp_char(who).religion.priest) != nil &&
		is_priest(rp_char(who).religion.priest) == relig {
		return true
	}

	return false
}

/*
 *  Has_Piety
 *  Thu Sep  5 12:23:25 1996 -- Scott Turner
 *
 *  Check to see if a priest has a certain amount of piety.
 *
 */
func has_piety(who, amount int) bool {
	if is_priest(who) == 0 {
		return false
	}
	return (rp_char(who).religion.piety >= amount)
}

/*
 *  Add_Piety
 *  Mon Oct  2 09:18:10 2000 -- Scott Turner
 *
 *  Add piety.
 *
 */
func add_piety(who, amount int, use_limit bool) {
	religion := is_priest(who)
	if religion == 0 {
		return
	}

	/*
	 *  Should we limit piety?
	 *
	 */
	if use_limit &&
		options.piety_limit != 0 &&
		options.head_priest_piety_limit != 0 {
		/*
		 *  Is he the high priest?
		 *
		 */
		var limit int
		if rp_relig_skill(religion).high_priest == who {
			limit = options.head_priest_piety_limit * len(rp_char(who).religion.followers)
		} else {
			limit = options.piety_limit * len(rp_char(who).religion.followers)
		}

		if rp_char(who).religion.piety >= limit {
			return
		}
	}

	rp_char(who).religion.piety += amount
}

/*
 *  Use_Piety
 *  Thu Sep  5 12:23:25 1996 -- Scott Turner
 *
 *  Use up some piety.
 *
 */
func use_piety(who, amount int) bool {
	if FALSE == is_priest(who) {
		return false
	}
	if rp_char(who).religion.piety < amount {
		return false
	}
	rp_char(who).religion.piety -= amount
	return true
}

func check_vision_target(c *command, target int) bool {
	switch kind(target) {
	case T_char, T_ship, T_loc:
		return true
	case T_item:
		if FALSE == item_unique(target) {
			wout(c.who, "%s is not a unique item.", box_code(target))
			return false
		}
		return true
	}
	wout(c.who, "Cannot receive a vision for %s.", box_code(target))
	return false
}

func v_reveal_vision(c *command) int {
	target := c.a

	p := rp_magic(c.who)
	if p != nil && test_bit(p.visions, target) {
		wout(c.who, "Already have received a vision of %s.", box_code(target))
		wout(c.who, "A vision may only be received once for a particular target.")
		return FALSE
	}

	return TRUE
}

func d_reveal_vision(c *command) int {
	target := c.a
	if !check_vision_target(c, target) {
		return FALSE
	}

	if has_artifact(target, ART_PROT_SKILL, sk_reveal_vision, 0, 0) != FALSE {
		wout(c.who, "A grey mist blocks your vision.")
		return FALSE
	}

	p := p_magic(c.who)
	p.visions = set_bit(p.visions, target)

	wout(c.who, "%s receives a vision of %s:", box_name(c.who), box_name(target))

	out(c.who, "")

	switch kind(target) {
	case T_loc, T_ship:
		show_loc(c.who, viewloc(target))
		alert_scry_attempt(c.who, target, " successfully")
		return TRUE
	case T_char:
		char_rep_sup(c.who, target)
		return TRUE
	case T_item:
		show_item_where(c.who, target)
		return TRUE
	}
	panic("!reached")
}

func v_resurrect(c *command) int {
	body := c.a

	if is_temple(subloc(c.who)) != is_priest(c.who) {
		wout(c.who, "You must be in a temple to perform a resurrection.")
		return FALSE
	}

	/*
	 *  Sun Jun  1 10:31:40 1997 -- Scott Turner
	 *
	 *  Might have a lost soul as a prisoner...
	 *
	 */
	if !valid_box(body) {
		wout(c.who, "Don't have %s.", box_code(body))
		return FALSE
	}

	if kind(body) != T_item {
		wout(c.who, "Don't have %s.", box_code(body))
	}

	if FALSE == is_priest(body) && rp_char(body).religion.priest != c.who {
		wout(c.who, "You may only resurrect your own followers.")
		return FALSE
	}

	if is_priest(body) != FALSE && is_priest(body) != is_priest(c.who) {
		wout(c.who, "You may only resurrect priests of your own religion.")
		return FALSE
	}

	if has_item(c.who, body) != FALSE {
		if !has_piety(c.who, 10) {
			wout(c.who, "You don't have the piety required to resurrect.")
		}
	} else {
		if !has_piety(c.who, 50) {
			wout(c.who, "You don't have the piety required to resurrect.")
		}
	}

	/*	assert(item_unique(body)); */

	return TRUE
}

func d_resurrect(c *command) int {
	body := c.a

	if is_temple(subloc(c.who)) != is_priest(c.who) {
		wout(c.who, "You must be in a temple to perform a resurrection.")
		return FALSE
	}

	/*
	 *  Sun Jun  1 10:31:40 1997 -- Scott Turner
	 *
	 *  Might have a lost soul as a prisoner...
	 *
	 */
	if !valid_box(body) {
		wout(c.who, "Don't have %s.", box_code(body))
		return FALSE
	}

	if kind(body) != T_item {
		wout(c.who, "Don't have %s.", box_code(body))
	}

	if FALSE == is_priest(body) && rp_char(body).religion.priest != c.who {
		wout(c.who, "You may only resurrect your own followers.")
		return FALSE
	}

	if is_priest(body) != FALSE && is_priest(body) != is_priest(c.who) {
		wout(c.who, "You may only resurrect priests of your own religion.")
		return FALSE
	}

	if has_item(c.who, body) != FALSE {
		if !use_piety(c.who, 10) {
			wout(c.who, "You don't have the piety required to resurrect.")
			return FALSE
		}
	} else {
		if !use_piety(c.who, 50) {
			wout(c.who, "You don't have the piety required to resurrect.")
			return FALSE
		}
	}

	if rp_misc(body) != nil {
		wout(c.who, "Brought %s back to life!", rp_misc(body).save_name)
	}

	restore_dead_body(c.who, body)

	return TRUE
}

func v_prep_ritual(c *command) int {
	if is_temple(subloc(c.who)) != is_priest(c.who) {
		wout(c.who, "You must be in a temple to pray.")
		return FALSE
	}

	return TRUE
}

func d_prep_ritual(c *command) int {
	chance := hinder_med_chance(c.who)

	p := p_magic(c.who)
	p.hinder_meditation = 0

	if rnd(1, 100) <= chance {
		wout(c.who, "Disturbing images and unquiet thoughts ruin your concentration, and your prayer fails.")
		return FALSE
	}

	wout(c.who, "Your devotions have been noted by %s.", rp_relig_skill(is_priest(c.who)).name)

	add_piety(c.who, 2, false)

	return TRUE
}

func v_last_rites(c *command) int {
	body := c.a

	if is_temple(subloc(c.who)) != is_priest(c.who) {
		wout(c.who, "You must be in a temple to perform last rites.")
		return FALSE
	}

	/*
	 *  Sun Jun  1 10:31:40 1997 -- Scott Turner
	 *
	 *  Might have a lost soul as a prisoner...
	 *
	 */
	if !valid_box(body) {
		wout(c.who, "Don't have %s.", box_code(body))
		return FALSE
	}

	follower := FALSE == is_priest(body) && rp_char(body).religion.priest == c.who

	if kind(body) == T_item {
		if subkind(body) != sub_dead_body {
			wout(c.who, "%s is not a dead body.", box_code(body))
			return FALSE
		}
		if FALSE == has_item(c.who, body) && !follower {
			wout(c.who, "You do not possess that body.")
			return FALSE
		}
	} else if kind(body) == T_char {
		if subkind(body) != sub_lost_soul {
			wout(c.who, "%s is not a lost soul.", box_code(body))
			return FALSE
		}
		if FALSE == has_prisoner(c.who, body) && !follower {
			wout(c.who, "You have not captured that lost soul.")
			return FALSE
		}
	}

	return TRUE
}

func d_last_rites(c *command) int {
	body := c.a
	var owner int

	if is_temple(subloc(c.who)) != is_priest(c.who) {
		wout(c.who, "You must be in a temple to perform last rites.")
		return FALSE
	}

	/*
	 *  Sun Jun  1 10:31:40 1997 -- Scott Turner
	 *
	 *  Might have a lost soul as a prisoner...
	 *
	 */
	if !valid_box(body) {
		wout(c.who, "Don't have %s.", box_code(body))
		return FALSE
	}

	follower := FALSE == is_priest(body) && rp_char(body).religion.priest == c.who

	if kind(body) == T_item {
		if subkind(body) != sub_dead_body {
			wout(c.who, "%s is not a dead body.", box_code(body))
			return FALSE
		}
		if FALSE == has_item(c.who, body) && !follower {
			wout(c.who, "You do not possess that body.")
			return FALSE
		}
	} else if kind(body) == T_char {
		if subkind(body) != sub_lost_soul {
			wout(c.who, "%s is not a lost soul.", box_code(body))
			return FALSE
		}
		if FALSE == has_prisoner(c.who, body) && !follower {
			wout(c.who, "You have not captured that lost soul.")
			return FALSE
		}
	}

	if kind(body) == T_item {
		owner = item_unique(body)
		assert(owner != 0)
	} else {
		owner = stack_leader(body)
	}

	if owner != 0 && kind(owner) == T_char {
		wout(owner, "%s vanished.", box_name(body))
	}

	old_name := p_misc(body).save_name
	if old_name == "" {
		old_name = box_code(body)
	}

	wout(c.who, "Last rites performed for %s.", old_name)

	/*
	 *  Wed Apr 30 11:04:08 1997 -- Scott Turner
	 *
	 *  Piety bonus for LRing magicians.
	 *
	 */
	if options.mp_antipathy != FALSE {
		if is_magician(body) {
			assert(is_priest(c.who) != FALSE)
			wout(c.who, "%s rewards you with extra piety for destroying a magician.", god_name(is_priest(c.who)))
			add_piety(c.who, 10, true)
		}
	}

	kill_char(body, 0, S_nothing)
	return TRUE
}

func v_remove_bless(c *command) int {
	target := c.a
	if target == 0 {
		target = c.who
		c.a = c.who
	}

	if !check_char_here(c.who, target) {
		return FALSE
	}

	return TRUE
}

func d_remove_bless(c *command) int {
	var has int
	target := c.a

	if !check_still_here(c.who, target) {
		return FALSE
	}

	has = has_item(target, item_blessed_soldier)
	if has < 1 {
		wout(c.who, "%s has no %s.", box_name(target), just_name_qty(item_blessed_soldier, 2))
		return FALSE
	}

	consume_item(target, item_blessed_soldier, has)
	gen_item(target, item_soldier, has)

	wout(c.who, "Removed blessing from %s.", just_name_qty(item_soldier, 2))

	if target != c.who {
		wout(target, "%s removed the blessing from %s of our soldiers!", box_name(c.who), comma_num(has))
	}

	return TRUE
}

/*
 *  Thu Sep  5 12:01:19 1996 -- Scott Turner
 *
 *  Healing is now a priestly duty.
 *
 */
func v_heal(c *command) int {
	target := c.a

	if !check_char_here(c.who, target) {
		return FALSE
	}

	if FALSE == char_sick(target) {
		wout(c.who, "%s is not sick.", box_name(target))
		return FALSE
	}

	if !is_member(target, is_priest(c.who)) {
		wout(c.who, "You may only heal members of your own religion.")
		return FALSE
	}

	return TRUE
}

func d_heal(c *command) int {
	target := c.a

	if !check_char_here(c.who, target) {
		return FALSE
	}

	if FALSE == char_sick(target) {
		wout(c.who, "%s is not sick.", box_name(target))
		return FALSE
	}

	if !is_member(target, is_priest(c.who)) {
		wout(c.who, "You may only heal members of your own religion.")
		return FALSE
	}

	if kind(target) != T_char {
		wout(c.who, "%s is no longer a character.", box_code(target))
		return FALSE
	}

	vector_clear()
	vector_add(c.who)
	vector_add(target)

	wout(VECT, "%s heals %s.", box_name(c.who), box_name(target))

	p_char(target).sick = FALSE

	wout(VECT, "%s has been cured, and should now recover.",
		box_name(target))

	return TRUE
}

/*
 *  Bless Follower
 *  Fri Sep  6 13:02:46 1996 -- Scott Turner
 *
 *  Give a blessing (+10 defense) to a follower.
 *
 */
func v_bless_follower(c *command) int {
	target := c.a

	if !check_char_here(c.who, target) {
		wout(c.who, "You must be with the noble you bless.")
		return FALSE
	}

	if !is_member(target, is_priest(c.who)) {
		wout(c.who, "You may only bless members of your own religion.")
		return FALSE
	}

	if kind(target) != T_char {
		wout(c.who, "%s is no longer a character.", box_code(target))
		return FALSE
	}

	/*
	 *  Don't bless mus or undead.
	 *
	 */
	if contains_mu_undead(target) != FALSE {
		wout(c.who, "The gods will not bless magicians or undead.")
		return FALSE
	}

	/*
	 *  Maybe he already has a blessing.
	 *
	 */
	if get_effect(target, ef_defense, 0, 0) != FALSE {
		wout(c.who, "%s has already been blessed.", box_name(target))
		return FALSE
	}

	/*
	 *  Give him the blessing...
	 *
	 */
	if FALSE == add_effect(target, ef_defense, 0, 30, 10) {
		wout(c.who, "For some reason, the blessing fails to take effect.")
		return FALSE
	}
	wout(c.who, "%s is now blessed for 30 days.", box_name(target))
	wout(target, "You have been blessed by %s for 30 days.", box_name(c.who))

	return TRUE
}

/*
 *  Mon Sep 16 15:32:26 1996 -- Scott Turner
 *
 *  Front end to proselytise.
 *
 */
func v_proselytise(c *command) int {
	ret := oly_parse_s(c, sout("collect %d %d %d", item_postulant, c.a, c.b))
	assert(ret)

	return v_collect(c)
}

/*
 *  Tue Sep 17 11:48:41 1996 -- Scott Turner
 *
 *  Holy symbols.
 *
 */
func v_create_holy_symbol(c *command) int {
	if is_temple(subloc(c.who)) != is_priest(c.who) {
		wout(c.who, "You must be in a temple to create a holy symbol.")
		return FALSE
	}

	return TRUE
}

func d_create_holy_symbol(c *command) int {
	if is_temple(subloc(c.who)) != is_priest(c.who) {
		wout(c.who, "You must be in a temple to create a holy symbol.")
		return FALSE
	}

	newt := create_unique_item(c.who, sub_holy_symbol)
	if newt < 0 {
		wout(c.who, "Unable to create holy symobl.")
		return FALSE
	}
	new_name := fmt.Sprintf("Holy symbol of %s", rp_relig_skill(is_priest(c.who)).name)
	set_name(newt, new_name)

	p := p_item(newt)
	p.weight = rnd(1, 3)

	pm := p_item_magic(newt)
	pm.creator = c.who
	pm.religion = is_priest(c.who)

	wout(c.who, "%s created.", box_name(newt))

	return TRUE
}

/*
 * Priest_In_Stack
 * Fri Sep 27 11:44:46 1996 -- Scott Turner
 *
 * Check to see if there's a priest of a particular sort in a stack.
 *
 */
func priest_in_stack(who, which int) bool {
	/*
	 *  In case we get passed a 0..
	 *
	 */
	if !valid_box(who) {
		return false
	}

	for _, i := range loop_stack(who) {
		if is_priest(i) == which {
			return true
		}
	}

	return false
}

/*
 *  DOL
 *
 */
// check_loc_for_beasts returns the number of beasts found
func check_loc_for_beasts(where, who int) (count int) {
	for _, i := range loop_here(where) {
		if is_npc(i) && beast_capturable(i) {
			count++
			wout(who, "  %s in %s.", liner_desc(i), box_name(where))
		}
		if is_loc_or_ship(i) {
			count += check_loc_for_beasts(i, who)
		}
	}

	return count
}

/*
 * Detect_Beasts
 * Thu Sep 26 11:57:00 1996 -- Scott Turner
 *
 * List all the wild beasts in this and all adjacent provinces.
 *
 */
func d_detect_beasts(c *command) int {
	/*
	 *  Detect beasts and print it out...
	 *
	 */
	where := province(c.who)
	wout(c.who, "Your prayer reveals:")

	found := check_loc_for_beasts(where, c.who)

	for _, j := range loop_prov_dest(where) {
		found += check_loc_for_beasts(j, c.who)
	}

	if found == 0 {
		wout(c.who, "  No beasts nearby.")
	}

	return TRUE
}

func v_detect_beasts(c *command) int {
	return TRUE
}

/*
 *  ANTEUS
 *
 */

/*
 *  Find_Mountain_Trail
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *  Hang an effect on the priest so that the next movement is reduced.
 *
 */
func v_find_mountain_trail(c *command) int {
	/*
	 *  Hang the effect.
	 *
	 */
	if FALSE == add_effect(c.who, ef_fast_move, 0, 2, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Anteus reveals to %s all the trails into adjacent mountains.",
		box_name(c.who))

	return TRUE
}

/*
 *  Obscure_Mountain_Trail
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *  Hang an effect on the province so that the next movement is increased.
 *
 */
func d_obscure_mountain_trail(c *command) int {
	if subkind(province(c.who)) != sub_mountain {
		wout(c.who, "You must be in a mountain province to obscure its trails.")
		return FALSE
	}

	// hang the effect.
	if FALSE == add_effect(province(c.who), ef_slow_move, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "The power of Anteus obscures all trails in %s.", box_name(province(c.who)))

	return TRUE
}

func v_obscure_mountain_trail(c *command) int {
	if subkind(province(c.who)) != sub_mountain {
		wout(c.who, "You must be in a mountain province to obscure its trails.")
		return FALSE
	}

	return TRUE
}

/*
 *  Improve_Mining
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *  Hang an effect on the mine so that production is increased.
 *
 */
func d_improve_mining(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_mine_shaft {
		wout(c.who, "The target of this prayer should be a mine.")
		return FALSE
	}

	/*
	 *  Hang the effect.
	 *
	 */
	if FALSE == add_effect(where, ef_improve_mine, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Production in %s improved for one month.", box_name(where))

	return TRUE
}

func v_improve_mining(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_mine_shaft {
		wout(c.who, "The target of this prayer should be a mine shaft.")
		return FALSE
	}

	return TRUE
}

/*
 *  Conceal_Mine
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *  Conceal a mine.
 *
 */
func d_conceal_mine(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_mine_shaft {
		wout(c.who, "The target of this prayer should be a mine shaft.")
		return FALSE
	}

	/*
	 *  Conceal the mine.
	 *
	 */
	if rp_loc(where) != nil {
		/*
		 *  Use up 1 holy plant.
		 *
		 */
		rp_loc(where).hidden = TRUE
		wout(c.who, "%s is now concealed.", box_name(where))
		set_known(c.who, where)
		return TRUE
	} else {
		wout(c.who, "For some reason, %s cannot be concealed.", box_name(c.a))
		return FALSE
	}

}

func v_conceal_mine(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_mine_shaft {
		wout(c.who, "The target of this prayer should be a mine shaft.")
		return FALSE
	}

	return TRUE
}

/*
 *  Protect_Mine
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *
 */
func d_protect_mine(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_mine_shaft {
		wout(c.who, "The target of this prayer should be a mine shaft.")
		return FALSE
	}

	/*
	 *  Hang the effect.
	 *
	 */
	if FALSE == add_effect(where, ef_protect_mine, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "%s protected against calamities for one month.",
		box_name(where))

	return TRUE
}

func v_protect_mine(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_mine_shaft {
		wout(c.who, "The target of this prayer should be a mine shaft.")
		return FALSE
	}

	return TRUE
}

/*
 *  Bless_Fort
 *  Mon Oct  7 15:06:38 1996 -- Scott Turner
 *
 */
func d_bless_fort(c *command) int {
	where := subloc(c.who)

	if loc_depth(c.a) != LOC_build {
		wout(c.who, "The target of this prayer should be a building.")
		return FALSE
	}

	if (where != c.a) && (province(c.who) != province(c.a)) {
		wout(c.who, "You must be in the same province as a building to bless it.")
		return FALSE
	}

	/*
	 *  Hang the effect.
	 *
	 */
	if FALSE == add_effect(c.a, ef_bless_fort, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "%s blessed for one month.", box_name(c.a))

	return TRUE
}

func v_bless_fort(c *command) int {
	where := subloc(c.who)

	if loc_depth(c.a) != LOC_build {
		wout(c.who, "The target of this prayer should be a building.")
		return FALSE
	}

	if (where != c.a) && (province(c.who) != province(c.a)) {
		wout(c.who, "You must be in the same province as a building to bless it.")
		return FALSE
	}

	return TRUE
}

/*
 *  Weaken_Fort
 *  Mon Oct  7 15:06:38 1996 -- Scott Turner
 *
 */
func d_weaken_fort(c *command) int {
	where := subloc(c.who)
	aura := c.b

	if loc_depth(c.a) != LOC_build {
		wout(c.who, "The target of this prayer should be a building.")
		return FALSE
	}

	if (where != c.a) && (subloc(c.who) != subloc(c.a)) {
		wout(c.who, "You must be in the same location as a building to weaken it.")
		return FALSE
	}

	/*
	 *  Might be thwarted if the building is blessed.
	 *
	 */
	if get_effect(c.a, ef_bless_fort, 0, 0) != FALSE {
		wout(c.who, "%s is blessed by Anteus and cannot be weakened.",
			box_name(c.a))
		return FALSE
	}

	if aura < 1 {
		aura = 5
	}

	if !use_piety(c.who, aura) {
		wout(c.who, "You don't have the piety required to use that prayer.")
		return FALSE
	}

	/*
	 *  Otherwise...
	 *
	 */
	add_structure_damage(c.a, aura*3)
	wout(c.who, "%s takes %d points of damage!", box_name(c.a), aura*3)

	return TRUE
}

func v_weaken_fort(c *command) int {
	where := subloc(c.who)
	aura := c.b

	if loc_depth(c.a) != LOC_build {
		wout(c.who, "The target of this prayer should be a building.")
		return FALSE
	}

	if (where != c.a) && (province(c.who) != province(c.a)) {
		wout(c.who, "You must be in the same province as a building to weaken it.")
		return FALSE
	}

	if aura < 1 {
		aura = 5
	}

	if !has_piety(c.who, aura) {
		wout(c.who, "You don't have the piety required to use that prayer.")
		return FALSE
	}

	return TRUE
}

/*
 * Generic Trap
 * Mon Oct  7 15:27:19 1996 -- Scott Turner
 *
 *
 */
func d_generic_trap(c *command) int {
	/*
	 *  Find the proper trap in the traps[] array.
	 *
	 */
	var i int
	for i = 0; traps[i].type_ != 0 && traps[i].type_ != is_priest(c.who); i++ {
		//
	}

	if FALSE == traps[i].type_ {
		wout(c.who, "Priests of %s cannot set traps.", god_name(is_priest(c.who)))
		return FALSE
	}

	if !is_holy_terrain(c.who, province(c.who)) &&
		!is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You may only create a %s trap in your holy terrain.",
			traps[i].name)
		return FALSE
	}
	/*
	 *  Create the trap
	 *
	 */
	if is_holy_terrain(c.who, subloc(c.who)) {
		if FALSE == add_effect(subloc(c.who), ef_religion_trap, is_priest(c.who), 300, 1) {
			wout(c.who, "For some reason, the prayer fails to take effect.")
			return FALSE
		}
		wout(c.who, "A %s trap is now set in %s.", traps[i].name,
			box_name(subloc(c.who)))
	} else if is_holy_terrain(c.who, province(c.who)) {
		if FALSE == add_effect(province(c.who), ef_religion_trap, is_priest(c.who), 300, 1) {
			wout(c.who, "For some reason, the prayer fails to take effect.")
			return FALSE
		}
		wout(c.who, "A %s trap is now set in %s.", traps[i].name,
			box_name(province(c.who)))
	}
	return TRUE
}

func v_generic_trap(c *command) int {
	var i int

	/*
	 *  Find the proper trap in the traps[] array.
	 *
	 */
	for i = 0; traps[i].type_ != 0 && traps[i].type_ != is_priest(c.who); i++ {
		//
	}

	if FALSE == traps[i].type_ {
		wout(c.who, "Priests of %s cannot set traps.", god_name(is_priest(c.who)))
		return FALSE
	}

	if !is_holy_terrain(c.who, province(c.who)) &&
		!is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You may only create a %s trap in your holy terrain.",
			traps[i].name)
		return FALSE
	}

	return TRUE
}

/*
 *  TIMEID
 *
 */

/*
 *  Find_Forest_Trail
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *  Hang an effect on the priest so that the next movement is reduced.
 *
 */
func v_find_forest_trail(c *command) int {
	/*
	 *  Hang the effect.
	 *
	 */
	if FALSE == add_effect(c.who, ef_fast_move, 0, 2, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Timeid reveals to %s all the trails into adjacent forests.",
		box_name(c.who))

	return TRUE
}

/*
 *  Obscure_Forest_Trail
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *  Hang an effect on the province so that the next movement is increased.
 *
 */
func d_obscure_forest_trail(c *command) int {
	if subkind(province(c.who)) != sub_forest {
		wout(c.who, "You must be in a forest province to obscure its trails.")
		return FALSE
	}

	/*
	 *  Hang the effect.
	 *
	 */
	if FALSE == add_effect(province(c.who), ef_slow_move, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "The power of Timeid obscures all trails in %s.",
		box_name(province(c.who)))

	return TRUE
}

func v_obscure_forest_trail(c *command) int {
	if subkind(province(c.who)) != sub_forest {
		wout(c.who, "You must be in a forest province to obscure its trails.")
		return FALSE
	}

	return TRUE
}

/*
 *  Improve_Logging
 *  Fri Oct  4 11:35:21 1996 -- Scott Turner
 *
 *  Hang an effect on the province so that production is increased.
 *
 */
func d_improve_logging(c *command) int {
	where := subloc(c.who)
	target := c.a | where // todo: really?

	if subkind(target) != sub_forest && subkind(target) != sub_yew_grove {
		wout(c.who, "You may only bless a forest or a yew grove with this prayer.")
		return FALSE
	}

	if province(target) != province(c.who) {
		wout(c.who, "You must be in the same province as the logging operation.")
		return FALSE
	}

	/*
	 *  Hang the effect(s).
	 *
	 */
	if FALSE == add_effect(target, ef_improve_production, 0, 30, item_lumber) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	if FALSE == add_effect(target, ef_improve_production, 0, 30, item_mallorn_wood) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	if FALSE == add_effect(target, ef_improve_production, 0, 30, item_yew) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Production in %s improved for one month.", box_name(c.a))

	return TRUE
}

func v_improve_logging(c *command) int {
	where := subloc(c.who)

	target := c.a | where

	if subkind(target) != sub_forest && subkind(target) != sub_yew_grove {
		wout(c.who, "You may only bless a forest or a yew grove with this prayer.")
		return FALSE
	}

	if province(target) != province(c.who) {
		wout(c.who, "You must be in the same province as the logging operation.")
		return FALSE
	}

	return TRUE
}

/*
 *  Find_Hidden_Features
 *  Wed Oct 16 12:14:25 1996 -- Scott Turner
 *
 *  Reveal to the priest (if he's in the right terrain) all hidden features.
 *
 */
func d_find_hidden_features(c *command) int {
	var l []*exit_view
	where := subloc(c.who)
	var hidden_exits, i int

	if !is_holy_terrain(c.who, where) {
		wout(c.who, "You must be in your holy terrain to reveal hidden features.")
		return FALSE
	}

	/*
	 *  Reveal all hidden features to the priest.
	 *
	 */
	l = exits_from_loc(c.who, where)
	hidden_exits = count_hidden_exits(l)
	/*
	 *  Nothing to find
	 */
	if hidden_exits <= 0 {
		wout(c.who, "There are no hidden features here.")
		return FALSE
	}
	/*
	 *  Otherwise reveal them all.
	 *
	 */
	for i = 0; i < len(l); i++ {
		if l[i].hidden != FALSE {
			find_hidden_exit(c.who, l, i)
		}
	}

	return TRUE
}

func v_find_hidden_features(c *command) int {
	if !is_holy_terrain(c.who, province(c.who)) {
		wout(c.who, "You must be in your holy terrain to reveal hidden features.")
		return FALSE
	}

	return TRUE
}

/*
 *  Improve_Fort
 *  Wed Oct 16 12:32:08 1996 -- Scott Turner
 *
 *  Add +25 to a fortification
 *
 */
func d_improve_fort(c *command) int {
	where := subloc(c.who)

	if loc_depth(c.a) != LOC_build {
		wout(c.who, "The target of this prayer should be a building.")
		return FALSE
	}

	if (where != c.a) && (province(c.who) != province(c.a)) {
		wout(c.who, "You must be in the same province as a building to improve it.")
		return FALSE
	}

	/*
	 *  Hang the effect.
	 *
	 */
	if FALSE == add_effect(c.a, ef_improve_fort, 0, 30, 25) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "%s improved for one month.", box_name(c.a))

	return TRUE
}

func v_improve_fort(c *command) int {
	where := subloc(c.who)

	if loc_depth(c.a) != LOC_build {
		wout(c.who, "The target of this prayer should be a building.")
		return FALSE
	}

	if (where != c.a) && (province(c.who) != province(c.a)) {
		wout(c.who, "You must be in the same province as a building to improve it.")
		return FALSE
	}

	return TRUE
}

/*
 *  Recruit_Elves
 *  Fri Oct 18 15:35:11 1996 -- Scott Turner
 *
 *  Entice 5-10 elves to join your stack.
 *
 */
func d_recruit_elves(c *command) int {
	if !is_holy_terrain(c.who, province(c.who)) {
		wout(c.who, "You must be in your holy terrain to recruit elves.")
		return FALSE
	}

	qty := rnd(5, 10)
	gen_item(c.who, item_elf, qty)
	wout(c.who, "%d elves join you.", qty)

	return TRUE
}

func v_recruit_elves(c *command) int {
	if !is_holy_terrain(c.who, province(c.who)) {
		wout(c.who, "You must be in your holy terrain to recruit elves.")
		return FALSE
	}

	return TRUE
}

/*
 *  KIREUS
 *  Mon Oct 21 12:04:25 1996 -- Scott Turner
 *
 */

/*
 *  Improve_Quarrying
 *  Mon Oct 21 12:05:00 1996 -- Scott Turner
 *
 *  Hang an effect on the province so that production is increased.
 *
 */
func d_improve_quarrying(c *command) int {
	target := province(c.who)

	/*
	 *  Hang the effect(s).
	 *
	 */
	if FALSE == add_effect(target, ef_improve_production, 0, 30, item_stone) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Production in %s improved for one month.", box_name(c.a))

	return TRUE
}

func v_improve_quarrying(c *command) int {
	return TRUE
}

/*
 *  Improve_Smithing
 *  Mon Oct 21 12:05:00 1996 -- Scott Turner
 *
 *  Hang an effect on the province so that production is increased.
 *
 */
func d_improve_smithing(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast improve smithing on a smith.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the smith you are blessing.")
		return FALSE
	}

	/*
	 *  Hang the effect(s).
	 *
	 */
	if FALSE == add_effect(target, ef_improve_make, 0, 30, item_longsword) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	if FALSE == add_effect(target, ef_improve_make, 0, 30, item_plate) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Smithing by %s improved for one month.", box_name(c.a))

	return TRUE
}

func v_improve_smithing(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast improve smithing on a smith.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the smith you are blessing.")
		return FALSE
	}

	return TRUE
}

/*
 *  Edge of Kireus
 *  Mon Oct 21 12:39:05 1996 -- Scott Turner
 *
 *  Priest's edged troops get a blessing.
 *
 */
func d_edge_of_kireus(c *command) int {
	/*
	 *  Add an effect to the priest.
	 *
	 */
	if FALSE == add_effect(c.who, ef_edge_of_kireus, 0, 7, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Kireus now blesses the edged weapons of %s.", box_name(c.who))

	return TRUE
}

func v_edge_of_kireus(c *command) int {
	return TRUE
}

/*
 *  Create Mithril
 *  Wed Oct 23 11:43:16 1996 -- Scott Turner
 *
 *  Create a unit of mithril from iron and gold
 *
 */
func d_create_mithril(c *command) int {
	if has_item(c.who, item_iron) < 10 {
		wout(c.who, "Creating mithril requires 10 units iron.")
		return FALSE
	}

	if !can_pay(c.who, 100) {
		wout(c.who, "Creating mithril requires 100 units gold.")
		return FALSE
	}

	/*
	 *  Do it.
	 *
	 */
	charge(c.who, 100)
	consume_item(c.who, item_iron, 10)
	gen_item(c.who, item_mithril, 1)
	wout(c.who, "One unit mithril created.")
	return TRUE

}

func v_create_mithril(c *command) int {
	if has_item(c.who, item_iron) < 10 {
		wout(c.who, "Creating mithril requires 10 units iron.")
		return FALSE
	}

	if !can_pay(c.who, 100) {
		wout(c.who, "Creating mithril requires 100 units gold.")
		return FALSE
	}

	return TRUE
}

/*
 *  HAM
 *  Thu Oct 24 14:28:07 1996 -- Scott Turner
 *
 */

/*
 *  Thu Oct 24 14:28:52 1996 -- Scott Turner
 *
 *  Oly already has something called "second_wait" which basically tacks
 *  a day's wait onto a unit if he's been involved in a combat.  To
 *  implement this prayer, I'm extending that notion to be any number
 *  of days.  This prayer just sets a unit's second_wait to 7.
 *
 */
func d_enchant_guard(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	/*
	 *  Mayhap he's another priest of Ham?
	 *
	 */
	if is_priest(target) == is_priest(c.who) {
		wout(c.who, "%s glibly talks the city guards out of detaining him.")
		return TRUE
	}

	/*
	 *  Do it.
	 *
	 */
	if rp_command(target).second_wait < 7 {
		rp_command(target).second_wait = 7
	}

	wout(c.who, "%s detained for questioning for 7 days.", box_name(target))
	wout(target, "City guards detain you for questioning for 7 days.")

	return TRUE

}

func v_enchant_guard(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	return TRUE
}

/*
 *  Thu Oct 24 15:11:08 1996 -- Scott Turner
 *
 */
func d_urchin_spy(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	if !charge(c.who, 50) {
		wout(c.who, "Can't afford 50 gold to hire an urchin.")
		return FALSE
	}

	/*
	 *  Mayhap he's another priest of Ham?
	 *
	 */
	if is_priest(target) == is_priest(c.who) {
		wout(c.who, "Somehow, %s talks the urchin out of spying on him.")
		return TRUE
	}

	/*
	 *  Do it.
	 *
	 */
	if FALSE == add_effect(target, ef_urchin_spy, c.who, 7, subloc(c.who)) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Urchin spy sent to report on %s.", box_name(target))

	return TRUE
}

func v_urchin_spy(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	if !can_pay(c.who, 50) {
		wout(c.who, "Can't afford 50 gold to hire an urchin.")
		return FALSE
	}

	return TRUE
}

func d_draw_crowds(c *command) int {

	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	/*
	 *  Do it.
	 *
	 */
	if FALSE == add_effect(subloc(c.who), ef_improve_production, 0, 30, item_peasant) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	if FALSE == add_effect(subloc(c.who), ef_improve_production, 0, 30, item_postulant) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	if FALSE == add_effect(subloc(c.who), ef_improve_production, 0, 30, item_petty_thief) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "You attract large crowds to the city.")
	show_to_garrison = true
	wout(subloc(c.who), "Unusual numbers of people crowd the city.")
	show_to_garrison = false

	return TRUE

}

func v_draw_crowds(c *command) int {
	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	return TRUE
}

/*
 *  Fri Nov 15 12:42:20 1996 -- Scott Turner
 *
 *  Priests of Ham need a way to acquire temples -- converting towers.
 *
 */

/*
 *  Thu Oct 24 16:25:07 1996 -- Scott Turner
 *
 */
func d_arrange_mugging(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	/*
	 *  Mayhap he's another priest of Ham?
	 *
	 */
	if is_priest(target) == is_priest(c.who) {
		wout(c.who, "Somehow, %s avoids the muggers.")
		return TRUE
	}

	wout(c.who, "Muggers dispatched to acost %s.", box_name(target))

	/*
	 *  -20-50% of his gold.
	 *
	 */
	loss := (rnd(20, 50) / 100.0) * has_item(target, item_gold)

	wout(target, "Muggers attack you!")

	if loss != 0 {
		wout(target, "They rob you of %d gold.", loss)
		charge(target, loss)
		show_to_garrison = true
		wout(subloc(c.who), "A mugging was reported to the city guard.")
		show_to_garrison = false
	}

	if rnd(1, 100) < 25 && char_health(target) > 25 {
		damage := rnd(5, 15)
		wout(target, "They beat you for %d points of damage!", damage)
		add_char_damage(target, damage, 0)
	}
	return TRUE
}

func v_arrange_mugging(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	return TRUE
}

/*
 *  HALON
 *
 */

func d_calm_peasants(c *command) int {
	/*
	 *  Go through everyone in this location.  If they're an angry
	 *  mob, then disperse them.  If their angry peasants in the
	 *  possession of a noble, turn them back into regular peasants.
	 *
	 */
	where := subloc(c.who)
	for _, i := range loop_here(where) {
		/*
		 *  Ignore priests of Halon altogether.
		 *
		 */
		if is_priest(i) == sk_halon {
			continue
		}
		/*
		 *  Disperse free-floating mobs, unless they're stacked
		 *  under a priest of Halon.
		 *
		 */
		if is_npc(i) && noble_item(i) == item_angry_peasant {
			if is_priest(stack_parent(i)) != sk_halon {
				wout(where, "%s is calmed by %s's speech.",
					box_name(i), box_name(c.who))
				char_reclaim(i)
			}
			continue
		}
		/*
		 *  And remove angry peasants from noble's inventory.
		 *
		 */
		for _, e := range loop_inventory(i) {
			if e.item == item_angry_peasant {
				/*
				 *  Eat the angry peasants and produce peasants.
				 *
				 */
				consume_item(i, item_angry_peasant, e.qty)
				gen_item(i, item_peasant, e.qty)
				wout(i, "%s calms your angry peasants.", box_name(c.who))
				continue
			}
		}

	}
	wout(c.who, "You calm the angry peasants in %s.", box_name(where))
	return TRUE
}

/*
 *  Sat Oct 26 19:19:01 1996 -- Scott Turner
 *
 */
func d_improve_charisma(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_member(target, is_priest(c.who)) {
		wout(c.who, "%s is not a follower of your god!", box_name(target))
		return FALSE
	}

	/*
	 *  Hang a 30 day effect on the target.
	 *
	 */
	if FALSE == add_effect(target, ef_charisma, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "%s blesses %s with great charisma!",
		god_name(is_priest(c.who)), box_name(target))
	wout(target, "%s blesses you with great charisma!",
		god_name(is_priest(c.who)))

	return TRUE
}

func v_improve_charisma(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if !is_member(target, is_priest(c.who)) {
		wout(c.who, "%s is not a follower of your god!", box_name(target))
		return FALSE
	}

	return TRUE
}

/*
 *  Mon Oct 28 10:37:59 1996 -- Scott Turner
 *
 *  Mesmerize
 *
 */
func d_mesmerize_crowd(c *command) int {
	/*
	 *  Go through everyone in this location.  If they're not
	 *  a priest or wild animals, then give them a 35% chance
	 *  of being mesmerized for 7 days.
	 *
	 */
	where := subloc(c.who)
	for _, i := range loop_here(where) {
		/*
		 *  Ignore priests of Halon altogether.
		 *
		 */
		if is_priest(i) == sk_halon {
			continue
		}
		/*
		 *  Wild animals?
		 *
		 */
		if is_npc(i) {
			continue
		}
		/*
		 *  Otherwise, a 35% chance of being mesmerized.
		 *
		 */
		if rnd(1, 100) < 35 {
			wout(i, "You find yourself captivated for 7 days by %s's mesmerizing speech!",
				box_name(c.who))
			rp_command(i).second_wait = 7
		}
	}
	wout(c.who, "You give a mesmerizing speech in %s.", box_name(where))
	rp_command(c.who).second_wait = 7
	return TRUE
}

/*
 *  Mon Oct 28 14:08:26 1996 -- Scott Turner
 *
 *  Improve Taxes
 *
 */
func d_improve_taxes(c *command) int {
	where := province(c.who)

	if FALSE == add_effect(where, ef_improve_taxes, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "Taxes in %s will be doubled this month.", box_name(where))
	return TRUE
}

/*
 *  Guard Loyalty
 *  Mon Oct 28 14:28:19 1996 -- Scott Turner
 *
 */

func d_guard_loyalty(c *command) int {
	target := c.a

	if !valid_box(target) || nil == rp_char(target) {
		wout(c.who, "You must cast this prayer on a noble.")
		return FALSE
	}

	if subloc(c.who) != subloc(target) {
		wout(c.who, "You must be in the same location as the target.")
		return FALSE
	}

	if FALSE == add_effect(target, ef_guard_loyalty, 0, 30, 1) {
		wout(c.who, "For some reason, the prayer fails to take effect.")
		return FALSE
	}

	wout(c.who, "%s is now unswervingly loyal for one month.", box_name(target))
	return TRUE
}

/*
 *  Mon Oct 28 15:08:00 1996 -- Scott Turner
 *
 */
func d_instill_fanaticism(c *command) int {
	num := c.a

	if num < 1 || num > 10 {
		num = 10
	}

	if has_item(c.who, item_soldier) < num {
		wout(c.who, "You do not have enough %s for this prayer.",
			plural_item_name(item_soldier, 2))
		return FALSE
	}

	/*
	 *  Do it.
	 *
	 */
	consume_item(c.who, item_soldier, num)
	gen_item(c.who, item_fanatic, num)
	wout(c.who, "You create %d fanatics.", num)
	return TRUE
}

func v_instill_fanaticism(c *command) int {
	num := c.a

	if num < 1 || num > 10 {
		num = 10
	}

	if has_item(c.who, item_soldier) < num {
		wout(c.who, "You do not have enough %s for this prayer.",
			plural_item_name(item_soldier, 2))
		return FALSE
	}

	return TRUE
}

/*
 *  DOMINGO
 *  Mon Oct 28 15:48:42 1996 -- Scott Turner
 *
 */

/*
 *  Find_All_Hidden_Features
 *  Mon Oct 28 15:48:59 1996 -- Scott Turner
 *
 *  Reveal to the priest all hidden features.
 *
 */
func d_find_all_hidden_features(c *command) int {
	var l []*exit_view
	where := subloc(c.who)
	var hidden_exits, i int

	/*
	 *  Reveal all hidden features to the priest.
	 *
	 */
	l = exits_from_loc(c.who, where)
	hidden_exits = count_hidden_exits(l)
	/*
	 *  Nothing to find
	 */
	if hidden_exits <= 0 {
		wout(c.who, "There are no hidden features here.")
		return FALSE
	}
	/*
	 *  Otherwise reveal them all.
	 *
	 */
	for i = 0; i < len(l); i++ {
		if l[i].hidden != FALSE {
			find_hidden_exit(c.who, l, i)
		}
	}

	return TRUE
}

/*
 *  Conceal_Location
 *  Mon Oct 28 15:50:18 1996 -- Scott Turner
 *
 *  Conceal a location
 *
 */
func d_conceal_location(c *command) int {
	where := subloc(c.who)

	if (where != c.a) && (province(c.who) != province(c.a)) {
		wout(c.who, "You must be in the same province as a location to conceal it.")
		return FALSE
	}

	if where == province(where) {
		wout(c.who, "You may only conceal a sublocation.")
		return FALSE
	}

	if has_item(province(where), item_peasant) >= 100 {
		wout(c.who, "You may only conceal a sublocation in a wilderness province.")
		return FALSE
	}

	/*
	 *  Conceal the subloc
	 *
	 */
	if rp_loc(where) != nil {
		rp_loc(where).hidden = TRUE
		wout(c.who, "%s is now concealed.", box_name(where))
		set_known(c.who, where)
		return TRUE
	} else {
		wout(c.who, "For some reason, %s cannot be concealed.", box_name(where))
		return FALSE
	}

}

func v_conceal_location(c *command) int {
	where := subloc(c.who)

	if (where != c.a) && (province(c.who) != province(c.a)) {
		wout(c.who, "You must be in the same province as a location to conceal it.")
		return FALSE
	}

	if where == province(where) {
		wout(c.who, "You may only conceal a sublocation.")
		return FALSE
	}

	if has_item(province(where), item_peasant) >= 100 {
		wout(c.who, "You may only conceal a sublocation in a wilderness province.")
		return FALSE
	}

	return TRUE
}

/*
 *  Mon Oct 28 15:56:14 1996 -- Scott Turner
 *
 */
func d_create_ninja(c *command) int {
	num := c.a

	if num < 1 || num > 10 {
		num = 10
	}

	if has_item(c.who, item_soldier) < num {
		wout(c.who, "You do not have enough %s for this prayer.",
			plural_item_name(item_soldier, 2))
		return FALSE
	}

	/*
	 *  Do it.
	 *
	 */
	consume_item(c.who, item_soldier, num)
	gen_item(c.who, item_ninja, num)
	wout(c.who, "You create %d ninjas.", num)
	return TRUE
}

func v_create_ninja(c *command) int {
	num := c.a

	if num < 1 || num > 10 {
		num = 10
	}

	if has_item(c.who, item_soldier) < num {
		wout(c.who, "You do not have enough %s for this prayer.",
			plural_item_name(item_soldier, 2))
		return FALSE
	}

	return TRUE
}

/*
 *  Tue Oct 29 08:23:01 1996 -- Scott Turner
 *
 */
func d_create_mist(c *command) int {
	/*
	 *  Do it.
	 *
	 */
	if FALSE == new_storm(0, sub_mist, 2, province(c.who)) {
		wout(c.who, "The mists fail to rise.")
	} else {
		wout(c.who, "Mists begin to seep from the ground.")
	}

	return TRUE
}

/*
 *  Wed Oct 30 11:47:14 1996 -- Scott Turner
 *
 */
func d_gather_holy_plant(c *command) int {
	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	if rnd(1, 100) < 50 {
		wout(c.who, "You find a %s and collect it.", box_name(holy_plant(c.who)))
		gen_item(c.who, holy_plant(c.who), 1)
	} else {
		wout(c.who, "You cannot find any %s.", box_name(holy_plant(c.who)))
	}

	return TRUE
}

func v_gather_holy_plant(c *command) int {
	if !is_holy_terrain(c.who, subloc(c.who)) {
		wout(c.who, "You must be in your holy terrain to use this prayer.")
		return FALSE
	}

	return TRUE
}

/*
 *  Dedicate Temple
 *  Fri Nov 15 13:44:30 1996 -- Scott Turner
 *
 *  Turn a generic temple into one for your god.
 *
 *  Mon Mar  1 09:21:23 1999 -- Scott Turner
 *
 *  Supposed to allow only one temple in a location.  Oh, taken care
 *  of in build.c
 *
 */
func v_dedicate_temple(c *command) int {
	where := subloc(c.who)
	var i int

	/*
	 *  Have to be in a temple.
	 *
	 */
	if subkind(where) != sub_temple {
		wout(c.who, "To dedicate a temple, you must be inside the temple.")
		return FALSE
	}

	/*
	 *  Have to be top dog.
	 *
	 */
	if building_owner(where) != c.who {
		wout(c.who, "Must be the first character inside a temple to dedicate.")
		return FALSE
	}

	/*
	 *  Need to be in proper terrain.
	 *
	 */
	if !is_holy_terrain(c.who, loc(where)) {
		wout(c.who, "Temples of the %s must be built in holy terrain.",
			box_name(is_priest(c.who)))
		return FALSE
	}

	/*
	 *  Can't build if there's already a temple here.
	 *
	 */
	for _, i = range loop_all_here(where) {
		if is_temple(i) == is_priest(c.who) {
			wout(c.who, "There is already a temple of your religion here.")
			return FALSE
		}
	}

	return TRUE
}

func d_dedicate_temple(c *command) int {
	where := subloc(c.who)
	p := rp_subloc(where)

	/*
	 *  Have to be in a temple.
	 *
	 */
	if subkind(where) != sub_temple {
		wout(c.who, "To dedicate a temple, you must be inside the temple.")
		return FALSE
	}

	/*
	 *  Have to be top dog.
	 *
	 */
	if building_owner(where) != c.who {
		wout(c.who, "Must be the first character inside a temple to dedicate.")
		return FALSE
	}

	/*
	 *  Need to be in proper terrain.
	 *
	 */
	if !is_holy_terrain(c.who, loc(where)) {
		wout(c.who, "Temples of the %s must be built in holy terrain.",
			box_name(is_priest(c.who)))
		return FALSE
	}

	/*
	 *  Can't build if there's already a temple here.
	 *
	 */
	for _, i := range loop_all_here(where) {
		if is_temple(i) == is_priest(c.who) {
			wout(c.who, "There is already a temple of your religion here.")
			return FALSE
		}
	}

	/*
	 *  Make it a temple.
	 *
	 */
	if p == nil {
		wout(c.who, "For some reason, your god refuses to dedicate this structure.")
		return FALSE
	} else {
		p.guild = is_priest(c.who)
		p.teaches = append(p.teaches, is_priest(c.who))
		wout(c.who, "%s now dedicated to %s.", box_name(where), god_name(is_priest(c.who)))
	}
	return TRUE
}

/*
 *  Wed Dec 30 11:39:36 1998 -- Scott Turner
 *
 *  Basic religion stuff, as necessary.
 *
 */

/*
 *  Thu Dec 31 11:46:42 1998 -- Scott Turner
 *
 *  Basic religion version of hinder meditation.  Affects any magician
 *  but can't cross oceans!
 *
 */
func v_hinder_med_b(c *command) int {
	target := c.a
	if crosses_ocean(c.who, target) {
		wout(c.who, "Something seems to block your vision.")
		return FALSE
	}

	wout(c.who, "Attempt to hinder attempts at meditation by %s.",
		box_code(target))

	return TRUE
}

func d_hinder_med_b(c *command) int {
	target := c.a
	var p *char_magic

	if crosses_ocean(c.who, target) {
		wout(c.who, "Something seems to block your vision.")
		return FALSE
	}

	wout(c.who, "Successfully cast %s on %s.",
		box_name(sk_hinder_med),
		box_code(target))

	p = p_magic(target)
	p.hinder_meditation = 3

	hinder_med_omen(target, c.who)

	return TRUE
}

func v_vision_reg(c *command) int {
	target := c.a
	if province(target) != target {
		wout(c.who, "The target of your vision must be a province.")
		return FALSE
	}

	if crosses_ocean(c.who, target) {
		wout(c.who, "Something seems to block your vision.")
		return FALSE
	}

	return TRUE
}

func d_vision_reg(c *command) int {
	target := c.a
	if province(target) != target {
		wout(c.who, "The target of your vision must be a province.")
		return FALSE
	}

	if crosses_ocean(c.who, target) {
		wout(c.who, "Something seems to block your vision.")
		return FALSE
	}

	if has_artifact(target, ART_PROT_SKILL, sk_reveal_vision, 0, 0) != FALSE {
		wout(c.who, "A grey mist blocks your vision.")
		alert_scry_attempt(c.who, target, " unsuccessfully")
		return FALSE
	}

	if loc_shroud(province(target)) >= rp_skill(c.use_skill).piety {
		wout(c.who, "%s is shrouded from your vision.",
			box_code(target))

		alert_scry_attempt(c.who, target, " unsuccessfully")

		return FALSE
	}

	/*
	 *  Check the proper sort of target based on the priest's religion.
	 *
	 *  Wed Mar 22 12:46:47 2000 -- Scott Turner
	 *
	 *  Ugh, have to fix this for Ham, Halon, Domingo & Eres.
	 *
	 */
	if !is_holy_terrain(c.who, target) {
		relig := is_priest(c.who)
		ht := rp_relig_skill(relig).terrain
		if (ht == sub_city && has_city(target) == FALSE) ||
			(ht == sub_ocean && sea_dist(target) > 1) ||
			(relig == sk_domingo && loc_contains_hidden(target) == FALSE) {
			wout(c.who, "May only vision holy terrain.")
			return FALSE
		}
	}
	wout(c.who, "You receive a vision from your god!")
	show_loc(c.who, viewloc(target))
	alert_scry_attempt(c.who, target, " unsuccessfully")
	return TRUE
}
