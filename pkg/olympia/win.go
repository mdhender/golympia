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

import "log"

const (
	MIN_TURNS = 12
)

// has one or the other nation won?
func check_nation_win() bool {
	// this is possibly redundant.
	calculate_nation_nps()

	// minimum # of turns.
	if sysclock.turn < MIN_TURNS {
		return false
	}

	for _, k := range loop_nation() {
		// ignore neutral nations
		if p_nation(k).neutral {
			continue
		}

		flag := 0
		for _, i := range loop_city() {
			ruler := player_controls_loc(i)
			if ruler != 0 && nation(ruler) != 0 && nation(ruler) != k {
				flag = 1
				break
			}
		}
		if flag != 0 {
			continue
		}

		for _, i := range loop_castle() {
			ruler := player_controls_loc(i)
			if ruler != 0 && nation(ruler) != 0 && nation(ruler) != k {
				flag = 1
				break
			}
		}
		if flag != 0 {
			continue
		}

		total := 0
		for _, j := range loop_nation() {
			if p_nation(j).neutral {
				continue
			}
			if k != j {
				total += rp_nation(j).nps
			}
		}

		// you haven't met the win conditions this turn, so we should zero you out.
		if total*2 >= rp_nation(k).nps {
			rp_nation(k).win = 0
			continue
		}

		// add another turn...
		rp_nation(k).win++

		// two nations cannot win simultaneously because of the last condition.
		if rp_nation(k).win == 2 {
			// we have a winner!
			// this really needs to go into the front of the Times, but how?
			return true
		}
	}

	return false
}

// is the game over?
func check_win_conditions() {
	// nation win
	if check_nation_win() {
		// do something here.
		log.Printf("Nation win!\n")
	}
}
