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
	RANK_lord     = 5
	RANK_knight   = 6
	RANK_baron    = 7
	RANK_count    = 8
	RANK_earl     = 9
	RANK_marquess = 10
	RANK_duke     = 11
	RANK_king     = 20
)

func garrison_gold()                        { panic("!implemented") }
func players_who_rule_here(where int) []int { panic("!implemented") }
func touch_garrison_locs()                  { panic("!implemented") }

/*
 *	garrison . castle . owner [ . char ]*
 *			     "admin"          "top ruler"
 */

/*

1-	touch loc on pledge of terroritory
2-	initial touch loc for both castle and ruler
3-	hook pledge into glob.c
4-	determine everyone's status at turn end
5-	display status in display.c
6.	who may get/give to a garrison
7-	maintenance of men in garrison
8-	forwarding of gold to castle
9-	circular pledge detector
10-	replace loc_owner with building_owner, province_ruler
11.	credit for castles as well as garrisoned provinces
12-	"garrison" keyword that matches the garrison in the current loc
13-	take -- leave behind 10
14-	there can't be an ni - 0, can there?  what if attacked?
15-	status = min(own status, pledge lord's status - 1)
16-	allow owner or admin to name province, sublocs
17-	garrison log for output

*/

/*
 *  Garrison should always be first; we should just have to look
 *  at the first character
 */

func garrison_here(where int) int {
	var n int

	n = first_character(where)

	if n != 0 && subkind(n) == sub_garrison {
		return n
	}

	return 0
}

func province_admin(n int) int {
	var garr int
	var castle int

	if kind(n) == T_loc {
		assert(loc_depth(n) == LOC_province)
		garr = garrison_here(n)
		if garr == 0 {
			return 0
		}
	} else {
		assert(subkind(n) == sub_garrison)
		garr = n
	}

	castle = garrison_castle(garr)

	if !valid_box(castle) {
		return 0
	}

	return building_owner(castle)
}

func top_ruler(n int) int {
	return province_admin(n)
}

//#if 0
///*
// *  is b pledged somewhere beneath a?
// */
//
//static int
//pledged_beneath(int a, int b)
//{
//
//    assert(kind(a) == T_char);
//    assert(kind(b) == T_char);
//
//    if (a == b)
//        return FALSE;
//
//    while (b > 0)
//    {
//        b = char_pledge(b);
//        if (a == b)
//            return TRUE;
//    }
//
//    return FALSE;
//}
//#endif

//#if 0
//int
//may_rule_here(who, where int)
//{
//  int pl = player(who);
//  var i int
//  int ret = FALSE;
//
//  if (is_loc_or_ship(where))
//    where = province(where);
//  else
//    assert(subkind(where) == sub_garrison);
//
//  loop_loc_owner(where, i)
//    {
//      if (player(i) == pl)
//    {
//      ret = TRUE;
//      break;
//    }
//    }
//  next_loc_owner;
//
//  return ret;
//}
//#endif

/*
 *  Fri May  7 18:34:54 1999 -- Scott Turner
 *
 *  The new rule is that you rule here if you have permission to
 *  stack with the castle ruler.
 *
 */
func may_rule_here(who, where int) int {
	return will_admit(province_admin(province(where)),
		province_admin(province(where)), who)
}

//#if 0
//ilist
//players_who_rule_here(where int)
//{
//    static var l []int
//    var i int
//    int pl;
//    int loop_check = 5000;
//
//    ilist_clear(&l);
//
//    loop_loc_owner(where, i)
//    {
//        pl = player(i);
//
//        if (pl && ilist_lookup(l, pl) < 0)
//            l = append(l, pl);
//
//        if (loop_check <= 0)
//        {
//            int j;
//            for (j = 0; j < len(l); j++)
//                fprintf(stderr, "l[%d] = %d\n", j, l[j]);
//            fprintf(stderr, "where = %d, i = %d\n", where, i);
//        }
//        assert(loop_check-- > 0);
//    }
//    next_loc_owner;
//
//    return l;
//}
//#endif

