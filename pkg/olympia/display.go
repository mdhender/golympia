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

import "fmt"

var (
	combat_ally         string
	show_combat_flag    = false
	show_display_string = FALSE
	show_loc_no_header  = false
)

func any_chars_here(where int) int { panic("!implemented") }

func loc_inside_string(where int) string {
	var name string
	var trail string

	if loc_depth(where) > LOC_region {
		trail = loc_inside_string(loc(where))
	}

	switch loc_depth(where) {
	case LOC_region:
		name = just_name(region(where))
		break
	case LOC_province:
		if !valid_box(where) {
			name = sout("adrift in the Cosmos")
		} else {
			name = box_name(where)
		}
		break
	default:
		name = box_name(where)
	}

	if len(trail) != 0 {
		return sout(", in %s%s", name, trail)
	}

	return sout(", in %s", name)
}

func show_loc_barrier(who, where int) int {
	if loc_barrier(where) != FALSE {
		tagout(who, "<tag type=loc_barrier id=%d>", where)
		wout(who, "A magical barrier surrounds %s.", box_name(where))
		tagout(who, "</tag type=loc_barrier>")
		return TRUE
	}

	return FALSE
}

func safe_haven_s(n int) string {

	if safe_haven(n) {
		return ", safe haven"
	}

	return ""
}

func ship_cap_s(n int) string {
	if sc := ship_cap(n); sc != FALSE {
		sw := ship_weight(n)
		return sout(", %d%% loaded", sw*100/sc)
	}

	return ""
}

func show_loc_stats(who, where int) {
	var sc, sw, n int
	first := true

	tagout(who, "<tag type=loc_stats loc=%d>", where)

	if n = loc_damage(where); n != FALSE {
		if first {
			out(who, "")
			first = false
		}
		tagout(who, "<tag type=loc_dam loc=%d damage=%d tot=%d>", where, n, loc_hp(where))
		out(who, "Damage: %d%%", (int)(n*100)/loc_hp(where))
	}

	if sc = ship_cap(where); sc != FALSE {
		sw = ship_weight(where)

		if first {
			out(who, "")
			first = false
		}

		tagout(who, "<tag type=loc_cap loc=%d cap=%d weight=%d>", where, sc, sw)
		out(who, "Ship capacity: %s/%s (%d%%)", comma_num(sw), comma_num(sc), sw*100/sc)
	}

	tagout(who, "</tag type=loc_stats>")
}

func loc_civ_s(where int) string {
	var n int

	if loc_depth(where) != LOC_province ||
		subkind(where) == sub_ocean ||
		in_faery(where) || in_hades(where) {
		return ""
	}

	n = has_item(where, item_peasant)

	if n < 100 {
		return ", wilderness"
	}

	return sout(", peasants: %d", n)
}

/*
 *  Thu Nov 12 11:20:28 1998 -- Scott Turner
 *
 *  Add display of shoring.
 *
 */
func show_loc_header(where int) string {
	var buf string

	buf = (box_name_kind(where))
	buf += (loc_inside_string(loc(where)))

	if subkind(where) == sub_mine_shaft {
		mi := get_mine_info(where)
		if mine_depth(where) != FALSE {
			buf += (sout(", depth~%d feet",
				(mine_depth(where)*100)+100))
		}
		if mi != nil {
			switch mi.shoring[mine_depth(where)] {
			case WOODEN_SHORING:
				buf += (sout(", wooden shoring"))
				break
			case IRON_SHORING:
				buf += (sout(", iron shoring"))
				break
			//case NO_SHORING:
			default:
				buf += (sout(", no shoring"))
				break
			}
		}
	}

	buf += (safe_haven_s(where))

	if loc_hidden(where) {
		buf += (", hidden")
	}

	buf += (loc_civ_s(where))

	return sout("%s", buf)
}

