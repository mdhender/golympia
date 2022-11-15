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
	"bytes"
	"log"
	"sort"
	"strings"
)

const (
	MAX_PRI = 5
)

var (
	cur_pri   = 0
	wait_list []int

	auto_attack_flag = false
	month_done       = false

	load_q [MAX_PRI]queue_l
	run_q  queue_l
)

/*
 *  Line parser.  Given an ilist of char *'s and a line,
 *  returns the ilist with new slices copied from the cut-up line.
 *
 */
func parse_line(s []byte) [][]byte {
	var l [][]byte

	for {
		s = bytes.TrimSpace(s) /* eat whitespace */

		if len(s) == 0 {
			break
		}

		var prev []byte
		if s[0] == '"' || s[0] == '\'' { /* handle 'single' or "double" quoted text */
			ch := s[0] // save the quote so that we know when to terminate
			s = s[1:]  // skip the quote
			for len(s) != 0 && s[0] != ch {
				prev = append(prev, s[0])
				s = s[1:]
			}
			// strip leading and trailing whitespace from quoted argument
			prev = bytes.TrimSpace(prev)
		} else { /* unquoted argument */
			for len(s) != 0 && !iswhite(s[0]) {
				prev = append(prev, s[0])
				s = s[1:]
			}
		}

		l = append(l, prev)
	}

	return l
}

func remove_comment(s []byte) []byte {
	var qMark byte // set to non-zero if we're in quoted text
	for i := 0; i < len(s); i++ {
		if qMark != 0 {
			if s[i] == '"' || s[i] == '\'' {
				qMark = 0
			}
		} else if s[i] == '"' || s[i] == '\'' {
			qMark = s[i]
		} else if s[i] == '#' {
			return s[:i]
		}
	}
	return s
}

func remove_ctrl_chars(s []byte) []byte {
	for i := 0; i < len(s); i++ {
		if (s[i] & 0x80) != 0 {
			s[i] = s[i] & 0x7F
		}
		if s[i] < 32 {
			s[i] = ' '
		}
	}
	return s
}

func parse_arg(who int, s []byte) int {
	/*
	 *  scode() will perform an atoi() if the string is a digits-only
	 *  number, or a code_to_int() if the string is a location code.
	 */
	n := scode(s)
	if n < 0 {
		n = 0
	}
	if n == 0 && who != 0 && subloc(who) != 0 && fuzzy_strcmp(s, []byte("garrison")) {
		if n = garrison_here(subloc(who)); n == 0 {
			n = garrison_magic
		}
	}
	return n
}

var fuzzy_find bool /* was string found with a fuzzy match? */

func find_command(s []byte) int {
	fuzzy_find = false
	if len(s) == 0 {
		return -1
	}
	for i := 1; cmd_tbl[i].name != ""; i++ {
		if i_strcmp([]byte(cmd_tbl[i].name), s) == 0 {
			return i
		}
	}
	for i := 1; cmd_tbl[i].name != ""; i++ {
		if fuzzy_strcmp([]byte(cmd_tbl[i].name), s) {
			fuzzy_find = true
			return i
		}
	}
	return -1
}

/*
 *  Call full oly_parse instead.  Don't call this.
 */
func oly_parse_cmd(c *command, s []byte) bool {
	c.cmd = 0
	c.fuzzy = false
	c.use_skill = 0

	s = bytes.TrimSpace(remove_comment(remove_ctrl_chars(s)))

	c.line = string(s)
	if len(s) != 0 {
		if s[0] == '&' {
			c.conditional = 1
			s = s[1:]
		} else if s[0] == '?' {
			c.conditional = 2
			s = s[1:]
		} else {
			c.conditional = 0
		}
	}
	c.parsed_line = str_save(s)

	c.parse = parse_line(c.parsed_line)

	if len(c.parse) != 0 {
		i := find_command(c.parse[0])
		if i < 0 {
			return false
		}
		c.cmd = i
		if fuzzy_find {
			c.fuzzy = true
		}
	}

	return true
}

/*
 *  Look up command and tokenize arguments
 */

