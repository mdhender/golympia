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
	"log"
	"sort"
)

/*
 *  10%  9
 *  40%  6
 *  40%  3
 *  10%  0
 */

func choose_city_prominence(city int) int {
	var n int

	if safe_haven(city) != FALSE || major_city(city) != FALSE {
		return 3
	}

	if loc_hidden(city) != FALSE || loc_hidden(province(city)) != FALSE {
		return 0
	}

	n = rnd(1, 100)

	if n <= 10 {
		return 0
	}
	if n <= 50 {
		return 1
	}
	if n <= 90 {
		return 2
	}
	return 3
}

func add_near_city(where int, city int) {
	var p *entity_subloc

	p = p_subloc(where)

	p.near_cities = append(p.near_cities, city)
}

func prop_city_near_list(city int) {
	var prom int
	var m int
	var i int
	var n int
	var dest int
	var where int
	var l []*exit_view

	clear_temps(T_loc)

	bx[province(city)].temp = 1
	prom = choose_city_prominence(city)
	p_subloc(city).prominence = prom
	prom *= 3

	for m = 1; m < prom; m++ {
		for _, where = range loop_loc() {
			if bx[where].temp != m {
				continue
			}

			l = exits_from_loc_nsew(0, where)

			for i = 0; i < len(l); i++ {
				dest = l[i].destination

				if loc_depth(dest) != LOC_province {
					continue
				}

				if bx[dest].temp == 0 {
					bx[dest].temp = m + 1
					if n = city_here(dest); n != FALSE {
						add_near_city(n, city)
					}
				}
			}
		}

	}
}

func seed_city_near_lists() {
	var city int

	stage("INIT: seed_city_near_lists()")

	for _, city = range loop_city() {
		p_subloc(city).near_cities = nil
	}

	for _, city = range loop_city() {
		prop_city_near_list(city)
	}

}

func seed_mob_cookies() {
	var i int

	for _, i = range loop_loc() {
		if subkind(i) != sub_city && loc_depth(i) != LOC_province {
			continue
		}

		if subkind(i) == sub_ocean {
			continue
		}

		gen_item(i, item_mob_cookie, 1)
	}

}

func seed_undead_cookies() {
	var i int

	for _, i = range loop_loc() {
		if subkind(i) != sub_graveyard {
			continue
		}

		gen_item(i, item_undead_cookie, 1)
	}

}

func seed_weather_cookies() {
	var i int

	for _, i = range loop_loc() {
		switch subkind(i) {
		case sub_forest:
			gen_item(i, item_rain_cookie, 1)
			gen_item(i, item_fog_cookie, 1)
			break

		case sub_plain, sub_desert, sub_mountain:
			gen_item(i, item_wind_cookie, 1)
			break

		case sub_swamp:
			gen_item(i, item_fog_cookie, 1)
			break

		case sub_ocean:
			gen_item(i, item_fog_cookie, 1)
			gen_item(i, item_wind_cookie, 1)
			gen_item(i, item_rain_cookie, 1)
			break
		}
	}

}

func seed_cookies() {

	stage("INIT: seed_cookies()")

	seed_mob_cookies()
	seed_undead_cookies()
	seed_weather_cookies()
}

/*
 *  Thu Mar 25 18:12:03 1999 -- Scott Turner
 *
 *  todo: apparently this doesn't properly wrap at both edges.
 *
 */
