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
	"path/filepath"
	"sort"
	"strings"
)

/*
 *  Mon Oct 11 07:06:02 1999 -- Scott Turner
 *
 *  This produces a "<tag type=tab col=44>" if tags are on,
 *  a null string elsewise.
 *
 */
func tab_to(where int) string {
	if options.output_tags != FALSE {
		return fmt.Sprintf("<tag type=tab col=%d>", where)
	}
	return ""
}

func output_order_comp(a, b int) int {
	if bx[a].output_order != bx[b].output_order {
		return bx[a].output_order - bx[b].output_order
	}
	return a - b
}

func sort_for_output(l []int) []int {
	//qsort(l, len(l), sizeof(int), output_order_comp);
	var cp []int
	for _, e := range l {
		cp = append(cp, e)
	}
	sort.Ints(cp)
	return cp
}

func determine_output_order() {
	count := 0
	var reg int
	var i int

	for _, reg = range loop_loc() {
		if loc_depth(reg) != LOC_region {
			continue
		}

		bx[reg].output_order = count
		count++

		for _, i = range loop_all_here(reg) {
			bx[i].output_order = count
			count++
		}

	}

	/*
	 *  Sort all player unit lists
	 */

	var pl int
	var p *entity_player

	for _, pl = range loop_player() {
		p = rp_player(pl)
		if p == nil {
			continue
		}

		p.units = sort_for_output(p.units)
	}

}

func show_carry_capacity(who, num int) {
	var w weights
	walk_percent := ""
	var buf string
	//mountains := FALSE;

	out(who, "")

	determine_unit_weights(num, &w, false)

	if w.land_cap > 0 {
		walk_percent = sout(" (%d%%)",
			w.land_weight*100/w.land_cap)
	}
	buf = fmt.Sprintf("%s/%s land%s",
		comma_num(w.land_weight), comma_num(w.land_cap), walk_percent)

	if w.ride_cap > 0 {
		buf += (sout(", %s/%s ride (%d%%)",
			comma_num(w.ride_weight), comma_num(w.ride_cap),
			w.ride_weight*100/w.ride_cap))
	}

	if w.fly_cap > 0 {
		buf += (sout(", %s/%s fly (%d%%)",
			comma_num(w.fly_weight), comma_num(w.fly_cap),
			w.fly_weight*100/w.fly_cap))
	}

	tagout(who, "<tag type=carry_capacity unit=%d weight=%d walk_weight=%d walk_cap=%d ride_weight=%d ride_cap=%d fly_weight=%d fly_cap=%d>",
		who,
		w.total_weight,
		w.land_weight,
		w.land_cap,
		w.ride_weight,
		w.ride_cap,
		w.fly_weight,
		w.fly_cap)

	wiout(who, len("Capacity:  "), "Capacity:  %s", buf)
	tagout(who, "</tag type=capacity>")
}

func show_item_skills_sup(who, item int, p *item_magic) {
	var i int
	var req_s string
	var sk int
	var parent int
	//#if 0
	//    int see_magic;
	//
	//    see_magic = is_magician(who);
	//#endif

	first := true
	for i = 0; i < len(p.may_study); i++ {
		sk = p.may_study[i]
		assert(valid_box(sk))
		parent = skill_school(sk)

		//#if 0
		//        if (magic_skill(sk) && !see_magic)
		//            continue;
		//#endif

		if sk == parent {
			req_s = ""
		} else {
			req_s = sout(" (requires %s)",
				just_name(parent))
		}

		//#if 0
		//        /*
		//         *  Output the days study remaining in the book.
		//         *
		//         */
		//        days = sout(", %d day%s instruction",
		//                p.orb_use_count,
		//                or_string(p.orb_use_count == 1 , "" , "s"));
		//#endif

		tagout(who, "<tag type=may_study id=%d skill=%d parent = %d days=%d>",
			item,
			or_int((sk != parent && has_skill(who, parent) < 1), 0, sk),
			parent,
			p.orb_use_count)

		if first {
			out(who, "")
			wout(who, "%s permits %d day%s study of the following skill%s:",
				box_name(item),
				p.orb_use_count,
				or_string(p.orb_use_count == 1, " of ", "s"),
				or_string(len(p.may_study) == 1, "", "s"))
			indent += 3
			first = false
		}

		if sk != parent && has_skill(who, parent) < 1 {
			wiout(who, 3, "???%s", req_s)
		} else {
			wiout(who, 3, "%s%s", box_name(sk), req_s)
		}

		tagout(who, "</tag type=may_study>")
	}

	if !first {
		indent -= 3
	}

	first = true
	for i = 0; i < len(p.may_use); i++ {
		sk = p.may_use[i]
		assert(valid_box(sk))
		parent = skill_school(sk)

		//#if 0
		//        if (magic_skill(sk) && !see_magic)
		//            continue;
		//#endif

		if first {
			out(who, "")
			wout(who, "%s permits use of the following skills:",
				box_name(item))
			indent += 3
			first = false
		}

		if sk == parent {
			req_s = ""
		} else {
			req_s = sout(" (requires %s)",
				just_name(parent))
		}

		tagout(who, "<tag type=may_use id=%d skill=%d parent = %d>",
			item,
			or_int((sk != parent && has_skill(who, parent) < 1), 0, sk),
			parent)

		if sk != parent && has_skill(who, parent) < 1 {
			wiout(who, 3, "???%s", req_s)
		} else {
			wiout(who, 3, "%s%s", box_name(sk), req_s)
		}

		tagout(who, "</tag type=may_use>")
	}

	if !first {
		indent -= 3
	}
}

func show_item_skills(who, num int) {
	var e *item_ent
	var p *item_magic
	first := true

	for _, e = range loop_inventory(num) {
		p = rp_item_magic(e.item)

		if p != nil {
			if first {
				tagout(who, "<tag type=item_skill_section id=%d>", who)
				first = false
			}
			show_item_skills_sup(who, e.item, p)
		}
	}

	if !first {
		tagout(who, "</tag type=item_skill_section>")
	}
}

func inv_item_comp(a, b *item_ent) int {
	return a.item - b.item
}

