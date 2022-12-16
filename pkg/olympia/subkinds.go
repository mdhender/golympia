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
	sub_ocean                  = 1
	sub_forest                 = 2
	sub_plain                  = 3
	sub_mountain               = 4
	sub_desert                 = 5
	sub_swamp                  = 6
	sub_under                  = 7  /* underground */
	sub_faery_hill             = 8  /* gateway to Faery */
	sub_island                 = 9  /* island subloc */
	sub_stone_cir              = 10 /* ring of stones */
	sub_mallorn_grove          = 11
	sub_bog                    = 12
	sub_cave                   = 13
	sub_city                   = 14
	sub_lair                   = 15 /* dragon lair */
	sub_graveyard              = 16
	sub_ruins                  = 17
	sub_battlefield            = 18
	sub_ench_forest            = 19 /* enchanted forest */
	sub_rocky_hill             = 20
	sub_tree_circle            = 21
	sub_pits                   = 22
	sub_pasture                = 23
	sub_oasis                  = 24
	sub_yew_grove              = 25
	sub_sand_pit               = 26
	sub_sacred_grove           = 27
	sub_poppy_field            = 28
	sub_temple                 = 29
	sub_galley                 = 30
	sub_roundship              = 31
	sub_castle                 = 32
	sub_galley_notdone         = 33
	sub_roundship_notdone      = 34
	sub_ghost_ship             = 35
	sub_temple_notdone         = 36
	sub_inn                    = 37
	sub_inn_notdone            = 38
	sub_castle_notdone         = 39
	sub_mine                   = 40
	sub_mine_notdone           = 41
	sub_scroll                 = 42 /* item is a scroll */
	sub_magic                  = 43 /* this skill is magical */
	sub_palantir               = 44
	sub_auraculum              = 45
	sub_tower                  = 46
	sub_tower_notdone          = 47
	sub_pl_system              = 48 /* system player */
	sub_pl_regular             = 49 /* regular player */
	sub_region                 = 50 /* region wrapper loc */
	sub_pl_savage              = 51 /* Savage King */
	sub_pl_npc                 = 52
	sub_mine_collapsed         = 53
	sub_ni                     = 54 /* ni=noble_item */
	sub_demon_lord             = 55 /* undead lord */
	sub_dead_body              = 56 /* dead noble's body */
	sub_fog                    = 57
	sub_wind                   = 58
	sub_rain                   = 59
	sub_hades_pit              = 60
	sub_artifact               = 61
	sub_pl_silent              = 62
	sub_npc_token              = 63 /* npc group control art */
	sub_garrison               = 64 /* npc group control art */
	sub_cloud                  = 65 /* cloud terrain type */
	sub_raft                   = 66 /* raft made out of flotsam */
	sub_raft_notdone           = 67
	sub_suffuse_ring           = 68
	sub_religion               = 69
	sub_holy_symbol            = 70 /* Holy symbol of some sort */
	sub_mist                   = 71
	sub_book                   = 72 /* Book */
	sub_guild                  = 73 /* Requires skill to enter */
	sub_trade_good             = 74
	sub_city_notdone           = 75
	sub_ship                   = 76
	sub_ship_notdone           = 77
	sub_mine_shaft             = 78
	sub_mine_shaft_notdone     = 79
	sub_orc_stronghold         = 80
	sub_orc_stronghold_notdone = 81
	sub_special_staff          = 82
	sub_lost_soul              = 83
	sub_undead                 = 84
	sub_pen_crown              = 85
	sub_animal_part            = 86
	sub_magic_artifact         = 87

	SUB_MAX = 88 /* one past highest sub_ */
)

var subKindStr = map[int]string{
	sub_ocean:                  "ocean",
	sub_forest:                 "forest",
	sub_plain:                  "plain",
	sub_mountain:               "mountain",
	sub_desert:                 "desert",
	sub_swamp:                  "swamp",
	sub_under:                  "underground",
	sub_faery_hill:             "faery hill",
	sub_island:                 "island",
	sub_stone_cir:              "ring of stones",
	sub_mallorn_grove:          "mallorn grove",
	sub_bog:                    "bog",
	sub_cave:                   "cave",
	sub_city:                   "city",
	sub_lair:                   "lair",
	sub_graveyard:              "graveyard",
	sub_ruins:                  "ruins",
	sub_battlefield:            "field",
	sub_ench_forest:            "enchanted forest",
	sub_rocky_hill:             "rocky hill",
	sub_tree_circle:            "circle of trees",
	sub_pits:                   "pits",
	sub_pasture:                "pasture",
	sub_oasis:                  "oasis",
	sub_yew_grove:              "yew grove",
	sub_sand_pit:               "sand pit",
	sub_sacred_grove:           "sacred grove",
	sub_poppy_field:            "poppy field",
	sub_temple:                 "temple",
	sub_galley:                 "galley",
	sub_roundship:              "roundship",
	sub_castle:                 "castle",
	sub_galley_notdone:         "galley-in-progress",
	sub_roundship_notdone:      "roundship-in-progress",
	sub_ghost_ship:             "ghost ship",
	sub_temple_notdone:         "temple-in-progress",
	sub_inn:                    "inn",
	sub_inn_notdone:            "inn-in-progress",
	sub_castle_notdone:         "castle-in-progress",
	sub_mine:                   "mine",
	sub_mine_notdone:           "mine-in-progress",
	sub_scroll:                 "scroll",
	sub_magic:                  "magic",
	sub_palantir:               "palantir",
	sub_auraculum:              "auraculum",
	sub_tower:                  "tower",
	sub_tower_notdone:          "tower-in-progress",
	sub_pl_system:              "pl_system",
	sub_pl_regular:             "pl_regular",
	sub_region:                 "region",
	sub_pl_savage:              "pl_savage",
	sub_pl_npc:                 "pl_npc",
	sub_mine_collapsed:         "collapsed mine",
	sub_ni:                     "ni",
	sub_demon_lord:             "demon lord",
	sub_dead_body:              "dead body",
	sub_fog:                    "fog",
	sub_wind:                   "wind",
	sub_rain:                   "rain",
	sub_hades_pit:              "pit",
	sub_artifact:               "artifact",
	sub_pl_silent:              "pl_silent",
	sub_npc_token:              "npc_token",
	sub_garrison:               "garrison",
	sub_cloud:                  "cloud",
	sub_raft:                   "raft",
	sub_raft_notdone:           "raft-in-progress",
	sub_suffuse_ring:           "suffuse_ring",
	sub_religion:               "religion",
	sub_holy_symbol:            "holy symbol",
	sub_mist:                   "mist",
	sub_book:                   "book",
	sub_guild:                  "guild", // mdhender: was sub_market
	sub_trade_good:             "trade_good",
	sub_city_notdone:           "city-in-progress",
	sub_ship:                   "ship",
	sub_ship_notdone:           "ship-in-progress",
	sub_mine_shaft:             "mine-shaft",
	sub_mine_shaft_notdone:     "mine-shaft-in-progress",
	sub_orc_stronghold:         "orc-stronghold",
	sub_orc_stronghold_notdone: "orc-stronghold-in-progress",
	sub_special_staff:          "Staff-of-the-Sun",
	sub_lost_soul:              "lost_soul",
	sub_undead:                 "undead",
	sub_pen_crown:              "pen-crown",
	sub_animal_part:            "animal-part",
	sub_magic_artifact:         "magical-artifact",
	0:                          "<no subkind>",
}

