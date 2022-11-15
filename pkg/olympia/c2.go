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

func v_discard(c *command) int {
	item := c.a
	qty := c.b
	have_left := c.c
	var ret int

	if kind(item) != T_item {
		wout(c.who, "%s is not an item.", box_code(item))
		return FALSE
	}

	if has_item(c.who, item) < 1 {
		wout(c.who, "%s does not have any %s.", box_name(c.who),
			box_code(item))
		return FALSE
	}

	qty = how_many(c.who, c.who, item, qty, have_left)

	if qty <= 0 {
		return FALSE
	}

	ret = drop_item(c.who, item, qty)
	if !(ret != 0) {
		panic("assert(ret!=0)")
	}

	wout(c.who, "Dropped.")

	return TRUE
}

func drop_player(pl int) {
	var who int
	s := ""
	email := ""
	var owner int
	var i int

	assert(kind(pl) == T_player)

	/*
	 *  Mon Apr 12 13:25:34 1999 -- Scott Turner
	 *
	 *  Remove any deserted nobles belonging to this player.
	 *
	 */
	for _, who = range loop_units(deserted_player) {
		if body_old_lord(who) == pl {
			if is_prisoner(who) && nation(stack_parent(who)) != nation(who) {
				rp_misc(who).old_lord = 0
			} else {
				wout(subloc(who), "%s melts into the ground and vanishes.",
					box_name(who))
				char_reclaim(who)
			}
		}
	}

	for _, who = range loop_units(pl) {
		if is_prisoner(who) && nation(stack_parent(who)) != nation(who) {
			unit_deserts(who, deserted_player, TRUE, LOY_UNCHANGED, 0)
			/*
			 *  We need to reset the old_lord at this point
			 *  because shortly the player is not going to
			 *  exist.
			 *
			 *  Mon Apr 12 13:16:49 1999 -- Scott Turner
			 *
			 *  This also permits you to capture this noble
			 *  at no NP cost.
			 *
			 */
			rp_misc(who).old_lord = 0
		} else {
			wout(subloc(who), "%s melts into the ground and vanishes.", box_name(who))
			char_reclaim(who)
		}
	}

	/*
	 *  Immediately rot any dead bodies belonging to the dropped player.
	 */
	for _, who = range loop_dead_body() {
		owner = item_unique(i)
		if owner == 0 {
			panic("assert(owner!=0)")
		}

		if rp_misc(i) == nil || rp_misc(i).old_lord != pl {
			continue
		}

		//#if 1
		p_misc(i).old_lord = 0
		//#else
		//                if (kind(owner) == T_char)
		//              wout(owner, "%s decomposed.", box_name(i));
		//
		//                destroy_unique_item(owner, i);
		//#endif
	}

	//#if 0
	//    loop_lost_soul(i)
	//      {
	//        if (rp_misc(i) == nil || rp_misc(i).old_lord != pl)
	//      continue;
	//
	//        p_misc(i).old_lord = 0;
	//      }
	//    next_lost_soul;
	//#endif

	if rp_player(pl) != nil {
		if rp_player(pl).email != "" {
			email = rp_player(pl).email
		}

		if rp_player(pl).full_name != "" {
			s = rp_player(pl).full_name
		}
	}

	log_output(LOG_DROP, "Dropped player %s", box_name(pl))
	log_output(LOG_DROP, "    %s <%s>", s, email)

	if save_flag && acct_flag {
		//sprintf(cmd, "%s -a %s -A %s-old -p %s -g tag%d -d", options.accounting_prog, options.accounting_dir, options.accounting_dir, box_code_less(pl), game_number);
		//system(cmd);
		panic("!implemented")
	}

	delete_box(pl)
}

func v_quit(c *command) int {
	target := c.a

	if target == 0 {
		if c.who == gm_player {
			return FALSE
		} else {
			target = c.who
		}
	}

	if target != c.who && c.who != gm_player {
		wout(c.who, "Not allowed to drop another player.")
		return FALSE
	}

	drop_player(target)

	/*
	 *  Don't call finish_command, we don't exist anymore
	 */

	return FALSE
}

