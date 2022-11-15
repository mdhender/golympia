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

//static int buyer_comp(const void *, const void *);
//static int seller_comp(const void *, const void *);

/* #define NEW */

/*
 *  Explanation of trade.cloak field:
 *
 *	0	normal -- open buy or sell, list in market report
 *	1	cloak trader, but list in market report
 *	2	invisible -- don't list in market report, cloak trader
 */

/*
 *  How it works
 *
 *  Each unit has a list of possible trades.
 *  Each trade is either a buy or a sell.
 *  When a trade is entered with the BUY or SELL command, the
 *	list of possible trades from other units in the city
 *	is consulted for a possible match.
 *  If no match is found, the trade is added to the unit's list
 *	of pending trades.
 *  When a unit enters the city, their trades are scanned to see if
 *	any match the pending trades of other units in the city.
 *  When an item is added to a unit, a check is made to see if the
 *	addition might validate a pending trade.  If so, we will try
 *	to match pending trades for the unit when the command is
 *	finished running.
 *	(perhaps this should happen at the end of the day?)
 */

func market_here(who int) int {

	for who > 0 && subkind(who) != sub_city && !(subkind(who) == sub_guild && rp_subloc(who) != nil && rp_subloc(who).guild == sk_trading) {
		if is_ship(who) {
			return 0
		}

		who = loc(who)
	}

	return who
}

func clear_all_trades(who int) {
	t := 0
	for _, t = range loop_trade(who) {
		my_free(t)
	}

	bx[who].trades = nil
}

func buyer_comp(a, b *trade) int {
	if a.cost == b.cost {
		return b.sort - a.sort
	}
	return b.cost - a.cost
}

func seller_comp(a, b *trade) int {
	if a.cost == b.cost {
		return a.sort - b.sort
	}
	return a.cost - b.cost
}

func seller_list(where, except int) []*trade {
	var l []*trade
	var i int
	count := 0

	l = nil

	for _, i = range loop_char_here(where) {
		/*
		 *  Don't trade with ourselves, moving characters, or prisoners
		 */
		if i == except || is_prisoner(i) || char_moving(i) != FALSE {
			continue
		}

		for _, t := range loop_trade(i) {
			if t.kind == SELL {
				t.sort, count = count, count+1
				l = append(l, t)
			}
		}

	}

	////#ifndef NEW_TRADE
	//if NEW_TRADE == FALSE {
	//    // sort before city
	//    // qsort(l, len(l), sizeof(struct trade *), seller_comp);
	//    sort.Slice(l, func(i, j int) bool {
	//        return seller_comp(l[i], l[j]) < 0
	//    })
	//}
	////#endif // NEW_TRADE

	for _, t := range loop_trade(where) {
		if t.kind == SELL {
			l = append(l, t)
		}
	}

	// sort after everyone, not excluding the city itself
	// qsort(l, len(l), sizeof(struct trade *), seller_comp);
	sort.Slice(l, func(i, j int) bool {
		return seller_comp(l[i], l[j]) < 0
	})

	return l
}

func buyer_list(where, except int) []*trade {
	var l []*trade
	var i int

	l = nil

	for _, i = range loop_char_here(where) {
		/*
		 *  Don't trade with ourselves, moving characters, or prisoners
		 */
		if i == except || is_prisoner(i) || char_moving(i) != FALSE {
			continue
		}

		for _, t := range loop_trade(i) {
			if t.kind == BUY {
				l = append(l, t)
			}
		}

	}

	for _, t := range loop_trade(where) {
		if t.kind == BUY {
			l = append(l, t)
		}
	}

	// sort everyone -- buyers high to low.
	// qsort(l, len(l), sizeof(struct trade *), buyer_comp);
	sort.Slice(l, func(i, j int) bool {
		return buyer_comp(l[i], l[j]) < 0
	})

	return l
}

/*
 *  If we're a buyer, reduce qty to how many we can afford, given cost
 *  If we're a seller, reduce qty to how many we actually have to sell
 *  If we're a city, leave qty alone, since we just gen the gold or item
 *
 *  Also reduce such that we don't exceed our own have_left requirement
 */

