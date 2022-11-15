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

// advanced sorcery

// give <who> <what> [qty] [have-left]
func v_teleport_item(c *command) int {
	return TRUE
}

func d_teleport_item(c *command) int {
	target := c.a
	item := c.b
	qty := c.c
	have_left := c.d

	if kind(target) != T_char {
		wout(c.who, "%s is not a character.", box_code(target))
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Prisoners may not be given anything.")
		return FALSE
	}

	if kind(item) != T_item {
		wout(c.who, "%s is not an item.", box_code(item))
		return FALSE
	}

	if has_item(c.who, item) < 1 {
		wout(c.who, "%s does not have any %s.", box_name(c.who), box_code(item))
		return FALSE
	}

	if rp_item(item).ungiveable != FALSE {
		wout(c.who, "You cannot teleport %s to another noble.", plural_item_name(item, 2))
		return FALSE
	}

	if crosses_ocean(target, c.who) {
		wout(c.who, "Something seems to block your magic.")
		return FALSE
	}

	qty = how_many(c.who, c.who, item, qty, have_left)
	if qty <= 0 {
		return FALSE
	}

	aura := 3 + item_weight(item)*qty/250
	if !check_aura(c.who, aura) {
		return FALSE
	}
	if !will_accept(target, item, c.who, qty) {
		return FALSE
	}
	charge_aura(c.who, aura)
	if !move_item(c.who, target, item, qty) {
		panic("assert(move_item(c.who, target, item, qty))")
	}

	wout(c.who, "Teleported %s to %s.", just_name_qty(item, qty), box_name(target))
	wout(target, "%s teleported %s to us.", box_name(c.who), just_name_qty(item, qty))

	return TRUE
}

// create iron golem. decays after one year.
func v_create_iron_golem(c *command) int {
	wout(c.who, "Begin construction of a iron golem.")
	return TRUE
}

func d_create_iron_golem(c *command) int {
	if !charge_aura(c.who, skill_piety(c.use_skill)) {
		return FALSE
	}

	gen_item(c.who, item_iron_golem, 1)
	wout(c.who, "You have created a iron golem.")

	return TRUE
}

func v_trance(c *command) int {
	if has_skill(c.who, sk_trance) < 1 {
		wout(c.who, "Requires knowledge of %s.", box_name(sk_trance))
		return FALSE
	}
	return TRUE
}

func d_trance(c *command) int {
	p := p_magic(c.who)
	p.cur_aura = max_eff_aura(c.who)
	wout(c.who, "Current aura is now %d.", p.cur_aura)

	if char_health(c.who) < 100 || char_sick(c.who) != FALSE {
		p_char(c.who).sick = FALSE
		rp_char(c.who).health = 100

		wout(c.who, "%s is fully healed.", box_name(c.who))
	}

	return TRUE
}
