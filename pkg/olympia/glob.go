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

var (
	from_host  = "moderator@olytag.com (Olympia Moderator)" /* UPDATE */
	reply_host = "tagtest@olytag.com (Olympia Orders)"      /* UPDATE */

	bx       []*box       /* all possible entities */
	box_head [T_MAX]int   /* heads of x_next_kind chain */
	sub_head [SUB_MAX]int /* heads of x_next_sub chain */

	//nations[MAX_NATIONS]entity_nation;  /* The nations */
	//num_nations = 0;

	libdir            = "lib"
	sysclock          olytime
	game_number       = 0   /* we hope to have many :-) */
	indep_player      = 100 /* independent player */
	gm_player         = 200 /* The Fates */
	deserted_player   = 201 /* Deserted nobles */
	skill_player      = 202 /* skill listing */
	eat_pl            = 203 /* Order scanner */
	npc_pl            = 206 /* Subloc monster player */
	garr_pl           = 207 /* Garrison unit owner */
	combat_pl         = 0   /* Combat log */
	show_day          = false
	post_has_been_run = FALSE
	seed_has_been_run = FALSE
	dist_sea_compute  = FALSE
	near_city_init    = FALSE
	cookie_init       = FALSE
	garrison_magic    = 999
	/* Map size, set in "system" */
	xsize   = 100
	ysize   = 100
	options options_struct
)

/*
 *  Allow field:
 *
 *	c	character
 *	p	player entity
 *	i	immediate mode only (debugging/maintenance)
 *	r	restricted -- for npc units under control
 *	g	garrison
 *	m	Gamemaster only
 */

