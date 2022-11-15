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

const (
	COMMON         = 1000
	UNUSUAL        = 10
	RARE           = 1
	RANDOM_BEAST   = -1
	RANDOM_SOLDIER = -2
	RANDOM_SKILL   = -3
	RANDOM_USE     = -4
)

/*
 *  Names of artifact types -- must align w/ oly.h
 *
 */
var artifact_names = []string{
	"No enchantment.",                             /* ART_NONE */
	"Combat Artifact, +%d, affects: %s.",          /* ART_COMBAT */
	"Leadership of Men, +%d%%.",                   /* ART_CTL_MEN */
	"Control Beasts, +%d%%.",                      /* ART_CTL_BEASTS */
	"Safety from Attack from %s.",                 /* ART_SAFETY */
	"Improved Attack, %s, +%d%%.",                 /* ART_IMPRV_ATT */
	"Improved Defense, %s, +%d%%.",                /* ART_IMPRV_DEF */
	"Safety at Sea.",                              /* ART_SAFE_SEA */
	"Defensive Terrain Enchantment in %s, +%d%%.", /* ART_TERRAIN */
	"Fast Terrain in %s, +%d day(s).",             /* ART_FAST_TERR */
	"Speed Use of %s by %d day(s).",               /* ART_SPEED_USE */
	"Hellring.",                                   /* ART_PROT_HADES */
	"Elfstone.",                                   /* ART_PROT_FAERY */
	"Hard Workers, +%d%% effort.",                 /* ART_WORKERS */
	"Increased Income in %s, +%d%%.",              /* ART_INCOME */
	"Fast Learning, +%d day(s).",                  /* ART_LEARNING */
	"Fast Teaching, +%d day(s).",                  /* ART_TEACHING */
	"Fast Training of %s.",                        /* ART_TRAINING */
	"Destroy Monster: %s, %d charges.",            /* ART_DESTROY */
	"Grant Skill: %s.",                            /* ART_SKILL */
	"Flying, %d weight.",                          /* ART_FLYING */
	"Protection from Skill: %s.",                  /* ART_PROT_SKILL */
	"Shield Location.",                            /* ART_SHIELD_PROV */
	"Riding, %d weight.",                          /* ART_RIDING */
	"Power Jewel, +%d aura/piety, %d charge(s).",  /* ART_POWER */
	"Summon Aid, %s, %d charges.",                 /* ART_SUMMON_AID */
	"Reduced Maintenance, %d%%.",                  /* ART_MAINTENANCE */
	"Improve Bargaining, %d%%.",                   /* ART_BARGAIN */
	"Weightlessness, %d weight.",                  /* ART_WEIGHTLESS */
	"Increased Healing, +%d health points.",       /* ART_HEALING */
	"Protection from Sickness.",                   /* ART_SICKNESS */
	"Restore Life, %d charges.",                   /* ART_RESTORE */
	"Teleport, %d weight, %d charges.",            /* ART_TELEPORT */
	"Orb of Scrying, %d charges.",                 /* ART_ORB */
	"Crown of Control over %s, %d charges.",       /* ART_CROWN */
	"Auraculum belonging to %s, +%s aura.",        /* ART_AURACULUM */
	"Carry Great Loads, +%d weight.",              /* ART_CARRY */
	"The Pen Crown: (+%d, +%d) in combat, %d uses as an Orb of Scrying.", /* ART_PEN */
	""}

type artifact_ent struct {
	what       int /* The ART_ number */
	rarity     int /* How rare is this? */
	min_param1 int
	max_param1 int /* Range for parameter #1 */
	min_param2 int
	max_param2 int /* Range for parameter #1 */
	min_uses   int
	max_uses   int /* Range for charges */
}

