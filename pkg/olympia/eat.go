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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MAX_ERR = 50

var (
	already_seen  = false
	cc_addr       string
	cmd_begin     = -1
	cmd_build     = -1
	cmd_email     = -1
	cmd_end       = -1
	cmd_format    = -1
	cmd_lore      = -1
	cmd_message   = -1
	cmd_notab     = -1
	cmd_passwd    = -1
	cmd_password  = -1
	cmd_players   = -1
	cmd_post      = -1
	cmd_press     = -1
	cmd_resend    = -1
	cmd_rumor     = -1
	cmd_set       = -1
	cmd_split     = -1
	cmd_stop      = -1
	cmd_unit      = -1
	cmd_vis_email = -1
	cmd_wait      = -1
	last_line     = -1
	line_count    = 0
	n_fail        = 0
	n_queued      = 0
	pl            = 0
	reply_addr    string
	save_line     []byte
	unit          = 0
	who_to        string
)

func find_command_s(s string) int { return find_command([]byte(s)) }

func find_meta_commands() {
	//extern int fuzzy_find;

	cmd_begin = find_command_s("begin")
	assert(cmd_begin > 0)
	assert(!fuzzy_find)

	cmd_end = find_command_s("end")
	assert(cmd_end > 0)
	assert(!fuzzy_find)

	cmd_unit = find_command_s("unit")
	assert(cmd_unit > 0)
	assert(!fuzzy_find)

	cmd_email = find_command_s("email")
	assert(cmd_email > 0)
	assert(!fuzzy_find)

	cmd_vis_email = find_command_s("vis_email")
	assert(cmd_vis_email > 0)
	assert(!fuzzy_find)

	cmd_lore = find_command_s("lore")
	assert(cmd_lore > 0)
	assert(!fuzzy_find)

	cmd_post = find_command_s("post")
	assert(cmd_post > 0)
	assert(!fuzzy_find)

	cmd_rumor = find_command_s("rumor")
	assert(cmd_rumor > 0)
	assert(!fuzzy_find)

	cmd_press = find_command_s("press")
	assert(cmd_press > 0)
	assert(!fuzzy_find)

	cmd_format = find_command_s("format")
	assert(cmd_format > 0)
	assert(!fuzzy_find)

	cmd_notab = find_command_s("notab")
	assert(cmd_notab > 0)
	assert(!fuzzy_find)

	cmd_message = find_command_s("message")
	assert(cmd_message > 0)
	assert(!fuzzy_find)

	cmd_resend = find_command_s("resend")
	assert(cmd_resend > 0)
	assert(!fuzzy_find)

	cmd_passwd = find_command_s("passwd")
	assert(cmd_passwd > 0)
	assert(!fuzzy_find)

	cmd_password = find_command_s("password")
	assert(cmd_password > 0)
	assert(!fuzzy_find)

	cmd_stop = find_command_s("stop")
	assert(cmd_stop > 0)
	assert(!fuzzy_find)

	cmd_players = find_command_s("players")
	assert(cmd_players > 0)
	assert(!fuzzy_find)

	cmd_split = find_command_s("split")
	assert(cmd_split > 0)
	assert(!fuzzy_find)

	cmd_wait = find_command_s("wait")
	assert(cmd_wait > 0)
	assert(!fuzzy_find)

	cmd_build = find_command_s("build")
	assert(cmd_build > 0)
	assert(!fuzzy_find)

	cmd_set = find_command_s("option")
	assert(cmd_set > 0)
	assert(!fuzzy_find)
}

func init_eat_vars() {

	if cmd_begin < 0 {
		find_meta_commands()
	}

	cc_addr = ""

	already_seen = false
	pl = 0
	unit = 0
	n_queued = 0
	n_fail = 0
	line_count = 0
	save_line = nil
}

// if the address contains "<.*>", return text
// otherwise trim address and return first word
//
//	rmatch("?*<(?*)>", s, &pat))
//	rmatch("(?*)[ \t]+\\(?**\\)", s, &pat))
func crack_address_sup(addr []byte) []byte {
	addr = bytes.TrimSpace(addr)
	if start := bytes.IndexByte(addr, '<'); start != -1 {
		if end := bytes.IndexByte(addr, '>'); end > start+1 {
			return addr[start+1 : end]
		}
	}
	if w := bytes.Fields(addr); len(w) != 0 {
		return str_save(w[0])
	}
	return str_save(addr)
}

func crack_address(addr []byte) []byte {
	if addr = crack_address_sup(addr); len(addr) != 0 {
		return addr
	}
	return nil
}

// returns address and pre-processor flag
//
//	false -> no pre-processing requested
//	 true -> something
func parse_reply(fp *os.File) ([]byte, bool) {
	s := getlin(fp)
	if i_strncmp(s, []byte("From "), 5) != 0 {
		return nil, false
	}
	s = s[5:] // skip the From:
	if sp := bytes.IndexAny(s, "\t "); sp != -1 {
		s = s[:sp]
	}
	if len(s) == 0 { // we did not get a From address line in the header
		return nil, false
	}

	var from_colon, reply_to []byte
	from_space := str_save(s)
	preFlag := false
	for s = getlin(fp); s != nil; s = getlin(fp) {
		if len(s) == 0 {
			break
		} else if i_strncmp(s, []byte("Subject:"), 8) == 0 {
			// subject line indicator asking for preprocessing?
			if bytes.Contains(s, []byte("cpp")) || bytes.Contains(s, []byte("preprocess")) {
				preFlag = true
			}
		} else if w := 5; i_strncmp(s, []byte("From:"), w) == 0 {
			from_colon = crack_address((s[w:]))
		} else if w = 9; i_strncmp(s, []byte("Reply-To:"), w) == 0 {
			reply_to = crack_address((s[w:]))
		} else if w = 6; i_strncmp(s, []byte("X-Loop"), w) == 0 {
			already_seen = true
		}
	}

	if reply_to != nil {
		return reply_to, preFlag
	} else if from_colon != nil {
		return from_colon, preFlag
	}
	return from_space, preFlag
}

func eat_line_2(fp *os.File, eat_white bool) []byte {
	var line []byte
	if eat_white {
		line = getlin_ew(fp)
	} else {
		line = getlin(fp)
	}
	if line == nil {
		return nil
	}
	remove_ctrl_chars(line)
	if eat_white {
		for len(line) != 0 && iswhite(line[0]) { // trim leading spaces
			line = line[1:]
		}
	}
	save_line = str_save(line)
	line_count++
	return line
}

func eat_next_line_sup(fp *os.File) []byte {
	line := getlin_ew(fp)
	if line == nil {
		return nil
	}
	remove_comment(line)
	remove_ctrl_chars(line)
	for len(line) != 0 && iswhite(line[0]) { // trim leading spaces
		line = line[1:]
	}
	save_line = str_save(line)
	line_count++
	return line
}

// consume empty lines
func eat_next_line(fp *os.File) []byte {
	for {
		if line := eat_next_line_sup(fp); !(line != nil && len(line) == 0) {
			return line
		}
	}
}

func err(k int, s string) {
	out_alt_who = k
	if k == EAT_ERR {
		n_fail++
	}
	if line_count < last_line {
		last_line = 0
	}
	if line_count > last_line {
		out(eat_pl, "line %d: %s: %q", line_count, s, string(save_line))
		last_line = line_count
	}
	indent += 3
	wiout(eat_pl, 2, "* %s", s)
	indent -= 3
}

