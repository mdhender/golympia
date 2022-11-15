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
	"log"
	"os"
	"reflect"
)

const (
	FK_fort  = -1 /* structure */
	FK_noble = -2

	/*
	 *  Combat Parts
	 *  Fri Nov  8 15:00:45 1996 -- Scott Turner
	 */
	SPECIAL = 1
	MISSILE = 2
	MELEE   = 3
	UNCANNY = 4 /* Against back row archers, only */

	NO_ATTACKS = -1 /* You didn't have any attacks in this phase */

	DEFEAT_DEFAULT = 0
	DEFEAT_FORCE   = 1
	DEFEAT_UNREADY = 2
)

var (
	combat_def_loc        = 0      /* Where the defenders are. */
	combat_rain           = false  /* bad for archers */
	combat_sea            = false  /* naval combat */
	combat_swampy         = false  /* bad for horses */
	combat_wind           = false  /* poor for archers */
	defense_side          []*fight /* The side "defending" */
	round                 int
	second_wait_list      []int
	special_attack_banner = false
)

type fight struct {
	unit    int /* what unit are we in */
	kind    int /* item type or FK_xxx */
	sav_num int /* original health or count */
	num     int /* num of item, health or fort rating */
	/* (will be 1 for nobles, since they */
	/* may only take 1 hit) */
	/* (damage for locs is stored in num) */
	/* Where hits are collected until after the round */
	tmp_num int
	/* Where hits are collected until after the battle */
	tmp_num2 int

	/*
	 *  New_Health is health after the battle for a noble.  This
	 *  is necessary because a unit drops out of battle when g.num
	 *  reaches zero.  So we can't track a noble's health in g.num
	 *  because then he won't drop out of the battle until he dies.
	 *  So num holds the psuedo-health (1 or 0, depending upon whether
	 *  the noble has personally broke yet) while new_health holds
	 *  his actual health...
	 *
	 */
	new_health     int
	attack         int
	defense        int
	behind         int  /* how behind we are, 0=front */
	missile        int  /* missile attack (0=can't fight from behind)*/
	protects       int  /* fighter struct we stand in front of */
	nprot          int  /* number of fighter structs protecting us */
	ally           int  /* unit pulled into fight by ally order */
	prisoner       int  /* unit lost and is candidate for prisoner */
	inside         bool /* fighter is inside a structure */ // mdhender: was int, hope this isn't location
	seize_slot     int  /* take loser's place if we win */
	survived_fatal int  /* survived a fatal wound with [9101] */
	attack_item    int  /* bonus item this character is wielding */
	defense_item   int
	missile_item   int
}

type wield struct {
	attack  int
	defense int
	missile int
}

func allied(a, b int) bool {
	return is_defend(a, b) != FALSE
}

// COMBAT_COMP
// Thu Jun 27 12:14:24 1996 -- Scott Turner
//
// Compare two inventory entries and determine which one has the higher
// attack+defense values, or in the case of a tie, the higher attack value.
// If a non-fighter, it always loses.
// todo: that "if non-fighter" does not seem to be implemented?
func combat_comp(a, b *item_ent) int {
	a_attack := item_attack(a.item)
	a_defense := item_defense(a.item)
	a_missile := item_missile(a.item)

	b_attack := item_attack(b.item)
	b_defense := item_defense(b.item)
	b_missile := item_missile(b.item)

	if a_attack+a_defense+a_missile == b_attack+b_defense+b_missile {
		return (b_attack + b_missile) - (a_attack + a_missile)
	}

	return ((b_attack + b_missile + b_defense) - (a_attack + a_missile + a_defense))
}

/*
 *  Thu Dec  3 08:25:50 1998 -- Scott Turner
 *
 *  Infrastructure for printing the special attack round banner
 *  only if an attack occurs.
 *
 */

func print_special_banner() {
	if !special_attack_banner {
		wout(VECT, "  Special attacks phase:")
		special_attack_banner = true
	}
}

func cannot_take_prisoners(who int) bool {
	if only_defeatable(who) != FALSE {
		return true
	}

	if subkind(who) == sub_garrison {
		return true
	}

	//#if 0
	//    if (is_npc(who))
	//        return TRUE;
	//#endif

	return false
}

func cannot_take_booty(who int) bool {

	//#if 0
	//    if (is_npc(who))
	//        return TRUE;
	//#endif

	if subkind(who) == sub_garrison {
		return true
	}

	return false
}

func v_behind(c *command) int {
	num := c.a
	if num < 0 {
		num = 0
	}
	if num > 9 {
		num = 9
	}

	p_char(c.who).behind = num

	var s string
	if num == 0 {
		s = " (front unit)"
	}

	wout(c.who, "Behind flag set to %d%s.", num, s)

	return TRUE
}

func execute_prisoner(who, pris int) {

	vector_stack(who, true)
	wout(VECT, "%s executes %s!", box_name(who), box_name(pris))

	if survive_fatal(pris) {
		wout(VECT, "%s miraculously is still alive!", box_name(pris))
		prisoner_escapes(pris)
	} else {
		kill_char(pris, who, S_body)
	}
}

func v_execute(c *command) int {
	pris := c.a
	first := true

	if numargs(c) < 1 {
		for _, pris := range loop_here(c.who) {
			if is_prisoner(pris) {
				execute_prisoner(c.who, pris)
				first = false
			}
		}

		if first {
			out(c.who, "No prisoners to execute.")
		}
		return TRUE
	}

	if FALSE == has_prisoner(c.who, pris) {
		wout(c.who, "Don't have a prisoner %s.", box_code(pris))
		return FALSE
	}

	execute_prisoner(c.who, pris)
	return TRUE
}

func fort_covers(n int) int {
	switch subkind(n) {
	case sub_castle:
		return 500
	case sub_castle_notdone:
		return 500
	case sub_tower:
		return 100
	case sub_tower_notdone:
		return 100
	case sub_galley:
		return 50
	case sub_galley_notdone:
		return 50
	case sub_roundship:
		return 50
	case sub_roundship_notdone:
		return 50
	case sub_temple:
		return 50
	case sub_temple_notdone:
		return 50
	case sub_inn:
		return 50
	case sub_inn_notdone:
		return 50
	case sub_mine:
		return 50
	case sub_mine_notdone:
		return 50
	case sub_guild:
		return 50
	case sub_orc_stronghold, sub_orc_stronghold_notdone:
		return 100
	case sub_ship:
		/*
		 *  Fri Jan  3 10:48:09 1997 -- Scott Turner
		 *
		 *  Depends on the # of fortifications.
		 *
		 */
		s := rp_ship(n)
		if s == nil {
			return 0
		}

		// todo: consider using lround() from math.h
		return int(float64(s.forts)/float64(s.hulls)*SHIP_FORTS_PROTECT + 0.5)

	default:
		fprintf(os.Stderr, "subkind is %s\n", subkind_s[subkind(n)])
		return 0
	}
}

func is_siege_engine(item int) bool {

	switch item {
	case item_battering_ram, item_catapult, item_siege_tower:
		return true
	}

	return false
}

func siege_engine_useful(l []*fight) bool {
	assert(len(l) > 0)
	if l[0].kind == FK_fort && l[0].num > 0 {
		return true
	}
	return false
}

func lead_char_pos(l []*fight) int {
	assert(len(l) > 0)
	if l[0].kind == FK_noble {
		return 0
	}
	assert(len(l) > 1)
	if l[1].kind == FK_noble {
		return 1
	}
	panic("!reached")
}

func lead_char(l []*fight) int {
	return l[lead_char_pos(l)].unit
}

/*
 *  Assumptions:
 *
 *	l[0] will be a fort, if there is one
 *	l[1] will be the lead noble in this case.
 *	Otherwise, l[0] is the lead noble.
 *
 *	Fighters inside a structure (those with .inside set) must
 *	be contiguous in the array, and start from the beginning.
 *
 *	(Defenders not in the structure will only be allies from the
 *	outside, which are added after the primary defenders).
 */

func dump_fighters(l []*fight) {
	out(combat_pl, "side:  %s", box_name(lead_char(l)))

	for i := 0; i < len(l); i++ {
		s := sout("bh=%d ms=%d ins=%d pris=%d sav=%d at=%d df=%d",
			l[i].behind, l[i].missile, l[i].inside,
			l[i].prisoner, l[i].sav_num, l[i].attack,
			l[i].defense)

		out(combat_pl, "   %s.%d n=%d prt=%d nprt=%d al=%d %s",
			box_code_less(l[i].unit), l[i].kind, l[i].num,
			l[i].protects, l[i].nprot, l[i].ally, s)
	}

	out(combat_pl, "")
}

/*
 *  Backup from a fighter slot until the previous noble
 *  they are part of is found.
 */

func who_protects(l []*fight, pos int) int {
	i := pos
	for i >= 0 && l[i].kind != FK_noble {
		i--
	}
	assert(i >= 0)
	assert(l[i].kind == FK_noble)
	assert(l[i].unit == l[pos].unit)
	return i
}

func init_prot(l []*fight) []*fight {
	for i := 0; i < len(l); i++ {
		if l[i].num > 0 {
			switch l[i].kind {
			case FK_fort:
				l[i].protects = -1
				break

			case FK_noble:
				l[i].protects = lead_char_pos(l)

				/*
				 *  Don't set protect field to point to ourself
				 */

				if l[i].protects == i {
					l[i].protects = -1
				}

				//#if 0
				//                    l[i].protects = find_unit(l, stack_parent(l[i].unit));
				//#endif
				break

			default:

				/*
				 *  Siege engines shouldn't count to protect the noble.
				 *  The attacker would rather claim them as booty than destroy them, too.
				 */

				if is_siege_engine(l[i].kind) {
					l[i].protects = -1
				} else {
					l[i].protects = who_protects(l, i)
				}

				break
			}
		}
	}

	for i := 0; i < len(l); i++ {
		l[i].nprot = 0
	}

	for i := 0; i < len(l); i++ {
		if l[i].protects >= 0 && l[i].num != 0 {
			l[l[i].protects].nprot++
		}
	}

	return l
}

//#if 0
///*
// *  Determine what the character is wielding and wearing, if
// *  anything.  There may be up to three items; an attack weapon,
// *  a missile weapon, and some sort of defensive garment.
// *
// *  If f is not nil, add in the bonuses to a fight struct.
// *
// *  Returns TRUE if the character is wearing or wielding something.
// *
// *  Mon Jun 28 12:12:07 1999 -- Scott Turner
// *
// *  Whoops, needs to be updated for new combat artifacts!
// *
// */
//
//int
//find_wield(struct wield *w, who int, f []*fight)
//{
//  struct item_ent *e;
//  int n;
//  attack_max := -1;
//  defense_max := -1;
//  missile_max := -1;
//  struct wield _w;
//  struct entity_artifact *a;
//  int item;
//
//  if (w == nil)
//    w = &_w;
//
//  if (combat_artifact_bonus(who, CA_N_MELEE, &item)) {
//    w.attack = item;
//  }
//  if (art = best_artifact(who, ART_IMPRV_DEF, f.kind, 0)) {
//    defense += (100 * rp_item_artifact(art).param1);
//  }
//
//  if (combat_artifact_bonus(who, CA_N_, &item)) {
//    w.defense = item;
//  }
//  if (combat_artifact_bonus(who, CA_N_MISSILE, &item)) {
//    w.missile = item;
//  }
//
//
////#if 0
////      if ((n = item_attack_bonus(e.item)) && (n > attack_max))
////    {
////      attack_max = n;
////      w.attack = e.item;
////    }
////
////      if ((n = item_defense_bonus(e.item)) && (n > defense_max))
////    {
////      defense_max = n;
////      w.defense = e.item;
////    }
////
////      if ((n = item_missile_bonus(e.item)) && (n > missile_max))
////    {
////      missile_max = n;
////      w.missile = e.item;
////    }
////#endif
//
//
//  if (f) {
//    if (w.attack)
//      f.attack += attack_max;
//
//    if (w.defense)
//      f.defense += defense_max;
//
//    if (w.missile)
//      f.missile += missile_max;
//  }
//
//  if (w.attack || w.defense || w.missile)
//    return TRUE;
//
//  return FALSE;
//}
//
//
//char *
//wield_s(who int)
//{
//    static char buf[LEN];
//    struct wield w;
//
//    if (FALSE == find_wield(&w, who, nil))
//        return "";
//
//    *buf = '\0';
//
///*
// *  Clear out multiple copies of the same item.
// *  This would happen if one weapon had multiple bonuses.
// *  We don't want to say "Wielding foo and foo, wearing foo."
// */
//
//    if (w.attack == w.missile)
//        w.missile = 0;
//
//    if (w.attack == w.defense)
//        w.defense = 0;
//
//    if (w.attack || w.missile)
//    {
//        if (w.attack == 0)
//            strcpy(buf, sout(", wielding %s",
//                box_name(w.missile)));
//        else if (w.missile == 0)
//            strcpy(buf, sout(", wielding %s",
//                box_name(w.attack)));
//        else
//            strcpy(buf, sout(", wielding %s and %s",
//                box_name(w.attack),
//                box_name(w.missile)));
//    }
//
//    if (w.defense)
//        strcat(buf, sout(", wearing %s", box_name(w.defense)));
//
//    return buf;
//}
//#endif

func init_attack_defense(l []*fight) []*fight {
	for i := 0; i < len(l); i++ {
		f := l[i]

		switch f.kind {
		case FK_fort:
			f.attack = 0
			f.missile = 0
			f.behind = 0
			f.defense = loc_defense(f.unit)
			/*
			 *  Wed Oct 16 12:44:40 1996 -- Scott Turner
			 *
			 *  This fort may be "improved".
			 *
			 */
			f.defense += get_effect(f.unit, ef_improve_fort, 0, 0)
			break

		case FK_noble:
			mk := noble_item(f.unit)
			if mk == 0 {
				f.attack = char_attack(f.unit)
				f.defense = char_defense(f.unit)
				f.missile = char_missile(f.unit)
			} else {
				f.attack = item_attack(mk)
				f.defense = item_defense(mk)
				f.missile = item_missile(mk)
			}

			f.behind = char_behind(f.unit)
			/*
			 *  Add in combat bonuses from items the noble is carrying
			 *
			 *  Mon Jun 28 12:20:30 1999 -- Scott Turner
			 *
			 *  Not used anymore, since these sorts of artifacts do not
			 *  exist.
			 *
			 */
			//#if 0
			//                find_wield(nil, f.unit, f);
			//#endif

			break

		default:
			f.attack = item_attack(f.kind)
			f.defense = item_defense(f.kind)
			f.missile = item_missile(f.kind)
			f.behind = char_behind(f.unit)
			if (combat_swampy || combat_sea) &&
				(f.kind == item_elite_guard ||
					f.kind == item_cavalier ||
					f.kind == item_knight) {
				f.attack -= 50
				f.defense -= 50
				assert(f.attack > 0)
				assert(f.defense > 0)
			}

			if combat_wind {
				f.missile /= 2
			}

			if combat_rain &&
				(f.kind == item_archer ||
					f.kind == item_horse_archer ||
					f.kind == item_elite_arch) {
				f.missile = 0
			}

			if f.kind == item_pirate {
				where := subloc(f.unit)

				if is_ship(where) || is_ship_notdone(where) {
					f.attack *= 3
					f.defense *= 3
				}
			}

			break
		}
	}

	return l
}

