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
 *  ADD_BUILD
 *  Tue Aug  6 12:10:25 1996 -- Scott Turner
 *
 *  Add an build to a thing.
 *
 */
func add_build(what, t, bm, er, eg int) bool {
	/*
	 *  Validity checks.
	 *
	 */
	if !valid_box(what) {
		return false
	}
	if rp_subloc(what) == nil {
		return false
	}
	/*
	 *  Allocate and fill in the new entity_build.
	 *
	 */
	newt := &entity_build{}
	newt.type_ = t
	newt.build_materials = bm
	newt.effort_required = er
	newt.effort_given = eg
	/*
	 *  Now append it to the entity_builds list.
	 *
	 */
	rp_subloc(what).builds = append(rp_subloc(what).builds, newt)
	return true
}

/*
 *  DELETE_ENTITY_BUILD
 *  Fri Sep 27 12:33:24 1996 -- Scott Turner
 *
 *  Delete the first entity_build of the given type.
 *
 */
func delete_build(what, type_ int) {
	/*
	 *  Validity checks.
	 *
	 */
	if !valid_box(what) {
		return
	}
	if rp_subloc(what) == nil {
		return
	}

	e := rp_subloc(what).builds
	if e == nil {
		return
	}

	for i := len(e) - 1; i >= 0; i-- {
		if e[i].type_ == type_ {
			rp_subloc(what).builds = rp_subloc(what).builds.delete(i)
			return
		}
	}
}

/*
 *  GET_BUILD
 *  Tue Aug  6 12:05:34 1996 -- Scott Turner
 *
 *  Get the first build of a type off of an build list.
 *
 */
func get_build(what, type_ int) *entity_build {
	/*
	 *  Validity checks.
	 *
	 */
	if !valid_box(what) {
		return nil
	}
	if rp_subloc(what) == nil {
		return nil
	}
	e := rp_subloc(what).builds
	/*
	 *  Possibly no builds, in which case we're done.
	 *
	 */
	if e == nil {
		return nil
	}
	/*
	 *  Look for the build.
	 *
	 */
	for i := 0; i < len(e); i++ {
		if e[i].type_ == type_ {
			return e[i]
		}
	}

	return nil
}

func fort_default_defense(sk int) int {
	switch sk {
	case sub_castle:
		return 40
	case sub_orc_stronghold:
		return 25
	case sub_tower:
		return 20
	case sub_galley:
		return 10
	}
	return 0
}

// if this is turned on you can find hidden sea routes with "build ship"!
func ship_loc_okay(c *command, where int) bool {
	// if has_ocean_access(where) {
	//     return true
	// }
	// wout(c.who, "%s is not an ocean port location.", box_name(where));
	// return false
	return true
}

/*
 *  Mon Aug 12 12:51:22 1996 -- Scott Turner
 *
 *  Modified to take into account terrain restrictions for temples.
 *
 */
func temple_loc_okay(c *command, where int) bool {
	if safe_haven(where) {
		wout(c.who, "Building is not permitted in safe havens.")
		return false
	}
	if loc_depth(where) == LOC_build {
		wout(c.who, "A temple may not be built inside another building.")
		return false
	}

	//#if 0
	//    /*
	//     *  Can't build if there's already a temple here.
	//     *
	//     */
	//    loop_all_here(where,i) {
	//      if (is_temple(i)) {
	//        wout(c.who, "There is already a temple here.");
	//        return FALSE;
	//      }
	//    } next_all_here;
	//#endif

	return true
}

/*
 *  Sun Mar  9 21:09:49 1997 -- Scott Turner
 *
 *  Orc stronghold
 *
 */
func real_orc_loc_okay(who, where int) bool {
	/* Gotta be an orc! */
	if !is_real_npc(who) || noble_item(who) != item_orc {
		wout(who, "Only orcs may build orc strongholds.")
		return false
	}

	if safe_haven(where) {
		return false
	}

	if loc_depth(where) == LOC_build {
		return false
	}

	if has_item(where, item_peasant) >= 100 {
		return false
	}

	/*
	 *  Can't build if there's already a stronghold here.
	 *
	 */
	for _, i := range loop_all_here(where) {
		if subkind(i) == sub_orc_stronghold {
			return false
		}
	}

	return true
}

func orc_loc_okay(c *command, where int) bool {
	return real_orc_loc_okay(c.who, where)
}

func tower_loc_okay(c *command, where int) bool {
	ld := loc_depth(where)

	if safe_haven(where) {
		wout(c.who, "Building is not permitted in safe havens.")
		return false
	}

	if ld != LOC_province &&
		ld != LOC_subloc &&
		subkind(where) != sub_castle &&
		subkind(where) != sub_castle_notdone {
		wout(c.who, "A tower may not be built here.")
		return false
	}

	if ld == LOC_build &&
		count_loc_structures(where, sub_tower, sub_tower_notdone) >= 6 {
		wout(c.who, "Six towers at most can be built within a %s.",
			subkind_s[subkind(where)])
		return false
	}

	return true
}

func mine_loc_okay(c *command, where int) bool {

	if subkind(where) != sub_mountain && subkind(where) != sub_rocky_hill {
		wout(c.who, "Mines may only be built in mountain provinces and rocky hills.")
		return false
	}

	if safe_haven(where) {
		wout(c.who, "Building is not permitted in safe havens.")
		return false
	}

	if count_loc_structures(where, sub_mine, sub_mine_notdone) != FALSE {
		wout(c.who, "A location may not have more than one mine.")
		return false
	}

	if count_loc_structures(where, sub_mine_collapsed, 0) != FALSE {
		wout(c.who, "Another mine may not be built here until the collapsed mine vanishes.")
		return false
	}

	return true
}

func mine_shaft_loc_okay(c *command, where int) bool {
	if subkind(where) != sub_mountain && subkind(where) != sub_mine_shaft {
		wout(c.who, "Mine shafts must be built in mountains or mine shafts.")
		return false
	}

	/*
	 *  Wed Aug 11 11:23:22 1999 -- Scott Turner
	 *
	 *  This is also used to count mine shafts, which aren't "here", they're "down".
	 *
	 */
	hasShaft := count_loc_structures(where, sub_mine_shaft, sub_mine_shaft_notdone) != FALSE
	if !hasShaft {
		if i := location_direction(where, DIR_DOWN); kind(i) == T_loc && ((subkind(i) == sub_mine_shaft) || (subkind(i) == sub_mine_shaft_notdone)) {
			hasShaft = true
		}
	}
	if hasShaft {
		wout(c.who, "There's already a mine shaft here.")
		return false
	}

	// can't go deeper than 20...
	if mine_depth(where)+1 >= MINE_MAX {
		wout(c.who, "It is impossible to dig any deeper at these great depths.")
		return false
	}

	return true
}

