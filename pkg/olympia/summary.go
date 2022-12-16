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

import "sort"

const (
	SHOW_TOP = false
)

var (
	ranks []int

	ncontrolled = 0
	nother      = 0
	nplayers    = 0
)

func collect_game_totals() {
	nplayers = 0
	ncontrolled = 0
	nother = 0

	clear_temps(T_char)

	for _, i := range loop_kind(T_player) {
		if subkind(i) != sub_pl_regular {
			continue
		}

		nplayers++
		for _, j := range loop_units(i) {
			ncontrolled++
			bx[j].temp = 1
		}
	}

	for _, i := range loop_kind(T_char) {
		if bx[i].temp == 0 {
			nother++
		}
	}
}

func out_rank_mine_top(who, num int, title string, total int) {
	var s string
	if total == 0 {
		s = comma_num(bx[ranks[num]].temp)
	} else {
		s = sout("%s/%s", comma_num(bx[ranks[num]].temp), comma_num(total))
	}

	var top string
	if SHOW_TOP {
		top = top_rank()
	}

	tagout(who, "<tag type=faction_entry title=\"%s\" value=\"%s\" rank=\"%s\" top=\"%s\">", title, s, ordinal(ranking(num)), top)

	j := 20 + 1 + 11 + 2 + 5 + 2 // total widths of columns
	wiout(who, j, "%-20s %16s  %-5s  %s", title, s, ordinal(ranking(num)), top)
}

func out_ranking(title string, total int) {
	ranks = nil

	for _, i := range loop_kind(T_player) {
		if subkind(i) != sub_pl_regular {
			continue
		}
		ranks = append(ranks, i)
	}

	// reverse sort ranks
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))

	for i := 0; i < len(ranks); i++ {
		out_rank_mine_top(ranks[i], i, title, total)
	}

	//if top_rank > 0 {		/* for the GM */
	//	out_rank_mine_top(gm_player, 0, title, total);
	//}
}

// comparison is reversed so maximum will be first in list
func rank_comp(a, b int) int {
	return bx[b].temp - bx[a].temp
}

// return the ranking of a player within ranks[], allowing for ties.
func ranking(n int) int {
	for n > 0 && bx[ranks[n]].temp == bx[ranks[n-1]].temp {
		n--
	}
	return n + 1
}

func summary_gold() {
	clear_temps(T_player)

	for _, i := range loop_kind(T_char) {
		bx[player(i)].temp += has_item(i, item_gold)
	}

	out_ranking("Gold:", 0)
}

func summary_land_owned() {
	clear_temps(T_player)

	for _, i := range loop_subkind(sub_garrison) {
		if player(province_admin(i)) != 0 {
			bx[player(province_admin(i))].temp++
		}
	}

	out_ranking("Land controlled:", 0)
}

func summary_men() {
	clear_temps(T_player)

	for _, i := range loop_kind(T_char) {
		for _, e := range inventory_loop(i) {
			if man_item(e.item) != FALSE {
				bx[player(i)].temp += e.qty
			}
		}
	}

	out_ranking("Men:", 0)
}

func summary_provinces() {
	nlocs := 0
	for _, i := range loop_loc() {
		if loc_depth(i) != LOC_province {
			continue
		}
		nlocs++
	}

	clear_temps(T_player)

	for _, pl := range loop_kind(T_player) {
		clear_temps(T_loc)

		for _, i := range known_sparse_loop(p_player(pl).Known) {
			bx[i].temp++
		}

		for _, i := range loop_loc() {
			if loc_depth(i) != LOC_province {
				continue
			}
			if bx[i].temp != 0 {
				bx[pl].temp++
			}
		}
	}

	out_ranking("Provinces visited:", nlocs)
}

