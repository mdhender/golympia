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

/*
 *  Trading-related skills and functions.
 *
 */

/*
 *  Dedicate Tower
 *  Fri Nov 15 13:44:30 1996 -- Scott Turner
 *
 *  Turn a generic tower into one for your guild.
 *
 */
func v_dedicate_tower(c *command) int {
	where := subloc(c.who)
	school := c.a
	var i int
	var p *entity_skill

	/*
	 *  School should be legitimate.
	 *
	 */
	if !valid_box(school) || kind(school) != T_skill ||
		skill_school(school) != school {
		wout(c.who, "To dedicate a tower, you must specify the skill school.")
		return FALSE
	}

	/*
	 *  You need to already know all the learnable and researchable
	 *  skills in this school.
	 *
	 */
	p = rp_skill(school)
	if p == nil {
		return FALSE
	}

	for i = 0; i < len(p.offered); i++ {
		if FALSE == has_skill(c.who, p.offered[i]) {
			wout(c.who, "You must master all available skills before creating a guild.")
			return FALSE
		}
	}
	for i = 0; i < len(p.research); i++ {
		if FALSE == has_skill(c.who, p.research[i]) {
			wout(c.who, "You must master all available skills before creating a guild.")
			return FALSE
		}
	}

	/*
	 *  You can't dedicate a tower if you're already in a guild and this
	 *  isn't that guild.
	 *
	 */
	if guild_member(c.who) != FALSE &&
		guild_member(c.who) != school {
		wout(c.who, "You cannot dedicate this tower, because you're already")
		wout(c.who, "a member of the %s Guild.", box_name(guild_member(c.who)))
		return FALSE
	}

	/*
	 *  Have to be in a tower.
	 *
	 */
	if subkind(where) != sub_tower {
		wout(c.who, "To dedicate a tower, you must be inside the tower.")
		return FALSE
	}

	/*
	 *  Have to be top dog.
	 *
	 *  Mon Oct 26 09:34:36 1998 -- Scott Turner
	 *
	 *  Maybe you really have to be alone, so you can't grandfather non-guild
	 *  members into a guild.
	 *
	 */
	if len(rp_loc_info(where).here_list) != 1 {
		wout(c.who, "Must be alone in a tower to dedicate.")
		return FALSE
	}

	/*
	 *  Need some gold for the dedication ceremony.
	 *
	 */
	if has_item(c.who, item_gold) < 250 {
		wout(c.who, "Dedicating a tower requires 250 gold.")
		return FALSE
	}

	/*
	 *  Must be in a city.
	 *
	 */
	if subkind(loc(where)) != sub_city {
		wout(c.who, "You may only create a guild in a city.")
		return FALSE
	}

	/*
	 *  Can't be any other guilds of the same sort here.
	 *
	 */
	for _, i = range loop_here(loc(where)) {
		if subkind(i) == sub_guild &&
			rp_subloc(i) != nil &&
			rp_subloc(i).guild == school {
			wout(c.who, "There is already a guild for %s in this city.",
				box_name(school))
			return FALSE
		}
	}

	/*
	 *  Tue Dec 29 11:33:45 1998 -- Scott Turner
	 *
	 *  Have to be a guild member already or ready to join the guild.
	 *
	 */
	if guild_member(c.who) != school &&
		FALSE == can_join_guild(c.who, school) {
		wout(c.who, "You must be a guild member or ready to join to dedicate a guild.")
		return FALSE
	}

	return TRUE
}

/*
 * Tue Apr 20 17:21:52 1999 -- Scott Turner
 *
 *  Encapsulated here because it is also called in seed_city.
 *
 */
func make_tower_guild(where, school int) {
	p := p_subloc(where)

	change_box_subkind(where, sub_guild)
	p.guild = school
	if options.guild_teaching != FALSE {
		p.teaches = append(p.teaches, school)
	}

}

