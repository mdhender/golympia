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
	"log"
	"strings"
)

/*
 *  Entity coding system
 *
 *  range           extent  use
 *       1-  9,999   9,999  reserved (items, skills)
 *  10,000- 19,999  10,000  provinces        (CCNN: AA00-DV99)
 *  20,000- 26,759   6,760  player entities  (CCN)
 *  26,760- 33,519   6,760  lucky characters (CNC)
 *  33,520- 36,119   2,600  lucky locs       (CNN)
 *  36,120-102,400  66,279  sublocs, runoff  (CCNC)
 */

var letters = []byte("abcdefghijklmnopqrstuvwxyz")

func int_to_code(l int) string {
	var n, a, b, c int

	if l < 10000 {
		return sout("%d", l)
	}

	if l < 20000 { /* CCNN */
		l -= 10000

		n = l % 100
		l /= 100

		a = l % 26
		b = l / 26

		return sout("%c%c%02d", letters[b], letters[a], n)
	}

	if l < 26760 { /* CCN */
		l -= 20000

		n = l % 10
		l /= 10

		a = l % 26
		b = l / 26

		return sout("%c%c%d", letters[b], letters[a], n)
	}

	if l < 33520 { /* CNC */
		l -= 26760

		n = l % 26
		l /= 26

		a = l % 10
		b = l / 10

		return sout("%c%d%c", letters[b], a, letters[n])
	}

	if l < 36120 { /* CNN */
		l -= 33520

		n = l % 10
		l /= 10

		a = l % 10
		b = l / 10

		return sout("%c%d%d", letters[b], a, n)
	}

	/* CCNC */
	l -= 36120

	a = l % 26
	l /= 26

	b = l % 10
	l /= 10

	c = l % 26
	l /= 26

	return sout("%c%c%d%c", letters[l], letters[c], b, letters[a])

}

func code_to_int(s []byte) int {
	if len(s) == 0 {
		return 0
	}

	if isdigit(s[0]) {
		return atoi(string(s))
	}

	if !isalpha(s[0]) {
		return 0
	}

	var a, b, c, d int
	switch len(s) {
	case 3:
		if isdigit(s[1]) && isalpha(s[2]) { /* CNC */
			a = int(tolower(s[0]) - 'a')
			b = int(s[1] - '0')
			c = int(tolower(s[2]) - 'a')

			return a*260 + b*26 + c + 26760
		}

		if isalpha(s[1]) && isdigit(s[2]) { /* CCN */
			a = int(tolower(s[0]) - 'a')
			b = int(tolower(s[1]) - 'a')
			c = int(s[2] - '0')

			return a*260 + b*10 + c + 20000
		}

		if isdigit(s[1]) && isdigit(s[2]) { /* CNN */
			a = int(tolower(s[0]) - 'a')
			b = int(s[1] - '0')
			c = int(s[2] - '0')

			return a*100 + b*10 + c + 33520
		}

		return 0

	case 4:
		if isalpha(s[1]) && isdigit(s[2]) && isdigit(s[3]) {
			a = int(tolower(s[0]) - 'a')
			b = int(tolower(s[1]) - 'a')
			c = int(s[2] - '0')
			d = int(s[3] - '0')

			return a*2600 + b*100 + c*10 + d + 10000
		}

		if isalpha(s[1]) && isdigit(s[2]) && isalpha(s[3]) {
			a = int(tolower(s[0]) - 'a')
			b = int(tolower(s[1]) - 'a')
			c = int(s[2] - '0')
			d = int(tolower(s[3]) - 'a')

			return a*6760 + b*260 + c*26 + d + 36120
		}
		return 0

	default:
		return 0
	}
}

func scode(s []byte) int {
	if len(s) == 0 {
		return 0
	}
	if s[0] == '[' || s[0] == '(' {
		s = s[1:]
	}
	return code_to_int(s)
}

func name(n int) string {
	assert(valid_box(n))
	return bx[n].name
}

func set_name(n int, s string) {
	assert(valid_box(n))
	bx[n].name = strings.ReplaceAll(strings.ReplaceAll(s, "[", "{"), "]", "}")
}

func set_banner(n int, s string) {
	p := p_misc(n)
	p.display = ""
	if len(s) > 50 {
		s = s[:50]
	}
	p.display = s
}

