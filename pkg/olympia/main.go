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

var (
	acct_flag     = FALSE
	flush_always  = false
	immed_after   = FALSE
	immed_see_all = FALSE /* override hidden-ness, for debugging */
	immediate     = TRUE  // can't be bool because it holds player and item id's, too
	/*
	 *  pretty_data_files:  include parenthesisted names in the data files,
	 *  to make them easier to read.
	 */
	pretty_data_files = FALSE
	save_flag         = false /* set in main? */
	time_self         = false /* print timing info */
)

func getopt(argc int, argv []string, flags string) int {
	panic("!implemented")
}

func RunOly(args ...string) int {
	add_flag := false
	art_flag := false
	combat_test_flag := false
	eat_flag := false
	errflag := false
	inhibit_add_flag := false
	lore_flag := false
	mail_now := false
	map_flag := false
	//map_test_flag := false
	run_flag := false
	unspool_first_flag := false
	var c int

	call_init_routines()

	argc := len(args)
	argv := args
	EOF := -1
	for {
		if c = getopt(argc, argv, "axefirmLl:pR?StMTAqXE"); c == EOF {
			break
		}
		switch byte(c) {
		case 'm':
			map_flag = true

		case 'a':
			add_flag = true
			immediate = FALSE

		case 'x':
			inhibit_add_flag = true

		case 'X':
			combat_test_flag = true

		case 'A':
			acct_flag = TRUE

		case 'f':
			flush_always = true

		case 'e':
			eat_flag = true
			immediate = FALSE

		case 'E':
			unspool_first_flag = true

		case 'i':
			immed_after = TRUE

		case 'l': /* set libdir */
			//libdir = string(optarg)
			panic("!implemented")

		case 'L':
			lore_flag = true

		case 'p':
			if pretty_data_files == FALSE {
				pretty_data_files = TRUE
			} else {
				pretty_data_files = FALSE
			}

		case 'q': /* test artifacts */
			art_flag = true

		case 'r': /* run a turn */
			immediate = FALSE
			run_flag = true

		case 'R': /* test random number generator */
			log.Printf("error: '%s -R' has been replaced with 'random/random_test'\n", argv[0])
			os.Exit(2)

		case 'S': /* save database when done */
			save_flag = true

		case 't':
			log.Printf("error: '%s -t' has been replaced with 'vectors/ilist_test'\n", argv[0])
			os.Exit(2)

		case 'T':
			time_self = true

		case 'M':
			mail_now = true

		default:
			errflag = true
		}
	}

	if errflag {
		log.Printf("usage: oly [options]\n")
		log.Printf("  -a        Add new players mode\n")
		log.Printf("  -x        Inhibit adding players during turn.\n")
		log.Printf("  -e        Eat orders from libdir/spool\n")
		log.Printf("  -f        Don't buffer files for debugging\n")
		log.Printf("  -i        Immediate mode\n")
		log.Printf("  -l dir    Specify libdir, default ./lib\n")
		log.Printf("  -L        Generate lore dictionary.\n")
		log.Printf("  -p        Don't make data files pretty\n")
		log.Printf("  -r        Run a turn\n")
		log.Printf("  -R        Test the random number generator\n")
		log.Printf("  -S        Save the database at completion\n")
		log.Printf("  -t        Test ilist code\n")
		log.Printf("  -T        Print timing info\n")
		log.Printf("  -M        Mail reports\n")
		log.Printf("  -A        Charge player accounts\n")
		log.Printf("  -X        Combat test\n")
		os.Exit(1)
	}

	/*
	 *  Sat Apr 15 11:55:36 2000 -- Scott Turner
	 *
	 *  Lock up; prevents multiple TAGs running simultaneously.
	 *
	 */
	lock_tag()

	load_db()

	/*
	 *  Create a couple of stacks and have them battle
	 *  it out a few hundred times and report the results.
	 *
	 */
	if combat_test_flag {
		var result int
		open_logfile()
		sum_a := 0
		for i := 0; i < 100; i++ {
			a := new_ent(T_char, 0)
			b := new_ent(T_char, 0)
			c := new_ent(T_char, 0)
			d := new_ent(T_char, 0)
			p := p_char(a)
			p.behind = 9
			p.health = 100
			p.attack = 80
			p.defense = 80
			p.break_point = 0
			s := p_skill_ent(a, sk_control_battle)
			s.know = SKILL_know
			s.experience = 500
			s = p_skill_ent(a, sk_use_beasts)
			s.know = SKILL_know
			s.experience = 500
			set_where(a, 10201)
			set_lord(a, gm_player, LOY_oath, 1)
			set_name(a, "A")
			p = p_char(b)
			p.behind = 0
			p.health = 100
			p.attack = 80
			p.defense = 80
			p.break_point = 0
			s = p_skill_ent(b, sk_control_battle)
			s.know = SKILL_know
			s.experience = 500
			s = p_skill_ent(b, sk_use_beasts)
			s.know = SKILL_know
			s.experience = 500
			set_where(b, 10201)
			set_lord(b, gm_player, LOY_oath, 1)
			set_name(b, "B")
			p = p_char(c)
			p.behind = 0
			p.health = 100
			p.attack = 360
			p.defense = 360
			p.break_point = 0
			s = p_skill_ent(c, sk_control_battle)
			s.know = SKILL_know
			s.experience = 500
			s = p_skill_ent(c, sk_use_beasts)
			s.know = SKILL_know
			s.experience = 500
			set_where(c, 10201)
			set_lord(c, gm_player, LOY_oath, 1)
			set_name(c, "C")
			p = p_char(d)
			p.behind = 9
			p.health = 100
			p.attack = 80
			p.defense = 80
			p.break_point = 0
			s = p_skill_ent(d, sk_control_battle)
			s.know = SKILL_know
			s.experience = 500
			s = p_skill_ent(d, sk_use_beasts)
			s.know = SKILL_know
			s.experience = 500
			set_where(d, 10201)
			set_lord(d, gm_player, LOY_oath, 1)
			set_name(d, "D")
			/*
			 *  A & B stacked together and C & D stacked
			 *  together.  A & D are behind.
			 *
			 */
			gen_item(a, item_soldier, 10)
			gen_item(a, item_angel, 15)
			gen_item(c, item_nazgul, 3)
			if result = regular_combat(c, a, 0, 0); result != 0 {
				sum_a++
			}
		}
		close_logfile()
		log.Printf("Sum = %d.\n", sum_a)
		os.Exit(1)
	}

	if map_flag {
		if asciiMap := load_cmap(); asciiMap != nil {
			for _, row := range asciiMap {
				fmt.Printf("%s\n", row)
			}
		}
		log.Println("")
		if asciiMap := load_cmap_players(); asciiMap != nil {
			for _, row := range asciiMap {
				fmt.Printf("%s\n", row)
			}
		}
		os.Exit(1)
	}

	if art_flag {
		var i int
		var c command
		c.who = gm_player
		for i = 0; i < 10; i++ {
			c.a = create_random_artifact(gm_player)
			v_identify(&c)
			save_box(os.Stderr, c.a)
		}
		/*
		 *  Look to see if anyone owns a crown, palantir, etc.
		 *
		 */
		for _, i = range loop_subkind(sub_npc_token) {
			if item_unique(i) != FALSE && !is_real_npc(player(item_unique(i))) {
				log.Printf("Crown %s held by %s.\n", box_name(i), box_name(player(item_unique(i))))
			}
		}
		for _, i = range loop_subkind(sub_palantir) {
			if item_unique(i) != FALSE && !is_real_npc(player(item_unique(i))) {
				log.Printf("Palantir %s held by %s.\n", box_name(i), box_name(player(item_unique(i))))
			}
		}
		for _, i = range loop_subkind(sub_scroll) {
			if item_unique(i) != FALSE && !is_real_npc(player(item_unique(i))) {
				log.Printf("Scroll %s held by %s.\n", box_name(i), box_name(player(item_unique(i))))
			}
		}
		for _, i = range loop_subkind(sub_suffuse_ring) {
			if item_unique(i) != FALSE && !is_real_npc(player(item_unique(i))) {
				log.Printf("Suffuse ring %s held by %s.\n", box_name(i), box_name(player(item_unique(i))))
			}
		}
		for _, i = range loop_subkind(sub_artifact) {
			if item_unique(i) != FALSE && !is_real_npc(player(item_unique(i))) {
				log.Printf("Artifact %s held by %s.\n", box_name(i), box_name(player(item_unique(i))))
			}
		}
		os.Exit(0)
	}

	if lore_flag {
		open_logfile()
		gm_show_all_skills(skill_player, true)
		os.Exit(0)
	}

	if eat_flag {
		if mail_now {
			eat_loop(true)
		} else {
			eat_loop(false)
		}
		os.Exit(0)
	}

	if run_flag {
		// if unspool_first_flag is on, then before running the turn, eat up any waiting mail.
		if unspool_first_flag {
			if err := mkdir(filepath.Join(libdir, "orders")); err != nil {
				panic(err)
			}
			if err := mkdir(filepath.Join(libdir, "spool")); err != nil {
				panic(err)
			}
			//chmod(filepath.Join(libdir, "spool"), 0777)
			if mail_now {
				read_spool(true)
			} else {
				read_spool(false)
			}
		}

		open_logfile()
		open_times()

		show_day = true
		pre_month()
		process_orders()
		post_month()
		//#if 0
		//        /* This was only to translate to the new system. */
		//        artifact_fixer();
		//#endif
		show_day = false

		determine_output_order()
		turn_end_loc_reports()
		list_order_templates()
		player_ent_info()
		character_report()

		player_banner()
		// if (acct_flag && !options.free)
		//    charge_account();
		// report_account();
		summary_report()
		player_report()

		scan_char_skill_lore()
		show_lore_sheets()
		//#if 0
		//        list_new_players();
		//#endif
		if !options.open_ended {
			check_win_conditions()
		}
		gm_report(gm_player)
		gm_show_all_skills(skill_player, true)
		if !inhibit_add_flag {
			add_new_players()
		}
		gen_include_section() /* must be last */
		close_logfile()

		write_player_list()
		write_nations_lists()
		write_email()
		write_totimes()
		write_forwards()
		write_factions()
	}

	if add_flag {
		if mail_now {
			new_player_top(TRUE)
		} else {
			new_player_top(FALSE)
		}
		mail_now = false
	}

	if immediate != FALSE || immed_after != FALSE {
		immediate = TRUE

		open_logfile()
		immediate_commands()
		close_logfile()
	}

	check_db() /* check database integrity */

	if save_flag {
		save_db()
	}

	if save_flag && run_flag {
		save_logdir()
	}

	do_times()

	if mail_now {
		mail_reports()
	}

	stage("")

	return 0
}