func d_dedicate_tower(c *command) int {
	where := subloc(c.who)
	var i int
	p := rp_subloc(where)
	school := c.a

	/*
	 *  Have to be in a tower.
	 *
	 */
	if subkind(where) != sub_tower {
		wout(c.who, "To dedicate a tower, you must be inside the tower.")
		return FALSE
	}

	/*
	 *  Have to be top dog.
	 *
	 *  Mon Oct 26 09:34:36 1998 -- Scott Turner
	 *
	 *  Maybe you really have to be alone, so you can't grandfather non-guild
	 *  members into a guild.
	 *
	 */
	if len(rp_loc_info(where).here_list) != 1 {
		wout(c.who, "Must be alone in a tower to dedicate.")
		return FALSE
	}

	/*
	 *  Must be in a city.
	 *
	 */
	if subkind(loc(where)) != sub_city {
		wout(c.who, "You may only create a guild in a city.")
		return FALSE
	}

	/*
	 *  Can't be any other guilds of the same sort here.
	 *
	 */
	for _, i = range loop_here(loc(where)) {
		if subkind(i) == sub_guild &&
			rp_subloc(i) != nil &&
			rp_subloc(i).guild == school {
			wout(c.who, "There is already a guild for %s in this city.",
				box_name(school))
			return FALSE
		}
	}

	/*
	 *  Tue Dec 29 11:33:45 1998 -- Scott Turner
	 *
	 *  Have to be a guild member already or ready to join the guild.
	 *
	 */
	if guild_member(c.who) != school &&
		FALSE == can_join_guild(c.who, school) {
		wout(c.who, "You must be a guild member or ready to join to dedicate a guild.")
		return FALSE
	}

	/*
	 *  Need some gold for the dedication ceremony.
	 *
	 */
	if FALSE == charge(c.who, 250) {
		wout(c.who, "Dedicating a tower requires 250 gold.")
		return FALSE
	}

	/*
	 *  Make it a tower.
	 *
	 */
	if p == nil {
		wout(c.who, "For some reason, you cannot dedicate this tower.")
		return FALSE
	} else {
		make_tower_guild(where, school)
		wout(c.who, "%s now dedicated to %s.",
			box_name(where),
			just_name(school))
		join_guild(c.who, school)
		/*
		 *  Have to "touch" the location so it generates a location report.
		 *
		 */
		touch_loc(c.who)
	}
	return TRUE
}

/*
 *  RANDOM_TRADE_GOOD
 *  Fri Nov 15 16:16:45 1996 -- Scott Turner
 *
 *  Select a trade good at random, based on the relative weights.
 *
 */
func random_trade_good() int {
	var t int
	sofar := 0
	selected := 0

	for _, t = range loop_subkind(sub_trade_good) {
		if rp_item(t) != nil {
			sofar += rp_item(t).trade_good
			if rnd(1, sofar) <= rp_item(t).trade_good {
				selected = t
			}
		}
	}

	return selected
}

/*
 *  ADD_TRADING_PRODUCTION
 *  Fri Nov 15 16:08:24 1996 -- Scott Turner
 *
 *  Add a production of a trading good to this market.  Don't pick
 *  one that is already there.
 *
 */
func add_trade(where, type_ int) {
	var new_tg int
	other := or_int((type_ == PRODUCE), CONSUME, PRODUCE)
	var t *trade

	for {
		new_tg = random_trade_good()
		if find_trade(where, type_, new_tg) != nil || find_trade(where, other, new_tg) != nil {
			continue
		}
		break
	}

	/*
	 *  Create a new trade and add it to the market.
	 *
	 */
	t = new_trade(where, type_, new_tg)
	t.qty = rp_item(new_tg).trade_good
	t.cost = rp_item(new_tg).base_price

	/*
	 *  Mon Feb  7 08:16:06 2000 -- Scott Turner
	 *  Put in the buy/sell as well.
	 *
	 *  Wed Feb  9 06:51:01 2000 -- Scott Turner
	 *
	 *  Whoops, only 1/8 the qty.
	 *
	 */
	other = or_int((type_ == PRODUCE), SELL, BUY)
	t = new_trade(where, other, new_tg)
	t.qty = rp_item(new_tg).trade_good / NUM_MONTHS
	t.cost = rp_item(new_tg).base_price
}

/*
 *  Big city trades.
 *
 */
func update_big_city_trades(where int) {
	/*
	 *  Remove small city trades.
	 *
	 */
	delete_city_trade(where, item_lumber)
	delete_city_trade(where, item_stone)
	/*
	 *  Update (or add) big city trades.
	 *
	 */
	update_city_trade(where, CONSUME, item_lana_bark, rnd(1, 5), rp_item(item_lana_bark).base_price, 0)
	update_city_trade(where, CONSUME, item_pretus_bones, rnd(1, 5), rp_item(item_pretus_bones).base_price, 0)
	update_city_trade(where, CONSUME, item_mallorn_wood, rnd(1, 5), rp_item(item_mallorn_wood).base_price, 0)
	update_city_trade(where, CONSUME, item_yew, rnd(1, 5), rp_item(item_yew).base_price, 0)
	update_city_trade(where, CONSUME, item_farrenstone, rnd(1, 5), rp_item(item_farrenstone).base_price, 0)
	update_city_trade(where, CONSUME, item_spiny_root, rnd(1, 5), rp_item(item_spiny_root).base_price, 0)
	update_city_trade(where, CONSUME, item_avinia_leaf, rnd(1, 5), rp_item(item_avinia_leaf).base_price, 0)
}

