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
	"path/filepath"
)

type fp_ent struct {
	fp           *os.File
	name         string
	player       int
	flush_always bool
	next         *fp_ent
}

const (
	// possible destinations of output:
	VECT   = (-1) // vector of recipients
	MASTER = (-2) // n >= 0: output to entity event log

	OUT_SUMMARY    = 0
	OUT_BANNER     = 0
	OUT_INCLUDE    = 1
	OUT_LORE       = 2
	OUT_NEW        = 3 // new player listing
	OUT_LOC        = 4 // location descriptions
	OUT_TEMPLATE   = 5 // order template
	OUT_GARR       = 6 // garrison log
	OUT_SHOW_POSTS = 7 // show what press and rumor look like

	// tags for log()
	LOG_CODE    = 10 // Code alerts
	LOG_SPECIAL = 11 // Special events
	LOG_DEATH   = 12 // Character deaths
	LOG_MISC    = 13 // Other junk
	LOG_DROP    = 14 // Player drops

	// tags for eat.c
	EAT_ERR     = 20 // Errors in orders submitted
	EAT_WARN    = 21 // Warnings in orders submitted
	EAT_QUEUE   = 22 // Current order queues
	EAT_HEADERS = 23 // Email headers bounced back
	EAT_OKAY    = 24 // Regular (non-error) output for scanner
	EAT_PLAYERS = 25 // Player list
)

var (
	out_alt_who      int // used if path == MASTER
	out_path         int // alternate sout directive
	out_vector       []int
	player_fp        map[int]*fp_ent
	second_indent    = 0
	show_to_garrison = false
	spaces           []byte // used for indenting
	spaces_len       int
	underlines       []byte
)

func alloc_fp(player int) *fp_ent {
	if player_fp == nil {
		player_fp = make(map[int]*fp_ent)
	}
	if _, ok := player_fp[player]; !ok {
		player_fp[player] = open_fp(&fp_ent{}, player)
	}
	return player_fp[player]
}

/*
 *  Output formatters
 */
func bottom_out(pl, who, unit int, s string) {
	var bits uint64

	var fp *os.File
	if immediate != FALSE {
		fp = os.Stdout
	} else {
		if pl == 0 {
			return
		}
		fp = grab_fp(pl)
		p_player(pl).output = set_bit((p_player(pl).output), who)
	}

	if show_day {
		bits |= 1
	}
	if second_indent != FALSE {
		bits |= 2
	}
	if unit != FALSE {
		bits |= 4
	}

	switch bits {
	case 0:
		fprintf(fp, "%d::%d::%s\n", who, indent, s)
	case 1:
		fprintf(fp, "%d::%d:%d:%s\n", who, indent, sysclock.day, s)
	case 2:
		fprintf(fp, "%d::%d/%d::%s\n", who, indent, second_indent, s)
	case 3:
		fprintf(fp, "%d::%d/%d:%d:%s\n", who, indent, second_indent, sysclock.day, s)
	case 4:
		fprintf(fp, "%d:%s:%d::%s\n", who, box_code_less(unit), indent, s)
	case 5:
		fprintf(fp, "%d:%s:%d:%d:%s\n", who, box_code_less(unit), indent, sysclock.day, s)
	case 6:
		fprintf(fp, "%d:%s:%d/%d::%s\n", who, box_code_less(unit), indent, second_indent, s)
	case 7:
		fprintf(fp, "%d:%s:%d/%d:%d:%s\n", who, box_code_less(unit), indent, second_indent, sysclock.day, s)
	default:
		panic("!reached")
	}
}

/*
 *  If output is sent to a location, show it to all characters in that
 *  location.  This includes characters one level deep in sublocations.
 *
 *  If we're not there, but we have a garrison there, and the event
 *  is one which a garrison would see, show it to them, unless it is
 *  in a hidden loc, which garrisons can't see into.
 */
func can_view_loc(pl int, p *EntityPlayer, where int, outer int) int {
	if test_bit(p.locs, where) {
		return TRUE
	} else if test_bit(p.locs, outer) {
		if loc_hidden(where) && test_known(pl, where) {
			return FALSE
		}
		return TRUE
	} else if show_to_garrison && loc_depth(outer) == LOC_province && player(province_admin(outer)) == pl {
		if where != outer && loc_hidden(where) {
			return FALSE
		}
		return TRUE
	}
	return FALSE
}

func close_logfile() {
	if immediate != FALSE {
		return
	}
	for _, fp := range player_fp {
		if fp.fp != nil {
			_ = fp.fp.Close()
		}
	}
	player_fp = make(map[int]*fp_ent)
}

func comma_append(s, t string) string {
	if len(s) == 0 {
		return t
	}
	return s + ", " + t
}

func grab_fp(player int) *os.File {
	if player_fp == nil {
		player_fp = make(map[int]*fp_ent)
	}
	if p, ok := player_fp[player]; ok {
		return p.fp
	}
	fp := alloc_fp(player)
	return fp.fp
}

func init_spaces() {
	spaces_len = 150
	for i := 0; i < spaces_len; i++ {
		spaces = append(spaces, ' ')
	}
}

func initialize_buffer() {
	// mdhender: do nothing
}

