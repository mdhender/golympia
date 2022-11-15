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

import "strings"

func keep_undead_check(c *command, check_bond int) int {
	target := c.a

	if kind(target) != T_char || subkind(target) != sub_demon_lord {
		wout(c.who, "%s is not a demon lord.", box_code(target))
		return FALSE
	}

	if subloc(target) != subloc(c.who) {
		wout(c.who, "%s is not here.", box_code(target))
		return FALSE
	}

	if check_bond != 0 && loyal_kind(target) != LOY_summon {
		wout(c.who, "%s is no longer bonded.", box_code(target))
		return FALSE
	}

	return TRUE
}

func v_keep_undead(c *command) int {
	//target := c.a;

	if keep_undead_check(c, TRUE) == 0 {
		return FALSE
	}

	if check_aura(c.who, 3) == 0 {
		return FALSE
	}

	return TRUE
}

func d_keep_undead(c *command) int {
	target := c.a

	if keep_undead_check(c, TRUE) == 0 {
		return FALSE
	}

	if charge_aura(c.who, 3) == 0 {
		return FALSE
	}

	set_loyal(target, LOY_summon, max(loyal_rate(target)+4, 8))

	wout(c.who, "%s will remain for %d months.",
		box_code(target),
		loyal_rate(target))
	return TRUE
}

func v_undead_lord(c *command) int {
	where := subloc(c.who)
	aura := c.a

	if aura < 3 {
		c.a, aura = 3, 3
	}
	if aura > 8 {
		c.a, aura = 8, 8
	}

	if may_cookie_npc(c.who, where, item_undead_cookie) == 0 {
		return FALSE
	}

	if check_aura(c.who, aura) == 0 {
		return FALSE
	}

	return TRUE
}

func d_undead_lord(c *command) int {
	where := subloc(c.who)
	aura := c.a

	if may_cookie_npc(c.who, where, item_undead_cookie) == 0 {
		return FALSE
	}

	if charge_aura(c.who, aura) == 0 {
		return FALSE
	}

	undead := do_cookie_npc(c.who, where, item_undead_cookie, c.who)

	if undead == 0 {
		log_output(LOG_CODE, "d_undead_lord: why not?")
		wout(c.who, "Unable to summon a demon lord.")
		return FALSE
	}

	var rating int
	switch aura {
	case 3:
		rating = 100
		break
	case 4:
		rating = 150
		break
	case 5:
		rating = 190
		break
	case 6:
		rating = 220
		break
	case 7:
		rating = 240
		break
	case 8:
		rating = 250
		break

	default:
		assert(false)
	}

	p_char(undead).attack = rating
	p_char(undead).defense = rating

	set_loyal(undead, LOY_summon, 5)

	wout(c.who, "Summoned %s.", box_name(undead))
	wout(where, "%s has summoned %s.",
		box_name(c.who),
		liner_desc(undead))

	return TRUE
}

func v_banish_undead(c *command) int {

	if keep_undead_check(c, FALSE) == 0 {
		return FALSE
	}

	if check_aura(c.who, 6) == 0 {
		return FALSE
	}

	if cast_check_char_here(c.who, c.a) == 0 {
		return FALSE
	}

	return TRUE
}

func d_banish_undead(c *command) int {
	target := c.a
	where := subloc(c.who)

	if keep_undead_check(c, FALSE) == 0 {
		return FALSE
	}

	if charge_aura(c.who, 6) == 0 {
		return FALSE
	}

	head := stack_leader(target)

	wout(head, "%s banishes %s!", box_name(c.who), box_name(target))
	wout(where, "%s banishes %s!", box_name(c.who), box_name(target))

	extract_stacked_unit(target)
	kill_char(target, 0, S_body)

	return TRUE
}

func v_eat_dead(c *command) int {
	body := c.a

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

	if kind(body) == T_item {
		if subkind(body) != sub_dead_body {
			wout(c.who, "%s is not a dead body.", box_code(body))
			return FALSE
		}
		if has_item(c.who, body) == 0 {
			wout(c.who, "You do not possess that body.")
			return FALSE
		}
	} else if kind(body) == T_char {
		if subkind(body) != sub_lost_soul {
			wout(c.who, "%s is not a lost soul.", box_code(body))
			return FALSE
		}
		if has_prisoner(c.who, body) == 0 {
			wout(c.who, "You have not captured that lost soul.")
			return FALSE
		}
	}

	return TRUE
}

func get_some_skills(who, body, chance int) {
	var parent int

	/*
	 *  First do category skills
	 */

	for _, e := range loop_char_skill_known(body) {
		parent = skill_school(e.skill)
		if parent != e.skill {
			continue
		}

		if has_skill(who, e.skill) != FALSE {
			continue
		}

		if e.skill == sk_adv_sorcery {
			continue
		}

		/*
		 *  Fri Sep 20 12:49:59 1996 -- Scott Turner
		 *
		 *  Can't learn religions by eating the dead!
		 *
		 */
		if rp_relig_skill(e.skill) != nil {
			continue
		}

		if rnd(1, 100) > chance {
			continue
		}

		learn_skill(who, e.skill)
	}

	/*
	 *  Now do subskills
	 *  Must know parent in order to pick up a subskill
	 */

	for _, e := range loop_char_skill_known(body) {
		parent = skill_school(e.skill)
		if parent == e.skill {
			continue
		}

		if has_skill(who, e.skill) != FALSE || has_skill(who, parent) == FALSE {
			continue
		}

		if rnd(1, 100) > chance {
			continue
		}

		learn_skill(who, e.skill)
	}
}

