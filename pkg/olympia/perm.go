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
	"strings"
)

type admit_l []*admit

func (a admit_l) Len() int {
	return len(a)
}

func (a admit_l) Less(i, j int) bool {
	return a[i].targ < a[j].targ
}

func (a admit_l) Swap(i, j int) {
	a[i].targ, a[j].targ = a[j].targ, a[i].targ
}

func rp_admit(pl, targ int) *admit {
	assert(kind(pl) == T_player)

	p := p_player(pl)
	for i := 0; i < len(p.admits); i++ {
		if p.admits[i].targ == targ {
			return p.admits[i]
		}
	}
	return nil
}

func p_admit(pl, targ int) *admit {
	assert(kind(pl) == T_player)

	p := p_player(pl)
	for i := 0; i < len(p.admits); i++ {
		if p.admits[i].targ == targ {
			return p.admits[i]
		}
	}
	a := &admit{targ: targ}
	p.admits = append(p.admits, a)
	return a
}

/*
 *  Will pl admit who into targ?
 */

func will_admit(pl, who, targ int) int {
	/*
	 *  Fri Nov  5 13:02:00 1999 -- Scott Turner
	 *
	 *  For purposes of admission, a garrison is treated as if it were
	 *  the ruler of the castle, e.g., the garrison will admit you if the
	 *  ruler of the castle would admit you.  This is a little odd, perhaps,
	 *  but that's the way the rules are currently written.
	 *
	 */
	if subkind(targ) == sub_garrison {
		targ = province_admin(targ)
		pl = targ
		if !valid_box(targ) {
			return FALSE
		}
	}

	pl = player(pl)

	if player(who) == pl {
		return TRUE
	}

	p := rp_admit(pl, targ)

	if p == nil {
		return FALSE
	}

	found := ilist_lookup(p.l, who) >= 0
	found_pl := ilist_lookup(p.l, player(who)) >= 0
	found_nation := ilist_lookup(p.l, nation(who)) >= 0

	/*
	 * Wed Jan 20 12:59:51 1999 -- Scott Turner
	 *
	 * If p.sense is true, then we have unit and player
	 * exclusion, e.g., if the unit or player is true, then
	 * don't admit them!.
	 *
	 */
	if p.sense != 0 {
		if found || found_pl || found_nation {
			return FALSE
		}
		return TRUE
	} else {
		if found || found_pl || found_nation {
			return TRUE
		}
		return FALSE
	}
}

/*
 *  Wed Jan 20 12:23:16 1999 -- Scott Turner
 *
 *  Add nation admits.
 *
 */
func v_admit(c *command) int {
	targ := c.a
	if !valid_box(targ) {
		wout(c.who, "Must specify an entity for admit.")
		return FALSE
	}
	pl := player(c.who)
	p := p_admit(pl, targ)

	cmd_shift(c)
	if numargs(c) == 0 {
		p.sense = FALSE
		p.l = nil
	}

	for numargs(c) > 0 {
		parse_s := string(c.parse[1])
		if strings.ToLower(parse_s) == "all" {
			p.sense = TRUE
		} else if nat := find_nation(parse_s); nat != 0 {
			/*
			 *  We can stick the nation # on there because we
			 *  can't have a box number that low (hopefully!).
			 *
			 */
			p.l = ilist_add(p.l, nat)
			wout(c.who, "Admitting '%s' to %s.", rp_nation(nat).name, box_code_less(targ))
		} else if kind(c.a) == T_char || kind(c.a) == T_player || kind(c.a) == T_unform {
			p.l = ilist_add(p.l, c.a)
		} else {
			wout(c.who, "%s isn't a valid entity to admit.", c.parse[1])
		}
		cmd_shift(c)
	}

	return TRUE
}

func print_admit_sup(pl int, p *admit) {
	count := 0

	buf := fmt.Sprintf("admit %4s", box_code_less(p.targ))

	if p.sense != 0 {
		buf += "  all"
		count++
	}

	for i := 0; i < len(p.l); i++ {
		if !valid_box(p.l[i]) {
			continue
		}
		if count = count + 1; count >= 12 {
			out(pl, "%s", buf)
			//buf = fmt.Sprintf("admit %4s", p.targ);
			buf += "          "
			count = 1
		}
		if kind(p.l[i]) == T_nation {
			buf += sout(" %s", rp_nation(p.l[i]).name)
		} else {
			buf += sout(" %4s", box_code_less(p.l[i]))
		}
	}

	if count != 0 {
		out(pl, "%s", buf)
	}
}

