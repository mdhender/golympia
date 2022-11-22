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
	"os"
)

type MasterData struct{}

func MasterDataLoad(name string) (*MasterData, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("MasterDataLoad: %w", err)
	}
	js := &SysData{}
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, fmt.Errorf("MasterDataLoad: %w", err)
	}

	return nil, nil
}

func MasterDataSave(name string) error {
	var js MasterData
	data, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return fmt.Errorf("MasterDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("MasterDataSave: %w", err)
	}

	return nil
}
