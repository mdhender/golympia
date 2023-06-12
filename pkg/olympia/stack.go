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

func here_preceeds(a, b int) int { panic("!implemented") }

func here_pos(who int) int {
	var p *loc_info
	var ret int

	p = rp_loc_info(loc(who))

	assert(p != nil)

	ret = ilist_lookup(p.here_list, who)

	assert(ret >= 0)

	return ret
}

/*
 *  Does a come before b in the here list?
 */

func here_precedes(a, b int) int {
	var p *loc_info
	var i int

	if loc(a) != loc(b) {
		return FALSE
	} /* they're in different here lists */

	p = rp_loc_info(loc(a))
	assert(p != nil)

	for i = 0; i < len(p.here_list); i++ {
		if p.here_list[i] == a {
			return TRUE
		} else if p.here_list[i] == b {
			return FALSE
		}
	}
	panic("!reached")
}

func first_prisoner_pos(where int) int {
	var p *loc_info
	var i int

	p = rp_loc_info(where)

	if p == nil {
		return -1
	}

	for i = 0; i < len(p.here_list); i++ {
		if kind(p.here_list[i]) == T_char &&
			is_prisoner(p.here_list[i]) {
			return i
		}
	}

	return -1
}

func stack_parent(who int) int {
	var n int

	n = loc(who)

	if kind(n) == T_char {
		return n
	}

	return 0
}

func stack_leader(who int) int {
	var n int
	count := 0

	assert(kind(who) == T_char)

	n = who
	for kind(n) == T_char {
		who = n
		n = stack_parent(n)

		// todo: count may not increment in release builds
		assert(count < 1000) /* infinite loop check */
		count++
	}

	return who
}

/*
 *  is b stacked somewhere beneath a?
 */

func stacked_beneath(a, b int) int {

	assert(kind(a) == T_char)
	assert(kind(b) == T_char)

	if a == b {
		return FALSE
	}

	for b > 0 {
		b = stack_parent(b)
		if a == b {
			return TRUE
		}
	}

	return FALSE
}

func promote(who, new_pos int) {
	var p *loc_info
	var i int
	var who_pos int

	p = rp_loc_info(loc(who))
	assert(p != nil)

	who_pos = ilist_lookup(p.here_list, who)

	assert(who_pos >= new_pos)

	for i = who_pos; i > new_pos; i-- {
		p.here_list[i] = p.here_list[i-1]
	}
	p.here_list[new_pos] = who
}

func unstack(who int) {
	leader := stack_leader(who)

	assert(valid_box(leader))

	/*
	 *  This assert is a late add-on, to convince myself that the
	 *  promote is correct for an unstack from multiple levels deep.
	 */
	assert(subloc(leader) == loc(leader))

	if release_swear(who) != FALSE {
		p_magic(who).swear_on_release = FALSE
	}

	set_where(who, subloc(leader))
	promote(who, here_pos(leader)+1)

	/*
	 *  If unstacking while moving, we have to call restore_stack_actions
	 *  on who after who has been unstacked.
	 */

	restore_stack_actions(who)

	if loyal_kind(who) == LOY_summon {
		set_loyal(who, LOY_npc, 0)
	}
}

func leave_stack(who int) {
	var leader int

	leader = stack_parent(who)
	if leader <= 0 {
		return
	}

	wout(leader, "%s unstacks from us.", box_name(who))
	wout(who, "%s unstacks from %s.", box_name(who), box_name(leader))

	vector_char_here(who)
	wout(VECT, "%s unstacks from %s.",
		box_name(who), box_name(leader), just_name(who))

	unstack(who)
}

func stack(who, target int) {
	var pos int

	assert(stack_parent(who) == 0)
	set_where(who, target)
	p_char(who).moving = char_moving(target)

	/*
	 *  Keep prisoners at the end of the stacking list
	 *  by promoting non-prisoners who join
	 */

	if !is_prisoner(who) {
		pos = first_prisoner_pos(target)

		if pos >= 0 {
			promote(who, pos)
		}
	}
}

func join_stack(who, target int) {

	assert(stacked_beneath(who, target) == 0)
	assert(!is_prisoner(target))

	leave_stack(who)

	assert(subloc(target) == subloc(who))

	wout(who, "%s stacks beneath %s.", box_name(who), box_name(target))
	wout(target, "%s stacks beneath us.", box_name(who))

	vector_char_here(who)
	wout(VECT, "%s stacks beneath %s.", box_name(who), box_name(target))

	stack(who, target)
}

