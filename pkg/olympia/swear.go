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

import "sort"

func swear_char(who, target int) { panic("!implemented") }

// lord returns the id of the lord of the character passed in
func lord(n int) int {
	if kind(n) == T_player { // lord of a player is the player
		return n
	}
	if rp_char(n) != nil {
		return rp_char(n).unit_lord
	}
	return 0 // item has no lord
}

// player returns the player that controls the character or item.
// it searches up through the hierarchy to find the player.
// returns 0 if the character isn't controlled by a player.
func player(item int) int {
	//#if 0
	//    int loop_check = 5000;
	//#endif

	for item != 0 && kind(item) != T_player {
		item = lord(item)
		//#if 0
		//        assert(loop_check-- > 0);
		//#endif
	}

	return item
}

func set_loyal(who int, k int, lev int) {
	p := p_char(who)
	p.loy_kind = k
	p.loy_rate = lev
}

/*
 *  is b sworn somewhere beneath a?
 */
func sworn_beneath(a, b int) bool {
	if a == b { // can't swear loyalty to self
		return false
	}
	for kind(b) != T_player {
		if b = lord(b); a == b {
			return true
		}
	}
	return false
}

func reswear_all_sworn(who int) {
	new_lord := lord(who)
	for _, i := range loop_units(player(who)) {
		if lord(i) == who {
			set_lord(i, new_lord, LOY_UNCHANGED, 0)
		}
	}
}

func set_lord(who int, new_lord int, k int, lev int) {
	new_pl, old_pl := 0, player(who)
	if new_lord != 0 {
		new_pl = player(new_lord)
	}

	if old_pl != 0 && old_pl != new_pl {
		flush_unit_orders(old_pl, who)
		interrupt_order(who)
		clear_all_att(who)
		reswear_all_sworn(who)

		//#if 0
		//        if (char_pledge(who))
		//            p_magic(who).pledge = 0;
		//#endif

		p_player(old_pl).units = rem_value(p_player(old_pl).units, who)
	}

	p_char(who).unit_lord = new_lord
	p_misc(who).old_lord = old_pl

	if new_lord != 0 && new_pl != old_pl {
		p_player(new_pl).units = append(p_player(new_pl).units, who)
		sort.Ints(p_player(new_pl).units)
		init_load_sup(who) /* load command from new owner */
		touch_loc(who)
	}

	// todo: figure out who pops out of what stack
	//       perhaps they don't?  Just ignore it?
	//       no auto-unstack?  Have to manually force them out?
	if k != LOY_UNCHANGED {
		set_loyal(who, k, lev)
	}
}

func v_swear(c *command) int {
	out(c.who, "The SWEAR order has been removed from the game.")
	return FALSE

	//target := c.a;
	//int targ_lord, targ_pl;
	//int old_lord, old_pl;
	//
	//
	//if (target == 0 && numargs(c) > 0 && *c.parse[1] == '0') {
	//    unit_deserts(c.who, deserted_player, FALSE, LOY_unsworn, 0);
	//    return TRUE;
	//}
	//
	//if (numargs(c) <= 0) {
	//    wout(c.who, "Must specify a character to swear fealty to.");
	//    return FALSE;
	//}
	//
	//if (FALSE == check_char_here(c.who, target)) {
	//    return FALSE;
	//}
	//
	//if (sworn_beneath(c.who, target)) {
	//    wout(c.who, "Cannot swear to a character beneath "
	//                 "you in the command hierarchy.");
	//    return FALSE;
	//}
	//
	//old_lord = lord(c.who);
	//old_pl = player(old_lord);
	//
	//targ_lord = lord(target);
	//targ_pl = player(targ_lord);
	//
	//if (old_lord == targ_lord) {
	//    wout(c.who, "Already sworn to %s.", box_name(old_lord));
	//    return FALSE;
	//}
	//
	//if (rp_player(old_pl).swear_this_turn) {
	//    wout(c.who, "Allowed at most one SWEAR per turn.");
	//    return FALSE;
	//}
	//
	//rp_player(old_pl).swear_this_turn++;
	//
	//if (old_pl != targ_pl) {
	//    wout(old_pl, "%s renounces loyalty.", box_name(c.who));
	//    wout(targ_pl, "%s swears loyalty.", box_name(c.who));
	//}
	//
	//wout(target, "%s swears loyalty to us.", box_name(c.who));
	//
	//{
	//    show_day = FALSE;
	//    out(c.who, "");
	//    show_day = TRUE;
	//}
	//
	//set_lord(c.who, targ_lord, LOY_UNCHANGED, 0);
	//
	//return TRUE;
}

