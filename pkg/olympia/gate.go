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

import "math"

func check_gate_here(who, gate int) bool {
	if kind(gate) != T_gate || subloc(gate) != subloc(who) {
		wout(who, "There is no gate %s here.", box_code(gate))
		return false
	}
	return true
}

func d_detect_gates(c *command) int {
	if charge_aura(c.who, 1) == FALSE {
		return FALSE
	}

	// Scott Turner: permit farcasting
	// ret := list_gates_here(c.who, subloc(c.who), true);
	ret := list_gates_here(c.who, cast_where(c.who), true)
	if !ret {
		ret = list_province_gates(c.who, province(cast_where(c.who))) != 0
		if !ret {
			list_nearby_gates(c.who, province(cast_where(c.who)))
		}
	}
	reset_cast_where(c.who)

	return TRUE
}

func d_notify_jump(c *command) int {
	gate := c.a
	where := cast_where(c.who)

	if charge_aura(c.who, 6) == FALSE {
		return FALSE
	}

	reset_cast_where(c.who)

	if kind(gate) != T_gate || subloc(gate) != where {
		wout(c.who, "There is no gate %s at %s.", box_code(gate), box_code(where))
		return FALSE
	}

	p_gate(gate).notify_jumps = c.who

	wout(c.who, "Notification spell successfully cast.")
	return TRUE
}

func d_notify_unseal(c *command) int {
	gate, key := c.a, c.b
	where := cast_where(c.who)

	if charge_aura(c.who, 5) == FALSE {
		return FALSE
	}

	reset_cast_where(c.who)

	if kind(gate) != T_gate || subloc(gate) != where {
		wout(c.who, "There is no gate %s at %s.", box_code(gate), box_code(where))
		return FALSE
	}

	sealed := gate_seal(gate)

	if sealed == 0 {
		wout(c.who, "The gate is not sealed.")
		return FALSE
	}

	if key != sealed {
		wout(c.who, "Incorrect gate key.  Spell fails.")
		return TRUE
	}

	if kind(gate) != T_gate {
		panic("assert(kind(gate) == T_gate)")
	}
	p_gate(gate).notify_unseal = c.who

	wout(c.who, "Notification spell successfully cast.")
	return TRUE
}

func d_rem_seal(c *command) int {
	gate := c.a
	where := cast_where(c.who)

	if charge_aura(c.who, 8) == FALSE {
		return FALSE
	}

	reset_cast_where(c.who)

	if kind(gate) != T_gate || subloc(gate) != where {
		wout(c.who, "There is no gate %s at %s.", box_code(gate), box_code(where))
		return FALSE
	}

	sealed := gate_seal(gate)
	if sealed == 0 {
		wout(c.who, "The gate is not sealed.")
		return FALSE
	}

	unseal_gate(c.who, gate)

	wout(c.who, "Unsealed %s.", box_name(gate))

	return TRUE
}

func d_reveal_key(c *command) int {
	gate := c.a

	where := cast_where(c.who)
	if charge_aura(c.who, 10) == FALSE {
		return FALSE
	}
	reset_cast_where(c.who)

	if kind(gate) != T_gate || subloc(gate) != where {
		wout(c.who, "There is no gate %s at %s.", box_code(gate), box_code(where))
		return FALSE
	}

	sealed := gate_seal(gate)
	if sealed == 0 {
		wout(c.who, "%s is not sealed.", box_name(gate))
		return FALSE
	}

	wout(c.who, "The key to %s is: %d", box_name(gate), sealed)
	return TRUE
}

func d_seal_gate(c *command) int {
	gate, key := c.a, c.b

	where := cast_where(c.who)
	if charge_aura(c.who, 6) == FALSE {
		return FALSE
	}
	reset_cast_where(c.who)

	if kind(gate) != T_gate || subloc(gate) != where {
		wout(c.who, "There is no gate %s at %s.", box_code(gate), box_code(where))
		return FALSE
	}

	sealed := gate_seal(gate)
	if sealed != 0 {
		wout(c.who, "%s has been sealed by someone else.", box_name(gate))
		return FALSE
	}

	p_gate(gate).seal_key = key
	wout(c.who, "Sealed %s with key %d.", box_name(gate), key)
	return TRUE
}