var artifact_tbl = []artifact_ent{
	{
		ART_COMBAT, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		5, 25,      /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_COMBAT, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		5, 25,      /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_COMBAT, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		5, 25,      /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_COMBAT, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		5, 25,      /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_COMBAT, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		5, 25,      /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_COMBAT, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		25, 50,     /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_COMBAT, /* What this is... */
		RARE,       /* COMMON, UNUSUAL, or RARE */
		50, 50,     /* Range for param1 */
		4095, 4095, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_CTL_MEN, /* What this is... */
		COMMON,      /* COMMON, UNUSUAL, or RARE */
		5, 25,       /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_CTL_BEASTS, /* What this is... */
		COMMON,         /* COMMON, UNUSUAL, or RARE */
		5, 25,          /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_SAFETY,      /* What this is... */
		COMMON,          /* COMMON, UNUSUAL, or RARE */
		RANDOM_BEAST, 0, /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_IMPRV_ATT, /* What this is... */
		COMMON,        /* COMMON, UNUSUAL, or RARE */
		5, 25,         /* Range for param1 */
		RANDOM_SOLDIER, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_IMPRV_DEF, /* What this is... */
		COMMON,        /* COMMON, UNUSUAL, or RARE */
		5, 25,         /* Range for param1 */
		RANDOM_SOLDIER, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_SAFE_SEA, /* What this is... */
		RARE,         /* COMMON, UNUSUAL, or RARE */
		1, 1,         /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_TERRAIN, /* What this is... */
		COMMON,      /* COMMON, UNUSUAL, or RARE */
		5, 25,       /* Range for param2 */
		sub_ocean, sub_poppy_field, /* Range for param1 */
		0, 0, /* Range for charges */
	},
	{
		ART_FAST_TERR, /* What this is... */
		RARE,          /* COMMON, UNUSUAL, or RARE */
		1, 2,          /* Range for param2 */
		sub_ocean, sub_swamp, /* Range for param1 */
		0, 0, /* Range for charges */
	},
	{
		ART_SPEED_USE, /* What this is... */
		UNUSUAL,       /* COMMON, UNUSUAL, or RARE */
		1, 2,          /* Range for param2 */
		RANDOM_USE, 0, /* Range for param1 */
		0, 0, /* Range for charges */
	},
	{
		ART_PROT_HADES, /* What this is... */
		UNUSUAL,        /* COMMON, UNUSUAL, or RARE */
		0, 0,           /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_WORKERS, /* What this is... */
		COMMON,      /* COMMON, UNUSUAL, or RARE */
		5, 25,       /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_INCOME, /* What this is... */
		UNUSUAL,    /* COMMON, UNUSUAL, or RARE */
		5, 25,      /* Range for param1 */
		sub_inn, sub_inn, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_INCOME, /* What this is... */
		UNUSUAL,    /* COMMON, UNUSUAL, or RARE */
		5, 25,      /* Range for param1 */
		sub_temple, sub_temple, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_LEARNING, /* What this is... */
		UNUSUAL,      /* COMMON, UNUSUAL, or RARE */
		1, 2,         /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_TEACHING, /* What this is... */
		RARE,         /* COMMON, UNUSUAL, or RARE */
		1, 2,         /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_TRAINING,      /* What this is... */
		UNUSUAL,           /* COMMON, UNUSUAL, or RARE */
		RANDOM_SOLDIER, 0, /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_DESTROY,     /* What this is... */
		COMMON,          /* COMMON, UNUSUAL, or RARE */
		RANDOM_BEAST, 0, /* Range for param1 */
		0, 0, /* Range for param2 */
		1, 3, /* Range for charges */
	},
	//#if 0
	//                {
	//                  ART_SKILL,	/* What this is... */
	//                  COMMON,		/* COMMON, UNUSUAL, or RARE */
	//                  RANDOM_SKILL, 0,	/* Range for param1 */
	//                  0,  0,		/* Range for param2 */
	//                  0,  0,		/* Range for charges */
	//                },
	//#endif
	{
		ART_FLYING, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		100, 500,   /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_CARRY, /* What this is... */
		COMMON,    /* COMMON, UNUSUAL, or RARE */
		100, 500,  /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_RIDING, /* What this is... */
		COMMON,     /* COMMON, UNUSUAL, or RARE */
		100, 500,   /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_PROT_SKILL,               /* What this is... */
		COMMON,                       /* COMMON, UNUSUAL, or RARE */
		sk_aura_blast, sk_aura_blast, /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_PROT_SKILL, /* What this is... */
		COMMON,         /* COMMON, UNUSUAL, or RARE */
		sk_reveal_vision, sk_reveal_vision,
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_SHIELD_PROV, /* What this is... */
		RARE,            /* COMMON, UNUSUAL, or RARE */
		0, 0,            /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_POWER, /* What this is... */
		UNUSUAL,   /* COMMON, UNUSUAL, or RARE */
		5, 25,     /* Range for param1 */
		0, 0, /* Range for param2 */
		1, 1, /* Range for charges */
	},
	{
		ART_SUMMON_AID,    /* What this is... */
		UNUSUAL,           /* COMMON, UNUSUAL, or RARE */
		RANDOM_SOLDIER, 0, /* Range for param1 */
		5, 25, /* Range for param2 */
		1, 2, /* Range for charges */
	},
	{
		ART_MAINTENANCE, /* What this is... */
		COMMON,          /* COMMON, UNUSUAL, or RARE */
		5, 25,           /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_BARGAIN, /* What this is... */
		COMMON,      /* COMMON, UNUSUAL, or RARE */
		5, 25,       /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_WEIGHTLESS, /* What this is... */
		COMMON,         /* COMMON, UNUSUAL, or RARE */
		100, 500,       /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_HEALING, /* What this is... */
		COMMON,      /* COMMON, UNUSUAL, or RARE */
		5, 15,       /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_SICKNESS, /* What this is... */
		COMMON,       /* COMMON, UNUSUAL, or RARE */
		0, 0,         /* Range for param1 */
		0, 0, /* Range for param2 */
		0, 0, /* Range for charges */
	},
	{
		ART_RESTORE, /* What this is... */
		RARE,        /* COMMON, UNUSUAL, or RARE */
		0, 0,        /* Range for param1 */
		0, 0, /* Range for param2 */
		1, 1, /* Range for charges */
	},
	{
		ART_TELEPORT, /* What this is... */
		UNUSUAL,      /* COMMON, UNUSUAL, or RARE */
		100, 500,     /* Range for param1 */
		0, 0, /* Range for param2 */
		1, 3, /* Range for charges */
	},
	{
		ART_ORB, /* What this is... */
		COMMON,  /* COMMON, UNUSUAL, or RARE */
		0, 0,    /* Range for param1 */
		0, 0, /* Range for param2 */
		3, 9, /* Range for charges */
	},
	{
		ART_CROWN,       /* What this is... */
		COMMON,          /* COMMON, UNUSUAL, or RARE */
		RANDOM_BEAST, 0, /* Range for param1 */
		0, 0, /* Range for param2 */
		1, 3, /* Range for charges */
	},
	{
		ART_PROT_FAERY, /* What this is... */
		COMMON,         /* COMMON, UNUSUAL, or RARE */
		0, 0,           /* Range for param1 */
		0, 0, /* Range for param2 */
		3, 9, /* Range for charges */
	},
	{
		ART_NONE, /* What this is... */
		RARE,     /* COMMON, UNUSUAL, or RARE */
		0, 0,     /* Range for param1 */
		0, 0, /* Range for param2 */
		2, 8, /* Range for charges */
	}}

