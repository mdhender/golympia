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

import "strings"

/*
 *  Fri Jan  3 09:25:12 1997 -- Scott Turner
 *
 *  Add/Remove sails...
 *
 */
func d_add_sails(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add a sail to it.")
		return FALSE
	}

	/*
	 *  Can this ship fit another sail?
	 *
	 */
	if ship.sails >= ship.hulls*SAILS_PER_HULL {
		wout(c.who, "There is no room to fit more sails to this ship.")
		return FALSE
	}

	/*
	 *  Must have the capacity for another sail
	 *
	 */
	if ship_cap(where) < SAIL_WEIGHT {
		wout(c.who, "This ship does not have the capacity to carry another sail.")
		return FALSE
	}

	/*
	 *  Materials?  Encoded in the skill under required items.  d_use
	 *  also takes care of consuming these things.  We can also encode
	 *  the workers there.  So at this point we just create the sail and
	 *  add it to the boat.
	 */
	ship.sails++
	wout(c.who, "%s now equipped with %s sail%s.",
		box_name(where), nice_num(ship.sails),
		or_string(ship.sails == 1, "", "s"))
	return TRUE
}

func v_add_sails(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add a sail to it.")
		return FALSE
	}

	/*
	 *  Can this ship fit another sail?
	 *
	 */
	if ship.sails >= ship.hulls*SAILS_PER_HULL {
		wout(c.who, "There is no room to fit more sails to this ship.")
		return FALSE
	}

	/*
	 *  Must have the capacity for another sail
	 *
	 */
	if ship_cap(where) < SAIL_WEIGHT {
		wout(c.who, "This ship does not have the capacity to carry another sail.")
		return FALSE
	}

	return TRUE
}

func d_remove_sails(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a sail.")
		return FALSE
	}

	/*
	 *  Does this ship have a sail to remove?
	 *
	 */
	if ship.sails == FALSE {
		wout(c.who, "There are no sails on this ship.")
		return FALSE
	}

	ship.sails--
	wout(c.who, "%s now equipped with %s sail%s.",
		box_name(where), nice_num(ship.sails),
		or_string(ship.sails == 1, "", "s"))
	return TRUE
}

func v_remove_sails(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a sail.")
		return FALSE
	}

	/*
	 *  Does this ship have a sail to remove?
	 *
	 */
	if ship.sails == FALSE {
		wout(c.who, "There are no sails on this ship.")
		return FALSE
	}

	return TRUE
}

/*
 *  Fri Jan  3 09:25:12 1997 -- Scott Turner
 *
 *  Add/Remove forts...
 *
 */
func d_add_forts(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add fortification to it.")
		return FALSE
	}

	/*
	 *  Must have the capacity for another fort
	 *
	 */
	if ship_cap(where) < FORT_WEIGHT {
		wout(c.who, "This ship does not have the capacity for more fortification.")
		return FALSE
	}

	/*
	 *  Materials?  Encoded in the skill under required items.  d_use
	 *  also takes care of consuming these things.  We can also encode
	 *  the workers there.  So at this point we just create the fort and
	 *  add it to the boat.
	 */
	ship.forts++
	wout(c.who, "%s now equipped with %s unit%s of fortification.",
		box_name(where), nice_num(ship.forts),
		or_string(ship.forts == 1, "", "s"))

	/*
	 *  Fix the fortification level for the ship.
	 *
	 */
	rp_subloc(where).defense = int(float64(ship.forts)/float64(ship.hulls)*SHIP_FORTS_BONUS + 0.5)
	wout(c.who, "New defense rating for %s is %d.", box_name(where), rp_subloc(where).defense)

	return TRUE
}

func v_add_forts(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add fortification to it.")
		return FALSE
	}

	/*
	 *  Must have the capacity for another fort
	 *
	 */
	if ship_cap(where) < FORT_WEIGHT {
		wout(c.who, "This ship does not have the capacity for more fortification.")
		return FALSE
	}

	return TRUE
}

