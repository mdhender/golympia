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
	"os"
	"path/filepath"
	"strings"
)

/*
 *  add.c  --  add new players to Olympia
 *
 *  oly -a will read data on new characters from stdin:
 *
 *	player number (provided by accounting system)
 *	faction name
 *	primary character name
 *	player's full name
 *	player's email address
 */

const (
	RANDOM_START = -1
)

var (
	new_players []int /* new players added this turn */
	new_chars   []int
)

// fetch_inp returns the next non-blank line from the input file
// after trimming spaces from it. returns an empty string on eof.
func fetch_inp(fp *os.File) string {
	if fp != nil {
		for s := getlin_ew(fp); s != nil; s = getlin_ew(fp) {
			if len(s) != 0 {
				return string(s)
			}
		}
	}
	return ""
}

/*
 *  Tue Apr  8 12:29:28 1997 -- Scott Turner
 *
 *  You can start in any city controlled by your nation; if not,
 *  you start in the default start city for your nation.
 *
 */
func pick_starting_city(nat, start_city int) int {
	if start_city != 0 && valid_box(start_city) && subkind(start_city) == sub_city && nation(player_controls_loc(start_city)) == nat {
		return start_city
	}

	/* RANDOM_START indicates we want to start in a random city. */
	if start_city == RANDOM_START {
		choice, sum := 0, 0
		for _, i := range loop_city() {
			if nation(player_controls_loc(i)) == nat {
				sum++
				if rnd(1, sum) == 1 {
					choice = i
				}
			}
		}

		if choice != 0 {
			return choice
		}
	}

	/*  return nations[nat].capital; */
	return rp_nation(nat).capital
}

/*
 *  Tue Apr  8 15:51:14 1997 -- Scott Turner
 *
 *  Total a unit's nps.
 *
 */
func nps_invested(who int) int {
	//int i, total = 1, categories = 0;
	//struct skill_ent *e;

	if kind(who) != T_char {
		return 0
	}

	categories, total := 0, 0
	for _, e := range loop_char_skill(who) {
		if e.skill != 0 && e.know == SKILL_know {
			total += skill_np_req(e.skill)
		}
		if skill_school(e.skill) == e.skill {
			categories++
		}
	}

	/*
	 *  Loyalty
	 *
	 */
	if p_char(who).loy_kind == LOY_oath {
		total += p_char(who).loy_rate
	}

	if categories > 3 {
		total += categories - 3
	}

	return total
}

/*
 *  Tue Apr  8 15:58:35 1997 -- Scott Turner
 *
 *  Do all the nations.
 *
 */
func calculate_nation_nps() {
	for _, i := range loop_nation() {
		rp_nation(i).nps = 0
		rp_nation(i).gold = 0
		rp_nation(i).players = 0
		rp_nation(i).nobles = 0
	}
	for _, pl := range loop_player() {
		if nation(pl) == 0 {
			continue
		}
		rp_nation(nation(pl)).players++
		rp_nation(nation(pl)).gold += has_item(pl, item_gold) /* CLAIM */
		for _, i := range loop_units(pl) {
			rp_nation(nation(pl)).nobles++
			rp_nation(nation(pl)).nps += nps_invested(i)
			rp_nation(nation(pl)).gold += has_item(i, item_gold)
		}
	}
}

/*
 *  Tue Apr  8 15:39:07 1997 -- Scott Turner
 *
 *  Calculate NP bonus for weak/strong nations.
 *
 */
func starting_noble_points(nation int) int {
	total_np, total_nations := 0.0, 0.0
	ratio := 0.0

	if rp_nation(nation).player_limit != 0 {
		return 12
	}

	if len(loop_nation()) == 1 {
		return 12
	}

	total_nations = 0
	for _, i := range loop_nation() {
		if i != nation && rp_nation(i).player_limit == 0 {
			total_np += float64(rp_nation(i).nps)
			total_nations++
		}
	}

	if total_np != 0 && total_nations != 0 && rp_nation(nation).nps != 0 {
		ratio = ((total_np / total_nations) / float64(rp_nation(nation).nps))
	} else if rp_nation(nation).nps == 0 {
		ratio = 1.0
	}

	if ratio >= 2.0 {
		return 20
	} else if ratio >= 1.75 {
		return 18
	} else if ratio >= 1.50 {
		return 16
	} else if ratio >= 1.25 {
		return 14
	} else if ratio >= 0.75 {
		return 12
	} else if ratio >= 0.50 {
		return 10
	} else if ratio >= 0.25 {
		return 8
	} else {
		return 6
	}
}