func text_list_free(l [][]byte) {
	var i int

	for i = 0; i < len(l); i++ {
		my_free(l[i])
	}
}

func line_length_check(l [][]byte) int {
	length := 0
	for i := 0; i < len(l); i++ {
		length = max(length, len(l[i]))
	}
	return length
}

func parse_text_list(c *command) [][]byte {
	lines := c.a
	pl := player(c.who)

	var l [][]byte
	if lines == 0 {
		done := true
		for !done {
			order := top_order(pl, c.who)
			if order == nil {
				wout(c.who, "Ran out of posting text.")
				text_list_free(l)
				return nil
			}

			t := eat_leading_trailing_whitespace(order)
			if i_strcmp(t, []byte("end")) == 0 {
				done = true
			} else {
				l = append(l, str_save(order))
			}
			pop_order(pl, c.who)
		}
	} else {
		for lines > 0 {
			order := top_order(pl, c.who)
			if order == nil {
				wout(c.who, "Ran out of posting text.")
				text_list_free(l)
				return nil
			}

			l = append(l, str_save(order))
			pop_order(pl, c.who)

			lines--
		}
	}

	return l
}

func v_post(c *command) int {
	l := parse_text_list(c)
	if l == nil {
		return FALSE
	}

	if line_length_check(l) > MAX_POST {
		wout(c.who, "Line length of posted text exceeds %d characters.", MAX_POST)
		wout(c.who, "Post rejected.")
		text_list_free(l)
		return FALSE
	}

	newEnt := new_ent(T_post, 0)
	p_item_magic(newEnt).creator = c.who
	set_where(newEnt, subloc(c.who))

	p_misc(newEnt).post_txt = nil
	for _, s := range l {
		p_misc(newEnt).post_txt = append(p_misc(newEnt).post_txt, string(s))
	}

	return TRUE
}

func v_message(c *command) int {
	targ := c.b
	var i int

	l := parse_text_list(c)
	if l == nil {
		return FALSE
	}

	if line_length_check(l) > MAX_POST {
		wout(c.who, "Line length of message text exceeds %d characters.", MAX_POST)
		wout(c.who, "Message rejected.")
		text_list_free(l)
		return FALSE
	}

	if !valid_box(targ) {
		wout(c.who, "%s is not a valid entity.", box_code(targ))
		text_list_free(l)
		return FALSE
	}

	if kind(targ) != T_char {
		wout(c.who, "May not send a message to %s.", box_code(targ))
		text_list_free(l)
		return FALSE
	}

	wout(targ, "Received a message from %s:", box_name(c.who))

	indent += 3
	first := true
	for i = 0; i < len(l); i++ {
		wout(targ, "%s%s%s",
			or_string(first, "\"", ""),
			l[i],
			or_string(i+1 == len(l), "\"", ""))

		if first {
			first = false
			indent += 1
		}
	}

	if !first {
		indent -= 1
	}

	indent -= 3

	text_list_free(l)

	wout(c.who, "Message delivered.")
	return TRUE
}

func v_improve_opium(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_poppy_field {
		wout(c.who, "Opium is produced only in poppy fields.")
		return FALSE
	}

	return TRUE
}

func d_improve_opium(c *command) int {
	where := subloc(c.who)

	if subkind(where) != sub_poppy_field {
		wout(c.who, "Not in a poppy field anymore.")
		return FALSE
	}

	p_misc(where).opium_double = TRUE

	return TRUE
}

func v_die(c *command) int {
	kill_char(c.who, MATES, S_body)
	return TRUE
}

//#if 0
//int
//v_format(c *command) int {
//    int pl;
//
//    pl = player(c.who);
//
//    p_player(pl).format = c.a;
//
//    wout(c.who, "Formatting for %s set to %d.", box_name(pl), c.a);
//
//    return TRUE;
//}
//#endif