func compute_dist_generic(terr int) {
	var where int
	var l []*exit_view
	var set_one int
	var i int
	var dest int
	var m int

	clear_temps(T_loc)

	for _, where = range loop_province() {
		if subkind(where) != schar(terr) {
			continue
		}

		l = exits_from_loc_nsew(0, where)

		for i = 0; i < len(l); i++ {
			if loc_depth(l[i].destination) != LOC_province {
				continue
			}

			if subkind(l[i].destination) != schar(terr) {
				bx[l[i].destination].temp = 1
			}
		}
	}

	m = 1

	for {
		set_one = FALSE

		for _, where = range loop_province() {
			if subkind(where) == schar(terr) || bx[where].temp != m {
				continue
			}

			l = exits_from_loc_nsew(0, where)

			for i = 0; i < len(l); i++ {
				dest = l[i].destination

				if loc_depth(dest) != LOC_province {
					continue
				}

				if subkind(dest) != schar(terr) && bx[dest].temp == 0 {
					bx[dest].temp = m + 1
					set_one = TRUE
				}
			}
		}

		m++
		if set_one != FALSE {
			continue
		}
		break
	}

	for _, where = range loop_province() {
		if region(where) == faery_region ||
			region(where) == hades_region ||
			region(where) == cloud_region {
			continue
		}

		if subkind(where) != schar(terr) && bx[where].temp < 1 {
			log.Printf("(2)error on %d, reg=%d\n",
				where, region(where))
		}
	}

}

/*
 *  Could be speeded up by saving the return from province_gate_here()
 *  in some temp field.  But this routine is only run once, when a new
 *  database is first read in, so it probably doesn't matter.
 */

func compute_dist_gate() {
	var where int
	var l []*exit_view
	var set_one int
	var i int
	var dest int
	var m int

	clear_temps(T_loc)

	for _, where = range loop_province() {
		if province_gate_here(where) == FALSE {
			continue
		}

		l = exits_from_loc_nsew(0, where)

		for i = 0; i < len(l); i++ {
			if loc_depth(l[i].destination) != LOC_province {
				continue
			}

			if province_gate_here(l[i].destination) == FALSE {
				bx[l[i].destination].temp = 1
			}
		}
	}

	m = 1

	for {
		set_one = FALSE

		for _, where = range loop_province() {
			if province_gate_here(where) != FALSE || bx[where].temp != m {
				continue
			}

			l = exits_from_loc_nsew(0, where)

			for i = 0; i < len(l); i++ {
				dest = l[i].destination

				if loc_depth(dest) != LOC_province {
					continue
				}

				if province_gate_here(dest) == FALSE && bx[dest].temp == 0 {
					bx[dest].temp = m + 1
					set_one = TRUE
				}
			}
		}

		m++
		if set_one != FALSE {
			continue
		}
		break
	}

	for _, where = range loop_province() {
		if region(where) == faery_region ||
			region(where) == hades_region ||
			region(where) == cloud_region {
			continue
		}

		if province_gate_here(where) == FALSE && bx[where].temp < 1 {
			log.Printf("(3)error on %d, reg=%d\n",
				where, region(where))
		}
	}

}

func compute_nearby_graves() {
	var where int
	var l []*exit_view
	set_one := TRUE
	var i int
	var dest int
	sequence := 0

	for _, i = range loop_province() {
		bx[i].temp = 0
		p_loc(i).near_grave = 0
	}

	for _, i = range loop_subkind(sub_graveyard) {
		where = province(i)

		p_loc(where).near_grave = i
		bx[where].temp = 1
	}

	for set_one != FALSE {
		set_one = FALSE
		sequence++

		for _, where = range loop_province() {
			if bx[where].temp == sequence {
				l = exits_from_loc(0, where)

				for i = 0; i < len(l); i++ {
					dest = province(l[i].destination)

					if l[i].water == FALSE && bx[dest].temp == 0 {
						bx[dest].temp = bx[where].temp + 1
						rp_loc(dest).near_grave = rp_loc(where).near_grave
						set_one = TRUE
					}
				}
			}
		}

	}

	/*
	 *  We skipped water on the first pass so that graveyards on the
	 *  same continent wouldn't short-cut across a water route.  Now
	 *  that all land graveyards should be set, do water ones.  These
	 *  should come out a straight-lines from the water location to the
	 *  nearest graveyard on nearby land.
	 */

	for _, i = range loop_province() {
		if bx[i].temp != FALSE {
			bx[i].temp = 1
		} /* reset sequence */
	}

	sequence = 0
	set_one = TRUE

	for set_one != FALSE {
		set_one = FALSE
		sequence++

		for _, where = range loop_province() {
			if bx[where].temp == sequence {
				l = exits_from_loc(0, where)

				for i = 0; i < len(l); i++ {
					dest = province(l[i].destination)

					if bx[dest].temp == 0 {
						bx[dest].temp = bx[where].temp + 1
						rp_loc(dest).near_grave = rp_loc(where).near_grave
						set_one = TRUE
					}
				}
			}
		}

	}
}