func with_inventory_string(who int) string {
	var with string
	var e *item_ent
	var mk int

	mk = noble_item(who)

	for _, e = range inventory_loop(who) {
		if mk == e.item || item_prominent(e.item) == FALSE {
			continue
		}

		if with == "" {
			with = (", ")
		} else {
			with = (", with ")
		}

		with += (just_name_qty(e.item, e.qty))
	}

	return sout("%s", with)
}

func display_with(who int) string {

	if rp_loc_info(who) != nil && len(rp_loc_info(who).here_list) > 0 {
		return ", accompanied~by:"
	}

	return ""
}

func display_owner(who int) string {

	if first_character(who) <= 0 {
		return ""
	}

	return ", owner:"
}

func incomplete_string(n int) string {
	var p *entity_subloc
	var b *entity_build

	p = rp_subloc(n)
	if p == nil {
		return ""
	}

	b = get_build(n, BT_BUILD)

	if b == nil || b.effort_required == 0 {
		return ""
	}

	return sout(", %d%% completed",
		b.effort_given*100/b.effort_required)
}

func liner_desc_ship(n int) string {
	var buf string
	ship := rp_ship(n)

	buf = fmt.Sprintf("%s%s",
		box_name_kind(n),
		incomplete_string(n))

	if ship != nil {
		if ship.hulls != FALSE {
			buf += (sout(", %d hull%s", ship.hulls,
				or_string(ship.hulls == 1, "", "s")))
		}
		if ship.ports != FALSE {
			buf += (sout(", %d rowing port%s", ship.ports,
				or_string(ship.ports == 1, "", "s")))
		}
		if ship.sails != FALSE {
			buf += (sout(", %d sail%s", ship.sails,
				or_string(ship.sails == 1, "", "s")))
		}
	}

	//#if 0
	//    buf += ( ship_cap_s(n));
	//
	//    if (loc_defense(n))
	//      buf += ( sout(", defense~%d", loc_defense(n)));
	//#endif

	if loc_damage(n) != FALSE {
		buf += (sout(", %d%%~damaged", (loc_damage(n)*100)/loc_hp(n)))
	}

	if show_display_string != FALSE {
		s := banner(n)
		if len(s) != 0 {
			buf += (sout(", \"%s\"", s))
		}
	}

	if loc_hidden(n) {
		buf += (", hidden")
	}

	if ship_has_ram(n) != FALSE {
		buf += (", with ram")
	}

	if rp_subloc(n) != nil {
		if rp_subloc(n).control.nobles != FALSE {
			buf += (sout(", boarding fee %s per noble",
				gold_s(rp_subloc(n).control.nobles)))
		}
		if rp_subloc(n).control.weight != FALSE {
			buf += (sout(", boarding fee %s per 1000 weight",
				gold_s(rp_subloc(n).control.weight)))
		}
		if rp_subloc(n).control.men != FALSE {
			buf += (sout(", boarding fee %s per 100 men",
				gold_s(rp_subloc(n).control.men)))
		}
	}

	//#if 0
	//    if (fee = board_fee(n))
	//      buf += ( ", %s per 100 wt. to board", gold_s(fee));
	//#endif

	return sout("%s", buf)
}

/*
 *	Name, mountain province, in region foo
 *	Name, port city, in province Name [, in region foo]
 *	Name, ocean, in Sea
 *	Name, island, in Ocean [, in Sea]
 *
 *	Mountain [aa01], mountain province, in region Tollus
 *	Island [aa01], island, in Ocean [bb01]
 *	City [aa01], port city, in province Mountain [aa01]
 *	Ocean [bb02], ocean, in South Sea
 */

