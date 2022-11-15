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

func v_look(c *command) int {

	if kind(c.who) != T_char {
		wout(c.who, "%s is not a character.", box_name(c.who))
		return FALSE
	}

	show_loc(c.who, subloc(c.who))

	return TRUE
}

func v_explore(c *command) int {
	return TRUE
}

func find_lost_items(who, where int) int {
	var e *item_ent
	item := 0
	var chance int

	for _, e = range loop_inventory(where) {
		if FALSE == item_unique(e.item) {
			continue
		}

		/*
		 *  Don't take dead bodies out of graveyards; that's what EXHUME is for
		 */
		if subkind(where) == sub_graveyard &&
			subkind(e.item) == sub_dead_body {
			continue
		}

		/*
		 *  Don't take magic rings from the market
		 *
		 *  Wed Apr 22 06:20:30 1998 -- Scott Turner
		 *
		 *  Don't find anything that is being sold.
		 */
		if find_trade(where, SELL, e.item) != nil {
			continue
		}

		item = e.item
		break
	}

	if loc_depth(where) >= LOC_subloc {
		chance = 100
	} else {
		chance = 40
	}

	if item == 0 || rnd(1, 100) > chance {
		return FALSE
	}

	move_item(where, who, item, 1)
	wout(who, "%s found one %s.", box_name(who), box_name(item))

	log_output(LOG_MISC, "%s found %s in %s.",
		box_name(who), box_name(item),
		char_rep_location(where))

	return TRUE
}

const HERO = false
const NEW_TRADE = true

/*
 *  Probability breakdown for explore
 *
 *	50%	fail
 *	33%	success
 *	17%	fail, but message indicating if there is something to find
 *
 *  Tue Oct 29 11:22:06 1996 -- Scott Turner
 *
 *  Priests of Domingo get an advantage:
 *
 *      30%     fail
 *      50%     success
 *	20%	fail, but message indicating if there is something to find
 */

func d_explore(c *command) int {
	var hidden_exits int
	var l []*exit_view
	where := subloc(c.who)
	var i, chance int
	var r int

	if find_lost_items(c.who, where) != FALSE {
		return TRUE
	}

	/*
	 *  Explore in a ship should explore the surrounding ocean region
	 */

	if is_ship(where) && subkind(loc(where)) == sub_ocean {
		where = loc(where)
		find_lost_items(c.who, where)
	}

	r = rnd(1, 100)

	if is_priest(c.who) == sk_domingo {
		chance = 30
	} else {
		chance = 50
	}

	if HERO { // #ifdef HERO
		//the Improved Explore skill.
		chance += min(50, skill_exp(c.who, sk_improved_explore))
	} // #endif // HERO

	if r <= chance {
		wout(c.who, "Exploration of %s uncovers no new features.",
			box_code(where))
		return FALSE
	}

	l = exits_from_loc(c.who, where)

	hidden_exits = count_hidden_exits(l)

	/*
	 *  Nothing to find
	 */

	if hidden_exits <= 0 {
		wout(c.who, "Exploration of %s uncovers no new features.",
			box_code(where))
		return FALSE
	}

	/*
	 *  Something to find, but a bad roll
	 */

	if is_priest(c.who) == sk_domingo {
		chance = 50
	} else {
		chance = 67
	}

	if r <= chance {
		switch rnd(1, 4) {
		case 1:
			wout(c.who, "Rumors speak of hidden features here, but none were found.")
			break

		case 2:
			wout(c.who, "We suspect something is hidden here, but did not find anything.")
			break

		case 3:
			wout(c.who, "Something may be hidden here.  Further exploration is needed.")
			break

		case 4:
			wout(c.who, "Nothing was found, but further exploration looks promising.")
			break

		default:
			panic("!reached")
		}
		return FALSE
	}

	/*
	 *  Choose what we found randomly
	 */

	i = rnd(1, hidden_exits)

	find_hidden_exit(c.who, l, hidden_count_to_index(i, l))

	return TRUE
}

func may_name(who, target int) bool {

	switch kind(target) {
	case T_char, T_player:
		return player(who) == player(target)

	case T_loc, T_ship:
		if loc_depth(target) == LOC_build {
			return player(who) == player(building_owner(target))
		}

		return may_rule_here(who, target) != FALSE

	case T_item:
		if has_auraculum(who) == target {
			return true
		}

		if item_creator(target) == who {
			return true
		}

		if is_artifact(target) != nil {
			return true
		}

		switch item_use_key(target) {
		case use_death_potion:
		case use_heal_potion:
		case use_slave_potion:
		case use_proj_cast:
		case use_quick_cast:
			return true
		}

		return false

	case T_storm:
		if npc_summoner(target) == who {
			return true
		}
	}

	return false
}

