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

// as always, thanks https://ueokande.github.io/go-slice-tricks/

// effect structure.
// only characters, locations, and sublocations have effects on them.
//
// effects hang off of nobles, locations or structures (*) and have the following properties:
//   - have a type (generally equal to the skill used to create it).
//   - have a duration in days.
//   - have data: single integer.
//
// This is used to implement things like a spell that gives a fortification a +25% resistance to attack, etc.
type effect struct {
	type_   int // type of effect, usually == to a sk_ number
	subtype int // a subtype, surprise!
	days    int // remaining days of the effect.
	data    int // generic data for effect.
}

// effects returns a boxed effect?
func effects(n int) []*effect {
	return bx[n].effects
}

// add an effect to a thing
func add_effect(what, type_, sub_type, duration, value int) int {
	// validity checks
	if !valid_box(what) {
		return FALSE
	}

	// allocate and fill in the new effect.
	e := &effect{}
	e.type_ = type_
	e.subtype = sub_type
	e.days = duration
	e.data = value

	// now append it to the effects list
	bx[what].effects = append(bx[what].effects, e)

	return TRUE
}

// delete the first effect of the given type.
// which seems to actually mean delete the righmost (newest) effect.
func delete_effect(what, type_, sub_type int) int {
	// validity checks.
	if !valid_box(what) {
		return FALSE
	}

	// only characters, locations, and sublocations have effects on them.
	el := effects(what)
	for i := len(el) - 1; i >= 0; i-- {
		if el[i].type_ == type_ && (el[i].subtype == sub_type || sub_type == 0) {
			// remove the effect from the list
			bx[what].effects = append(bx[what].effects[:i], bx[what].effects[i+1:]...)
			return TRUE
		}
	}

	return FALSE
}

func delete_all_effects(what int, type_ int, sub_type int) int {
	// validity checks.
	if !valid_box(what) {
		return FALSE
	}

	// only characters, locations, and sublocations have effects on them.
	found := false
	for _, e := range effects(what) {
		if e.type_ == type_ && (e.subtype == sub_type || sub_type == 0) {
			found = true
			break
		}
	}
	if !found {
		return FALSE
	}

	var el []*effect
	for _, e := range effects(what) {
		if e.type_ == type_ && (e.subtype == sub_type || sub_type == 0) {
			// remove the effect from the list
		}
		el = append(el, e)
	}
	bx[what].effects = el

	return TRUE
}

// like get_effect, only sums multiple cumulative effects.
func get_all_effects(what, type_, sub_type, v int) int {
	// validity checks.
	if !valid_box(what) {
		return 0
	}
	// only characters, locations, and sublocations have effects on them.

	// look for the effect
	answer := 0
	for _, e := range effects(what) {
		if e.type_ == type_ && (e.subtype == sub_type || sub_type == 0) && (e.data == v || v == 0) {
			answer += e.data
		}
	}

	return answer
}

// get the first effect of a type off of an effect list
func get_effect(what, type_, sub_type, v int) int {
	// validity checks.
	if !valid_box(what) {
		return 0
	}
	// only characters, locations, and sublocations have effects on them.

	// look for the effect
	for _, e := range effects(what) {
		if e.type_ == type_ && (e.subtype == sub_type || sub_type == 0) && (e.data == v || v == 0) {
			return e.data
		}
	}

	return 0
}

// update something's effects for the passage of a day.
func update_effects(what int) {
	// validity checks.
	if !valid_box(what) {
		return
	}

	// only characters, locations, and sublocations have effects on them.
	el := effects(what)

	// this is mean - if you hide something for too long, you can lose it permanently,
	// so check hidden stuff every 30 days...
	for _, e := range el {
		if !(e.type_ == ef_hide_money || e.type_ == ef_hide_item) {
			continue
		}
		if e.days%30 == 0 && rnd(1, 100) < (e.days/30) {
			// we've permanently forgotten this hidden item.
			wout(what, "You have a nagging feeling you've forgotten something.")
			// set the days so that the main loop will delete it.
			e.days = 0
			continue
		}
		// bump days to generate a reminder later
		e.days += 2 // two because the main loop decrements by one
	}

	// go through the list, decrement the "days" for each effect, and delete any that are now expired.
	// we go through the list backwards, because deleting removes elements and shifts the list.
	for i := len(el) - 1; i >= 0; i-- {
		// disappearing golems.
		if el[i].type_ == ef_kill_dirt_golem && el[i].days < 2 {
			// remove a dirt golem from this guy.
			if consume_item(what, item_dirt_golem, 1) {
				wout(what, "A dirt golem suddenly crumbles to dust!")
			} else {
				wout(gm_player, "No dirt golem available to crumble for %s.", box_name(what))
			}
		} else if el[i].type_ == ef_kill_flesh_golem && el[i].days < 2 {
			// remove a flesh golem from this guy.
			if consume_item(what, item_flesh_golem, 1) {
				wout(what, "A flesh golem suddenly crumbles to dust!")
			} else {
				wout(gm_player, "No flesh golem available to crumble for %s.", box_name(what))
			}
		}

		// when tap health runs out, you suffer a wound
		if el[i].type_ == ef_tap_wound && el[i].days == 1 {
			wout(what, "You feel a aura-rending blast hit your body!")
			add_char_damage(what, el[i].data, MATES)
		}

		// generic case - delete when the days have expired
		el[i].days--
		if el[i].days <= 0 {
			el = append(el[:i], el[i+1:]...)
		}
	}

	// we might now have an empty list;
	// if so we could reclaim it, but for the moment we leave it alone.
	bx[what].effects = el
}
