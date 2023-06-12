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

package maps

import (
	"bytes"
	"fmt"
	"github.com/mdhender/golympia/pkg/io"
	"log"
	"unicode/utf8"
)

type Cell struct {
	Code rune
}

var (
	minRow, maxRow, minCol, maxCol = 0, 0, 0, 0
)

func Read(name string) ([][]rune, error) {
	cells, err := read(name)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return cells, nil
}

func read(name string) ([][]rune, error) {
	// trim trailing spaces from every line, stopping at the first empty line
	var lines [][]byte
	data, err := io.ReadLines(name)
	if err != nil {
		panic(err)
	}
	for _, line := range data {
		line = bytes.TrimRight(line, " \t\r\n")
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}

	// find the size of the map
	maxRow = len(lines)
	for _, line := range lines {
		col := 0
		for ; len(line) != 0; col++ {
			_, w := utf8.DecodeRune(line)
			line = line[w:]
		}
		if maxCol < col {
			maxCol = col
		}
	}

	// create the cells and initialize them
	cells := make([][]rune, maxRow, maxRow)
	for row := 0; row < maxRow; row++ {
		cells[row] = make([]rune, maxCol, maxCol)
		for col := 0; col < maxCol; col++ {
			cells[row][col] = '#'
		}
	}

	for row, line := range lines {
		for col := 0; len(line) != 0; col++ {
			r, w := utf8.DecodeRune(line)
			if r != utf8.RuneError {
				cells[row][col] = r
			}
			line = line[w:]
		}
	}

	log.Printf("read: min (%d, %d): max(%d, %d)\n", minRow, minCol, maxRow, maxCol)

	return cells, nil
}