func oly_parse(c *command, s []byte) bool {
	c.a, c.b, c.c, c.d, c.e, c.f, c.g, c.h = 0, 0, 0, 0, 0, 0, 0, 0
	if !oly_parse_cmd(c, s) {
		return false
	}

	switch min(len(c.parse), 9) {
	case 9:
		c.h = parse_arg(c.who, c.parse[8])
		fallthrough
	case 8:
		c.g = parse_arg(c.who, c.parse[7])
		fallthrough
	case 7:
		c.f = parse_arg(c.who, c.parse[6])
		fallthrough
	case 6:
		c.e = parse_arg(c.who, c.parse[5])
		fallthrough
	case 5:
		c.d = parse_arg(c.who, c.parse[4])
		fallthrough
	case 4:
		c.c = parse_arg(c.who, c.parse[3])
		fallthrough
	case 3:
		c.b = parse_arg(c.who, c.parse[2])
		fallthrough
	case 2:
		c.a = parse_arg(c.who, c.parse[1])
	}

	return true
}

// todo: unravel the byte/string stuff
func oly_parse_s(c *command, s string) bool {
	return oly_parse(c, []byte(s))
}

func cmd_shift(c *command) {
	if len(c.parse) > 1 {
		/*
		 *  Deleted argument need not be freed, since it's just a
		 *  pointer into another string.  It was never allocated
		 *  itself.
		 */
		c.parse = c.parse.delete(1)
	}

	// shift arguments
	c.a, c.b, c.c, c.d, c.e, c.f, c.g, c.h = c.b, c.c, c.d, c.e, c.f, c.g, c.h, 0
	if numargs(c) >= 8 {
		c.h = parse_arg(c.who, c.parse[8])
	}
}

func check_allow(c *command, allow []byte) bool {
	if len(allow) == 0 {
		return true
	} else if immediate != FALSE && bytes.LastIndexByte(allow, 'i') != -1 {
		return true
	}

	//#if 0
	//    if (allow == nil || immediate)
	//        return TRUE;			/* don't check */
	//#endif

	var t byte
	switch bx[c.who].kind {
	case T_player:
		t = 'p'
		break

	case T_char:
		t = restricted_control(c.who)
		if t == 0 {
			t = 'c'
		}
		break

	default:
		log.Printf("check_allow: bad kind: %s\n", box_name(c.who))
		assert(false)
	}

	if bytes.IndexByte(allow, 'm') != -1 && player(c.who) == gm_player {
		return true
	} else if bytes.IndexByte(allow, t) == -1 {
		wout(c.who, "%s may not issue that order.", box_name(c.who))
		return false
	}

	return true
}

func get_command(who int) []byte {
	if !valid_box(who) || kind(who) == T_deadchar || subkind(who) == sub_dead_body {
		return nil
	}

	// who controls us now?
	fact := player(who)

	// if we don't have any orders, then fail.
	order := top_order(fact, who)
	if len(order) == 0 {
		return nil
	} else if len(order) >= LEN {
		panic("assert(len(order) < LEN)")
	}

	// update the player's last turn field so we know if he's been playing or not.
	p_player(fact).last_order_turn = sysclock.turn

	pop_order(fact, who)

	// return a copy of the command
	return append([]byte{}, order...)
}

func exec_precedence(who int) int {
	if kind(who) != T_char {
		return 0
	}

	stack_depth := 0
	for n := who; n != 0; n = stack_parent(n) {
		stack_depth++
	}
	if !(stack_depth < 100) {
		panic("assert(stack_depth < 100)")
	}

	pos := here_pos(who)
	if !(pos < 10000) {
		panic("assert(pos < 10000)")
	}

	return stack_depth*10000 + pos
}

func exec_comp(a, b *int) int {
	return bx[*a].temp - bx[*b].temp
}

func sort_load_queue(q queue_l) queue_l {
	for i := 0; i < len(q); i++ {
		bx[q[i]].temp = exec_precedence(q[i])
	}
	sort.Sort(q)
	return q
}