func extra_item_info(who, item, qty int) string {
	var buf string
	var lc, rc, fc int
	var at, df, mi int

	buf = ""

	lc = item_land_cap(item)
	rc = item_ride_cap(item)
	fc = item_fly_cap(item)

	if fc > 0 {
		buf = fmt.Sprintf("fly %s", nice_num(fc*qty))
	} else if rc > 0 {
		buf = fmt.Sprintf("ride %s", nice_num(rc*qty))
	} else if lc > 0 {
		buf = fmt.Sprintf("cap %s", nice_num(lc*qty))
	}

	at = item_attack(item)
	df = item_defense(item)
	mi = item_missile(item)

	if is_fighter(item) != FALSE {
		buf += (sout(" (%d,%d,%d)", at, df, mi))
	}

	if n := item_attack_bonus(item); n != FALSE {
		buf += (sout("+%d attack", n))
	}

	if n := item_defense_bonus(item); n != FALSE {
		buf += (sout("+%d defense", n))
	}

	if n := item_missile_bonus(item); n != FALSE {
		buf += (sout("+%d missile", n))
	}

	if n := item_aura_bonus(item); n != FALSE {
		if who != FALSE && is_magician(who) {
			buf += (sout("+%d aura", n))
		}
	}

	return sout("%s", buf)
}

func show_char_inventory(who, num int, prefix string) {
	first := true
	var e *item_ent
	var weight int
	count := 0
	total_weight := 0

	if len(bx[num].items) > 0 {
		//qsort(bx[num].items, len(bx[num].items), sizeof(int), inv_item_comp);
		sort.Slice(bx[num].items, func(i, j int) bool {
			return inv_item_comp(bx[num].items[i], bx[num].items[j]) < 0
		})
	}

	if len(prefix) == 0 {
		tagout(who, "<tag type=inventory_section unit=%d>", who)
	}

	for _, e = range loop_inventory(num) {
		weight = item_weight(e.item) * e.qty

		if first {
			if len(prefix) == 0 {
				out(who, "")
			}
			out(who, "%sInventory:", prefix)
			out(who, "%s%9s  %-30s %9s", prefix,
				"qty", "name", "weight")
			out(who, "%s%9s  %-30s %9s", prefix,
				"---", "----", "------")
			first = false
		}

		if len(prefix) == 0 {
			tagout(who, "<tag type=inventory unit=%d item=%d qty=%d weight=%d extra=\"%s\">",
				who, e.item, e.qty, item_weight(e.item)*e.qty,
				extra_item_info(who, e.item, e.qty))
		}

		out(who, "%s%9s  %-30s %s%9s  %s",
			prefix,
			comma_num(e.qty),
			plural_item_box(e.item, e.qty),
			tab_to(45),
			comma_num(weight),
			extra_item_info(who, e.item, e.qty))

		if len(prefix) == 0 {
			tagout(who, "</tag type=inventory>")
		}

		count++
		total_weight += weight
	}

	if count > 0 {
		out(who, "%s%9s  %-30s %9s", prefix,
			"", "", "======")
		out(who, "%s%9s  %-30s %9s", prefix,
			"", "", comma_num(total_weight))
	}

	if first {
		if len(prefix) == 0 {
			out(who, "")
		}
		out(who, "%s has no possessions.", box_name(num))
	}

	if len(prefix) == 0 {
		tagout(who, "</tag type=inventory_section unit=%d>", who)
	}
}

/*
 *	1.  building		%s, in
 *
 *	2.  land subloc		%s, in province
 *	3.  ocean subloc	%s, in
 *
 *	4.  land province	%s, in region %s
 *	5.  ocean province	%s, in %s
 */

func char_rep_location(who int) string {
	where := subloc(who)
	var s string
	var reg_name string

	if where == 0 {
		return "nowhere"
	}

	for where != FALSE && loc_depth(where) > LOC_province {
		if len(s) != 0 {
			s = sout("%s, in %s", s, box_name(where))
		} else {
			s = box_name(where)
		}
		where = loc(where)
	}

	if where == FALSE {
		s = sout("%s, adrift in the Cosmos", s)
		return s
	}

	if subkind(province(where)) == sub_ocean {
		if len(s) != 0 {
			s = sout("%s, in %s", s, box_name(where))
		} else {
			s = box_name(where)
		}

		reg_name = name(region(where))

		if len(reg_name) != 0 {
			s = sout("%s, in %s", s, reg_name)
		}
	} else {
		if len(s) != 0 {
			s = sout("%s, in province %s", s, box_name(where))
		} else {
			s = box_name(where)
		}

		reg_name = name(region(where))

		//#if 0
		//        if (reg_name && *reg_name)
		//            s = sout("%s, in region %s", s, reg_name);
		//#else
		if len(reg_name) != 0 {
			s = sout("%s, in %s", s, reg_name)
		}
		//#endif
	}

	return s
}

func char_rep_stack_info(who, num int) {
	var n int
	first := true

	if n = stack_parent(num); n != FALSE {
		tagout(who, "<tag type=stack id=%d under=%d>", num, n)
		wiout(who, 16, "Stacked under:  %s", box_name(n))
	}

	for _, n = range loop_here(num) {
		if kind(n) == T_char && !is_prisoner(n) {
			if first {
				out(who, "Stacked over:   %s", box_name(n))
				first = first
			} else {
				out(who, "                %s", box_name(n))
			}
			tagout(who, "<tag type=stack id=%d over=%d>", num, n)
		}
	}

}

//#if 0
//static int pledge_backlinks = FALSE;
//
//static void
//collect_pledge_backlinks()
//{
//    var i int
//    var n int
//    var p *char_magic
//
//    pledge_backlinks = TRUE;
//
//    loop_char(i)
//    {
//        if (n = char_pledge(i))
//        {
//            p = p_magic(n);
//            ilist_append(&p.pledged_to_us, i);
//        }
//    }
//    next_char;
//}
//
//
//static void
//show_pledged(who, num int)
//{
//    int i, n;
//    first := TRUE;
//    var p *char_magic
//
//    if (!pledge_backlinks)
//        collect_pledge_backlinks();
//
//    if (n = char_pledge(num))
//        wiout(who, 16,  "Pledged to:     %s", box_name(n));
//
//    p = rp_magic(num);
//
//    if (p)
//    {
//        for i = 0; i < len(p.pledged_to_us); i++
//        {
//            n = p.pledged_to_us[i];
//
//            if (first)
//            {
//                out(who, "Pledged to us:  %s", box_name(n));
//                first = FALSE;
//            }
//            else
//                out(who, "                %s", box_name(n));
//        }
//    }
//}
//#endif