/*
 *  Small city trades.
 *
 */
func update_small_city_trades(where int) {
	amount := (10000 - has_item(province(where), item_peasant)) / 400
	update_city_trade(where, CONSUME, item_lumber, amount/2, rp_item(item_lumber).base_price, 0)
	update_city_trade(where, CONSUME, item_stone, amount, rp_item(item_stone).base_price, 0)
}

/*
 *  Override causes cities which only produce a good once per year
 *  to produce it now anyway.  This is useful for epoch city trade
 *  seeding.
 *
 *  Thu Dec  2 05:52:22 1999 -- Scott Turner
 *
 *  Renamed from loc_trade_sup to do_production.
 *
 */
func do_production(where int, override bool) {
	var t *trade

	for _, t = range loop_trade(where) {
		okay := true

		if t.month_prod != FALSE && !override {
			this_month := oly_month(&sysclock) - 1
			next_month := (this_month + 1) % NUM_MONTHS
			prod_month := t.month_prod - 1

			if next_month != prod_month {
				okay = false
			}
		}

		if t.kind == PRODUCE && okay {
			newTrade := new_trade(where, SELL, t.item)

			if newTrade.qty < t.qty {
				newTrade.qty = t.qty
			}

			newTrade.cost = t.cost
			newTrade.cloak = t.cloak
		} else if t.kind == CONSUME {
			newTrade := new_trade(where, BUY, t.item)
			if newTrade.qty < t.qty {
				newTrade.qty = t.qty
			}

			newTrade.cost = t.cost
			newTrade.cloak = t.cloak
		}
	}
}

/*
 *  Update_market
 *  Tue Jan 18 08:43:27 2000 -- Scott Turner
 *
 *  Update one market.
 *
 *  Tue Jan 18 10:16:25 2000 -- Scott Turner
 *
 *  Refactor BUY/SELL.  Modify "change" to negative for SELL transactions.
 *  Check to see if the item has been sold or bought with a != instead
 *  of a >= or a <=.  Everything else should be the same.
 *
 */