func v_name(c *command) int {
	target := c.who
	var new_name string
	var old_name string
	var l int

	if numargs(c) >= 2 && c.a > 0 {
		target = c.a
		cmd_shift(c)
	}

	new_name = rest_name(c, 1)

	if len(new_name) == 0 {
		wout(c.who, "No new name given.")
		return FALSE
	}

	if !may_name(c.who, target) {
		wout(c.who, "Not allowed to change the name of %s.",
			box_code(target))
		return FALSE
	}

	switch kind(target) {
	case T_char:
		l = 35
		break

	default:
		l = 25
	}

	if len(new_name) > l {
		wout(c.who, "Name is longer than %d characters.", l)
		return FALSE
	}

	old_name = (box_name(target))
	set_name(target, new_name)

	wout(c.who, "%s will now be known as %s.",
		old_name, box_name(target))

	if target != c.who &&
		(kind(target) == T_char || is_loc_or_ship(target)) {
		wout(target, "%s will now be known as %s.",
			old_name, box_name(target))
	}

	return TRUE
}

func v_times(c *command) int {
	var p *entity_player

	if c.a == 1 {
		c.a = 2
	} else if c.a == 2 {
		c.a = 1
	}

	p = p_player(player(c.who))
	p.compuserve = c.a != FALSE

	if p.compuserve {
		wout(c.who, "Will not receive the paper.")
	} else {
		wout(c.who, "Will receive the paper.")
	}

	return TRUE
}

func v_fullname(c *command) int {
	var new_name string
	var p *entity_player

	new_name = rest_name(c, 1)

	if len(new_name) == 0 {
		wout(c.who, "No new name given.")
		return FALSE
	}

	if len(new_name) > 60 {
		wout(c.who, "Name is longer than %d characters.", 60)
		return FALSE
	}

	p = p_player(player(c.who))
	if len(p.full_name) != 0 {
		my_free(p.full_name)
	}
	p.full_name = new_name

	return TRUE
}

func v_banner(c *command) int {
	target := c.who
	var new_name string

	if numargs(c) >= 2 && c.a > 0 {
		target = c.a
		cmd_shift(c)

		if !valid_box(target) {
			wout(c.who, "%s is not a valid entity.",
				box_code(target))
			return FALSE
		}

		if !may_name(c.who, target) {
			wout(c.who, "You do not control %s.",
				box_code(target))
			return FALSE
		}
	}

	new_name = rest_name(c, 1)

	//#if 0
	//    if (kind(target) != T_char) {
	//        wout(c.who, "Cannot set the banner of %s.", box_name(target));
	//        return FALSE;
	//    }
	//#endif

	if len(new_name) > 50 {
		wout(c.who, "Banner is longer than 50 characters.")
		return FALSE
	}

	set_banner(target, new_name)

	if len(new_name) != 0 {
		out(c.who, "Banner set.")
	} else {
		out(c.who, "Banner cleared.")
	}

	return TRUE
}

func how_many(who, from_who, item, qty, have_left int) int {
	var num_has int

	num_has = has_item(from_who, item)

	if num_has <= 0 {
		wout(who, "%s has no %s.",
			just_name(from_who),
			just_name(item))

		return 0
	}

	if num_has <= have_left {
		wout(who, "%s has only %s.",
			just_name(from_who),
			just_name_qty(item, num_has))
		return 0
	}

	if qty == 0 {
		qty = num_has
	}

	qty = min(num_has-have_left, qty)

	assert(qty > 0)

	return qty
}

func v_accept(c *command) int {
	from_who := c.a
	item := c.b
	qty := c.c
	var p *entity_char
	var i int

	p = p_char(c.who)

	if item != FALSE && ((kind(item) == T_player || kind(item) == T_char) &&
		kind(from_who) != T_char) { /* fixed VLN */
		/* problem was:
		    wanted to accept a prisoner [an5o] from noble f1k
		   so I wrote: accept f1k an5o 1

		   but got this notice on the report:

		   1: > ACCEPT f1k an5o 1
		   1: Accept expects "accept <from-who> <item> <qty>"
		   1: It looks like you've swapped the first two. I'll swap them back for you.

		   So I was instead accepting the noble from the prisoner!!!

		   VLN added second line to check.
		*/
		var tmp int
		wout(c.who, "Accept expects \"accept <from-who> <item> <qty>\"")
		wout(c.who, "It looks like you've swapped the first two.  I'll swap them back for you.")
		tmp = item
		item = from_who
		from_who = tmp
	}

	/*
	 *  Wed Jan  6 06:19:54 1999 -- Scott Turner
	 *
	 *  Special cases... accept clear && nation names.
	 *
	 */
	if from_who == 0 && strncasecmp([]byte("clear"), c.parse[1], 5) == 0 {
		var i int

		for i = 0; i < len(p.accept); i++ {
			my_free(p.accept[i])
		}

		p.accept = nil

		wout(c.who, "Accept list cleared.")
		return TRUE
	}

	/*
	 *  Is it a nation name?
	 *
	 */
	if from_who == 0 {
		n := find_nation(string(c.parse[1]))

		if n != FALSE {
			wout(c.who, "Accepting from nation '%s'.", rp_nation(n).name)
			from_who = n
		}
	}

	/*
	 *  Maybe you're already doing this?
	 *
	 */
	for i = 0; i < len(p.accept); i++ {
		if p.accept[i].item == item &&
			p.accept[i].from_who == from_who &&
			p.accept[i].qty == qty {
			wout(c.who, "You already have that accept order active.")
			return FALSE
		}
	}

	new_ae := &accept_ent{}

	new_ae.item = item
	new_ae.from_who = from_who
	new_ae.qty = qty

	p.accept = append(p.accept, new_ae)

	return TRUE
}