func prisoner_health(who int) string {
	health := char_health(who)

	assert(health != 0)

	if health < 0 {
		return ""
	}

	return sout(", health %d", health)
}

func char_rep_prisoners(who, num int) {
	var n int
	first := true

	for _, n = range loop_here(num) {
		if kind(n) == T_char && is_prisoner(n) {
			if first {
				out(who, "Prisoners:      %s%s%s%s",
					box_name(n),
					nation_s(n),
					deserted_s(n),
					prisoner_health(n))
				first = false
			} else {
				out(who, "                %s%s%s%s",
					box_name(n),
					nation_s(n),
					deserted_s(n),
					prisoner_health(n))
			}
			tagout(who, "<tag type=prisoner id=%d prisoner=%d nation=\"%s\" deserted=\"%s\" health=%d> ",
				num,
				n,
				nation_s(n),
				deserted_s(n),
				char_health(n))

		}
	}

}

func char_rep_health(who, num int, prefix string) {
	var n int
	var s string

	n = char_health(num)

	if n == -1 {
		out(who, "%sHealth:         n/a", prefix)
		return
	}

	if n > 0 && n < 100 {
		if char_sick(num) != FALSE {
			s = " (getting worse)"
		} else {
			s = " (getting better)"
		}
	}

	if len(prefix) == 0 {
		tagout(who, "<tag type=health id=%d health=%d sick=%d>",
			num, n, char_sick(num))
	}
	out(who, "%sHealth:         %d%%%s", prefix, n, s)
}

func char_rep_combat(who, num int) {
	var n int
	var s string
	var mk, att, def, mis int

	mk = noble_item(num)
	if mk == 0 {
		att = char_attack(num)
		def = char_defense(num)
		mis = char_missile(num)
	} else {
		att = item_attack(mk)
		def = item_defense(mk)
		mis = item_missile(mk)
	}

	n = char_behind(num)

	if n == 0 {
		s = " (front line in combat)"
	} else {
		s = " (stay behind in combat)"
	}

	if char_break(num) == 0 {
		s = " (fight to the death)"
	} else if char_break(num) == 100 {
		s = " (break almost immediately)"
	} else {
		s = ""
	}

	tagout(who, "<tag type=combat id=%d attack=%d defense=%d missile=%d behind=%d break=%d personal_break=%d>",
		num, att, def, mis, n, char_break(num), personal_break(num))
	out(who, "Combat:         attack %d, defense %d, missile %d",
		att, def, mis)
	out(who, "                behind %d %s", n, s)
	out(who, "Break point:    %d%%%s", char_break(num), s)

	if has_skill(num, sk_personal_fttd) != FALSE {
		out(who, "Personal break point:	When health reaches %d.",
			personal_break(num))
	}

}

func char_rep_misc(who, num int) {
	var s string
	var p *char_magic

	p = rp_magic(num)
	if p != nil && p.ability_shroud != FALSE {
		out(who, "Ability shroud: %d aura", p.ability_shroud)
	}

	if has_skill(num, sk_hide_self) != FALSE {
		if char_hidden(num) != FALSE {
			if char_alone_stealth(num) {
				s = "concealing self"
			} else {
				s = "concealing self, but not alone"
			}
		} else {
			s = "not concealing self"
		}

		out(who, "use %4d %d      (%s)",
			sk_hide_self, char_hidden(num), s)
	}
	/*
	 *  Smuggling.
	 *
	 */
	if has_skill(num, sk_smuggle_goods) != FALSE ||
		has_skill(num, sk_smuggle_men) != FALSE {
		s = " nothing"
		if get_effect(num, ef_smuggle_goods, 0, 0) != FALSE {
			if get_effect(num, ef_smuggle_men, 0, 0) != FALSE {
				s = " men and goods"
			} else {
				s = " goods"
			}
		} else if get_effect(num, ef_smuggle_men, 0, 0) != FALSE {
			s = " men"
		}
		out(who, "Smuggling:      %s.", s)
	}
}

func char_rep_magic(who, num int, prefix string) {
	var ca, ma, mea int

	ca = char_cur_aura(num)
	ma = char_max_aura(num)
	mea = max_eff_aura(num)

	if len(prefix) == 0 {
		out(who, "")
	}
	out(who, "%sCurrent aura:   %d", prefix, ca)

	if ma < mea {
		out(who, "%sMaximum aura:   %d (%d+%d)",
			prefix, mea, ma, mea-ma)
	} else {
		out(who, "%sMaximum aura:   %d", prefix, ma)
	}

	if char_abil_shroud(num) > 0 {
		out(who, "%sAbility shroud: %d aura",
			prefix, char_abil_shroud(num))
	}

	if is_loc_or_ship(char_proj_cast(num)) {
		out(who, "%sProject cast:   %s",
			prefix,
			box_name(char_proj_cast(num)))
	}

	if char_quick_cast(num) != FALSE {
		out(who, "%sQuicken cast:   %d",
			prefix,
			char_quick_cast(num))
	}

	if len(prefix) == 0 {
		tagout(who, "<tag type=magic id=%d cur_aura=%d max_aura=%d max_eff_aura=%d abil_shroud=%d project=%d quick=%d>",
			num, ca, ma, mea, char_abil_shroud(num),
			char_proj_cast(num), char_quick_cast(num))
	}

}

func char_rep_religion(who, num int) {
	var i int

	if is_priest(num) != FALSE {
		out(who, "Current piety:  %d", rp_char(num).religion.piety)
		if len(rp_char(num).religion.followers) > 0 {
			out(who, "Followers:  %s", box_name(rp_char(num).religion.followers[0]))
			for i = 1; i < len(rp_char(num).religion.followers); i++ {
				out(who, "            %s", box_name(rp_char(num).religion.followers[i]))
			}
		}
	} else if is_follower(num) != FALSE {
		out(who, "Dedicated to:   %s (%s)",
			box_name(is_priest(is_follower(num))),
			box_name(is_follower(num)))
	}
	tagout(who, "<tag type=religion id=%d followers=%d priest=%d>",
		num, or_int(is_priest(num) != FALSE, len(rp_char(num).religion.followers), 0),
		is_follower(num))
}

