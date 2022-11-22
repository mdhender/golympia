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
)

// BoxAlloc replaces alloc_box()
func BoxAlloc(id, kind, skind int) {
	if !(0 < id && id < MAX_BOXES) {
		panic(fmt.Sprintf("assert(0 < %d < MAX_BOXES)", id))
	} else if !(bx[id] == nil) {
		panic(fmt.Sprintf("assert(bx[%]d == nil)", id))
	}

	bx[id] = &box{
		kind:  schar(kind),
		skind: schar(skind),
	}
}