func compute_dist() {
	var i int

	stage("INIT: compute_dist()")
	compute_dist_generic(sub_ocean)

	for _, i = range loop_province() {
		p_loc(i).dist_from_sea = bx[i].temp
	}

	compute_dist_generic(sub_swamp)

	for _, i = range loop_province() {
		p_loc(i).dist_from_swamp = bx[i].temp
	}

	compute_dist_gate()

	for _, i = range loop_province() {
		p_loc(i).dist_from_gate = bx[i].temp
	}

	compute_nearby_graves()
}

//int int_comp(const void *q1, const void *q2) {
//    int *a = (int *)q1;
//    int *b = (int *)q2;
//
//    return *a - *b;
//}

/*
 *  Fri Apr 17 11:46:42 1998 -- Scott Turner
 *
 *  This should be changed for g3 -- maximum of 3 skills,
 *  and no "uncommon" skills.
 *
 */
func seed_city_skill(where int) {
	var p *entity_subloc
	num, skill, count := 0, 0, 0

	p = p_subloc(where)
	p.teaches = nil

	if in_faery(where) { /* taught only in Faery city */
		p.teaches = append(p.teaches, sk_scry)
		p.teaches = append(p.teaches, sk_artifact)
		return
	}

	if in_clouds(where) /* taught only in the Cloudlands */ {
		p.teaches = append(p.teaches, sk_weather)
		p.teaches = append(p.teaches, sk_artifact)
		return
	}

	if in_hades(where) /* taught only in Hades */ {
		p.teaches = append(p.teaches, sk_necromancy)
		p.teaches = append(p.teaches, sk_artifact)
		return
	}

	/*
	 *  Regular cities get 1-3 skills, selected from a short list.
	 *
	 */
	for num = rnd(1, 3); num > 0; num-- {
		/*
		 *  Select the skill to be used.
		 *
		 */
		for _, newSkill := range loop_skill() {
			if skill_school(newSkill) == newSkill &&
				!magic_skill(newSkill) &&
				!religion_skill(newSkill) &&
				learn_time(newSkill) < 14 {
				count++
				if rnd(1, count) == 1 {
					skill = newSkill
				}
			}
		}

		if ilist_lookup(p.teaches, skill) == -1 {
			p.teaches = append(p.teaches, skill)
		}
	}

	if len(p.teaches) > 0 {
		sort.Ints(p.teaches) // qsort(p.teaches, len(p.teaches), sizeof(int), int_comp);
	}
}