func char_rep_sup(who, num int) {
	tagout(who, "<tag type=char_report id=%d location=%d loc_string=\"%s\" loyalty=%s loy_type=%d loy_num=%d guild=%d>", who, subloc(num),
		char_rep_location(num), cap_(loyal_s(num)),
		loyal_kind(num), loyal_rate(num), guild_member(who))

	wiout(who, 16, "Location:       %s", char_rep_location(num))
	out(who, "Loyalty:        %s", cap_(loyal_s(num)))
	if guild_member(who) != FALSE {
		out(who, "Guild:          %s", box_name(guild_member(who)))
	}

	char_rep_stack_info(who, num)
	char_rep_health(who, num, "")
	char_rep_combat(who, num)
	char_rep_misc(who, num)

	if tmp := banner(num); len(tmp) != 0 {
		tagout(who, "<tag type=banner id=%d value=\"%s\">", num, tmp)
		out(who, "Banner:         %s", tmp)
	}
	//#if 0
	//    show_pledged(who, num);
	//#endif

	char_rep_religion(who, num)

	if is_magician(num) {
		char_rep_magic(who, num, "")
	}
	char_rep_prisoners(who, num)

	print_att(who, num)
	list_skills(who, num, "")
	list_partial_skills(who, num, "")
	list_accepts(who, num)
	show_char_inventory(who, num, "")
	show_carry_capacity(who, num)
	show_item_skills(who, num)
	list_pending_trades(who, num)

	out(who, "")
	tagout(who, "</tag type=char_report id=%d>", who)
}

func character_report() {
	var who int

	stage("character_report()")

	indent += 3

	for _, who = range loop_char() {
		if subkind(player(who)) == sub_pl_silent {
			continue
		}

		tagout(who, "</tag type=unit_report>")

		if is_prisoner(who) {
			p_char(who).prisoner = FALSE /* turn output on */
			tagout(who, "<tag type=char_report id=%d>", who)
			out(who, "%s is being held prisoner.", box_name(who))
			tagout(who, "<tag type=header>")
			out(who, "")
			tagout(who, "</tag type=header>")
			tagout(who, "</tag type=char_report id=%d>", who)
			p_char(who).prisoner = TRUE /* output off again */
		} else {
			out(who, "")
			out(who, "%s", box_name(who))
			out(who, "")
			char_rep_sup(who, who)
		}
	}

	indent -= 3
}

func show_unclaimed(who, num int) {
	first := true
	var e *item_ent
	var weight int

	if len(bx[num].items) > 0 {
		//qsort(bx[num].items, len(bx[num].items), sizeof(int), inv_item_comp);
		sort.Slice(bx[num].items, func(i, j int) bool {
			return inv_item_comp(bx[num].items[i], bx[num].items[j]) < 0
		})
	}

	for _, e = range loop_inventory(num) {
		weight = item_weight(e.item) * e.qty

		if first {
			out(who, "")
			out(who, "Unclaimed items:")
			out(who, "")
			out(who, "%9s  %-30s %9s", "qty", "name", "weight")
			out(who, "%9s  %-30s %9s", "---", "----", "------")
			first = false
		}

		out(who, "%9s  %-30s %s%9s  %s",
			comma_num(e.qty),
			plural_item_box(e.item, e.qty),
			tab_to(42),
			comma_num(weight),
			extra_item_info(0, e.item, e.qty))
	}

	if rp_player(who).first_tower == FALSE {
		if first {
			out(who, "")
			out(who, "Unclaimed items:")
			out(who, "")
		}
		out(who, "")
		out(who, "   You have not yet built your 'free' tower.")
	}
}

func player_ent_info() {
	var pl int

	for _, pl = range loop_player() {
		if subkind(pl) == sub_pl_silent {
			continue
		}

		tagout(pl, "</tag type=unit_report>")
		print_admit(pl)
		print_att(pl, pl)
		list_accepts(pl, pl)
		show_unclaimed(pl, pl)
	}
}

func sum_fighters(who int) int {
	sum := 0
	var t *item_ent

	for _, t = range loop_inventory(who) {
		if man_item(t.item) != FALSE && is_fighter(t.item) != FALSE {
			val := max(item_attack(t.item), max(item_defense(t.item), item_missile(t.item)))
			if val > 1 {
				sum += t.qty
			}
		}
	}

	return sum
}

const TRUNC_NAME = 15

var stupid_words = []string{"a", "the", "of", "de", "des", "la", "and", "du", "aux", "et", "ses", "avec", "un", "van", "von", "-", "--", ""}

func strip_leading_stupid_word(s string) string {
	for i := 0; stupid_words[i] != ""; i++ {
		if word := stupid_words[i]; strings.HasPrefix(strings.ToLower(s), word) {
			if t := strings.TrimSpace(s[len(word):]); len(t) != 0 {
				return t
			}
			break
		}
	}
	return s
}

func stupid_word(s string) bool {
	return lookup_ss(stupid_words, s) >= 0
}

func prev_word(s, t string) string {
	//for t > s && *t != ' ' {
	//    t--;
	//}
	//if (t > s) {
	//    return t;
	//}
	//return "";
	panic("!implemented")
}

func summary_trunc_name(who int) string {
	//
	//var s, t string
	//
	//s = sout("%s", just_name(who));
	//
	//if (len(s) <= TRUNC_NAME) {
	//    return s;
	//}
	//
	//s = strip_leading_stupid_word(s);
	//
	//if (len(s) <= TRUNC_NAME) {
	//    return s;
	//}
	//
	//t = prev_word(s, &s[TRUNC_NAME]);
	//if (t) {
	//    *t = '\0';
	//}
	//
	//for t = prev_word(s, t); t != nil && *t == ' ' && stupid_word(t + 1); t = prev_words(s, t) {
	//    *t = '\0';
	//}
	//
	//s[TRUNC_NAME] = '\0';        /* catches a case */
	//
	//return s;

	panic("!implemented")
}

var sum_gold int
var sum_peas int
var sum_work int
var sum_sail int
var sum_fight int

var loyal_chars = "ucofns"

