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
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	DASH_LINE = "===-===-===-===-===-===-===-===-===-===-===-===-===-===-===-===-===-===-=\n"
	months    = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	press_fp  *os.File
	rumor_fp  *os.File
)

/*
 *  Fri Mar  6 09:53:39 1998 -- Scott Turner
 *
 *  Everything to do with the newsletter.
 *
 */

func close_times() {
	press_fp = fclose(press_fp)
	rumor_fp = fclose(rumor_fp)
}

func copy_file(src, dst string) {
	r, err := os.Open(src)
	if err != nil {
		return // ignore errors?
	}

	w, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return // ignore errors?
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return // ignore errors?
	}
}

func copy_file_slow(src io.Reader, dst io.Writer) {
	if src == nil || dst == nil {
		return
	}
	_, err := io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
}

func do_times() {
	times_goal_info()
	times_masthead()
	close_times()
	times_index()

	w, err := os.OpenFile(filepath.Join(libdir, "Times"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	// for filename in libdir
	//   if filename like *~ then skip
	//   if filename not like times* then skip
	//   append filename to libdir/Times
	files, err := os.ReadDir(libdir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "times") && !strings.HasSuffix(f.Name(), "~") {
			if r, err := os.Open(f.Name()); err != nil {
				log.Printf("do_times: %+v\n", err)
			} else if _, err = io.Copy(w, r); err != nil {
				log.Printf("do_times: %+v\n", err)
			} else {
				_ = r.Close()
			}
		}
	}
}

func open_times() {
	var err error

	if press_fp == nil {
		press_fp, err = os.Create(filepath.Join(libdir, "times_press"))
		if err != nil {
			panic(err)
		}
	}
	fprintf(press_fp, "\n%s\n", DASH_LINE)
	fprintf(press_fp, "                        Player-contributed press\n\n")
	fprintf(press_fp, "%s\n", DASH_LINE)

	if rumor_fp == nil {
		rumor_fp, err = os.Create(filepath.Join(libdir, "times_rumor"))
		if err != nil {
			panic(err)
		}
	}

	fprintf(rumor_fp, "                                 Rumors\n\n")
	fprintf(rumor_fp, "%s\n", DASH_LINE)
}

// this generates the section of the Times that covers how the various game goals are doing.
func times_goal_info() {
	fp, err := fopen(filepath.Join(libdir, "times_info"), "w")
	if err != nil {
		panic(err)
	}

	fprintf(fp, "\nGame Information\n")
	fprintf(fp, "================\n\n")

	fprintf(fp, "  Nations Summary\n")
	fprintf(fp, "  ---------------\n\n")

	calculate_nation_nps()
	for _, i := range loop_kind(T_nation) {
		fprintf(fp, "    %s:\n", rp_nation(i).name)
		fprintf(fp, "        * Next player joining starts with %d NPs.\n",
			starting_noble_points(i))
		fprintf(fp, "        * Next player joining starts with %d gold.\n", starting_gold(i))
		fprintf(fp, "        * Players: %d.\n", rp_nation(i).players)
		fprintf(fp, "        * Nobles (and controlled units): %d.\n", rp_nation(i).nobles)
		fprintf(fp, "        * Noble points: %d.\n", rp_nation(i).nps)

		total1 := 0
		for _, j := range loop_province() {
			if nation(player_controls_loc(j)) == i {
				total1++
			}
		}
		fprintf(fp, "        * Provinces controlled: %d.\n", total1)
		total1 = 0
		for _, j := range loop_city() {
			if nation(player_controls_loc(j)) == i {
				total1++
			}
		}
		fprintf(fp, "        * Cities controlled: %d.\n", total1)
		total1 = 0
		for _, j := range loop_castle() {
			if nation(player_controls_loc(j)) == i {
				total1++
			}
		}
		fprintf(fp, "        * Castles controlled: %d.\n", total1)
	}

	if !options.open_ended {
		for _, i := range loop_nation() {
			if rp_nation(i).win == 1 {
				fprintf(fp, "\n  *********************************************************\n")
				fprintf(fp, "      The %s will win at the end of next turn if\n", rp_nation(i).name)
				fprintf(fp, "      they maintain the win conditions.\n\n")
				fprintf(fp, "  *********************************************************\n")
			} else if rp_nation(i).win == 2 {
				fprintf(fp, "  ***********************************************\n")
				fprintf(fp, "  ***********************************************\n")
				fprintf(fp, "    Congratulations to the %s!  They have\n", rp_nation(i).name)
				fprintf(fp, "    conquered all enemies and have become the most\n")
				fprintf(fp, "    powerful nation in the history of Olympia!\n")
				fprintf(fp, "    Final NP Statistics:\n")
				fprintf(fp, "       %s: %d NPs.\n", rp_nation(i).name, rp_nation(i).nps)
				for _, j := range loop_nation() {
					if j != i {
						fprintf(fp, "       %s: %d NPs.\n", rp_nation(j).name,
							rp_nation(j).nps)
					}
				}
				fprintf(fp, "\n    Winning Players:\n")
				for _, j := range loop_pl_regular() {
					if nation(j) == i {
						fprintf(fp, "      %s\n", rp_player(j).full_name)
					}
				}
				fprintf(fp, "  ***********************************************\n")
				fprintf(fp, "  ***********************************************\n")
			}
		}
	}

	if options.mp_antipathy {
		fprintf(fp, "\n  Staff of the Sun Summary\n")
		fprintf(fp, "  ------------------------\n\n")

		/*
		 *  How many pieces are out there?
		 *
		 */
		n, unfound, priests, mus, others := 0, 0, 0, 0, 0
		for _, i := range loop_subkind(sub_special_staff) {
			n++
			if kind(item_unique(i)) != T_char {
				unfound++
			} else {
				if is_priest(item_unique(i)) != FALSE {
					priests++
				} else if is_magician(item_unique(i)) {
					mus++
				}
				others++
			}
		}

		if priests == 0 && mus == 0 && others == 0 {
			fprintf(fp, "    * No pieces of the Staff of the Sun have been found.\n")
		} else {
			if others > 0 {
				fprintf(fp, "  * %s pieces of the Staff of the Sun are at large.\n", cap_(nice_num(others)))
				for _, i := range loop_subkind(sub_special_staff) {
					if kind(item_unique(i)) == T_char && is_magician(item_unique(i)) && is_priest(item_unique(i)) == FALSE {
						fprintf(fp, "    * One piece is held somewhere on %s.\n", just_name(region(item_unique(i))))
					}
				}
			}
			if priests > 0 {
				if priests == 1 {
					fprintf(fp, "  * %s pieces of the Staff of the Sun is held by priests.\n", nice_num(priests))
				} else {
					fprintf(fp, "  * %s pieces of the Staff of the Sun are held by priests.\n", nice_num(priests))
				}
				for _, i := range loop_subkind(sub_special_staff) {
					if kind(item_unique(i)) == T_char && is_priest(item_unique(i)) != FALSE {
						fprintf(fp, "    * One piece is held on %s.\n", just_name(region(item_unique(i))))
					}
				}
			}
			if mus > 0 {
				if mus == 1 {
					fprintf(fp, "  * %s piece of the Staff of the Sun is held by magicians.\n", nice_num(mus))
				} else {
					fprintf(fp, "  * %s pieces of the Staff of the Sun are held by magicians.\n", nice_num(mus))
				}
				for _, i := range loop_subkind(sub_special_staff) {
					if kind(item_unique(i)) == T_char && is_magician(item_unique(i)) {
						fprintf(fp, "    * One piece is held on %s.\n", just_name(region(item_unique(i))))
					}
				}
			}
		}
	}

	fprintf(fp, "\n                                *  *  *\n\n")
	fprintf(fp, "                                *  *  *\n\n")

	fp = fclose(fp)
}

func times_credit(c *command) bool {
	if options.times_pay == FALSE {
		return true
	}

	pl := player(c.who)
	if times_paid(pl) {
		wout(c.who, "The Times has already paid faction %s this month.", box_name(pl))
		return false
	}

	p_player(pl).times_paid = TRUE
	wout(pl, "The Times adds %s to your CLAIM.", gold_s(options.times_pay))
	gen_item(pl, item_gold, options.times_pay)
	gold_times += options.times_pay

	return true
}

func times_index() {
	/*
	 *  Open the output file.
	 *
	 */
	fp, err := fopen(filepath.Join(libdir, "index_times_middle"), "w")
	if err != nil {
		panic(err)
	}
	/*
	 *  Put in the list of back issues.
	 *
	 */
	for i := sysclock.turn - 1; i > 0; i-- {
		fprintf(fp, "<A HREF=\"%02d.html\">%02d</A>\n", i, i)
	}
	fclose(fp)

}

func times_masthead() {
	fp, err := os.Create(filepath.Join(libdir, "times_0"))
	if err != nil {
		panic(err)
	}

	now := time.Now().UTC()
	date := now.Format("2006 01 02")

	for _, i := range loop_kind(T_player) {
		if subkind(i) == sub_pl_regular {
			nplayers++
		}
	}

	turn_s := fmt.Sprintf("Turn %d  %d Players", sysclock.turn, nplayers)

	month := oly_month(&sysclock)
	year := oly_year(&sysclock)

	issue_s := fmt.Sprintf("Game %d, Season \"%s\", month %d, in the year %d.", game_number, month_names[month], month+1, year+1)

	fprintf(fp, "From: %s\n", from_host)
	fprintf(fp, "Subject: The Gods Speak (Game %d, Issue %d)\n", game_number, sysclock.turn)
	/* 	fprintf(fp, "To: UPDATE email address here\n\n",game_number); */
	fprintf(fp, "\n   +----------------------------------------------------------------------+\n")
	fprintf(fp, "   | The Gods Speak %53s |\n", date)
	fprintf(fp, "   | %-68s |\n", issue_s)
	fprintf(fp, "   |                                                                      |\n")
	fprintf(fp, "   | %-40s http://olytag.com           |\n", turn_s)
	fprintf(fp, "   |                    Send orders to: tagtest (at) olytag (dot) com     |\n", game_number) /* UPDATE */
	fprintf(fp, "   +----------------------------------------------------------------------+\n\n")
	fprintf(fp, "           Questions, comments, to play:  moderator (at) olytag (dot) com\n\n") /*UPDATE*/
	fprintf(fp, "                             Olympia PBEM\n\n")
	fprintf(fp, "                                *  *  *\n\n")
	fprintf(fp, "                                *  *  *\n\n")

	fp = fclose(fp)
}

func v_press(c *command) int {
	l := parse_text_list(c)
	if len(l) == 0 {
		return FALSE
	}

	if line_length_check(l) > 78 {
		wout(c.who, "Line length of message text exceeds %d characters.", 78)
		wout(c.who, "Post rejected.")
		return FALSE
	}

	for i := 0; i < len(l); i++ {
		if bytes.HasPrefix(l[i], []byte("===-")) {
			fprintf(press_fp, "> %s\n", l[i])
		} else {
			fprintf(press_fp, "%s\n", l[i])
		}
	}

	/*
	 *  Turn off tags in the Times?
	 *
	 */
	tags_off()
	attrib := fmt.Sprintf("-- %s", box_name(player(c.who)))
	tags_on()
	attrib = strings.ReplaceAll(attrib, "~", " ")
	fprintf(press_fp, "\n%55s\n\n", attrib)

	fprintf(press_fp, DASH_LINE)
	fprintf(press_fp, "\n")
	fflush(press_fp)

	wout(c.who, "Press posted.")
	times_credit(c)
	return TRUE
}

func v_rumor(c *command) int {
	l := parse_text_list(c)
	if len(l) == 0 {
		return FALSE
	}

	if line_length_check(l) > 78 {
		wout(c.who, "Line length of message text exceeds %d characters.", 78)
		wout(c.who, "Post rejected.")
		return FALSE
	}

	for i := 0; i < len(l); i++ {
		if bytes.HasPrefix(l[i], []byte("===-")) {
			fprintf(rumor_fp, "> %s\n", l[i])
		} else {
			fprintf(rumor_fp, "%s\n", l[i])
		}
	}

	fprintf(rumor_fp, "\n")
	fprintf(rumor_fp, DASH_LINE)
	fprintf(rumor_fp, "\n")
	fflush(rumor_fp)

	wout(c.who, "Rumor posted.")
	times_credit(c)
	return TRUE
}
