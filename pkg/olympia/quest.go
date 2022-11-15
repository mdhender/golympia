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

const (
	/*
	 *  Tue Oct 13 13:14:57 1998 -- Scott Turner
	 *
	 *  This borrows the old "monstermark" idea to grade the toughness of
	 *  a monster stack, to let us give it appropriate treasure.
	 *
	 *  Tue Apr 27 06:40:02 1999 -- Scott Turner
	 *
	 *  MM provides a monstermark for a beast type, monstermark() for a stack.
	 *
	 *  Depending upon how MM is calculated, you can vary the relative numbers
	 *  of beasts.  If MM = att * def, then a peasant gets a 1 and a dragon gets
	 *  2250000 -- which means that dragons are going to be essentially unknown.
	 *  If MM = att + def, then the range is 0 to 3000, which perhaps makes more
	 *  sense.
	 *
	 *  MAX_MM is a calculated maximum MM (in io.c)
	 *
	 */
	MAX_MM = 1000
)

var (
	art_att_s = []string{"sword", "dagger", "longsword"}
	art_def_s = []string{"helm", "shield", "armor"}
	art_mis_s = []string{"spear", "bow", "javelin", "dart"}
	art_mag_s = []string{"ring", "staff", "amulet"}
	of_names  = []string{
		"Achilles", "Darkness", "Justice", "Truth", "Norbus", "Dirbrand",
		"Pyrellica", "Halhere", "Eadak", "Faelgrar", "Napich", "Renfast",
		"Ociera", "Shavnor", "Dezarne", "Roshun", "Areth Lorbin", "Anarth",
		"Vernolt", "Pentara", "Gravecarn", "Sardis", "Lethrys", "Habyn",
		"Obraed", "Beebas", "Bayarth", "Haim", "Balatea", "Bobbiek", "Moldarth",
		"Grindor", "Sallen", "Ferenth", "Rhonius", "Ragnar", "Pallia", "Kior",
		"Baraxes", "Coinbalth", "Raskold", "Lassan", "Haemfrith", "Earnberict",
		"Sorale", "Lorbin", "Osgea", "Fornil", "Kuneack", "Davchar", "Urvil",
		"Pantarastar", "Cyllenedos", "Echaliatic", "Iniera", "Norgar", "Broen",
		"Estbeorn", "Claunecar", "Salamus", "Rhovanth", "Illinod", "Pictar",
		"Elakain", "Antresk", "Kichea", "Raigor", "Pactra", "Aethelarn",
		"Descarq", "Plagcath", "Nuncarth", "Petelinus", "Cospera", "Sarindor",
		"Albrand", "Evinob", "Dafarik", "Haemin", "Resh", "Tarvik", "Odasgunn",
		"Areth Pirn", "Miranth", "Dorenth", "Arkaune", "Kircarth", "Perendor",
		"Syssale", "Aelbarik", "Drassa", "Pirn", "Maire", "Lebrus", "Birdan",
		"Fistrock", "Shotluth", "Aldain", "Nantasarn", "Carim", "Ollayos",
		"Hamish", "Sudabuk", "Belgarth", "Woodhead",
		""}
	pref          = []string{"magic", "golden", "crystal", "enchanted", "elven"}
	subloc_player = 0
)

func MM(item int) int {
	if kind(item) == T_char && noble_item(item) != FALSE {
		item = noble_item(item)
	}
	return (item_attack(item) + item_defense(item))
}

func monstermark(unit int) int {
	val := 0

	for _, e := range loop_inventory(unit) {
		if item_attack(e.item) != FALSE && item_defense(e.item) != FALSE {
			val += e.qty * MM(e.item)
		}
	}

	val += char_attack(unit) + char_defense(unit)

	return val
}

/*
 *  Tue Oct 13 13:22:30 1998 -- Scott Turner
 *
 *  Generate a single treasure for a monster.
 *
 */