// we want priority n-1 commands to finish before priority n commands.
// we must fish the priority out of the command structure and give it msb status in the sorting tag.
func sort_run_queue(q queue_l) queue_l {
	for i := 0; i < len(q); i++ {
		c := rp_command(q[i])
		if c == nil {
			panic("assert(c != nil)")
		}
		bx[q[i]].temp = c.pri*1_000_000 + exec_precedence(q[i])
	}
	sort.Sort(q)
	return q
}

func set_state(c *command, state, new_pri int) {
	switch c.state {
	case RUN:
		run_q = run_q.rem_value_uniq(c.who)
	case LOAD:
		load_q[new_pri] = load_q[new_pri].rem_value_uniq(c.who)
	}
	c.state = state
	switch c.state {
	case RUN:
		run_q = append(run_q, c.who)
	case LOAD:
		load_q[new_pri] = append(load_q[new_pri], c.who)
	}
}

/*
 *  Load the command structure with a new command.
 *  Returns TRUE or FALSE depending on whether a command was present
 *  to load.
 *
 *  Sets c.state:
 *
 *	DONE	no more commands remain in the queue
 *	LOAD	command loaded and ready to run
 *	ERROR	player command has an error.
 *
 *	Note that both LOAD and ERROR states are passed to do_command.
 *	do_commmand will report error states to the player.
 */
func load_command(c *command) bool {
	if !(c != nil) {
		panic("assert(c != nil)")
	}

	buf := get_command(c.who)
	if buf == nil {
		set_state(c, DONE, 0)
		return false
	}

	if !oly_parse(c, buf) {
		set_state(c, ERROR, 0)
		// return true since there was a command, even though it was invalid
		return true
	}

	pri := cmd_tbl[c.cmd].pri
	set_state(c, LOAD, pri)
	c.pri = pri
	c.wait = cmd_tbl[c.cmd].time
	c.poll = cmd_tbl[c.cmd].poll
	c.days_executing = 0
	return true
}

func command_done(c *command) {
	if immediate != FALSE {
		set_state(c, DONE, 0)
		return
	}
	if load_command(c) { /* sets c.state */
		if c.pri < cur_pri {
			cur_pri = c.pri
		}
	}
}

// finish_command can't be bool because it also returns the command status.
func finish_command(c *command) int {
	//extern int cmd_wait;

	if !(c != nil) {
		panic("assert(c != nil)")
	}
	if !(valid_box(c.who)) {
		panic("assert(valid_box(c.who))")
	}

	if kind(c.who) == T_deadchar {
		command_done(c)
		return FALSE
	}

	/*
	 *  Characters stacked under units engaged in movement have
	 *  their commands suspended until they get to their destination,
	 *  except for wait completion checks.
	 *
	 */

	if char_gone(c.who) != FALSE && stack_leader(c.who) != c.who && c.cmd != cmd_wait {
		return TRUE
	}

	if c.wait > 0 {
		c.wait--
	}

	if c.wait > 0 && c.debug == sysclock.day {
		out(c.who, "finish_command called twice, wait=%d", c.wait)
	}
	c.debug = sysclock.day

	/*
	 *  Call the finish routine once, when the command is done waiting,
	 *  or every evening if the poll flag is set.
	 */

	if c.wait <= 0 || c.poll != 0 {
		if cmd_tbl[c.cmd].finish != nil && !c.inhibit_finish {
			c.status = cmd_tbl[c.cmd].finish(c)
		}
	}

	if c.state == RUN && (c.status == FALSE || c.wait == 0) {
		command_done(c)
	}

	return c.status
}

/*
 *  Thu Oct 24 15:52:04 1996 -- Scott Turner
 *
 *  If there's any urchin spies on this guy, report his doings.
 *
 */
func do_urchin_spies(c *command) {
	for _, e := range loop_effects(c.who) {
		if e.type_ == ef_urchin_spy && e.data == subloc(c.who) && valid_box(e.subtype) && e.data == subloc(e.subtype) {
			wout(e.subtype, "An urchin reports that %s does \"%s\".", box_name(c.who), c.line)
		}
	}
}

