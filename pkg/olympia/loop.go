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

type LOOPCTL int

const (
	LOOP_BREAK LOOPCTL = iota
	LOOP_CONTINUE
)

//
//  loop.h -- abstracted loops
//
//  or, "Too bad this language doesn't have generators"
//
//
//  Q:  Why abstract loops?  These defines are really gross.
//
//  A:  We use abstract data types so that we can change the representation
//      without having to change all of the code that uses the types.
//
//      It is easy to abstract add, delete and fetch operations.  But the
//      most common operation is to iterate over the elements of a collection.
//
//      The method of iteration is almost certain to change when switching
//      implementations, too.  List-to-tree, tree-to-bit-array, etc.
//
//      Abstracting loops makes the code cleaner and easier to read, and
//      easier to change.

//  Loops below should be free of the "delete problem", i.e.
//
//	loop_something(i)
//	{
//		delete_something(i);
//	}
//	next_something;
//
//  should work, i.e. it shouldn't core dump because next(i) is no
//  longer defined, or go over an element twice, or miss an element.

//  break and continue should work inside these loops, but don't
//  return out of them.  Return will bypass the end-of-loop cleanup.

//  Thanks to the X Window system for going ahead of us and making sure
//  that most vendor's compilers can handle unreasonably large defines.

// #define    loop_kind(kind, i)
// { int ll_i, ll_next;
//
//	int ll_check = 5;
//	ll_next = kind_first(kind);
//	while ((ll_i = ll_next) > 0) {
//	  ll_next = kind_next(ll_i);
//	  i = ll_i;
//
// #define    next_kind    } assert(ll_check == 5); }
func loop_kind(kind int) []int {
	var ll_l []int
	for ll_next := kind_first(kind); ll_next > 0; ll_next = kind_next(ll_next) {
		ll_l = append(ll_l, ll_next)
	}
	return ll_l
}

// loop_nation(i) loop_kind(T_nation, i)
func loop_nation() []int {
	return loop_kind(T_nation)
}

// #define    loop_subkind(sk, i) \
// int ll_i, ll_next
// int ll_check = 26
// ll_next = sub_first(sk)
//
//	 while ((ll_i = ll_next) > 0) {
//	   ll_next = sub_next(ll_i);
//	   i = ll_i;
//
//	#define    next_subkind    } assert(ll_check == 26); }
func loop_subkind(sk int) []int {
	var ll_l []int
	for ll_next := sub_first(sk); ll_next > 0; ll_next = sub_next(ll_next) {
		ll_l = append(ll_l, ll_next)
	}
	return ll_l
}

/*
#define    loop_all_here(where, i) \
{ int ll_i;

	int ll_check = 2;
	ilist ll_l = NULL;
	all_here(where, &ll_l);
	for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) {
	  i = ll_l[ll_i];

#define    next_all_here        } assert(ll_check == 2); ilist_reclaim(&ll_l); }
*/
func loop_all_here(where int) []int {
	return all_here(where, nil)
}

/*
#define    loop_char_here(where, i) \

	{ int ll_i; \
	  int ll_check = 13; \
	  ilist ll_l = NULL; \
	  all_char_here(where, &ll_l); \
	  for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
	    i = ll_l[ll_i];

#define    exit_char_here    { assert(ll_check == 13); ilist_reclaim(&ll_l); }
#define    next_char_here    } assert(ll_check == 13); ilist_reclaim(&ll_l); }
*/
func loop_char_here(where int) []int {
	return all_char_here(where, nil)
}

/*
#define    loop_stack(who, i) \

	{ int ll_i; \
	  int ll_check = 20; \
	  ilist ll_l = NULL; \
	    all_stack(who, &ll_l); \
	    for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
	    { \
	        i = ll_l[ll_i];

#define    next_stack    } assert(ll_check == 20); ilist_reclaim(&ll_l); }
*/
func loop_stack(who int) []int {
	var ll_l []int
	for _, el := range all_stack(who, nil) {
		ll_l = append(ll_l, el)
	}
	return ll_l
}

