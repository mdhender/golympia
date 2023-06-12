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

type EntityBuildList []*EntityBuild

// EntityBuild describes what kind of build is going on in a location.
type EntityBuild struct {
	Type           int // What work is going on?
	BuildMaterials int // fifths of materials we've used
	EffortGiven    int //
	EffortRequired int // not finished if nonzero
}

type entity_build_l []*entity_build

// describes what kind of build is going on in a location.
type entity_build struct {
	type_           int // what work is going on?
	build_materials int // fifths of materials we've used
	effort_given    int //
	effort_required int // not finished if nonzero
}

func (e *entity_build) IsZero() bool {
	return e == nil || (e.type_ == 0 && e.build_materials == 0 && e.effort_required == 0 && e.effort_given == 0)
}

func (e *entity_build) ToEntityBuild() *EntityBuild {
	if e.IsZero() {
		return nil
	}
	return &EntityBuild{
		Type:           e.type_,
		BuildMaterials: e.build_materials,
		EffortGiven:    e.effort_given,
		EffortRequired: e.effort_required,
	}
}

func (l entity_build_l) delete(index int) entity_build_l {
	var cp entity_build_l
	for i, e := range l {
		if i != index {
			cp = append(cp, e)
		}
	}
	return cp
}

func (l entity_build_l) ToEntityBuildList() (el EntityBuildList) {
	for _, e := range l {
		eb := e.ToEntityBuild()
		if eb == nil {
			continue
		}
		el = append(el, eb)
	}
	return el
}
