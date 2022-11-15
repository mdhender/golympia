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
	"sort"
)

var (
	gold_common_magic = 0
	gold_lead_to_gold = 0
	gold_pot_basket   = 0
	gold_trade        = 0
	gold_fish         = 0
	gold_inn          = 0
	gold_taxes        = 0
	gold_tariffs      = 0
	gold_fees         = 0
	gold_combat       = 0
	gold_combat_indep = 0
	gold_petty_thief  = 0
	gold_temple       = 0
	gold_pillage      = 0
	gold_ferry        = 0
	gold_opium        = 0
	gold_claim        = 0
	gold_times        = 0
)

func skill_use_comp(a, b *int) int {
	return rp_skill(*a).use_count - rp_skill(*b).use_count
}

func gm_show_skill_use_counts(pl int) {
	var sk int
	var i int
	var p *entity_skill

	var l skills_l

	for _, sk = range loop_skill() {
		if skill_school(sk) == sk { /* skip category skills */
			continue
		}

		if rp_skill(sk) != nil && rp_skill(sk).use_count != 0 {
			l = append(l, sk)
		}
	}

	sort.Sort(l)

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "")
	out(pl, "Skill use counts:")
	out(pl, "-----------------")
	out(pl, "")
	out(pl, "%5s  %4s  %5s  %s", "count", "who", "skill", "name")
	out(pl, "%5s  %4s  %5s  %s", "-----", "---", "-----", "----")

	for i = 0; i < len(l); i++ {
		first_use := ' '

		/*  Note if this is the first use ever. */
		if !test_known(gm_player, l[i]) {
			set_known(gm_player, l[i])
			first_use = '*'
		}

		p = rp_skill(l[i])

		out(pl, "%4d%1c  %4s  %5s  %s",
			p.use_count,
			first_use,
			box_code_less(player(p.last_use_who)),
			box_code_less(l[i]),
			just_name(l[i]))
	}

	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

// mdhender: this is a reverse sort?
func skills_known_comp(a, b *int) int {
	return bx[*b].temp - bx[*a].temp
}

func gm_show_skills_known(pl int) {
	var sk int
	var i int
	var e *skill_ent
	var p *entity_skill

	var l skills_l
	clear_temps(T_skill)

	for _, i = range loop_char() {
		for _, e = range loop_char_skill_known(i) {
			bx[e.skill].temp++
		}
	}

	for _, sk = range loop_skill() {
		if skill_school(sk) == sk { /* skip category skills */
			continue
		}

		//#if 0
		//                if (bx[sk].temp)
		//#endif
		l = append(l, sk)
	}

	// reverse sort skills by something?
	l.sort_known_comp()

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "")
	out(pl, "Skills known by players:")
	out(pl, "------------------------")
	out(pl, "")
	out(pl, "%5s  %5s  %4s  %s", "count", "skill", "use", "name")
	out(pl, "%5s  %5s  %4s  %s", "-----", "-----", "---", "----")

	for i = 0; i < len(l); i++ {
		p = rp_skill(l[i])

		if bx[l[i]].temp == 0 {
			break
		}

		out(pl, "%4d   %5s  %4s  %s", bx[l[i]].temp, box_code_less(l[i]), or_string(p.use_count != 0, comma_num(p.use_count), ""), just_name(l[i]))
	}

	out(pl, "")

	out_path = 0
	out_alt_who = 0
	l = nil
}

func gm_count_priests_mages(pl int) {
	priests := 0
	mages := 0
	var i int
	max_aura := 0
	cur_aura := 0
	cur_piety := 0

	for _, i = range loop_char() {
		if is_priest(i) != FALSE {
			priests++
			if rp_char(i).religion.piety > cur_piety {
				cur_piety = rp_char(i).religion.piety
			}
		}
		if is_wizard(i) != FALSE {
			mages++
			if char_max_aura(i) > max_aura {
				max_aura = char_max_aura(i)
			}
			if char_cur_aura(i) > cur_aura {
				cur_aura = char_cur_aura(i)
			}
		}
	}

	out_path = MASTER
	out_alt_who = OUT_LORE
	out(pl, "")
	out(pl, "Mages and Priests")
	out(pl, "-----------------")
	out(pl, "")
	out(pl, " Mages:   %s, Max aura: %s, Cur aura: %s", comma_num(mages), comma_num(max_aura), comma_num(cur_aura))
	out(pl, " Priests: %s, Max piety: %s", comma_num(priests), comma_num(cur_piety))
	out(pl, "")
	out_path = 0
}