func update_market(where int) {
	var t *trade
	var change, other, bp int

	for _, t = range loop_trade(where) {
		/*
		 *  Ignore opium!
		 *
		 */
		if t.item == item_opium {
			continue
		}

		/*
		 * Calculate what the change in price, if any, will be.
		 *
		 */
		change = (t.cost * (9 + rnd(1, 11))) / 100
		if change < 1 {
			change = 1
		}

		/*
		 *  If this isn't a sell or buy, skip over it.
		 *
		 */
		if t.kind != SELL && t.kind != BUY {
			continue
		}

		/*
		 *  other is the corresponding other part of the
		 *  transaction, i.e., SELL-PRODUCE and BUY-CONSUME.
		 *
		 */
		other = CONSUME

		/*
		 *  If we're a SELL we switch the other and the "polarity" of the price change.
		 *
		 */
		if t.kind == SELL {
			change = -change
			other = PRODUCE
		}

		/*
		 *  Find the corresponding other.
		 */
		newTrade := new_trade(where, other, t.item)
		assert(newTrade != nil)

		/*
		 *  Add in the price delta if quantity unchanged from last month;
		 *  otherwise subtract it.
		 *
		 *  Tue Jan 18 11:47:12 2000 -- Scott Turner
		 *
		 *  City goods don't always have an "old_qty" at first; we'll assume
		 *  those are unsold.
		 *
		 */
		if FALSE == t.old_qty || t.old_qty == t.qty {
			t.counter++
			if rp_item(t.item) == nil {
				t.cost += change
				newTrade.cost += change
			}
		} else {
			t.counter = 0
			if rp_item(t.item) == nil {
				t.cost -= change
				newTrade.cost -= change
			}
		}

		/*
			                     *  If a good has reached the end of its "counter", then delete the two trades.
			                     *
			                     *  Tue Jan 18 12:46:59 2000 -- Scott Turner
			                     *
			                     *  This is a bad thing for the "constant" city goods, like fish, etc.
								 *  It's a good thing for trade goods, unique items, etc.
			                     *  Hmm.
		*/
		if (rp_item(t.item).trade_good != FALSE || item_unique(t.item) != FALSE) &&
			t.counter > options.market_age {
			wout(gm_player, "Deleting good %s from %s.",
				box_name(t.item),
				box_name(where))
			bx[where].trades = bx[where].trades.rem_value(t)
			bx[where].trades = bx[where].trades.rem_value(newTrade)
			continue
		}

		/*
		 *  Make sure our price is legitimate.
		 *
		 *  Tue Jan 18 11:11:20 2000 -- Scott Turner
		 *
		 *  Some items may not have a base price; in which case
		 *  we'll only insist that the price stay positive.
		 *
		 *  Thu Jun 15 18:58:46 2000 -- Scott Turner
		 *
		 *  Uh, let's keep the price under 50 if it doesn't have
		 *  a bp.
		 *
		 */
		bp = rp_item(t.item).base_price

		if bp != FALSE && t.cost > bp*2 {
			t.cost = bp * 2
			newTrade.cost = t.cost
		}

		if bp != FALSE && t.cost < bp/2 {
			t.cost = bp / 2
			newTrade.cost = t.cost
		}

		if bp == FALSE && t.cost < 1 {
			t.cost = 1
			newTrade.cost = 1
		}

		if bp == FALSE && t.cost > 50 {
			t.cost = 50
			newTrade.cost = 50
		}

		/*
		 *  Now save the current quantity for next month.
		 *
		 */
		t.old_qty = t.qty

	}

	/*
	 *   Tue Jan 18 10:51:16 2000 -- Scott Turner
	 *
	 *   Add in test for prod_month (unused?).
	 *
	 *   Tue Jan 18 10:53:32 2000 -- Scott Turner
	 *
	 *   In the trade guilds, qty is the yearly production, and
	 *   you can't exceed that (generated qty/8 per month).  Unfortunately,
	 *   in the city markets, qty is the monthly production *and* the limit.
	 *   So we need a special case test in here (or we could modify the
	 *   databases :-().
	 *
	 */
	for _, t = range loop_trade(where) {
		other := SELL

		if t.kind != CONSUME && t.kind != PRODUCE {
			continue
		}
		if t.kind == CONSUME {
			other = BUY
		}

		/*
		 *  A production or consumption with no qty can be deleted
		 *  and skipped.
		 *
		 */
		if t.qty == 0 {
			bx[where].trades = bx[where].trades.rem_value(t)
			continue
		}

		/*
		 *  Some productions only happen 1x year.
		 *
		 */
		if t.kind == PRODUCE && t.month_prod != FALSE {
			this_month := oly_month(&sysclock) - 1
			next_month := (this_month + 1) % NUM_MONTHS
			prod_month := t.month_prod - 1

			if next_month != prod_month {
				continue
			}
		}

		/*
		 *  Produce or consume the monthly amount, up to the
		 *  yearly limit.
		 *
		 */
		newTrade := new_trade(where, other, t.item)
		if newTrade.qty < t.qty {
			if rp_item(t.item).trade_good != FALSE {
				newTrade.qty += (t.qty / 8) /* per month */
			} else {
				newTrade.qty = t.qty
			}
			newTrade.old_qty = newTrade.qty
		}
		newTrade.cost = t.cost
		newTrade.cloak = t.cloak
	}
}

func add_trade_goods(where int) {
	var t *trade
	var i int
	produce, consume := 0, 0

	for _, t = range loop_trade(where) {
		/*
		 *  Track how many of each; we may need to add some
		 *  at the end.
		 *
		 *  Produce or consume the monthly amount, up to the
		 *  yearly limit.
		 *
		 */
		if t.kind == PRODUCE {
			produce++
		}

		if t.kind == CONSUME {
			consume++
		}
	}

	/*
	 *  We need to have 3 productions and 3 consumptions.
	 *
	 */
	if produce < 3 {
		/*
		 *  Add a production.
		 *
		 */
		for i = 0; i < (3 - produce); i++ {
			add_trade(where, PRODUCE)
		}
	}
	if consume < 3 {
		/*
		 *  Add a consumption.
		 *
		 */
		for i = 0; i < (3 - consume); i++ {
			add_trade(where, CONSUME)
		}
	}
}

/*
 *  Update_Markets
 *  Tue Nov 12 12:33:30 1996 -- Scott Turner
 *
 *  Fri Sep 25 08:20:12 1998 -- Scott Turner
 *
 *  Split the buy/sell price updates and the consume/produce
 *  quantity updates into two parts.
 *
 *  Tue Jan 18 08:41:30 2000 -- Scott Turner
 *
 *  Extend this to cover all markets; need to add special code to
 *  distinguish trade good special cases.  Non-trade goods should not
 *  be "cycled".
 *
 */
