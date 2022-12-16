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
	"os"
	"sort"
)

const (
	NOT_TAUGHT         = 0
	TAUGHT_SPECIFIC    = -1
	TAUGHT_GENERIC     = -2
	TAUGHT_STUDYPOINTS = -3
)

type use_tbl_ent struct {
	allow string /* who may execute the command */
	skill int

	start     func(*command) int /* initiator */
	finish    func(*command) int /* conclusion */
	interrupt func(*command) int /* interrupted order */

	time int /* how long command takes */
	poll int /* call finish each day, not just at end */
}

var (
	use_tbl []use_tbl_ent
)

func rep_skill_comp(a, b *skill_ent) int {
	//#if 1
	if a.know != SKILL_know && b.know == SKILL_know {
		return -1
	} else if b.know != SKILL_know && a.know == SKILL_know {
		return 1
	}
	//#else
	//    if ((*a).know == 0)
	//        return (*a).level - (*b).level;
	//
	//    if ((*b).level == 0)
	//        return (*b).level - (*a).level;
	//#endif

	pa := skill_school(a.skill) /* parent skill of a */
	pb := skill_school(b.skill) /* parent skill of b */
	if pa != pb {
		return pa - pb
	}
	return a.skill - b.skill
}

func flat_skill_comp(a, b *skill_ent) int {
	return a.skill - b.skill
}

func v_implicit(c *command) int {
	wout(c.who, "Use of this skill is automatic when appropriate.")
	wout(c.who, "No direct USE function exists.")
	return FALSE
}

func v_use_cs(c *command) int {
	sk := c.use_skill
	if c.a == 0 {
		wout(c.who, "Cleared all uses of %s.", box_name(sk))
		delete_all_effects(c.who, ef_cs, sk)
		return TRUE
	}

	var rounds string
	num := numargs(c)
	for i := 1; i <= num; i++ {
		arg := parse_arg(c.who, c.parse[i])
		if arg == 1 {
			wout(c.who, "You cannot use combat spells during the first round of combat.")
		} else if arg > 10 {
			wout(c.who, "You can only specify combat spell use up to round ten.")
		} else {
			add_effect(c.who, ef_cs, sk, -1, arg)
			rounds = sout("%s %d", rounds, arg)
		}
	}

	if rounds != "" {
		wout(c.who, "Use %s in rounds:%s.", box_name(sk), rounds)
	}

	return TRUE
}

func v_shipbuild(c *command) int {
	wout(c.who, "Use the BUILD order to build ships.")
	return FALSE
}

func find_use_entry(skill int) int {
	for i := 1; use_tbl[i].skill != 0; i++ {
		if use_tbl[i].skill == skill {
			return i
		}
	}
	return -1
}

/*
 *  If we know the skill, return the skill number
 *
 *  If an artifact grants us the ability to use the skill, return the item number of the artifact.
 *
 *  If a one-shot scroll lets us use the skill, return the item number of the scroll.
 *
 *  if there are multiple items, return the number of the last one found.
 */
func may_use_skill(who, sk int) int {
	if has_skill(who, sk) > 0 {
		return sk
	}

	// items other than scrolls should take precedence, to preserve the one-shot scrolls.
	ret, scroll := 0, 0
	for _, e := range loop_inventory(who) {
		p := rp_item_magic(e.item)
		if p != nil && ilist_lookup(p.MayUse, sk) >= 0 {
			if subkind(e.item) == sub_scroll || subkind(e.item) == sub_book {
				scroll = e.item
			} else {
				ret = e.item
			}
		}
	}

	if ret != 0 {
		return ret
	} else if scroll != 0 {
		return scroll
	}
	return 0
}

func magically_speed_casting(c *command, sk int) {
	if magic_skill(sk) && char_quick_cast(c.who) != FALSE && sk != sk_save_quick && sk != sk_trance {
		p := p_magic(c.who)
		var n int /* amount speeded by */
		if c.wait == 0 {
			n = 0 /* don't do anything */
		} else if p.quick_cast < c.wait {
			n = p.quick_cast
			c.wait -= p.quick_cast
			p.quick_cast = 0
		} else {
			n = c.wait - 1
			p.quick_cast = 0
			c.wait = 1
		}
		wout(c.who, "(speeded cast by %d day%s)", n, add_s(n))
	}
}

/*
 *  To use a spell through a scroll, the character should USE the
 *  spell number, and not USE the scroll.  If they use the scroll,
 *  guess which spell they meant to use from within the scroll,
 *  and use it.
 */
func correct_use_item(c *command) int {
	item := c.a
	if item_use_key(item) != FALSE {
		return item
	}

	p := rp_item_magic(item)
	if p == nil || len(p.MayUse) < 1 {
		return item
	}

	c.a = p.MayUse[0]
	return c.a
}

func meets_requirements(who, skill int) bool {
	p := rp_skill(skill)
	if p == nil {
		return true
	}

	l := p.req
	for i := 0; i < len(l); i++ { // yes, we want the increment here even though the body increments i, too
		for has_item(who, l[i].item) < l[i].qty && l[i].consume == REQ_OR {
			// if this assertion fails then a req list ended with a REQ_OR instead of a REQ_YES or a REQ_NO
			if i = i + 1; i >= len(l) {
				panic(fmt.Sprint("assert(skill != %d)", skill))
			}
		}
		if has_item(who, l[i].item) < l[i].qty {
			wout(who, "%s does not have %s.", just_name(who), box_name_qty(l[i].item, l[i].qty))
			return false
		}
		for i < len(l) && l[i].consume == REQ_OR {
			i++
		}
	}

	return true
}

func consume_requirements(who, skill int) {
	p := rp_skill(skill)
	if p == nil {
		return
	}

	l := p.req
	for i := 0; i < len(l); i++ { // yes, we want the increment here even though the body increments i, too
		for has_item(who, l[i].item) < l[i].qty && l[i].consume == REQ_OR {
			// if this assertion fails then a req list ended with a REQ_OR instead of a REQ_YES or a REQ_NO
			if i = i + 1; i >= len(l) {
				panic(fmt.Sprintf("assert(req list !end REQ_OR and skill != %d)", skill))
			}
		}
		item, qty := l[i].item, l[i].qty
		for i < len(l) && l[i].consume == REQ_OR {
			i++
		}
		if l[i].consume == REQ_YES {
			consume_item(who, item, qty)
		}
	}
}

func consume_scroll(who, basis, amount int) {
	if kind(basis) == T_item && p_item_magic(basis) != nil {
		// use up some days
		p_item_magic(basis).OrbUseCount -= amount

		// perhaps it is used up by now
		if p_item_magic(basis).OrbUseCount < 1 {
			wout(who, "%s crumbles into dust.", box_name(basis))
			if item_unique(basis) != FALSE {
				destroy_unique_item(who, basis)
			} else {
				consume_item(who, basis, 1)
			}
		}
	}
}

// shorten use of some skills based on experience
func experience_use_speedup(c *command) {
	exp := max(c.use_exp-1, 0)
	if exp != 0 && c.wait >= 7 {
		if c.wait >= 14 {
			c.wait -= exp
		} else if c.wait >= 10 {
			c.wait -= exp / 2
		} else if exp >= 2 {
			c.wait--
		}
	}
}

/*
 *  Wed Oct 30 13:22:42 1996 -- Scott Turner
 *
 *  Speed up strength skills by 1 day; slow weakness skills
 *  by 2 days.
 *
 */
func religion_use_speedup(c *command) {
	if is_follower(c.who) == FALSE {
		return
	}
	if rp_relig_skill(is_priest(is_follower(c.who))).strength == skill_school(c.use_skill) && c.wait >= 5 {
		c.wait--
	}
	if rp_relig_skill(is_priest(is_follower(c.who))).weakness == skill_school(c.use_skill) {
		c.wait += 2
	}
}

/*
 *  Tue Oct  6 18:18:27 1998 -- Scott Turner
 *
 *  Some artifacts speed skill use.
 *
 */
func artifact_use_speedup(c *command) {
	if art := best_artifact(c.who, ART_SPEED_USE, c.use_skill, 0); art != 0 {
		c.wait -= rp_item_artifact(art).Param1
		wout(c.who, "Using this skill is magically sped by %s day%s.", nice_num(rp_item_artifact(art).Param1), or_string(rp_item_artifact(art).Param1 > 1, "s", ""))
		if c.wait < 1 {
			c.wait = 1
		}
	}
}

func v_use(c *command) int {
	sk := c.a
	if !valid_box(sk) {
		wout(c.who, "%s is not a valid skill to use.", c.parse[1])
		c.use_skill = 0
		return FALSE
	}
	c.use_skill = sk

	if kind(sk) == T_item {
		sk = correct_use_item(c)
		if kind(sk) == T_item {
			return v_use_item(c)
		}
	}

	if kind(sk) != T_skill {
		wout(c.who, "%s is not a valid skill to use.", c.parse[1])
		return FALSE
	}

	parent := skill_school(sk)
	if parent == sk {
		wout(c.who, "Skill schools have no direct use.  Only subskills within a school may be used.")
		return FALSE
	}

	basis := may_use_skill(c.who, sk) /* what our skill ability is based upon */
	if basis == 0 {
		wout(c.who, "%s does not know %s.", just_name(c.who), box_code(sk))
		return FALSE
	}

	// if this was (has_skill(c.who, parent) < 1) then a category skill couldn't be used from an item.
	if FALSE == may_use_skill(c.who, parent) {
		wout(c.who, "Knowledge of %s is first required before %s may be used.", box_name(parent), box_code(sk))
		return FALSE
	}

	ent := find_use_entry(sk)
	if ent <= 0 {
		wout(c.who, "There's no way to use %s.", box_code(sk))
		return FALSE
	}

	/*
	 *  Mon May  5 12:31:15 1997 -- Scott Turner
	 *  Special checks for magic skills.
	 *
	 */
	if magic_skill(sk) {
		if in_safe_now(c.who) != FALSE {
			wout(c.who, "Magic may not be used in safe havens.")
			return FALSE
		}
		/*
		 *  Aura
		 *
		 */
		if skill_aura(sk) != 0 && !check_aura(c.who, skill_piety(sk)) {
			wout(c.who, "You don't have the aura required to use that spell.")
			return FALSE
		}
	}
	/*
	 *  Special checks for religion skills
	 *  Mon Oct 21 15:32:53 1996 -- Scott Turner
	 *
	 *  + Can't use religion in safe havens.
	 *
	 *  + Check piety, holy symbol, holy plant, etc.
	 *
	 */
	if religion_skill(sk) {
		/*
		 *  Safe Havens.
		 *
		 */
		if in_safe_now(c.who) != FALSE {
			wout(c.who, "The gods cannot answer prayers in safe havens.")
			return FALSE
		}
		/*
		 *  Right kind of priest.
		 *
		 */
		if (skill_school(sk) != sk_basic_religion) && is_priest(c.who) != skill_school(sk) {
			wout(c.who, "You must be a priest of %s to use that prayer.", god_name(skill_school(sk)))
			return FALSE
		}

		/*
		 *  Piety
		 *
		 */
		if !has_piety(c.who, skill_piety(sk)) {
			wout(c.who, "You don't have the piety required to use that prayer.")
			return FALSE
		}

		/*
		 *  Holy symbol
		 *
		 */
		if (skill_flags(sk)&REQ_HOLY_SYMBOL) != 0 && !has_holy_symbol(c.who) {
			wout(c.who, "You must have a holy symbol to use that prayer.")
			return FALSE
		}

		/*
		 *  Holy plant
		 *
		 */
		if (skill_flags(sk)&REQ_HOLY_PLANT) != 0 && !has_holy_plant(c.who) {
			wout(c.who, "You must have a holy plant to use that prayer.")
			return FALSE
		}
	}

	if !meets_requirements(c.who, sk) {
		return FALSE
	}

	cmd_shift(c)
	c.use_ent = ent
	c.use_skill = sk
	c.use_exp = has_skill(c.who, sk)
	/* c.poll = use_tbl[ent].poll; */
	c.poll = skill_flags(sk) & IS_POLLED
	/*	c.wait = use_tbl[ent].time; */
	c.wait = skill_time_to_use(sk)
	c.h = basis

	experience_use_speedup(c)
	religion_use_speedup(c)
	artifact_use_speedup(c)

	if use_tbl[ent].start != nil {
		ret := use_tbl[ent].start(c)
		if ret != 0 {
			magically_speed_casting(c, sk)
		}
		return ret
	}

	// perhaps an assertion here that we are indeed a production skill use, to catch skills without implementations
	if n := skill_produce(sk); n != 0 {
		wout(c.who, "Work to produce one %s.", just_name(n))
	} else if use_tbl[ent].finish == nil {
		panic(fmt.Sprintf("assert(sk != %q)", box_name(sk)))
	}

	return TRUE
}

