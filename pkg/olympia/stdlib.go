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

var stdout, stderr *io.FILE

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type AccessFlag int

const (
	R_OK AccessFlag = iota
)

func access(path string, flag AccessFlag) int {
	panic(path)
}

func assert(t bool) {
	if !t {
		panic("assert() failed")
	}
}

func bcopy(src []byte) (cp []byte) {
	return append(cp, src...)
}

func bzero(dst []byte, length int) {
	panic("!implemented")
}

func isalpha(c byte) bool {
	return (((c) >= 'a' && (c) <= 'z') || ((c) >= 'A' && (c) <= 'Z'))
}

func isdigit(c byte) bool {
	return ((c) >= '0' && (c) <= '9')
}

func iswhite(c byte) bool {
	return ((c) == ' ' || (c) == '\t')
}

func system(cmd string) int {
	panic(cmd)
}

func tolower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c = c - 'A' + 'a'
	}
	return c
}

func toupper(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c = c - 'a' + 'A'
	}
	return c
}

type UMaskFlag int

const (
	S_IWGRP UMaskFlag = iota
	S_IXGRP
	S_IWOTH
	S_IXOTH
	S_IRWXO
)

func umask(flag UMaskFlag) {
	panic("!implemented")
}

func unlink(path string) int {
	panic(path)
}
