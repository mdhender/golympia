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

type RoadList []*Road

func (r RoadList) Len() int {
	return len(r)
}

func (r RoadList) Less(i, j int) bool {
	return r[i].Id < r[j].Id
}

func (r RoadList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type Road struct {
	Id     int    `json:"id"`             // identity of the road
	Name   string `json:"name,omitempty"` // name of the road
	Where  int    `json:"where"`          // where this road is located (region or location)
	To     int    `json:"to"`             // identity of connected destination
	Hidden bool   `json:"hidden,omitempty"`
}

type RoadLink struct {
	Key int `json:"key,omitempty"` // identity of key required to use road
}

// RoadsFromMapGen loads roads from the globals created by the map generator.
func RoadsFromMapGen() (roads RoadList) {
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] == nil {
				continue
			}
			for j := 0; j < len(map_[row][col].roads); j++ {
				roads = append(roads, &Road{
					Id:     map_[row][col].roads[j].ent_num,
					Name:   map_[row][col].roads[j].name,
					Where:  map_[row][col].region,
					To:     map_[row][col].roads[j].to_loc,
					Hidden: map_[row][col].roads[j].hidden != FALSE,
				})
			}
		}
	}

	for _, tile := range subloc_mg {
		if tile == nil {
			continue
		}
		for _, road := range tile.roads {
			roads = append(roads, &Road{
				Id:     road.ent_num,
				Name:   road.name,
				Where:  tile.region,
				To:     road.to_loc,
				Hidden: road.hidden != FALSE,
			})
		}
	}

	sort.Sort(roads)

	return roads
}