/*
 *  Increment the experience counts for a skill and its parent
 */

func add_skill_experience(who, sk int) {
	/*
	 *  Don't increase the experience if we don't actually know the
	 *  skill.  For instance, use through a scroll or book should
	 *  not add experience to the character, unless the character
	 *  knows the skill himself.
	 */
	p := rp_skill_ent(who, sk)
	if p == nil {
		return
	}
	if p.exp_this_month == FALSE {
		p.experience++
		p.exp_this_month = TRUE
	}
}

func d_use(c *command) int {
	sk := c.use_skill
	var n int
	basis := c.h
	ret := TRUE

	if kind(sk) == T_item {
		return d_use_item(c)
	}

	/*
	 *  c.use_ent is not saved
	 *  if it is zero here, look it up again
	 *  This is so that it will be re-looked-up across turn boundaries
	 */
	ent := c.use_ent
	if ent <= 0 {
		ent = find_use_entry(sk)
		if ent <= 0 {
			fprintf(os.Stderr, "d_use: no use table entry for %s\n", c.parse[1])
			out(c.who, "Internal error.")
			return FALSE
		}
	}

	/*
	 *  Don't call poll routine for ordinary delays
	 */
	if c.wait > 0 && c.poll == 0 {
		return TRUE
	}
	if c.poll == 0 && !meets_requirements(c.who, sk) {
		return FALSE
	}

	/*
	 *  Maintain count of how many times each skill is used during
	 *  a turn, for informational purposes only.
	 *
	 *  Wed Dec 23 19:10:39 1998 -- Scott Turner
	 *
	 *  This counts each day of a polled use as 1!  not good.
	 */

	if sk != sk_breed_beasts && (c.poll == 0 || c.wait == 0) {
		p_skill(sk).use_count++
	}
	p_skill(sk).last_use_who = c.who

	/*
	 *  Special checks for religion skills
	 *  Mon Oct 21 15:32:53 1996 -- Scott Turner
	 *
	 *  + Use up piety and plants.
	 *
	 */
	if religion_skill(sk) {
		/*
		 *  Safe Havens.
		 *
		 */
		if in_safe_now(c.who) != FALSE {
			wout(c.who, "The gods cannot answer prayers in safe havens.")
			return FALSE
		}
		/*
		 *  Right kind of priest.
		 *
		 */
		if skill_school(sk) != sk_basic_religion && is_priest(c.who) != skill_school(sk) {
			wout(c.who, "You must be a priest of %s to use that prayer.", god_name(skill_school(sk)))
			return FALSE
		}

		/*
		 *  Holy symbol
		 *
		 */
		if (skill_flags(sk)&REQ_HOLY_SYMBOL) != 0 && !has_holy_symbol(c.who) {
			wout(c.who, "You must have a holy symbol to use that prayer.")
			return FALSE
		}

		/*
		 *  Piety
		 *
		 */
		if !has_piety(c.who, skill_piety(sk)) {
			wout(c.who, "You don't have the piety required to use that prayer.")
			return FALSE
		}

		/*
		 *  Holy plant
		 *
		 */
		if (skill_flags(sk)&REQ_HOLY_PLANT) != 0 && !has_holy_plant(c.who) {
			wout(c.who, "You must have a holy plant to use that prayer.")
			return FALSE
		}
	}

	/*
	 *  Thu Oct 24 15:18:44 1996 -- Scott Turner
	 *
	 *  Calculate return if there's a finish function.
	 *  Otherwise, it automatically worked.
	 *
	 */
	if use_tbl[ent].finish != nil {
		ret = use_tbl[ent].finish(c)
	}

	/*
	 *  Were we successful?
	 *
	 */
	if ret != 0 {
		/*
		 *  Special checks for religion skills
		 *  Mon Oct 21 15:32:53 1996 -- Scott Turner
		 *
		 *  + Use up piety and plants.
		 *
		 */
		if religion_skill(sk) {
			/*
			 *  Holy plant
			 *
			 */
			if (skill_flags(sk) & REQ_HOLY_PLANT) != 0 {
				move_item(c.who, 0, holy_plant(c.who), 1)
				wout(c.who, "Used one %s.", box_name(holy_plant(c.who)))
			}
			/*
			 *  Piety
			 *
			 */
			if !use_piety(c.who, skill_piety(sk)) {
				wout(c.who, "That was exceedingly strange.  Mention this to the DM, please.")
			} else if skill_piety(sk) != 0 {
				wout(c.who, "Used %s piety.", nice_num(skill_piety(sk)))
			}
		}

		/*
		 *  Experience.
		 *
		 *  Mon Nov 25 11:39:13 1996 -- Scott Turner
		 *
		 *  Modified to only augment for "make" skills.
		 *
		 */
		if c.wait == 0 && c.use_exp != 0 && ret != 0 {
			add_skill_experience(c.who, sk)
		}

		/*
		 *  Consumables
		 *
		 */
		consume_scroll(c.who, basis, c.wait)
		consume_requirements(c.who, sk)

		/*
		 *  Production
		 *
		 */
		if n = skill_produce(sk); n != 0 {
			gen_item(c.who, n, 1)
			wout(c.who, "Produced one %s.", box_name(n))
		}
	}

	return ret
}

func i_use(c *command) int {
	ent := c.use_ent
	if ent < 0 {
		out(c.who, "Internal error.")
		fprintf(os.Stderr, "i_use: c.use_ent is %d\n", c.use_ent)
		return FALSE
	}
	if use_tbl[ent].interrupt != nil {
		return use_tbl[ent].interrupt(c)
	}
	return 0 // todo: should return a value
}

/*
 *  Tue Dec  7 17:47:13 1999 -- Scott Turner
 *
 *  Use a special staff to find more staves.
 *
 */
func v_use_special_staff(c *command) int {
	var dist int
	where := province(c.who)
	reg := region(c.who)
	flag := false

	if (char_cur_aura(c.who) < 1 && !has_piety(c.who, 1)) ||
		in_faery(where) || in_hades(where) || in_clouds(where) {
		wout(c.who, "Nothing special happens.")
		return FALSE
	}

	for _, i := range loop_subkind(sub_special_staff) {
		/*
		 *  Don't report on items held by the same noble or not
		 *  in this region.
		 *
		 */
		if item_unique(i) == c.who {
			continue
		}
		if region(item_unique(i)) != reg {
			continue
		}
		dist = los_province_distance(c.who, item_unique(i))
		wout(c.who, "You sense another part of the staff %s province%s away.", nice_num(dist), or_string(dist == 1, "", "s"))
		flag = true
	}

	if !flag {
		wout(c.who, "Nothing special happens.")
		return FALSE
	}

	return TRUE
}

func v_use_item(c *command) int {
	c.poll = FALSE
	c.wait = 0

	item := c.a
	if has_item(c.who, item) < 1 {
		wout(c.who, "%s has no %s.", just_name(c.who), box_code(item))
		return FALSE
	}

	n := item_use_key(item)
	if n == 0 && subkind(item) == sub_special_staff {
		n = use_special_staff
	}
	if n == 0 && is_artifact(item) != nil {
		n = 100 + is_artifact(item).Type
	}
	if n == 0 {
		wout(c.who, "Nothing special happens.")
		return FALSE
	}

	//#if 0
	//    cmd_shift(c);
	//#endif

	/*
	 *  If they use a magical object, and we're in a safe
	 *  haven, don't allow it.
	 */

	switch n {
	case use_palantir, use_proj_cast, use_quick_cast, use_orb, use_barbarian_kill, use_savage_kill, use_corpse_kill, use_orc_kill, use_skeleton_kill, use_ancient_aura:
		if in_safe_now(c.who) != FALSE {
			wout(c.who, "Magic may not be used in safe havens.")
			c.wait = 0
			c.inhibit_finish = true
			return FALSE
		}
	}

	var ret bool
	switch n {
	case use_special_staff:
		ret = v_use_special_staff(c) != FALSE
	case use_heal_potion:
		ret = v_use_heal(c) != FALSE
	case use_slave_potion:
		ret = v_use_slave(c) != FALSE
	case use_death_potion:
		ret = v_use_death(c) != FALSE
	case use_fiery_potion:
		ret = v_use_fiery(c) != FALSE
	case use_proj_cast:
		ret = v_use_proj_cast(c) != FALSE
	case use_quick_cast:
		ret = v_use_quick_cast(c) != FALSE
	case use_drum:
		ret = v_use_drum(c) != FALSE
	case use_orb:
		ret = v_use_orb(c) != FALSE
	case use_weightlessness_potion:
		ret = v_use_weightlessness(c) != FALSE
	case use_nothing:
		wout(c.who, "Nothing happens.")
		ret = false
	case ART_PROT_FAERY + 100:
		ret = v_use_faery_artifact(c) != FALSE
	case ART_DESTROY + 100:
		ret = v_art_destroy(c) != FALSE
	case ART_POWER + 100:
		ret = v_power_jewel(c) != FALSE
	case ART_SUMMON_AID + 100:
		ret = v_summon_aid(c) != FALSE
	case ART_TELEPORT + 100:
		ret = v_art_teleport(c) != FALSE
	case ART_ORB + 100, ART_PEN + 100:
		ret = v_art_orb(c) != FALSE
	case ART_CROWN + 100:
		ret = v_art_crown(c) != FALSE
	default:
		wout(c.who, "Nothing happens.")
		ret = false
	}

	if ret || c.wait == 0 {
		c.wait = 0
		c.inhibit_finish = true
	}

	if ret {
		return TRUE
	}
	return FALSE
}