func unit_summary_sup(pl, who int) {
	var nam string
	var health_s string
	var under_s string
	var loy_s string
	var cur_aura_s string
	pr := is_prisoner(who)
	var gold, peas, work, sail, fight int
	var buf string
	var n int

	nam = summary_trunc_name(who)

	n = char_health(who)

	if n == 100 {
		health_s = "100 "
	} else if n == -1 {
		health_s = "n/a "
	} else if char_sick(who) != FALSE {
		health_s = sout("%d-", char_health(who))
	} else {
		health_s = sout("%d+", char_health(who))
	}

	if pr {
		under_s = " ?? "
	} else if stack_parent(who) != FALSE {
		under_s = box_code_less(stack_parent(who))
	} else {
		under_s = ""
	}

	if is_magician(who) {
		cur_aura_s = sout("%d", char_cur_aura(who))
	} else if is_priest(who) != FALSE {
		cur_aura_s = sout("%d", rp_char(who).religion.piety)
	} else {
		cur_aura_s = ""
	}

	gold = has_item(who, item_gold)
	peas = has_item(who, item_peasant)
	work = has_item(who, item_worker)
	sail = has_item(who, item_sailor)
	fight = sum_fighters(who)

	loy_s = sout("%c%s", loyal_chars[loyal_kind(who)],
		knum(loyal_rate(who), false))

	if options.output_tags > 0 {
		buf = fmt.Sprintf(
			"<tag type=unit_summary unit=%d loc=%d loyal=%s health=%s behind=%d aura=%d gold=%d peas=%d work=%d sail=%d fight=%d under=%d name=\"%s\">",
			who,
			or_int(pr, -1, subloc(who)),
			loy_s, health_s,
			char_behind(who), char_cur_aura(who),
			has_item(who, item_gold),
			has_item(who, item_peasant),
			has_item(who, item_worker),
			has_item(who, item_sailor),
			fight,
			or_int(pr, -1, stack_parent(who)),
			nam)
		tagout(pl, buf)
	}

	buf = fmt.Sprintf("%-*s %-*s %-5s %4s%2d%5s %4s %4s %4s %4s %4s %-*s %s",
		CHAR_FIELD, box_code_less(who),
		CHAR_FIELD, or_string(pr, " ?? ", box_code_less(subloc(who))),
		loy_s,
		health_s,
		char_behind(who),
		cur_aura_s,
		knum(gold, true),
		knum(peas, true),
		knum(work, true),
		knum(sail, true),
		knum(fight, true),
		CHAR_FIELD, under_s,
		nam)

	out(pl, "%s", buf)

	tagout(pl, "</tag type=unit_summary>")

	sum_gold += gold
	sum_peas += peas
	sum_work += work
	sum_sail += sail
	sum_fight += fight
}

func unit_summary(pl int) {
	var i int
	count := 0

	clear_temps(T_char)

	sum_gold = 0
	sum_peas = 0
	sum_work = 0
	sum_sail = 0
	sum_fight = 0

	count = len(p_player(pl).units)

	if count <= 0 {
		return
	}

	tagout(pl, "<tag type=unit_section pl=%d>", pl)
	out(pl, "")
	out(pl, "Unit Summary:")
	out(pl, "")
	out(pl, "%-*s %-*s loyal heal B  CA  gold peas work sail figh %-*s name",
		CHAR_FIELD, "unit",
		CHAR_FIELD, "where",
		CHAR_FIELD, "under")
	out(pl, "%-*s %-*s ----- ---- - ---- ---- ---- ---- ---- ---- %-*s ----",
		CHAR_FIELD, "----",
		CHAR_FIELD, "-----",
		CHAR_FIELD, "-----")

	for _, i = range loop_units(pl) {
		unit_summary_sup(pl, i)
	}

	if count > 1 {
		out(pl, "%*s %-*s                   ==== ==== ==== ==== ====",
			CHAR_FIELD, "",
			CHAR_FIELD, "")
		out(pl, "%*s %-*s                   %4s %4s %4s %4s %4s",
			CHAR_FIELD, "",
			CHAR_FIELD, "",
			knum(sum_gold, true),
			knum(sum_peas, true),
			knum(sum_work, true),
			knum(sum_sail, true),
			knum(sum_fight, true))
	}
	tagout(pl, "</tag type=unit_section>")
}

func loc_ind_s(where int) string {
	var ld int

	ld = loc_depth(where) - 1

	if ld <= 0 {
		return just_name(where)
	}

	return sout("%s%s", &spaces[spaces_len-(ld*2)], box_name(where))
}

func loc_stack_catchup(pl, where int) {

	if where == 0 || bx[where].temp != FALSE {
		return
	}

	loc_stack_catchup(pl, loc(where))
	out(pl, "%s", loc_ind_s(where))
	bx[where].temp = -1
}

var loc_stack_explain int

func loc_stack_rep_sup(pl, where, who int) {
	var where_s, star, ind string

	if where != FALSE {
		loc_stack_catchup(pl, loc(where))
		where_s = loc_ind_s(where)
	}

	if player(who) != pl {
		star = " *"
		loc_stack_explain = TRUE
	}

	if kind(loc(who)) == T_char {
		ind = "  "
	}

	out(pl, "%-34s %s%s%s%s", where_s, tab_to(35),
		ind, box_name(who), star)
	if where != FALSE {
		bx[loc(where)].temp = 0
	}
}

func loc_stack_report(pl int) {
	var i, j int
	var l []int
	var locs []int

	clear_temps(T_loc)
	clear_temps(T_ship)
	clear_temps(T_char)

	loc_stack_explain = FALSE

	tagout(pl, "<tag type=loc_summary pl=%d>", pl)

	for _, i = range loop_units(pl) {
		if is_prisoner(i) {
			continue
		}

		l = append(l, i)
	}

	if len(l) < 1 {
		return
	}

	tagout(pl, "<tag type=header>")
	out(pl, "")
	out(pl, "Stack Locations:")
	out(pl, "")
	out(pl, "%-34s %s", "Location", "Stack")
	out(pl, "%-34s %s", "--------", "-----")
	tagout(pl, "</tag type=header>")

	sort_for_output(l)

	for i = len(l) - 1; i >= 0; i-- {
		where := subloc(l[i])

		if bx[where].temp == 0 {
			locs = append(locs, where)
		}

		bx[l[i]].temp = bx[where].temp
		bx[where].temp = l[i]
	}

	sort_for_output(locs)

	for i = 0; i < len(locs); i++ {
		j = bx[locs[i]].temp
		assert(valid_box(j))

		loc_stack_rep_sup(pl, locs[i], j)

		for j = bx[j].temp; j != FALSE; j = bx[j].temp {
			assert(valid_box(j))
			loc_stack_rep_sup(pl, 0, j)
		}
	}

	if loc_stack_explain != FALSE {
		out(pl, "")
		out(pl, "%-34s    %s%s", "",
			tab_to(38),
			"* -- unit belongs to another faction")
	}

	tagout(pl, "</tag type=loc_summary>", pl)
}