func liner_desc_loc(n int) string {
	var buf string

	buf = fmt.Sprintf("%s%s%s",
		box_name_kind(n),
		safe_haven_s(n),
		incomplete_string(n))

	if loc_depth(n) == LOC_province &&
		rp_loc(n) != nil && province_admin(n) != FALSE {
		if rp_loc(n).control.nobles != FALSE {
			buf += (sout(", entrance fee %s per noble",
				gold_s(rp_loc(n).control.nobles)))
		}
		if rp_loc(n).control.weight != FALSE {
			buf += (sout(", entrance fee %s per 1000 weight",
				gold_s(rp_loc(n).control.weight)))
		}
		if rp_loc(n).control.men != FALSE {
			buf += (sout(", entrance fee %s per 100 men",
				gold_s(rp_loc(n).control.men)))
		}
	}

	if rp_subloc(n) != nil && first_character(n) != FALSE {
		if rp_subloc(n).control.nobles != FALSE {
			buf += (sout(", entrance fee %s per noble",
				gold_s(rp_subloc(n).control.nobles)))
		}
		if rp_subloc(n).control.weight != FALSE {
			buf += (sout(", entrance fee %s per 1000 weight",
				gold_s(rp_subloc(n).control.weight)))
		}
		if rp_subloc(n).control.men != FALSE {
			buf += (sout(", entrance fee %s per 100 men",
				gold_s(rp_subloc(n).control.men)))
		}
	}

	if loc_depth(n) == LOC_province &&
		rp_loc(n) != nil && province_admin(n) != FALSE {
		if rp_loc(n).control.closed {
			buf += (sout(", border closed"))
		}
	}

	if rp_subloc(n) != nil && first_character(n) != FALSE {
		if rp_subloc(n).control.closed {
			buf += (sout(", closed"))
		}
	}

	if entrance_size(n) != FALSE {
		buf += (sout(", entrance size: %d", entrance_size(n)))
	}

	if loc_depth(n) == LOC_build {
		if loc_defense(n) != FALSE {
			buf += (sout(", defense~%d", loc_defense(n)))
		}

		if get_effect(n, ef_improve_fort, 0, 0) != FALSE {
			buf += (sout(", magical defense~%d",
				get_effect(n, ef_improve_fort, 0, 0)))
		}

		if loc_moat(n) != FALSE {
			buf += (", with moat")
		}

		if loc_damage(n) != FALSE {
			buf += (sout(", %d%%~damaged", (100*loc_damage(n))/loc_hp(n)))
		}
	}

	if show_display_string != FALSE {
		s := banner(n)

		if len(s) != 0 {
			buf += (sout(", \"%s\"", s))
		}
	}

	if loc_hidden(n) {
		buf += (", hidden")
	}

	if subkind(n) == sub_mine_shaft {
		mi := get_mine_info(n)
		if mine_depth(n) != FALSE {
			buf += (sout(", depth~%d feet",
				(mine_depth(n)*100)+100))
		}
		if mi != nil {
			switch mi.shoring[mine_depth(n)] {
			case WOODEN_SHORING:
				buf += (sout(", wooden shoring"))
				break
			case IRON_SHORING:
				buf += (sout(", iron shoring"))
				break
			case NO_SHORING:
			default:
				buf += (sout(", no shoring"))
				break
			}
		}
	}

	if loc_depth(n) == LOC_subloc {
		if subkind(n) == sub_hades_pit {
			buf += (", 28 days")
		} else if subkind(n) != sub_mine_shaft &&
			subkind(n) != sub_mine_shaft_notdone {
			buf += (", 1 day")
		}
	}

	return sout("%s", buf)
}

func mage_s(n int) string {
	var a int

	if is_magician(n) || char_hide_mage(n) != FALSE {
		return ""
	}

	a = max_eff_aura(n)

	if a <= 5 {
		return ""
	}
	if a <= 10 {
		return ", conjurer"
	}
	if a <= 15 {
		return ", mage"
	}
	if a <= 20 {
		return ", wizard"
	}
	if a <= 30 {
		return ", sorcerer"
	}
	if a <= 40 {
		return ", 6th black circle"
	}
	if a <= 50 {
		return ", 5th black circle"
	}
	if a <= 60 {
		return ", 4th black circle"
	}
	if a <= 70 {
		return ", 3rd black circle"
	}
	if a <= 80 {
		return ", 2nd black circle"
	}

	return ", master of the black arts"
}