func d_use_item(c *command) int {
	item := c.a
	if has_item(c.who, item) < 1 {
		wout(c.who, "%s no longer has %s.", just_name(c.who), box_code(item))
		return FALSE
	}

	n := item_use_key(item)
	if n == 0 {
		wout(c.who, "Nothing special happens.")
		return FALSE
	}

	// todo: why?
	switch n {
	//#if 0
	//        case use_palantir:		return d_use_palantir(c);
	//        case use_ancient_aura:		return d_use_ancient_aura(c);
	//#endif
	default:
		panic(fmt.Sprintf("assert(key != %d)", n))
	}

	return TRUE
}

func exp_level(exp int) int {
	switch exp {
	case 0, 1:
		return exp_novice
	case 2, 3:
		return exp_journeyman
	case 4, 5:
		return exp_teacher
	case 6, 7, 8:
		return exp_master
	default:
		return exp_grand
	}
}

func exp_s(level int) string {
	switch level {
	case exp_novice:
		return "apprentice"
	case exp_journeyman:
		return "journeyman"
	case exp_teacher:
		return "adept"
	case exp_master:
		return "master"
	case exp_grand:
		return "grand master"
	}
	panic(fmt.Sprintf("assert(level != %d)", level))
}

func rp_skill_ent(who, skill int) *skill_ent {
	p := rp_char(who)
	if p == nil {
		return nil
	}
	for i := 0; i < len(p.skills); i++ {
		if p.skills[i].skill == skill {
			return p.skills[i]
		}
	}
	return nil
}

func p_skill_ent(who, skill int) *skill_ent {
	p := p_char(who)
	if p == nil {
		return nil
	}
	for i := 0; i < len(p.skills); i++ {
		if p.skills[i].skill == skill {
			return p.skills[i]
		}
	}
	newt := &skill_ent{}
	newt.skill = skill

	p.skills = append(p.skills, newt)

	return newt
}

/*
 *  Mon Dec 21 06:20:44 1998 -- Scott Turner
 *
 *  I believe that forgetting a category skill should cause you to also
 *  forget all the subskills of that category.
 *
 */
func forget_skill(who, skill int) bool {
	p := rp_char(who)
	if p == nil {
		return false
	}

	t := rp_skill_ent(who, skill)
	if t == nil {
		return false
	}

	p.skills = p.skills.rem_value(t)
	wout(who, "Forgot %s.", box_code(skill))

	/*
	 *  Recursively forget if a category skill.
	 *
	 */
	if skill_school(skill) == skill {
		for _, i := range loop_skill() {
			if i != skill && skill_school(i) == skill && p_skill_ent(who, i) != nil {
				forget_skill(who, i)
			}
		}
	}

	if magic_skill(skill) {
		ch := p_magic(who)
		ch.max_aura--
		if ch.max_aura < 0 {
			ch.max_aura = 0
		}
	}

	return true
}

func v_forget(c *command) int {
	skill := c.a
	if kind(skill) != T_skill {
		wout(c.who, "%s is not a skill.", box_code(skill))
		return FALSE
	} else if !forget_skill(c.who, skill) {
		wout(c.who, "Don't know any %s.", box_code(skill))
		return FALSE
	}
	return TRUE
}

/*
 *  Fri Aug  9 12:13:04 1996 -- Scott Turner
 *
 *  This is needed because we sometimes want to know if a char has
 *  any religion skill.
 *
 *  Fri Oct  9 18:28:33 1998 -- Scott Turner
 *
 *  We need to check all artifacts?  Blech.
 *
 *  Tue May 23 06:46:04 2000 -- Scott Turner
 *
 *  Modified to check all skills, not just the "known" ones.  That
 *  will pick up priests & magicians half-way through their training.
 *
 *  Fri Jun  2 11:57:26 2000 -- Scott Turner
 *
 *  Modified it back... causes problems elsewhere.
 *
 */
func has_subskill(who, subskill int) int {
	for _, e := range loop_char_skill_known(who) {
		if subkind(e.skill) == schar(subskill) {
			return e.skill
		}
	}

	/*
	 *  Run through all his ART_SKILL artifacts :-(
	 *
	 */
	for _, i := range loop_inventory(who) {
		if a := is_artifact(i.item); a != nil {
			if a.Type == ART_SKILL && subkind(a.Param1) == schar(subskill) {
				return a.Param1
			}
		}
	}

	return 0
}

/*
 *  Fri Oct  9 18:25:33 1998 -- Scott Turner
 *
 *  You might also have this skill via an artifact!
 *
 */
func has_skill(who, skill int) int {
	p := rp_skill_ent(who, skill)
	if p == nil || p.know != SKILL_know {
		// check if an artifact gives the skill
		if a := has_artifact(who, ART_SKILL, skill, 0, 0); a != FALSE {
			return exp_teacher
		}
		return 0
	}
	return exp_level(p.experience)
}

/*
 *  Use learn_skill() to grant a character a skill
 */
func set_skill(who, skill, know int) {
	p_skill_ent(who, skill).know = know
}

func skill_school(sk int) int {
	for count := 0; count < 1_000; count++ {
		if !valid_box(sk) || nil == rp_skill(sk) {
			return FALSE
		}
		n := req_skill(sk)
		if n == 0 {
			return sk
		}
		sk = n
	}
	panic("assert(count < 1000)")
}

/*
 *  Order skills for display
 *
 *  Subskills follow their parent
 *  Skills we don't know are pushed to the end
 */

//#if 0
//void
//list_skill_sup(int who, struct skill_ent *e)
//{
//    char *exp;
//
//    if (skill_no_exp(e.skill) || skill_school(e.skill) == e.skill)
//        wiout(who, CHAR_FIELD+2, "%*s  %s",
//            CHAR_FIELD,
//            box_code_less(e.skill),
//            cap(just_name(e.skill)));
//    else
//        wiout(who, CHAR_FIELD+2, "%*s  %s, %s %s",
//              CHAR_FIELD,
//              box_code_less(e.skill),
//              cap(just_name(e.skill)),
//              exp_s(exp_level(e.experience)),
//              practice_s(e));
//}
//
//#else

func list_skill_sup(who int, e *skill_ent, prefix string) {
	if len(prefix) == 0 {
		tagout(who, "<tag type=skill id=%d skill=%d exp=%d>", who, e.skill, e.experience)
	}

	if skill_no_exp(e.skill) != FALSE || skill_school(e.skill) == e.skill {
		wout(who, "%s%s", prefix, box_name(e.skill))
	} else {
		wout(who, "%s%s, (Level %d)", prefix, box_name(e.skill), e.experience)
	}

	/*
	 *  Fri Apr 27 10:55:26 2001 -- Scott Turner
	 *
	 *  Add a line for combat skills that says what round you will use them.
	 *
	 */
	if kind(who) == T_char && combat_skill(e.skill) != FALSE && get_effect(who, ef_cs, e.skill, 0) != FALSE {
		var rounds string
		/*
		 *  Walk through all the effects and gather up the rounds.
		 *
		 */
		ef := effects(who)
		for i := 0; i < len(ef); i++ {
			if ef[i].type_ == ef_cs && ef[i].subtype == e.skill {
				rounds = fmt.Sprintf("%s %d", rounds, ef[i].data)
			}
		}
		wout(who, "  Use in rounds: %s", rounds)
	}

	if len(prefix) == 0 {
		tagout(who, "</tag type=skill>")
	}
}

//#endif

func list_skills(who, num int, prefix string) {
	assert(valid_box(num))

	if len(prefix) == 0 {
		out(who, "")
	}
	out(who, "%sSkills known:", prefix)
	indent += 3

	flag := true
	if rp_char(num) != nil && len(rp_char(num).skills) >= 1 {
		l := rp_char(num).skills.copy()
		sort.Slice(l, func(i, j int) bool {
			return rep_skill_comp(l[i], l[j]) < 0
		})

		for i := 0; i < len(l); i++ {
			if l[i].know != SKILL_know {
				continue
			}

			flag = false

			//#if 0
			//  if (i > 0 && skill_school(l[i].skill) != skill_school(l[i-1].skill)) {
			//    out(who, "");
			//  }
			//#endif
			if req_skill(l[i].skill) != FALSE {
				indent += 6
				list_skill_sup(who, l[i], prefix)
				indent -= 6
			} else {
				//#if 0
				//  if (i > 0 && req_skill(l[i-1].skill)) {
				//    out(who, "");
				//  }
				//#endif
				list_skill_sup(who, l[i], prefix)
			}
		}
	}

	if flag {
		out(who, "%snone", prefix)
	}

	indent -= 3
}

/*
 *  Archery, 0/7
 *  Archery, 1/7
 *  Archery, 0/7, 1 NP req'd
 */
func fractional_skill_qualifier(p *skill_ent) string {
	assert(p.know != SKILL_know)

	if p.know == SKILL_dont {
		return sout("0/%d%s", learn_time(p.skill), np_req_s(p.skill))
	}

	var explanation string

	assert(p.know == SKILL_learning)
	if p.days_studied/TOUGH_NUM >= learn_time(p.skill) {
		explanation = sout("(Religion weakness skill)")
	}

	frac := p.days_studied % TOUGH_NUM
	if frac == 0 {
		return sout("%d/%d %s", p.days_studied/TOUGH_NUM, learn_time(p.skill), explanation)
	}

	return sout("%d.%d/%d %s", p.days_studied/TOUGH_NUM, frac*100/TOUGH_NUM, learn_time(p.skill), explanation)
}