func seed_city_trade(where int) {
	prov := province(where)
	prov_kind := subkind(prov)
	p := rp_subloc(where)

	clear_all_trades(where)

	if in_hades(where) {
		return
	}

	if in_clouds(where) {
		return
	}

	if in_faery(where) /* seed Faery city trade */ {
		add_city_trade(where, PRODUCE, item_pegasus, 1, 1000, 0)

		if rnd(1, 2) == 1 {
			add_city_trade(where, PRODUCE, item_lana_bark, 3, 50, 0)
		} else {
			add_city_trade(where, PRODUCE, item_avinia_leaf, 10, 35, 0)
		}

		if rnd(1, 2) == 1 {
			add_city_trade(where, PRODUCE, item_yew, 5, 100, 0)
		} else {
			add_city_trade(where, PRODUCE, item_mallorn_wood, 5, 200, 0)
		}

		add_city_trade(where, CONSUME, item_mithril, 10, 500, 0)

		if rnd(1, 2) == 1 {
			add_city_trade(where, CONSUME, item_gate_crystal, 2, 1000, 0)
		} else {
			add_city_trade(where, CONSUME, item_jewel, 5, 100, 0)
		}

		do_production(where, TRUE)
		return
	}

	if is_port_city(where) {
		add_city_trade(where, CONSUME, item_fish, 100, 2, 0)
		if rnd(1, 2) == 1 {
			add_city_trade(where, PRODUCE, item_fish_oil, 25, 5, 0)
		} else {
			add_city_trade(where, PRODUCE, item_dried_fish, 20, 6, 0)
		}
		add_city_trade(where, PRODUCE, item_glue, 10, 50, 0)
	} else if rnd(1, 3) == 1 {
		if rnd(1, 2) == 1 {
			add_city_trade(where, CONSUME, item_fish_oil, 15,
				9+min(7, sea_dist(prov))*2, 0)
		} else {
			add_city_trade(where, CONSUME, item_dried_fish, 12,
				10+min(7, sea_dist(prov))*2, 0)
		}
	}

	if rnd(1, 2) == 1 {
		add_city_trade(where, CONSUME, item_pot, 9, 7, 0)
	} else {
		add_city_trade(where, CONSUME, item_basket, 15, 4, 0)
	}

	if prov_kind == sub_plain {
		add_city_trade(where, PRODUCE, item_ox, 5, 100, 0)
		add_city_trade(where, PRODUCE, item_riding_horse, rnd(2, 3),
			rnd(20, 30)*5, 0)
	} else if rnd(1, 3) == 1 {
		add_city_trade(where, CONSUME, item_leather, rnd(3, 6),
			rnd(125, 135), 0)
	}

	if prov_kind == sub_mountain {
		add_city_trade(where, PRODUCE, item_iron, rnd(1, 2),
			rnd(75, 200), 0)
	}

	if prov_kind == sub_forest {
		add_city_trade(where, PRODUCE, item_lumber, 25,
			rnd(11, 15), 0)
	}

	if p != nil && ilist_lookup(p.teaches, sk_alchemy) >= 0 {
		add_city_trade(where, PRODUCE, item_lead, 50, 10, 0)
	}

	do_production(where, TRUE)
}

func base_price(n int) int {
	assert(kind(n) == T_item)

	t := p_item(n)
	if t.base_price == 0 {
		t.base_price = rnd(9, 56)
	}

	return t.base_price
}

