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
	"path/filepath"
	"sort"
	"strings"
)

type CharMagic struct {
	AbilityShroud    int    `json:"ability-shroud,omitempty"`    //
	AuraReflect      int    `json:"aura-reflect,omitempty"`      // reflect aura blast
	Auraculum        int    `json:"auraculum,omitempty"`         // char created an auraculum
	CurAura          int    `json:"cur-aura,omitempty"`          // current aura level for magician
	HideMage         int    `json:"hide-mage,omitempty"`         // number of points hiding the magician
	HideSelf         bool   `json:"hide-self,omitempty"`         // character is hidden
	HinderMeditation int    `json:"hinder-meditation,omitempty"` // number of points to hinder, usually 0...3
	KnowsWeather     bool   `json:"knows-weather,omitempty"`     // knows weather magic
	Magician         bool   `json:"magician,omitempty"`          // is a magician
	MaxAura          int    `json:"max-aura,omitempty"`          // maximum aura level for magician
	ProjectCast      int    `json:"project-cast,omitempty"`      // project next cast
	QuickCast        int    `json:"quick-cast,omitempty"`        // speed next cast
	SwearOnRelease   int    `json:"swear-on-release,omitempty"`  // swear to one who frees us
	Token            int    `json:"token,omitempty"`             // we are controlled by this art
	Visions          sparse `json:"visions,omitempty"`           // visions revealed

	// Pledge  int    // lands are pledged to another // not used?

	// the following are not saved
	ferry_flag    bool   // ferry has tooted its horn
	mage_worked   int    // worked this month
	pledged_to_us ints_l // temp
}

type CharReligion struct {
	Priest    int    /* Who this noble is dedicated to, if anyone. */
	Piety     int    /* Our current piety. */
	Followers ints_l /* Who is dedicated to us, if anyone. */
}

type EntityChar struct {
	Attack             int           `json:"attack,omitempty"`               // fighter attack rating
	Behind             int           `json:"behind,omitempty"`               // are we behind in combat?
	BreakPoint         int           `json:"breakPoint,omitempty"`           // break point when fighting
	Contact            ints_l        `json:"contact,omitempty"`              // who have we contacted, also, who has found us
	DeathTime          *OlyTime      `json:"death-time,omitempty"`           // when was character killed
	Defense            int           `json:"defense,omitempty"`              // fighter defense rating
	Guard              int           `json:"guard,omitempty"`                // character is guarding the loc
	Guild              int           `json:"guild,omitempty"`                // This is the guild we belong to.
	Health             int           `json:"health,omitempty"`               // current health
	LoyKind            int           `json:"loy-kind,omitempty"`             // LOY_xxx
	LoyRate            int           `json:"loy-rate,omitempty"`             // level with kind of loyalty
	Missile            int           `json:"missile,omitempty"`              // capable of missile attacks?
	Moving             int           `json:"moving,omitempty"`               // daystamp of beginning of movement
	NpcProg            int           `json:"npc-prog,omitempty"`             // npc program
	Pay                int           `json:"pay,omitempty"`                  // How much will you pay to enter?
	PersonalBreakPoint int           `json:"personal-break-point,omitempty"` // personal break point when fighting
	Prisoner           int           `json:"prisoner,omitempty"`             // is this character a prisoner?
	Rank               int           `json:"rank,omitempty"`                 // noble peerage status
	Religion           *CharReligion `json:"religion,omitempty"`             // Our religion info...
	Sick               int           `json:"sick,omitempty"`                 // 1=character is getting worse
	Skills             skill_ent_l   `json:"skills,omitempty"`               // skills known by char
	TimeFlying         int           `json:"time-flying,omitempty"`          // time airborne over ocean
	UnitItem           int           `json:"unit-item,omitempty"`            // unit is made of this kind of item
	UnitLord           int           `json:"unit-lord,omitempty"`            // who is our owner?
	XCharMagic         *CharMagic    `json:"x-char-magic,omitempty"`

	// Effects []*effect // list of effects on char // not used?

	// the following are not saved
	accept     accept_ent_l // what we can be given
	fresh_hire int          // don't erode loyalty
	melt_me    int          // in process of melting away
	new_lord   int          // got a new lord this turn
	studied    int          // num days we studied
}

type CharacterList []*Character