/*
 *  Fri Oct  2 18:31:10 1998 -- Scott Turner
 *
 *  Select (in one pass) something from the artifact table.
 *
 */
func get_random_artifact() int {
	i, choice := 0, 0
	sum := 0

	for i = 0; artifact_tbl[i].what != ART_NONE; i++ {
		sum += artifact_tbl[i].rarity
		if rnd(1, sum) <= artifact_tbl[i].rarity {
			choice = i
		}
	}

	return choice
}

/*
 *  Sat Oct  3 18:22:54 1998 -- Scott Turner
 *
 *  Select a random soldier-type unit.
 *
 */
func random_soldier() int {
	i, choice := 0, 0
	sum := 0

	for _, i = range loop_item() {
		if item_attack(i) != FALSE &&
			item_defense(i) != FALSE &&
			man_item(i) != FALSE &&
			rp_item(i).maintenance != FALSE {
			sum++
			if rnd(1, sum) == 1 {
				choice = i
			}
		}
	}

	return choice

}

/*
 *  Sat Oct  3 18:29:00 1998 -- Scott Turner
 *
 *  Select a random skill -- not category skills.
 *
 */
func random_skill() int {
	i, choice := 0, 0
	sum := 0

	for _, i = range loop_skill() {
		if i != skill_school(i) {
			sum++
			if rnd(1, sum) == 1 {
				choice = i
			}
		}
	}

	return choice
}

/*
 *  Sat Oct  3 18:31:24 1998 -- Scott Turner
 *
 *  Select a random use skill.
 *
 */
func random_use() int {
	i, choice := 0, 0
	sum := 0

	for _, i = range loop_skill() {
		if i != skill_school(i) && find_use_entry(i) != FALSE {
			sum++
			if rnd(1, sum) == 1 {
				choice = i
			}
		}
	}

	return choice
}

/*
 *  Sat Oct  3 18:33:15 1998 -- Scott Turner
 *
 *  Select a random beast.  This should be inversely weighted by the
 *  combat prowess of the beast.
 *
 */
func random_beast(sk int) int {
	i, choice := 0, 0
	sum := 0
	var val int

	for _, i = range loop_item() {
		if item_attack(i) != FALSE &&
			item_defense(i) != FALSE &&
			item_wild(i) != FALSE &&
			(sk == FALSE || subkind(i) == schar(sk)) &&
			rp_item(i).maintenance == FALSE {
			val = MAX_MM - MM(i) + 1
			if val < 0 {
				val = 1
			}
			sum += val
			if rnd(1, sum) <= val {
				choice = i
			}
		}
	}

	assert(choice != FALSE)
	return choice
}

/*
 *  Tue Oct  6 12:33:26 1998 -- Scott Turner
 *
 *  Combat artifacts take some special consideration.
 *
 */
func create_combat_artifact(piece int) {
	/*
	 *  Set some number of the flags.  Usually just one,
	 *  but a chance for more.
	 *
	 */
	rp_item_artifact(piece).param2 = 0
	for {
		rp_item_artifact(piece).param2 |= 1 << rnd(0, 11)
		if rnd(1, 100) < 30 {
			continue
		}
		break
	}
}

/*
 *  Fri Oct  2 18:34:25 1998 -- Scott Turner
 *
 *  Create and return a random artifact.
 *
 *  (1) About 20% ART_NONE
 *  (2) Otherwise, randomly instantiated.
 *
 */