/*
#define known_entity_player_loop(ep, i) \
{ int ll_check = 3; \
  qsort(ep, entity_player_list_len(kn), sizeof(int), int_comp); \
  for (int ll_i = 0; ll_i < entity_player_list_len(ep); ll_i++) { \
    (i) = (ep)[ll_i];

#define known_entity_player_next \
      } assert(ll_check == 3); }
*/

/*
#define known_sparse_loop(kn, i) \
{ int ll_check = 3; \
  qsort(kn, ilist_len(kn), sizeof(int), int_comp); \
  for (int ll_i = 0; ll_i < ilist_len(kn); ll_i++) { \
    (i) = (kn)[ll_i];

#define known_sparse_next \
      } assert(ll_check == 3); }
*/
// todo: fix bug with sparse vs entity_player *
func known_sparse_loop(kn []int) []int {
	// copy the list
	var ll_l []int
	for _, i := range kn {
		ll_l = append(ll_l, i)
	}
	sort.Ints(ll_l)
	return ll_l
}

/*
//#define    loop_known(kn, i) \
//{ int ll_check = 3; \
//    qsort(kn, entity_player_list_len(kn), sizeof(int), int_comp); \
//    for (int ll_i = 0; ll_i < entity_player_list_len(kn); ll_i++) { \
//        (i) = (kn)[ll_i];
//
//#define    next_known    } assert(ll_check == 3); }
*/

/*
//
//  Iterate over all valid boxes.  i is instantiated with the entity
//  numbers.
///

#define    loop_boxes(i) \

	{ int ll_i; \
	  int ll_check = 4; \
	    for (ll_i = 1; ll_i < MAX_BOXES; ll_i++) \
	        if (kind(ll_i) != T_deleted) \
	        { \
	            i = ll_i;

#define    next_box    } assert(ll_check == 4); }
*/
func loop_boxes() []int {
	var ll_l []int
	for i := 1; i < MAX_BOXES; i++ {
		if kind(i) != T_deleted {
			ll_l = append(ll_l, i)
		}
	}
	return ll_l
}

/*
#define    loop_char(i)    loop_kind(T_char, i)
#define    next_char    next_kind
*/
func loop_char() []int {
	return loop_kind(T_char)
}

/*
#define    loop_player(i)    loop_kind(T_player, i)
#define    next_player    next_kind
*/
func loop_player() []int {
	return loop_kind(T_player)
}

/*
#define    loop_loc(i)    loop_kind(T_loc, i)
#define    next_loc    next_kind
*/
func loop_loc() []int {
	return loop_kind(T_loc)
}

/*
#define    loop_item(i)    loop_kind(T_item, i)
#define    next_item    next_kind
*/
func loop_item() []int {
	return loop_kind(T_item)
}

/*
#define    loop_exit(i)    loop_kind(T_exit, i)
#define    next_exit    next_kind
*/
func loop_exit() []int {
	//return loop_kind(T_exit)
	panic("!implemented")
}

/*
#define    loop_skill(i)    loop_kind(T_skill, i)
#define    next_skill    next_kind
*/
func loop_skill() []int {
	return loop_kind(T_skill)
}

/*
#define    loop_gate(i)    loop_kind(T_gate, i)
#define    next_gate    next_kind
*/
func loop_gate() []int {
	return loop_kind(T_gate)
}

/*
#define    loop_ship(i)    loop_kind(T_ship, i)
#define    next_ship    next_kind
*/
func loop_ship() []int {
	return loop_kind(T_ship)
}

/*
#define    loop_post(i)    loop_kind(T_post, i)
#define    next_post    next_kind
*/
func loop_post() []int {
	return loop_kind(T_post)
}

/*
#define    loop_storm(i)    loop_kind(T_storm, i)
#define    next_storm    next_kind
*/
func loop_storm() []int {
	return loop_kind(T_storm)
}

// loop_castle(i)  loop_subkind(sub_castle, i)
func loop_castle() []int {
	return loop_subkind(sub_castle)
}