func check_prisoner_escape(who, chance int) int {

	if rnd(1, 100) > chance {
		return FALSE
	}

	prisoner_escapes(who)
	return TRUE
}

func prisoner_escapes(who int) {
	var leader int
	var where, out_one int

	leader = stack_parent(who)

	wout(leader, "Prisoner %s escaped!", box_name(who))
	p_char(who).prisoner = FALSE
	p_magic(who).swear_on_release = FALSE
	unstack(who)
	touch_loc(who)

	wout(who, "We escaped!")

	where = subloc(who)

	if loc_depth(where) <= LOC_province {
		return
	}

	out_one = loc(where)

	if is_ship(where) && subkind(out_one) == sub_ocean {
		out_one = find_nearest_land(out_one)

		wout(who, "After jumping over the side of the boat and enduring a long, grueling, swim, we finally washed ashore at %s.", box_name(out_one))

		wout(leader, "%s jumped overboard and presumably drowned.", just_name(who))

		log_output(LOG_SPECIAL, "!! Someone swam ashore, who=%s", box_code_less(who))
	}

	move_stack(who, out_one)
}

func prisoner_movement_escape_check(who int) {

	for _, i := range loop_char_here(who) {
		if is_prisoner(i) {
			check_prisoner_escape(i, 2)
		}
	}
}

func weekly_prisoner_escape_check() {
	var chance int

	for _, who := range loop_char() {
		if is_prisoner(who) {
			continue
		}

		if subkind(subloc(who)) == sub_ocean {
			continue
		} /* they're flying */

		for _, i := range loop_here(who) {
			if kind(i) == T_char && is_prisoner(i) && release_swear(i) == FALSE {
				chance = or_int(loc_depth(subloc(who)) >= LOC_build, 1, 2)
				check_prisoner_escape(i, chance)
			}
		}
	}
}

func drop_stack(who, to_drop int) {
	release_swear_flag := FALSE

	assert(stack_parent(to_drop) == who)

	if is_prisoner(to_drop) {
		p_char(to_drop).prisoner = FALSE
		touch_loc(to_drop)
		wout(who, "Freed prisoner %s.", box_name(to_drop))
		wout(to_drop, "%s set us free.", box_name(who))

		if release_swear(to_drop) != FALSE {
			release_swear_flag = TRUE
		}
	} else {
		wout(who, "Dropped %s from stack.", box_name(to_drop))
		wout(to_drop, "%s dropped us from the stack.", box_name(who))

		vector_char_here(to_drop)
		wout(VECT, "%s dropped %s from the stack.",
			box_name(who),
			box_name(to_drop),
			just_name(to_drop))
	}

	unstack(to_drop)

	if release_swear_flag != FALSE {
		log_output(LOG_SPECIAL, "%s frees a swear_on_release prisoner", box_name(who))

		if rnd(1, 5) < 5 {
			wout(who, "%s is grateful for your gallantry.",
				box_name(to_drop))
			wout(who, "%s pledges fealty to us.",
				box_name(to_drop))

			set_lord(to_drop, player(who), LOY_oath, 1)
		} else {
			switch rnd(1, 3) {
			case 1:
				wout(who, "%s spits on you, and vanishes in a cloud of orange smoke.", box_name(to_drop))
				break
			case 2:
				wout(who, "%s cackles wildly and vanishes.", box_name(to_drop))
				break
			case 3:
				wout(who, "%s smiles briefly at you, then vanishes.", box_name(to_drop))
				break
			}

			unit_deserts(to_drop, 0, TRUE, LOY_unsworn, 0)
			put_back_cookie(to_drop)
			set_where(to_drop, 0)
			change_box_kind(to_drop, T_deadchar)
		}
	}
}

func free_all_prisoners(who int) {
	for _, i := range loop_here(who) {
		if kind(i) == T_char && is_prisoner(i) {
			drop_stack(who, i)
		}
	}
}

/*
 *  Remove who from a stack, leaving those above and below him behind.
 *
 *  This routine used to let who's prisoners go free, but now it is
 *  possible to use this routine manually, from the UNSTACK command.
 */