func add_to_fight_list(l *[]*fight, unit int, kind int, num int, ally int, inside bool) {
	/*
	 *  Siege engines not engaged for naval combat
	 */

	if combat_sea && is_siege_engine(kind) {
		return
	}

	/*
	 *  Maybe none?
	 *
	 */
	if num == 0 {
		return
	}

	newt := &fight{}
	newt.unit = unit
	newt.kind = kind
	newt.sav_num = num
	newt.new_health = num

	if kind == FK_noble {
		newt.num = 1
	} else {
		newt.num = num
	}

	newt.ally = ally
	newt.inside = inside

	*l = append(*l, newt)
}

/*
 *  Determine how many men this entity (who) can control in
 *  battle.
 *
 *  (1) Nobles can control equal to their sk_control_battle
 *      skill, or 10, whichever is greater.
 *
 *  (2) All others can control a virtually unlimited amount.
 *
 *  Mon Mar  1 11:49:20 1999 -- Scott Turner
 *
 *  Whoops, we need to limit the number of men a controlled monster
 *  can control (do you follow that?).
 *
 *  Tue May 30 17:53:06 2000 -- Scott Turner
 *
 *  Need to add in the bonus for being a castle-owner.
 *
 */
func calc_man_limit(who, defender int) int {
	man_limit := 0

	if !is_real_npc(who) {
		/*
		 *  Two types for player-controlled units: nobles and controlled
		 *  monster stacks.
		 *
		 */
		if subkind(who) != 0 {
			man_limit = GARRISON_CONTROLLED
		} else {
			p := rp_skill_ent(who, sk_control_battle)
			/*
			 *  This skill will be the number of men he can
			 *  control in battle.
			 *
			 */
			if p == nil {
				man_limit = DEFAULT_CONTROLLED
			} else {
				man_limit = p.experience + DEFAULT_CONTROLLED
			}
		}
	} else if subkind(who) == sub_garrison {
		man_limit = GARRISON_CONTROLLED
	} else {
		man_limit = 10000
	}
	/*
	 *  You may get a bonus...
	 *
	 */
	if defender != 0 {
		man_limit += DEFENDER_CONTROL_BONUS
	}
	/*
	 *  Castle-owner bonus.
	 *
	 */
	if !is_real_npc(who) {
		man_limit += char_rank(who)
	}
	/*
	 *  Charisma effect doubles your man limit.
	 *
	 */
	if get_effect(who, ef_charisma, 0, 0) == 0 {
		man_limit += man_limit
	}
	/*
	 *  Mayhap you have an artifact.
	 *
	 */
	if art := best_artifact(who, ART_CTL_MEN, 0, 0); art != 0 {
		man_limit = ((100 + rp_item_artifact(art).param1) * man_limit) / 100
	}

	return man_limit
}

/*
 *  Determine how many beasts this entity (who) can control in
 *  battle.
 *
 *  (1) Nobles can control equal to their sk_use_beasts skill.
 *
 *  (2) All others can control a virtually unlimited amount.
 *
 */
func calc_beast_limit(who, defender int) int {
	var beast_limit int

	if !is_real_npc(who) {
		p := rp_skill_ent(who, sk_use_beasts)
		/*
		 *  This skill will be the number of men he can
		 *  control in battle.
		 *
		 */
		if p == nil {
			beast_limit = 0
		} else {
			beast_limit = p.experience
		}
	} else if subkind(who) == sub_garrison {
		beast_limit = 0
	} else {
		beast_limit = 10000
	}

	/*
	 *  Mayhap you have an artifact.
	 *
	 */
	if art := best_artifact(who, ART_CTL_BEASTS, 0, 0); art != 0 {
		beast_limit = ((100 + rp_item_artifact(art).param1) * beast_limit) / 100
	}

	return beast_limit
}

func add_fighters(l *[]*fight, who int, ally int, inside bool, defender int) {
	man_limit := 0
	beast_limit := 0

	assert(kind(who) == T_char)

	if is_prisoner(who) {
		return
	}

	add_to_fight_list(l, who, FK_noble, or_int(char_health(who) > 0, char_health(who), 1), ally, inside)

	// mdhender: use_beasts is not used
	// use_beasts := false
	// if (is_real_npc(who) || has_skill(who, sk_use_beasts)!=FALSE) {
	//    use_beasts = true
	// }

	man_limit = calc_man_limit(who, defender)
	beast_limit = calc_beast_limit(who, defender)
	for _, e := range loop_sorted_inv(who) {
		if subkind(e.item) == sub_dead_body {
			continue
		}

		/*
		 *  Fri Nov  8 13:25:09 1996 -- Scott Turner
		 *
		 *  Note: added in order of their strength!
		 *
		 *  Mon Mar  8 12:30:03 1999 -- Scott Turner
		 *
		 *  Special exemption needed for non-animals, too, eg., bandits.
		 *
		 */
		if is_fighter(e.item) != FALSE {
			/*
			 *  Mon Mar  1 12:09:05 1999 -- Scott Turner
			 *
			 *  Special exemption for controlled NPCs.
			 *
			 */
			if !is_real_npc(who) && subkind(who) == sub_ni &&
				e.item == noble_item(who) {
				add_to_fight_list(l, who, e.item, e.qty, ally, inside)
				continue
			}
			if FALSE == item_animal(e.item) {
				num := min(e.qty, man_limit)
				add_to_fight_list(l, who, e.item, num, ally, inside)
				man_limit -= num
			} else {
				num := min(e.qty, beast_limit)
				add_to_fight_list(l, who, e.item, num, ally, inside)
				beast_limit -= num
			}
		}
		if man_limit < 1 && beast_limit < 1 {
			break
		}
	}
	return
}

func add_fight_stack(l *[]*fight, who int, ally int, defender int) {
	assert(kind(who) == T_char)

	inside := len(*l) > 0 && (*l)[0].kind == FK_fort && somewhere_inside((*l)[0].unit, who) != FALSE
	add_fighters(l, who, ally, inside, defender)
	for _, i := range loop_char_here(who) {
		add_fighters(l, i, ally, inside, defender)
	}

	return
}

/*
 *  Sat Sep 30 11:51:27 2000 -- Scott Turner
 *
 *  Side contains a unit?
 *
 */
func contains_unit(l []*fight, unit int) bool {
	for i := 0; i < len(l); i++ {
		if l[i].kind == FK_noble && l[i].unit == unit {
			return true
		}
	}

	return false
}

/*
 *  Look for any characters in where that are allied to def1 or def2
 */
func look_for_allies(l *[]*fight, where int, def1 int, def2 int, attacker int, defender int, l_a []*fight) {
	for _, i := range loop_here(where) {
		if kind(i) == T_char &&
			i != def1 && i != def2 &&
			(allied(i, def1) || allied(i, def2)) &&
			FALSE == char_gone(i) &&
			stack_leader(i) == i &&
			attacker != i &&
			player(attacker) != player(i) &&
			!contains_unit(*l, i) {
			assert(stack_leader(i) == i)
			/*
			 * Wed May 27 12:04:12 1998 -- Scott Turner
			 *
			 * Give this guy a clue...
			 *
			 */
			wout(i, "  %s joins the battle in defense of %s!", box_name(i), box_name(def1))
			wout(VECT, "  %s joins the battle in defense of %s!", box_name(i), box_name(def1))
			wout(combat_pl, "  %s joins the battle in defense of %s!", box_name(i), box_name(def1))
			add_fight_stack(l, i, TRUE, defender)
		}
	}

	return
}

/*
 *  Wed Jun  2 19:10:37 1999 -- Scott Turner
 *
 *  When attacking a loc, who is the real target?
 *
 */
func loc_target(target int) int {
	if subkind(target) == sub_city && garrison_here(province(target)) != FALSE && province_admin(province(target)) != FALSE {
		return garrison_here(province(target))
	} else if kind(target) == T_loc && loc_depth(target) == LOC_province {
		if garrison_here(province(target)) != FALSE && province_admin(province(target)) != FALSE {
			return garrison_here(province(target))
		}
	} else if loc_depth(target) == LOC_build {
		return building_owner(target)
	} else {
		return first_character(target)
	}

	return 0
}

/*
 *  Given the target, fill in the fighter array with the
 *  members of the side's fighting force.
 *
 *  A protecting building will be the first element of the
 *  array.  Nobles and each kind of fighter each get their
 *  own slot.
 *
 *  Then fill in attack and defense and protection values
 *  for the slots.
 */

/*
 *  If target is a location, add the protecting structure as
 *  the first element of the array, then start filling in the
 *  rest with the location owner, and work down through the stack.
 *
 *  If there is no location owner, return an empty array, since
 *  there is no one to fight.
 *
 *  If there is no protecting location, the target will be first
 *  in the array, then men and stackmates follow.
 */
func construct_fight_list(target int, attacker int, defender int) []*fight {
	var l []*fight
	who := 0
	if is_loc_or_ship(target) {
		who = loc_target(target)
		if who == 0 {
			return nil
		}
		if loc_depth(target) == LOC_build {
			hp := or_int(loc_hp(target) != 0, loc_hp(target), 100)
			rating := hp - loc_damage(target)
			assert(rating > 0)
			add_to_fight_list(&l, target, FK_fort, rating, FALSE, false)
		}
	} else {
		who = target
	}
	add_fight_stack(&l, who, FALSE, defender)
	return l
}

// todo: this is just not right...
//
//	the list should hold fighter pointers, not integers
func fight_list_lookup(l []*fight, i int) int { panic("!implemented") }

func construct_guard_fight_list(target int, attacker int, l_a []*fight, defender int) []*fight {
	var l []*fight

	where := subloc(target)
	for _, i := range loop_here(where) {
		if kind(i) != T_char || FALSE == char_guard(i) {
			continue
		}

		/*
		 *  Don't count a guarding unit if they are stacked with the
		 *  pillagers, or if they are part of the pillager's faction.
		 */

		if player(i) == player(attacker) {
			continue
		}

		// todo: this is just not right...
		//       the list should hold fighter pointers, not integers
		if fight_list_lookup(l_a, i) >= 0 {
			continue
		}

		add_fight_stack(&l, i, FALSE, defender)
	}

	/*
	   look_for_allies(&l, subloc(target), target, 0, attacker, defender);
	*/

	return l
}

func ready_fight_list(l []*fight) {
	l = init_prot(l)
	l = init_attack_defense(l)
}

func advance_behind(l []*fight) int {
	least := 0

	for i := 0; i < len(l); i++ {
		if l[i].behind != 0 && l[i].kind != FK_fort && (l[i].behind < least || least == 0) {
			least = l[i].behind
		}
	}

	if least != 0 {
		for i := 0; i < len(l); i++ {
			if l[i].behind == least {
				l[i].behind = 0

				out(combat_pl, "    advancing unit %s.%d", box_code_less(l[i].unit), l[i].kind)

				dump_fighters(l)
			}
		}
	}

	return least
}

func num_targets(f *fight, enemy []*fight) int {
	if f.kind == FK_fort && f.num > 0 {
		return 1
	} else if is_siege_engine(f.kind) && !siege_engine_useful(enemy) {
		return 0
	}
	return f.num
}

func num_valid_targets(f *fight, enemy []*fight) int {
	if f.nprot > 0 || f.behind != 0 {
		return 0
	}
	return num_targets(f, enemy)
}

func num_non_damage(f *fight) int {

	if f.kind == FK_fort {
		return 0
	}

	return f.num
}

func total_non_damage(l []*fight) int {
	sum := 0
	for i := 0; i < len(l); i++ {
		sum += num_non_damage(l[i])
	}
	return sum
}

func combat_sum(f *fight) int {
	if f.kind == FK_fort || is_siege_engine(f.kind) {
		return 0
	}
	assert(f.num >= 0)
	return (max(f.attack, f.missile) + f.defense) * f.num
}

func total_combat_sum(l []*fight) int {
	sum := 0
	for i := 0; i < len(l); i++ {
		sum += combat_sum(l[i])
	}
	return sum
}

/*
 *  Mon Feb 24 15:21:32 1997 -- Scott Turner
 *
 *  The "special" versions of these functions are called after allies join
 *  in the second round, and they calculate the combat sums as if it was
 *  the beginning of the battle, i.e., they use sav_num instead of num.
 *
 *  Wed Apr 15 20:59:17 1998 -- Scott Turner
 *
 *  Except sav_num should properly be the noble's health if it's
 *  a noble...
 *
 */
func special_combat_sum(f *fight) int {
	if f.kind == FK_fort || is_siege_engine(f.kind) {
		return 0
	}
	assert(f.sav_num >= 0)
	if f.kind == FK_noble {
		return (max(f.attack, f.missile) + f.defense)
	} else {
		return (max(f.attack, f.missile) + f.defense) * f.sav_num
	}
}

func special_total_combat_sum(l []*fight) int {
	sum := 0

	for i := 0; i < len(l); i++ {
		sum += special_combat_sum(l[i])
	}

	return sum
}

/*
 *  Register a hit on a fighter
 *
 *  l		side list for g
 *  attacker	fighter doing the hitting
 *  g		fighter that has been hit
 *
 *  Thu Sep  7 06:12:06 2000 -- Scott Turner
 *
 *  Doesn't properly handle multiple hits upon the same target.  To do that,
 *  we need to randomly target between the already dead (g.tmp_num) and the
 *  rest (g.num - g.tmp_num).
 *
 */
func decrement_num(l []*fight, attacker *fight, g *fight) {
	/*
	 *  How can we be hit if we're supposed to be protected?
	 */
	assert(g.nprot == 0)

	/*
	 *  Perhaps we swung against someone already dead?
	 *
	 */
	if rnd(1, g.num) <= g.tmp_num {
		wout(combat_pl, "     (%s.%d already dead.)", box_code_less(g.unit), g.kind)
		return
	}

	switch g.kind {
	case FK_noble:
		g.tmp_num++
		break

	case FK_fort:
		/*
		 *  Mon Oct  7 15:05:25 1996 -- Scott Turner
		 *  A fortification may be protected by a prayer.
		 *
		 *  Fri Jan 10 11:05:39 1997 -- Scott Turner
		 *  No hits against the defense anymore.
		 *
		 */
		if FALSE == get_effect(g.unit, ef_bless_fort, 0, 0) || rnd(1, 2) == 1 {
			var hit int
			if attacker != nil && is_siege_engine(attacker.kind) {
				hit = rnd(5, 10)
			} else {
				hit = 1
			}

			g.tmp_num += hit
		}

		break

	case item_blessed_soldier:
		if rnd(1, 2) == 1 { /* 50% chance of surviving a hit */
			g.tmp_num++
		}
		break

	default:
		g.tmp_num++
	}
}

/*
 *  Wed Aug 21 13:46:58 1996 -- Scott Turner
 *
 *  Check to see if a side contains a priest.
 *
 */
func contains_priest(l []*fight) bool {
	for i := 0; i < len(l); i++ {
		if l[i].kind == FK_noble && is_priest(l[i].unit) != FALSE {
			return true
		}
	}
	return false
}

