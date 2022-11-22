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
)

type item_ent_l []*item_ent

func (ie item_ent_l) Len() int {
	return len(ie)
}

func (ie item_ent_l) Less(i, j int) bool {
	return ie[i].item < ie[j].item
}

func (ie item_ent_l) Swap(i, j int) {
	ie[i], ie[j] = ie[j], ie[i]
}

type ItemList []*Item
type Item struct {
	Id   int    `json:"id"`             // identity of the item
	Name string `json:"name,omitempty"` // name of the item
}

func ItemDataLoad(name string) (ItemList, error) {
	log.Printf("ItemDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("ItemDataLoad: %w", err)
	}
	var js ItemList
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, fmt.Errorf("ItemDataLoad: %w", err)
	}
	return nil, nil
}

func ItemDataSave(name string) error {
	var js ItemList
	data, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return fmt.Errorf("ItemDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("ItemDataSave: %w", err)
	}
	return nil
}