func d_unseal_gate(c *command) int {
	gate, key := c.a, c.b
	where := cast_where(c.who)

	if charge_aura(c.who, 3) == FALSE {
		return FALSE
	}

	reset_cast_where(c.who)

	if kind(gate) != T_gate || subloc(gate) != where {
		wout(c.who, "There is no gate %s at %s.", box_code(gate),
			box_code(where))
		return FALSE
	}

	sealed := gate_seal(gate)
	if sealed == 0 {
		wout(c.who, "The gate is not sealed.")
		return FALSE
	}

	if key != sealed {
		wout(c.who, "Incorrect gate key.  Unseal fails.")
		return TRUE
	}

	unseal_gate(c.who, gate)

	wout(c.who, "Successfully unsealed %s.", box_name(gate))

	return TRUE
}

func do_jump(who, dest, gate int, backwards bool) {
	leave_stack(who)
	wout(who, "Successful jump to %s.", box_name(dest))
	move_stack(who, dest)

	// todo: departure message? arrival message?

	clear_guard_flag(who)

	if gate != 0 {
		if kind(gate) != T_gate {
			panic("assert(kind(gate) == T_gate)")
		}
		p := rp_gate(gate)
		if p != nil && kind(p.notify_jumps) == T_char {
			if backwards {
				wout(p.notify_jumps, "%s has jumped backwards through %s.", box_name(who), box_name(gate))
			} else {
				wout(p.notify_jumps, "%s has jumped through %s.", box_name(who), box_name(gate))
			}
		}
	}

}

func list_gates_here(who, where int, show_dest bool) bool {
	first := true

	for _, gate := range loop_gates_here(where) {
		if first {
			out(who, "Gates here:")
			indent += 3
			first = false
		}

		var sealed, dest string
		if gate_seal(gate) != 0 {
			sealed = ", sealed"
		}
		if show_dest {
			dest = sout(", to %s", box_name(gate_dest(gate)))
		}

		out(who, "%s%s%s", box_name(gate), sealed, dest)
		set_known(who, gate)
	}

	if first {
		out(who, "There are no gates here.")
		return false
	}

	indent -= 3
	return true
}

func list_nearby_gates(who, where int) {
	dist := gate_dist(where)
	if dist == 0 {
		wout(who, "There are no nearby gates.")
	} else {
		wout(who, "The nearest gate is %s province%s away.", nice_num(dist), add_s(dist))
	}
}

func list_province_gates(who, where int) int {
	gate := province_gate_here(where)
	if gate != 0 {
		wout(who, "A gate exists somewhere in this province.")
	}
	return gate
}

func province_gate_here(where int) int {
	for _, i := range loop_all_here(where) {
		if kind(i) == T_gate {
			return i
		}
	}
	return 0
}

func unseal_gate(who, gate int) {
	if kind(gate) != T_gate {
		panic("assert(kind(gate) == T_gate)")
	}

	p := p_gate(gate)
	p.seal_key = 0
	if kind(p.notify_unseal) == T_char {
		wout(p.notify_unseal, "%s has been unsealed by %s.",
			box_name(gate), box_name(who))
	}
	p.notify_unseal = 0
}

func v_detect_gates(c *command) int {
	if check_aura(c.who, 1) == FALSE {
		return FALSE
	}
	return TRUE
}