func player_report_sup(pl int) {
	var p *entity_player

	if subkind(pl) == sub_pl_system {
		return
	}

	p = p_player(pl)

	tagout(pl, "<tag type=player_summary pl=%d np=%d sp=%d>",
		pl, p.noble_points, p.jump_start)
	out(pl, "Noble points:  %d     (%d gained, %d spent)",
		p.noble_points,
		p.np_gained,
		p.np_spent)
	out(pl, "Study points:  %d",
		p.jump_start)

	print_hiring_status(pl)
	print_unformed(pl)

	tagout(pl, "</tag type=player_summary>")
}

func stack_capacity_report(pl int) {
	var w weights
	var who int
	var walk_s, ride_s, fly_s string
	var s string
	first := true
	var n int

	for _, who = range loop_units(pl) {
		if first {
			tagout(pl, "<tag type=capacity_summary pl=%d>", pl)
			tagout(pl, "<tag type=header >", pl)
			out(pl, "")
			out(pl, "Stack Capacities:")
			out(pl, "  - First number is additional weight you can walk, ride or fly with.")
			out(pl, "  - If you're overloaded, it shows the excess in parentheses.")
			out(pl, "  - Second number is the percentage of your total capacity used.")
			out(pl, "")
			out(pl, "%*s  %10s %15s %15s %15s",
				CHAR_FIELD, "stack",
				"total wt",
				"   walk   ",
				"   ride   ",
				"   fly    ")
			out(pl, "%*s  %10s %15s %15s %15s",
				CHAR_FIELD, "-----",
				"--------",
				"-----------",
				"-----------",
				"-----------")
			tagout(pl, "</tag type=header >", pl)
			first = false
		}

		determine_stack_weights(who, &w, false)

		if w.land_cap > 0 {
			n = w.land_weight * 100 / w.land_cap

			if n > 999 {
				s = " -- "
			} else {
				s = sout("%3d%%", n)
			}

			if w.land_weight > w.land_cap {
				walk_s = sout("(%s) %s",
					comma_num(w.land_weight-w.land_cap),
					s)
			} else {
				walk_s = sout("%s %s",
					comma_num(w.land_cap-w.land_weight),
					s)
			}
		} else {
			walk_s = ""
		}

		if w.ride_cap > 0 {
			n = w.ride_weight * 100 / w.ride_cap

			if n > 999 {
				s = " -- "
			} else {
				s = sout("%3d%%", n)
			}

			if w.ride_weight > w.ride_cap {
				ride_s = sout("(%s) %s",
					comma_num(w.ride_weight-w.ride_cap),
					s)
			} else {
				ride_s = sout("%s %s",
					comma_num(w.ride_cap-w.ride_weight),
					s)
			}
		} else {
			ride_s = ""
		}

		if w.fly_cap > 0 {
			n = w.fly_weight * 100 / w.fly_cap

			if n > 999 {
				s = " -- "
			} else {
				s = sout("%3d%%", n)
			}

			if w.fly_weight > w.fly_cap {
				fly_s = sout("(%s) %s",
					comma_num(w.fly_weight-w.fly_cap),
					s)
			} else {
				fly_s = sout("%s %s",
					comma_num(w.fly_cap-w.fly_weight),
					s)
			}
		} else {
			fly_s = ""
		}

		tagout(pl, "<tag type=capacity_unit unit=%d weight=%d walk_weight=%d walk_cap=%d ride_weight=%d ride_cap=%d fly_weight=%d fly_cap=%d>",
			who,
			w.total_weight,
			w.land_weight,
			w.land_cap,
			w.ride_weight,
			w.ride_cap,
			w.fly_weight,
			w.fly_cap)

		out(pl, "%*s  %10s %15s %15s %15s",
			CHAR_FIELD, box_code_less(who),
			comma_num(w.total_weight),
			walk_s,
			ride_s,
			fly_s)

	}

	if !first {
		tagout(pl, "</tag type=capacity_summary>", pl)
	}
}

func player_report() {
	var pl int

	stage("player_report()")

	out_path = MASTER
	out_alt_who = OUT_BANNER

	for _, pl = range loop_player() {
		if subkind(pl) == sub_pl_silent {
			continue
		}

		player_report_sup(pl)
		unit_summary(pl)
		loc_stack_report(pl)
		stack_capacity_report(pl)
		storm_report(pl)
		ship_summary(pl)
		garrison_summary(pl)
		out(pl, "")
	}

	out_path = 0
	out_alt_who = 0
}

func rep_player(pl int) {
	var s string

	tags_off()
	s = box_name(pl)
	tags_on()

	tagout(pl, "<tag type=unit_report id=%d name=\"%s\">", pl,
		box_name(pl))
	tagout(pl, "<tag type=header>")
	lines(pl, s)
	tagout(pl, "</tag type=header>")

	out(pl, "#include %d", pl)
	/*
	    *  Don't do this; it has been done in player_ent_report.
	    *
	   tagout(pl,"</tag type=unit_report>");
	*/
	tagout(pl, "<tag type=header>")
	out(pl, "")
	tagout(pl, "</tag type=header>")
}

func rep_char(pl int, l []int) {
	var i int
	var s, t string

	sort_for_output(l)

	for i = 0; i < len(l); i++ {
		if subkind(l[i]) == sub_dead_body ||
			subkind(l[i]) == sub_lost_soul {
			tags_off()
			s = sout("%s~%s",
				p_misc(l[i]).save_name, box_code(l[i]))
			tags_on()
			t = sout("%s~%s",
				p_misc(l[i]).save_name, box_code(l[i]))
		} else {
			tags_off()
			s = box_name(l[i])
			tags_on()
			t = box_name(l[i])
		}

		tagout(pl, "<tag type=unit_report id=%d name=\"%s\">", l[i], t)
		tagout(pl, "<tag type=header>")
		out(pl, "")
		lines(pl, s)
		tagout(pl, "</tag type=header>")
		out(pl, "#include %d", l[i])
	}
}