/*
#define    loop_garrison(i)    loop_subkind(sub_garrison, i)
#define    next_garrison        next_subkind
*/
func loop_garrison() []int {
	return loop_subkind(sub_garrison)
}

/*
#define    loop_city(i)        loop_subkind(sub_city, i)
#define    next_city        next_subkind
*/
func loop_city() []int {
	return loop_subkind(sub_city)
}

/*
#define    loop_guild(i)        loop_subkind(sub_guild, i)
#define    next_guild        next_subkind
*/
func loop_guild() []int {
	return loop_subkind(sub_guild)
}

/*
#define    loop_mountain(i)    loop_subkind(sub_mountain, i)
#define    next_mountain        next_subkind
*/
func loop_mountain() []int {
	return loop_subkind(sub_mountain)
}

/*
#define    loop_inn(i)        loop_subkind(sub_inn, i)
#define    next_inn        next_subkind
*/
func loop_inn() []int {
	return loop_subkind(sub_inn)
}

/*
#define    loop_temple(i)        loop_subkind(sub_temple, i)
#define    next_temple        next_subkind
*/
func loop_temple() []int {
	return loop_subkind(sub_temple)
}

/*
#define    loop_collapsed_mine(i)    loop_subkind(sub_mine_collapsed, i)
#define    next_collapsed_mine    next_subkind
*/
func loop_collapsed_mine() []int {
	return loop_subkind(sub_mine_collapsed)
}
func loop_mine_collapsed() []int {
	return loop_subkind(sub_mine_collapsed)
}

/*
#define    loop_dead_body(i)    loop_subkind(sub_dead_body, i)
#define    next_dead_body        next_subkind
*/
func loop_dead_body() []int {
	return loop_subkind(sub_dead_body)
}

/*
#define    loop_lost_soul(i)    loop_subkind(sub_lost_soul, i)
#define    next_lost_soul        next_subkind
*/
func loop_lost_soul() []int {
	return loop_subkind(sub_lost_soul)
}

/*
#define    loop_pl_regular(i)    loop_subkind(sub_pl_regular, i)
#define    next_pl_regular        next_subkind
*/
func loop_pl_regular() []int {
	return loop_subkind(sub_pl_regular)
}

/*
#define    loop_artifact(i)    loop_subkind(sub_magic_artifact, i)
#define    next_artifact        next_subkind
*/
func loop_artifact() []int {
	return loop_subkind(sub_magic_artifact)
}

/*
#define    loop_loc_or_ship(i) \

	{ int ll_i; \
	  int ll_check = 17; \
	  int ll_state = 1; \
	    ll_i = kind_first(T_ship); \
	    if (ll_i <= 0) { ll_i = kind_first(T_loc); ll_state = 0; } \
	    while (ll_i > 0) { \
	        i = ll_i;

	#define    next_loc_or_ship \
	    ll_i = kind_next(ll_i); \
	    if (ll_i <= 0 && ll_state) \
	     { ll_i = kind_first(T_loc); ll_state = 0; } \
	    } assert(ll_check == 17); }
*/
func loop_loc_or_ship() []int {
	panic("!implemented")
}

/*
#define    loop_province(i) \

	{ int ll_i; \
	  int ll_check = 6; \
	    for (ll_i = kind_first(T_loc); ll_i > 0; ll_i = kind_next(ll_i)) \
	        if (loc_depth(ll_i) == LOC_province) { \
	            i = ll_i;

#define    next_province    } assert(ll_check == 6); }
*/
func loop_province() []int {
	var ll_l []int
	for _, i := range loop_kind(T_loc) {
		if loc_depth(i) == LOC_province {
			ll_l = append(ll_l, i)
		}
	}
	return ll_l
}