func create_random_artifact(monster int) int {
	var select_ int
	piece := create_unique_item(monster, sub_magic_artifact)

	set_name(piece, "Unknown artifact")
	p_item(piece).weight = 5
	p_item_artifact(piece).type_ = ART_NONE
	rp_item_artifact(piece).param1 = 0
	rp_item_artifact(piece).param2 = 0
	rp_item_artifact(piece).uses = 0

	/*
	 *  Possibly nothing.
	 *
	 */
	if rnd(1, 100) < 20 {
		return piece
	}
	/*
	 *  No, so select something.
	 *
	 */
	select_ = get_random_artifact()
	p_item_artifact(piece).type_ = artifact_tbl[select_].what
	/*
	 *  Set parameter one, which might be special.
	 *
	 */
	if artifact_tbl[select_].min_param1 == artifact_tbl[select_].max_param1 {
		rp_item_artifact(piece).param1 = artifact_tbl[select_].min_param1
	} else if artifact_tbl[select_].min_param1 == RANDOM_SOLDIER {
		rp_item_artifact(piece).param1 = random_soldier()
	} else if artifact_tbl[select_].min_param1 == RANDOM_SKILL {
		rp_item_artifact(piece).param1 = random_skill()
	} else if artifact_tbl[select_].min_param1 == RANDOM_BEAST {
		rp_item_artifact(piece).param1 = random_beast(0)
	} else if artifact_tbl[select_].min_param1 == RANDOM_USE {
		rp_item_artifact(piece).param1 = random_use()
	} else {
		rp_item_artifact(piece).param1 =
			rnd(artifact_tbl[select_].min_param1, artifact_tbl[select_].max_param1)
	}
	/*
	 *  Set parameter two, no specials
	 *
	 */
	if artifact_tbl[select_].min_param2 == artifact_tbl[select_].max_param2 {
		rp_item_artifact(piece).param2 = artifact_tbl[select_].min_param2
	} else if artifact_tbl[select_].min_param2 == RANDOM_SOLDIER {
		rp_item_artifact(piece).param2 = random_soldier()
	} else if artifact_tbl[select_].min_param2 == RANDOM_SKILL {
		rp_item_artifact(piece).param2 = random_skill()
	} else if artifact_tbl[select_].min_param2 == RANDOM_BEAST {
		rp_item_artifact(piece).param2 = random_beast(0)
	} else if artifact_tbl[select_].min_param2 == RANDOM_USE {
		rp_item_artifact(piece).param2 = random_use()
	} else {
		rp_item_artifact(piece).param2 =
			rnd(artifact_tbl[select_].min_param2, artifact_tbl[select_].max_param2)
	}
	/*
	 *  Set uses, no specials
	 *
	 */
	if artifact_tbl[select_].min_uses == artifact_tbl[select_].max_uses {
		rp_item_artifact(piece).uses = artifact_tbl[select_].min_uses
	} else {
		rp_item_artifact(piece).uses =
			rnd(artifact_tbl[select_].min_uses, artifact_tbl[select_].max_uses)
	}
	/*
	 *  And special case for combat.
	 *
	 */
	if rp_item_artifact(piece).type_ == ART_COMBAT &&
		rp_item_artifact(piece).param2 == FALSE {
		create_combat_artifact(piece)
	}
	/*
	 *  And return the artifact.
	 *
	 */
	return piece
}

func v_make_artifact(c *command) int {

	create_random_artifact(c.who)
	return TRUE
}

/*
 *  Sun Oct  4 17:29:15 1998 -- Scott Turner
 *
 *  Find the best artifact of a given type on a noble.
 *
 */
func best_artifact(who int, type_ int, param2 int, uses int) int {
	var e *item_ent
	best := 0
	best_val := 0

	for _, e = range loop_inventory(who) {
		if is_artifact(e.item) != nil {
			if rp_item_artifact(e.item).type_ == type_ &&
				rp_item_artifact(e.item).param1 > best_val &&
				(param2 == FALSE || rp_item_artifact(e.item).param2 == param2) &&
				(uses == FALSE || rp_item_artifact(e.item).uses > 0) {
				best_val = rp_item_artifact(e.item).param1
				best = e.item
			}
		}
	}

	return best
}

/*
 *  Sun Oct  4 17:36:03 1998 -- Scott Turner
 *
 *  Find any artifact matching.
 *
 */
func has_artifact(who, type_, p1, p2, charges int) int {
	var e *item_ent

	for _, e = range loop_inventory(who) {
		a := is_artifact(e.item)
		if a != nil {
			if a.type_ == type_ &&
				(p1 == FALSE || a.param1 == p1) &&
				(p2 == FALSE || a.param2 == p2) &&
				(charges == FALSE || a.uses > 0) {
				return e.item
			}
		}
	}

	return 0
}

/*
 *  Thu Oct  8 17:41:24 1998 -- Scott Turner
 *
 *  Find a combat bonus.
 *
 */
func combat_artifact_bonus(who int, part int, unused *int) int {
	var e *item_ent
	best := 0

	for _, e = range loop_inventory(who) {
		a := is_artifact(e.item)
		if a != nil {
			if a.type_ == ART_COMBAT && (a.param2&part) != FALSE &&
				a.param1 > best {
				best = a.param1
			}
		}
	}

	return best

}

/*
 *  Thu Oct  8 18:30:46 1998 -- Scott Turner
 *
 *  Calculate someone's effective workforce, including artifacts.
 *
 */
func effective_workers(who int) int {
	w := has_item(who, item_worker)
	a := best_artifact(who, ART_WORKERS, 0, 0)

	if a != FALSE {
		w = (w * (100 + rp_item_artifact(a).param1)) / 100
	}

	return w

}

/*
 *  Fri Oct  9 18:19:22 1998 -- Scott Turner
 *
 *  Destroying monster.
 *
 */