func rep_loc(pl int, l []int) {
	var i int

	sort_for_output(l)

	for i = 0; i < len(l); i++ {
		tagout(pl, "<tag type=loc_report id=%d name=\"%s\">", l[i],
			show_loc_header(l[i]))
		tagout(pl, "<tag type=header>")
		tags_off()
		lines(pl, show_loc_header(l[i]))
		tags_on()
		tagout(pl, "</tag type=header>")
		out(pl, "#include %d", l[i])
		tagout(pl, "<tag type=header>")
		out(pl, "")
		tagout(pl, "</tag type=header>")
	}
}

func inc(pl, code int, s string) {
	tagout(pl, "<tag type=unit_report id=%d name=\"%s\">", code, s)
	tagout(pl, "<tag type=header>")
	lines(pl, s)
	tagout(pl, "</tag type=header>")
	out(pl, "#include %d", code)
	tagout(pl, "</tag type=unit_report>")
	/*VLN this might be the spot to "#include 2" in the GM report */
	tagout(pl, "<tag type=header>")
	out(pl, "")
	tagout(pl, "</tag type=header>")
}

func gen_include_sup(pl int) {
	var char_l, loc_l []int
	var n int
	var player_output, new_flag, loc_flag, code_flag, special_flag, death_flag, misc_flag, eat_queue, eat_warn, eat_error, eat_headers, eat_okay, eat_players, template_flag, garr_flag, drop_flag, show_post bool

	for _, n = range known_sparse_loop(p_player(pl).output) {
		switch n {
		case OUT_BANNER, OUT_INCLUDE:
			continue
		case OUT_LORE:
			//lore_flag = true
			log.Printf("todo: lore_flag set but not used\n")
			continue
		case OUT_NEW:
			new_flag = true
			continue
		case OUT_LOC:
			loc_flag = true
			continue
		case OUT_TEMPLATE:
			template_flag = true
			continue
		case OUT_GARR:
			garr_flag = true
			continue
		case OUT_SHOW_POSTS:
			show_post = true
			continue
		case LOG_CODE:
			code_flag = true
			continue
		case LOG_SPECIAL:
			special_flag = true
			continue
		case LOG_DROP:
			drop_flag = true
			continue
		case LOG_DEATH:
			death_flag = true
			continue
		case LOG_MISC:
			misc_flag = true
			continue
		case EAT_ERR:
			eat_error = true
			continue
		case EAT_WARN:
			eat_warn = true
			continue
		case EAT_QUEUE:
			eat_queue = true
			continue
		case EAT_HEADERS:
			eat_headers = true
			continue
		case EAT_OKAY:
			eat_okay = true
			continue
		case EAT_PLAYERS:
			eat_players = true
			continue
		}

		if !valid_box(n) { /* doesn't exist anymore */
			continue
		}

		switch kind(n) {
		case T_char, T_deadchar:
			char_l = append(char_l, n)
			break

		case T_loc, T_ship:
			loc_l = append(loc_l, n)
			break

		case T_player:
			assert(n == pl)
			player_output = true
			break

		case T_item:
			if subkind(n) == sub_dead_body {
				char_l = append(char_l, n)
			} else {
				panic("!reached")
			}
			break

		default:
			panic("!reached")
		}
	}

	out(pl, "#include %d", OUT_BANNER)
	out(pl, "")

	if eat_okay {
		out(pl, "#include %d", EAT_OKAY)
		out(pl, "")
	}

	if drop_flag {
		inc(pl, LOG_DROP, "Player drops")
	}
	if code_flag {
		inc(pl, LOG_CODE, "Code alerts")
	}
	if special_flag {
		inc(pl, LOG_SPECIAL, "Special events")
	}
	if misc_flag {
		inc(pl, LOG_MISC, "Miscellaneous")
	}
	if death_flag {
		inc(pl, LOG_DEATH, "Character deaths")
	}
	if eat_error {
		inc(pl, EAT_ERR, "Errors")
	}
	if eat_warn {
		inc(pl, EAT_WARN, "Warnings")
	}
	if show_post {
		inc(pl, OUT_SHOW_POSTS, "Press and rumors")
	}
	if eat_queue {
		inc(pl, EAT_QUEUE, "Current order queues")
	}

	if pl != eat_pl && player_output {
		rep_player(pl)
	}

	if garr_flag {
		inc(pl, OUT_GARR, "Garrison log")
	}

	if loc_flag {
		out(pl, "#include %d", OUT_LOC)
	}

	rep_char(pl, char_l)
	rep_loc(pl, loc_l)

	/* if (lore_flag)	inc(pl, OUT_LORE, ""); */

	if new_flag {
		inc(pl, OUT_NEW, "New players")
	}

	if template_flag {
		inc(pl, OUT_TEMPLATE, "Order template")
	}

	if eat_players {
		inc(pl, EAT_PLAYERS, "Current player list")
	}

	if eat_headers {
		inc(pl, EAT_HEADERS, "Original message")
	}
}

func gen_include_section() {
	var pl int

	out_path = MASTER
	out_alt_who = OUT_INCLUDE

	for _, pl = range loop_player() {
		if subkind(pl) != sub_pl_silent {
			gen_include_sup(pl)
		}
		/*
		 *  Any tags needed to finish everything up.
		 *
		 */
		tagout(pl, "</tag type=turn>")
	}

	out_path = 0
	out_alt_who = 0
}

func turn_end_loc_reports() {
	var pl int
	var i int
	var p *entity_player

	stage("turn_end_loc_reports()")

	out_path = MASTER
	show_loc_no_header = true

	for _, pl = range loop_player() {
		if subkind(pl) == sub_pl_silent {
			continue
		}

		separate := (player_format(pl) & ALT) != FALSE

		p = p_player(pl)

		for _, i = range known_sparse_loop(p.locs) {
			if !valid_box(i) { /* loc doesn't exit anymore */
				continue
			} /* ex: mine has collapsed */

			if separate {
				out_alt_who = OUT_LOC

				tags_off()
				lines(pl, show_loc_header(i))
				tags_on()
				show_loc(pl, i)
				out(pl, "")
			} else {
				out_alt_who = i
				out(pl, "")
				show_loc(pl, i)
			}
		}
	}

	out_path = 0
	out_alt_who = 0
	show_loc_no_header = false
}

