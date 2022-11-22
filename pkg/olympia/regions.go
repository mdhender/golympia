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

import "sort"

type RegionList []*Region

func (r RegionList) Len() int {
	return len(r)
}

func (r RegionList) Less(i, j int) bool {
	return r[i].Id < r[j].Id
}

func (r RegionList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type Region struct {
	Id       int    `json:"id"`                 // identity of the region
	Name     string `json:"name,omitempty"`     // name of the region
	Contains ints_l `json:"contains,omitempty"` // here list of inside locations
}

// ContinentsFromMapGen loads continents from the globals created by the map generator.
func ContinentsFromMapGen() (regions RegionList) {
	for i := 1; i <= inside_top; i++ {
		r := &Region{
			Id:   REGION_OFF + i,
			Name: inside_names[i],
		}
		for _, e := range inside_list[i] {
			if e == nil {
				continue
			}
			r.Contains = append(r.Contains, e.region)
		}
		sort.Ints(r.Contains)
		regions = append(regions, r)
	}

	sort.Sort(regions)

	return regions
}
