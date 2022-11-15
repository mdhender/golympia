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
	// defines to turn on/control the Oly 3 changes
	ATTACK_LIMIT           = 4    // The number of attacks in a row permitted.
	CONTROL_SKILL          = true // New skill to control men in battle.
	DEFAULT_CONTROLLED     = 10   // Number controlled if you don't have the skill.
	GARRISON_CONTROLLED    = 20   // Default controlled by garrison
	DEFENDER_CONTROL_BONUS = 20   // On defense, control extra men.
	TACTICS_FACTOR         = 0.02 // % given with each week of study.
	TACTICS_LIMIT          = 2.0  // Maximum bonus from tactics.
	CITY_DEFENSE_BONUS     = 1.25 // Terrain bonuses to defense
	FOREST_DEFENSE_BONUS   = 1.50
	MOUNTAIN_DEFENSE_BONUS = 2.00
	SWAMP_DEFENSE_BONUS    = 0.75

	// scrolls & Books
	SCROLL_CHANCE = 50 // % for a scroll or book to appear in a city

	// shipcraft stuff
	HULL_CAPACITY = 2500
	FORT_WEIGHT   = 250
	SAIL_WEIGHT   = 250
	KEEL_WEIGHT   = 100

	SAILS_PER_HULL = 3
	KEELS_PER_HULL = 1
	PORTS_PER_HULL = 3

	SHIP_FORTS_PROTECT = 10
	SHIP_FORTS_BONUS   = 5

	ROWERS_PER_PORT = 4

	// moat Constants
	MOAT_MATERIAL = 200
	MOAT_EFFORT   = 4000
)