func nation_s(n int) string {
	var nation_title string
	var neutral_title string

	/*
	 *  Wed Apr 16 11:29:27 1997 -- Scott Turner
	 *
	 *  Check for a concealment.
	 *
	 */
	if get_effect(n, ef_conceal_nation, 0, 0) != FALSE {
		new_nation := get_effect(n, ef_conceal_nation, 0, 0)
		assert(new_nation >= 1 && kind(new_nation) == T_nation)
		if rp_nation(new_nation).neutral {
			neutral_title = (", neutral")
		}
		nation_title = fmt.Sprintf(", %s%s", rp_nation(new_nation).name, neutral_title)
	} else if nation(n) != FALSE {
		if rp_nation(nation(n)).neutral {
			neutral_title = (", neutral")
		}
		nation_title = fmt.Sprintf(", %s%s", rp_nation(nation(n)).name, neutral_title)
	} else if subkind(n) == sub_garrison {
		/*
		 *  Special case for an uncommanded garrison?
		 *
		 */
		nation_title = (", uncontrolled")
	} else if refugee(n) {
		nation_title = (", refugee")
	}

	return nation_title

}

func deserted_s(n int) string {
	var deserted string
	pl := player(n)

	/*
	 *  Wed Apr 14 10:50:34 1999 -- Scott Turner
	 *
	 *  Print "deserted" if a noble belongs to an NPC faction
	 *  but has a body_old_lord.
	 *
	 */
	if is_real_npc(pl) && body_old_lord(n) != FALSE {
		deserted = (", deserted")
	}

	return deserted
}

func priest_s(n int) string {
	var priest_title string
	var e *entity_religion_skill

	if is_priest(n) == FALSE {
		return ""
	}

	e = rp_relig_skill(is_priest(n))

	if e != nil && e.high_priest == n {
		priest_title = fmt.Sprintf(", the High Priest of the %s", box_name(is_priest(n)))
	} else if e != nil && (e.bishops[0] == n || e.bishops[1] == n) {
		priest_title = fmt.Sprintf(", Bishop of the %s", box_name(is_priest(n)))
	} else {
		priest_title = fmt.Sprintf(", priest of the %s", box_name(is_priest(n)))
	}
	return priest_title

}

func liner_desc_char(n int) string {
	//extern int show_combat_flag;
	var buf string
	var s string

	buf = (box_name(n))

	sk := subkind(n)

	if sk == sub_ni {
		mk := noble_item(n)
		num := has_item(n, mk) + 1

		if num == 1 {
			buf += (sout(", %s", plural_item_name(mk, num)))
		} else {
			buf += (sout(", %s, number:~%s",
				plural_item_name(mk, num),
				comma_num(num)))
		}
	} else if sk != FALSE {
		if sk == sub_temple {
			if is_temple(n) != FALSE {
				buf += (sout(", Temple of %s", god_name(is_temple(n))))
			} else {
				buf += (sout(", undedicated temple"))
			}
		} else if sk == sub_guild {
			buf += (sout(", %s Guild", box_name(is_guild(n))))
		} else {
			buf += (sout(", %s", subkind_s[sk]))
		}
	}

	buf += (nation_s(n))
	buf += (deserted_s(n))
	buf += (rank_s(n))
	buf += (mage_s(n))
	buf += (priest_s(n))
	//#if 0
	//    buf += ( wield_s(n));
	//#endif

	if show_combat_flag {
		if char_behind(n) != FALSE {
			buf += (sout(", behind~%d%s",
				char_behind(n), combat_ally))
		} else {
			buf += (combat_ally)
		}
	} else if char_guard(n) != FALSE && stack_leader(n) == n &&
		subkind(n) != sub_garrison {
		buf += (", on guard")
	}

	//#if 0
	//    if (subkind(n) == 0)	/* only show lord for regular players */
	//    {
	//        int sp = lord(n);
	//
	//        if (sp != indep_player && !cloak_lord(n))
	//            buf += ( sout(", of~%s", box_code_less(sp)));
	//    }
	//#endif

	if show_display_string != FALSE {
		s = banner(n)

		if len(s) != 0 {
			buf += (sout(", \"%s\"", s))
		}
	}

	buf += (with_inventory_string(n))

	if is_prisoner(n) {
		buf += (", prisoner")
	}

	return sout("%s", buf)
}

