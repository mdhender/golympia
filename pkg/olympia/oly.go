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

// oly.h

const (
	DEFAULT_PASSWORD = "defpwd123"

	MAX_BOXES  = 102400
	MONTH_DAYS = 30
	NUM_MONTHS = 8

	T_deleted  = 0 /* forget on save */
	T_player   = 1
	T_char     = 2
	T_loc      = 3
	T_item     = 4
	T_skill    = 5
	T_gate     = 6
	T_road     = 7
	T_deadchar = 8
	T_ship     = 9
	T_post     = 10
	T_storm    = 11
	T_unform   = 12 /* unformed noble */
	T_lore     = 13
	T_nation   = 14
	T_MAX      = 15 /* one past highest T_xxx define */

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

	item_gold = 1

	item_peasant         = 10
	item_worker          = 11
	item_soldier         = 12
	item_archer          = 13
	item_knight          = 14
	item_elite_guard     = 15
	item_pikeman         = 16
	item_blessed_soldier = 17
	item_ghost_warrior   = 18
	item_sailor          = 19
	item_swordsman       = 20
	item_crossbowman     = 21
	item_elite_arch      = 22
	item_angry_peasant   = 23
	item_pirate          = 24
	item_elf             = 25
	item_spirit          = 26
	item_postulant       = 27
	item_fanatic         = 28
	item_ninja           = 29
	item_angel           = 30

	item_corpse       = 31
	item_savage       = 32
	item_skeleton     = 33
	item_barbarian    = 34
	item_wagon        = 35
	item_skirmisher   = 36
	item_hvy_foot     = 37
	item_hvy_xbowman  = 38
	item_elvish_arrow = 39
	item_hvy_xbow     = 40
	item_horse_archer = 41
	item_pctg_token   = 42
	item_cavalier     = 43
	item_new_wagon    = 44
	item_hvy_wagon    = 45
	item_war_wagon    = 46

	item_wild_horse   = 51
	item_riding_horse = 52
	item_warmount     = 53
	item_pegasus      = 54
	item_nazgul       = 55

	item_flotsam         = 59
	item_battering_ram   = 60
	item_catapult        = 61
	item_siege_tower     = 62
	item_ratspider_venom = 63
	item_lana_bark       = 64
	item_avinia_leaf     = 65
	item_spiny_root      = 66
	item_farrenstone     = 67
	item_yew             = 68
	item_elfstone        = 69
	item_mallorn_wood    = 70
	item_pretus_bones    = 71
	item_longbow         = 72
	item_plate           = 73
	item_longsword       = 74
	item_pike            = 75
	item_ox              = 76
	item_lumber          = 77
	item_stone           = 78
	item_iron            = 79
	item_leather         = 80
	item_ratspider       = 81
	item_mithril         = 82
	item_gate_crystal    = 83
	item_blank_scroll    = 84
	item_crossbow        = 85
	item_rug             = 86
	item_fish            = 87
	item_pepper          = 88
	item_pipeweed        = 89
	item_ale             = 90
	item_wool            = 91
	item_jewel           = 92
	item_opium           = 93
	item_basket          = 94 /* woven basket */
	item_pot             = 95 /* clay pot */
	item_tax_cookie      = 96
	item_fish_oil        = 97
	item_drum            = 98
	item_hide            = 99
	item_mob_cookie      = 101
	item_lead            = 102
	item_fine_cloak      = 103
	item_chocolate       = 104
	item_ivory           = 105
	item_cardamom        = 106
	item_honey           = 107
	item_ink             = 108
	item_licorice        = 109
	item_soap            = 110
	item_old_book        = 111
	item_jade_idol       = 112
	item_purple_cloth    = 113
	item_rose_perfume    = 114
	item_silk            = 115
	item_incense         = 116
	item_ochre           = 117
	item_jeweled_egg     = 118
	item_obsidian        = 119
	item_orange          = 251
	item_cinnabar        = 252
	item_myrhh           = 253
	item_saffron         = 254
	item_dried_fish      = 255
	item_tallow          = 256
	item_candles         = 257
	item_wax             = 258
	item_sugar           = 259
	item_salt            = 260
	item_glue            = 261
	item_linen           = 262
	item_beans           = 263
	item_walnuts         = 264
	item_flax            = 265
	item_flutes          = 266
	item_cassava         = 267
	item_plum_wine       = 268
	item_vinegar         = 269
	item_tea             = 270
	item_centaur         = 271
	item_minotaur        = 272
	item_undead_cookie   = 273
	item_fog_cookie      = 274
	item_wind_cookie     = 275
	item_rain_cookie     = 276
	item_mage_menial     = 277 /* mage menial labor cookie */
	item_spider          = 278 /* giant spider */
	item_rat             = 279 /* horde of rats */
	item_lion            = 280
	item_bird            = 281 /* giant bird */
	item_lizard          = 282
	item_bandit          = 283
	item_chimera         = 284
	item_harpie          = 285
	item_dragon          = 286
	item_orc             = 287
	item_gorgon          = 288
	item_wolf            = 289
	item_orb             = 290
	item_cyclops         = 291
	item_giant           = 292
	item_faery           = 293
	item_petty_thief     = 294
	item_seagrass        = 295
	item_firewort        = 296
	item_beastnip        = 297
	item_elf_poppy       = 298
	item_ironwood        = 299
	item_kings_fern      = 300
	item_moon_palms      = 301
	item_otter           = 302
	item_mole            = 303
	item_bull            = 304
	item_eagle           = 305
	item_monkey          = 306
	item_hare            = 307
	item_wardog          = 308
	item_sand_rat        = 309
	item_balrog          = 310
	item_dirt_golem      = 311
	item_flesh_golem     = 312
	item_iron_golem      = 313
	item_lesser_demon    = 314
	item_greater_demon   = 315
	item_green_rose      = 316
	item_elf_ear         = 317
	item_savage_ear      = 318
	item_nazgul_tail     = 319
	item_centaur_hide    = 320
	item_minotaur_hide   = 321
	item_spider_eye      = 322
	item_rat_tail        = 323
	item_lion_mane       = 324
	item_bird_feather    = 325
	item_lizard_tail     = 326
	item_bandit_ear      = 327
	item_chimera_eye     = 328
	item_harpie_feather  = 329
	item_dragon_scale    = 330
	item_orc_scalp       = 331
	item_gorgon_liver    = 332
	item_wolf_hide       = 333
	item_cyclops_eye     = 334
	item_giant_tongue    = 335
	item_balrog_horn     = 336

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

	PROG_bandit         = 1 /* wilderness spice */
	PROG_subloc_monster = 2
	PROG_npc_token      = 3
	PROG_balrog         = 4
	PROG_dumb_monster   = 5
	PROG_smart_monster  = 6
	PROG_orc            = 7
	PROG_elf            = 8
	PROG_daemon         = 9

	use_death_potion          = 1
	use_heal_potion           = 2
	use_slave_potion          = 3
	use_palantir              = 4
	use_proj_cast             = 5 /* stored projected cast */
	use_quick_cast            = 6 /* stored cast speedup */
	use_drum                  = 7 /* beat savage's drum */
	use_faery_stone           = 8 /* Faery gate opener */
	use_orb                   = 9 /* crystal orb */
	use_barbarian_kill        = 10
	use_savage_kill           = 11
	use_corpse_kill           = 12
	use_orc_kill              = 13
	use_skeleton_kill         = 14
	use_ancient_aura          = 15 /* bta's auraculum */
	use_weightlessness_potion = 16
	use_fiery_potion          = 17 /* Fiery Death */
	use_nothing               = 18 /* Something that does nothing */
	use_special_staff         = 19

	DIR_N    = 1
	DIR_E    = 2
	DIR_S    = 3
	DIR_W    = 4
	DIR_UP   = 5
	DIR_DOWN = 6
	DIR_IN   = 7
	DIR_OUT  = 8
	MAX_DIR  = 9 /* one past highest direction */

	LOC_region   = 1 /* top most continent/island group */
	LOC_province = 2 /* main location area */
	LOC_subloc   = 3 /* inner sublocation */
	LOC_build    = 4 /* building, structure, etc. */

	LOY_UNCHANGED = (-1)
	LOY_unsworn   = 0
	LOY_contract  = 1
	LOY_oath      = 2
	LOY_fear      = 3
	LOY_npc       = 4
	LOY_summon    = 5

	exp_novice     = 1 /* apprentice */
	exp_journeyman = 2
	exp_teacher    = 3
	exp_master     = 4
	exp_grand      = 5 /* grand master */

	S_body    = 1 /* kill to a dead body */
	S_soul    = 2 /* kill to a lost soul */
	S_nothing = 3 /* kill completely */

	MONSTER_ATT = -1 /* phony id for "all monsters" */
	ATT_NONE    = 0  /* no attitude -- default */
	NEUTRAL     = 1  /* explicitly neutral */
	HOSTILE     = 2  /* attack on sight */
	DEFEND      = 3  /* defend if attacked */

	// effect identifiers
	ef_defense            = 1001 /* Add to defense of stack. */
	ef_religion_trap      = 1002 /* A religion trap. */
	ef_fast_move          = 1004 /* Fast move into the next province */
	ef_slow_move          = 1005 /* Slow moves into this province */
	ef_improve_mine       = 1006 /* Improve a mine's production */
	ef_protect_mine       = 1007 /* Protect a mine against calamities */
	ef_bless_fort         = 1008 /* Bless a fort */
	ef_improve_production = 1010 /* Generic production improvement +50%*/
	ef_improve_make       = 1011 /* Generic make improvement +100% */
	ef_improve_fort       = 1012 /* Increase a fortification */
	ef_edge_of_kireus     = 1013 /* Edged weapons +25 attack */
	ef_urchin_spy         = 1014 /* Report unit's doings for 7 days */
	ef_charisma           = 1015 /* Double leadership */
	ef_improve_taxes      = 1016 /* Double taxes */
	ef_guard_loyalty      = 1017 /* Unit immune to loyalty checks */
	ef_hide_item          = 1018 /* Hide an item. */
	ef_hide_money         = 1019 /* Hide money. */
	ef_smuggle_goods      = 1020 /* Smuggling. */
	ef_smuggle_men        = 1021 /* Ditto. */
	ef_grow               = 1022 /* Encourage pop to grow. */
	ef_weightlessness     = 1023 /* Negative weight */
	ef_scry_offense       = 1024 /* + to offense in battle */
	ef_conceal_nation     = 1025 /* conceal your nation */
	ef_kill_dirt_golem    = 1026 /* Timer on golems */
	ef_kill_flesh_golem   = 1027 /* Timer on golems */
	ef_food_found         = 1028 /* Food found by a ranger */
	ef_conceal_artifacts  = 1029 /* Conceal someone's artifacts */
	ef_obscure_artifact   = 1030 /* Obscure artifact's identity */
	ef_forced_march       = 1031 /* Next move at riding speed */
	ef_faery_warning      = 1032 /* Give people a month grace. */
	ef_tap_wound          = 1033 /* The wound from Tap Health. */
	ef_magic_barrier      = 1034 /* New implementation of magic barrier */
	ef_inhibit_barrier    = 1035 /* After magic barrier falls. */
	ef_cs                 = 1036 /* Combat spell control. */

	// artifacts - generalized magic items -- not unlike effects.
	ART_NONE        = 0  /* No effect */
	ART_COMBAT      = 1  /* + to combat */
	ART_CTL_MEN     = 2  /* + to controlled men */
	ART_CTL_BEASTS  = 3  /* + to controlled beasts */
	ART_SAFETY      = 4  /* safety from one monster type */
	ART_IMPRV_ATT   = 5  /* improve attack of particular unit */
	ART_IMPRV_DEF   = 6  /* improve defense of particular unit */
	ART_SAFE_SEA    = 7  /* safety at sea */
	ART_TERRAIN     = 8  /* improved defense in particular terrain */
	ART_FAST_TERR   = 9  /* fast movement into particular terrain */
	ART_SPEED_USE   = 10 /* speed use of skill */
	ART_PROT_HADES  = 11 /* protection in Hades */
	ART_PROT_FAERY  = 12 /* protection in Faery */
	ART_WORKERS     = 13 /* improve productivity of workers */
	ART_INCOME      = 14 /* improved income from X */
	ART_LEARNING    = 15 /* improved learning */
	ART_TEACHING    = 16 /* improved teaching */
	ART_TRAINING    = 17 /* improved training (of X?) */
	ART_DESTROY     = 18 /* destroy monster X (uses) */
	ART_SKILL       = 19 /* grants a skill */
	ART_FLYING      = 20 /* permit user to fly */
	ART_PROT_SKILL  = 21 /* protection from a skill (aura blast,scry) */
	ART_SHIELD_PROV = 22 /* protect entire province from scrying */
	ART_RIDING      = 23 /* permit user riding pace */
	ART_POWER       = 24 /* increase piety/aura (1 use) */
	ART_SUMMON_AID  = 25 /* summon help (1 use) */
	ART_MAINTENANCE = 26 /* reduce maintenance costs */
	ART_BARGAIN     = 27 /* better market prices */
	ART_WEIGHTLESS  = 28 /* weightlessness */
	ART_HEALING     = 29 /* faster healing */
	ART_SICKNESS    = 30 /* protection from sickness */
	ART_RESTORE     = 31 /* restore life */
	ART_TELEPORT    = 32 /* teleportation (uses) */
	ART_ORB         = 33 /* orb */
	ART_CROWN       = 34 /* crown */
	ART_AURACULUM   = 35 /* auraculum */
	ART_CARRY       = 36 /* increase land carry capacity */
	ART_PEN         = 37 /* the Pen Crown */
	ART_LAST        = 38

	// bit-masks to pull out the appropriate bits of combat artifacts.
	CA_N_MELEE     = (1 << 0)
	CA_N_MISSILE   = (1 << 1)
	CA_N_SPECIAL   = (1 << 2)
	CA_M_MELEE     = (1 << 3)
	CA_M_MISSILE   = (1 << 4)
	CA_M_SPECIAL   = (1 << 5)
	CA_N_MELEE_D   = (1 << 6)
	CA_N_MISSILE_D = (1 << 7)
	CA_N_SPECIAL_D = (1 << 8)
	CA_M_MELEE_D   = (1 << 9)
	CA_M_MISSILE_D = (1 << 10)
	CA_M_SPECIAL_D = (1 << 11)

	TOUGH_NUM = 2520

	SKILL_dont     = 0 /* don't know the skill */
	SKILL_learning = 1 /* in the process of learning it */
	SKILL_know     = 2 /* know it */

	// build descriptor - describes what kind of build is going on in a location.
	BT_BUILD      = 1
	BT_FORTIFY    = 2
	BT_STRENGTHEN = 3
	BT_MOAT       = 4

	// mine descriptors
	MINE_MAX       = 20
	NO_SHORING     = 0
	WOODEN_SHORING = 1
	IRON_SHORING   = 2

	/*
	 *  Skill flags.
	 *
	 */
	IS_POLLED       = (1 << 0)
	REQ_HOLY_SYMBOL = (1 << 1)
	REQ_HOLY_PLANT  = (1 << 2)
	COMBAT_SKILL    = (1 << 3)
	MAX_FLAGS       = 4

	REQ_NO  = 0 /* don't consume item */
	REQ_YES = 1 /* consume item */
	REQ_OR  = 2 /* or with next */

	// in-process command structure
	DONE  = 0
	LOAD  = 1
	RUN   = 2
	ERROR = 3

	// types of command arguments
	CMD_undef    = 0
	CMD_unit     = 1
	CMD_item     = 2
	CMD_skill    = 3
	CMD_days     = 4
	CMD_qty      = 5
	CMD_gold     = 6
	CMD_use      = 7
	CMD_practice = 8

	// trade flags
	BUY     = 1
	SELL    = 2
	PRODUCE = 3
	CONSUME = 4

	CHAR_FIELD = 5  /* field length for box_code_less */
	MAX_POST   = 60 /* max line length for posts and messages */

	// used in combat.c and move.c
	A_WON     = 1
	B_WON     = 2
	TIE       = 3
	NO_COMBAT = 4

	// defines a faction as either MUs or Priests.
	MU_FACTION     = 1
	PRIEST_FACTION = 2

	// define format flags
	NONE = 0
	HTML = (1 << 0)
	TEXT = (1 << 1)
	TAGS = (1 << 2)
	RAW  = (1 << 3)
	ALT  = (1 << 4)
)