func update_markets() {
	var where int
	stage("update_markets()")

	for _, where = range loop_guild() {
		if is_guild(where) == sk_trading {
			update_market(where)
			add_trade_goods(where)
		}
	}

	/*
	 *  Wed Sep 15 17:39:39 1999 -- Scott Turner
	 *
	 *  Now adjust city markets for common goods.
	 *
	 *  Harbor city: buys fish
	 *  City < 10000: buys 10 wood per (10k-pop)
	 *                at bp = ?
	 *              : buys 10 stone per (10k-pop)
	 *  City == 10K : buys 1-5 lana bark, avinia,
	 *                spiny root, farrenstone, yew,
	 *                mallor, pretus bones
	 *  All cities  : buys 1-5 wild horses
	 *  add_city_trade(int where, int kind, int item,
	 *                 int qty, int cost, int month)
	 *  delete_city_trade(int where, int item)
	 *
	 *  Wed Nov 10 08:23:13 1999 -- Scott Turner
	 *
	 *  We should do these in the same manner as fish, etc.
	 *  (e.g., with PRODUCE/CONSUME) to make everything consistent.
	 *
	 *  Thu Dec  2 05:28:53 1999 -- Scott Turner
	 *
	 *  Folded in the calls from location_trades.
	 *
	 */
	for _, where = range loop_city() {
		/*
		 *  City-size specific trades.
		 *
		 */
		if has_item(where, item_peasant) > 9500 {
			update_big_city_trades(where)
		} else {
			update_small_city_trades(where)
		}

		/*
		 *  Everyone buys wild horses.
		 *
		 */
		update_city_trade(where, CONSUME, item_wild_horse, rnd(1, 10),
			rp_item(item_wild_horse).base_price, 0)

		/*
		 *  Update the markets.
		 *
		 */
		update_market(where)

		/*
		 *
		 * Opium
		 *
		 */
		opium_market_delta(where)

		if in_faery(where) || in_clouds(where) {
			trade_suffuse_ring(where)
		}

		if !in_faery(where) &&
			!in_clouds(where) &&
			!in_hades(where) {
			add_scrolls(where)
		}

	}
}

/*
 *  Mon Nov 25 11:41:48 1996 -- Scott Turner
 *
 *  Smuggle goods hangs (or removes) a "smuggling" effect on the user.
 *  Then, based on how much smuggling experience he has, he may avoid
 *  the costs to enter a city.
 *
 */
func v_smuggle_goods(c *command) int {
	flag := c.a

	if flag != FALSE && get_effect(c.who, ef_smuggle_goods, 0, 0) != FALSE {
		wout(c.who, "You are already prepared to smuggle goods.")
		return FALSE
	}

	if FALSE == flag && FALSE == get_effect(c.who, ef_smuggle_goods, 0, 0) {
		wout(c.who, "You are not smuggling at this time.")
		return FALSE
	}

	return TRUE
}

func d_smuggle_goods(c *command) int {
	flag := c.a

	if flag != FALSE {
		if get_effect(c.who, ef_smuggle_goods, 0, 0) != FALSE {
			wout(c.who, "You are already prepared to smuggle goods.")
			return FALSE
		}
		add_effect(c.who, ef_smuggle_goods, 0, -1, 1)
		wout(c.who, "You are now prepared to smuggle goods.")
		return TRUE
	}

	if FALSE == flag {
		if FALSE == get_effect(c.who, ef_smuggle_goods, 0, 0) {
			wout(c.who, "You are not smuggling at this time.")
			return FALSE
		}
		delete_effect(c.who, ef_smuggle_goods, 0)
		wout(c.who, "You are no longer smuggling goods.")
		return TRUE
	}

	return 0 // todo: should this return something?
}

/*
 *
 */
func v_smuggle_men(c *command) int {
	flag := c.a

	if flag != FALSE && get_effect(c.who, ef_smuggle_men, 0, 0) != FALSE {
		wout(c.who, "You are already prepared to smuggle goods.")
		return FALSE
	}

	if FALSE == flag && FALSE == get_effect(c.who, ef_smuggle_men, 0, 0) {
		wout(c.who, "You are not smuggling at this time.")
		return FALSE
	}

	return TRUE
}

func d_smuggle_men(c *command) int {
	flag := c.a

	if flag != FALSE {
		if get_effect(c.who, ef_smuggle_men, 0, 0) != FALSE {
			wout(c.who, "You are already prepared to smuggle goods.")
			return FALSE
		}
		add_effect(c.who, ef_smuggle_men, 0, -1, 1)
		wout(c.who, "You are now prepared to smuggle goods.")
		return TRUE
	}

	if FALSE == flag {
		if FALSE == get_effect(c.who, ef_smuggle_men, 0, 0) {
			wout(c.who, "You are not smuggling at this time.")
			return FALSE
		}
		delete_effect(c.who, ef_smuggle_men, 0)
		wout(c.who, "You are no longer smuggling goods.")
		return TRUE
	}

	return TRUE
}