func inn_loc_okay(c *command, where int) bool {
	if safe_haven(where) {
		wout(c.who, "Building is not permitted in safe havens.")
		return false
	}

	if loc_depth(where) != LOC_province && subkind(where) != sub_city {
		wout(c.who, "Inns may only be built in cities and provinces.")
		return false
	}

	return true
}

// return a sublocation if it can be found in the province or in the city.
func province_subloc(where, sk int) int {
	prov := province(where)
	city := city_here(prov)
	if n := subloc_here(prov, sk); n != FALSE {
		return n
	} else if city != FALSE {
		return subloc_here(city, sk)
	}
	return 0
}

func castle_loc_okay(c *command, where int) bool {
	if safe_haven(where) {
		wout(c.who, "Building is not permitted in safe havens.")
		return false
	} else if loc_depth(where) != LOC_province &&
		subkind(where) != sub_city {
		wout(c.who, "A castle must be built in a province or a city.")
		return false
	} else if province_subloc(where, sub_castle) != FALSE || province_subloc(where, sub_castle_notdone) != FALSE {
		wout(c.who, "This province already contains a castle.  Another may not be built here.")
		return false
	}
	return true
}

/*
 *  Wed Dec  4 13:45:00 1996 -- Scott Turner
 *
 *  Rules for building a city:
 *    (1) 1000 pop
 *    (2) 100+ pop in all adjacent (non-ocean) provinces
 *    (3) No city w/in 5 provinces.
 *
 * int los_province_distance(a,b)
 */
func habitable(n int) bool {
	return (valid_box(n) && ((subkind(n) >= sub_forest && subkind(n) <= sub_under) || subkind(n) == sub_cloud))
}

func city_loc_okay(c *command, where int) bool {
	if province(where) != where {
		wout(c.who, "You must build a city in a province.")
		return false
	}

	/*
	 *  1000 pop
	 *
	 */
	if has_item(where, item_peasant) < 1000 {
		wout(c.who, "There is not enough population here to support a city.")
		return false
	}

	/*
	 *  100+ pop in adjacent provinces.
	 *
	 */
	l := exits_from_loc_nsew(0, where)
	for i := 0; i < len(l); i++ {
		here := l[i].destination
		if loc_depth(here) != LOC_province {
			continue
		}
		if habitable(here) && has_item(here, item_peasant) < 100 {
			wout(c.who, "All adjacent habitable provinces must have a population", "of at least 100 peasants.")
			return false
		}
	}

	// no cities within 5 provinces walking.
	for _, here := range loop_city() {
		if los_province_distance(where, province(here)) < 5 {
			wout(c.who, "Too near a city to build another city.")
			return false
		}
	}

	return true
}

type build_ent struct {
	what               string /* what are we building? */
	skill_req, skill2  int    /* one or the other */
	kind               int
	loc_ok             func(c *command, where int) bool
	unfinished_subkind int
	finished_subkind   int
	min_workers        int /* min # of workers to begin */
	worker_days        int /* time to complete */
	min_days           int /* soonest can be completed */
	req_item           int
	req_qty            int
	req_item2          int
	req_qty2           int
	default_name       string
	direction          int /* Direction to new location */
	num                int /* Multiplier for ship hulls */
}

var build_tbl = []build_ent{
	//#if 0
	//                {
	//                "galley",
	//                sk_shipbuilding, 0,
	//                T_ship,
	//                ship_loc_okay,				/* can we build here? */
	//                sub_galley_notdone, sub_galley,		/* ship types */
	//                3,					/* minimum # of workers */
	//                250,					/* worker-days to complete */
	//                1,					/* at least n days */
	//                item_lumber, 10,			/* required item, 1/5 qty */
	//                0, 0,					/* required item #2, 1/5 qty */
	//                "New galley"				/* default name */
	//                },
	//                {
	//                "roundship",
	//                sk_shipbuilding, 0,
	//                T_ship,
	//                ship_loc_okay,				/* can we build here? */
	//                sub_roundship_notdone, sub_roundship,	/* ship types */
	//                3,					/* minimum # of workers */
	//                500,					/* worker-days to complete */
	//                1,					/* at least n days */
	//                item_lumber, 20,			/* required item, 1/5 qty */
	//                0, 0,					/* required item #2, 1/5 qty */
	//                "New roundship"				/* default name */
	//                },
	//                {
	//                "raft",
	//                0, 0,
	//                T_ship,
	//                ship_loc_okay,				/* can we build here? */
	//                sub_raft_notdone, sub_raft,		/* ship types */
	//                0,					/* minimum # of workers */
	//                45,					/* worker-days to complete */
	//                1,					/* at least n days */
	//                item_flotsam, 5,			/* required item, 1/5 qty */
	//                0, 0,					/* required item #2, 1/5 qty */
	//                "New raft"				/* default name */
	//                },
	//#endif
	build_ent{
		"ship",
		sk_shipbuilding, 0,
		T_ship,
		ship_loc_okay,              /* can we build here? */
		sub_ship_notdone, sub_ship, /* ship types */
		3,              /* minimum # of workers */
		100,            /* worker-days to complete */
		1,              /* at least n days */
		item_lumber, 2, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New ship", /* default name */
		FALSE, 0},
	{
		/*  Synonym for "build ship" */
		"hull",
		sk_shipbuilding, 0,
		T_ship,
		ship_loc_okay,              /* can we build here? */
		sub_ship_notdone, sub_ship, /* ship types */
		3,              /* minimum # of workers */
		100,            /* worker-days to complete */
		1,              /* at least n days */
		item_lumber, 2, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New ship", /* default name */
		FALSE, 0},
	{
		"temple",
		sk_construction, 0,
		T_loc,
		temple_loc_okay, /* can we build here? */
		sub_temple_notdone, sub_temple,
		3,              /* minimum # of workers */
		1000,           /* worker-days to complete */
		1,              /* at least n days */
		item_stone, 10, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New temple", /* default name */
		FALSE, 0},
	{
		"inn",
		sk_construction, 0,
		T_loc,
		inn_loc_okay, /* can we build here? */
		sub_inn_notdone, sub_inn,
		3,               /* minimum # of workers */
		300,             /* worker-days to complete */
		1,               /* at least n days */
		item_lumber, 15, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New inn", /* default name */
		FALSE, 0},
	{
		"castle",
		sk_construction, 0,
		T_loc,
		castle_loc_okay, /* can we build here? */
		sub_castle_notdone, sub_castle,
		3,               /* minimum # of workers */
		10000,           /* worker-days to complete */
		1,               /* at least n days */
		item_stone, 100, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New castle", /* default name */
		FALSE, 0},
	{
		"stronghold",
		0, 0,
		T_loc,
		orc_loc_okay, /* can we build here? */
		sub_orc_stronghold_notdone, sub_orc_stronghold,
		3,    /* minimum # of workers */
		2500, /* worker-days to complete */
		1,    /* at least n days */
		0, 0, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"An Orc Stronghold", /* default name */
		FALSE, 0},
	{
		"city",
		sk_build_city, 0,
		T_loc,
		city_loc_okay, /* can we build here? */
		sub_city_notdone, sub_city,
		3,               /* minimum # of workers */
		50000,           /* worker-days to complete */
		14,              /* at least n days */
		item_stone, 400, /* required item, 1/5 qty */
		item_lumber, 100, /* required item #2, 1/5 qty */
		"New City", /* default name */
		FALSE, 0},
	//#if 0
	//                {
	//                "mine",
	//                sk_construction, sk_mining,
	//                T_loc,
	//                mine_loc_okay,				/* can we build here? */
	//                sub_mine_notdone, sub_mine,
	//                3,					/* minimum # of workers */
	//                500,					/* worker-days to complete */
	//                1,					/* at least n days */
	//                item_lumber, 5,				/* required item, 1/5 qty */
	//                0, 0,					/* required item #2, 1/5 qty */
	//                "New mine"				/* default name */
	//                },
	//#endif
	{
		"shaft",
		sk_deepen_mine, 0,
		T_loc,
		mine_shaft_loc_okay, /* can we build here? */
		sub_mine_shaft_notdone, sub_mine_shaft,
		3,    /* minimum # of workers */
		500,  /* worker-days to complete */
		1,    /* at least n days */
		0, 0, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New mine shaft", /* default name */
		DIR_DOWN, 0},
	{
		/* Synonym for "BUILD SHAFT" */
		"mine",
		sk_deepen_mine, 0,
		T_loc,
		mine_shaft_loc_okay, /* can we build here? */
		sub_mine_shaft_notdone, sub_mine_shaft,
		3,    /* minimum # of workers */
		500,  /* worker-days to complete */
		1,    /* at least n days */
		0, 0, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New mine shaft", /* default name */
		DIR_DOWN, 0},
	{
		"tower",
		sk_construction, 0,
		T_loc,
		tower_loc_okay, /* can we build here? */
		sub_tower_notdone, sub_tower,
		3,              /* minimum # of workers */
		2000,           /* worker-days to complete */
		1,              /* at least n days */
		item_stone, 20, /* required item, 1/5 qty */
		0, 0, /* required item #2, 1/5 qty */
		"New tower", /* default name */
		FALSE, 0,
	},
	{}}

