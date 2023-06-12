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

import "log"

// create a couple of stacks and have them battle it out a few hundred times and report the results.
func TestCombat(dirLibrary string) error {
	libdir = dirLibrary

	// lock up; prevents multiple TAGs running simultaneously.
	lock_tag()

	if err := call_init_routines(); err != nil {
		return err
	}
	if err := load_db(); err != nil {
		return err
	}

	open_logfile()
	sum_a := 0

	for i := 0; i < 100; i++ {
		a := new_ent(T_char, 0)
		b := new_ent(T_char, 0)
		c := new_ent(T_char, 0)
		d := new_ent(T_char, 0)

		p := p_char(a)
		p.behind = 9
		p.health = 100
		p.attack = 80
		p.defense = 80
		p.break_point = 0
		s := p_skill_ent(a, sk_control_battle)
		s.know = SKILL_know
		s.experience = 500
		s = p_skill_ent(a, sk_use_beasts)
		s.know = SKILL_know
		s.experience = 500
		set_where(a, 10201)
		set_lord(a, gm_player, LOY_oath, 1)
		set_name(a, "A")

		p = p_char(b)
		p.behind = 0
		p.health = 100
		p.attack = 80
		p.defense = 80
		p.break_point = 0
		s = p_skill_ent(b, sk_control_battle)
		s.know = SKILL_know
		s.experience = 500
		s = p_skill_ent(b, sk_use_beasts)
		s.know = SKILL_know
		s.experience = 500
		set_where(b, 10201)
		set_lord(b, gm_player, LOY_oath, 1)
		set_name(b, "B")

		p = p_char(c)
		p.behind = 0
		p.health = 100
		p.attack = 360
		p.defense = 360
		p.break_point = 0
		s = p_skill_ent(c, sk_control_battle)
		s.know = SKILL_know
		s.experience = 500
		s = p_skill_ent(c, sk_use_beasts)
		s.know = SKILL_know
		s.experience = 500
		set_where(c, 10201)
		set_lord(c, gm_player, LOY_oath, 1)
		set_name(c, "C")

		p = p_char(d)
		p.behind = 9
		p.health = 100
		p.attack = 80
		p.defense = 80
		p.break_point = 0
		s = p_skill_ent(d, sk_control_battle)
		s.know = SKILL_know
		s.experience = 500
		s = p_skill_ent(d, sk_use_beasts)
		s.know = SKILL_know
		s.experience = 500
		set_where(d, 10201)
		set_lord(d, gm_player, LOY_oath, 1)
		set_name(d, "D")

		// A & B stacked together and C & D stacked together.
		// A & D are behind.

		gen_item(a, item_soldier, 10)
		gen_item(a, item_angel, 15)
		gen_item(c, item_nazgul, 3)

		if result := regular_combat(c, a, 0, 0); result != 0 {
			sum_a++
		}
	}
	close_logfile()
	log.Printf("Sum = %d.\n", sum_a)

	return nil
}