func print_admit(pl int) {
	first := TRUE

	assert(kind(pl) == T_player)

	p := p_player(pl)

	if len(p.admits) > 0 {
		sort.Sort(p.admits)
	}

	for i := 0; i < len(p.admits); i++ {
		if valid_box(p.admits[i].targ) {
			if first != FALSE {
				tagout(pl, "<tag type=header>")
				out(pl, "")
				tagout(pl, "</tag type=header>")
				out(pl, "Admit permissions:")
				out(pl, "")
				indent += 3
				first = FALSE
			}

			print_admit_sup(pl, p.admits[i])
		}
	}

	if first == FALSE {
		indent -= 3
	}
}

func clear_all_att(who int) {
	p := rp_disp(who)
	if p == nil {
		return
	}
	p.neutral = nil
	p.hostile = nil
	p.defend = nil
}

func clear_att(who int, disp int) {
	p := rp_disp(who)
	if p == nil {
		return
	}
	switch disp {
	case NEUTRAL:
		p.neutral = nil
	case HOSTILE:
		p.hostile = nil
	case DEFEND:
		p.defend = nil
	case ATT_NONE:
	default:
		assert(false)
	}
}

func set_att(who int, targ int, disp int) {
	p := p_disp(who)

	p.neutral = rem_value(p.neutral, targ)
	p.hostile = rem_value(p.hostile, targ)
	p.defend = rem_value(p.defend, targ)

	switch disp {
	case NEUTRAL:
		p.neutral = append(p.neutral, targ)
		sort.Ints(p.neutral)
	case HOSTILE:
		p.hostile = append(p.hostile, targ)
		sort.Ints(p.hostile)
	case DEFEND:
		p.defend = append(p.defend, targ)
		sort.Ints(p.defend)
	case ATT_NONE:
	default:
		assert(false)
	}
}

/*
  - Mon May 18 19:07:03 1998 -- Scott Turner
    *
  - Macro doesn't work because of conceal_nation_ef...
    *
    #define nation(n)	(n && player(n) && rp_player(player(n)) ?
    rp_player(player(n)).nation : 0)
    *
*/
func nation(who int) int {
	/*
	 *  Sanity checks.
	 *
	 */
	if !valid_box(who) {
		return 0
	}
	/*
	 *  Return the phony nation, if any!
	 *
	 */
	n := get_effect(who, ef_conceal_nation, 0, 0)
	if n != 0 {
		assert(kind(n) == T_nation)
		return n
	}
	/*
	 *  A garrison ought to be considered to be of the nation
	 *  of its lord.
	 *
	 */
	if subkind(who) == sub_garrison {
		if ruler := province_admin(who); ruler != 0 && rp_player(player(ruler)) != nil {
			return rp_player(player(ruler)).nation
		}
	}
	/*
	 *  A deserted noble ought to be considered still of the nation
	 *  of his old lord, if he has one.
	 *
	 */
	pl := player(who)
	if is_real_npc(pl) && body_old_lord(who) != 0 && rp_player(player(body_old_lord(who))) != nil {
		return rp_player(player(body_old_lord(who))).nation
	}
	/*
	 *  Otherwise...
	 *
	 */
	if pl != 0 && rp_player(pl) != nil {
		return rp_player(pl).nation
	}

	return 0
}

/*
 *  Try to find a nation.
 *
 */
func find_nation(name string) int {
	for _, i := range loop_nation() {
		if fuzzy_strcmp([]byte(rp_nation(i).name), []byte(name)) || strings.HasPrefix(strings.ToLower(rp_nation(i).name), strings.ToLower(name)) {
			return i
		}
	}
	return 0
}

func find_nation_b(b []byte) int {
	return find_nation(string(b))
}

/*
 *  Tue Jan 12 12:11:32 1999 -- Scott Turner
 *
 *  Added support for hostile to monsters.
 *
 */
