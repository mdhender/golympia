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
	"sort"
)

/*  check.c -- check database integrity and effect minor repairs */

/*
 *  1.  Go through every box.  If box claims it's in a location,
 *	but it doesn't show up in the here list of the location,
 *	then add it to the here list.
 *
 *  2.	Go through every location.  If the loc has a here list, see
 *	if each box in the list claims that its in that location.
 *	If not, remove it from the here list.
 *
 *  This scheme gives precendence to the location the unit claims
 *  to be in over the location's here-list of units.  If they
 *  disagree, we correct the database based on where the unit
 *  claims to be.
 */

func check_here() {
	for _, i := range loop_boxes() {
		where := loc(i)
		/*
		 *  Don't fix regions, because we play some games
		 *  with the cloudlands region's "where".
		 *
		 */
		if where > 0 &&
			subkind(i) != sub_region &&
			!in_here_list(where, i) {
			log.Printf("\tcheck_here: adding [%d] to here list of [%d]\n", i, where)

			add_to_here_list(where, i)
		}
	}

	for _, i := range loop_boxes() {
		for _, j := range loop_here(i) {
			where := loc(j)
			if where != i {
				log.Printf("\tcheck_here: removing [%d] from here list of [%d]\n", j, i)

				remove_from_here_list(i, j)
			}
		}
	}
}

/*
 *  1.	For every box, if that box is in a faction, then make sure
 *	the box appears in the faction's unit list.
 *
 *  2.	For every faction, go through the unit list seeing if
 *	the units actually claim to be in the faction.  If not,
 *	remove them from the faction's unit list.
 *
 *  This scheme gives precedence to the Faction attribute
 *  over the player_ent's chars list.  If they disagree,
 *  the database will be corrected according to the what faction
 *  the unit claims to be in, over the faction's list of units.
 */

func check_swear() {
	for _, i := range loop_char() {
		over := player(i)
		if over > 0 && !is_unit(over, i) {
			log.Printf("\tcheck_swear: adding [%d] to player [%d]\n", i, over)
			p_player(over).Units = append(p_player(over).Units, i)

			sort.Ints(p_player(over).Units)
		}
	}

	for _, i := range loop_player() {
		for _, j := range loop_units(i) {
			over := player(j)
			if over != i {
				log.Printf("\tcheck_swear: removing [%d] from player list of [%d]\n", j, i)

				p_player(i).Units = rem_value(p_player(i).Units, j)
			}
		}
	}
}

func check_indep() {
	if bx[indep_player] == nil {
		log.Printf("\tcheck_indep: creating independent player [%d]\n", indep_player)

		alloc_box(indep_player, T_player, sub_pl_npc)
	}

	assert(kind(indep_player) == T_player)

	if name(indep_player) == "" {
		set_name(indep_player, "Independent player")
	}

	for _, i := range loop_char() {
		if player(i) == 0 {
			log.Printf("\tcheck_indep: swearing unit [%d] to %s\n", i, box_name(indep_player))

			set_lord(i, indep_player, LOY_unsworn, 0)
		}
	}
}

func check_gm() {

	if bx[gm_player] == nil {
		log.Printf("\tcheck_gm: creating gm player [%d]\n", gm_player)

		alloc_box(gm_player, T_player, sub_pl_system)
	}

	assert(kind(gm_player) == T_player)

	if name(gm_player) == "" {
		set_name(gm_player, "Gamemaster")
	}
}

func check_deserted() {

	if bx[deserted_player] == nil {
		log.Printf("\tcheck_deserted: creating deserted player [%d]\n", deserted_player)

		alloc_box(deserted_player, T_player, sub_pl_system)
	}

	assert(kind(deserted_player) == T_player)

	if name(deserted_player) == "" {
		set_name(deserted_player, "Deserted Nobles")
	}
}

func check_skill_player() {

	if bx[skill_player] == nil {
		log.Printf("\tcheck_skill_player: creating skill player [%d]\n", skill_player)

		alloc_box(skill_player, T_player, sub_pl_system)
	}

	assert(kind(skill_player) == T_player)

	if name(skill_player) == "" {
		set_name(skill_player, "Skill list")
	}
}

func check_eat_player() {

	if bx[eat_pl] == nil {
		log.Printf("\tcheck_eat_player: creating eat player [%d]\n", eat_pl)

		alloc_box(eat_pl, T_player, sub_pl_system)
	}

	assert(kind(eat_pl) == T_player)

	if name(eat_pl) == "" {
		set_name(eat_pl, "Order eater")
	}
}

func check_npc_player() {

	if bx[npc_pl] == nil {
		log.Printf("\tcheck_npc_player: creating npc player [%d]\n", npc_pl)

		alloc_box(npc_pl, T_player, sub_pl_silent)
	}

	assert(kind(npc_pl) == T_player)

	if name(npc_pl) == "" {
		set_name(npc_pl, "NPC control")
	}
}