func gm_show_interesting_attributes(pl int) {
	out_path = MASTER
	out_alt_who = OUT_LORE

	ability_shroud := 0
	hinder_meditation := 0
	project_cast := 0
	quick_cast := 0
	for _, i := range loop_char() {
		pc := rp_magic(i)
		if pc == nil {
			continue
		}
		if pc.ability_shroud != FALSE {
			ability_shroud++
		}
		if pc.hinder_meditation != FALSE {
			hinder_meditation++
		}
		if pc.project_cast != FALSE {
			project_cast++
		}
		if pc.quick_cast != FALSE {
			quick_cast++
		}
	}

	loc_shroud := 0
	loc_barrier := 0
	loc_opiums := 0
	for _, i := range loop_loc() {
		lc := rp_loc(i)
		if lc != nil && lc.shroud != FALSE {
			loc_shroud++
		}

		if lc != nil && get_effect(i, ef_magic_barrier, 0, 0) != FALSE {
			loc_barrier++
		}

		if loc_opium(i) != FALSE {
			loc_opiums++
		}
	}

	ngal := 0
	rams := 0
	for _, i := range loop_ship() {
		if subkind(i) != sub_galley {
			continue
		}
		ngal++
		if ship_has_ram(i) != FALSE {
			rams++
		}
	}

	nplay := 0
	format_one := 0
	for _, i := range loop_player() {
		if subkind(i) != sub_pl_regular {
			continue
		}
		nplay++
		if player_format(i) != FALSE {
			format_one++
		}
	}

	out(pl, "Interesting attribute counts")
	out(pl, "----------------------------")
	out(pl, "")

	out(pl, "char ability shroud:    %d", ability_shroud)
	out(pl, "char hinder meditate:   %d", hinder_meditation)
	out(pl, "char project cast:      %d", project_cast)
	out(pl, "char quick cast:        %d", quick_cast)
	out(pl, "loc shroud:             %d", loc_shroud)
	out(pl, "loc barrier:            %d", loc_barrier)
	out(pl, "loc opium:              %d", loc_opiums)
	out(pl, "format one:             %d/%d", format_one, nplay)

	out(pl, "")

	if ngal != FALSE {
		out(pl, "galleys with rams:      %d (%d%%)", rams, rams*100/ngal)
	}

	out_path = 0
	out_alt_who = 0
}

func gm_list_animate_items(pl int) {
	var i int

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "")
	out(pl, "Animate items list")
	out(pl, "------------------")
	out(pl, "")

	out(pl, "%25s %8s %8s %8s  %s", "name", "swamp", "man-like", "beast", "fighter")
	out(pl, "%25s %8s %8s %8s  %s", "----", "-----", "--------", "-----", "-------")

	for _, i = range loop_item() {
		if FALSE == is_fighter(i) && FALSE == man_item(i) && !beast_capturable(i) && FALSE == item_animal(i) {
			continue
		}

		var buf string
		if is_fighter(i) != FALSE {
			buf = fmt.Sprintf("(%d,%d,%d)", item_attack(i), item_defense(i), item_missile(i))
		} else {
			buf = " -"
		}

		out(pl, "%25s %8s %8s %8s  %s",
			box_name(i),
			or_string(item_animal(i) != FALSE, "yes ", "no  "),
			or_string(man_item(i) != FALSE, "yes ", "no  "),
			or_string(item_capturable(i) != FALSE, "yes ", "no  "),
			buf)
	}

	out_path = 0
	out_alt_who = 0
}

