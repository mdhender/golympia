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
	"fmt"
)

// BoxAlloc replaces alloc_box()
func BoxAlloc(id, kind, skind int) {
	if bx == nil {
		panic("assert(bx != nil)")
	} else if !(bx[id] == nil) {
		panic(fmt.Sprintf("assert(bx[%d] == nil)", id))
	}
	bx[id] = &box{
		kind:  schar(kind),
		skind: schar(skind),
	}
	add_next_chain(id)
	add_sub_chain(id)
}

type Box struct {
	Id             int             `json:"id"`             // identity of the thing
	Name           string          `json:"name,omitempty"` // name of the thing
	Kind           int             `json:"kind,omitempty"`
	SubKind        int             `json:"sub-kind,omitempty"`
	Attitudes      *Attitudes      `json:"attitudes,omitempty"`
	CharMagic      *CharMagic      `json:"char-magic,omitempty"`
	Effects        EffectList      `json:"effects,omitempty"`
	EntityArtifact *EntityArtifact `json:"entity-artifact,omitempty"`
	EntityChar     *EntityChar     `json:"entity-char,omitempty"`
	EntityItem     *EntityItem     `json:"entity-item,omitempty"`
	EntityLoc      *EntityLoc      `json:"entity-loc,omitempty"`
	EntityPlayer   *EntityPlayer   `json:"entity-player,omitempty"`
	EntitySubLoc   *EntitySubLoc   `json:"entity-subloc,omitempty"`
	ItemMagic      *ItemMagic      `json:"item-magic,omitempty"`
	Items          InventoryList   `json:"items,omitempty"`
	LocationInfo   *LocationInfo   `json:"location-info,omitempty"`
	Trades         TradeList       `json:"trades,omitempty"`
}

// ToBoxList replaces boxlist_print
func (l ints_l) ToBoxList() (il ints_l) {
	for _, e := range l {
		// todo: why carve out the monster attitude?
		if !(valid_box(e) || e == MONSTER_ATT) {
			continue
		}
		il = append(il, e)
	}
	return il
}

func (b *box) ToAttitudes() *Attitudes {
	if b == nil {
		return nil
	}
	return b.x_disp.ToAttitudes()
}

func (b *box) ToCharMagic() *CharMagic {
	if b == nil || b.x_char == nil {
		return nil
	}
	return b.x_char.x_char_magic.ToCharMagic()
}

func (b *box) ToEffectList() EffectList {
	if b == nil {
		return nil
	}
	return b.effects.ToEffectList()
}

func (b *box) ToEntityArtifact(id int) *EntityArtifact {
	if b == nil || b.x_item == nil {
		return nil
	}
	return b.x_item.ToEntityArtifact(id)
}

func (b *box) ToEntityChar() *EntityChar {
	if b == nil {
		return nil
	}
	return b.x_char.ToEntityChar()
}

func (b *box) ToEntityItem(id int) *EntityItem {
	if b == nil {
		return nil
	}
	return b.x_item.ToEntityItem(id)
}

func (b *box) ToEntityLoc() *EntityLoc {
	if b == nil {
		return nil
	}
	return b.x_loc.ToEntityLoc()
}

func (b *box) ToEntityPlayer() *EntityPlayer {
	if b == nil {
		return nil
	}
	return b.x_player
}

func (b *box) ToEntitySubLoc() *EntitySubLoc {
	if b == nil {
		return nil
	}
	return b.x_subloc.ToEntitySubLoc()
}

func (b *box) ToInventoryList() InventoryList {
	if b == nil {
		return nil
	}
	return b.items.ToInventoryList()
}

func (b *box) ToItemMagic(id int) *ItemMagic {
	if b == nil || b.x_item == nil {
		return nil
	}
	return b.x_item.ToItemMagic(id)
}

func (b *box) ToTradeList() TradeList {
	if b == nil {
		return nil
	}
	return b.trades.ToTradeList()
}