func do_command(c *command) {
	if !(c != nil) {
		panic("assert(c != nil)")
	}

	if immediate == FALSE {
		if options.output_tags < 1 {
			out(c.who, "> %s", c.line)
		} else {
			rest := strings.Index(c.line, " ")
			if rest == -1 { // no spaces in line?
				out(c.who, "> <tag type=command name=%s>%s</tag type=command>", cmd_tbl[c.cmd].name, c.parse[0])
			} else {
				out(c.who, "> <tag type=command name=%s>%s</tag type=command>%s", cmd_tbl[c.cmd].name, c.parse[0], c.line[rest:])
			}
		}

		if c.fuzzy {
			out(c.who, "(assuming you meant '%s')",
				cmd_tbl[c.cmd].name)
		}
	}

	if c.state == ERROR {
		out(c.who, "Unrecognized command.")
		c.status = FALSE
	} else if !check_allow(c, []byte(cmd_tbl[c.cmd].allow)) {
		c.status = FALSE
	} else if cmd_tbl[c.cmd].start == nil {
		out(c.who, "Unimplemented command.")
		c.status = FALSE
	} else {
		/*
		 *  Increment count of commands started this turn
		 */
		p_player(player(c.who)).cmd_count++

		set_state(c, RUN, 0)

		assert(c.days_executing == 0)

		c.debug = 0
		c.inhibit_finish = false
		c.status = cmd_tbl[c.cmd].start(c)

		/*
		 *  Thu Oct 24 15:48:47 1996 -- Scott Turner
		 *
		 *  Let the urchins who are spying on you report
		 *  back to their masters.
		 *
		 */
		do_urchin_spies(c)
	}

	if c.status == FALSE {
		command_done(c)
	} else if c.wait == 0 && c.state == RUN {
		c.status = finish_command(c)
	}

	//#ifndef NEW_TRADE
	//    if (len(trades_to_check) > 0)
	//        check_validated_trades();
	//#endif
}

func init_load_sup(who int) {
	/*
	 *  All characters should have a command structure.
	 *  Create one if they don't.
	 */
	c := rp_command(who)
	if c == nil {
		c = p_command(who)
		c.who = who
		c.state = DONE
	}
	if !(who == c.who) {
		panic("assert(who == c.who)")
	}

	switch c.state {
	case LOAD:
		load_q[c.pri] = append(load_q[c.pri], c.who)
	case RUN:
		run_q = append(run_q, c.who)
	case DONE:
		load_command(c)
	default:
		panic("!reached")
	}
}

func initial_command_load() {
	for _, i := range loop_char() {
		init_load_sup(i)
	}
	for _, i := range loop_player() {
		init_load_sup(i)
	}
}

func min_pri_ready() int {
	for pri := 0; pri < MAX_PRI; pri++ {
		for i := 0; i < len(load_q[pri]); i++ {
			c := rp_command(load_q[pri][i])
			if !(c != nil) {
				panic("assert(c != nil)")
			}
			if !(c.state == LOAD) {
				panic("assert(c.state == LOAD)")
			}
			if !(c.pri == pri) {
				panic("assert(c.pri == pri)")
			}
			if !is_prisoner(c.who) && char_moving(c.who) == FALSE && c.second_wait == FALSE {
				return pri
			}
		}
	}

	return 99
}

func init_wait_list() {
	for _, i := range loop_char() {
		c := rp_command(i)
		if c != nil && c.state == RUN && c.cmd == cmd_wait {
			wait_list = append(wait_list, i)
		}
	}
}

func check_all_waits() {
	for i := 0; i < len(wait_list); i++ {
		c := rp_command(wait_list[i])
		if c != nil && c.state == RUN && c.cmd == cmd_wait {
			assert(c.wait == -1)
			finish_command(c)
		}
	}
}