func reduce_qty(t *trade, cost int) int {
	var has int

	if kind(t.who) == T_loc {
		return t.qty
	}

	if t.kind == BUY {
		has = has_item(t.who, item_gold) - t.have_left
		if has < 0 {
			has = 0
		}

		return min(t.qty, has/cost)
	}

	if t.kind == SELL {
		has = has_item(t.who, t.item) - t.have_left
		if has < 0 {
			has = 0
		}

		return min(t.qty, has)
	}

	panic("!reached")
}

func attempt_trade(buyer, seller *trade, cost int) int {
	item := buyer.item
	var qty int
	var art int
	var item_cost int
	tariff := 0
	var where, garrison int
	var seller_s string
	var buyer_s string
	//extern int gold_pot_basket;
	//extern int gold_trade;
	//extern int gold_opium;
	//extern int gold_fish;
	//extern int gold_tariffs;

	assert(buyer != nil)
	assert(seller != nil)
	assert(buyer.item == seller.item)

	/*
	 *  seller.cost is the minimum price the seller is willing
	 *  	to accept.
	 *
	 *  buyer.cost is the maximum price the buyer is willing to
	 *  	pay.
	 *
	 *  cost is initially the "next best offer" for this good,
	 *      plus one gold, i.e., the price of the transaction has
	 *      to be at least this or the seller would not accept.
	 *      If there is no initial cost, we set this to be the
	 *      seller's minimum.
	 *
	 *  So we can reject this trade if:
	 *
	 *  (1) the maximum the buyer will pay is less than "cost".
	 *  (2) the minimum the seller will accept is more than the
	 *      maximum the buyer will pay.
	 */

	/*
	 *  Temporarily this holds the cost per item.
	 *
	 */
	if cost == FALSE || cost < seller.cost {
		cost = seller.cost
	}

	/*
	 *  Buyer.cost too low.
	 */
	if buyer.cost < cost {
		return 0
	}
	/*
	 *  Seller too high.
	 *
	 */
	if seller.cost > buyer.cost {
		return 0
	}
	var buyer_qty int
	var seller_qty int

	buyer_qty = reduce_qty(buyer, seller.cost)
	seller_qty = reduce_qty(seller, seller.cost)

	qty = min(buyer_qty, seller_qty)

	if qty <= 0 {
		return 0
	}

	item_cost = cost
	cost = cost * qty

	/*
	 *  Determine where
	 *
	 */
	if kind(buyer.who) == T_loc {
		where = buyer.who
	} else if kind(seller.who) == T_loc {
		where = seller.who
	} else {
		where = subloc(buyer.who)
	}

	/*
	 *  Fri Nov 22 11:16:18 1996 -- Scott Turner
	 *
	 *  Don't charge a tariff to anyone who rules here.
	 *
	 *  Thu Nov 12 07:07:26 1998 -- Scott Turner
	 *
	 *  Don't charge a tariff if there's no garrison!
	 *
	 *  Fri Nov 13 13:21:24 1998 -- Scott Turner
	 *
	 *  Halve the tariff if the seller has sk_avoid_taxes.
	 *
	 */
	garrison = garrison_here(province(where))
	if rp_subloc(where) != nil &&
		garrison != FALSE &&
		may_rule_here(seller.who, where) == FALSE {
		tariff = (cost * rp_subloc(where).tax_market) / 100

		if has_skill(seller.who, sk_avoid_taxes) != FALSE {
			tariff /= 2
		}
	}

	buyer.qty -= qty
	seller.qty -= qty

	assert(buyer.qty >= 0)
	assert(seller.qty >= 0)

	if kind(buyer.who) == T_loc {
		gen_item(seller.who, item_gold, cost-tariff)
		consume_item(seller.who, item, qty)

		if item == item_pot || item == item_basket {
			gold_pot_basket += (cost - tariff)
		} else if item == item_fish {
			gold_fish += (cost - tariff)
		} else if item == item_opium {
			treasury := cost
			log_output(LOG_SPECIAL, "%s earned %s selling opium.",
				box_name(seller.who), gold_s(cost-tariff))
			gold_opium += cost - tariff
			/*
			 *  Sun Dec  1 20:01:12 1996 -- Scott Turner
			 *
			 *  Opium kills peasants.
			 *
			 *  Sat Apr 18 17:46:14 1998 -- Scott Turner
			 *
			 *  More complex: pulls out of province gold until
			 *  that runs out, then kills off peasants.
			 */
			if treasury > has_item(province(buyer.who), item_gold) {
				treasury = has_item(province(buyer.who), item_gold)
			}
			if cost-treasury > 0 {
				// todo: or_int is naive and won't prevent both branches from executing
				if !consume_item(province(buyer.who), item_peasant, or_int((has_item(province(buyer.who), item_peasant) > qty), qty, has_item(province(buyer.who), item_peasant))) {
					panic("assert(consume_item(...) != false)")
				}
			}
		} else {
			gold_trade += cost - tariff
		}
	} else if kind(seller.who) == T_loc {
		/*
		 *  Sun Oct 11 18:42:53 1998 -- Scott Turner
		 *
		 *  Bargaining artifact comes into play only here.
		 *
		 */
		if art = has_artifact(buyer.who, ART_BARGAIN, 0, 0, 0); art != FALSE {
			discount := (cost * rp_item_artifact(art).param1) / 100
			wout(buyer.who,
				"The merchant admires your artifact and gives you a %d discount.",
				discount)
			cost -= discount
		}
		consume_item(buyer.who, item_gold, cost)

		if item_unique(item) != FALSE {
			move_item(seller.who, buyer.who, item, qty)
		} else {
			gen_item(buyer.who, item, qty)
		}
	} else {
		gen_item(seller.who, item_gold, cost-tariff)
		consume_item(buyer.who, item_gold, cost)
		move_item(seller.who, buyer.who, item, qty)
	}

	if seller.cloak != FALSE {
		seller_s = ""
	} else {
		seller_s = sout(" from %s", box_name(seller.who))
	}

	if buyer.cloak != FALSE {
		buyer_s = ""
	} else {
		buyer_s = sout(" to %s", box_name(buyer.who))
	}

	if kind(buyer.who) != T_loc {
		wout(buyer.who, "Bought %s%s for %s.",
			box_name_qty(item, qty),
			seller_s,
			gold_s(cost))
	}

	if kind(seller.who) != T_loc {
		if tariff != FALSE {
			wout(seller.who, "Sold %s%s for %s each, total of %s (%s minus %s tariff)",
				box_name_qty(item, qty),
				buyer_s,
				gold_s(item_cost),
				gold_s(cost-tariff),
				gold_s(cost),
				gold_s(tariff))
		} else {
			wout(seller.who, "Sold %s%s for %s each, total of %s.",
				box_name_qty(item, qty),
				buyer_s,
				gold_s(item_cost),
				gold_s(cost-tariff))
		}
	}

	if seller.cloak == FALSE && buyer.cloak == FALSE {
		wout(where, "%s bought %s from %s for %s.",
			box_name(buyer.who),
			box_name_qty(item, qty),
			box_name(seller.who),
			gold_s(cost))
	} else if seller.cloak != FALSE {
		wout(where, "%s bought %s for %s.",
			box_name(buyer.who),
			box_name_qty(item, qty),
			gold_s(cost))
	} else if buyer.cloak != FALSE {
		wout(where, "%s sold %s for %s.",
			box_name(seller.who),
			box_name_qty(item, qty),
			gold_s(cost))
	}
	/*
	 *  Where does the tariff go?
	 *
	 */
	if garrison != FALSE && tariff != FALSE {
		gen_item(garrison, item_gold, tariff)
		wout(seller.who, "Paid a market tax of %s to %s.",
			gold_s(tariff),
			box_name(where))
		wout(garrison, "Collected %s in tariffs.",
			gold_s(tariff))
		gold_tariffs += tariff
	}

	return TRUE

}

