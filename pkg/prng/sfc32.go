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

// sfc32_state holds the state for our PRNG.
type sfc32_state struct {
	a uint32
	b uint32
	c uint32
	d uint32
}

var (
	state        *sfc32_state
	privateState sfc32_state
	defaultSeed  = [4]uint32{0, 12345, 0, 1}
)

// sfc32_init returns an initialized PRNG.
// it is equivalent to calling sfc32_seed(0, a, b, c, d).
func sfc32_init(a, b, c, d uint32) *sfc32_state {
	state := &sfc32_state{}
	return state.seed(a, b, c, d)
}

// seed seeds a PRNG.
func (state *sfc32_state) seed(a, b, c, d uint32) *sfc32_state {
	state.a = a
	state.b = b
	state.c = c
	state.d = d

	// source recommends running 12 iterations before using state
	for i := 0; i < 12; i++ {
		state.next()
	}
	return state
}

// next returns the next value from the PRNG.
func (state *sfc32_state) next() uint32 {
	var t uint32
	state.a |= 0
	state.b |= 0
	state.c |= 0
	state.d |= 0
	t = (state.a + state.b | 0) + state.d | 0
	state.d = state.d + 1 | 0
	state.a = state.b ^ state.b>>9
	state.b = state.c + (state.c << 3) | 0
	state.c = state.c<<21 | state.c>>11
	state.c = state.c + t | 0
	return (t >> 0)
}