/*
 *  Wed Nov 27 12:25:46 1996 -- Scott Turner
 *
 *  Build wagons is just a production skill, so everything necessary
 *  is encoded in lib/skill
 *
 */

/*
 *  Mon Nov 25 10:50:57 1996 -- Scott Turner
 *
 *  Is a good traded by this place?
 *
 */
func traded_here(where, good int) *trade {
	var t *trade

	/*
	 *  Fri Nov 13 12:52:55 1998 -- Scott Turner
	 *
	 *  Loop_trade only runs through the trades hung on "where";
	 *  it doesn't pick up any from the traders in where.
	 *
	 */
	for _, t = range loop_trade(where) {
		if t.item == good &&
			(t.kind == SELL || t.kind == BUY) {
			return t
		}
	}
	return nil
}

/*
 *  Mon Nov 25 10:50:57 1996 -- Scott Turner
 *
 *  Is a good traded by this place?
 *
 */
func produced_here(where, good int) *trade {
	var t *trade

	/*
	 *  Fri Nov 13 12:52:55 1998 -- Scott Turner
	 *
	 *  Loop_trade only runs through the trades hung on "where";
	 *  it doesn't pick up any from the traders in where.
	 *
	 */
	for _, t = range loop_trade(where) {
		if t.item == good &&
			(t.kind == CONSUME || t.kind == PRODUCE) {
			return t
		}
	}
	return nil
}

/*
 *  Mon Nov 25 10:41:35 1996 -- Scott Turner
 *
 *  Increase the demand for a good -- raise it's price 5-8%.  Note that
 *  you can't raise a price over base*2.
 *
 */
func v_increase_demand(c *command) int {
	where := subloc(c.who)
	good := c.a

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good traded in this guild (by the guild).
	 *
	 */
	if FALSE == good || nil == traded_here(where, good) {
		wout(c.who, "That good is not traded here.")
		return FALSE
	}

	return TRUE
}

func d_increase_demand(c *command) int {
	where := subloc(c.who)
	good := c.a
	var change int
	var t *trade

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good traded in this guild (by the guild).
	 *
	 *  Fri Nov 13 12:54:20 1998 -- Scott Turner
	 *
	 *  Note that all these functions implicitly assume that a loc
	 *  has only one BUY or SELL for a particular good.
	 *
	 */
	tradedHere := FALSE == good
	if tradedHere {
		t = traded_here(where, good)
		tradedHere = t != nil
	}
	if !tradedHere {
		wout(c.who, "That good is not traded here.")
		return FALSE
	}

	/*
	 *  Maybe a problem?
	 *
	 */
	if nil == rp_item(t.item) {
		wout(c.who, "I'm confused about that trade good, tell the GM.")
		return FALSE
	}

	/*
	 *  No more than 2*base_price
	 *
	 */
	if rp_item(t.item) != nil && rp_item(t.item).base_price*2 <= t.cost {
		wout(c.who, "The demand for %s is straining the market and cannot be further increased.", box_name(t.item))
		return FALSE
	}

	/*
	 *  Up the price, etc.
	 *
	 */
	change = (t.cost * (5 + rnd(1, 3))) / 100
	if change < 1 {
		change = 1
	}
	t.cost += change
	wout(c.who, "The demand for %s increases; the new price is %s.",
		box_name(t.item),
		gold_s(t.cost))

	/*
	 *  Also modify the consume/produce price, as necessary.
	 *
	 */
	if t = produced_here(where, good); t != nil {
		t.cost += change
	}

	return TRUE
}

/*
 *
 */
func v_decrease_demand(c *command) int {
	where := subloc(c.who)
	good := c.a

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good traded in this guild (by the guild).
	 *
	 */
	if FALSE == good || nil == traded_here(where, good) {
		wout(c.who, "That good is not traded here.")
		return FALSE
	}

	return TRUE
}

func d_decrease_demand(c *command) int {
	where := subloc(c.who)
	good := c.a
	var change int
	var t *trade

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good traded in this guild (by the guild).
	 *
	 */
	tradedHere := FALSE == good
	if tradedHere {
		t = traded_here(where, good)
		tradedHere = t != nil
	}
	if !tradedHere {
		wout(c.who, "That good is not traded here.")
		return FALSE
	}

	/*
	 *  Maybe a problem?
	 *
	 */
	if nil == rp_item(t.item) {
		wout(c.who, "I'm confused about that trade good, tell the GM.")
		return FALSE
	}

	/*
	 *  No more than 2*base_price
	 *
	 */
	if rp_item(t.item) != nil && rp_item(t.item).base_price/2 >= t.cost {
		wout(c.who, "The demand for %s has bottomed out and cannot be further decreased.", box_name(t.item))
		return FALSE
	}

	/*
	 *  Up the price, etc.
	 *
	 */
	change = (t.cost * (5 + rnd(1, 3))) / 100
	if change < 1 {
		change = 1
	}
	t.cost -= change
	wout(c.who, "The demand for %s decreases; the new price is %s.",
		box_name(t.item),
		gold_s(t.cost))

	/*
	 *  Also modify the consume/produce price, as necessary.
	 *
	 */
	if t = produced_here(where, good); t != nil {
		t.cost -= change
	}

	return TRUE
}