func will_accept_sup(who, item, from, qty int) bool {
	var p *entity_char
	var i int

	p = rp_char(who)

	if p != nil {
		for i = 0; i < len(p.accept); i++ {
			item_match := (p.accept[i].item == item ||
				p.accept[i].item == 0)
			from_match := (p.accept[i].from_who == from ||
				nation(from) == p.accept[i].from_who ||
				p.accept[i].from_who == 0)
			qty_match := (p.accept[i].qty >= qty ||
				p.accept[i].qty == 0)

			if item_match && from_match && qty_match {
				if p.accept[i].qty != FALSE {
					p.accept[i].qty -= qty
				}

				return true
			}
		}
	}

	return false
}

func will_accept(who, item, from, qty int) bool {

	if item == item_gold {
		return true
	}

	if player(who) == player(from) {
		return true
	}

	if subkind(who) == sub_garrison {
		if may_rule_here(from, who) != FALSE {
			return true
		}

		wout(from, "%s is not under your control.",
			box_name(who))
		return false
	}

	if will_accept_sup(who, item, from, qty) ||
		will_accept_sup(player(who), item, from, qty) ||
		will_accept_sup(who, item, player(from), qty) ||
		will_accept_sup(player(who), item, player(from), qty) {
		return true
	}

	wout(who, "Refused %s from %s.", just_name_qty(item, qty),
		box_name(from))
	wout(from, "Refused by %s.", just_name(who))

	return true
}

/*
 *  give <who> <what> [qty] [have-left]
 *
 *  Fri Mar 30 18:15:28 2001 -- Scott Turner
 *
 *  Hacking this so that giving to your faction id puts stuff in your claim.
 *
 */

func v_give(c *command) int {
	target := c.a
	item := c.b
	qty := c.c
	have_left := c.d
	var ret int
	var fee int

	/*
	 *  Try to correct the arguments if they got the order wrong
	 */

	if numargs(c) >= 2 &&
		(kind(target) == T_item || has_prisoner(c.who, target) != FALSE) &&
		(kind(item) == T_char && subloc(c.who) == subloc(item) &&
			FALSE == has_prisoner(c.who, item)) {
		var tmp int

		tmp = c.a
		c.a = c.b
		c.b = tmp

		target = c.a
		item = c.b

		switch numargs(c) {
		case 2:
			wout(c.who, "(assuming you meant 'give %d %d')",
				target, item)
			break

		case 3:
			wout(c.who, "(assuming you meant 'give %d %d %d')",
				target, item, qty)
			break

		default:
			wout(c.who, "(assuming you meant 'give %d %d %d %d')",
				target, item, qty, have_left)
		}
	}

	/*
	 *  Fri Mar 30 18:17:32 2001 -- Scott Turner
	 *
	 *  Permit giving to your faction if allowed.
	 *
	 */
	if options.claim_give != FALSE &&
		target == player(c.who) {
		if subkind(subloc(c.who)) != sub_city {
			wout(c.who, "You can only deposit into CLAIM in a city.")
			return FALSE
		}
		if item != item_gold {
			wout(c.who, "You cannot deposit that item into CLAIM.")
			return FALSE
		}
	} else if FALSE == check_char_here(c.who, target) {
		return FALSE
	} else if FALSE == check_char_gone(c.who, target) {
		return FALSE
	}

	if loyal_kind(target) == LOY_summon {
		wout(c.who, "Summoned entities may not be given anything.")
		return FALSE
	}

	if is_prisoner(target) {
		wout(c.who, "Prisoners cannot accept anything.")
		return FALSE
	}

	if has_prisoner(c.who, item) != FALSE {
		if !will_accept(target, item, c.who, 1) {
			return FALSE
		}

		return give_prisoner(c.who, target, item)
	}

	if kind(item) != T_item {
		wout(c.who, "%s is not an item or a prisoner.",
			box_code(item))
		return FALSE
	}

	if has_item(c.who, item) < 1 {
		wout(c.who, "%s does not have any %s.", box_name(c.who),
			box_code(item))
		return FALSE
	}

	if rp_item(item).ungiveable != FALSE {
		wout(c.who, "You cannot transfer %s to another noble.",
			plural_item_name(item, 2))
		return FALSE
	}

	qty = how_many(c.who, c.who, item, qty, have_left)

	if qty <= 0 {
		return FALSE
	}

	if !will_accept(target, item, c.who, qty) {
		return FALSE
	}

	if options.claim_give != FALSE && target == player(c.who) {
		fee = qty / 10
		if fee < 1 {
			fee = 1
		}
		sub_item(c.who, item, fee)
		qty -= fee
	}

	ret = move_item(c.who, target, item, qty)
	assert(ret != FALSE)

	if target == player(c.who) {
		wout(c.who, "Deposited %s into CLAIM.", just_name_qty(item, qty))
		wout(c.who, "Paid fee of %s.", just_name_qty(item, fee))
	} else {
		wout(c.who, "Gave %s to %s.", just_name_qty(item, qty),
			box_name(target))
		wout(target, "Received %s from %s.", box_name_qty(item, qty),
			box_name(c.who))
	}

	return TRUE
}

