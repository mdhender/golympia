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

const (
	terr_land        = 1
	terr_ocean       = 2
	terr_forest      = 3
	terr_swamp       = 4
	terr_mountain    = 5
	terr_plain       = 6
	terr_desert      = 7
	terr_water       = 8
	terr_island      = 9
	terr_stone_cir   = 10 /* circle of stones */
	terr_grove       = 11 /* mallorn grove */
	terr_bog         = 12
	terr_cave        = 13
	terr_city        = 14
	terr_guild       = 15
	terr_grave       = 16
	terr_ruins       = 17
	terr_battlefield = 18
	terr_ench_for    = 19 /* enchanted forest */
	terr_rocky_hill  = 20
	terr_tree_cir    = 21
	terr_pits        = 22
	terr_pasture     = 23
	terr_oasis       = 24
	terr_yew_grove   = 25
	terr_sand_pit    = 26
	terr_sac_grove   = 27 /* sacred grove */
	terr_pop_field   = 28 /* poppy field */
	terr_temple      = 29
	terr_lair        = 30 /* dragon lair */
)

var terrainStr = map[int]string{
	0:                "",
	terr_land:        "land",
	terr_ocean:       "ocean",
	terr_forest:      "forest",
	terr_swamp:       "swamp",
	terr_mountain:    "mountain",
	terr_plain:       "plain",
	terr_desert:      "desert",
	terr_water:       "water",
	terr_island:      "island",
	terr_stone_cir:   "ring of stones",
	terr_grove:       "mallorn grove",
	terr_bog:         "bog",
	terr_cave:        "cave",
	terr_city:        "city",
	terr_guild:       "guild",
	terr_grave:       "graveyard",
	terr_ruins:       "ruins",
	terr_battlefield: "field",
	terr_ench_for:    "enchanted forest",
	terr_rocky_hill:  "rocky hill",
	terr_tree_cir:    "circle of trees",
	terr_pits:        "pits",
	terr_pasture:     "pasture",
	terr_oasis:       "oasis",
	terr_yew_grove:   "yew grove",
	terr_sand_pit:    "sand pit",
	terr_sac_grove:   "sacred grove",
	terr_pop_field:   "poppy field",
	terr_temple:      "temple",
	terr_lair:        "lair",
}

var strTerrain = map[string]int{
	"":                 0,
	"<null>":           0,
	"land":             terr_land,
	"ocean":            terr_ocean,
	"forest":           terr_forest,
	"swamp":            terr_swamp,
	"mountain":         terr_mountain,
	"plain":            terr_plain,
	"desert":           terr_desert,
	"water":            terr_water,
	"island":           terr_island,
	"ring of stones":   terr_stone_cir,
	"mallorn grove":    terr_grove,
	"bog":              terr_bog,
	"cave":             terr_cave,
	"city":             terr_city,
	"guild":            terr_guild,
	"graveyard":        terr_grave,
	"ruins":            terr_ruins,
	"field":            terr_battlefield,
	"enchanted forest": terr_ench_for,
	"rocky hill":       terr_rocky_hill,
	"circle of trees":  terr_tree_cir,
	"pits":             terr_pits,
	"pasture":          terr_pasture,
	"oasis":            terr_oasis,
	"yew grove":        terr_yew_grove,
	"sand pit":         terr_sand_pit,
	"sacred grove":     terr_sac_grove,
	"poppy field":      terr_pop_field,
	"temple":           terr_temple,
	"lair":             terr_lair,
}