/*
 *  Tue Jul 21 09:51:48 1998 -- Scott R. Turner
 *
 *  A new trading philosophy -- the auction.
 *
 *  (1) Sort all from low to high.
 *  (2) Iterate through the buyers:
 *      (a) Highest purchases however many he wants, paying
 *          $1 more than the next highest buyer, or the asking
 *          price.
 *      (b) qtys for buyer & seller are reduced.
 *  (3) Eliminate trades for buyers & sellers reduced to zero.
 *
 */

/*
 *  We can still use "attempt_trade" or rather the NEW version
 *  of it (which takes a "cost" parameter)
 *
 *
 *  See if we can get away with only using resolve_trades_here...
 */

/*
 *  Tue Jul 21 13:05:16 1998 --
 *
 *  Make sure someone is not making an unbacked "BUY" to simply drive up
 *  the market cost.
 *
 */
func legitimate_offer(t *trade) int {
	var has int

	if t.kind == BUY {
		has = has_item(t.who, item_gold) - t.have_left
		if has < 0 {
			has = 0
		}
		return min(t.qty, has/t.cost)
	}

	if t.kind == SELL {
		has = has_item(t.who, t.item) - t.have_left
		if has < 0 {
			has = 0
		}
		return min(t.qty, has)
	}

	panic("!reached")
}

