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

/*
 *  Wed Jun 11 14:26:47 1997 -- Scott Turner
 *
 *  Functions for killing off a character.  Transitioning it
 *  through a dead body to a dead soul to nothing.  Properly handles
 *  npc characters.
 *
 *  That's the theory, anyway.
 *
 */
var dead_body_np = TRUE

func nearby_grave(where int) int {
	var p *entity_loc

	where = province(where)
	p = rp_loc(where)

	if p != nil && p.near_grave != 0 {
		return p.near_grave
	}

	log_output(LOG_CODE, "%s has no nearby grave", box_name(where))

	var l []int
	for _, i := range loop_subkind(sub_graveyard) {
		l = append(l, i)
	}

	assert(len(l) > 0)

	ilist_scramble(l)

	return l[rnd(0, len(l)-1)]
}

func remove_follower(who int) {
	/*
	 *  Wed Jun 11 15:32:10 1997 -- Scott Turner
	 *
	 *  When a body finally is destroyed for good, we need to free
	 *  up all the priest related "stuff".
	 *
	 */
	if rp_char(who) != nil && rp_char(who).religion.priest != 0 {
		/*
		 *  You are no longer following someone.
		 *
		 */
		rp_char(who).religion.priest = 0

		/*
		 *  Remove it from the list of followers.
		 *
		 */
		if valid_box(rp_char(who).religion.priest) &&
			is_priest(rp_char(who).religion.priest) != 0 {
			wout(rp_char(who).religion.priest,
				"An angel informs you that your follower %s~%s has passed on.",
				rp_misc(who).save_name, box_code(who))
			rp_char(rp_char(who).religion.priest).religion.followers = rem_value(rp_char(rp_char(who).religion.priest).religion.followers, who)
		}
	}
}

/*
 *  Wed Jun 11 15:29:07 1997 -- Scott Turner
 *
 *  Destroy what you get, based on what it is.
 *
 */
func convert_to_nothing(cur_pl int, who int) {
	var name string
	/*
	 *  If it has an old lord, and can receive the old NPs, then
	 *  return them.
	 *
	 */
	pl := body_old_lord(who)
	if pl == 0 {
		pl = player(who)
	}

	/*
	 *  Figure out a name.  It may not be save_name if the bugger
	 *  never became a corpse.
	 *
	 */
	if rp_misc(who) != nil && rp_misc(who).save_name != "" {
		name = rp_misc(who).save_name
	} else {
		name = box_name(who)
	}

	if valid_box(pl) && dead_body_np != FALSE {
		if options.death_nps == 1 {
			add_np(pl, 1)
			out(pl, "%s~%s has passed on.  Gained 1 NP.",
				name, box_code(who))
		} else if options.death_nps == 2 {
			num := nps_invested(who)
			add_np(pl, num)
			if num == 1 {
				out(pl, "%s~%s has passed on.  Gained %s NP.", name, box_code(who), nice_num(num))
			} else {
				out(pl, "%s~%s has passed on.  Gained %s NP%ss.", name, box_code(who), nice_num(num))
			}
		}
		out(pl, "%s~%s has passed on.", name, box_code(who))
	} else if valid_box(pl) {
		out(pl, "The soul of %s~%s has passed on.  No NPs reclaimed.",
			name, box_code(who))
	}

	/*
	 *  Perhaps he was a follower.
	 *
	 */
	remove_follower(who)

	/*
	 *  If he was a priest, free up all his followers.
	 *
	 */
	if rp_char(who) != nil && is_priest(who) != 0 {
		for _, i := range loop_followers(who) {
			if valid_box(i) && rp_char(i) != nil {
				wout(i, "An angel informs you that your priest %s has passed on.",
					name)
				rp_char(i).religion.priest = 0
			}
		}
		rp_char(who).religion.followers = nil
	}

	/*
	 *  Now we actually destroy the object.
	 *
	 */
	if item_unique(who) != 0 {
		destroy_unique_item(item_unique(who), who)
	} else {
		if kind(who) == T_char && player(who) != 0 {
			set_lord(who, 0, LOY_UNCHANGED, 0)
		}
		set_where(who, 0)
		delete_box(who)
	}
}

