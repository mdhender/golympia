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

type NationList []*Nation

func (l NationList) Len() int {
	return len(l)
}

func (l NationList) Less(i, j int) bool {
	return l[i].Id < l[j].Id
}

func (l NationList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Nation struct {
	Id      int    `json:"id"`             // identity of the item
	Name    string `json:"name,omitempty"` // name of the item
	Kind    string `json:"kind,omitempty"`
	SubKind string `json:"sub-kind,omitempty"`
}

func NationDataLoad(name string, scanOnly bool) (NationList, error) {
	log.Printf("NationDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("NationDataLoad: %w", err)
	}
	var list NationList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("NationDataLoad: %w", err)
	}
	if scanOnly {
		return nil, nil
	}
	for _, e := range list {
		BoxAlloc(e.Id, strKind[e.Kind], strSubKind[e.SubKind])
	}
	return nil, nil
}

func NationDataSave(name string) error {
	list := NationList{}
	sort.Sort(list)
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("NationDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("NationDataSave: %w", err)
	}
	return nil
}
