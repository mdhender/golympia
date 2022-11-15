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

type queue_l []int

func (q queue_l) Len() int {
	return len(q)
}

func (q queue_l) Less(i, j int) bool {
	return bx[q[i]].temp < bx[q[j]].temp
}

func (q queue_l) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q queue_l) copy() queue_l {
	var cp queue_l
	return append(cp, q...)
}

func (q queue_l) delete(index int) queue_l {
	var cp queue_l
	for i, e := range q {
		if i == index {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}

// rem_value_uniq removes the rightmost element in the list that matches the value
func (q queue_l) rem_value_uniq(value int) queue_l {
	for i := len(q) - 1; i >= 0; i-- {
		if e := q[i]; e == value {
			return q.delete(i)
		}
	}
	return q
}