func v_stop(c *command) int {

	return TRUE
}

func v_archery(c *command) int {

	return TRUE
}

func d_archery(c *command) int {
	var p *entity_char

	p = p_char(c.who)

	p.missile++

	wout(c.who, "Missile rating now %d.", p.missile)
	return TRUE
}

func v_defense(c *command) int {

	return TRUE
}

func d_defense(c *command) int {
	var p *entity_char

	p = p_char(c.who)

	p.defense++

	wout(c.who, "Defense rating now %d.", p.defense)
	return TRUE
}

func v_swordplay(c *command) int {

	return TRUE
}

func d_swordplay(c *command) int {
	var p *entity_char

	p = p_char(c.who)

	p.attack++

	wout(c.who, "Attack rating now %d.", p.attack)
	return TRUE
}

func v_claim(c *command) int {
	item := c.a
	qty := c.b
	have_left := c.c
	pl := player(c.who)
	var ret int

	/*
	 *  Common mistake checker!
	 *
	 *  If they said CLAIM 500, assume that they meant CLAIM 1 500.
	 */

	if numargs(c) < 2 &&
		(kind(item) != T_item || has_item(pl, item) < 1) &&
		qty == 0 {
		log_output(LOG_CODE, "correcting CLAIM for %s:  %s",
			box_code_less(player(c.who)), c.line)

		wout(c.who, "(assuming you meant CLAIM %d %d)", item_gold, item)
		qty = item
		item = item_gold
	}

	if kind(item) != T_item {
		wout(c.who, "%s is not an item.", box_code(item))
		return FALSE
	}

	if has_item(pl, item) < 1 {
		wout(c.who, "No %s for you to claim.", box_code(item))
		return FALSE
	}

	qty = how_many(c.who, pl, item, qty, have_left)

	if qty <= 0 {
		return FALSE
	}

	ret = move_item(pl, c.who, item, qty)
	if !(ret != 0) {
		panic("assert(ret!=0)")
	}

	wout(c.who, "Claimed %s.", just_name_qty(item, qty))

	/*
	 *  Track gold from claim.
	 *
	 */
	if item == item_gold {
		//extern int gold_claim;
		gold_claim += qty
	}
	return TRUE
}

func v_fight_to_death(c *command) int {
	flag := c.a

	if flag < 0 {
		flag = 0
	}
	if flag > 100 {
		flag = 100
	}

	p_char(c.who).break_point = flag
	wout(c.who, "Troops led by %s will break at %d%%%.",
		box_name(c.who), flag)

	return TRUE
}

func v_fee(c *command) int {
	where := c.c
	type_ := c.parse[1]
	amount := c.b
	var control *loc_control_ent

	/*
	 *  If no "where" argument, then assume the noble's subloc.
	 *
	 */
	if where != 0 {
		where = subloc(c.who)
	}

	if !valid_box(where) || (nil == rp_loc(where) && nil == rp_subloc(where)) {
		wout(c.who, "That's not a place where you can set an entrance fee.")
		return FALSE
	}

	if amount < 0 {
		wout(c.who, "Fee must be non-negative.")
		return FALSE
	}

	if kind(where) == T_loc && loc_depth(where) == LOC_province {
		/*
		 *  Requires a garrison.
		 *
		 */
		if FALSE == garrison_here(where) || province_admin(province(where)) != c.who {
			wout(c.who, "You do not rule that province.")
			return FALSE
		}
		control = &p_loc(where).control2
	} else {
		if nil == p_subloc(where) || first_character(where) != c.who {
			wout(c.who, "You do not control that location.")
			return FALSE
		}
		control = &p_subloc(where).control2
	}

	if strcasecmp(type_, []byte("weight")) == 0 {
		control.weight = amount
		wout(c.who, "Entrance fee for %s will be set to %s per 1000 weight.",
			box_name(where), gold_s(amount))
	} else if strcasecmp(type_, []byte("men")) == 0 {
		control.men = amount
		wout(c.who, "Entrance fee for %s will be set to %s per 100 men.",
			box_name(where), gold_s(amount))
	} else if strcasecmp(type_, []byte("nobles")) == 0 {
		control.nobles = amount
		wout(c.who, "Entrance fee for %s will be set to %s per noble.",
			box_name(where), gold_s(amount))
	} else {
		wout(c.who, "The type of fee must be WEIGHT, MEN or NOBLES.")
		return FALSE
	}

	return TRUE
}

