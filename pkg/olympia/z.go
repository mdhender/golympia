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
	"io"
	"log"
	"os"
)

// z.h

const TRUE = 1
const FALSE = 0

const LEN = 512 /* generic string max length */

var (
	line_fd     *os.File
	lower_array [256]byte
)

func asfail(file string, line int, cond string) {
	log.Printf("assertion failure: %s (%d): %s\n", file, line, cond)
	os.Exit(2)
}

// copy_fp copies data from src to dst, ignoring all errors
func copy_fp(src io.Reader, dst io.Writer) {
	if src != nil && dst != nil {
		_, _ = io.Copy(dst, src)
	}
}

// return slice without leading and trailing whitespace.
// does not update any data in original slice.
func eat_leading_trailing_whitespace(s []byte) []byte {
	for len(s) != 0 && iswhite(s[0]) { // trim leading whitespace
		s = s[1:]
	}
	for len(s) != 0 && iswhite(s[len(s)-1]) { // trim trailing whitespace
		s = s[:len(s)-1]
	}
	return s
}

// fuzzy_one_bad returns true if s1 and s2 are the same except for one character.
// panics if l1 or l2 are longer than the length of s1 or s2.
func fuzzy_one_bad(s1, s2 []byte, l1, l2 int) bool {
	if l1 != l2 {
		return false
	}
	count := 0
	for i := 0; i < l2; i++ {
		if tolower(s1[i]) != tolower(s2[i]) {
			if count != 0 {
				return false
			}
			count++
		}
	}
	return true
}

// fuzzy_one_extra returns true if s1 and s2 are the same except s2 is has one extra character.
// panics if l1 or l2 are longer than the length of s1 or s2.
func fuzzy_one_extra(s1, s2 []byte, l1, l2 int) bool {
	if l1 != l2-1 {
		return false
	}
	count := 0
	for i, j := 0, 0; i < l1; i, j = i+1, j+1 {
		if tolower(s1[i]) != tolower(s2[j]) {
			if count != 0 {
				return false
			}
			i, count = i-1, count+1
		}
	}
	return true
}

// fuzzy_one_less returns true if s1 and s2 are the same except s2 is missing one character.
// panics if l1 or l2 are longer than the length of s1 or s2.
func fuzzy_one_less(s1, s2 []byte, l1, l2 int) bool {
	if l1 != l2+1 {
		return false
	}
	count := 0
	for i, j := 0, 0; j < l2; i, j = i+1, j+1 {
		if tolower(s1[i]) != tolower(s2[j]) {
			if count != 0 {
				return false
			}
			j, count = j-1, count+1
		}
	}
	return true
}

// returns true if s1 and s2 are sort of the same.
func fuzzy_strcmp(s1, s2 []byte) bool {
	if bytes.Compare(s1, s2) == 0 {
		return true
	} else if l1, l2 := len(s1), len(s2); l2 >= 4 && fuzzy_transpose(s1, s2, l1, l2) {
		return true
	} else if l2 >= 5 && fuzzy_one_less(s1, s2, l1, l2) {
		return true
	} else if l2 >= 5 && fuzzy_one_extra(s1, s2, l1, l2) {
		return true
	} else if l2 >= 5 && fuzzy_one_bad(s1, s2, l1, l2) {
		return true
	}
	return false
}

func fuzzy_strcmp_bs(b []byte, s string) bool {
	return fuzzy_strcmp(b, []byte(s))
}

// fuzzy_transpose returns true if the strings are the same length
// and are the same except for a transposition of a single character.
// panics if l1 or l2 are longer than the length of s1 or s2.
func fuzzy_transpose(s1, s2 []byte, l1, l2 int) bool {
	same := false
	if l1 == l2 {
		for i := 0; !same && i < l2-1; i++ {
			s2[i], s2[i+1] = s2[i+1], s2[i]
			same = i_strcmp(s1, s2) == 0
			s2[i], s2[i+1] = s2[i+1], s2[i]
		}
	}
	return same

}

// line reader with no size limits.
// ignores carriage-return and strips newline off end of line.
// returns nil on end of input.
func getlin(r io.Reader) []byte {
	var buf []byte
	if r != nil {
		p := []byte{0}
		for {
			if n, err := r.Read(p); err != nil {
				break
			} else if n == 0 {
				panic("assert(n && !err)")
			} else if p[0] == '\n' {
				break
			} else if p[0] == '\r' {
				continue
			}
			buf = append(buf, p[0])
		}
	}
	return buf
}

// get line, remove leading and trailing whitespace,
// and changing control-characters to spaces.
func getlin_ew(r io.Reader) []byte {
	s := getlin(r)
	for i := 0; i < len(s); i++ { // remove ctrl chars
		if s[i] < 32 || s[i] == '\r' || s[i] == '\t' {
			s[i] = ' '
		}
	}
	return eat_leading_trailing_whitespace(s)
}

// i_strcmp compares s and t, ignoring case.
// it returns -1 if s < t, 0 if s == t, and 1 if s > t
func i_strcmp(s, t []byte) int {
	for {
		if len(s) == 0 && len(t) == 0 {
			return 0
		} else if len(s) == 0 {
			return -1
		} else if len(t) == 0 {
			return 1
		} else if d := tolower(s[0]) - tolower(t[0]); d < 0 {
			return -1
		} else if d > 0 {
			return 1
		}
		s, t = s[1:], t[1:]
	}
}

// i_strcmp compares the first n bytes of s and t, ignoring case.
// it returns -1 if s < t, 0 if s == t, and 1 if s > t
func i_strncmp(s, t []byte, n int) int {
	if len(s) > n {
		s = s[:n]
	}
	if len(t) > n {
		t = t[:n]
	}
	return i_strcmp(s, t)
}

func strcasecmp(s, t []byte) int {
	return i_strcmp(s, t)
}

func strcasecmp_bs(b []byte, s string) int {
	return strcasecmp(b, []byte(s))
}

func strcasecmp_sb(s string, b []byte) int {
	return strcasecmp([]byte(s), b)
}

func strcmp(s, t []byte) int {
	return i_strcmp(s, t)
}

func strcmp_sb(s string, b []byte) int {
	return i_strcmp([]byte(s), b)
}

func strncasecmp(s, t []byte, n int) int {
	return i_strncmp(s, t, n)
}

func strncasecmp_bs(b []byte, s string, n int) int {
	return i_strncmp(b, []byte(s), n)
}

func strncmp(s, t string, n int) int {
	return i_strncmp([]byte(s), []byte(t), n)
}

func strncmp_bs(s []byte, t string, n int) int {
	return i_strncmp(s, []byte(t), n)
}

func init_lower() {
	for i := range lower_array {
		if 'A' <= i && i <= 'Z' {
			lower_array[i] = byte(i + 'a' - 'A')
			continue
		}
		lower_array[i] = byte(i)
	}
}

func init_random() {
	// mdhender: do nothing
}

func lcase(s []byte) []byte {
	for i, ch := range s {
		s[i] = tolower(ch)
	}
	return s
}

func readfile(path string) bool {
	var err error
	if line_fd != nil {
		_ = line_fd.Close()
	}
	line_fd, err = os.Open(path)
	if err != nil {
		log.Printf("readfile: can't open %q: %v", path, err)
		panic("!")
		return false
	}

	return true
}

// readlin reads the next line from the global file pointer `line_fd`.
// ignores carriage-return and strips newline off end of line.
// returns nil on end of input.
func readlin() []byte {
	return getlin(line_fd)
}

func readlin_ew() []byte {
	return getlin_ew(line_fd)
}

func str_save(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