func list_partial_skills(who, num int, prefix string) {
	assert(valid_box(num))
	if rp_char(num) == nil {
		return
	} else if len(rp_char(num).skills) < 1 {
		return
	}

	l := rp_char(num).skills.copy()
	sort.Slice(l, func(i, j int) bool {
		return flat_skill_comp(l[i], l[j]) < 0
	})

	flag := true
	for i := 0; i < len(l); i++ {
		if l[i].know == SKILL_know {
			continue
		}

		if flag {
			if len(prefix) == 0 {
				out(who, "")
			}
			out(who, "%sPartially known skills:", prefix)
			indent += 3
			flag = false
		}

		wiout(who, 6, "%s%s, %s", prefix, box_name(l[i].skill), fractional_skill_qualifier(l[i]))
	}

	if !flag {
		indent -= 3
	}
}

func skill_cost(sk int) int {
	return 100 // todo: should this be a constant value?
}

/*
 *  Thu Dec 24 08:12:51 1998 -- Scott Turner
 *
 *  Can "teach" in a tower or a guild, if you're teaching a skill
 *  from that guild.
 *
 */
func teachable_place(where, sk int) bool {
	return (subkind(where) == sub_tower || is_guild(where) == skill_school(sk))
}

/*
 *  Sat Nov  2 09:00:07 1996 -- Scott Turner
 *
 *  Returns one of the following flags or the id of the scroll/book
 *  being used.
 *
 *  We want to try to find TAUGHT_SPECIFIC, then scroll/book, then
 *  TAUGHT_GENERIC.
 *
 *  Thu Apr 16 09:04:12 1998 -- Scott Turner
 *
 *  Added a separate return for the item being used.
 *
 */
func being_taught(who, sk int, item, teach_bonus *int) int {
	school := skill_school(sk)
	where := subloc(who)
	p := rp_subloc(where)
	var teacher int
	var c, c2 *command
	ret_specific := 0
	ret_generic := 0

	*item = 0

	// possibly being taught in this location.
	taught_specific := (p != nil && ilist_lookup(p.teaches, sk) >= 0)
	taught_generic := (p != nil && ilist_lookup(p.teaches, school) >= 0)

	/*
	 *  Are you in a teaching tower?
	 *
	 */
	if teachable_place(where, sk) {
		/*
		 *  Mon Dec 21 07:23:15 1998 -- Scott Turner
		 *
		 *  Teacher is the first noble in the tower using the
		 *  "teach" command.
		 *
		 *  Thu Mar 23 12:51:41 2000 -- Scott Turner
		 *
		 *  Hmm, must check to make sure he's actually able to teach
		 *  this command!
		 *
		 */
		for _, i := range loop_stack(where) {
			c = rp_command(i)
			if c != nil && c.state == RUN && fmt.Sprintf("%p", cmd_tbl[c.cmd].start) == fmt.Sprintf("%p", v_teach) {
				// todo: the function comparison hack above is not 100% safe
				teacher = i
				break
			}
		}

		/*  Is he teaching? */
		if teacher != 0 && c != nil && fmt.Sprintf("%p", cmd_tbl[c.cmd].start) == fmt.Sprintf("%p", v_teach) {
			/*  Are we close to him? */
			/*
			 *  Mon Dec 21 06:55:29 1998 -- Scott Turner
			 *
			 *  Modified to say that you only have to be one of the first
			 *  five students that could benefit from his teaching.
			 *
			 */
			loc := 0
			for _, i := range loop_stack(where) {
				if i == who {
					break
				} /* you. */
				/*
				 *  Get this person's command.
				 *
				 */
				c2 = rp_command(i)
				/*
				 *  Is this person studying and benefiting
				 *  from the teacher?
				 *
				 */
				if c2 != nil && fmt.Sprintf("%p", cmd_tbl[c2.cmd].start) == fmt.Sprintf("%p", v_study) && (c2.a == c.a || skill_school(c2.a) == c.a) {
					// todo: the function comparison hack above is not 100% safe
					loc++
				}
			}

			if loc <= 6 {
				/*  What is he teaching? */
				if c.a == sk {
					taught_specific = true
				}
				if c.a == school {
					taught_generic = true
				}
				/*
				 *  Maybe he has an artifact bonus.
				 *
				 */
				if a := best_artifact(c.who, ART_TEACHING, 0, 0); a != FALSE {
					*teach_bonus += rp_item_artifact(a).Param1
				}
			}
		}
	}

	/*
	 *  Maybe you have a book?
	 *
	 */

	for _, e := range loop_inventory(who) {
		p := rp_item_magic(e.item)
		if p != nil && ilist_lookup(p.MayStudy, sk) >= 0 && p.OrbUseCount > 0 && subkind(e.item) == sub_book {
			ret_specific = e.item
			break
		} else if p != nil && ilist_lookup(p.MayStudy, school) >= 0 && p.OrbUseCount > 0 && subkind(e.item) == sub_book {
			ret_generic = e.item
		}
	}

	/*
	 *  Wed Apr 15 11:58:33 1998 -- Scott Turner
	 *
	 *  Another version of study points, that provides instruction
	 *  albeit only in the "common" skills.
	 *
	 */
	taught_studypoints := false
	if (sk == int(sk_combat) ||
		sk == int(sk_shipcraft) ||
		sk == int(sk_construction) ||
		sk == int(sk_forestry) ||
		sk == int(sk_ranger) ||
		sk == int(sk_mining) ||
		sk == int(sk_trading)) && player_js(player(who)) != 0 {
		taught_studypoints = true
	}

	if taught_specific {
		return TAUGHT_SPECIFIC
	}
	if ret_specific != 0 {
		*item = ret_specific
		return TAUGHT_SPECIFIC
	}
	if taught_generic {
		return TAUGHT_GENERIC
	}
	if ret_generic != 0 {
		*item = ret_generic
		return TAUGHT_GENERIC
	}
	if taught_studypoints {
		return TAUGHT_STUDYPOINTS
	}
	return NOT_TAUGHT
}

//#if 0
///*
// *  If it's taught by the location
// *  If it's offered by a skill that we know
// *  If it's offered by an item that we have
// */
//
//static int
//may_study(who, sk int)
//{
//    where := subloc(who);
//
///*
// *  Does the location offer the skill?
// */
//
//    {
//        struct entity_subloc *p;
//
//        p = rp_subloc(where);
//
//        if (p && ilist_lookup(p.teaches, sk) >= 0)
//            return where;
//    }
//
///*
// *  Is the skill offered by a skill that we already know?
// */
//
//    {
//        var q *entity_skill
//        var e *skill_ent
//        ret := 0;
//
//        for _, e = range loop_char_skill_known(who, e)
//        {
//            q = rp_skill(e.skill);
//
//            if (q && ilist_lookup(q.offered, sk) >= 0)
//            {
//                ret = e.skill;
//                break;
//            }
//        }
//
//
//        if (ret)
//            return ret;
//    }
//
///*
// *  Is instruction offered by a scroll or a book?
// *
// *  Items other than scrolls should take precedence, to preserve
// *  the one-shot scrolls.
// */
//
//    {
//        var e *item_ent
//        var p *ItemMagic
//        ret := 0;
//        scroll := 0;
//
//        for _, e = range loop_inventory(who, e)
//        {
//            p = rp_item_magic(e.item);
//            if (p && ilist_lookup(p.may_study, sk) >= 0)
//            {
//                if (subkind(e.item) == sub_scroll)
//                    scroll = e.item;
//                else
//                    ret = e.item;
//            }
//        }
//
//
//        if (ret)
//            return ret;
//
//        if (scroll)
//            return scroll;
//    }
//
//    return 0;
//}
//#endif

func begin_study(c *command, sk int) bool {
	cost := skill_cost(sk)

	/*
	 *  Fri Nov  1 13:02:38 1996 -- Scott Turner
	 *
	 *  Add on cost if this is the 4th+ skill category this guy
	 *  is learning.
	 *
	 */
	num_base_skills := 0
	for _, e := range loop_char_skill(c.who) {
		if e.skill == sk_basic_religion || e.skill == sk_adv_sorcery {
			continue
		}
		if FALSE == rp_skill(e.skill).required_skill {
			num_base_skills++
		}
	}

	np_req := skill_np_req(sk)
	if num_base_skills >= 3 && FALSE == rp_skill(sk).required_skill {
		np_req++
	}

	if !is_npc(c.who) && np_req > 0 {
		wout(c.who, "It will require %s noble point%s to complete the study of %s.",
			cap_(nice_num(np_req)),
			add_s(np_req),
			box_code(sk))
	}

	if !is_npc(c.who) && cost > 0 {
		if !charge(c.who, cost) {
			wout(c.who, "Cannot afford %s to begin study.", gold_s(cost))
			return false
		}
		wout(c.who, "Paid %s to begin study.", gold_s(cost))
	}

	p := p_skill_ent(c.who, sk)
	p.know = SKILL_learning

	return true
}

func end_study(c *command, sk int) bool {
	/*
	 *  Fri Nov  1 13:02:38 1996 -- Scott Turner
	 *
	 *  Add on cost if this is the 4th+ skill category this guy
	 *  is learning.
	 *
	 */
	num_base_skills := 0
	for _, e := range loop_char_skill(c.who) {
		if e.know != SKILL_know {
			continue
		}
		if e.skill == sk_basic_religion || e.skill == sk_adv_sorcery {
			continue
		}
		if FALSE == rp_skill(e.skill).required_skill {
			num_base_skills++
		}
	}

	np_req := skill_np_req(sk)
	if num_base_skills >= 3 && FALSE == rp_skill(sk).required_skill {
		np_req++
	}

	if !is_npc(c.who) && np_req > 0 {
		wout(c.who, "Deducted %s noble point%s to complete study.", nice_num(np_req), add_s(np_req))
		deduct_np(player(c.who), np_req)
	} else if np_req > 0 {
		wout(c.who, "%s noble point%s %s required to complete the study of %s.", cap_(nice_num(np_req)), add_s(np_req), or_string(np_req == 1, "is", "are"), box_code(sk))
		return false
	}

	return true
}

//#if 0
///*
// *  To study through a scroll, the character should STUDY the
// *  skill number, not study the scroll.  If they study the scroll,
// *  guess which skill they meant to learn from within the scroll.
// */
//static int
//correct_study_item(c *command)
//{
//  item := c.a;
//  var p *ItemMagic
//
//  p = rp_item_magic(item);
//
//  if (p == nil || len(p.may_study) < 1)
//    return item;
//
//  c.c = item;
//  c.a = p.may_study[0];
//  return c.a;
//}
//#endif

/*
 *  Thu Apr  9 08:48:40 1998 -- Scott Turner
 *
 *  Code to check to see if we can study something.
 *
 */