/*
 *  Wed Aug 21 13:50:11 1996 -- Scott Turner
 *
 *  Check to see if a side contains a magician or undead.
 *
 */
func contains_mu_or_undead(l []*fight) bool {
	for i := 0; i < len(l); i++ {
		if l[i].kind == FK_fort {
			continue
		}
		if l[i].kind == FK_noble {
			if is_magician(l[i].unit) && FALSE == char_hide_mage(l[i].unit) {
				return true
			}
			continue
		}
		if subkind(l[i].kind) == sub_undead {
			return true
		}
		if subkind(l[i].kind) == sub_demon_lord {
			return true
		}
	}
	return false
}

/*
 *  Fri Nov  8 15:07:19 1996 -- Scott Turner
 *
 *  Revised version with types.
 *
 */
func num_defenders(f *fight, enemy []*fight, type_ int) int {
	/*
	 *  Protected don't defend.
	 *
	 */
	if f.nprot > 0 {
		return 0
	}

	/*
	 *  Behind are safe except from SPECIAL & UNCANNY
	 *
	 */
	if type_ != SPECIAL && type_ != UNCANNY && f.behind != 0 {
		return 0
	}

	/*
	 *  UNCANNY only attacks non-noble missile weapons behind.
	 *
	 */
	if type_ == UNCANNY && (f.kind == FK_noble || f.behind == 0 || f.missile == 0) {
		return 0
	}

	/*
	 *  A fort is only 1 target.
	 *
	 */
	if f.kind == FK_fort && f.num > 0 {
		return 1
	}

	/*
	 *  A siege engine is a target only if it is useful.
	 *
	 */
	if is_siege_engine(f.kind) && !siege_engine_useful(enemy) {
		return 0
	}

	return f.num
}

/*
 *  Fri Nov  8 15:03:59 1996 -- Scott Turner
 *
 *  Total defenders in a given category.
 *
 */
func total_defenders(l []*fight, enemy []*fight, type_ int) int {
	sum := 0
	for i := 0; i < len(l); i++ {
		sum += num_defenders(l[i], enemy, type_)
	}

	/*
	 *  Fri Nov  8 15:05:41 1996 -- Scott Turner
	 *
	 *  If there are no melee defenders, then we need to advance
	 *  some people to the front line.
	 *
	 */
	for type_ != SPECIAL && type_ != UNCANNY && sum == 0 && advance_behind(l) != 0 {
		for i := 0; i < len(l); i++ {
			sum += num_defenders(l[i], enemy, type_)
		}
	}

	assert(sum != 0 || type_ == UNCANNY)

	return sum
}

/*
 *  Mon Feb 17 21:07:05 1997 -- Scott Turner
 *
 *  Does he have any combat spell?
 *
 */
func has_magic_combat_attack(who int) int {
	num_attacks := 0

	for _, e := range loop_char_skill(who) {
		if e.skill != 0 && e.know == SKILL_know && (skill_flags(e.skill)&COMBAT_SKILL) != 0 {
			num_attacks++
		}
	}

	return num_attacks
}

/*
 *  Fri Nov  8 15:28:30 1996 -- Scott Turner
 *
 */
func toughest_defender(l []*fight, enemy []*fight, type_ int) *fight {
	var ret *fight

	//toughest := 0
	best_attack := 0

	/*
	 *  Thu Mar  1 05:53:24 2001 -- Scott Turner
	 *
	 *  This call is just to force some defenders to advance, if
	 *  necessary
	 *
	 */
	_ = total_defenders(l, enemy, type_)

	for i := 0; i < len(l); i++ {
		if num_defenders(l[i], enemy, type_) != 0 && l[i].attack > best_attack {
			best_attack = l[i].attack
			ret = l[i]
		}
	}

	return ret
}

/*
 *  Fri Nov  8 15:28:30 1996 -- Scott Turner
 *
 */
func find_defender(l []*fight, man int, enemy []*fight, type_ int) *fight {
	for i := 0; i < len(l); i++ {
		man -= num_defenders(l[i], enemy, type_)
		if man <= 0 {
			return l[i]
		}
	}

	panic("!reached")
}

/*
 *  Tue Feb 18 15:53:06 1997 -- Scott Turner
 *
 *  Dispatch for various special magic attacks.
 *
 *  Fri Apr 27 10:58:51 2001 -- Scott Turner
 *
 *  This is where we hook in to add ef_cs control of combat spells.
 *
 */
// returns the amount of magical damage done
func do_magic_combat_attack(l_a []*fight, f *fight, l_b []*fight, skill int) int {
	sum := 0
	max_raise := 10

	if !check_aura(f.unit, skill_piety(skill)) {
		return 0
	}

	/*
	 *  Fri Apr 27 11:00:11 2001 -- Scott Turner
	 *
	 *  If he doesn't want to use this skill this round, punt.
	 *
	 */
	if FALSE == get_effect(f.unit, ef_cs, skill, or_int((round >= 10), 10, round)) {
		out(combat_pl, "%s skips using %s in round: %d.", box_name(f.unit), box_name(skill), round)
		return 0
	}

	switch skill {
	case sk_lightning_bolt:
		g := toughest_defender(l_b, l_a, MELEE)
		if g == nil {
			return 0
		}
		if !charge_aura(f.unit, skill_piety(skill)) {
			return 0
		}
		if g.defense < 80 {
			return 0
		}
		sum = 0 // todo: lightning bolts are only 1 point of damage?
		if real_resolve_hit(l_a, f, l_b, g, 1, SPECIAL, 500, g.defense) {
			sum++
		}
		out(combat_pl, "%s uses a lightning bolt, killing %d.", box_name(f.unit), sum)
		print_special_banner()
		if sum != 0 {
			wout(VECT, "    %s uses a lightning bolt, striking %s.", box_name(f.unit), nice_num(sum))
		} else {
			wout(VECT, "    %s uses a lightning bolt, but misses!", box_name(f.unit))
		}
		return sum

	case sk_fireball:
		/*
		 *  Thu Mar  1 06:13:32 2001 -- Scott Turner
		 *
		 *  Don't waste a fireball?
		 *
		 */
		if total_defenders(l_b, l_a, MELEE) < 10 && total_defenders(l_a, l_b, MELEE) > 10 {
			return 0
		}

		if !charge_aura(f.unit, skill_piety(skill)) {
			return 0
		}
		for i := 0; i < 50; i++ {
			num_defend := total_defenders(l_b, l_a, MELEE)
			assert(num_defend > 0)
			man := rnd(1, num_defend)
			g := find_defender(l_b, man, l_a, MELEE)
			if real_resolve_hit(l_a, f, l_b, g, man, SPECIAL, 500, g.defense) {
				sum++ // todo: fireballs are only 1 point of damage?
			}
		}
		out(combat_pl, "%s uses a fireball, killing %d.", box_name(f.unit), sum)
		print_special_banner()
		wout(VECT, "    %s casts a fireball and engulfs %s victims!", box_name(f.unit), nice_num(sum))
		return sum

	case sk_drain_mana:
		if rnd(1, 100) > 75 ||
			contains_mu_or_undead(l_b) || contains_priest(l_b) {
			for i := 0; i < len(l_b); i++ {
				if l_b[i].kind == FK_noble {
					amount := char_cur_aura(l_b[i].unit)
					if amount > skill_piety(skill) {
						amount = skill_piety(skill)
					}
					if charge_aura(l_b[i].unit, amount) {
						sum++
					}
				}
			}
		}
		print_special_banner()
		wout(VECT, "    %s opens a mana-draining gate!", box_name(f.unit))
		if sum != 0 {
			out(combat_pl, "%s drains %d magicians of aura.", box_name(f.unit), sum)
			add_aura(f.unit, sum/10)
		}
		return sum

	case sk_raise_soldiers:
		if !charge_aura(f.unit, skill_piety(skill)) {
			return 0
		}
		for i := 0; i < len(l_a) && max_raise > 0; i++ {
			if l_a[i].kind != FK_noble && l_a[i].num < l_a[i].sav_num {
				num_to_raise := l_a[i].sav_num - l_a[i].num
				if num_to_raise > max_raise {
					num_to_raise = max_raise
				}
				l_a[i].num += num_to_raise
				l_a[i].tmp_num2 += num_to_raise
				max_raise -= num_to_raise
				print_special_banner()
				wout(VECT, "    %s raises %s %s.", box_name(f.unit), nice_num(num_to_raise), plural_item_name(l_a[i].kind, num_to_raise))
				out(combat_pl, "%s raises %d %s.", box_name(f.unit), num_to_raise, box_name(l_a[i].unit))
			}
		}
		return (10 - max_raise)

	case sk_foresee_defense:
		if !charge_aura(f.unit, skill_piety(skill)) {
			return 0
		}
		if FALSE == add_effect(f.unit, ef_scry_offense, 0, 0, 0) {
			return 0
		}
		print_special_banner()
		wout(VECT, "    %s foresees the next moments of battle.", box_name(f.unit))
		out(combat_pl, "%s foresees the tide of battle.", box_name(f.unit))
		return 1

	default:
		break
	}
	return 0
}

/*
 *  Tue Feb 18 12:13:19 1997 -- Scott Turner
 *
 *  Run through a noble's magic attacks.
 *
 */
func do_magic_combat_attacks(l_a []*fight, f *fight, l_b []*fight) int {
	sum := 0
	for _, e := range loop_char_skill(f.unit) {
		if e.skill != 0 && e.know == SKILL_know && (skill_flags(e.skill)&COMBAT_SKILL) != 0 {
			sum += do_magic_combat_attack(l_a, f, l_b, e.skill)
		}
	}
	return sum
}

/*
 *  Wed Sep 23 06:48:09 1998 -- Scott Turner
 *
 *  An alchemist attacking with "fiery death".
 *
 */
// returns the amount of fiery death done
func do_fiery_attack(l_a []*fight, f *fight, l_b []*fight) int {
	sum := 0
	for _, e := range inventory_loop(f.unit) {
		if item_use_key(e.item) == use_fiery_potion {
			print_special_banner()
			out(VECT, "%s throws a potion of Fiery Death!", box_name(f.unit))
			if rnd(1, 100) < 6 {
				out(combat_pl, "The potion lands harmlessly among the defenders.")
				break
			}
			for i := 0; i < 5; i++ {
				num_defend := total_defenders(l_b, l_a, MELEE)
				assert(num_defend > 0)
				man := rnd(1, num_defend)
				g := find_defender(l_b, man, l_a, MELEE)
				if real_resolve_hit(l_a, f, l_b, g, man, SPECIAL, 30, g.defense) {
					sum++
				}
			}
			if sum != 0 {
				out(combat_pl, "The potion explodes, killing %s!", nice_num(sum))
			} else {
				out(combat_pl, "The explosion does no damage.")
			}
			break
		}
	}
	return sum
}

/*
 *  Fri Nov  8 15:07:19 1996 -- Scott Turner
 *
 *  Revised version with types.
 *
 */
// returns the number of attackers
func num_attackers(f *fight, enemy []*fight, type_ int) int {
	/*
	 *  Forts never attack
	 *
	 */
	if f.kind == FK_fort {
		return 0
	}

	/*
	 *  Behinds never melee.
	 *
	 */
	if f.behind != 0 && type_ == MELEE {
		return 0
	}

	/*
	 *  Fronts never missile.
	 *
	 */
	if f.behind == 0 && type_ == MISSILE {
		return 0
	}

	/*
	 *  Need an attack if you're a meleer
	 *
	 */
	if type_ == MELEE && f.attack == 0 {
		return 0
	}

	/*
	 *  Ditto for missiles.
	 *
	 */
	if f.missile == 0 && type_ == MISSILE {
		return 0
	}

	/*
	 *  Siege engine may not be useful
	 *
	 */
	if is_siege_engine(f.kind) && !siege_engine_useful(enemy) {
		return 0
	}

	return f.num
}

/*
 *  Fri Nov  8 15:03:59 1996 -- Scott Turner
 *
 *  Total attackers in a given category.
 *
 */
func total_attackers(l []*fight, enemy []*fight, type_ int) int {
	sum := 0
	for i := 0; i < len(l); i++ {
		sum += num_attackers(l[i], enemy, type_)
	}

	/*
	 *  Fri Nov  8 15:05:41 1996 -- Scott Turner
	 *
	 *  If there are no melee attackers, then we need to advance
	 *  some people to the front line.
	 *
	 */
	if type_ == MELEE && sum == 0 && advance_behind(l) != 0 {
		for i := 0; i < len(l); i++ {
			sum += num_attackers(l[i], enemy, type_)
		}
	}

	return sum
}

/*
 *  Fri Nov  8 15:27:32 1996 -- Scott Turner
 *
 */
func find_attacker(l []*fight, man int, enemy []*fight, type_ int) *fight {
	for i := 0; i < len(l); i++ {
		man -= num_attackers(l[i], enemy, type_)
		if man <= 0 {
			return l[i]
		}
	}
	panic("!reached")
}

/*
 *  Tue Apr 28 08:53:57 1998 -- Scott Turner
 *
 */
// returns the number of siege towers
func count_siege_towers(l []*fight) int {
	num := 0
	for i := 0; i < len(l); i++ {
		if l[i].kind == item_siege_tower {
			num += l[i].num
		}
	}
	return num
}

// returns the attack(?) factors of the siege tower
func siege_tower_factor(l *fight) int {
	switch subkind(l.unit) {
	case sub_castle:
		return 1
	case sub_tower:
		return 5
	}
	return 10
}

/*
 *  For some reason Rich has min() defined for ints.
 *
 */
func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

/*
 *  Tue Apr  7 06:56:41 1998 -- Scott Turner
 *
 *  Re-arrange so + first and then * for attack/defense effects.
 *
 */