func is_unit(pl, v int) bool {
	if !(kind(pl) == T_player) {
		panic("assert(kind(pl) == T_player)")
	}
	return ilist_lookup(p_player(pl).units, v) >= 0
}

func unit_deserts(who int, to_who int, loy_check int, k int, lev int) bool {
	must_pay := false
	sp := player(who)
	wout(gm_player, "Attempted desertion %s to %s.",
		box_name(who), box_name(to_who))

	// first calculate the real player of "who".
	// if he's a deserted noble, then his real player is body_old_lord.
	// if he's a refugee, then his body_old_lord == 0
	if is_real_npc(sp) && body_old_lord(who) != 0 {
		sp = body_old_lord(who)
	}

	// now determine if we must pay NP transfer cost.

	// if it is one nation to another, then you don't pay, so long as neither nation is neutral.
	if sp != 0 && nation(sp) != 0 && !rp_nation(nation(sp)).neutral && to_who != 0 && nation(to_who) != 0 && !rp_nation(nation(to_who)).neutral && nation(sp) != nation(to_who) {
		must_pay = true
	}

	// don't pay if the sp has no nation.
	// this is the case for "old men" prisoners and so on.
	if is_real_npc(sp) && nation(sp) == 0 {
		must_pay = false
	}

	// the only other time we don't pay is if we're reclaiming our own.
	if nation(sp) == nation(to_who) {
		must_pay = false
	}

	// don't make NPCs pay.
	if to_who == 0 || is_real_npc(to_who) {
		must_pay = false
	}

	if must_pay {
		num := nps_invested(who)
		if deduct_np(to_who, num) == 0 {
			return false
		}
		if num != 0 {
			wout(to_who, "Paying %s noble points for %s.",
				nice_num(num), box_name(who))
		}
		// the player who is losing the noble receives those points back only if he is losing the noble to someone from his nation.
		add_np(sp, num)
		if num != 0 {
			wout(sp, "Receiving %s noble points for the loss of %s.", nice_num(num), box_name(who))
		}
	}
	// wipe out proscribed skills.
	if to_who != 0 && !is_real_npc(to_who) && nation(to_who) != 0 {
		for _, e := range loop_char_skill_known(who) {
			if ilist_lookup(rp_nation(nation(to_who)).proscribed_skills, e.skill) != -1 ||
				ilist_lookup(rp_nation(nation(to_who)).proscribed_skills, skill_school(e.skill)) != -1 {
				forget_skill(who, e.skill)
			}
		}
	}

	/*
	 *  Mon Jan 22 18:14:56 2001 -- Scott Turner
	 *
	 *  Might affect your aura/piety.
	 *
	 */
	if char_cur_aura(who) > char_max_aura(who) {
		p_magic(who).cur_aura = char_max_aura(who)
	}
	if is_priest(who) == FALSE && char_piety(who) != FALSE {
		p_char(who).religion.piety = 0
	}

	if to_who != 0 && sp != 0 {
		wout(sp, "%s renounces loyalty to us.", box_name(who))
		wout(who, "%s renounces loyalty.", box_name(who))
		show_day = false
		out(who, "")
		show_day = false
	}

	// if a prisoner deserts to the faction of the unit holding it prisoner, don't extract it from the stack.
	// instead, simply clear the prisoner bit.
	if to_who != 0 && is_prisoner(who) && player(to_who) == player(stack_parent(who)) {
		p_char(who).prisoner = FALSE
	} else if !is_prisoner(who) {
		extract_stacked_unit(who)
	}

	set_lord(who, to_who, k, lev)

	if to_who != 0 {
		wout(who, "%s pledges fealty to us.", box_name(who))
		wout(to_who, "%s pledges fealty to us.", box_name(who))

	}
	p_char(who).new_lord = 1

	return true
}