func d_remove_forts(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a fort.")
		return FALSE
	}

	/*
	 *  Does this ship have a fort to remove?
	 *
	 */
	if ship.forts == FALSE {
		wout(c.who, "There are no fortifications on this ship.")
		return FALSE
	}

	ship.forts--
	wout(c.who, "%s now equipped with %s unit%s of fortification.",
		box_name(where), nice_num(ship.forts),
		or_string(ship.forts == 1, "", "s"))

	/*
	 *  Fix the fortification level for the ship.
	 *
	 */
	rp_subloc(where).defense = int(float64(ship.forts)/float64(ship.hulls)*SHIP_FORTS_BONUS + 0.5)
	wout(c.who, "New defense rating for %s is %d.", box_name(where), rp_subloc(where).defense)

	return TRUE
}

func v_remove_forts(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a fort.")
		return FALSE
	}

	/*
	 *  Does this ship have a fort to remove?
	 *
	 */
	if ship.forts == FALSE {
		wout(c.who, "There are no fortifications on this ship.")
		return FALSE
	}

	return TRUE
}

/*
 *  Fri Jan  3 09:25:12 1997 -- Scott Turner
 *
 *  Add/Remove keels...
 *
 */
func d_add_keels(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add a keel to it.")
		return FALSE
	}

	/*
	 *  Can this ship fit another keel?
	 *
	 */
	if ship.keels >= ship.hulls*KEELS_PER_HULL {
		wout(c.who, "There is no room to fit more keels to this ship.")
		return FALSE
	}

	/*
	 *  Must have the capacity for another keel
	 *
	 */
	if ship_cap(where) < KEEL_WEIGHT {
		wout(c.who, "This ship does not have the capacity for more keel.")
		return FALSE
	}

	/*
	 *  Materials?  Encoded in the skill under required items.  d_use
	 *  also takes care of consuming these things.  We can also encode
	 *  the workers there.  So at this point we just create the keel and
	 *  add it to the boat.
	 */
	ship.keels++
	wout(c.who, "%s now equipped with %s unit%s of keel.",
		box_name(where), nice_num(ship.keels),
		or_string(ship.keels == 1, "", "s"))

	return TRUE
}

func v_add_keels(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add a keel to it.")
		return FALSE
	}

	/*
	 *  Can this ship fit another keel?
	 *
	 */
	if ship.keels >= ship.hulls*KEELS_PER_HULL {
		wout(c.who, "There is no room to fit more keels to this ship.")
		return FALSE
	}

	/*
	 *  Must have the capacity for another keel
	 *
	 */
	if ship_cap(where) < KEEL_WEIGHT {
		wout(c.who, "This ship does not have the capacity for more keel.")
		return FALSE
	}

	return TRUE
}

func d_remove_keels(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a keel.")
		return FALSE
	}

	/*
	 *  Does this ship have a keel to remove?
	 *
	 */
	if ship.keels == FALSE {
		wout(c.who, "There are no keels on this ship.")
		return FALSE
	}

	ship.keels--
	wout(c.who, "%s now equipped with %s unit%s of keel.",
		box_name(where), nice_num(ship.keels),
		or_string(ship.keels == 1, "", "s"))
	return TRUE
}

func v_remove_keels(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a keel.")
		return FALSE
	}

	/*
	 *  Does this ship have a keel to remove?
	 *
	 */
	if ship.keels == FALSE {
		wout(c.who, "There are no keels on this ship.")
		return FALSE
	}

	return TRUE
}

/*
 *  Fri Jan  3 09:25:12 1997 -- Scott Turner
 *
 *  Add/Remove ports...
 *
 */
