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

// mining-related skills and functions.

const (
	MINE_PRODUCTS = 4
)

type mine_product struct {
	mins [MINE_PRODUCTS]int
	maxs [MINE_PRODUCTS]int
}

var (
	mine_products = [MINE_PRODUCTS]int{item_iron, item_gold, item_mithril, item_gate_crystal}
	mine_qties    = [MINE_MAX + 1]mine_product{
		/*        iron  gold  mithril gate_crystal  */
		/*        ----  ----  ------- ------------  */
		/* 0 */ {[MINE_PRODUCTS]int{0, 0, 0, 0}, [MINE_PRODUCTS]int{0, 0, 0, 0}},
		/* 1 */ {[MINE_PRODUCTS]int{0, 0, 0, 0}, [MINE_PRODUCTS]int{10, 0, 0, 0}},
		/* 2 */ {[MINE_PRODUCTS]int{10, 0, 0, 0}, [MINE_PRODUCTS]int{20, 0, 0, 0}},
		/* 3 */ {[MINE_PRODUCTS]int{10, 0, 0, 0}, [MINE_PRODUCTS]int{30, 0, 0, 0}},
		/* 4 */ {[MINE_PRODUCTS]int{10, 0, 0, 0}, [MINE_PRODUCTS]int{40, 0, 0, 0}},
		/* 5 */ {[MINE_PRODUCTS]int{10, 0, 0, 0}, [MINE_PRODUCTS]int{30, 100, 0, 0}},
		/* 6 */ {[MINE_PRODUCTS]int{0, 0, 0, 0}, [MINE_PRODUCTS]int{20, 300, 0, 0}},
		/* 7 */ {[MINE_PRODUCTS]int{0, 100, 0, 0}, [MINE_PRODUCTS]int{20, 450, 0, 0}},
		/* 8 */ {[MINE_PRODUCTS]int{0, 120, 0, 0}, [MINE_PRODUCTS]int{20, 600, 0, 0}},
		/* 9 */ {[MINE_PRODUCTS]int{0, 140, 0, 0}, [MINE_PRODUCTS]int{20, 750, 0, 0}},
		/*10 */ {[MINE_PRODUCTS]int{0, 160, 0, 0}, [MINE_PRODUCTS]int{20, 900, 1, 0}},
		/*11 */ {[MINE_PRODUCTS]int{0, 180, 0, 0}, [MINE_PRODUCTS]int{0, 1000, 3, 0}},
		/*12 */ {[MINE_PRODUCTS]int{0, 200, 0, 0}, [MINE_PRODUCTS]int{0, 1200, 6, 0}},
		/*13 */ {[MINE_PRODUCTS]int{0, 220, 0, 0}, [MINE_PRODUCTS]int{0, 1300, 9, 0}},
		/*14 */ {[MINE_PRODUCTS]int{0, 240, 0, 0}, [MINE_PRODUCTS]int{0, 1500, 12, 0}},
		/*15 */ {[MINE_PRODUCTS]int{0, 200, 0, 0}, [MINE_PRODUCTS]int{0, 1200, 15, 1}},
		/*16 */ {[MINE_PRODUCTS]int{0, 160, 0, 0}, [MINE_PRODUCTS]int{0, 900, 10, 1}},
		/*17 */ {[MINE_PRODUCTS]int{0, 120, 0, 0}, [MINE_PRODUCTS]int{0, 600, 5, 1}},
		/*18 */ {[MINE_PRODUCTS]int{0, 80, 0, 0}, [MINE_PRODUCTS]int{0, 300, 0, 1}},
		/*19 */ {[MINE_PRODUCTS]int{0, 20, 0, 1}, [MINE_PRODUCTS]int{0, 100, 0, 1}},
		/*20 */ {[MINE_PRODUCTS]int{0, 0, 0, 1}, [MINE_PRODUCTS]int{0, 100, 0, 3}},
	}
)

