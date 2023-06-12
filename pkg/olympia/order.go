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
	"github.com/mdhender/golympia/pkg/io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// manage list of unit orders for each faction

var (
	_static_auto_comment_c command
)

func autocomment(who int, line []byte) []byte {
	c := &_static_auto_comment_c
	/*
	 *  Set c.who after the parse, so that "garrison" doesn't
	 *  get de-referenced to the local garrison.
	 *
	 */
	c.who = 0
	if !oly_parse(c, []byte(line)) {
		return nil
	}
	c.who = who

	if c.cmd == FALSE {
		return nil
	}

	if cmd_tbl[c.cmd].cmd_comment != nil {
		ret := (cmd_tbl[c.cmd].cmd_comment)(c)
		if ret == "" {
			return nil
		}
		/*
		 *  We need to go through the comment and
		 *  make it all unbreakable spaces.
		 *
		 */
		return append([]byte{'#', '~'}, bytes.ReplaceAll([]byte(ret), []byte{' '}, []byte{'~'})...)
	}

	return nil
}

func flush_unit_orders(pl, who int) {
	for top_order(pl, who) != nil {
		pop_order(pl, who)
	}

	if player(who) == pl {
		c := rp_command(who)
		if c != nil && c.state == LOAD {
			command_done(c)
		}
	}
}

func is_stop_order(s []byte) bool {
	if s = bytes.TrimSpace(s); len(s) == 0 {
		panic("assert(len(s) != 0)")
	}

	fields := strings.Fields(string(s))
	if strings.ToLower(fields[0]) == "stop" {
		return true
	} else if fuzzy_strcmp([]byte(fields[0]), []byte("stop")) {
		return true
	}

	return false
}

func list_order_templates() {
	out_path = MASTER
	out_alt_who = OUT_TEMPLATE
	for _, pl := range loop_player() {
		if subkind(pl) == sub_pl_system || subkind(pl) == sub_pl_silent {
			continue
		}
		orders_template(pl, pl)
	}
	out_path = 0
	out_alt_who = 0
}

func list_pending_orders_sup(who, num, show_empty int) { panic("!implemented") }

func load_orders() error {
	dirOrders := filepath.Join(libdir, "orders")
	files, err := os.ReadDir(dirOrders)
	if err != nil {
		log.Printf("load_orders: can't open %q: %v\n", dirOrders, err)
		return err
	}
	for _, f := range files {
		if isdigit(f.Name()[0]) && !strings.HasSuffix(f.Name(), "~") {
			fact := atoi(f.Name())
			if !valid_box(fact) {
				log.Printf("ERROR: orders/%d but no box [%d]\n", fact, fact)
				continue
			}
			err := load_player_orders(fact)
			if err != nil {
				log.Printf("load_orders: %v\n", err)
			}
		}
	}
	return nil
}

func load_player_orders(pl int) error {
	if !valid_box(pl) {
		panic("assert(valid_box(pl))")
	} else if rp_player(pl).Orders != nil {
		panic("assert(rp_player(pl).orders == nil)")
	} else if !io.ReadFile(filepath.Join(libdir, "orders", fmt.Sprintf("%d", pl))) {
		return nil
	}

	for line, ok := io.ReadLine(); ok; line, ok = io.ReadLine() {
		var unit, p string
		if i := strings.Index(line, ":"); i != -1 {
			unit, p = line[:i], line[i+1:]
		} else {
			unit = line
		}
		queue_order(pl, atoi(unit), p)
	}

	return nil
}