/*
 *  Thu Mar 20 11:34:02 1997 -- Scott Turner
 *
 *  Close or open a border
 *
 *  border <where> {open|close}
 *
 */

func v_border(c *command) int {
	where := c.a
	var closed *int

	if !valid_box(where) ||
		(nil == rp_loc(where) && nil == rp_subloc(where)) {
		wout(c.who, "That's not a place with a border.")
		return FALSE
	}

	if kind(where) == T_loc &&
		loc_depth(where) == LOC_province {
		/*
		 *  Province borders require a garrison.
		 *
		 */
		if FALSE == garrison_here(where) ||
			player(province_admin(where)) != player(c.who) {
			wout(c.who, "You do not control that garrison.")
			return FALSE
		}
		closed = &rp_loc(where).control.closed
	} else {
		if nil == p_subloc(where) || first_character(where) != c.who {
			wout(c.who, "You do not control that location.")
			return FALSE
		}
		closed = &rp_subloc(where).control.closed
	}

	if strncasecmp(c.parse[2], []byte("open"), 4) == 0 {
		*closed = 0
		wout(c.who, "The border at %s is now open.", box_name(where))
		return TRUE
	}

	if strncasecmp(c.parse[2], []byte("close"), 5) == 0 {
		*closed = 1
		wout(c.who, "The border at %s is now closed.", box_name(where))
		return TRUE
	}

	wout(c.who, "Must specify OPEN or CLOSE.")
	return FALSE
}

func board_message(who, ship int) {
	where := subloc(ship)

	if char_really_hidden(who) {
		return
	}

	if weather_here(where, sub_fog) != FALSE && is_priest(who) != sk_domingo {
		return
	}

	desc := liner_desc(who)
	comma := ""
	if strings.Index(desc, ",") != -1 {
		comma = ","
	}

	with := display_with(who)
	if len(with) == 0 {
		with = "."
	}

	wout(where, "%s%s%s boarded %s%s", desc, comma,
		box_name(ship), with)
	show_chars_below(where, who)
}