func gm_show_gold(pl int) {
	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "")
	out(pl, "Gold report")
	out(pl, "-----------")
	out(pl, "")

	sum := gold_common_magic + gold_lead_to_gold + gold_pot_basket + gold_trade + gold_opium + gold_inn +
		gold_taxes + gold_tariffs + gold_combat_indep + gold_petty_thief + gold_combat +
		gold_temple + gold_pillage + gold_ferry + gold_claim + gold_fish + gold_times

	if sum != 0 {
		out(pl, "Common magic:         %10s %3d%%", comma_num(gold_common_magic), gold_common_magic*100/sum)
		out(pl, "Lead to gold:         %10s %3d%%", comma_num(gold_lead_to_gold), gold_lead_to_gold*100/sum)
		out(pl, "Pots and baskets:     %10s %3d%%", comma_num(gold_pot_basket), gold_pot_basket*100/sum)
		out(pl, "Opium:                %10s %3d%%", comma_num(gold_opium), gold_opium*100/sum)
		out(pl, "Fish:                 %10s %3d%%", comma_num(gold_fish), gold_fish*100/sum)
		out(pl, "Trade to cities:      %10s %3d%%", comma_num(gold_trade), gold_trade*100/sum)
		out(pl, "Inn income:           %10s %3d%%", comma_num(gold_inn), gold_inn*100/sum)
		out(pl, "Taxes:                %10s %3d%%", comma_num(gold_taxes), gold_taxes*100/sum)
		out(pl, "Tariffs:              %10s %3d%%", comma_num(gold_tariffs), gold_tariffs*100/sum)
		out(pl, "Fees:                 %10s %3d%%", comma_num(gold_fees), gold_fees*100/sum)
		out(pl, "Combat with indeps:   %10s %3d%%", comma_num(gold_combat_indep), gold_combat_indep*100/sum)
		out(pl, "Combat with players:  %10s %3d%%", comma_num(gold_combat), gold_combat*100/sum)
		out(pl, "Petty thievery:       %10s %3d%%", comma_num(gold_petty_thief), gold_petty_thief*100/sum)
		out(pl, "Temple income:        %10s %3d%%", comma_num(gold_temple), gold_temple*100/sum)
		out(pl, "Pillaging:            %10s %3d%%", comma_num(gold_pillage), gold_pillage*100/sum)
		out(pl, "Ferry boarding:       %10s %3d%%", comma_num(gold_ferry), gold_ferry*100/sum)
		out(pl, "Claim:                %10s %3d%%", comma_num(gold_claim), gold_claim*100/sum)
		out(pl, "Times:                %10s %3d%%", comma_num(gold_times), gold_times*100/sum)
		out(pl, "                      %10s %4s", "", "----")
	}
	out(pl, "Total:                %10s", comma_num(sum))
	out_path = 0
	out_alt_who = 0
}

func gm_show_control_arts(pl int) {
	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "")
	out(pl, "Captured control artifacts")
	out(pl, "--------------------------")
	out(pl, "")

	for _, item := range loop_subkind(sub_magic_artifact) {
		if !is_real_npc(player(item_unique(item))) {
			var c command
			c.who = pl
			c.a = item
			out(pl, "%-33s  %s", box_name(item), box_name(player(item_unique(item))))
			artifact_identify("    ", &c)
		}
	}
	out(pl, "")
	out_path = 0
	out_alt_who = 0
}

func gm_count_stuff(pl int) int {
	castle := 0
	castle_notdone := 0
	tower := 0
	tower_notdone := 0
	mine := 0
	mine_notdone := 0
	temple := 0
	temple_notdone := 0
	galley := 0
	galley_notdone := 0
	round := 0
	round_notdone := 0
	inn := 0
	inn_notdone := 0

	for _, i := range loop_loc_or_ship() {
		switch subkind(i) {
		case sub_castle:
			castle++
			break
		case sub_castle_notdone:
			castle_notdone++
			break
		case sub_tower:
			tower++
			break
		case sub_tower_notdone:
			tower_notdone++
			break
		case sub_galley:
			galley++
			break
		case sub_galley_notdone:
			galley_notdone++
			break
		case sub_roundship:
			round++
			break
		case sub_roundship_notdone:
			round_notdone++
			break
		case sub_temple:
			temple++
			break
		case sub_temple_notdone:
			temple_notdone++
			break
		case sub_inn:
			inn++
			break
		case sub_inn_notdone:
			inn_notdone++
			break
		case sub_mine_shaft, sub_mine:
			mine++
			break
		case sub_mine_shaft_notdone, sub_mine_notdone:
			mine_notdone++
			break
		}
	}

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%10s  %9s  %s", "", "finished", "unfinished")
	out(pl, "%10s +----------------------", "")
	out(pl, "%10s |%9s  %8s", "galley", comma_num(galley), comma_num(galley_notdone))
	out(pl, "%10s |%9s  %8s", "roundship", comma_num(round), comma_num(round_notdone))
	out(pl, "%10s |%9s  %8s", "tower", comma_num(tower), comma_num(tower_notdone))
	out(pl, "%10s |%9s  %8s", "castle", comma_num(castle), comma_num(castle_notdone))
	out(pl, "%10s |%9s  %8s", "mine", comma_num(mine), comma_num(mine_notdone))
	out(pl, "%10s |%9s  %8s", "inn", comma_num(inn), comma_num(inn_notdone))
	out(pl, "%10s |%9s  %8s", "temple", comma_num(temple), comma_num(temple_notdone))
	out(pl, "")

	out_path = 0
	out_alt_who = 0
	return 0
}