func check_garr_player() {

	if bx[garr_pl] == nil {
		log.Printf("\tcheck_garr_player: creating garrison player [%d]\n", garr_pl)

		alloc_box(garr_pl, T_player, sub_pl_silent)
	}

	assert(kind(garr_pl) == T_player)

	if name(garr_pl) == "" {
		set_name(garr_pl, "Garrison units")
	}
}

/*
 *  1.	Check that T_MAX and kind_s agree
 *  2.  Check that SUB_MAX and subkind_s agre
 */

func check_glob() {
	var i int

	for i = 1; kind_s[i] != ""; i++ {
		//
	}
	assert(i == T_MAX)

	for i = 1; subkind_s[i] != ""; i++ {
		//
	}
	assert(i == SUB_MAX)
}

func check_nowhere() {
	/*
	 *  Not thorough enough?  What about other entity types?  sublocs, etc.
	 */

	for _, i := range loop_char() {
		if loc(i) == 0 {
			log.Printf("\twarning: unit %s is nowhere\n",
				box_code(i))
		}
	}

	for _, i := range loop_loc_or_ship() {
		if loc_depth(i) > LOC_region && loc(i) == 0 {
			log.Printf("\twarning: loc %s is nowhere\n",
				box_code(i))
		}
	}
}

func check_skills() {
	var i int

	//#if 0
	//    /*
	//     *  Wed Sep  3 16:10:07 1997 -- Scott Turner
	//     *
	//     *  Really only desired for development, to check that times
	//     *  got translated correctly into the lib/skill db from the
	//     *  table in use.c.  Newer skills don't have values in the
	//     *  use.c table, so this complains about all of them.
	//     *
	//     */
	//check_skill_times();
	//#endif

	for _, sk := range loop_skill() {
		if sk >= 9000 && skill_school(sk) == sk {
			log.Printf("\twarning: orphaned subskill %s\n",
				box_code(sk))
		}
		bx[sk].temp = 0
	}

	for _, sk := range loop_skill() {
		p := rp_skill(sk)

		if learn_time(sk) == 0 {
			log.Printf("\twarning: learn time of %s is 0\n",
				box_name(sk))
		}

		if p == nil {
			continue
		}

		for i = 0; i < len(p.offered); i++ {
			if bx[p.offered[i]].temp != 0 {
				log.Printf("\twarning: both %s and %s offer skill %d\n",
					box_name(sk),
					box_name(bx[p.offered[i]].temp),
					p.offered[i])
			} else {
				bx[p.offered[i]].temp = sk
			}

			if skill_school(p.offered[i]) != sk {
				log.Printf("\twarning: %s offers %d, but %d is in school %d\n",
					box_name(sk), p.offered[i],
					p.offered[i],
					skill_school(p.offered[i]))
			}
		}

		for i = 0; i < len(p.research); i++ {
			if bx[p.research[i]].temp != 0 {
				log.Printf("\twarning: both %s and %s offer skill %d\n",
					box_name(sk),
					box_name(bx[p.research[i]].temp),
					p.research[i])
			} else {
				bx[p.research[i]].temp = sk
			}

			if skill_school(p.research[i]) != sk {
				log.Printf("\twarning: %s offers %d, but %d is in school %d\n",
					box_name(sk), p.research[i],
					p.research[i],
					skill_school(p.research[i]))
			}
		}

		for i = 0; i < len(p.guild); i++ {
			if bx[p.guild[i]].temp != 0 {
				log.Printf("\twarning: both %s and %s offer skill %d\n",
					box_name(sk),
					box_name(bx[p.guild[i]].temp),
					p.guild[i])
			} else {
				bx[p.guild[i]].temp = sk
			}

			if skill_school(p.guild[i]) != sk {
				log.Printf("\twarning: %s offers %d, but %d is in school %d\n",
					box_name(sk), p.guild[i],
					p.guild[i],
					skill_school(p.guild[i]))
			}
		}

	}

	for _, sk := range loop_skill() {
		if skill_school(sk) == sk {
			continue
		}

		if bx[sk].temp == 0 {
			log.Printf("\twarning: non-offered skill %s\n",
				box_name(sk))
		}
	}
}

func check_item_counts() {
	clear_temps(T_item)

	for _, i := range loop_boxes() {
		for _, e := range loop_inventory(i) {
			if kind(e.item) != T_item {
				log.Printf("\t%s has non-item %s\n",
					box_name(i),
					box_name(e.item))
				continue
			}

			if item_unique(e.item) == FALSE {
				continue
			}

			if item_unique(e.item) != i {
				log.Printf("\tunique item %s: whohas=%s, actual=%s\n",
					box_name(e.item),
					box_name(item_unique(e.item)),
					box_name(i))
				p_item(e.item).who_has = i
			}

			if e.qty != 1 {
				log.Printf("\t%s has qty %d of unique item %s\n",
					box_name(i),
					e.qty,
					box_name(e.item))
			}

			bx[e.item].temp += e.qty
		}
	}

	for _, i := range loop_item() {
		if item_unique(i) != FALSE {
			if bx[i].temp != 1 {
				log.Printf("\tunique item %s count %d\n",
					box_name(i),
					bx[i].temp)
			}
		}
	}
}

