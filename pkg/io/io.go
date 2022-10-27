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

// Package io is my attempt at the global input/output for mapgen.
package io

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadLines(name string) ([][]byte, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("io::ReadLines: %w", err)
	}

	var lines [][]byte
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		lines = append(lines, bytes.Trim(line, " \t\r\n"))
	}

	return lines, nil
}

func Open(name string) (*FILE, error) {
	panic("!")
}

func NewReader(name string) (*FILE, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return &FILE{rw: bufio.NewReadWriter(bufio.NewReader(bytes.NewReader(buf)), nil)}, nil
}

func NewWriter(name string) (*FILE, error) {
	fp, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &FILE{rw: bufio.NewReadWriter(nil, bufio.NewWriter(fp))}, nil
}

const GETLIN_ALLOC = 255

type FILE struct {
	rw *bufio.ReadWriter
}

func (fp *FILE) Close() {
	if fp == nil || fp.rw == nil {
		return
	}
	// todo: write out updates?
	fp.rw = nil
}

// GetLine is z::getlin()
// Line reader with no size limits strips newline off end of line.
func (fp *FILE) GetLine() (string, bool) {
	if fp == nil || fp.rw == nil {
		return "", false
	}

	buf, err := fp.rw.ReadBytes('\n')
	if err != nil {
		panic(err)
	}

	// strip cr, lf, trailing spaces and tabs
	return string(bytes.TrimRight(buf, " \r\t\n")), true
}

// GetLineEW is z::getlin_ew()
// Returns the next line from the file with leading and trailing spaces trimmed
func (fp *FILE) GetLineEW() (string, bool) {
	s, ok := line_fd.GetLine()
	return strings.TrimSpace(s), ok
}

func (fp *FILE) Printf(format string, args ...interface{}) {
	if fp == nil {
		return
	}
	// todo: output the string...
	//fp.fputs(fmt.Sprintf(format, args...))
}

var (
	line_fd *FILE
	nread   int
	point   []byte
)

// ReadFile uses, reuses, and abuses a global file descriptor.
func ReadFile(name string) bool {
	if line_fd != nil {
		line_fd.Close()
	}

	var err error
	line_fd, err = NewReader(name)
	if err != nil {
		log.Printf("io::ReadFile: %+v\n", err)
		return false
	}

	return true
}

// ReadLine is z::readlin()
// ReadLine returns the next line from the global line_fd file.
func ReadLine() (string, bool) {
	return line_fd.GetLine()
}

// ReadLineEW is z::readlin_ew()
// ReadLineEW returns the next line from the global line_fd file with leading and trailing spaces trimmed
func ReadLineEW() (string, bool) {
	return line_fd.GetLineEW()
}