func check_study(c *command, requires_instruction bool) bool {
	//#if 0
	//    /*
	//     *  Thu Jan  7 18:30:50 1999 -- Scott Turner
	//     *
	//     *  I believe this is no longer needed.
	//     *
	//     */
	//
	//    /*
	//     *  In case they're studying from a book...
	//     *
	//     */
	//    c.c = 0;
	//    if (kind(sk) == T_item) sk = correct_study_item(c);
	//#endif

	sk := c.a
	if kind(sk) != T_skill {
		wout(c.who, "%s is not a valid skill.", c.parse[1])
		return false
	}

	parent := skill_school(sk)

	if nation(c.who) != 0 &&
		(ilist_lookup(rp_nation(nation(c.who)).proscribed_skills, sk) != -1 ||
			ilist_lookup(rp_nation(nation(c.who)).proscribed_skills, parent) != -1) {
		wout(c.who, "%ss may not learn that skill.", rp_nation(nation(c.who)).citizen)
		return false
	}

	if parent != sk && has_skill(c.who, parent) < 1 {
		wout(c.who, "%s must be learned before %s can be known.",
			cap_(box_name(parent)),
			box_code(sk))
		return false
	}

	p := rp_skill_ent(c.who, sk)
	if p != nil && p.know == SKILL_know {
		wout(c.who, "Already know %s.", box_name(sk))
		return false
	}

	/*
	 *  Tue Feb 22 09:00:39 2000 -- Scott Turner
	 *
	 *  You cannot study Basic Religion or Advanced Sorcery
	 *
	 */
	if sk == sk_basic_religion || sk == sk_adv_sorcery {
		wout(c.who, "You cannot study that skill.")
		return false
	}

	/*
	 *  Fri Aug  9 12:17:47 1996 -- Scott Turner
	 *
	 *  Can't learn a (top-level) religion skill if you already
	 *  know one.
	 *
	 */
	if is_priest(c.who) != FALSE &&
		religion_skill(sk) &&
		skill_school(sk) != sk_basic_religion &&
		skill_school(sk) != is_priest(c.who) {
		wout(c.who, "A priest may not learn a second religion.")
		return false
	}

	/*
	 *  Fri Nov 29 09:24:28 1996 -- Scott Turner
	 *
	 *  Priests must be guild members of their strength skill.
	 *
	 */
	if religion_skill(sk) &&
		skill_school(sk) != sk_basic_religion &&
		guild_member(c.who) != 0 &&
		guild_member(c.who) != rp_relig_skill(skill_school(sk)).strength {
		wout(c.who, "Because you're a member of the %s Guild, the %s will not accept you.", box_name(guild_member(c.who)), box_name(sk))
		return false
	}

	/*
	 *  Fri Aug  9 12:19:36 1996 -- Scott Turner
	 *
	 *  Priests can't learn magic.
	 */
	if magic_skill(sk) && is_priest(c.who) != FALSE {
		wout(c.who, "Priests may not study magic.")
		return false
	}

	/*
	 *  Tue Apr  7 11:35:34 1998 -- Scott Turner
	 *
	 *  Followers can't learn magic.
	 */
	if magic_skill(sk) && is_follower(c.who) != FALSE {
		wout(c.who, "Your dedication to %s prevents you from learning magic.", box_name(is_priest(is_follower(c.who))))
		return false
	}

	/*
	 *  Fri Aug  9 12:19:36 1996 -- Scott Turner
	 *
	 *  Priests can't learn their weakness skill.
	 *
	 */
	if is_priest(c.who) != FALSE && rp_skill(is_priest(c.who)).religion_skill.weakness == sk {
		wout(c.who, "Priests may not study their weakness skill.")
		return false
	}

	/*
	 *  Fri Nov 29 09:27:29 1996 -- Scott Turner
	 *
	 *  Uh, magicians can't learn religion, either.
	 *
	 */
	if religion_skill(sk) && is_wizard(c.who) != FALSE {
		wout(c.who, "Magicians may not study religion.")
		return false
	}

	/*
	 *  Thu Apr  9 08:55:53 1998 -- Scott Turner
	 *
	 *  If we don't require instruction, then we're done.
	 *
	 */
	if is_npc(c.who) || !requires_instruction {
		return true
	}

	/*
	 *  What kind of skill are they learning?
	 *
	 *  Fri Nov 29 09:53:08 1996 -- Scott Turner
	 *
	 *  Added guild skills.
	 */
	category, guild, teachable, unteachable := (parent == sk), false, false, false
	if !category {
		guild = ilist_lookup(rp_skill(parent).guild, sk) != -1
		teachable = ilist_lookup(rp_skill(parent).offered, sk) != -1
	}
	if !category && !teachable && !guild {
		unteachable = true
	}

	/*
	 *  What's the teaching situation?
	 *
	 */
	item := 0
	bonus := 0
	bt := being_taught(c.who, sk, &item, &bonus)

	/*
	 *  Correct for missing flag.
	 *
	 */
	if category && bt == TAUGHT_STUDYPOINTS && c.b == 0 {
		wout(c.who, "You have no source of instruction.")
		wout(c.who, "To use study points, add a flag to your study command, e.g.,")
		wout(c.who, "    study %d 1", c.a)
		return false
	}

	/*
	 *  Category skills must be taught.
	 *
	 */
	if category && bt == 0 {
		wout(c.who, "Category skills may only be learned with a teacher.")
		return false
	}

	/*
	 *  You must be a guild member to learn a guild skill.
	 *
	 */
	if guild && guild_member(c.who) != skill_school(sk) {
		wout(c.who, "You must be a guild member to learn that skill.")
		return false
	}

	/*
	 *  Guild skills must be learned in a guild.
	 *
	 */
	if guild && is_guild(subloc(c.who)) != skill_school(sk) {
		wout(c.who, "Guild skills may only be studied in a guild.")
		return false
	}

	/*
	 *  Unteachable skills do not benefit from a teacher.
	 *
	 */
	if unteachable && bt != 0 {
		wout(c.who, "Instruction does not benefit in learning %s.", box_name(sk))
	}

	return true
}

/*
 *	if we know it already, then error
 *	if we have already started studying, then continue
 *	if we may not study it, then error
 *
 *	if we don't have enough money, then error
 *
 *	start studying skill:
 *		init skill entry
 *		deduct money
 *
 *   Fri Aug  9 12:11:36 1996 -- Scott Turner
 *
 *   Added restrictions due to religion.
 *
 *   Sat Nov  2 08:59:33 1996 -- Scott Turner
 *
 *   New teaching changes
 *
 *   Tue Jul  1 16:24:00 1997 -- Scott Turner
 *
 *   Can't learn a skill on your nation's "proscribed_skill" list.
 *
 *  Wed Apr 22 11:09:09 1998 -- Scott Turner
 *
 *  NPCs don't need instruction or money, but learn slowly
 *
 *  Wed Oct 28 07:18:37 1998 -- Scott Turner
 *
 *  Need to save the "solitary" status somewhere in "c" so that we can
 *  check that it has held the entire time of the study, not just at the
 *  end.
 *
 */
func v_study(c *command) int {
	sk := c.a

	//#if 0
	//    /*
	//     *  In case they're studying from a book...
	//     *
	//     */
	//    c.c = 0;
	//    if (kind(sk) == T_item) sk = correct_study_item(c);
	//#endif

	if numargs(c) < 1 {
		wout(c.who, "Must specify a skill to study.")
		return FALSE
	}

	if !check_study(c, true) {
		return FALSE
	}

	/*
	 *  Give him the skill and charge him.
	 *
	 */
	p := rp_skill_ent(c.who, sk)
	if p == nil {
		if !begin_study(c, sk) {
			return FALSE
		}
	}

	/*
	 *  If you're alone in a tower, then get a +1 day bonus.
	 *  We keep the status for this in c.c.
	 */
	c.c = 1
	where := subloc(c.who)
	if subkind(where) != sub_tower || !alone_here(c.who) {
		c.c = 0
	}

	wout(c.who, "Study %s for %s day%s.", just_name(sk), nice_num(c.wait), add_s(c.wait))

	return TRUE
}

/*
 *  Use learn_skill() to grant a character a skill
 *
 *  Fri Aug  9 11:47:15 1996 -- Scott Turner
 *
 *  Modified for religions.
 *
 */

func learn_skill(who, sk int) bool {
	if nation(who) != 0 &&
		(ilist_lookup(rp_nation(nation(who)).proscribed_skills, sk) != -1 ||
			ilist_lookup(rp_nation(nation(who)).proscribed_skills, skill_school(sk)) != -1) {
		return false
	}

	p := p_skill_ent(who, sk)

	/*
	 *  Archery grants a missile attack.
	 *
	 */
	if sk == sk_archery {
		pc := p_char(who)
		if pc.missile < 50 {
			pc.missile += 50
		}
	}

	/*
	 *  Personal FTTD needs to set the personal break point
	 *
	 */
	if sk == sk_personal_fttd {
		p_char(who).personal_break_point = 100
	}

	ch := p_magic(who)

	if magic_skill(sk) {
		ch.max_aura++
		add_aura(who, 1)

		wout(who, "Maximum aura now %d.", ch.max_aura)

		p_magic(who).magician = TRUE

		if sk == sk_weather {
			p_magic(who).knows_weather = 1
		}

		/*
		 *  Undedicate him, if he's dedicated to someone.
		 *
		 */
		if rp_char(who).religion.priest != 0 {
			wout(rp_char(who).religion.priest, "An angel informs you that your follower %s has become a heretic.", box_name(who))
			rp_char(rp_char(who).religion.priest).religion.followers = rem_value(rp_char(rp_char(who).religion.priest).religion.followers, who)
			rp_char(who).religion.priest = 0
		}
	}

	/*
	 *  It's his first religion skill, so let's do all the
	 *  special tasks for religion skills.
	 *
	 */
	if religion_skill(sk) && is_priest(who) == FALSE {
		/*
		 *  Undedicate him, if he's dedicated to someone.
		 *
		 */
		if rp_char(who).religion.priest != FALSE {
			wout(rp_char(who).religion.priest, "An angel informs you that your follower %s has become a priest of %s.", box_name(who), god_name(sk))
			rp_char(rp_char(who).religion.priest).religion.followers = rem_value((rp_char(rp_char(who).religion.priest).religion.followers), who)
			rp_char(who).religion.priest = 0
		}

		/*
		 *  Eliminate any magic skills.
		 *
		 */
		for _, e := range loop_char_skill_known(who) {
			if magic_skill(e.skill) {
				forget_skill(who, e.skill)
			}
		}
		p_magic(who).max_aura = 0
		p_magic(who).cur_aura = 0

		/*
		 *  Eliminate the "weakness" skill.
		 *
		 */
		if rp_skill(sk).religion_skill.weakness != FALSE &&
			rp_skill(rp_skill(sk).religion_skill.weakness) != nil {
			forget_skill(who, rp_skill(sk).religion_skill.weakness)
		}
		/*
		 *  Add the "strength skill" -- careful of a loop here!
		 *
		 */
		if rp_skill(sk).religion_skill.strength != FALSE &&
			rp_skill(rp_skill(sk).religion_skill.strength) != nil {
			learn_skill(who, rp_skill(sk).religion_skill.strength)
		}

		/*
		 *  Tue Dec 29 12:43:33 1998 -- Scott Turner
		 *
		 *  He also learns "Basic Religion"
		 *
		 *  Thu Feb 10 19:27:13 2000 -- Scott Turner
		 *
		 *   Make sure he knows the religion skill, so we don't
		 *   get into an infinite loop.
		 *
		 */
		p.know = SKILL_know
		p.days_studied = 0
		learn_skill(who, sk_basic_religion)

		/*
		 *  Mon Apr  2 08:43:40 2001 -- Scott Turner
		 *
		 *  Give him enough piety to dedicate a temple.
		 *
		 */
		add_piety(who, skill_piety(sk_dedicate_temple_b), false)
		wout(who, "%s grants you %s piety for your new devotion as priest.", god_name(sk), nice_num(skill_piety(sk_dedicate_temple_b)))
	}

	wout(who, "Learned %s.", box_name(sk))
	p.know = SKILL_know
	p.days_studied = 0

	return true
}