var fuzzy_build_match bool

func find_build(s string) *build_ent {
	fuzzy_build_match = false

	for i := 0; build_tbl[i].what != ""; i++ {
		if i_strcmp([]byte(build_tbl[i].what), []byte(s)) == 0 {
			return &build_tbl[i]
		}
	}

	for i := 0; build_tbl[i].what != ""; i++ {
		if fuzzy_strcmp([]byte(build_tbl[i].what), []byte(s)) {
			fuzzy_build_match = true // mdhender: moved here, why?
			return &build_tbl[i]
		}
	}

	return nil
}

func build_materials_check(c *command, bi *build_ent) bool {

	if bi.skill_req != FALSE { /* if a skill is required... */
		if bi.skill2 != FALSE { /* either one of two skills */
			if has_skill(c.who, bi.skill_req) < 1 &&
				has_skill(c.who, bi.skill2) < 1 {
				wout(c.who, "Building a %s requires either %s or %s.",
					bi.what,
					box_name(bi.skill_req),
					box_name(bi.skill2))
				return false
			}
		} else { /* single skill requirement */
			if has_skill(c.who, bi.skill_req) < 1 {
				wout(c.who, "Building a %s requires %s.",
					bi.what,
					box_name(bi.skill_req))
				return false
			}
		}
	}

	/*
	 *  Materials check
	 */

	if bi.req_item > 0 && has_item(c.who, bi.req_item) < bi.req_qty*bi.num {
		wout(c.who, "Need %s to start.",
			box_name_qty(bi.req_item, bi.req_qty*bi.num))
		return false
	}

	if bi.req_item2 > 0 && has_item(c.who, bi.req_item2) < bi.req_qty2*bi.num {
		wout(c.who, "Need %s to start.",
			box_name_qty(bi.req_item2, bi.req_qty2*bi.num))
		return false
	}

	if !is_real_npc(c.who) &&
		effective_workers(c.who) < bi.min_workers {
		wout(c.who, "Need at least %s for construction.",
			box_name_qty(item_worker, bi.min_workers))
		return false
	}

	if is_real_npc(c.who) &&
		noble_item(c.who) != FALSE &&
		has_item(c.who, noble_item(c.who)) < bi.min_workers {
		wout(c.who, "Hey npc you need at least %s for construction!",
			box_name_qty(noble_item(c.who), bi.min_workers))
		return false
	}
	/*
	 *  Materials deduct
	 */

	if bi.req_item > 0 {
		consume_item(c.who, bi.req_item, bi.req_qty*bi.num)
	}

	if bi.req_item2 > 0 {
		consume_item(c.who, bi.req_item2, bi.req_qty2*bi.num)
	}

	return true
}

/*
 *  Wed Feb 12 11:28:14 1997 -- Scott Turner
 *
 *  Connection loc1 to loc2.
 *
 */
func connect_locations(loc1, dir1 int, loc2, dir2 int) {
	p1 := p_loc(loc1)
	p2 := p_loc(loc2)

	assert(dir1 != FALSE && dir1 < MAX_DIR)
	assert(dir2 != FALSE && dir2 < MAX_DIR)
	assert(p1 != nil)
	assert(p2 != nil)

	// make sure both loc1 && loc2 have enough prov_dest.
	for len(p1.prov_dest) < dir1 {
		p1.prov_dest = append(p1.prov_dest, 0)
	}
	for len(p2.prov_dest) < dir2 {
		p1.prov_dest = append(p2.prov_dest, 0)
	}

	// make sure both locations don't already have something in that direction.
	assert(p1.prov_dest[dir1-1] == FALSE)
	assert(p2.prov_dest[dir2-1] == FALSE)

	// add accordingly
	p1.prov_dest[dir1-1] = loc2
	p2.prov_dest[dir2-1] = loc1
}

/*
 *  Wed Feb 12 11:28:14 1997 -- Scott Turner
 *
 *  Unconnect loc1 from the map...
 *
 */
func unconnect_location(loc1 int) {
	var i, j int
	p1 := p_loc(loc1)
	assert(p1 != nil)

	// go through all of loc1's prov_dests and remove loc2.
	for i = 0; i < len(p1.prov_dest); i++ {
		if p1.prov_dest[i] == FALSE {
			continue
		}
		p2 := p_loc(p1.prov_dest[i])
		// remove loc1 from p2's connections...
		for j = 0; j < len(p2.prov_dest); j++ {
			if p2.prov_dest[j] == loc1 {
				p2.prov_dest[j] = 0
			}
		}
		// now remove it p2 from p1
		p1.prov_dest[i] = 0
		return
	}
}