//#if 0
//
//#define	MAX_NEAR		3	/* provinces considered "near" */
//#define	LONG_ROUTE_MIN		15	/* min length of a "long" route */
//#define	LONG_ROUTE_MAX		50	/* max length of a "long" route */
//
//
//static int
//nearby_city(ilist cities, int a)
//{
//    var i int
//    int m = 999999;
//    int dist;
//    int save = 0;
//    int reg = region(a);
//
//    for i = 0; i < len(cities); i++
//    {
//        if (cities[i] == a || reg != region(cities[i]))
//            continue;
//
//        dist = los_province_distance(cities[i], a);
//
//        if (dist < m)
//        {
//            m = dist;
//            save = cities[i];
//        }
//    }
//
//    if (m > MAX_NEAR)
//        return 0;
//
//    return save;
//}
//
//
//static int
//nearby_city_two(ilist cities, int a, int b)
//{
//    var i int
//    int m = 999999;
//    int dist;
//    int save = 0;
//    int reg = region(a);
//
//    for i = 0; i < len(cities); i++
//    {
//        if (cities[i] == a ||
//            cities[i] == b ||
//            reg != region(cities[i]))
//            continue;
//
//        dist = max(los_province_distance(cities[i], a),
//                los_province_distance(cities[i], b));
//
//        if (dist < m)
//        {
//            m = dist;
//            save = cities[i];
//        }
//    }
//
//    if (m > MAX_NEAR)
//        return 0;
//
//    return save;
//}
//
//
//static void
//long_route_sup(int item, ilist source, ilist consume, int distance)
//{
//    int qty;
//    int profit;
//    int premium;
//    int m, n;
//    int q, c;
//    var i int
//
//    assert(kind(item) == T_item);
//
//    n = len(source);
//    m = len(consume);
//
//    assert(n > 0);
//    assert(m > 0);
//
//    qty = 100 + (n-1) * 25;
//
//    profit = (135 + 2 * distance) * distance;
//    premium = profit / qty;
//
//    for i = 0; i < n; i++
//    {
//        q = qty/n;
//        q += rnd(-q/10, q/10);
//
//        c = base_price(item);
//        c -= rnd(0, c/20);
//
//        add_city_trade(source[i], PRODUCE, item, q, c, rnd(1, NUM_MONTHS));
//    }
//
//    for i = 0; i < m; i++
//    {
//        q = qty/m;
//        q += rnd(-q/10, q/10);
//
//        c = base_price(item) + premium;
//        c += rnd(0, c/20);
//
//        add_city_trade(consume[i], CONSUME, item, q, c, 0);
//    }
//}
//
//
//static void
//city_cluster(ilist cities, ilist *l, int where)
//{
//    int one, two;
//
//    ilist_append(l, where);
//
//    if (rnd(1,2) == 1)
//        return;
//
//    one = nearby_city(cities, where);
//
//    if (one)
//        ilist_append(l, one);
//
//    if (rnd(1,2) == 1)
//        return;
//
//    two = nearby_city_two(cities, where, one);
//
//    if (two)
//        ilist_append(l, two);
//}
//
//
//static int
//long_route(ilist cities, int item)
//{
//    static ilist source = nil;
//    static ilist consume = nil;
//    int count = 0;
//    int dist, dist_r, dist_g;
//    int one, two;		/* source, consume cities */
//
//    ilist_clear(&source);
//    ilist_clear(&consume);
//
//    assert(len(cities) >= 2);
//
//    while (1)
//    {
//        if (count++ > 50)
//            return FALSE;
//
//        ilist_scramble(cities);
//
//        one = cities[0];
//        two = cities[1];
//
//        dist_g = distance(one, two, TRUE);
//
//        if (dist_g < LONG_ROUTE_MIN)
//            continue;
//
//        dist_r = distance(one, two, FALSE);
//
//        if (dist_r < LONG_ROUTE_MIN || dist_r > LONG_ROUTE_MAX)
//            continue;
//
//        if (dist_g > dist_r || dist_r - dist_g <= 10)
//            break;
//    }
//
//    dist = min(dist_r, dist_g);
//
//    city_cluster(cities, &source, one);
//    city_cluster(cities, &consume, two);
//
//    printf("%s from %s-%d <%s> to %s-%d <%s>, d=%d g=%d\n",
//            just_name(item),
//            just_name(region(one)),
//            len(source),
//            box_code_less(one),
//            just_name(region(two)),
//            len(consume),
//            box_code_less(two),
//            dist_r, dist_g);
//
//    long_route_sup(item, source, consume, dist);
//
//    return TRUE;
//}
//
//#endif

var rare_trade_items = []int{
	item_fine_cloak,
	item_chocolate,
	item_ivory,
	item_rug,
	item_honey,
	item_ink,
	item_licorice,
	item_soap,
	item_old_book,
	item_jade_idol,
	item_purple_cloth,
	item_rose_perfume,
	item_silk,
	item_incense,
	item_ochre,
	item_jeweled_egg,
	item_obsidian,
	item_pepper,
	item_cardamom,
	item_orange,
	item_cinnabar,
	item_myrhh,
	item_saffron,
	0}

var common_trade_items = []int{
	item_pipeweed,
	item_ale,
	item_tallow,
	item_candles,
	item_wool,
	item_vinegar,
	item_wax,
	item_sugar,
	item_salt,
	item_linen,
	item_beans,
	item_walnuts,
	item_flax,
	item_flutes,
	item_cassava,
	item_plum_wine,
	item_tea,
	0}

