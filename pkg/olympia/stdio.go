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
	"github.com/mdhender/golympia/pkg/io"
)

func fclose(fp *io.FILE) {
	fp.Close()
}

func fopen(name string, mode string) (*io.FILE, error) {
	switch mode {
	case "w":
		return io.NewWriter(name)
	}
	return nil, fmt.Errorf("fopen: unknown mode %q", mode)
}
func fprintf(fp *io.FILE, format string, args ...interface{}) {
	fp.Printf(format, args...)
}