func v_art_destroy(c *command) int {
	item := c.use_skill
	where := province(subloc(c.who))
	var num int
	var t *item_ent
	var kind int

	assert(rp_item_artifact(item) != nil)
	kind = rp_item_artifact(item).param1

	if rp_item_artifact(item).uses < 1 {
		wout(c.who, "Nothing happens.")
		wout(c.who, "%s vanishes!", box_name(item))
		destroy_unique_item(c.who, item)
		return TRUE
	}

	log_output(LOG_SPECIAL, "Destroy monster artifact %s used by %s",
		box_code_less(item), box_code_less(player(c.who)))

	wout(c.who, "A golden glow suffuses the province.")
	wout(where, "A golden glow suffuses the province.")

	for _, num = range loop_all_here(where) {
		wout(num, "A golden glow suffuses the province.")

		for _, t = range loop_inventory(num) {
			if t.item == kind {
				wout(num, "%s vanished!", box_name_qty(t.item, t.qty))
				consume_item(num, t.item, t.qty)
			}
		}

		if subkind(num) == sub_ni && noble_item(num) == kind {
			kill_char(num, MATES, S_body)
		}
	}

	rp_item_artifact(item).uses--
	if rp_item_artifact(item).uses == FALSE {
		wout(c.who, "%s vanishes.", box_name(item))
		destroy_unique_item(c.who, item)
	}

	return TRUE
}

/*
 *  Fri Oct  9 19:01:32 1998 -- Scott Turner
 *
 *  Power Jewel
 *
 */
func v_power_jewel(c *command) int {
	item := c.use_skill

	assert(rp_item_artifact(item) != nil)
	//kind := rp_item_artifact(item).param1;

	if rp_item_artifact(item).uses < 1 {
		wout(c.who, "Nothing happens.")
		wout(c.who, "%s vanishes!", box_name(item))
		destroy_unique_item(c.who, item)
		return TRUE
	}

	log_output(LOG_SPECIAL, "Power jewel %s used by %s",
		box_code_less(item), box_code_less(player(c.who)))

	wout(c.who, "A golden glow suffuses your being.")

	if is_priest(c.who) != FALSE {
		wout(c.who, "You feel the hand of %s.", god_name(is_priest(c.who)))
		rp_char(c.who).religion.piety += rp_item_artifact(item).param1
	} else if is_magician(c.who) {
		wout(c.who, "You feel charged with power!")
		add_aura(c.who, rp_item_artifact(item).param1)
	} else if p_char(c.who).health < 100 {
		p_char(c.who).health += rp_item_artifact(item).param1
		if p_char(c.who).health > 100 {
			p_char(c.who).health = 100
		}
		wout(c.who, "You feel healing suffuse your body!")
	} else {
		wout(c.who, "You feel a vague sense of loss.")
	}

	rp_item_artifact(item).uses--
	if rp_item_artifact(item).uses == FALSE {
		wout(c.who, "%s vanishes.", box_name(item))
		destroy_unique_item(c.who, item)
	}

	return TRUE
}

/*
 *  Sun Oct 11 18:32:02 1998 -- Scott Turner
 *
 *  Summon Aid
 *
 */
func v_summon_aid(c *command) int {
	item := c.use_skill
	var kind, num int

	assert(rp_item_artifact(item) != nil)
	kind = rp_item_artifact(item).param1
	num = rp_item_artifact(item).param2
	assert(kind > 0 && num > 0)

	if rp_item_artifact(item).uses < 1 {
		wout(c.who, "Nothing happens.")
		wout(c.who, "%s vanishes!", box_name(item))
		destroy_unique_item(c.who, item)
		return TRUE
	}

	log_output(LOG_SPECIAL, "Summon aid %s used by %s",
		box_code_less(item), box_code_less(player(c.who)))

	wout(loc(c.who), "There is a momentary flash of yellow light.")
	wout(c.who, "A bright yellow light momentarily blinds you.")
	wout(c.who, "You are now accompanied by %s.",
		box_name_qty(kind, num))
	gen_item(c.who, kind, num)

	rp_item_artifact(item).uses--
	if rp_item_artifact(item).uses == FALSE {
		wout(c.who, "%s vanishes.", box_name(item))
		destroy_unique_item(c.who, item)
	}

	return TRUE
}

/*
 *  Tue Oct 13 12:42:52 1998 -- Scott Turner
 *
 *  Teleport using an artifact.
 *
 */
func v_art_teleport(c *command) int {
	item := c.use_skill
	dest := c.b
	var w weights

	assert(rp_item_artifact(item) != nil)

	if !is_loc_or_ship(dest) {
		wout(c.who, "There is no location %s.", c.parse[1])
		return FALSE
	}

	if rp_item_artifact(item).uses < 1 {
		wout(c.who, "Nothing happens.")
		wout(c.who, "%s vanishes!", box_name(item))
		destroy_unique_item(c.who, item)
		return TRUE
	}

	determine_stack_weights(c.who, &w, false)

	if w.total_weight > rp_item_artifact(item).param1 {
		wout(c.who, "%s hums briefly but nothing happens.")
		return FALSE
	}

	wout(loc(c.who), "There is a momentary flash of yellow light.")
	wout(c.who, "A bright yellow light momentarily blinds you.")
	do_jump(c.who, dest, 0, false)

	log_output(LOG_SPECIAL, "Teleport artifact %s used by %s",
		box_code_less(item), box_code_less(player(c.who)))

	rp_item_artifact(item).uses--
	if rp_item_artifact(item).uses == FALSE {
		wout(c.who, "%s vanishes.", box_name(item))
		destroy_unique_item(c.who, item)
	}

	return TRUE
}

/*
 *  Tue Oct 13 12:53:57 1998 -- Scott Turner
 *
 *  Use a scrying artifact.
 *
 */