/*
 *  Wed Apr 15 12:10:44 1998 -- Scott Turner
 *
 *  Use up a study point.
 *
 */
func use_studypoint(who int) {
	pl := player(who)
	if pl != 0 {
		rp_player(pl).JumpStart--
		if rp_player(pl).JumpStart < 0 {
			rp_player(pl).JumpStart = 0
		}
	}
}

/*
 *  Mon Dec 13 18:24:28 1999 -- Scott Turner
 *
 *  To be alone here, you must be the only unit in the location, and not
 *  have any men in your inventory.
 *
 */
func alone_here(who int) bool {
	if len(all_char_here(subloc(who), nil)) == 1 {
		return true
	}
	return false
}

/*
 *  Note:  d_study is polled daily
 */
func d_study(c *command) int {
	sk := c.a
	if kind(sk) != T_skill {
		log_output(LOG_CODE, "d_study: skill %d is gone, who=%d\n", sk, c.who)
		out(c.who, "Internal error: skill %s is gone", box_code(sk))
		return FALSE
	}

	ch := p_char(c.who)
	ch.studied++

	var diminish int
	if ch.studied <= 7 {
		diminish = 14
	} else if ch.studied <= 14 {
		diminish = 8
	} else if ch.studied <= 21 {
		diminish = 4
	} else {
		diminish = 2
	}

	/*
	 *  What kind of skill are they learning?
	 *
	 */
	parent := skill_school(sk)
	category := (parent == sk)
	teachable := !category && ilist_lookup(rp_skill(parent).offered, sk) != -1
	//unteachable := !category && !teachable

	/*
	 *  Sun Nov  3 19:08:28 1996 -- Scott Turner
	 *
	 *  Adjust diminish for the state of being taught.
	 *
	 */
	var artifact_bonus, item int
	bt := being_taught(c.who, sk, &item, &artifact_bonus)
	if !is_npc(c.who) && category && bt == 0 {
		wout(c.who, "Category skills may only be learned with a teacher.")
		return FALSE
	}

	if (teachable || category) && bt == TAUGHT_STUDYPOINTS && c.b == 0 {
		if category {
			wout(c.who, "You have no source of instruction.")
			wout(c.who, "To use study points, add a flag to your study command, e.g.,")
			wout(c.who, "    study %d 1", c.a)
			return FALSE
		}
		bt = 0
	}

	/*
	 *  Unteachable skills do not benefit from a teacher.
	 *
	 */
	if teachable && bt != 0 {
		if bt == TAUGHT_GENERIC {
			diminish *= 2
		} else if bt == TAUGHT_SPECIFIC {
			diminish *= 4
		} else if bt == TAUGHT_STUDYPOINTS {
			use_studypoint(c.who)
		}
	}

	if category && bt == TAUGHT_STUDYPOINTS {
		use_studypoint(c.who)
	}

	if (teachable || category) && bt != 0 && item != 0 {
		if ip := rp_item_magic(item); ip != nil {
			consume_scroll(c.who, item, 1)
		}
	}

	p := p_skill_ent(c.who, sk)
	npc_penalty := or_int(is_npc(c.who), 2, 1)
	p.days_studied += (TOUGH_NUM * diminish) / (14 * npc_penalty)

	/*
	 *  Wed Oct 30 13:00:44 1996 -- Scott Turner
	 *
	 *  If you're a follower, then you get a bonus of 1 day in
	 *  learning skills in your strength skill category, and -2 days
	 *  in your weakness category.
	 *
	 */
	religion_bonus := 0
	if is_follower(c.who) != 0 &&
		rp_relig_skill(is_priest(is_follower(c.who))).strength == skill_school(sk) {
		religion_bonus = TOUGH_NUM
	}

	if is_follower(c.who) != 0 &&
		rp_relig_skill(is_priest(is_follower(c.who))).weakness == skill_school(sk) {
		religion_bonus = -2 * TOUGH_NUM
	}

	/*
	 *  If you're alone in a tower, then get a +1 day bonus.
	 *
	 */
	where := subloc(c.who)
	if subkind(where) != sub_tower || !alone_here(c.who) {
		c.c = 0
	}
	building_bonus := 0
	if c.c != 0 {
		building_bonus = TOUGH_NUM
	}

	/*
	 *  Bonuses from artifacts.
	 *
	 */
	if a := best_artifact(c.who, ART_LEARNING, 0, 0); a != 0 {
		artifact_bonus += rp_item_artifact(a).Param1
	}

	if p.days_studied+religion_bonus+building_bonus+artifact_bonus >= learn_time(sk)*TOUGH_NUM {
		/*
		 *  Report bonuses.
		 *
		 */
		if building_bonus > 0 {
			wout(c.who, "Your study benefits from the solitude of your tower.")
		}
		if religion_bonus > 0 {
			wout(c.who, "Your study is inspired by holy energy.")
		}
		if religion_bonus < 0 {
			wout(c.who, "This topic offends your values and is difficult to learn.")
		}

		/*
		 *  Now try to end our study
		 *
		 */
		if !end_study(c, sk) {
			c.wait = 0
			return FALSE
		}
		/*
		 *  Otherwise learn the skill and be done.
		 *
		 */
		learn_skill(c.who, sk)
		c.wait = 0
		return TRUE
	}

	/* message here if done? */
	return TRUE
}

//#if 0
//int
//v_acquire(c *command)
//{
//  var p *skill_ent
//  sk := c.a;
//
//  if (numargs(c) < 1) {
//    wout(c.who, "Must specify a skill to acquire.");
//    return FALSE;
//  }
//
//  if (!check_study(c,FALSE)) return FALSE;
//
//  p = rp_skill_ent(c.who, sk);
//
//  /*
//   *  Give him the skill and charge him.
//   *
//   */
//  if (!begin_study(c, sk)) return FALSE;
//  learn_skill(c.who, sk);
//  return TRUE;
//}
//#endif

/*
 *  Thu Jul  3 15:16:22 1997 -- Scott Turner
 *
 *  Practice a skill.
 *
 */
func v_practice(c *command) int {
	sk := c.a
	//days := c.b;

	if numargs(c) < 1 {
		wout(c.who, "Must specify a skill to practice.")
		return FALSE
	}

	if kind(sk) != T_skill {
		wout(c.who, "%s is not a valid skill.", c.parse[1])
		return FALSE
	}

	/*
	 *  They must know the skill.
	 *
	 */
	if FALSE == has_skill(c.who, sk) {
		wout(c.who, "Oddly enough, you must know a skill before practicing it.")
		return FALSE
	}
	/*
	 *  It must be "practiceable"
	 *
	 */
	if FALSE == practice_time(sk) {
		wout(c.who, "Cannot practice '%s'.", c.parse[1])
		return FALSE
	}

	/*
	 *  Need the practice cost.
	 *
	 */
	if practice_cost(sk) > 0 && !charge(c.who, practice_cost(sk)) {
		wout(c.who, "Can't afford %s to practice %s.",
			gold_s(practice_cost(sk)), box_name(sk))
		return FALSE
	}

	/*
	 *  Days to practice?
	 *
	 */
	c.wait = practice_time(sk)

	/*
	 *  Wed May 27 13:09:00 1998 -- Scott Turner
	 *
	 *  It might be progressively harder.
	 *
	 */
	if practice_progressive(sk) != FALSE {
		c.wait += (p_skill_ent(c.who, sk).experience) / practice_progressive(sk)
	}

	wout(c.who, "Practice %s for %s day%s.",
		just_name(sk),
		nice_num(c.wait), add_s(c.wait))

	return TRUE
}

/*
 *  Note:  d_practice is polled daily
 *
 *  Wed May 13 19:00:16 1998 -- Scott Turner
 *
 *  Not any more!
 *
 */
func d_practice(c *command) int {
	sk := c.a

	if kind(sk) != T_skill {
		log_output(LOG_CODE, "d_study: skill %d is gone, who=%d\n", sk, c.who)
		out(c.who, "Internal error: skill %s is gone", box_code(sk))
		return FALSE
	}

	if numargs(c) < 1 {
		wout(c.who, "Must specify a skill to study.")
		return FALSE
	}

	if kind(sk) != T_skill {
		wout(c.who, "%s is not a valid skill.", c.parse[1])
		return FALSE
	}

	/*
	 *  They must know the skill.
	 *
	 */
	if FALSE == has_skill(c.who, sk) {
		wout(c.who, "Oddly enough, you must know a skill before practicing it.")
		return FALSE
	}
	/*
	 *  It must be "practiceable"
	 *
	 */
	if FALSE == practice_time(sk) {
		wout(c.who, "Cannot practice '%s'.", c.parse[1])
		return FALSE
	}

	p := p_skill_ent(c.who, sk)
	p.experience++
	wout(c.who, "You are now level %s in experience for %s.",
		nice_num(p.experience), box_name(sk))
	return TRUE
}