//#if 0
//void
//seed_long_routes()
//{
//    var i int
//    static ilist all_cities = nil;
//    var cities []int
//    static ilist regions = nil;
//    int reg = 0;
//    int reg_len;
//    int one, two;
//    var item int
//
//    ilist_clear(&all_cities);
//    ilist_clear(&regions);
//    clear_temps(T_loc);
//
//    for _, city = range loop_city(i)
//    {
//        ilist_append(&all_cities, i);
//        bx[region(i)].temp++;
//    }
//
//
//    for _, i = range loop_loc(i)
//    {
//        if (bx[i].temp)
//            ilist_append(&regions, i);
//    }
//
//
//    ilist_scramble(regions);
//    reg_len = len(regions);
//
//    for (item = 0; rare_trade_items[item] && reg+1 < reg_len; item++)
//    {
//        one = regions[reg];
//        two = regions[reg+1];
//        ilist_clear(&cities);
//
//        for i = 0; i < len(all_cities); i++
//            if (region(all_cities[i]) == one ||
//                region(all_cities[i]) == two)
//                ilist_append(&cities, all_cities[i]);
//
//        if (!long_route(cities, rare_trade_items[item]))
//                if (!long_route(all_cities, rare_trade_items[item]))
//                {
//                log.Printf( ">>> Couldn't route %s!\n",
//                just_name(rare_trade_items[item]));
//                }
//
//        reg += 2;
//    }
//
//    for (; rare_trade_items[item]; item++)
//    {
//        if (!long_route(all_cities, rare_trade_items[item]))
//        {
//            log.Printf( ">>> Couldn't route %s!\n",
//                just_name(rare_trade_items[item]));
//        }
//    }
//}
//#endif

func seed_common_tradegoods() {
	var reg int
	var i, j int
	var cities []int
	var goods []int
	var source, consume int
	var item int
	count := 0
	qty, premium := 0, 0

	for _, reg = range loop_subkind(sub_region) {
		if loc_depth(reg) != LOC_region {
			continue
		}

		cities = nil
		goods = nil

		for _, i = range loop_city() {
			if region(i) == reg {
				cities = append(cities, i)
			}
		}

		for i = 0; common_trade_items[i] != FALSE; i++ {
			goods = append(goods, common_trade_items[i])
		}

		cities = shuffle_ints(cities)
		goods = shuffle_ints(goods)

		for i, j = 0, 0; i < len(cities) && j < len(goods); i, j = i+1, j+1 {
			source = cities[i]
			item = goods[j]

			qty = (base_price(item) * item_weight(item) % 40) +
				rnd(6, 12)

			add_city_trade(source, PRODUCE, item, qty,
				base_price(item)-rnd(0, 1), 0)

			qty += rnd(-qty/10, qty/10)
			premium = rnd(100, 200) / qty

			count = 0
			for {
				consume = cities[rnd(0, len(cities)-1)]
				flag := consume == source && count < 5
				count++
				if flag {
					continue
				}
				break
			}

			if consume != source {
				add_city_trade(consume, CONSUME, item, qty,
					base_price(item)+premium, 0)
			}
		}
	}

}

func seed_rare_tradegoods() {
	var reg int
	var i, j int
	var cities []int
	var goods []int
	//var source, consume int
	var item int
	//count := 0
	qty, premium := 0, 0

	for _, reg = range loop_subkind(sub_region) {
		if loc_depth(reg) != LOC_region {
			continue
		}

		cities = nil
		goods = nil

		for _, i = range loop_city() {
			if region(i) == reg {
				cities = append(cities, i)
			}
		}

		for i = 0; rare_trade_items[i] != FALSE; i++ {
			goods = append(goods, rare_trade_items[i])
		}

		cities = shuffle_ints(cities)
		goods = shuffle_ints(goods)

		for i, j = 0, 0; i < len(cities) && j < len(goods); i++ {
			if rnd(0, 1) != FALSE {
				continue
			}

			item, j = goods[j], j+1
			qty = 100 + (item+base_price(item))%250
			qty += rnd(-qty/10, qty/10)

			premium = rnd(1500, 3000) / qty

			if rnd(0, 1) == 0 {
				log.Printf("%s sold in %s for %s\n",
					box_name_qty(item, qty),
					box_name(cities[i]),
					comma_num(base_price(item)))

				add_city_trade(cities[i], PRODUCE, item, qty,
					base_price(item)-rnd(0, 1),
					rnd(1, NUM_MONTHS))
			} else {
				if rnd(0, 1) != FALSE {
					qty /= 2
				}

				add_city_trade(cities[i], CONSUME, item, qty,
					base_price(item)+premium, 0)

				log.Printf("%s bought in %s for %s\n",
					box_name_qty(item, qty),
					box_name(cities[i]),
					comma_num(base_price(item)+premium))
			}
		}
	}

}