/*
 *  Fri Jan 10 16:28:24 1997 -- Scott Turner
 *
 *  Modified to use build entities.
 *
 *  Wed Feb 12 09:25:43 1997 -- Scott Turner
 *
 *  If "direction" exists, then we make this not a subloc but a full-fledged
 *  location pointing off in the required direction.
 *
 *  Fri Apr 23 12:56:51 1999 -- Scott Turner
 *
 *  Need to do a "touch_loc" if this creates a new location.
 *
 */
func create_new_building(c *command, bi *build_ent, where int) {
	var p *entity_subloc
	//var b *entity_build

	change_box_subkind(where, bi.finished_subkind)

	/*
	 *  Is there a build structure?
	 *
	 */
	if get_build(where, BT_BUILD) != nil {
		delete_build(where, BT_BUILD)
	}

	/*
	 *  Modified to teach the priest's religion.
	 *
	 *  Fri Nov 15 13:42:43 1996 -- Scott Turner
	 *
	 *  Don't set the "guild" -- do that through a priest spell.
	 *
	 */

	if bi.skill_req != FALSE {
		add_skill_experience(c.who, bi.skill_req)
	}

	/*
	 *  Wed Feb 12 09:29:15 1997 -- Scott Turner
	 *
	 *  Now either fix up the subloc info, or set it up as a location.
	 *
	 */
	if bi.direction != FALSE {
		parent := province(loc(where))
		/*
		 *  It's a location, now.  So remove it from where it is (a subloc)
		 *  and put it in the parent location's region.
		 *
		 */
		set_where(where, loc(parent))
		/*
		 *  And touch everyone who is now in here.
		 *
		 */
		touch_loc_after_move(c.who, where)
		/*
		 *  Now connect it to this province in the requested direction.
		 *
		 */
		connect_locations(parent, bi.direction, where, exit_opposite[bi.direction])
		/*
		 *  We need to also set its "near_grave"
		 *
		 */
		assert(rp_loc(parent) != nil)
		assert(rp_loc(parent).near_grave != FALSE)
		p_loc(where).near_grave = rp_loc(parent).near_grave
	} else {
		/*
		 *  It's a sublocation.
		 *
		 */
		p = p_subloc(where)

		p.hp = 100
		p.defense = fort_default_defense(bi.finished_subkind)

		//#if 0
		//        if (bi.finished_subkind == sub_mine)
		//          {
		//        p.shaft_depth = 3;
		//        if (rnd(1,5) == 1)
		//          gen_item(where, item_gate_crystal, 1);
		//        mine_production(where);
		//          }
		//#endif

		/*
		 *  Thu Jan  2 14:46:05 1997 -- Scott Turner
		 *
		 *  Finished ships get a "free" rowing port.
		 *
		 */
		if bi.finished_subkind == sub_ship {
			p_ship(where).ports = 1
		}

		/*
		 *  Thu Feb  6 11:44:45 1997 -- Scott Turner
		 *
		 *  New mines get a "mine_info"
		 *
		 */
		if bi.finished_subkind == sub_mine {
			create_mine_info(where)
		}
		/*
		 *  Mon Feb 26 17:12:02 2001 -- Scott Turner
		 *
		 *  Seed new cities.
		 *
		 */
		if bi.finished_subkind == sub_city {
			seed_city(where)
		}
	}
}

func start_build(c *command, bi *build_ent, id int) bool {
	var new_name string
	var p *entity_subloc
	var b *entity_build
	pl := player(c.who)
	instant_build := false /* build takes 1 day and no effort */
	where := subloc(c.who)
	//extern int new_ent_prime;        /* allocate short numbers */

	/*
	 *  Thu Jan  2 14:10:40 1997 -- Scott Turner
	 *
	 *  If you're building a "ship", need to specify the number
	 *  of hulls.
	 *
	 */
	bi.num = 1
	if bi.finished_subkind == sub_ship {
		if c.c == 0 {
			wout(c.who, "When building a ship, you must specify the number of hulls after the ship name.  For this build I'll assume a 1 hull ship.")
		} else {
			bi.num = c.c
		}
	}

	/*
	 *  First tower a player builds take no workers, materials or time
	 */
	if bi.finished_subkind == sub_tower &&
		p_player(pl).first_tower == FALSE {
		instant_build = true
		p_player(pl).first_tower = TRUE
	}

	if !instant_build && !build_materials_check(c, bi) {
		return false
	}

	if id != FALSE {
		if kind(id) != T_unform ||
			ilist_lookup(p_player(player(c.who)).unformed, id) < 0 {
			wout(c.who, "%s is not a valid unformed entity.", box_code(id))
			wout(c.who, "I will use a random identifier.")
			id = 0
		}
	}

	var newt int
	if id != FALSE {
		newt = id
		change_box_kind(id, bi.kind)
		change_box_subkind(id, bi.unfinished_subkind)
	} else {
		new_ent_prime = true
		newt = new_ent(bi.kind, bi.unfinished_subkind)
		new_ent_prime = false
	}

	p_player(player(c.who)).unformed = rem_value((p_player(player(c.who)).unformed), id)
	set_where(newt, where)

	if numargs(c) < 2 || c.parse[2] == nil || len(c.parse[2]) == 0 {
		new_name = bi.default_name
	} else {
		new_name = string(c.parse[2])
	}

	if len(new_name) > 25 {
		wout(c.who, "The name you gave is too long.  Place names must be 25 characters or less.  Please use the NAME order to set a shorter name next turn.")

		new_name = bi.default_name
	}

	set_name(newt, new_name)

	p = p_subloc(newt)

	/*
	 *  Thu Jan  2 14:44:39 1997 -- Scott Turner
	 *
	 *  If this is a ship, then set the number of hulls now, so
	 *  that it can't be changed later.
	 *
	 */
	if bi.finished_subkind == sub_ship {
		p_ship(newt).hulls = bi.num
	}

	if instant_build {
		change_box_subkind(newt, bi.finished_subkind)
		wout(c.who, "Built %s.", box_name_kind(newt))

		show_to_garrison = true
		wout(where, "%s built %s in %s.",
			box_name(c.who),
			box_name_kind(newt),
			box_name(where))
		show_to_garrison = false
	} else {
		wout(c.who, "Created %s.", box_name_kind(newt))
		show_to_garrison = true
		wout(where, "%s began construction of %s in %s.",
			box_name(c.who),
			box_name_kind(newt),
			box_name(where))
		show_to_garrison = false
	}

	move_stack(c.who, newt)

	if instant_build {
		create_new_building(c, bi, newt)

		c.wait = 1
		c.inhibit_finish = true
		return true
	}

	/*
	 *  Try to find the build record.  If nothing, create one.
	 *
	 */
	b = get_build(newt, BT_BUILD)
	if b == nil {
		add_build(newt, BT_BUILD, 0, 0, 0)
		b = get_build(newt, BT_BUILD)
	}
	b.effort_required = bi.worker_days * 100 * bi.num
	b.effort_given = 0
	p.damage = 0
	b.build_materials = 0

	return true
}