func call_init_routines() error {
	init_lower()
	dir_assert()
	glob_init()         /* initialize global tables */
	initialize_buffer() /* used by sout() */
	init_spaces()
	init_random() /* seed random number generator */

	return nil
}

func write_totimes() {
	var fp *os.File
	var fnam string
	var pl int

	fnam = filepath.Join(libdir, "totimes")

	fp, err := fopen(fnam, "w")
	if err != nil {
		log.Printf("can't write %s: %v\n", fnam, err)
		return
	}

	for _, pl = range loop_player() {
		if rp_player(pl) != nil && rp_player(pl).email != "" && !player_compuserve(pl) {
			fprintf(fp, "%s\n", rp_player(pl).email)
		}
	}

	_ = fp.Close()
}

func write_email() {
	var fnam string
	var pl int

	fnam = filepath.Join(libdir, "email")

	fp, err := fopen(fnam, "w")
	if err != nil {
		log.Printf("can't write %s: %v", fnam, err)
		return
	}

	for _, pl = range loop_player() {
		if rp_player(pl) != nil && rp_player(pl).email != "" {
			fprintf(fp, "%s\n", rp_player(pl).email)
		}
	}

	_ = fp.Close()
}

func fix_email(email string) string {
	return strings.Replace(email, "@", "(at)", 1)
}

