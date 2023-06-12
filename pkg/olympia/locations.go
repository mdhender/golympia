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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type LocationInfo struct {
	Where    int    `json:"where,omitempty"`
	HereList ints_l `json:"hereList,omitempty"`
}

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
	Id      int    `json:"id"`             // identity of the location
	Name    string `json:"name,omitempty"` // name of the location
	Kind    string `json:"kind,omitempty"`
	SubKind string `json:"sub-kind,omitempty"`
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

func LocationDataLoad(name string, scanOnly bool) (LocationList, error) {
	log.Printf("LocationDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("LocationDataLoad: %w", err)
	}
	var list LocationList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("LocationDataLoad: %w", err)
	}
	if scanOnly {
		return nil, nil
	}
	for _, e := range list {
		BoxAlloc(e.Id, strKind[e.Kind], strSubKind[e.SubKind])
	}
	return nil, nil
}

func LocationDataSave(name string) error {
	list := LocationList{}
	sort.Sort(list)
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("LocationDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("LocationDataSave: %w", err)
	}
	return nil
}

func (b *box) ToLocationInfo() *LocationInfo {
	if b == nil || b.x_loc_info.IsZero() {
		return nil
	}
	li := &LocationInfo{Where: b.x_loc_info.where}
	if len(b.x_loc_info.here_list) != 0 {
		li.HereList = append(li.HereList, b.x_loc_info.here_list...)
	}
	return li
}

func (li loc_info) IsZero() bool {
	return li.where == 0 && len(li.here_list) == 0
}

type EntityLoc struct {
	Barrier       int            `json:"barrier,omitempty"`         // magical barrier
	Control       *LocControlEnt `json:"control,omitempty"`         //
	DistFromGate  int            `json:"dist-from-gate,omitempty"`  //
	DistFromSea   int            `json:"dist-from-sea,omitempty"`   // provinces to sea province
	DistFromSwamp int            `json:"dist-from-swamp,omitempty"` //
	Hidden        int            `json:"hidden,omitempty"`          // is location hidden?
	MineInfo      *EntityMine    `json:"mine-info,omitempty"`       // If there's a mine.
	NearGrave     int            `json:"near-grave,omitempty"`      // nearest graveyard
	ProvDest      ints_l         `json:"prov-dest,omitempty"`       // cached exits for flood fills
	Recruited     int            `json:"recruited,omitempty"`       // How many recruited this month
	SeaLane       int            `json:"sea-lane,omitempty"`        // fast ocean travel here, also "tracks" for npc ferries
	Shroud        int            `json:"shroud,omitempty"`          // magical scry shroud
	TaxRate       int            `json:"tax-rate,omitempty"`        // Tax rate for this loc.

	// Effects EffectList        // list of effects on location // not used?

	// location control -- need two so that we can only change fees at the end of the month.
	control2 *LocControlEnt // doesn't need to be saved
}

func (e *entity_loc) ToEntityLoc() *EntityLoc {
	if e == nil {
		return nil
	}
	el := &EntityLoc{
		Barrier:       e.barrier,
		Control:       e.control.ToLocControlEnt(),
		DistFromGate:  e.dist_from_gate,
		DistFromSea:   e.dist_from_sea,
		DistFromSwamp: e.dist_from_swamp,
		Hidden:        e.hidden,
		MineInfo:      e.mine_info.ToEntityMine(),
		NearGrave:     e.near_grave,
		ProvDest:      e.prov_dest,
		Recruited:     e.recruited,
		SeaLane:       e.sea_lane,
		Shroud:        e.shroud,
		TaxRate:       e.tax_rate,
		control2:      e.control2.ToLocControlEnt(),
	}
	return el
}

// LocControlEnt is fees plus open/closed
type LocControlEnt struct {
	Open   bool `json:"open,omitempty"`
	Men    int  `json:"men,omitempty"`    // fee per person
	Nobles int  `json:"nobles,omitempty"` // fee per noble
	Weight int  `json:"weight,omitempty"` // fee per measure of weight
}

