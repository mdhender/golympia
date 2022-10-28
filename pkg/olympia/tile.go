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
 * aint with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package olympia

// tile.h

type tile struct {
	save_char         byte
	region            int
	name              string
	terrain           int
	hidden            int
	city              int
	mark              int
	inside            int
	color             int /* map coloring for */
	row, col          int /* map tile we're inside */
	depth             int
	safe_haven        int
	sea_lane          int
	uldim_flag        int
	summerbridge_flag int
	region_boundary   int
	major_city        int

	subs       ilist
	gates_dest ilist /* gates from here */
	gates_num  ilist /* gates from here */
	gates_key  ilist
	roads      rlist
}

type tlist []*tile
type tile_t = tile