/*
 *  Thu Jul  3 13:55:17 1997 -- Scott Turner
 *
 *  Calculate gold bonus for weak/strong nations.
 *
 */
func starting_gold(nation int) int {
	ratio, total_gold, total_nations := 0.0, 0.0, 0.0

	if rp_nation(nation).player_limit != 0 {
		return 5000
	}

	if len(loop_nation()) == 1 {
		return 5000
	}

	total_nations = 0
	for _, i := range loop_nation() {
		if i != nation && rp_nation(i).player_limit == 0 {
			total_gold += float64(rp_nation(i).gold)
			total_nations++
		}
	}

	if total_gold != 0 && rp_nation(nation).gold != 0 {
		ratio = (total_gold / total_nations) / float64(rp_nation(nation).gold)
	} else if rp_nation(nation).nps == 0 {
		ratio = 1.0
	}

	if ratio >= 2.0 {
		ratio = 2.0
	} else if ratio < 0.25 {
		ratio = 0.25
	}
	return int(ratio * 5000)
}

func add_new_player(pl int, faction, character, full_name, email string, nation, start_city int) int {
	var who int
	//extern int new_ent_prime;        /* allocate short numbers */

	new_ent_prime = true
	who = new_ent(T_char, 0)
	new_ent_prime = false

	if who < 0 {
		return 0
	}

	set_name(pl, faction)
	set_name(who, character)

	pp := p_player(pl)
	cp := p_char(who)

	pp.FullName = full_name
	pp.EMail = email
	/*
	 *  Tue Apr  8 12:40:56 1997 -- Scott Turner
	 *
	 *  Noble points need to depend upon how the other nations
	 *  are doing.
	 *
	 */
	pp.NoblePoints = starting_noble_points(nation)
	pp.FirstTurn = sysclock.turn + 1
	pp.LastOrderTurn = sysclock.turn
	pp.Nation = nation
	/*
	 *  Thu Apr  9 08:41:42 1998 -- Scott Turner
	 *
	 *  Jump start points come from the nation... plus some
	 *  per turn?
	 *
	 */
	pp.JumpStart = rp_nation(nation).jump_start + (sysclock.turn / 5)
	if pp.JumpStart > 56 {
		pp.JumpStart = 56
	}

	if strings.HasSuffix(email, "@compuserve.com") {
		pp.CompuServe = true
	}

	cp.health = 100
	cp.break_point = 50
	cp.attack = 80
	cp.defense = 80

	set_where(who, pick_starting_city(nation, start_city))
	set_lord(who, pl, LOY_oath, 2)

	gen_item(who, item_peasant, 25)
	gen_item(who, item_gold, 200)
	gen_item(pl, item_gold, starting_gold(nation)) /* CLAIM item */
	gen_item(pl, item_lumber, 50)                  /* CLAIM item */

	new_players = append(new_players, pl)
	new_chars = append(new_chars, who)

	add_unformed_sup(pl)

	return pl
}

func failed_join(email, reason string) {
	if !save_flag {
		return
	}

	//tmpfile := fmt.Sprintf("failedjoin%d", get_process_id());
	//tmp, err := os.Create(tmpfile);
	//if err != nil {
	//	panic(err)
	//}
	//fprintf(tmp, "From: %s\n", from_host);
	//fprintf(tmp, "Subject: Failed attempt to join Olympia: The Age of Gods, Game %d\n", game_number);
	//fprintf(tmp, "To: %s, %s\n\n", email, rp_player(gm_player).email);
	//fprintf(tmp, "%s", reason);
	//fprintf(tmp, "\n\n -- The Game Master\n\n");
	//fclose(tmp);
	///*VLN cmd = sout("sendmail -t -odq < %s", tmpfile);*/
	//cmd := sout("msmtp -t < %s", tmpfile);
	//if (system(cmd)) {
	//    log.Printf( "Failed to send 'failed join' mail?\n");
	//} else {
	//    unlink(tmpfile);
	//}

	panic("!implemented")
}

