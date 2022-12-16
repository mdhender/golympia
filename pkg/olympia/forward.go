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
	"os"
	"strconv"
)

func d_build_wagons(c *command) int                 { panic("!implemented") }
func d_cloak_creat(c *command) int                  { panic("!implemented") }
func d_cloak_reg(c *command) int                    { panic("!implemented") }
func d_curse_noncreat(c *command) int               { panic("!implemented") }
func d_exhume(c *command) int                       { panic("!implemented") }
func d_move(c *command) int                         { panic("!implemented") }
func d_show_art_creat(c *command) int               { panic("!implemented") }
func d_show_art_reg(c *command) int                 { panic("!implemented") }
func d_teach(c *command) int                        { panic("!implemented") }
func delta_loyalty(who int, amount int, silent int) { panic("!implemented") }
func distance(orig, dest, gate int) int             { panic("!implemented") }
func fetch_inside_name() string                     { panic("!implemented") }
func first_char_here(where int) int                 { panic("!implemented") }
func get_process_id() int                           { panic("!implemented") }
func get_rid_of_building(fort int)                  { panic("!implemented") }
func i_petty_thief(c *command) int                  { panic("!implemented") }
func immediate_commands()                           { panic("!implemented") }
func int_comp(q1, q2 interface{}) int      { panic("!implemented") }
func is_artifact(item int) *EntityArtifact { panic("!implemented") }
func is_port_city_where(where int) int     { panic("!implemented") }
func my_free(ptr interface{})                       { panic("!implemented") }
func my_malloc(size int)                            { panic("!implemented") }
func my_realloc(ptr interface{}, size int)          { panic("!implemented") }
func mylog(base int, num int) int                   { panic("!implemented") }
func random() int                                   { panic("!implemented") }
func rename(from, to string) error                  { panic("!implemented") }
func srandom(seed uint32)                           { panic("!implemented") }
func test_random() int                              { panic("!implemented") }
func v_acquire(c *command) int                      { panic("!implemented") }
func v_attack(c *command) int                       { panic("!implemented") }
func v_build_wagons(c *command) int                 { panic("!implemented") }
func v_cloak_creat(c *command) int                  { panic("!implemented") }
func v_cloak_reg(c *command) int                    { panic("!implemented") }
func v_curse_noncreat(c *command) int               { panic("!implemented") }
func v_exhume(c *command) int                       { panic("!implemented") }
func v_move(c *command) int                         { panic("!implemented") }
func v_plugh(c *command) int                        { panic("!implemented") }
func v_rem_art_cloak(c *command) int                { panic("!implemented") }
func v_show_art_creat(c *command) int               { panic("!implemented") }
func v_show_art_reg(c *command) int                 { panic("!implemented") }
func wield_s(who int) string                        { panic("!implemented") }

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func atoi_b(b []byte) int {
	return atoi(string(b))
}

func mkdir(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if err = os.Mkdir(path, 0755); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("not a directory")
	}
	return nil
}

func or_float(t bool, a, b float64) float64 {
	if t {
		return a
	}
	return b
}

func or_int(t bool, a, b int) int {
	if t {
		return a
	}
	return b
}

func or_string(t bool, a, b string) string {
	if t {
		return a
	}
	return b
}

func rmdir(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		// todo: this should handle errors
		return nil
	} else if !fi.IsDir() {
		return fmt.Errorf("rmdir: %q: not a directory", path)
	} else if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("rmdir: %w", err)
	}
	return nil
}