func seed_city(where int) {
	var i, num int
	var p *entity_subloc
	//extern int new_ent_prime;
	flag := 0

	/*
	 *  Wed Dec 18 12:25:00 1996 -- Scott Turner
	 *
	 *  Cities don't teach skills any more, except rarely...
	 *
	 */
	if rnd(1, 20) == 10 {
		seed_city_skill(where)
	}
	seed_city_trade(where)
	/*
	 *  Build some random empty towers into each city.
	 *
	 *  Tue Apr 20 17:18:22 1999 -- Scott Turner
	 *
	 *  Turn the first tower into a Trading Guild.
	 */
	num = rnd(1, 4)
	for i = 0; i < num; i++ {
		new_ent_prime = TRUE
		newEnt := new_ent(T_loc, sub_tower)
		new_ent_prime = FALSE
		set_where(newEnt, where)
		set_name(newEnt, "Public tower")
		p = p_subloc(newEnt)
		p.hp = 100
		p.defense = fort_default_defense(sub_tower)
		if flag == FALSE {
			flag = 1
			set_name(newEnt, sout("%s Trader's Guild", just_name(where)))
			make_tower_guild(newEnt, sk_trading)
		}
	}
}

func seed_population() {
	var where int

	stage("seed_population()")

	for _, where = range loop_loc() {
		if loc_depth(where) != LOC_province {
			continue
		}

		if in_faery(where) || in_hades(where) || in_clouds(where) {
			continue
		}

		if subkind(where) != sub_forest &&
			subkind(where) != sub_plain {
			continue
		}
		/*
		 *  1-20 peasants everywhere...
		 *
		 */
		sub_item(where, item_peasant, has_item(where, item_peasant))
		gen_item(where, item_peasant, rnd(0, 20))

	}

	/* Cities need 1000 pop */
	for _, where = range loop_city() {
		gen_item(province(where), item_peasant, rnd(800, 1200))
	}

}

func seed_initial_locations() {

	var i int

	for _, i = range loop_city() {
		seed_city(i)
	}

	/*
	 *  Now update the markets to put initial trades into
	 *  all the trading guilds.
	 *
	 */
	update_markets()

	/* seed_common_tradegoods(); */
	/* seed_rare_tradegoods(); */

	for _, i = range loop_city() {
		do_production(i, TRUE)
	}

	seed_orcs()

	seed_has_been_run = TRUE
}

/*
 *  Sun Dec  1 19:27:52 1996 -- Scott Turner
 *
 *  Tax stuff removed.
 *
 */
func seed_taxes() {
	var where int
	//var base int
	//var pil int

	for _, where = range loop_loc() {
		if loc_depth(where) != LOC_province &&
			subkind(where) != sub_city {
			continue
		}

		if subkind(where) == sub_ocean {
			continue
		}

		if subkind(where) == sub_city {
			consume_item(where, item_petty_thief,
				has_item(where, item_petty_thief))

			gen_item(where, item_petty_thief, 1)
		}

		/*
		 *  Magician menial labor cookies
		 */

		consume_item(where, item_mage_menial,
			has_item(where, item_mage_menial))

		gen_item(where, item_mage_menial, 1)

	}
}