// todo: attacks do from 0...1 point of damage? is that right?
func real_resolve_hit(l_a []*fight, f *fight, l_b []*fight, g *fight, man int, type_ int, attack int, defense int) bool {
	/*
	 *  Tue Feb 18 20:35:11 1997 -- Scott Turner
	 *
	 *  Multiply by 100 so the % effects are reasonable.
	 *
	 */
	attack *= 100
	defense *= 100

	/**** Beginning of additive attack bonuses ****/

	/*
	 *  If there's a priest on the attacking side, and undead or
	 *  magicians on the other side, then the attacker gets +10
	 *  to his AF.
	 *
	 *  Tue Nov 10 11:48:03 1998 -- Scott Turner
	 *
	 *  Only if we have "antipathy" turned on.
	 */
	if options.mp_antipathy &&
		contains_priest(l_a) && contains_mu_or_undead(l_b) {
		/*
		 *  Note that +10 is actually +10*100 == +1000
		 *
		 */
		attack += 1000
	}

	/*
	 *  If the unit attacking is stacked under a priest of Kireus
	 *  with the edge of Kireus in effect, give him +25(00) to hit
	 *  if he's wielding an edged weapon.
	 *
	 */
	if get_effect(f.unit, ef_edge_of_kireus, 0, 0) != 0 &&
		type_ == MELEE &&
		(f.kind == item_swordsman ||
			f.kind == item_hvy_foot ||
			f.kind == item_knight ||
			f.kind == item_elite_guard ||
			f.kind == item_cavalier ||
			f.kind == item_pirate) {
		attack += 2500
	}

	/*
	 *  Artifact-based bonuses.  For all the "combat artifact" bonuses,
	 *  we can just add in the skill we get back from combat_artifact_bonus,
	 *  since that will be zero if he doesn't have one.
	 *
	 */
	var item int
	if f.kind == FK_noble {
		switch type_ {
		case MELEE:
			attack += (100 * combat_artifact_bonus(f.unit, CA_N_MELEE, &item))
			break
		case MISSILE:
			attack += (100 * combat_artifact_bonus(f.unit, CA_N_MISSILE, &item))
			break
		case SPECIAL:
			attack += (100 * combat_artifact_bonus(f.unit, CA_N_SPECIAL, &item))
			break
		default:
			break
		}
	} else {
		switch type_ {
		case MELEE:
			attack += (100 * combat_artifact_bonus(f.unit, CA_M_MELEE, &item))
			break
		case MISSILE:
			attack += (100 * combat_artifact_bonus(f.unit, CA_M_MISSILE, &item))
			break
		case SPECIAL:
			attack += (100 * combat_artifact_bonus(f.unit, CA_M_SPECIAL, &item))
			break
		default:
			break
		}
	}

	if art := best_artifact(f.unit, ART_IMPRV_ATT, f.kind, 0); art != 0 {
		attack += (100 * rp_item_artifact(art).param1)
	}

	/**** Beginning of mulplicative attack bonuses ****/

	/*
	 *  Adjust attack for leader noble's (possible) specials.
	 *  The leader is "unit", as far as I can tell.
	 *
	 */
	p := rp_skill_ent(f.unit, sk_attack_tactics)
	if p != nil && p.know == SKILL_know {
		attack = int(float64(attack) * minf(1.0+float64(p.experience)*TACTICS_FACTOR, 2.0))
	}

	/*
	 *  Adjust attack for scry offense.
	 *
	 */
	if get_effect(f.unit, ef_scry_offense, 0, 0) != 0 {
		attack = int(float64(attack) * 1.1)
	}

	/**** Beginning of additive defense bonuses ****/

	/*
	 *  Fri Jan 10 11:14:51 1997 -- Scott Turner
	 *
	 *  Defense modified by damage...
	 *
	 */
	defense_bonus := 0
	if g.inside && man <= fort_covers(l_b[0].unit) {
		hp := or_float(loc_hp(l_b[0].unit) != 0, float64(loc_hp(l_b[0].unit)), 100.0)
		defense_bonus = int(float64(l_b[0].defense*l_b[0].num) / hp)
		/*
		 *  Tue Apr 28 08:52:03 1998 -- Scott Turner
		 *
		 *  Defensive bonus can be reduced by siege towers on the
		 *  opposing side.
		 *
		 */
		siege_towers := count_siege_towers(l_a) * siege_tower_factor(l_b[0])

		out(combat_pl, "%s.%d gets fort bonus of %d, minus %d siege towers.",
			box_code_less(g.unit), g.kind, defense_bonus, siege_towers)

		if siege_towers > defense_bonus {
			siege_towers = defense_bonus
		}
		defense_bonus -= siege_towers

		defense += defense_bonus
	}

	/*
	 *  If there's a moat, then defenders against missile attacks
	 *  get a +100, although that doesn't apply to siege weapons.
	 *
	 */
	if g.inside &&
		man <= fort_covers(l_b[0].unit) &&
		loc_moat(l_b[0].unit) != 0 &&
		type_ == MISSILE &&
		!is_siege_engine(f.kind) {
		out(combat_pl, "%s.%d gets moat bonus.", box_code_less(g.unit), g.kind)
		defense += 100 * 100
	}

	/*
	 *  Might be a "Bless Follower" on the guy getting hit.  Note that
	 *  the bless won't work if the defender has undead or is a magician.
	 *
	 */
	if get_effect(g.unit, ef_defense, 0, 0) != 0 && (!options.mp_antipathy || FALSE == contains_mu_undead(g.unit)) {
		defense += get_effect(g.unit, ef_defense, 0, 0) * 100
	}

	/*
	 *  Artifact-based bonuses.  For all the "combat artifact" bonuses,
	 *  we can just add in the skill we get back from combat_artifact_bonus,
	 *  since that will be zero if he doesn't have one.
	 *
	 */
	if f.kind == FK_noble {
		switch type_ {
		case MELEE:
			defense += (100 * combat_artifact_bonus(f.unit, CA_N_MELEE_D, &item))
			break
		case MISSILE:
			defense += (100 * combat_artifact_bonus(f.unit, CA_N_MISSILE_D, &item))
			break
		case SPECIAL:
			defense += (100 * combat_artifact_bonus(f.unit, CA_N_SPECIAL_D, &item))
			break
		default:
			break
		}
	} else {
		switch type_ {
		case MELEE:
			defense += (100 * combat_artifact_bonus(f.unit, CA_M_MELEE_D, &item))
			break
		case MISSILE:
			defense += (100 * combat_artifact_bonus(f.unit, CA_M_MISSILE_D, &item))
			break
		case SPECIAL:
			defense += (100 * combat_artifact_bonus(f.unit, CA_M_SPECIAL_D, &item))
			break
		default:
			break
		}
	}

	if art := best_artifact(g.unit, ART_IMPRV_DEF, f.kind, 0); art != 0 {
		defense += (100 * rp_item_artifact(art).param1)
	}

	/**** Beginning of mulplicative defense bonuses ****/

	/*
	 *  Artifacts.
	 *
	 */
	if art := best_artifact(g.unit, ART_TERRAIN, int(subkind(combat_def_loc)), 0); art != 0 {
		defense = ((100 + rp_item_artifact(art).param1) * defense) / 100
	}

	/*
	 *  Adjust defense for leader noble's (possible) specials.
	 *
	 */
	p = rp_skill_ent(g.unit, sk_defense_tactics)
	if p != nil && p.know == SKILL_know {
		defense = int(minf(1.0+float64(p.experience)*TACTICS_FACTOR, 2.0) * float64(defense))
	}

	/*
	 *  Adjust defense for terrain (location) bonuses -- but only
	 *  if the unit being attacked is on the defense.
	 *
	 */
	// todo: C pointer compares don't work here
	log.Printf("real_resolve_hit: C pointer compares don't work here\n")
	sameSide := reflect.DeepEqual(l_b, defense_side)
	if sameSide && combat_def_loc != 0 && subkind(combat_def_loc) != 0 && type_ != SPECIAL {
		switch subkind(combat_def_loc) {
		case sub_forest:
			defense = int(float64(defense) * FOREST_DEFENSE_BONUS)
			break
		case sub_mountain:
			defense = int(float64(defense) * MOUNTAIN_DEFENSE_BONUS)
			break
		case sub_city:
			defense = int(float64(defense) * CITY_DEFENSE_BONUS)
			break
		case sub_swamp:
			defense = int(float64(defense) * SWAMP_DEFENSE_BONUS)
			break
		}
	}

	/*
	 *  Pikeman double their defense against mounted troops.
	 *
	 */
	if g.kind == item_pikeman &&
		(f.kind == item_knight || f.kind == item_elite_guard ||
			f.kind == item_cavalier || f.kind == item_centaur) {
		defense *= 2
	}

	n := rnd(1, attack+defense)
	if n > attack {
		out(combat_pl, "    %s.%d failed to hit %s.%d (%d/%d)", box_code_less(f.unit), f.kind, box_code_less(g.unit), g.kind, attack, attack+defense)
		return false
	}

	out(combat_pl, "    %s.%d hit %s.%d (%d/%d)", box_code_less(f.unit), f.kind, box_code_less(g.unit), g.kind, attack, attack+defense)

	decrement_num(l_b, f, g) /* f scores against g */
	return true
}

// todo: attacks do 0...1 point of damage?
func resolve_hit(l_a []*fight, f *fight, l_b []*fight, g *fight, man int, type_ int) bool {
	if type_ == MISSILE {
		assert(f.missile != 0)
		return real_resolve_hit(l_a, f, l_b, g, man, type_, f.missile, g.defense)
	}
	assert(f.attack != 0)
	return real_resolve_hit(l_a, f, l_b, g, man, type_, f.attack, g.defense)
}

/*
 *  Mon Nov 30 08:48:48 1998 -- Scott Turner
 *
 *  The number of attacks an attacker gets -- one, unless he's a hero.
 *
 */
func num_attacks(f *fight, type_ int) int {
	if f.kind != FK_noble || count_any_real(f.unit, false, false) > 1 {
		return 1
	}

	var sk int
	if type_ == MELEE {
		sk = sk_extra_attacks
	} else if type_ == MISSILE {
		sk = sk_extra_missile_attacks
	} else {
		return 1
	}

	if FALSE == has_skill(f.unit, sk) {
		return 1
	}

	p := rp_skill_ent(f.unit, sk)
	if p == nil {
		return 1
	}

	return max(1, p.experience/8)
}

/*
 *  Fri Nov  8 15:24:48 1996 -- Scott Turner
 *
 *  Tue Feb 18 12:06:22 1997 -- Scott Turner
 *
 *
 */
func do_attack(attacker int, l_a []*fight, l_b []*fight, type_ int) int {
	man := -1

	f := find_attacker(l_a, attacker, l_b, type_)

	/*
	 *  Mon Nov 30 08:46:32 1998 -- Scott Turner
	 *
	 *  Heroes can have multiple attacks.  That's handled by looping
	 *  over the actual attack.
	 *
	 */
	attacks, uncanny, kills := num_attacks(f, type_), 0, 0
	for i := 0; i < attacks; i++ {
		var g *fight
		if is_siege_engine(f.kind) && siege_engine_useful(l_b) {
			g = l_b[0] /* set defender to structure */
			assert(g.kind == FK_fort)
			assert(g.num > 0)
		} else {
			/*
			 *  Mon Nov 30 17:40:49 1998 -- Scott Turner
			 *
			 *  If a hero has uncanny accuracy and he's in the back row
			 *  shooting missiles, then he gets to target the back row missile
			 *  troops of the other side.  If there aren't any, he acts like
			 *  a normal missiler.
			 *
			 */
			if type_ == MISSILE && has_skill(f.unit, sk_uncanny_accuracy) != FALSE && total_defenders(l_b, l_a, UNCANNY) > 0 {
				type_ = UNCANNY
				uncanny++
			}
			num_defend := total_defenders(l_b, l_a, type_)
			assert(num_defend > 0)
			man = rnd(1, num_defend)
			g = find_defender(l_b, man, l_a, type_)
		}
		if resolve_hit(l_a, f, l_b, g, man, type_) /* f tries to hit g */ {
			kills++
		}
	}

	if uncanny == 1 {
		wout(VECT, "    With uncanny accuracy, %s fires into the back row!", box_name(f.unit))
	} else if uncanny != 0 {
		wout(VECT, "    With uncanny accuracy, %s fires into the back row %s times!", box_name(f.unit), nice_num(uncanny))
	}
	return kills
}

/*
 *  Fri Nov  8 13:28:10 1996 -- Scott Turner
 *
 *  New version: A attacks B
 *
 *  Fri Nov  8 15:01:57 1996 -- Scott Turner
 *
 *  type_ == SPECIAL, MISSILE or MELEE
 *
 *  Wed May 13 08:19:47 1998 -- Scott Turner
 *
 *  Limit the number of missile attacks to the number of men in the
 *  front row.
 *
 */
func attack_round(l_a []*fight, l_b []*fight, type_ int) int {
	assert(type_ == MELEE || type_ == MISSILE)

	/*
	 *  Wed May 13 18:45:57 1998 -- Scott Turner
	 *
	 *  We need to force troops into the front line for the attacker if
	 *  this is a missile attack, since they need people in front of them.
	 *
	 */
	num_defend_a := total_defenders(l_a, l_b, MELEE)
	order_size := total_attackers(l_a, l_b, type_)
	total_attacks := order_size

	/*
	 *  You might well have no attacks, in which case do nothing.
	 * */
	if total_attacks < 1 {
		return NO_ATTACKS
	}

	num_defend_b := total_defenders(l_b, l_a, type_)

	/*
	 *  Wed May 13 18:08:50 1998 -- Scott Turner
	 *
	 *  In missile attacks, you cannot have more archers than you
	 *  have protection in the front row, which turns out to be the
	 *  number of defenders *you* have against a MELEE attack.
	 *
	 */
	if type_ == MISSILE {
		if total_attacks > num_defend_a {
			total_attacks = num_defend_a
		}
	}

	/*
	 *  Fri Nov  8 15:19:31 1996 -- Scott Turner
	 *
	 *  In melee, you're not allowed to hide more people in the
	 *  back than you have in the front.  To figure this out, we
	 *  get the total defense for a special attack (which hits
	 *  all ranks) and subtract out the melee defenders.
	 *
	 */
	if type_ == MELEE {
		num_defend_special := total_defenders(l_b, l_a, SPECIAL)
		if num_defend_special-num_defend_b > num_defend_b {
			num_defend_b = num_defend_special / 2
		}
	}

	/*
	 *  Fri Nov  8 15:02:39 1996 -- Scott Turner
	 *
	 *  In melee, only (4x the defenders) attacks are allowed.
	 *
	 */
	if type_ == MELEE && total_attacks > num_defend_b*4 {
		total_attacks = num_defend_b * 4
	}

	out(combat_pl, "  %s phase:  Attacks = %d", or_string(type_ == MELEE, "Melee", "Missile"), total_attacks)

	out(combat_pl, "  %s phase:  Defenders = %d", or_string(type_ == MELEE, "Melee", "Missile"), num_defend_b)

	assert(num_defend_b > 0)

	/*
	 *  Sun Dec  6 19:24:53 1998 -- Scott Turner
	 *
	 *  If there are less attacks then there are men, then you really
	 *  want the attacks to be selected randomly (w/o replacement) rather
	 *  than to just take the first "n" attacks off the front of l_a.
	 *
	 *  We malloc up an array of the right size, fill it with 1..n, scramble
	 *  it, and then do attacks in that order.  Sigh.
	 */

	/*
	 *  Create, populate and scramble the array.
	 *
	 */
	var order []int
	for i := 0; i < order_size; i++ {
		order = append(order, i+1)
	} /* Men are 1 based */
	order = shuffle_ints(order)

	sum := 0
	for man := 0; man < total_attacks; man++ {
		sum += do_attack(order[man], l_a, l_b, type_) /* b is the defender */
	}

	out(combat_pl, "  %s phase: %d hits out of %d attacks.", or_string(type_ == MELEE, "Melee", "Missile"), sum, total_attacks)

	return sum
}

/*
 *  Tue Feb 18 15:43:22 1997 -- Scott Turner
 *
 *  Special attacks code.
 *
 */

func do_special_attack(l_a []*fight, f *fight, l_b []*fight) int {
	sum := 0
	if f.kind == item_dragon && (round%3) == 1 {
		for i := 0; i < 10*f.num; i++ {
			num_defend := total_defenders(l_b, l_a, MELEE)
			assert(num_defend > 0)
			man := rnd(1, num_defend)
			g := find_defender(l_b, man, l_a, MELEE)
			if real_resolve_hit(l_a, f, l_b, g, man, SPECIAL, 200, g.defense) {
				sum++
			}
		}
		out(combat_pl, "%s breathes dragonfire, killing %d.", box_name(f.unit), sum)
		print_special_banner()
		if f.num > 1 {
			wout(VECT, "    %s's dragons release their dragonfire, engulfing %s!", box_name(f.unit), nice_num(sum))
		} else {
			wout(VECT, "    %s's dragon releases its dragonfire, engulfing %s!", box_name(f.unit), nice_num(sum))
		}
		return sum
	}
	return 0
}