/*
 *  Thu May 29 12:10:18 1997 -- Scott Turner
 *
 *  Turn something into a lost soul.  Assumes your character
 *  name is in save_name.
 *
 *  Wed Apr 29 11:19:44 1998 -- Scott Turner
 *
 *  Need to fix the "unique item" status of this object when we
 *  un-bodify-it :-)
 */
func convert_to_soul(pl int, who int) {
	dest, sum := 0, 0

	assert(kind(who) == T_char ||
		(kind(who) == T_item && subkind(who) == sub_dead_body))

	/*
	 *  Remove this thing from wherever it was at.
	 *
	 */
	if kind(who) == T_char {
		set_where(who, 0)
	} else if item_unique(who) != 0 {
		sub_item(item_unique(who), who, 1)
		p_item(who).who_has = 0
	} else {
		panic("!reached")
	}

	/*
	 *  Change it into the appropriate kind of thing.
	 *
	 */
	change_box_kind(who, T_char)
	change_box_subkind(who, sub_lost_soul)

	/*
	 *  Set it's name correctly.
	 *
	 */
	if p_misc(who) != nil && p_misc(who).save_name == "" {
		p_misc(who).save_name = bx[who].name
		bx[who].name = ""
	}
	new_name := fmt.Sprintf("Lost soul of %s", p_misc(who).save_name)
	bx[who].name = new_name

	p_item(who).weight = item_weight(item_peasant)
	p_item(who).plural_name = "lost souls"

	/*
	 *  Save the old lord properly, if possible.
	 *
	 */
	if rp_misc(who).old_lord == 0 && player(who) != 0 {
		rp_misc(who).old_lord = player(who)
	}

	/*
	 *  Make sure it won't fight.
	 *
	 */
	rp_char(who).break_point = 100 /* Don't fight */

	/*
	 *  Select a location from all the provinces in Hades.
	 *
	 */
	for _, i := range loop_province() {
		if !in_hades(i) {
			continue
		}
		if province_subloc(i, sub_city) != 0 {
			continue
		}
		sum++
		if rnd(1, sum) == 1 {
			dest = i
		}
	}
	if dest == 0 {
		panic("assert(dest)")
	}

	set_where(who, dest)
	set_lord(who, indep_player, LOY_UNCHANGED, 0) /* will this work? */
}

/*
 *  Wed Jun 11 15:11:40 1997 -- Scott Turner
 *
 *  Convert to a dead body.
 *
 *  Tue May 18 06:59:34 1999 -- Scott Turner
 *
 *  This should now drop the body wherever it was.  The only place
 *  we don't want to drop a body is an ocean province.
 *
 */
func convert_to_body(pl int, who int) {
	where := province(who)
	// grave := nearby_grave(where);
	//var p *entity_item

	/*
	 *  Should be a character.
	 *
	 */
	if !(kind(who) == T_char && subkind(who) == 0) {
		panic("assert(kind(who) == T_char && !subkind(who))")
	}

	/*
	 *  Let's not drop this into the ocean.
	 *
	 */
	if subkind(where) == sub_ocean {
		where = find_nearest_land(where)
		assert(where != 0 && subkind(where) != sub_ocean)
	}

	///*
	// * If we couldn't find a nearby grave, we might have to go
	// * directly to being a lost soul.
	// *
	// */
	//if (!grave) {
	//  convert_to_soul(pl, who);
	//  return;
	//};

	/*
	 *  Save the old "lord"
	 *
	 */
	p_misc(who).old_lord = pl

	/*
	 *  Remove it from the player's list of units...
	 *
	 */
	set_lord(who, 0, LOY_UNCHANGED, 0)

	/*
	 *  Remove this from the world.
	 *
	 */
	set_where(who, 0)

	/*
	 *  Make it a dead body.
	 *
	 */
	change_box_kind(who, T_item)
	change_box_subkind(who, sub_dead_body)

	/*
	 *  Name it appropriately.
	 *
	 */
	p_misc(who).save_name = bx[who].name
	bx[who].name = "dead body"
	p_item(who).plural_name = "dead bodies"
	p_item(who).weight = item_weight(item_peasant)

	/*
	 *  Finally, stick it in the grave.
	 *
	 */
	hack_unique_item(who, where)
}

