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
	"bytes"
	"log"
	"unicode"
	"unicode/utf8"
)

type OlyTime struct {
	Day            int `json:"day"`            /* day of month */
	Turn           int `json:"turn"`           /* turn number */
	DaysSinceEpoch int `json:"DaysSinceEpoch"` /* days since game begin */
}

// SysData is the json version of the system data.
type SysData struct {
	SysClock             OlyTime `json:"sys-clock"`
	GameNumber           int     `json:"game-number,omitempty"`
	OpenEnded            bool    `json:"open-ended,omitempty"`
	TurnLimit            int     `json:"turn-limit,omitempty"`
	TurnCharge           string  `json:"turn-charge,omitempty"`
	AutoDrop             bool    `json:"auto-drop,omitempty"`
	FromHost             string  `json:"from-host,omitempty"`
	ReplyHost            string  `json:"reply-host,omitempty"`
	XSize                int     `json:"x-size,omitempty"`
	YSize                int     `json:"y-size,omitempty"`
	GMPlayer             int     `json:"gm-player,omitempty"`
	CombatPlayer         int     `json:"combat-player,omitempty"`
	DesertedPlayer       int     `json:"deserted-player,omitempty"`
	IndepPlayer          int     `json:"indep-player,omitempty"`
	SkillPlayer          int     `json:"skill-player,omitempty"`
	Seed                 [3]int  `json:"seed,omitempty"`
	AccountingDir        string  `json:"accounting-dir,omitempty"`
	AccountingProg       string  `json:"accounting-prog,omitempty"`
	CPP                  string  `json:"cpp,omitempty"`
	HTMLPath             string  `json:"html-path,omitempty"`
	HTMLPasswords        string  `json:"html-passwords,omitempty"`
	Free                 bool    `json:"free,omitempty"`
	FullMarkets          bool    `json:"full-markets,omitempty"`
	GuildTeaching        bool    `json:"guild-teaching,omitempty"`
	MonsterSublocInit    bool    `json:"monster-subloc-init,omitempty"`
	MPAntipathy          bool    `json:"mp-antipathy,omitempty"`
	PopulationInit       bool    `json:"population-init,omitempty"`
	PostHasBeenRun       bool    `json:"post-has-been-run,omitempty"`
	SeedHasBeenRun       bool    `json:"seed-has-been-run,omitempty"`
	SurviveNP            bool    `json:"survive-np,omitempty"`
	CheckBalance         int     `json:"check-balance,omitempty"`
	ClaimGive            int     `json:"claim-give,omitempty"`
	CloudRegion          int     `json:"cloud-region,omitempty"`
	CookieInit           int     `json:"cookie-init,omitempty"`
	DeathNPs             int     `json:"death-n-ps,omitempty"`
	DistSeaCompute       int     `json:"dist-sea-compute,omitempty"`
	FaeryPlayer          int     `json:"faery-player,omitempty"`
	FaeryRegion          int     `json:"faery-region,omitempty"`
	FreeNPLimit          int     `json:"free-np-limit,omitempty"`
	HadesPit             int     `json:"hades-pit,omitempty"`
	HadesPlayer          int     `json:"hades-player,omitempty"`
	HadesRegion          int     `json:"hades-region,omitempty"`
	MarketAge            int     `json:"market-age,omitempty"`
	NearCityInit         int     `json:"near-city-init,omitempty"`
	NPCPlayer            int     `json:"npc-player,omitempty"`
	NumBooks             int     `json:"num-books,omitempty"`
	OutputTags           int     `json:"output-tags,omitempty"`
	TimesPay             int     `json:"times-pay,omitempty"`
	PietyLimit           int     `json:"piety-limit,omitempty"`
	MinPiety             int     `json:"min-piety,omitempty"`
	TopPiety             int     `json:"top-piety,omitempty"`
	MiddlePiety          int     `json:"middle-piety,omitempty"`
	BottomPiety          int     `json:"bottom-piety,omitempty"`
	HeadPriestPietyLimit int     `json:"head-priest-piety-limit,omitempty"`

	NL int `json:"nl,omitempty"`
	NR int `json:"nr,omitempty"`
	TR int `json:"tr,omitempty"`
	UR int `json:"ur,omitempty"`
}

