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

type ShipList []*Ship

func (l ShipList) Len() int {
	return len(l)
}

func (l ShipList) Less(i, j int) bool {
	return l[i].Id < l[j].Id
}

func (l ShipList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Ship struct {
	Id      int    `json:"id"`             // identity of the ship
	Name    string `json:"name,omitempty"` // name of the ship
	Kind    string `json:"kind,omitempty"`
	SubKind string `json:"sub-kind,omitempty"`
}

func ShipDataLoad(name string, scanOnly bool) (ShipList, error) {
	log.Printf("ShipDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("ShipDataLoad: %w", err)
	}
	var list ShipList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("ShipDataLoad: %w", err)
	}
	if scanOnly {
		return nil, nil
	}
	for _, e := range list {
		BoxAlloc(e.Id, strKind[e.Kind], strSubKind[e.SubKind])
	}
	return nil, nil
}

func ShipDataSave(name string) error {
	list := ShipList{}
	sort.Sort(list)
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("ShipDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("ShipDataSave: %w", err)
	}
	return nil
}

type EntityShip struct {
	Forts     int `json:"forts,omitempty"`      //
	GalleyRam int `json:"galley-ram,omitempty"` // galley is fitted with a ram
	Hulls     int `json:"hulls,omitempty"`      // Various ship parts
	Keels     int `json:"keels,omitempty"`      //
	Ports     int `json:"ports,omitempty"`      //
	Sails     int `json:"sails,omitempty"`      //
}

func (e *entity_ship) IsZero() bool {
	return e == nil || (e.forts == 0 && e.galley_ram == 0 && e.hulls == 0 && e.keels == 0 && e.ports == 0 && e.sails == 0)
}

func (e *entity_ship) ToEntityShip() *EntityShip {
	if e.IsZero() {
		return nil
	}
	return &EntityShip{
		Forts:     e.forts,
		GalleyRam: e.galley_ram,
		Hulls:     e.hulls,
		Keels:     e.keels,
		Ports:     e.ports,
		Sails:     e.sails,
	}
}
