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

/*
 * Heroism: Skills to make an individual noble an adventurer.
 * Many suggested by Jeremy Maiden (JeremyM@ctxuk.citrix.com)
 *
 * re Survive Fatal Wound (return NP version)
 * of Defense
 * of Swordplay
 *
 * of Avoid wounds (decrease amount of wounds by %) +
 * of Avoid Illness (reduce chance of illness)
 * of Improved Recovery (recover from illness)
 * of Personal Fight to the Death (don't quit at first wound) +
 * of Forced Marching (improved foot speed) +
 * re Acute Senses (avoid traps) +
 * re Improved Explore (subdivided for terrain types?) +
 * re Uncanny Accuracy (archery to the back rank) +
 * re Extra Attacks (increase # of attacks each combat round) +
 * re Blinding Speed (1 attack during Special Attack round)
 *
 * Skills marked with (+) only usable if you are not commanding men.
 *
 */

/*
 *  Extra attacks.  Get an extra attack after every 8 weeks of practice.
 *
 *  I don't think we actually need to have a v_ and a d_ for this,
 *  since using "practice"  should do the trick.  Just need to have
 *  the skill set up appropriately in Lib/skill.  Note practice cost
 *  of $50.
 */

/*
 *  Avoid wounds.  Get -% to wound damage up to a maximum of
 *  -80%, -2% per week practiced.
 *
 *  I don't think we actually need to have a v_ and a d_ for this,
 *  since using "practice"  should do the trick.  Just need to have
 *  the skill set up appropriately in Lib/skill.  Note practice cost
 *  of $50.
 */

/*
 *  Avoid illness is exactly the same as avoid wounds.  Get 2% to
 *  avoid illness up to a maximum of 80%, 2% per week practiced.
 */

/*
 *  Improved Recovery gives you an additional +1% to recover per
 *  week for each level of expertise, up to a maximum of +50%
 */

/*
 *  Acute Senses gives you a 2% chance to avoid a trap in a province
 *  for every week's practice to a maximum of 80%.
 */

/*
 *  Improved Explore gives a +1% chance to explore for each week's
 *  practice, up to a maximum of +50%.
 */

/*
 *  Uncanny Accuracy permits the hero to target his missile attacks
 *  against any missile troops in the back row of the opponent.  Is this
 *  easily doable?
 */

/*
 *  Blinding Speed gives the Hero a melee attack during the special
 *  attacks phase.
 */

/*
 *  Forced Marching makes your next move be at riding speed, at the
 *  cost of that many days travel in health.
 */
func v_forced_march(c *command) int {
	/*  Has a $10 cost */
	if charge(c.who, 10) == FALSE {
		wout(c.who, "Can't afford %s to prepare for a forced march.", gold_s(10))
		return FALSE
	}

	if get_effect(c.who, ef_forced_march, 0, 0) != FALSE {
		wout(c.who, "You are already prepared for a forced march.")
	} else if add_effect(c.who, ef_forced_march, 0, -1, 1) == FALSE {
		wout(c.who, "For some reason, you cannot prepare for a forced march!")
		return FALSE
	} else {
		wout(c.who, "Now prepared to make the next move as a forced march.")
	}

	return TRUE
}

/*
 *  Personal Fight to the Death allows you to set the level of wounds
 *  that will cause you to quit a fight, rather than automatically
 *  quit after the first wound.
 */
func v_personal_fight_to_death(c *command) int {
	flag := c.a
	if flag < 0 {
		flag = 0
	} else if flag > 100 {
		flag = 100
	}
	p_char(c.who).personal_break_point = schar(flag)
	wout(c.who, "%s will now leave battle when he has %d health remaining", box_name(c.who), flag)

	return TRUE
}