func list_a_player(fp *os.File, pl int, flag *int) {
	var p *entity_player
	var s string
	//var n int
	var c byte
	var email string

	p = p_player(pl)
	if p.email != "" || p.vis_email != "" {
		if p.vis_email != "" {
			email = p.vis_email
		} else {
			email = p.email
		}

		if p.full_name != "" {
			s = sout("%s &lt;%s&gt;", p.full_name, fix_email(email))
		} else {
			s = sout("&lt;%s&gt", fix_email(email))
		}
	} else if p.full_name != "" {
		s = p.full_name
	} else {
		s = ""
	}

	if ilist_lookup(new_players, pl) >= 0 {
		c = '*'
		*flag = TRUE
	} else {
		c = ' '
	}

	fprintf(fp, "<TR><TD>%4s %c</TD>", box_code_less(pl), c)
	fprintf(fp, "<TD>%s</TD>", just_name(pl))
	fprintf(fp, "<TD>%s</TD>", rp_nation(nation(pl)).name)

	if len(s) != 0 {
		fprintf(fp, "<TD>%s</TD>", s)
	}
	fprintf(fp, "</TR>\n")
}

/*
 *  Mon Nov  9 18:09:30 1998 -- Scott Turner
 *
 *  Drop a person (or resubscribe him) to the nation's mailing list.
 *
 */