func daily_build(c *command, bi *build_ent) bool {
	var nworkers int
	inside := subloc(c.who)
	var b *entity_build
	var bonus int
	var effort_given int

	/*
	 *  todo:  apply building energy to repair structure if damaged
	 *		currently, damage figure gets erased when
	 *		the structure is completed.
	 */

	if subkind(inside) == schar(bi.finished_subkind) {
		wout(c.who, "%s is finished!",
			box_name(inside))
		c.wait = 0
		return true
	}

	if subkind(inside) != schar(bi.unfinished_subkind) {
		wout(c.who, "%s is no longer in a %s.  Construction halts.",
			just_name(c.who), bi.what)
		return false
	}

	if bi.min_workers == 0 {
		nworkers = 1
	} else {
		nworkers = effective_workers(c.who)
	}

	/*
	 *  Mon Mar 10 08:48:26 1997 -- Scott Turner
	 *
	 *  All orcs act as workers, at least if they have an orc
	 *  leader, and they can build orc strongholds.
	 *
	 */
	if is_real_npc(c.who) && noble_item(c.who) != FALSE {
		nworkers += has_item(c.who, noble_item(c.who))
	}

	if nworkers <= 0 {
		wout(c.who, "%s has no workers.  Construction halts.",
			just_name(c.who))
		return false
	}

	//p := p_subloc(inside);

	/*
	 *  Give a 5% speed bonus for each experience level of the
	 *  construction skill
	 */
	bonus = 5 * c.use_exp * nworkers
	effort_given = nworkers*100 + bonus

	/*
	 *  Materials check
	 *
	 *  Thu Jan  2 14:49:03 1997 -- Scott Turner
	 *
	 *  First adjust bi.num if this is a ship.
	 *
	 */
	bi.num = 1
	if rp_ship(inside) != nil {
		bi.num = rp_ship(inside).hulls
	}

	/*
	 *  Try to find the build record.  If nothing, problem!
	 *
	 */
	b = get_build(inside, BT_BUILD)
	if b == nil {
		wout(c.who, "Some mysterious force prevents you from completing this building.")
		wout(c.who, "Report this to the GM.")
		return false
	}

	if bi.req_item > 0 {
		fifth := (b.effort_given + effort_given) * 5 /
			b.effort_required

		for fifth < 5 && b.build_materials < fifth {
			if !consume_item(c.who, bi.req_item, bi.req_qty*bi.num) {
				wout(c.who, "Need another %s to continue work.  Construction halted.",
					box_name_qty(bi.req_item, bi.req_qty*bi.num))
				return false
			}

			if bi.req_item2 > 0 &&
				!consume_item(c.who, bi.req_item2, bi.req_qty2*bi.num) {
				wout(c.who, "Need another %s to continue work.  Construction halted.",
					box_name_qty(bi.req_item2, bi.req_qty2*bi.num))
				return false
			}

			b.build_materials++
		}
	}

	b.effort_given += effort_given

	if b.effort_given < b.effort_required ||
		command_days(c) < bi.min_days {
		return true
	}

	/*
	 *  It's done
	 */
	create_new_building(c, bi, inside)

	wout(c.who, "%s is finished!", box_name(inside))

	c.wait = 0
	return true
}

func build_structure(c *command, bi *build_ent, id int) bool {
	who := c.who
	where := subloc(who)

	if loc_depth(where) == LOC_build {
		if subkind(where) == schar(bi.finished_subkind) {
			wout(who, "%s is already finished.", box_name(where))
			return false
		}

		if subkind(where) == schar(bi.unfinished_subkind) {
			wout(who, "Continuing work on %s.", box_name(where))
			return true
		}
	}

	if subkind(where) == sub_ocean {
		wout(who, "Construction may not take place at sea.")
		return false
	}

	if !bi.loc_ok(c, where) {
		return false
	}

	return start_build(c, bi, id)
}

func unfinished_building(who int) string {
	where := subloc(who)
	switch subkind(where) {
	case sub_castle_notdone:
		return "castle"
	case sub_tower_notdone:
		return "tower"
	case sub_temple_notdone:
		return "temple"
	case sub_galley_notdone:
		return "galley"
	case sub_roundship_notdone:
		return "roundship"
	case sub_inn_notdone:
		return "inn"
	case sub_mine_notdone:
		return "mine"
	}

	return ""
}

func v_build(c *command) int {
	var t *build_ent
	var s string
	var days, id int
	where := subloc(c.who)
	partial := false

	if numargs(c) < 1 {
		if s = unfinished_building(c.who); s != "" {
			wout(c.who, "(assuming you meant 'build %s')", s)
			ret := oly_parse(c, []byte(sout("build %s", s)))
			assert(ret)
			return v_build(c)
		}
	}

	if numargs(c) < 1 {
		wout(c.who, "Must specify what to build.")
		return FALSE
	}

	t = find_build(string(c.parse[1]))

	if t == nil {
		wout(c.who, "Don't know how to build '%s'.",
			c.parse[1])
		return FALSE
	}

	if fuzzy_build_match {
		wout(c.who, "(assuming you meant 'build %s')", t.what)
	}

	/*
	 *  Fri Jan 29 13:59:30 1999 -- Scott Turner
	 *
	 *  This is very ugly.  We have at least the following cases:
	 *
	 *  build castle "Foo" 7 ae4
	 *  build castle 7 ae4  -- continuation or default name!
	 *  build castle 7 -- continuation or default name!
	 *          a    b   c d  e
	 *  build ship "Foo" 3 7 ae4
	 *  build ship 3 7 ae4 -- start w/ default name
	 *  build ship 3 0 ae4 -- continuation.
	 *  build ship 3
	 *
	 *  Mon Feb  1 12:52:55 1999 -- Scott Turner
	 *
	 *  We can distinguish 3 cases:
	 *
	 *  -- starting a structure
	 *  -- starting a ship
	 *  -- continuing anything
	 *
	 *  I think each of these have fixed arguments.  However, to do
	 *  this we need to know if we're in a partially completed
	 *  structure.
	 *
	 */
	if loc_depth(where) == LOC_build &&
		subkind(where) == schar(t.unfinished_subkind) {
		partial = true
	}

	if partial {
		/*
		 *  Partial is optionally followed by max days.  Try to
		 *  "fix" things if someone has specified a name.
		 */
		if c.b == 0 && c.c > 0 {
			wout(c.who, "Do not specify a name when continuing a build.  I'll try to ignore the name and figure out what you meant.")
			/* is this even close? :-) */
			c.b = c.c
			c.parse[2] = c.parse[3]
			c.c = c.d
			c.parse[3] = c.parse[4]
		}
		days = c.b
	} else {
		/*
		 *  Depends on if it is a ship, in which case
		 *  we also have to skip over the number of hulls.
		 */
		if t.finished_subkind == sub_ship {
			days = c.d
			id = c.e
		} else {
			days = c.c
			id = c.d
		}
	}

	if days != FALSE {
		c.wait = days
	}
	if build_structure(c, t, id) {
		return TRUE
	}
	return FALSE
}