//#if 0
//static int
//research_notknown(who, sk int)
//{
//    static ilist l = nil;
//    var i int
//    var p *entity_skill
//    int new;
//
//    p = rp_skill(sk);
//
//    if (p == nil)
//        return 0;
//
//    ilist_clear(&l);
//
//    for i = 0; i < len(p.research); i++
//    {
//        new = p.research[i];
//
//        if (rp_skill_ent(who, new) == nil &&
//            rp_skill_ent(who, req_skill(new)) != nil)
//        {
//            ilist_append(&l, new);
//        }
//    }
//
//    if (len(l) <= 0)
//        return 0;
//
//    i = rnd(0, len(l)-1);
//
//    return l[i];
//}
//
//
//int
//v_research(c *command)
//{
//  sk := c.a;
//  where := subloc(c.who);
//  int parent, category = 0, teachable = 0, unteachable = 0;
//
//  /*
//   *  In case they're researching from a scroll or book...
//   *
//   */
//  c.b = 0;
//  if (kind(sk) == T_item) sk = correct_study_item(c);
//
//  if (numargs(c) < 1) {
//    wout(c.who, "Must specify skill to research.");
//    return FALSE;
//  }
//
//  if (kind(sk) != T_skill) {
//    wout(c.who, "%s is not a valid skill.", c.parse[1]);
//    return FALSE;
//  }
//
//  if (has_skill(c.who, sk)) {
//    wout(c.who, "%s already knows %s",
//     box_name(c.who), box_name(sk));
//    return FALSE;
//  }
//
//  /*
//   *  What kind of skill are they learning?
//   *
//   */
//  parent = skill_school(sk);
//  category = (parent == sk);
//  if (!category)
//    teachable = (ilist_lookup(rp_skill(parent).offered,sk) != -1);
//  if (!category && !teachable) unteachable = 1;
//
//  if (teachable) {
//    wout(c.who, "%s does not require research.  You may study that skill at any time.", box_name(sk));
//    return FALSE;
//  }
//
//  /*
//   *  If it's a unteachable skill, they have to already know
//   *  the parent skill.
//   *
//   */
//  if (unteachable && FALSE == has_skill(c.who, parent)) {
//    wout(c.who, "Before you can research %s, you must learn %s.",
//     box_name(sk), box_name(parent));
//    return FALSE;
//  }
//
//  if (c.b) {
//    (void) p_skill_ent(c.who, sk);
//    wout(c.who,"You research %s immediately from %s.",
//     box_name(sk), box_name(c.b));
//    consume_scroll(c.who, c.b, 7);
//    return FALSE;  /* Stops the research from continuing. */
//  }
//
//  if (being_taught(c.who, sk)) {
//    (void) p_skill_ent(c.who, sk);
//    wout(c.who,"Your teacher helps you research %s immediately.",
//     box_name(sk));
//    return FALSE;  /* Stops the research from continuing. */
//  }
//
//  if (!can_pay(c.who, 25)) {
//    wout(c.who, "Can't afford 25 gold to research.");
//    return FALSE;
//  }
//
//  if (sk == sk_religion) {
//    if (subkind(where) != sub_temple) {
//      wout(c.who, "%s may only be researched in a temple.", box_name(sk_religion));
//      return FALSE;
//    }
//
//    if (building_owner(where) != c.who) {
//      wout(c.who, "Must be the first character inside the temple to research.");
//      return FALSE;
//    }
//  } else {
//    if (subkind(where) != sub_tower) {
//      wout(c.who, "Research must be performed in a tower.");
//      return FALSE;
//    }
//
//    if (building_owner(where) != c.who) {
//      wout(c.who, "Must be the first character inside the tower to research.");
//      return FALSE;
//    }
//  }
//
//  wout(c.who, "Research %s.", box_name(sk));
//  return TRUE;
//}
//
//
//int
//d_research(c *command)
//{
//  sk := c.a;
//  int new_skill;
//  int chance = 20;
//
//  if (kind(sk) != T_skill) {
//    wout(c.who, "Internal error.");
//    fprintf(os.Stderr, "d_research: skill %d is gone, who=%d\n",
//        sk, c.who);
//    return FALSE;
//  }
//
//  if (!charge(c.who, 25)) {
//    wout(c.who, "Can't afford 25 gold to research.");
//    return FALSE;
//  }
//
//  /*
//   *  10% bonus to research in your strength skill, if any.
//   *
//   */
//  if (is_follower(c.who) &&
//      rp_relig_skill(is_priest(is_follower(c.who))).strength == sk) {
//    chance += 10;
//  }
//
//  /*
//   *  -10% bonus to research in your weakness skill, if any.
//   *
//   */
//  if (is_follower(c.who) &&
//      rp_relig_skill(is_priest(is_follower(c.who))).weakness == sk) {
//    chance -= 10;
//  }
//
//  if (rnd(1,100) > chance) {
//    wout(c.who, "Research is unsuccessful.");
//    return FALSE;
//  }
//
//  /*
//   *  Cause the new skill to be partially known
//   */
//
//  (void) p_skill_ent(c.who, sk);
//
//    wout(c.who, "Research successful for %s",
//     box_name(sk));
//
//    wout(c.who, "To begin learning this skill, order 'study %s'.",
//     box_code_less(sk));
//
//    return TRUE;
//
//#endif

func check_skill_times() {
	for _, sk := range loop_skill() {
		if sk >= 9000 {
			entry := find_use_entry(sk)
			if entry != 0 && skill_time_to_use(sk) != use_tbl[entry].time {
				fprintf(os.Stderr, "Mismatched time to use for %s should be %d.\n", box_code(sk), use_tbl[entry].time)
			}
			if entry != 0 && (skill_flags(sk)&IS_POLLED) != use_tbl[entry].poll {
				fprintf(os.Stderr, "Mismatched poll for %s should be %d.\n", box_code(sk), use_tbl[entry].poll)
			}
		}
	}

}

/*
 *  V_TEACH
 *  Fri Nov  1 16:02:14 1996 -- Scott Turner
 *
 *  Teach -- should be top dog in a tower.  Command takes the skill
 *  to teach as an arg.
 *
 */
func v_teach(c *command) int {
	sk := c.a
	days := c.b
	where := subloc(c.who)

	if numargs(c) < 1 {
		wout(c.who, "Must specify a skill to teach.")
		return FALSE
	}

	if kind(sk) != T_skill {
		wout(c.who, "%s is not a valid skill.", c.parse[1])
		return FALSE
	}

	p := rp_skill_ent(c.who, sk)
	if p == nil || p.know != SKILL_know {
		wout(c.who, "You cannot teach a skill you don't know.")
		return FALSE
	}

	/*
	 *  In the top of a tower?
	 *
	 *  Fri Nov 29 09:56:12 1996 -- Scott Turner
	 *
	 *  Can teach skills in a guild.
	 *
	 */
	if !teachable_place(where, sk) {
		wout(c.who, "Teaching must be performed in a tower or (appropriate) guild.")
		return FALSE
	}

	/*
	 *  Thu Dec 24 08:10:09 1998 -- Scott Turner
	 *
	 *  You no longer have to be the tower owner, but if someone else
	 *  is already teaching we should mention it to you!
	 *
	 */
	for _, i := range loop_stack(where) {
		if i == c.who {
			break
		}
		c2 := rp_command(i)
		if c2 != nil && fmt.Sprintf("%p", cmd_tbl[c2.cmd].start) == fmt.Sprintf("%p", v_teach) {
			// todo: the function comparison hack above is not 100% safe
			wout(c.who, "Someone else is already teaching here.")
			wout(c.who, "Your teaching won't be effective until they stop.")
			break
		}
	}

	if days != 0 {
		c.wait = days
	}

	wout(c.who, "Teaching %s for %d days.", box_name(sk), c.wait)

	/*
	 *  Wed Apr 14 19:33:28 1999 -- Scott Turner
	 *
	 *  Lower the priority of teaching so it finishes *after* any
	 *  "study".
	 *
	 */
	c.pri += 2
	return TRUE

}