func liner_desc_road(n int) string {
	var dest int
	hid := ""
	var dist int

	dest = road_dest(n)

	if road_hidden(n) != FALSE {
		hid = ", hidden"
	}

	dist = exit_distance(loc(n), dest)

	return sout("%s, to %s%s, %d~day%s",
		box_name(n), box_name(dest), hid,
		add_ds(dist))
}

func liner_desc_storm(n int) string {
	var buf string
	var owner int
	var p *entity_misc

	buf = fmt.Sprintf("%s", box_name_kind(n))

	p = rp_misc(n)

	owner = npc_summoner(n)
	if owner != FALSE {
		buf += (sout(", owner~%s", box_code_less(owner)))
	}

	buf += (sout(", strength~%s", comma_num(storm_strength(n))))

	if p != nil && p.npc_dir != FALSE {
		buf += (sout(", heading %s", full_dir_s[p.npc_dir]))
	}

	return sout("%s", buf)
}

/*
 *  Viewed from outside
 */

func liner_desc(n int) string {

	switch kind(n) {
	case T_ship:
		return liner_desc_ship(n)
	case T_loc:
		return liner_desc_loc(n)
	case T_char:
		return liner_desc_char(n)
	case T_road:
		return liner_desc_road(n)
	case T_storm:
		return liner_desc_storm(n)
	}
	panic("!reached")
}

func highlight_units(who int, n int, depth int) string {

	assert(depth >= 3)
	assert(indent == 0)

	if kind(who) == T_player && player(n) == who {
		return sout(" *%s", string(spaces[:spaces_len-(depth-2)]))
	}

	return string(spaces[:spaces_len-depth])
}

func show_chars_below(who int, n int) {
	var i int

	assert(valid_box(who))

	indent += 3
	for _, i = range loop_char_here(n) {
		assert(valid_box(who))
		wiout(who, 3, "%s", liner_desc(i))
	}

	indent -= 3
}

func show_chars_below_highlight(who int, n int, depth int, where int) {
	var i int

	depth += 3

	for _, i = range loop_char_here(n) {
		tagout(who, "<tag type=char_here id=%d where=%d under=%d>",
			i, where, n)
		wiout(who, depth, "%s%s",
			highlight_units(who, i, depth),
			liner_desc(i))
		tagout(who, "</tag type=char_here>")
	}

}

func show_owner_stack(who int, n int) {
	var i int
	var depth int
	first := TRUE

	depth = indent + 3
	indent = 0

	for _, i = range loop_here(n) {
		if kind(i) == T_char {
			if first == FALSE && char_really_hidden(i) {
				continue
			}

			wiout(who, depth, "%s%s%s",
				highlight_units(who, i, depth),
				liner_desc(i),
				display_with(i))

			show_chars_below_highlight(who, i, depth, n)

			first = FALSE
		}
	}

	indent = depth - 3
}

func show_chars_here(who, where int) {
	first := TRUE
	var i int
	var depth int
	flying := ""

	tagout(who, "<tag type=chars_here where=%d>", where)
	if loc_depth(where) == LOC_province &&
		weather_here(where, sub_fog) != FALSE &&
		is_priest(who) != sk_domingo {
		out(who, "")
		out(who, "No one can be seen through the fog.")
		return
	}

	if loc_depth(where) == LOC_province &&
		weather_here(where, sub_mist) != FALSE &&
		is_priest(who) != sk_domingo {
		out(who, "")
		out(who, "A dank, unhealthy mist conceals everything.")
		return
	}

	depth = indent
	indent = 0

	if subkind(where) == sub_ocean {
		flying = ", flying"
	}

	for _, i = range loop_here(where) {
		if kind(i) == T_char {
			if char_really_hidden(i) {
				continue
			}

			if first != FALSE {
				first = FALSE
				out(who, "")
				out(who, "Seen here:")
				depth += 3
			}

			tagout(who, "<tag type=char_here id=%d where=%d under=0>",
				i, where)
			wiout(who, depth, "%s%s%s%s",
				highlight_units(who, i, depth),
				liner_desc(i),
				flying,
				display_with(i))
			tagout(who, "</tag type=char_here>")

			show_chars_below_highlight(who, i, depth, where)
		}
	}

	if first == FALSE {
		depth -= 3
	}

	indent = depth
	tagout(who, "</tag type=chars_here>")
}

