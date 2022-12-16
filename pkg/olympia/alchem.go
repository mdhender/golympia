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

func d_brew_death(c *command) int {
	potion := new_potion(c.who)
	if potion < 0 {
		wout(c.who, "Attempt to brew potion failed.")
		return FALSE
	}
	p_item_magic(potion).UseKey = use_death_potion
	return TRUE
}

func d_brew_fiery(c *command) int {
	potion := new_potion(c.who)
	if potion < 0 {
		wout(c.who, "Attempt to brew potion failed.")
		return FALSE
	}
	p_item_magic(potion).UseKey = use_fiery_potion
	return TRUE
}

func d_brew_heal(c *command) int {
	potion := new_potion(c.who)
	if potion < 0 {
		wout(c.who, "Attempt to brew potion failed.")
		return FALSE
	}
	p_item_magic(potion).UseKey = use_heal_potion
	return TRUE
}

func d_brew_slave(c *command) int {
	potion := new_potion(c.who)
	if potion < 0 {
		wout(c.who, "Attempt to brew potion failed.")
		return FALSE
	}
	p_item_magic(potion).UseKey = use_slave_potion
	return TRUE
}

func d_brew_weightlessness(c *command) int {
	potion := new_potion(c.who)
	if potion < 0 {
		wout(c.who, "Attempt to brew potion failed.")
		return FALSE
	}
	p_item_magic(potion).UseKey = use_weightlessness_potion
	return TRUE
}

//extern int gold_lead_to_gold;
func d_lead_to_gold(c *command) int {
	qty := c.d
	has := has_item(c.who, item_lead)

	if has_item(c.who, item_farrenstone) < 1 {
		wout(c.who, "Requires %s.", box_name_qty(item_farrenstone, 1))
		return FALSE
	}

	if has < qty {
		qty = has
	}
	if qty == 0 {
		wout(c.who, "Don't have any %s.", box_name(item_lead))
		return FALSE
	}

	wout(c.who, "Turned %s into %s.", just_name_qty(item_lead, qty), just_name_qty(item_gold, qty*10))

	consume_item(c.who, item_lead, qty)
	consume_item(c.who, item_farrenstone, 1)
	gen_item(c.who, item_gold, qty*10)

	gold_lead_to_gold += 100

	return TRUE
}

func new_potion(who int) int {
	potion := create_unique_item(who, 0)
	if potion < 0 {
		return -1
	}

	switch rnd(1, 2) {
	case 1:
		set_name(potion, "Magic potion")
	case 2:
		set_name(potion, "Strange potion")
	default:
		panic("!reached")
	}

	p := p_item_magic(potion)
	p.Creator = who
	p_item(potion).weight = 1

	wout(who, "Produced one %q", box_name(potion))

	return potion
}

func v_brew(c *command) int {
	return TRUE
}

func v_lead_to_gold(c *command) int {
	amount := c.a

	if has_item(c.who, item_farrenstone) < 1 {
		wout(c.who, "Requires %s.", box_name_qty(item_farrenstone, 1))
		return FALSE
	}

	qty := has_item(c.who, item_lead)
	if amount < 1 || amount > qty {
		amount = qty
	}
	qty = min(qty, 20)
	if qty == 0 {
		wout(c.who, "Don't have any %s.", box_name(item_lead))
		return FALSE
	}

	c.d = qty

	return TRUE
}

func v_use_death(c *command) int {
	item := c.a
	if kind(item) != T_item {
		panic("assert(kind(item) == T_item)")
	}

	wout(c.who, "%s drinks the potion...", just_name(c.who))
	destroy_unique_item(c.who, item)

	wout(c.who, "It's poison!")

	p_char(c.who).sick = TRUE
	add_char_damage(c.who, 50, MATES)

	return TRUE
}

func v_use_fiery(c *command) int {
	item := c.a
	if kind(item) != T_item {
		panic("assert(kind(item) == T_item)")
	}

	wout(c.who, "%s drinks the potion...", just_name(c.who))
	destroy_unique_item(c.who, item)

	wout(c.who, "It burns horribly!")

	add_char_damage(c.who, 10+rnd(1, 10), MATES)

	return TRUE
}

func v_use_heal(c *command) int {
	item := c.a
	if kind(item) != T_item {
		panic("assert(kind(item) == T_item)")
	}

	wout(c.who, "%s drinks the potion...", just_name(c.who))
	destroy_unique_item(c.who, item)

	if char_health(c.who) == 100 {
		if p_char(c.who).sick != FALSE {
			// todo: can a sick character drink a healing potion
			panic("p_char(c.who).sick == FALSE")
		}
		wout(c.who, "Nothing happens.")
		return TRUE
	}

	wout(c.who, "%s is immediately healed of all wounds!", just_name(c.who))

	p_char(c.who).sick = FALSE
	rp_char(c.who).health = 100

	return TRUE
}

func v_use_slave(c *command) int {
	item := c.a
	if kind(item) != T_item {
		panic("assert(kind(item) == T_item)")
	}

	creator := item_creator(item)

	// todo: must take into account different loyalties, percentage chance?
	//      5	0
	//      4  15 -- ?or death? --
	//      3  30
	//      2  60
	//      1  90

	// todo: should be log_code, not log_output?
	log_output(LOG_SPECIAL, "%s drinks a slavery potion to %s\n", box_name(c.who), box_name(creator))

	wout(c.who, "%s drinks the potion...", just_name(c.who))
	destroy_unique_item(c.who, item)

	if !valid_box(creator) || kind(creator) != T_char || get_effect(c.who, ef_guard_loyalty, 0, 0) != FALSE {
		wout(c.who, "Nothing happens.")
	} else if unit_deserts(c.who, creator, TRUE, LOY_contract, 250) {
		wout(c.who, "%s is suddenly overcome with an irresistible desire to serve %s.", just_name(c.who), box_name(creator))
	} else {
		wout(c.who, "Nothing happens.")
	}

	return TRUE
}

func v_use_weightlessness(c *command) int {
	item := c.a
	if kind(item) != T_item {
		panic("assert(kind(item) == T_item)")
	}

	wout(c.who, "%s drinks the potion...", just_name(c.who))
	destroy_unique_item(c.who, item)

	if get_effect(c.who, ef_weightlessness, 0, 0) != FALSE {
		wout(c.who, "%s is already weightless.", just_name(c.who))
		destroy_unique_item(c.who, item)
		return TRUE
	}

	if add_effect(c.who, ef_weightlessness, 0, 7, 1) == FALSE {
		wout(c.who, "Oddly enough, the potion has no effect.")
	} else {
		wout(c.who, "%s feels himself become weightless!", just_name(c.who))
	}

	return TRUE
}