/*
 *  Mon Nov 25 11:01:03 1996 -- Scott Turner
 *
 *  Increase the # of a good demanded or offered.
 */
func v_increase_supply(c *command) int {
	where := subloc(c.who)
	good := c.a

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good traded in this guild (by the guild).
	 *
	 */
	if FALSE == good || nil == traded_here(where, good) || produced_here(where, good) == nil {
		wout(c.who, "That good is not traded here.")
		return FALSE
	}

	return TRUE
}

func d_increase_supply(c *command) int {
	where := subloc(c.who)
	good := c.a
	var p, t *trade

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good produced in this guild (by the guild).
	 *
	 */
	tradedHere := FALSE == good
	if tradedHere {
		t = traded_here(where, good)
		tradedHere = t != nil
		if tradedHere {
			p = produced_here(where, good)
			tradedHere = p != nil
		}
	}
	if tradedHere {
		wout(c.who, "That good is not produced here.")
		return FALSE
	}

	/*
	 *  Maybe a problem?
	 *
	 */
	if nil == rp_item(t.item) {
		wout(c.who, "I'm confused about that trade good, tell the GM.")
		return FALSE
	}

	/*
	 *  No more than the year's supply
	 *
	 */
	if rp_item(t.item) != nil && rp_item(t.item).trade_good <= t.qty {
		wout(c.who, "The supply of %s is exhausted and cannot be further increased.", box_name(t.item))
		return FALSE
	}

	/*
	 *  Add up to 15-33% more of the good; handle amounts < 1 as
	 *  chance...
	 *
	 */
	change := ((float64(p.qty) / 8) * float64(14+rnd(1, 19))) / 100.0
	if change < 1 && rnd(1, 100) < int(change*100) {
		change = 1.0
	}
	t.qty += int(change)
	wout(c.who, "The supply of %s increases; the new amount available is %d.",
		box_name(t.item), t.qty)

	return TRUE
}

/*
 *
 */
func v_decrease_supply(c *command) int {
	where := subloc(c.who)
	good := c.a

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good traded in this guild (by the guild).
	 *
	 */
	if FALSE == good || nil == traded_here(where, good) || nil == produced_here(where, good) {
		wout(c.who, "That good is not traded here.")
		return FALSE
	}

	return TRUE
}

func d_decrease_supply(c *command) int {
	where := subloc(c.who)
	good := c.a
	var change float64
	var p, t *trade

	/*
	 *  You need to be in the guild where the good is being traded.
	 *
	 *  Fri Nov 13 12:49:25 1998 -- Scott Turner
	 *
	 *  Now works in any market, but only on goods sold/bought by the
	 *  market.
	 */
	if FALSE == market_here(where) {
		wout(c.who, "You must be in a market to use this skill.")
		return FALSE
	}

	/*
	 *  Needs to be a good traded in this guild (by the guild).
	 *
	 */
	tradedHere := FALSE == good
	if tradedHere {
		t = traded_here(where, good)
		tradedHere = t != nil
		if tradedHere {
			p = produced_here(where, good)
			tradedHere = p != nil
		}
	}
	if tradedHere {
		wout(c.who, "That good is not traded here.")
		return FALSE
	}

	/*
	 *  Maybe a problem?
	 *
	 */
	if nil == rp_item(t.item) {
		wout(c.who, "I'm confused about that trade good, tell the GM.")
		return FALSE
	}

	/*
	 *  No more than the year's supply
	 *
	 */
	if t.qty < 1 {
		wout(c.who, "The supply of %s cannot be further decreased.", box_name(t.item))
		return FALSE
	}

	/*
	 *  Add up to 15-33% more of the good; handle amounts < 1 as
	 *  chance...
	 *
	 */
	change = ((float64(p.qty) / 8) * float64(14+rnd(1, 19))) / 100.0
	if change < 1 && rnd(1, 100) < int(change*100) {
		change = 1.0
	}
	t.qty -= int(change)
	if t.qty < 0 {
		t.qty = 0
	}
	wout(c.who, "The supply of %s decreases; the new amount available is %d.",
		box_name(t.item), t.qty)

	return TRUE
}