func v_pay(c *command) int {
	target := c.a
	qty := c.b
	have_left := c.c

	ret := oly_parse(c, []byte(sout("give %d 1 %d %d", target, qty, have_left)))
	assert(ret)

	return v_give(c)
}

func may_take(who, target int) bool {

	if FALSE == check_char_here(who, target) {
		return false
	}
	if FALSE == check_char_gone(who, target) {
		return false
	}

	if subkind(target) == sub_garrison {
		if may_rule_here(who, target) != FALSE {
			return true
		}

		wout(who, "%s is not under your control.",
			box_name(target))
		return false
	}

	if FALSE == my_prisoner(who, target) &&
		player(target) != player(who) {
		wout(who, "May only take items from other units in your faction.")
		return false
	}

	return true
}

/*
 *  get <who> <what> <qty> <have-left>
 */

func v_get(c *command) int {
	target := c.a
	item := c.b
	qty := c.c
	have_left := c.d
	var ret int

	if !may_take(c.who, target) {
		return FALSE
	}

	if has_prisoner(target, item) != FALSE {
		return give_prisoner(target, c.who, item)
	}

	if kind(item) != T_item {
		wout(c.who, "%s is not an item or a prisoner.",
			box_code(item))
		return FALSE
	}

	if has_item(target, item) < 1 {
		wout(c.who, "%s does not have any %s.", box_name(target),
			box_code(item))
		return FALSE
	}

	qty = how_many(c.who, target, item, qty, have_left)

	if qty <= 0 {
		return FALSE
	}

	if subkind(target) == sub_garrison && man_item(item) != FALSE {
		garr_men := count_man_items(target)

		garr_men -= qty

		if garr_men < 10 {
			wout(c.who, "Garrisons must be left with a minimum of ten men.")
			return FALSE
		}
	}

	if rp_item(item).ungiveable != FALSE {
		wout(c.who, "You cannot transfer %s between nobles.",
			plural_item_name(item, 2))
		return FALSE
	}

	ret = move_item(target, c.who, item, qty)
	assert(ret != FALSE)

	wout(c.who, "Took %s from %s.", just_name_qty(item, qty),
		box_name(target))

	wout(target, "%s took %s from us.", box_name(c.who),
		box_name_qty(item, qty))

	if item == item_sailor || item == item_pirate {
		check_captain_loses_sailors(qty, target, c.who)
	}

	return TRUE
}

// todo: really?
func noble_cost(pl int) int { return 1 }

func next_np_turn(pl int) int {
	var p *entity_player
	var ft, ct int
	var n int

	p = p_player(pl)

	ct = (7 - (sysclock.turn+1)%NUM_MONTHS)
	/* ct = (7 - (sysclock.turn + 1)) % NUM_MONTHS; */
	ft = p.first_turn % NUM_MONTHS
	n = (ft + ct) % NUM_MONTHS

	return n
}

func print_hiring_status(pl int) {
	var n int

	assert(kind(pl) == T_player)

	if subkind(pl) != sub_pl_regular {
		return
	}

	n = next_np_turn(pl)

	if n == 0 {
		n += NUM_MONTHS
	}

	wout(pl, "The next NP will be received at the end of turn %d.",
		sysclock.turn+n)
}

func print_unformed(pl int) {
	p := rp_player(pl)
	var n int
	var buf string
	var i int

	if p == nil || len(p.unformed) < 1 {
		return
	}

	n = len(p.unformed)

	for i = 0; i < n && i < 5; i++ {
		buf += (sout(" %s", box_code_less(p.unformed[i])))
	}

	out(pl, "")
	wout(pl, "The next %s nobles formed will be: %s", nice_num(n), buf)
}

