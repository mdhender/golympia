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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type skills_l []int

func (l skills_l) Len() int {
	return len(l)
}

func (l skills_l) Less(i, j int) bool {
	return rp_skill(l[i]).use_count < rp_skill(l[j]).use_count
}

func (l skills_l) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l skills_l) delete(index int) skills_l {
	var cp skills_l
	for i, e := range l {
		if i == index {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}

// rem_value removes all elements that match the value
func (l skills_l) rem_value(value int) skills_l {
	cp := l
	for i := len(cp) - 1; i >= 0; i-- {
		if e := cp[i]; e == value {
			cp = cp.delete(i)
		}
	}
	return cp
}

// rem_value_uniq removes the rightmost element in the list that matches the value
func (l skills_l) rem_value_uniq(value int) skills_l {
	for i := len(l) - 1; i >= 0; i-- {
		if e := l[i]; e == value {
			return l.delete(i)
		}
	}
	return l
}

// reverse sorts the skills by bx.temp
func (l skills_l) sort_known_comp() {
	sort.Sort(bxtmp_l(l))
}

type skill_ent_l []*skill_ent

func (l skill_ent_l) copy() skill_ent_l {
	var cp skill_ent_l
	for _, e := range l {
		cp = append(cp, e)
	}
	return cp
}

func (l skill_ent_l) delete(index int) skill_ent_l {
	var cp skill_ent_l
	for i, e := range l {
		if i == index {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}

// rem_value removes all elements that match the value
func (l skill_ent_l) rem_value(value *skill_ent) skill_ent_l {
	var cp skill_ent_l
	for _, e := range l {
		if e != value {
			continue
		}
		cp = append(cp, e)
	}
	return cp
}

type SkillList []*Skill
type Skill struct {
	Id      int    `json:"id"`             // identity of the skill
	Name    string `json:"name,omitempty"` // name of the skill
	Kind    string `json:"kind,omitempty"`
	SubKind string `json:"sub-kind,omitempty"`
}

func SkillDataLoad(name string, scanOnly bool) (SkillList, error) {
	log.Printf("SkillDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("SkillDataLoad: %w", err)
	}
	var list SkillList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("SkillDataLoad: %w", err)
	}
	if scanOnly {
		return nil, nil
	}
	for _, e := range list {
		BoxAlloc(e.Id, strKind[e.Kind], strSubKind[e.SubKind])
	}
	return nil, nil
}

func SkillDataSave(name string) error {
	var list SkillList
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("SkillDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("SkillDataSave: %w", err)
	}
	return nil
}

const (
	sk_shipcraft      = 1000
	sk_combat         = 1100
	sk_stealth        = 1200
	sk_beast          = 123
	sk_persuasion     = 1300
	sk_construction   = 1400
	sk_alchemy        = 1500
	sk_forestry       = 1600
	sk_mining         = 1700
	sk_trading        = 1800
	sk_ranger         = 1900
	sk_religion       = 150
	sk_eres           = 2000
	sk_anteus         = 2100
	sk_dol            = 2200
	sk_timeid         = 2300
	sk_ham            = 2400
	sk_kireus         = 2500
	sk_halon          = 2600
	sk_domingo        = 2700
	sk_basic          = 2800
	sk_weather        = 161
	sk_scry           = 2900
	sk_gate           = 3000
	sk_artifact       = 3100
	sk_necromancy     = 3200
	sk_adv_sorcery    = 3300
	sk_heroism        = 3400
	sk_basic_religion = 3500

	lore_skeleton_npc_token  = 9001
	lore_orc_npc_token       = 9002
	lore_undead_npc_token    = 9003
	lore_savage_npc_token    = 9004
	lore_barbarian_npc_token = 9005
	lore_orb                 = 9006
	lore_faery_stone         = 9007
	lore_barbarian_kill      = 9008
	lore_savage_kill         = 9009
	lore_undead_kill         = 9010
	lore_orc_kill            = 9011
	lore_skeleton_kill       = 9012
	lore_pen_crown           = 9013

	sk_meditate          = 2801
	sk_forge_aura        = 3130 /* forge auraculum */
	sk_mage_menial       = 2802 /* menial labor for mages */
	sk_appear_common     = 2803
	sk_view_aura         = 2804
	sk_quick_cast        = 2830 /* speed next cast */
	sk_fortify_castle    = 1491
	sk_detect_artifacts  = 3101
	sk_reveal_artifacts  = 3102
	sk_mutate_artifact   = 3131
	sk_conceal_artifacts = 3132
	sk_teleport          = 3030
	sk_obscure_artifact  = 3133
	sk_strengthen_castle = 1492
	sk_detect_gates      = 3001
	sk_jump_gate         = 3002
	sk_seal_gate         = 3031
	sk_unseal_gate       = 3032
	sk_notify_unseal     = 3033
	sk_rem_seal          = 3034 /* forcefully unseal gate */
	sk_reveal_key        = 3035
	sk_notify_jump       = 3036
	sk_reveal_mage       = 2831 /* reveal abilities of mage */
	sk_fierce_wind       = 2030
	sk_transcend_death   = 3238
	sk_tap_health        = 2832
	sk_moat_castle       = 1493
	sk_widen_entrance    = 1494
	sk_deepen_mine       = 1703
	sk_wooden_shoring    = 1790
	sk_iron_shoring      = 1791
	sk_forge_weapon      = 3134
	sk_forge_armor       = 3135
	sk_forge_bow         = 3136
	sk_bind_storm        = 9134
	sk_lightning_bolt    = 2840
	sk_foresee_defense   = 2940
	sk_drain_mana        = 3039
	sk_raise_soldiers    = 3236
	sk_fireball          = 3332
	sk_conceal_nation    = 1291
	sk_scry_region       = 2901
	sk_shroud_region     = 2930
	sk_dispel_region     = 2931 /* dispel region shroud */
	sk_remove_obscurity  = 3137
	sk_spot_hidden       = 1901
	sk_protect_noble     = 1990
	sk_write_basic       = 2833
	sk_assassinate       = 1292
	sk_find_food         = 1903
	sk_write_scry        = 2932
	sk_write_gate        = 3037
	sk_write_art         = 3138
	sk_write_necro       = 3230
	sk_prot_blast_1      = 2042
	sk_prot_blast_2      = 2139
	sk_prot_blast_3      = 2237
	sk_prot_blast_4      = 2338
	sk_prot_blast_5      = 2437
	sk_bar_loc           = 2933 /* create location barrier */
	sk_unbar_loc         = 2934
	sk_prot_blast_6      = 2536
	sk_prot_blast_7      = 2637
	sk_destroy_art       = 3103
	sk_rev_jump          = 3038
	sk_prot_blast_8      = 2735
	sk_locate_char       = 2935
	sk_deep_identify     = 3104
	sk_shroud_abil       = 2834 /* ability shroud */
	sk_detect_abil       = 2835 /* detect ability scry */
	sk_detect_scry       = 2936 /* detect region scry */
	sk_proj_cast         = 2937 /* project next cast */
	sk_dispel_abil       = 2836 /* dispel ability shroud */
	sk_adv_med           = 2837 /* advanced meditation */
	sk_hinder_med        = 2838 /* hinder meditation */
	sk_forge_palantir    = 3139
	sk_save_proj         = 2938 /* save projected cast */
	sk_save_quick        = 2839 /* save speeded cast */
	sk_summon_ghost      = 3201 /* summon ghost warriors */
	sk_raise_corpses     = 3202 /* summon undead corpses */
	sk_undead_lord       = 3231 /* summon undead unit */
	sk_renew_undead      = 3232
	sk_banish_undead     = 3233
	sk_eat_dead          = 3234
	sk_aura_blast        = 3203
	sk_absorb_blast      = 3235
	sk_summon_rain       = 2036
	sk_summon_wind       = 2037
	sk_summon_fog        = 2038
	sk_direct_storm      = 2039
	sk_dissipate         = 2031
	sk_renew_storm       = 2032
	sk_lightning         = 2033
	sk_seize_storm       = 2034
	sk_death_fog         = 2035
	sk_banish_corpses    = 2939
	sk_trance            = 3330
	sk_teleport_item     = 3331

	/* 2000 Eres */
	sk_resurrect              = 2001
	sk_pray                   = 2002
	sk_last_rites             = 2003
	sk_gather_holy_plant      = 2004
	sk_bless_follower         = 2005
	sk_proselytise            = 2006
	sk_create_holy_symbol     = 2007
	sk_heal                   = 2008
	sk_summon_water_elemental = 2040
	sk_write_weather          = 2041
	sk_dedicate_temple        = 2009

	/* 2100 Anteus */
	sk_find_mtn_trail    = 2130
	sk_obscure_mtn_trail = 2131
	sk_improve_mining    = 2132
	sk_conceal_mine      = 2133
	sk_protect_mine      = 2134
	sk_bless_fort        = 2135
	sk_weaken_fort       = 2136
	sk_boulder_trap      = 2137
	sk_write_anteus      = 2138

	/* 2200 */
	sk_detect_beasts = 2234
	sk_snake_trap    = 2235
	sk_write_dol     = 2236

	/* 2300 */
	sk_find_forest_trail    = 2330
	sk_obscure_forest_trail = 2331
	sk_improve_forestry     = 2332
	sk_reveal_forest        = 2333
	sk_improve_fort         = 2334
	sk_create_deadfall      = 2335
	sk_recruit_elves        = 2336
	sk_write_timeid         = 2337

	/* 2400 */
	sk_reveal_vision   = 2434
	sk_enchant_guard   = 2430
	sk_urchin_spy      = 2431
	sk_draw_crowds     = 2432
	sk_arrange_mugging = 2433
	sk_write_ham       = 2435
	sk_pr_shroud_loc   = 2436

	/* 2500 */
	sk_improve_quarry   = 2530
	sk_improve_smithing = 2531
	sk_edge_of_kireus   = 2532
	sk_create_mithril   = 2533
	sk_quicksand_trap   = 2534
	sk_write_kireus     = 2535

	/* 2600 */
	sk_calm_ap            = 2630
	sk_improve_charisma   = 2631
	sk_mesmerize_crowd    = 2632
	sk_improve_taxes      = 2633
	sk_guard_loyalty      = 2634
	sk_instill_fanaticism = 2635
	sk_write_halon        = 2636

	/* 2700 */
	sk_find_hidden      = 2730
	sk_conceal_loc      = 2731
	sk_mists_of_conceal = 2732
	sk_create_ninja     = 2733
	sk_write_domingo    = 2734

	sk_survive_fatal       = 3483
	sk_pilot_ship          = 1001
	sk_shipbuilding        = 1002
	sk_bird_spy            = 2231
	sk_fight_to_death      = 1102
	sk_capture_beasts      = 2232
	sk_use_beasts          = 1195 /* use beasts in battle */
	sk_breed_beasts        = 2233
	sk_petty_thief         = 1201
	sk_deep_sea            = 1094
	sk_bribe_noble         = 1301
	sk_catch_horse         = 1930
	sk_spy_inv             = 1202 /* determine char inventory */
	sk_spy_skills          = 1203 /* determine char skill */
	sk_spy_lord            = 1204 /* determine char's lord */
	sk_find_rich           = 1230
	sk_torture             = 1231
	sk_train_wild          = 1931
	sk_train_warmount      = 1932
	sk_persuade_oath       = 1330
	sk_raise_mob           = 1302
	sk_rally_mob           = 1303
	sk_incite_mob          = 1331
	sk_make_ram            = 1601 /* make battering ram */
	sk_make_catapult       = 1194
	sk_make_siege          = 1401
	sk_extract_venom       = 1530 /* from ratspider */
	sk_brew_slave          = 1531 /* potion of slavery */
	sk_brew_heal           = 1501
	sk_brew_death          = 1502
	sk_brew_weightlessness = 1590 /* potion of weightlessness */
	sk_add_ram             = 1095 /* add ram to galley */
	sk_cloak_trade         = 1232
	sk_mine_iron           = 1701
	sk_mine_gold           = 1702
	sk_mine_mithril        = 1730
	sk_quarry_stone        = 1402
	sk_mine_crystal        = 1731
	sk_harvest_lumber      = 1602
	sk_harvest_yew         = 1603
	sk_defense             = 1104
	sk_record_skill        = 1532
	sk_sneak_build         = 1233
	sk_archery             = 1902
	sk_swordplay           = 1105
	sk_weaponsmith         = 1106
	sk_fishing             = 1004
	sk_collect_foliage     = 1604
	sk_collect_elem        = 1503
	sk_summon_savage       = 1332
	sk_keep_savage         = 1333
	sk_harvest_mallorn     = 1630
	sk_harvest_opium       = 1605
	sk_improve_opium       = 1631
	sk_lead_to_gold        = 1533
	sk_hide_lord           = 1205
	sk_train_angry         = 1304
	sk_hide_self           = 1290
	sk_control_battle      = 1107
	sk_attack_tactics      = 1131
	sk_defense_tactics     = 1132
	sk_combat_discipline   = 1108
	sk_train_armor         = 1192

	sk_smuggle_goods   = 1830
	sk_smuggle_men     = 1831
	sk_avoid_taxes     = 1832
	sk_build_wagons    = 1801
	sk_increase_demand = 1802
	sk_decrease_demand = 1803
	sk_increase_supply = 1804
	sk_decrease_supply = 1805
	sk_hide_money      = 1890
	sk_hide_item       = 1891
	sk_grow_pop        = 1390
	sk_build_city      = 1490

	sk_train_knight  = 1190
	sk_train_paladin = 1191

	sk_add_sails = 1005
	sk_add_forts = 1091
	sk_add_ports = 1006
	sk_add_keels = 1090

	sk_remove_sails = 1096
	sk_remove_forts = 1093
	sk_remove_ports = 1097
	sk_remove_keels = 1092
	sk_remove_ram   = 1098
	sk_brew_fiery   = 1591

	sk_dirt_golem  = 2841
	sk_flesh_golem = 3237
	sk_iron_golem  = 3333

	/* 3400 Heroism */
	sk_swordplay2            = 3401
	sk_defense2              = 3402
	sk_survive_fatal2        = 3483
	sk_avoid_wounds          = 3480
	sk_avoid_illness         = 3481
	sk_improved_recovery     = 3482
	sk_personal_fttd         = 3403
	sk_forced_march          = 3484
	sk_extra_attacks         = 3485
	sk_extra_missile_attacks = 3486
	sk_acute_senses          = 3487
	sk_improved_explore      = 3488
	sk_uncanny_accuracy      = 3489
	sk_blinding_speed        = 3490

	/* 3500 Basic Religion */
	sk_heal_b              = 3501
	sk_last_rites_b        = 3502
	sk_resurrect_b         = 3530
	sk_create_holy_b       = 3503
	sk_dedicate_temple_b   = 3504
	sk_pray_b              = 3505
	sk_bless_b             = 3506
	sk_gather_holy_plant_b = 3507
	sk_write_religion_b    = 3508
	sk_proselytise_b       = 3509
	sk_banish_undead_b     = 3510
	sk_prot_blast_b        = 3531
	sk_hinder_med_b        = 3532
	sk_scry_b              = 3533
)

var convSkill = map[int]int{
	120:  1000,
	121:  1100,
	122:  1200,
	124:  1300,
	125:  1400,
	126:  1500,
	128:  1600,
	129:  1700,
	130:  1800,
	131:  1900,
	151:  2000,
	152:  2100,
	153:  2200,
	154:  2300,
	155:  2400,
	156:  2500,
	157:  2600,
	158:  2700,
	160:  2800,
	162:  2900,
	163:  3000,
	164:  3100,
	165:  3200,
	170:  3300,
	9101: 2801,
	9102: 3130,
	9103: 2802,
	9104: 2803,
	9105: 2804,
	9106: 2830,
	9107: 1491,
	9108: 3101,
	9109: 3102,
	9110: 3131,
	9111: 3132,
	9112: 3030,
	9113: 3133,
	9114: 1492,
	9115: 3001,
	9116: 3002,
	9117: 3031,
	9118: 3032,
	9119: 3033,
	9120: 3034,
	9121: 3035,
	9122: 3036,
	9123: 2831,
	9124: 2030,
	9125: 3238,
	9126: 2832,
	9127: 1493,
	9128: 1703,
	9129: 1790,
	9130: 1791,
	9131: 3134,
	9132: 3135,
	9133: 3136,
	9135: 2840,
	9136: 2940,
	9137: 3039,
	9138: 3236,
	9139: 3332,
	9140: 1291,
	9141: 2901,
	9142: 2930,
	9143: 2931,
	9144: 3137,
	9145: 1901,
	9146: 1990,
	9147: 2833,
	9148: 2041,
	9149: 1292,
	9150: 1903,
	9151: 2932,
	9152: 3037,
	9153: 3138,
	9154: 3230,
	9155: 2042,
	9156: 2139,
	9157: 2237,
	9158: 2338,
	9159: 2437,
	9160: 2933,
	9161: 2934,
	9162: 2536,
	9163: 2637,
	9164: 3103,
	9165: 3038,
	9166: 2735,
	9167: 2935,
	9168: 3104,
	9169: 2834,
	9170: 2835,
	9171: 2936,
	9172: 2937,
	9173: 2836,
	9174: 2837,
	9175: 2838,
	9176: 3139,
	9177: 2938,
	9178: 2839,
	9179: 3201,
	9180: 3202,
	9181: 3231,
	9182: 3232,
	9183: 3233,
	9184: 3234,
	9185: 3203,
	9186: 3235,
	9187: 2036,
	9188: 2037,
	9189: 2038,
	9190: 2039,
	9191: 2031,
	9193: 2032,
	9194: 2033,
	9195: 2034,
	9196: 2035,
	9197: 2939,
	9201: 3330,
	9202: 3331,
	9302: 2001,
	9303: 2002,
	9304: 2003,
	9305: 2004,
	9306: 2005,
	9307: 2006,
	9308: 2007,
	9309: 2008,
	9310: 2040,
	9311: 2130,
	9312: 2101,
	9313: 2102,
	9314: 2103,
	9315: 2104,
	9316: 2105,
	9317: 2106,
	9318: 2107,
	9319: 2108,
	9320: 2234,
	9321: 2235,
	9322: 2201,
	9323: 2202,
	9324: 2203,
	9325: 2204,
	9326: 2205,
	9327: 2206,
	9328: 2207,
	9329: 2208,
	9332: 2301,
	9333: 2302,
	9334: 2303,
	9335: 2304,
	9336: 2305,
	9337: 2306,
	9338: 2307,
	9339: 2308,
	9341: 2434,
	9342: 2401,
	9343: 2402,
	9344: 2403,
	9345: 2404,
	9346: 2405,
	9347: 2406,
	9348: 2407,
	9349: 2408,
	9352: 2501,
	9353: 2502,
	9354: 2503,
	9355: 2504,
	9356: 2505,
	9357: 2506,
	9358: 2507,
	9359: 2508,
	9362: 2601,
	9363: 2602,
	9364: 2603,
	9365: 2604,
	9366: 2605,
	9367: 2606,
	9368: 2607,
	9369: 2608,
	9372: 2701,
	9373: 2702,
	9374: 2703,
	9375: 2704,
	9376: 2705,
	9377: 2706,
	9378: 2707,
	9379: 2708,
	9400: 2131,
	9401: 2132,
	9402: 2133,
	9403: 2134,
	9404: 2135,
	9405: 2136,
	9406: 2137,
	9407: 2330,
	9408: 2331,
	9409: 2332,
	9410: 2333,
	9411: 2334,
	9412: 2335,
	9413: 2336,
	9414: 2530,
	9415: 2531,
	9416: 2532,
	9417: 2533,
	9418: 2534,
	9419: 2430,
	9420: 2431,
	9421: 2432,
	9422: 2433,
	9423: 2630,
	9424: 2631,
	9425: 2632,
	9426: 2633,
	9427: 2634,
	9428: 2635,
	9429: 2730,
	9430: 2731,
	9431: 2732,
	9432: 2733,
	9433: 2138,
	9434: 2236,
	9435: 2337,
	9436: 2435,
	9437: 2535,
	9438: 2636,
	9439: 2734,
	9440: 2109,
	9441: 2209,
	9442: 2309,
	9443: 2409,
	9444: 2509,
	9445: 2609,
	9446: 2709,
	9447: 2009,
	9448: 2436,
	9501: 1193,
	9502: 1001,
	9503: 1002,
	9504: 2231,
	9505: 1102,
	9506: 2232,
	9507: 1195,
	9508: 2233,
	9509: 1201,
	9510: 1094,
	9515: 1301,
	9517: 1930,
	9519: 1202,
	9520: 1203,
	9521: 1204,
	9522: 1230,
	9523: 1231,
	9529: 1931,
	9530: 1932,
	9536: 1330,
	9537: 1302,
	9538: 1303,
	9539: 1331,
	9540: 1601,
	9541: 1194,
	9542: 1401,
	9549: 1530,
	9550: 1531,
	9551: 1501,
	9552: 1502,
	9553: 1590,
	9554: 1095,
	9562: 1232,
	9563: 1701,
	9564: 1702,
	9565: 1730,
	9566: 1402,
	9567: 1731,
	9568: 1602,
	9569: 1603,
	9570: 1104,
	9573: 1532,
	9574: 1233,
	9579: 1902,
	9580: 1105,
	9581: 1106,
	9582: 1004,
	9583: 1604,
	9584: 1503,
	9585: 1332,
	9586: 1333,
	9587: 1630,
	9588: 1605,
	9589: 1631,
	9590: 1533,
	9591: 1205,
	9592: 1304,
	9593: 1290,
	9594: 1107,
	9595: 1131,
	9596: 1132,
	9598: 1108,
	9599: 1192,
	9600: 1830,
	9601: 1831,
	9602: 1832,
	9603: 1801,
	9604: 1802,
	9605: 1803,
	9606: 1804,
	9607: 1805,
	9608: 1890,
	9609: 1891,
	9610: 1390,
	9611: 1490,
	9612: 1190,
	9613: 1191,
	9614: 1005,
	9615: 1091,
	9616: 1006,
	9617: 1090,
	9618: 1096,
	9619: 1093,
	9620: 1097,
	9621: 1092,
	9622: 1098,
	9623: 1591,
	9624: 2841,
	9625: 3237,
	9626: 3333,
}