func (e *loc_control_ent) IsZero() bool {
	return e == nil || (e.closed && e.men == 0 && e.nobles == 0 && e.weight == 0)
}

func (e *loc_control_ent) ToLocControlEnt() *LocControlEnt {
	if e == nil {
		return nil
	}
	return &LocControlEnt{
		Open:   !e.closed,
		Men:    e.men,
		Nobles: e.nobles,
		Weight: e.weight,
	}
}

type EntitySubLoc struct {
	BoundStorms  ints_l          `json:"bound-storms,omitempty"`  // storms bound to this ship
	Builds       EntityBuildList `json:"builds,omitempty"`        // Ongoing builds here.
	Damage       int             `json:"damage,omitempty"`        // 0=none, hp=fully destroyed
	Defense      int             `json:"defense,omitempty"`       // defense rating of structure
	EntranceSize int             `json:"entrance-size,omitempty"` // size of entrance to subloc
	Guild        int             `json:"guild,omitempty"`         // what skill, if a sub_guild
	Hp           int             `json:"hp,omitempty"`            // "hit points"
	LinkFrom     ints_l          `json:"link-from,omitempty"`     // where we are linked from
	LinkTo       ints_l          `json:"link-to,omitempty"`       // where we are linked to
	Loot         int             `json:"loot,omitempty"`          // loot & pillage level
	Major        int             `json:"major,omitempty"`         // major city
	Moat         int             `json:"moat,omitempty"`          // Has a moat?
	Moving       int             `json:"moving,omitempty"`        // daystamp of beginning of movement
	NearCities   ints_l          `json:"near-cities,omitempty"`   // cities rumored to be nearby
	OpiumEcon    int             `json:"opium-econ,omitempty"`    // addiction level of city
	Prominence   int             `json:"prominence,omitempty"`    // prominence of city
	Safe         bool            `json:"safe,omitempty"`          // safe haven
	TaxMarket    int             `json:"tax-market,omitempty"`    // Market tax rate.
	TaxMarket2   int             `json:"tax-market-2,omitempty"`  // Temporary until end of month
	Teaches      ints_l          `json:"teaches,omitempty"`       // skills location offers
	XShip        *EntityShip     `json:"x-ship,omitempty"`        // Maybe a ship?

	// location control -- either here or loc
	Control  *LocControlEnt `json:"control,omitempty"`   //
	Control2 *LocControlEnt `json:"control-2,omitempty"` //

	//short shaft_depth;		/* depth of mine shaft */
	//int capacity;			/* capacity of ship */
	//schar galley_ram;		/* galley is fitted with a ram */
	//schar link_when;		/* month link is open, -1 = never */
	//schar link_open;		/* link is open now */
	//struct effect **effects;        /* ilist of effects on sub-location */
}

func (e *entity_subloc) ToEntitySubLoc() *EntitySubLoc {
	if e == nil {
		return nil
	}
	return &EntitySubLoc{
		BoundStorms:  e.bound_storms.ToList(),
		Builds:       e.builds.ToEntityBuildList(),
		Damage:       e.damage,
		Defense:      e.defense,
		EntranceSize: e.entrance_size,
		Guild:        e.guild,
		Hp:           e.hp,
		LinkFrom:     e.link_from.ToList(),
		LinkTo:       e.link_to.ToList(),
		Loot:         e.loot,
		Major:        e.major,
		Moat:         e.moat,
		Moving:       e.moving,
		NearCities:   e.near_cities.ToList(),
		OpiumEcon:    e.opium_econ,
		Prominence:   e.prominence,
		Safe:         e.safe,
		TaxMarket:    e.tax_market,
		TaxMarket2:   e.tax_market2,
		Teaches:      e.teaches.ToList(),
		XShip:        e.x_ship.ToEntityShip(),
		Control:      e.control.ToLocControlEnt(),
		Control2:     e.control2.ToLocControlEnt(),
	}
}