func v_nationlist(c *command) int {
	var p *entity_player

	p = p_player(player(c.who))

	if p.nationlist != FALSE {
		p.nationlist = FALSE
		wout(c.who, "Will receive the nation mailing list.")
	} else {
		p.nationlist = TRUE
		wout(c.who, "Will not receive the nation mailing list.")
	}

	return TRUE

}

/*
 *  Fri Nov  6 16:01:01 1998 -- Scott Turner
 *
 *  Write out the files that serve as nation mailing lists.
 *
 *  Tue Nov 10 05:55:45 1998 -- Scott Turner
 *
 *  Add DM to all the lists.
 *
 */
func write_nations_lists() {
	stage("write_nations_lists()")

	for _, i := range loop_nation() {
		fnam := filepath.Join(libdir, rp_nation(i).citizen)
		fp, err := fopen(fnam, "w")

		if fp == nil {
			log.Printf("can't write %s: %v\n", fnam, err)
			continue
		}

		/*
		 *  Add regular players who belong to this nation, have
		 *  an email address, and haven't "opted-out".
		 *
		 */
		for _, pl := range loop_player() {
			if rp_player(pl) != nil && subkind(pl) == sub_pl_regular && nation(pl) == i && rp_player(pl).email != "" && rp_player(pl).nationlist == FALSE {
				fprintf(fp, "%s\n", rp_player(pl).email)
			}
		}

		/*
		 *  Add the DM...
		 *
		 */
		if rp_player(gm_player) != nil && rp_player(gm_player).nationlist == FALSE {
			fprintf(fp, "%s\n", rp_player(gm_player).email)
		}
		/*
		 *  Close the file.
		 *
		 */
		_ = fp.Close()
	}
}

func write_player_list() {
	var fnam string
	var pl int
	flag := false

	stage("write_player_list()")

	fnam = filepath.Join(libdir, "players.html")

	fp, err := fopen(fnam, "w")
	if err != nil {
		log.Printf("can't write %s: %v", fnam, err)
		return
	}

	fprintf(fp, "<HTML>\n")
	fprintf(fp, "<HEAD>\n")
	fprintf(fp, "<TITLE>Olympia Game %d Player List</TITLE>\n", game_number)
	fprintf(fp, "</HEAD>\n")
	fprintf(fp, "<BODY>\n")
	fprintf(fp, "<TABLE ALIGN=ABSCENTER CELLSPACING=0 CELLPADDING=5 WIDTH=\"100%%\" BGCOLOR=\"#48D1CC\" >\n")
	fprintf(fp, "<TR><TD><B><FONT SIZE=+1>Num</FONT></B></TD><TD><B><FONT SIZE=+1>Faction</FONT></B></TD><TD><B><FONT SIZE=+1>Nation</FONT></B></TD><TD><B><FONT SIZE=+1>Email Address</FONT></B></TD></TR>\n")

	for _, pl = range loop_player() {
		if rp_player(pl) != nil && rp_player(pl).email != "" && subkind(pl) == sub_pl_regular {
			var yint int
			list_a_player(fp, pl, &yint)
			flag = yint != FALSE
		}
	}

	fprintf(fp, "</TABLE>\n")

	if flag {
		fprintf(fp, "* -- New player this turn")
	}

	fprintf(fp, "</BODY></HTML>\n")

	_ = fp.Close()
}