func v_board(c *command) int {
	ship := c.a
	max_fee := will_pay(c.who)
	var v *exit_view
	var ship_fee int /* fee the captain is charging */
	var owner int    /* captain of the ship */
	var w weights    /* how much we weigh */
	var amount int   /* how much we have to pay */
	var sw, sc int   /* ship weight, capacity */
	//extern int gold_ferry;

	if !is_ship(ship) {
		wout(c.who, "%s is not a ship.", box_code(ship))
		return FALSE
	}

	log_output(LOG_SPECIAL, "BOARD for %s", box_name(player(c.who)))

	v = parse_exit_dir(c, subloc(c.who), "board")

	if v == nil {
		return FALSE
	}

	assert(v.destination == ship)

	if v.in_transit != FALSE {
		wout(c.who, "%s is underway.  Boarding is not possible.", box_name(v.destination))
		return FALSE
	}

	owner = building_owner(ship)

	if nil == rp_subloc(ship) {
		wout(c.who, "Something prevents boarding that ship.")
		return FALSE
	}

	// todo: confusing
	noFerry := !valid_box(owner)
	if !noFerry {
		if ship_fee = calc_entrance_fee(&rp_subloc(ship).control, c, owner); ship_fee != 0 {
			noFerry = false
		}
	}
	if noFerry {
		wout(c.who, "%s is not being operated as a ferry (no boarding FEE is set).", box_name(ship))
		return FALSE
	}

	determine_stack_weights(c.who, &w, FALSE)

	/*
	 *  Check that the ship isn't already overloaded, and that it won't
	 *  be overloaded if we board.
	 */

	sc = ship_cap(ship)
	if sc != 0 {
		sw = ship_weight(ship)

		if sw > sc {
			wout(c.who, "%s is already overloaded.  It can take no more passengers.", box_name(ship))
			wout(owner, "Refused to let %s board because we are overloaded.", box_name(c.who))
			return FALSE
		}

		if sw+w.total_weight > sc {
			wout(c.who, "%s would be overloaded with us.  We can't board.", box_name(ship))
			wout(owner, "Refused to let %s board because then we would be overloaded.", box_name(c.who))
			return FALSE
		}
	}

	if max_fee != 0 && amount > max_fee {
		wout(c.who, "Refused to pay a boarding fee of %s.",
			gold_s(amount))
		wout(owner, "%s refused to pay a boarding fee of %s.",
			box_name(c.who),
			gold_s(amount))
		return FALSE
	}

	if FALSE == charge_entrance_fees(c.who, owner, amount) {
		wout(owner, "%s couldn't afford a boarding fee of %s.",
			box_name(c.who),
			gold_s(amount))
		return FALSE
	}

	wout(c.who, "Paid %s to board %s.", gold_s(amount), box_name(ship))
	wout(owner, "%s paid %s to board.", box_name(c.who), gold_s(amount))
	board_message(c.who, ship)
	gold_ferry += amount
	move_stack(c.who, ship)

	return TRUE
}

func unboard_message(who, ship int) {
	where := subloc(ship)

	if char_really_hidden(who) {
		return
	}
	if weather_here(where, sub_fog) != FALSE && is_priest(who) != sk_domingo {
		return
	}

	var comma string
	desc := liner_desc(who)
	if strings.Index(desc, ",") != -1 {
		comma = ","
	}

	with := display_with(who)
	if len(with) == 0 {
		with = "."
	}

	wout(where, "%s%s disembarked from %s%s",
		desc, comma,
		box_name(ship), with)
	show_chars_below(where, who)
}

/*
 *  Unload passengers from a ferry
 */

func v_unload(c *command) int {
	ship := subloc(c.who)
	var where int
	var i int

	if !is_ship(ship) || building_owner(ship) != c.who {
		wout(c.who, "%s is not the captain of a ship.",
			box_name(c.who))
		return FALSE
	}

	where = subloc(ship)

	if subkind(where) == sub_ocean {
		wout(c.who, "Can't unload passengers at sea.  They won't go.")
		return FALSE
	}

	/*
	 *  Tue Feb 27 10:07:09 2001 -- Scott Turner
	 *
	 *  This ends up unstacking nobles because it doesn't check whether
	 *  they're already out of the ship (due to a previous move_stack).
	 *  Avoid this by only moving stack leaders.
	 *
	 */
	unloadedPassengers := false
	for _, i = range loop_char_here(ship) {
		if i == c.who {
			continue
		}
		if stack_leader(i) != i {
			continue
		}

		wout(c.who, "%s disembarks.", box_name(i))
		wout(i, "%s disembarks.", box_name(i))
		unboard_message(i, ship)

		move_stack(i, where)
		unloadedPassengers = true
	}

	if unloadedPassengers {
		wout(c.who, "All passengers unloaded.")
	} else {
		wout(c.who, "No passengers to unload.")
	}

	return TRUE
}

/*
 *  Toot our horn -- wakeup any WAIT FERRY's in the port
 */

func v_ferry(c *command) int {
	ship := subloc(c.who)
	var where int

	if !is_ship(ship) || building_owner(ship) != c.who {
		wout(c.who, "%s is not the captain of a ship.",
			box_name(c.who))
		return FALSE
	}

	where = subloc(ship)

	wout(where, "%s sounds a blast on its horn.", box_name(ship))
	log_output(LOG_SPECIAL, "FERRY for %s", box_name(player(c.who)))

	p_magic(ship).ferry_flag = TRUE

	return TRUE
}