func display_name(n int) string {
	if !valid_box(n) {
		return ""
	}

	s := name(n)
	if len(s) != 0 {
		return s
	}

	switch kind(n) {
	case T_player:
		return "Player"
	case T_gate:
		return "Gate"
	case T_post:
		return "Sign"
	}

	if i := noble_item(n); i != FALSE {
		return cap_(plural_item_name(i, 1))
	}

	return cap_(subkind_s[subkind(n)])
}

func display_kind(n int) string {
	switch subkind(n) {
	case sub_city:
		if is_port_city_where(n) != FALSE {
			return "port city"
		}
		return "city"
	case sub_fog, sub_mist, sub_rain, sub_wind:
		return "storm"
	case sub_guild:
		if rp_subloc(n) != nil && rp_subloc(n).guild != 0 && valid_box(rp_subloc(n).guild) {
			return sout("%s Guild", bx[rp_subloc(n).guild].name)
		}
		return ""
	}
	return subkind_s[subkind(n)]
}

/*
 *  Same as box code, less the brackets
 */

func box_code_less(n int) string {

	return int_to_code(n)
}

func box_code(n int) string {
	if n == garrison_magic {
		return "Garrison"
	}

	return sout("[%s]", int_to_code(n))
}

func box_name(n int) string {
	if n == garrison_magic {
		return "Garrison"
	}

	if valid_box(n) {
		if s := display_name(n); len(s) != 0 {
			if options.output_tags > 0 {
				return sout("<tag type=box id=%d>%s~%s</tag type=box id=%d>", n, s, box_code(n), n)
			}
			return sout("%s~%s", s, box_code(n))
		}
	}

	return box_code(n)
}

func just_name(n int) string {
	if n == garrison_magic {
		return "Garrison"
	}

	if valid_box(n) {
		if s := display_name(n); len(s) != 0 {
			return s
		}
	}

	return box_code(n)
}

func plural_item_name(item, qty int) string {
	if qty == 1 {
		return display_name(item)
	}

	var s string
	if rp_item(item) != nil {
		s = rp_item(item).plural_name
	}
	if len(s) == 0 {
		log.Printf("warning: plural name not set for item %s\n", box_code(item))
		s = display_name(item)
	}

	return s
}

func plural_item_box(item, qty int) string {
	if qty == 1 {
		return box_name(item)
	}
	s := plural_item_name(item, qty)
	if options.output_tags > 0 {
		return sout("<tag type=box id=%d link=%d>%s~%s</tag type=box id=%d>", item, item, s, box_code(item), item)
	}
	return sout("%s~%s", s, box_code(item))
}

func just_name_qty(item, qty int) string {

	return sout("%s~%s", nice_num(qty), plural_item_name(item, qty))
}

func box_name_qty(item, qty int) string {

	return sout("%s~%s", nice_num(qty), plural_item_box(item, qty))
}

func box_name_kind(n int) string {

	return sout("%s, %s", box_name(n), display_kind(n))
}

/*
 *  Routines for allocating entities and threading like entities
 *  together (kind_first, kind_next)
 */
var (
	next_chain = struct {
		cache_last int
		cache_kind int
	}{0, 0}
	sub_chain = struct {
		cache_last int
		cache_kind int
	}{0, -1}
)

func add_next_chain(n int) {
	assert(bx[n] != nil)
	kind := bx[n].kind
	if kind == 0 {
		return
	}

	/*  optim! */

	if next_chain.cache_kind == int(kind) && n > next_chain.cache_last && bx[next_chain.cache_last].x_next_kind == 0 {
		bx[next_chain.cache_last].x_next_kind = n
		bx[n].x_next_kind = 0
		next_chain.cache_last = n
		return
	}

	next_chain.cache_last = n
	next_chain.cache_kind = int(kind)

	if box_head[kind] == 0 {
		box_head[kind] = n
		bx[n].x_next_kind = 0
		return
	}

	if n < box_head[kind] {
		bx[n].x_next_kind = box_head[kind]
		box_head[kind] = n
		return
	}

	i := box_head[kind]
	for bx[i].x_next_kind > 0 && bx[i].x_next_kind < n {
		i = bx[i].x_next_kind
	}

	bx[n].x_next_kind = bx[i].x_next_kind
	bx[i].x_next_kind = n
}

func remove_next_chain(n int) {
	assert(bx[n] != nil)

	i := box_head[bx[n].kind]
	if i == n {
		box_head[bx[n].kind] = bx[n].x_next_kind
	} else {
		for i > 0 && bx[i].x_next_kind != n {
			i = bx[i].x_next_kind
		}
		assert(i > 0)
		bx[i].x_next_kind = bx[n].x_next_kind
	}
	bx[n].x_next_kind = 0
}