func v_bribe(c *command) int {
	target := c.a
	amount := c.b

	if has_skill(c.who, sk_bribe_noble) < 1 {
		wout(c.who, "BRIBE requires knowledge of %s.",
			cap_(box_name(sk_bribe_noble)))
		return FALSE
	}

	if FALSE == check_char_here(c.who, target) {
		return FALSE
	}

	if char_new_lord(target) != FALSE {
		wout(c.who, "%s just switched employers this month, and is not looking for a new one so soon.", box_name(target))
		return FALSE
	}

	if is_npc(target) {
		wout(c.who, "NPC's cannot be bribed.")
		return FALSE
	}

	if player(target) == player(c.who) {
		wout(c.who, "%s already belongs to our faction.",
			box_name(target))
		return FALSE
	}

	if amount == 0 {
		wout(c.who, "Must specify an amount of gold to use as a bribe.")
		return FALSE
	}

	if FALSE == can_pay(c.who, amount) {
		wout(c.who, "Don't have %s for a bribe.", gold_s(amount))
		return FALSE
	}

	wout(c.who, "Attempt to bribe %s with a gift of %s.",
		box_name(target), gold_s(amount))

	return TRUE
}

func thanks_for_gift(who, target int) {
	switch rnd(1, 3) {
	case 1:
		wout(who, "%s graciously accepts our gift.",
			box_name(target))
		break

	case 2:
		wout(who, "%s thanks us for the gift.",
			box_name(target))
		break

	case 3:
		wout(who, "%s pockets the gold.", box_name(target))
		break

	default:
		panic("!reached")
	}
}

/*
 *	over threshold			under threshold
 *	--------------			---------------
 *	35%	switch			50% pocket
 *	30%	pocket			50% report
 *	25%	report bribe
 *	10%	go independent
 */
const (
	SWITCH         = 1
	POCKET         = 2
	REPORT         = 3
	HEAD_FOR_HILLS = 4
)

func d_bribe(c *command) int {
	target := c.a
	amount := c.b
	flag := c.c != FALSE
	bribe_thresh := 0
	var outcome int

	if FALSE == check_still_here(c.who, target) {
		return FALSE
	}

	if char_new_lord(c.who) != FALSE {
		wout(c.who, "%s just switched employers this month, and is not looking for a new one so soon.", box_name(target))
		return FALSE
	}

	if FALSE == charge(c.who, amount) {
		wout(c.who, "Don't have %s for a bribe.", gold_s(amount))
		return FALSE
	}

	switch loyal_kind(target) {
	case LOY_unsworn, LOY_contract:
		bribe_thresh = loyal_rate(target)
		if bribe_thresh < 249 {
			bribe_thresh = 249
		}
	case LOY_fear:
		bribe_thresh = 249
	case LOY_oath:
		break
	default:
		panic("!reached")
	}

	//#if 0
	//    /*
	//     *  Special case -- if you are a "refugee" then
	//     *  you can be bribed for 249 in any case.
	//     *
	//     *  Thu Dec 23 06:24:22 1999 -- Scott Turner
	//     *
	//     *  Nah, this causes more problems than it is worth.  For
	//     *  one thing, you get charged the NP for the oath-bound
	//     *  noble, but he ends up in your faction with contract
	//     *  loyalty.
	//     *
	//     */
	//    if (refugee(target)) bribe_thresh = 249;
	//#endif

	if get_effect(target, ef_guard_loyalty, 0, 0) != FALSE {
		outcome = REPORT
	} else if bribe_thresh <= 0 || amount < bribe_thresh {
		if rnd(1, 2) == 1 {
			outcome = POCKET
		} else {
			outcome = REPORT
		}
	} else if amount == bribe_thresh {
		outcome = HEAD_FOR_HILLS
	} else {
		outcome = SWITCH
		//#if 0
		//        int n = rnd(1,100);
		//
		//        if (n <= 35)
		//            outcome = SWITCH;
		//        else if (n <= 65)
		//            outcome = POCKET;
		//        else if (n <= 90)
		//            outcome = REPORT;
		//        else
		//            outcome = HEAD_FOR_HILLS;
		//#endif
	}

	switch outcome {
	case SWITCH:
		if unit_deserts(target, player(c.who), TRUE, LOY_contract, (amount-bribe_thresh)/2) {
			wout(c.who, "%s accepts the gift, and has decided to join us.", box_name(target))
			p_char(target).fresh_hire = TRUE
			if flag {
				join_stack(target, c.who)
			}
		} else {
			thanks_for_gift(c.who, target)
		}
		break

	case HEAD_FOR_HILLS:
		thanks_for_gift(c.who, target)
		if player(target) != deserted_player && unit_deserts(target, deserted_player, TRUE, LOY_unsworn, 0) {
			wout(c.who, "%s left the service of %s, but didn't join us.", box_name(target), box_name(player(target)))
		}
		break

	case POCKET:
		thanks_for_gift(c.who, target)
		break

	case REPORT:
		thanks_for_gift(c.who, target)
		gen_item(c.who, item_gold, amount)
		wout(target, "%s tried to bribe us with %s.", box_name(c.who), gold_s(amount))
		break

	default:
		panic("!reached")
	}

	return TRUE
}