func gm_show_gate_stats(pl int) {
	n_gates := 0
	n_found := 0
	ngate_seal := 0
	ngate_jump := 0
	ngate_unseal := 0

	clear_temps(T_gate)

	for _, i := range loop_player() {
		for _, j := range known_sparse_loop(p_player(i).known) {
			if kind(j) != T_gate {
				continue
			}
			bx[j].temp++
		}
	}

	for _, i := range loop_gate() {
		n_gates++
		if bx[i].temp != 0 {
			n_found++
		}

		if g := rp_gate(i); g != nil {
			if g.seal_key != 0 {
				ngate_seal++
			}
			if g.notify_jumps != 0 {
				ngate_jump++
			}
			if g.notify_unseal != 0 {
				ngate_unseal++
			}
		}
	}

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%d/%d gates found (%d%%)", n_found, n_gates, n_found*100/n_gates)
	out(pl, "    %d sealed (%d%%), %d notify jump (%d%%), %d notify unseal (%d%%)",
		ngate_seal, ngate_seal*100/n_gates, ngate_jump, ngate_jump*100/n_gates, ngate_unseal, ngate_unseal*100/n_gates)
	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

func gm_show_locs_visited(pl int) {
	n_prov := 0
	n_prov_visit := 0
	n_sub := 0
	n_sub_visit := 0
	hid := 0
	vis := 0
	nf := 0
	nt := 0
	nf_vis := 0
	nf_hid := 0

	clear_temps(T_loc)

	for _, i := range loop_player() {
		for _, j := range known_sparse_loop(p_player(i).known) {
			if kind(j) != T_loc {
				continue
			}
			bx[j].temp++
		}
	}

	for _, i := range loop_loc() {
		if loc_depth(i) == LOC_province {
			n_prov++
			if bx[i].temp != 0 {
				n_prov_visit++
			}
		} else if loc_depth(i) == LOC_subloc {
			n_sub++
			if bx[i].temp != 0 {
				n_sub_visit++
				if loc_hidden(i) {
					hid++
				} else {
					vis++
				}
			}
		}
	}

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%d/%d provinces visited (%d%%)", n_prov_visit, n_prov, n_prov_visit*100/n_prov)
	out(pl, "%d/%d sublocs visited (%d%%)", n_sub_visit, n_sub, n_sub_visit*100/n_sub)
	out(pl, "    %d%% visible, %d%% hidden", vis*100/n_sub_visit, hid*100/n_sub_visit)

	for _, i := range loop_loc() {
		if loc_depth(i) != LOC_province || bx[i].temp == 0 {
			continue
		}
		for _, j := range loop_here(i) {
			if kind(j) != T_loc || loc_depth(j) != LOC_subloc {
				continue
			}
			nt++
			if bx[j].temp != 0 {
				nf++
				if loc_hidden(j) {
					nf_hid++
				} else {
					nf_vis++
				}
			}
		}
	}

	out(pl, "    %d%% of visisted province's sublocs found", nf*100/nt)
	out(pl, "    %d%% visible, %d%% hidden", nf_vis*100/nf, nf_hid*100/nf)
	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

func gm_loyalty_stats(pl int) {
	tot := 0
	oath := 0
	fear := 0
	cont := 0
	unsw := 0
	npc := 0

	for _, i := range loop_char() {
		if subkind(i) != 0 {
			continue
		}

		tot++
		switch loyal_kind(i) {
		case LOY_oath:
			oath++
		case LOY_fear:
			fear++
		case LOY_contract:
			cont++
		case LOY_unsworn:
			unsw++
		case LOY_npc:
			npc++
		default:
			panic(fmt.Sprintf("assert(loy_kind != %d)", loyal_kind(i)))
		}
	}

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%d chars: %d oath (%d%%), %d fear (%d%%), %d contract (%d%%), %d unsworn (%d%%)",
		tot, oath, oath*100/tot, fear, fear*100/tot, cont, cont*100/tot, unsw, unsw*100/tot)
	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

//static int region_occupy_comp(const void *q1, const void *q2) {
//    int *a = (int *)q1;
//    int *b = (int *)q2;
//    return bx[*b].temp - bx[*a].temp;
//}

func gm_land_stats(pl int) {
	clear_temps(T_loc)

	n_chars := 0
	n_beasts := 0
	for _, i := range loop_char() {
		if region(i) != 0 {
			if player(i) < 1000 {
				n_beasts++
			} else {
				bx[region(i)].temp++
				n_chars++
			}
		}
	}

	var l bxtmp_l
	for _, i := range loop_loc() {
		if bx[i].temp != 0 {
			l = append(l, i)
		}
	}
	sort.Sort(l)

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%10s  %6s  %s", "nobles", "beasts", "region")
	out(pl, "%10s  %6s  %s", "------", "------", "------")
	for i := 0; i < len(l); i++ {
		out(pl, "%10s  %6s  %s", comma_num(bx[l[i]].temp), 0, just_name(l[i]))
	}
	out(pl, "%10s  %6s  %s", "======", "======", "")
	out(pl, "%10s  %6s  %s", comma_num(n_chars), comma_num(n_beasts), "")
	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

//static int wealth_list_comp(const void *q1, const void *q2) {
//    int *a = (int *)q1;
//    int *b = (int *)q2;
//    return bx[*b].temp - bx[*a].temp;
//}

func gm_faction_wealth(pl int) {
	clear_temps(T_player)

	for _, i := range loop_char() {
		bx[player(i)].temp += has_item(i, item_gold)
	}

	var l bxtmp_l
	for _, i := range loop_player() {
		if subkind(i) == sub_pl_regular {
			l = append(l, i)
		}
	}
	sort.Sort(l)

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%4s %11s  %s", "rank", "gold", "faction")
	out(pl, "%4s %11s  %s", "----", "----", "-------")
	for i := 0; i < len(l); i++ {
		out(pl, "%4d %11s  %s", i+1, comma_num(bx[l[i]].temp), box_name(l[i]))
	}
	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

//static int nobles_list_comp(const void *q1, const void *q2) {
//    int *a = (int *)q1;
//    int *b = (int *)q2;
//    return bx[*b].temp - bx[*a].temp;
//}

func gm_nobles_list(pl int) {
	clear_temps(T_player)

	for _, i := range loop_char() {
		bx[player(i)].temp++
	}

	var l bxtmp_l
	for _, i := range loop_player() {
		if subkind(i) == sub_pl_regular {
			l = append(l, i)
		}
	}
	sort.Sort(l)

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%4s %11s  %s", "rank", "nobles", "faction")
	out(pl, "%4s %11s  %s", "----", "------", "-------")
	for i := 0; i < len(l); i++ {
		out(pl, "%4d %11s  %s", i+1, comma_num(bx[l[i]].temp), box_name(l[i]))
	}
	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

func gm_player_details(pl int) {
	var sum_gold, sum_units, sum_bld, sum_subloc, sum_ship, sum_skills int
	var age string

	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%4s %3s %5s %5s %4s %6s %5s %6s", "who", "age", "gold", "units", "bld", "subloc", "ships", "skills")
	out(pl, "%4s %3s %5s %5s %4s %6s %5s %6s", "---", "---", "----", "-----", "---", "------", "-----", "------")

	for _, i := range loop_player() {
		if subkind(i) == sub_pl_system {
			continue
		}

		sum_gold = 0
		for _, j := range loop_units(i) {
			sum_gold += has_item(j, item_gold)
		}

		sum_units = len(loop_units(i))

		sum_bld = 0
		for _, j := range loop_units(i) {
			where := subloc(j)
			if loc_depth(where) == LOC_build && building_owner(where) == j && !is_ship(where) && !is_ship_notdone(where) {
				sum_bld++
			}
		}

		sum_subloc = 0
		for _, j := range loop_units(i) {
			where := subloc(j)
			if loc_depth(where) == LOC_subloc && first_character(where) == j {
				sum_subloc++
			}
		}

		sum_ship = 0
		for _, j := range loop_units(i) {
			where := subloc(j)
			if (is_ship(where) || is_ship_notdone(where)) && building_owner(where) == j {
				sum_ship++
			}
		}

		sum_skills = 0
		clear_temps(T_skill)
		for _, j := range loop_units(i) {
			for _, e := range loop_char_skill_known(j) {
				bx[e.skill].temp++
			}
		}

		for _, j := range loop_skill() {
			if bx[j].temp != 0 {
				sum_skills++
			}
		}

		if rp_player(i) != nil {
			age = sout("%d", sysclock.turn-rp_player(i).first_turn)
		} else {
			age = "???"
		}

		out(pl, "%4s %3s %5s %5s %4s %6s %5s %6s  %s",
			box_code_less(i), age,
			knum(sum_gold, false),
			knum(sum_units, false),
			knum(sum_bld, true),
			knum(sum_subloc, true),
			knum(sum_ship, true),
			knum(sum_skills, true),
			just_name(i))
	}

	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

func list_all_notices(pl int) {
	clear_temps(T_loc)

	for _, i := range loop_post() {
		bx[subloc(i)].temp++
	}

	out_path = MASTER
	out_alt_who = OUT_LORE

	for _, i := range loop_loc() {
		if bx[i].temp == 0 {
			continue
		}
		show_loc_posts(pl, i, TRUE)
	}

	out_path = 0
	out_alt_who = 0
}

func list_all_items(pl int) {
	out_path = MASTER
	out_alt_who = OUT_LORE

	out(pl, "%4s %-24s %3s %6s %4s %4s %4s %4s %3s", "item", "name", "mnt", "weight", "land", "ride", "fly", "cost", "cap")
	out(pl, "%4s %-24s %3s %6s %4s %4s %4s %4s %3s", "----", "----", "---", "------", "----", "----", "---", "----", "---")

	for _, i := range loop_item() {
		var buf string
		if item_attack(i) != 0 {
			buf = fmt.Sprintf("%s (%d,%d,%d)", just_name(i), item_attack(i), item_defense(i), item_missile(i))
		} else {
			buf = fmt.Sprintf("%s", just_name(i))
		}
		/* VLN */
		/*
		   printf("player = %d\n",pl);
		   printf("box_code_less = %s\n",box_code_less(i));
		   printf("buf = %s\n",buf);
		   printf("maint = %d\n",rp_item(i).maintenance);
		   printf("wt = %d\n",item_weight(i));
		   printf("land = %d\n",item_land_cap(i));
		   printf("ride = %d\n",item_ride_cap(i));
		   printf("fly = %d\n",item_fly_cap(i));
		   printf("price = %d\n",rp_item(i).base_price);
		   printf("animal? = %s\n",item_animal(i) ? "yes" : "no");
		*/
		out(pl, "%4s %-24s %3d %4d %4d %4d %4d %4d %3s",
			box_code_less(i), buf, rp_item(i).maintenance,
			item_weight(i), item_land_cap(i), item_ride_cap(i), item_fly_cap(i),
			rp_item(i).base_price, or_string(item_animal(i) != 0, "yes", "no"))
	}
	out(pl, "")

	out_path = 0
	out_alt_who = 0
}

func gm_report(pl int) {
	stage("gm_report()")
	gm_show_gold(pl)
	gm_show_control_arts(pl)
	gm_count_priests_mages(pl)
	gm_show_skill_use_counts(pl)
	gm_count_stuff(pl)
	gm_land_stats(pl)
	gm_show_gate_stats(pl)
	gm_show_locs_visited(pl)
	gm_loyalty_stats(pl)
	gm_show_skills_known(pl)
	gm_faction_wealth(pl)
	gm_nobles_list(pl)
	gm_player_details(pl)
	list_all_notices(pl)
	gm_show_interesting_attributes(pl)
	gm_list_animate_items(pl)

	list_all_items(skill_player)
}
