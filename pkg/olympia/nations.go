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

type NationList []*Nation
type Nation struct {
	Id   int    `json:"id"`             // identity of the item
	Name string `json:"name,omitempty"` // name of the item
}

func NationDataLoad(name string) (NationList, error) {
	log.Printf("NationDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("NationDataLoad: %w", err)
	}
	var js NationList
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, fmt.Errorf("NationDataLoad: %w", err)
	}
	return nil, nil
}

func NationDataSave(name string) error {
	var js NationList
	data, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return fmt.Errorf("NationDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("NationDataSave: %w", err)
	}
	return nil
}