func System(buf []byte) *SysData {
	lines := bytes.Split(buf, []byte{'\n'})

	system := &SysData{
		TopPiety:    12,
		MiddlePiety: 6,
		BottomPiety: 3,
	}

	for _, line := range lines {
		line = bytes.TrimRightFunc(line, func(r rune) bool {
			return r == utf8.RuneError || unicode.IsSpace(r)
		})
		// ignore comments and blank lines
		if len(line) == 0 || bytes.HasPrefix(line, []byte{'#'}) {
			continue
		}
		// olytime is special
		if bytes.HasPrefix(line, []byte("sysclock:")) {
			system.SysClock = olytimeFromSlice(line[9:])
			continue
		}
		if bytes.IndexByte(line, '=') == -1 {
			log.Printf("system: invalid line: %q\n", string(line))
			continue
		}
		f := bytes.FieldsFunc(line, func(r rune) bool {
			return r == '='
		})
		if len(f) != 2 {
			log.Printf("system: invalid line: %q\n", string(line))
			continue
		}
		switch string(f[0]) {
		case "accounting_dir":
			system.AccountingDir = string(f[1])
		case "accounting_prog":
			system.AccountingProg = string(f[1])
		case "autodrop":
			system.AutoDrop = btoi(f[1]) != 0
		case "bottom_piety":
			system.BottomPiety = btoi(f[1])
		case "check_balance":
			system.CheckBalance = btoi(f[1])
		case "claim_give":
			system.ClaimGive = btoi(f[1])
		case "cp":
			system.CombatPlayer = btoi(f[1])
		case "cpp":
			system.CPP = string(f[1])
		case "cr":
			system.CloudRegion = btoi(f[1])
		case "death_nps":
			system.DeathNPs = btoi(f[1])
		case "deserted_player":
			system.DesertedPlayer = btoi(f[1])
		case "ds":
			system.DistSeaCompute = btoi(f[1])
		case "fp":
			system.FaeryPlayer = btoi(f[1])
		case "full_markets":
			system.FullMarkets = btoi(f[1]) != 0
		case "fr":
			system.FaeryRegion = btoi(f[1])
		case "free":
			system.Free = btoi(f[1]) != 0
		case "free_np_limit":
			system.FreeNPLimit = btoi(f[1])
		case "from_host":
			system.FromHost = string(f[1])
		case "game_num":
			system.GameNumber = btoi(f[1])
		case "gm_player":
			system.GMPlayer = btoi(f[1])
		case "guild_teaching":
			system.GuildTeaching = btoi(f[1]) != 0
		case "head_priest_piety_limit":
			system.HeadPriestPietyLimit = btoi(f[1])
		case "hl":
			system.HadesPlayer = btoi(f[1])
		case "hp":
			system.HadesPit = btoi(f[1])
		case "hr":
			system.HadesRegion = btoi(f[1])
		case "html_path":
			system.HTMLPath = string(f[1])
		case "html_passwords":
			system.HTMLPasswords = string(f[1])
		case "indep_player":
			system.IndepPlayer = btoi(f[1])
		case "init":
			system.SeedHasBeenRun = btoi(f[1]) != 0
		case "market_age":
			system.MarketAge = btoi(f[1])
		case "mi":
			system.CookieInit = btoi(f[1])
		case "middle_piety":
			system.MiddlePiety = btoi(f[1])
		case "min_piety":
			system.MinPiety = btoi(f[1])
		case "mp_antipathy":
			system.MPAntipathy = btoi(f[1]) != 0
		case "ms":
			system.MonsterSublocInit = btoi(f[1]) != 0
		case "nc":
			system.NearCityInit = btoi(f[1])
		case "nl":
			system.NL = btoi(f[1])
		case "np":
			system.NPCPlayer = btoi(f[1])
		case "nr":
			system.NR = btoi(f[1])
		case "num_books":
			system.NumBooks = btoi(f[1])
		case "open_ended":
			system.OpenEnded = btoi(f[1]) != 0
		case "output_tags":
			system.OutputTags = btoi(f[1])
		case "pi":
			system.PopulationInit = btoi(f[1]) != 0
		case "piety_limit":
			system.PietyLimit = btoi(f[1])
		case "post":
			system.PostHasBeenRun = btoi(f[1]) != 0
		case "reply_host":
			system.ReplyHost = string(f[1])
		case "seed0":
			system.Seed[0] = btoi(f[1])
		case "seed1":
			system.Seed[1] = btoi(f[1])
		case "seed2":
			system.Seed[2] = btoi(f[1])
		case "skill_player":
			system.SkillPlayer = btoi(f[1])
		case "survive_np":
			system.SurviveNP = btoi(f[1]) != 0
		case "times_pay":
			system.TimesPay = btoi(f[1])
		case "top_piety":
			system.TopPiety = btoi(f[1])
		case "tr":
			system.TR = btoi(f[1])
		case "turn_charge":
			system.TurnCharge = string(f[1])
		case "turn_limit":
			system.TurnLimit = btoi(f[1])
		case "ur":
			system.UR = btoi(f[1])
		case "xsize":
			system.XSize = btoi(f[1])
		case "ysize":
			system.YSize = btoi(f[1])
		default:
			log.Printf("system: unrecognized field: %q\n", string(f[0]))
		}
	}

	// reset values for piety if the user set them all to zero.
	if system.TopPiety == 0 && system.MiddlePiety == 0 && system.BottomPiety == 0 {
		system.TopPiety = 12
		system.MiddlePiety = 6
		system.BottomPiety = 3
	}

	// enforce ordering on the piety levels
	if system.TopPiety < system.MiddlePiety {
		system.TopPiety, system.MiddlePiety = system.MiddlePiety, system.TopPiety
	}
	if system.TopPiety < system.BottomPiety {
		system.TopPiety, system.BottomPiety = system.BottomPiety, system.TopPiety
	}
	if system.MiddlePiety < system.BottomPiety {
		system.MiddlePiety, system.BottomPiety = system.BottomPiety, system.MiddlePiety
	}

	return system
}

func btoi(b []byte) int {
	return atoi_b(b)
}