//#if 0
///*
// *  Some micromodeling nonsense to randomly equip a new noble
// *  with a few items or skills
// */
//
//static void
//equip_new_noble(int who, int new)
//{
//    where := subloc(who);
//    var n int
//    int qty;
//
//    if (rnd(1,4) == 1)	/* appropriate region skill */
//    {
//        if (is_port_city(where) && rnd(1,5) < 5)
//        {
//            n = sk_shipcraft;
//        }
//        else if (has_ocean_access(where) && rnd(1,2) == 1)
//        {
//            n = sk_shipcraft;
//        }
//        else
//        {
//            switch (rnd(1,4))
//            {
//            case 1:
//            case 2:
//                n = sk_combat;
//                break;
//
//            case 3:
//                n = sk_construction;
//                break;
//
//            case 4:
//                n = sk_stealth;
//                break;
//
//            default:
//                panic("!reached")
//            }
//        }
//
//        set_skill(new, n, SKILL_know);
//        wout(who, "%s knows %s.", just_name(new), box_name(n));
//    }
//
//    if (rnd(1,5) == 1)	/* a possession */
//    {
//        qty = 1;
//
//        switch (rnd(1,5))
//        {
//        case 1:		/* gold */
//            qty = rnd(50, 550);
//            n = item_gold;
//            break;
//
//        case 2:
//            n = item_riding_horse;
//            break;
//
//        case 3:
//            n = item_longsword;
//            break;
//
//        case 4:
//            n = item_longbow;
//            break;
//
//        case 5:
//            n = item_warmount;
//            break;
//
//        default:
//            panic("!reached")
//        }
//
//        gen_item(new, n, qty);
//        wout(who, "%s has %s.", just_name(new), just_name_qty(n, qty));
//    }
//}
//#endif

func form_new_noble(who int, name string, new_noble int) {
	var p *entity_char
	var op *entity_char

	assert(kind(new_noble) == T_unform)

	change_box_kind(new_noble, T_char)

	p = p_char(new_noble)
	op = p_char(who)

	p.behind = op.behind
	p.fresh_hire = TRUE
	p.health = 100

	p_char(new_noble).attack = 80
	p_char(new_noble).defense = 80
	p_char(new_noble).break_point = 50

	set_name(new_noble, name)

	set_where(new_noble, subloc(who))
	set_lord(new_noble, player(who), LOY_contract, 500)

	join_stack(new_noble, who)

	//#if 0
	//    equip_new_noble(who, new);
	//#endif
}

func v_form(c *command) int {
	var pl int
	var cost int

	if subkind(subloc(c.who)) != sub_city {
		wout(c.who, "Nobles may only be formed in cities.")
		return FALSE
	}

	pl = player(c.who)
	cost = noble_cost(pl)

	if player_np(pl) < cost {
		wout(c.who, "To form another noble requires %d Noble Point%s.", cost, add_s(cost))
		return FALSE
	}

	return TRUE
}

func d_form(c *command) int {
	var new_name string
	var pl int
	var cost int
	new_noble := c.a
	var p *entity_player

	pl = player(c.who)
	cost = noble_cost(pl)

	if player_np(pl) < cost {
		wout(c.who, "To form another noble requires %d Noble Point%s.",
			cost, add_s(cost))
		return FALSE
	}

	p = p_player(player(c.who))

	if new_noble != FALSE {
		if kind(new_noble) == T_char && player(new_noble) == pl {
			wout(c.who, "You've already created that noble.")
			return FALSE
		}
		if kind(new_noble) != T_unform ||
			ilist_lookup(p.unformed, new_noble) < 0 {
			wout(c.who, "%s is not a valid unformed noble entity.", box_code(new_noble))
			wout(c.who, "I will use one of your unformed noble codes at random.")
			new_noble = 0
		}
	}

	if new_noble == 0 && len(p.unformed) > 0 {
		new_noble = p.unformed[0]
	}

	//#if 1
	if new_noble == 0 {
		wout(c.who, "No further nobles may be formed this turn.")
		return FALSE
	}
	//#else
	//    if (new == 0){
	//		new = new_ent(T_unform);
	//	}
	//    if (new < 0) {
	//        wout(c.who, "No nobles were interested in joining.");
	//        return FALSE;
	//    }
	//#endif

	if numargs(c) < 2 || c.parse[2] == nil || len(c.parse[2]) == 0 {
		new_name = "New noble"
	} else {
		new_name = string(c.parse[2])
	}

	form_new_noble(c.who, new_name, new_noble)

	p.unformed = rem_value(p.unformed, new_noble)
	deduct_np(pl, cost)

	return TRUE
}

type flag_ent struct {
	who  int
	flag string
}

var flags []*flag_ent

func flag_raised(who int, flag string) int {
	var i int

	for i = 0; i < len(flags); i++ {
		if who != 0 &&
			player(flags[i].who) != who &&
			flags[i].who != who {
			continue
		}

		if i_strcmp([]byte(flags[i].flag), []byte(flag)) == 0 {
			return i
		}
	}

	return -1
}

