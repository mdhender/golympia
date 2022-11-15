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

func d_find_food(c *command) int {
	if has_item(province(c.who), item_peasant) >= 100 {
		wout(c.who, "You may only search for food in wilderness provinces.")
		return FALSE
	}
	// add the experience level of this guy into any effect...
	food_found := get_effect(c.who, ef_food_found, 0, 0)
	delete_effect(c.who, ef_food_found, 0)
	food_found += 10
	if subkind(province(c.who)) == sub_mountain {
		food_found += skill_exp(c.who, sk_find_food) / 2
	} else if subkind(province(c.who)) == sub_desert {
		food_found += skill_exp(c.who, sk_find_food) / 4
	} else {
		food_found += skill_exp(c.who, sk_find_food)
	}
	if add_effect(c.who, ef_food_found, 0, 35, food_found) == 0 {
		wout(c.who, "Through some odd circumstance you cannot find any food!")
		return FALSE
	}
	return TRUE
}

func v_find_food(c *command) int {
	if has_item(province(c.who), item_peasant) >= 100 {
		wout(c.who, "You may only search for food in wilderness provinces.")
		return FALSE
	}
	return TRUE
}