func orders_other(who, pl int) {
	if kind(pl) != T_player {
		panic("assert(kind(pl) == T_player)")
	}

	p := rp_player(pl)
	if p == nil {
		return
	}

	first := true
	for i := 0; i < len(p.Orders); i++ {
		if pl == p.Orders[i].unit || !valid_box(p.Orders[i].unit) || kind(p.Orders[i].unit) == T_deadchar {
			continue
		} else if p.Units.lookup(p.Orders[i].unit) >= 0 {
			continue
		}

		/*
		 *  Don't output an empty template for a unit we swore away this turn
		 */
		c := rp_command(p.Orders[i].unit)
		l := rp_order_head(pl, p.Orders[i].unit)
		if (c == nil || c.state == DONE) && (l == nil || len(l.l) == 0) {
			continue
		}

		if first {
			out(who, "#")
			out(who, "# Orders for units you do not control as of now")
			out(who, "#")
			first = false
		}
		orders_template_sup(who, p.Orders[i].unit, pl)
	}
}

func orders_template(who, pl int) {
	var pass string
	p := rp_player(pl)
	if p != nil {
		if p.Password != "" {
			pass = sout(" \"*******\"")
		}
		/*	pass = sout(" \"%s\"", p.password); */
	}

	tagout(who, "<tag type=order_template id=%d>", who)
	tags_off()
	out(who, "")
	out(who, "# Note: Fill in your password below!")
	if player_broken_mailer(pl) != FALSE {
		out(who, " begin %s%s  # %s", box_code_less(pl), pass, box_name(pl))
	} else {
		out(who, "begin %s%s  # %s", box_code_less(pl), pass, box_name(pl))
	}

	if player(who) == pl && kind(who) == T_char && loyal_kind(who) == LOY_contract && loyal_rate(who) < 51 {
		out(who, "")
		out(who, "# WARNING -- loyalty is %s", loyal_s(who))
	}

	out(who, "")

	orders_template_sup(who, pl, pl)

	for _, i := range loop_units(pl) {
		orders_template_sup(who, i, pl)
	}

	orders_other(who, pl)

	out(who, "end")

	tags_on()
	tagout(who, "</tag type=order_template>")
}

func orders_template_sup(who, num, pl int) {
	var nam, time_left string
	if pl == player(num) {
		if is_prisoner(num) || kind(num) != T_char {
			nam = sout("  # %s", box_name(num))
		} else {
			nam = sout("  # %s in %s", box_name(num), box_name(subloc(num)))
		}
	}

	out(who, "unit %s%s", box_code_less(num), nam)
	indent += 3

	if pl == player(num) {
		c := rp_command(num)

		if loyal_kind(num) == LOY_contract && loyal_rate(num) <= 50 {
			out(who, "#")
			out(who, "# %s has loyalty %s and will renounce",
				box_code_less(num), loyal_s(num))
			out(who, "# loyalty at the end of this turn.")
			out(who, "#")
		}

		/*
		 *  Tue Dec 26 12:32:29 2000 -- Scott Turner
		 *
		 *  Provide some helpful information?
		 *
		 *  Tue Dec 26 13:03:30 2000 -- Scott Turner
		 *
		 *  Hmm, stuff gets wrapped by rep.  And other problems.
		 *
		 */
		//out(who, "#");
		//out(who, "# Location: %s", char_rep_location(num));
		//if (stack_parent(num)) {out(who, "# Stacked under: %s", box_name(stack_parent(num)));}
		//char_rep_health(who,num,"# ");
		//if (is_priest(num)) {out(who,"# Current piety:  %d", rp_char(num).religion.piety);}
		//if (is_magician(num)) {char_rep_magic(who, num, "# ");}
		//list_skills(who, num, "# ");
		//list_partial_skills(who, num, "# ");
		//show_char_inventory(who, num, "# ");

		if c != nil && (c.state == RUN || c.state == LOAD) {
			if c.state == RUN {
				if c.wait < 0 {
					time_left = " (still~executing)"
				} else {
					time_left = sout(" (executing for %s more day%s)", nice_num(c.wait), add_s(c.wait))
				}
				out(who, "# > %s%s", c.line, time_left)
			} else { /* command has loaded, but not started yet */
				if valid_box(num) && player(num) != 0 && player(num) == pl {
					out(who, "%-20s%s", c.line, autocomment(num, []byte(c.line)))
				} else {
					out(who, "%-20s", c.line)
				}
			}
		}
	}

	l := rp_order_head(pl, num)
	if l != nil && len(l.l) > 0 {
		for i := 0; i < len(l.l); i++ {
			if valid_box(num) && player(num) != 0 && player(num) == pl {
				out(who, "%-20s%s", eat_leading_trailing_whitespace([]byte(l.l[i])), autocomment(num, l.l[i]))
			} else {
				out(who, "%-20s", eat_leading_trailing_whitespace([]byte(l.l[i])))
			}
		}
	}

	indent -= 3
	out(who, "")
}