/*
#define    loop_loc_teach(where, i) \

	{ int ll_i; \
	  int ll_check = 7; \
	  ilist ll_l = NULL; \
	    assert(valid_box(where)); \
	    if (rp_subloc(where)) \
	        ll_l = ilist_copy(rp_subloc(where)->teaches); \
	        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
	            i = ll_l[ll_i];

#define    next_loc_teach    } assert(ll_check == 7); ilist_reclaim(&ll_l); }
*/
func loop_loc_teach(where int) []int {
	var ll_l []int
	if !valid_box(where) {
		panic("assert(valid_box(where))")
	}
	if rp_subloc(where) != nil {
		for _, e := range rp_subloc(where).teaches {
			ll_l = append(ll_l, e)
		}
	}
	return ll_l
}

/*
#define    loop_units(pl, i) \
{ int ll_i; \
  ilist ll_l = NULL; \
  int ll_check = 21; \
    if (rp_player(pl)) \
        ll_l = ilist_copy(rp_player(pl)->units); \
        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
            i = ll_l[ll_i];

#define    next_unit    } assert(ll_check == 21); ilist_reclaim(&ll_l); }
*/

func loop_units(pl int) []int {
	var ll_l []int
	if rp_player(hades_player) != nil {
		ll_l = append(ll_l, rp_player(hades_player).units...)
	}
	return ll_l
}

/*
#define    loop_here(where, i) \

	{ int ll_i; \
	  ilist ll_l = NULL; \
	  int ll_check = 8; \
	    assert(valid_box(where)); \
	    if (rp_loc_info(where)) \
	        ll_l = ilist_copy(rp_loc_info(where)->here_list); \
	        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
	            i = ll_l[ll_i];

#define exit_here    { assert(ll_check == 8); ilist_reclaim(&ll_l); }
#define    next_here    } assert(ll_check == 8); ilist_reclaim(&ll_l); }
*/
func loop_here(where int) []int {
	if !valid_box(where) {
		panic("assert(valid_box(where))")
	}
	if rp_loc_info(where) == nil {
		return nil
	}
	var ll_l []int
	for _, l := range rp_loc_info(where).here_list {
		ll_l = append(ll_l, l)
	}
	return ll_l
}

/*
#define    loop_gates_here(where, i) \

	{ int ll_i; \
	  ilist ll_l = NULL; \
	  int ll_check = 18; \
	    assert(valid_box(where)); \
	    if (rp_loc_info(where)) \
	        ll_l = ilist_copy(rp_loc_info(where)->here_list); \
	        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
	            if (kind(ll_l[ll_i]) == T_gate) { \
	                i = ll_l[ll_i];

#define    next_gate_here    } assert(ll_check == 18); ilist_reclaim(&ll_l); }
*/
func loop_gates_here(where int) []int {
	if !valid_box(where) {
		panic("assert(valid_box(where))")
	}
	var ll_l []int
	if rp_loc_info(where) != nil {
		for _, i := range rp_loc_info(where).here_list {
			if kind(i) == T_gate {
				ll_l = append(ll_l, i)
			}
		}
	}
	return ll_l
}

/*
#define    loop_exits_here(where, i) \
{ int ll_i; \
  int ll_check = 10; \
  ilist ll_l = NULL; \
    assert(valid_box(where)); \
    if (rp_loc(where)) \
        ll_l = ilist_copy(rp_loc(where)->exits_here); \
        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
            i = ll_l[ll_i];

#define    next_exit_here    } assert(ll_check == 10); ilist_reclaim(&ll_l); }
*/

/*
//
//  Iterate struct item_ent//e over who's inventory
///

#define inventory_loop(who, e) \

	{   int ie_check = 11; \
	    assert(valid_box(who)); \
	    struct item_ent **ie_l = ie_list_copy(bx[who]->items); \
	    for (int ie_i = 0; ie_i < ie_list_len(ie_l); ie_i++) { \
	        if (valid_box(ie_l[ie_i]->item) && ie_l[ie_i]->qty > 0) { \
	            struct item_ent ie_copy = **ie_l[ie_i]; \
	            (e) = ie_l[ie_i];

	#define inventory_next \
	        } \
	    } \
	    assert(ie_check == 11); \
	    ie_list_reclaim(&ie_l); \
	}
*/
func inventory_loop(who int) []*item_ent {
	return loop_inventory(who)
}
func loop_inventory(who int) []*item_ent {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	var ll_l []*item_ent
	for _, e := range bx[who].items {
		if valid_box(e.item) && e.qty > 0 {
			ll_l = append(ll_l, e)
		}
	}
	return ll_l
}