func v_flag(c *command) int {
	var new_flag *flag_ent
	var flag string

	if numargs(c) < 1 {
		wout(c.who, "Must specify what message to signal.")
		return FALSE
	}

	flag = string(c.parse[1])

	if flag_raised(c.who, flag) >= 0 {
		wout(c.who, "%s has already given that signal this month.",
			box_name(c.who))
		return FALSE
	}

	new_flag = &flag_ent{}
	new_flag.who = c.who
	new_flag.flag = flag

	flags = append(flags, new_flag)

	return TRUE
}

var wait_tags = []string{
	"time",  /* 0 */
	"day",   /* 1 */
	"unit",  /* 2 */
	"gold",  /* 3 */
	"item",  /* 4 */
	"flag",  /* 5 */
	"loc",   /* 6 */
	"stack", /* 7 */
	"top",   /* 8 */
	"ferry", /* 9 */
	"ship",  /* 10 */
	"rain",  /* 11 */
	"fog",   /* 12 */
	"wind",  /* 13 */
	"not",   /* 14 */
	"owner", /* 15 */

	"raining", /* 16 . 11 */
	"foggy",   /* 17 . 12 */
	"windy",   /* 18 . 13 */

	"clear",   /* 19 */
	"shiploc", /* 20 */
	"month",   /* 21 */
	"turn",    /* 22 */

	"mist",    /* 23 */
	"misty",   /* 24 */
	"alone",   /* 25 */
	"teacher", /* 26 */
	"teach",   /* 27 */

	""}

func clear_wait_parse(c *command) {
	var i int

	for i = 0; i < len(c.wait_parse); i++ {
		my_free(c.wait_parse[i])
		c.wait_parse[i] = nil
	}

	c.wait_parse = nil
}

func parse_wait_args(c *command) string {
	var tag int
	var i int
	var tag_s string
	var new_wa *wait_arg

	assert(len(c.wait_parse) == 0)

	i = 1
	for i <= numargs(c) {
		tag_s = string(c.parse[i])
		tag = lookup(wait_tags, tag_s)

		switch tag {
		case 16:
			tag = 11
			break
		case 17:
			tag = 12
			break
		case 18:
			tag = 13
			break
		case 24:
			tag = 23
			break
		case 27:
			tag = 26
			break
		}

		if tag >= 0 {
			i++
		} else {
			return sout("Unknown condition '%s'.", tag_s)
		}

		new_wa = &wait_arg{}
		c.wait_parse = append(c.wait_parse, new_wa)
		new_wa.tag = tag
		new_wa.a1 = 0
		new_wa.a2 = 0
		new_wa.flag = ""

		switch tag {
		case 0, /* time n */
			1,  /* day n */
			21, /* month */
			22, /* turn */
			2,  /* unit n */
			3,  /* gold n */
			6,  /* loc n */
			20, /* shiploc */
			7,  /* stack n */
			9,  /* ferry n */
			10, /* ship n */
			26: /* teacher n */
			if i <= numargs(c) {
				new_wa.a1 = parse_arg(c.who, c.parse[i])
				i++
			} else {
				return sout("Argument missing for '%s'.", tag_s)
			}
			break

		case 4: /* item n q */
			if i <= numargs(c) {
				new_wa.a1 = parse_arg(c.who, c.parse[i])
				i++
			} else {
				return sout("Argument missing for '%s'.", tag_s)
			}

			if i <= numargs(c) {
				new_wa.a2 = parse_arg(c.who, c.parse[i])
				i++
			} else {
				new_wa.a2 = 1
			} /* missing arg, really */
			break

		case 5: /* flag f [n] */
			if i <= numargs(c) {
				new_wa.flag = string(c.parse[i])
				i++
			} else {
				return sout("Flag missing.")
			}

			new_wa.a1 = player(c.who) /* special default */

			if i <= numargs(c) && (isdigit(c.parse[i][0]) || parse_arg(c.who, c.parse[i]) != FALSE) {
				new_wa.a1 = parse_arg(c.who, c.parse[i])
				i++
			}
			break

		case 8, /* top */
			11, /* rain */
			12, /* fog */
			13, /* wind */
			14, /* not */
			15, /* owner */
			19, /* clear */
			23, /* mist */
			25: /* alone */
			break

		default:
			panic(fmt.Sprintf("assert(tag != %d)", tag))
		}
	}

	return ""
}