func p_order_head(pl, who int) *orders_list {
	p := p_player(pl)
	for i := 0; i < len(p.Orders); i++ {
		if p.Orders[i].unit == who {
			return p.Orders[i]
		}
	}

	ol := &orders_list{unit: who}
	p.Orders = append(p.Orders, ol)

	return ol
}

func pop_order(player, who int) {
	p := rp_order_head(player, who)
	if p == nil {
		panic("assert(p != nil)")
	} else if len(p.l) == 0 {
		panic("assert(len(p.l) != 0)")
	}
	p.l = p.l[1:]
}

func prepend_order(pl, who int, s string) {
	p := p_order_head(pl, who)
	p.l = append(orders_l{[]byte(s)}, p.l...)
}

// loose, convenient interface for queue_order()
func queue(who int, format string, args ...interface{}) {
	queue_order(player(who), who, fmt.Sprintf(format, args...))
}

func queue_order(player, who int, s string) {
	p := p_order_head(player, who)
	p.l = append(p.l, []byte(s))
}

func queue_stop(pl, who int) {
	if stop_order(pl, who) {
		return
	}
	prepend_order(pl, who, "stop")
}

func rp_order_head(pl, who int) *orders_list {
	p := rp_player(pl)
	if p == nil {
		return nil
	}
	for i := 0; i < len(p.Orders); i++ {
		if p.Orders[i].unit == who {
			return p.Orders[i]
		}
	}
	return nil
}

func save_orders() error {
	if err := rmdir(filepath.Join(libdir, "orders")); err != nil {
		return fmt.Errorf("save_orders: %w", err)
	} else if err = os.Mkdir(filepath.Join(libdir, "orders"), 0755); err != nil {
		return fmt.Errorf("save_orders: %w", err)
	}
	for _, i := range loop_player() {
		if err := save_player_orders(i); err != nil {
			return fmt.Errorf("save_orders: player %d: %w", i, err)
		}
	}
	return nil
}

func save_player_orders(pl int) error {
	if !valid_box(pl) {
		panic("assert(valid_box(pl))")
	}

	p := rp_player(pl)
	if p == nil {
		return nil
	}

	var fp *os.File
	var err error
	for i := 0; i < len(p.Orders); i++ {
		if !valid_box(p.Orders[i].unit) || kind(p.Orders[i].unit) == T_deadchar {
			continue
		}
		for j := 0; j < len(p.Orders[i].l); j++ {
			if fp == nil {
				fname := filepath.Join(libdir, "orders", fmt.Sprintf("%d", pl))
				fp, err = fopen(fname, "w")
				if err != nil {
					return fmt.Errorf("save_player_orders: %q: %w\n", fname, err)
				}
			}
			fprintf(fp, "%d:%s\n", p.Orders[i].unit, p.Orders[i].l[j])
		}
	}
	fp = fclose(fp)
	return nil
}

/*
 *  Return TRUE if a stop order is queue for the given unit.
 *  STOP orders must be the first command in the order queue.
 */
func stop_order(pl, who int) bool {
	s := top_order(pl, who)
	if len(s) == 0 {
		return false
	}
	return is_stop_order(s)
}

func top_order(player, who int) []byte {
	p := rp_order_head(player, who)
	if p != nil && len(p.l) > 0 {
		return p.l[0]
	}
	return nil
}