var (
	// global box

	LINES = "---------------------------------------------"

	dir_s           []string
	evening         bool /* are we in the evening phase? */
	indent          int
	trades_to_check []int
)

type accept_ent struct {
	item     int /* 0 = any item */
	from_who int /* 0 = anyone, else char or player */
	qty      int /* 0 = any qty */
}

type admit struct {
	targ  int /* char or loc admit is declared for */
	sense int /* 0=default no, 1=all but.. */
	l     []int
	flag  int /* first time set this turn -- not saved */
}

type att_ent struct {
	neutral []int
	hostile []int
	defend  []int
}

type box struct {
	kind       schar
	skind      schar
	name       string
	x_loc_info loc_info
	x_player   *entity_player
	x_char     *entity_char
	x_loc      *entity_loc
	x_subloc   *entity_subloc
	x_item     *entity_item
	x_skill    *entity_skill
	x_nation   *entity_nation
	x_gate     *entity_gate
	x_misc     *entity_misc
	x_disp     *att_ent

	cmd     *command
	items   item_ent_l /* ilist of items held */
	trades  trade_l    /* pending buys/sells */
	effects []*effect  /* ilist of effects */

	temp         int /* scratch space */
	output_order int /* for report ordering -- not saved */

	x_next_kind int /* link to next entity of same type */
	x_next_sub  int /* link to next of same subkind */
}