func write_forward_sup(who_for int, target int, fp *os.File) {
	var pl int
	var s string

	pl = player(who_for)
	s = player_email(pl)

	if len(s) != 0 {
		fprintf(fp, "%s|%s\n", box_code_less(target), s)
	}
}

func write_forwards() {
	var fnam string

	fnam = filepath.Join(libdir, "forward")

	fp, err := fopen(fnam, "w")
	if err != nil {
		log.Printf("can't write %s: %v", fnam, err)
		return
	}

	for _, pl := range loop_player() {
		write_forward_sup(pl, pl, fp)
	}

	for _, npc := range loop_char() {
		write_forward_sup(npc, npc, fp)
	}

	for _, garr := range loop_garrison() {
		if player(province_admin(garr)) != 0 {
			write_forward_sup(player(province_admin(garr)), garr, fp)
		}
	}

	_ = fp.Close()
}

func write_faction_sup(who_for int, target int, fp *os.File) {
	var pl int
	var s string

	pl = player(who_for)
	s = player_email(pl)

	if len(s) != 0 {
		fprintf(fp, "%s|%s\n", box_code_less(target), s)
	}
}

func write_factions() {
	var fnam string

	fnam = filepath.Join(libdir, "factions")

	fp, err := fopen(fnam, "w")
	if err != nil {
		log.Printf("can't write %s: %v", fnam, err)
		return
	}

	for _, pl := range loop_player() {
		write_faction_sup(pl, pl, fp)
	}

	_ = fp.Close()
}

/*
 *  Make a report from a raw file into the specified file.
 *
 */
func make_report(format int, fnam string, report string, pl int) int {
	var form, entab, tags string
	switch format {
	case HTML:
		form = "-html"
	case RAW:
		form = "-raw"
	case TAGS:
		form = "-text"
		tags = "-tags"
	default:
		form = "-text"
	}

	if player_notab(pl) {
		entab = "-notab"
	}

	tmp := fmt.Sprintf("rep %s %s %s %s >> %s", form, tags, entab, fnam, report)
	return system(tmp)
}

/*
 *  Wed Nov 10 05:50:57 1999 -- Scott Turner
 *
 *  A little function to return a format string.
 *
 */
func format_string(i int) string {
	switch i {
	case HTML:
		return "HTML"
	case TEXT:
		return "TEXT"
	case TAGS:
		return "TAGS"
	case RAW:
		return "RAW"
	}
	return "UNKNOWN"
}

/*
 *  Sat Oct 30 13:13:45 1999 -- Scott Turner
 *
 *  Added a loop over all the formats.
 *
 */