func resolve_trades_here(where int) {
	var sellers []*trade
	var buyers []*trade
	var cur_buyer, cur_seller, next_buyer, cost, success int

	/*
	 *  Do nothing if this is not a market.
	 *
	 */
	if market_here(where) == FALSE {
		return
	}

	sellers = seller_list(where, 0)
	buyers = buyer_list(where, 0)

	/*
	 *  Note that we need to revisit the entire buyers list each time a
	 *  sale occurs, because that sale might lower the auction price and
	 *  permit a buyer with a marginal amount of money to make a buy he
	 *  couldn't afford before.
	 *
	 */
	cur_buyer = 0
	for cur_buyer < len(buyers) {
		/*
		 *  Success tracks whether we actually made a trade.
		 *
		 */
		success = FALSE
		/*
		 *  Maybe he's already satisfied?  In that case, go on
		 *  to the next buyer.
		 *
		 */
		if buyers[cur_buyer].qty < 1 {
			cur_buyer++
			continue
		}
		/*
		 *  See if there is a seller.  Note we can stop looking when this
		 *  buyer runs out of stuff to buy.
		 *
		 */
		for cur_seller = 0; cur_seller < len(sellers) && buyers[cur_buyer].qty > 0; cur_seller++ {
			if sellers[cur_seller].who != buyers[cur_buyer].who &&
				sellers[cur_seller].item == buyers[cur_buyer].item &&
				sellers[cur_seller].qty > 0 {
				/*
				 *  First see if anyone else is buying this item to establish
				 *  a "cost".
				 *
				 */
				cost = 0
				for next_buyer = cur_buyer + 1; next_buyer < len(buyers); next_buyer++ {
					if buyers[next_buyer].who != sellers[cur_seller].who &&
						buyers[next_buyer].item == buyers[cur_buyer].item &&
						buyers[next_buyer].qty > 0 &&
						buyers[next_buyer].cost > cost &&
						legitimate_offer(buyers[next_buyer]) != FALSE {
						/*
						 *  This is the best cost offered by any other buyer,
						 *  so keep it.
						 */
						cost = buyers[next_buyer].cost
					}
				}
				/*
				 *  If we got a "next best bid", then mark that cost up by
				 *  one (the minimum).
				 *
				 */
				if cost != FALSE {
					cost++
				}
				/*
				 *  And now attempt the trade.  If it is successful, note that
				 *  because it means we'll have to try all the buyers again.
				 *
				 */
				if attempt_trade(buyers[cur_buyer], sellers[cur_seller], cost) != FALSE {
					success = TRUE
				}
			}
		}
		/*
		 *  If we have a sale, peruse the list again.
		 *
		 */
		if success != FALSE {
			cur_buyer = 0
		} else {
			cur_buyer++
		}
	}
}

/*
 *  Note move.c needs to be fixed for efficiency once this works.
 *
 */
/*
void
match_trades(who int)
{
resolve_trades_here(subloc(who));
};
*/

func match_all_trades() {
	var where int

	for _, where = range loop_loc() {
		resolve_trades_here(where)
	}

}

