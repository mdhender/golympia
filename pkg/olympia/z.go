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
	"github.com/mdhender/golympia/pkg/io"
)

// z.h

const TRUE = 1
const FALSE = 0

const LEN = 512 /* generic string max length */

func fuzzy_strcmp(s, t string) int {
	panic("!implemented")
}

func i_strcmp(s, t string) int {
	panic("!implemented")
}

func i_strncmp(s, t string, n int) int {
	panic("!implemented")
}

// z.c

// Line reader with no size limits strips newline off end of line
func getlin(fp *io.FILE) (string, bool) {
	return fp.GetLine()
}

// Get line, remove leading and trailing whitespace
func getlin_ew(fp *io.FILE) (string, bool) {
	return fp.GetLineEW()
}

func readlin() (string, bool) {
	return io.ReadLine()
}

func readlin_ew() (string, bool) {
	return io.ReadLineEW()
}

func str_save(s string) string {
	return s
}
