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

package cli

import (
	"fmt"
	"github.com/mdhender/golympia/pkg/olympia"
	"github.com/spf13/cobra"
	"log"
)

// cmdGenerateMap runs the map generator command
var cmdGenerateMap = &cobra.Command{
	Use:   "map",
	Short: "generate a new map",
	RunE: func(cmd *cobra.Command, args []string) error {
		if argsRoot.libdir == "" {
			return fmt.Errorf("missing lib-dir parameter")
		} else if argsGenerateMap.mapFileName == "" {
			return fmt.Errorf("missing map-data parameter")
		} else if argsGenerateMap.continentFileName == "" {
			return fmt.Errorf("missing continent-data parameter")
		} else if argsGenerateMap.landFileName == "" {
			return fmt.Errorf("missing land-data parameter")
		} else if argsGenerateMap.locationFileName == "" {
			return fmt.Errorf("missing location-data parameter")
		} else if argsGenerateMap.regionFileName == "" {
			return fmt.Errorf("missing region-data parameter")
		} else if argsGenerateMap.seedFileName == "" {
			return fmt.Errorf("missing seed-data parameter")
		}

		var options []olympia.Option
		options = append(options, olympia.WithLibPath(argsRoot.libdir))
		options = append(options, olympia.WithMapData(argsGenerateMap.mapFileName))
		options = append(options, olympia.WithCityData(argsGenerateMap.cityFileName))
		options = append(options, olympia.WithContinentData(argsGenerateMap.continentFileName))
		options = append(options, olympia.WithGateData(argsGenerateMap.gateFileName))
		options = append(options, olympia.WithLandData(argsGenerateMap.landFileName))
		options = append(options, olympia.WithLocationData(argsGenerateMap.locationFileName))
		options = append(options, olympia.WithRegionData(argsGenerateMap.regionFileName))
		options = append(options, olympia.WithRoadData(argsGenerateMap.roadFileName))
		options = append(options, olympia.WithSeedData(argsGenerateMap.seedFileName))

		if err := olympia.GenerateMap(options...); err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

var argsGenerateMap struct {
	mapFileName       string
	cityFileName      string
	continentFileName string
	gateFileName      string
	landFileName      string
	locationFileName  string
	regionFileName    string
	roadFileName      string
	seedFileName      string
}

func init() {
	cmdGenerate.AddCommand(cmdGenerateMap)
	// inputs
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.mapFileName, "map-data", "map-data.txt", "map data to import")
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.cityFileName, "city-data", "cities.json", "city name data to import")
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.landFileName, "land-data", "lands.json", "land data to import")
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.regionFileName, "region-data", "regions.json", "region data to import")

	// outputs
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.continentFileName, "continent-data", "continents.json", "continent data to export")
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.locationFileName, "location-data", "locations.json", "location data to export")
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.gateFileName, "gate-data", "gates.json", "gate data to export")
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.roadFileName, "road-data", "roads.json", "road data to export")
	cmdGenerateMap.Flags().StringVar(&argsGenerateMap.seedFileName, "seed-data", "randseed.json", "random seed data to export")

	//if err := cmdGenerateMap.MarkFlagRequired("map-data"); err != nil {
	//	panic(err)
	//}
}
