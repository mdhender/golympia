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

const (
	T_deleted  = 0 /* forget on save */
	T_player   = 1
	T_char     = 2
	T_loc      = 3
	T_item     = 4
	T_skill    = 5
	T_gate     = 6
	T_road     = 7
	T_deadchar = 8
	T_ship     = 9
	T_post     = 10
	T_storm    = 11
	T_unform   = 12 /* unformed noble */
	T_lore     = 13
	T_nation   = 14
	T_MAX      = 15 /* one past highest T_xxx define */
)

var strKind = map[string]int{
	"deleted":  T_deleted,
	"player":   T_player,
	"char":     T_char,
	"loc":      T_loc,
	"item":     T_item,
	"skill":    T_skill,
	"gate":     T_gate,
	"road":     T_road,
	"deadchar": T_deadchar,
	"ship":     T_ship,
	"post":     T_post,
	"storm":    T_storm,
	"unform":   T_unform,
	"lore":     T_lore,
	"nation":   T_nation,
}

var kindStr = map[int]string{
	T_deleted:  "deleted",
	T_player:   "player",
	T_char:     "char",
	T_loc:      "loc",
	T_item:     "item",
	T_skill:    "skill",
	T_gate:     "gate",
	T_road:     "road",
	T_deadchar: "deadchar",
	T_ship:     "ship",
	T_post:     "post",
	T_storm:    "storm",
	T_unform:   "unform",
	T_lore:     "lore",
	T_nation:   "nation",
}