func d_add_ports(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add a port to it.")
		return FALSE
	}

	/*
	 *  Can this ship fit another port?
	 *
	 */
	if ship.ports >= ship.hulls*PORTS_PER_HULL {
		wout(c.who, "There is no room to fit more ports to this ship.")
		return FALSE
	}

	/*
	 *  Materials?  Encoded in the skill under required items.  d_use
	 *  also takes care of consuming these things.  We can also encode
	 *  the workers there.  So at this point we just create the port and
	 *  add it to the boat.
	 */
	ship.ports++
	wout(c.who, "%s now equipped with %s rowing port%s.",
		box_name(where), nice_num(ship.ports),
		or_string(ship.ports == 1, "", "s"))
	return TRUE
}

func v_add_ports(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to add a port to it.")
		return FALSE
	}

	/*
	 *  Can this ship fit another port?
	 *
	 */
	if ship.ports >= ship.hulls*PORTS_PER_HULL {
		wout(c.who, "There is no room to fit more ports to this ship.")
		return FALSE
	}

	return TRUE
}

func d_remove_ports(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a port.")
		return FALSE
	}

	/*
	 *  Does this ship have a port to remove?
	 *
	 */
	if ship.ports == FALSE {
		wout(c.who, "There are no ports on this ship.")
		return FALSE
	}

	ship.ports--
	wout(c.who, "%s now equipped with %s rowing port%s.",
		box_name(where), nice_num(ship.ports),
		or_string(ship.ports == 1, "", "s"))

	return TRUE
}

func v_remove_ports(c *command) int {
	where := subloc(c.who)
	ship := rp_ship(where)

	/*
	 *  In a ship?
	 *
	 */
	if subkind(where) != sub_ship || ship == nil {
		wout(c.who, "You must be in a finished ship to remove a port.")
		return FALSE
	}

	/*
	 *  Does this ship have a port to remove?
	 *
	 */
	if ship.ports == FALSE {
		wout(c.who, "There are no ports on this ship.")
		return FALSE
	}

	return TRUE
}

/*
 *  Mon Jan  6 11:56:21 1997 -- Scott Turner
 *
 *  Ship summary at end of turn reveals everything about the ships
 *  you've been in this turn...
 *
 */
func ship_summary(pl int) {
	var first, sc, sw = 0, 0, 0

	for _, i := range known_sparse_loop(p_player(pl).locs) {
		if valid_box(i) && subkind(i) == sub_ship {
			ship := rp_ship(i)
			if ship == nil {
				continue
			}
			if first == 0 {
				tagout(pl, "<tag type=ship_report pl=%d>", pl)
				tagout(pl, "<tag type=header>", pl)
				out(pl, "")
				out(pl, "Ship Report:")
				out(pl, "")
				out(pl, "ship  hulls  ports  sails  keels  forts  damge  capacity")
				out(pl, "----  -----  -----  -----  -----  -----  -----  --------")
				tagout(pl, "</tag type=header>", pl)

				first = 1
			}
			sw = ship_weight(i)
			sc = ship_cap(i)

			tagout(pl, "<tag type=ship ship=%d hulls=%d ports=%d sails=%d keels=%d forts=%d damage=%d weight=%d capacity=%d>",
				i,
				ship.hulls,
				ship.ports,
				ship.sails,
				ship.keels,
				ship.forts,
				loc_damage(i),
				sw, sc)

			out(pl, "%4s  %5d  %5d  %5d  %5d  %5d  %4d%%  %s/%s (%d%%)",
				box_code_less(i),
				ship.hulls,
				ship.ports,
				ship.sails,
				ship.keels,
				ship.forts,
				loc_damage(i),
				comma_num(sw), comma_num(sc), sw*100/sc)

			tagout(pl, "</tag type=ship>")

		}
	}

	if first != 0 {
		tagout(pl, "</tag type=ship_report pl=%d>", pl)
	}

}

/*
 *  Tue Jan  7 11:58:39 1997 -- Scott Turner
 *
 *  Check a ship to see if it is going to sustain some storm damage.
 *
 */