/*
 *  Mine_Chance
 *  Fri Jan 24 12:51:08 1997 -- Scott Turner
 *
 *  This function returns the % chance per day that a particular miner finds a particular item.
 *  The chance depends upon:
 *   (1) the miner's experience (skill_exp == months practiced)
 *   (2) the item being sought
 *   (3) the number of workers
 *
 *  According to the following formula:
 *   (1) Initial chance from item.
 *   (2) +5.0% for each month experience (max +50%) // was +2.0%
 *   (3) +0.2% for each worker           (max +20%) // was +0.1%
 *
 *  Per week of mining, so that we divide by 7 at the end.
 */
func mine_chance(item, who, skill int) float64 {
	var initial_chance float64
	switch item {
	case item_iron:
		initial_chance = 70.0
	case item_gold:
		initial_chance = 50.0
	case item_mithril:
		initial_chance = 30.0
	case item_gate_crystal:
		initial_chance = 10.0
	default:
		return 0
	}
	// add experience bonus
	exp_bonus := float64(skill_exp(who, skill)) * 5.0
	if exp_bonus > 50 {
		exp_bonus = 50.0
	}
	// worker bonus
	worker_bonus := float64(effective_workers(who)) * 0.20
	if worker_bonus > 20.0 {
		worker_bonus = 20.0
	}

	// add it up and return the daily %.
	return (initial_chance + exp_bonus + worker_bonus) / 7.0
}

/*
 *  Mon Jan 20 11:54:43 1997 -- Scott Turner
 *
 *  Mine depth is now how deep underground your mine shaft location is (nested).
 *
 *  Wed May  7 12:22:12 1997 -- Scott Turner
 *
 *  Mine shafts are no longer nested -they are up/down.
 *  So you want to follow the "up link" to determine how deep you are.
 *  Ugh.
 */
func mine_depth(where int) int {
	depth := -1
	for where > 0 && subkind(where) == sub_mine_shaft {
		where = location_direction(where, DIR_UP)
		depth++
	}
	return depth
}

// Get the mine that heads up this mine shaft.
// Scott Turner - Ugghhhhhhhh. This is considerably uglier since
// mine shafts are no longer sublocs. we need to go "up" until we
// find a mine_info.
func get_mine_info(where int) *entity_mine {
	if subkind(where) != sub_mine_shaft {
		return nil
	}
	for where > 0 && subkind(where) == sub_mine_shaft {
		if where = location_direction(where, DIR_UP); where == 0 {
			panic("assert(where)")
		}
	}
	if p_loc(where).mine_info == nil {
		create_mine_info(where)
	}
	return p_loc(where).mine_info
}

/*
 *  Fri Jan 31 12:47:13 1997 -- Scott Turner
 *
 *  Create mine info when first accessed.
 *
 */
func create_mine_info(mine int) {
	if p_loc(mine).mine_info != nil {
		panic("assert(!p_loc(mine).mine_info)")
	}
	// allocate the entity_mine
	p_loc(mine).mine_info = &entity_mine{}
	// go through each level of the mine and add in the goods for that level...
	for i := 0; i < MINE_MAX; i++ {
		el := p_loc(mine).mine_info.mc[i].items
		p_loc(mine).mine_info.shoring[i] = NO_SHORING
		for j := 0; j < MINE_PRODUCTS; j++ {
			qty := rnd(mine_qties[i].mins[j], mine_qties[i].maxs[j])
			if qty != 0 {
				el = append(el, &item_ent{item: mine_products[j], qty: qty})
			}
		}
	}
}

/*
 *  Sat Jan 25 11:30:04 1997 -- Scott Turner
 *
 *  Return the amount of some item at some mine level.
 *
 */
func mine_has_item(mine, depth, item int) int {
	mi := get_mine_info(mine)

	if mi == nil {
		return 0
	}

	assert(depth >= 0)
	if depth >= MINE_MAX {
		depth = MINE_MAX - 1
	}

	for i := 0; i < len(mi.mc[depth].items); i++ {
		if mi.mc[depth].items[i].item == item {
			return mi.mc[depth].items[i].qty
		}
	}

	return 0
}

/*
 *  Sat Jan 25 11:30:04 1997 -- Scott Turner
 *
 *  Decrement the amount of some item at some mine level.
 *
 */