//#ifndef NEW_TRADE
//static void
//scan_trades(struct trade *t, struct trade **l)
//{
//  var i int
//
//  for i =  0; i < len(l) && t.qty > 0; i++
//    {
//      if (l[i].item != t.item)
//    continue;
//
//      if (l[i].who == t.who)
//    continue;	/* don't trade with ourself */
//
//      if (t.kind == BUY)
//    attempt_trade(t, l[i]);
//      else if (t.kind == SELL)
//    attempt_trade(l[i], t);
//    }
//}
//
//
//void
//match_trades(who int)
//{
//    var t *trade
//    where := subloc(who);
//    int first_buy = TRUE;
//    int first_sell = TRUE;
//    var sellers []*trade
//    var buyers []*trade
//
//    if (!market_here(who))
//        return;
//
//    for _, t = range loop_trade(who, t)
//    {
//        assert(t.who == who);
//
//        if (t.kind == BUY)
//        {
//            if (first_buy)
//            {
//                first_buy = FALSE;
//                sellers = seller_list(where, who);
//            }
//
//            scan_trades(t, sellers);
//        }
//        else if (t.kind == SELL)
//        {
//            if (first_sell)
//            {
//                first_sell = FALSE;
//                buyers = buyer_list(where, who);
//            }
//
//            scan_trades(t, buyers);
//        }
//    }
//
//}
//
//void
//match_all_trades()
//{
//    var where int
//    var sellers []*trade
//    var buyers []*trade
//    var i int
//
//    for _, where = range loop_loc(where)
//    {
//        if (!market_here(where))
//            continue;
//
//        sellers = seller_list(where, 0);
//        buyers = buyer_list(where, 0);
//
//        if (len(buyers) <= 0 || len(sellers) <= 0)
//            continue;
//
//        for i =  0; i < len(buyers); i++
//            scan_trades(buyers[i], sellers);
//    }
//
//}
//#endif // NEW_TRADE

//#ifndef NEW_TRADE
//ilist trades_to_check = nil;
//
//
///*
// *  This may also become hideously inefficient with NEW, since it may
// *  check the same loc over and over.
// *
// */
//void
//check_validated_trades()
//{
//    var i int
//
//    for i =  0; i < len(trades_to_check); i++
//    {
//        match_trades(trades_to_check[i]);
//    }
//
//    ilist_clear(&trades_to_check);
//}
//
///*
// *  Who has been given some item.  See if this validates
// *  any pending trades.
// *
// *  Our strategy is to see if any pending trades weren't already
// *  active at the old quantity of the item.  If so, we'll attempt
// *  a match soon.
// *
// *  We don't fire the trade inside of add_item, since it's too
// *  dangerous.  We want the command to complete, and be able to
// *  assert that a unit actually has an item after add_item has
// *  been called.
// */
//
//void
//investigate_possible_trade(int who, int item, int old_has)
//{
//    var t *trade
//    int check = FALSE;
//
//    if (item == item_gold)
//    {
//        for _, t = range loop_trade(who, t)
//        {
//            if (t.kind != BUY)
//                continue;
//
//            if ((old_has - t.have_left) / t.cost < t.qty)
//            {
//                check = TRUE;
//                break;
//            }
//        }
//
//    }
//    else
//    {
//        for _, t = range loop_trade(who, t)
//        {
//#if 0
//            if (t.kind != SELL)
//#else
//            if (t.kind != SELL || t.item != item)
//#endif
//                continue;
//
//            if (old_has - t.have_left < t.qty)
//            {
//                check = TRUE;
//                break;
//            }
//        }
//
//    }
//
//    if (check)
//        ilist_append(&trades_to_check, who);
//}
//#endif // NEW_TRADE

func find_trade(who, kind, item int) *trade {
	for _, t := range loop_trade(who) {
		if t.kind == kind && t.item == item {
			return t
		}
	}
	return nil
}

func new_trade(who, kind, item int) *trade {
	var ret *trade

	ret = find_trade(who, kind, item)

	if ret == nil {
		ret = &trade{}

		ret.who = who
		ret.kind = kind
		ret.item = item

		bx[who].trades = append(bx[who].trades, ret)
	}

	return ret
}

func gold_each(cost, qty int) string {

	if qty == 1 {
		return gold_s(cost)
	}

	return sout("%s each", gold_s(cost))
}

