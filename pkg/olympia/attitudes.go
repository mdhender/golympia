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

type att_ent struct {
	neutral ints_l
	hostile ints_l
	defend  ints_l
}

type Attitudes struct {
	Neutral []int `json:"neutral,omitempty"`
	Defend  []int `json:"defend,omitempty"`
	Hostile []int `json:"hostile,omitempty"`
}

func (a *att_ent) ToAttitudes() *Attitudes {
	if a == nil || (len(a.neutral) == 0 && len(a.defend) == 0 || len(a.hostile) == 0) {
		return nil
	}
	return &Attitudes{
		Neutral: a.neutral.ToBoxList(),
		Defend:  a.defend.ToBoxList(),
		Hostile: a.hostile.ToBoxList(),
	}
}
