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
	"github.com/mdhender/golympia/pkg/io"
	"log"
	"os"
)

type olyfile struct {
	name  string
	pos   int
	lines []*olyline
}

type olyline struct {
	no   int
	line []byte
}

// readfile opens the file, loads it, splits it into lines,
// trims trailing spaces, tabs and carriage-returns from each line, and returns those lines.
// if there is an error loading the file, nil and false are returned.
func oly_readfile(filename string) (*olyfile, bool) { // src/z.c
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("readfile: %+v\n", err)
		return nil, false
	}
	of := &olyfile{name: filename}
	for no, line := range bytes.Split(data, []byte{'\n'}) {
		of.lines = append(of.lines, &olyline{no: no + 1, line: bcopy(bytes.TrimRight(line, " \r\t"))})
	}
	return of, true
}

func oly_readbytes(filename string, data []byte) (*olyfile, bool) {
	of := &olyfile{name: filename}
	for no, line := range bytes.Split(data, []byte{'\n'}) {
		of.lines = append(of.lines, &olyline{no: no + 1, line: bcopy(bytes.TrimRight(line, " \r\t"))})
	}
	return of, true
}

// fgets returns the next line from the buffer
func (of *olyfile) fgets() *olyline {
	if of == nil || len(of.lines) == 0 || of.pos >= len(of.lines) {
		return nil
	}
	line := of.lines[of.pos]
	of.pos++
	return line
}

func (of *olyfile) fputs(s string) {
	panic("!implemented")
}

func (of *olyfile) readlin() (string, bool) {
	return io.ReadLine()
}

// ungets backs up a line in the file, if possible.
func (of *olyfile) ungets() {
	if of == nil || of.pos == 0 {
		return
	}
	of.pos--
}

func (l *olyline) String() string {
	if l == nil {
		return ""
	}
	return string(l.line)
}
