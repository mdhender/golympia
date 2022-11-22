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

type GateList []*Gate

func (g GateList) Len() int {
	return len(g)
}

func (g GateList) Less(i, j int) bool {
	return g[i].Id < g[j].Id
}

func (g GateList) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

type GateLink struct {
	To  int `json:"to"`            // identity of connected destination
	Key int `json:"key,omitempty"` // identity of key required to use gate
}

type Gate struct {
	Id    int         `json:"id"`    // identity of the gate
	Where int         `json:"where"` // where this gate is located (region or location)
	Links []*GateLink `json:"links,omitempty"`
}

// GatesFromMapGen loads gates from the globals created by the map generator.
func GatesFromMapGen() (gates GateList) {
	for row := 0; row < MAX_ROW; row++ {
		for col := 0; col < MAX_COL; col++ {
			if map_[row][col] == nil {
				continue
			}
			for j := 0; j < len(map_[row][col].gates_dest); j++ {
				gates = append(gates, &Gate{
					Id:    map_[row][col].gates_num[j],
					Where: map_[row][col].region,
					Links: []*GateLink{{
						To:  map_[row][col].gates_dest[j],
						Key: map_[row][col].gates_key[j],
					}},
				})
			}
		}
	}

	for _, tile := range subloc_mg {
		if tile == nil {
			continue
		}
		for j := 0; j < len(tile.gates_num); j++ {
			gates = append(gates, &Gate{
				Id:    tile.gates_num[j],
				Where: tile.region,
				Links: []*GateLink{{
					To:  tile.gates_dest[j],
					Key: tile.gates_key[j],
				}},
			})
		}
	}

	sort.Sort(gates)
	return gates
}

func GateDataLoad(name string) (GateList, error) {
	log.Printf("GateDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("GateDataLoad: %w", err)
	}
	var js GateList
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, fmt.Errorf("GateDataLoad: %w", err)
	}
	return nil, nil
}

func GateDataSave(name string) error {
	if buf, err := json.MarshalIndent(GatesFromMapGen(), "", "  "); err != nil {
		return fmt.Errorf("GateDataSave: %w", err)
	} else if err = os.WriteFile(name, buf, 0666); err != nil {
		return fmt.Errorf("GateDataSave: %w", err)
	}
	log.Printf("GateDataSave: created %s\n", name)
	return nil
}