func v_honor(c *command) int {
	amount := c.a
	if amount == 0 {
		wout(c.who, "Must specify an amount of gold to use as a gift.")
		return FALSE
	}

	if loyal_kind(c.who) == LOY_oath {
		wout(c.who, "%s graciously declines the offer.", box_name(c.who))
		return FALSE
	}

	if FALSE == charge(c.who, amount) {
		wout(c.who, "Do not have %s.", gold_s(amount))
		return FALSE
	}

	if loyal_kind(c.who) != LOY_contract {
		p_char(c.who).loy_kind = LOY_contract
		p_char(c.who).loy_rate = 0
	}

	p_char(c.who).loy_rate += amount
	wout(c.who, "%s now bound with %s.", box_name(c.who), loyal_s(c.who))

	return TRUE
}

func v_oath(c *command) int {
	flag := c.a
	if flag < 1 {
		flag = 1
	} else if flag > 2 {
		flag = 2
	}

	pl := player(c.who)
	lk := loyal_kind(c.who)
	lr := loyal_rate(c.who)
	if lk == LOY_oath && lr >= 2 {
		wout(c.who, "%s already is at %s, the maximum loyalty.", box_name(c.who), loyal_s(c.who))
		return FALSE
	}

	if flag == 2 && lk == LOY_oath && lr == 1 {
		flag = 1
	}
	np_cost := flag

	if !(np_cost > 0) {
		panic("assert(np_cost > 0)")
	}

	if player_np(pl) < 1 {
		wout(c.who, "Player %s has no Noble Points.", box_code(pl))
		return FALSE
	}

	if player_np(pl) < np_cost {
		wout(c.who, "Player %s only has %d Noble Points.", box_code(pl), player_np(pl))
		np_cost = player_np(pl)
	}

	if lk != LOY_oath {
		p_char(c.who).loy_kind = LOY_oath
		p_char(c.who).loy_rate = 0
	}

	p_char(c.who).loy_rate += np_cost
	deduct_np(pl, np_cost)

	wout(c.who, "%s now bound with %s.", box_name(c.who), loyal_s(c.who))

	return TRUE
}

func terrorize_vassal(c *command) bool {
	target := c.a
	severity := c.b
	if severity < 1 {
		severity = 1
	}
	add_char_damage(target, severity, c.who)
	if !alive(target) {
		return false
	}
	if loyal_kind(target) != LOY_fear {
		p_char(target).loy_kind = LOY_fear
		p_char(target).loy_rate = 0
	}
	p_char(target).loy_rate += severity
	wout(c.who, "%s now bound with %s.", box_name(target), loyal_s(target))
	return true
}

func terrorize_prisoner(c *command) bool {
	target := c.a
	severity := c.b
	if severity < 1 {
		severity = 1
	}
	add_char_damage(target, severity, c.who)
	if !alive(target) {
		return false
	}
	wout(c.who, "Health of %s is now %d.", box_name(target), char_health(target))
	if (loyal_kind(target) != LOY_oath) && rnd(1, 100) <= severity && unit_deserts(target, player(c.who), TRUE, LOY_fear, severity) {
		wout(c.who, "%s has been convinced to join us.", box_name(target))
		p_char(target).fresh_hire = TRUE
		return false
	}
	wout(c.who, "%s refuses to swear fealty to us.", box_name(target))
	return false
}