func extract_stacked_unit(who int) {
	first := 0

	/*
	 *  locate first stacked non-pris char
	 */

	for _, i := range loop_here(who) {
		if kind(i) == T_char && !is_prisoner(i) {
			first = i
			break
		}
	}

	/*
	 *  move all other chars beneath
	 *  move up & out one level, position just after us
	 */

	if first != 0 {
		for _, i := range loop_here(who) {
			if i != first && kind(i) == T_char {
				set_where(i, first)
			}
		}

		set_where(first, loc(who))
		promote(first, here_pos(who)+1)
	}

	/*
	 *  Free any prisoners left
	 */

	for _, i := range loop_here(who) {
		if kind(i) == T_char && is_prisoner(i) {
			prisoner_escapes(i)
		}
	}

	vector_char_here(who)
	wout(VECT, "%s unstacks.", box_name(who))

	leave_stack(who)
}

/*
 *  Promote lower to be before higher
 */

func promote_stack(lower, higher int) {
	var p *loc_info
	var pos int

	assert(stacked_beneath(lower, higher) == 0)

	set_where(lower, loc(higher))

	/*
	 *  Now the here list has higher at some point, at or after the
	 *  beginning, and lower as the last element.  We want to move
	 *  lower to be before higher.
	 */

	p = rp_loc_info(loc(higher))
	assert(p != nil)

	assert(p.here_list[len(p.here_list)-1] == lower)
	pos = ilist_lookup(p.here_list, higher)

	promote(lower, pos)

	wout(higher, "Promoted %s.", box_name(lower))
	wout(lower, "%s promoted us.", box_name(higher))
	//#if 0
	//    vector_char_here(lower);
	//    wout(VECT, "%s promoted %s.",            box_name(higher), box_name(lower));
	//#endif
}

func take_prisoner(who, target int) {
	ni := false

	assert(who != target)
	assert(kind(who) == T_char)
	assert(kind(target) == T_char)

	if subkind(target) == sub_ni && beast_capturable(target) {
		ni = true
	}

	vector_stack(stack_leader(target), true)
	vector_add(who)

	if ni {
		wout(VECT, "%s disbands.", box_name(target))
	} else {
		wout(VECT, "%s is %s by %s.", box_name(target), "taken prisoner", box_name(who))
	}

	//#if 0
	//    {
	//        show_day = FALSE;
	//        out(target, "");
	//        show_day = TRUE;
	//    }
	//#endif

	/*
	 *  Suppresses further output to target
	 */
	p_char(target).prisoner = TRUE

	take_unit_items(target, who, or_int(ni, TAKE_NI, TAKE_ALL))
	extract_stacked_unit(target)
	interrupt_order(target)

	if ni {
		unit_deserts(target, 0, TRUE, LOY_unsworn, 0)
		put_back_cookie(target)
		set_where(target, 0)
		change_box_kind(target, T_deadchar)
	} else {
		stack(target, who)
	}
}

func has_prisoner(who, pris int) int {
	ret := FALSE

	for _, i := range loop_here(who) {
		if i == pris && is_prisoner(i) {
			ret = TRUE
			break
		}
	}

	return ret
}

func move_prisoner(who, target, pris int) int {
	rs := release_swear(pris)

	unstack(pris)
	stack(pris, target)

	if rs != 0 {
		p_magic(pris).swear_on_release = TRUE
	}

	assert(rp_char(pris).prisoner != 0)

	return TRUE
}

func give_prisoner(who, target, pris int) int {

	if check_prisoner_escape(pris, 2) != FALSE {
		return FALSE
	}

	move_prisoner(who, target, pris)

	wout(who, "Transferred prisoner %s to %s.",
		box_name(pris), box_name(target))

	wout(target, "%s transferred the prisoner %s to us.",
		box_name(who), box_name(pris))

	return TRUE
}

/*
 *  Test for priests, mu/undead for stacking restraints
 *
 */
func stack_contains_priest(who int) int {
	for _, i := range loop_stack(who) {
		if is_priest(i) != FALSE && !is_prisoner(i) {
			return 1
		}
	}
	return 0
}

func contains_mu_undead(i int) int {
	if is_magician(i) && char_hide_mage(i) == FALSE {
		return 1
	}
	/*
	 *  These checks are for "independent" undead.
	 *
	 */
	if noble_item(i) != FALSE && subkind(noble_item(i)) == sub_undead {
		return 1
	}
	if noble_item(i) != FALSE && subkind(noble_item(i)) == sub_demon_lord {
		return 1
	}
	/*
	 *  Check inventory for undead.
	 *
	 */
	for _, j := range inventory_loop(i) {
		if subkind(j.item) == sub_undead || subkind(j.item) == sub_demon_lord {
			return 1
		}
	}

	return 0
}