func v_buy(c *command) int {
	//where := subloc(c.who);
	item := c.a
	qty := c.b
	cost := c.c
	have_left := c.d
	hide_me := c.e
	var t *trade

	if kind(item) != T_item {
		wout(c.who, "%s is not an item.", box_code(item))
		return FALSE
	}

	if item == item_gold {
		wout(c.who, "Can't buy or sell gold.")
		return FALSE
	}

	if rp_item(item).ungiveable != FALSE {
		wout(c.who, "Cannot buy an item than cannot be transferred.")
		return FALSE
	}

	if hide_me != FALSE {
		if has_skill(c.who, sk_cloak_trade) != FALSE {
			hide_me = 1
		} else {
			wout(c.who, "Must have %s to conceal trades.",
				box_code_less(sk_cloak_trade))
			return FALSE
		}
	}

	if cost < 1 {
		cost = 1
	}

	t = new_trade(c.who, BUY, item)
	assert(t.who == c.who)

	if qty <= 0 {
		if t.qty <= 0 {
			wout(c.who, "No pending buy for %s.", box_name(item))
		} else {
			wout(c.who, "Cleared pending buy for %s.", box_name(item))
		}
	}

	t.qty = qty
	t.cost = cost
	t.cloak = hide_me
	t.have_left = have_left

	if qty > 0 {
		wout(c.who, "Try to buy %s for %s.",
			box_name_qty(item, qty),
			gold_each(cost, qty))

		//#ifndef NEW_TRADE
		//        if (market_here(c.who))
		//            scan_trades(t, seller_list(where, c.who));
		//#endif
	}

	return TRUE
}

func v_sell(c *command) int {
	//where := subloc(c.who);
	item := c.a
	qty := c.b
	cost := c.c
	have_left := c.d
	hide_me := c.e
	var t *trade

	if kind(item) != T_item {
		wout(c.who, "%s is not an item.", box_code(item))
		return FALSE
	}

	if item == item_gold {
		wout(c.who, "Can't buy or sell gold.")
		return FALSE
	}

	if rp_item(item).ungiveable != FALSE {
		wout(c.who, "Cannot sell an item than cannot be transferred.")
		return FALSE
	}

	if hide_me != FALSE {
		if has_skill(c.who, sk_cloak_trade) != FALSE {
			hide_me = 1
		} else {
			wout(c.who, "Must have %s to conceal trades.",
				box_code_less(sk_cloak_trade))
			return FALSE
		}
	}

	if cost < 1 {
		cost = 1
	}

	t = new_trade(c.who, SELL, item)
	assert(t.who == c.who)

	if qty <= 0 {
		if t.qty <= 0 {
			wout(c.who, "No pending sell for %s.", box_name(item))
		} else {
			wout(c.who, "Cleared pending sell for %s.", box_name(item))
		}
	}

	t.qty = qty
	t.cost = cost
	t.cloak = hide_me
	t.have_left = have_left

	if qty > 0 {
		wout(c.who, "Try to sell %s for %s.",
			box_name_qty(item, qty),
			gold_each(cost, qty))

		//#ifndef NEW_TRADE
		//        if (market_here(c.who))
		//            scan_trades(t, buyer_list(where, c.who));
		//#endif // NEW_TRADE
	}

	return TRUE
}

func list_market_items(who int, l []*trade, first int) int {
	var i int
	var qty int

	for i = 0; i < len(l); i++ {
		if l[i].cloak >= 2 {
			continue
		}

		qty = reduce_qty(l[i], l[i].cost)

		if qty <= 0 {
			continue
		}

		if first != FALSE {
			out(who, "")
			out(who, "%5s %*s %7s %6s %9s   %-25s",
				"trade", CHAR_FIELD, "who", "price",
				"qty", "wt/ea", "item")
			out(who, "%5s %*s %7s %6s %9s   %-25s",
				"-----", CHAR_FIELD, "---", "-----", "---",
				"-----", "----")

			first = FALSE
		}

		tagout(who, "<tag type=trade type=%d who=%d cost=%d qty=%d item=%d>",
			l[i].kind,
			or_int(l[i].cloak != FALSE, 0, l[i].who),
			l[i].cost,
			qty,
			l[i].item)
		out(who, "%5s %*s %7s %6s %9s   %-25s",
			or_string(l[i].kind == BUY, "buy", "sell"),
			CHAR_FIELD,
			or_string(l[i].cloak != FALSE, "?", box_code_less(l[i].who)),
			comma_num(l[i].cost),
			comma_num(qty),
			comma_num(item_weight(l[i].item)),
			plural_item_box(l[i].item, qty))
		tagout(who, "</tag type=trade>")
	}

	return first
}