type char_magic struct {
	max_aura  int /* maximum aura level for magician */
	cur_aura  int /* current aura level for magician */
	auraculum int /* char created an auraculum */

	visions sparse /* visions revealed */
	// pledge  int    /* lands are pledged to another */ // not used?
	token int /* we are controlled by this art */

	project_cast      int /* project next cast */
	quick_cast        int /* speed next cast */
	ability_shroud    int
	hide_mage         int // number of points hiding the magician
	hinder_meditation int // number of points to hinder, usually 0...3
	magician          int /* is a magician */
	aura_reflect      int /* reflect aura blast */
	hide_self         int /* character is hidden */
	swear_on_release  int /* swear to one who frees us */
	knows_weather     int /* knows weather magic */

	mage_worked   int   /* worked this month -- not saved */
	ferry_flag    bool  /* ferry has tooted its horn -- ns */
	pledged_to_us []int /* temp -- not saved */
}

// character religion - c aptures a nobles religious standing, such as it is.
type char_religion struct {
	priest    int   /* Who this noble is dedicated to, if anyone. */
	piety     int   /* Our current piety. */
	followers []int /* Who is dedicated to us, if anyone. */
}

type cmd_tbl_ent struct {
	allow string /* who may execute the command */
	name  string /* name of command */

	start     func(c *command) int /* initiator */
	finish    func(c *command) int /* conclusion */
	interrupt func(c *command) int /* interrupted order */

	time int /* how long command takes */
	poll int /* call finish each day, not just at end */
	pri  int /* command priority or precedence */

	// mods to add some command checking during the eat phase:
	//     -- Num_args_required
	//     -- Max_args
	//     -- Arg_types[]
	//     -- cmd_comment()
	//     -- cmd_check()
	num_args_required int
	max_args          int
	arg_types         [5]int
	cmd_comment       func(c *command) string
	cmd_check         func(c *command)
}

type command struct {
	who                          int         // entity this is under (redundant)
	wait                         int         // time until completion
	cmd                          int         // index into cmd_tbl
	use_skill                    int         // skill we are using, if any
	use_ent                      int         // index into use_tbl[] for skill usage
	use_exp                      int         // experience level at using this skill
	days_executing               int         // how long has this command been running
	a, b, c, d, e, f, g, h, i, j int         // command arguments
	line                         string      // original command line
	parsed_line                  []byte      // cut-up line, pointed to by parse
	parse                        args_l      // ilist of parsed arguments
	state                        int         // LOAD, RUN, ERROR, DONE
	status                       int         // success or failure; can this be bool?
	poll                         int         // call finish routine each day?
	pri                          int         // command priority or precedence
	conditional                  schar       // 0=none 1=last succeeded 2=last failed
	inhibit_finish               bool        // don't call d_xxx
	fuzzy                        bool        // command matched fuzzy -- not saved
	second_wait                  int         // delay resulting from auto attacks -- saved
	wait_parse                   []*wait_arg // not saved
	// debug is not a bool because it is used to hold multiple values
	debug int // debugging check -- not saved
}

type entity_artifact struct {
	type_  int
	param1 int
	param2 int
	uses   int
}

// describes what kind of build is going on in a location.
type entity_build struct {
	type_           int /* What work is going on? */
	build_materials int /* fifths of materials we've used */
	effort_required int /* not finished if nonzero */
	effort_given    int
}