var cmd_tbl = []cmd_tbl_ent{
	/*
	   allow  name         start         finish      intr      time poll pri
	   num_args_req max_args arg_types[] cmd_checker
	*/

	{"", "", nil, nil, nil, 0, 0, 3, 0, 0, [5]int{}, nil, nil},
	//{"m",  "acquire",   v_acquire,    nil,       nil,       1,  0,  3,         0, 0, [5]int{}, nil},
	{"cpr", "accept", v_accept, nil, nil, 0, 0, 0, 1, 3, [5]int{0, 0, CMD_qty, 0, 0}, accept_comment, nil},
	{"cpr", "admit", v_admit, nil, nil, 0, 0, 0, 1, 0, [5]int{}, admit_comment, admit_check},
	{"cr", "attack", v_move_attack, d_move_attack, nil, -1, 0, 3, 1, 0, [5]int{}, attack_comment, nil},
	{"cr", "banner", v_banner, nil, nil, 0, 0, 1, 1, 2, [5]int{}, nil, nil},
	{"cr", "behind", v_behind, nil, nil, 0, 0, 1, 1, 1, [5]int{}, nil, nil},
	//{"c",  "bind",      v_bind_storm, d_bind_storm, nil,     7,  0,  3,         0, 0, [5]int{}, nil},
	//{"c",  "board",     v_board,      nil,       nil,       0,  0,  2,         0, 0, [5]int{}, nil},
	{"c", "border", v_border, nil, nil, 0, 0, 0, 2, 2, [5]int{}, nil, nil},
	{"c", "breed", v_breed, d_breed, nil, 7, 0, 3, 2, 2, [5]int{CMD_item, CMD_item, 0, 0, 0}, nil, nil},
	{"c", "bribe", v_bribe, d_bribe, nil, 7, 0, 3, 2, 3, [5]int{CMD_unit, CMD_gold, 0, 0, 0}, nil, nil},
	{"c", "build", v_build, d_build, nil, -1, 1, 3, 1, 4, [5]int{}, nil, nil},
	{"c", "buy", v_buy, nil, nil, 0, 0, 1, 1, 4, [5]int{CMD_item, CMD_qty, 0, 0, 0}, buy_comment, nil},
	{"c", "catch", v_catch, nil, nil, -1, 1, 3, 0, 2, [5]int{CMD_qty, CMD_days, 0, 0, 0}, catch_comment, nil},
	{"c", "claim", v_claim, nil, nil, 0, 0, 1, 1, 2, [5]int{CMD_item, CMD_qty, 0, 0, 0}, nil, nil},
	{"c", "collect", v_collect, d_collect, i_collect, -1, 1, 3, 1, 3, [5]int{CMD_item, CMD_qty, CMD_days, 0, 0}, collect_comment, nil},
	{"cr", "contact", v_contact, nil, nil, 0, 0, 0,
		1, 1, [5]int{CMD_unit, 0, 0, 0, 0}, nil, nil},
	{"m", "credit", v_credit, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"c", "decree", v_decree, nil, nil, 0, 0, 0,
		2, 2, [5]int{}, nil, nil},
	{"c", "dedicate", v_dedicate, d_dedicate, nil, 7, 0, 3,
		1, 1, [5]int{CMD_unit, 0, 0, 0, 0}, default_comment, nil},
	{"cpr", "default", v_att_clear, nil, nil, 0, 0, 0,
		0, 0, [5]int{CMD_unit, 0, 0, 0, 0}, attitude_comment, nil},
	{"cpr", "defend", v_defend, nil, nil, 0, 0, 0,
		0, 0, [5]int{CMD_unit, 0, 0, 0, 0}, attitude_comment, nil},
	{"c", "die", v_die, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"c", "discard", v_discard, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"cr", "drop", v_discard, nil, nil, 0, 0, 1,
		1, 3, [5]int{CMD_item, CMD_qty, 0, 0, 0}, drop_comment, nil},
	{"m", "emote", v_emote, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"cr", "execute", v_execute, nil, nil, 0, 0, 1,
		0, 1, [5]int{CMD_unit, 0, 0, 0, 0}, default_comment, nil},
	//{"c",  "exhume",    v_exhume,     d_exhume,   nil,       7,  0,  , 3,         0, 0, [5]int{}, nil},
	{"c", "explore", v_explore, d_explore, nil, 7, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "fee", v_fee, nil, nil, 0, 0, 1,
		2, 3, [5]int{}, nil, nil},
	{"c", "ferry", v_ferry, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"c", "fish", v_fish, nil, nil, -1, 1, 3,
		0, 2, [5]int{0, CMD_days, 0, 0, 0}, nil, nil},
	{"cr", "flag", v_flag, nil, nil, 0, 0, 1,
		1, 1, [5]int{}, nil, nil},
	{"c", "fly", v_fly, d_fly, nil, -1, 1, 2,
		1, 0, [5]int{}, move_comment, nil},
	{"c", "forget", v_forget, nil, nil, 0, 0, 1,
		1, 1, [5]int{CMD_skill, 0, 0, 0, 0}, default_comment, nil},
	{"c", "form", v_form, d_form, nil, 7, 0, 3,
		2, 2, [5]int{}, nil, nil},
	{"cp", "format", v_format, nil, nil, 0, 0, 1,
		1, 2, [5]int{}, nil, nil},
	{"c", "garison", v_garrison, nil, nil, 1, 0, 3,
		1, 1, [5]int{}, nil, nil},
	{"c", "garrison", v_garrison, nil, nil, 1, 0, 3,
		1, 1, [5]int{}, nil, nil},
	{"cr", "get", v_get, nil, nil, 0, 0, 1,
		2, 4, [5]int{CMD_unit, CMD_item, CMD_qty, 0, 0}, get_comment, nil},
	{"cr", "give", v_give, nil, nil, 0, 0, 1,
		2, 4, [5]int{CMD_unit, CMD_item, CMD_qty, 0, 0}, give_comment, nil},
	{"cr", "go", v_move_attack, d_move_attack, nil, -1, 0, 2,
		0, 0, [5]int{}, move_comment, nil},
	{"c", "guard", v_guard, nil, nil, 0, 0, 1,
		1, 1, [5]int{}, default_comment, nil},
	{"c", "guild", v_dedicate_tower, d_dedicate_tower, nil, 7, 0, 3,
		1, 1, [5]int{CMD_skill, 0, 0, 0, 0}, default_comment, nil},
	{"c", "hide", v_hide, d_hide, nil, 3, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "honor", v_honor, nil, nil, 0, 0, 3,
		1, 1, [5]int{CMD_gold, 0, 0, 0, 0}, nil, nil},
	{"c", "honour", v_honor, nil, nil, 1, 0, 3,
		1, 1, [5]int{CMD_gold, 0, 0, 0, 0}, nil, nil},
	{"cpr", "hostile", v_hostile, nil, nil, 0, 0, 0,
		0, 0, [5]int{CMD_unit, 0, 0, 0, 0}, attitude_comment, nil},
	{"c", "identify", v_identify, nil, nil, 0, 0, 3,
		1, 1, [5]int{CMD_item, 0, 0, 0, 0}, default_comment, nil},
	{"c", "incite", v_incite, nil, nil, 7, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "make", v_make, d_make, i_make, -1, 1, 3,
		1, 2, [5]int{CMD_item, 0, 0, 0, 0}, nil, nil},
	{"c", "mallorn", v_mallorn, nil, nil, -1, 1, 3,
		0, 0, [5]int{}, nil, nil},
	{"cp", "message", v_message, nil, nil, 1, 0, 3,
		2, 2, [5]int{0, CMD_unit, 0, 0, 0}, nil, nil},
	{"cr", "move", v_move_attack, d_move_attack, nil, -1, 0, 2,
		1, 0, [5]int{}, move_comment, nil},
	{"cpr", "maxpay", v_maxpay, nil, nil, 0, 0, 1,
		0, 1, [5]int{CMD_gold, 0, 0, 0, 0}, nil, nil},
	{"cpr", "name", v_name, nil, nil, 0, 0, 1,
		1, 2, [5]int{}, nil, nil},
	{"cp", "nationlist", v_nationlist, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"cpr", "neutral", v_neutral, nil, nil, 0, 0, 0,
		0, 0, [5]int{CMD_unit, 0, 0, 0, 0}, attitude_comment, nil},
	{"cp", "notab", v_notab, nil, nil, 0, 0, 1,
		1, 1, [5]int{}, nil, nil},
	{"c", "oath", v_oath, nil, nil, 1, 0, 3,
		1, 1, [5]int{}, nil, nil},
	{"c", "opium", v_opium, nil, nil, -1, 1, 3,
		0, 0, [5]int{}, nil, nil},
	{"cr", "pay", v_pay, nil, nil, 0, 0, 1,
		1, 3, [5]int{CMD_unit, CMD_gold, 0, 0, 0}, nil, nil},
	{"cr", "pillage", v_pillage, d_pillage, nil, 7, 1, 3,
		0, 1, [5]int{}, nil, nil},
	{"c", "post", v_post, nil, nil, 1, 0, 3,
		1, 1, [5]int{}, nil, nil},
	{"c", "practice", v_practice, d_practice, nil, 7, 0, 3,
		1, 1, [5]int{CMD_practice, 0, 0, 0, 0}, study_comment, nil},
	{"cp", "press", v_press, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"cr", "promote", v_promote, nil, nil, 0, 0, 1,
		1, 1, [5]int{CMD_unit, 0, 0, 0, 0}, default_comment, nil},
	{"c", "quarry", v_quarry, nil, nil, -1, 1, 3,
		0, 2, [5]int{0, CMD_days, 0, 0, 0}, nil, nil},
	{"p", "quit", v_quit, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, quit_check},
	{"c", "raise", v_raise, d_raise, nil, 7, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "rally", v_rally, d_rally, nil, 7, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"cr", "raze", v_raze, d_raze, nil, -1, 1, 3,
		0, 2, [5]int{0, CMD_days, 0, 0, 0}, default_comment, nil},
	{"cpr", "realname", v_fullname, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"c", "reclaim", v_reclaim, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"c", "recruit", v_recruit, nil, nil, -1, 1, 3,
		0, 1, [5]int{CMD_days, 0, 0, 0, 0}, nil, nil},
	{"c", "repair", v_repair, d_repair, i_repair, -1, 1, 3,
		0, 1, [5]int{CMD_days, 0, 0, 0, 0}, nil, nil},
	{"cp", "rumor", v_rumor, nil, nil, 0, 0, 1,
		0, 1, [5]int{}, nil, nil},
	{"cp", "rumors", v_rumor, nil, nil, 0, 0, 1,
		0, 1, [5]int{}, nil, nil},
	{"c", "sail", v_sail, d_sail, i_sail, -1, 0, 4,
		1, 0, [5]int{}, move_comment, nil},
	{"c", "sell", v_sell, nil, nil, 0, 0, 1,
		3, 4, [5]int{CMD_item, CMD_qty, 0, 0, 0}, buy_comment, nil},
	{"cr", "seek", v_seek, d_seek, nil, 7, 1, 3,
		0, 1, [5]int{CMD_unit, 0, 0, 0, 0}, default_comment, nil},
	{"c", "sneak", v_sneak, d_sneak, nil, 3, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"cp", "split", v_split, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"cr", "stack", v_stack, nil, nil, 0, 0, 1,
		1, 1, [5]int{CMD_unit, 0, 0, 0, 0}, default_comment, nil},
	{"c", "stone", v_quarry, nil, nil, -1, 1, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "study", v_study, d_study, nil, 7, 1, 3,
		1, 2, [5]int{CMD_skill, 0, 0, 0, 0}, study_comment, nil},
	/* {"c",  "surrender", v_surrender,  nil,       nil,       1,  0,  1}, ,
	   0, 0, [5]int{}, nil, nil*/
	{"c", "swear", v_swear, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"cr", "take", v_get, nil, nil, 0, 0, 1,
		2, 4, [5]int{CMD_unit, CMD_item, 0, 0, 0}, get_comment, nil},
	{"cp", "tax", v_tax, nil, nil, 0, 0, 1,
		3, 3, [5]int{}, nil, nil},
	{"c", "teach", v_teach, nil, nil, 7, 1, 2,
		1, 2, [5]int{CMD_skill, CMD_days, 0, 0, 0}, default_comment, nil},
	{"cp", "tell", v_tell, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"c", "think", v_think, nil, nil, 1, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"cp", "times", v_times, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"c", "train", v_make, d_make, i_make, -1, 1, 3,
		1, 2, [5]int{CMD_item, CMD_qty, 0, 0, 0}, nil, nil},
	{"c", "trance", v_trance, d_trance, nil, 28, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"cr", "terrorize", v_terrorize, d_terrorize, nil, 7, 0, 3,
		2, 2, [5]int{CMD_unit, 0, 0, 0, 0}, nil, nil},
	{"c", "torture", v_torture, d_torture, nil, 7, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "unload", v_unload, nil, nil, 0, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "ungarrison", v_ungarrison, nil, nil, 1, 0, 3,
		0, 1, [5]int{}, nil, nil},
	{"cr", "unstack", v_unstack, nil, nil, 0, 0, 1,
		0, 1, [5]int{}, nil, nil},
	{"c", "use", v_use, d_use, i_use, -1, 1, 3,
		1, 0, [5]int{CMD_use, 0, 0, 0, 0}, study_comment, nil},
	{"crm", "wait", v_wait, d_wait, i_wait, -1, 1, 1,
		1, 0, [5]int{}, nil, nil},
	{"c", "wood", v_wood, nil, nil, -1, 1, 3,
		0, 2, [5]int{0, CMD_days, 0, 0, 0}, nil, nil},
	{"cr", "xyzzy", v_xyzzy, nil, nil, 0, 0, 3,
		0, 0, [5]int{}, nil, nil},
	{"c", "yew", v_yew, nil, nil, -1, 1, 3,
		0, 2, [5]int{0, CMD_days, 0, 0, 0}, nil, nil},

	{"cr", "north", v_north, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "n", v_north, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "s", v_south, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "south", v_south, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "east", v_east, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "e", v_east, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "west", v_west, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "w", v_west, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "enter", v_enter, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "exit", v_exit, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "in", v_enter, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},
	{"cr", "out", v_exit, nil, nil, -1, 0, 2,
		0, 0, [5]int{}, nil, nil},

	{"", "begin", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "unit", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "email", nil, nil, nil, 0, 0, 0,
		1, 1, [5]int{}, nil, nil},
	{"", "vis_email", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "end", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "flush", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "lore", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "passwd", nil, nil, nil, 0, 0, 0,
		0, 1, [5]int{}, nil, nil},
	{"", "password", nil, nil, nil, 0, 0, 0,
		0, 1, [5]int{}, nil, nil},
	{"", "players", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "resend", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"", "option", nil, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},
	{"cpr", "stop", v_stop, nil, nil, 0, 0, 0,
		0, 0, [5]int{}, nil, nil},

	{"i", "look", v_look, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "l", v_look, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "ct", v_ct, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "be", v_be, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "additem", v_add_item, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "subitem", v_sub_item, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "artifact", v_make_artifact, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "h", v_listcmds, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "dump", v_dump, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "i", v_invent, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "fix", v_fix, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "fix2", v_fix2, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "kill", v_kill, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "los", v_los, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"m", "relore", v_relore, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "sk", v_skills, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "know", v_know, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "seed", v_seed, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "seedorc", v_seedorc, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "seedmarket", v_seedmarket, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "sheet", v_lore, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "poof", v_poof, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "postproc", v_postproc, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "save", v_save, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "seeall", v_see_all, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "tp", v_take_pris, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"i", "makeloc", v_makeloc, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil},
	{"", "", nil, nil, nil, 0, 0, 1,
		0, 0, [5]int{}, nil, nil}}