func summary_report() {
	stage("summary_report()")

	collect_game_totals()

	out_path = MASTER
	out_alt_who = OUT_SUMMARY

	for _, pl := range loop_kind(T_player) {
		if subkind(pl) != sub_pl_regular {
			continue
		}

		tagout(pl, "<tag type=game_totals players=%s controlled=%s other=%s>", comma_num(nplayers), comma_num(ncontrolled), comma_num(nother))
		out(pl, "")
		out(pl, "Game totals:")
		indent += 3
		out(pl, "%-20s  %5s", "Players:", comma_num(nplayers))
		out(pl, "%-20s  %5s", "Controlled units:", comma_num(ncontrolled))
		out(pl, "%-20s  %5s", "Other units:", comma_num(nother))
		indent -= 3
		out(pl, "")
		tagout(pl, "</tag type=game_totals>")

		tagout(pl, "<tag type=faction_summary pl=%d>", pl)
		s := sout("Faction %s", box_code(pl))
		if SHOW_TOP {
			out(pl, "%-20s %16s  %-5s  %s", s, "", "rank", "top faction")
		} else {
			out(pl, "%-20s %16s  %-5s  %s", s, "", "rank", "")
		}

		var p string
		for i := 0; i < len(s); i++ {
			p += "-"
		}
		if SHOW_TOP {
			out(pl, "%-20s %16s  %-5s  %s", p, "", "----", "-----------")
		} else {
			out(pl, "%-20s %16s  %-5s  %s", p, "", "----", "")
		}
	}

	summary_units()
	summary_men()
	summary_gold()
	summary_land_owned()
	summary_skills()
	summary_spells()
	summary_provinces()

	for _, pl := range loop_kind(T_player) {
		if subkind(pl) != sub_pl_regular {
			continue
		}
		out(pl, "")
		out(pl, "")
		tagout(pl, "</tag type=faction_summary pl=%d>", pl)
	}

	out_path = 0
	out_alt_who = 0
}

func summary_skills() {
	nskills := 0

	for _, i := range loop_kind(T_skill) {
		if magic_skill(i) {
			continue
		}
		nskills++
	}

	clear_temps(T_player)

	for _, pl := range loop_kind(T_player) {
		if subkind(pl) != sub_pl_regular {
			continue
		}

		clear_temps(T_skill)

		for _, who := range loop_units(pl) {
			for _, e := range loop_char_skill_known(who) {
				bx[e.skill].temp++
			}
		}

		for _, i := range loop_kind(T_skill) {
			if magic_skill(i) {
				continue
			}

			if bx[i].temp != 0 {
				bx[pl].temp++
			}
		}
	}

	out_ranking("Skills known:", nskills)
}

func summary_spells() {
	nskills := 0
	for _, i := range loop_kind(T_skill) {
		if !magic_skill(i) {
			continue
		}
		nskills++
	}

	clear_temps(T_player)

	for _, pl := range loop_kind(T_player) {
		if subkind(pl) != sub_pl_regular {
			continue
		}

		clear_temps(T_skill)

		for _, who := range loop_units(pl) {
			for _, e := range loop_char_skill_known(who) {
				bx[e.skill].temp++
			}
		}

		for _, i := range loop_kind(T_skill) {
			if !magic_skill(i) {
				continue
			}
			if bx[i].temp != 0 {
				bx[pl].temp++
			}
		}
	}

	out_ranking("Spells known:", nskills)
}

func summary_sublocs() {
	nlocs := 0
	for _, i := range loop_loc() {
		if loc_depth(i) != LOC_subloc {
			continue
		}
		nlocs++
	}

	clear_temps(T_player)

	for _, pl := range loop_kind(T_player) {
		clear_temps(T_loc)

		for _, i := range known_sparse_loop(p_player(pl).Known) {
			bx[i].temp++
		}

		for _, i := range loop_loc() {
			if loc_depth(i) != LOC_subloc {
				continue
			}
			if bx[i].temp != 0 {
				bx[pl].temp++
			}
		}
	}

	out_ranking("Sublocs found:", nlocs)
}

func summary_units() {
	clear_temps(T_player)
	for _, i := range loop_kind(T_char) {
		bx[player(i)].temp++
	}
	out_ranking("Characters:", 0)
}

func top_rank() string {
	var s string

	for i, n := 0, 0; i < len(ranks); i, n = i+1, n+1 {
		s = comma_append(s, box_code(ranks[i]))
		if n > 3 {
			return "-"
		} else if i >= len(ranks) {
			break
		} else if bx[ranks[i]].temp != bx[ranks[i-1]].temp {
			break
		}
	}

	return s
}