func v_art_orb(c *command) int {
	item := c.use_skill
	target := c.b
	owner, where := 0, 0

	assert(rp_item_artifact(item) != nil)

	if rp_item_artifact(item).uses < 1 {
		wout(c.who, "Nothing happens.")
		wout(c.who, "%s vanishes!", box_name(item))
		destroy_unique_item(c.who, item)
		return TRUE
	}

	switch kind(target) {
	case T_loc, T_ship:
		where = province(target)
		break

	case T_char:
		where = province(target)
		break

	case T_item:
		if owner = item_unique(target); owner != FALSE {
			where = province(owner)
		}
		break
	}

	if where == 0 {
		wout(c.who, "%s hums briefly but nothing happens.", box_name(item))
		return FALSE
	} else {
		wout(loc(c.who), "There is a momentary flash of yellow light.")
		wout(c.who, "A bright yellow light momentarily blinds you and then a vision of %s appears:", box_name(where))
		show_loc(c.who, where)
		alert_scry_generic(c.who, where)
	}

	log_output(LOG_SPECIAL, "Scry artifact %s used by %s",
		box_code_less(item), box_code_less(player(c.who)))

	rp_item_artifact(item).uses--
	if rp_item_artifact(item).uses == FALSE {
		wout(c.who, "%s vanishes.", box_name(item))
		destroy_unique_item(c.who, item)
	}

	return TRUE
}

/*
 *  Tue Oct 13 13:00:58 1998 -- Scott Turner
 *
 *  A crown artifact.
 *
 */
func v_art_crown(c *command) int {
	item := c.use_skill
	target := c.b

	assert(rp_item_artifact(item) != nil)

	if rp_item_artifact(item).uses < 1 {
		wout(c.who, "Nothing happens.")
		wout(c.who, "%s vanishes!", box_name(item))
		destroy_unique_item(c.who, item)
		return TRUE
	}

	wout(c.who, "You are momentarily blinded by a flash of yellow light.")
	wout(loc(c.who), "There is a momentary flash of yellow light.")

	/*
	 *  Need to be in the same location.
	 *
	 */
	if subloc(c.who) != subloc(target) {
		wout(c.who, "Nothing happens.")
		return FALSE
	}

	if noble_item(target) == rp_item_artifact(item).param1 {
		var p *entity_player
		/*
		 *  Make player the lord of this unit, restrict its commands,
		 *  and add it to the list of the player's units.
		 *
		 */
		wout(c.who, "%s joins your faction.", box_name(target))
		set_lord(target, c.who, LOY_UNCHANGED, 0)
		p_misc(target).cmd_allow = 'r'
		p_char(target).break_point = 0
		p = p_player(c.who)
		p.units = append(p.units, target)
	}

	log_output(LOG_SPECIAL, "Crown artifact %s used by %s",
		box_code_less(item), box_code_less(player(c.who)))

	rp_item_artifact(item).uses--
	if rp_item_artifact(item).uses == FALSE {
		wout(c.who, "%s vanishes.", box_name(item))
		destroy_unique_item(c.who, item)
	}

	return TRUE
}

/*
 *  Fri Oct 16 09:11:24 1998 -- Scott Turner
 *
 *  Special routine particularly for combat artifacts.
 *
 */
func describe_combat_artifact(who, target int, header string) {
	first := 1
	var i, val int
	var buf, total string

	/*
	 *  This steps through all the possible combinations of
	 *  enchantments.
	 *
	 */
	for i = 0; i < 12; i++ {
		val = (1 << i)
		buf = ""
		if (rp_item_artifact(target).param2&CA_N_MELEE) != FALSE && val == CA_N_MELEE {
			buf = sout("personal melee attack")
		}
		if (rp_item_artifact(target).param2&CA_N_MISSILE) != FALSE && val == CA_N_MISSILE {
			buf = sout("personal missile attack")
		}
		if (rp_item_artifact(target).param2&CA_N_SPECIAL) != FALSE && val == CA_N_SPECIAL {
			buf = sout("personal special attack")
		}
		if (rp_item_artifact(target).param2&CA_N_MELEE_D) != FALSE &&
			val == CA_N_MELEE_D {
			buf = sout("personal melee defense")
		}
		if (rp_item_artifact(target).param2&CA_N_MISSILE_D) != FALSE &&
			val == CA_N_MISSILE_D {
			buf = sout("personal missile defense")
		}
		if (rp_item_artifact(target).param2&CA_N_SPECIAL_D) != FALSE &&
			val == CA_N_SPECIAL_D {
			buf = sout("personal special defense")
		}
		if (rp_item_artifact(target).param2&CA_M_MELEE) != FALSE && val == CA_M_MELEE {
			buf = sout("commanded men melee attack")
		}
		if (rp_item_artifact(target).param2&CA_M_MISSILE) != FALSE && val == CA_M_MISSILE {
			buf = sout("commanded men missile attack")
		}
		if (rp_item_artifact(target).param2&CA_M_SPECIAL) != FALSE && val == CA_M_SPECIAL {
			buf = sout("commanded men special attack")
		}
		if (rp_item_artifact(target).param2&CA_M_MELEE_D) != FALSE &&
			val == CA_M_MELEE_D {
			buf = sout("commanded men melee defense")
		}
		if (rp_item_artifact(target).param2&CA_M_MISSILE_D) != FALSE &&
			val == CA_M_MISSILE_D {
			buf = sout("commanded men missile defense")
		}
		if (rp_item_artifact(target).param2&CA_M_SPECIAL_D) != FALSE &&
			val == CA_M_SPECIAL_D {
			buf = sout("commanded men special defense")
		}

		if buf != "" {
			if first != FALSE {
				total = buf
				first = 0
			} else {
				total = comma_append(total, buf)
			}
		}
	}

	buf = sout("%s %s",
		header,
		artifact_names[rp_item_artifact(target).type_])
	wout(who, buf,
		rp_item_artifact(target).param1,
		total)
}