func d_build(c *command) int {
	var t *build_ent

	t = find_build(string(c.parse[1]))

	if t == nil {
		log_output(LOG_CODE, "d_build: t is nil (%s)", c.parse[1])
		out(c.who, "Internal error.")
		return FALSE
	}

	if daily_build(c, t) {
		return TRUE
	}
	return FALSE
}

/*
 *  Worker-days to repair a structure, for different structures
 */

func repair_points(k int) int {
	switch k {
	case sub_castle:
		return 3
	case sub_castle_notdone:
		return 3
	case sub_tower:
		return 2
	case sub_tower_notdone:
		return 2
	case sub_temple:
		return 2
	case sub_temple_notdone:
		return 2
	case sub_inn:
		return 2
	case sub_inn_notdone:
		return 2
	case sub_mine:
		return 2
	case sub_mine_notdone:
		return 2
	case sub_mine_shaft:
		return 2
	case sub_mine_shaft_notdone:
		return 2
	case sub_guild:
		return 2
	case sub_galley:
		return 1
	case sub_galley_notdone:
		return 1
	case sub_roundship:
		return 1
	case sub_roundship_notdone:
		return 1
	case sub_ship_notdone:
		return 1
	case sub_ship:
		return 1
	case sub_orc_stronghold_notdone:
		return 1
	case sub_orc_stronghold:
		return 1
	}

	panic("!reached")
}

func v_repair(c *command) int {
	days := c.a
	where := subloc(c.who)
	var workers int
	req_item := 0
	//fort_def := fort_default_defense(int(subkind(where)));

	if days < 1 {
		days = -1
	}

	if loc_depth(where) != LOC_build {
		wout(c.who, "%s may not be repaired.", box_name(where))
		return FALSE
	}

	if loc_damage(where) < 1 {
		wout(c.who, "%s is not damaged.", box_name(where))
		return FALSE
	}

	workers = effective_workers(c.who)

	if workers < 1 {
		wout(c.who, "Need at least one %s.", box_name(item_worker))
		return FALSE
	}

	switch subkind(where) {
	case sub_galley, sub_roundship, sub_ship:
		req_item = item_glue
		break
	}

	if req_item != FALSE && !consume_item(c.who, req_item, 1) {
		wout(c.who, "%s repair requires %s.",
			cap_(subkind_s[subkind(where)]),
			box_name_qty(req_item, 1))
		return FALSE
	}

	c.d = 0 /* remainder */
	c.e = 0
	c.f = 0

	c.wait = days
	return TRUE
}

func d_repair(c *command) int {
	where := subloc(c.who)
	var workers int
	var per_point int
	var p *entity_subloc
	var points int

	if loc_depth(where) != LOC_build {
		wout(c.who, "No longer in a repairable structure.")
		return FALSE
	}

	if loc_damage(where) < 1 {
		wout(c.who, "%s has been fully repaired.", box_name(where))
		c.wait = 0
		return TRUE
	}

	workers = effective_workers(c.who)

	if workers < 1 {
		wout(c.who, "No longer have at least one %s.",
			box_name(item_worker))
		return FALSE
	}

	per_point = repair_points(int(subkind(where)))

	workers += c.d
	c.d = workers % per_point
	points = workers / per_point

	p = p_subloc(where)

	if p.damage > 0 {
		if points > p.damage {
			points -= p.damage
			c.e += p.damage
			p.damage = 0
		} else {
			p.damage -= points
			c.e += points
			points = 0
		}
	}

	if p.damage < 1 {
		wout(c.who, "%s has been fully repaired.", box_name(where))
		i_repair(c)
		c.wait = 0
		return TRUE
	}

	if c.wait == 0 {
		i_repair(c)
	}

	return TRUE
}

func i_repair(c *command) int {
	where := subloc(c.who)

	vector_char_here(where)

	wout(VECT, "%s repaired %s damage to %s.",
		box_name(c.who),
		nice_num(c.e), box_name(where))

	return TRUE
}

/*
 *  Tue Aug 13 09:19:55 1996 -- Scott Turner
 *
 *  Modified to permit razing of a building from outside, if the
 *  building is empty.
 *
 *  If there's no target, then assume that he's razing
 *  the building he's in.
 */
func v_raze(c *command) int {
	where := subloc(c.who)
	target := or_int(c.a != FALSE, c.a, where)
	time := c.b
	//var men int
	//var per_point int

	/*
	 *  Best have a target.
	 *
	 */
	if !valid_box(target) {
		wout(c.who, "You must specify a building to RAZE.")
		return FALSE
	}

	/*
	 *  Better be a structure of some sort.
	 *
	 */
	if loc_depth(target) != LOC_build {
		wout(c.who, "You can only RAZE buildings.")
		return FALSE
	}

	/*
	 *  Case 1: We're in the building
	 *
	 */
	if where == target {
		if building_owner(where) != c.who {
			wout(c.who, "Must be the owner of a structure to RAZE.")
			return FALSE
		}
		/*
		 *  Set the target correctly for d_raze.
		 *
		 */
		c.a = where

		/*
		 *  Case 2: We're outside of the building.
		 *
		 */
	} else if (subloc(c.who) == subloc(target)) ||
		(is_ship(target) &&
			is_ship(subloc(c.who)) &&
			province(c.who) == province(target)) {
		/*
		 *  The building must be empty.
		 *
		 */
		if first_character(target) != FALSE {
			wout(c.who, "Cannot RAZE an occupied structure.")
			return FALSE
		}
	} else {
		wout(c.who, "Must be adjacent to a structure to RAZE it.")
		return FALSE
	}

	//per_point := repair_points(int(subkind(target))); // mdhender: not used?
	//men := count_man_items(c.who) + 1; // mdhender: not used?

	//#if 0
	//    if (men < per_point)
	//    {
	//        wout(c.who, "Need at least %d men to harm a %s.",
	//                    per_point, subkind_s[subkind(target)]);
	//        return FALSE;
	//    }
	//#endif

	if time != FALSE {
		c.wait = time
		wout(c.who, "Razing %s for %s days.",
			box_name(c.a), nice_num(c.b))
	}

	c.d = 0 /* remainder */

	return TRUE
}