/*
 *  Tue Feb 18 15:46:47 1997 -- Scott Turner
 *
 *  Do a special attack from f.  Right now the only option is
 *  noble magic; eventually there will be monster attacks as well.
 *
 */
func special_attack(l_a []*fight, i int, l_b []*fight, type_ int) int {
	f := l_a[i]

	/*
	 *  Spell attacks.
	 *
	 */
	if round > 1 &&
		f.kind == FK_noble &&
		f.behind != 0 &&
		has_magic_combat_attack(f.unit) != 0 {
		return do_magic_combat_attacks(l_a, f, l_b)
	}

	/*
	 *  Alchemist attack.
	 *
	 */
	if round > 1 && f.kind == FK_noble && f.behind != 0 && has_skill(f.unit, sk_brew_fiery) != 0 {
		return do_fiery_attack(l_a, f, l_b)
	}

	/*
	 *  Hero attack.
	 *
	 */
	if f.kind == FK_noble && count_any_real(f.unit, false, false) == 1 && has_skill(f.unit, sk_blinding_speed) != 0 && has_item(f.unit, item_warmount) > 0 && f.behind == 0 {
		print_special_banner()
		out(VECT, "      %s dashes out and attacks!", box_name(f.unit))
		num_defend := total_defenders(l_b, l_a, MELEE)
		assert(num_defend > 0)
		man := rnd(1, num_defend)
		g := find_defender(l_b, man, l_a, MELEE)
		// todo: should we not do something with this result?
		resolve_hit(l_a, f, l_b, g, man, MELEE) /* f tries to hit g */
	}

	/*
	 *  Other special attacks.
	 *
	 */
	return do_special_attack(l_a, f, l_b)
}

func special_attacks(l_a []*fight, l_b []*fight, type_ int) int {
	total_attacks := 0
	for i := 0; i < len(l_a); i++ {
		total_attacks += special_attack(l_a, i, l_b, type_)
	}
	return total_attacks
}

/*
 *  Fri Nov  8 12:04:49 1996 -- Scott Turner
 *
 *  Mark the dead men "for real"
 *
 *  Mon Jun  8 18:16:44 1998 -- Scott Turner
 *
 *  This is where code to do "breaking" for peasants & skirmishers
 *  will go.
 *
 */
func resolve_dead(l_a []*fight) {
	for i := 0; i < len(l_a); i++ {
		g := l_a[i]

		/*
		 *  Ignore fighters that need no resolving.
		 *
		 */
		if g.num <= 0 {
			continue
		}
		if g.tmp_num <= 0 {
			continue
		}

		if g.tmp_num > g.num {
			g.tmp_num = g.num
		}
		if g.kind == FK_fort {
			g.num -= g.tmp_num
			wout(VECT, "      %s takes %s point%s damage.", box_name(g.unit),
				nice_num(g.tmp_num),
				or_string(g.tmp_num > 1, "s", ""))
		} else if g.kind == FK_noble {
			/*
			 *  Wed Nov 25 10:36:23 1998 -- Scott Turner
			 *
			 *  We need to calculate the health hit here, because the new
			 *  heroism skills mean a noble might stay in a battle and suffer
			 *  multiple hits.  Also, we need the option to reset "num" so that
			 *  he doesn't leave the battle.
			 *
			 */
			wout(VECT, "      %s suffers a hit!", box_name(g.unit))
			damage := rnd(1, 100)
			/*
			 *  Damage possibly reduced by skill.
			 *
			 */
			if has_skill(g.unit, sk_avoid_wounds) != 0 && count_any(g.unit, FALSE, FALSE) == 1 {
				damage -= (damage * min(2*skill_exp(g.unit, sk_avoid_wounds), 80)) / 100
			}
			g.new_health = max(g.new_health-damage, 0)
			/*
			 *  Now (possibly) adjust g.num to take the noble out of the
			 *  battle.
			 *
			 */
			if FALSE == has_skill(g.unit, sk_personal_fttd) ||
				count_any(g.unit, FALSE, FALSE) > 1 ||
				g.new_health <= p_char(g.unit).personal_break_point {
				g.num -= g.tmp_num
			} else {
				wout(VECT, "      %s heroically continues to fight!", box_name(g.unit))
			}
		} else {
			g.num -= g.tmp_num
			wout(VECT, "      %s loses %s %s.", box_name(g.unit), nice_num(g.tmp_num), plural_item_box(g.kind, g.tmp_num))
			/*
			 *  Peasants, workers & skirmishers will "break".
			 *
			 */
			if g.kind == item_skirmisher || g.kind == item_peasant || g.kind == item_worker {
				brk := 0
				limit := 2
				if g.kind == item_skirmisher {
					limit = 3
				}
				for j := 0; i < g.num; j++ {
					if rnd(1, limit) == 1 {
						brk++
					}
				}
				if brk != 0 {
					wout(VECT, "        %d %s break and leave the battle!", brk, plural_item_box(g.kind, brk))
					g.num -= brk
				}
			}
		}
		g.tmp_num = 0

		if g.num <= 0 && g.protects >= 0 {
			assert(l_a[g.protects].nprot > 0)

			l_a[g.protects].nprot--

			out(combat_pl, "    %s.%d no longer protects %s", box_code_less(g.unit), g.kind, box_code_less(l_a[g.protects].unit))

			dump_fighters(l_a)
		}
	}
}

func side_has_skill(l []*fight, sk int) bool {
	for i := 0; i < len(l); i++ {
		if l[i].kind == FK_noble && has_skill(l[i].unit, sk) != FALSE {
			return true
		}
	}
	return false
}

/*
 *  Fri Sep 11 09:42:51 1998 -- Scott Turner
 *
 *  Generate animal parts and give them to the noble.
 *
 */
func generate_animal_parts(who int, what int, how_many int) {
	/*
	 *  Maybe this doesn't generate an animal part?
	 *
	 */
	part := item_animal_part(what)
	if part == 0 {
		return
	}
	/*
	 *  Otherwise, give him that many!
	 *
	 *  Fri Sep 25 06:58:08 1998 -- Scott Turner
	 *
	 *  He gets randomly somewhere between 0 and 1/2.
	 *
	 */
	half := how_many / 2
	if half < 1 {
		half = 1
	}
	total := rnd(0, half)

	gen_item(who, part, total)
	indent += 3
	if total != 0 {
		wout(who, "%s gathers %s.", box_name(who), box_name_qty(part, total))
	}
	indent -= 3
}

/*
 *  Deduct men-as-inventory who have been killed.
 *  Then kill or wound the noble.
 *
 *  The inherit parameter determines who gets the booty from dead nobles.
 *
 *  Fri Sep 27 13:04:49 1996 -- Scott Turner
 *
 *  Modified for a nil l_b for anonymous attacks.
 *
 *  Fri Sep 11 09:40:28 1998 -- Scott Turner
 *
 *  Deduct dead should now generate "animal parts" and give them to the
 *  other stack as necessary.
 *
 *
 */
func deduct_dead(l_a []*fight, l_b []*fight, inherit int) {
	if len(l_b) == 0 || cannot_take_booty(lead_char(l_b)) {
		inherit = 0
	}

	/*
	 *  First deduct all of the dead men
	 */

	for i := 0; i < len(l_a); i++ {
		unit := l_a[i].unit
		item := l_a[i].kind

		if item > 0 {
			num_to_kill := l_a[i].sav_num - l_a[i].num

			//#if 0
			//            /*
			//             *  Wed Feb 26 08:43:59 1997 -- Scott Turner
			//             *
			//             *  You can no longer capture beasts in battle.
			//             *
			//             *
			//             */
			//            if (inherit != MATES_SILENT &&
			//                beast_capturable(l_a[i].unit) &&
			//            l_b &&
			//            side_has_skill(l_b, sk_capture_beasts))
			//            {
			//            num_to_kill = rnd(0, num_to_kill/2);
			//            }
			//#endif

			/*
			 *  The following assert has failed in the past when a unit has
			 *  someone become a member of both the attacking and defending
			 *  party.  What happens is that this routine is called twice
			 *  for the same unit.  The second time, the unit has lost thus,
			 *  and the following consistency check fails.
			 *
			 *                  assert(has_item(unit, item) == l_a[i].sav_num);
			 *
			 *  Wed Jul 10 15:28:00 1996 -- Scott Turner
			 *
			 *  This assert is no longer true, because you may not be able to use
			 *  all your men in battle, so sav_num could be less than the number
			 *  you have.  However, you should still have at least that many.
			 */

			assert(has_item(unit, item) >= l_a[i].sav_num)
			consume_item(unit, item, num_to_kill)
			/*
			 *  Fri Sep 11 09:42:12 1998 -- Scott Turner
			 *
			 *  Now possibly transfer animal parts.
			 *
			 *  Fri Sep 25 06:57:47 1998 -- Scott Turner
			 *
			 *  The loser doesn't get to gather parts.
			 *
			 */
			if inherit != MATES_SILENT {
				generate_animal_parts(lead_char(l_b), item, num_to_kill)
			}
		}
	}

	/*
	 *  If a garrison units loses all of its men, terminate it
	 *
	 *  It's not clear that this case can ever happen.
	 *  Since the garrison unit itself dies after one hit, and
	 *  generates no hits itself, if all the men die, it should
	 *  die too.
	 *
	 *  Mon Jul  1 12:43:49 1996 -- Scott Turner
	 *
	 *  We need to destroy the garrison if all the men it could control
	 *  die -- since it may have men it wasn't able to use in the battle.
	 *
	 *  Sat Dec  5 20:43:49 1998 -- Scott Turner
	 *
	 *  Actually, the garrison needs to die if no one is protecting it.
	 *
	 */

	for i := 0; i < len(l_a); i++ {
		if l_a[i].kind == FK_noble && subkind(l_a[i].unit) == sub_garrison && l_a[i].nprot < 1 {
			if l_a[i].new_health != 0 || l_a[i].num != 0 {
				log_output(LOG_CODE, "%s lost all men, zeroed out", box_name(l_a[i].unit))
			}
			l_a[i].new_health = 0
			l_a[i].num = 0
		}
	}

	/*
	 *  Now apply any wounds the nobles received, possibly killing them.
	 */

	for i := 0; i < len(l_a); i++ {
		if l_a[i].kind != FK_noble { /* Not a noble. */
			continue
		}
		who := l_a[i].unit
		if l_a[i].new_health == 0 {
			kill_char(who, inherit, S_body)
			l_a[i].num = 0
		} else if l_a[i].new_health < char_health(who) {
			assert(l_a[i].new_health > 0)
			amount := char_health(who) - l_a[i].new_health
			add_char_damage(who, amount, inherit)
			l_a[i].num = 0
		}
		/*
		 *  else:
		 *	either new_health == current_health, meaning the noble wasn't hit
		 *	or if it's greater, then they survived a fatal wound
		 */

	}
}

/*
 *  Hit wounded units with 1-100 health loss
 *  Store new health for unit (0 = dead) in .new_health.
 */
//#ifdef HERO
func determine_noble_wounds(l []*fight) {}

//#endif // HERO
//#ifndef HERO
//func determine_noble_wounds(l []*fight) {
//    for i := 0; i < len(l); i++ {
//        if (l[i].kind != FK_noble) {
//            continue;
//        }
//
//        assert(l[i].num == 0 || l[i].num == 1);
//
//        if (l[i].num)        /* not hit */
//        {
//            l[i].new_health = l[i].sav_num;
//            continue;
//        }
//
///*
// *  Already dead or undead when we started.  One hit kills an undead.
// */
//
//        if (l[i].sav_num <= 0) {
//            l[i].new_health = 0;
//        } else {
//            l[i].new_health = max(l[i].sav_num - rnd(1, 100), 0);
//        }
//    }
//}
//#endif // HERO

/*
 *  See if any would-be dead nobles have Survive fatal wound [9101].
 */

func check_fatal_survive(l []*fight) {
	for i := 0; i < len(l); i++ {
		if l[i].kind != FK_noble {
			continue
		}
		if l[i].new_health == 0 && survive_fatal(l[i].unit) {
			l[i].survived_fatal = TRUE
			l[i].new_health = 100
		}
	}
}

/*
 *  Deduct any structure damage.  If the structure is completely
 *  destroyed, it may vanish, expelling the units.  If it is a ship
 *  it may sink, killing anyone left on board.
 */
func structure_damage(l []*fight) {
	if len(l) > 0 && l[0].kind == FK_fort {
		damage := l[0].sav_num - l[0].num
		unit = l[0].unit
		assert(damage >= 0)
		if damage != 0 {
			add_structure_damage(unit, damage)
		}

		/*
		 *  Wed Oct 16 12:58:42 1996 -- Scott Turner
		 *
		 *  Take care not to permanently improve this
		 *  fortification.  If it's current defense is
		 *  greater than originally, reset it.
		 *
		 */
		if l[0].defense < p_subloc(unit).defense {
			p_subloc(unit).defense = l[0].defense
		}
	}
}

/*
 *	1:1	25%
 *	2:1	50%
 *	3:1	75%
 */
func determine_prisoners(l_a []*fight, l_b []*fight) {
	lead_a := lead_char(l_a)
	//lead_b := lead_char(l_b);

	//#if 0
	//    no_take = cannot_take_prisoners(lead_a) || char_break(lead_b) == 0;
	//#else
	no_take := cannot_take_prisoners(lead_a)
	//#endif

	//#if 0
	//    /*
	//     *  Can no longer capture beasts in battle.
	//     *
	//     */
	//    capture_beasts = side_has_skill(l_a, sk_capture_beasts);
	//#endif
	capture_beasts := false

	num_a := total_non_damage(l_a)
	num_b := total_non_damage(l_b)

	/*
	 *  If we're on a ship on the ocean, there is nowhere for the
	 *  losers to flee to, so capture them all.
	 */

	chance := 0 // chance that a unit is taken prisoner
	if l_b[0].kind == FK_fort && is_ship(l_b[0].unit) && subkind(subloc(l_b[0].unit)) == sub_ocean && l_a[lead_char_pos(l_a)].seize_slot != 0 {
		chance = 100
	} else {
		if num_b <= 0 {
			chance = 100
		} else {
			chance = 25 * num_a / num_b
		}

		if chance < 25 {
			chance = 25
		}
		if chance > 75 {
			chance = 75
		}
	}

	/*
	 *  Set prisoner flag based on chance for all nobles left alive.
	 *
	 *  If a unit can't be taken prisoner, or the winner cannot
	 *  take prisoners, then kill the would-be prisoner.
	 */

	for i := 0; i < len(l_b); i++ {
		if l_b[i].kind != FK_noble {
			continue
		}

		/*
		 *  if the attackers has the skill Capture beasts in battle [9506],
		 *  capturable beasts (ni chars whose items have the capturable flag)
		 *  shouldn't be killed, they should become the property of the victor.
		 *  take_prisoners will later remove the unit wrapper and add the ni
		 *  items to the winning character.
		 */

		if beast_capturable(l_b[i].unit) && l_b[i].new_health == 0 {
			//#if 0
			//            /*
			//             *  All capturable beasts should have their break point set at 0,
			//             *  so they will fight to the death.  Otherwise, non-beastmasters
			//             *  could take them prisoner.  The assert below will fail if the
			//             *  beast's break point was not zero.
			//             *
			//             *  We could simply force new_health to be zero here to kill the
			//             *  beast, but then it wouldn't have fought to its fullest.
			//             */
			//
			//            /*
			//             *  Ah, this is wrong in the following case:
			//             *  A human player is leading a stack containing a capturable beast.
			//             */
			//
			//                        if (l_b[i].new_health != 0)
			//                            fprintf(os.Stderr,  "%s\n", box_name_kind(l_b[i].unit));
			//                        assert(char_break(l_b[i].unit) == 0);
			//                        assert(l_b[i].new_health == 0);
			//#endif

			// todo: if we can no longer capture beasts in battle, should we remove this?
			if capture_beasts {
				l_b[i].prisoner = TRUE
				l_b[i].new_health = l_b[i].sav_num
				continue
			}
		}

		if l_b[i].new_health == 0 {
			continue
		}

		if rnd(1, 100) > chance {
			continue
		}

		if no_take {
			l_b[i].new_health = 0
		} else {
			l_b[i].prisoner = TRUE
		}
	}
}

