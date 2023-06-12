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
	"sort"
)

const (
	item_gold = 1

	item_peasant         = 10
	item_worker          = 11
	item_soldier         = 12
	item_archer          = 13
	item_knight          = 14
	item_elite_guard     = 15
	item_pikeman         = 16
	item_blessed_soldier = 17
	item_ghost_warrior   = 18
	item_sailor          = 19
	item_swordsman       = 20
	item_crossbowman     = 21
	item_elite_arch      = 22
	item_angry_peasant   = 23
	item_pirate          = 24
	item_elf             = 25
	item_spirit          = 26
	item_postulant       = 27
	item_fanatic         = 28
	item_ninja           = 29
	item_angel           = 30

	item_corpse       = 31
	item_savage       = 32
	item_skeleton     = 33
	item_barbarian    = 34
	item_wagon        = 35
	item_skirmisher   = 36
	item_hvy_foot     = 37
	item_hvy_xbowman  = 38
	item_elvish_arrow = 39
	item_hvy_xbow     = 40
	item_horse_archer = 41
	item_pctg_token   = 42
	item_cavalier     = 43
	item_new_wagon    = 44
	item_hvy_wagon    = 45
	item_war_wagon    = 46

	item_wild_horse   = 51
	item_riding_horse = 52
	item_warmount     = 53
	item_pegasus      = 54
	item_nazgul       = 55

	item_flotsam         = 59
	item_battering_ram   = 60
	item_catapult        = 61
	item_siege_tower     = 62
	item_ratspider_venom = 63
	item_lana_bark       = 64
	item_avinia_leaf     = 65
	item_spiny_root      = 66
	item_farrenstone     = 67
	item_yew             = 68
	item_elfstone        = 69
	item_mallorn_wood    = 70
	item_pretus_bones    = 71
	item_longbow         = 72
	item_plate           = 73
	item_longsword       = 74
	item_pike            = 75
	item_ox              = 76
	item_lumber          = 77
	item_stone           = 78
	item_iron            = 79
	item_leather         = 80
	item_ratspider       = 81
	item_mithril         = 82
	item_gate_crystal    = 83
	item_blank_scroll    = 84
	item_crossbow        = 85
	item_rug             = 86
	item_fish            = 87
	item_pepper          = 88
	item_pipeweed        = 89
	item_ale             = 90
	item_wool            = 91
	item_jewel           = 92
	item_opium           = 93
	item_basket          = 94 /* woven basket */
	item_pot             = 95 /* clay pot */
	item_tax_cookie      = 96
	item_fish_oil        = 97
	item_drum            = 98
	item_hide            = 99
	item_mob_cookie      = 101
	item_lead            = 102
	item_fine_cloak      = 103
	item_chocolate       = 104
	item_ivory           = 105
	item_cardamom        = 106
	item_honey           = 107
	item_ink             = 108
	item_licorice        = 109
	item_soap            = 110
	item_old_book        = 111
	item_jade_idol       = 112
	item_purple_cloth    = 113
	item_rose_perfume    = 114
	item_silk            = 115
	item_incense         = 116
	item_ochre           = 117
	item_jeweled_egg     = 118
	item_obsidian        = 119
	item_orange          = 251
	item_cinnabar        = 252
	item_myrhh           = 253
	item_saffron         = 254
	item_dried_fish      = 255
	item_tallow          = 256
	item_candles         = 257
	item_wax             = 258
	item_sugar           = 259
	item_salt            = 260
	item_glue            = 261
	item_linen           = 262
	item_beans           = 263
	item_walnuts         = 264
	item_flax            = 265
	item_flutes          = 266
	item_cassava         = 267
	item_plum_wine       = 268
	item_vinegar         = 269
	item_tea             = 270
	item_centaur         = 271
	item_minotaur        = 272
	item_undead_cookie   = 273
	item_fog_cookie      = 274
	item_wind_cookie     = 275
	item_rain_cookie     = 276
	item_mage_menial     = 277 /* mage menial labor cookie */
	item_spider          = 278 /* giant spider */
	item_rat             = 279 /* horde of rats */
	item_lion            = 280
	item_bird            = 281 /* giant bird */
	item_lizard          = 282
	item_bandit          = 283
	item_chimera         = 284
	item_harpie          = 285
	item_dragon          = 286
	item_orc             = 287
	item_gorgon          = 288
	item_wolf            = 289
	item_orb             = 290
	item_cyclops         = 291
	item_giant           = 292
	item_faery           = 293
	item_petty_thief     = 294
	item_seagrass        = 295
	item_firewort        = 296
	item_beastnip        = 297
	item_elf_poppy       = 298
	item_ironwood        = 299
	item_kings_fern      = 300
	item_moon_palms      = 301
	item_otter           = 302
	item_mole            = 303
	item_bull            = 304
	item_eagle           = 305
	item_monkey          = 306
	item_hare            = 307
	item_wardog          = 308
	item_sand_rat        = 309
	item_balrog          = 310
	item_dirt_golem      = 311
	item_flesh_golem     = 312
	item_iron_golem      = 313
	item_lesser_demon    = 314
	item_greater_demon   = 315
	item_green_rose      = 316
	item_elf_ear         = 317
	item_savage_ear      = 318
	item_nazgul_tail     = 319
	item_centaur_hide    = 320
	item_minotaur_hide   = 321
	item_spider_eye      = 322
	item_rat_tail        = 323
	item_lion_mane       = 324
	item_bird_feather    = 325
	item_lizard_tail     = 326
	item_bandit_ear      = 327
	item_chimera_eye     = 328
	item_harpie_feather  = 329
	item_dragon_scale    = 330
	item_orc_scalp       = 331
	item_gorgon_liver    = 332
	item_wolf_hide       = 333
	item_cyclops_eye     = 334
	item_giant_tongue    = 335
	item_balrog_horn     = 336
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

type EntityItemList []*EntityItem

func (l EntityItemList) Len() int {
	return len(l)
}

func (l EntityItemList) Less(i, j int) bool {
	return l[i].Id < l[j].Id
}

func (l EntityItemList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type EntityItem struct {
	Id            int             `json:"id"`                      // identity of the item
	Name          string          `json:"name,omitempty"`          // name of the item
	PluralName    string          `json:"plural-name,omitempty"`   //
	Kind          string          `json:"kind,omitempty"`          //
	SubKind       string          `json:"sub-kind,omitempty"`      //
	Animal        bool            `json:"animal,omitempty"`        // unit is or contains a horse or an ox
	AnimalPart    int             `json:"animal_part,omitempty"`   // Produces this when killed.
	Attack        int             `json:"attack,omitempty"`        // fighter attack rating
	BasePrice     int             `json:"base_price,omitempty"`    // base price of item for market seeding
	Capturable    bool            `json:"capturable,omitempty"`    // ni-char contents are capturable
	Defense       int             `json:"defense,omitempty"`       // fighter defense rating
	FlyCap        int             `json:"fly-cap,omitempty"`       //
	IsManItem     int             `json:"is_man_item,omitempty"`   // unit is a character like thing
	LandCap       int             `json:"land-cap,omitempty"`      //
	Maintenance   int             `json:"maintenance,omitempty"`   // Maintenance cost
	Missile       int             `json:"missle,omitempty"`        // capable of missile attacks?
	NpcSplit      int             `json:"npc_split,omitempty"`     // Size to "split" at...
	Prominent     bool            `json:"prominent,omitempty"`     // big things that everyone sees
	RideCap       int             `json:"ride-cap,omitempty"`      //
	TradeGood     bool            `json:"trade_good,omitempty"`    // Is this thing a trade good? & how much
	Ungiveable    bool            `json:"ungiveable,omitempty"`    //  Can't be transferred between nobles.
	Weight        int             `json:"weight,omitempty"`        //
	WhoHas        int             `json:"who_has,omitempty"`       // who has this unique item
	Wild          bool            `json:"wild,omitempty"`          // appears in the wild as a random encounter. (value is actually the NPC_prog.)
	XItemArtifact *EntityArtifact `json:"item-artifact,omitempty"` // eventually will replace XItemMagic
	XItemMagic    *ItemMagic      `json:"item-magic,omitempty"`    // will be replaced by XItemArtifact
}

// EntityItemDataLoad loads items from a JSON file and converts it
// to the in-memory data store (a/k/a, global boxes).
func EntityItemDataLoad(name string, scanOnly bool) (EntityItemList, error) {
	log.Printf("EntityItemDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("EntityItemDataLoad: %w", err)
	}
	var list EntityItemList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("EntityItemDataLoad: %w", err)
	}
	if scanOnly {
		return nil, nil
	}
	for _, e := range list {
		if e.Kind == "" {
			e.Kind = "item"
		}
		BoxAlloc(e.Id, strKind[e.Kind], strSubKind[e.SubKind])
		bx[e.Id].x_item = e.toBox()
	}
	return nil, nil
}

// EntityItemDataSave scans the in-memory data store (a/k/a, global boxes),
// converts to the JSON model, and saves the data to a file.
func EntityItemDataSave(name string) error {
	var list EntityItemList
	for id, box := range bx {
		if box == nil || box.x_item == nil {
			continue
		}
		list = append(list, box.ToEntityItem(id))
	}
	sort.Sort(list)
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("EntityItemDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("EntityItemDataSave: %w", err)
	}
	log.Printf("EntityItemDataSave: wrote %d/%d to %s\n", len(list), len(data), name)
	return nil
}

func (e *EntityItem) toBox() *entity_item {
	ei := &entity_item{
		id:          e.Id,
		animal_part: e.AnimalPart,
		attack:      e.Attack,
		base_price:  e.BasePrice,
		defense:     e.Defense,
		fly_cap:     e.FlyCap,
		is_man_item: e.IsManItem,
		land_cap:    e.LandCap,
		maintenance: e.Maintenance,
		missile:     e.Missile,
		npc_split:   e.NpcSplit,
		plural_name: e.PluralName,
		ride_cap:    e.RideCap,
		weight:      e.Weight,
		who_has:     e.WhoHas,
	}
	if e.Animal {
		ei.animal = TRUE
	}
	if e.Prominent {
		ei.prominent = TRUE
	}
	if e.Capturable {
		ei.capturable = TRUE
	}
	if e.Ungiveable {
		ei.ungiveable = TRUE
	}
	if e.Wild {
		ei.wild = TRUE
	}
	if e.TradeGood {
		ei.trade_good = TRUE
	}
	return ei
}

func (l item_ent_l) ToInventoryList() (il InventoryList) {
	for _, e := range l {
		if !valid_box(e.item) || e.qty <= 0 {
			continue
		}
		il = append(il, Inventory{Id: e.item, Qty: e.qty})
	}
	return il
}