func check_wait_conditions(c *command) string {
	var i int
	var p *wait_arg
	var ret string
	not, setnot, cond := false, false, false
	var where_ship int

	where_ship = subloc(c.who)
	if is_ship_either(where_ship) {
		where_ship = subloc(where_ship)
	}

	if len(c.wait_parse) < 1 {
		if ret = parse_wait_args(c); ret != "" {
			return ret
		}

		assert(len(c.wait_parse) > 0)
	}

	for i = 0; i < len(c.wait_parse); i++ {
		p = c.wait_parse[i]

		if setnot {
			setnot = false
		} else if not {
			not = false
		}

		switch p.tag {
		case 0: /* time n */
			cond = (command_days(c) >= p.a1)
			if not {
				cond = !cond
			}
			if cond {
				return sout("%s day%s%s passed.",
					nice_num(p.a1),
					or_string(p.a1 == 1, " has", "s have"),
					or_string(not, " not", ""))
			}
			break

		case 1: /* day n */
			cond = (sysclock.day >= p.a1)
			if not {
				cond = !cond
			}
			if cond {
				if not {
					return sout("today is%s day %d.", or_string(not, " not", ""), p.a1)
				} else {
					return sout("today is day %d.", sysclock.day)
				}
			}
			break

		case 21, /* month n */
			22: /* turn n */
			cond = (sysclock.turn >= p.a1)
			if not {
				cond = !cond
			}
			if cond {
				if not {
					return sout("it is%s turn %d.", or_string(not, " not", ""), p.a1)
				} else {
					return sout("it is turn %d.", sysclock.turn)
				}
			}
			break

		case 2: /* unit n */
			if !valid_box(p.a1) {
				return sout("%s does not exist.", box_code(p.a1))
			}

			//#if 0
			//                cond = (subloc(c.who) == subloc(p.a1));
			//#else
			cond = char_here(c.who, p.a1)
			//#endif
			if not {
				cond = !cond
			}
			if cond {
				return sout("%s is%s here.", box_code(p.a1), or_string(not, " not", ""))
			}
			break

		case 3: /* gold n */
			cond = (has_item(c.who, item_gold) >= p.a1)
			if not {
				cond = !cond
			}
			if cond {
				if not {
					return sout("%s doesn't have %s.", just_name(c.who), gold_s(p.a1))
				}
				return sout("%s has %s.", just_name(c.who), gold_s(has_item(c.who, item_gold)))
			}
			break

		case 4: /* item n q */
			//#if 0
			//                /* fails if we wait for a corpse */
			//
			//                if (kind(p.a1) != T_item)
			//                  return sout("%s is not an item.", box_code(p.a1));
			//#endif

			cond = (kind(p.a1) == T_item && has_item(c.who, p.a1) >= p.a2)
			if not {
				cond = !cond
			}
			if cond {
				if not {
					return sout("%s doesn't have %s.", just_name(c.who), just_name_qty(p.a1, p.a2))
				}
				return sout("%s has %s.", just_name(c.who), just_name_qty(p.a1, has_item(c.who, p.a1)))
			}
			break

		case 5: /* flag */
			if p.a1 != FALSE && !valid_box(p.a1) {
				return sout("%s does not exist.", box_code(p.a1))
			}
			j := flag_raised(p.a1, p.flag)
			if not {
				if j < 0 {
					return "received no signal"
				}
			} else {
				if j >= 0 {
					return sout("%s signaled '%s'", box_name(flags[j].who), flags[j].flag)
				}
			}
			break

		case 6: /* loc */
			if !is_loc_or_ship(p.a1) && c.f == FALSE {
				wout(c.who, "Warning: %s is not a location or ship.", box_code(p.a1))
				c.f = 1
			}
			cond = (subloc(c.who) == p.a1)
			if not {
				cond = !cond
			}
			if cond {
				return sout("%sat %s.", or_string(not, "not ", ""), box_name(p.a1))
			}
			break

		case 20: /* shiploc */
			ship := subloc(c.who)
			if !is_ship(ship) && !is_ship_notdone(ship) {
				return sout("%s is not on a ship.", box_name(c.who))
			} else if !is_loc_or_ship(p.a1) {
				return sout("%s is not a location or ship.",
					box_code(p.a1))
			}
			where := subloc(ship)
			cond = (where == p.a1)
			if not {
				cond = !cond
			}
			if cond {
				return sout("%sat %s.", or_string(not, "not ", ""), box_name(p.a1))
			}
			break

		case 7: /* stack */
			//#if 0
			//                if (kind(p.a1) != T_char){
			//					return sout("%s is not a live character.", box_code(p.a1));
			//				}
			//#else
			if kind(p.a1) != T_char {
				break
			} /* just hang */
			//#endif

			cond = (stack_leader(c.who) == stack_leader(p.a1))
			if not {
				cond = !cond
			}
			if cond {
				return sout("%s is%s stacked with us.", box_name(p.a1), or_string(not, " not", ""))
			}
			break

		case 8: /* top */
			cond = (stack_leader(c.who) == c.who)
			if not {
				cond = !cond
			}
			if cond {
				return sout("we are%s the stack leader", or_string(not, " not", ""))
			}
			break

		case 9: /* ferry n */
			if !is_ship(p.a1) {
				return sout("%s is not a ship", box_code(p.a1))
			}
			cond = (subloc(p.a1) == subloc(c.who) && ferry_horn(p.a1) != FALSE)
			if not {
				cond = !cond
			}
			if cond {
				return sout("the ferry has%s signaled.", or_string(not, " not", ""))
			}
			break

		case 10: /* ship n */
			if kind(p.a1) != T_ship {
				return sout("%s is not a ship.", box_code(p.a1))
			}
			cond = (where_ship == subloc(p.a1))
			if not {
				cond = !cond
			}
			if cond {
				return sout("%s is%s here.", box_code(p.a1),
					or_string(not, " not", ""))
			}
			break

		case 11: /* rain */
			cond = (weather_here(province(c.who), sub_rain) != FALSE)
			if not {
				cond = !cond
			}
			if cond {
				return sout("it is%s raining.", or_string(not, " not", ""))
			}
			break

		case 12: /* fog */
			cond = (weather_here(province(c.who), sub_fog) != FALSE)
			if not {
				cond = !cond
			}
			if cond {
				return sout("it is%s foggy.", or_string(not, " not", ""))
			}
			break

		case 13: /* wind */
			cond = (weather_here(province(c.who), sub_wind) != FALSE)
			if not {
				cond = !cond
			}
			if cond {
				return sout("it is%s windy.", or_string(not, " not", ""))
			}
			break

		case 23: /* mist */
			cond = (weather_here(province(c.who), sub_mist) != FALSE)
			if not {
				cond = !cond
			}
			if cond {
				return sout("it is%s misty.", or_string(not, " not", ""))
			}
			break
		case 14: /* not */
			not = true
			setnot = true
			break

		case 15: /* owner */
			cond = (first_character(subloc(c.who)) == c.who)
			if not {
				cond = !cond
			}
			if cond {
				return sout("we are%s the first character here", or_string(not, " not", ""))
			}
			break

		case 19: /* clear */
			cond = (FALSE == weather_here(subloc(c.who), sub_fog) && FALSE == weather_here(subloc(c.who), sub_rain) && FALSE == weather_here(subloc(c.who), sub_wind))
			if not {
				cond = !cond
			}
			if cond {
				return sout("it is%s clear.", or_string(not, " not", ""))
			}
			break

		case 25: /* alone */
			cond = len(rp_loc_info(subloc(c.who)).here_list) == 1
			if not {
				cond = !cond
			}
			if cond {
				return sout("you are%s alone.", or_string(not, " not", ""))
			}
			break

		case 26: /* teach <n> */
			var item, bonus int
			cond = being_taught(c.who, p.a1, &item, &bonus) == TAUGHT_SPECIFIC
			if not {
				cond = !cond
			}
			if cond {
				return sout("instruction for %s is%s available.", box_code(p.a1), or_string(not, " not", ""))
			}
			break

		default:
			panic(fmt.Sprintf("assert(tag != %d)", p.tag))
		}
	}

	return ""
}

