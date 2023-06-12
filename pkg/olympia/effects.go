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

type EffectList []*Effect

// Effect - only characters, locations, and sublocations have effects on them.
//
// effects hang off of nobles, locations or structures (*) and have the following properties:
//   - have a type (generally equal to the skill used to create it).
//   - have a duration in days.
//   - have data: single integer.
//
// This is used to implement things like a spell that gives a fortification a +25% resistance to attack, etc.
type Effect struct {
	Type    int `json:"type,omitempty"`    // type of effect, usually == to a sk_ number
	SubType int `json:"subType,omitempty"` // a subtype, surprise!
	Days    int `json:"days,omitempty"`    // remaining days of the effect.
	Data    int `json:"data,omitempty"`    // generic data for effect.
}

func (l effect_l) ToEffectList() (el EffectList) {
	for _, e := range l {
		el = append(el, &Effect{
			Type:    e.type_,
			SubType: e.subtype,
			Days:    e.days,
			Data:    e.data,
		})
	}
	return el
}