func market_report(who, where int) {
	var l []*trade
	first := TRUE
	taxes := ""

	if rp_subloc(where) != nil && rp_subloc(where).tax_market != FALSE {
		taxes = sout(" (tax of %d%%)", rp_subloc(where).tax_market)
	}

	tagout(who, "<tag type=market_report loc=%d tax=%d>",
		where, or_int(rp_subloc(where) != nil, rp_subloc(where).tax_market, 0))
	out(who, "")
	out(who, "Market report%s:", taxes)
	indent += 3

	var t *trade
	flag := true

	for _, t = range loop_trade(where) {
		if t.kind == PRODUCE && t.month_prod != FALSE {
			if flag {
				out(who, "")
				flag = false
			}

			tagout(who, "<tag type=produces loc=%d item=%d month=%d>",
				where, t.item, t.month_prod)
			wout(who, "%s produces %s on month %d.",
				just_name(where),
				plural_item_name(t.item, 2),
				t.month_prod)
		}
	}

	l = buyer_list(where, 0)

	if len(l) > 0 {
		first = list_market_items(who, l, first)
	}

	l = seller_list(where, 0)

	if len(l) > 0 {
		first = list_market_items(who, l, first)
	}

	if first != FALSE {
		out(who, "No goods offered for trade.")
	}

	indent -= 3

	tagout(who, "</tag type=market_report>")
}

func list_pending_trades(who, num int) {
	first := TRUE
	var t *trade

	for _, t = range loop_trade(num) {
		if t.kind != BUY && t.kind != SELL {
			continue
		}

		if first != FALSE {
			tagout(who, "<tag type=pending_trades id=%d>", who)

			out(who, "")
			out(who, "Pending trades:")
			out(who, "")
			indent += 3
			first = FALSE

			out(who, "%5s  %7s  %5s   %s",
				"trade", "price", "qty", "item")
			out(who, "%5s  %7s  %5s   %s",
				"-----", "-----", "---", "----")
		}

		tagout(who, "<tag type=trade id=%d kind=%d cost=%d num=%d item=%d>",
			who,
			t.kind,
			(t.cost),
			(t.qty),
			(t.item))

		out(who, "%5s  %7s  %5s   %s",
			or_string(t.kind == BUY, "buy", "sell"),
			comma_num(t.cost),
			comma_num(t.qty),
			box_name(t.item))
		tagout(who, "</tag type=trade>")
	}

	if first == FALSE {
		indent -= 3
		tagout(who, "</tag type=pending_trades>")
	}

}

func add_city_trade(where, kind, item, qty, cost, month int) {
	var t *trade

	t = new_trade(where, kind, item)
	t.qty = qty
	t.cost = cost
	t.month_prod = month
}

/*
 *  Wed Nov 10 06:33:48 1999 -- Scott Turner
 *
 *  If a trade does not exist, it is created.
 *
 *  Otherwise the trade's qty is updated.
 *
 *  Tue Jan 18 08:28:45 2000 -- Scott Turner
 *
 *  Prices are not getting reset for common trade items.
 *
 *  Tue Jan 18 12:22:20 2000 -- Scott Turner
 *
 *  With update_markets, we should have this just add the big city
 *  trades if they don't already exist.
 *
 *
 */
func update_city_trade(where, kind, item, qty, cost, month int) {
	if find_trade(where, kind, item) != nil {
		return
	}

	/*
	 *  Otherwise the trade wasn't found.
	 *
	 */
	add_city_trade(where, kind, item, qty, cost, month)
}

/*
 *  Wed Sep 15 17:48:27 1999 -- Scott Turner
 *
 *  Delete all trades having to do with an item.
 *
 */
func delete_city_trade(where, item int) {
	var e *trade
	for _, e = range loop_trade(where) {
		if e.item == item {
			bx[where].trades = bx[where].trades.rem_value(e)
		}
	}

}

/*
 *  Opium model
 *
 *	Every city has a status indicating its level of opium economic
 *	development, i.e. how addicted is the local populace.
 *	This status is 1-8.  Any sale maintains, saturation causes
 *	rise, no sale decay.
 *
 *	level	profit	qty	price
 *	-----	------	---	-----
 *	  8	 800	 80	 10
 *	  7	 700	 70	 10
 *	  6	 600	 66	  9
 *	  5	 500	 55	  9
 *	  4	 400	 50	  8
 *	  3	 300	 37	  8
 *	  2	 200	 28	  7
 *	  1	 100	 15	  7
 */