type entity_char struct {
	unit_item int /* unit is made of this kind of item */

	health int
	sick   int /* 1=character is getting worse */

	guard    int /* character is guarding the loc */
	loy_kind int /* LOY_xxx */
	loy_rate int /* level with kind of loyalty */

	death_time olytime /* when was character killed */

	skills skill_ent_l /* ilist of skills known by char */

	// effects []*effect /* ilist of effects on char */ // not used?

	moving    int /* daystamp of beginning of movement */
	unit_lord int /* who is our owner? */

	contact []int /* who have we contacted, also, who has found us */

	x_char_magic         *char_magic
	prisoner             int /* is this character a prisoner? */
	behind               int /* are we behind in combat? */
	time_flying          int /* time airborne over ocean */
	break_point          int /* break point when fighting */
	personal_break_point int /* personal break point when fighting */
	rank                 int /* noble peerage status */
	npc_prog             int /* npc program */

	guild int /* This is the guild we belong to. */

	attack  int /* fighter attack rating */
	defense int /* fighter defense rating */
	missile int /* capable of missile attacks? */

	religion char_religion /* Our religion info... */

	pay int /* How much will you pay to enter? */

	// the following are not saved by io.c:
	melt_me    int           /* in process of melting away */
	fresh_hire int           /* don't erode loyalty */
	new_lord   int           /* got a new lord this turn */
	studied    int           /* num days we studied */
	accept     []*accept_ent /* what we can be given */
}

type entity_gate struct {
	to_loc        int /* destination of gate */
	notify_jumps  int /* whom to notify */
	notify_unseal int /* whom to notify */
	seal_key      int /* numeric gate password */
	road_hidden   int /* this is a hidden road or passage */
}

type entity_item struct {
	weight      int
	land_cap    int
	ride_cap    int
	fly_cap     int
	attack      int /* fighter attack rating */
	defense     int /* fighter defense rating */
	missile     int /* capable of missile attacks? */
	maintenance int /* Maintenance cost */
	npc_split   int /* Size to "split" at... */
	animal_part int /* Produces this when killed. */

	is_man_item int /* unit is a character like thing */
	animal      int /* unit is or contains a horse or an ox */
	prominent   int /* big things that everyone sees */
	capturable  int /* ni-char contents are capturable */
	ungiveable  int /* Can't be transferred between nobles. */
	wild        int /* Appears in the wild as a random encounter. */
	/* Value is actually the NPC_prog. */

	plural_name string
	base_price  int /* base price of item for market seeding */
	trade_good  int /* Is this thing a trade good? & how much*/
	who_has     int /* who has this unique item */

	x_item_magic    *item_magic
	x_item_artifact *entity_artifact /* Eventually replace item_magic */
}

type entity_loc struct {
	prov_dest       []int /* cached exits for flood fills */
	near_grave      int   /* nearest graveyard */
	shroud          int   /* magical scry shroud */
	barrier         int   /* magical barrier */
	tax_rate        int   /* Tax rate for this loc. */
	recruited       int   /* How many recruited this month */
	hidden          int   /* is location hidden? */
	dist_from_sea   int   /* provinces to sea province */
	dist_from_swamp int
	dist_from_gate  int
	sea_lane        int /* fast ocean travel here, also "tracks" for npc ferries */
	// effects []*effect        /* ilist of effects on location */ // not used?
	mine_info *entity_mine /* If there's a mine. */
	// location control -- need two so that we can only change fees at the end of the month.
	control  loc_control_ent
	control2 loc_control_ent // doesn't need to be saved
}

type entity_mine struct {
	mc      [MINE_MAX]mine_contents // todo: shouldn't be static size
	shoring [MINE_MAX]int           // todo: shouldn't be static size
}

type entity_misc struct {
	display         string /* entity display banner */
	npc_created     int    /* turn peasant mob created */
	npc_home        int    /* where npc was created */
	npc_cookie      int    /* allocation cookie item for us */
	summoned_by     int    /* who summoned us? */
	save_name       string /* orig name of noble for dead bodies */
	old_lord        int    /* who did this dead body used to belong to */
	npc_memory      sparse
	only_vulnerable int  /* only defeatable with this rare artifact */
	garr_castle     int  /* castle which owns this garrison */
	border_open     int  /* Is the garrison's border open? */
	bind_storm      int  /* storm bound to this ship */
	storm_str       int  /* storm strength */
	npc_dir         int  /* last direction npc moved */
	mine_delay      int  /* time until collapsed mine vanishes */
	cmd_allow       byte /* unit under restricted control */

	// not saved:
	opium_double schar    // improved opium production
	post_txt     []string // text of posted sign
	storm_move   int      // next loc storm will move to
	garr_watch   []int    // units garrison watches for
	garr_host    []int    // units garrison will attack
	garr_tax     int      // garrison taxes collected
	garr_forward int      // garrison taxes forwarded
}

// nations structure
type entity_nation struct {
	name              string /* Name of the nation, e.g., Mandor Empire */
	citizen           string /* Name of a citizen, e.g., Mandorean */
	nobles            int    /* Total nobles */
	nps               int    /* Total NPS */
	gold              int    /* Total gold */
	players           int    /* Num players */
	win               int    /* Win? */
	proscribed_skills []int  /* Skills you can't have... */
	player_limit      int    /* Limit to # of players */
	capital           int    /* Capital city. */
	jump_start        int    /* Jump start points to start */
	neutral           bool   /* Can't capture/lose NPs. */
}

type entity_player struct {
	full_name       string
	email           string
	vis_email       string /* address to put in player list */
	password        string
	first_turn      int           /* which turn was their first? */
	last_order_turn int           /* last turn orders were submitted */
	orders          []*order_list /* ilist of orders for units in this faction */
	known           sparse        /* visited, lore seen, encountered */

	units    []int   /* what units are in our faction? */
	admits   admit_l /* admit permissions list */
	unformed []int   /* nobles as yet unformed */

	split_lines int /* split mail at this many lines */
	split_bytes int /* split mail at this many bytes */

	nation       int /* Player's nation */
	magic        int /* MUs or Priests? */
	noble_points int /* how many NP's the player has */
	jump_start   int /* Jump start points */

	format        int    /* turn report formatting control */
	rules_path    string /* external path for HTML */
	db_path       string /* external path for HTML */
	notab         bool   /* player can't tolerate tabs */
	first_tower   int    /* has player built first tower yet? */
	sent_orders   int    /* sent in orders this turn? */
	dont_remind   int    /* don't send a reminder */
	compuserve    bool   /* get Times from CIS */
	nationlist    int    /* Receive the nation mailing list? */
	broken_mailer int    /* quote begin lines */

	/* not saved: */
	times_paid      int    /* Times paid this month? */
	swear_this_turn int    /* have we used SWEAR this turn? */
	cmd_count       int    /* count of cmds started this turn */
	np_gained       int    /* np's added this turn -- not saved */
	np_spent        int    /* np's lost this turn -- not saved */
	deliver_lore    []int  /* show these to player -- not saved */
	weather_seen    sparse /* locs we've viewed the weather */
	output          sparse /* units with output -- not saved */
	locs            sparse /* locs we touched -- not saved */

}