var kind_s = []string{
	"deleted",  /* T_deleted */
	"player",   /* T_player */
	"char",     /* T_char */
	"loc",      /* T_loc */
	"item",     /* T_item */
	"skill",    /* T_skill */
	"gate",     /* T_gate */
	"road",     /* T_road */
	"deadchar", /* T_deadchar */
	"ship",     /* T_ship */
	"post",     /* T_post */
	"storm",    /* T_storm */
	"unform",   /* T_unform */
	"lore",     /* T_lore */
	"nation",   /* T_nation */
	""}

var subkind_s = []string{
	"<no subkind>",
	"ocean",                      /* sub_ocean */
	"forest",                     /* sub_forest */
	"plain",                      /* sub_plain */
	"mountain",                   /* sub_mountain */
	"desert",                     /* sub_desert */
	"swamp",                      /* sub_swamp */
	"underground",                /* sub_under */
	"faery hill",                 /* sub_faery_hill */
	"island",                     /* sub_island */
	"ring of stones",             /* sub_stone_cir */
	"mallorn grove",              /* sub_mallorn_grove */
	"bog",                        /* sub_bog */
	"cave",                       /* sub_cave */
	"city",                       /* sub_city */
	"lair",                       /* sub_lair */
	"graveyard",                  /* sub_graveyard */
	"ruins",                      /* sub_ruins */
	"field",                      /* sub_battlefield */
	"enchanted forest",           /* sub_ench_forest */
	"rocky hill",                 /* sub_rocky_hill */
	"circle of trees",            /* sub_tree_cir */
	"pits",                       /* sub_pits */
	"pasture",                    /* sub_pasture */
	"oasis",                      /* sub_oasis */
	"yew grove",                  /* sub_yew_grove */
	"sand pit",                   /* sub_sand_pit */
	"sacred grove",               /* sub_sacred_grove */
	"poppy field",                /* sub_poppy_field */
	"temple",                     /* sub_temple */
	"galley",                     /* sub_galley */
	"roundship",                  /* sub_roundship */
	"castle",                     /* sub_castle */
	"galley-in-progress",         /* sub_galley_notdone */
	"roundship-in-progress",      /* sub_roundship_notdone */
	"ghost ship",                 /* sub_ghost_ship */
	"temple-in-progress",         /* sub_temple_notdone */
	"inn",                        /* sub_inn */
	"inn-in-progress",            /* sub_inn_notdone */
	"castle-in-progress",         /* sub_castle_notdone */
	"mine",                       /* sub_mine */
	"mine-in-progress",           /* sub_mine_notdone */
	"scroll",                     /* sub_scroll */
	"magic",                      /* sub_magic */
	"palantir",                   /* sub_palantir */
	"auraculum",                  /* sub_auraculum */
	"tower",                      /* sub_tower */
	"tower-in-progress",          /* sub_tower_notdone */
	"pl_system",                  /* sub_pl_system */
	"pl_regular",                 /* sub_pl_regular */
	"region",                     /* sub_region */
	"pl_savage",                  /* sub_pl_savage */
	"pl_npc",                     /* sub_pl_npc */
	"collapsed mine",             /* sub_mine_collapsed */
	"ni",                         /* sub_ni */
	"demon lord",                 /* sub_demon_lord */
	"dead body",                  /* sub_dead_body */
	"fog",                        /* sub_fog */
	"wind",                       /* sub_wind */
	"rain",                       /* sub_rain */
	"pit",                        /* sub_hades_pit */
	"artifact",                   /* sub_artifact */
	"pl_silent",                  /* sub_pl_silent */
	"npc_token",                  /* sub_npc_token */
	"garrison",                   /* sub_garrison */
	"cloud",                      /* sub_cloud */
	"raft",                       /* sub_raft */
	"raft-in-progress",           /* sub_raft_notdone */
	"suffuse_ring",               /* sub_suffuse_ring */
	"religion",                   /* sub_religion */
	"holy symbol",                /* sub_holy_symbol */
	"mist",                       /* sub_mist */
	"book",                       /* sub_book */
	"guild",                      /* sub_market */
	"trade_good",                 /* sub_trade_good */
	"city-in-progress",           /* sub_city_notdone */
	"ship",                       /* sub_ship */
	"ship-in-progress",           /* sub_ship_notdone */
	"mine-shaft",                 /* sub_mine_shaft */
	"mine-shaft-in-progress",     /* sub_mine_shaft_notdone */
	"orc-stronghold",             /* sub_orc_stronghold */
	"orc-stronghold-in-progress", /* sub_orc_stronghold_notdone */
	"Staff-of-the-Sun",           /* sub_special_staff */
	"lost_soul",                  /* sub_lost_soul */
	"undead",                     /* sub_undead */
	"pen-crown",                  /* sub_pen_crown */
	"animal-part",                /* sub_animal_part */
	"magical-artifact",           /* sub_magic_artifact */
	""}

var short_dir_s = []string{
	"<no dir>",
	"n",
	"e",
	"s",
	"w",
	"u",
	"d",
	"i",
	"o",
	""}

var full_dir_s = []string{
	"<no dir>",
	"north",
	"east",
	"south",
	"west",
	"up",
	"down",
	"in",
	"out",
	""}

var exit_opposite = []int{
	0,
	DIR_S,
	DIR_W,
	DIR_N,
	DIR_E,
	DIR_DOWN,
	DIR_UP,
	DIR_OUT,
	DIR_IN,
	0}

var loc_depth_s = []string{
	"<no depth>",
	"region",
	"province",
	"subloc",
	""}

var month_names = []string{
	"Fierce winds",
	"Snowmelt",
	"Blossom bloom",
	"Sunsear",
	"Thunder and rain",
	"Harvest",
	"Waning days",
	"Dark night",
	""}

func glob_init() {
	//for i := 0; i < MAX_BOXES; i++ {
	//	bx[i] = nil
	//}

	for i := 0; i < T_MAX; i++ {
		box_head[i] = 0
	}

	bx = make([]*box, MAX_BOXES)
}
