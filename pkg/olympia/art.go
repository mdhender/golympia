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
 *  Mon Oct 26 18:14:08 1998 -- Scott Turner
 *
 *  Revised Artifact Construction skills, using the new-style artifacts.
 *
 *  Detect Artifact
 *  Reveal Artifact
 *
 *  Mutate Artifact
 *  Destroy Artifact
 *
 *  Conceal Artifact
 *  Obscure Artifact
 *  Remove Obscurity
 *  Disguise Artifact?
 *  Remove Disguise?
 *
 *  Curse Artifact?
 *  Detect Cursed Artifact?
 *  Remove Curse?
 *
 *  Teleport Through Artifact
 *
 */

func has_auraculum(who int) int {
	var ac int

	ac = char_auraculum(who)

	if ac != 0 && has_item(who, ac) > 0 {
		return ac
	}

	return 0
}

/*
 *  Maximum aura, innate plus the auraculum bonus
 */

func max_eff_aura(who int) int {
	var a int  /* aura */
	var ac int /* auraculum */

	a = char_max_aura(who)
	if ac = has_auraculum(who); ac != 0 {
		a += p_item_artifact(ac).param2
	}

	{
		var e *item_ent
		var n int

		for _, e = range inventory_loop(who) {
			if n = int(item_aura_bonus(e.item)); n != 0 {
				a += n
			}
		}
	}

	return a
}

/*
 *  Tue Oct 27 06:55:54 1998 -- Scott Turner
 *
 *  Modified for new-style artifacts.
 *
 */
func v_forge_palantir(c *command) int {
	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	wout(c.who, "Attempt to create a palantir.")
	return TRUE
}

