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

// basic magic

import "math"

func add_aura(who, aura int) {
	p_magic(who).cur_aura += aura
	if char_max_aura(who) != 0 && char_cur_aura(who) > 5*max_eff_aura(who) {
		p_magic(who).cur_aura = 5 * max_eff_aura(who)
	}
}

func v_meditate(c *command) int {
	if !char_alone(c.who) {
		wout(c.who, "You cannot meditate unless completely alone.")
		return FALSE
	}

	if char_cur_aura(c.who) >= max_eff_aura(c.who) {
		wout(c.who, "Current aura is already %d.  It may not be increased further via meditation.", char_cur_aura(c.who))
		return FALSE
	}

	wout(c.who, "Meditate for %s.", weeks(c.wait))
	return TRUE
}

func hinder_med_chance(who int) int {
	p := rp_magic(who)
	if p == nil || p.hinder_meditation < 1 {
		return 0
	}
	switch p.hinder_meditation {
	case 1:
		return 50
	case 2:
		return 75
	case 3:
		return 90
	}
	panic("!reached")
}

func d_meditate(c *command) int {
	if !char_alone(c.who) {
		wout(c.who, "You cannot meditate unless completely alone.")
		return FALSE
	}

	if char_cur_aura(c.who) >= max_eff_aura(c.who) {
		wout(c.who, "Current aura is already %d.  It may not be increased further via meditation.", char_cur_aura(c.who))
		return FALSE
	}

	if c.wait != 0 {
		return TRUE
	}

	chance := hinder_med_chance(c.who)

	p := p_magic(c.who)
	p.hinder_meditation = 0
	if rnd(1, 100) <= chance {
		wout(c.who, "Disturbing images and unquiet thoughts ruin the meditative trance.  Meditation fails.")
		return FALSE
	}

	// how much should we add?
	// 2, 4 if alone in a tower, and 2 additional if they have an auraculum.
	add_aura(c.who, 2)
	if subkind(subloc(c.who)) == sub_tower && alone_here(c.who) != FALSE {
		add_aura(c.who, 2)
	}
	if has_auraculum(c.who) != FALSE {
		add_aura(c.who, 2)
	}

	wout(c.who, "Current aura is now %d.", p.cur_aura)
	return TRUE
}

func v_adv_med(c *command) int {
	if subkind(subloc(c.who)) != sub_tower || FALSE == alone_here(c.who) {
		wout(c.who, "You must be alone in a tower to use Advanced Meditation.")
		return FALSE
	}
	wout(c.who, "Advanced meditation for %s.", weeks(c.wait))
	return TRUE
}

func d_adv_med(c *command) int {
	if subkind(subloc(c.who)) != sub_tower || FALSE == alone_here(c.who) {
		wout(c.who, "You must be alone in a tower to use Advanced Meditation.")
		return FALSE
	}

	if c.wait != 0 {
		return TRUE
	}

	chance := hinder_med_chance(c.who)
	p := p_magic(c.who)
	p.hinder_meditation = 0
	// todo: should we use the max effective aura?
	//m_a := max_eff_aura(c.who);

	if rnd(1, 100) <= chance {
		wout(c.who, "Disturbing images and unquiet thoughts hamper the meditative trance!")
		return FALSE
	}

	p.max_aura++
	wout(c.who, "Maximum aura is now %d.", p.max_aura)
	return TRUE
}

func v_hinder_med(c *command) int {
	target := c.a
	if c.b < 1 {
		c.b = 1
	} else if c.b > 3 {
		c.b = 3
	}
	aura := c.b

	if FALSE == cast_check_char_here(c.who, target) {
		return FALSE
	} else if FALSE == check_aura(c.who, aura) {
		return FALSE
	}

	wout(c.who, "Attempt to hinder attempts at meditation by %s.", box_code(target))

	return TRUE
}

