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

// monsters are used to initialize the items in a new map

type monster_t EntityItemList

var monster_tbl = EntityItemList{
	{Id: 10,
		Kind:       "item",
		Name:       "peasant",
		PluralName: "peasants",
		Attack:     1,
		Defense:    1,
		IsManItem:  TRUE,
		LandCap:    100,
		Prominent:  true,
		Weight:     100,
	},
	{Id: 55,
		Kind:       "item",
		Name:       "nazgul",
		Animal:     true,
		Attack:     80,
		Capturable: true,
		Defense:    80,
		FlyCap:     150,
		LandCap:    150,
		PluralName: "nazgul",
		Prominent:  true,
		RideCap:    150,
		Weight:     1500},
	{Id: 101,
		Kind:       "item",
		Name:       "mob cookie",
		PluralName: "mob cookies",
	},
	{Id: 273,
		Kind:       "item",
		Name:       "undead cookie",
		PluralName: "undead cookies",
	},
	{Id: 274,
		Kind:       "item",
		Name:       "fog cookie",
		PluralName: "fog cookies",
	},
	{Id: 275,
		Kind:       "item",
		Name:       "wind cookie",
		PluralName: "wind cookies",
	},
	{Id: 276,
		Kind:       "item",
		Name:       "rain cookie",
		PluralName: "rain cookies",
	},
	{Id: 277,
		Kind:       "item",
		Name:       "mage menial cookie",
		PluralName: "mage menial cookies",
	},
	{Id: 287,
		Kind:       "item",
		Name:       "orc",
		PluralName: "orcs",
		Attack:     20,
		Capturable: true,
		Defense:    15,
		FlyCap:     0,
		IsManItem:  TRUE,
		LandCap:    100,
		Prominent:  true,
		RideCap:    1,
		Weight:     100,
	},
	{Id: 294,
		Kind:       "item",
		Name:       "petty thief cookie",
		PluralName: "petty thief cookies",
	},
}
