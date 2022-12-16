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

import (
	"fmt"
	"log"
	"sort"
)

func deliver_lore(who, num int) {
	switch kind(num) {
	case T_skill:
		// turned off -- in rules now.
		//deliver_skill_lore(who, num, FALSE, FALSE);
	case T_item:
		deliver_lore_sheet(who, item_lore(num), num, false)
	default:
		deliver_lore_sheet(who, num, num, false)
	}

	out(who, "")
}

func deliver_lore_sheet(who, num, display_num int, use_texi bool) {
	out(who, "")
	if use_texi {
		wout(who, "@node %s", box_name(display_num))
		if skill_school(num) == num {
			wout(who, "@appendixsec %s", box_name(display_num))
		} else {
			wout(who, "@appendixsubsec %s", box_name(display_num))
		}
	} else {
		match_lines(who, box_name(display_num))
	}

	//if kind(num) == T_skill && skill_school(num) != num {
	//	match_lines(who, sout("Skill:  %s", box_name(num)))
	//} else {
	//	match_lines(who, sout("Lore for %s", box_name(num)))
	//}
	var fnam string
	if kind(num) == T_skill {
		fnam = sout("%s/lore/%d/%d", libdir, skill_school(num), num)
	} else {
		fnam = sout("%s/lore/etc/%d", libdir, num)
	}
	fp, err := fopen(fnam, "r")
	if err != nil {
		out(who, "<lore sheet not available>")
		log.Printf("can't open %s: ", fnam)
		return
	}

	/*  This is ended in do_skill_header */
	if use_texi {
		wout(who, "@cartouche")
		wout(who, "@display")
		wout(who, "@group")
	}

	first_blank := false
	for line := getlin(fp); line != nil; line = getlin(fp) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if line[0] == '$' {
			lore_function(who, string(line[1:]))
			continue
		}

		if len(line) < 2 && !first_blank && kind(num) == T_skill {
			// do the auto-header for skills before the first blank line.
			do_skill_header(who, num, use_texi)
			first_blank = true
		}

		// copy line into new buffer, substituting something for '@'
		var q []byte
		for p := line; len(p) != 0; {
			if p[0] == '@' {
				q, p = append(q, []byte(box_code_less(display_num))...), p[1:]
				for len(p) != 0 && p[0] == '@' {
					p = p[1:]
				}
				continue
			}
			q, p = append(q, p[0]), p[1:]
		}
		out(who, "%s", string(q))
	}

	fclose(fp)
}

func deliver_skill_lore(who, sk int, show_research, use_texi bool) {
	deliver_lore_sheet(who, sk, sk, use_texi)

	p := rp_skill(sk)
	if p == nil {
		return
	}

	if len(p.offered) != 0 {
		out(who, "")
		wout(who, "The following teachable skills may be studied directly once %s is known:", just_name(sk))
		out(who, "")

		if use_texi {
			wout(who, "@example")
			wout(who, "@group")
		}
		indent += 3
		out(who, "%-*s  %-34s %13s", CHAR_FIELD, "num", "skill", "time to learn")
		out(who, "%-*s  %-34s %13s", CHAR_FIELD, "---", "-----", "-------------")

		for i := 0; i < len(p.offered); i++ {
			out_skill_line(who, p.offered[i])
		}

		indent -= 3
		if use_texi {
			wout(who, "@end group")
			wout(who, "@end example")
		}

		out(who, "")
	}

	if len(p.research) != 0 {
		wout(who, "The following unteachable skills may be studied directly once %s is known:", just_name(sk))
		out(who, "")

		if use_texi {
			wout(who, "@example")
			wout(who, "@group")
		}
		indent += 3
		out(who, "%-*s  %-34s %13s", CHAR_FIELD, "num", "skill", "time to learn")
		out(who, "%-*s  %-34s %13s", CHAR_FIELD, "---", "-----", "-------------")

		for i := 0; i < len(p.research); i++ {
			out_skill_line(who, p.research[i])
		}

		indent -= 3
		if use_texi {
			wout(who, "@end group")
			wout(who, "@end example")
		}

		out(who, "")
	}

	if len(p.guild) != 0 {
		wout(who, "The following guild skills may be studied by guild members:")
		out(who, "")

		if use_texi {
			wout(who, "@example")
			wout(who, "@group")
		}
		indent += 3
		out(who, "%-*s  %-34s %13s", CHAR_FIELD, "num", "skill", "time to learn")
		out(who, "%-*s  %-34s %13s", CHAR_FIELD, "---", "-----", "-------------")

		for i := 0; i < len(p.guild); i++ {
			out_skill_line(who, p.guild[i])
		}

		indent -= 3
		if use_texi {
			wout(who, "@end group")
			wout(who, "@end example")
		}

		out(who, "")
	}
}