func v_tell(c *command) int {
	wout(c.who, "The TELL order has been removed as of turn 50.")
	return FALSE

	//target := c.a;
	//what := c.b;
	//if (numargs(c) < 2) {
	//    wout(c.who, "Usage: TELL <who> <what>");
	//    return FALSE;
	//}
	//
	//if (kind(target) != T_char && kind(target) != T_player) {
	//    wout(c.who, "%s is not a character or player.",
	//         c.parse[1]);
	//
	//    return FALSE;
	//}
	//
	//if (i_strcmp(c.parse[2], "all") == 0) {
	//    what = -1;
	//} else if (!valid_box(what) || !test_known(c.who, what)) {
	//    wout(c.who, "%s doesn't know anything about %s.",
	//         box_name(c.who), c.parse[2]);
	//    return FALSE;
	//}
	//
	//if (what > 0 && !is_loc_or_ship(what)) {
	//    wout(c.who, "Can only TELL about hidden locations.");
	//    return FALSE;
	//}
	//
	//if (what == -1) {
	//    wout(c.who, "Not implemented.");
	//    return FALSE;
	//}
	//
	//set_known(target, what);
	//
	//wout(c.who, "Told %s about %s.", box_name(target), box_name(what));
	//wout(target, "%s told us about %s.", box_name(c.who), box_name(what));
	//
	//return TRUE;
}

/*
 *  Religion
 *  Wed Aug 14 14:39:05 1996 -- Scott Turner
 *
 *  Dedicating yourself.
 *
 */
func v_dedicate(c *command) int {
	/*
	 *  Priests can't dedicate themselves.
	 *
	 */
	if is_priest(c.who) != FALSE {
		wout(c.who, "Priests are already dedicated to the church.")
		return FALSE
	}
	/*
	 *  Are you a magic user?
	 *
	 *  Tue May 23 06:44:37 2000 -- Scott Turner
	 *
	 *  Should also stop anyone who has partially learned magic, sigh.
	 *
	 */
	if is_wizard(c.who) != FALSE {
		wout(c.who, "Wizards may not dedicate themselves to the church.")
		return FALSE
	}

	/*
	 *  Already dedicated?
	 *
	 */
	if rp_char(c.who).religion.priest != FALSE {
		wout(c.who, "You are already dedicated to %s.",
			box_name(rp_char(c.who).religion.priest))
		return FALSE
	}
	/*
	 *  Did you name a priest?
	 *
	 */
	if !valid_box(c.a) || is_priest(c.a) == FALSE || subloc(c.a) != subloc(c.who) {
		wout(c.who, "You must be with a priest to dedicate yourself.")
		return FALSE
	}
	/*
	 *  You must be in the same location as the priest.
	 *
	 */
	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must spend the entire dedication time with %s.",
			box_name(c.a))
		return FALSE
	}
	return TRUE
}

func d_dedicate(c *command) int {
	/*
	 *  Priests can't dedicate themselves.
	 *
	 */
	if is_priest(c.who) != FALSE {
		wout(c.who, "Priests are already dedicated to the church.")
		return FALSE
	}
	/*
	 *  Already dedicated?
	 *
	 */
	if rp_char(c.who).religion.priest != FALSE {
		wout(c.who, "You are already dedicated to %s.",
			box_name(rp_char(c.who).religion.priest))
		return FALSE
	}
	/*
	 *  Did you name a priest?
	 *
	 */
	if !valid_box(c.a) || is_priest(c.a) == FALSE || subloc(c.a) != subloc(c.who) {
		wout(c.who, "You must be with a priest to dedicate yourself.")
		return FALSE
	}
	/*
	 *  You must be in the same location as the priest.
	 *
	 */
	if subloc(c.who) != subloc(c.a) {
		wout(c.who, "You must spend the entire dedication time with %s.",
			box_name(c.a))
		return FALSE
	}
	/*
	 *  Otherwise success
	 *
	 */
	wout(c.who, "You are now dedicated to %s of %s.",
		box_name(c.a), box_name(is_priest(c.a)))
	wout(c.a, "%s has dedicated himself to you.", box_name(c.who))
	rp_char(c.who).religion.priest = c.a
	rp_char(c.a).religion.followers = append(rp_char(c.a).religion.followers, c.who)

	return TRUE
}