func d_raze(c *command) int {
	where := subloc(c.who)
	target := or_int(c.a != FALSE, c.a, where)
	var men int
	var per_point int
	var points int

	/*
	 *  Best have a target.
	 *
	 */
	if !valid_box(target) {
		wout(c.who, "That building no longer exists!")
		wout(c.who, "Maybe someone else finished razing it before you.")
		return FALSE
	}

	/*
	 *  Better be a structure of some sort.
	 *
	 */
	if loc_depth(target) != LOC_build {
		wout(c.who, "You can only RAZE buildings.")
		return FALSE
	}

	/*
	 *  Case 1: We're in the building
	 *
	 */
	if where == target {
		if building_owner(where) != c.who {
			wout(c.who, "Must be the owner of a structure to RAZE.")
			return FALSE
		}
		/*
		 *  Case 2: We're outside of the building.
		 *
		 *  Mon Dec 11 09:21:30 2000 -- Scott Turner
		 *
		 *  Permit ship-to-ship razing, if empty.
		 *
		 */
	} else if (subloc(c.who) == subloc(target)) ||
		(is_ship(target) &&
			is_ship(subloc(c.who)) &&
			province(c.who) == province(target)) {
		/*
		 *  The building must be empty.
		 *
		 */
		if len(rp_loc_info(target).here_list) > 0 {
			wout(c.who, "%s is now occupied.", box_name(target))
			return FALSE
		}
	} else {
		wout(c.who, "Must be adjacent to a structure to RAZE it.")
		return FALSE
	}

	per_point = repair_points(int(subkind(target)))
	men = count_man_items(c.who) + 1

	men += c.d
	c.d = men % per_point
	points = men / per_point

	if points > 100 {
		points = 100
	}

	/*
	 *  NOTYET:  first erode defense points before going on to structure
	 *	     damage, as with combat damage against structures?
	 */
	if add_structure_damage(target, points) {
		c.wait = 0
	}
	return TRUE
}

/*
 *  Thu Jan  9 12:14:05 1997 -- Scott Turner
 *
 *  Add fortification to a castle.  Interruptable, etc.  I believe we can use
 *  build_materials, effort_required and effort_given at this point, since
 *  the castle is already complete.
 *
 */
func v_fortify_castle(c *command) int {
	days := c.a
	where := subloc(c.who)
	var workers int
	req_item := 0
	current_defense := loc_defense(where)
	var p *entity_subloc
	var b *entity_build

	if days < 1 {
		days = -1
	}
	c.wait = days

	if subkind(where) != sub_castle {
		wout(c.who, "Only castles may be further fortified.")
		return FALSE
	}

	workers = effective_workers(c.who)

	if workers < 1 {
		wout(c.who, "Need at least one %s.", box_name(item_worker))
		return FALSE
	}

	p = rp_subloc(where)
	assert(p != nil)

	if p.damage > 0 {
		wout(c.who, "You cannot add fortification to a damaged building.")
		return FALSE
	}

	/*
	 *  Try to find the build record.  If nothing, create one.  If there
	 *  is one, then you should just help out...
	 *
	 */
	b = get_build(where, BT_FORTIFY)
	if b != nil {
		wout(c.who, "Renewing fortification of %s.", box_name(where))
		return TRUE
	}

	/*
	 *  Requires a certain amount of materiel.
	 *
	 */
	req_item = current_defense + 1

	if has_item(c.who, item_stone) < req_item {
		wout(c.who, "Adding to the fortification of this castle requires %d stone.", req_item)
		return FALSE
	}

	/*
	 *  Otherwise, begin fortifying.
	 *
	 */
	add_build(where, BT_FORTIFY, 0, 0, 0)
	b = get_build(where, BT_FORTIFY)
	b.effort_required = current_defense * 500
	b.effort_given = 0
	b.build_materials = 0

	wout(c.who, "Begun fortification of %s.", box_name(where))

	return TRUE
}

func d_fortify_castle(c *command) int {
	nworkers := effective_workers(c.who)
	inside := subloc(c.who)
	var p *entity_subloc
	var bonus, fifth, req_stone int
	var effort_given int
	var b *entity_build

	if nworkers <= 0 {
		wout(c.who, "%s has no workers.  Fortification halts.",
			just_name(c.who))
		return FALSE
	}

	p = p_subloc(inside)

	if p.damage > 0 {
		wout(c.who, "You cannot add fortification to a damaged building.")
		return FALSE
	}

	b = get_build(inside, BT_FORTIFY)
	if b == nil {
		wout(c.who, "Fortification at the %d level is complete.", p.defense)
		return FALSE
	}

	/*
	 *  Give a 5% speed bonus for each experience level of the
	 *  construction skill
	 */
	bonus = 5 * c.use_exp * nworkers
	effort_given = nworkers*100 + bonus

	/*
	 *  Materials check
	 *
	 */
	fifth = (b.effort_given + effort_given) * 5 / b.effort_required
	req_stone = loc_defense(inside) / 5

	for fifth < 5 && b.build_materials < fifth {
		if !consume_item(c.who, item_stone, req_stone) {
			wout(c.who, "Need another %s to continue work.  Construction halted.",
				box_name_qty(item_stone, req_stone))
			return FALSE
		}
		b.build_materials++
	}

	b.effort_given += effort_given

	if b.effort_given < b.effort_required {
		return TRUE
	}

	/*
	 *  Otherwise done.
	 *
	 */
	delete_build(inside, BT_FORTIFY)
	p.defense++
	c.wait = 0
	wout(c.who, "Fortification at the %d level is complete.", p.defense)

	return TRUE
}

/*
 *  Thu Jan 16 12:41:48 1997 -- Scott Turner
 *
 *  Strengthen walls -- i.e., add hp to a structure
 *
 */
func v_strengthen_castle(c *command) int {
	days := c.a
	where := subloc(c.who)
	var workers int
	req_item := 0
	current_hp := loc_hp(where)
	var p *entity_subloc
	var b *entity_build

	if days < 1 {
		days = -1
	}
	c.wait = days

	if subkind(where) != sub_castle {
		wout(c.who, "Only castle walls may be strengthened.")
		return FALSE
	}

	workers = effective_workers(c.who)

	if workers < 1 {
		wout(c.who, "Need at least one %s.", box_name(item_worker))
		return FALSE
	}

	p = rp_subloc(where)
	assert(p != nil)

	if p.damage > 0 {
		wout(c.who, "You cannot strengthen a damaged building.")
		return FALSE
	}

	/*
	 *  Try to find the build record.  If nothing, create one.  If there
	 *  is one, then you should just help out...
	 *
	 */
	b = get_build(where, BT_STRENGTHEN)
	if b != nil {
		wout(c.who, "Renewing strengthening of %s.", box_name(where))
		return TRUE
	}

	/*
	 *  Requires a certain amount of materiel.
	 *
	 */
	req_item = current_hp

	if has_item(c.who, item_stone) < req_item {
		wout(c.who, "Strengthening this castle requires %d stone.", req_item)
		return FALSE
	}

	/*
	 *  Otherwise, begin strengthening.
	 *
	 */
	add_build(where, BT_STRENGTHEN, 0, 0, 0)
	b = get_build(where, BT_STRENGTHEN)
	b.effort_required = current_hp * 500
	b.effort_given = 0
	b.build_materials = 0

	wout(c.who, "Begun strengthening %s.", box_name(where))

	return TRUE
}