/*
 *  Sat Nov 23 10:38:41 1996 -- Scott Turner
 *
 *  Implement this by hanging an effect on the guy.  If hiding nothing,
 *  then unhide the hidden item.
 */
func v_hide_item(c *command) int {
	item := c.a

	/*
	 *  Is he unhiding?
	 *
	 */
	if FALSE == item && FALSE == get_effect(c.who, ef_hide_item, 0, 0) {
		wout(c.who, "You don't have anything to unhide.")
		return FALSE
	}

	/*
	 *  Is it an item?
	 *
	 */
	if item != FALSE && (!valid_box(item) || kind(item) != T_item) {
		wout(c.who, "You can't hide such a thing.")
		return FALSE
	}

	/*
	 *  Does he have the item?
	 *
	 */
	if item != FALSE && FALSE == has_item(c.who, item) {
		wout(c.who, "You do not have that item to hide.")
		return FALSE
	}

	return TRUE
}

/*
 *  Tue Nov 26 16:02:34 1996 -- Scott Turner
 *
 *  For a unique item, I guess we want to move it onto some
 *  appropriate entity while it is "hidden".
 *
 */
func d_hide_item(c *command) int {
	item := c.a
	what := get_effect(c.who, ef_hide_item, 0, 0)

	/*
	 *  Is he unhiding?
	 *
	 */
	if FALSE == item {
		if FALSE == what {
			wout(c.who, "You don't have anything to unhide.")
			return FALSE
		}
		/*
		 *  Otherwise unhide it.
		 *
		 */
		delete_effect(c.who, ef_hide_item, 0)
		if item_unique(what) != FALSE {
			move_item(indep_player, c.who, what, 1)
		} else {
			gen_item(c.who, what, 1)
		}
		wout(c.who, "You unhide one %s.", box_name(what))
		return TRUE
	}

	/*
	 *  Is it an item?
	 *
	 */
	if !valid_box(item) || kind(item) != T_item {
		wout(c.who, "You can't hide such a thing.")
		return FALSE
	}

	/*
	 *  Does he have the item?
	 *
	 */
	if FALSE == has_item(c.who, item) {
		wout(c.who, "You do not have that item to hide.")
		return FALSE
	}

	/*
	 *  Possibly unhide.
	 *
	 */
	if what != FALSE {
		delete_effect(c.who, ef_hide_item, 0)
		if item_unique(what) != FALSE {
			move_item(indep_player, c.who, what, 1)
		} else {
			gen_item(c.who, what, 1)
		}
		wout(c.who, "You unhide one %s.", box_name(what))
	}

	/*
	 *  Hang the effect and delete the item from his possession.
	 *
	 */
	if FALSE == item_unique(item) {
		consume_item(c.who, item, 1)
	} else {
		move_item(c.who, indep_player, item, 1)
	}
	add_effect(c.who, ef_hide_item, 0, 1, item)
	wout(c.who, "You hide one %s.", box_name(item))

	return TRUE
}

/*
 *  Sun Nov 24 11:29:12 1996 -- Scott Turner
 *
 *  Similar, but put it in the hidden money slot.
 *
 *  The argument is the total amount of money hidden.
 *
 */
func v_hide_money(c *command) int {
	amount := c.a
	what := get_effect(c.who, ef_hide_money, 0, 0)

	/*
	 *  None to unhide.
	 *
	 */
	if FALSE == amount && FALSE == what {
		wout(c.who, "You don't have any money to unhide.")
		return FALSE
	}

	return TRUE
}

func d_hide_money(c *command) int {
	amount := c.a
	what := get_effect(c.who, ef_hide_money, 0, 0)

	/*
	 *  None to unhide.
	 *
	 */
	if FALSE == amount && FALSE == what {
		wout(c.who, "You don't have any money to unhide.")
		return FALSE
	}

	/*
	 *  Does he have enough money?
	 *
	 */
	if what+has_item(c.who, item_gold) < amount {
		wout(c.who, "You don't have %s to hide, hiding %s.",
			gold_s(amount-what), has_item(c.who, item_gold))
		what = has_item(c.who, item_gold)
	}

	/*
	 *  Adjust the levels.
	 *
	 */
	if amount-what > 0 {
		charge(c.who, amount-what)
		wout(c.who, "Adding %s to the amount in hiding.", gold_s(amount-what))
	} else {
		gen_item(c.who, item_gold, what-amount)
		wout(c.who, "Removing %s from the amount in hiding.",
			gold_s(what-amount))
	}
	delete_effect(c.who, ef_hide_money, 0)
	add_effect(c.who, ef_hide_money, 0, 1, amount)
	wout(c.who, "You now have %s in hiding.", gold_s(amount))

	return TRUE
}