func send_rep(pl, turn int) int {
	var p *entity_player
	var ret int
	var zfnam string
	var fnam string
	split_lines := player_split_lines(pl)
	split_bytes := player_split_bytes(pl)
	formats := p_player(pl).format
	var i int
	var email string

	p = rp_player(pl)

	if p == nil || p.email == "" {
		return FALSE
	}

	/*
	 *  Default for text format if nothing is set.
	 *
	 */
	if formats == 0 {
		formats = TEXT
	}

	/*
	 *  Prepare the report file and the input file before the
	 *  format loop.
	 *
	 */
	report := fmt.Sprintf("/tmp/sendrep%d.%s", get_process_id(), box_code_less(pl))

	fnam = filepath.Join(libdir, "save", fmt.Sprintf("%d", turn), fmt.Sprintf("%d", pl))

	if access(fnam, R_OK) < 0 {
		zfnam = filepath.Join(libdir, "save", fmt.Sprintf("%d", turn), fmt.Sprintf("%d.gz", pl))

		if access(zfnam, R_OK) < 0 {
			unlink(report)
			return FALSE
		}

		fnam = sout("/tmp/zrep.%d", pl)

		ret = system(sout("gzcat %s > %s", zfnam, fnam))
		if ret == 0 {
			log.Printf("couldn't unpack %s\n", zfnam)
			unlink(fnam)
			unlink(report)
			return FALSE
		}
	}

	/*
	 *  Here's the format loop.
	 *
	 */
	for i = 1; i < ALT; i = i << 1 {
		if (formats & i) != 0 {

			fp, err := fopen(report, "w")
			if err != nil {
				log.Printf("send_rep: can't write %s: %v", report, err)
				return FALSE
			}

			fprintf(fp, "From: %s\n", from_host)
			if reply_host != "" {
				fprintf(fp, "Reply-To: %s\n", reply_host)
			}
			fprintf(fp, "To: %s (%s)\n", p.email, or_string(p.full_name != "", p.full_name, "???"))
			fprintf(fp, "Subject: Olympia:TAG game %d turn %d report [%s]\n", game_number, turn, format_string(i))
			fprintf(fp, "\n")
			_ = fp.Close()

			ret = make_report(i, fnam, report, pl)
			if ret != 0 {
				log.Printf("send_rep: failed: make_report(%d, %q, %q, %d)\n", i, fnam, report, pl)
				unlink(report)
				if zfnam != "" {
					unlink(fnam)
				}
				return FALSE
			}

			var cmd string
			if split_lines == 0 && split_bytes == 0 {
				/* VLN cmd = sout("sendmail -t -odq < %s", report); */
				cmd = sout("msmtp -t < %s", report)
			} else {
				/* VLN cmd = sout("mailsplit -s %d -l %d -c 'sendmail -t -odq' < %s", split_bytes, split_lines, report); */
				cmd = sout("mailsplit -s %d -l %d -c 'msmtp -t' < %s", split_bytes, split_lines, report)
			}

			log.Printf("   %s\n", cmd)
			ret = system(cmd)
			if ret != 0 {
				log.Printf("send_rep: mail to %s failed: %s\n", p.email, cmd)
			}
			unlink(report)
		}
	}

	if zfnam != "" {
		unlink(fnam)
	}

	/*
	 *  Sun Dec 31 17:45:28 2000 -- Scott Turner
	 *
	 *  Send the Times also.
	 *
	 */
	/* VLN cmd = sout("sendmail -odq %s < %s/Times", p.email, libdir); */

	/* VLN remove "," from email adddress string */
	email = strings.ReplaceAll(p.email, ",", " ")

	/*VLN cmd = sout("msmtp %s < %s/Times", p.email, libdir); */
	cmd := sout("msmtp %s < %s/Times", email, libdir)
	ret = system(cmd)
	if ret != 0 {
		log.Printf("send_rep: mail to %s failed: %q\n", p.email, cmd)
	}

	return TRUE
}

func mail_reports() {
	var pl int

	stage("mail_reports()")

	/*
	 *  Refactor -- move all the formats down into send_rep
	 *
	 */
	for _, pl = range loop_player() {
		send_rep(pl, sysclock.turn)
	}

	setup_html_all()
}

func preprocess(in string, out string, args string) int {
	var buf string
	err := 0

	buf = fmt.Sprintf("%s %s < %s > %s", options.cpp, args, in, out)
	err = system(buf)
	if err != 0 {
		log.Printf("Error preprocessing %s.\n", in)
	}
	return err
}

/*
 *  Rich's on-line web stuff
 *
 *  Fri Oct 15 10:22:40 1999 -- Scott Turner
 *
 *  Modify our umask and permissions so that any file created is
 *  readable by all.  This is necessary for the reports-on-line stuff.
 *
 *  Thu Nov 11 06:56:58 1999 -- Scott Turner
 *
 *  Add stuff to write out the %d.pre file.
 *
 */