// religion sub-structure
type entity_religion_skill struct {
	name        string /* Of the god.  */
	strength    int    /* Related strength skill */
	weakness    int    /* Related weakness skill */
	plant       int    /* The holy plant. */
	animal      int    /* The holy animal. */
	terrain     int    /* Holy terrain. */
	high_priest int    /* The high priest */
	bishops     [2]int /* The two bishops */
}

// ship descriptor
type entity_ship struct {
	hulls      int /* Various ship parts */
	forts      int
	sails      int
	ports      int
	keels      int
	galley_ram int /* galley is fitted with a ram */
}

// the specific religion skills require quite a bit of addition to the basic entity_skill structure.
// this stuff is captured by making a new skill subkind ("religion") and putting in a pointer to a "religion" structure.
type entity_skill struct {
	time_to_learn  int   /* days of study req'd to learn skill */
	time_to_use    int   /* days to use this skill */
	flags          int   /* Flags such as IS_POLLED, REQ_HOLY_SYMBOL, etc. */
	required_skill int   /* skill required to learn this skill */
	np_req         int   /* noble points required to learn */
	offered        []int /* skills learnable after this one */
	research       []int /* skills researable with this one */
	guild          []int /* skills offered if you're a guild member. */

	req            []*req_ent             /* ilist of items required for use or cast */
	produced       int                    /* simple production skill result */
	no_exp         int                    /* this skill not rated for experience */
	practice_cost  int                    /* cost to practice this skill. */
	practice_time  int                    /* time to practice this skill. */
	practice_prog  int                    /* A day longer to practice each N experience levels. */
	religion_skill *entity_religion_skill /* Possible religion pointer. */
	piety          int                    /* Piety required to cast religious skill. Or aura to cast magic skill; not used for all. */
	// not saved:
	use_count    int /* times skill used during turn */
	last_use_who int /* who last used the skill (this turn) */
}

type entity_subloc struct {
	teaches    []int /* skills location offers */
	opium_econ int   /* addiction level of city */
	defense    int   /* defense rating of structure */

	loot   int /* loot & pillage level */
	hp     int /* "hit points" */
	damage int /* 0=none, hp=fully destroyed */
	moat   int /* Has a moat? */
	//short shaft_depth;		/* depth of mine shaft */

	builds entity_build_l /* Ongoing builds here. */

	moving int          /* daystamp of beginning of movement */
	x_ship *entity_ship /* Maybe a ship? */

	//int capacity;			/* capacity of ship */
	//schar galley_ram;		/* galley is fitted with a ram */

	near_cities []int /* cities rumored to be nearby */
	safe        bool  /* safe haven */
	major       int   /* major city */
	prominence  int   /* prominence of city */

	//schar link_when;		/* month link is open, -1 = never */
	//schar link_open;		/* link is open now */

	link_to   []int /* where we are linked to */
	link_from []int /* where we are linked from */

	bound_storms []int /* storms bound to this ship */

	guild int /* what skill, if a sub_guild */

	//struct effect **effects;        /* ilist of effects on sub-location */

	tax_market  int /* Market tax rate. */
	tax_market2 int /* Temporary until end of month */

	/* Location control -- either here or loc */
	control       loc_control_ent
	control2      loc_control_ent
	entrance_size int /* Size of entrance to subloc */
}

type harvest struct {
	item      int
	skill     int
	worker    int
	chance    int // chance to get one each day, if nonzero
	got_em    string
	none_now  string
	none_ever string
	task_desc string
	public    int // 3rd party view, yes/no
	piety     int // does it use piety?
}

type item_ent struct {
	item int
	qty  int
}

type item_magic struct {
	creator       int
	lore          int /* deliver this lore for the item */
	religion      int /* Might be a religious artifact */
	curse_loyalty int /* curse noncreator loyalty */
	cloak_region  int
	cloak_creator int
	use_key       int   /* special use action */
	may_use       []int /* list of usable skills via this */
	may_study     []int /* list of skills studying from this */
	project_cast  int   /* stored projected cast */
	token_ni      int   /* ni for controlled npc units */
	quick_cast    int   /* stored quick cast */
	attack_bonus  int
	defense_bonus int
	missile_bonus int
	aura_bonus    int
	aura          int /* auraculum aura */
	token_num     int /* how many token controlled units */
	orb_use_count int /* how many uses left in the orb */

	// not saved:
	one_turn_use schar /* flag for one use per turn */
}

// location control: whether or not we're open, fees.
type loc_control_ent struct {
	closed bool
	nobles int // fee per noble
	men    int // fee per person
	weight int // fee per measure of weight
}

type loc_info struct {
	where     int
	here_list []int
}

type make_ struct {
	item        int
	inp1        int
	inp1_factor int // # of inp1 needed to make 1
	inp2        int
	inp2_factor int // # of inp2 needed to make 1
	inp3        int
	inp3_factor int // # of inp3 needed to make 1
	req_skill   int
	worker      int // worker needed
	got_em      string
	public      int // does everyone see us make this
	where       int // place required for production
	aura        int // aura per unit required
	factor      int // multiplying qty factor, usually 1
	days        int // days to make each thing
}

type mine_contents struct {
	items []*item_ent /* ilist of items held */
	// iron, gold, mithril, gate_crystals int // not used?
}

type olytime struct {
	day              int /* day of month */
	turn             int /* turn number */
	days_since_epoch int /* days since game begin */
}

// this structure holds game "options" for various different flavors of TAG.
type options_struct struct {
	turn_limit              int    /* Limit players to a certain # of turns. */
	auto_drop               bool   /* Drop non-responsive players. */
	free                    bool   /* Don't charge for this game. */
	turn_charge             string /* How much to charge per turn. */
	mp_antipathy            bool   /* Do mages & priests hate each other? */
	survive_np              bool   /* Does SFW return NPs when forgotten? */
	death_nps               int    /* What NPs get returned at death? */
	guild_teaching          bool   /* Do guilds teach guild skills? */
	accounting_dir          string /* Directory to "join" from. */
	accounting_prog         string /* Path of the accounting program. */
	html_path               string /* Path to html directories */
	html_passwords          string /* Path to html passwords */
	times_pay               int    /* What the Times pays for an article. */
	cpp                     string /* Path of cpp */
	full_markets            bool   /* City markets buy wood, etc. */
	output_tags             int    /* include <tag> in output */
	open_ended              bool   /* No end to game. */
	num_books               int    /* Number of teaching books in city */
	market_age              int    /* Months untouched in market before removal. */
	min_piety               int    /* Any priest can have this much piety. */
	piety_limit             int    /* Normal priest limited to piety_limit * num_followers */
	head_priest_piety_limit int    /* Head priest limited to head_priest_piety_limit * num_followers */
	top_piety               int    /* Monthly +piety for head priest */
	middle_piety            int    /* Monthly +piety for junior priests */
	bottom_piety            int    /* Monthly +piety for everyone else */
	claim_give              int    /* Allow putting gold in claim? */
	check_balance           int    /* No orders w/o positive balance. */
	free_np_limit           int    /* Play for free with this many NPs. */
}