/*
 *  Wed Jun 11 14:34:06 1997 -- Scott Turner
 *
 *  This function should take who to kill, inherit, and
 *  eventual status (body, soul, or nothing).  It should be able to
 *  transition anything to any of those states.
 *
 *  Tue May 18 06:54:02 1999 -- Scott Turner
 *
 *  Simplifying death:
 *	- Upon death, body is dropped where killed.
 *	- Bodies decay for 12 months before disintegrating.
 *
 */
var verbs = []string{"", "died", "transmigrated", "permanently ascended"}

func kill_char(who int, inherit int, status int) {
	//where := subloc(who);
	pl := player(who)

	assert(kind(who) == T_char || (kind(who) == T_item && subkind(who) == sub_dead_body))

	/*
	 *  Tue May 18 06:56:44 1999 -- Scott Turner
	 *
	 *  Treat S_soul as S_nothing...
	 *
	 */
	if status == S_soul {
		status = S_nothing
	}

	/*
	 *  Don't do anything to a noble who has survival_fatal.
	 *
	 */
	if char_melt_me(who) == 0 && survive_fatal(who) {
		return
	}

	/*
	 *  If you're being melted, you should go to nothing no matter what.
	 *
	 */
	if char_melt_me(who) != 0 {
		status = S_nothing
	}

	///*
	// *  The only non-character thing that can transition is a dead
	// *  body; everything else becomes nothing.
	// *
	// */
	//if (kind(who) != T_char) {
	//  if (subkind(who) == sub_dead_body)
	//    status = S_soul;
	//  else
	//    status = S_nothing;
	//} else {
	//  /*
	//   *  Anything with a subkind has to go to nothing.
	//   *
	//   */
	//  if (subkind(who)) status = S_nothing;
	//};

	if kind(who) != T_char || subkind(who) != 0 {
		status = S_nothing
	}

	/*
	 *  Regardless of what is eventually going to happen to you,
	 *  the following things must be done:
	 *
	 */

	/*
	 *  Inform you and your stack parent what has happened, or whoever
	 *  is holding you, if you're a unique item.
	 *
	 */
	if subkind(who) == sub_garrison {
		wout(who, "Garrison disbanded.")
	} else if item_unique(who) != 0 && rp_char(item_unique(who)) != nil {
		wout(item_unique(who), "*** %s has %s ***", just_name(who), verbs[status])
	} else if p_char(who) != nil {
		p_char(who).prisoner = FALSE
		wout(who, "*** %s has %s ***", just_name(who), verbs[status])

		sp := stack_parent(who)
		if sp != 0 {
			wout(sp, "%s has %s.", box_name(who), verbs[status])
		}

		p_char(who).prisoner = TRUE /* suppress output */
	} else {
		/* You should be either a character or a unique item. */
		panic("!reached")
	}

	/*
	 *  Inform the death log.
	 *
	 */
	log_output(LOG_DEATH, "%s %s in %s.",
		box_name(who), verbs[status], char_rep_location(who))

	/*
	 *  If the dying thing has any possessions, inherit them out.
	 *
	 */
	if kind(who) == T_char {
		take_unit_items(who, inherit, TAKE_SOME)
	}

	/*
	 *  If the dying thing is the leader of a stack and is
	 *  moving, we need to restore the stacks actions!
	 *
	 */
	if kind(who) == T_char && stack_leader(who) == who && char_moving(who) != 0 {
		restore_stack_actions(who)
	}

	/*
	 *  If the dying thing is in a stack, then remove it from the
	 *  stack.
	 *
	 */
	if kind(who) == T_char {
		extract_stacked_unit(who)
	}

	/*
	 *  Fix its orders.
	 *
	 */
	if kind(who) == T_char && player(who) != 0 {
		flush_unit_orders(player(who), who)
		interrupt_order(who)
	}

	/*
	 *  Delete all its aura/piety if has any.
	 *
	 */
	if rp_magic(who) != nil {
		rp_magic(who).cur_aura = 0
	}

	/*
	 *  If it's newly dead, then it's not a prisoner.
	 *
	 */
	if rp_char(who) != nil {
		rp_char(who).prisoner = FALSE
	}

	/*
	 *  If it can cheat death, then let it.
	 *
	 */
	if char_melt_me(who) == 0 && has_skill(who, sk_transcend_death) != 0 {
		log_output(LOG_SPECIAL, "%s transcends death", box_name(who))
		log_output(LOG_SPECIAL, "...%s moved to %s",
			box_name(who), box_name(hades_pit))
		p_char(who).prisoner = FALSE
		p_char(who).sick = FALSE
		p_char(who).health = 100
		move_stack(who, hades_pit)
		wout(who, "%s appears at %s.",
			box_name(who), box_name(hades_pit))
		return
	}

	/*
	 *  It might have an artifact.
	 *
	 */
	if char_melt_me(who) == 0 {
		if a := has_artifact(who, ART_RESTORE, 0, 0, 1); a != 0 {
			where := pick_starting_city(nation(who), 0)
			log_output(LOG_SPECIAL, "%s uses a Restore Life artifact %s.",
				box_name(who), box_name(a))
			log_output(LOG_SPECIAL, "Transferring %s to %s.", box_name(who),
				box_name(where))
			p_char(who).prisoner = FALSE
			p_char(who).sick = FALSE
			p_char(who).health = 100
			move_stack(who, where)
			wout(who, "%s appears at %s.",
				box_name(who), box_name(where))
			/*
			 *  Now fix up the artifact.
			 *
			 */
			rp_item_artifact(a).Uses--
			if rp_item_artifact(a).Uses == 0 {
				wout(who, "%s vanishes.", box_name(a))
				destroy_unique_item(who, a)
			}
			return
		}
	}

	/*
	 *  It might be a "token item", meaning we have to take care
	 *  of resetting the token, etc.
	 */
	{
		token_item := our_token(who)
		if token_item != 0 {
			who_has := item_unique(token_item)
			token_pl := p_player(token_item)
			token_pl.Units = rem_value(token_pl.Units, who)

			if char_melt_me(who) == 0 {
				p_item_magic(token_item).TokenNum--
			}

			if item_token_num(token_item) <= 0 {
				if player(who_has) == sub_pl_regular {
					wout(who_has, "%s vanishes.",
						box_name(token_item))
				}
				destroy_unique_item(who_has, token_item)
			}
			assert(status == S_nothing) /* Should be disappearing? */
		}
	}

	/*
	 *  If the unit is going to change players, it needs
	 *  to update it's loyalties.
	 *
	 */
	//if (kind(who) == T_char) {
	//	unit_deserts(who, 0, TRUE, LOY_UNCHANGED, 0);
	//}

	/*
	 *  Put back the cookie...does nothing if inappropriate.
	 *
	 */
	put_back_cookie(who)

	switch status {
	case S_body:
		p_char(who).death_time = sysclock /* record time of death */
		convert_to_body(pl, who)
	case S_soul:
		//convert_to_soul(pl, who);
		panic("!reached")
	case S_nothing:
		convert_to_nothing(pl, who)
	default:
		panic("!reached")
	}

}