func generate_one_treasure(monster int) {
	/*
	 *  Tue Oct 13 17:58:12 1998 -- Scott Turner
	 *
	 *  45% -- gold (100 to 500)
	 *  25% -- artifact
	 *  10% -- rare trade items
	 *  10% -- horses
	 *   5% -- prisoner
	 *   5% -- book
	 */
	choice := rnd(1, 100)

	if choice < 45 {
		gen_item(monster, item_gold, rnd(100, 500))
	} else if choice < 70 {
		_ = create_random_artifact(monster)
	} else if choice < 80 {
		item := random_trade_good()
		gen_item(monster, item,
			rnd(rp_item(item).trade_good/2, rp_item(item).trade_good*2))
	} else if choice < 90 {
		var item int

		switch rnd(1, 4) {
		case 1:
			item = item_pegasus
			break
		case 2:
			item = item_wild_horse
			break
		case 3:
			item = item_riding_horse
			break
		case 4:
			item = item_warmount
			break
		}
		gen_item(monster, item, rnd(3, 10))

	} else if choice < 95 {
		var pris int
		var name string

		switch rnd(1, 8) {
		case 1:
			name = "Old man"
			break
		case 2:
			name = "Old man"
			break
		case 3:
			name = "Knight"
			break
		case 4:
			name = "Princess"
			break
		case 5:
			name = "King's daughter"
			break
		case 6:
			name = "Nobleman"
			break
		case 7:
			name = "Merchant"
			break
		case 8:
			name = "Distressed Lady"
			break

		default:
			panic("!reached")
		}

		pris = new_char(0, 0, monster, 100, indep_player,
			LOY_unsworn, 0, name)

		p_magic(pris).swear_on_release = TRUE
		p_char(pris).prisoner = TRUE
	} else {
		make_teach_book(monster, 0, 0, sub_book)
	}
}

/*
 *  Tue Oct 13 13:18:39 1998 -- Scott Turner
 *
 *  Generate treasures for a monster stack, based on the monstermark.
 *  There should be a guarantee of 1 treasure for each "50 orc equivalent".
 *  Lesser strength has a % chance of getting a treasure.
 *
 */
func generate_treasure(unit, divisor int) {
	unit_mm := float64(monstermark(unit)) / float64(divisor)
	one_treasure := float64(25 * MM(item_orc))
	count, i := 0, 0

	if unit_mm < one_treasure*0.80 {
		unit_mm = one_treasure * 0.80
	}

	/*
	 *  +/- 20% will give us a "range" of treasures.
	 *
	 */
	unit_mm = (unit_mm * float64((80 + rnd(1, 40)))) / 100
	/*
	 *  And we have a small chance of being really out there.
	 *
	 */
	if rnd(1, 100) < 10 {
		unit_mm *= float64(rnd(5, 20))
	}

	for unit_mm > 0 {
		if float64(rnd(1, int(one_treasure))) < unit_mm {
			generate_one_treasure(unit)
			count++
		}
		unit_mm -= one_treasure
		one_treasure += one_treasure
	}

	/*
	 *  Tue Dec 15 12:14:02 1998 -- Scott Turner
	 *
	 *  Some monsters may have specific treasures.
	 *
	 *  Dragons always have 1-6 power jewels.
	 *
	 */
	if noble_item(unit) == item_dragon {
		for i = 1; i < rnd(1, 6); i++ {
			create_specific_artifact(unit, ART_POWER)
			count++
		}
	} else if noble_item(unit) == item_balrog {
		create_specific_artifact(unit, ART_COMBAT)
		count++
	}
}