func hinder_med_omen(who, other int) {
	if rnd(1, 100) < 50 {
		return
	}
	switch rnd(1, 3) {
	case 1:
		wout(who, "A disturbing image of %s appeared last night in a dream.", box_name(other))
	case 2:
		wout(who, "As a cloud drifts across the moon, it seems for an instant that it takes the shape of a ghoulish face, looking straight at you.")
	case 3:
		wout(who, "You are shocked out of your slumber in the middle of the night by cold fingers touching your neck, but when you glance about, there is no one to be seen.")
	}
	panic("!reached")
}

func d_hinder_med(c *command) int {
	target := c.a
	aura := c.b
	var p *char_magic

	if FALSE == charge_aura(c.who, aura) {
		return FALSE
	}

	wout(c.who, "Successfully cast %s on %s.",
		box_name(sk_hinder_med),
		box_name(target))

	p = p_magic(target)
	p.hinder_meditation += aura
	if p.hinder_meditation > 3 {
		p.hinder_meditation = 3
	}
	hinder_med_omen(target, c.who)

	return TRUE
}

func v_reveal_mage(c *command) int {
	target := c.a
	category := c.b
	if c.c < 1 {
		c.c = 1
	}
	aura := c.c

	if FALSE == cast_check_char_here(c.who, target) {
		return FALSE
	}

	if FALSE == check_aura(c.who, aura) {
		return FALSE
	}

	if category == 0 || !magic_skill(category) || skill_school(category) != category {
		wout(c.who, "%s is not a magical skill category.", box_code(category))
		wout(c.who, "Assuming %s.", box_name(sk_basic))

		c.b = sk_basic
		category = sk_basic
	}

	wout(c.who, "Attempt to scry the magical abilities of %s within %s.", box_name(target), box_name(category))

	return TRUE
}

func d_reveal_mage(c *command) int {
	target := c.a
	category := c.b
	aura := c.c

	if FALSE == charge_aura(c.who, aura) {
		return FALSE
	}

	assert(valid_box(category))
	assert(skill_school(category) == category && magic_skill(category))

	var source string
	has_detect := has_skill(target, sk_detect_abil)
	if has_detect > exp_novice {
		source = box_name(c.who)
	} else {
		source = "Someone"
	}

	if aura <= char_abil_shroud(target) {
		wout(c.who, "The abilities of %s are shrouded from your scry.", box_name(target))
		if has_detect != 0 {
			wout(target, "%s cast %s on us, but failed to learn anything.", source, box_name(sk_reveal_mage))
		}
		if has_detect > exp_teacher {
			wout(target, "They sought to learn what we know of %s.", box_name(category))
		}
		return FALSE
	}

	first := true
	for _, e := range loop_char_skill_known(target) {
		if skill_school(e.skill) != category ||
			e.skill == category {
			continue
		}

		if first {
			wout(c.who, "%s knows the following %s spells:", box_name(target), box_name(category))
			indent += 3
			first = false
		}

		if c.use_exp > exp_journeyman {
			list_skill_sup(c.who, e, "")
		} else {
			wout(c.who, "%s", box_name(e.skill))
		}
	}

	if first {
		wout(c.who, "%s knowns no %s spells.", box_name(target), box_name(category))
	} else {
		indent -= 3
	}

	if has_detect != 0 {
		wout(target, "%s successfully cast %s on us.",
			source, box_name(sk_reveal_mage))

		if has_detect > exp_teacher {
			wout(target, "Our knowledge of %s was revealed.",
				box_name(category))
		}
	}

	return TRUE
}