func restore_dead_body(owner, who int) {
	var pm *entity_misc
	var pi *entity_item
	var pc *entity_char

	log_output(LOG_CODE, "dead body revived: who=%s, owner=%s, player=%s",
		box_code_less(who),
		box_code_less(owner),
		box_code_less(player(owner)))

	/*
	 *  If it's an item (body), remove it...
	 *
	 *  Thu May 17 07:33:23 2001 -- Scott Turner
	 *
	 *  Owner is not necessarily the holder of the body, since
	 *  you can now resurrect at a distance.
	 *
	 */
	if kind(who) == T_item {
		ret := sub_item(item_unique(who), who, 1)
		assert(ret)
		p_item(who).who_has = 0
	}

	change_box_kind(who, T_char)
	change_box_subkind(who, 0)

	pm = p_misc(who)
	pi = p_item(who)
	pc = p_char(who)

	pi.weight = 0
	my_free(pi.plural_name)
	pi.plural_name = ""

	if pm.save_name != "" {
		set_name(who, pm.save_name)
		pm.save_name = ""
	}

	pc.health = 100
	pc.sick = FALSE
	pc.break_point = 50

	set_where(who, owner)

	if kind(pm.old_lord) == T_player {
		wout(pm.old_lord, "%s has been brought back to life.",
			box_name(who))

		//if (pm.old_lord != player(owner)) {
		//	p_char(who).prisoner = TRUE;
		//}

		set_lord(who, pm.old_lord, LOY_UNCHANGED, 0)
	} else {
		set_lord(who, indep_player, LOY_UNCHANGED, 0)
	}

	pm.old_lord = 0

	wout(owner, "%s stacks beneath us.", box_name(who))
}
