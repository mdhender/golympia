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
	"os"
)

func immed_commands() {
	log.Printf("Olympia immediate mode\n")

	immediate = gm_player
	out(immediate, "You are now %s.", box_name(immediate))
	init_locs_touched()

	show_day = true

	for {
		c := p_command(immediate)
		c.who = immediate
		c.wait = 0

		// todo: how to flush stdout?
		log.Printf("%d> ", immediate)

		line := getlin(os.Stdin)
		if line == nil {
			break
		}
		buf := str_save(line)

		if !oly_parse(c, buf) {
			log.Printf("Unrecognized command.\n")
			continue
		}
		if c.fuzzy {
			out(immediate, "(assuming you meant '%s')", cmd_tbl[c.cmd].name)
		}

		c.pri = cmd_tbl[c.cmd].pri
		c.wait = cmd_tbl[c.cmd].time
		c.poll = cmd_tbl[c.cmd].poll
		c.days_executing = 0
		c.state = LOAD

		do_command(c)

		for c.state == RUN {
			evening = true
			finish_command(c)
			evening = false
			olytime_increment(&sysclock)
		}

		// ifndef NEW_TRADE
		//   if len(trades_to_check) > 0 {
		//	   check_validated_trades()
		//   }

		// show_day = false
	}

	log.Println("")
}

func v_add_item(c *command) int {
	if kind(c.a) == T_item {
		if kind(c.who) != T_char {
			out(c.who, "Warning: %s not a character", box_name(c.who))
		}
		gen_item(c.who, c.a, c.b)
		return TRUE
	}
	wout(c.who, "%s is not a valid item.", c.parse[1])
	return FALSE
}

func v_be(c *command) int {
	if valid_box(c.a) {
		immediate = c.a
		out(immediate, "You are now %s.", box_name(c.a))
		return TRUE
	}
	out(c.who, "'%s' not a valid box.", c.parse[1])
	return FALSE
}

func v_dump(c *command) int {
	if valid_box(c.a) {
		bx[c.a].temp = 0
		save_box(os.Stdout, c.a)
		return TRUE
	}
	return FALSE
}

func v_listcmds(c *command) int {
	indent += 4

	var buf []byte
	for i := 1; cmd_tbl[i].name != ""; i++ {
		buf = append(buf, []byte(fmt.Sprintf("%-12s", cmd_tbl[i].name))...)
		if i%5 == 0 {
			out(c.who, "%s", string(buf))
			buf = nil
		}
	}
	if len(buf) != 0 {
		out(c.who, "%s", string(buf))
	}
	indent -= 4

	return TRUE
}

func v_poof(c *command) int {
	if !is_loc_or_ship(c.a) {
		wout(c.who, "%s is not a location.", c.parse[1])
		return FALSE
	}
	move_stack(c.who, c.a)
	wout(c.who, ">poof!<  A cloud of orange smoke appears and whisks you away...")
	out(c.who, "")
	show_loc(c.who, loc(c.who))
	return TRUE
}

func v_see_all(c *command) int {
	if len(c.parse[1]) == 0 {
		immed_see_all = 1
	} else {
		immed_see_all = c.a
	}
	if immed_see_all != FALSE {
		out(c.who, "Will reveal all hidden features.")
	} else {
		out(c.who, "Hidden features will operate normally.")
	}
	return TRUE
}

func v_sub_item(c *command) int {
	consume_item(c.who, c.a, c.b)
	return TRUE
}

func v_makeloc(c *command) int {
	sk := lookup_sb(subkind_s, c.parse[1])
	if sk < 0 {
		wout(c.who, "Unknown subkind.")
		return FALSE
	}
	var kind int
	if sk == sub_galley || sk == sub_roundship {
		kind = T_ship
	} else {
		kind = T_loc
	}

	n := new_ent(kind, sk)
	if n < 0 {
		wout(c.who, "Out of boxes.")
		return FALSE
	}

	set_where(n, subloc(c.who))
	wout(c.who, "Created %s", box_name(n))
	if sk == sub_temple {
		p_subloc(n).teaches = append(p_subloc(n).teaches, sk_religion)
	}

	return TRUE
}

func v_invent(c *command) int {
	show_char_inventory(c.who, c.who, "")
	show_carry_capacity(c.who, c.who)
	show_item_skills(c.who, c.who)
	return TRUE
}

func v_know(c *command) int {
	if kind(c.a) != T_skill {
		wout(c.who, "%s is not a skill.", c.parse[1])
		return FALSE
	}
	learn_skill(c.who, c.a)
	// set_skill(c.who, c.a, SKILL_know);
	return TRUE
}

func v_skills(c *command) int {
	list_skills(c.who, c.who, "")
	list_partial_skills(c.who, c.who, "")
	return TRUE
}

func v_save(c *command) int {
	save_db()
	return TRUE
}

func v_los(c *command) int {
	target := c.a
	if !is_loc_or_ship(target) {
		wout(c.who, "%s is not a location.", box_code(target))
		return FALSE
	}
	d := los_province_distance(subloc(c.who), target)
	wout(c.who, "distance=%d", d)
	return TRUE
}

func v_kill(c *command) int {
	kill_char(c.a, MATES, S_body)
	return TRUE
}

func v_take_pris(c *command) int {
	if !check_char_here(c.who, c.a) {
		return FALSE
	}
	take_prisoner(c.who, c.a)
	return TRUE
}