func (l CharacterList) Len() int {
	return len(l)
}

func (l CharacterList) Less(i, j int) bool {
	return l[i].Id < l[j].Id
}

func (l CharacterList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Character struct {
	Id   int    `json:"id"`             // identity of the character
	Name string `json:"name,omitempty"` // name of the character
}

func CharactersLoad(scanOnly bool) error {
	path := filepath.Join(libdir, "characters")
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("CharactersLoad: scan: %w", err)
	}

	for _, f := range files {
		jsonFile := isdigit(f.Name()[0]) && strings.HasSuffix(f.Name(), ".json")
		if !jsonFile {
			continue
		}
		//scan_boxes(filepath.Join("fact", f.Name()))
		_, _ = CharacterDataLoad(filepath.Join(path, f.Name()), scanOnly)
	}

	return nil
}

func CharactersSave() error {
	return fmt.Errorf("CharactersSave: not implemented")
}

func CharacterDataLoad(name string, scanOnly bool) (CharacterList, error) {
	log.Printf("CharacterDataLoad: loading %s\n", name)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("CharacterDataLoad: %w", err)
	}
	var list CharacterList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("CharacterDataLoad: %w", err)
	}
	if scanOnly {
		return nil, nil
	}
	return nil, nil
}

func CharacterDataSave(name string) error {
	list := CharacterList{}
	sort.Sort(list)
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("CharacterDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("CharacterDataSave: %w", err)
	}
	return nil
}

func (e *entity_char) ToEntityChar() *EntityChar {
	if e == nil {
		return nil
	}
	ec := &EntityChar{
		Attack:             e.attack,
		Behind:             e.behind,
		BreakPoint:         e.break_point,
		Contact:            e.contact,
		DeathTime:          e.death_time.ToOlyTime(),
		Defense:            e.defense,
		Guard:              e.guard,
		Guild:              e.guild,
		Health:             e.health,
		LoyKind:            e.loy_kind,
		LoyRate:            e.loy_rate,
		Missile:            e.missile,
		Moving:             e.moving,
		NpcProg:            e.npc_prog,
		Pay:                e.pay,
		PersonalBreakPoint: e.personal_break_point,
		Prisoner:           e.prisoner,
		Rank:               e.rank,
		Religion:           e.religion.ToCharReligion(),
		Sick:               e.sick,
		Skills:             e.skills,
		TimeFlying:         e.time_flying,
		UnitItem:           e.unit_item,
		UnitLord:           e.unit_lord,
		XCharMagic:         e.x_char_magic.ToCharMagic(),
		fresh_hire:         e.fresh_hire,
		melt_me:            e.melt_me,
		new_lord:           e.new_lord,
		studied:            e.studied,
	}
	if len(e.accept) != 0 {
		ec.accept = append(ec.accept, e.accept...)
	}

	return ec
}

func (r char_religion) ToCharReligion() *CharReligion {
	if r.IsZero() {
		return nil
	}
	c := &CharReligion{
		Priest: r.priest,
		Piety:  r.piety,
	}
	if len(r.followers) != 0 {
		c.Followers = append(c.Followers, r.followers...)
	}
	return c
}

func (x *char_magic) ToCharMagic() *CharMagic {
	if x == nil {
		return nil
	}
	c := &CharMagic{
		AbilityShroud:    x.ability_shroud,
		AuraReflect:      x.aura_reflect,
		Auraculum:        x.auraculum,
		CurAura:          x.cur_aura,
		HideMage:         x.hide_mage,
		HideSelf:         x.hide_self != FALSE,
		HinderMeditation: x.hinder_meditation,
		KnowsWeather:     x.knows_weather != FALSE,
		Magician:         x.magician != FALSE,
		MaxAura:          x.max_aura,
		ProjectCast:      x.project_cast,
		QuickCast:        x.quick_cast,
		SwearOnRelease:   x.swear_on_release,
		Token:            x.token,
		ferry_flag:       x.ferry_flag,
		mage_worked:      x.mage_worked,
	}
	if len(x.pledged_to_us) != 0 {
		c.pledged_to_us = append(c.pledged_to_us, x.pledged_to_us...)
	}
	if len(x.visions) != 0 {
		c.Visions = append(c.Visions, x.visions...)
	}
	return c
}
