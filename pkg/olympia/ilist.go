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
 * but WITHOUT ANY WARRANTY{panic("!implemented")} without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package olympia

// vectors/ilist.h

type ilist []int

type ints_l []int

// vectors/ilist.c

func ilist_add(l []int, n int) []int {
	if ilist_lookup(l, n) == -1 {
		return append(l, n)
	}
	return l
}

func ilist_copy(l ilist) ilist {
	panic("!implemented")
}

func ilist_delete(l []int, index int) []int {
	var cp []int
	for i, e := range l {
		if i == index {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}

func (l ints_l) lookup(n int) int {
	for i, v := range l {
		if n == v {
			return i
		}
	}
	return -1
}

func ilist_lookup(l ilist, n int) int {
	for i, v := range l {
		if n == v {
			return i
		}
	}
	return -1
}

func ilist_prepend(l *ilist, n int) *ilist {
	panic("!implemented")
}

func ilist_reclaim(l *ilist) {
	panic("!implemented")
}

func ilist_rem_value(l *ilist, n int) *ilist {
	panic("!implemented")
}

func ilist_rem_value_uniq(l *ilist, n int) *ilist {
	panic("!implemented")
}

func ilist_scramble(l ilist) ilist {
	panic("!implemented")
}

func rem_value(l []int, n int) []int {
	panic("!implemented")
}

func (l ints_l) ToList() (list ints_l) {
	for _, e := range l {
		if e != 0 {
			list = append(list, e)
		}
	}
	return list
}