/*
 *  Create an Old book which offers instruction in a rare skill
 *
 *  Mon Nov  4 09:55:21 1996 -- Scott Turner
 *
 *  Genericized so it can produce any skill.  If rare flag is set,
 *  only produce the "rare" books - magic and religion.  Otherwise,
 *  produce the other books with twice the chance of the rare
 *  books.  If category is set, produce only category books.
 *
 *  Wed Apr 15 12:42:44 1998 -- Scott Turner
 *
 *  Give most books a descriptive title.
 *
 *  Thu Nov 19 18:37:38 1998 -- Scott Turner
 *
 *  Prevent any city controlled by a nation (or a capital) from
 *  producing a forbidden book.
 *
 *  Mon Dec  7 06:41:05 1998 -- Scott Turner
 *
 *  Can we improve this?  Run through all the cities in this region.
 *  Create the intersection of all proscribed books for nations that
 *  have capitals or controlled cities in the region.  Then use that
 *  proscribed list to avoid books.
 *
 *  Mon Dec  7 13:05:40 1998 -- Scott Turner
 *
 *  Better yet.  If a city is controlled, use that proscribed list.
 *  If it is a capital, use that proscribed list.  Otherwise do the
 *  intersection business.
 */

func make_teach_book(who, rare, category, subkind int) int {
	count := 0
	var s string
	var p *item_magic
	skill := 0
	first, nat := 1, 0
	var proscribed_skills []int

	assert(subkind == sub_book || subkind == sub_scroll)

	/*
	 *  If it is a controlled city, then only allow that nation's
	 *  books.
	 *
	 */
	if nation(player_controls_loc(who)) != 0 {
		proscribed_skills =
			ilist_copy(rp_nation(nation(player_controls_loc(who))).proscribed_skills)
	} else {
		nat = 0
		for _, j := range loop_nation() {
			if rp_nation(j).capital == who {
				nat = j
				break
			}
		}
		//#if 0
		//        for(j=1;j<=num_nations;j++)
		//          if (nations[j].capital == who) {
		//        nat = j;
		//        break;
		//          };
		//#endif
		/*
		 *  If it is an uncontrolled capital...
		 *
		 */
		if nat != 0 {
			proscribed_skills =
				ilist_copy(rp_nation(nat).proscribed_skills)
		} else {
			/*
			 *  Otherwise, create an ilist of the intersection of all
			 *  proscribed skills.
			 *
			 */
			for _, i := range loop_city() {
				/*
				 *  Only care about cities in this region, that are controlled or
				 *  capitals.
				 */
				if region(i) == region(who) {
					/*
					 *  Is it a capital?
					 *
					 */
					nat = 0
					for _, j := range loop_nation() {
						if rp_nation(j).capital == i {
							nat = j
							break
						}
					}
					/*
					 *  Is it controlled?
					 *
					 */
					if nat == 0 && player_controls_loc(who) != 0 {
						nat = nation(player_controls_loc(who))
					}
					/*
					 *  Nat is the nation, so if it exists intersect it into
					 *  our proscribed list.  If it's our first nation, then
					 *  just use the whole proscribed list.
					 *
					 */
					if nat != 0 {
						if first != 0 {
							proscribed_skills = ilist_copy(rp_nation(nat).proscribed_skills)
							first = 0
						} else {
							/*
							 *  Intersect by removing anything from proscribed skills
							 *  that appears in the nation's proscribed skills.
							 *
							 */
							tmp := ilist_copy(proscribed_skills)
							for i = 0; i < len(tmp); i++ {
								if ilist_lookup(rp_nation(nat).proscribed_skills,
									tmp[i]) == -1 {
									proscribed_skills = rem_value(proscribed_skills, tmp[i])
								}
							}
							ilist_reclaim(&tmp)
						}
					}
				}
			}
		}
	}
	/*
	 *  At this point proscribed_skills are books we don't want to create.
	 *
	 */
	for _, newSkill := range loop_skill() {
		if newSkill == sk_adv_sorcery {
			continue
		}
		if newSkill == sk_basic_religion {
			continue
		}
		/*
		 *  Sun Feb  6 10:31:08 2000 -- Scott Turner
		 *
		 *  No point in making books for skills that won't benefit from
		 *  instruction.
		 *
		 */
		if newSkill != skill_school(newSkill) &&
			ilist_lookup(rp_skill(skill_school(newSkill)).offered, newSkill) == -1 {
			continue
		}
		if category == 0 || skill_school(newSkill) == newSkill {
			/*
			 *  Choose only rare skills if rare is set.
			 *
			 */
			if rare != FALSE &&
				!magic_skill(newSkill) &&
				!religion_skill(newSkill) {
				continue
			}
			/*
			 *   Don't pick something if it violates proscribed_skills.
			 *
			 */
			if nat != 0 && ilist_lookup(proscribed_skills, newSkill) != -1 {
				continue
			}
			/*
			 *  Would this line be picked if last?
			 *
			 */
			count++
			if rnd(1, count) == 1 {
				skill = newSkill
			}
			/*
			 *  If this is a non-rare line, pretend it is in the list twice,
			 *  thereby doubling (roughly) it's chance of being selected.
			 *
			 */
			if !magic_skill(newSkill) && !religion_skill(newSkill) {
				count++
				if rnd(1, count) == 1 {
					skill = newSkill
				}
			}
		}
	}

	newSkill := create_unique_item(who, subkind)

	/*
	 *  Thu Jun 22 06:46:09 2000 -- Scott Turner
	 *
	 *  Set an arbitrary base price of 200 gold for books.
	 *
	 */
	p_item(newSkill).base_price = 200

	p = p_item_magic(newSkill)

	if subkind == sub_book {
		p_item(newSkill).weight = 5
		p.orb_use_count = rnd(7, 28)
	} else {
		p_item(newSkill).weight = 1
		p.orb_use_count = 1
	}

	/*
	 *  Add the skill...
	 *
	 */
	p.may_study = append(p.may_study, skill)
	if p.orb_use_count < learn_time(skill) {
		p.orb_use_count = learn_time(skill)
	}

	if subkind == sub_book {
		/*
		 *  Might be a very special book.
		 *
		 */
		chance := rnd(1, 300)
		if chance == 1 {
			/*  Tome of magic... */
			set_name(newSkill, "Tome of Magic")
			p.may_study = nil
			for _, skill := range loop_skill() {
				if skill == sk_adv_sorcery {
					continue
				}
				if newSkill == sk_basic_religion {
					continue
				}
				if magic_skill(skill) && skill_school(skill) == skill {
					p.may_study = append(p.may_study, skill)
				}
			}
		} else if chance == 2 {
			/*  Bible */
			set_name(newSkill, "Great Bible")
			p.may_study = nil
			for _, skill := range loop_skill() {
				if religion_skill(skill) && skill_school(skill) == skill {
					p.may_study = append(p.may_study, skill)
				}
			}
		} else if chance == 3 {
			/*  Tome of All Knowledge... */
			set_name(newSkill, "Tome of All Knowledge")
			p.may_study = nil
			for _, skill := range loop_skill() {
				if skill == sk_adv_sorcery {
					continue
				}
				if newSkill == sk_basic_religion {
					continue
				}
				if ilist_lookup(proscribed_skills, skill) != -1 {
					continue
				}
				if skill_school(skill) == skill {
					p.may_study = append(p.may_study, skill)
				}
			}
		} else if rnd(1, 100) < 75 {
			s = fmt.Sprintf("Manual of %s", just_name(skill))
			set_name(newSkill, s)
		} else {
			switch rnd(1, 3) {
			case 1:
				s = "old book"
				break
			case 2:
				s = "rare book"
				break
			case 3:
				s = "strange tome"
				break
			case 4:
				s = "ancient manual"
				break

			default:
				panic("!reached")
			}
			set_name(newSkill, s)
		}
	}

	proscribed_skills = nil
	return newSkill
}