type order_list struct {
	unit int      // unit orders are for
	l    orders_l // ilist of orders for unit
}
type orders_l [][]byte

type req_ent struct {
	item    int /* item required to use */
	qty     int /* quantity required */
	consume int /* REQ_xx */
}

type schar = byte

type skill_ent struct {
	skill        int
	days_studied int /* days studied * TOUGH_NUM */
	experience   int /* experience level with skill */
	know         int /* SKILL_xxx */
	// not saved:
	exp_this_month byte /* flag for add_skill_experience() */

}

type sparse = []int

type trade struct {
	kind       int /* BUY or SELL */
	item       int
	qty        int
	cost       int
	cloak      int /* don't reveal identity of trader */
	have_left  int
	month_prod int /* month city produces item */
	who        int /* redundant -- not saved */
	sort       int /* temp key for sorting -- not saved */
	old_qty    int /* qty at beginning of month, for trade goods */
	counter    int /* Counter to age and lose untraded goods. */
}

// traps - encodes some common traps...
type trap_struct struct {
	type_         int    /* Type of trap. */
	religion      int    /* Religion that ignores this trap. */
	num_attacks   int    /* Number of attacks */
	attack_chance int    /* Chance of attack killing someone */
	name          string /* Name of trap. */
	ignored       string /* Message if you can ignore this trap. */
	flying        string /* Message if you fly over the trap. */
	attack        string /* Message if it attacks you. */
}

type wait_arg struct {
	tag  int
	a1   int
	a2   int
	flag string
}

type uchar = byte

// how long a command has been running
func command_days(c *command) int { return c.days_executing }

// todo: is the "len - 1" right?
func numargs(c *command) int   { return len(c.parse) - 1 }
func oly_month(a *olytime) int { return (a.turn - 1) % NUM_MONTHS }
func oly_year(a *olytime) int  { return (a.turn - 1) / NUM_MONTHS }

// malloc-on-demand substructure references
func p_char(n int) *entity_char {
	if bx[n].x_char == nil {
		bx[n].x_char = &entity_char{}
	}
	return bx[n].x_char
}
func p_command(n int) *command {
	if bx[n].cmd == nil {
		bx[n].cmd = &command{}
	}
	return bx[n].cmd
}
func p_disp(n int) *att_ent {
	if bx[n].x_disp == nil {
		bx[n].x_disp = &att_ent{}
	}
	return bx[n].x_disp
}
func p_gate(n int) *entity_gate {
	if bx[n].x_gate == nil {
		bx[n].x_gate = &entity_gate{}
	}
	return bx[n].x_gate
}
func p_item(n int) *entity_item {
	if bx[n].x_item == nil {
		bx[n].x_item = &entity_item{}
	}
	return bx[n].x_item
}
func p_item_artifact(n int) *entity_artifact {
	if p_item(n).x_item_artifact == nil {
		p_item(n).x_item_artifact = &entity_artifact{}
	}
	return p_item(n).x_item_artifact
}
func p_item_magic(n int) *item_magic {
	if p_item(n).x_item_magic == nil {
		p_item(n).x_item_magic = &item_magic{}
	}
	return p_item(n).x_item_magic
}
func p_loc(n int) *entity_loc {
	if bx[n].x_loc == nil {
		bx[n].x_loc = &entity_loc{}
	}
	return bx[n].x_loc
}
func p_loc_info(n int) *loc_info {
	return &bx[n].x_loc_info
}
func p_magic(n int) *char_magic {
	if p_char(n).x_char_magic == nil {
		p_char(n).x_char_magic = &char_magic{}
	}
	return p_char(n).x_char_magic
}
func p_misc(n int) *entity_misc {
	if bx[n].x_misc == nil {
		bx[n].x_misc = &entity_misc{}
	}
	return bx[n].x_misc
}
func p_nation(n int) *entity_nation {
	if bx[n].x_nation == nil {
		bx[n].x_nation = &entity_nation{}
	}
	return bx[n].x_nation
}
func p_player(n int) *entity_player {
	if bx[n].x_player == nil {
		bx[n].x_player = &entity_player{}
	}
	return bx[n].x_player
}
func p_ship(n int) *entity_ship {
	if p_subloc(n).x_ship == nil {
		p_subloc(n).x_ship = &entity_ship{}
	}
	return p_subloc(n).x_ship
}
func p_skill(n int) *entity_skill {
	if bx[n].x_skill == nil {
		bx[n].x_skill = &entity_skill{}
	}
	return bx[n].x_skill
}
func p_subloc(n int) *entity_subloc {
	if bx[n].x_subloc == nil {
		bx[n].x_subloc = &entity_subloc{}
	}
	return bx[n].x_subloc
}

// "raw" pointers to substructures, may be NULL
func rp_char(n int) *entity_char { return bx[n].x_char }
func rp_command(n int) *command  { return bx[n].cmd }
func rp_disp(n int) *att_ent     { return bx[n].x_disp }
func rp_gate(n int) *entity_gate { return bx[n].x_gate }
func rp_item(n int) *entity_item { return bx[n].x_item }
func rp_item_artifact(n int) *entity_artifact {
	if rp_item(n) != nil {
		return rp_item(n).x_item_artifact
	} else {
		return nil
	}
}
func rp_item_magic(n int) *item_magic {
	if rp_item(n) != nil {
		return rp_item(n).x_item_magic
	} else {
		return nil
	}
}
func rp_loc(n int) *entity_loc    { return bx[n].x_loc }
func rp_loc_info(n int) *loc_info { return &bx[n].x_loc_info }
func rp_magic(n int) *char_magic {
	if rp_char(n) != nil {
		return rp_char(n).x_char_magic
	} else {
		return nil
	}
}
func rp_misc(n int) *entity_misc     { return bx[n].x_misc }
func rp_nation(n int) *entity_nation { return bx[n].x_nation }
func rp_player(n int) *entity_player { return bx[n].x_player }
func rp_relig_skill(n int) *entity_religion_skill {
	if bx[n].x_skill != nil {
		return bx[n].x_skill.religion_skill
	}
	return nil
}
func rp_ship(n int) *entity_ship {
	if rp_subloc(n) != nil {
		return rp_subloc(n).x_ship
	} else {
		return nil
	}
}
func rp_skill(n int) *entity_skill   { return bx[n].x_skill }
func rp_subloc(n int) *entity_subloc { return bx[n].x_subloc }