func show_inner_locs(who, where int) {
	first := TRUE
	var i int

	for _, i = range loop_here(where) {
		if is_loc_or_ship(i) {
			if loc_hidden(i) &&
				!test_known(who, i) &&
				see_all(who) == FALSE {
				continue
			}

			if first != FALSE {
				first = FALSE
				if subkind(who) != sub_mine_shaft {
					indent += 3
				}
			}

			wout(who, "%s%s", liner_desc(i), display_owner(i))
			show_owner_stack(who, i)
			show_loc_barrier(who, i)
			show_inner_locs(who, i)
		}
	}

	if first == FALSE {
		indent -= 3
	}
}

func show_sublocs_here(who, where int) {
	first := TRUE
	var i int
	var p *entity_subloc

	tagout(who, "<tag type=sublocs_here loc=%d>", where)
	for _, i = range loop_here(where) {
		if kind(i) == T_loc {
			if loc_hidden(i) &&
				!test_known(who, i) &&
				see_all(who) == FALSE {
				continue
			}

			if first != FALSE {
				first = FALSE
				out(who, "")
				out(who, "Inner locations:")
				indent += 3
			}

			tagout(who, "<tag type=subloc id=%d where=%d>",
				i, where)
			if subkind(i) == sub_city {
				wout(who, "%s", liner_desc(i))
				show_loc_barrier(who, i)
			} else {
				wout(who, "%s%s", liner_desc(i),
					display_owner(i))
				show_owner_stack(who, i)
				show_loc_barrier(who, i)
				show_inner_locs(who, i)
			}
			tagout(who, "</tag type=subloc>")
		}
	}

	p = rp_subloc(where)

	if p != nil {
		for i = 0; i < len(p.link_from); i++ {
			if loc_hidden(p.link_from[i]) &&
				!test_known(who, p.link_from[i]) &&
				see_all(who) == FALSE {
				continue
			}

			/*
			 *  Mon Dec  9 15:19:58 1996 -- Scott Turner
			 *
			 *  All links are "open" now...
			 *
			 */
			/*		if (loc_link_open(p.link_from[i])) */
			if p.link_from[i] != FALSE {
				if first != FALSE {
					first = FALSE
					out(who, "")
					out(who, "Inner locations:")
					indent += 3
				}

				tagout(who, "<tag type=subloc id=%d where=%d>",
					p.link_from[i], where)
				wout(who, "%s", liner_desc(p.link_from[i]))
				tagout(who, "</tag type=subloc>")
			}

		}
	}

	if first == FALSE {
		indent -= 3
	}
	tagout(who, "<tag type=sublocs_here>")
}

func show_ships_here(who, where int) {
	first := TRUE
	var i int

	tagout(who, "<tag type=ships_here where=%d>", where)

	for _, i = range loop_here(where) {
		if kind(i) == T_ship {
			if loc_hidden(i) &&
				!test_known(who, i) &&
				see_all(who) == FALSE {
				continue
			}

			if first != FALSE {
				first = FALSE

				out(who, "")
				if subkind(where) == sub_ocean {
					out(who, "Ships sighted:")
				} else {
					out(who, "Ships docked at port:")
				}

				indent += 3
			}

			tagout(who, "<tag type=ship_here id=%d where=%d owner=%d>",
				i, where,
				first_character(i))
			wiout(who, 3, "%s%s", liner_desc(i), display_owner(i))
			show_owner_stack(who, i)
			show_loc_barrier(who, i)
			tagout(who, "</tag type=ship_here>")
		}
	}

	if first == FALSE {
		indent -= 3
	}

	tagout(who, "</tag type=ships_here>")
}