/*
 *  Find an artifact in our region held by a subloc monster
 *  which is not only-defeatable by another artifact.
 */

func free_artifact(where int) int {
	reg := region(where)
	var owner int
	var l []int
	var ret int

	for _, i := range loop_item() {
		if subkind(i) != sub_artifact {
			continue
		}

		owner = item_unique(i)
		assert(owner != 0)

		if region(owner) != reg {
			continue
		}

		if !is_npc(owner) ||
			npc_program(owner) != PROG_subloc_monster {
			continue
		}

		if only_defeatable(owner) != FALSE {
			continue
		}

		l = append(l, i)
	}

	if len(l) == 0 {
		return 0
	}

	ret = l[rnd(0, len(l)-1)]

	l = nil

	return ret
}

//#if 0
//static int
//new_artifact(int who)
//{
//    int new;
//    char *s;
//
//    new = create_unique_item(who, sub_artifact);
//
//    switch (rnd(1,4))
//    {
//    case 1:
//        s = art_att_s[rnd(0,2)];
//        p_item_magic(new).attack_bonus = rnd(1,10) * 5;
//        break;
//
//    case 2:
//        s = art_def_s[rnd(0,2)];
//        p_item_magic(new).defense_bonus = rnd(1,10) * 5;
//        break;
//
//    case 3:
//        s = art_mis_s[rnd(0,3)];
//        p_item_magic(new).missile_bonus = rnd(1,10) * 5;
//        break;
//
//    case 4:
//        s = art_mag_s[rnd(0,2)];
//        p_item_magic(new).aura_bonus = rnd(1,3);
//        break;
//
//    default:
//        panic("!reached");
//    }
//
//    if (rnd(1,3) < 3)
//    {
//        s = sout("%s %s", pref[rnd(0,4)], s);
//    }
//    else
//    {
//        var i int
//
//        for i =  0; of_names[i]; i++
//            ;
//        i = rnd(0, i-1);
//
//        s = sout("%s of %s", cap_(s), of_names[i]);
//    }
//
//    p_item(new).weight = 10;
//    set_name(new, s);
//
//    return new;
//}
//#endif