func do_skill_header(who, num int, use_texi bool) {
	wout(who, "Time to learn: %d days.", rp_skill(num).time_to_learn)

	if rp_skill(num).np_req != 0 {
		wout(who, "NPs to learn: %d.", rp_skill(num).np_req)
	}
	if (rp_skill(num).flags & COMBAT_SKILL) != 0 {
		wout(who, "Time to use: Automatic in combat")
	} else if rp_skill(num).time_to_use == -1 {
		wout(who, "Time to use: Variable.")
	} else if rp_skill(num).time_to_use != 0 {
		var s string
		if rp_skill(num).time_to_use > 1 {
			s = "s"
		}
		wout(who, "Time to use: %d day%s.", rp_skill(num).time_to_use, s)
	}

	if rp_skill(num).practice_time != 0 {
		var s string
		if rp_skill(num).practice_time > 1 {
			s = "s"
		}
		wout(who, "Time to practice: %d day%s.", rp_skill(num).practice_time, s)
		wout(who, "Cost to practice: %s.", gold_s(rp_skill(num).practice_cost))
	}

	if rp_skill(num).piety != 0 {
		if magic_skill(num) {
			wout(who, "Aura to use: %d.", rp_skill(num).piety)
		} else {
			wout(who, "Piety to use: %d.", rp_skill(num).piety)
		}
	}

	if (rp_skill(num).flags&REQ_HOLY_SYMBOL) != 0 && (rp_skill(num).flags&REQ_HOLY_PLANT) != 0 {
		wout(who, "Required to use: Holy symbol, holy plant")
	} else if (rp_skill(num).flags & REQ_HOLY_SYMBOL) != 0 {
		wout(who, "Required to use: Holy symbol")
	} else if (rp_skill(num).flags & REQ_HOLY_PLANT) != 0 {
		wout(who, "Required to use: Holy plant")
	}

	if (rp_skill(num).flags & COMBAT_SKILL) != 0 {
		wout(who, "This is a combat spell.")
	}
	/*
	 *  "Requires"
	 *
	 */
	var consume_string, continuation_string string
	first := true
	if rp_skill(num) != nil && len(rp_skill(num).req) != 0 {
		l := rp_skill(num).req
		for i := 0; i < len(l); i++ {
			consume_string = ""
			continuation_string = ""
			if l[i].consume == REQ_NO {
				consume_string = " (not consumed)"
			}
			if l[i].consume == REQ_OR {
				continuation_string = " *OR*"
			}
			if first {
				wout(who, "Requires: %s %s%s%s", nice_num(l[i].qty), plural_item_box(l[i].item, l[i].qty), consume_string, continuation_string)
				first = false
			} else {
				wout(who, "          %s %s%s%s", nice_num(l[i].qty), plural_item_box(l[i].item, l[i].qty), consume_string, continuation_string)
			}
		}
	}
	if rp_skill(num) != nil && rp_skill(num).produced != 0 {
		wout(who, "Produces: %s", box_name(rp_skill(num).produced))
	}

	if use_texi {
		wout(who, "@end group")
		wout(who, "@end display")
		wout(who, "@end cartouche")
	}
}

