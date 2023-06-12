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

type Tile struct {
	Row, Col int // map tile we're inside

	City     int
	Color    int // map coloring for output
	Depth    int
	Gates    []*Gate
	Inside   int
	Mark     int
	Name     string
	Region   int
	Roads    []*Road
	SaveChar byte
	Subs     []int
	Terrain  int
	Is       struct {
		Hidden         bool
		MajorCity      bool
		RegionBoundary bool
		SafeHaven      bool
		SeaLane        bool
		Summerbridge   bool
		Uldim          bool
	}
}

