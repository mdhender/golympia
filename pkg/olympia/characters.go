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
	"path/filepath"
	"strings"
)

type CharacterList []*Character
type Character struct {
	Id   int    `json:"id"`             // identity of the character
	Name string `json:"name,omitempty"` // name of the character
}

func CharactersLoad() error {
	path := filepath.Join(libdir, "characters")
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("CharactersLoad: scan: %w", err)
	}

	for _, f := range files {
		jsonFile := isdigit(f.Name()[0]) && strings.HasSuffix(f.Name(), ".json")
		if !jsonFile {
			continue
		}
		//scan_boxes(filepath.Join("fact", f.Name()))
		_, _ = CharacterDataLoad(filepath.Join(path, f.Name()))
	}

	return nil
}

func CharactersSave() error {
	return fmt.Errorf("CharactersSave: not implemented")
}

func CharacterDataLoad(name string) (CharacterList, error) {
	log.Printf("CharacterDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("CharacterDataLoad: %w", err)
	}
	var js CharacterList
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, fmt.Errorf("CharacterDataLoad: %w", err)
	}
	return nil, nil
}

func CharacterDataSave(name string) error {
	var js struct{}
	data, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return fmt.Errorf("CharacterDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("CharacterDataSave: %w", err)
	}
	return nil
}