func player_banner() {
	var pl int

	stage("player_banner()")

	out_path = MASTER
	out_alt_who = OUT_BANNER

	for _, pl = range loop_player() {
		if subkind(pl) == sub_pl_silent {
			continue
		}

		//p := p_player(pl); // mdhender: p not used

		//#if 0
		//                out(pl, "From: %s", from_host);
		//                out(pl, "Reply-To: %s", reply_host);
		//                if (p.email)
		//                    out(pl, "To: %s (%s)", p.email,
		//                        or_string(p.full_name!="" , p.full_name , "???"));
		//                out(pl, "Subject: Olympia:TAG game %d, turn %d report", game_number, sysclock.turn);
		//                out(pl, "");
		//#endif

		//#if 0
		//                switch (player_compuserve(pl))
		//                {
		//                case 1:
		//                    indent += 3;
		//                    wout(pl, "Note:  Please download the Olympia paper for this turn from Library 3 in the PBMGAMES forum on Compuserve.  It should be available in file OT%d.OLY as soon as it is released by the sysops.", sysclock.turn);
		//                    indent -= 3;
		//                    out(pl, "");
		//                    break;
		//
		//                case 2:
		//                    indent += 3;
		//                    wout(pl, "Note:  Olympia paper mailing is turned off for you.To begin receiving the paper again, issue the order TIMES 0");
		//                    indent -= 3;
		//                    out(pl, "");
		//                    break;
		//                }
		//#else
		if player_compuserve(pl) {
			indent += 3
			wout(pl, "Note:  Olympia paper mailing is turned off for you.To begin receiving the paper via email again, issue the order TIMES 0")
			indent -= 3
			out(pl, "")
			break
		}
		//#endif

		tagout(pl, "<tag type=turn pl=%d game=%d turn=%d>",
			pl, game_number, sysclock.turn)
		wout(pl, "Olympia:TAG game %d turn %d report for %s.",
			game_number,
			sysclock.turn, box_name(pl))

		month := oly_month(&sysclock)
		year := oly_year(&sysclock)

		wout(pl, "Season \"%s\", month %d, in the year %d.",
			month_names[month],
			month+1,
			year+1)

		out(pl, "")
	}

	out_path = 0
	out_alt_who = 0
}

func report_account_out(pl, who int) int {
	var fnam string
	var cmd string

	fnam = filepath.Join("/", "tmp", "oly-acct.tmp")

	cmd = fmt.Sprintf("%s -a %s -A %s-old -g tag%d -p %s -s 4 > %s",
		options.accounting_prog,
		options.accounting_dir,
		options.accounting_dir,
		game_number,
		box_code_less(pl), fnam)
	system(cmd)

	fp, err := fopen(fnam, "r")
	if err != nil {
		out(who, "<account summary not available>")
		out(who, "")
		log.Printf("can't open %s: ", fnam)
		unlink(fnam)
		return 0
	}

	for line := getlin(fp); line != nil; line = getlin(fp) {
		out(who, "%s", string(line))
	}

	fclose(fp)
	unlink(fnam)
	return 1

}

func report_account_sup(pl int) {
	tagout(pl, "<tag type=account_summary>")
	tagout(pl, "<tag type=header>")
	out(pl, "Account summary")
	out(pl, "---------------")
	tagout(pl, "</tag type=header>")

	indent += 3
	report_account_out(pl, pl)
	indent -= 3
	tagout(pl, "<tag type=header>")
	out(pl, "---------------")
	out(pl, "")
	tagout(pl, "</tag type=header>")
	tagout(pl, "</tag type=account_summary>")
}

func report_account() {
	var pl int

	stage("report_account()")
	close_logfile()

	out_path = MASTER
	out_alt_who = OUT_BANNER

	for _, pl = range loop_player() {
		if subkind(pl) != sub_pl_regular {
			continue
		}

		report_account_sup(pl)
	}

	out_path = 0
	out_alt_who = 0
}

func charge_account() {
	var pl int
	var p *entity_player
	var cmd string
	//var val_s string

	stage("charge_account()")

	out_path = MASTER
	out_alt_who = OUT_BANNER

	for _, pl = range loop_player() {
		if subkind(pl) != sub_pl_regular {
			continue
		}

		p = rp_player(pl)
		if p == nil {
			continue
		}

		cmd = fmt.Sprintf("%s -g tag%d -a %s -A %s-old -p %s -t%s -y \"Olympia TAG turn %d\"",
			options.accounting_prog,
			game_number,
			options.accounting_dir,
			options.accounting_dir,
			box_code_less(pl),
			turn_charge(pl),
			sysclock.turn)

		if system(cmd) != FALSE {
			out(gm_player, "Could not charge %s.",
				box_code(pl))
		}
	}

	out_path = 0
	out_alt_who = 0
}

func print_accept_sup(who int, a *accept_ent, parens bool) {
	var buf string

	buf = fmt.Sprintf("Accepting ")

	if a.item == 0 {
		buf += ("any amount of anything ")
	} else {
		if a.qty == 0 {
			buf += ("any amount of ")
		} else {
			buf += ("up to ")
			buf += (nice_num(a.qty))
		}
		buf += (box_name(a.item))
		buf += (" ")
	}

	buf += ("from ")

	if a.from_who == 0 {
		buf += ("anyone")
	} else if kind(a.from_who) == T_nation {
		buf += ("any member of ")
		buf += (rp_nation(a.from_who).name)
	} else {
		buf += (box_name(a.from_who))
	}

	if parens {
		out(who, "(%s.)", buf)
	} else {
		out(who, "%s.", buf)
	}

}

func list_accepts(who, num int) {
	var i int
	var accept []*accept_ent
	flag := true

	assert(valid_box(who) && p_char(num) != nil)
	accept = p_char(who).accept

	out(who, "")
	out(who, "Accepting:")
	out(who, "")
	indent += 3

	for i = 0; i < len(accept); i++ {
		print_accept_sup(num, accept[i], false)
		flag = false
	}

	/*
	 *  Wed Jan 13 12:57:35 1999 -- Scott Turner
	 *
	 *  And print the "player" accepts too.
	 *
	 */
	if player(num) != num {
		accept = p_char(player(num)).accept
		for i = 0; i < len(accept); i++ {
			print_accept_sup(num, accept[i], true)
			flag = false
		}
	}

	if flag {
		out(num, "Nothing")
	}

	indent -= 3
}