func v_jump_gate(c *command) int {
	gate, key := c.a, c.b
	if !check_gate_here(c.who, gate) {
		return FALSE
	}
	set_known(c.who, gate)

	sealed, dest := gate_seal(gate), gate_dest(gate)
	if sealed > 0 && key == 0 {
		wout(c.who, "The gate is sealed.")
		return FALSE
	} else if sealed > 0 && key != sealed {
		wout(c.who, "Incorrect gate key.  Jump fails.")
		return TRUE
	}

	var w weights
	determine_stack_weights(c.who, &w, FALSE)
	cost := int(math.Ceil(float64(w.total_weight) / 250.0))
	if charge_aura(c.who, cost) == FALSE {
		return FALSE
	}

	if !valid_box(dest) {
		panic("assert(valid_box(dest))")
	} else if !is_loc_or_ship(dest) {
		panic("assert(is_loc_or_ship(dest))")
	}

	do_jump(c.who, dest, gate, false)
	wout(c.who, "Cost %s aura to jump %s weight.", nice_num(cost), nice_num(w.total_weight))

	return TRUE
}

func v_notify_jump(c *command) int {
	if check_aura(c.who, 6) == FALSE {
		return FALSE
	}
	return TRUE
}

func v_notify_unseal(c *command) int {
	if c.b == 0 {
		wout(c.who, "Must specify the gate seal.")
		return FALSE
	} else if check_aura(c.who, 5) == FALSE {
		return FALSE
	}

	return TRUE
}

func v_rem_seal(c *command) int {
	if check_aura(c.who, 8) == FALSE {
		return FALSE
	}
	return TRUE
}

func v_reveal_key(c *command) int {
	if check_aura(c.who, 10) == FALSE {
		return FALSE
	}
	wout(c.who, "Attempt to learn the key for %s.", box_name(c.a))
	return TRUE
}

func v_reverse_jump(c *command) int {
	gate, key := c.a, c.b
	if kind(gate) != T_gate || gate_dest(gate) != subloc(c.who) {
		wout(c.who, "No gate %s leads here.", box_code(gate))
		return FALSE
	}

	sealed, dest := gate_seal(gate), subloc(gate)
	if sealed > 0 && key == 0 {
		wout(c.who, "The gate is sealed.")
		return FALSE
	} else if sealed > 0 && key != sealed {
		wout(c.who, "Incorrect gate key.  Jump fails.")
		return TRUE
	}

	var w weights
	determine_stack_weights(c.who, &w, FALSE)
	cost := 2 * int(math.Ceil(float64(w.total_weight)/250.0))
	if charge_aura(c.who, cost) == FALSE {
		return FALSE
	}

	if !valid_box(dest) {
		panic("assert(valid_box(dest))")
	} else if !is_loc_or_ship(dest) {
		panic("assert(is_loc_or_ship(dest))")
	}

	do_jump(c.who, dest, gate, true)
	wout(c.who, "Cost %s aura to jump %s weight.", nice_num(cost), nice_num(w.total_weight))

	return TRUE
}

func v_seal_gate(c *command) int {
	if check_aura(c.who, 6) == FALSE {
		return FALSE
	}

	key := c.b
	if key < 1 {
		wout(c.who, "Must specify a key to seal the gate with.")
		return FALSE
	} else if key > 999 {
		wout(c.who, "The key must be between 1 and 999.")
		return FALSE
	}

	return TRUE
}

func v_teleport(c *command) int {
	dest := c.a

	if !is_loc_or_ship(dest) {
		wout(c.who, "There is no location %s.", c.parse[1])
		return FALSE
	}

	var w weights
	determine_stack_weights(c.who, &w, FALSE)
	cost := int(math.Ceil(float64(w.total_weight) / 250))

	if has_item(c.who, item_gate_crystal) < 1 {
		wout(c.who, "Teleportation requires %s.", box_name_qty(item_gate_crystal, 1))
		return FALSE
	}

	if charge_aura(c.who, cost) == FALSE {
		return FALSE
	}

	consume_item(c.who, item_gate_crystal, 1)

	do_jump(c.who, dest, 0, true)

	return TRUE
}

func v_unseal_gate(c *command) int {
	if check_aura(c.who, 3) == FALSE {
		return FALSE
	}
	key := c.b
	if key == 0 {
		wout(c.who, "Must specify a key to unseal the gate.")
		return FALSE
	}

	return TRUE
}