func new_monster(where int) int {
	item := 0

	switch subkind(where) {
	case sub_graveyard, sub_battlefield:
		switch rnd(1, 3) {
		case 1:
			item = item_corpse
			break
		case 2:
			item = item_skeleton
			break
		case 3:
			item = item_spirit
			break
		default:
			panic("!reached")
		}
		break

	case sub_ench_forest:
		switch rnd(1, 2) {
		case 1:
			item = item_elf
			break
		case 2:
			item = item_faery
			break
		default:
			panic("!reached")
		}
		break

	case sub_island:
		if rnd(1, 2) == 1 {
			switch rnd(1, 3) {
			case 1:
				item = item_pirate
				break
			case 2:
				item = item_spider
				break
			case 3:
				item = item_cyclops
				break
			default:
				panic("!reached")
			}
			break
		}

	default:
		item = random_beast(0)
	}

	/*
	 *  Base this on the split size.
	 *
	 */
	// todo: can rnd(2, n) ever be 0?
	newStack := create_monster_stack(item, or_int(rnd(2, item_split(item)) != 0, 10, item_split(item)/2), where)
	p_char(newStack).npc_prog = PROG_subloc_monster
	return newStack
}

func seed_subloc_with_monster(where, limit int) int {
	var monster int

	monster = new_monster(where)
	/*
	 *  Keep track of how many of each kind of monster we
	 *  generate.
	 *
	 */
	bx[noble_item(monster)].temp++
	generate_treasure(monster, 1)

	if rnd(1, 6) == 1 {
		var item int
		/*
		 *  Temporarily set only_vulnerable for ourselves so we don't
		 *  have a circular problem.  free_artifact() will take care of
		 *  skipping over other only_vulnerable's.
		 */
		p_misc(monster).only_vulnerable = 1
		item = free_artifact(monster)

		if item != 0 {
			rp_misc(monster).only_vulnerable = item
		} else {
			rp_misc(monster).only_vulnerable = 0
		}
	}
	/*
	 *  Fri Nov 26 17:54:59 1999 -- Scott Turner
	 *
	 *  Make him a subloc monster.
	 *
	 */
	p_char(monster).npc_prog = PROG_subloc_monster

	/*
	 *  Tue May 25 11:02:29 1999 -- Scott Turner
	 *
	 *  Close the border.
	 *
	 */
	p_subloc(where).control.closed = true

	if limit != 0 && rnd(1, 6) < 3 {
		p_subloc(where).entrance_size = rnd(1, 4) + rnd(1, 4)
	}

	return monster
}