func mine_sub_item(mine, depth, item, amount int) {
	mi := get_mine_info(mine)
	if mi == nil {
		return
	}

	if depth < 0 {
		panic("assert(depth >= 0)")
	}
	if depth >= MINE_MAX {
		depth = MINE_MAX - 1
	}

	for i := 0; i < len(mi.mc[depth].items); i++ {
		if mi.mc[depth].items[i].item == item {
			mi.mc[depth].items[i].qty -= amount
			return
		}
	}

	return
}

func start_generic_mine(c *command, item, skill int) bool {
	where := subloc(c.who)
	days := c.a
	if days == 0 {
		days = 7
	}

	if subkind(where) != sub_mine_shaft {
		wout(c.who, "Must be in a mine shaft to extract %s.", just_name(item))
		return false
	}

	nworkers := effective_workers(c.who)
	if nworkers < 10 {
		wout(c.who, "Mining requires at least ten workers.")
		return false
	}

	c.wait = days

	wout(c.who, "Will mine %s for the next %s days.", just_name(item), nice_num(c.wait))

	return true
}

func finish_generic_mine(c *command, item, skill int) bool {
	where := subloc(c.who)
	if subkind(where) != sub_mine_shaft {
		wout(c.who, "%s is no longer in a mine shaft.", box_name(c.who))
		return false
	}

	nworkers := effective_workers(c.who)
	if nworkers < 10 {
		wout(c.who, "%s no longer has ten workers.", box_name(c.who))
		return false
	}

	/*
	 *  Calculate the chance that we found some of whatever we're
	 *  looking for...
	 *
	 */
	chance := int(mine_chance(item, c.who, skill) * float64(c.days_executing))

	if chance < rnd(1, 100) {
		wout(c.who, "Mining yielded no %s.", just_name(item))
		return false
	}

	/*
	 *  Now see if there's any of that left here at this level.
	 *
	 */
	depth := mine_depth(where)
	has := mine_has_item(where, depth, item)
	if has == FALSE {
		wout(c.who, "Mining yielded no %s.", just_name(item))
		if rnd(1, 100) < 70 {
			wout(c.who, "This level appears to be mined out.")
		}
		return false
	}

	/*
	 *  How much of this do we actually get?  We can recover from
	 *  40-60% of this, depending upon how many men we have as workers.
	 *
	 */
	if nworkers > 100 {
		nworkers = 100
	}
	qty := (has * rnd(40, 40+(nworkers/5))) / 100

	/*
	 *  Guarantee at least 1 for gate crystals, mithril, etc.
	 *
	 */
	if qty == 0 {
		qty = 1
	}

	/*
	 *  Remove that much from the mine; any bonuses don't affect
	 *  the actual amount in the mine.
	 *
	 */
	mine_sub_item(where, depth, item, qty)

	/*
	 *  This mine might be blessed.
	 *
	 */
	if get_effect(where, ef_improve_mine, 0, 0) != 0 {
		qty += int(float64(qty)*0.50 + 0.50)
		wout(c.who, "%s is unusually productive.", box_name(where))
	}

	/*
	 *  We might have hit a rich lode.
	 *
	 */
	lode := rnd(1, 5000)
	if lode < c.days_executing {
		wout(c.who, "You hit an incredibly rich lode of %s!", just_name(item))
		qty *= 100
	} else if lode < c.days_executing*10 {
		wout(c.who, "You hit a very rich lode of %s!", just_name(item))
		qty *= 10
	} else if lode < c.days_executing*100 {
		wout(c.who, "You hit a rich lode of %s!", just_name(item))
		qty *= 2
	}

	gen_item(c.who, item, qty)

	wout(c.who, "Mining yielded %s.", box_name_qty(item, qty))
	return true
}

func v_mine_iron(c *command) int {
	if start_generic_mine(c, item_iron, sk_mine_iron) {
		return TRUE
	}
	return FALSE
}

func d_mine_iron(c *command) int {
	if finish_generic_mine(c, item_iron, sk_mine_iron) {
		return TRUE
	}
	return FALSE
}

func v_mine_gold(c *command) int {
	if start_generic_mine(c, item_gold, sk_mine_gold) {
		return TRUE
	}
	return FALSE
}

func d_mine_gold(c *command) int {
	if finish_generic_mine(c, item_gold, sk_mine_gold) {
		return TRUE
	}
	return FALSE
}