func v_terrorize(c *command) int {
	target := c.a
	//severity := c.b;

	if loyal_kind(c.who) == LOY_fear {
		wout(c.who, "Units of fear loyalty may not terrorize.")
		return FALSE
	}

	if FALSE == check_char_here(c.who, target) {
		return FALSE
	}

	if is_prisoner(target) {
		if stack_leader(target) != stack_leader(c.who) {
			wout(c.who, "%s is not a prisoner of %s.", box_code(target), box_name(c.who))
			return FALSE
		}

		if is_npc(target) {
			wout(c.who, "NPC's cannot swear to player factions.")
			return FALSE
		}

		wout(c.who, "Attempt to gain the loyalty of %s through terror.", box_name(target))
		return TRUE
	}

	if player(target) != player(c.who) {
		wout(c.who, "%s does not belong to our faction.",
			box_code(target))
		return FALSE
	}

	if FALSE == stacked_beneath(c.who, target) {
		wout(c.who, "Unit to be terrorized must be stacked beneath us.")
		return FALSE
	}

	if loyal_kind(target) == LOY_oath {
		wout(c.who, "Oathbound units do not need to have their loyalty reinforced through terror.")
		return FALSE
	}

	wout(c.who, "Increase the loyalty of %s through terror.", box_name(target))

	return TRUE
}

func d_terrorize(c *command) int {
	target := c.a
	//severity := c.b;

	if FALSE == check_still_here(c.who, target) {
		return FALSE
	} else if is_prisoner(target) {
		if stack_leader(target) != stack_leader(c.who) {
			wout(c.who, "%s is not a prisoner of %s.", box_code(target), box_name(c.who))
			return FALSE
		} else if terrorize_prisoner(c) {
			return TRUE
		}
		return FALSE
	} else if player(target) != player(c.who) {
		wout(c.who, "%s does not belong to our faction.", box_code(target))
		return FALSE
	} else if FALSE == stacked_beneath(c.who, target) {
		wout(c.who, "Unit to be terrorized must be stacked beneath us.")
		return FALSE
	} else if loyal_kind(target) == LOY_oath {
		wout(c.who, "Oathbound units do not need to have their loyalty reinforced through terror.")
		return FALSE
	} else if terrorize_vassal(c) {
		return TRUE
	}
	return FALSE
}

func v_raise(c *command) int {
	where := subloc(c.who)

	if FALSE == check_skill(c.who, sk_raise_mob) {
		return FALSE
	} else if FALSE == may_cookie_npc(c.who, where, item_mob_cookie) {
		return FALSE
	} else if has_item(where, item_peasant) < 100 {
		wout(c.who, "Insufficient peasantry to raise a mob.")
		return FALSE
	}

	return TRUE
}

func d_raise(c *command) int {
	where := subloc(c.who)
	mob := do_cookie_npc(c.who, where, item_mob_cookie, where)

	if mob <= 0 {
		log_output(LOG_CODE, "d_raise mob <= 0")
		wout(c.who, "Failed to raise peasant mob.")
		return FALSE
	}

	add_skill_experience(c.who, sk_raise_mob)

	queue(mob, "guard 1")
	init_load_sup(mob) /* make ready to execute commands immediately */

	wout(c.who, "Raised %s.", box_name(mob))
	wout(where, "A speech by %s has raised %s.", box_name(c.who), liner_desc(mob))
	// deduct those peasants from the province.
	sub_item(where, item_peasant, has_item(mob, item_peasant)+1)

	return TRUE
}

func v_rally(c *command) int {
	mob := c.a

	if FALSE == check_skill(c.who, sk_rally_mob) {
		return FALSE
	} else if FALSE == check_char_here(c.who, mob) {
		return FALSE
	} else if noble_item(mob) != item_peasant && noble_item(mob) != item_angry_peasant {
		wout(c.who, "%s is not a peasant mob.", box_name(mob))
		return FALSE
	}

	/*
	 *  Fri Sep 18 12:01:45 1998 -- Scott Turner
	 *
	 *  A limit on how many mobs you can summon, since they
	 *  are orthogonal to "control men in battle"?
	 *
	 */
	num_mobs := 0
	for _, i := range loop_stack(c.who) {
		if noble_item(i) == item_peasant || noble_item(i) == item_angry_peasant {
			num_mobs++
		}
	}
	if num_mobs > 2 {
		wout(c.who, "No more peasant mobs will stack with you.")
		return FALSE
	}
	return TRUE
}

func d_rally(c *command) int {
	mob := c.a

	if FALSE == check_char_gone(c.who, mob) {
		return FALSE
	} else if noble_item(mob) != item_peasant && noble_item(mob) != item_angry_peasant {
		wout(c.who, "%s is not a peasant mob.", box_name(mob))
		return FALSE
	}

	add_skill_experience(c.who, sk_rally_mob)

	if n := stack_parent(mob); n != 0 {
		set_loyal(mob, LOY_summon, min(loyal_rate(mob)+3, 5))
		wout(c.who, "Renewed enthusiasm of %s for %s.", box_name(mob), box_name(n))
		wout(c.who, "The peasants will stay spirited for %d months.", loyal_rate(mob))
		return TRUE
	}

	join_stack(mob, c.who)
	set_loyal(mob, LOY_summon, 3)

	/*
	 *  auto_mob() may have queued some orders, with a preceeding wait.
	 *  Get rid of them now that the mob is LOY_summon
	 */

	flush_unit_orders(player(mob), mob)
	interrupt_order(mob)

	return TRUE // todo: mdhender changed this from return; to return TRUE;
}