func setup_html_all() {
	var pl int
	var fnam, fnam2 string

	stage("setup_html()")

	fnam = filepath.Join(options.html_path, fmt.Sprintf("%d", game_number), fmt.Sprintf("%d.pre", game_number))
	fp, err := fopen(fnam, "w")
	if err != nil {
		log.Printf("Can't open %s for writing?", fnam)
		fp = os.Stderr
	}
	fprintf(fp, "<TABLE CELLSPACING=0 CELLPADDING=5 WIDTH=\"100%%\" BGCOLOR=\"#48D1CC\" >\n")

	/* write and execute only for self. */
	umask(S_IWGRP | S_IXGRP | S_IWOTH | S_IXOTH)

	for _, pl = range loop_player() {
		if !is_real_npc(pl) {
			setup_html_dir(pl)
			set_html_pass(pl)
			/* VLN output_html_rep(pl);  because these are broken */
			/*VLN output_html_map(pl); */
			fprintf(fp, "<TR><TD><A HREF=\"http://UPDATE.com/oly/tag/reports/%d/%s/tag%d-%d.html.gz\">[%s] %s</A></TD></TR>\n",
				game_number,
				box_code_less(pl),
				game_number,
				sysclock.turn,
				box_code_less(pl),
				just_name(pl))
		}
	}

	/* Now switch back. */
	umask(S_IRWXO)

	fprintf(fp, "</TABLE>")

	if fp != os.Stderr {
		_ = fp.Close()
	}

	// now we need to call the C pre-processor on the reports.pre file to incorporate this latest set of changes.
	fnam = filepath.Join(options.html_path, fmt.Sprintf("%d", game_number), "reports.pre")
	fnam2 = filepath.Join(options.html_path, fmt.Sprintf("%d", game_number), "reports.html")
	preprocess(fnam, fnam2, "-P")

	/* copy_public_turns(); */
}

func setup_html_dir(pl int) {
	var fnam string
	var fnam2 string
	var fp *os.File

	/* read and execute for all. */
	umask(S_IWGRP | S_IWOTH)
	fnam = filepath.Join(options.html_path, fmt.Sprintf("%d", game_number), box_code_less(pl))
	if err := mkdir(fnam); err != nil {
		panic(err)
	}
	umask(S_IWGRP | S_IXGRP | S_IWOTH | S_IXOTH)

	fnam2 = filepath.Join(fnam, ".htaccess")
	fp, err := fopen(fnam2, "w")
	if err != nil {
		log.Printf("can't write %s: %v", fnam2, err)
		return
	}

	fprintf(fp, "AuthUserFile %s.%d\n", options.html_passwords, game_number)
	fprintf(fp, "AuthGroupFile /dev/null\n")
	fprintf(fp, "AuthName ByPassword\n")
	fprintf(fp, "AuthType Basic\n")
	fprintf(fp, "\n")
	fprintf(fp, "<Limit GET>\n")
	fprintf(fp, "require user %s admin\n", box_code_less(pl))
	fprintf(fp, "</Limit>\n")

	_ = fp.Close()
}

func set_html_pass(pl int) {
	var buf string
	var p *entity_player

	p = rp_player(pl)
	if p == nil {
		return
	}

	pw := p.password
	if len(pw) == 0 {
		pw = DEFAULT_PASSWORD
	}
	/*
	   Usage:
	   	htpasswd [-cmdpsD] passwordfile username
	   	htpasswd -b[cmdpsD] passwordfile username password

	   	htpasswd -n[mdps] username
	   	htpasswd -nb[mdps] username password
	*/
	buf = fmt.Sprintf("htpasswd -b %s.%d %s \"%s\"", options.html_passwords, game_number, box_code_less(pl), pw)

	system(buf)
}