func is_hostile(who int, targ int) int {
	if player(who) == player(targ) {
		return FALSE
	}

	if subkind(who) == sub_garrison {
		if p := rp_misc(who); p != nil && ilist_lookup(p.garr_host, targ) >= 0 {
			return TRUE
		}
	}

	if p := rp_disp(who); p != nil {
		if ilist_lookup(p.hostile, targ) >= 0 {
			return TRUE
		}
		/*
		 *  Mon May 18 19:04:22 1998 -- Scott Turner
		 *
		 *  Might be a nation...
		 *
		 */
		if nation(targ) != 0 && ilist_lookup(p.hostile, nation(targ)) >= 0 {
			return TRUE
		}
		/*
		 *  Tue Jan 12 12:09:53 1999 -- Scott Turner
		 *
		 *  Might be a "monster"
		 *
		 */
		if !is_real_npc(who) &&
			is_real_npc(targ) &&
			kind(targ) == T_char &&
			subkind(targ) == sub_ni &&
			ilist_lookup(p.hostile, MONSTER_ATT) >= 0 {
			return TRUE
		}
	}

	if p := rp_disp(player(who)); p != nil {
		if ilist_lookup(p.hostile, targ) >= 0 {
			return TRUE
		}
		/*
		 *  Mon May 18 19:04:22 1998 -- Scott Turner
		 *
		 *  Might be a nation...
		 *
		 */
		if nation(targ) != 0 && ilist_lookup(p.hostile, nation(targ)) >= 0 {
			return TRUE
		}
		/*
		 *  Tue Jan 12 12:09:53 1999 -- Scott Turner
		 *
		 *  Might be a "monster"
		 *
		 */
		if !is_real_npc(who) &&
			is_real_npc(targ) &&
			kind(targ) == T_char &&
			subkind(targ) == sub_ni &&
			ilist_lookup(p.hostile, MONSTER_ATT) >= 0 {
			return TRUE
		}
	}
	return FALSE
}

func is_defend(who int, targ int) int {
	/*
	 *  Mon Mar  3 13:24:58 1997 -- Scott Turner
	 *
	 *  All npcs defend each other!
	 *
	 *  Sun Mar  9 20:57:06 1997 -- Scott Turner
	 *
	 *  A little simplistic.  But we should have all intelligent
	 *  NPCs defend each other, and all animals of the same type.
	 *
	 */
	if is_real_npc(who) && is_real_npc(targ) &&
		npc_program(who) != 0 &&
		npc_program(who) != PROG_dumb_monster &&
		npc_program(targ) == npc_program(who) {
		wout(who, "Smart enough to help %s in battle.", box_name(targ))
		return TRUE
	}

	if is_real_npc(who) && is_real_npc(targ) &&
		subkind(who) == sub_ni &&
		subkind(targ) == sub_ni &&
		noble_item(who) == noble_item(targ) {
		wout(who, "Rushing to the defense of similar beast %s.", box_name(targ))
		return TRUE
	}

	if is_hostile(who, targ) != FALSE {
		return FALSE
	}

	if p := rp_disp(who); p != nil {
		if ilist_lookup(p.defend, targ) >= 0 {
			return TRUE
		}
		if ilist_lookup(p.neutral, targ) >= 0 {
			return FALSE
		}

		if ilist_lookup(p.defend, player(targ)) >= 0 {
			return TRUE
		}
		if ilist_lookup(p.neutral, player(targ)) >= 0 {
			return FALSE
		}
		/*
		 *  Mon May 18 19:04:22 1998 -- Scott Turner
		 *
		 *  Might be a nation...
		 *
		 */
		if nation(targ) != 0 && ilist_lookup(p.defend, nation(targ)) >= 0 {
			return TRUE
		}
		if nation(targ) != 0 && ilist_lookup(p.neutral, nation(targ)) >= 0 {
			return FALSE
		}
	}

	pl := player(who)
	if p := rp_disp(pl); p != nil {
		if ilist_lookup(p.defend, targ) >= 0 {
			return TRUE
		}
		if ilist_lookup(p.neutral, targ) >= 0 {
			return FALSE
		}

		if ilist_lookup(p.defend, player(targ)) >= 0 {
			return TRUE
		}
		if ilist_lookup(p.neutral, player(targ)) >= 0 {
			return FALSE
		}

		/*
		 *  Mon May 18 19:04:22 1998 -- Scott Turner
		 *
		 *  Might be a nation...
		 *
		 */
		if nation(targ) != 0 && ilist_lookup(p.defend, nation(targ)) >= 0 {
			return TRUE
		}
		if nation(targ) != 0 && ilist_lookup(p.neutral, nation(targ)) >= 0 {
			return FALSE
		}
	}

	if pl == player(targ) && pl != indep_player {
		if cloak_lord(who) != FALSE {
			return FALSE
		}
		return TRUE
	}

	return FALSE
}

