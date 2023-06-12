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
	"os"
)

type Option func() error

var (
	cityDataFilename      string
	continentDataFilename string
	gateDataFilename      string
	landDataFilename      string
	locationDataFilename  string
	mapDataFilename       string
	regionDataFilename    string
	roadDataFilename      string
	seedDataFilename      string
)

func WithCityData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("city data: missing file name")
		} else if fi, err := os.Stat(name); err != nil {
			return fmt.Errorf("city data: %w", err)
		} else if !fi.Mode().IsRegular() {
			return fmt.Errorf("city data: %w", fmt.Errorf("not a file"))
		}
		cityDataFilename = name
		return nil
	}
}

func WithContinentData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("continent data: missing file name")
			//} else if fi, err := os.Stat(name); err != nil {
			//	return fmt.Errorf("continent data: %w", err)
			//} else if !fi.Mode().IsRegular() {
			//	return fmt.Errorf("continent data: %w", fmt.Errorf("not a file"))
		}
		continentDataFilename = name
		return nil
	}
}

func WithGateData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("gate data: missing file name")
			//} else if fi, err := os.Stat(name); err != nil {
			//	return fmt.Errorf("gate data: %w", err)
			//} else if !fi.Mode().IsRegular() {
			//	return fmt.Errorf("gate data: %w", fmt.Errorf("not a file"))
		}
		gateDataFilename = name
		return nil
	}
}

func WithLandData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("land data: missing file name")
			//} else if fi, err := os.Stat(name); err != nil {
			//	return fmt.Errorf("land data: %w", err)
			//} else if !fi.Mode().IsRegular() {
			//	return fmt.Errorf("land data: %w", fmt.Errorf("not a file"))
		}
		landDataFilename = name
		return nil
	}
}

func WithLibPath(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("lib path: missing path name")
		} else if fi, err := os.Stat(name); err != nil {
			return fmt.Errorf("lib path: %w", err)
		} else if !fi.Mode().IsDir() {
			return fmt.Errorf("lib path: %w", fmt.Errorf("not a directory"))
		}
		locationDataFilename = name
		return nil
	}
}

func WithLocationData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("location data: missing file name")
			//} else if fi, err := os.Stat(name); err != nil {
			//	return fmt.Errorf("location data: %w", err)
			//} else if !fi.Mode().IsRegular() {
			//	return fmt.Errorf("location data: %w", fmt.Errorf("not a file"))
		}
		locationDataFilename = name
		return nil
	}
}

func WithMapData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("map data: missing file name")
		} else if fi, err := os.Stat(name); err != nil {
			return fmt.Errorf("map data: %w", err)
		} else if !fi.Mode().IsRegular() {
			return fmt.Errorf("map data: %w", fmt.Errorf("not a file"))
		}
		mapDataFilename = name
		return nil
	}
}

func WithRegionData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("region data: missing file name")
		} else if fi, err := os.Stat(name); err != nil {
			return fmt.Errorf("region data: %w", err)
		} else if !fi.Mode().IsRegular() {
			return fmt.Errorf("region data: %w", fmt.Errorf("not a file"))
		}
		regionDataFilename = name
		return nil
	}
}

func WithRoadData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("road data: missing file name")
			//} else if fi, err := os.Stat(name); err != nil {
			//	return fmt.Errorf("road data: %w", err)
			//} else if !fi.Mode().IsRegular() {
			//	return fmt.Errorf("road data: %w", fmt.Errorf("not a file"))
		}
		roadDataFilename = name
		return nil
	}
}

func WithSeedData(name string) func() error {
	return func() error {
		if name == "" {
			return fmt.Errorf("seed data: missing file name")
			//} else if fi, err := os.Stat(name); err != nil {
			//	return fmt.Errorf("seed data: %w", err)
			//} else if !fi.Mode().IsRegular() {
			//	return fmt.Errorf("seed data: %w", fmt.Errorf("not a file"))
		}
		seedDataFilename = name
		return nil
	}
}