//#if 0
//void
//touch_garrison_locs()
//{
//    var garr int
//    var where int
//    var owner int
//
//    for _, garr := range loop_garrison()
//    {
//        where = subloc(garr);
//
//        loop_loc_owner(garr, owner)
//        {
//            touch_loc_pl(player(owner), where);
//        }
//        next_loc_owner;
//    }
//
//}
//#endif

func new_province_garrison(where, castle int) int {
	newChar := new_char(sub_garrison,
		0, where, -1, garr_pl, LOY_npc,
		0, "Garrison")

	if newChar < 0 {
		return -1
	}

	gen_item(newChar, item_soldier, 10)
	p_misc(newChar).cmd_allow = 'g'
	p_misc(newChar).garr_castle = castle
	p_char(newChar).guard = TRUE
	p_char(newChar).break_point = 0

	promote(newChar, 0)

	out(where, "%s now guards %s.", liner_desc(newChar), box_name(where))

	return newChar
}

func v_garrison(c *command) int {
	castle := c.a
	where := subloc(c.who)

	if loc_depth(where) != LOC_province {
		out(c.who, "Garrisons may only be installed at province level.")
		return FALSE
	}

	if garrison_here(where) != FALSE {
		out(c.who, "There is already a garrison here.")
		return FALSE
	}

	//#if 0
	//    if (first_character(where) != c.who)
	//    {
	//        out(c.who, "Must be the first unit in the location to install a garrison.");
	//        return FALSE;
	//    }
	//#endif

	if numargs(c) < 1 {
		out(c.who, "Must specify a castle to claim the province in the name of.")
		return FALSE
	}

	if subkind(castle) == sub_castle_notdone {
		out(c.who, "%s is not finished.  Garrisons may only be bound to completed castles.",
			box_name(castle))
		return FALSE
	}

	if subkind(castle) != sub_castle {
		out(c.who, "%s is not a castle.", c.parse[1])
		return FALSE
	}

	if region(castle) != region(where) {
		out(c.who, "%s is not in this region.", box_name(castle))
		return FALSE
	}

	//#if 0
	//    {
	//        out(c.who, "A garrison here must be bound to %s.",
	//                    box_name(castle));
	//        return FALSE;
	//    }
	//#endif

	if has_item(c.who, item_soldier) < 10 {
		out(c.who, "Must have %s to establish a new garrison.",
			box_name_qty(item_soldier, 10))
		return FALSE
	}

	newGarrison := new_province_garrison(where, castle)

	if newGarrison < 0 {
		out(c.who, "Failed to install garrison.")
		return FALSE
	}

	/*
	 *  Reset the tax rate for this province.
	 *
	 */
	rp_loc(where).tax_rate = 0

	consume_item(c.who, item_soldier, 10)

	out(c.who, "Installed %s", liner_desc(newGarrison))
	out(c.who, "Local tax rate set to zero percent.")

	return TRUE
}

//#if 0
//int
//v_pledge(struct command *c)
//{
//    target := c.a;
//
//    if (target == c.who)
//    {
//        wout(c.who, "Can't pledge to yourself.");
//        return FALSE;
//    }
//
//    if (target == 0)
//    {
//        p_magic(c.who).pledge = 0;
//        out(c.who, "Pledge cleared.  "
//                "Lands will be claimed for ourselves.");
//        return TRUE;
//    }
//
//    if (kind(target) != T_char)
//    {
//        out(c.who, "%s is not a character.", c.parse[1]);
//        return FALSE;
//    }
//
//    if (is_npc(target))
//    {
//        out(c.who, "May not pledge land to %s.", c.parse[1]);
//        return FALSE;
//    }
//
//    if (pledged_beneath(c.who, target))
//    {
//        wout(c.who, "Cannot pledge to %s since %s is pledged to you.",
//                box_name(target), just_name(target));
//        return FALSE;
//    }
//
//    out(c.who, "Lands are now pledged to %s.", box_name(target));
//
//    out(target, "%s pledges to us.", box_name(c.who));
//
//    p_magic(c.who).pledge = target;
//
//    return TRUE;
//}
//#endif

