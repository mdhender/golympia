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

package prng

import (
	"encoding/json"
	"os"
)

// Range returns a pseudo-random value in the range of low...high
func Range(low, high int) int {
	if state == nil {
		state = &privateState
		state.seed(defaultSeed[0], defaultSeed[1], defaultSeed[2], defaultSeed[3])
	}
	if high < low {
		high, low = low, high
	}
	n := (int)(state.next())
	if n < 0 {
		n = -n
	}
	return low + (n % (high - low + 1))
}

func LoadSeed(name string) {
	var data struct {
		State struct {
			A, B, C, D uint32
		} `json:"state"`
	}

	state = &privateState
	if buf, err := os.ReadFile(name); err == nil {
		if err = json.Unmarshal(buf, &data); err == nil {
			state = state.seed(data.State.A, data.State.B, data.State.C, data.State.C)
		} else {
			state = state.seed(defaultSeed[0], defaultSeed[1], defaultSeed[2], defaultSeed[3])
		}
	} else {
		state = state.seed(defaultSeed[0], defaultSeed[1], defaultSeed[2], defaultSeed[3])
	}
}

func SaveSeed(name string) {
	if state == nil {
		state = &privateState
		state = state.seed(defaultSeed[0], defaultSeed[1], defaultSeed[2], defaultSeed[3])
	}
	var data struct {
		State struct {
			A, B, C, D uint32
		} `json:"state"`
	}
	data.State.A = state.a
	data.State.B = state.b
	data.State.C = state.c
	data.State.D = state.d

	if buf, err := json.MarshalIndent(data, "", "  "); err != nil {
		panic(err)
	} else if err = os.WriteFile(name, buf, 0666); err != nil {
		panic(err)
	}
}