func seed_monster_sublocs(all bool) {
	stage("seed_monster_sublocs()")

	clear_temps(T_item)

	for _, where := range loop_loc() {
		if loc_depth(where) != LOC_subloc {
			continue
		}

		if in_faery(where) || in_hades(where) {
			continue
		}

		if subkind(where) == sub_city {
			continue
		}

		if controlled_humans_here(province(where)) != FALSE {
			continue
		}

		seed_subloc_with_monster(where, 1)
	}

	for _, i := range loop_item() {
		if bx[i].temp != 0 {
			wout(gm_player, "Generated %d stacks of %s.",
				bx[i].temp, box_name(i))
		}
	}
}

/*
 *  Tue Oct 13 18:50:29 1998 -- Scott Turner
 *
 *  Randomly select an appropriate subloc for reseeding.
 *
 */
func add_lair_monster() {
	sum := 0
	var where, owner int

	for _, where = range loop_loc() {
		if loc_depth(where) != LOC_subloc ||
			in_faery(where) || in_hades(where) ||
			subkind(where) == sub_city ||
			controlled_humans_here(province(where)) != FALSE ||
			has_item(province(where), item_peasant) > 100 {
			continue
		}

		/*
		 *  A laired beast is one that owns the location and
		 *  has the subloc monster program.
		 *
		 */
		owner = first_character(where)
		if owner != 0 && is_real_npc(owner) {
			continue
		}
		/*
		 *  Possibly select.
		 *
		 */
		sum++
		// this check set choice, but choice isn't used anywhere
		//if rnd(1, sum) == 1 {
		//	choice := where
		//}
	}

	wout(gm_player, "Seeding %s with a monster.", box_name(where))
	seed_subloc_with_monster(where, 0)
}

/*
 *  Thu Sep 24 13:40:03 1998 -- Scott Turner
 *
 *  Like seed, except only sometimes and never where there
 *  is civilization.
 *
 *  Tue Oct 13 18:38:17 1998 -- Scott Turner
 *
 *  What this really needs to do is to try to keep the sublocs partially
 *  stocked.  What we need to do is count the number of sublocs, the number
 *  of sublocs with laired beasts and go from there...
 *
 */
func reseed_monster_sublocs() {
	var i int
	var where int
	num_sublocs := 0
	num_laired := 0
	var needed, owner int

	stage("reseed_monster_sublocs()")

	for _, where = range loop_loc() {

		if loc_depth(where) != LOC_subloc {
			continue
		}

		if in_faery(where) || in_hades(where) {
			continue
		}

		if subkind(where) == sub_city {
			continue
		}

		if controlled_humans_here(province(where)) != FALSE {
			continue
		}

		if has_item(province(where), item_peasant) > 100 {
			continue
		}

		num_sublocs++

		/*
		 *  A laired beast is one that owns the location and
		 *  has the subloc monster program.
		 *
		 */
		owner = first_character(where)
		if owner != 0 && is_npc(owner) && npc_program(owner) == PROG_subloc_monster {
			num_laired++
		}
	}

	/*  We want 33% of the sublocs to have laired monsters. */
	needed = ((num_sublocs * 33) / 100) - num_laired

	if needed > 0 {
		log_output(LOG_MISC, "Reseeding %d sublocs.", needed)
	} else {
		log_output(LOG_MISC, "No need to reseed sublocs.")
	}

	for i = 0; i < needed; i++ {
		add_lair_monster()
	}
}