func check_loc_name_lengths() {

	for _, i := range loop_loc() {
		if len(just_name(i)) > 25 {
			log.Printf("\twarning: %s name too long\n", box_name(i))
		}
	}
}

func check_moving() {
	var c *command
	var leader int

	for _, i := range loop_char() {
		if stack_leader(i) != i || char_moving(i) == FALSE {
			continue
		}

		c = rp_command(i)

		if c == nil || c.state != RUN {
			log.Printf("\t%s moving but no command\n",
				box_name(i))

			restore_stack_actions(i)
		}
	}

	for _, i := range loop_char() {
		leader = stack_leader(i)

		if leader == i || char_moving(i) == char_moving(leader) {
			continue
		}

		log.Printf("\t%s moving disagrees with leader\n",
			box_name(i))
		p_char(i).moving = char_moving(leader)
	}
}

func check_prisoner() {
	for _, who := range loop_char() {
		if !is_prisoner(who) {
			continue
		}

		if stack_parent(who) == 0 {
			log.Printf("\t%s prisoner but unstacked\n", box_name(who))
			p_char(who).prisoner = FALSE
		}
	}
}

func check_city() {
	for _, city := range loop_city() {
		for _, t := range loop_trade(city) {
			if t.kind == SELL &&
				item_unique(t.item) != FALSE &&
				has_item(city, t.item) == FALSE {
				log.Printf("%s trying to sell %s which it doesn't have.\n",
					box_name(city), box_name(t.item))
			}
		}
	}
}

/*
 *  Tue Sep 22 13:25:11 1998 -- Scott Turner
 *
 *  A hack -- keep peasants out of Faery.
 *
 */
func check_peasants() {
	for _, i := range loop_province() {
		if region(i) == faery_region && has_item(i, item_peasant) != FALSE {
			log.Printf("Eliminating %s peasants from %s.\n", nice_num(has_item(i, item_peasant)), box_name(i))
			sub_item(i, item_peasant, has_item(i, item_peasant))
		}
	}
}

/*
 *  Tue Jul  6 12:57:13 1999 -- Scott Turner
 *
 *  Magical artifacts should have x_item.EntityArtifact
 *
 */
func check_magical_artifacts() {
	for _, item := range loop_subkind(sub_magic_artifact) {
		if rp_item_artifact(item) == nil {
			log.Printf("Problem with %s.\n", box_name(item))
			/*
			 *  Wed Oct 13 07:30:18 1999 -- Scott Turner
			 *
			 *  Try to replace this artifact with a new one.
			 *
			 */
			if item_unique(item) != FALSE && is_real_npc(item_unique(item)) {
				create_random_artifact(item_unique(item))
				destroy_unique_item(item_unique(item), item)
			}
		}
	}
}

/*
 *  Check database integrity.  Fixes minor problem in backlinks and lists.
 *  Always notes a database correction with a message to strerr.
 */

func check_db() error {
	stage("check_db()")

	/*
	 *  Turn off tags for db; this stuff all goes out to stderr.
	 */
	tags_off()

	check_glob()
	check_here()
	check_swear()
	check_indep()
	check_gm()
	check_deserted()
	check_skill_player()
	check_eat_player()
	check_npc_player()
	check_garr_player()
	check_nowhere()
	check_skills()
	check_item_counts()
	check_loc_name_lengths()
	check_moving()
	check_prisoner()
	check_city()
	check_peasants()
	check_magical_artifacts()

	// don't leave unique items lying about in cities; destroy them.
	for _, i := range loop_city() {
		for _, e := range loop_inventory(i) {
			if item_unique(e.item) != FALSE &&
				find_trade(i, SELL, e.item) == nil &&
				find_trade(i, BUY, e.item) == nil {
				log.Printf("Deleting %s from %s.\n", box_name(e.item), box_name(i))
				destroy_unique_item(i, e.item)
			}
		}
	}

	//#if 0
	//    /*
	//     *  Temporary check for super-rich provinces.
	//     *
	//     */
	//    for _, i := range loop_province() {
	//      if (has_item(i, item_gold) > 2*has_item(i, item_peasant))
	//        log.Printf( "%d has %d gold and %d peasants.\n",
	//            i, has_item(i, item_gold), has_item(i, item_peasant));
	//    } next_province;
	//#endif

	if bx[garrison_magic] != nil {
		log.Printf("\twarning: %s should not be allocated,\n", box_name(garrison_magic))
		log.Printf("\t\treserved for garrison_magic\n")
	}
	tags_on()

	return nil
}