func nprovs_to_rank(n int) int {

	if n == 0 {
		return 0
	}
	if n <= 3 {
		return RANK_lord
	}
	if n <= 7 {
		return RANK_knight
	}
	if n <= 11 {
		return RANK_baron
	}
	if n <= 15 {
		return RANK_count
	}
	if n <= 19 {
		return RANK_earl
	}
	if n <= 22 {
		return RANK_marquess
	}

	return RANK_duke
}

func rank_s(who int) string {
	n := char_rank(who)

	switch n {
	case 0:
		return ""
	case RANK_lord:
		return ", lord"
	case RANK_knight:
		return ", knight"
	case RANK_baron:
		return ", baron"
	case RANK_count:
		return ", count"
	case RANK_earl:
		return ", earl"
	case RANK_marquess:
		return ", marquess"
	case RANK_duke:
		return ", duke"
	case RANK_king:
		return ", king"
	default:
		/* Temp fix for old ranks. */
		return ""
	}
}

func find_kings() {
	var ruler int
	var nprovs int

	for _, reg := range loop_loc() {
		if loc_depth(reg) != LOC_region {
			continue
		}

		ruler = -1
		nprovs = 0

		for _, where := range loop_here(reg) {
			if kind(where) != T_loc {
				continue
			}

			nprovs++

			if ruler == -1 {
				ruler = top_ruler(where)
				if ruler == 0 {
					break
				} /* fail */
			} else {
				if ruler != top_ruler(where) {
					ruler = 0
					break /* fail */
				}
			}
		}

		if ruler != 0 && nprovs >= 25 {
			p_char(ruler).rank = RANK_king
		}
	}

}

/*
 *  A noble's status is:
 *
 *	min(status by own provinces, lord's status - 1)
 *
 *  Fri May  7 18:11:34 1999 -- Scott Turner
 *
 *  Only by your own provinces.
 *
 */

func determine_noble_ranks() {

	stage("determine_noble_ranks()")

	clear_temps(T_player)
	clear_temps(T_char)

	/*
	 *  Cal # of controlled provinces for every castle-sitter.
	 *
	 */
	for _, garr := range loop_garrison() {
		if valid_box(province_admin(garr)) {
			bx[province_admin(garr)].temp++
		}
	}

	/*
	 *  Now find the best castle-sitter per each player.
	 *
	 */
	for _, i := range loop_player() {
		for _, j := range loop_units(i) {
			if bx[j].temp > bx[i].temp {
				bx[i].temp = bx[j].temp
			}
		}
	}

	/*
	 *  Now give each char the rank of the best castle-sitter in his faction.
	 *
	 */
	for _, who := range loop_char() {
		if char_rank(who) != 0 {
			rp_char(who).rank = 0
		}

		if bx[player(who)].temp == 0 {
			continue
		}

		p_char(who).rank = nprovs_to_rank(bx[player(who)].temp)
	}

	find_kings()
}

func garrison_notices(garr, target int) int {

	if is_npc(target) ||
		count_stack_units(target) >= 5 ||
		count_stack_figures(target) >= 20 {
		return TRUE
	}

	p := rp_misc(garr)

	if p != nil && p.garr_watch.lookup(target) >= 0 {
		return TRUE
	}

	return FALSE
}

func garrison_spot_check(garr, target int) int {
	found := FALSE

	assert(valid_box(garr))

	p := rp_misc(garr)
	if p == nil {
		return FALSE
	}

	for _, i := range loop_stack(target) {
		if p.garr_watch.lookup(i) >= 0 {
			found = TRUE
			break
		}
	}

	if found != FALSE {
		wout(garr, "Spotted in %s:", box_name(province(garr)))
	}

	return found
}