func init_use_tbl() {
	if len(use_tbl) != 0 {
		return
	}

	use_tbl = make([]use_tbl_ent, 245)
	use_tbl = append(use_tbl, use_tbl_ent{})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_meditate, v_meditate, d_meditate, nil, 7, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_detect_gates, v_detect_gates, d_detect_gates, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_jump_gate, v_jump_gate, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_teleport, v_teleport, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_seal_gate, v_seal_gate, d_seal_gate, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_unseal_gate, v_unseal_gate, d_unseal_gate, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_notify_unseal, v_notify_unseal, d_notify_unseal, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_rem_seal, v_rem_seal, d_rem_seal, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_reveal_key, v_reveal_key, d_reveal_key, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_notify_jump, v_notify_jump, d_notify_jump, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_rev_jump, v_reverse_jump, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_reveal_mage, v_reveal_mage, d_reveal_mage, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_view_aura, v_view_aura, d_view_aura, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_shroud_abil, v_shroud_abil, d_shroud_abil, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_detect_abil, v_detect_abil, d_detect_abil, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_scry_region, v_scry_region, d_scry_region, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_shroud_region, v_shroud_region, d_shroud_region, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_pr_shroud_loc, v_shroud_region, d_shroud_region, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_detect_scry, v_detect_scry, d_detect_scry, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_dispel_region, v_dispel_region, d_dispel_region, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_dispel_abil, v_dispel_abil, d_dispel_abil, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_adv_med, v_adv_med, d_adv_med, nil, 7, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_hinder_med, v_hinder_med, d_hinder_med, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_proj_cast, v_proj_cast, d_proj_cast, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_locate_char, v_locate_char, d_locate_char, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_bar_loc, v_bar_loc, d_bar_loc, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_unbar_loc, v_unbar_loc, d_unbar_loc, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_forge_palantir, v_forge_palantir, d_forge_palantir, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_destroy_art, v_destroy_art, d_destroy_art, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_save_proj, v_save_proj, d_save_proj, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_save_quick, v_save_quick, d_save_quick, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_quick_cast, v_quick_cast, d_quick_cast, nil, 4, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_basic, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_weather, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_scry, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_gate, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_art, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_necro, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_dirt_golem, v_create_dirt_golem, d_create_dirt_golem, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_flesh_golem, v_create_flesh_golem, d_create_flesh_golem, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_iron_golem, v_create_iron_golem, d_create_iron_golem, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_forge_aura, v_forge_aura, d_forge_aura, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mutate_artifact, v_mutate_art, d_mutate_art, nil, 30, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_conceal_artifacts, v_conceal_arts, d_conceal_arts, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_detect_artifacts, v_detect_arts, d_detect_arts, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_obscure_artifact, v_obscure_art, d_obscure_art, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_remove_obscurity, v_unobscure_art, d_unobscure_art, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_reveal_artifacts, v_reveal_arts, d_reveal_arts, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_deep_identify, v_deep_identify, nil, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_shipbuilding, v_shipbuild, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_pilot_ship, v_sail, d_sail, i_sail, -1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_train_wild, v_use_train_riding, nil, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_train_warmount, v_use_train_war, nil, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_make_ram, nil, nil, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_make_catapult, nil, nil, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_make_siege, nil, nil, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_brew_slave, v_brew, d_brew_slave, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_brew_heal, v_brew, d_brew_heal, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_brew_death, v_brew, d_brew_death, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_brew_weightlessness, v_brew, d_brew_weightlessness, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_brew_fiery, v_brew, d_brew_fiery, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mine_iron, v_mine_iron, d_mine_iron, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mine_gold, v_mine_gold, d_mine_gold, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mine_mithril, v_mine_mithril, d_mine_mithril, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mine_crystal, v_mine_gate_crystal, d_mine_gate_crystal, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_quarry_stone, v_quarry, nil, nil, -1, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_catch_horse, v_catch, nil, nil, -1, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_extract_venom, nil, nil, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_harvest_lumber, v_wood, nil, nil, -1, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_harvest_yew, v_yew, nil, nil, -1, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_add_ram, v_add_ram, d_add_ram, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_remove_ram, v_remove_ram, d_remove_ram, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_assassinate, v_assassinate, d_assassinate, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_find_food, v_find_food, d_find_food, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_spy_inv, v_spy_inv, d_spy_inv, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_spy_skills, v_spy_skills, d_spy_skills, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_spy_lord, v_spy_lord, d_spy_lord, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_record_skill, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_bribe_noble, v_bribe, d_bribe, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_summon_savage, v_summon_savage, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_keep_savage, v_keep_savage, d_keep_savage, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_opium, v_improve_opium, d_improve_opium, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_raise_mob, v_raise, d_raise, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_rally_mob, v_rally, d_rally, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_incite_mob, v_incite, d_incite, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_bird_spy, v_bird_spy, d_bird_spy, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_lead_to_gold, v_lead_to_gold, d_lead_to_gold, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_raise_corpses, v_raise_corpses, nil, nil, -1, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_undead_lord, v_undead_lord, d_undead_lord, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_banish_undead, v_banish_undead, d_banish_undead, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_renew_undead, v_keep_undead, d_keep_undead, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_eat_dead, v_eat_dead, d_eat_dead, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_aura_blast, v_aura_blast, d_aura_blast, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_absorb_blast, v_aura_reflect, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_summon_rain, v_summon_rain, d_summon_rain, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_summon_wind, v_summon_wind, d_summon_wind, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_summon_fog, v_summon_fog, d_summon_fog, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_direct_storm, v_direct_storm, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_renew_storm, v_renew_storm, d_renew_storm, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_dissipate, v_dissipate, d_dissipate, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_lightning, v_lightning, d_lightning, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_fierce_wind, v_fierce_wind, d_fierce_wind, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_seize_storm, v_seize_storm, d_seize_storm, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_death_fog, v_death_fog, d_death_fog, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_banish_corpses, v_banish_corpses, d_banish_corpses, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_hide_self, v_hide, d_hide, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_conceal_nation, v_conceal_nation, d_conceal_nation, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_sneak_build, v_sneak, d_sneak, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mage_menial, v_mage_menial, d_mage_menial, nil, 7, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_petty_thief, v_petty_thief, d_petty_thief, nil, 7, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_appear_common, v_appear_common, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_defense, v_defense, d_defense, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_defense2, v_defense, d_defense, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_archery, v_archery, d_archery, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_swordplay, v_swordplay, d_swordplay, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_swordplay2, v_swordplay, d_swordplay, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_find_rich, v_find_rich, d_find_rich, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_harvest_opium, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_train_angry, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_weaponsmith, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_hide_lord, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_transcend_death, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_collect_foliage, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_fishing, v_fish, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_summon_ghost, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_capture_beasts, v_capture_beasts, d_capture_beasts, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_use_beasts, v_use_beasts, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_collect_elem, v_implicit, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_torture, v_torture, d_torture, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_fight_to_death, v_fight_to_death, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_breed_beasts, v_breed, d_breed, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_persuade_oath, v_persuade_oath, d_persuade_oath, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_forge_weapon, v_forge_art_x, d_forge_art_x, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_forge_armor, v_forge_art_x, d_forge_art_x, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_forge_bow, v_forge_art_x, d_forge_art_x, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_trance, v_trance, d_trance, nil, 28, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_teleport_item, v_teleport_item, d_teleport_item, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_tap_health, v_tap_health, d_tap_health, nil, 7, 0})
	//use_tbl = append(use_tbl, use_tbl_ent{"c", sk_bind_storm,     v_bind_storm,     d_bind_storm,     nil,     7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_control_battle, v_prac_control, nil, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_protect_noble, v_prac_protect, nil, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_attack_tactics, v_attack_tactics, nil, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_defense_tactics, v_defense_tactics, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_resurrect, v_resurrect, d_resurrect, nil, 30, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_pray, v_prep_ritual, d_prep_ritual, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_last_rites, v_last_rites, d_last_rites, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_gather_holy_plant, v_gather_holy_plant, d_gather_holy_plant, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_bless_follower, v_bless_follower, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_proselytise, v_proselytise, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_create_holy_symbol, v_create_holy_symbol, d_create_holy_symbol, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_heal, v_heal, d_heal, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_summon_water_elemental, v_generic_trap, d_generic_trap, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_find_mtn_trail, v_find_mountain_trail, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_obscure_mtn_trail, v_obscure_mountain_trail, d_obscure_mountain_trail, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_mining, v_improve_mining, d_improve_mining, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_conceal_mine, v_conceal_mine, d_conceal_mine, nil, 30, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_protect_mine, v_protect_mine, d_protect_mine, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_bless_fort, v_bless_fort, d_bless_fort, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_weaken_fort, v_weaken_fort, d_weaken_fort, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_boulder_trap, v_generic_trap, d_generic_trap, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_anteus, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_detect_beasts, v_detect_beasts, d_detect_beasts, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_capture_beasts, v_capture_beasts, d_capture_beasts, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_snake_trap, v_generic_trap, d_generic_trap, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_dol, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_find_forest_trail, v_find_forest_trail, nil, nil, 1, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_obscure_forest_trail, v_obscure_forest_trail, d_obscure_forest_trail, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_forestry, v_improve_logging, d_improve_logging, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_reveal_forest, v_find_hidden_features, d_find_hidden_features, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_fort, v_improve_fort, d_improve_fort, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_create_deadfall, v_generic_trap, d_generic_trap, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_recruit_elves, v_recruit_elves, d_recruit_elves, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_timeid, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_reveal_vision, v_reveal_vision, d_reveal_vision, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_enchant_guard, v_enchant_guard, d_enchant_guard, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_urchin_spy, v_urchin_spy, d_urchin_spy, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_draw_crowds, v_draw_crowds, d_draw_crowds, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_arrange_mugging, v_arrange_mugging, d_arrange_mugging, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_ham, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_quarry, v_improve_quarrying, d_improve_quarrying, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_smithing, v_improve_smithing, d_improve_smithing, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_edge_of_kireus, v_edge_of_kireus, d_edge_of_kireus, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_create_mithril, v_create_mithril, d_create_mithril, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_quicksand_trap, v_generic_trap, d_generic_trap, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_kireus, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_calm_ap, nil, d_calm_peasants, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_charisma, v_improve_charisma, d_improve_charisma, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mesmerize_crowd, nil, d_mesmerize_crowd, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_improve_taxes, nil, d_improve_taxes, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_guard_loyalty, nil, d_guard_loyalty, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_instill_fanaticism, nil, d_instill_fanaticism, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_halon, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_find_hidden, nil, d_find_all_hidden_features, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_conceal_loc, v_conceal_location, d_conceal_location, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_create_ninja, v_create_ninja, d_create_ninja, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_mists_of_conceal, nil, d_create_mist, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_domingo, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_dedicate_temple, v_dedicate_temple, d_dedicate_temple, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_smuggle_goods, v_smuggle_goods, d_smuggle_goods, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_smuggle_men, v_smuggle_men, d_smuggle_men, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_build_wagons, nil, nil, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_increase_demand, v_increase_demand, d_increase_demand, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_decrease_demand, v_decrease_demand, d_decrease_demand, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_increase_supply, v_increase_supply, d_increase_supply, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_decrease_supply, v_decrease_supply, d_decrease_supply, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_hide_money, v_hide_money, d_hide_money, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_hide_item, v_hide_item, d_hide_item, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_grow_pop, v_grow_pop, d_grow_pop, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_add_sails, v_add_sails, d_add_sails, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_remove_sails, v_remove_sails, d_remove_sails, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_add_forts, v_add_forts, d_add_forts, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_remove_forts, v_remove_forts, d_remove_forts, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_add_ports, v_add_ports, d_add_ports, nil, 4, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_remove_ports, v_remove_ports, d_remove_ports, nil, 4, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_add_keels, v_add_keels, d_add_keels, nil, 5, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_remove_keels, v_remove_keels, d_remove_keels, nil, 5, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_fortify_castle, v_fortify_castle, d_fortify_castle, nil, 5, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_strengthen_castle, v_strengthen_castle, d_strengthen_castle, nil, 5, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_moat_castle, v_moat_castle, d_moat_castle, nil, 5, 1})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_widen_entrance, v_widen_entrance, d_widen_entrance, nil, 5, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_wooden_shoring, v_add_wooden_shoring, d_add_wooden_shoring, nil, 5, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_iron_shoring, v_add_iron_shoring, d_add_iron_shoring, nil, 5, 0})

	// Combat Spells
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_lightning_bolt, v_use_cs, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_foresee_defense, v_use_cs, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_drain_mana, v_use_cs, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_raise_soldiers, v_use_cs, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_fireball, v_use_cs, nil, nil, 0, 0})

	// Heroism
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_personal_fttd, v_personal_fight_to_death, nil, nil, 0, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_forced_march, v_forced_march, nil, nil, 0, 0})

	// Basic religion
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_resurrect_b, v_resurrect, d_resurrect, nil, 30, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_pray_b, v_prep_ritual, d_prep_ritual, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_last_rites_b, v_last_rites, d_last_rites, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_gather_holy_plant_b, v_gather_holy_plant, d_gather_holy_plant, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_bless_b, v_bless_follower, nil, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_proselytise_b, v_proselytise, nil, nil, 3, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_create_holy_b, v_create_holy_symbol, d_create_holy_symbol, nil, 14, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_heal_b, v_heal, d_heal, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_dedicate_temple_b, v_dedicate_temple, d_dedicate_temple, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_write_religion_b, v_write_spell, d_write_spell, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_hinder_med_b, v_hinder_med_b, d_hinder_med_b, nil, 10, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_scry_b, v_vision_reg, d_vision_reg, nil, 7, 0})
	use_tbl = append(use_tbl, use_tbl_ent{"c", sk_banish_undead_b, v_banish_undead, d_banish_undead, nil, 7, 0})

	use_tbl = append(use_tbl, use_tbl_ent{})
}
