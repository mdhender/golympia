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
	"time"
)

type OlyTime struct {
	Day            int `json:"day"`            /* day of month */
	Turn           int `json:"turn"`           /* turn number */
	DaysSinceEpoch int `json:"DaysSinceEpoch"` /* days since game begin */
}

// SysData is the json version of the system data.
type SysData struct {
	CreatedAt            time.Time `json:"created-at,omitempty"`
	UpdatedAt            time.Time `json:"updated-at,omitempty"`
	SysClock             OlyTime   `json:"sys-clock"`
	GameNumber           int       `json:"game-number,omitempty"`
	OpenEnded            bool      `json:"open-ended,omitempty"`
	TurnLimit            int       `json:"turn-limit,omitempty"`
	TurnCharge           string    `json:"turn-charge,omitempty"`
	AutoDrop             bool      `json:"auto-drop,omitempty"`
	FromHost             string    `json:"from-host,omitempty"`
	ReplyHost            string    `json:"reply-host,omitempty"`
	XSize                int       `json:"x-size,omitempty"`
	YSize                int       `json:"y-size,omitempty"`
	GMPlayer             int       `json:"gm-player,omitempty"`
	CombatPlayer         int       `json:"combat-player,omitempty"`
	DesertedPlayer       int       `json:"deserted-player,omitempty"`
	IndepPlayer          int       `json:"indep-player,omitempty"`
	SkillPlayer          int       `json:"skill-player,omitempty"`
	Seed                 [3]int    `json:"seed,omitempty"`
	AccountingDir        string    `json:"accounting-dir,omitempty"`
	AccountingProg       string    `json:"accounting-prog,omitempty"`
	CPP                  string    `json:"cpp,omitempty"`
	HTMLPath             string    `json:"html-path,omitempty"`
	HTMLPasswords        string    `json:"html-passwords,omitempty"`
	Free                 bool      `json:"free,omitempty"`
	FullMarkets          bool      `json:"full-markets,omitempty"`
	GuildTeaching        bool      `json:"guild-teaching,omitempty"`
	MonsterSublocInit    bool      `json:"monster-subloc-init,omitempty"`
	MPAntipathy          bool      `json:"mp-antipathy,omitempty"`
	PopulationInit       bool      `json:"population-init,omitempty"`
	PostHasBeenRun       bool      `json:"post-has-been-run,omitempty"`
	SeedHasBeenRun       bool      `json:"seed-has-been-run,omitempty"`
	SurviveNP            bool      `json:"survive-np,omitempty"`
	CheckBalance         int       `json:"check-balance,omitempty"`
	ClaimGive            int       `json:"claim-give,omitempty"`
	CloudRegion          int       `json:"cloud-region,omitempty"`
	CookieInit           int       `json:"cookie-init,omitempty"`
	DeathNPs             int       `json:"death-n-ps,omitempty"`
	DistSeaCompute       int       `json:"dist-sea-compute,omitempty"`
	FaeryPlayer          int       `json:"faery-player,omitempty"`
	FaeryRegion          int       `json:"faery-region,omitempty"`
	FreeNPLimit          int       `json:"free-np-limit,omitempty"`
	HadesPit             int       `json:"hades-pit,omitempty"`
	HadesPlayer          int       `json:"hades-player,omitempty"`
	HadesRegion          int       `json:"hades-region,omitempty"`
	MarketAge            int       `json:"market-age,omitempty"`
	NearCityInit         int       `json:"near-city-init,omitempty"`
	NPCPlayer            int       `json:"npc-player,omitempty"`
	NumBooks             int       `json:"num-books,omitempty"`
	OutputTags           int       `json:"output-tags,omitempty"`
	TimesPay             int       `json:"times-pay,omitempty"`
	PietyLimit           int       `json:"piety-limit,omitempty"`
	MinPiety             int       `json:"min-piety,omitempty"`
	TopPiety             int       `json:"top-piety,omitempty"`
	MiddlePiety          int       `json:"middle-piety,omitempty"`
	BottomPiety          int       `json:"bottom-piety,omitempty"`
	HeadPriestPietyLimit int       `json:"head-priest-piety-limit,omitempty"`

	NL int `json:"nl,omitempty"`
	NR int `json:"nr,omitempty"`
	TR int `json:"tr,omitempty"`
	UR int `json:"ur,omitempty"`
}