var strSubKind = map[string]int{
	"ocean":                      sub_ocean,
	"forest":                     sub_forest,
	"plain":                      sub_plain,
	"mountain":                   sub_mountain,
	"desert":                     sub_desert,
	"swamp":                      sub_swamp,
	"underground":                sub_under,
	"faery hill":                 sub_faery_hill,
	"island":                     sub_island,
	"ring of stones":             sub_stone_cir,
	"mallorn grove":              sub_mallorn_grove,
	"bog":                        sub_bog,
	"cave":                       sub_cave,
	"city":                       sub_city,
	"lair":                       sub_lair,
	"graveyard":                  sub_graveyard,
	"ruins":                      sub_ruins,
	"field":                      sub_battlefield,
	"enchanted forest":           sub_ench_forest,
	"rocky hill":                 sub_rocky_hill,
	"circle of trees":            sub_tree_circle,
	"pits":                       sub_pits,
	"pasture":                    sub_pasture,
	"oasis":                      sub_oasis,
	"yew grove":                  sub_yew_grove,
	"sand pit":                   sub_sand_pit,
	"sacred grove":               sub_sacred_grove,
	"poppy field":                sub_poppy_field,
	"temple":                     sub_temple,
	"galley":                     sub_galley,
	"roundship":                  sub_roundship,
	"castle":                     sub_castle,
	"galley-in-progress":         sub_galley_notdone,
	"roundship-in-progress":      sub_roundship_notdone,
	"ghost ship":                 sub_ghost_ship,
	"temple-in-progress":         sub_temple_notdone,
	"inn":                        sub_inn,
	"inn-in-progress":            sub_inn_notdone,
	"castle-in-progress":         sub_castle_notdone,
	"mine":                       sub_mine,
	"mine-in-progress":           sub_mine_notdone,
	"scroll":                     sub_scroll,
	"magic":                      sub_magic,
	"palantir":                   sub_palantir,
	"auraculum":                  sub_auraculum,
	"tower":                      sub_tower,
	"tower-in-progress":          sub_tower_notdone,
	"pl_system":                  sub_pl_system,
	"pl_regular":                 sub_pl_regular,
	"region":                     sub_region,
	"pl_savage":                  sub_pl_savage,
	"pl_npc":                     sub_pl_npc,
	"collapsed mine":             sub_mine_collapsed,
	"ni":                         sub_ni,
	"demon lord":                 sub_demon_lord,
	"dead body":                  sub_dead_body,
	"fog":                        sub_fog,
	"wind":                       sub_wind,
	"rain":                       sub_rain,
	"pit":                        sub_hades_pit,
	"artifact":                   sub_artifact,
	"pl_silent":                  sub_pl_silent,
	"npc_token":                  sub_npc_token,
	"garrison":                   sub_garrison,
	"cloud":                      sub_cloud,
	"raft":                       sub_raft,
	"raft-in-progress":           sub_raft_notdone,
	"suffuse_ring":               sub_suffuse_ring,
	"religion":                   sub_religion,
	"holy symbol":                sub_holy_symbol,
	"mist":                       sub_mist,
	"book":                       sub_book,
	"guild":                      sub_guild, // mdhender: was sub_market
	"trade_good":                 sub_trade_good,
	"city-in-progress":           sub_city_notdone,
	"ship":                       sub_ship,
	"ship-in-progress":           sub_ship_notdone,
	"mine-shaft":                 sub_mine_shaft,
	"mine-shaft-in-progress":     sub_mine_shaft_notdone,
	"orc-stronghold":             sub_orc_stronghold,
	"orc-stronghold-in-progress": sub_orc_stronghold_notdone,
	"Staff-of-the-Sun":           sub_special_staff,
	"lost_soul":                  sub_lost_soul,
	"undead":                     sub_undead,
	"pen-crown":                  sub_pen_crown,
	"animal-part":                sub_animal_part,
	"magical-artifact":           sub_magic_artifact,
	"<no subkind>":               0,
}