/*
 *  Mon May 18 18:47:41 1998 -- Scott Turner
 *
 *  Accept nation names as well.
 *
 *  Tue Jan 12 11:58:09 1999 -- Scott Turner
 *
 *  Accept "monster" as well?
 *
 */
var verbs_perm = []string{
	"no attitude",
	"neutral",
	"hostile",
	"defend"}

func v_set_att(c *command, k int) int {
	var n int

	if numargs(c) == 0 {
		/*
		 *  Clear a list.
		 *
		 */
		wout(c.who, "Cleared %s list.", verbs_perm[k])
		clear_att(c.who, k)
		return TRUE
	}

	for numargs(c) > 0 {
		if !valid_box(c.a) {
			/*
			 *  Look for a nation name.
			 *
			 */
			n = find_nation(string(c.parse[1]))
			if n != 0 {
				set_att(c.who, n, k)
				wout(c.who, "Declared %s toward nation %s.", verbs_perm[k], rp_nation(n).name)
			} else {
				/*
				 *  Might be "monster" or "monsters"
				 *
				 */
				if fuzzy_strcmp(c.parse[1], []byte("monster")) || fuzzy_strcmp(c.parse[1], []byte("monsters")) {
					set_att(c.who, MONSTER_ATT, k)
				} else {
					wout(c.who, "%s is not a valid entity.", c.parse[1])
				}
			}
		} else if k == HOSTILE && player(c.who) == player(c.a) &&
			player(c.who) != indep_player {
			wout(c.who, "Can't be hostile to a unit in the same faction.")
		} else {
			set_att(c.who, c.a, k)
			wout(c.who, "Declared %s towards %s.", verbs_perm[k], box_code(c.a))
		}
		cmd_shift(c)
	}
	return TRUE
}

func v_hostile(c *command) int {
	return v_set_att(c, HOSTILE)
}

func v_defend(c *command) int {
	return v_set_att(c, DEFEND)
}

func v_neutral(c *command) int {
	return v_set_att(c, NEUTRAL)
}

func v_att_clear(c *command) int {
	return v_set_att(c, ATT_NONE)
}

func print_att_sup(who int, l []int, header string, first *int) {
	if len(l) == 0 {
		return
	}
	sort.Ints(l)

	buf := header
	count := 0

	for i := 0; i < len(l); i++ {
		if l[i] != MONSTER_ATT && !valid_box(l[i]) {
			continue
		}

		if *first != FALSE {
			out(who, "")
			out(who, "Declared attitudes:")
			out(who, "")
			indent += 3
			*first = FALSE
		}

		if count = count + 1; count >= 12 {
			out(who, "%s", buf)
			buf = string(spaces[:len(header)])
			count = 1
		}

		if l[i] == MONSTER_ATT {
			buf += " Monsters "
		} else if kind(l[i]) == T_nation {
			buf += sout(" %s", rp_nation(l[i]).name)
		} else {
			buf += sout(" %4s", box_code_less(l[i]))
		}
	}
	if count != 0 {
		out(who, "%s", buf)
	}
}

func print_att(who int, n int) {
	first := TRUE
	var p *att_ent

	p = rp_disp(n)

	if p == nil {
		return
	}

	print_att_sup(who, p.hostile, "hostile", &first)
	print_att_sup(who, p.neutral, "neutral", &first)
	print_att_sup(who, p.defend, "defend ", &first)

	if first == FALSE {
		indent -= 3
	}
}