func d_eat_dead(c *command) int {
	body := c.a

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

	if kind(body) == T_item {
		if subkind(body) != sub_dead_body {
			wout(c.who, "%s is not a dead body.", box_code(body))
			return FALSE
		}
		if has_item(c.who, body) == FALSE {
			wout(c.who, "You do not possess that body.")
			return FALSE
		}
	} else if kind(body) == T_char {
		if subkind(body) != sub_lost_soul {
			wout(c.who, "%s is not a lost soul.", box_code(body))
			return FALSE
		}
		if has_prisoner(c.who, body) == FALSE {
			wout(c.who, "You have not captured that lost soul.")
			return FALSE
		}
	}

	if charge_aura(c.who, 5) == 0 {
		return FALSE
	}

	wout(c.who, "Consumed %s.", box_name(body))
	get_some_skills(c.who, body, 100)

	pl := body_old_lord(body)

	if valid_box(pl) {
		out(pl, "The spirit of %s~%s has been defiled by Necromancy.",
			rp_misc(body).save_name, box_code(body))

		rp_misc(body).old_lord = 0 /* inhibit NP return */
	}

	/*
	 *  If the body was a priest, well, that's good for your mana!
	 *
	 */
	if options.mp_antipathy != FALSE && is_priest(body) != FALSE && rp_magic(c.who) != nil {
		//p := p_magic(c.who);
		wout(c.who, "Your mana grows stronger on the soul of a priest!")
		add_aura(c.who, 15)
	}

	dead_body_np = FALSE
	kill_char(body, 0, S_nothing)
	dead_body_np = TRUE

	if rnd(1, 100) <= 25 && char_sick(c.who) == FALSE && has_artifact(c.who, ART_SICKNESS, 0, 0, 0) == FALSE {
		p_char(c.who).sick = TRUE
		wout(c.who, "%s has fallen ill.", box_name(c.who))
	}

	return TRUE
}

func random_body_here(where int) int {
	var l []int

	for _, e := range inventory_loop(where) {
		if subkind(e.item) == sub_dead_body &&
			sysclock.turn > p_char(e.item).death_time.turn {
			l = append(l, e.item)
		}
	}

	if len(l) == 0 {
		return 0
	}

	ilist_scramble(l)

	return l[0]
}

////#if 0
//int
//v_exhume(c *command)
//{
//    where := subloc(c.who);
//    int targ = c.a;
//    int n;
//
//    if (subkind(where) != sub_graveyard)
//    {
//        wout(c.who, "Bodies may only be exhumed in graveyards.");
//        return FALSE;
//    }
//
//    if (targ &&
//        (!valid_box(targ) ||
//        subkind(targ) != sub_dead_body ||
//        has_item(where, targ) == 0))
//    {
//        wout(c.who, "No body %s is buried here.", box_code(targ));
//        return FALSE;
//    }
//
//    if (!targ && (n = random_body_here(where)) == 0)
//    {
//        wout(c.who, "There are no fresh graves here to dig up.");
//        return FALSE;
//    }
//
//    if (!targ)
//        targ = n;
//
//    if (sysclock.turn == p_char(targ).death_time.turn)
//    {
//        wout(c.who, "%s may not be exhumed until next month.",
//                        cap_(box_name(targ)));
//        return FALSE;
//    }
//
//    return TRUE;
//}
//
//
//int
//d_exhume(c *command)
//{
//    where := subloc(c.who);
//    int targ = c.a;
//    int n;
//
//    if (subkind(where) != sub_graveyard)
//    {
//        wout(c.who, "Bodies may only be exhumed in graveyards.");
//        return FALSE;
//    }
//
//    if (targ &&
//        (!valid_box(targ) ||
//        subkind(targ) != sub_dead_body ||
//        has_item(where, targ) == 0))
//    {
//        wout(c.who, "No body %s is buried here.", box_code(targ));
//        return FALSE;
//    }
//
//    if (!targ && (n = random_body_here(where)) == 0)
//    {
//        wout(c.who, "There are no fresh graves here to dig up.");
//        return FALSE;
//    }
//
//    if (!targ)
//        targ = n;
//
//    if (sysclock.turn == p_char(targ).death_time.turn)
//    {
//        wout(c.who, "%s may not be exhumed until next month.",
//                        cap_(box_name(targ)));
//        return FALSE;
//    }
//
//    move_item(where, c.who, targ, 1);
//
//    wout(c.who, "Exhumed %s.", box_name(targ));
//    wout(where, "%s exhumed %s.", box_name(c.who), box_name(targ));
//
//    return TRUE;
//}
////#endif

