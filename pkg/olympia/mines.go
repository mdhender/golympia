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

type EntityMine struct {
	Contents []*MineContents `json:"contents,omitempty"`
	Shoring  ints_l          `json:"shoring,omitempty"`
}

type MineContents struct {
	Items InventoryList `json:"items,omitempty"` // list of items held
	// iron, gold, mithril, gate_crystals int // not used?
}

func (m *entity_mine) IsZero() bool {
	if m == nil {
		return true
	}
	for _, e := range m.mc {
		if !e.IsZero() {
			return false
		}
	}
	for _, e := range m.shoring {
		if e != 0 {
			return false
		}
	}
	return true
}

func (m *entity_mine) ToEntityMine() *EntityMine {
	if m == nil || m.IsZero() {
		return nil
	}
	em := &EntityMine{}
	for _, e := range m.mc {
		if e.IsZero() {
			continue
		}
		em.Contents = append(em.Contents, e.ToMineContents())
	}
	for _, e := range m.shoring {
		if e == 0 {
			continue
		}
		em.Shoring = append(em.Shoring, e)
	}
	if len(em.Contents) == 0 && len(em.Shoring) == 0 {
		return nil
	}
	return em
}

func (m *MineContents) IsZero() bool {
	return m == nil || len(m.Items) == 0
}

func (c *mine_contents) IsZero() bool {
	return c == nil || len(c.items) == 0
}

func (c *mine_contents) ToMineContents() *MineContents {
	if c == nil || c.IsZero() {
		return nil
	}
	m := &MineContents{}
	for _, e := range c.items {
		if e.item == 0 || e.qty <= 0 {
			continue
		}
		m.Items = append(m.Items, Inventory{Id: e.item, Qty: e.qty})
	}
	if m.IsZero() {
		return nil
	}
	return m
}