func v_incite(c *command) int {
	mob := c.a
	target := c.b

	if FALSE == check_skill(c.who, sk_incite_mob) {
		return FALSE
	}

	if FALSE == check_char_here(c.who, mob) {
		return FALSE
	}

	if subloc(target) != subloc(c.who) {
		wout(c.who, "%s is not here.", box_code(target))
		return FALSE
	}

	if noble_item(mob) != item_peasant &&
		noble_item(mob) != item_angry_peasant {
		wout(c.who, "%s is not a peasant mob.", box_name(mob))
		return FALSE
	}

	if stack_parent(mob) != FALSE {
		wout(c.who, "%s is stacked under a leader.",
			box_name(mob))
		return FALSE
	}

	return TRUE
}

func d_incite(c *command) int {
	mob := c.a
	target := c.b
	where := subloc(c.who)

	if FALSE == check_char_gone(c.who, mob) {
		return FALSE
	}

	if noble_item(mob) != item_peasant &&
		noble_item(mob) != item_angry_peasant {
		wout(c.who, "%s is not a peasant mob.", box_name(mob))
		return FALSE
	}

	if subloc(target) != where {
		wout(c.who, "%s is no longer here.", box_name(target))
		return FALSE
	}

	if stack_parent(mob) != FALSE {
		wout(c.who, "%s is stacked under a leader.",
			box_name(mob))
		return FALSE
	}

	add_skill_experience(c.who, sk_incite_mob)

	if rnd(1, 3) == 1 {
		for _, i := range loop_here(where) {
			if kind(i) != T_loc || subkind(i) != sub_inn {
				continue
			}
			wout(i, "Rumors claim that %s is trying to incite a mob to attack %s.", box_name(c.who), box_name(target))
		}
	}

	if rnd(1, 2) == 1 {
		wout(c.who, "Failed to incite the mob to violence.")
		return FALSE
	}

	flush_unit_orders(player(mob), mob)
	interrupt_order(mob)
	queue(mob, "attack %s", box_code_less(target))
	init_load_sup(mob) /* make ready to execute commands immediately */

	wout(c.who, "%s will attack %s!", box_name(mob), box_name(target))

	return TRUE
}

func v_persuade_oath(c *command) int {
	target := c.a

	if FALSE == check_char_here(c.who, target) {
		return FALSE
	}

	if char_new_lord(target) != FALSE {
		wout(c.who, "%s just switched employers this month, and is not looking for a new one so soon.", box_name(target))
		return FALSE
	}

	if FALSE == can_pay(c.who, 25) {
		wout(c.who, "Don't have %s.", gold_s(25))
		return FALSE
	}

	return TRUE
}

func d_persuade_oath(c *command) int {
	target := c.a
	flag := c.b != FALSE

	if FALSE == check_still_here(c.who, target) {
		return FALSE
	}

	if char_new_lord(c.who) != FALSE {
		wout(c.who, "%s just switched employers this month, and is not looking for a new one so soon.", box_name(target))
		return FALSE
	}

	if loyal_kind(target) != LOY_oath {
		wout(c.who, "%s does not have oath loyalty.", box_name(target))
		return FALSE
	}

	if FALSE == charge(c.who, 25) {
		wout(c.who, "Don't have %s.", gold_s(25))
		return FALSE
	}

	if loyal_rate(target) != 1 || rnd(1, 100) > 2 {
		wout(c.who, "Failed to convince %s to join us.", box_name(target))
		return TRUE
	}

	if unit_deserts(target, player(c.who), TRUE, LOY_contract, 250) {
		wout(c.who, "%s has been convinced to join us!", box_name(target))
		p_char(target).fresh_hire = TRUE
		if flag {
			join_stack(target, c.who)
		}
	} else {
		wout(c.who, "Failed to convince %s to join us.",
			box_name(target))
	}

	return TRUE
}