//#if 0
//static void
//init_gate_dests()
//{
//    var i int
//    var where int
//    var dest int
//    var l []*exit_view
//    int j;
//
//    loop_gate(i)
//    {
//        where = subloc(i);
//        dest = gate_dest(i);
//
//        if (in_faery(where) || in_faery(dest))
//            continue;
//
//        ilist_append(&p_misc(where).gate_dest, dest);
//        ilist_append(&p_misc(dest).gate_dest, where);
//    }
//    next_gate;
//
//    for _, i = range loop_loc(i)
//    {
//        if (loc_depth(i) == LOC_region)
//            continue;
//
//        l = exits_from_loc(0, i);
//
//        for (j = 0; j < len(l); j++)
//            ilist_append(&p_misc(i).prov_dest, l[j].destination);
//    }
//    next_loc(i);
//}
//
//
//int
//distance(int orig, int dest, int gate)
//{
//    var i, j int
//    var where int
//    flag := TRUE;
//    static int gate_init = FALSE;
//    struct entity_misc *p;
//
//    if (!gate_init)
//    {
//        init_gate_dests();
//        gate_init = TRUE;
//    }
//
//    clear_temps(T_loc);
//
//    bx[orig].temp = 1;
//
//    while (bx[dest].temp == 0 && flag)
//    {
//        flag = FALSE;
//
//        for _, where = for _, i = loop_loc()
//        {
//        if (bx[where].temp > 0)
//        {
//            p = rp_misc(where);
//
//            for i = 0; i < len(p.prov_dest); i++
//            {
//            j = p.prov_dest[i];
//            if (bx[j].temp == 0)
//            {
//                bx[j].temp = bx[where].temp + 1;
//                flag = TRUE;
//            }
//            }
//
//            if (gate)
//            {
//            for i = 0; i < len(p.gate_dest); i++
//            {
//                j = p.gate_dest[i];
//                if (bx[j].temp == 0)
//                {
//                bx[j].temp = bx[where].temp + 1;
//                flag = TRUE;
//                }
//            }
//            }
//
//            if (where != dest)
//            bx[where].temp = -1;
//        }
//        }
//
//    }
//
//    return bx[dest].temp;
//}
//#endif

/*
 *  Fri Apr 10 08:27:16 1998 -- Scott Turner
 *
 *  Seed orcs into the mountains -- this is a DM command.
 *  Should also be called when the map is initialized.
 *
 */
func seed_orcs() {
	var where, who, found int

	stage("seed_orcs()")
	for _, where = range loop_province() {
		/*
		 *  Only in the mountains
		 *
		 */
		if subkind(where) != sub_mountain {
			continue
		}
		/*
		 *  Already an orc stack here?
		 *
		 */
		found = FALSE
		for _, who = range loop_here(where) {
			//item := noble_item(who);
			if is_npc(who) && noble_item(who) == item_orc {
				found = TRUE
				break
			}
		}
		/*
		 *  About 15% of the mountains filled w/ orcs
		 *
		 */
		if found == FALSE && rnd(1, 100) < 15 {
			total := rnd(10, 20)
			var name string
			for total < 100 && rnd(1, 6) == 1 {
				total += rnd(10, 20)
			}
			/*
			 *  Now create it...
			 *
			 */
			name = fmt.Sprintf("Stack of %s", plural_item_name(item_orc, total))
			newOrc := new_char(sub_ni, item_orc, where, -1, indep_player,
				LOY_npc, 0, name)
			p_char(newOrc).break_point = 0
			rp_char(newOrc).npc_prog = PROG_orc
			gen_item(newOrc, item_orc, total)
			do_npc_orders(newOrc, 0, 0)
			wout(gm_player,
				"Created %d size new orc stack in %s.",
				total, box_name(where))
		}
	}
}

func v_seedorc(c *command) int {
	seed_orcs()
	wout(c.who, "Orcs seeded!")
	return TRUE
}