func v_seed(c *command) int {
	seed_initial_locations()
	return TRUE
}

func v_postproc(c *command) int {
	for _, i := range loop_char() {
		ch := rp_char(i)
		if ch != nil {
			ch.studied = 0
		}
		for _, e := range loop_char_skill(i) {
			e.exp_this_month = FALSE
		}
	}

	post_month()
	olytime_turn_change(&sysclock)
	return TRUE
}

func v_lore(c *command) int {
	if valid_box(c.a) {
		deliver_lore(c.who, c.a)
	}
	return TRUE
}

/*
 *  Clear city trades
 */
func v_ct(c *command) int {
	for _, i := range loop_loc() {
		if subkind(i) == sub_city {
			bx[i].trades = nil
		}
	}
	update_markets()
	return TRUE
}

func v_seedmarket(c *command) int {
	for _, i := range loop_city() {
		seed_city_trade(i)
	}
	return TRUE /* ??? */
	// todo: huh, what?
	seed_common_tradegoods()
	seed_rare_tradegoods()
	for _, i := range loop_city() {
		do_production(i, true)
	}
	return TRUE
}

/*
 *  credit <who> <amount> <what -- defaults to gold, can be np>
 *
 */
func v_credit(c *command) int {
	target, amount, item := c.a, c.b, c.c
	if amount == 0 {
		wout(c.who, "You didn't specify an amount and/or item.")
		return FALSE
	}
	if kind(target) != T_char && kind(target) != T_player {
		wout(c.who, "%s not a character or player.", c.parse[1])
		return FALSE
	}

	if numargs(c) >= 3 && i_strcmp(c.parse[3], []byte("np")) == 0 {
		if kind(target) != T_player {
			wout(c.who, "%s not a player.", box_code(target))
			return FALSE
		}
		add_np(target, amount)
		wout(c.who, "Credited %s %d NP.", box_name(target), amount)
		wout(target, "Received GM credit of %d NP.", amount)
		return TRUE
	}

	if item == 0 {
		item = item_gold
	}

	gen_item(target, item, amount)
	wout(c.who, "Credited %s %s.", box_name(target), box_name_qty(item, amount))
	wout(target, "Received CLAIM credit of %s.", box_name_qty(item, amount))
	return TRUE
}

func v_relore(c *command) int {
	skill := c.a
	if !valid_box(skill) || kind(skill) != T_skill {
		wout(c.who, "%s is not a skill.", c.parse[1])
		return FALSE
	}

	for _, i := range loop_char() {
		if has_skill(i, skill) == 0 {
			queue_lore(i, skill, true)
		}
	}

	return TRUE
}

// todo: what does this do?
func v_xyzzy(c *command) int {
	if sysclock.turn != 13 {
		wout(c.who, "Only may be used on turn 13.")
		return FALSE
	}
	// item, targ := 50912, 27624;
	log_output(LOG_SPECIAL, "XYZZY")
	return TRUE
}

func v_fix2(c *command) int {
	for _, i := range loop_char() {
		if char_auraculum(i) != 0 {
			learn_skill(i, sk_adv_sorcery)
		}
	}
	return TRUE
}

func fix_gates() {
	clear_temps(T_loc)

	for _, where := range loop_province() {
		if !in_hades(where) && !in_clouds(where) && !in_faery(where) {
			continue
		} else if province_gate_here(where) == 0 {
			continue
		}
		log.Printf("Gate in %s\n", box_name(where))
		l := exits_from_loc_nsew(0, where)
		for i := 0; i < len(l); i++ {
			if loc_depth(l[i].destination) != LOC_province {
				continue
			}
			if province_gate_here(l[i].destination) == 0 {
				bx[l[i].destination].temp = 1
			}
		}
	}

	m := 1
	for {
		set_one := false

		for _, where := range loop_province() {
			if !in_hades(where) && !in_clouds(where) && !in_faery(where) {
				continue
			} else if province_gate_here(where) != 0 || bx[where].temp != m {
				continue
			}
			l := exits_from_loc_nsew(0, where)
			for i := 0; i < len(l); i++ {
				dest := l[i].destination
				if loc_depth(dest) != LOC_province {
					continue
				} else if province_gate_here(dest) == 0 && bx[dest].temp == 0 {
					bx[dest].temp = m + 1
					set_one = true
				}
			}
		}
		m++
		if set_one {
			continue
		}
		break
	}

	for _, where := range loop_province() {
		if !in_hades(where) && !in_clouds(where) && !in_faery(where) {
			continue
		} else if province_gate_here(where) == 0 && bx[where].temp < 1 {
			log.Printf("(1)error on %d\n", where)
		}
	}

	for _, where := range loop_province() {
		if !in_hades(where) && !in_clouds(where) && !in_faery(where) {
			continue
		}
		p_loc(where).dist_from_gate = bx[where].temp
	}
}

func v_fix(c *command) int {
	for _, i := range loop_char() {
		if has_skill(i, sk_trance) != 0 && has_skill(i, sk_quick_cast) != 0 {
			wout(c.who, "%s (%s)", box_name(i), box_name(player(i)))
			if has_skill(i, sk_aura_blast) != 0 {
				wout(c.who, "   ...also has aura blast")
			}
		}
	}
	return TRUE
}