/*
#define    loop_inv(who, e) \
{ int ll_i; \
 int ll_check = 11; \
 struct item_ent ll_copy; \
 struct item_ent **ll_l = NULL; \
   assert(valid_box(who)); \
   ll_l = (struct item_ent **) ilist_copy((ilist) bx[who]->items); \
   for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) \
       if (valid_box(ll_l[ll_i]->item) && ll_l[ll_i]->qty > 0) { \
           ll_copy =//ll_l[ll_i]; \
           e = &ll_copy;

#define    next_inv   } assert(ll_check == 11); ilist_reclaim((ilist//) &ll_l); }
*/

/*
	#define    loop_char_skill(who, e) \
		{ int ll_check = 15;               \
		  struct skill_ent **ll_l = 0;     \
		  assert(valid_box(who)); \
		  if (rp_char(who)) {ll_l = skill_ent_list_copy(rp_char(who)->skills);} \
		    for (int ll_i = 0; ll_i < skill_ent_list_len(ll_l); ll_i++) { \
		      (e) = ll_l[ll_i];

#define    next_char_skill } assert(ll_check == 15); skill_ent_list_reclaim(&ll_l); }
*/
func loop_char_skill(who int) []*skill_ent {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	var ll_l []*skill_ent
	if rp_char(who) != nil {
		for _, e := range rp_char(who).skills {
			ll_l = append(ll_l, e)
		}
	}
	return ll_l
}

/*
#define    loop_char_skill_known(who, e) \

	{ int ll_check = 16; \
	  struct skill_ent **ll_l = NULL; \
	  assert(valid_box(who)); \
	    if (rp_char(who)) \
	       ll_l = skill_ent_list_copy(rp_char(who)->skills); \
	       for (int ll_i = 0; ll_i < skill_ent_list_len(ll_l); ll_i++) \
	          if (ll_l[ll_i]->know == SKILL_know) { \
	            (e) = ll_l[ll_i];

#define    exit_char_skill_known { assert(ll_check == 16); skill_ent_list_reclaim(&ll_l); }
#define    next_char_skill_known } assert(ll_check == 16); skill_ent_list_reclaim(&ll_l); }
*/
func loop_char_skill_known(who int) []*skill_ent {
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	var ll_l []*skill_ent
	if rp_char(who) != nil {
		for _, se := range rp_char(who).skills {
			if se.know == SKILL_know {
				ll_l = append(ll_l, se)
			}
		}
	}
	return ll_l
}

/*
#define trade_loop(who, e) \
{ int tr_check = 19; \
  assert(valid_box(who)); \
  struct trade//*tr_l = tr_list_copy(bx[who]->trades); \
    for (int tr_i = 0; tr_i < tr_list_len(tr_l); tr_i++) { \
        if (valid_box(tr_l[tr_i]->item) && tr_l[tr_i]->qty > 0) { \
            (e) = tr_l[tr_i];

#define trade_next \
        } \
    } \
    assert(tr_check == 19); \
    tr_list_reclaim(&tr_l); \
}
*/

/*
#define    loop_trade(who, e) \

	{ int ll_check = 19; \
	  assert(valid_box(who)); \
	  struct trade **ll_l = tr_list_copy(bx[who]->trades); \
	    for (int ll_i = 0; ll_i < tr_list_len(ll_l); ll_i++) \
	        if (valid_box(ll_l[ll_i]->item) && ll_l[ll_i]->qty > 0) { \
	            (e) = ll_l[ll_i];

#define    next_trade  } assert(ll_check == 19); tr_list_reclaim(&ll_l); }
*/
func loop_trade(who int) []*trade {
	panic("!implemented")
}