func ship_storm_check(ship int) {
	where := province(ship)
	//coastal := 0
	damage_chance := 15
	storm_strength := weather_here(where, sub_rain)

	if !is_ship(ship) {
		return
	}
	if storm_strength < 1 {
		return
	}
	if storm_strength > 50 {
		storm_strength = 50
	}

	s := rp_ship(ship)

	/*
	 *  Storm damage depends on whether we're coastal.
	 *  3 == non-coastal
	 */
	if near_rocky_coast(where) == 3 {
		damage_chance = 35
	}

	/*
	 *  Keel protection.
	 *
	 */
	if s != nil && s.keels != FALSE {
		damage_chance -= int((float64(s.keels) / float64(s.hulls)) * 30)
	}

	/*
	 *  Protection from a Priest of Eres.
	 *
	 */
	if priest_in_stack(first_character(ship), sk_eres) {
		damage_chance -= 5
	}

	if storm_strength < 1 || damage_chance < 1 || rnd(1, 100) > damage_chance {
		return
	}

	/*
	 *  Pen sailors take 1/2 damage.
	 *
	 */
	if nation(first_character(ship)) != 0 && strings.HasPrefix(rp_nation(nation(first_character(ship))).name, "Pen") {
		wout(first_character(ship), "You sail cleanly through the heart of the storm.")
		return
	}

	wout(ship, "A storm batters the ship and causes some damage.")
	add_structure_damage(ship, rnd(1, storm_strength))

}

/*
 *  Rams
 *  Fri Oct 30 08:44:55 1998 -- Scott Turner
 *
 *  These probably don't work right now.
 *
 */
func v_add_ram(c *command) int {
	ship := subloc(c.who)

	/*
	 *  In a ship?
	 *
	 */
	if ship == 0 || subkind(ship) != sub_ship {
		wout(c.who, "You must be in a finished ship to remove a ram.")
		return FALSE
	}

	if ship_has_ram(ship) != FALSE {
		wout(c.who, "%s already has a ram.", box_name(ship))
	}

	wout(c.who, "Work to add an iron-tipped ram to this vessel.")
	return TRUE
}

func d_add_ram(c *command) int {
	ship := subloc(c.who)

	/*
	 *  In a ship?
	 *
	 */
	if ship == 0 || subkind(ship) != sub_ship {
		wout(c.who, "You must be in a finished ship to remove a ram.")
		return FALSE
	}

	if ship_has_ram(ship) != FALSE {
		wout(c.who, "%s already has a ram.", box_name(ship))
		return FALSE
	}

	p_ship(ship).galley_ram = 1

	wout(c.who, "%s has been fitted with a ram!", box_name(ship))
	wout(ship, "%s has been fitted with a ram!", box_name(ship))

	return TRUE
}

func v_remove_ram(c *command) int {
	ship := subloc(c.who)

	/*
	 *  In a ship?
	 *
	 */
	if ship == 0 || subkind(ship) != sub_ship {
		wout(c.who, "You must be in a finished ship to remove a ram.")
		return FALSE
	}

	if ship_has_ram(ship) != FALSE {
		wout(c.who, "%s already has a ram.", box_name(ship))
	}

	wout(c.who, "Work to add an iron-tipped ram to this vessel.")
	return TRUE
}

func d_remove_ram(c *command) int {
	ship := subloc(c.who)

	/*
	 *  In a ship?
	 *
	 */
	if ship == 0 || subkind(ship) != sub_ship {
		wout(c.who, "You must be in a finished ship to remove a ram.")
		return FALSE
	}

	if ship_has_ram(ship) != FALSE {
		wout(c.who, "%s already has a ram.", box_name(ship))
		return FALSE
	}

	p_ship(ship).galley_ram = 1

	wout(c.who, "%s has been fitted with a ram!", box_name(ship))
	wout(ship, "%s has been fitted with a ram!", box_name(ship))

	return TRUE
}