func next_cmd(fp *os.File, c *command) {
	c.cmd = 0
	for {
		line := eat_next_line(fp)
		if line == nil {
			c.cmd = cmd_end
			return
		} else if !oly_parse(c, line) {
			err(EAT_ERR, "unrecognized command")
			if n_fail > MAX_ERR {
				err(EAT_ERR, "too many errors, aborting")
				c.cmd = cmd_end
				return
			}
			continue
		}

		if c.fuzzy {
			err(EAT_WARN, sout("assuming you meant '%s'", cmd_tbl[c.cmd].name))
		}

		return
	}
}

/*
 *  Tue Apr 17 12:23:03 2001 -- Scott Turner
 *
 *  How much should we charge someone?
 *
 */
func turn_charge(pl int) []byte {
	nps := 0
	for _, i := range loop_units(pl) {
		nps += nps_invested(i)
	}
	if nps <= options.free_np_limit {
		return str_save([]byte("0.00"))
	}
	return str_save([]byte(options.turn_charge))
}

func do_begin(c *command) bool {
	var pl_pass []byte

	if numargs(c) < 1 {
		err(EAT_ERR, "No player specified on BEGIN line")
		return false
	}

	if kind(c.a) != T_player {
		err(EAT_ERR, "No such player")
		return false
	}

	/*
	 *  Tue Apr 17 12:09:14 2001 -- Scott Turner
	 *
	 *  Check for a low balance, and reject this order set if
	 *  they can't afford to pay for the next turn.
	 *
	 */
	if options.check_balance != FALSE {
		cmd := sout("%s -p %s -g tag%d -T %s > /dev/null", options.accounting_prog, box_code_less(pl), game_number, turn_charge(pl))
		result := system(cmd)
		if result != -1 && result != 0 {
			err(EAT_ERR, "*********************************************************")
			err(EAT_ERR, "**                                                     **")
			err(EAT_ERR, "** Warning: Low account balance                        **")
			err(EAT_ERR, "**                                                     **")
			err(EAT_ERR, "** Your account has a low balance and you cannot       **")
			err(EAT_ERR, "** afford to pay for your next turn.                   **")
			err(EAT_ERR, "**                                                     **")
			err(EAT_ERR, "** See FIXME                                           **")
			err(EAT_ERR, "**                                                     **")
			err(EAT_ERR, "*********************************************************")
			return false
		}
	}

	pl_pass = str_save([]byte(p_player(c.a).Password))

	if numargs(c) > 1 {
		if len(pl_pass) == 0 {
			err(EAT_WARN, "No password is currently set")
		} else if i_strcmp(pl_pass, c.parse[2]) != 0 {
			err(EAT_ERR, "Incorrect password")
			return false
		}
	} else if len(pl_pass) != 0 {
		err(EAT_ERR, "Incorrect password")
		err(EAT_ERR, "Must give password on BEGIN line.")
		return false
	}

	pl = c.a

	p_player(pl).SentOrders = 1 /* okay, they sent something in */

	//#if 0
	//    /*
	//     *  We set unit here in case they forget the UNIT command for the
	//     *  player entity.  If they do, they lose the auto-flush ability,
	//     *  but at least their command will get queued, and echoed back
	//     *  in the confirmation.
	//     */
	//
	//    unit = pl;
	//#endif

	return true
}

func valid_char_or_player(who int) bool {
	if kind(who) == T_char || kind(who) == T_player {
		return true
	} else if kind(who) == T_item && subkind(who) == sub_dead_body {
		return true
	}
	return false
}

func do_unit(c *command) bool {
	unit = -1 /* ignore following unit commands */

	if pl == 0 {
		err(EAT_ERR, "BEGIN must appear before UNIT")
		out(eat_pl, "      rest of commands for unit ignored")
		return true
	}

	if kind(c.a) == T_unform {
		if ilist_lookup(p_player(pl).Unformed, c.a) < 0 {
			err(EAT_WARN, "Not an unformed unit of yours.")
		}
	} else if !valid_char_or_player(c.a) {
		err(EAT_ERR, "Not a character or unformed unit.")
		return true
	} else if player(c.a) != pl {
		err(EAT_WARN, "Not one of your controlled characters.")
	}

	unit = c.a
	flush_unit_orders(pl, unit)
	return true
}

func do_email(c *command) bool {

	if cc_addr != "" {
		err(EAT_ERR, "no more than one EMAIL order per message")
		out(eat_pl, "      new email address not set")
		return true
	}

	if pl == 0 {
		err(EAT_ERR, "BEGIN must come before EMAIL")
		out(eat_pl, "      new email address not set")
		return true
	}

	if numargs(c) < 1 || len(c.parse[1]) == 0 {
		err(EAT_ERR, "no new email address given")
		out(eat_pl, "      new email address not set")
		return true
	}

	cc_addr = rp_player(pl).EMail
	p_player(pl).EMail = string(c.parse[1])

	return true
}

func do_vis_email(c *command) bool {
	if pl == 0 {
		err(EAT_ERR, "BEGIN must come before VIS_EMAIL")
		out(eat_pl, "      new address not set")
		return true
	}

	if numargs(c) < 1 || len(c.parse[1]) == 0 {
		p_player(pl).VisEMail = ""
		return true
	}

	p_player(pl).VisEMail = string(c.parse[1])

	return true
}

func do_lore(c *command) bool {
	sheet := c.a

	if pl == 0 {
		err(EAT_ERR, "BEGIN must appear before LORE")
		return true
	}

	if kind(sheet) == T_item {
		sheet = item_lore(sheet)
	}

	if !valid_box(sheet) {
		err(EAT_ERR, "no such lore sheet")
		return true
	}

	if !test_known(pl, sheet) {
		err(EAT_ERR, "you haven't seen that lore sheet before")
		return true
	}

	out_alt_who = OUT_LORE
	deliver_lore(eat_pl, c.a)

	return true
}

func do_players(c *command) bool {
	fnam := filepath.Join(libdir, "save", fmt.Sprintf("%d", sysclock.turn), "players.html")
	fp, errr := fopen(fnam, "r")
	if errr != nil {
		err(EAT_ERR, sout("Sorry, couldn't find the player list."))
		return true
	}
	out_alt_who = EAT_PLAYERS
	for s := getlin(fp); s != nil; s = getlin(fp) {
		out(eat_pl, "%s", string(s))
	}

	fclose(fp)

	return true
}

func do_resend(c *command) bool {
	if pl == 0 {
		err(EAT_ERR, "BEGIN must appear before RESEND")
		return true
	}
	turn := c.a
	if turn == 0 {
		turn = sysclock.turn
	}
	if send_rep(pl, turn) != FALSE {
		out_alt_who = EAT_OKAY
		wout(eat_pl, "Turn %d report has been mailed to you in a separate message.", turn)
	} else {
		err(EAT_ERR, sout("Sorry, couldn't find your turn %d report", turn))
	}

	return true
}

func do_format(c *command) bool {
	if pl == 0 {
		err(EAT_ERR, "BEGIN must appear before FORMAT")
		return true
	}

	c.who = pl
	v_format(c)
	err(EAT_WARN, "Formatting set.")

	return true
}