/*
 *  Fri Nov 22 15:21:05 1996 -- Scott Turner
 *
 *  Tax on a market in a city.
 *
 */
func do_market_tax(c *command) int {
	city := c.a
	rate := c.c

	if !(valid_box(city) &&
		(subkind(city) == sub_city || is_guild(city) == sk_trading)) {
		wout(c.who, "Must name a city or trading guild.", box_name(c.a))
		return FALSE
	}

	if player(province_admin(province(city))) != player(c.who) {
		wout(c.who, "You may only tax areas you rule.")
		return FALSE
	}

	if rate < 0 {
		wout(c.who, "The market tax should be from 0 to 25.")
		wout(c.who, "I'll assume you meant 'tax 0'.")
		rate = 0
	}

	if rate > 25 {
		wout(c.who, "The market tax should be from 0 to 25.")
		wout(c.who, "I'll assume you meant 'tax 25'.")
		rate = 25
	}

	rp_subloc(city).tax_market2 = rate
	wout(c.who, "%s market tax rate will be set to %d%% at the end of the month.",
		box_name(city), rate)
	return TRUE

}

/*
 *  Sun Dec  1 19:13:45 1996 -- Scott Turner
 *
 *  Province tax.
 *
 */
func do_province_tax(c *command) int {
	province := c.a
	rate := c.c

	if !valid_box(province) || kind(province) != T_loc ||
		loc_depth(province) != LOC_province {
		wout(c.who, "%s is not a province.", box_name(province))
		return FALSE
	}

	if player(province_admin(province)) != player(c.who) {
		wout(c.who, "You may only tax areas you rule.")
		return FALSE
	}

	if rate < 0 {
		wout(c.who, "The province tax should be from 0 to 10.")
		wout(c.who, "I'll assume you meant 'tax 0'.")
		rate = 0
	}

	if rate > 10 {
		wout(c.who, "The province tax should be from 0 to 10.")
		wout(c.who, "I'll assume you meant 'tax 10'.")
		rate = 10
	}

	if FALSE == garrison_here(province) {
		wout(c.who, "Can only tax garrisoned provinces.")
		return FALSE
	}

	rp_loc(province).tax_rate = rate
	wout(c.who, "%s tax rate set to %d%%.",
		box_name(province), rate)
	wout(garrison_here(province), "%s set the province tax rate to %d%%.",
		box_name(c.who), rate)
	return TRUE

}

/*
 *  Fri Nov 22 15:16:00 1996 -- Scott Turner
 *
 *  Set tax rates.
 *
 */
func v_tax(c *command) int {
	where := c.a

	if subkind(where) != sub_city &&
		subkind(where) != sub_guild &&
		subkind(where) != sub_forest &&
		subkind(where) != sub_plain &&
		subkind(where) != sub_swamp &&
		subkind(where) != sub_mountain &&
		subkind(where) != sub_desert &&
		subkind(where) != sub_island {
		wout(c.who, "%s is not a valid city, province, or trading guild.",
			box_code(where))
		return FALSE
	} else if subkind(where) == sub_city || is_guild(where) == sk_trading {
		return do_market_tax(c)
	} else if loc_depth(where) == LOC_province {
		return do_province_tax(c)
	} else {
		wout(c.who, "%s is not a valid city, province or trading guild.",
			box_name(where))
		return FALSE
	}
}