func stack_contains_mu_undead(who int) int {
	for _, i := range loop_stack(who) {
		if contains_mu_undead(i) != FALSE && !is_prisoner(i) {
			return 1
		}
	}
	return 0
}

func v_stack(c *command) int {
	target := c.a

	if check_char_gone(c.who, target) {
		return FALSE
	}

	if target == c.who {
		wout(c.who, "Can't stack beneath oneself.")
		return FALSE
	}

	if stacked_beneath(c.who, target) != FALSE {
		wout(c.who, "Cannot stack beneath %s since %s is stacked under you.",
			box_name(target), just_name(target))
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Can't stack beneath prisoners.")
		return FALSE
	}

	if will_admit(target, c.who, target) == FALSE {
		wout(c.who, "%s refuses to let us stack.", box_name(target))
		wout(target, "Refused to let %s stack with us.",
			box_name(c.who))
		return FALSE
	}

	/*
	 *  Clerics can't stack w/ magicians/undead and vice versa
	 *
	 */
	if options.mp_antipathy &&
		((stack_contains_priest(c.who) != FALSE &&
			stack_contains_mu_undead(target) != FALSE) ||
			(stack_contains_priest(target) != FALSE &&
				stack_contains_mu_undead(c.who) != FALSE)) {
		wout(c.who, "Priests cannot stack with magicians or undead, and vice versa.")
		return FALSE
	}
	join_stack(c.who, target)
	return TRUE
}

func v_unstack(c *command) int {
	target := c.a

	if numargs(c) < 1 {
		if stack_parent(c.who) <= 0 {
			wout(c.who, "Not stacked under anyone.")
			return FALSE
		}

		leave_stack(c.who)
		return TRUE
	}

	if c.who == target {
		//#if 1
		extract_stacked_unit(c.who)
		//#else
		//        leave_stack(c.who);
		//
		//        loop_here(c.who, i)
		//        {
		//            drop_stack(c.who, i);
		//        }
		//        next_here;
		//#endif

		return TRUE
	}

	if !valid_box(target) || stack_parent(target) != c.who {
		wout(c.who, "%s is not stacked beneath us.", c.parse[1])
		return FALSE
	}

	drop_stack(c.who, target)
	return TRUE
}

func v_surrender(c *command) int {
	target := c.a

	if !check_char_gone(c.who, target) {
		return FALSE
	}

	if player(target) == player(c.who) {
		wout(c.who, "Can't surrender to oneself.")
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Can't surrender to a prisoner.")
		return FALSE
	}

	log_output(LOG_SPECIAL, "Player %s surrenders %s",
		box_code_less(player(c.who)),
		box_name(c.who))

	vector_stack(stack_leader(c.who), true)
	vector_stack(stack_leader(target), false)

	wout(VECT, "%s surrenders to %s.", box_name(c.who), box_name(target))

	take_prisoner(target, c.who)
	return TRUE
}

/*
 *  Does b appear later in location order than a
 */
func promote_after(a, b int) int {
	where := subloc(a)
	ret := 0

	assert(subloc(b) == where)

	for _, i := range loop_char_here(where) {
		if i == a {
			ret = a
			break
		} else if i == b {
			ret = b
			break
		}

	}

	if ret == a {
		return TRUE
	} else if ret == b {
		return FALSE
	}

	panic("!reached")
}

func v_promote(c *command) int {
	target := c.a
	var targ_par int

	if numargs(c) < 1 {
		wout(c.who, "Must specify which character to promote.")
		return FALSE
	}

	if kind(target) != T_char {
		wout(c.who, "%s is not a character.", c.parse[1])
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Can't promote prisoners.")
		return FALSE
	}

	targ_par = stack_parent(target)

	if !check_char_here(c.who, target) {
		return FALSE
	}

	if target == c.who {
		wout(c.who, "Can't promote oneself.")
		return FALSE
	}

	/*
	 *  Only do the strict check if the unit belongs to another character.
	 *  If its one of ours, just promote it, since the units are in the
	 *  same location.
	 */

	if player(c.who) == player(target) {
		if promote_after(c.who, target) == FALSE {
			wout(c.who, "%s already comes before us in location order.", box_name(target))
			return FALSE
		}
	} else if targ_par != c.who && here_precedes(c.who, target) == FALSE {
		wout(c.who, "Characters to be promoted must be stacked immediately beneath the promoter, or be listed after the promoter at the same level.")

		return FALSE
	}

	promote_stack(target, c.who)
	return TRUE
}