/*
 *  Thu Oct 15 18:39:28 1998 -- Scott Turner
 *
 *  Identify tells you -- if possible -- what your artifact really is.
 *
 */
func artifact_identify(header string, c *command) int {
	target := c.a
	var f string
	//var type_ int

	if rp_item_artifact(target) == nil {
		wout(c.who, "%s No enchantment", header)
		return TRUE
	}

	f = sout("%s %s",
		header,
		artifact_names[rp_item_artifact(target).type_])

	switch rp_item_artifact(target).type_ {
	case ART_COMBAT:
		describe_combat_artifact(c.who, target, header)
		break
		/*
		 *  Skills that have a box_name for param1 and nought else.
		 *
		 */
	case ART_SAFETY, ART_TRAINING, ART_SKILL, ART_PROT_SKILL:
		wout(c.who, f,
			box_name(rp_item_artifact(target).param1))
		break
		/*
		 *  Skills that have a box_name for param2 and a numeric param1
		 *
		 */
	case ART_IMPRV_DEF, ART_IMPRV_ATT, ART_SPEED_USE:
		wout(c.who, f,
			box_name(rp_item_artifact(target).param2),
			rp_item_artifact(target).param1)
		break
		/*
		 *  A subkind as param2, a numeric as param1
		 *
		 */
	case ART_TERRAIN, ART_FAST_TERR, ART_INCOME:
		wout(c.who, f,
			subkind_s[rp_item_artifact(target).param2],
			rp_item_artifact(target).param1)
		break
		/*
		 *  Skills that have a box_name for param1 and charges.
		 *
		 */
	case ART_DESTROY, ART_CROWN:
		wout(c.who, f,
			box_name(rp_item_artifact(target).param1),
			rp_item_artifact(target).uses)
		break
		/*
		 *  Skills that have a numeric for param1 and charges.
		 *
		 */
	case ART_TELEPORT:
		wout(c.who, f,
			rp_item_artifact(target).param1,
			rp_item_artifact(target).uses)
		break
		/*
		 *  Box name, numeric, charges
		 *
		 */
	case ART_SUMMON_AID:
		wout(c.who, f,
			box_name_qty(rp_item_artifact(target).param1,
				rp_item_artifact(target).param2),
			rp_item_artifact(target).uses)
		break
		/*
		 *  Just charges
		 *
		 */
	case ART_RESTORE, ART_ORB:
		wout(c.who, f,
			rp_item_artifact(target).uses)
		break
		/*
		 *  Auraculum is a special case.
		 *
		 */
	case ART_AURACULUM:
		wout(c.who, f,
			box_name(rp_item_artifact(target).param1),
			nice_num(rp_item_artifact(target).param2))
		break

		/*
		 *  Pen Crown.
		 *
		 */
	case ART_PEN:
		wout(c.who, f,
			rp_item_artifact(target).param1,
			rp_item_artifact(target).param2, rp_item_artifact(target).uses)
		break

	default:
		wout(c.who, f,
			rp_item_artifact(target).param1,
			rp_item_artifact(target).param2,
			rp_item_artifact(target).uses)
	}
	return TRUE
}

func v_identify(c *command) int {
	target := c.a

	if !valid_box(target) ||
		is_artifact(target) == nil ||
		has_item(c.who, target) == FALSE ||
		is_artifact(target).type_ == ART_AURACULUM ||
		get_effect(target, ef_obscure_artifact, 0, 0) != FALSE ||
		(target%5) == 0 {
		wout(c.who, "You are unable to identify that item.")
		return TRUE
	}

	artifact_identify("You carefully read the runes on this artifact and identify it as: ", c)

	return 0 // todo: should this return something?
}

func fix_quests(old, newQuest int) {
	var i int
	for _, i = range loop_char() {
		if only_defeatable(i) == old {
			rp_misc(i).only_vulnerable = newQuest
		}
	}
}