func take_prisoners(winner int, l []*fight) {
	/*
	 *  Wed Dec  9 18:15:11 1998 -- Scott Turner
	 *
	 *  Since l[i].num == 0 no longer indicates a dead noble,
	 *  we need to change this condition.
	 *
	 */
	for i := 0; i < len(l); i++ {
		if l[i].prisoner != 0 && l[i].new_health > 0 {
			take_prisoner(winner, l[i].unit)
		}
	}
}

/*
 *  Fri Nov 26 13:26:09 1999 -- Scott Turner
 *
 *  This should only do a "promote".  Seizing a location
 *  is handled in d_move_attack.
 *
 */
func seize_position(winner int, loser_where int, loser_pos int) {

	/*
	 *  If loser is in a better spot than the winner, take it.
	 *	Loser is in a structure we want
	 *	Loser is before us in the location list
	 */

	if loser_where != subloc(winner) {
		return
	}

	if loser_pos < here_pos(winner) {
		promote(winner, loser_pos)
	}
}

func stack_flee(who, winner int) {
	assert(kind(who) == T_char)

	where := subloc(who)

	/*
	 *  If we're on a ship on an island, flee to the island
	 *  Otherwise, flee to the province.
	 */

	var to_where int
	if is_ship(where) && subkind(subloc(where)) == sub_island {
		to_where = subloc(where)
	} else {
		to_where = province(who)
	}

	/*
	 *  Since there's nowhere to flee on a ship on the ocean,
	 *  if someone is taking the ship (seize_slot) the prisoner
	 *  percentage should have been set at 100%.  No one should
	 *  be left to flee.
	 */

	//#if 1
	if subkind(to_where) == sub_ocean {
		return
	}
	//#else
	//    assert(subkind(to_where) != sub_ocean);
	//#endif

	/*
	 *  If we're already in the province, just move us to the end
	 *  of the list.
	 */

	if to_where == where {
		set_where(who, where)
		return
	}

	/*
	 *  Combat has already set up the out_vector, so squirrel it
	 *  away for this stack output
	 */

	tmp := save_output_vector()
	vector_stack(who, true)
	wout(VECT, "We flee to %s.", box_name(to_where))

	restore_output_vector(tmp)

	wout(winner, "%s flees to %s.", box_name(who), box_name(to_where))
	wout(who, "%s flees to %s.", box_name(who), box_name(to_where))

	leave_stack(who)
	move_stack(who, to_where)
}

func demote_units(winner int, l []*fight) {
	for i := 0; i < len(l); i++ {
		if l[i].kind != FK_noble ||
			l[i].prisoner != 0 ||
			!alive(l[i].unit) ||
			!(stack_leader(l[i].unit) == l[i].unit) { // todo: should be !(... == ...), maybe?
			continue
		}

		stack_flee(l[i].unit, winner)
	}
}

func combat_display_with(l []*fight) string {
	for i := 1; i < len(l); i++ {
		if l[i].kind == FK_noble {
			if l[0].kind == FK_fort {
				return ", owner:"
			}
			return ", accompanied~by:"
		}
	}

	return ""
}

func show_side_units(l []*fight) {
	show_combat_flag = true

	out(VECT, "")
	wout(VECT, "%s%s", liner_desc(l[0].unit), combat_display_with(l))

	indent += 3
	for i := 1; i < len(l); i++ {
		if l[i].kind == FK_noble {
			if l[i].ally != 0 {
				combat_ally = ", ally"
			} else {
				combat_ally = ""
			}
			wiout(VECT, 3, "%s", liner_desc(l[i].unit))
		}
	}
	indent -= 3

	show_combat_flag = false
}

func out_side(l []*fight, s string) {
	for i := 0; i < len(l); i++ {
		if l[i].kind == FK_noble {
			wout(l[i].unit, "%s", s)
		}
	}
}

/*
 *  We attack Foo!
 *  Bar leads us in an attack against Foo!
 *  Bar attacks us!
 *  Bar attacks Castle!
 */
func combat_banner(l_a []*fight, l_b []*fight) {
	assert(len(l_a) > 0)
	assert(len(l_b) > 0)

	wout(VECT, "%s attacks %s!", box_name(lead_char(l_a)), box_name(l_b[0].unit))

	indent += 3
	show_side_units(l_a)
	show_side_units(l_b)
	indent -= 3

	wout(lead_char(l_a), "Attack %s.", box_name(l_b[0].unit))

	for i := lead_char_pos(l_a) + 1; i < len(l_a); i++ {
		if l_a[i].kind == FK_noble {
			wout(l_a[i].unit, "%s leads us in an attack against %s.", box_name(lead_char(l_a)), box_name(l_b[0].unit))
		} else {
			wout(l_a[i].unit, "  Commanding %s.", box_name_qty(l_a[i].kind, l_a[i].num))
		}
	}

	for i := 0; i < len(l_b); i++ {
		if l_b[i].kind == FK_noble {
			if l_b[i].ally != 0 {
				wout(l_b[i].unit, "%s attacks %s!  We rush to the defense.", box_name(lead_char(l_a)), box_name(l_b[0].unit))
			} else {
				wout(l_b[i].unit, "%s attacks us!", box_name(lead_char(l_a)))
			}
		} else if l_b[i].kind > 0 {
			wout(l_b[i].unit, "  Commanding %s.", box_name_qty(l_b[i].kind, l_b[i].num))
		}
	}
}

/*
 *  (Foo is fine.)
 *
 *  Foo was killed.
 *  Foo was wiped out.
 *  Foo was taken prisoner.
 *
 *  Foo lost n men.
 *  N soldiers, m archers, ... were killed.
 *
 *  Leader lost n soldiers, m archers, ...
 */
func tally_side_losses(l []*fight) string {
	clear_temps(T_item)

	for i := 0; i < len(l); i++ {
		if l[i].kind > 0 {
			bx[l[i].kind].temp += l[i].sav_num - l[i].num
		}
	}

	var s string
	for _, i := range loop_item() {
		if bx[i].temp > 0 {
			s = comma_append(s, just_name_qty(i, bx[i].temp))
		}
	}
	return s
}

func tally_personal_losses(l []*fight, pos int) string {
	assert(len(l) > pos)
	assert(l[pos].kind == FK_noble)

	var s string
	for i := pos + 1; i < len(l) && l[i].unit == l[pos].unit; i++ {
		if l[i].kind > 0 && l[i].num < l[i].sav_num {
			s = comma_append(s, just_name_qty(l[i].kind, l[i].sav_num-l[i].num))
		}
	}
	return s
}

func what_happened_to_noble(l []*fight, pos int) string {
	if l[pos].prisoner != 0 && l[pos].new_health != 0 {
		/*
		   if (subkind(l[pos].unit) == sub_ni &&
		       beast_capturable(l[pos].unit))
		       return  "was captured";
		   else
		*/
		return "was taken prisoner."
	} else if l[pos].num == 0 && l[pos].new_health == 0 {
		//#if 0
		//        if (len(l) > pos+1 && l[pos+1].unit == l[pos].unit)
		//        {
		//            s = "was wiped out.";	/* noble + others */
		//        }
		//        else
		//#endif
		if l[pos].sav_num > 0 {
			return "was killed." /* noble alone */
		}
		return "was destroyed." /* "killed" for undead */
	} else if l[pos].survived_fatal != 0 {
		return "survived a fatal wound!"
	} else if l[pos].new_health == -1 {
		// undead who was not injured in battle
	} else if l[pos].num == 0 {
		return "was wounded."
	}

	return ""
}

func show_side_results(l []*fight) {
	first := true

	lead := lead_char(l)
	tally := tally_side_losses(l)

	if tally != "" {
		indent += 3
		wout(VECT, "%s lost %s.", just_name(lead), tally)
		indent -= 3
		first = false
	}

	for i := 0; i < len(l); i++ {
		if l[i].kind != FK_noble {
			continue
		}

		if s := tally_personal_losses(l, i); s != "" {
			wout(l[i].unit, "%s lost %s.", box_name(l[i].unit), s)
		}

		if s := what_happened_to_noble(l, i); s != "" {
			indent += 3
			wout(VECT, "%s %s", box_name(l[i].unit), s)
			indent -= 3
			first = false
		}
	}

	/* NOTYET:  show structure damage here */

	if !first {
		out(VECT, "")
	}
}

func best_here_pos(l []*fight, where int) int {
	best := 99999
	for i := 0; i < len(l); i++ {
		if l[i].kind != FK_noble {
			continue
		}

		if subloc(l[i].unit) != where {
			continue
		}

		assert(subloc(l[i].unit) == where)

		n := here_pos(l[i].unit)
		if n < best {
			best = n
		}
	}

	if best == 99999 {
		log_output(LOG_CODE, "best_here_pos: best == 99999, day=%d, l[0]=%s", sysclock.day, box_code_less(lead_char(l)))
	}

	return best
}

func combat_stop_movement(who int, l []*fight) {
	ship := subloc(who)
	if is_ship(ship) && ship_moving(ship) != 0 {
		interrupt_order(who)
		p_subloc(ship).moving = 0

		tmp := save_output_vector()
		vector_char_here(ship)
		wout(VECT, "Loss in battle cancels movement.")
		restore_output_vector(tmp)

		log_output(LOG_CODE, "battle interrupts sailing, who=%d, where=%d", ship, subloc(who))
		return
	}

	for i := 0; i < len(l); i++ {
		if l[i].kind == FK_noble && char_moving(l[i].unit) != 0 {
			interrupt_order(l[i].unit)

			tmp := save_output_vector()
			vector_stack(l[i].unit, true)
			wout(VECT, "Loss in battle cancels movement.")
			restore_output_vector(tmp)

			//#if 0
			//            log_output(LOG_CODE,
			//                "battle interrupts movement, who=%d, where=%d",
			//                l[i].unit, subloc(l[i].unit));
			//#endif
		}
	}
}

func reconcile(winner int, l_a []*fight, l_b []*fight) {
	determine_noble_wounds(l_a)
	determine_noble_wounds(l_b)

	if winner != 0 {
		determine_prisoners(l_a, l_b)
	}

	check_fatal_survive(l_a)
	check_fatal_survive(l_b)

	out(VECT, "")
	if winner != 0 {
		wout(VECT, "%s is victorious!", box_name(lead_char(l_a)))
		out_side(l_a, "We won!")
		out_side(l_b, "We lost!")
	} else {
		wout(VECT, "No victor emerges from the fight.")
		out_side(l_a, "Neither side prevailed.")
		out_side(l_b, "Neither side prevailed.")
	}
	out(VECT, "")

	/*
		show_side_results(l_a);
		show_side_results(l_b);
	*/

	if winner != 0 {
		winner = lead_char(l_a)
		loser := lead_char(l_b)

		loser_where := subloc(loser)
		loser_pos := best_here_pos(l_b, loser_where)

		combat_stop_movement(loser, l_b)
		clear_guard_flag(loser)

		show_side_results(l_a)
		deduct_dead(l_a, l_b, MATES_SILENT)

		show_side_results(l_b)
		deduct_dead(l_b, l_a, winner)

		take_prisoners(winner, l_b)

		if l_a[lead_char_pos(l_a)].seize_slot != 0 {
			seize_position(winner, loser_where, loser_pos)
		}

		demote_units(winner, l_b)

		if combat_sea {
			log_output(LOG_CODE, "sea combat unchecked NOTYET case, who=%s", box_name(winner))
		}

		/*
		 *  NOTYET:
		 *
		 *	If this is a sea battle:
		 *		If the loser's ship is sinking, and the victor is
		 *		on it, the victor goes back to his own ship.
		 *		(or doesn't go in the first place).
		 *
		 *		If the victor's ship is sinking, and the loser's
		 *		isn't, then the victors go over to the loser's ship.
		 *
		 *		What about other, non-combat stacks on the victor's
		 *		or the loser's ship?  Perhaps everyone on the ship
		 *		should be involved in the battle?
		 *
		 *		What do other stacks on the loser ship do when the
		 *		victor takes it over?
		 */

	} else {
		deduct_dead(l_a, l_b, MATES_SILENT)
		deduct_dead(l_b, l_a, MATES_SILENT)
	}

	/*
	 *  Sink ships, destroy castles, etc.
	 */

	structure_damage(l_a)
	structure_damage(l_b)
}

//#if 0
//static void
//increase_attack(l []*fight)
//{
//    int i;
//
//    for i := 0; i < len(l); i++
//        if (l[i].kind == FK_noble)
//            p_char(l[i].unit).attack++;
//}
//
//
//static void
//increase_defense(l []*fight)
//{
//    int i;
//
//    for i := 0; i < len(l); i++
//        if (l[i].kind == FK_noble)
//            p_char(l[i].unit).defense++;
//}
//#endif

/*
 *  Tue Jun 25 11:53:57 1996 -- Scott Turner
 *
 *  Modified to give each side only 4 hits in a row, at most.
 *
 *  Mon Feb 24 11:40:06 1997 -- Scott Turner
 *
 *  Note that l_b should always be the defender!
 *
 *  Mon Feb 24 11:56:03 1997 -- Scott Turner
 *
 *  Allies added in second round.
 *
 *  Tue Feb 25 11:42:14 1997 -- Scott Turner
 *
 *  Since run_combat can now modify l_a and l_b, they need to come
 *  in as pointers...
 *
 */