func start_phase() {
	for {
		cur_pri = min_pri_ready()
		pri := cur_pri

		/*
		 *  Auto-attacks (declare hostile) should occur once per day.
		 *  We must make sure they don't run before a command of lesser
		 *  priority than attack.
		 */

		if auto_attack_flag && pri >= 3 {
			check_all_auto_attacks()
			auto_attack_flag = false
		}
		if pri == 99 {
			return
		}

		sort_load_queue(load_q[pri])
		l := load_q[pri].copy()
		for j := 0; j < len(l); j++ {
			i := l[j]
			c := rp_command(i)
			if !(c != nil) {
				panic("assert(c != nil)")
			}

			if c.state == LOAD && c.pri == pri && !is_prisoner(i) && char_moving(i) == FALSE && c.second_wait == FALSE {
				do_command(c)
				check_all_waits()
				if pri != cur_pri {
					break
				}
			}
		}
		l = nil
	}
}

func evening_phase() {
	evening = true
	run_q = sort_run_queue(run_q)

	l := run_q.copy()
	for j := 0; j < len(l); j++ {
		i := l[j]
		c := rp_command(i)
		if !(c != nil) {
			panic("assert(c != nil)")
		}
		if c.state != RUN {
			continue
		} else if c.second_wait != FALSE {
			continue
		}

		if !(c.state == RUN) {
			panic("assert(c.state == RUN)")
		}

		c.days_executing++
		finish_command(c)

		//#ifndef NEW_TRADE
		//        if (len(trades_to_check) > 0)
		//            check_validated_trades();
		//#endif // NEW_TRADE
	}
	l = nil

	check_all_waits()

	evening = false
}

/*
 *  repeat until we don't queue anything new:
 *  	for every character not doing something
 *  		check for a command, but avoid
 *		IDLE commands induced by wait events
 *
 *  check for satisfied wait events
 *
 *  repeat until we don't queue anything new:
 *  	for every character not doing something
 *		check for a command, including IDLE
 *
 *  "evening" command completion check phase
 */

func daily_command_loop() {
	auto_attack_flag = true
	start_phase()
	evening_phase()
	clear_second_waits()
}

func process_player_orders() {
	for _, pl := range loop_player() {
		c := rp_command(pl)
		if c == nil {
			continue
		}

		/*
		 *  pl can switch to T_deleted as a result of a quit order
		 */
		for kind(pl) == T_player && c.state == LOAD {
			do_command(c)

			//#if 0
			//                    /*
			//                     *  All player orders should be zero time
			//                     */
			//
			//                                assert(c.state != RUN);
			//#endif
		}
	}
}

/*
 *  Interrupt an executing order, if any.
 */
func interrupt_order(who int) {
	if stack_leader(who) == who { /* not moving anymore */
		restore_stack_actions(who)
	}

	c := rp_command(who)
	if c == nil {
		return
	}

	if c.state == RUN {
		if cmd_tbl[c.cmd].interrupt != nil {
			c.status = cmd_tbl[c.cmd].interrupt(c)
		}
		command_done(c)
		if !(c.state != RUN) {
			panic("assert(c.state != RUN)")
		}
	}
}

func process_interrupted_units() {
	for _, who := range loop_char() {
		if !stop_order(player(who), who) {
			continue
		}

		pop_order(player(who), who)
		out(who, "> stop")

		who_c := rp_command(who)
		if who_c != nil && who_c.cmd != 0 && who_c.wait != 0 {
			out(who, "Interrupt current order.")
			interrupt_order(who)
		} else {
			out(who, "No order is currently executing.")
		}
	}
}

func process_orders() {
	stage("process_orders()")

	cmd_wait = find_command([]byte("wait"))
	assert(cmd_wait > 0)

	init_locs_touched()
	init_weather_views()
	olytime_turn_change(&sysclock)

	//#if 0
	//    player_accounting();
	//#endif

	init_wait_list()
	init_collect_list()
	queue_npc_orders()
	initial_command_load()
	ping_garrisons()
	//#if 0
	//    check_token_units();			/* XXX/NOTYET -- temp fix */
	//#endif

	process_interrupted_units() /* happens on day 0 */
	process_player_orders()
	scan_char_item_lore()

	stage("")

	for sysclock.day < MONTH_DAYS {
		olytime_increment(&sysclock)
		if sysclock.day == 1 {
			match_all_trades()
		}

		print_dot('.')

		daily_command_loop()
		daily_events()
	}
	log.Printf("\n")

	month_done = true
}