//#if 0
///*
// *  Tue Oct 27 10:24:08 1998 -- Scott Turner
// *
// *  This exists only to transition to the new system.
// *
// */
//void
//artifact_fixer()
//{
//  int i, new;
//  loop_subkind(sub_npc_token, i) {
//    new = create_specific_artifact(item_unique(i), ART_CROWN);
//    wout(item_unique(i),"The gods swap your %s for %s.",
//     box_name(i), box_name(new));
//    destroy_unique_item(item_unique(i), i);
//  } next_subkind;
//  loop_subkind(sub_palantir, i) {
//    new = create_specific_artifact(item_unique(i), ART_ORB);
//    wout(item_unique(i),"The gods swap your %s for %s.",
//     box_name(i), box_name(new));
//    destroy_unique_item(item_unique(i), i);
//  } next_subkind;
//  for _, i = range loop_item(i) {
//    if (i > 1000 && item_use_key(i) == use_orb) {
//      new = create_specific_artifact(item_unique(i), ART_ORB);
//      wout(item_unique(i),"The gods swap your %s for %s.",
//       box_name(i), box_name(new));
//      destroy_unique_item(item_unique(i), i);
//    };
//  }
//  loop_subkind(sub_suffuse_ring, i) {
//    new = new_suffuse_ring(item_unique(i));
//    wout(item_unique(i),"The gods swap your %s for %s.",
//     box_name(i), box_name(new));
//    destroy_unique_item(item_unique(i), i);
//  } next_subkind;
//  for _, i = range loop_item(i) {
//    if (i > 1000 &&
//    item_use_key(i) >= use_barbarian_kill &&
//    item_use_key(i) <= use_skeleton_kill) {
//      new = create_specific_artifact(item_unique(i), ART_DESTROY);
//      wout(item_unique(i),"The gods swap your %s for %s.",
//       box_name(i), box_name(new));
//      destroy_unique_item(item_unique(i), i);
//    };
//  }
//  loop_subkind(sub_artifact, i) {
//    new = create_specific_artifact(item_unique(i), ART_COMBAT);
//    wout(item_unique(i),"The gods swap your %s for %s.",
//     box_name(i), box_name(new));
//    fix_quests(i, new);
//    destroy_unique_item(item_unique(i), i);
//  } next_subkind;
//  for _, i = range loop_item(i) {
//    if (i > 1000 &&
//    item_use_key(i) == use_faery_stone) {
//      new = create_specific_artifact(item_unique(i), ART_PROT_FAERY);
//      wout(item_unique(i),"The gods swap your %s for %s.",
//       box_name(i), box_name(new));
//      destroy_unique_item(item_unique(i), i);
//    };
//  }
//};
//#endif

/*
 *  Tue Oct 27 10:31:09 1998 -- Scott Turner
 *
 *  Create and return a specific artifact.
 */
func create_specific_artifact(monster, t int) int {
	var select_ int

	for select_ = 0; artifact_tbl[select_].what != ART_NONE && artifact_tbl[select_].what != t; select_++ {
		//
	}
	assert(artifact_tbl[select_].what != ART_NONE)

	var piece int
	piece = create_unique_item(monster, sub_magic_artifact)
	set_name(piece, "Unknown ring")
	p_item(piece).weight = 5
	p_item_artifact(piece).type_ = ART_NONE
	rp_item_artifact(piece).param1 = 0
	rp_item_artifact(piece).param2 = 0
	rp_item_artifact(piece).uses = 0

	/*
	 *  No, so select something.
	 *
	 */
	p_item_artifact(piece).type_ = artifact_tbl[select_].what
	/*
	 *  Set parameter one, which might be special.
	 *
	 */
	if artifact_tbl[select_].min_param1 == artifact_tbl[select_].max_param1 {
		rp_item_artifact(piece).param1 = artifact_tbl[select_].min_param1
	} else if artifact_tbl[select_].min_param1 == RANDOM_SOLDIER {
		rp_item_artifact(piece).param1 = random_soldier()
	} else if artifact_tbl[select_].min_param1 == RANDOM_SKILL {
		rp_item_artifact(piece).param1 = random_skill()
	} else if artifact_tbl[select_].min_param1 == RANDOM_BEAST {
		rp_item_artifact(piece).param1 = random_beast(0)
	} else if artifact_tbl[select_].min_param1 == RANDOM_USE {
		rp_item_artifact(piece).param1 = random_use()
	} else {
		rp_item_artifact(piece).param1 =
			rnd(artifact_tbl[select_].min_param1, artifact_tbl[select_].max_param1)
	}
	/*
	 *  Set parameter two, no specials
	 *
	 */
	if artifact_tbl[select_].min_param2 == artifact_tbl[select_].max_param2 {
		rp_item_artifact(piece).param2 = artifact_tbl[select_].min_param2
	} else if artifact_tbl[select_].min_param2 == RANDOM_SOLDIER {
		rp_item_artifact(piece).param2 = random_soldier()
	} else if artifact_tbl[select_].min_param2 == RANDOM_SKILL {
		rp_item_artifact(piece).param2 = random_skill()
	} else if artifact_tbl[select_].min_param2 == RANDOM_BEAST {
		rp_item_artifact(piece).param2 = random_beast(0)
	} else if artifact_tbl[select_].min_param2 == RANDOM_USE {
		rp_item_artifact(piece).param2 = random_use()
	} else {
		rp_item_artifact(piece).param2 =
			rnd(artifact_tbl[select_].min_param2, artifact_tbl[select_].max_param2)
	}
	/*
	 *  Set uses, no specials
	 *
	 */
	if artifact_tbl[select_].min_uses == artifact_tbl[select_].max_uses {
		rp_item_artifact(piece).uses = artifact_tbl[select_].min_uses
	} else {
		rp_item_artifact(piece).uses =
			rnd(artifact_tbl[select_].min_uses, artifact_tbl[select_].max_uses)
	}
	/*
	 *  And special case for combat.
	 *
	 */
	if rp_item_artifact(piece).type_ == ART_COMBAT && rp_item_artifact(piece).param2 == FALSE {
		create_combat_artifact(piece)
	}
	/*
	 *  And return the artifact.
	 *
	 */
	return piece
}