func v_view_aura(c *command) int {
	if c.a < 1 {
		c.a = 1
	}
	aura := c.a
	if FALSE == check_aura(c.who, aura) {
		return FALSE
	} else if crosses_ocean(cast_where(c.who), c.who) != FALSE {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	where := reset_cast_where(c.who)
	c.d = where

	wout(c.who, "Will scry the current aura ratings of other mages in %s.", box_name(where))

	return TRUE
}

func d_view_aura(c *command) int {
	aura := c.a
	where := c.d

	if !is_loc_or_ship(where) {
		wout(c.who, "%s is no longer a valid location.", box_code(where))
		return FALSE
	} else if FALSE == charge_aura(c.who, aura) {
		return FALSE
	}

	first := true
	for _, n := range loop_char_here(where) {
		if is_magician(n) {
			first = false

			// does the viewed magician have Detect ability scry?
			var s string
			level := char_cur_aura(n)
			learned := !(aura <= char_abil_shroud(n)) // todo: very confusing
			if learned {
				s = sout("%d", level)
			} else {
				s = "???"
			}

			wout(c.who, "%s, current aura: %s", box_name(n), s)

			var source string
			has_detect := has_skill(n, sk_detect_abil)
			if has_detect > exp_novice {
				source = box_name(c.who)
			} else {
				source = "Someone"
			}

			if has_detect != 0 {
				wout(n, "%s cast View aura here.", source)
			}
			if has_detect > exp_journeyman {
				if learned {
					wout(n, "Our current aura rating was learned.")
				} else {
					wout(n, "Our current aura rating was not revealed.")
				}
			}
		}
	}

	if first {
		wout(c.who, "No mages are seen here.")
		log_output(LOG_CODE, "d_view_aura: not a mage?\n")
	}

	return TRUE
}

func v_shroud_abil(c *command) int {
	if c.a < 1 {
		c.a = 1
	}
	// todo: this always succeeds?
	//aura := c.a;

	wout(c.who, "Attempt to create a magical shroud to conceal our abilities.")

	return TRUE
}

func d_shroud_abil(c *command) int {
	aura := c.a
	if FALSE == charge_aura(c.who, aura) {
		return FALSE
	}

	p := p_magic(c.who)
	p.ability_shroud += aura

	wout(c.who, "Now cloaked in an aura %s ability shroud.", nice_num(p.ability_shroud))

	return TRUE
}

func v_detect_abil(c *command) int {
	if FALSE == check_aura(c.who, 1) {
		return FALSE
	}

	wout(c.who, "Will practice ability scry detection.")
	return TRUE
}

func d_detect_abil(c *command) int {
	if FALSE == charge_aura(c.who, 1) {
		return FALSE
	}

	return TRUE
}

func v_dispel_abil(c *command) int {
	target := c.a

	if FALSE == cast_check_char_here(c.who, target) {
		return FALSE
	}

	if FALSE == check_aura(c.who, 3) {
		return FALSE
	}

	wout(c.who, "Attempt to dispel any ability shroud from %s.",
		box_name(target))

	return TRUE
}

func d_dispel_abil(c *command) int {
	target := c.a

	p := rp_magic(target)
	if p != nil && p.ability_shroud > 0 {
		if FALSE == charge_aura(c.who, 3) {
			return FALSE
		}

		wout(c.who, "Dispeled an aura %s ability shroud from %s.", nice_num(p.ability_shroud), box_name(target))
		p.ability_shroud = 0
		wout(target, "The magical ability shroud has dissipated.")
	} else {
		wout(c.who, "%s had no ability shroud.", box_name(target))
	}

	return TRUE
}

func v_quick_cast(c *command) int {
	if c.a < 1 {
		c.a = 1
	}
	aura := c.a
	if aura > 3 {
		wout(c.who, "You may only speed casting by 3 days.")
		aura = 3
	}
	if FALSE == check_aura(c.who, aura) {
		return FALSE
	}

	wout(c.who, "Attempt to speed next spell cast.")

	return TRUE
}

func d_quick_cast(c *command) int {
	aura := c.a
	if FALSE == charge_aura(c.who, aura) {
		return FALSE
	}

	p := p_magic(c.who)
	p.quick_cast += aura

	if p.quick_cast > 3 {
		p.quick_cast = 3
	}

	wout(c.who, "Spell cast speedup now %d.", p.quick_cast)

	return TRUE
}

func v_save_quick(c *command) int {
	if char_quick_cast(c.who) < 1 {
		wout(c.who, "No stored spell cast speedup.")
		return FALSE
	}
	if FALSE == check_aura(c.who, 3) {
		return FALSE
	}

	wout(c.who, "Attempt to save speeded cast state.")
	return TRUE
}

func d_save_quick(c *command) int {
	if char_quick_cast(c.who) < 1 {
		wout(c.who, "No stored spell cast speedup.")
		return FALSE
	}

	if FALSE == charge_aura(c.who, 3) {
		return FALSE
	}

	newPotion := new_potion(c.who)
	if newPotion < 0 {
		wout(c.who, "Spell failed.")
		return FALSE
	}

	p := p_magic(c.who)
	im := p_item_magic(newPotion)
	im.use_key = use_quick_cast
	im.quick_cast = p.quick_cast

	p.quick_cast = 0

	return TRUE
}

func v_use_quick_cast(c *command) int {
	item := c.a
	assert(kind(item) == T_item)

	wout(c.who, "%s drinks the potion...", just_name(c.who))

	im := rp_item_magic(item)
	if im == nil || im.quick_cast < 1 || !is_magician(c.who) {
		wout(c.who, "Nothing happens.")
		destroy_unique_item(c.who, item)
		return FALSE
	}

	p_magic(c.who).quick_cast += im.quick_cast

	wout(c.who, "Spell cast speedup now %d.", char_quick_cast(c.who))
	destroy_unique_item(c.who, item)

	return TRUE
}

/*
 *  Mon Feb 14 10:58:03 2000 -- Scott Turner
 *
 *  Charge mana for religious spells!
 *
 */
const (
	DAYS_PER_SCROLL = 7.0
)

func scrolls_needed(x int) int {
	return int(math.Ceil(float64(x) / DAYS_PER_SCROLL))
}

func v_write_spell(c *command) int {
	spell := c.a
	days := c.b
	book := c.c

	know := has_skill(c.who, spell)
	if know < 1 {
		wout(c.who, "%s does not know %s.", box_name(c.who),
			box_code(spell))
		return FALSE
	}

	if !magic_skill(c.use_skill) && magic_skill(spell) {
		wout(c.who, "Magical skills may not be scribed with %s.", box_name(c.use_skill))
		return FALSE
	}

	if magic_skill(c.use_skill) && skill_school(spell) != skill_school(c.use_skill) {
		wout(c.who, "%s only allows %s spells to be scribed.", box_code(c.use_skill), box_name(skill_school(c.use_skill)))
		return FALSE
	}

	if !religion_skill(c.use_skill) && religion_skill(spell) {
		wout(c.who, "Religion skills may not be scribed with %s.", box_name(c.use_skill))
		return FALSE
	}

	if religion_skill(c.use_skill) && skill_school(spell) != skill_school(c.use_skill) {
		wout(c.who, "%s only allows %s spells to be scribed.", box_code(c.use_skill), box_name(skill_school(c.use_skill)))
		return FALSE
	}

	if magic_skill(c.use_skill) && FALSE == check_aura(c.who, 2*days) {
		wout(c.who, "Recording a spell requires 2 aura for each day recorded.")
		return FALSE
	}

	if religion_skill(c.use_skill) && FALSE == check_aura(c.who, 2*days) {
		wout(c.who, "Recording a spell requires 2 piety for each day recorded.")
		return FALSE
	}

	/*
	 *  Only write category skills or teachable subskills.
	 *
	 */
	parent := skill_school(spell)
	assert(parent != 0)
	category := (skill_school(spell) == spell)
	var teachable, unteachable bool
	if !category {
		teachable = ilist_lookup(rp_skill(parent).offered, spell) != -1
	}
	if !category && !teachable {
		unteachable = true
	}
	if unteachable {
		wout(c.who, "%s is an unteachable subskill, so you cannot write it down.", box_name(spell))
		return FALSE
	}

	// now the wait is calculated from the command.
	if days == 0 {
		wout(c.who, "You should put the number of days to record as the second argument to this command.  I'll assume that you mean 7 days.")
		days = 7
	}

	// needs to have that many blank scrolls.
	if has_item(c.who, item_blank_scroll) < scrolls_needed(days) {
		wout(c.who, "You need %d scroll%s to record %d day%s.",
			scrolls_needed(days), or_string(scrolls_needed(days) > 1, "s", ""),
			days, or_string(days > 1, "s", ""))
		return FALSE
	}

	/*
	 *  Maybe he's adding on to a book?
	 *
	 *  Mon Dec 21 08:20:33 1998 -- Scott Turner
	 *
	 *  Now that this is not polled, use up all the scrolls at
	 *  the start of the command!
	 *
	 */
	// todo: this is confusing
	isAdding := book != 0 && valid_box(book) && has_item(c.who, book) != FALSE
	if isAdding {
		if p := rp_item_magic(book); p != nil {
			isAdding = ilist_lookup(p.may_study, spell) != -1
		}
	}
	if isAdding {
		consume_item(c.who, item_blank_scroll, scrolls_needed(days))
		wout(c.who, "Adding pages to %s.", box_name(book))
	} else {
		if FALSE == can_pay(c.who, 100) {
			wout(c.who, "You cannot afford to start a new book.")
			return FALSE
		}
		charge(c.who, 100)
		consume_item(c.who, item_blank_scroll, scrolls_needed(days))
		newItem := create_unique_item(c.who, sub_book)
		p := p_item_magic(newItem)
		set_name(newItem, "Study Guide")
		p.may_study = append(p.may_study, spell)
		p.orb_use_count = 0
		p_item(newItem).weight = 5
		c.c = newItem
	}

	if (magic_skill(c.use_skill) || religion_skill(c.use_skill)) && FALSE == charge_aura(c.who, 2*days) {
		wout(c.who, "Some odd warp in the space-time continuum aborts this command.")
		return FALSE
	}

	c.wait = days

	wout(c.who, "Spend %s writing %s into %s.", weeks(c.wait), box_name(spell), box_name(c.c))

	return TRUE
}

func new_scroll(who int) int {
	newScroll := create_unique_item(who, sub_scroll)
	if newScroll < 0 {
		wout(who, "Scroll creation failed.")
		return FALSE
	}
	set_name(newScroll, "Scroll")

	p := p_item_magic(newScroll)
	p.creator = who
	p.orb_use_count = 1 /* Let the scroll get one use. */
	p_item(newScroll).weight = 1

	wout(who, "Produced %s.", box_name(newScroll))

	return newScroll
}

/*
 *  Sun Nov 10 20:03:56 1996 -- Scott Turner
 *
 *  Each day write one more page into the book.  (Need to make
 *  this skill "polled".)
 *
 *  Mon Dec 21 08:25:33 1998 -- Scott Turner
 *
 *  No longer used?  Ah, no wait, let's put the instruction a day
 *  at a time.
 *
 */
func d_write_spell(c *command) int {
	spell := c.a
	// todo: ignore c.b?
	book := c.c

	var p *item_magic
	if book == 0 || valid_box(book) || has_item(c.who, book) == FALSE {
		p = rp_item_magic(book)
		if p == nil || ilist_lookup(p.may_study, spell) == -1 {
			wout(c.who, "You seem to have lost your book.")
			return FALSE
		}
	}

	p.orb_use_count++
	return TRUE

}

func v_mage_menial(c *command) int {
	where := province(c.who)

	if has_item(where, item_peasant) < 100 {
		wout(c.who, "Not enough peasantry here to use this skill.")
		return FALSE
	}

	/*
	 *  maybe a temple?
	 *
	 */
	for _, j := range loop_all_here(where) {
		if subkind(j) == sub_temple {
			wout(c.who, "No peasants will be seen with you in this province.")
			return FALSE
		}
	}

	/*
	 *  Mon Dec 14 11:24:43 1998 -- Scott Turner
	 *
	 *  Alternate second argument is days to spend.
	 *
	 */
	c.wait = 7

	/*
	 *  Running count of gold earned... h should be unused
	 *
	 *  Wed Sep  9 18:45:44 1998 -- Scott Turner
	 *
	 *  Whoops, h is used for "basis"!
	 *
	 */
	c.g = 0

	return TRUE

}

func mage_menial_how() string {
	switch rnd(1, 9) {
	case 1:
		return " curing runny noses"
	case 2:
		return " dowsing for water"
	case 3:
		return " selling love potions"
	case 4:
		return " selling good luck charms"
	case 5:
		return " predicting the future"
	case 6:
		return " reading palms"
	case 7, 8, 9:
		return ""
	}
	panic("!reached")
}

func d_mage_menial(c *command) int {
	where := province(c.who)
	if has_item(where, item_peasant) < 100 {
		wout(c.who, "Not enough peasantry here to use this skill.")
		wout(c.who, "Earned a total of %s gold.", nice_num(c.g))
		return FALSE
	}

	/*
	 *  maybe a temple?
	 *
	 */
	for _, j := range loop_all_here(where) {
		if subkind(j) == sub_temple {
			wout(c.who, "No peasants will be seen with you in this province.")
			wout(c.who, "Earned a total of %s gold.", nice_num(c.g))
			return FALSE
		}
	}

	/*
	 *  Maybe no money?
	 *
	 */
	if has_item(where, item_gold) < 1 {
		wout(c.who, "The peasants here have no more gold to spend.")
		return FALSE
	}

	/*
	 *  Earn gold for today, if possible!
	 *
	 */
	amount := 0
	switch c.use_exp {
	case exp_novice:
		amount = 5
	case exp_journeyman:
		amount = 8
	case exp_teacher:
		amount = 10
	case exp_master:
		amount = 12
	case exp_grand:
		amount = 14
	}
	if amount > has_item(where, item_gold) {
		amount = has_item(where, item_gold)
	}

	assert(move_item(where, c.who, item_gold, amount) != 0)
	c.g += amount

	if c.wait == FALSE {
		gold_common_magic += c.g
		wout(c.who, "Earned %s%s.", gold_s(c.g), mage_menial_how())
		show_to_garrison = true
		wout(where, "%s earned %s working at common magic.", box_name(c.who), gold_s(c.g))
		if subloc(c.who) != where {
			wout(subloc(c.who), "%s earned %s working at common magic.", box_name(c.who), gold_s(c.g))
		}
		show_to_garrison = false
	}

	/*
	 *  Fri Sep 11 09:03:02 1998 -- Scott Turner
	 *
	 *  Add some peasants.
	 *
	 */
	gen_item(where, item_peasant, rnd(1, 5))

	return TRUE
}

func v_appear_common(c *command) int {
	aura := c.a
	if aura < 1 {
		aura = 1
	}

	if FALSE == charge_aura(c.who, aura) {
		return FALSE
	}

	p := p_magic(c.who)
	if p.hide_mage == 0 {
		p.hide_mage = 1
	}
	p.hide_mage += aura

	wout(c.who, "Will appear common until the end of turn %d.", sysclock.turn+p.hide_mage-1)

	return TRUE
}

func v_tap_health(c *command) int {
	return TRUE
}

func d_tap_health(c *command) int {
	amount := c.a
	health := char_health(c.who)
	if amount > health/5 {
		amount = health / 5
	}
	add_aura(c.who, amount)

	wout(c.who, "Current aura is now %d.", char_cur_aura(c.who))

	add_effect(c.who, ef_tap_wound, 0, rnd(1, 60), amount*5)
	/* add_char_damage(c.who, amount * 5, MATES); */

	return TRUE
}

/*
 *  Mon May  5 12:25:36 1997 -- Scott Turner
 *
 *  Create dirt golem.  Decays after one year.
 *
 */
func v_create_dirt_golem(c *command) int {
	wout(c.who, "Begin construction of a dirt golem.")
	return TRUE
}

func d_create_dirt_golem(c *command) int {
	if FALSE == charge_aura(c.who, skill_piety(c.use_skill)) {
		return FALSE
	}

	// add an effect to destroy this golem in a year.
	if add_effect(c.who, ef_kill_dirt_golem, 0, 150+rnd(1, 60), 1) == 0 {
		wout(c.who, "For some reason, the blessing fails to take effect.")
		return FALSE
	}

	gen_item(c.who, item_dirt_golem, 1)
	wout(c.who, "You have created a dirt golem.")

	return TRUE
}