func make_new_players_sup(acct string, fp *os.File) bool {
	n, sc := 0, 0
	var faction, character, full_name, email, nat string

	if faction = fetch_inp(fp); faction == "" {
		log.Printf("%s: Unable to add %s <%s>n", acct, full_name, email) // todo: uninitialized values
		log.Printf("    partial read of faction.")
		return false
	}

	if character = fetch_inp(fp); character == "" {
		log.Printf("%s: Unable to add %s <%s>n", acct, full_name, email) // todo: uninitialized values
		log.Printf("    partial read of character.")
		return false
	}

	if full_name = fetch_inp(fp); full_name == "" {
		log.Printf("%s: Unable to add %s <%s>n", acct, full_name, email) // todo: uninitialized values
		log.Printf("    partial read of full_name.")
		return false
	}

	if email = fetch_inp(fp); email == "" {
		log.Printf("%s: Unable to add %s <%s>n", acct, full_name, email)
		log.Printf("    partial read of email.")
		return false
	}

	if nat = fetch_inp(fp); nat == "" {
		log.Printf("%s: Unable to add %s <%s>n", acct, full_name, email)
		log.Printf("    partial read of nation.")
		return false
	}

	/*
	 *  Start_City can be null (blank)
	 *
	 */
	start_city := fetch_inp(fp)

	pl := scode([]byte(acct))
	assert(pl > 0 && pl < MAX_BOXES)

	/*
	 * Maybe he's already in the game.
	 *
	 */
	if bx[pl] != nil {
		fail_buf := fmt.Sprintf("Olympia was unable to add you to this game because\nthere is already a faction assigned to your account.\n")
		failed_join(email, fail_buf)
		return true
	}

	/*
	 *  Figure out the nation
	 *
	 */
	n = find_nation(nat)
	if n != 0 {
		wout(gm_player, "Couldn't add player %s: bad nation.", acct)
		fail_buf := fmt.Sprintf("Olympia was unable to add you to this game because\nwe could not decipher the nation name (%s) that\nyou provided.  Please try to join again using a valid nation.\n", nat)
		failed_join(email, fail_buf)
		return true
	}
	/*
	 *  Wed Jul  2 11:32:01 1997 -- Scott Turner
	 *
	 *  Check player limits on the nation they've chosen.
	 *
	 */
	if rp_nation(n).player_limit != 0 {
		total := 0
		for _, i := range loop_player() {
			if nation(i) == n {
				total++
			}
		}
		if total >= rp_nation(n).player_limit {
			fail_buf := fmt.Sprintf("Olympia was unable to add you to this game because\nthe %s nation has already reached its limit of %d players.\n  Please try again using a different nation.\n", rp_nation(n).name, rp_nation(n).player_limit)
			failed_join(email, fail_buf)
			return true
		}
	}
	/*
	 *  Check to make sure the city someone is trying to
	 *  start in is controlled by the appropriate nation.
	 *  If not, give him the default starting city.
	 *
	 */
	if strings.HasPrefix(strings.ToLower(start_city), "rand") {
		sc = RANDOM_START
	} else {
		sc = code_to_int([]byte(start_city))
	}

	alloc_box(pl, T_player, sub_pl_regular)

	add_new_player(pl, faction, character, full_name, email, n, sc)
	if start_city != "" {
		my_free(start_city)
	}
	log.Printf("\tadded player %s\n", box_name(pl))

	return true
}

func make_new_players() {
	files, err := os.ReadDir(libdir)
	if err != nil {
		log.Printf("make_new_players: can't open %s: \n", options.accounting_dir)
		return
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		acct := f.Name()

		log.Printf("VLN account = %s\n", acct)
		fnam := filepath.Join(options.accounting_dir, acct, fmt.Sprintf("Join-tag-%d", game_number))

		fp, err := os.Open(fnam)
		if err != nil {
			continue
		}

		if !make_new_players_sup(acct, fp) {
			// this should generate some notice of a failed add and alert the GM/user.
			log.Printf("Failed to add new player %q\n", acct)
		}
		_ = fp.Close()
	}
}

func rename_act_join_files() error {
	for i := 0; i < len(new_players); i++ {
		pl := new_players[i]
		acct := fmt.Sprintf("%s", box_code_less(pl))

		old_name := filepath.Join(options.accounting_dir, acct, fmt.Sprintf("Join-tag-%d", game_number))
		new_name := filepath.Join(options.accounting_dir, acct, fmt.Sprintf("Join-tag-%d-", game_number))

		if err := rename(old_name, new_name); err != nil {
			return fmt.Errorf("rename(%s, %s): %w", old_name, new_name)
		}
	}
	return nil
}