func v_wait(c *command) int {
	var s string

	if numargs(c) < 1 {
		wout(c.who, "Must say what condition to wait for.")
		return FALSE
	}

	clear_wait_parse(c)

	if s = check_wait_conditions(c); s != "" {
		wout(c.who, "Wait finished: %s", s)

		c.wait = 0
		c.inhibit_finish = true /* don't call d_wait */
		return TRUE
	}

	wait_list = append(wait_list, c.who)
	return TRUE
}

func d_wait(c *command) int {
	var s string

	if s = check_wait_conditions(c); s != "" {
		wout(c.who, "Wait finished: %s", s)
		wait_list = rem_value(wait_list, c.who)

		c.wait = 0
		c.inhibit_finish = true
		clear_wait_parse(c)
		return TRUE
	}

	return TRUE
}

func i_wait(c *command) int {

	wait_list = rem_value(wait_list, c.who)
	return TRUE
}

func v_split(c *command) int {
	lines := c.a
	bytes := c.b
	var p *entity_player
	var pl int

	pl = player(c.who)
	p = p_player(pl)

	if lines > 0 && lines < 1000 {
		lines = 1000
		out(c.who, "Minimum lines to split at is 1,000")
	}

	if bytes > 0 && bytes < 10000 {
		bytes = 10000
		out(c.who, "Minimum bytes to split at is 10,000")
	}

	p.split_lines = lines
	p.split_bytes = bytes

	if lines == 0 && bytes == 0 {
		out(c.who, "Reports will not be split when mailed.")
	} else if lines != FALSE && bytes != FALSE {
		out(c.who, "Reports will be split at %d lines or %d bytes, whichever limit is hit first.",
			lines, bytes)
	} else if lines != FALSE {
		out(c.who, "Reports will be split at %d lines.", lines)
	} else {
		out(c.who, "Reports will be split at %d bytes.", bytes)
	}

	return 0
}

func v_emote(c *command) int {
	target := c.a

	if numargs(c) < 2 {
		wout(c.who, "Usage: EMOTE <target> <message>")
		return FALSE
	}

	wout(target, string(c.parse[2]), box_name(target))

	return TRUE
}