func run_combat(lap, lbp *[]*fight) int {
	l_a, l_b := *lap, *lbp
	lead_a, lead_b := lead_char(l_a), lead_char(l_b)
	first := false

	out(combat_pl, "")
	out(combat_pl, "Combat between %s and %s on day %d", box_name(l_a[0].unit), box_name(l_b[0].unit), sysclock.day)
	out(combat_pl, "")

	dump_fighters(l_a)
	dump_fighters(l_b)

	thresh_a := total_combat_sum(l_a) * char_break(lead_a) / 100
	thresh_b := total_combat_sum(l_b) * char_break(lead_b) / 100

	// the defender gets an initial missile round.
	round = 0
	defense_side = l_b
	wout(VECT, " ")
	if total_attackers(l_b, l_a, MISSILE) > 0 {
		wout(VECT, "Initial missile round for the defense (%s):", box_name(lead_char(l_b)))
		tmp := attack_round(l_b, l_a, MISSILE)
		wout(VECT, "    %s missile%s hit%s!", cap_(nice_num(tmp)), or_string(tmp > 1, "s", ""), or_string(tmp > 1, "", "s"))
		resolve_dead(l_a)
		wout(VECT, " ")
	}

	// perhaps we're done?
	num_a := total_combat_sum(l_a)
	num_b := total_combat_sum(l_b)

	for num_a > thresh_a && num_b > thresh_b {
		round++
		out(combat_pl, "Combat round: %d", round)
		wout(VECT, "Combat round: %s", nice_num(round))

		/*
		 *  Mon Feb 24 11:51:52 1997 -- Scott Turner
		 *
		 *  On round 2, add the defender's allies, if any.
		 *
		 */
		if round == 2 {
			// attacker is the attacking unit.
			attacker := l_a[lead_char_pos(l_a)].unit

			// If a fort, call in the defenders for both the fort and the owner
			// of the fort.  Otherwise, just call in the unit's defenders.
			if l_b[0].kind == FK_fort {
				look_for_allies(lbp, subloc(l_b[0].unit), l_b[0].unit, l_b[1].unit, attacker, TRUE, l_a)
				l_b = *lbp
				look_for_allies(lbp, subloc(l_b[1].unit), l_b[0].unit, l_b[1].unit, attacker, TRUE, l_a)
			} else {
				look_for_allies(lbp, subloc(l_b[0].unit), l_b[0].unit, l_b[0].unit, attacker, TRUE, l_a)
			}
			// reset the pointers.
			l_a, l_b = *lap, *lbp

			ready_fight_list(l_b)
			wout(combat_pl, "  Old threshold: %d.", thresh_b)
			thresh_b = special_total_combat_sum(l_b) * char_break(lead_b) / 100
			wout(combat_pl, "  New threshold w/ allies: %d.", thresh_b)
		}

		/*
		 *  Special weapons.
		 *
		 *  Tue Feb 18 15:41:27 1997 -- Scott Turner
		 *
		 *  Random order, since these may do more than just kill people.
		 *
		 *  Thu Jun 19 15:24:54 1997 -- Scott Turner
		 *
		 *  Don't start special attacks till round 2.
		 *
		 */
		special_attack_banner = false
		if rnd(1, 2) == 1 {
			out(combat_pl, "  Special Phase: %s", box_name(l_a[0].unit))
			_ = special_attacks(l_a, l_b, SPECIAL)
			out(combat_pl, "  Special Phase: %s", box_name(l_b[0].unit))
			_ = special_attacks(l_b, l_a, SPECIAL)
		} else {
			out(combat_pl, "  Special Phase: %s", box_name(l_b[0].unit))
			_ = special_attacks(l_b, l_a, SPECIAL)
			out(combat_pl, "  Special Phase: %s", box_name(l_a[0].unit))
			_ = special_attacks(l_a, l_b, SPECIAL)
		}
		resolve_dead(l_a)
		resolve_dead(l_b)
		//perhaps we're done?
		num_a = total_combat_sum(l_a)
		num_b = total_combat_sum(l_b)
		if num_a <= thresh_a || num_b <= thresh_b {
			break
		}

		// missile weapons.
		first = true
		if tmp := attack_round(l_a, l_b, MISSILE); tmp != NO_ATTACKS {
			if first {
				first = false
				wout(VECT, "  Missile phase:")
			}
			wout(VECT, "    %s's forces hit with %s missile%s.", just_name(lead_char(l_a)), nice_num(tmp), or_string(tmp > 1, "s", ""))
		}
		if tmp := attack_round(l_b, l_a, MISSILE); tmp != NO_ATTACKS {
			if first {
				first = false
				wout(VECT, "  Missile phase:")
			}
			wout(VECT, "    %s's forces hit with %s missile%s.", just_name(lead_char(l_b)), nice_num(tmp), or_string(tmp > 1, "s", ""))
		}
		resolve_dead(l_a)
		resolve_dead(l_b)

		// perhaps we're done?
		num_a = total_combat_sum(l_a)
		num_b = total_combat_sum(l_b)
		if num_a <= thresh_a || num_b <= thresh_b {
			break
		}

		// melee weapons.
		wout(VECT, "  Melee phase:")
		if tmp := attack_round(l_a, l_b, MELEE); tmp != NO_ATTACKS {
			wout(VECT, "    %s's forces hit %s time%s.", just_name(lead_char(l_a)), nice_num(tmp), or_string(tmp > 1, "s", ""))
		} else {
			wout(VECT, "    %s's forces miss!", just_name(lead_char(l_a)))
		}
		if tmp := attack_round(l_b, l_a, MELEE); tmp != NO_ATTACKS {
			wout(VECT, "    %s's forces hit %s time%s.", just_name(lead_char(l_b)), nice_num(tmp), or_string(tmp > 1, "s", ""))
		} else {
			wout(VECT, "    %s's forces miss!", just_name(lead_char(l_b)))
		}
		resolve_dead(l_a)
		resolve_dead(l_b)

		num_a = total_combat_sum(l_a)
		num_b = total_combat_sum(l_b)

		// now remove any scry_offense effects
		for i := 0; i < len(l_a); i++ {
			if l_a[i].kind == FK_noble {
				delete_effect(l_a[i].unit, ef_scry_offense, 0)
			}
		}
		for i := 0; i < len(l_b); i++ {
			if l_b[i].kind == FK_noble {
				delete_effect(l_b[i].unit, ef_scry_offense, 0)
			}
		}
	}
	/*
	 *  Wed Feb 19 09:27:38 1997 -- Scott Turner
	 *
	 *  Now we might have to fix up any soldiers who were "raised" and should be dead.
	 *
	 */
	for i := 0; i < len(l_a); i++ {
		if l_a[i].tmp_num2 != 0 {
			if l_a[i].num < l_a[i].tmp_num2 {
				l_a[i].tmp_num2 = l_a[i].num
			}
			l_a[i].num -= l_a[i].tmp_num2
			out(combat_pl, "%d %s return to the dead.", l_a[i].tmp_num2, box_name(l_a[i].unit))
		}
	}
	for i := 0; i < len(l_b); i++ {
		if l_b[i].tmp_num2 != 0 {
			if l_b[i].num < l_b[i].tmp_num2 {
				l_b[i].tmp_num2 = l_b[i].num
			}
			l_b[i].num -= l_b[i].tmp_num2
			out(combat_pl, "%d %s return to the dead.", l_b[i].tmp_num2, box_name(l_b[i].unit))
		}
	}

	/*
	 *  Who won?
	 *
	 *	The loser must be below their break threshold,
	 *	the winner must still be above theirs,
	 *	and the lead noble of the winner must not have been wounded.
	 *
	 *	(The lead noble will be the last figure to take a wound,
	 *	so if the lead noble is wounded, it means there's no one
	 *	left unhurt on the winner's side)
	 *
	 *	If both sides are below their break threshold, it's a draw.
	 */
	if num_a <= thresh_a && num_b > thresh_b {
		assert(l_b[lead_char_pos(l_b)].num > 0)
		return B_WON
	}

	if num_b <= thresh_b && num_a > thresh_a {
		assert(l_a[lead_char_pos(l_a)].num > 0)
		return A_WON
	}

	return TIE
}

func combat_top(lap, lbp *[]*fight) bool {
	l_a, l_b := *lap, *lbp

	print_dot('*')

	//#if 0
	//        where := subloc(l_b[0].unit);
	//#else
	where := viewloc(subloc(l_a[0].unit))
	where2 := viewloc(subloc(l_b[0].unit))
	//#endif

	vector_clear()
	vector_add(where)
	if where2 != where {
		vector_add(where2)
	}

	show_to_garrison = true

	assert(len(l_a) > 0)
	assert(len(l_b) > 0)

	combat_banner(l_a, l_b)

	var result int
	result = run_combat(lap, lbp)

	// reset the pointers.
	l_a, l_b = *lap, *lbp

	switch result {
	case A_WON:
		reconcile(TRUE, l_a, l_b)
		//#if 0
		//            increase_attack(l_a);
		//#endif
		break

	case B_WON:
		reconcile(TRUE, l_b, l_a)
		//#if 0
		//            increase_defense(l_b);
		//#endif
		break

	case TIE:
		reconcile(FALSE, l_a, l_b)
		break

	default:
		panic("!reached")
	}

	show_to_garrison = false

	return result == A_WON
}

func fail_defeat_check(a int, l_b []*fight) int {
	lead_b := lead_char(l_b)

	n := only_defeatable(lead_b)
	if n == 0 {
		return DEFEAT_DEFAULT
	}

	if has_item(a, n) != FALSE {
		return DEFEAT_FORCE
	}

	out(a, "Only one who possesses %s can defeat %s.", box_name(n), box_name(lead_b))

	return DEFEAT_UNREADY
}

/*
 *  Thu Oct 24 14:43:16 1996 -- Scott Turner
 *
 *  Modified to handle waits of more than 1 day.  Note that we now have
 *  to copy and discard to avoid the delete problem.
 *
 */

func clear_second_waits() {
	ll_l := ilist_copy(second_wait_list)

	for i := 0; i < len(ll_l); i++ {
		rp_command(ll_l[i]).second_wait--
		if rp_command(ll_l[i]).second_wait < 1 {
			rp_command(ll_l[i]).second_wait = 0
			second_wait_list = rem_value(second_wait_list, ll_l[i])
		}
	}
}

func set_second_waits(l []*fight, already_waiting int) {
	/*
	 *  already_waiting is the id of a character which issued
	 *  the ATTACK order.  It is already paying the day cost for
	 *  the attack.  Its stackmates and the defenders are not, so
	 *  we set second_wait for them.
	 */

	for i := 0; i < len(l); i++ {
		if l[i].kind == FK_noble && l[i].unit != already_waiting {
			c := p_command(l[i].unit)
			if c.second_wait == FALSE {
				c.second_wait = TRUE
				second_wait_list = append(second_wait_list, l[i].unit)
			}
		}
	}
}

func set_weather(where, where2 int) {
	combat_rain = weather_here(where, sub_rain) != FALSE
	combat_sea = where != where2 && subkind(province(where)) == sub_ocean
	combat_swampy = subkind(province(where)) == sub_swamp
	combat_wind = weather_here(where, sub_wind) != FALSE
}

func regular_combat(a int, b int, seize_slot int, already_waiting int) int {
	assert(a != b)

	l_a := construct_fight_list(a, b, FALSE)
	l_b := construct_fight_list(b, a, TRUE)
	if l_b == nil {
		return A_WON
	}

	set_weather(subloc(lead_char(l_b)), subloc(lead_char(l_a)))
	ready_fight_list(l_a)
	ready_fight_list(l_b)

	set_second_waits(l_a, already_waiting)

	/*
	 *  Set the location of the defenders, which is used in
	 *  resolve_hit to determine terrain defense bonuses.
	 *
	 *  Note that a castle in the mountains will return the
	 *  castle, and so you won't get the terrain bonus.  This
	 *  seems reasonable (consider a city in the woods) but if
	 *  desired, changing this to province() would give the
	 *  other behavior.
	 *
	 */
	combat_def_loc = subloc(b)

	/*
	   *  We don't force side b to wait as well.
	    *  This is on purpose (a playability choice)
	   *  set_second_waits(l_b, already_waiting);
	*/

	lead_a := lead_char(l_a)

	if is_loc_or_ship(b) && len(l_b) < 1 {
		out(lead_a, "%s is unoccupied.", box_name(b))

		if seize_slot != 0 {
			seize_position(lead_a, b, 0)
		}

		return NO_COMBAT
	}

	if defeat_flag := fail_defeat_check(a, l_b); defeat_flag == DEFEAT_UNREADY {
		return TIE
	}

	assert(len(l_a) > 0)
	assert(len(l_b) > 0)

	if seize_slot != 0 {
		l_a[lead_char_pos(l_a)].seize_slot = TRUE
	}

	/*
	 *  Note that seize_slot is always false for the defender.
	 */

	if combat_top(&l_a, &l_b) {
		return A_WON
	}
	return FALSE
}

/*
 *  Who can we attack?
 *
 *	another character here
 *	a building or subloc of distance <= 1 that we can get to
 *	    this should cover castles, ships, and sublocs
 *	    while the distance requirement prevents inter-province attacks.
 *
 *  Thus, we can't attack outside characters from inside of a building.
 */

func select_target(c *command) int {
	target := c.a
	//where := subloc(c.who)

	if kind(target) == T_deadchar {
		wout(c.who, "%s is not here.", c.parse[1])
		return 0
	}

	if kind(target) == T_char {
		//#if 0
		//        if (subloc(c.who) != subloc(target))
		//        {
		//            wout(c.who, "%s is not here.", c.parse[1]);
		//            return 0;
		//        }
		//#else
		if !check_char_here(c.who, target) {
			return 0
		}
		//#endif

		if is_prisoner(target) {
			wout(c.who, "Cannot attack prisoners.")
			return 0
		}

		if c.who == target {
			wout(c.who, "Can't attack oneself.")
			return 0
		}

		if stack_leader(c.who) == stack_leader(target) {
			wout(c.who, "Can't attack a member of the same stack.")
			return 0
		}

		return stack_leader(target)
	}

	v := parse_exit_dir(c, subloc(c.who), "attack")
	if v == nil {
		return 0
	} else if v.direction != DIR_IN {
		wout(c.who, "%s cannot be attacked from here.",
			box_name(v.destination))
		return 0
	} else if v.distance > 1 {
		wout(c.who, "%s is too far away to attack from here.",
			box_name(v.destination))
		return 0
	} else if v.impassable != FALSE {
		wout(c.who, "%s is impassable.", box_name(v.destination))
		return 0
	}
	return v.destination
}

func select_attacker(who, target int) int {
	where := subloc(who)
	if loc_depth(where) == LOC_build && building_owner(where) == who && FALSE == somewhere_inside(where, target) {
		/*
		 *  This code is apparently for naval combat.
		 *
		 *  One would attack another ship by specifying the direction link to the other ship.
		 */
		who = where
	}

	return who
}