func garr_own_s(garr int) string {
	return box_code_less(province_admin(garr))
}

func garrison_summary(pl int) {
	var garr int
	var l []int
	var i int
	var taxr int
	first := TRUE

	/*
	 *  Fri May  7 18:41:08 1999 -- Scott Turner
	 *
	 *  Only the castle owner gets the garrison report.
	 *
	 */
	for _, garr = range loop_garrison() {
		if player(province_admin(garr)) != pl {
			continue
		}

		l = append(l, garr)
	}

	if len(l) == 0 {
		return
	}

	tagout(pl, "<tag type=garr_report pl=%d>", pl)
	tagout(pl, "<tag type=header>")

	out(pl, "")
	out(pl, "Garrison Report:")
	out(pl, "")
	out(pl, "%4s %5s %6s %4s %4s %4s %4s %6s %s",
		"garr", "where", "border", "men", "cost", "tax", "rate", "castle", "rulers")
	out(pl, "%4s %5s %6s %4s %4s %4s %4s %6s %s",
		"----", "-----", "------", "---", "----", "---", "----", "------", "------")
	tagout(pl, "</tag type=header>")

	sort_for_output(l)

	for i = 0; i < len(l); i++ {
		garr = l[i]

		taxr = rp_loc(subloc(garr)).tax_rate

		tagout(pl, "<tag type=garrison unit=%d loc=%d men=%d cost=%d tax=%d rate=%d castle=%d ruler=%d closed=%d>",
			garr,
			subloc(garr),
			count_stack_figures(garr),
			unit_maint_cost(garr, 0),
			rp_misc(garr).garr_tax,
			taxr,
			garrison_castle(garr),
			province_admin(garr),
			p_loc(subloc(garr)).control.closed)

		out(pl, "%4s %5s %6s %4d %4d %4d %4d %5s  %s",
			box_code_less(garr),
			box_code_less(subloc(garr)),
			or_string(rp_loc(subloc(garr)).control.closed, "closed", "open"),
			count_stack_figures(garr),
			unit_maint_cost(garr, 0),
			rp_misc(garr).garr_tax,
			taxr,
			box_code_less(garrison_castle(garr)),
			garr_own_s(garr))

		tagout(pl, "</tag type=garrison>")
	}
	tagout(pl, "</tag type=garr_report pl=%d>", pl)

	out(pl, "")
	for i = 0; i < len(l); i++ {
		/* Output inventories for any garrisons w/ inventory. */
		first = TRUE
		garr = l[i]
		if len(bx[garr].items) > 0 {
			sort.Sort(bx[garr].items)
			for _, e := range loop_inventory(garr) {
				if first != 0 {
					tagout(pl, "<tag type=garr_inv id=%d>", garr)
					out(pl, "%s in %s<tag type=tab col=48>weight",
						box_name(garr),
						box_name(province(garr)))
					out(pl, "-----------------------------------<tag type=tab col=48>------")
					first = FALSE
				}
				tagout(pl, "<tag type=inventory unit=%d item=%d qty=%d weight=%d extra=\"%s\">",
					garr, e.item, e.qty, item_weight(e.item)*e.qty,
					extra_item_info(garr, e.item, e.qty))
				tagout(pl, "%9s  %-30s <tag type=tab col=45>%9s  %s",
					comma_num(e.qty),
					plural_item_box(e.item, e.qty),
					comma_num(item_weight(e.item)*e.qty),
					extra_item_info(i, e.item, e.qty))
				tagout(pl, "</tag type=inventory>")
			}
		}
		if first == 0 {
			out(pl, "")
			tagout(pl, "</tag type=garr_inv>")
		}
	}
}