func d_forge_palantir(c *command) int {
	var p *entity_item
	//var pm *item_magic

	newItem := create_unique_item(c.who, sub_magic_artifact)

	if newItem < 0 {
		wout(c.who, "Spell failed.")
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	set_name(newItem, "Palantir")
	p = p_item(newItem)
	p.weight = 2
	p_item_artifact(newItem).type_ = ART_ORB
	p_item_artifact(newItem).param1 = 0
	p_item_artifact(newItem).param2 = 0
	p_item_artifact(newItem).uses = rnd(3, 9)

	wout(c.who, "Created %s.", box_name(newItem))
	log_output(LOG_SPECIAL, "%s created %s.", box_name(c.who), box_name(newItem))
	return TRUE
}

/*
 *  v_use_palantir is replaced with v_art_orb.
 *
 */
func v_use_palantir(c *command) int {
	item := c.a
	target := c.b
	var p *item_magic

	if !is_loc_or_ship(target) {
		wout(c.who, "%s is not a location.", box_code(target))
		return FALSE
	}

	p = rp_item_magic(item)

	if p != nil && p.one_turn_use != FALSE {
		wout(c.who, "The palantir may only be used once per month.")
		return FALSE
	}

	wout(c.who, "Will attempt to view %s with the palantir.",
		box_code(target))

	c.wait = 7

	return TRUE
}

func d_use_palantir(c *command) int {
	item := c.a
	target := c.b

	if !is_loc_or_ship(target) {
		wout(c.who, "%s is not a location.", box_code(target))
		return FALSE
	}

	if loc_shroud(target) != FALSE {
		log_output(LOG_CODE, "Murky palantir result, who=%s, targ=%s",
			box_code_less(c.who), box_code_less(target))
		wout(c.who, "Only murky, indistinct images are seen in the palantir.")
		return FALSE
	}

	log_output(LOG_CODE, "Palantir scry, who=%s, targ=%s",
		box_code_less(c.who), box_code_less(target))

	p_item_magic(item).one_turn_use++

	wout(c.who, "A vision of %s appears:", box_name(target))
	out(c.who, "")
	show_loc(c.who, target)

	alert_palantir_scry(c.who, target)

	return TRUE
}

func v_destroy_art(c *command) int {
	item := c.a

	if has_item(c.who, item) < 1 {
		wout(c.who, "%s does not have %s.",
			box_name(c.who),
			box_code(item))
		return FALSE
	}

	if nil == is_artifact(item) {
		wout(c.who, "Cannot destroy %s with this spell.",
			box_name(item))
		return FALSE
	}

	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	wout(c.who, "Attempt to destroy %s.", box_name(item))
	return TRUE
}

func d_destroy_art(c *command) int {
	item := c.a
	var aura int

	if has_item(c.who, item) < 0 {
		wout(c.who, "%s does not have %s.",
			box_name(c.who),
			box_code(item))
		return FALSE
	}

	if nil == is_artifact(item) {
		wout(c.who, "Cannot destroy %s with this spell.",
			box_name(item))
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	aura = rnd(1, 20)
	if (rp_item_artifact(item).param2&CA_N_MELEE) != 0 ||
		(rp_item_artifact(item).param2&CA_N_MISSILE) != 0 ||
		(rp_item_artifact(item).param2&CA_N_SPECIAL) != 0 ||
		(rp_item_artifact(item).param2&CA_M_MELEE) != 0 ||
		(rp_item_artifact(item).param2&CA_M_SPECIAL) != 0 ||
		(rp_item_artifact(item).param2&CA_N_MELEE_D) != 0 ||
		(rp_item_artifact(item).param2&CA_N_MISSILE_D) != 0 ||
		(rp_item_artifact(item).param2&CA_N_SPECIAL_D) != 0 ||
		(rp_item_artifact(item).param2&CA_M_MELEE_D) != 0 ||
		(rp_item_artifact(item).param2&CA_M_MISSILE_D) != 0 ||
		(rp_item_artifact(item).param2&CA_M_SPECIAL_D) != 0 {
		aura = rp_item_artifact(item).param1 / 5
	}
	if rp_item_artifact(item).type_ == ART_AURACULUM {
		aura = rp_item_artifact(item).param2 / 2
	}
	if rp_item_artifact(item).type_ == ART_ORB {
		aura = rnd(1, 8)
	}

	if aura > 20 {
		aura = 20
	}

	add_aura(c.who, aura)
	wout(c.who, "You gain %s aura from the destroyed artifact.",
		nice_num(char_cur_aura(c.who)))

	log_output(LOG_SPECIAL, "%s destroyed %s.",
		box_name(c.who), box_name(item))

	destroy_unique_item(c.who, item)

	return 0 // todo: should this return something?
}

func v_mutate_art(c *command) int {
	item := c.a

	if has_item(c.who, item) < 1 {
		wout(c.who, "%s does not have %s.",
			box_name(c.who),
			box_code(item))
		return FALSE
	}

	if nil == is_artifact(item) {
		wout(c.who, "Cannot mutate %s with this spell.",
			box_name(item))
		return FALSE
	}

	if (rp_item_artifact(item).type_ == ART_COMBAT) ||
		(rp_item_artifact(item).type_ == ART_AURACULUM) ||
		(rp_item_artifact(item).type_ == ART_ORB) {
		wout(c.who, "%s is not a mutable artifact.",
			box_name(item))
		return FALSE
	}

	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	return TRUE
}

func d_mutate_art(c *command) int {
	item := c.a
	//var aura int

	if has_item(c.who, item) < 0 {
		wout(c.who, "%s does not have %s.",
			box_name(c.who),
			box_code(item))
		return FALSE
	}

	if nil == is_artifact(item) {
		wout(c.who, "Cannot mutate %s with this spell.",
			box_name(item))
		return FALSE
	}

	if rp_item_artifact(item).type_ == ART_COMBAT ||
		rp_item_artifact(item).type_ == ART_AURACULUM ||
		rp_item_artifact(item).type_ == ART_ORB {
		wout(c.who, "%s is not a mutable artifact.",
			box_name(item))
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	newArt := create_random_artifact(c.who)
	wout(c.who, "%s mutates into %s!", box_name(item), box_name(newArt))
	destroy_unique_item(c.who, item)
	return TRUE
}

/*
 *  Tue Oct 27 11:26:28 1998 -- Scott Turner
 *
 *  Conceal artifacts.
 *
 */
func v_conceal_arts(c *command) int {
	target := c.a

	if FALSE == cast_check_char_here(c.who, target) {
		wout(c.who, "Cannot cast on that target.")
		return FALSE
	}

	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	return TRUE
}

func d_conceal_arts(c *command) int {
	target := c.a

	if FALSE == cast_check_char_here(c.who, target) {
		wout(c.who, "Cannot cast on that target.")
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	if FALSE == add_effect(target, ef_conceal_artifacts, 0, 30, 1) {
		wout(c.who, "For some odd reason, your spell fails.")
		return FALSE
	}
	wout(c.who, "%s now has artifacts concealed for 30 days.",
		box_name(target))

	reset_cast_where(c.who)
	return TRUE
}

/*
 *  Tue Oct 27 11:26:28 1998 -- Scott Turner
 *
 *  Reveal artifacts.
 *
 */
func v_reveal_arts(c *command) int {
	target := c.a

	if FALSE == cast_check_char_here(c.who, target) {
		wout(c.who, "Cannot cast on that target.")
		return FALSE
	}

	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	return TRUE
}

func d_reveal_arts(c *command) int {
	target := c.a
	num := 0
	var e *item_ent

	if FALSE == cast_check_char_here(c.who, target) {
		wout(c.who, "Cannot cast on that target.")
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	if get_effect(target, ef_conceal_artifacts, 0, 0) != FALSE {
		wout(c.who, "%s is carrying no artifacts.", box_name(target))
		return TRUE
	}

	for _, e = range inventory_loop(target) {
		if item_unique(e.item) != FALSE &&
			is_artifact(e.item) != nil {
			wout(c.who, "%s is carrying an artifact %s.", box_name(target),
				box_name(e.item))
			num++
		}
	}

	if num == 0 {
		wout(c.who, "%s is carrying no artifacts.", box_name(target))
		return TRUE
	}

	reset_cast_where(c.who)
	return TRUE
}

/*
 *  Tue Oct 27 12:19:41 1998 -- Scott Turner
 *
 *  Deep identify ignores obscurity.
 *
 */
func v_deep_identify(c *command) int {
	target := c.a

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	if !valid_box(target) ||
		nil == is_artifact(target) ||
		FALSE == has_item(c.who, target) ||
		get_effect(target, ef_obscure_artifact, 0, 0) != FALSE {
		wout(c.who, "You are unable to identify that item.")
		return TRUE
	}

	artifact_identify("You study the aura of this artifact and identify it as: ", c)

	return 0 // todo: should this return something?
}

/*
 *  Tue Oct 27 11:26:28 1998 -- Scott Turner
 *
 *  Obscure an artifact.
 *
 */
func v_obscure_art(c *command) int {
	target := c.a

	if FALSE == has_item(c.who, target) {
		wout(c.who, "You must possess an artifact to obscure it.")
		return FALSE
	}

	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	return TRUE
}

func d_obscure_art(c *command) int {
	target := c.a

	if FALSE == has_item(c.who, target) {
		wout(c.who, "You must possess an artifact to obscure it.")
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	if FALSE == add_effect(target, ef_obscure_artifact, 0, -1, 1) {
		wout(c.who, "For some odd reason, your spell fails.")
		return FALSE
	}

	wout(c.who, "%s is now permanently obscured.",
		box_name(target))

	return TRUE
}

/*
 *  Tue Oct 27 11:26:28 1998 -- Scott Turner
 *
 *  Remove an obscurity.
 *
 */
func v_unobscure_art(c *command) int {
	target := c.a

	if FALSE == has_item(c.who, target) {
		wout(c.who, "You must possess an artifact to remove an obscurity.")
		return FALSE
	}

	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	return TRUE
}

func d_unobscure_art(c *command) int {
	target := c.a

	if FALSE == has_item(c.who, target) {
		wout(c.who, "You must possess an artifact to remove an obscurity.")
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	delete_effect(target, ef_obscure_artifact, 0)

	wout(c.who, "Remove obscurity cast upon %s.",
		box_name(target))

	return TRUE
}

/*
 *  Tue Oct 27 11:34:53 1998 -- Scott Turner
 *
 *  This function is needed for detecting artifacts.
 *  Ignore concealed artifacts.
 */
func find_nearest_artifact(who int) int {
	distance := 9999
	var d, i int
	where := province(who)

	if region(where) == faery_region ||
		region(where) == hades_region ||
		region(where) == cloud_region {
		return -1
	}

	for _, i = range loop_artifact() {
		if region(item_unique(i)) == faery_region ||
			region(item_unique(i)) == hades_region ||
			region(item_unique(i)) == cloud_region {
			continue
		}
		/*
		 *  Might be concealed.
		 *
		 */
		if item_unique(i) != FALSE &&
			get_effect(item_unique(i), ef_conceal_artifacts, 0, 0) != FALSE {
			continue
		}
		d = los_province_distance(item_unique(i), where)
		if d < distance {
			distance = d
		}
	}

	return distance
}

/*
 *  Tue Oct 27 11:26:28 1998 -- Scott Turner
 *
 *  Detect artifacts.
 *
 */
func v_detect_arts(c *command) int {
	if region(c.who) == faery_region ||
		region(c.who) == hades_region ||
		region(c.who) == cloud_region {
		wout(c.who, "Your magic does not work in this place.")
		return FALSE
	}

	if FALSE == check_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}

	return TRUE
}

func d_detect_arts(c *command) int {
	var distance int

	if region(c.who) == faery_region ||
		region(c.who) == hades_region ||
		region(c.who) == cloud_region {
		wout(c.who, "Your magic does not work in this place.")
		return FALSE
	}

	if FALSE == charge_aura(c.who, skill_aura(c.use_skill)) {
		return FALSE
	}
	wout(c.who, "Used %s aura casting this spell.", nice_num(skill_aura(c.use_skill)))

	distance = find_nearest_artifact(c.who)
	if distance == -1 {
		wout(c.who, "The nearest artifact is very distant.")
	} else {
		wout(c.who, "The nearest artifact is %s province%s away.",
			nice_num(distance), or_string(distance == 1, "", "s"))
	}

	return TRUE
}

func v_forge_aura(c *command) int {
	var aura int

	if char_auraculum(c.who) != FALSE {
		wout(c.who, "%s may only be used once.",
			box_name(c.use_skill))
		return FALSE
	}

	if c.a < 1 {
		c.a = 1
	}
	aura = c.a

	if FALSE == check_aura(c.who, aura) {
		return FALSE
	}

	if aura > char_max_aura(c.who) {
		wout(c.who, "The specified amount of aura exceeds the maximum aura level of %s.", box_name(c.who))
		return FALSE
	}

	wout(c.who, "Attempt to forge an auraculum.")
	return TRUE
}

func notify_others_auraculum(who, item int) {
	var n int

	for _, n = range loop_char() {
		if n != who && is_magician(n) != FALSE && has_auraculum(n) != FALSE {
			wout(n, "%s has constructed an auraculum.",
				box_name(who))
		}
	}

	log_output(LOG_SPECIAL, "%s created %s, %s.",
		box_name(who),
		box_name(item),
		subkind_s[subkind(item)])
}

/*
 *  Mon Oct 26 10:31:46 1998 -- Scott Turner
 *
 *  Modified to create a new-style artifact.
 *
 */
func d_forge_aura(c *command) int {
	aura := c.a
	var new_name string
	var p *entity_item
	var cm *char_magic

	if aura > char_max_aura(c.who) {
		wout(c.who, "The specified amount of aura exceeds the maximum aura level of %s.", box_name(c.who))
		return FALSE
	}

	if FALSE == charge_aura(c.who, aura) {
		return FALSE
	}

	if numargs(c) < 2 {
		switch rnd(1, 3) {
		case 1:
			new_name = "Gold ring"
			break
		case 2:
			new_name = "Wooden staff"
			break
		case 3:
			new_name = "Jeweled crown"
			break
		default:
			panic("!reached")
		}
	} else {
		new_name = string(c.parse[2])
	}

	newItem := create_unique_item(c.who, sub_magic_artifact)

	if newItem < 0 {
		wout(c.who, "Spell failed.")
		return FALSE
	}

	set_name(newItem, new_name)

	p = p_item(newItem)
	p.weight = rnd(1, 3)

	p_item_artifact(newItem).type_ = ART_AURACULUM
	p_item_artifact(newItem).param1 = c.who    /* creator */
	p_item_artifact(newItem).param2 = aura * 2 /* aura */

	cm = p_magic(c.who)
	cm.auraculum = newItem
	cm.max_aura -= aura

	wout(c.who, "Created %s.", box_name(newItem))
	notify_others_auraculum(c.who, newItem)

	learn_skill(c.who, sk_adv_sorcery)

	return TRUE
}

func v_forge_art_x(c *command) int {
	aura := c.a
	var rare_item int

	if aura < 1 {
		c.a = 1
		aura = 1
	}
	if aura > 20 {
		c.a = 1
		aura = 20
	}

	if FALSE == check_aura(c.who, aura) {
		return FALSE
	}

	if FALSE == can_pay(c.who, 500) {
		wout(c.who, "Requires %s.", gold_s(500))
		return FALSE
	}

	switch c.use_skill {
	case sk_forge_weapon:
	case sk_forge_armor:
		rare_item = item_mithril
		break

	case sk_forge_bow:
		rare_item = item_mallorn_wood
		break

	default:
		panic("!reached")
	}
	c.d = rare_item

	if !(has_item(c.who, rare_item) >= 1) { // todo: should be !(... >= 1), maybe?
		wout(c.who, "Requires %s.", box_name_qty(rare_item, 1))
		return FALSE
	}

	return TRUE
}

func d_forge_art_x(c *command) int {
	aura := c.a
	rare_item := c.d
	var new_name string
	//var pm *item_magic

	if FALSE == check_aura(c.who, aura) {
		return FALSE
	}

	if FALSE == charge(c.who, 500) {
		wout(c.who, "Requires %s.", gold_s(500))
		return FALSE
	}

	if !(has_item(c.who, rare_item) >= 1) { // todo: should be !(... >= 1), maybe?
		wout(c.who, "Requires %s.", box_name_qty(rare_item, 1))
		return FALSE
	}

	charge_aura(c.who, aura)
	charge(c.who, 500)
	consume_item(c.who, rare_item, 1)

	newItem := create_unique_item(c.who, sub_magic_artifact)
	p_item(newItem).weight = 10
	p_item_artifact(newItem).type_ = ART_COMBAT

	switch c.use_skill {
	case sk_forge_weapon:
		rp_item_artifact(newItem).param2 = CA_N_MELEE
		rp_item_artifact(newItem).param1 = aura * 5
		new_name = "enchanted sword"
		break

	case sk_forge_armor:
		rp_item_artifact(newItem).param2 = CA_N_MELEE_D | CA_N_MISSILE_D | CA_N_SPECIAL_D
		rp_item_artifact(newItem).param1 = aura * 5
		new_name = "enchanted armor"
		break

	case sk_forge_bow:
		rp_item_artifact(newItem).param2 = CA_N_MISSILE
		rp_item_artifact(newItem).param1 = aura * 5
		new_name = "enchanted bow"
		break

	default:
		panic("!reached")
	}

	if numargs(c) >= 2 && len(c.parse[2]) != 0 {
		new_name = string(c.parse[2])
	}

	set_name(newItem, new_name)
	wout(c.who, "Created %s.", box_name(newItem))

	return TRUE
}

func new_suffuse_ring(who int) int {
	//var ni int
	//var lore int

	newItem := create_unique_item(who, sub_magic_artifact)
	set_name(newItem, "Golden ring")
	p_item(newItem).weight = 1
	p_item_artifact(newItem).type_ = ART_DESTROY
	rp_item_artifact(newItem).uses = 1
	rp_item_artifact(newItem).param1 = random_beast(0)
	return newItem
}

/* Temporary */

var orb_used_this_month []int

func v_use_orb(c *command) int {
	item := c.a
	target := c.b
	where := 0
	var owner int
	var p *item_magic

	if ilist_lookup(orb_used_this_month, item) >= 0 {
		wout(c.who, "The orb may only be used once per month.")
		wout(c.who, "Only murky, indistinct images are seen in the orb.")
		return FALSE
	}

	orb_used_this_month = append(orb_used_this_month, item)

	if rnd(1, 3) == 1 {
		wout(c.who, "Only murky, indistinct images are seen in the orb.")
		return FALSE
	}

	switch kind(target) {
	case T_loc, T_ship:
		where = province(target)
		break

	case T_char:
		where = province(target)
		break

	case T_item:
		if owner = item_unique(target); owner != 0 {
			where = province(owner)
		}
		break
	}

	if where == 0 {
		wout(c.who, "The orb is unsure what location is meant to be scried.")
	} else if loc_shroud(where) != FALSE {
		wout(c.who, "The orb is unable to penetrate a shroud over %s.",
			box_name(where))
	} else {
		wout(c.who, "A vision of %s appears:", box_name(where))
		show_loc(c.who, where)
		alert_scry_generic(c.who, where)
	}

	p = p_item_magic(item)

	p.orb_use_count--
	if p.orb_use_count <= 0 {
		wout(c.who, "After the vision fades, the orb grows dark, and shatters.  The orb is gone")
		destroy_unique_item(c.who, item)
	}

	return TRUE
}