func do_split(c *command) bool {
	lines := c.a
	chars := c.b

	if pl == 0 {
		err(EAT_ERR, "BEGIN must appear before SPLIT")
		return true
	}

	out_alt_who = EAT_OKAY

	if lines > 0 && lines < 500 {
		lines = 500
		out(eat_pl, "Minimum lines to split at is 500")
	}

	if chars > 0 && chars < 10000 {
		chars = 10000
		out(eat_pl, "Minimum bytes to split at is 10,000")
	}

	p_player(pl).SplitLines = lines
	p_player(pl).SplitBytes = chars

	if lines == 0 && chars == 0 {
		out(eat_pl, "Reports will not be split when mailed.")
	} else if lines != 0 && chars != 0 {
		out(eat_pl, "Reports will be split at %d lines or %d bytes, whichever limit is hit first.",
			lines, chars)
	} else if lines != 0 {
		out(eat_pl, "Reports will be split at %d lines.", lines)
	} else {
		out(eat_pl, "Reports will be split at %d bytes.", chars)
	}

	return true
}

func do_notab(c *command) bool {
	if pl == 0 {
		err(EAT_ERR, "BEGIN must appear before NOTAB")
		return true
	}

	p_player(pl).NoTab = c.a != FALSE

	out_alt_who = EAT_OKAY
	if c.a != FALSE {
		wout(eat_pl, "No TAB characters will appear in turn reports.")
	} else {
		wout(eat_pl, "TAB characters may appear in turn reports.")
	}

	return true
}

func do_password(c *command) bool {
	if pl == 0 {
		err(EAT_ERR, "BEGIN must appear before PASSWORD")
		return true
	}

	if numargs(c) < 1 || len(c.parse[1]) != 0 {
		p_player(pl).Password = ""

		out_alt_who = EAT_OKAY
		wout(eat_pl, "Password cleared.")
		return true
	}

	p_player(pl).Password = string(c.parse[1])

	out_alt_who = EAT_OKAY
	wout(eat_pl, "Password set to \"%s\".", c.parse[1])

	return true
}