func v_mine_mithril(c *command) int {
	if start_generic_mine(c, item_mithril, sk_mine_mithril) {
		return TRUE
	}
	return FALSE
}

func d_mine_mithril(c *command) int {
	if finish_generic_mine(c, item_mithril, sk_mine_mithril) {
		return TRUE
	}
	return FALSE
}

func v_mine_gate_crystal(c *command) int {
	if start_generic_mine(c, item_gate_crystal, sk_mine_crystal) {
		return TRUE
	}
	return FALSE
}

func d_mine_gate_crystal(c *command) int {
	if finish_generic_mine(c, item_gate_crystal, sk_mine_crystal) {
		return TRUE
	}
	return FALSE
}

/*
 *  Add_Wooden_Shoring
 *  Tue Feb  4 08:45:40 1997 -- Scott Turner
 *
 *  Modify the appropriate level of a mine shaft with wooden shoring.
 *
 */
func d_add_wooden_shoring(c *command) int {
	where := subloc(c.who)
	if subkind(where) != sub_mine_shaft {
		wout(c.who, "You must be in a mine shaft to add wooden shoring.")
		return FALSE
	}
	/*
	 *  Maybe already shored?
	 *
	 */
	mi := get_mine_info(where)
	if mi == nil {
		panic("assert(mi)")
	}
	depth := mine_depth(where)
	if mi.shoring[depth] >= WOODEN_SHORING {
		wout(c.who, "This mine shaft already has wooden shoring.")
		return FALSE
	}
	mi.shoring[depth] = WOODEN_SHORING
	wout(c.who, "Wooden shoring added to %s.", box_name(where))
	return TRUE
}

func v_add_wooden_shoring(c *command) int {
	where := subloc(c.who)
	if subkind(where) != sub_mine_shaft {
		wout(c.who, "You must be in a mine shaft to add wooden shoring.")
		return FALSE
	}
	/*
	 *  Maybe already shored?
	 *
	 */
	mi := get_mine_info(where)
	if mi == nil {
		panic("assert(mi)")
	}
	depth := mine_depth(where)
	if mi.shoring[depth] >= WOODEN_SHORING {
		wout(c.who, "This mine shaft already has wooden shoring.")
		return FALSE
	}
	/*
	 *  Needs wood.
	 *
	 */
	if has_item(c.who, item_lumber) < 25 {
		wout(c.who, "You do not have sufficient wood to shore this shaft.")
		return FALSE
	}
	return TRUE
}

/*
 *  Add_Iron_Shoring
 *  Tue Feb  4 08:45:40 1997 -- Scott Turner
 *
 *  Modify the appropriate level of a mine shaft with iron shoring.
 *
 */
func d_add_iron_shoring(c *command) int {
	where := subloc(c.who)
	if subkind(where) != sub_mine_shaft {
		wout(c.who, "You must be in a mine shaft to add iron shoring.")
		return FALSE
	}

	/*
	 *  Maybe already shored?
	 *
	 */
	mi := get_mine_info(where)
	if mi == nil {
		panic("assert(mi)")
	}
	depth := mine_depth(where)
	if mi.shoring[depth] >= IRON_SHORING {
		wout(c.who, "This mine shaft already has iron shoring.")
		return FALSE
	}
	mi.shoring[depth] = IRON_SHORING
	wout(c.who, "Iron shoring added to %s.", box_name(where))
	return TRUE
}

func v_add_iron_shoring(c *command) int {
	where := subloc(c.who)
	if subkind(where) != sub_mine_shaft {
		wout(c.who, "You must be in a mine shaft to add iron shoring.")
		return FALSE
	}
	/*
	 *  Maybe already shored?
	 *
	 */
	mi := get_mine_info(where)
	if mi == nil {
		panic("assert(mi)")
	}
	depth := mine_depth(where)
	if mi.shoring[depth] >= IRON_SHORING {
		wout(c.who, "This mine shaft already has iron shoring.")
		return FALSE
	}
	/*
	 *  Needs iron.
	 *
	 */
	if has_item(c.who, item_iron) < 5 {
		wout(c.who, "You do not have sufficient iron to shore this shaft.")
		return FALSE
	}
	return TRUE
}