func SysDataLoad(name string) (*SysData, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("SysDataLoad: %w", err)
	}
	js := &SysData{}
	if err := json.Unmarshal(data, &js); err != nil {
		log.Printf("load_system: system: %v\n", err)
		return nil, fmt.Errorf("SysDataLoad: %w", err)
	}

	// reset values for piety if the user set them all to zero.
	if js.TopPiety == 0 && js.MiddlePiety == 0 && js.BottomPiety == 0 {
		js.TopPiety = 12
		js.MiddlePiety = 6
		js.BottomPiety = 3
	}

	// enforce ordering on the piety levels
	if js.TopPiety < js.MiddlePiety {
		js.TopPiety, js.MiddlePiety = js.MiddlePiety, js.TopPiety
	}
	if js.TopPiety < js.BottomPiety {
		js.TopPiety, js.BottomPiety = js.BottomPiety, js.TopPiety
	}
	if js.MiddlePiety < js.BottomPiety {
		js.MiddlePiety, js.BottomPiety = js.BottomPiety, js.MiddlePiety
	}

	cloud_region = js.CloudRegion
	combat_pl = js.CombatPlayer
	cookie_init = js.CookieInit
	deserted_player = js.DesertedPlayer
	dist_sea_compute = js.DistSeaCompute
	faery_player = js.FaeryPlayer
	faery_region = js.FaeryRegion
	from_host = js.FromHost
	game_number = js.GameNumber
	gm_player = js.GMPlayer
	hades_pit = js.HadesPit
	hades_player = js.HadesPlayer
	hades_region = js.HadesRegion
	indep_player = js.IndepPlayer
	monster_subloc_init = js.MonsterSublocInit
	near_city_init = js.NearCityInit
	npc_pl = js.NPCPlayer
	sysclock.day = js.SysClock.Day
	sysclock.turn = js.SysClock.Turn
	sysclock.days_since_epoch = js.SysClock.DaysSinceEpoch
	options.created_at = js.CreatedAt
	options.updated_at = js.UpdatedAt
	options.accounting_dir = js.AccountingDir
	options.accounting_prog = js.AccountingProg
	options.auto_drop = js.AutoDrop
	options.bottom_piety = js.BottomPiety
	options.check_balance = js.CheckBalance
	options.claim_give = js.ClaimGive
	options.cpp = js.CPP
	options.death_nps = js.DeathNPs
	options.free = js.Free
	options.free_np_limit = js.FreeNPLimit
	options.full_markets = js.FullMarkets
	options.guild_teaching = js.GuildTeaching
	options.head_priest_piety_limit = js.HeadPriestPietyLimit
	options.html_passwords = js.HTMLPasswords
	options.html_path = js.HTMLPath
	options.market_age = js.MarketAge
	options.middle_piety = js.MiddlePiety
	options.min_piety = js.MinPiety
	options.mp_antipathy = js.MPAntipathy
	options.num_books = js.NumBooks
	options.open_ended = js.OpenEnded
	options.output_tags = js.OutputTags
	options.piety_limit = js.PietyLimit
	options.survive_np = js.SurviveNP
	options.times_pay = js.TimesPay
	options.top_piety = js.TopPiety
	options.turn_charge = js.TurnCharge
	options.turn_limit = js.TurnLimit
	population_init = js.PopulationInit
	if js.PostHasBeenRun {
		post_has_been_run = TRUE
	} else {
		post_has_been_run = FALSE
	}
	reply_host = js.ReplyHost
	seed[0] = js.Seed[0]
	seed[1] = js.Seed[1]
	seed[2] = js.Seed[2]
	if js.SeedHasBeenRun {
		seed_has_been_run = TRUE
	} else {
		seed_has_been_run = FALSE
	}
	skill_player = js.SkillPlayer
	xsize = js.XSize
	ysize = js.YSize

	return js, nil
}

func SysDataSave(name string) error {
	var js SysData

	js.CreatedAt = options.created_at
	if js.CreatedAt.IsZero() {
		js.CreatedAt = time.Now().UTC()
	}
	js.UpdatedAt = time.Now().UTC()
	js.CloudRegion = cloud_region
	js.CombatPlayer = combat_pl
	js.CookieInit = cookie_init
	js.DesertedPlayer = deserted_player
	js.DistSeaCompute = dist_sea_compute
	js.FaeryPlayer = faery_player
	js.FaeryRegion = faery_region
	js.FromHost = from_host
	js.GameNumber = game_number
	js.GMPlayer = gm_player
	js.HadesPit = hades_pit
	js.HadesPlayer = hades_player
	js.HadesRegion = hades_region
	js.IndepPlayer = indep_player
	js.MonsterSublocInit = monster_subloc_init
	js.NearCityInit = near_city_init
	js.NPCPlayer = npc_pl
	js.SysClock.Day = sysclock.day
	js.SysClock.Turn = sysclock.turn
	js.SysClock.DaysSinceEpoch = sysclock.days_since_epoch
	js.AccountingDir = options.accounting_dir
	js.AccountingProg = options.accounting_prog
	js.AutoDrop = options.auto_drop
	js.BottomPiety = options.bottom_piety
	js.CheckBalance = options.check_balance
	js.ClaimGive = options.claim_give
	js.CPP = options.cpp
	js.DeathNPs = options.death_nps
	js.Free = options.free
	js.FreeNPLimit = options.free_np_limit
	js.FullMarkets = options.full_markets
	js.GuildTeaching = options.guild_teaching
	js.HeadPriestPietyLimit = options.head_priest_piety_limit
	js.HTMLPasswords = options.html_passwords
	js.HTMLPath = options.html_path
	js.MarketAge = options.market_age
	js.MiddlePiety = options.middle_piety
	js.MinPiety = options.min_piety
	js.MPAntipathy = options.mp_antipathy
	js.NumBooks = options.num_books
	js.OpenEnded = options.open_ended
	js.OutputTags = options.output_tags
	js.PietyLimit = options.piety_limit
	js.SurviveNP = options.survive_np
	js.TimesPay = options.times_pay
	js.TopPiety = options.top_piety
	js.TurnCharge = options.turn_charge
	js.TurnLimit = options.turn_limit
	js.PopulationInit = population_init
	js.PostHasBeenRun = post_has_been_run != FALSE
	js.ReplyHost = reply_host
	js.Seed[0] = seed[0]
	js.Seed[1] = seed[1]
	js.Seed[2] = seed[2]
	js.SeedHasBeenRun = seed_has_been_run != FALSE
	js.SkillPlayer = skill_player
	js.XSize = xsize
	js.YSize = ysize

	data, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return fmt.Errorf("SysDataSave: %w", err)
	} else if err := os.WriteFile(name, data, 0666); err != nil {
		return fmt.Errorf("SysDataSave: %w", err)
	}

	return nil
}