/*
#define    loop_prov_dest(where, i) \

	{ int ll_i; \
	  int ll_check = 23; \
	  struct entity_loc//ll_p; \
	    assert(loc_depth(where) == LOC_province); \
	    ll_p = rp_loc(where); \
	    assert(ll_p); \
	    for (ll_i = 0; ll_i < ilist_len(ll_p->prov_dest); ll_i++) \
	    { \
	        i = ll_p->prov_dest[ll_i];

#define    next_prov_dest        } assert(ll_check == 23); }
*/
func loop_prov_dest(where int) []int {
	assert(loc_depth(where) == LOC_province)
	loc := rp_loc(where)
	assert(loc != nil)
	var ll_l []int
	for _, e := range loc.prov_dest {
		ll_l = append(ll_l, e)
	}
	return ll_l
}

/*
#if 0
#define	loop_loc_owner(where, i) \
{ int ll_i, ll_next; \
  int ll_check = 25; \
    ll_next = province_admin(where); \
    while ((ll_i = ll_next) > 0) { \
        ll_next = char_pledge(ll_i); \
        i = ll_i;

#define	next_loc_owner	} assert(ll_check == 25); }
#endif
*/

/*
//
//  This loop sorts the player's inventory according to the
//  attack && defense values, so that we pick up the toughest
//  fighters first when going through the inventory.
//
//  Note that it really belongs in loop.h.
//
///
#define    loop_sorted_inv(who, e) \

	{ int ll_check = 27; \
	    assert(valid_box(who)); \
	    struct item_ent//*ll_l = ie_list_copy(bx[who]->items); \
	    qsort(ll_l, ie_list_len(ll_l), sizeof(int), combat_comp); \
	    for (int ll_i = 0; ll_i < ie_list_len(ll_l); ll_i++) \
	        if (valid_box(ll_l[ll_i]->item) && ll_l[ll_i]->qty > 0) { \
	            struct item_ent ll_copy =//ll_l[ll_i]; \
	            e = &ll_copy;

#define    next_sorted_inv   } assert(ll_check == 27); ie_list_reclaim(&ll_l); }
*/
func loop_sorted_inv(who int) []*item_ent {
	var l []*item_ent
	if !valid_box(who) {
		panic("assert(valid_box(who))")
	}
	for _, e := range bx[who].items {
		if valid_box(e.item) && e.qty > 0 {
			l = append(l, e)
		}
	}
	sort.Slice(l, func(i, j int) bool {
		return combat_comp(l[i], l[j]) < 0
	})
	return l
}

/*
//
//  Loop over a priest's followers.
//
///
#define    loop_followers(priest, i) \

	{ int ll_i; \
	  ilist ll_l = NULL; \
	  int ll_check = 28; \
	    if (rp_char(priest) && is_priest(priest)) \
	        ll_l = ilist_copy(rp_char(priest)->religion.followers); \
	        for (ll_i = 0; ll_i < ilist_len(ll_l); ll_i++) { \
	            i = ll_l[ll_i];

#define    next_follower    } assert(ll_check == 28); ilist_reclaim(&ll_l); }
*/
func loop_followers(priest int) []int {
	panic("!implemented")
}

/*
 *  Loop over the effects on something.
 *
 *  What = int
 *  e = struct effect *;
 *
  if (rp_char(what)) { \
    el = (ilist) rp_char(what).effects; \
  } else if (rp_loc(what)) { \
    el = (ilist) rp_loc(what).effects; \
  } else if (rp_subloc(what)) { \
    el = (ilist) rp_subloc(what).effects; \
  }; \
 *
*/
//#define    loop_effects(what, e) \
//{ int ll_check = 29; \
//  effect_list_t el = effects(what); \
//  effect_list_t ll_l = effect_list_copy(el); \
//  for (int ll_i = 0; ll_i < effect_list_len(ll_l); ll_i++) {  \
//    (e) = ll_l[ll_i];
//#define    next_effect    } assert(ll_check == 29); effect_list_reclaim(&ll_l); }
func loop_effects(what int) []*effect {
	return effects(what)
}