//#if 0
//int
//v_attack(struct command *c)
//{
//    int attacker;
//    int target;
//    int targ_who;
//    flag := c.b
//    int seize_slot;
//
//    if (in_safe_now(c.who))
//    {
//        wout(c.who, "Combat is not permitted in safe havens.");
//        return FALSE;
//    }
//
//    if (stack_leader(c.who) != c.who)
//    {
//        wout(c.who, "Only the stack leader may initiate combat.");
//        return FALSE;
//    }
//
//    target = select_target(c);
//    if (target <= 0)
//        return FALSE;
//
//    attacker = select_attacker(c.who, target);
//    if (attacker <= 0)
//        return FALSE;
//
//    if (is_loc_or_ship(target))
//    {
//        if (loc_depth(target) == LOC_build)
//            targ_who = building_owner(target);
//        else
//            targ_who = first_character(target);
//    }
//    else
//        targ_who = target;
//
//    if (player(c.who) == player(targ_who))
//    {
//        wout(c.who, "Units in the same faction may not engage "
//                "in combat.");
//        return FALSE;
//    }
//
//    if (flag)
//        seize_slot = FALSE;
//    else
//        seize_slot = TRUE;
//
//    regular_combat(attacker, target, seize_slot, c.who);
//
//    return TRUE;
//}
//#endif

func loc_guarded(where, except int) int {
	for _, i := range loop_here(where) {
		if kind(i) == T_char && char_guard(i) != 0 && player(i) != except {
			return i
		}
	}
	return 0
}

func attack_guard_units(a, b int) int {
	l_a := construct_fight_list(a, b, FALSE)
	l_b := construct_guard_fight_list(b, a, l_a, TRUE)

	set_weather(subloc(lead_char(l_b)), subloc(lead_char(l_a)))

	ready_fight_list(l_a)
	ready_fight_list(l_b)

	if len(l_b) <= 0 { /* no guards */
		return TRUE
	}

	assert(len(l_a) > 0)

	if combat_top(&l_a, &l_b) {
		return TRUE
	}
	return FALSE
}

func v_pillage(c *command) int {
	where := subloc(c.who)
	//has := has_item(where, item_peasant);
	men := count_fighters(c.who, item_attack(item_skirmisher))
	flag := c.a != FALSE

	if men < 10 {
		wout(c.who, "At least 10 fighters are needed to pillage.")
		return FALSE
	}

	if in_safe_now(c.who) != FALSE {
		wout(c.who, "Pillaging is not permitted in safe havens.")
		return FALSE
	}
	if stack_leader(c.who) != c.who {
		wout(c.who, "Only the stack leader may pillage.")
		return FALSE
	}

	if guard := loc_guarded(where, player(c.who)); guard != 0 {
		wout(c.who, "%s is protected by guards.", box_name(where))
		if !flag || in_safe_now(c.who) != FALSE {
			return FALSE
		}
		if FALSE == attack_guard_units(c.who, guard) {
			return FALSE
		}
	}

	c.h = 0

	return TRUE
}

func end_pillage(c *command) {
	where := province(subloc(c.who))
	mob := 0
	wout(c.who, "Pillaging yielded %s.", gold_s(c.h))

	if subkind(where) == sub_city {
		wout(where, "%s loots the city.", box_name(c.who))
	} else {
		wout(where, "%s loots the countryside.", box_name(c.who))
	}

	if rnd(1, 3) == 1 {
		mob = create_peasant_mob(where)
	}

	if mob != 0 {
		wout(c.who, "%s has formed to resist pillaging.", liner_desc(mob))
		wout(where, "%s has formed to resist pillaging.", liner_desc(mob))

		if rnd(1, 3) == 1 {
			queue(mob, "attack %s", box_code_less(c.who))
		}
	}
}

func d_pillage(c *command) int {
	where := province(subloc(c.who))
	flag := c.a != FALSE
	//int guard;
	//extern int gold_pillage;
	men := count_fighters(c.who, item_attack(item_peasant))
	//int amount;

	if men < 10 {
		wout(c.who, "At least 10 fighters are needed to pillage.")
		end_pillage(c)
		return FALSE
	}

	if has_item(where, item_peasant) < 100 {
		wout(c.who, "Too few peasants to pillage.")
		end_pillage(c)
		return FALSE
	}

	if in_safe_now(c.who) != FALSE {
		wout(c.who, "Pillaging is not permitted in safe havens.")
		end_pillage(c)
		return FALSE
	}

	if stack_leader(c.who) != c.who {
		wout(c.who, "Only the stack leader may pillage.")
		end_pillage(c)
		return FALSE
	}

	if guard := loc_guarded(where, player(c.who)); guard != 0 {
		wout(c.who, "%s is protected by guards.", box_name(where))
		if !flag || in_safe_now(c.who) != FALSE {
			end_pillage(c)
			return FALSE
		}
		if FALSE == attack_guard_units(c.who, guard) {
			end_pillage(c)
			return FALSE
		}
	}

	amount := men / 10
	if has_item(where, item_peasant)-amount < 100 {
		amount = has_item(where, item_peasant) - 99
	}
	assert(consume_item(where, item_peasant, amount))

	amount *= 2
	if has_item(where, item_gold) < amount {
		amount = has_item(where, item_gold)
	}
	assert(move_item(where, c.who, item_gold, amount))
	c.h += amount
	gold_pillage += amount

	p_subloc(province(where)).loot = 1

	if c.wait == FALSE {
		end_pillage(c)
	}
	return TRUE
}

func v_guard(c *command) int {
	flag := c.a != FALSE
	where := subloc(c.who)

	if flag {
		p_char(c.who).guard = TRUE
		wout(c.who, "Will guard %s.", box_name(where))
		return TRUE
	}

	p_char(c.who).guard = FALSE
	wout(c.who, "Will not guard %s.", box_name(where))

	return TRUE
}

func auto_attack(who, target int) {
	out(who, "> [auto-attack %s]", box_code_less(target))
	regular_combat(who, target, FALSE, 0)
}

func check_auto_attack_sup(who int) {
	where := subloc(who)
	okay := true

	if in_safe_now(who) != FALSE { /* safe haven, no combat permitted */
		return
	}

	if loc_depth(where) == LOC_province && weather_here(where, sub_fog) != FALSE {
		return
	}

	if char_gone(who) != FALSE {
		return
	}

	for _, i := range loop_all_here(where) {
		if FALSE == is_hostile(who, i) {
			continue
		}

		var target int
		if kind(i) == T_char {
			if is_prisoner(i) {
				continue
			}

			//#if 1
			if !char_here(who, i) {
				continue
			}
			//#else
			//                    if (char_really_hidden(i))
			//                        continue;
			//#endif

			target = stack_leader(i)
			if stack_leader(who) == target {
				continue
			}

			n := only_defeatable(target)
			if n != FALSE && FALSE == has_item(who, n) {
				continue
			}
		} else if is_loc_or_ship(i) {
			if kind(i) == T_loc &&
				loc_hidden(i) &&
				!test_known(who, i) {
				continue
			}

			target = i

			var targ_who int
			if loc_depth(i) == LOC_build {
				targ_who = building_owner(target)
			} else {
				targ_who = first_character(target)
			}

			if targ_who == 0 ||
				player(who) == player(targ_who) {
				continue
			}
		} else {
			continue
		}

		auto_attack(who, target)
		okay = false
		break /* only get to attack first target */
	}

	if okay && is_ship(where) {
		outer := subloc(where)

		for _, i := range loop_here(outer) {
			if i == where || !is_ship(i) {
				continue
			}

			if FALSE == is_hostile(who, i) {
				continue
			}

			if building_owner(i) == 0 {
				continue
			}

			auto_attack(who, i)
			okay = false
			break /* only get to attack first target */
		}
	}
}

func check_all_auto_attacks() {
	for _, i := range loop_char() {
		if stack_parent(i) != 0 { /* must be stack leader to initiate */
			continue
		} /* auto-attack */

		if char_health(i) != 100 && char_health(i) != -1 {
			continue
		}

		_ = rp_command(i)
		check_auto_attack_sup(i)
	}
}

/*
 *  PRAC_CONTROL
 *  Thu Jun 27 11:34:04 1996 -- Scott Turner
 *
 *  Practicing your "control men in battle" skill, i.e., explicit
 *  USE.
 *
 */
func v_prac_control(c *command) int {
	ret := oly_parse_s(c, sout("practice %d", sk_control_battle))
	assert(ret)

	return v_practice(c)

	//#if 0
	//    struct skill_ent *p;
	//    /*
	//     *  This keeps track of the uses of the skill?
	//     *
	//     */
	//    p_skill(sk_control_battle).use_count++;
	//
	//    /*
	//     *  Increment the user's skill.  We should like to use
	//     *  add_skill_experience here, but that only allows one increment
	//     *  per month.  Why is that?
	//     *
	//     */
	//
	//    p = rp_skill_ent(c.who, sk_control_battle);
	//
	//    /*
	//     *  Don't increase the experience if we don't actually know the
	//     *  skill.  For instance, use through a scroll or book should
	//     *  not add experience to the character, unless the character
	//     *  knows the skill himself.
	//     */
	//
	//    if (p == nil) return FALSE;
	//
	//    p.experience++;
	//    p.exp_this_month = TRUE;
	//
	//    wout(c.who, "You may now control %d men in battle.", p.experience);
	//    return TRUE;
	//#endif
}

/*
 *  ATTACK_TACTICS
 *  Thu Jun 27 11:34:04 1996 -- Scott Turner
 *
 *  Practicing your attack_tactics skill
 *
 */
func v_attack_tactics(c *command) int {
	p := rp_skill_ent(c.who, sk_attack_tactics)

	/*
	 *  Maybe he's too experienced already to benefit from more study?
	 *
	 */
	if (1.0 + float64(p.experience)*TACTICS_FACTOR) >= TACTICS_LIMIT {
		wout(c.who, "You will not benefit from further practice of this skill.")
		return FALSE
	}

	ret := oly_parse_s(c, sout("practice %d", sk_attack_tactics))
	assert(ret)

	return v_practice(c)
}

//#if 0
//int d_attack_tactics(struct command *c)
//{
//  struct skill_ent *p;
//
//  /*
//   *  This keeps track of the uses of the skill?
//   *
//   */
//  p_skill(sk_attack_tactics).use_count++;
//
//  /*
//   *  Increment the user's skill.  We should like to use
//   *  add_skill_experience here, but that only allows one increment
//   *  per month.  Why is that?
//   *
//   */
//
//  p = rp_skill_ent(c.who, sk_attack_tactics);
//
//  /*
//   *  Don't increase the experience if we don't actually know the
//   *  skill.  For instance, use through a scroll or book should
//   *  not add experience to the character, unless the character
//   *  knows the skill himself.
//   */
//
//  if (p == nil) return FALSE;
//
//  p.experience++;
//  p.exp_this_month = TRUE;
//
//  wout(c.who, "Your attack tactics multiplier is now %f.",
//       p.experience * TACTICS_LIMIT);
//  return TRUE;
//}
//#endif

/*
 *  DEFENSE_TACTICS
 *  Thu Jun 27 11:34:04 1996 -- Scott Turner
 *
 *  Practicing your defense_tactics skill
 *
 */
//#if 0
//int v_defense_tactics(struct command *c)
//{
//  struct skill_ent *p;
//  p = rp_skill_ent(c.who, sk_defense_tactics);
//
//  /*
//   *  Maybe he's too experienced already to benefit from more
//   *  study?
//   *
//   */
//  if ((1.0 + p.experience * TACTICS_FACTOR) >= TACTICS_LIMIT) {
//    wout(c.who, "You will not benefit from further practice of this skill.");
//    return FALSE;
//  }
//
//  return TRUE;
//}
//#endif

func v_defense_tactics(c *command) int {
	p := rp_skill_ent(c.who, sk_defense_tactics)
	if (1.0 + float64(p.experience)*TACTICS_FACTOR) >= TACTICS_LIMIT {
		wout(c.who, "You will not benefit from further practice of this skill.")
		return FALSE
	}

	ret := oly_parse_s(c, sout("practice %d", sk_defense_tactics))
	assert(ret)

	return v_practice(c)

	//#if 0
	//    struct skill_ent *p;
	//
	//    /*
	//     *  This keeps track of the uses of the skill?
	//     *
	//     */
	//    p_skill(sk_defense_tactics).use_count++;
	//
	//    /*
	//     *  Increment the user's skill.  We should like to use
	//     *  add_skill_experience here, but that only allows one increment
	//     *  per month.  Why is that?
	//     *
	//     */
	//
	//    p = rp_skill_ent(c.who, sk_defense_tactics);
	//
	//    /*
	//     *  Don't increase the experience if we don't actually know the
	//     *  skill.  For instance, use through a scroll or book should
	//     *  not add experience to the character, unless the character
	//     *  knows the skill himself.
	//     */
	//
	//    if (p == nil) return FALSE;
	//
	//    p.experience++;
	//    p.exp_this_month = TRUE;
	//
	//    wout(c.who, "Your defense tactics multiplier is now %f.",
	//         p.experience * TACTICS_LIMIT);
	//    return TRUE;
	//#endif

}

/*
 *  Trap_Attack
 *  Fri Sep 27 12:52:00 1996 -- Scott Turner
 *
 *  Water elemental, boulder trap, etc.
 *
 */
func do_trap_attack(who int, num_attacks int, chance int) int {
	//int i, lost = 0, num_defend, man;
	//struct fight **l_a;
	//struct fight *g;

	l_a := construct_fight_list(who, 0, TRUE)
	ready_fight_list(l_a)

	for i := 0; i < num_attacks; i++ {
		num_defend := total_defenders(l_a, nil, MELEE)
		assert(num_defend > 0)
		man := rnd(1, num_defend)
		g := find_defender(l_a, man, nil, MELEE)
		/*
		 *  Did we get him?
		 *
		 */
		if rnd(1, 100) < chance {
			decrement_num(l_a, nil, g)
		}
	}

	determine_noble_wounds(l_a)
	check_fatal_survive(l_a)
	show_side_results(l_a)
	deduct_dead(l_a, nil, 0)

	return 0 // todo: should this return a value?
}

/*
 *  USE_BEASTS
 *  Thu Oct  3 12:00:26 1996 -- Scott Turner
 *
 *  Practicing your "control beasts in battle" skill
 *
 */
func v_use_beasts(c *command) int {
	ret := oly_parse_s(c, sout("practice %d", sk_use_beasts))
	assert(ret)

	return v_practice(c)

	//#if 0
	//    struct skill_ent *p;
	//    /*
	//     *  This keeps track of the uses of the skill?
	//     *
	//     */
	//    p_skill(sk_use_beasts).use_count++;
	//
	//    /*
	//     *  Increment the user's skill.  We should like to use
	//     *  add_skill_experience here, but that only allows one increment
	//     *  per month.  Why is that?
	//     *
	//     */
	//    p = rp_skill_ent(c.who, sk_use_beasts);
	//
	//    /*
	//     *  Don't increase the experience if we don't actually know the
	//     *  skill.  For instance, use through a scroll or book should
	//     *  not add experience to the character, unless the character
	//     *  knows the skill himself.
	//     */
	//    if (p == nil) return FALSE;
	//
	//    p.experience++;
	//    p.exp_this_month = TRUE;
	//    wout(c.who, "You may now control %d beasts in battle.", p.experience);
	//
	//    return TRUE;
	//#endif
}

/*
 *  PRAC_PROTECT
 *  Tue May 23 06:55:55 2000 -- Scott Turner
 *
 *  Synonym for "practice 1990"
 *
 */
func v_prac_protect(c *command) int {
	ret := oly_parse_s(c, sout("practice %d", sk_protect_noble))
	assert(ret)

	return v_practice(c)
}