// func	item_creat_loc(n int) int {if rp_item_magic(n) != nil {return rp_item_magic(n).region_created} else {return 0}}
// func	loc_link_open(n int) int {if rp_subloc(n) != nil {return rp_subloc(n).link_open} else {return 0}}
// func	mine_depth(n int) int {if rp_subloc(n) != nil {return rp_subloc(n).shaft_depth / 3} else {return 0}}
// func char_pledge(n int) int {if rp_magic(n) != nil {return rp_magic(n).pledge} else {return 0}}
// func loc_barrier(n int) int {if rp_loc(n) != nil {return rp_loc(n).barrier} else {return 0}}
func banner(n int) string {
	if rp_misc(n) != nil {
		return rp_misc(n).display
	}
	return ""
}
func board_fee(n int) int {
	if rp_magic(n) != nil {
		panic("return rp_magic(n).fee")
	}
	return 0
}
func body_old_lord(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).old_lord
	}
	return 0
}
func border_open(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).border_open
	}
	return 0
}
func char_abil_shroud(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).ability_shroud
	}
	return 0
}
func char_attack(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).attack
	}
	return 0
}
func char_auraculum(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).auraculum
	}
	return 0
}
func char_behind(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).behind
	}
	return 0
}
func char_break(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).break_point
	}
	return 0
}
func char_cur_aura(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).cur_aura
	}
	return 0
}
func char_defense(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).defense
	}
	return 0
}
func char_guard(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).guard
	}
	return 0
}
func char_health(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).health
	}
	return 0
}
func char_hidden(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).hide_self
	}
	return 0
}
func char_hide_mage(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).hide_mage
	}
	return 0
}
func char_max_aura(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).max_aura
	}
	return 0
}
func char_melt_me(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).melt_me
	}
	return 0
}
func char_missile(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).missile
	}
	return 0
}
func char_new_lord(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).new_lord
	}
	return 0
}
func char_persuaded(n int) int {
	if rp_char(n) != nil {
		panic("return rp_char(n).persuaded")
	}
	return 0
}
func char_piety(n int) int {
	if is_priest(n) != 0 && rp_char(n) != nil {
		return rp_char(n).religion.piety
	}
	return 0
}
func char_proj_cast(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).project_cast
	}
	return 0
}
func char_quick_cast(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).quick_cast
	}
	return 0
}
func char_rank(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).rank
	}
	return 0
}
func char_sick(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).sick
	}
	return 0
}
func combat_skill(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).flags & COMBAT_SKILL
	}
	return 0
}
func entrance_size(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).entrance_size
	}
	return 0
}
func ferry_horn(n int) bool {
	if rp_magic(n) != nil {
		return rp_magic(n).ferry_flag
	}
	return false
}
func garrison_castle(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).garr_castle
	}
	return 0
}
func gate_dest(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).to_loc
	}
	return 0
}
func gate_dist(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).dist_from_gate
	}
	return 0
}
func gate_seal(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).seal_key
	} else {
		return 0
	}
}
func in_faery(n int) bool  { return region(n) == faery_region }
func in_hades(n int) bool  { return region(n) == hades_region }
func in_clouds(n int) bool { return region(n) == cloud_region }
func is_fighter(n int) int {
	if item_attack(n) != 0 || item_defense(n) != 0 || item_missile(n) != 0 || n == item_ghost_warrior {
		return TRUE
	} else {
		return FALSE
	}
}
func is_magician(n int) bool {
	if rp_magic(n) != nil {
		return rp_magic(n).magician != 0
	}
	return false
}
func item_animal(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).animal
	}
	return 0
}
func item_animal_part(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).animal_part
	}
	return 0
}
func item_attack(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).attack
	}
	return 0
}
func item_attack_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).attack_bonus
	}
	return 0
}
func item_aura(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).aura
	}
	return 0
}
func item_aura_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).aura_bonus
	}
	return 0
}
func item_capturable(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).capturable
	}
	return 0
}
func item_creat_cloak(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).cloak_creator
	}
	return 0
}
func item_creator(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).creator
	}
	return 0
}
func item_curse_non(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).curse_loyalty
	}
	return 0
}
func item_defense(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).defense
	}
	return 0
}
func item_defense_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).defense_bonus
	}
	return 0
}
func item_lore(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).lore
	} else {
		return 0
	}
}
func item_missile(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).missile
	} else {
		return 0
	}
}
func item_missile_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).missile_bonus
	} else {
		return 0
	}
}
func item_prog(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).wild
	} else {
		return 0
	}
}
func item_prominent(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).prominent
	} else {
		return 0
	}
}
func item_reg_cloak(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).cloak_region
	} else {
		return 0
	}
}
func item_split(n int) int {
	if rp_item(n) != nil {
		if rp_item(n).npc_split != 0 {
			return rp_item(n).npc_split
		} else {
			return 50
		}
	} else {
		return 50
	}
}
func item_token_ni(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).token_ni
	} else {
		return 0
	}
}
func item_token_num(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).token_num
	} else {
		return 0
	}
}
func item_unique(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).who_has
	} else {
		return 0
	}
}
func item_use_key(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).use_key
	}
	return 0
}
func item_wild(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).wild
	} else {
		return 0
	}
}
func learn_time(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).time_to_learn
	} else {
		return 0
	}
}
func loc_barrier(n int) int { return get_effect(n, ef_magic_barrier, 0, 0) }
func loc_damage(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).damage
	}
	return 0
}
func loc_defense(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).defense
	} else {
		return 0
	}
}
func loc_hp(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).hp
	}
	return 0
}
func loc_moat(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).moat
	}
	return 0
}
func loc_opium(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).opium_econ
	}
	return 0
}
func loc_pillage(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).loot
	}
	return 0
}
func loc_prominence(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).prominence
	}
	return 0
}
func loc_sea_lane(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).sea_lane
	}
	return 0
}
func loc_shroud(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).shroud
	}
	return 0
}
func loyal_kind(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).loy_kind
	}
	return 0
}
func loyal_rate(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).loy_rate
	}
	return 0
}
func major_city(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).major
	}
	return 0
}
func man_item(n int) int {
	// todo: this feels wrong
	if rp_item(n) != nil {
		return rp_item(n).is_man_item
	}
	return 0
}
func no_barrier(n int) int { return get_effect(n, ef_inhibit_barrier, 0, 0) }
func noble_item(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).unit_item
	} else {
		return 0
	}
}
func npc_last_dir(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).npc_dir
	}
	return 0
}
func npc_program(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).npc_prog
	} else {
		return 0
	}
}
func npc_summoner(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).summoned_by
	} else {
		return 0
	}
}
func only_defeatable(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).only_vulnerable
	} else {
		return 0
	}
}
func our_token(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).token
	} else {
		return 0
	}
}
func personal_break(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).personal_break_point
	} else {
		return 0
	}
}
func player_broken_mailer(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).broken_mailer
	}
	return 0
}
func player_compuserve(n int) bool {
	if rp_player(n) != nil {
		return rp_player(n).compuserve
	}
	return false
}
func player_email(n int) string {
	if rp_player(n) != nil {
		return rp_player(n).email
	} else {
		return ""
	}
}
func player_format(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).format
	}
	return 0
}
func player_js(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).jump_start
	} else {
		return 0
	}
}
func player_notab(n int) bool {
	if rp_player(n) != nil {
		return rp_player(n).notab
	}
	return false
}
func player_np(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).noble_points
	} else {
		return 0
	}
}
func player_split_bytes(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).split_bytes
	} else {
		return 0
	}
}
func player_split_lines(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).split_lines
	} else {
		return 0
	}
}
func practice_cost(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).practice_cost
	} else {
		return 0
	}
}
func practice_progressive(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).practice_prog
	} else {
		return 0
	}
}
func practice_time(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).practice_time
	} else {
		return 0
	}
}
func recent_pillage(n int) int {
	if rp_subloc(n) != nil {
		panic("return rp_subloc(n).recent_loot")
	}
	return 0
}
func reflect_blast(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).aura_reflect
	} else {
		return 0
	}
}
func release_swear(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).swear_on_release
	} else {
		return 0
	}
}
func req_skill(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).required_skill
	} else {
		return 0
	}
}
func restricted_control(n int) schar {
	if rp_misc(n) != nil {
		return rp_misc(n).cmd_allow
	} else {
		return 0
	}
}
func road_dest(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).to_loc
	} else {
		return 0
	}
}
func road_hidden(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).road_hidden
	} else {
		return 0
	}
}
func safe_haven(n int) bool {
	if rp_subloc(n) != nil {
		return rp_subloc(n).safe
	}
	return false
}
func sea_dist(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).dist_from_sea
	}
	return 0
}