func new_player_banners() {
	out_path = MASTER
	out_alt_who = OUT_BANNER

	for i := 0; i < len(new_players); i++ {
		pl := new_players[i]

		//#if 0
		//			p := p_player(pl);
		//        out(pl, "From: %s", from_host);
		//        out(pl, "Reply-To: %s", reply_host);
		//        if (p.email)
		//            out(pl, "To: %s (%s)", p.email,
		//                p.full_name ? p.full_name : "???");
		//        out(pl, "Subject: Welcome to Olympia");
		//        out(pl, "");
		//#endif

		wout(pl, "Welcome to Olympia!")
		wout(pl, "")
		wout(pl, "This is an initial position report for your new faction.")

		wout(pl, "You are player %s, \"%s\".", box_code_less(pl), just_name(pl))
		wout(pl, "")

		wout(pl, "The next turn will be turn %d.", sysclock.turn+1)

		month := (sysclock.turn) % NUM_MONTHS
		year := (sysclock.turn + 1) / NUM_MONTHS

		wout(pl, "It is season \"%s\", month %d, in the year %d.", month_names[month], month+1, year+1)
		out(pl, "")

		// report_account_sup(pl)
	}

	out_path = 0
	out_alt_who = 0
}

func show_new_char_locs() {
	out_path = MASTER
	show_loc_no_header = true

	for i := 0; i < len(new_chars); i++ {
		who := new_chars[i]
		where := subloc(who)

		out_alt_who = where
		show_loc(player(who), where)

		where = loc(where)
		if loc_depth(where) == LOC_province {
			out_alt_who = where
			show_loc(player(who), where)
		}
		mark_loc_stack_known(who, where)
	}

	show_loc_no_header = false
	out_path = 0
	out_alt_who = 0
}

func new_player_report() {
	var i int

	out_path = MASTER
	out_alt_who = OUT_BANNER

	for i = 0; i < len(new_players); i++ {
		player_report_sup(new_players[i])
	}

	out_path = 0
	out_alt_who = 0

	for i = 0; i < len(new_players); i++ {
		show_unclaimed(new_players[i], new_players[i])
	}
}

func new_char_report() {
	var i int

	indent += 3

	for i = 0; i < len(new_chars); i++ {
		char_rep_sup(new_chars[i], new_chars[i])
	}

	indent -= 3
}

func mail_initial_reports() {
	var i int
	var s, t string
	var pl int

	for i = 0; i < len(new_players); i++ {
		pl = new_players[i]

		s = filepath.Join(libdir, "log", libdir, fmt.Sprintf("%d", pl))
		t = filepath.Join(libdir, "save", fmt.Sprintf("%d", sysclock.turn), fmt.Sprintf("%d", pl))

		if err := rename(s, t); err != nil {
			log.Printf("couldn't rename %s to %s: %v\n", s, t, err)
		}

		send_rep(pl, sysclock.turn)
	}
}

func new_order_templates() {
	var pl, i int

	out_path = MASTER
	out_alt_who = OUT_TEMPLATE

	for i = 0; i < len(new_players); i++ {
		pl = new_players[i]
		orders_template(pl, pl)
	}

	out_path = 0
	out_alt_who = 0
}

func new_player_list_sup(who int, pl int) {
	var p *EntityPlayer
	var s string

	p = p_player(pl)

	if p.EMail != "" {
		if p.FullName != "" {
			s = sout("%s <%s>", p.FullName, p.EMail)
		} else {
			s = sout("<%s>", p.EMail)
		}
	} else if p.FullName != "" {
		s = p.FullName
	} else {
		s = ""
	}

	out(who, "%4s   %s  (%s)", box_code_less(pl), just_name(pl), rp_nation(nation(pl)).name)
	if s != "" {
		out(who, "       %s", s)
	}
	out(who, "")
}

func new_player_list() {
	var pl int
	var i int

	stage("new_player_list()")

	out_path = MASTER
	out_alt_who = OUT_NEW

	vector_players()

	//#if 0
	//    for i =  0; i < len(new_players); i++ {
	//        pl = new_players[i];
	//        ilist_rem_value(&out_vector, pl);
	//    }
	//#endif

	for i = 0; i < len(new_players); i++ {
		pl = new_players[i]
		new_player_list_sup(VECT, pl)
	}

	out_path = 0
	out_alt_who = 0
}

func new_player_top(mail int) {

	stage("new_player_top()")

	open_logfile()
	/* Need to do this before making a new player! */
	calculate_nation_nps()
	make_new_players()
	show_new_char_locs()
	new_char_report()
	new_player_banners()
	new_player_report()
	new_order_templates()
	gen_include_section() /* must be last */
	close_logfile()

	if mail != FALSE {
		mail_initial_reports()
	}
}

func add_new_players() {

	stage("add_new_players()")

	calculate_nation_nps()
	make_new_players()
	show_new_char_locs()
	new_char_report()
	new_player_banners()
	new_player_report()
	new_order_templates()
	new_player_list() /* show new players to the old players */
}