func add_sub_chain(n int) {
	assert(bx[n] != nil)
	kind := bx[n].skind

	/*  optim! */

	if sub_chain.cache_kind == int(kind) && n > sub_chain.cache_last && bx[sub_chain.cache_last].x_next_sub == 0 {
		bx[sub_chain.cache_last].x_next_sub = n
		bx[n].x_next_sub = 0
		sub_chain.cache_last = n
		return
	}

	sub_chain.cache_last = n
	sub_chain.cache_kind = int(kind)

	if sub_head[kind] == 0 {
		sub_head[kind] = n
		bx[n].x_next_sub = 0
		return
	}

	if n < sub_head[kind] {
		bx[n].x_next_sub = sub_head[kind]
		sub_head[kind] = n
		return
	}

	i := sub_head[kind]
	for bx[i].x_next_sub > 0 && bx[i].x_next_sub < n {
		i = bx[i].x_next_sub
	}

	bx[n].x_next_sub = bx[i].x_next_sub
	bx[i].x_next_sub = n
}

func remove_sub_chain(n int) {
	assert(bx[n] != nil)

	i := sub_head[bx[n].skind]
	if i == n {
		sub_head[bx[n].skind] = bx[n].x_next_sub
	} else {
		for i > 0 && bx[i].x_next_sub != n {
			i = bx[i].x_next_sub
		}

		assert(i > 0)

		bx[i].x_next_sub = bx[n].x_next_sub
	}

	bx[n].x_next_sub = 0
}

func delete_box(n int) {
	remove_next_chain(n)
	remove_sub_chain(n)
	if bx[n].kind == T_char {
		bx[n].kind = T_deadchar
	} else {
		bx[n].kind = T_deleted
	}
}

func change_box_kind(n int, kind int) {
	remove_next_chain(n)
	bx[n].kind = schar(kind)
	add_next_chain(n)
}

func change_box_subkind(n int, sk int) {
	if subkind(n) == schar(sk) {
		return
	}
	remove_sub_chain(n)
	bx[n].skind = schar(sk)
	add_sub_chain(n)
}

func alloc_box(n int, kind int, sk int) {
	assert(n > 0 && n < MAX_BOXES)

	if bx[n] != nil {
		log.Printf("alloc_box: DUP box %d\n", n)
		assert(false)
	}

	bx[n] = &box{}
	bx[n].kind = schar(kind)
	bx[n].skind = schar(sk)
	add_next_chain(n)
	add_sub_chain(n)
}

// rnd_alloc_num allocates a number in the range low...high.
// it panics if it can't find an available number in that range.
func rnd_alloc_num(low, high int) int {
	n := rnd(low, high)
	for i := n; i <= high; i++ {
		if bx[i] == nil {
			return i
		}
	}
	for i := low; i < n; i++ {
		if bx[i] == nil {
			return i
		}
	}
	log.Printf("rnd_alloc_num(%d,%d) failed\n", low, high)
	return -1
}

/*
 *  Entity coding system
 *
 *  range           extent  use
 *       1-  9,999   9,999  reserved (items, skills)
 *  10,000- 19,999  10,000  provinces        (CCNN: AA00-DV99)
 *  20,000- 26,759   6,760  player entities  (CCN)
 *  26,760- 33,519   6,760  lucky characters (CNC)
 *  33,520- 36,119   2,600  lucky locs       (CNN)
 *  36,120-102,400  66,279  sublocs, runoff  (CCNC)
 */

var new_ent_prime = false /* allocate from lower end of the range (prime real estate) */

func new_ent(kind, sk int) int {
	n := -1
	switch kind {
	case T_player:
		n = rnd_alloc_num(20000, 26759)
		if n < 0 {
			n = rnd_alloc_num(36120, MAX_BOXES-1)
		}
	case T_char, T_unform:
		if new_ent_prime {
			n = rnd_alloc_num(26760, 33519)
		}
		if n < 0 {
			n = rnd_alloc_num(36120, MAX_BOXES-1)
		}
	case T_skill:
		//n = rnd_alloc_num(8000, 8999);
		assert(false)
	case T_loc:
		if new_ent_prime {
			n = rnd_alloc_num(33520, 36119)
		}
		if n < 0 {
			n = rnd_alloc_num(36120, MAX_BOXES-1)
		}
	default:
		n = rnd_alloc_num(36120, MAX_BOXES-1)
	}
	alloc_box(n, kind, sk)
	return n
}
