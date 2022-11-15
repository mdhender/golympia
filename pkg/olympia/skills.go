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

type skills_l []int

func (l skills_l) Len() int {
	return len(l)
}

func (l skills_l) Less(i, j int) bool {
	return rp_skill(l[i]).use_count < rp_skill(l[j]).use_count
}

func (l skills_l) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l skills_l) delete(index int) skills_l {
	var cp skills_l
	for i, e := range l {
		if i == index {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}

// rem_value removes all elements that match the value
func (l skills_l) rem_value(value int) skills_l {
	cp := l
	for i := len(cp) - 1; i >= 0; i-- {
		if e := cp[i]; e == value {
			cp = cp.delete(i)
		}
	}
	return cp
}

// rem_value_uniq removes the rightmost element in the list that matches the value
func (l skills_l) rem_value_uniq(value int) skills_l {
	for i := len(l) - 1; i >= 0; i-- {
		if e := l[i]; e == value {
			return l.delete(i)
		}
	}
	return l
}

// reverse sorts the skills by bx.temp
func (l skills_l) sort_known_comp() {
	sort.Sort(bxtmp_l(l))
}

type skill_ent_l []*skill_ent

func (l skill_ent_l) copy() skill_ent_l {
	var cp skill_ent_l
	for _, e := range l {
		cp = append(cp, e)
	}
	return cp
}

func (l skill_ent_l) delete(index int) skill_ent_l {
	var cp skill_ent_l
	for i, e := range l {
		if i == index {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}

// rem_value removes all elements that match the value
func (l skill_ent_l) rem_value(value *skill_ent) skill_ent_l {
	var cp skill_ent_l
	for _, e := range l {
		if e != value {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}