func d_strengthen_castle(c *command) int {
	nworkers := effective_workers(c.who)
	inside := subloc(c.who)
	var p *entity_subloc
	var bonus, fifth, req_stone int
	var effort_given int
	var b *entity_build

	if nworkers <= 0 {
		wout(c.who, "%s has no workers.  Strengthening halts.",
			just_name(c.who))
		return FALSE
	}

	p = p_subloc(inside)

	if p.damage > 0 {
		wout(c.who, "You cannot strengthen a damaged building.")
		return FALSE
	}

	b = get_build(inside, BT_STRENGTHEN)
	if b == nil {
		wout(c.who, "Walls of %s built to %d strength.",
			box_name(inside), loc_hp(inside))
		return FALSE
	}

	/*
	 *  Give a 5% speed bonus for each experience level of the
	 *  construction skill
	 */
	bonus = 5 * c.use_exp * nworkers
	effort_given = nworkers*100 + bonus

	/*
	 *  Materials check
	 *
	 */
	fifth = (b.effort_given + effort_given) * 5 / b.effort_required
	req_stone = loc_defense(inside) / 5

	for fifth < 5 && b.build_materials < fifth {
		if !consume_item(c.who, item_stone, req_stone) {
			wout(c.who, "Need another %s to continue work.  Construction halted.",
				box_name_qty(item_stone, req_stone))
			return FALSE
		}
		b.build_materials++
	}

	b.effort_given += effort_given

	if b.effort_given < b.effort_required {
		return TRUE
	}

	/*
	 *  Otherwise done.
	 *
	 */
	delete_build(inside, BT_STRENGTHEN)
	if p.hp != loc_hp(inside) {
		p.hp = loc_hp(inside)
	}
	p.hp += 5
	c.wait = 0
	wout(c.who, "Walls of %s built to %d strength.",
		box_name(inside), loc_hp(inside))

	return TRUE
}

/*
 *  Thu Jan 16 14:31:25 1997 -- Scott Turner
 *
 *  Moat prevents missile attacks against a castle.
 *
 */
func v_moat_castle(c *command) int {
	days := c.a
	where := subloc(c.who)
	var workers int
	//req_item := 0;
	var p *entity_subloc
	var b *entity_build

	if days < 1 {
		days = -1
	}
	c.wait = days

	if subkind(where) != sub_castle {
		wout(c.who, "Only castle may have a moat.")
		return FALSE
	}

	workers = effective_workers(c.who)

	if workers < 1 {
		wout(c.who, "Need at least one %s.", box_name(item_worker))
		return FALSE
	}

	p = rp_subloc(where)
	assert(p != nil)

	if p.moat > 0 {
		wout(c.who, "Hmm.  There's already a moat here.  It's old, but looks like it still works.")
		return FALSE
	}

	if p.damage > 0 {
		wout(c.who, "You cannot add a moat to a damaged building.")
		return FALSE
	}

	/*
	 *  Try to find the build record.  If nothing, create one.  If there
	 *  is one, then you should just help out...
	 *
	 */
	b = get_build(where, BT_MOAT)
	if b != nil {
		wout(c.who, "Renewing moating of %s.", box_name(where))
		return TRUE
	}

	/*
	 *  Requires a certain amount of materiel.
	 *
	 */
	if has_item(c.who, item_stone) < (MOAT_MATERIAL / 5) {
		wout(c.who, "Moating this castle requires %d stone.", (MOAT_MATERIAL / 5))
		return FALSE
	}

	/*
	 *  Otherwise, begin moating.
	 *
	 */
	add_build(where, BT_MOAT, 0, 0, 0)
	b = get_build(where, BT_MOAT)
	b.effort_required = MOAT_EFFORT * 100
	b.effort_given = 0
	b.build_materials = 0

	wout(c.who, "Begun adding a moat to %s.", box_name(where))

	return TRUE
}

func d_moat_castle(c *command) int {
	nworkers := effective_workers(c.who)
	inside := subloc(c.who)
	var p *entity_subloc
	var bonus, fifth, req_stone int
	var effort_given int
	var b *entity_build

	if nworkers <= 0 {
		wout(c.who, "%s has no workers.  Moat building halts.",
			just_name(c.who))
		return FALSE
	}

	p = p_subloc(inside)

	if p.damage > 0 {
		wout(c.who, "You cannot add a moat to a damaged building.")
		return FALSE
	}

	b = get_build(inside, BT_MOAT)
	if b == nil {
		if p.moat != FALSE {
			wout(c.who, "Moat added to %s.", box_name(inside))
		}
		return FALSE
	}

	/*
	 *  Give a 5% speed bonus for each experience level of the
	 *  construction skill
	 */
	bonus = 5 * c.use_exp * nworkers
	effort_given = nworkers*100 + bonus

	/*
	 *  Materials check
	 *
	 */
	fifth = (b.effort_given + effort_given) * 5 / b.effort_required
	req_stone = MOAT_MATERIAL / 5

	for fifth < 5 && b.build_materials < fifth {
		if !consume_item(c.who, item_stone, req_stone) {
			wout(c.who, "Need another %s to continue work.  Construction halted.",
				box_name_qty(item_stone, req_stone))
			return FALSE
		}
		b.build_materials++
	}

	b.effort_given += effort_given

	if b.effort_given < b.effort_required {
		return TRUE
	}

	/*
	 *  Otherwise done.
	 *
	 */
	delete_build(inside, BT_MOAT)
	p.moat = 1
	c.wait = 0
	wout(c.who, "Moat added to %s.", box_name(inside))

	return TRUE
}

func v_widen_entrance(c *command) int {
	where := c.a
	var workers int
	//var p *entity_subloc

	if !valid_box(where) {
		wout(c.who, "You must specify a location to widen.")
		return FALSE
	}

	if rp_subloc(where) == nil {
		wout(c.who, "You canot widen that.")
		return FALSE
	}

	if subloc(where) != subloc(c.who) {
		wout(c.who, "You must be outside the location to widen it.")
		return FALSE
	}

	workers = effective_workers(c.who)

	if workers < 10 {
		wout(c.who, "Need at least ten %s.", plural_item_name(item_worker, 10))
		return FALSE
	}

	if FALSE == entrance_size(where) {
		wout(c.who, "You cannot widen that.")
		return FALSE
	}

	wout(c.who, "Begun widening %s.", box_name(where))

	return TRUE
}

func d_widen_entrance(c *command) int {
	where := c.a
	if subloc(where) != subloc(c.who) {
		wout(c.who, "You must be outside the location to widen it.")
		return FALSE
	}

	workers := effective_workers(c.who)
	if workers < 10 {
		wout(c.who, "Need at least ten %s.", plural_item_name(item_worker, 10))
		return FALSE
	}

	if FALSE == entrance_size(where) {
		wout(c.who, "You cannot widen that.")
		return FALSE
	}

	change := workers / 10
	wout(c.who, "Widened entrance by %s.", nice_num(change))
	rp_subloc(where).entrance_size += change

	return TRUE
}