// todo: what does this do and why?
func see_all(n int) int {
	return immed_see_all
}
func ship_has_ram(n int) int {
	if rp_ship(n) != nil {
		return rp_ship(n).galley_ram
	}
	return 0
}
func skill_exp(who, sk int) int {
	if rp_skill_ent(who, sk) != nil {
		return rp_skill_ent(who, sk).experience
	}
	return 0
}
func skill_flags(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).flags
	}
	return 0
}
func skill_no_exp(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).no_exp
	}
	return 0
}
func skill_np_req(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).np_req
	}
	return 0
}
func skill_produce(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).produced
	}
	return 0
}
func skill_time_to_use(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).time_to_use
	}
	return 0
}
func storm_bind(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).bind_storm
	}
	return 0
}
func storm_strength(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).storm_str
	}
	return 0
}
func swamp_dist(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).dist_from_swamp
	}
	return 0
}

func times_paid(n int) bool {
	if rp_player(n) != nil {
		return rp_player(n).times_paid != FALSE
	}
	return false
}

func weather_mage(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).knows_weather
	}
	return 0
}

func is_loc_or_ship(n int) bool { return (kind(n) == T_loc || kind(n) == T_ship) }
func is_ship(n int) bool {
	return (subkind(n) == sub_galley || subkind(n) == sub_roundship || subkind(n) == sub_raft || subkind(n) == sub_ship)
}
func is_ship_notdone(n int) bool {
	return (subkind(n) == sub_galley_notdone || subkind(n) == sub_roundship_notdone || subkind(n) == sub_raft_notdone || subkind(n) == sub_ship_notdone)
}
func is_ship_either(n int) bool { return (is_ship(n) || is_ship_notdone(n)) }
func kind(n int) schar {
	if n > 0 && n < MAX_BOXES && bx[n] != nil {
		return bx[n].kind
	} else {
		return T_deleted
	}
}
func kind_first(n int) int { return (box_head[(n)]) }
func kind_next(n int) int  { return (bx[(n)].x_next_kind) }

// return exactly where a unit is.
// May point to another character, a structure, or a region.
func loc(n int) int       { return rp_loc_info(n).where }
func sub_first(n int) int { return (sub_head[(n)]) }
func sub_next(n int) int  { return (bx[(n)].x_next_sub) }
func subkind(n int) schar {
	if bx[n] != nil {
		return bx[n].skind
	}
	return 0
}
func valid_box(n int) bool { return kind(n) != T_deleted }

// guild stuff
func is_guild(n int) int {
	if subkind(n) == sub_guild && rp_subloc(n).guild != 0 {
		return rp_subloc(n).guild
	}
	return 0
}
func guild_member(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).guild
	}
	return 0
}

// religion stuff
func god_name(n int) string { return rp_relig_skill(n).name }
func holy_plant(n int) int {
	if is_priest(n) != 0 {
		return rp_relig_skill(is_priest(n)).plant
	} else {
		return 0
	}
}
func is_follower(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).religion.priest
	} else {
		return 0
	}
}
func is_priest(n int) int { return skill_school(has_subskill((n), sub_religion)) }
func is_temple(n int) int {
	if subkind(n) == sub_temple && rp_subloc(n).guild != 0 {
		return rp_subloc(n).guild
	} else {
		return 0
	}
}
func is_wizard(n int) int { return has_subskill((n), sub_magic) }
func skill_aura(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).piety
	}
	return 0
}
func skill_piety(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).piety
	}
	return 0
}

// movement tests
// _moving indicates that the unit has initiated movement
// _gone indicates that the unit has actually left the locations, and should not be interacted with anymore.
//
// The distinction allows zero time commands to interact with the entity on the day movement is begun.
func char_moving(n int) int {
	if rp_char(n) == nil {
		return 0
	}
	return rp_char(n).moving
}
func char_gone(n int) int {
	if char_moving(n) == 0 {
		return 0
	}
	if evening {
		return sysclock.days_since_epoch - char_moving(n) + 1
	}
	return sysclock.days_since_epoch - char_moving(n)
}
func ship_moving(n int) int {
	if rp_subloc(n) == nil {
		return 0
	}
	return rp_subloc(n).moving
}
func ship_gone(n int) int {
	if ship_moving(n) == 0 {
		return 0
	}
	if evening {
		return sysclock.days_since_epoch - ship_moving(n) + 1
	}
	return sysclock.days_since_epoch - ship_moving(n)
}

func add_s(n int) string {
	if n == 0 {
		return ""
	}
	return "s"
}

func add_ds(n int) string { panic(`n, ((n) == 1 ? "" : "s"`) }
func alive(n int) bool    { return kind(n) == T_char }
func char_alone(n int) bool {
	// alone except for angels & ninjas
	return count_stack_any_real(n, true, true) == 1
}
func char_alone_stealth(n int) bool {
	// alone except for ninjas (permit angels)
	return count_stack_any_real(n, true, true) == 1
}
func char_really_alone(n int) bool {
	// alone counting everything
	return count_stack_any_real(n, false, false) == 1
}
func char_really_hidden(n int) bool {
	// hidden & alone (w/ angels & ninjas)
	return (char_hidden(n) != 0 && char_alone_stealth(n) && !is_prisoner(n))
}
func is_npc(n int) bool {
	// subkind(n) == sub_ni || loyal_kind(n) == LOY_npc)
	return subkind(n) != 0 || loyal_kind(n) == LOY_npc || loyal_kind(n) == LOY_summon
}
func is_prisoner(n int) bool {
	if rp_char(n) == nil {
		return false
	}
	return rp_char(n).prisoner != 0
}
func is_real_npc(n int) bool    { return player(n) < 1000 }
func magic_skill(n int) bool    { return subkind(skill_school(n)) == sub_magic }
func refugee(n int) bool        { return nation(n) == 0 && subkind(n) != sub_ni }
func religion_skill(n int) bool { return subkind(skill_school(n)) == sub_religion }
func wait(n int) int {
	if rp_command(n) != nil {
		return rp_command(n).wait
	} else {
		return 0
	}
}
func will_pay(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).pay
	} else {
		return 0
	}
}