func v_decree_watch(c *command) int {
	target := c.a
	ncontrol := 0
	nordered := 0

	if kind(target) != T_char {
		wout(c.who, "%s is not a character.", c.parse[1])
		return FALSE
	}

	for _, garr := range loop_garrison() {
		if c.who != province_admin(garr) {
			continue
		}

		ncontrol++
		p := p_misc(garr)

		if len(p.garr_watch) < 3 {
			p.garr_watch = append(p.garr_watch, target)
			wout(garr, "%s orders us to watch for %s.",
				box_name(c.who), box_code(target))

			nordered++
		}
	}

	if ncontrol == 0 {
		wout(c.who, "We rule over no garrisons.")
		return FALSE
	}

	if nordered == 0 {
		wout(c.who, "Garrisons may only watch for up to three units per month.")
		return FALSE
	}

	wout(c.who, "Watch decree given to %s garrison%s.",
		nice_num(nordered), add_s(nordered))

	return TRUE
}

func v_decree_hostile(c *command) int {
	var garr int
	target := c.a
	ncontrol := 0
	nordered := 0

	if kind(target) != T_char {
		wout(c.who, "%s is not a character.", c.parse[1])
		return FALSE
	}

	for _, garr = range loop_garrison() {
		if c.who != province_admin(garr) {
			continue
		}

		ncontrol++
		p := p_misc(garr)

		if len(p.garr_host) < 3 {
			p.garr_host = append(p.garr_host, target)
			wout(garr, "%s orders us to attack %s on sight.",
				box_name(c.who), box_code(target))

			nordered++
		}
	}

	if ncontrol == 0 {
		wout(c.who, "We rule over no garrisons.")
		return FALSE
	}

	if nordered == 0 {
		wout(c.who, "Garrisons may be hostile to at most three units.")
		return FALSE
	}

	wout(c.who, "Hostile decree given to %s garrison%s.",
		nice_num(nordered), add_s(nordered))

	return TRUE
}

var decree_tags = []string{
	"watch",   /* 0 */
	"hostile", /* 1 */
	""}

func v_decree(c *command) int {
	var tag int

	if numargs(c) < 1 {
		wout(c.who, "Must specify what to decree.")
		return FALSE
	}

	tag = lookup_sb(decree_tags, c.parse[1])

	if tag < 0 {
		wout(c.who, "Unknown decree '%s'.", c.parse[1])
		return FALSE
	}

	cmd_shift(c)

	switch tag {
	case 0:
		return v_decree_watch(c)
	case 1:
		return v_decree_hostile(c)
	}
	panic("!reached")
}

func ping_garrisons() {
	var garr int
	var where int

	show_to_garrison = true

	for _, garr = range loop_garrison() {
		where = subloc(garr)

		assert(rp_char(garr) != nil)

		rp_char(garr).guard = FALSE

		wout(where, "%s guards %s.",
			liner_desc(garr), box_name(where))

		rp_char(garr).guard = TRUE
	}

	show_to_garrison = false
}

func v_ungarrison(c *command) int {
	garr := c.a
	where := subloc(c.who)
	first := TRUE

	if garr == 0 {
		garr = garrison_here(where)

		if garr == 0 {
			wout(c.who, "There is no garrison here.")
			return FALSE
		}
	} else if garrison_here(where) != garr {
		wout(c.who, "No garrison %s is here.", c.parse[1])
		return FALSE
	}

	if player(c.who) != player(province_admin(garr)) {
		wout(c.who, "%s does not rule over %s.",
			box_name(c.who), box_name(garr))
		return FALSE
	}

	wout(c.who, "%s disbands.", box_name(garr))

	vector_clear()
	vector_add(garr)
	vector_add(where)
	wout(VECT, "%s is disbanded by %s.", box_name(garr), box_name(c.who))

	for _, e := range loop_inventory(garr) {
		if first != 0 {
			first = FALSE
			wout(c.who, "Received from %s:", box_name(garr))
			indent += 3
		}

		wout(c.who, "%s", box_name_qty(e.item, e.qty))

		move_item(garr, c.who, e.item, e.qty)
	}

	if first == 0 {
		indent -= 3
	}

	p_misc(garr).garr_castle = 0 /* become silent */
	kill_char(garr, 0, S_nothing)

	return TRUE
}