type opium_t struct {
	qty, cost int
}

var opium_data = []opium_t{
	{15, 17},
	{28, 17},
	{37, 18},
	{50, 18},
	{55, 19},
	{66, 19},
	{70, 20},
	{80, 20},
}

const MAX_OPIUM_ECON = 7

func opium_market_delta(where int) {
	var p *entity_subloc
	var t *trade

	assert(subkind(where) == sub_city)

	t = find_trade(where, BUY, item_opium)
	p = p_subloc(where)

	if t != nil {
		if t.qty < 1 /* sold everything */ {
			p.opium_econ++
		} else if t.qty == opium_data[p.opium_econ].qty {
			p.opium_econ-- /* sold none */
		}
	}

	if p.opium_econ > MAX_OPIUM_ECON {
		p.opium_econ = MAX_OPIUM_ECON
	}
	if p.opium_econ < 0 {
		p.opium_econ = 0
	}

	t = new_trade(where, CONSUME, item_opium)

	t.qty = opium_data[p.opium_econ].qty
	t.cost = opium_data[p.opium_econ].cost

	if p.opium_econ > 0 {
		t.cloak = 1
	} else {
		t.cloak = 2
	}
}

func trade_suffuse_ring(where int) {
	var t *trade
	found := FALSE
	var item int

	for _, t = range loop_trade(where) {
		if subkind(t.item) == sub_magic_artifact &&
			rp_item_artifact(t.item).type_ == ART_DESTROY &&
			t.kind == SELL && t.qty > 0 {
			found = TRUE
		}
	}

	if found != FALSE || rnd(1, 3) < 3 {
		return
	}

	item = new_suffuse_ring(where)

	newTrade := new_trade(where, SELL, item)

	newTrade.qty = 1
	newTrade.cost = 450 + rnd(0, 12)*50
	newTrade.cloak = FALSE
}

/*
 *  Tue Aug 20 16:12:17 1996 -- Scott Turner
 *
 *  With a small probability, generate a scroll that teaches a
 *  religion in a city.  This permits "bootstrapping" clerics.
 *
 *  Mon Nov  4 15:32:54 1996 -- Scott Turner
 *
 *  Modified to add scrolls & books to a city.  Only add them if
 *  they're aren't some there already.  Even then only add them 50% of
 *  the time.
 *
 *  Sat Nov  9 17:04:17 1996 -- Scott Turner
 *
 *  No more scrolls.
 *
 *  Tue Apr 21 19:30:10 1998 -- Scott Turner
 *
 *  Don't actually put the scrolls in the city -- that way they
 *  can't be "explored".
 *
 */
func add_scrolls(where int) {
	found := 0
	for _, t := range loop_trade(where) {
		if rp_item_magic(t.item) != nil && len(rp_item_magic(t.item).may_study) != 0 {
			found++
		}
	}

	if found >= options.num_books {
		return
	}

	/*
	 *  Otherwise, possibly create one.
	 *
	 */
	if rnd(1, 100) > SCROLL_CHANCE {
		return
	}

	/*
	 *  Okay, make a random book/scroll
	 *
	 */
	newBook := make_teach_book(where, 0, 1, sub_book)

	t := new_trade(where, SELL, newBook)
	t.qty = 1
	t.cost = 100 + p_item_magic(newBook).orb_use_count*25 + rnd(0, 6)*50
	p_item(newBook).base_price = t.cost
	t.cloak = FALSE
	t.counter = 0
	wout(gm_player, "Created new teaching book %s at %s.",
		box_name(t.item), box_name(where))

}

//#if 0
///*
// *  NOw called in update_markets
// *
// */
//void
//location_trades()
//{
//  var where int
//
//  loop_city(where) {
//    opium_market_delta(where);
//    loc_trade_sup(where, FALSE);
//
//    if (in_faery(where)) {
//      trade_suffuse_ring(where);
//    };
//
//    if (in_clouds(where)) {
//      trade_suffuse_ring(where);
//    };
//
//    if (!in_faery(where) &&
//    !in_clouds(where) &&
//    !in_hades(where))
//      add_scrolls(where);
//
//  } next_city;
//}
//#endif