func lines(who int, s string) {
	if underlines == nil {
		for i := 0; i < 72; i++ { /* used to be wrap_pos */
			underlines = append(underlines, '-')
		}
	}

	wout(who, "%s", s)
	tagout(who, "<tag type=header>")
	out(who, "%s", string(underlines))
	tagout(who, "</tag type=header>")
}

func log_output(k int, format string, args ...interface{}) {
	save_out_path := out_path
	save_out_alt_who := out_alt_who

	if !(k >= 10 && k <= 20) {
		panic("assert(k >= 10 && k <= 20)")
	}

	out_path = MASTER
	out_alt_who = k

	out_sup(gm_player, fmt.Sprintf(format, args...))

	out_path = save_out_path
	out_alt_who = save_out_alt_who
}

func match_lines(who int, s string) {
	out(who, "%s", s)
	var buf []byte
	for i := 0; i < len(s); i++ {
		buf = append(buf, '-')
	}
	out(who, "%s", string(buf))
}

/*
 *  Fan-out decollating	fp allocater
 */
func open_fp(fp *fp_ent, player int) *fp_ent {
	var err error
	fp.name = filepath.Join(libdir, "log", fmt.Sprintf("%d", player))
	fp.fp, err = os.OpenFile(fp.name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("open_fp: can't open %q: %v\n", fp.name, err)
		panic(err)
	}
	fp.player = player
	fp.flush_always = flush_always
	return fp
}

func open_logfile() {
	if immediate != FALSE {
		return
	} else if err := rmdir(filepath.Join(libdir, "log")); err != nil {
		panic(err)
	} else if err = os.Mkdir(filepath.Join(libdir, "log"), 0755); err != nil {
		panic(err)
	}
}

func open_logfile_nondestruct() {
	if err := mkdir(filepath.Join(libdir, "log")); err != nil {
		panic(fmt.Sprintf("%q: %v\n", filepath.Join(libdir, "log"), err))
	}
}

func out(who int, format string, args ...interface{}) {
	out_sup(who, fmt.Sprintf(format, args...))
}

func out_location(where int, s string) {
	outer := viewloc(where)

	for _, pl := range loop_player() {
		p := p_player(pl)

		/*
		 *  If we can see the subloc, and we've touched it or its outer
		 *  viewloc, then we can see what goes on there.
		 */

		if can_view_loc(pl, p, where, outer) != FALSE {
			bottom_out(pl, outer, 0, s)
		}
	}
}

func out_garrison(garr int, s string) {
	pl := player(province_admin(garr))
	if pl != 0 && subkind(pl) != sub_pl_silent {
		bottom_out(pl, OUT_GARR, garr, s)
	}
}

func out_sup(who int, s string) {
	if who == VECT {
		for i := 0; i < len(out_vector); i++ {
			if out_vector[i] == VECT {
				panic("assert(out_vector[i] != VECT)")
			}
			out_sup(out_vector[i], s)
		}
		return
	}

	if is_prisoner(who) { /* prisoners don't report anything */
		return
	}

	if out_path == MASTER {
		bottom_out(who, out_alt_who, 0, s)
	} else if subkind(who) == sub_garrison {
		out_garrison(who, s)
	} else if is_loc_or_ship(who) {
		out_location(who, s)
	} else {
		bottom_out(player(who), who, 0, s)
	}
}

func restore_output_vector(t []int) {
	out_vector = nil
	for _, i := range t {
		out_vector = append(out_vector, i)
	}
}

func save_output_vector() []int {
	tmp := out_vector
	out_vector = nil
	return tmp
}

/*
 *  printf which returns a char	* which	does not need to
 *  be freed or	otherwise explicitly reclaimed.
 *
 *  Grabs the next buffer from the circular list, printfs
 *  into it, and returns it.  Buffer will eventually be	reused.
 */
func sout(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func tags_off() {
	options.output_tags--
}

func tags_on() {
	options.output_tags++
}

func tagout(who int, format string, args ...interface{}) {
	if options.output_tags < 1 {
		return
	}
	out_sup(who, fmt.Sprintf(format, args...))
}

func vector_add(who int) {
	out_vector = append(out_vector, who)
}

func vector_char_here(where int) {
	out_vector = nil
	for _, i := range loop_char_here(where) {
		out_vector = append(out_vector, i)
	}
}

func vector_clear() {
	out_vector = nil
}

func vector_players() {
	out_vector = nil
	for _, pl := range loop_player() {
		if pl != eat_pl && pl != skill_player {
			out_vector = append(out_vector, pl)
		}
	}
}

func vector_stack(who int, clear bool) {
	if clear {
		out_vector = nil
	}
	if clear || ilist_lookup(out_vector, who) < 0 {
		out_vector = append(out_vector, who)
	}
	for _, i := range loop_char_here(who) {
		if clear || ilist_lookup(out_vector, i) < 0 {
			out_vector = append(out_vector, i)
		}
	}
}

func wiout(who, ind int, format string, args ...interface{}) {
	second_indent = ind
	out_sup(who, fmt.Sprintf(format, args...))
	second_indent = 0
}

func wout(who int, format string, args ...interface{}) {
	out_sup(who, fmt.Sprintf(format, args...))
}

func wrap_done()       { panic("!initialized") }
func wrap_set(who int) { panic("!initialized") }