func show_nearby_cities(who, where int) {
	var p *entity_subloc
	var i int
	var s string

	p = rp_subloc(where)

	if p == nil || len(p.near_cities) < 1 {
		return
	}

	tagout(who, "<tag type=nearby_cities loc=%d>", where)
	out(who, "")
	out(who, "Cities rumored to be nearby:")
	indent += 3
	for i = 0; i < len(p.near_cities); i++ {
		if safe_haven(p.near_cities[i]) {
			s = ", safe haven"
		} else {
			s = ""
		}

		tagout(who, "<tag type=nearby_city loc=%d city=%d name=\"%s\" province=%d province_name=\"%s\">",
			where,
			p.near_cities[i],
			box_name(p.near_cities[i]),
			province(p.near_cities[i]),
			box_name(province(p.near_cities[i])))
		out(who, "%s, in %s%s",
			box_name(p.near_cities[i]),
			box_name(province(p.near_cities[i])), s)
	}
	indent -= 3
	tagout(who, "</tag type=nearby_cities>")
}

func show_loc_skills(who, where int) {
	var i int
	s := ""

	tagout(who, "<tag type=taught_here loc=%d>", where)

	for _, i = range loop_loc_teach(where) {
		tagout(who, "<tag type=taught_skill id=%d name=\"%s\">",
			i, box_name(i))
		s = comma_append(s, box_name(i))
	}

	if len(s) != 0 {
		out(who, "")
		out(who, "Skills taught here:")
		indent += 3
		wout(who, "%s", s)
		indent -= 3
	}

	tagout(who, "</tag type=taught_here>")
}

func show_loc_posts(who int, where int, show_full_loc int) {
	var post int
	var i int
	flag := TRUE
	var first int
	var l []string

	tagout(who, "<tag type=loc_posts>")
	for _, post = range loop_here(where) {
		if kind(post) != T_post {
			continue
		}

		if rp_misc(post) == nil ||
			len(rp_misc(post).post_txt) < 1 {
			//continue
			panic("!reached") /* what happened to the post? */
		}

		l = rp_misc(post).post_txt

		if flag != FALSE {
			out(who, "")
			wout(who, "Posted in %s:",
				or_string(show_full_loc != FALSE, box_name(where), just_name(where)))
			flag = FALSE
			indent += 3
		} else {
			out(who, "")
		}

		if item_creator(post) != FALSE {
			wout(who, "Posted by %s:",
				box_name(item_creator(post)))
		} else {
			wout(who, "Posted:")
		}

		indent += 3

		first = TRUE

		for i = 0; i < len(l); i++ {
			wout(who, "%s%s%s",
				or_string(first != FALSE, "\"", ""),
				l[i],
				or_string(i+1 == len(l), "\"", ""))

			if first != FALSE {
				first = FALSE
				indent += 1
			}
		}

		if first == FALSE {
			indent -= 1
		}

		indent -= 3
	}

	if flag == FALSE {
		indent -= 3
	}
	tagout(who, "</tag type=loc_posts>")
}

func show_weather(who, where int) {
	rain := weather_here(where, sub_rain)
	wind := weather_here(where, sub_wind)
	fog := weather_here(where, sub_fog)
	mist := weather_here(where, sub_mist)

	if rain == FALSE && wind == FALSE && fog == FALSE && mist == FALSE {
		return
	}

	tagout(who, "<tag type=weather loc=%d rain=%d wind=%d fog=%d mist=%d>",
		where, rain, wind, fog, mist)

	out(who, "")

	if rain != FALSE {
		out(who, "It is raining.")
	}

	if wind != FALSE {
		out(who, "It is windy.")
	}

	if fog != FALSE {
		out(who, "The province is blanketed in fog.")
	}

	if mist != FALSE {
		out(who, "The province is covered with a dank mist.")
	}

	if can_see_weather_here(who, where) {
		var i int
		first := TRUE

		for _, i = range loop_here(where) {
			if kind(i) != T_storm {
				continue
			}

			if first != FALSE {
				indent += 3
				first = FALSE
			}

			tagout(who, "<tag type=storm id=%d loc=%d>",
				i, where)
			wout(who, "%s", liner_desc(i))
			tagout(who, "<tag type=storm>")
		}

		if first == FALSE {
			indent -= 3
		}
	}
	tagout(who, "</tag type=weather>")
}