func auto_undead(who int) {
	var master int
	where := subloc(who)

	master = npc_summoner(who)

	if master != 0 && subloc(who) == subloc(master) {
		queue(who, "attack %s", box_code_less(master))
		p_misc(who).summoned_by = 0
		return
	}

	if loc_depth(where) != LOC_province || rnd(1, 2) == 1 {
		npc_move(who)
		return
	}

	queue(who, "pillage 1")
}

func v_aura_blast(c *command) int {
	//target := c.a;
	//aura := c.b;
	//have_left := c.c;
	where := subloc(c.who)

	if in_safe_now(where) != FALSE {
		wout(c.who, "Not allowed in a safe haven.")
		return FALSE
	}

	return TRUE
}

func d_aura_blast(c *command) int {
	target := c.a
	aura := c.b
	have_left := c.c
	where := subloc(c.who)
	has_protection := 0

	if cast_check_char_here(c.who, target) == FALSE {
		return FALSE
	}

	if in_safe_now(where) != FALSE || in_safe_now(target) != FALSE {
		wout(c.who, "Not allowed in a safe haven.")
		return FALSE
	}

	if aura < 1 {
		aura = char_cur_aura(c.who)
	}

	if have_left != 0 {
		m := max(char_cur_aura(c.who)-have_left, 0)

		if aura > m {
			aura = m
		}
	}

	if aura == 0 {
		wout(c.who, "No aura available for blast.")
		return FALSE
	}

	if charge_aura(c.who, aura) == FALSE {
		return FALSE
	}

	vector_clear()
	vector_add(c.who)
	vector_add(target)
	vector_add(where)

	wout(VECT, "%s blasts %s with a burst of aura!",
		box_name(c.who), box_name(target))

	log_output(LOG_SPECIAL, "%s blasts %s with a burst of aura!",
		box_name(c.who), box_name(target))

	/*
	 *  Wed Sep 30 13:18:37 1998 -- Scott Turner
	 *
	 *  A priest might have a "Protection from Aura Blast" prayer.
	 *
	 */
	for _, e := range loop_char_skill(target) {
		if strings.HasPrefix(strings.ToLower(bx[e.skill].name), "protection from aura blast") && e.know == SKILL_know {
			has_protection = 1
			break
		}
	}

	if has_skill(target, sk_absorb_blast) != FALSE || has_skill(target, sk_prot_blast_b) != FALSE {
		if reflect_blast(target) != FALSE {
			wout(VECT, "%s reflected the blast back to %s!",
				just_name(target), just_name(c.who))

			add_char_damage(c.who, aura*2, MATES)
		} else {
			wout(VECT, "%s absorbed the blast!", just_name(target))

			add_aura(target, aura/2)
			wout(target, "Current aura is now %d.",
				rp_magic(target).cur_aura)
		}
	} else if has_protection != FALSE {
		wout(VECT, "The blast dissipates harmlessly against %s's holy aura!",
			just_name(target))
	} else if has_artifact(target, ART_PROT_SKILL, sk_aura_blast, 0, 0) != FALSE {
		wout(target, "%s briefly gives off a blinding light.")
	} else {
		add_char_damage(target, aura*2, MATES)
	}
	return TRUE
}

func v_aura_reflect(c *command) int {
	flag := c.a

	p_magic(c.who).aura_reflect = flag

	if flag != FALSE {
		wout(c.who, "Will reflect aura blasts back at the attacker.")
	} else {
		wout(c.who, "Will absorb aura blasts.")
	}

	return TRUE
}

/*
 *  Mon May  5 12:25:36 1997 -- Scott Turner
 *
 *  Create flesh golem.  Decays after one year.
 *
 */
func v_create_flesh_golem(c *command) int {
	body := c.a

	if kind(body) != T_item || subkind(body) != sub_dead_body {
		wout(c.who, "%s is not the dead body of a noble.",
			box_code(body))
		return FALSE
	}

	if has_item(c.who, body) == 0 {
		wout(c.who, "Don't have %s.", box_code(body))
		return FALSE
	}

	wout(c.who, "Begin construction of a flesh golem.")
	return TRUE
}

func d_create_flesh_golem(c *command) int {
	body := c.a

	if charge_aura(c.who, skill_piety(c.use_skill)) == FALSE {
		return FALSE
	}

	if kind(body) != T_item || subkind(body) != sub_dead_body {
		wout(c.who, "%s is not the dead body of a noble.",
			box_code(body))
		return FALSE
	}

	if has_item(c.who, body) == 0 {
		wout(c.who, "Don't have %s.", box_code(body))
		return FALSE
	}

	/*
	 *  Add an effect to destroy this golem in a year.
	 *
	 */
	if add_effect(c.who, ef_kill_flesh_golem, 0, 360+rnd(1, 120), 1) == FALSE {
		wout(c.who, "For some reason, the spell fails to take effect.")
		return FALSE
	}

	gen_item(c.who, item_flesh_golem, 1)
	wout(c.who, "You have created a flesh golem.")
	dead_body_np = TRUE
	destroy_unique_item(c.who, body)
	return TRUE
}