func output_html_rep(pl int) {
	var fnam, report string
	var ret int

	fnam = filepath.Join(libdir, "save", fmt.Sprintf("%d", sysclock.turn), fmt.Sprintf("%d", pl))
	report = filepath.Join(options.html_path, fmt.Sprintf("%d", game_number), fmt.Sprintf("tag%d-%d.html", box_code_less(pl), game_number, sysclock.turn))
	ret = make_report(HTML, fnam, report, pl)
	if ret != 0 {
		log.Printf("Cannot make HTML report for %s?", box_code_less(pl))
		return
	}

	/*
	 *  Compress the output.
	 *
	 */
	system(sout("gzip -f %s", report))
	/* VLN

	     system(sout("rm %s/%d/%s/index.html; ln -s %s.gz %s/%d/%s/index.html",
	   	      options.html_path, game_number, box_code_less(pl),
	   	      report, options.html_path, game_number, box_code_less(pl)));
	*/
}

/*
 *  Mon May 14 13:53:57 2001 -- Scott Turner
 *
 *  Add in what we need to create the info subdirectory and the
 *  map.
 *
 */
func output_html_map(pl int) {
	var info, fnam string
	minx := 1000
	miny := 1000
	maxx := 0
	maxy := 0
	var i int

	fnam = filepath.Join(libdir, "save", fmt.Sprintf("%d", sysclock.turn), fmt.Sprintf("%d", pl))
	info = filepath.Join(options.html_path, fmt.Sprintf("%d", game_number), box_code_less(pl), "info")
	// create info subdir if necessary.
	umask(S_IWGRP | S_IWOTH)
	if err := mkdir(info); err != nil {
		panic(err)
	}
	umask(S_IWOTH | S_IXOTH)

	/*
	 *  Run the analyze script on the current turn report, outputting
	 *  to the info directory.
	 *
	 */
	if system(sout("anlz.pl %s %s", fnam, info)) != 0 {
		log.Printf("Cannot run anlyz for %s?", box_code_less(pl))
		umask(S_IRWXO)
		return
	}

	/*
	 *  Now run map3.pl.
	 *
	 *  Mon May 14 17:47:14 2001 -- Scott Turner
	 *
	 *  Have to figure out min/max values for x and y.
	 */
	for _, i = range known_sparse_loop(p_player(pl).known) {
		if kind(i) == T_loc &&
			loc_depth(i) == LOC_province &&
			region(i) != hades_region &&
			region(i) != faery_region &&
			region(i) != cloud_region {
			if region_col(i) < minx {
				minx = region_col(i)
			}
			if region_col(i) > maxx {
				maxx = region_col(i)
			}
			if region_row(i) < miny {
				miny = region_row(i)
			}
			if region_row(i) > maxy {
				maxy = region_row(i)
			}
		}
	}

	for _, i = range known_sparse_loop(p_player(pl).locs) {
		if kind(i) == T_loc &&
			loc_depth(i) == LOC_province &&
			region(i) != hades_region &&
			region(i) != faery_region &&
			region(i) != cloud_region {
			if region_col(i) < minx {
				minx = region_col(i)
			}
			if region_col(i) > maxx {
				maxx = region_col(i)
			}
			if region_row(i) < miny {
				miny = region_row(i)
			}
			if region_row(i) > maxy {
				maxy = region_row(i)
			}
		}
	}

	if system(sout("cd %s;map3.pl map.html %d %d %d %d full", info, minx, maxx, miny, maxy)) != 0 {
		log.Printf("Cannot run map3 for %s?", box_code_less(pl))
		umask(S_IRWXO)
		return
	}
	umask(S_IRWXO)
}

//#if 0
//copy_public_turns()
//{
//    var fnam string
//    char cmd[LEN];
//    var pl int
//
//    for _, pl = range loop_player(pl)
//    {
//        if (!player_public_turn(pl))
//            continue;
//
//        fnam = filepath.Join( "%s/html/%s", libdir, box_code_less(pl));
//
//        sprintf(cmd, "sed -e '/Account summary/,/Balance/d' -e 's/^begin %s.*$/begin %s/' %s/index.html > %s.html",
//                box_code_less(pl), box_code_less(pl),
//                fnam, fnam);
//        system(cmd);
//    }
//
//}
//#endif
