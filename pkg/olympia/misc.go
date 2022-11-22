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

type MiscList []*Misc
type Misc struct {
	Id   int    `json:"id"`             // identity of the thing
	Name string `json:"name,omitempty"` // name of the thing
}

func MiscDataLoad(name string) (MiscList, error) {
	log.Printf("MiscDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("MiscDataLoad: %w", err)
	}
	var js MiscList
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, fmt.Errorf("MiscDataLoad: %w", err)
	}
	return nil, nil
}

func MiscDataSave(name string) error {
	var js MiscList
	data, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return fmt.Errorf("MiscDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("MiscDataSave: %w", err)
	}
	return nil
}