func show_loc_ruler(who, where int) {
	var garr int
	var castle int
	var ruler int
	var fee int
	var buf string

	if loc_depth(where) != LOC_province || subkind(where) == sub_ocean {
		return
	}

	garr = garrison_here(where)

	if garr == 0 {
		return
	}

	castle = garrison_castle(garr)
	ruler = top_ruler(garr)

	if castle != FALSE {
		prov := province(castle)

		tagout(who, "<tag type=ruler id=%d castle=%d province=%d ruler=%d>",
			castle, prov, ruler)

		out(who, "Province controlled by %s, in %s",
			box_name_kind(castle), box_name(prov))

		if ruler != FALSE {
			out(who, "Ruled by %s%s",
				box_name(ruler), rank_s(ruler))
		}

		tagout(who, "</tag type=ruler>")
	}

	if ruler == 0 {
		return
	}

	if rp_loc(where).control.closed {
		tagout(who, "<tag type=border_closed id=%d>", where)
		out(who, "Border closed.")
		tagout(who, "</tag type=border_closed>")
	}

	if rp_loc(where).control.nobles != FALSE {
		buf += (sout("%s per noble", gold_s(rp_loc(where).control.nobles)))
		fee = 1
	}

	if rp_loc(where).control.weight != FALSE {
		if fee != FALSE {
			buf += (", ")
		}
		buf += (sout("entrance fee %s per 1000 weight",
			gold_s(rp_loc(where).control.weight)))
		fee = 1
	}
	if rp_loc(where).control.men != FALSE {
		if fee != FALSE {
			buf += (", ")
		}
		buf += (sout(", entrance fee %s per 100 men",
			gold_s(rp_loc(where).control.men)))
		fee = 1
	}

	if fee != FALSE {
		tagout(who, "<tag type=fees id=%d nobles=%d weight=%d men=%d>",
			where, rp_loc(where).control.nobles,
			rp_loc(where).control.weight,
			rp_loc(where).control.men)
		buf = ("Entrance fees: ")
		out(who, buf)
		tagout(who, "</tag type=fees>")
	}

	out(who, "")
}

/*
 *  Don't include the leading location name, kind, etc. on the location
 *  report.  Handled with a global since we don't have default parameters.
 */
func show_loc(who, where int) {
	var pil int

	assert(valid_box(where))
	if loc_depth(where) < LOC_province {
		return
	}

	show_display_string = TRUE /* previously only for show_chars_here */

	tagout(who, "</tag type=loc_report>")
	tagout(who, "<tag type=location id=%d>", where)

	if !show_loc_no_header {
		wout(who, "%s", show_loc_header(where))
		out(who, "")
	}

	if show_loc_barrier(who, where) != FALSE {
		out(who, "")
	}

	if pil = loc_pillage(where); pil != FALSE {
		tagout(who, "<tag type=pillaged id=%d>", where)
		wout(who, "Recent pillaging has frightened the peasants.")
		out(who, "")
		tagout(who, "</tag type=pillaged>")
	}

	show_loc_ruler(who, where)

	list_sailable_routes(who, where)
	list_exits(who, where)
	show_loc_stats(who, where)
	show_nearby_cities(who, where)
	show_loc_skills(who, where)

	if subkind(where) == sub_city ||
		is_guild(where) == sk_trading {
		market_report(who, where)
	}

	show_sublocs_here(who, where)
	show_loc_posts(who, where, FALSE)
	show_weather(who, where)
	show_chars_here(who, where)
	show_ships_here(who, where)

	show_display_string = FALSE

	tagout(who, "</tag type=location>")

}