func show_post(l [][]byte, cmd int) {
	var i int
	sav := out_alt_who

	out_alt_who = OUT_SHOW_POSTS

	for i = 0; i < len(l); i++ {
		if bytes.HasPrefix(l[i], []byte("=-=-")) {
			out(eat_pl, "> %s", l[i])
		} else {
			out(eat_pl, "%s", l[i])
		}
	}

	out(eat_pl, "")

	if cmd == cmd_press {
		attrib := strings.ReplaceAll(fmt.Sprintf("-- %s", box_name(player(unit))), "~", " ")
		out(eat_pl, "%55s", attrib)
		out(eat_pl, "")
	}

	out(eat_pl, "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

	out_alt_who = sav
}

/*
 *  Mon Dec 20 12:06:30 1999 -- Scott Turner
 *
 *  The "SET" command, which will incorporate a number of exiting
 *  commands.  It is intended to capture "meta-level" stuff, not
 *  actual noble actions.  All these should be "immediate" actions.
 *
 *      email
 *      format
 *  	notab
 *  	password
 *      vis_email
 * 	rules_url
 *	db_url
 *
 */
func do_set(c *command) bool {
	cmd := c.parse[1]

	/*
	 *  Shift the "set" out of the way.
	 *
	 */
	cmd_shift(c)

	if strcasecmp_bs(cmd, "email") == 0 {
		return do_email(c)
	} else if strcasecmp_bs(cmd, "format") == 0 {
		return v_format(c) != FALSE
	} else if strcasecmp_bs(cmd, "notab") == 0 {
		return v_notab(c) != FALSE
	} else if strcasecmp_bs(cmd, "password") == 0 {
		return do_password(c) /* ? */
	} else if strcasecmp_bs(cmd, "split") == 0 {
		return do_split(c) /* ? */
	} else if strcasecmp_bs(cmd, "vis_email") == 0 {
		return do_vis_email(c)
	} else if strcasecmp_bs(cmd, "rules_url") == 0 {
		return do_rules_url(c)
	} else if strcasecmp_bs(cmd, "db_url") == 0 {
		return do_db_url(c)
	} else {
		out_alt_who = EAT_OKAY
		wout(or_int(c.who != FALSE, c.who, eat_pl), "I don't know how to set option \"%s\".", cmd)
		return false
	}
}

/*
 *  Wed Dec 27 12:50:15 2000 -- Scott Turner
 *
 *  Command commenters.
 *
 */

/*
 *  It's safe to give a full name for something if it's
 *  either known to the player or a common item or skill.
 *
 */
func is_safe(who, n int) bool {
	/*
	 *  Might not be something.
	 *
	 */
	if !valid_box(n) {
		return false
	}
	/*
	 *  Common item.
	 *
	 */
	if kind(n) == T_item && FALSE == item_unique(n) {
		return true
	}
	/*
	 *  Skill
	 *
	 */
	if kind(n) == T_skill {
		return true
	}
	/*
	 *  Might be who
	 *
	 */
	if who == n {
		return true
	}
	/*
	 *  Might be a unit of this player.
	 *
	 */
	if kind(n) == T_char && player(who) != FALSE && is_unit(player(who), n) {
		return true
	}
	/*
	 *  Otherwise, if they "know" about it.
	 *
	 */
	if valid_box(who) && valid_box(n) && player(who) != FALSE {
		return test_known(who, n)
	}

	return false
}

func safe_name_qty(who, n, qty int, str []byte) string {
	if n == 0 {
		return string(str)
	} else if !is_safe(who, n) {
		if qty == 1 {
			return box_code(n)
		}
		return sout("%s %s", nice_num(qty), box_code(n))
	}
	return box_name_qty(n, qty)
}

func safe_name_s(who, n int, s string) string {
	return safe_name(who, n, []byte(s))
}

func safe_name(who, n int, b []byte) string {
	if n == 0 {
		return string(b)
	} else if !is_safe(who, n) {
		return box_code(n)
	}
	return box_name(n)
}

// todo: unwrap the byte/string
func buy_comment(c *command) string {
	if numargs(c) < 1 {
		return sout("probably incorrect")
	} else if numargs(c) == 1 || (numargs(c) == 2 && c.b == 0) || (numargs(c) == 3 && c.b == 0 && c.c == 0) {
		return sout("clear %s for %s", c.parse[0], safe_name(c.who, c.a, c.parse[1]))
	} else if numargs(c) == 2 {
		return sout("probably incorrect: specify an item, quantity and price")
	} else if numargs(c) == 3 {
		return sout("%s %s for %s gold each", c.parse[0], safe_name_qty(c.who, c.a, c.b, c.parse[1]), nice_num(c.c))
	} else {
		return sout("%s %s for %s gold each, keeping %s", c.parse[0], safe_name_qty(c.who, c.a, c.b, c.parse[1]), nice_num(c.c), nice_num(c.d))
	}
}

func drop_comment(c *command) string {
	if numargs(c) < 2 {
		return sout("probably incorrect")
	} else if numargs(c) < 3 {
		if c.b == 0 {
			return sout("drop all %s", safe_name(c.who, c.a, c.parse[1]))
		} else {
			return sout("drop %s %s", nice_num(c.b), safe_name(c.who, c.a, c.parse[1]))
		}
	} else {
		if c.b == 0 {
			return sout("drop all %s (keeping %s)", safe_name(c.who, c.a, c.parse[1]), nice_num(c.c))
		} else {
			return sout("drop %s %s (keeping %s)", nice_num(c.b), safe_name(c.who, c.a, c.parse[1]), nice_num(c.c))
		}
	}
}

func collect_comment(c *command) string {
	if numargs(c) < 1 {
		return sout("probably incorrect")
	} else if numargs(c) == 1 || (numargs(c) == 2 && c.b == 0) || (numargs(c) == 3 && c.b == 0 && c.c == 0) {
		return sout("collect all %s", safe_name(c.who, c.a, c.parse[1]))
	} else if numargs(c) == 2 || (numargs(c) == 3 && c.c == 0) {
		return sout("collect %s %s", nice_num(c.b), safe_name(c.who, c.a, c.parse[1]))
	} else if numargs(c) == 3 {
		return sout("collect %s %s for %s day%s", nice_num(c.b), safe_name(c.who, c.a, c.parse[1]), nice_num(c.c), or_string((c.c == 1), "", "s"))
	}
	return ""
}

func catch_comment(c *command) string {
	var cp command
	ret := oly_parse_s(&cp, sout("collect %d %d %d", item_wild_horse, c.a, c.b))
	assert(ret)
	return collect_comment(&cp)
}

func give_comment(c *command) string {
	if numargs(c) < 2 {
		return sout("probably incorrect")
	} else if numargs(c) < 3 {
		return sout("give all %s to %s",
			safe_name(c.who, c.b, c.parse[2]),
			safe_name(c.who, c.a, c.parse[1]))
	} else if numargs(c) < 4 {
		if c.c == 0 {
			return sout("give all %s to %s", safe_name(c.who, c.b, c.parse[2]), safe_name(c.who, c.a, c.parse[1]))
		} else {
			return sout("give %s to %s", safe_name_qty(c.who, c.b, c.c, c.parse[2]), safe_name(c.who, c.a, c.parse[1]))
		}
	} else {
		if c.c == 0 {
			return sout("give all %s to %s, keeping %s", safe_name(c.who, c.b, c.parse[2]), safe_name(c.who, c.a, c.parse[1]), nice_num(c.d))
		} else {
			return sout("give %s to %s, keeping %s", safe_name_qty(c.who, c.b, c.c, c.parse[2]), safe_name(c.who, c.a, c.parse[1]), nice_num(c.d))
		}
	}
}

func get_comment(c *command) string {
	if numargs(c) < 2 {
		return sout("probably incorrect")
	} else if numargs(c) < 3 {
		return sout("get all %s from %s", safe_name(c.who, c.b, c.parse[2]), safe_name(c.who, c.a, c.parse[1]))
	} else if numargs(c) < 4 {
		if c.c == 0 {
			return sout("get all %s from %s", safe_name(c.who, c.b, c.parse[2]), safe_name(c.who, c.a, c.parse[1]))
		} else {
			return sout("get %s from %s", safe_name_qty(c.who, c.b, c.c, c.parse[2]), safe_name(c.who, c.a, c.parse[1]))
		}
	} else {
		if c.c == 0 {
			return sout("get all %s from %s, leaving %s", safe_name(c.who, c.b, c.parse[2]), safe_name(c.who, c.a, c.parse[1]), nice_num(c.d))
		} else {
			return sout("get %s from %s, leaving %s", safe_name_qty(c.who, c.b, c.c, c.parse[2]), safe_name(c.who, c.a, c.parse[1]), nice_num(c.d))
		}
	}
}

func study_comment(c *command) string {
	if numargs(c) < 1 {
		return sout("probably incorrect")
	} else {
		return sout("%s %s", c.parse[0], safe_name(c.who, c.a, c.parse[1]))
	}
}

func admit_comment(c *command) string {
	if numargs(c) < 1 {
		return sout("probably incorrect")
	} else if numargs(c) == 1 {
		return sout("clear admits for %s", safe_name(c.who, c.a, c.parse[1]))
	} else {
		arg := sout("Admit to %s: %s", safe_name(c.who, c.a, c.parse[1]), safe_name(c.who, c.b, c.parse[2]))
		var args [6]int
		args[3] = c.c
		args[4] = c.d
		args[5] = c.e
		if strncasecmp_bs(c.parse[1], "all", len(c.parse[1])) == 0 {
			arg = sout("%s EXCEPT", arg)
		}
		for i := 3; i <= min(numargs(c), 5); i++ {
			arg = sout("%s, %s", arg, safe_name(c.who, args[i], c.parse[i]))
		}
		return arg
	}
	// return ""; // mdhender: not reached?
}

var first_admit_check = true

func admit_check(c *command) {
	if numargs(c) < 1 {
		return
	}
	/*
	 *  Look for a nation or monster name.
	 *
	 */
	if FALSE == find_nation_b(c.parse[1]) &&
		!fuzzy_strcmp_bs(c.parse[1], "monster") &&
		!fuzzy_strcmp_bs(c.parse[1], "monsters") &&
		!fuzzy_strcmp_bs(c.parse[1], "garrison") &&
		(!valid_box(c.a) || kind(c.a) != T_char) {
		/*
		 *  Possibly an "unformed" unit.
		 *
		 */
		if kind(c.a) == T_unform {
			err(EAT_WARN, sout("Note: %s is currently an unformed unit.", safe_name(c.who, c.a, c.parse[1])))
		} else {
			if first_admit_check {
				err(EAT_WARN, "The first argument to Admit should be a noble.  If you're trying to Admit someone to a location (such as a province or castle) then you should Admit them to the noble that controls that location, not the location itself.")
				first_admit_check = false
			} else {
				err(EAT_WARN, "The first argument to Admit should be a noble.")
			}
		}
	}
}

func quit_check(c *command) {
	moderator_email := "moderator@olytag.com" // todo: update the email from the environment!
	err(EAT_WARN, "This command will drop you from the game.")
	err(EAT_WARN, "If you are quitting, please send "+moderator_email+" a quick email to indicate why you are dropping.  Your feedback is important and helps to keep improving The Age of Gods.  Thanks.")
}

func accept_comment(c *command) string {
	if numargs(c) < 1 {
		return sout("probably incorrect")
	} else if numargs(c) == 1 {
		if strcasecmp_sb("clear", c.parse[1]) == 0 {
			return sout("clearing all accepts")
		}
		return sout("accept all from %s", safe_name(c.who, c.a, c.parse[1]))
	}
	var s1 string
	if c.a != FALSE {
		s1 = safe_name(c.who, c.a, c.parse[1])
	} else {
		s1 = "anyone"
	}
	var s2 string
	if c.b != FALSE {
		s2 = safe_name(c.who, c.b, c.parse[2])
	} else {
		s2 = "anything"
	}
	if numargs(c) == 2 {
		return sout("accept %s from %s", s2, s1)
	}
	return sout("accept up to %s %s from %s", nice_num(c.c), s2, s1)
}

func destination_name(c *command, i, num int) string {
	dir := lookup_sb(full_dir_s, c.parse[i])
	if dir < 0 {
		dir = lookup_sb(short_dir_s, c.parse[i])
		if dir < 0 {
			return sout("to %s", safe_name(c.who, num, c.parse[i]))
		}
	}
	return sout("%s", full_dir_s[dir])
}

func move_comment(c *command) string {
	var args [6]int
	args[1] = c.a
	args[2] = c.b
	args[3] = c.c
	args[4] = c.d

	if numargs(c) < 1 {
		return "probably incorrect"
	}
	arg := sout("%s %s", c.parse[0], destination_name(c, 1, args[1]))
	if numargs(c) > 1 {
		arg = sout("%s (or %s", arg, destination_name(c, 2, args[2]))
		for i := 3; i <= min(numargs(c), 4); i++ {
			arg = sout("%s, %s", arg, destination_name(c, i, args[i]))
		}
		if numargs(c) > 4 {
			arg = sout("%s...)", arg)
		} else {
			arg = sout("%s)", arg)
		}
	}
	return arg
}

func attack_comment(c *command) string {
	/*
	 *  Possibly a flag at the end.
	 *
	 */
	num, nomove := numargs(c), false
	if num > 1 && strcmp_sb("1", c.parse[num]) == 0 {
		nomove = true
		num--
	}
	if num < 1 {
		return "probably incorrect"
	}

	var args [6]int
	args[1] = c.a
	args[2] = c.b
	args[3] = c.c
	args[4] = c.d

	arg := sout("%s %s", c.parse[0], destination_name(c, 1, args[1]))
	if num > 1 {
		arg = sout("%s (or %s", arg, destination_name(c, 2, args[2]))
		for i := 3; i <= min(num, 4); i++ {
			arg = sout("%s, %s", arg, destination_name(c, i, args[i]))
		}
		if num > 4 {
			arg = sout("%s...)", arg)
		} else {
			arg = sout("%s)", arg)
		}
	}
	if nomove {
		arg = sout("%s (without entering)", arg)
	}
	return arg
}

func default_comment(c *command) string {
	if numargs(c) != 1 {
		return ""
	}
	return sout("%s %s", c.parse[0], safe_name(c.who, c.a, c.parse[1]))
}

func attitude_comment(c *command) string {
	if numargs(c) == 0 {
		return sout("Clear %s list.", c.parse[0])
	}
	arg := sout("%s to", c.parse[0])
	for i := 1; i <= numargs(c); i++ {
		arg = sout("%s %s", arg, safe_name(c.who, c.a, c.parse[i]))
	}
	return arg
}

/*
 *  Wed Dec 27 08:20:27 2000 -- Scott Turner
 *
 *  Check an argument against a "type".
 *
 */
func check_arg(c *command, i, t int) {
	if i < 1 || i > 5 {
		return
	}

	var args [6]int
	args[1] = c.a
	args[2] = c.b
	args[3] = c.c
	args[4] = c.d
	args[5] = c.e

	switch t {
	case CMD_undef:
		return
	case CMD_unit:
		/*
		 *  Look for a nation or monster name.
		 *
		 */
		if find_nation_b(c.parse[i]) != FALSE {
			return
		}
		if fuzzy_strcmp_bs(c.parse[i], "monster") {
			return
		}
		if fuzzy_strcmp_bs(c.parse[i], "monsters") {
			return
		}
		if fuzzy_strcmp_bs(c.parse[i], "garrison") {
			return
		}
		if fuzzy_strcmp_bs(c.parse[i], "garison") {
			return
		}
		/*
		 *  Otherwise, it should be a unit, player, etc.
		 *
		 */
		if valid_box(args[i]) && (kind(args[i]) == T_player || kind(args[i]) == T_char) {
			return
		}
		/*
		 *  Possibly an "unformed" unit.
		 *
		 */
		if valid_box(args[i]) && kind(args[i]) == T_unform {
			err(EAT_WARN, sout("Note: %s is currently an unformed unit.", safe_name(c.who, args[i], c.parse[i])))
			return
		}

		/*
		 *  Hmmm.
		 *
		 */
		err(EAT_WARN, sout("The %s argument of this command should be a unit (player, noble, or nation).  %s does not appear to be a unit.",
			ordinal(i), safe_name(c.who, args[i], c.parse[i])))
		return
	case CMD_item:
		if !valid_box(args[i]) || kind(args[i]) != T_item {
			err(EAT_WARN, sout("The %s argument of this command should be an item. %s does not appear to be an item.",
				ordinal(i), safe_name(c.who, args[i], c.parse[i])))
		}
		return
	case CMD_skill:
		if !valid_box(args[i]) || kind(args[i]) != T_skill {
			err(EAT_WARN, sout("The %s argument of this command should be a skill. %s does not appear to be an skill.",
				ordinal(i), safe_name(c.who, args[i], c.parse[i])))
		}
		return
	case CMD_days:
		if args[i] > 30 {
			err(EAT_WARN, "Note: This command may last more than a month.")
		}
		return
	case CMD_qty:
		return
	case CMD_gold:
		if unit != FALSE && has_item(unit, item_gold) < args[i] {
			err(EAT_WARN, sout("This command uses %s gold, and this unit currently has only %s gold.",
				nice_num(args[i]), nice_num(has_item(unit, item_gold))))
		}
		return
		/*
		 *  An item or skill to use.
		 *
		 */
	case CMD_use:
		if !valid_box(args[i]) ||
			(kind(args[i]) != T_skill && kind(args[i]) != T_item) {
			err(EAT_WARN, sout("The %s argument of this command should be a skill or an item. %s does not appear to be either.",
				ordinal(i), safe_name(c.who, args[i], c.parse[i])))
		} else if kind(args[i]) == T_skill && find_use_entry(args[i]) == -1 {
			err(EAT_WARN, sout("%s doesn't appear to be a skill you can 'use'.", safe_name(c.who, args[i], c.parse[i])))
		} else if kind(args[i]) == T_item && FALSE == item_unique(args[i]) {
			err(EAT_WARN, sout("%s doesn't appear to be an item you can 'use'.", safe_name(c.who, args[i], c.parse[i])))
		}
		return
		/*
		 *  A skill to practice.
		 *
		 */
	case CMD_practice:
		if !valid_box(args[i]) || kind(args[i]) != T_skill {
			err(EAT_WARN, sout("The %s argument of this command should be a skill. %s does not appear to be an skill.",
				ordinal(i), safe_name(c.who, args[i], c.parse[i])))
		} else if FALSE == has_skill(c.who, args[i]) {
			err(EAT_WARN, sout("Note: %s doesn't currently know %s.", box_name(c.who), safe_name(c.who, args[i], c.parse[i])))
		} else if FALSE == practice_time(args[i]) {
			err(EAT_WARN, sout("You cannot practice %s.", safe_name(c.who, args[i], c.parse[i])))
		} else if has_item(c.who, item_gold) < practice_cost(args[i]) {
			err(EAT_WARN, sout("Note: %s doesn't currently have the %s gold required to practice %s.",
				box_name(c.who), nice_num(practice_cost(args[i])), safe_name(c.who, args[i], c.parse[i])))
		}
		return
	default:
		return
	}
}

/*
 *  Wed Dec 27 07:41:27 2000 -- Scott Turner
 *
 *  Command checker.
 *
 */
func check_cmd(c *command) {
	/*
	 *  If we didn't get an entry into the command table,
	 *  we'd better bail.
	 *
	 */
	if c.cmd == FALSE {
		return
	}
	/*
	 *  Check to make sure this entity is allowed to use this
	 *  command, e.g., some commands cannot be issued by the player.
	 *
	 */
	if cmd_tbl[c.cmd].allow != "" && is_safe(pl, c.who) {
		var t byte
		switch bx[c.who].kind {
		case T_player:
			t = 'p'
			break
		case T_char:
			if t = restricted_control(c.who); t == 0 {
				t = 'c'
			}
			break
		default:
			t = 'c'
		}

		if strings.IndexByte(cmd_tbl[c.cmd].allow, t) == -1 &&
			!(strings.IndexByte(cmd_tbl[c.cmd].allow, 'm') != -1 && player(c.who) == gm_player) {
			err(EAT_WARN, sout("%s may not issue that order.",
				safe_name_s(c.who, c.who, "This unit")))
		}
	}

	/*
	 *  Required args?
	 *
	 */
	if cmd_tbl[c.cmd].num_args_required != 0 &&
		numargs(c) < cmd_tbl[c.cmd].num_args_required {
		err(EAT_WARN, sout("This command requires %s argument%s, and you've only provided %s.",
			nice_num(cmd_tbl[c.cmd].num_args_required), or_string(cmd_tbl[c.cmd].num_args_required == 1, "", "s"), nice_num(numargs(c))))
	}
	/*
	 *  Extra args?
	 *
	 */
	if cmd_tbl[c.cmd].max_args != 0 &&
		numargs(c) > cmd_tbl[c.cmd].max_args {
		err(EAT_WARN, sout("This command takes only %s argument%s, and you have %s.",
			nice_num(cmd_tbl[c.cmd].max_args), or_string((cmd_tbl[c.cmd].max_args == 1), "", "s"), nice_num(numargs(c))))
	}

	/*
	 *  Step through the arguments and do some checking on them.  We are
	 *  currently limited to the first five arguments.
	 *
	 */
	for i := 1; i <= min(numargs(c), 5); i++ {
		check_arg(c, i, cmd_tbl[c.cmd].arg_types[i-1])
	}

	/*
	 *  Possibly a specific error checker.
	 *
	 */
	if cmd_tbl[c.cmd].cmd_check != nil {
		cmd_tbl[c.cmd].cmd_check(c)
	}
}

func do_eat_command(c *command, fp *os.File) bool {

	assert(c.cmd != cmd_end)

	if c.cmd == cmd_begin {
		return do_begin(c)
	}
	if c.cmd == cmd_unit {
		return do_unit(c)
	}
	if c.cmd == cmd_email {
		return do_email(c)
	}
	if c.cmd == cmd_vis_email {
		return do_vis_email(c)
	}
	if c.cmd == cmd_lore {
		return do_lore(c)
	}
	if c.cmd == cmd_resend {
		return do_resend(c)
	}
	if c.cmd == cmd_format {
		return do_format(c)
	}
	if c.cmd == cmd_notab {
		return do_notab(c)
	}
	if c.cmd == cmd_split {
		return do_split(c)
	}
	if c.cmd == cmd_set {
		return do_set(c)
	}
	if c.cmd == cmd_players {
		return do_players(c)
	}
	if c.cmd == cmd_passwd || c.cmd == cmd_password {
		return do_password(c)
	}

	if unit == 0 {
		err(EAT_ERR, "can't queue orders, missing UNIT command")
		unit = -1
		return true
	}

	if unit == -1 {
		n_fail++
		return true
	}

	if c.cmd == cmd_stop {
		queue_stop(pl, unit)
	} else {
		queue_order(pl, unit, c.line)
	}
	n_queued++

	if c.cmd == cmd_wait {
		s := parse_wait_args(c)
		if len(s) != 0 {
			err(EAT_ERR, sout("Bad WAIT: %s", s))
		}
		clear_wait_parse(c)
	}

	if c.cmd == cmd_post || c.cmd == cmd_message || c.cmd == cmd_rumor || c.cmd == cmd_press {
		count := c.a
		reject_flag := false
		max_len := MAX_POST
		var l [][]byte

		if c.cmd == cmd_rumor || c.cmd == cmd_press {
			max_len = 78
		}

		for {
			s := eat_line_2(fp, c.cmd == cmd_post || c.cmd == cmd_message)
			if len(s) == 0 {
				err(EAT_ERR, "End of input reached before end of post!")
				break
			}

			length := len(s)
			if length > max_len {
				err(EAT_ERR, sout("Line length exceeds %d characters", max_len))
				reject_flag = true
			}

			queue_order(pl, unit, string(s))

			if count == 0 {
				t := eat_leading_trailing_whitespace(s)
				if i_strcmp(t, []byte("end")) == 0 {
					break
				}
				l = append(l, str_save(s))
			} else {
				l = append(l, str_save(s))
				count--
				if count <= 0 {
					break
				}
			}
		}

		if reject_flag {
			err(EAT_ERR, "Post will be rejected.")
		} else if c.cmd == cmd_press || c.cmd == cmd_rumor {
			show_post(l, c.cmd)
		}

		text_list_free(l)
	}

	/*
	 *  Wed Dec 27 07:40:41 2000 -- Scott Turner
	 *
	 *  Do some checking on the command, if possible.
	 *
	 */
	if unit != 0 && pl != 0 && player(unit) != 0 && player(unit) == pl {
		c.who = unit
		check_cmd(c)
	}

	return true
}

func parse_and_munch(fp *os.File) {
	c := &command{}
	first_admit_check = true
	next_cmd(fp, c)
	for c.cmd != cmd_end {
		if !do_eat_command(c, fp) {
			return
		}
		next_cmd(fp, c)
	}
}

func no_spaces(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\t", " "), " ", "~")
}

/*
 *  Thu Dec  2 12:18:28 1999 -- Scott Turner
 *
 *  Make sure there are no spaces in the mail header so that
 *  those lines don't get wrapped!
 *
 */
func eat_banner() {
	out_alt_who = OUT_BANNER

	out(eat_pl, "From: %s", no_spaces(from_host))
	out(eat_pl, "Reply-To: %s", no_spaces(reply_host))

	var to string
	if pl != 0 && len(rp_player(pl).EMail) != 0 {
		to = rp_player(pl).EMail
	} else {
		to = reply_addr
	}

	var full_name string
	if valid_box(pl) {
		if p := rp_player(pl); p != nil && len(p.FullName) != 0 {
			full_name = sout(" (%s)", p.FullName)
		}
	}

	if already_seen {
		to = "moderator@olytag.com" /*UPDATE*/
		full_name = " (Error Watcher)"
		cc_addr = ""
	}

	who_to = to

	out(eat_pl, "To:~%s%s", no_spaces(to), no_spaces(full_name))

	if cc_addr != "" {
		out(eat_pl, "Cc:~%s", no_spaces(cc_addr))
		who_to += (" ")
		who_to += (cc_addr)
	}

	out(eat_pl, "Subject: Acknowledge")
	out(eat_pl, "X-Loop: moderator@olytag.com")
	out(eat_pl, "Bcc: moderator@olytag.com") /* VLN: UPDATE */
	out(eat_pl, "")
	out(eat_pl, "     - Olympia order scanner -")
	out(eat_pl, "")
	if pl != 0 {
		out(eat_pl, "Hello, %s", box_name(pl))
	}
	out(eat_pl, "")
	/*
	 *  Tue Jan  2 07:32:17 2001 -- Scott Turner
	 *
	 *  Warning about low balances.
	 *
	 */
	if pl != 0 {

		cmd := sout("%s -p %s -g tag%d -T %s > /dev/null",
			options.accounting_prog,
			box_code_less(pl),
			game_number,
			options.turn_charge)
		result := system(cmd)

		if result != -1 && result != 0 {
			out(eat_pl, "*********************************************************")
			out(eat_pl, "**                                                     **")
			out(eat_pl, "** Warning: Low account balance                        **")
			out(eat_pl, "**                                                     **")
			out(eat_pl, "** Your account has a low balance and you cannot       **")
			out(eat_pl, "** afford to pay for your next turn.                   **")
			out(eat_pl, "**                                                     **")
			out(eat_pl, "**                                                     **")
			out(eat_pl, "*********************************************************")
			report_account_out(pl, eat_pl)
			out(eat_pl, "*********************************************************")
		}
	}

	out(eat_pl, "%d queued, %d error%s.", n_queued, n_fail, or_string(n_fail == 1, "", "s"))
}

func include_orig(fp *os.File) {
	out_alt_who = EAT_HEADERS

	// rewind(fp);
	if i, err := fp.Seek(0, 0); err != nil {
		panic(err)
	} else if i != 0 {
		panic("assert(fp.Seek(0,0) == 0)")
	}

	for s := getlin(fp); s != nil; s = getlin(fp) {
		/*
		 *  Tue Nov  7 08:23:19 2000 -- Scott Turner
		 *
		 *  Don't output the password in lines where it might appear.
		 *
		 */
		if bytes.Index(s, []byte("begin")) != -1 || bytes.Index(s, []byte("BEGIN")) != -1 {
			/*
			 *  begin ke4 "password"
			 *
			 *  Terminate at 2nd space.  Obviously can be easily fooled.
			 *
			 */
			i := 0
			for i < len(s) && !iswhite(s[i]) { // skip first word
				i++
			}
			for i < len(s) && iswhite(s[i]) { // skip first set of space
				i++
			}
			for i < len(s) && !iswhite(s[i]) { // skip second word
				i++
			}
			if i != 0 { // trim at this point
				s = s[:i]
			}
		} else if bytes.Index(s, []byte("password")) != -1 || bytes.Index(s, []byte("PASSWORD")) != -1 {
			// terminate at first space.
			if i := bytes.IndexByte(s, ' '); i != -1 {
				s = s[:i]
			} else if i = bytes.IndexByte(s, '\t'); i != -1 {
				s = s[:i]
			}
		}
		out(eat_pl, "%s", s)
	}
}

func show_pending() {
	out_alt_who = EAT_QUEUE
	orders_template(eat_pl, pl)
}

var eat_queue_mode = false

func eat(fnam string, mail_now bool) {
	ret := 0
	var fnam_cpp string

	init_eat_vars()
	eat_queue_mode = false

	fp, err := fopen(fnam, "r")
	if err != nil {
		log.Printf("can't open %s: %v", fnam, err)
		return
	}

	b, cpp := parse_reply(fp)
	reply_addr = string(b)

	if cpp {
		/*
		 *  We need to call the GCC preprocessor directly to avoid
		 *  the fact that GCC wants a .c ending on a source file.  We
		 *  want the GCC preprocessor because it let's us do -imacros.
		 *
		 */
		fclose(fp)
		buf := fmt.Sprintf("egrep -v '#include' %s > /tmp/cpp1.%d; %s -P -imacros %s/defines /tmp/cpp1.%d > /tmp/cpp.%d",
			fnam, get_process_id(), options.cpp, libdir, get_process_id(), get_process_id())
		system(buf)
		fnam_cpp = fmt.Sprintf("/tmp/cpp.%d", get_process_id())
		fp, err = fopen(fnam_cpp, "r")
		/*
		 *  Mon Mar  1 12:28:58 1999 -- Scott Turner
		 *
		 *  How to handle an error?
		 *
		 */
		if err != nil {
			log.Printf("Cannot open cpp file %s: %v\n", fnam_cpp, err)
			fp, err = fopen(fnam, "r")
			if err != nil {
				log.Printf("can't open %s: %v", fnam, err)
				return
			}
		}
		/*
		 *  Skip ahead to beginning of message.
		 *
		 */
		for s := getlin(fp); s != nil; s = getlin(fp) {
			if len(s) == 0 {
				break
			}
		}
	}

	if reply_addr != "" {
		if i_strncmp([]byte(reply_addr), []byte("postmaster"), 10) == 0 ||
			i_strncmp([]byte(reply_addr), []byte("mailer-daemon"), 13) == 0 ||
			i_strncmp([]byte(reply_addr), []byte("mail-daemon"), 11) == 0 {
			already_seen = true
		}

		unlink(sout("%s/log/%d", libdir, eat_pl))
		open_logfile_nondestruct()
		p_player(eat_pl).output = clear_know_rec(p_player(eat_pl).output)
		out_path = MASTER

		parse_and_munch(fp)
		eat_banner()

		if pl != 0 {
			show_pending()
		}

		include_orig(fp)

		out_alt_who = OUT_INCLUDE
		gen_include_sup(eat_pl) /* must be last */

		out_path = 0
		out_alt_who = 0

		if pl != 0 {
			unlink(sout("%s/orders/%d", libdir, pl))
			if err := save_player_orders(pl); err != nil {
				log.Printf("eat: %v\n", err)
			}
			unlink(sout("%s/fact/%d", libdir, pl))
			write_player(pl)
		}

		close_logfile()

		if mail_now {
			/* VLN		  ret = system(sout("rep %s/log/%d | sendmail -t", libdir, eat_pl)); */
			ret = system(sout("rep %s/log/%d | msmtp -t", libdir, eat_pl))
			if ret != 0 {
				log.Printf("error: couldn't mail ack to %s\n", who_to)
				/* VLN			log.Printf( "command was: %s\n", sout("rep %s/log/%d | sendmail -t", libdir, eat_pl)); */
				log.Printf("command was: %s\n", sout("rep %s/log/%d | msmtp -t", libdir, eat_pl))
				log.Printf("ret was = %d\n", ret)
			} else {
				// let's not overwhelm the system with a bunch of rapid-fire mail responses.
				time.Sleep(5 * time.Second)
			}
		}

	}

	fclose(fp)

	if cpp {
		var buf string
		buf = fmt.Sprintf("/tmp/cpp1.%d", get_process_id())
		unlink(buf)
		buf = fmt.Sprintf("/tmp/cpp.%d", get_process_id())
		unlink(buf)
	}

	eat_queue_mode = false
}

func write_remind_list() {
	fnam := sout("%s/remind", libdir)
	fp, err := fopen(fnam, "w")
	if err != nil {
		log.Printf("can't write %s: %v", fnam, err)
		return
	}
	for _, pl := range loop_player() {
		if subkind(pl) != sub_pl_regular {
			continue
		}
		p := rp_player(pl)
		if p == nil || p.SentOrders != FALSE || p.DontRemind != FALSE {
			continue
		} else if p.EMail == "" {
			log.Printf("player %s has no email address\n", box_code(pl))
			continue
		}
		fprintf(fp, "%s\n", p.EMail)
	}
	fclose(fp)
}

func read_spool(mail_now bool) bool {
	dirSpool := filepath.Join(libdir, "spool")
	files, err := os.ReadDir(dirSpool)
	if err != nil {
		log.Printf("read_spool: can't open %q: %v\n", dirSpool, err)
		return false
	}

	// check for stop file
	// stop is a special file that tells us to stop processing the directory
	for _, f := range files {
		fname := f.Name()
		if strings.HasPrefix(fname, "stop") {
			log.Printf("read_spool: stop file found: %q\n", fname)
			return false
		}
	}

	remindMe := false
	for _, f := range files {
		// ignore hidden and editor temp files
		fname := f.Name()
		if strings.HasPrefix(fname, ".") || strings.HasSuffix(f.Name(), "~") || strings.HasSuffix(f.Name(), ".swp") {
			log.Printf("read_spool: ignoring %q\n", fname)
			continue
		}
		if strings.HasPrefix(fname, "m") {
			log.Printf("read_spool: processing %q: mail_now %v\n", fname, mail_now)
			mailFile := filepath.Join(libdir, "spool", fname)
			eat(mailFile, mail_now)
			if mail_now { // remove the spooled file if we're actually replying
				// todo: brave of us to assume that we had no errors processing the file before deleting it
				unlink(mailFile)
			}
			remindMe = true
		}
	}

	if remindMe {
		write_remind_list()
	}

	return true
}

func eat_loop(mail_now bool) {
	if err := mkdir(filepath.Join(libdir, "orders")); err != nil {
		panic(err)
	}
	if err := mkdir(filepath.Join(libdir, "spool")); err != nil {
		panic(err)
	}
	//chmod(sout("%s/spool", libdir), 0777);

	write_remind_list()

	for read_spool(mail_now) {
		time.Sleep(10 * time.Second)
	}
}

func v_format(c *command) int {
	/*
	 *  Make this work in the BEGIN section as well.
	 *
	 */
	var plyr int
	if c.who != 0 {
		plyr = player(c.who)
	} else if pl != 0 {
		plyr = pl
	} else {
		err(EAT_ERR, "BEGIN must come before EMAIL")
		out(eat_pl, "      Notab not set.")
		return TRUE
	}

	/*
	 *  Format can have:
	 *
	 *  CLEAR
	 *  HTML
	 *  TEXT
	 *  RAW
	 *  TAGS
	 *  ALT
	 *
	 */
	for i := 1; i < len(c.parse) && len(c.parse[i]) != 0; i++ {
		out_alt_who = EAT_OKAY
		if strcasecmp_sb("HTML", c.parse[i]) == 0 {
			p_player(plyr).Format |= HTML
			wout(or_int(c.who != FALSE, c.who, eat_pl), "HTML format added.")
		} else if strcasecmp_sb("CLEAR", c.parse[i]) == 0 {
			p_player(plyr).Format = 0
			wout(or_int(c.who != FALSE, c.who, eat_pl), "All formats cleared (text only).")
		} else if strcasecmp_sb("TEXT", c.parse[i]) == 0 {
			p_player(plyr).Format |= TEXT
			wout(or_int(c.who != FALSE, c.who, eat_pl), "TEXT format added.")
		} else if strcasecmp_sb("TAGS", c.parse[i]) == 0 {
			p_player(plyr).Format |= TAGS
			wout(or_int(c.who != FALSE, c.who, eat_pl), "TAGS format added.")
		} else if strcasecmp_sb("ALT", c.parse[i]) == 0 {
			p_player(plyr).Format |= ALT
			wout(or_int(c.who != FALSE, c.who, eat_pl), "ALT format added.")
		} else if strcasecmp_sb("RAW", c.parse[i]) == 0 {
			p_player(plyr).Format |= RAW
			wout(or_int(c.who != FALSE, c.who, eat_pl), "RAW format added.")
		}
	}
	return TRUE
}

func do_rules_url(c *command) bool {
	//cmd := c.parse[1];

	/*
	 *  Make this work in the BEGIN section as well.
	 *
	 */
	var plyr int
	if c.who != 0 {
		plyr = player(c.who)
	} else if pl != 0 {
		plyr = pl
	} else {
		err(EAT_ERR, "BEGIN must come before EMAIL")
		out(eat_pl, "      Notab not set.")
		return true
	}

	if len(c.parse[2]) != 0 {
		if len(c.parse[2]) > 255 {
			out_alt_who = EAT_OKAY
			wout(or_int(c.who != FALSE, c.who, eat_pl), "Rules HTML path too long.  Maximum length 255 chars.")
		} else {
			if p_player(plyr).RulesPath != "" {
				my_free(p_player(plyr).RulesPath)
			}
			p_player(plyr).RulesPath = string(c.parse[2])
			out_alt_who = EAT_OKAY
			wout(or_int(c.who != FALSE, c.who, eat_pl), "Rules HTML path \"%s\" set.", p_player(plyr).RulesPath)
		}
	} else {
		if p_player(plyr).RulesPath != "" {
			my_free(p_player(plyr).RulesPath)
		}
		p_player(plyr).RulesPath = ""
		out_alt_who = EAT_OKAY
		wout(or_int(c.who != FALSE, c.who, eat_pl), "Rules HTML path cleared.")
	}
	return true
}

func do_db_url(c *command) bool {
	//cmd := c.parse[1];

	/*
	 *  Make this work in the BEGIN section as well.
	 *
	 */
	var plyr int
	if c.who != 0 {
		plyr = player(c.who)
	} else if pl != 0 {
		plyr = pl
	} else {
		err(EAT_ERR, "BEGIN must come before EMAIL")
		out(eat_pl, "      Notab not set.")
		return true
	}

	if len(c.parse[2]) != 0 {
		if len(c.parse[2]) > 255 {
			out_alt_who = EAT_OKAY
			wout(or_int(c.who != FALSE, c.who, eat_pl),
				"DB HTML path too long.  Maximum length 255 chars.")
		} else {
			if p_player(plyr).DBPath != "" {
				my_free(p_player(plyr).DBPath)
			}
			p_player(plyr).DBPath = string(c.parse[2])
			out_alt_who = EAT_OKAY
			wout(or_int(c.who != FALSE, c.who, eat_pl), "DB HTML path \"%s\" set.", p_player(plyr).DBPath)
		}
	} else {
		if p_player(plyr).DBPath != "" {
			my_free(p_player(plyr).DBPath)
		}
		p_player(plyr).DBPath = ""
		out_alt_who = EAT_OKAY
		wout(or_int(c.who != FALSE, c.who, eat_pl), "DB HTML path cleared.")
	}
	return true
}

func v_notab(c *command) int {
	/*
	 *  Make this work in the BEGIN section as well.
	 *
	 */
	var plyr int
	if c.who != 0 {
		plyr = player(c.who)
	} else if pl != 0 {
		plyr = pl
	} else {
		err(EAT_ERR, "BEGIN must come before EMAIL")
		out(eat_pl, "      Notab not set.")
		return TRUE
	}

	p_player(plyr).NoTab = c.a != FALSE

	out_alt_who = EAT_OKAY
	if c.a != 0 {
		if c.who != 0 {
			wout(c.who, "No TAB characters will appear in turn reports.")
		} else {
			wout(eat_pl, "No TAB characters will appear in turn reports.")
		}
	} else if c.who != 0 {
		wout(c.who, "TAB characters may appear in turn reports.")
	} else {
		wout(eat_pl, "TAB characters may appear in turn reports.")
	}
	return TRUE
}