func gm_show_all_skills(pl int, use_texi bool) {
	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "")
	out(pl, "Skill listing:")
	out(pl, "--------------")

	out(pl, "")
	out(pl, "Skill schools:")
	out(pl, "")
	if use_texi {
		out(pl, "@itemize")
	}

	for _, sk := range loop_skill() {
		if skill_school(sk) == sk {
			if use_texi {
				out(pl, "@item")
			}
			out_skill_line(pl, sk)
		}
	}

	if use_texi {
		out(pl, "@end itemize")
	}
	out(pl, "")

	if use_texi {
		out(pl, "@table @asis")
	}
	for _, sk := range loop_skill() {
		if skill_school(sk) == sk {
			if use_texi {
				out(pl, "@item %s", box_name(sk))
				out(pl, "@example")
			} else {
				out(pl, "%s", box_name(sk))
			}
			indent += 3

			for _, i := range loop_skill() {
				if skill_school(i) == sk && i != sk {
					out_skill_line(pl, i)
				}
			}

			indent -= 3
			out(pl, "")
			if use_texi {
				out(pl, "@end example")
				out(pl, "")
			}
		}
	}

	if use_texi {
		out(pl, "@end table")
	}

	//// output lore sheets for all skills
	//for _, sk := range loop_skill() {
	//	if skill_school(sk) == sk {
	//		deliver_skill_lore(pl, sk, true, use_texi)
	//		out(pl, "")
	//		for _, subsk := range loop_skill() {
	//			if subsk != sk && skill_school(subsk) == sk {
	//				deliver_skill_lore(pl, subsk, true, use_texi)
	//				out(pl, "")
	//			}
	//		}
	//	}
	//}
	//for _, sk := range loop_skill() {
	//	if skill_school(sk) != sk {
	//		deliver_skill_lore(pl, sk, true, use_texi)
	//		out(pl, "")
	//	}
	//}

	out_path = 0
	out_alt_who = 0
}

func lore_comp(a, b int) int {
	return a - b
}

func lore_function(who int, does string) {
	switch does {
	case "animal_fighters":
		for _, i := range loop_item() {
			if item_animal(i) != FALSE && is_fighter(i) != FALSE {
				out(who, "    %-20s %s", box_name(i), sout("(%d,%d,%d)", item_attack(i), item_defense(i), item_missile(i)))
			}
		}
	case "capturable_animals":
		for _, i := range loop_item() {
			if item_animal(i) != FALSE && item_capturable(i) != FALSE {
				out(who, "    %s", box_name(i))
			}
		}
	default:
		panic(fmt.Sprintf("bad lore sheet function: %q", does))
	}
}

func np_req_s(skill int) string {
	np := skill_np_req(skill)
	if np < 1 {
		return ""
	} else if np == 1 {
		return ", 1 NP req'd"
	}
	return sout(", %d NP req'd", np)
}

func out_skill_line(who, sk int) {
	out(who, "%-*s  %-34s %s%s", CHAR_FIELD, box_code_less(sk), just_name(sk), weeks(learn_time(sk)), np_req_s(sk))
}

// show a player a lore sheet.
// set anyway to show them the lore sheet even if they've seen it before.
func queue_lore(who, num int, anyway bool) {
	pl := player(who)
	if kind(pl) != T_player {
		panic("assert(kind(pl) == T_player)")
	}
	if test_known(pl, num) && !anyway {
		return
	}
	p := p_player(pl)
	p.deliverLore = append(p.deliverLore, num)
	set_known(pl, num)
}

func scan_char_item_lore() {
	for _, who := range loop_char() {
		for _, e := range inventory_loop(who) {
			lore := item_lore(e.item)
			if lore != 0 && !test_known(who, e.item) {
				queue_lore(who, e.item, false)
			}
		}
	}
}

func scan_char_skill_lore() {
	for _, who := range loop_char() {
		for _, e := range loop_char_skill(who) {
			queue_lore(who, e.skill, false)
		}
	}
}

func show_lore_sheets() {
	stage("show_lore_sheets()")

	out_path = MASTER
	out_alt_who = OUT_LORE

	for _, pl := range loop_player() {
		p := rp_player(pl)
		if p == nil || len(p.deliverLore) <= 0 {
			continue
		}

		sort.Ints(p.deliverLore)
		for i := 0; i < len(p.deliverLore); i++ {
			// weed out duplicates in p.deliver_lore
			if i > 0 && p.deliverLore[i] == p.deliverLore[i-1] {
				continue
			}
			deliver_lore(pl, p.deliverLore[i])
		}
	}

	out_path = 0
	out_alt_who = 0
}
