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

type LocationList []*Location

func (l LocationList) Len() int {
	return len(l)
}

func (l LocationList) Less(i, j int) bool {
	return l[i].Id < l[j].Id
}

func (l LocationList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Location struct {
	Id      int    `json:"id"`                // identity of the location
	Name    string `json:"name,omitempty"`    // name of the location
	Terrain string `json:"terrain,omitempty"` // type of terrain in location
	Is      struct {
		Hidden       bool `json:"hidden,omitempty"`        // was `LO > hi`
		MajorCity    bool `json:"major-city,omitempty"`    // was `SL > mc`
		SafeHaven    bool `json:"safe-haven,omitempty"`    // was `SL > sh`
		SeaLane      bool `json:"sea-lane,omitempty"`      // was `LO > sl`
		SummerBridge bool `json:"summer-bridge,omitempty"` // was `SL > sf`
		Uldim        bool `json:"uldim,omitempty"`         // was `SL > uf`
	} `json:"is,omitempty"`
	Dest struct { // was `LO > pd`
		North int `json:"north,omitempty"`
		East  int `json:"east,omitempty"`
		South int `json:"south,omitempty"`
		West  int `json:"west,omitempty"`
	} `json:"dest,omitempty"`
	Inside struct {
		Where        int    `json:"where,omitempty"` // was `LI > wh`
		Gates        ints_l `json:"gates,omitempty"` // here list of inside locations
		Roads        ints_l `json:"roads,omitempty"` // here list of inside locations
		SubLocations ints_l `json:"sub-locations,omitempty"`
	} `json:"inside,omitempty"`
}

// MapLocationsFromMapGen loads map locations from the globals created by the map generator.
func MapLocationsFromMapGen() (locations LocationList) {
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] == nil { // hole in map?
				continue
			}

			loc := &Location{
				Id:      map_[row][col].region,
				Name:    map_[row][col].name,
				Terrain: terrainStr[map_[row][col].terrain],
			}

			if loc.Name == "Unnamed" {
				loc.Name = ""
			}

			loc.Is.Hidden = map_[row][col].hidden != FALSE
			loc.Is.SafeHaven = map_[row][col].safe_haven != FALSE
			loc.Is.SeaLane = map_[row][col].sea_lane != FALSE // untested
			loc.Is.SummerBridge = map_[row][col].summerbridge_flag != FALSE
			loc.Is.Uldim = map_[row][col].uldim_flag != FALSE

			if map_[row][col].inside != 0 {
				loc.Inside.Where = map_[row][col].inside + REGION_OFF
			}

			for _, gateId := range map_[row][col].gates_num {
				loc.Inside.Gates = append(loc.Inside.Gates, gateId)
			}
			sort.Ints(loc.Inside.Gates)
			for _, road := range map_[row][col].roads {
				if road == nil {
					continue
				}
				loc.Inside.Roads = append(loc.Inside.Roads, road.ent_num)
			}
			sort.Ints(loc.Inside.Roads)
			for _, locId := range map_[row][col].subs {
				loc.Inside.SubLocations = append(loc.Inside.SubLocations, locId)
			}
			sort.Ints(loc.Inside.SubLocations)

			loc.Dest.North = prov_dest(map_[row][col], MG_DIR_N)
			loc.Dest.East = prov_dest(map_[row][col], MG_DIR_E)
			loc.Dest.South = prov_dest(map_[row][col], MG_DIR_S)
			loc.Dest.West = prov_dest(map_[row][col], MG_DIR_W)

			locations = append(locations, loc)
		}
	}

	sort.Sort(locations)
	return locations
}

// SubLocationsFromMapGen loads sub-locations from the globals created by the map generator.
func SubLocationsFromMapGen() (locations LocationList) {
	for _, smg := range subloc_mg {
		if smg == nil {
			continue
		}

		loc := &Location{
			Id:      smg.region,
			Name:    smg.name,
			Terrain: terrainStr[smg.terrain],
		}
		if loc.Name == "Unnamed" {
			loc.Name = ""
		}

		if smg.inside == 0 {
			panic("assert(subloc[i].inside != 0)")
		}

		loc.Inside.Where = smg.inside

		loc.Is.Hidden = smg.hidden != FALSE
		loc.Is.MajorCity = smg.major_city != FALSE
		loc.Is.SafeHaven = smg.safe_haven != FALSE

		for _, gateId := range smg.gates_num {
			loc.Inside.Gates = append(loc.Inside.Gates, gateId)
		}
		sort.Ints(loc.Inside.Gates)
		for _, road := range smg.roads {
			if road == nil {
				continue
			}
			loc.Inside.Roads = append(loc.Inside.Roads, road.ent_num)
		}
		sort.Ints(loc.Inside.Roads)
		for _, locId := range smg.subs {
			loc.Inside.SubLocations = append(loc.Inside.SubLocations, locId)
		}
		sort.Ints(loc.Inside.SubLocations)

		locations = append(locations, loc)
	}

	sort.Sort(locations)
	return locations
}
