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
	"time"
)

// oly.h

const (
	DEFAULT_PASSWORD = "defpwd123"

	MAX_BOXES  = 102400
	MONTH_DAYS = 30
	NUM_MONTHS = 8

	PROG_bandit         = 1 /* wilderness spice */
	PROG_subloc_monster = 2
	PROG_npc_token      = 3
	PROG_balrog         = 4
	PROG_dumb_monster   = 5
	PROG_smart_monster  = 6
	PROG_orc            = 7
	PROG_elf            = 8
	PROG_daemon         = 9

	use_death_potion          = 1
	use_heal_potion           = 2
	use_slave_potion          = 3
	use_palantir              = 4
	use_proj_cast             = 5 /* stored projected cast */
	use_quick_cast            = 6 /* stored cast speedup */
	use_drum                  = 7 /* beat savage's drum */
	use_faery_stone           = 8 /* Faery gate opener */
	use_orb                   = 9 /* crystal orb */
	use_barbarian_kill        = 10
	use_savage_kill           = 11
	use_corpse_kill           = 12
	use_orc_kill              = 13
	use_skeleton_kill         = 14
	use_ancient_aura          = 15 /* bta's auraculum */
	use_weightlessness_potion = 16
	use_fiery_potion          = 17 /* Fiery Death */
	use_nothing               = 18 /* Something that does nothing */
	use_special_staff         = 19

	DIR_N    = 1
	DIR_E    = 2
	DIR_S    = 3
	DIR_W    = 4
	DIR_UP   = 5
	DIR_DOWN = 6
	DIR_IN   = 7
	DIR_OUT  = 8
	MAX_DIR  = 9 /* one past highest direction */

	LOC_region   = 1 /* top most continent/island group */
	LOC_province = 2 /* main location area */
	LOC_subloc   = 3 /* inner sublocation */
	LOC_build    = 4 /* building, structure, etc. */

	LOY_UNCHANGED = (-1)
	LOY_unsworn   = 0
	LOY_contract  = 1
	LOY_oath      = 2
	LOY_fear      = 3
	LOY_npc       = 4
	LOY_summon    = 5

	exp_novice     = 1 /* apprentice */
	exp_journeyman = 2
	exp_teacher    = 3
	exp_master     = 4
	exp_grand      = 5 /* grand master */

	S_body    = 1 /* kill to a dead body */
	S_soul    = 2 /* kill to a lost soul */
	S_nothing = 3 /* kill completely */

	MONSTER_ATT = -1 /* phony id for "all monsters" */
	ATT_NONE    = 0  /* no attitude -- default */
	NEUTRAL     = 1  /* explicitly neutral */
	HOSTILE     = 2  /* attack on sight */
	DEFEND      = 3  /* defend if attacked */

	// effect identifiers
	ef_defense            = 1001 /* Add to defense of stack. */
	ef_religion_trap      = 1002 /* A religion trap. */
	ef_fast_move          = 1004 /* Fast move into the next province */
	ef_slow_move          = 1005 /* Slow moves into this province */
	ef_improve_mine       = 1006 /* Improve a mine's production */
	ef_protect_mine       = 1007 /* Protect a mine against calamities */
	ef_bless_fort         = 1008 /* Bless a fort */
	ef_improve_production = 1010 /* Generic production improvement +50%*/
	ef_improve_make       = 1011 /* Generic make improvement +100% */
	ef_improve_fort       = 1012 /* Increase a fortification */
	ef_edge_of_kireus     = 1013 /* Edged weapons +25 attack */
	ef_urchin_spy         = 1014 /* Report unit's doings for 7 days */
	ef_charisma           = 1015 /* Double leadership */
	ef_improve_taxes      = 1016 /* Double taxes */
	ef_guard_loyalty      = 1017 /* Unit immune to loyalty checks */
	ef_hide_item          = 1018 /* Hide an item. */
	ef_hide_money         = 1019 /* Hide money. */
	ef_smuggle_goods      = 1020 /* Smuggling. */
	ef_smuggle_men        = 1021 /* Ditto. */
	ef_grow               = 1022 /* Encourage pop to grow. */
	ef_weightlessness     = 1023 /* Negative weight */
	ef_scry_offense       = 1024 /* + to offense in battle */
	ef_conceal_nation     = 1025 /* conceal your nation */
	ef_kill_dirt_golem    = 1026 /* Timer on golems */
	ef_kill_flesh_golem   = 1027 /* Timer on golems */
	ef_food_found         = 1028 /* Food found by a ranger */
	ef_conceal_artifacts  = 1029 /* Conceal someone's artifacts */
	ef_obscure_artifact   = 1030 /* Obscure artifact's identity */
	ef_forced_march       = 1031 /* Next move at riding speed */
	ef_faery_warning      = 1032 /* Give people a month grace. */
	ef_tap_wound          = 1033 /* The wound from Tap Health. */
	ef_magic_barrier      = 1034 /* New implementation of magic barrier */
	ef_inhibit_barrier    = 1035 /* After magic barrier falls. */
	ef_cs                 = 1036 /* Combat spell control. */

	// artifacts - generalized magic items -- not unlike effects.
	ART_NONE        = 0  /* No effect */
	ART_COMBAT      = 1  /* + to combat */
	ART_CTL_MEN     = 2  /* + to controlled men */
	ART_CTL_BEASTS  = 3  /* + to controlled beasts */
	ART_SAFETY      = 4  /* safety from one monster type */
	ART_IMPRV_ATT   = 5  /* improve attack of particular unit */
	ART_IMPRV_DEF   = 6  /* improve defense of particular unit */
	ART_SAFE_SEA    = 7  /* safety at sea */
	ART_TERRAIN     = 8  /* improved defense in particular terrain */
	ART_FAST_TERR   = 9  /* fast movement into particular terrain */
	ART_SPEED_USE   = 10 /* speed use of skill */
	ART_PROT_HADES  = 11 /* protection in Hades */
	ART_PROT_FAERY  = 12 /* protection in Faery */
	ART_WORKERS     = 13 /* improve productivity of workers */
	ART_INCOME      = 14 /* improved income from X */
	ART_LEARNING    = 15 /* improved learning */
	ART_TEACHING    = 16 /* improved teaching */
	ART_TRAINING    = 17 /* improved training (of X?) */
	ART_DESTROY     = 18 /* destroy monster X (uses) */
	ART_SKILL       = 19 /* grants a skill */
	ART_FLYING      = 20 /* permit user to fly */
	ART_PROT_SKILL  = 21 /* protection from a skill (aura blast,scry) */
	ART_SHIELD_PROV = 22 /* protect entire province from scrying */
	ART_RIDING      = 23 /* permit user riding pace */
	ART_POWER       = 24 /* increase piety/aura (1 use) */
	ART_SUMMON_AID  = 25 /* summon help (1 use) */
	ART_MAINTENANCE = 26 /* reduce maintenance costs */
	ART_BARGAIN     = 27 /* better market prices */
	ART_WEIGHTLESS  = 28 /* weightlessness */
	ART_HEALING     = 29 /* faster healing */
	ART_SICKNESS    = 30 /* protection from sickness */
	ART_RESTORE     = 31 /* restore life */
	ART_TELEPORT    = 32 /* teleportation (uses) */
	ART_ORB         = 33 /* orb */
	ART_CROWN       = 34 /* crown */
	ART_AURACULUM   = 35 /* auraculum */
	ART_CARRY       = 36 /* increase land carry capacity */
	ART_PEN         = 37 /* the Pen Crown */
	ART_LAST        = 38

	// bit-masks to pull out the appropriate bits of combat artifacts.
	CA_N_MELEE     = (1 << 0)
	CA_N_MISSILE   = (1 << 1)
	CA_N_SPECIAL   = (1 << 2)
	CA_M_MELEE     = (1 << 3)
	CA_M_MISSILE   = (1 << 4)
	CA_M_SPECIAL   = (1 << 5)
	CA_N_MELEE_D   = (1 << 6)
	CA_N_MISSILE_D = (1 << 7)
	CA_N_SPECIAL_D = (1 << 8)
	CA_M_MELEE_D   = (1 << 9)
	CA_M_MISSILE_D = (1 << 10)
	CA_M_SPECIAL_D = (1 << 11)

	TOUGH_NUM = 2520

	SKILL_dont     = 0 /* don't know the skill */
	SKILL_learning = 1 /* in the process of learning it */
	SKILL_know     = 2 /* know it */

	// build descriptor - describes what kind of build is going on in a location.
	BT_BUILD      = 1
	BT_FORTIFY    = 2
	BT_STRENGTHEN = 3
	BT_MOAT       = 4

	// mine descriptors
	MINE_MAX       = 20
	NO_SHORING     = 0
	WOODEN_SHORING = 1
	IRON_SHORING   = 2

	/*
	 *  Skill flags.
	 *
	 */
	IS_POLLED       = (1 << 0)
	REQ_HOLY_SYMBOL = (1 << 1)
	REQ_HOLY_PLANT  = (1 << 2)
	COMBAT_SKILL    = (1 << 3)
	MAX_FLAGS       = 4

	REQ_NO  = 0 /* don't consume item */
	REQ_YES = 1 /* consume item */
	REQ_OR  = 2 /* or with next */

	// in-process command structure
	DONE  = 0
	LOAD  = 1
	RUN   = 2
	ERROR = 3

	// types of command arguments
	CMD_undef    = 0
	CMD_unit     = 1
	CMD_item     = 2
	CMD_skill    = 3
	CMD_days     = 4
	CMD_qty      = 5
	CMD_gold     = 6
	CMD_use      = 7
	CMD_practice = 8

	// trade flags
	BUY     = 1
	SELL    = 2
	PRODUCE = 3
	CONSUME = 4

	CHAR_FIELD = 5  /* field length for box_code_less */
	MAX_POST   = 60 /* max line length for posts and messages */

	// used in combat.c and move.c
	A_WON     = 1
	B_WON     = 2
	TIE       = 3
	NO_COMBAT = 4

	// defines a faction as either MUs or Priests.
	MU_FACTION     = 1
	PRIEST_FACTION = 2

	// define format flags
	NONE = 0
	HTML = (1 << 0)
	TEXT = (1 << 1)
	TAGS = (1 << 2)
	RAW  = (1 << 3)
	ALT  = (1 << 4)
)

var (
	// global box

	LINES = "---------------------------------------------"

	dir_s           []string
	evening         bool /* are we in the evening phase? */
	indent          int
	trades_to_check []int
)

type accept_ent_l []*accept_ent

type accept_ent struct {
	item     int /* 0 = any item */
	from_who int /* 0 = anyone, else char or player */
	qty      int /* 0 = any qty */
}

type Admit struct {
	Flag  int    // first time set this turn -- not saved
	List  ints_l //
	Sense int    // 0=default no, 1=all but..
	Targ  int    // char or loc Admit is declared for
}

type box struct {
	kind       schar
	skind      schar
	name       string
	x_loc_info loc_info
	x_player   *EntityPlayer
	x_char     *entity_char
	x_loc      *entity_loc
	x_subloc   *entity_subloc
	x_item     *entity_item
	x_skill    *entity_skill
	x_nation   *entity_nation
	x_gate     *entity_gate
	x_misc     *entity_misc
	x_disp     *att_ent

	cmd     *command
	items   item_ent_l /* ilist of items held */
	trades  trade_l    /* pending buys/sells */
	effects effect_l   /* ilist of effects */

	temp         int /* scratch space */
	output_order int /* for report ordering -- not saved */

	x_next_kind int /* link to next entity of same type */
	x_next_sub  int /* link to next of same subkind */
}

type char_magic struct {
	max_aura  int /* maximum aura level for magician */
	cur_aura  int /* current aura level for magician */
	auraculum int /* char created an auraculum */

	visions sparse /* visions revealed */
	// pledge  int    /* lands are pledged to another */ // not used?
	token int /* we are controlled by this art */

	project_cast      int /* project next cast */
	quick_cast        int /* speed next cast */
	ability_shroud    int
	hide_mage         int // number of points hiding the magician
	hinder_meditation int // number of points to hinder, usually 0...3
	magician          int /* is a magician */
	aura_reflect      int /* reflect aura blast */
	hide_self         int /* character is hidden */
	swear_on_release  int /* swear to one who frees us */
	knows_weather     int /* knows weather magic */

	mage_worked   int    /* worked this month -- not saved */
	ferry_flag    bool   /* ferry has tooted its horn -- ns */
	pledged_to_us ints_l /* temp -- not saved */
}

// character religion - c aptures a nobles religious standing, such as it is.
type char_religion struct {
	priest    int    /* Who this noble is dedicated to, if anyone. */
	piety     int    /* Our current piety. */
	followers ints_l /* Who is dedicated to us, if anyone. */
}

func (cr char_religion) IsZero() bool {
	return cr.priest == 0 && cr.piety == 0 && len(cr.followers) == 0
}

type cmd_tbl_ent struct {
	allow string /* who may execute the command */
	name  string /* name of command */

	start     func(c *command) int /* initiator */
	finish    func(c *command) int /* conclusion */
	interrupt func(c *command) int /* interrupted order */

	time int /* how long command takes */
	poll int /* call finish each day, not just at end */
	pri  int /* command priority or precedence */

	// mods to add some command checking during the eat phase:
	//     -- Num_args_required
	//     -- Max_args
	//     -- Arg_types[]
	//     -- cmd_comment()
	//     -- cmd_check()
	num_args_required int
	max_args          int
	arg_types         [5]int
	cmd_comment       func(c *command) string
	cmd_check         func(c *command)
}

type command struct {
	who                          int         // entity this is under (redundant)
	wait                         int         // time until completion
	cmd                          int         // index into cmd_tbl
	use_skill                    int         // skill we are using, if any
	use_ent                      int         // index into use_tbl[] for skill usage
	use_exp                      int         // experience level at using this skill
	days_executing               int         // how long has this command been running
	a, b, c, d, e, f, g, h, i, j int         // command arguments
	line                         string      // original command line
	parsed_line                  []byte      // cut-up line, pointed to by parse
	parse                        args_l      // ilist of parsed arguments
	state                        int         // LOAD, RUN, ERROR, DONE
	status                       int         // success or failure; can this be bool?
	poll                         int         // call finish routine each day?
	pri                          int         // command priority or precedence
	conditional                  schar       // 0=none 1=last succeeded 2=last failed
	inhibit_finish               bool        // don't call d_xxx
	fuzzy                        bool        // command matched fuzzy -- not saved
	second_wait                  int         // delay resulting from auto attacks -- saved
	wait_parse                   []*wait_arg // not saved
	// debug is not a bool because it is used to hold multiple values
	debug int // debugging check -- not saved
}

type EntityArtifact struct {
	Id     int `json:"id,omitempty"`
	Type   int `json:"type,omitempty"`
	Param1 int `json:"p1,omitempty"`
	Param2 int `json:"p2,omitempty"`
	Uses   int `json:"uses,omitempty"`
}

func (ea *EntityArtifact) IsZero() bool {
	// https://freshman.tech/snippets/go/check-empty-struct/
	return zero_check(ea)
}

type entity_char struct {
	unit_item int /* unit is made of this kind of item */

	health int
	sick   int /* 1=character is getting worse */

	guard    int /* character is guarding the loc */
	loy_kind int /* LOY_xxx */
	loy_rate int /* level with kind of loyalty */

	death_time olytime /* when was character killed */

	skills skill_ent_l /* ilist of skills known by char */

	// effects []*effect /* ilist of effects on char */ // not used?

	moving    int /* daystamp of beginning of movement */
	unit_lord int /* who is our owner? */

	contact []int /* who have we contacted, also, who has found us */

	x_char_magic         *char_magic
	prisoner             int /* is this character a prisoner? */
	behind               int /* are we behind in combat? */
	time_flying          int /* time airborne over ocean */
	break_point          int /* break point when fighting */
	personal_break_point int /* personal break point when fighting */
	rank                 int /* noble peerage status */
	npc_prog             int /* npc program */

	guild int /* This is the guild we belong to. */

	attack  int /* fighter attack rating */
	defense int /* fighter defense rating */
	missile int /* capable of missile attacks? */

	religion char_religion /* Our religion info... */

	pay int /* How much will you pay to enter? */

	// the following are not saved by io.c:
	melt_me    int          /* in process of melting away */
	fresh_hire int          /* don't erode loyalty */
	new_lord   int          /* got a new lord this turn */
	studied    int          /* num days we studied */
	accept     accept_ent_l /* what we can be given */
}

type entity_gate struct {
	to_loc        int /* destination of gate */
	notify_jumps  int /* whom to notify */
	notify_unseal int /* whom to notify */
	seal_key      int /* numeric gate password */
	road_hidden   int /* this is a hidden road or passage */
}

type entity_item struct {
	id int // unique identifier for this thing

	animal      int /* unit is or contains a horse or an ox */
	animal_part int /* Produces this when killed. */
	attack      int /* fighter attack rating */
	base_price  int /* base price of item for market seeding */
	capturable  int /* ni-char contents are capturable */
	defense     int /* fighter defense rating */
	fly_cap     int
	is_man_item int /* unit is a character like thing */
	land_cap    int
	maintenance int /* Maintenance cost */
	missile     int /* capable of missile attacks? */
	npc_split   int /* Size to "split" at... */
	plural_name string
	prominent   int /* big things that everyone sees */
	ride_cap    int
	trade_good  int /* Is this thing a trade good? & how much*/
	ungiveable  int /* Can't be transferred between nobles. */
	weight      int
	who_has     int /* who has this unique item */

	// appears in the wild as a random encounter.
	// value is actually the NPC_prog.
	wild int

	x_item_magic    *ItemMagic
	x_item_artifact *EntityArtifact /* Eventually replace x_item_magic */
}

func (ei *entity_item) IsZero() bool {
	// https://freshman.tech/snippets/go/check-empty-struct/
	return zero_check(ei)
}

// todo: must move id into EntityArtifact
func (ei *entity_item) ToEntityArtifact(id int) *EntityArtifact {
	if ei == nil || ei.x_item_artifact == nil || ei.x_item_artifact.IsZero() {
		return nil
	}
	if ei.x_item_artifact.Id != id {
		panic(fmt.Sprintf("assert(%d == %d)", ei.x_item_artifact.Id, id))
	}
	return ei.x_item_artifact
}

// todo: must move id into entity_item
func (ei *entity_item) ToEntityItem(id int) *EntityItem {
	if ei == nil || ei.IsZero() {
		return nil
	}
	return &EntityItem{
		Id:          id,
		Animal:      ei.animal == TRUE,
		AnimalPart:  ei.animal_part,
		Attack:      ei.attack,
		BasePrice:   ei.base_price,
		Capturable:  ei.capturable == TRUE,
		Defense:     ei.defense,
		FlyCap:      ei.fly_cap,
		IsManItem:   ei.is_man_item,
		LandCap:     ei.land_cap,
		Maintenance: ei.maintenance,
		Missile:     ei.missile,
		NpcSplit:    ei.npc_split,
		PluralName:  ei.plural_name,
		Prominent:   ei.prominent == TRUE,
		RideCap:     ei.ride_cap,
		TradeGood:   ei.trade_good == TRUE,
		Ungiveable:  ei.ungiveable == TRUE,
		Weight:      ei.weight,
		WhoHas:      ei.who_has,
		Wild:        ei.wild == TRUE,
	}
}

func (ei *entity_item) ToItemMagic(id int) *ItemMagic {
	if ei == nil || ei.x_item_magic == nil || ei.x_item_magic.IsZero() {
		return nil
	}
	if ei.x_item_magic.Id != id {
		panic(fmt.Sprintf("assert(%d == %d)", ei.x_item_magic.Id, id))
	}
	return ei.x_item_magic
}

type entity_loc struct {
	prov_dest       []int /* cached exits for flood fills */
	near_grave      int   /* nearest graveyard */
	shroud          int   /* magical scry shroud */
	barrier         int   /* magical barrier */
	tax_rate        int   /* Tax rate for this loc. */
	recruited       int   /* How many recruited this month */
	hidden          int   /* is location hidden? */
	dist_from_sea   int   /* provinces to sea province */
	dist_from_swamp int
	dist_from_gate  int
	sea_lane        int /* fast ocean travel here, also "tracks" for npc ferries */
	// effects []*effect        /* ilist of effects on location */ // not used?
	mine_info *entity_mine /* If there's a mine. */
	// location control -- need two so that we can only change fees at the end of the month.
	control  loc_control_ent
	control2 loc_control_ent // doesn't need to be saved
}

type entity_mine struct {
	mc      [MINE_MAX]mine_contents // todo: shouldn't be static size
	shoring [MINE_MAX]int           // todo: shouldn't be static size
}

type entity_misc struct {
	display         string /* entity display banner */
	npc_created     int    /* turn peasant mob created */
	npc_home        int    /* where npc was created */
	npc_cookie      int    /* allocation cookie item for us */
	summoned_by     int    /* who summoned us? */
	save_name       string /* orig name of noble for dead bodies */
	old_lord        int    /* who did this dead body used to belong to */
	npc_memory      sparse
	only_vulnerable int  /* only defeatable with this rare artifact */
	garr_castle     int  /* castle which owns this garrison */
	border_open     int  /* Is the garrison's border open? */
	bind_storm      int  /* storm bound to this ship */
	storm_str       int  /* storm strength */
	npc_dir         int  /* last direction npc moved */
	mine_delay      int  /* time until collapsed mine vanishes */
	cmd_allow       byte /* unit under restricted control */

	// not saved:
	opium_double schar    // improved opium production
	post_txt     []string // text of posted sign
	storm_move   int      // next loc storm will move to
	garr_watch   ints_l    // units garrison watches for
	garr_host    ints_l    // units garrison will attack
	garr_tax     int      // garrison taxes collected
	garr_forward int      // garrison taxes forwarded
}

// nations structure
type entity_nation struct {
	name              string /* Name of the nation, e.g., Mandor Empire */
	citizen           string /* Name of a citizen, e.g., Mandorean */
	nobles            int    /* Total nobles */
	nps               int    /* Total NPS */
	gold              int    /* Total gold */
	players           int    /* Num players */
	win               int    /* Win? */
	proscribed_skills []int  /* Skills you can't have... */
	player_limit      int    /* Limit to # of players */
	capital           int    /* Capital city. */
	jump_start        int    /* Jump start points to start */
	neutral           bool   /* Can't capture/lose NPs. */
}

type EntityPlayer struct {
	Admits        admit_l        `json:"admits,omitempty"`          // Admit permissions list
	BrokenMailer  int            `json:"broken-mailer,omitempty"`   // quote begin lines
	CompuServe    bool           `json:"compuserve,omitempty"`      // get Times from CIS
	DBPath        string         `json:"db-path,omitempty"`         // external path for HTML
	DontRemind    int            `json:"dont-remind,omitempty"`     // don't send a reminder
	EMail         string         `json:"e-mail,omitempty"`          //
	FirstTower    int            `json:"first-tower,omitempty"`     // has player built first tower yet?
	FirstTurn     int            `json:"first-turn,omitempty"`      // which turn was their first?
	Format        int            `json:"format,omitempty"`          // turn report formatting control
	FullName      string         `json:"full-name,omitempty"`       //
	JumpStart     int            `json:"jump-start,omitempty"`      // Jump start points
	Known         sparse         `json:"known,omitempty"`           // visited, lore seen, encountered
	LastOrderTurn int            `json:"last-order-turn,omitempty"` // last turn orders were submitted
	Magic         int            `json:"magic,omitempty"`           // MUs or Priests?
	Nation        int            `json:"nation,omitempty"`          // Player's Nation
	NationList    int            `json:"nation-list,omitempty"`     // Receive the Nation mailing list?
	NoblePoints   int            `json:"noble-points,omitempty"`    // how many NP's the player has
	NoTab         bool           `json:"no-tab,omitempty"`          // player can't tolerate tabs
	Orders        []*orders_list `json:"orders,omitempty"`          // list of Orders for units in this faction
	Password      string         `json:"password,omitempty"`        //
	RulesPath     string         `json:"rules-path,omitempty"`      // external path for HTML
	SentOrders    int            `json:"sent-orders,omitempty"`     // sent in Orders this turn?
	SplitBytes    int            `json:"split-bytes,omitempty"`     // split mail at this many bytes
	SplitLines    int            `json:"split-lines,omitempty"`     // split mail at this many lines
	Unformed      ints_l         `json:"unformed,omitempty"`        // nobles as yet Unformed
	Units         ints_l         `json:"units,omitempty"`           // what Units are in our faction?
	VisEMail      string         `json:"vis-e-mail,omitempty"`      // address to put in player list

	// not saved:
	cmdCount      int    // count of cmds started this turn
	deliverLore   ints_l // show these to player
	locs          sparse // locs we touched
	npGained      int    // np's added this turn
	npSpent       int    // np's lost this turn
	output        sparse // Units with output
	swearThisTurn int    // have we used SWEAR this turn?
	timesPaid     int    // Times paid this month?
	weatherSeen   sparse // locs we've viewed the weather
}

func (p *EntityPlayer) IsZero() bool {
	// https://freshman.tech/snippets/go/check-empty-struct/
	return zero_check(p)
}

// religion sub-structure
type entity_religion_skill struct {
	name        string /* Of the god.  */
	strength    int    /* Related strength skill */
	weakness    int    /* Related weakness skill */
	plant       int    /* The holy plant. */
	animal      int    /* The holy animal. */
	terrain     int    /* Holy terrain. */
	high_priest int    /* The high priest */
	bishops     [2]int /* The two bishops */
}

// ship descriptor
type entity_ship struct {
	hulls      int /* Various ship parts */
	forts      int
	sails      int
	ports      int
	keels      int
	galley_ram int /* galley is fitted with a ram */
}

// the specific religion skills require quite a bit of addition to the basic entity_skill structure.
// this stuff is captured by making a new skill subkind ("religion") and putting in a pointer to a "religion" structure.
type entity_skill struct {
	time_to_learn  int   /* days of study req'd to learn skill */
	time_to_use    int   /* days to use this skill */
	flags          int   /* Flags such as IS_POLLED, REQ_HOLY_SYMBOL, etc. */
	required_skill int   /* skill required to learn this skill */
	np_req         int   /* noble points required to learn */
	offered        []int /* skills learnable after this one */
	research       []int /* skills researable with this one */
	guild          []int /* skills offered if you're a guild member. */

	req            []*req_ent             /* ilist of items required for use or cast */
	produced       int                    /* simple production skill result */
	no_exp         int                    /* this skill not rated for experience */
	practice_cost  int                    /* cost to practice this skill. */
	practice_time  int                    /* time to practice this skill. */
	practice_prog  int                    /* A day longer to practice each N experience levels. */
	religion_skill *entity_religion_skill /* Possible religion pointer. */
	piety          int                    /* Piety required to cast religious skill. Or aura to cast magic skill; not used for all. */
	// not saved:
	use_count    int /* times skill used during turn */
	last_use_who int /* who last used the skill (this turn) */
}

type entity_subloc struct {
	teaches    ints_l /* skills location offers */
	opium_econ int    /* addiction level of city */
	defense    int    /* defense rating of structure */

	loot   int /* loot & pillage level */
	hp     int /* "hit points" */
	damage int /* 0=none, hp=fully destroyed */
	moat   int /* Has a moat? */
	//short shaft_depth;		/* depth of mine shaft */

	builds entity_build_l /* Ongoing builds here. */

	moving int          /* daystamp of beginning of movement */
	x_ship *entity_ship /* Maybe a ship? */

	//int capacity;			/* capacity of ship */
	//schar galley_ram;		/* galley is fitted with a ram */

	near_cities ints_l /* cities rumored to be nearby */
	safe        bool   /* safe haven */
	major       int    /* major city */
	prominence  int    /* prominence of city */

	//schar link_when;		/* month link is open, -1 = never */
	//schar link_open;		/* link is open now */

	link_to   ints_l /* where we are linked to */
	link_from ints_l /* where we are linked from */

	bound_storms ints_l /* storms bound to this ship */

	guild int /* what skill, if a sub_guild */

	//struct effect **effects;        /* ilist of effects on sub-location */

	tax_market  int /* Market tax rate. */
	tax_market2 int /* Temporary until end of month */

	/* Location control -- either here or loc */
	control       loc_control_ent
	control2      loc_control_ent
	entrance_size int /* Size of entrance to subloc */
}

type harvest struct {
	item      int
	skill     int
	worker    int
	chance    int // chance to get one each day, if nonzero
	got_em    string
	none_now  string
	none_ever string
	task_desc string
	public    int // 3rd party view, yes/no
	piety     int // does it use piety?
}

type item_ent struct {
	item int
	qty  int
}

// todo: carry Id as part of struct
type ItemMagic struct {
	Id           int    `json:"id,omitempty"`            //
	AttackBonus  int    `json:"attack-bonus,omitempty"`  //
	Aura         int    `json:"aura,omitempty"`          // auraculum aura
	AuraBonus    int    `json:"aura-bonus,omitempty"`    //
	CloakCreator int    `json:"cloak-creator,omitempty"` //
	CloakRegion  int    `json:"cloak-region,omitempty"`  //
	Creator      int    `json:"creator,omitempty"`       //
	CurseLoyalty int    `json:"curse-loyalty,omitempty"` // curse noncreator loyalty
	DefenseBonus int    `json:"defense-bonus,omitempty"` //
	Lore         int    `json:"lore,omitempty"`          // deliver this lore for the item
	MayStudy     ints_l `json:"may-study,omitempty"`     // list of skills studying from this
	MayUse       ints_l `json:"may-use,omitempty"`       // list of usable skills via this
	MissileBonus int    `json:"missile-bonus,omitempty"` //
	OrbUseCount  int    `json:"orb-use-count,omitempty"` // how many uses left in the orb
	ProjectCast  int    `json:"project-cast,omitempty"`  // stored projected cast
	QuickCast    int    `json:"quick-cast,omitempty"`    // stored quick cast
	Religion     int    `json:"religion,omitempty"`      // Might be a religious artifact
	TokenNI      int    `json:"token-ni,omitempty"`      // ni for controlled npc units
	TokenNum     int    `json:"token-num,omitempty"`     // how many token controlled units
	UseKey       int    `json:"use-key,omitempty"`       // special use action

	// not saved:
	one_turn_use schar /* flag for one use per turn */
}

func (im *ItemMagic) IsZero() bool {
	// https://freshman.tech/snippets/go/check-empty-struct/
	return zero_check(im)
}

// location control: whether or not we're open, fees.
type loc_control_ent struct {
	closed bool
	nobles int // fee per noble
	men    int // fee per person
	weight int // fee per measure of weight
}

type loc_info struct {
	where     int
	here_list []int
}

type make_ struct {
	item        int
	inp1        int
	inp1_factor int // # of inp1 needed to make 1
	inp2        int
	inp2_factor int // # of inp2 needed to make 1
	inp3        int
	inp3_factor int // # of inp3 needed to make 1
	req_skill   int
	worker      int // worker needed
	got_em      string
	public      int // does everyone see us make this
	where       int // place required for production
	aura        int // aura per unit required
	factor      int // multiplying qty factor, usually 1
	days        int // days to make each thing
}

type mine_contents struct {
	items item_ent_l /* ilist of items held */
	// iron, gold, mithril, gate_crystals int // not used?
}

type olytime struct {
	day              int /* day of month */
	turn             int /* turn number */
	days_since_epoch int /* days since game begin */
}

func (ot olytime) IsZero() bool {
	return ot.day == 0 && ot.turn == 0 && ot.days_since_epoch == 0
}

func (ot olytime) ToOlyTime() *OlyTime {
	if ot.IsZero() {
		return nil
	}
	return &OlyTime{
		Day:            ot.day,
		Turn:           ot.turn,
		DaysSinceEpoch: ot.days_since_epoch,
	}
}

// this structure holds game "options" for various different flavors of TAG.
type options_struct struct {
	created_at              time.Time // moment the system data file was first created
	updated_at              time.Time // moment the system data file was last updated
	accounting_dir          string    /* Directory to "join" from. */
	accounting_prog         string    /* Path of the accounting program. */
	auto_drop               bool      /* Drop non-responsive players. */
	bottom_piety            int       /* Monthly +piety for everyone else */
	check_balance           int       /* No orders w/o positive balance. */
	claim_give              int       /* Allow putting gold in claim? */
	cpp                     string    /* Path of cpp */
	death_nps               int       /* What NPs get returned at death? */
	free                    bool      /* Don't charge for this game. */
	free_np_limit           int       /* Play for free with this many NPs. */
	full_markets            bool      /* City markets buy wood, etc. */
	guild_teaching          bool      /* Do guilds teach guild skills? */
	head_priest_piety_limit int       /* Head priest limited to head_priest_piety_limit * num_followers */
	html_passwords          string    /* Path to html passwords */
	html_path               string    /* Path to html directories */
	market_age              int       /* Months untouched in market before removal. */
	middle_piety            int       /* Monthly +piety for junior priests */
	min_piety               int       /* Any priest can have this much piety. */
	mp_antipathy            bool      /* Do mages & priests hate each other? */
	num_books               int       /* Number of teaching books in city */
	open_ended              bool      /* No end to game. */
	output_tags             int       /* include <tag> in output */
	piety_limit             int       /* Normal priest limited to piety_limit * num_followers */
	survive_np              bool      /* Does SFW return NPs when forgotten? */
	times_pay               int       /* What the Times pays for an article. */
	top_piety               int       /* Monthly +piety for head priest */
	turn_charge             string    /* How much to charge per turn. */
	turn_limit              int       /* Limit players to a certain # of turns. */
}

type orders_list struct {
	unit int      // unit orders are for
	l    orders_l // ilist of orders for unit
}
type orders_l [][]byte

type req_ent struct {
	item    int /* item required to use */
	qty     int /* quantity required */
	consume int /* REQ_xx */
}

type schar = int

type skill_ent struct {
	skill        int
	days_studied int /* days studied * TOUGH_NUM */
	experience   int /* experience level with skill */
	know         int /* SKILL_xxx */
	// not saved:
	exp_this_month byte /* flag for add_skill_experience() */

}

type sparse = []int

type trade struct {
	kind       int /* BUY or SELL */
	item       int
	qty        int
	cost       int
	cloak      int /* don't reveal identity of trader */
	have_left  int
	month_prod int /* month city produces item */
	who        int /* redundant -- not saved */
	sort       int /* temp key for sorting -- not saved */
	old_qty    int /* qty at beginning of month, for trade goods */
	counter    int /* Counter to age and lose untraded goods. */
}

// traps - encodes some common traps...
type trap_struct struct {
	type_         int    /* Type of trap. */
	religion      int    /* Religion that ignores this trap. */
	num_attacks   int    /* Number of attacks */
	attack_chance int    /* Chance of attack killing someone */
	name          string /* Name of trap. */
	ignored       string /* Message if you can ignore this trap. */
	flying        string /* Message if you fly over the trap. */
	attack        string /* Message if it attacks you. */
}

type wait_arg struct {
	tag  int
	a1   int
	a2   int
	flag string
}

type uchar = byte

// how long a command has been running
func command_days(c *command) int { return c.days_executing }

// todo: is the "len - 1" right?
func numargs(c *command) int   { return len(c.parse) - 1 }
func oly_month(a *olytime) int { return (a.turn - 1) % NUM_MONTHS }
func oly_year(a *olytime) int  { return (a.turn - 1) / NUM_MONTHS }

// malloc-on-demand substructure references
func p_char(n int) *entity_char {
	if bx[n].x_char == nil {
		bx[n].x_char = &entity_char{}
	}
	return bx[n].x_char
}
func p_command(n int) *command {
	if bx[n].cmd == nil {
		bx[n].cmd = &command{}
	}
	return bx[n].cmd
}
func p_disp(n int) *att_ent {
	if bx[n].x_disp == nil {
		bx[n].x_disp = &att_ent{}
	}
	return bx[n].x_disp
}
func p_gate(n int) *entity_gate {
	if bx[n].x_gate == nil {
		bx[n].x_gate = &entity_gate{}
	}
	return bx[n].x_gate
}
func p_item(n int) *entity_item {
	if bx[n].x_item == nil {
		bx[n].x_item = &entity_item{}
	}
	return bx[n].x_item
}
func p_item_artifact(n int) *EntityArtifact {
	if p_item(n).x_item_artifact == nil {
		p_item(n).x_item_artifact = &EntityArtifact{}
	}
	return p_item(n).x_item_artifact
}
func p_item_magic(n int) *ItemMagic {
	if p_item(n).x_item_magic == nil {
		p_item(n).x_item_magic = &ItemMagic{}
	}
	return p_item(n).x_item_magic
}
func p_loc(n int) *entity_loc {
	if bx[n].x_loc == nil {
		bx[n].x_loc = &entity_loc{}
	}
	return bx[n].x_loc
}
func p_loc_info(n int) *loc_info {
	if _, ok := bx[n]; !ok {
		return nil
	}
	return &bx[n].x_loc_info
}
func p_magic(n int) *char_magic {
	if p_char(n).x_char_magic == nil {
		p_char(n).x_char_magic = &char_magic{}
	}
	return p_char(n).x_char_magic
}
func p_misc(n int) *entity_misc {
	if bx[n].x_misc == nil {
		bx[n].x_misc = &entity_misc{}
	}
	return bx[n].x_misc
}
func p_nation(n int) *entity_nation {
	if bx[n].x_nation == nil {
		bx[n].x_nation = &entity_nation{}
	}
	return bx[n].x_nation
}
func p_player(n int) *EntityPlayer {
	if bx[n].x_player == nil {
		bx[n].x_player = &EntityPlayer{}
	}
	return bx[n].x_player
}
func p_ship(n int) *entity_ship {
	if p_subloc(n).x_ship == nil {
		p_subloc(n).x_ship = &entity_ship{}
	}
	return p_subloc(n).x_ship
}
func p_skill(n int) *entity_skill {
	if bx[n].x_skill == nil {
		bx[n].x_skill = &entity_skill{}
	}
	return bx[n].x_skill
}
func p_subloc(n int) *entity_subloc {
	if bx[n].x_subloc == nil {
		bx[n].x_subloc = &entity_subloc{}
	}
	return bx[n].x_subloc
}

// "raw" pointers to substructures, may be NULL
func rp_char(n int) *entity_char {
	if box, ok := bx[n]; ok {
		return box.x_char
	}
	return nil
}
func rp_command(n int) *command {
	if box, ok := bx[n]; ok {
		return box.cmd
	}
	return nil
}
func rp_disp(n int) *att_ent {
	if box, ok := bx[n]; ok {
		return box.x_disp
	}
	return nil
}
func rp_gate(n int) *entity_gate {
	if box, ok := bx[n]; ok {
		return box.x_gate
	}
	return nil
}
func rp_item(n int) *entity_item {
	if box, ok := bx[n]; ok {
		return box.x_item
	}
	return nil
}
func rp_item_artifact(n int) *EntityArtifact {
	if rp_item(n) != nil {
		return rp_item(n).x_item_artifact
	}
	return nil
}
func rp_item_magic(n int) *ItemMagic {
	if rp_item(n) != nil {
		return rp_item(n).x_item_magic
	}
	return nil
}
func rp_loc(n int) *entity_loc {
	if box, ok := bx[n]; ok {
		return box.x_loc
	}
	return nil
}
func rp_loc_info(n int) *loc_info {
	if box, ok := bx[n]; ok {
		return &box.x_loc_info
	}
	return nil
}
func rp_magic(n int) *char_magic {
	if rp_char(n) != nil {
		return rp_char(n).x_char_magic
	}
	return nil
}
func rp_misc(n int) *entity_misc {
	if box, ok := bx[n]; ok {
		return box.x_misc
	}
	return nil
}
func rp_nation(n int) *entity_nation {
	if box, ok := bx[n]; ok {
		return box.x_nation
	}
	return nil
}
func rp_player(n int) *EntityPlayer {
	if box, ok := bx[n]; ok {
		return box.x_player
	}
	return nil
}
func rp_relig_skill(n int) *entity_religion_skill {
	if bx[n].x_skill != nil {
		return bx[n].x_skill.religion_skill
	}
	return nil
}
func rp_ship(n int) *entity_ship {
	if rp_subloc(n) != nil {
		return rp_subloc(n).x_ship
	}
	return nil
}
func rp_skill(n int) *entity_skill {
	if box, ok := bx[n]; ok {
		return box.x_skill
	}
	return nil
}
func rp_subloc(n int) *entity_subloc {
	if box, ok := bx[n]; ok {
		return box.x_subloc
	}
	return nil
}

// func	item_creat_loc(n int) int {if rp_item_magic(n) != nil {return rp_item_magic(n).region_created} else {return 0}}
// func	loc_link_open(n int) int {if rp_subloc(n) != nil {return rp_subloc(n).link_open} else {return 0}}
// func	mine_depth(n int) int {if rp_subloc(n) != nil {return rp_subloc(n).shaft_depth / 3} else {return 0}}
// func char_pledge(n int) int {if rp_magic(n) != nil {return rp_magic(n).pledge} else {return 0}}
// func loc_barrier(n int) int {if rp_loc(n) != nil {return rp_loc(n).barrier} else {return 0}}
func banner(n int) string {
	if rp_misc(n) != nil {
		return rp_misc(n).display
	}
	return ""
}
func board_fee(n int) int {
	if rp_magic(n) != nil {
		panic("return rp_magic(n).fee")
	}
	return 0
}
func body_old_lord(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).old_lord
	}
	return 0
}
func border_open(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).border_open
	}
	return 0
}
func char_abil_shroud(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).ability_shroud
	}
	return 0
}
func char_attack(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).attack
	}
	return 0
}
func char_auraculum(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).auraculum
	}
	return 0
}
func char_behind(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).behind
	}
	return 0
}
func char_break(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).break_point
	}
	return 0
}
func char_cur_aura(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).cur_aura
	}
	return 0
}
func char_defense(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).defense
	}
	return 0
}
func char_guard(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).guard
	}
	return 0
}
func char_health(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).health
	}
	return 0
}
func char_hidden(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).hide_self
	}
	return 0
}
func char_hide_mage(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).hide_mage
	}
	return 0
}
func char_max_aura(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).max_aura
	}
	return 0
}
func char_melt_me(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).melt_me
	}
	return 0
}
func char_missile(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).missile
	}
	return 0
}
func char_new_lord(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).new_lord
	}
	return 0
}
func char_persuaded(n int) int {
	if rp_char(n) != nil {
		panic("return rp_char(n).persuaded")
	}
	return 0
}
func char_piety(n int) int {
	if is_priest(n) != 0 && rp_char(n) != nil {
		return rp_char(n).religion.piety
	}
	return 0
}
func char_proj_cast(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).project_cast
	}
	return 0
}
func char_quick_cast(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).quick_cast
	}
	return 0
}
func char_rank(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).rank
	}
	return 0
}
func char_sick(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).sick
	}
	return 0
}
func combat_skill(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).flags & COMBAT_SKILL
	}
	return 0
}
func entrance_size(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).entrance_size
	}
	return 0
}
func ferry_horn(n int) bool {
	if rp_magic(n) != nil {
		return rp_magic(n).ferry_flag
	}
	return false
}
func garrison_castle(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).garr_castle
	}
	return 0
}
func gate_dest(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).to_loc
	}
	return 0
}
func gate_dist(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).dist_from_gate
	}
	return 0
}
func gate_seal(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).seal_key
	}
	return 0
}
func in_faery(n int) bool  { return region(n) == faery_region }
func in_hades(n int) bool  { return region(n) == hades_region }
func in_clouds(n int) bool { return region(n) == cloud_region }
func is_fighter(n int) int {
	if item_attack(n) != 0 || item_defense(n) != 0 || item_missile(n) != 0 || n == item_ghost_warrior {
		return TRUE
	}
	return FALSE
}
func is_magician(n int) bool {
	if rp_magic(n) != nil {
		return rp_magic(n).magician != 0
	}
	return false
}
func item_animal(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).animal
	}
	return 0
}
func item_animal_part(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).animal_part
	}
	return 0
}
func item_attack(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).attack
	}
	return 0
}
func item_attack_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).AttackBonus
	}
	return 0
}
func item_aura(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).Aura
	}
	return 0
}
func item_aura_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).AuraBonus
	}
	return 0
}
func item_capturable(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).capturable
	}
	return 0
}
func item_creat_cloak(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).CloakCreator
	}
	return 0
}
func item_creator(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).Creator
	}
	return 0
}
func item_curse_non(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).CurseLoyalty
	}
	return 0
}
func item_defense(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).defense
	}
	return 0
}
func item_defense_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).DefenseBonus
	}
	return 0
}
func item_lore(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).Lore
	}
	return 0
}
func item_missile(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).missile
	}
	return 0
}
func item_missile_bonus(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).MissileBonus
	}
	return 0
}
func item_prog(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).wild
	}
	return 0
}
func item_prominent(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).prominent
	}
	return 0
}
func item_reg_cloak(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).CloakRegion
	}
	return 0
}
func item_split(n int) int {
	if rp_item(n) != nil {
		if rp_item(n).npc_split != 0 {
			return rp_item(n).npc_split
		}
	}
	return 50
}
func item_token_ni(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).TokenNI
	}
	return 0
}
func item_token_num(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).TokenNum
	}
	return 0
}
func item_unique(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).who_has
	}
	return 0
}
func item_use_key(n int) int {
	if rp_item_magic(n) != nil {
		return rp_item_magic(n).UseKey
	}
	return 0
}
func item_wild(n int) int {
	if rp_item(n) != nil {
		return rp_item(n).wild
	}
	return 0
}
func learn_time(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).time_to_learn
	}
	return 0
}
func loc_barrier(n int) int { return get_effect(n, ef_magic_barrier, 0, 0) }
func loc_damage(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).damage
	}
	return 0
}
func loc_defense(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).defense
	}
	return 0
}
func loc_hp(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).hp
	}
	return 0
}
func loc_moat(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).moat
	}
	return 0
}
func loc_opium(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).opium_econ
	}
	return 0
}
func loc_pillage(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).loot
	}
	return 0
}
func loc_prominence(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).prominence
	}
	return 0
}
func loc_sea_lane(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).sea_lane
	}
	return 0
}
func loc_shroud(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).shroud
	}
	return 0
}
func loyal_kind(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).loy_kind
	}
	return 0
}
func loyal_rate(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).loy_rate
	}
	return 0
}
func major_city(n int) int {
	if rp_subloc(n) != nil {
		return rp_subloc(n).major
	}
	return 0
}
func man_item(n int) int {
	// todo: this feels wrong
	if rp_item(n) != nil {
		return rp_item(n).is_man_item
	}
	return 0
}
func no_barrier(n int) int { return get_effect(n, ef_inhibit_barrier, 0, 0) }
func noble_item(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).unit_item
	} else {
		return 0
	}
}
func npc_last_dir(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).npc_dir
	}
	return 0
}
func npc_program(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).npc_prog
	}
	return 0
}
func npc_summoner(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).summoned_by
	}
	return 0
}
func only_defeatable(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).only_vulnerable
	}
	return 0
}
func our_token(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).token
	}
	return 0
}
func personal_break(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).personal_break_point
	}
	return 0
}
func player_broken_mailer(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).BrokenMailer
	}
	return 0
}
func player_compuserve(n int) bool {
	if rp_player(n) != nil {
		return rp_player(n).CompuServe
	}
	return false
}
func player_email(n int) string {
	if rp_player(n) != nil {
		return rp_player(n).EMail
	}
	return ""
}
func player_format(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).Format
	}
	return 0
}
func player_js(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).JumpStart
	}
	return 0
}
func player_notab(n int) bool {
	if rp_player(n) != nil {
		return rp_player(n).NoTab
	}
	return false
}
func player_np(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).NoblePoints
	}
	return 0
}
func player_split_bytes(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).SplitBytes
	}
	return 0
}
func player_split_lines(n int) int {
	if rp_player(n) != nil {
		return rp_player(n).SplitLines
	}
	return 0
}
func practice_cost(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).practice_cost
	}
	return 0
}
func practice_progressive(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).practice_prog
	}
	return 0
}
func practice_time(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).practice_time
	}
	return 0
}
func recent_pillage(n int) int {
	if rp_subloc(n) != nil {
		panic("return rp_subloc(n).recent_loot")
	}
	return 0
}
func reflect_blast(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).aura_reflect
	}
	return 0
}
func release_swear(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).swear_on_release
	}
	return 0
}
func req_skill(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).required_skill
	}
	return 0
}
func restricted_control(n int) byte {
	if rp_misc(n) != nil {
		return rp_misc(n).cmd_allow
	}
	return 0
}
func road_dest(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).to_loc
	}
	return 0
}
func road_hidden(n int) int {
	if rp_gate(n) != nil {
		return rp_gate(n).road_hidden
	}
	return 0
}
func safe_haven(n int) bool {
	if rp_subloc(n) != nil {
		return rp_subloc(n).safe
	}
	return false
}
func sea_dist(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).dist_from_sea
	}
	return 0
}

// todo: what does this do and why?
func see_all(n int) int {
	return immed_see_all
}
func ship_has_ram(n int) int {
	if rp_ship(n) != nil {
		return rp_ship(n).galley_ram
	}
	return 0
}
func skill_exp(who, sk int) int {
	if rp_skill_ent(who, sk) != nil {
		return rp_skill_ent(who, sk).experience
	}
	return 0
}
func skill_flags(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).flags
	}
	return 0
}
func skill_no_exp(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).no_exp
	}
	return 0
}
func skill_np_req(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).np_req
	}
	return 0
}
func skill_produce(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).produced
	}
	return 0
}
func skill_time_to_use(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).time_to_use
	}
	return 0
}
func storm_bind(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).bind_storm
	}
	return 0
}
func storm_strength(n int) int {
	if rp_misc(n) != nil {
		return rp_misc(n).storm_str
	}
	return 0
}
func swamp_dist(n int) int {
	if rp_loc(n) != nil {
		return rp_loc(n).dist_from_swamp
	}
	return 0
}

func times_paid(n int) bool {
	if rp_player(n) != nil {
		return rp_player(n).timesPaid != FALSE
	}
	return false
}

func weather_mage(n int) int {
	if rp_magic(n) != nil {
		return rp_magic(n).knows_weather
	}
	return 0
}

func is_loc_or_ship(n int) bool { return (kind(n) == T_loc || kind(n) == T_ship) }
func is_ship(n int) bool {
	return (subkind(n) == sub_galley || subkind(n) == sub_roundship || subkind(n) == sub_raft || subkind(n) == sub_ship)
}
func is_ship_notdone(n int) bool {
	return (subkind(n) == sub_galley_notdone || subkind(n) == sub_roundship_notdone || subkind(n) == sub_raft_notdone || subkind(n) == sub_ship_notdone)
}
func is_ship_either(n int) bool { return (is_ship(n) || is_ship_notdone(n)) }
func kind(n int) schar {
	if box, ok := bx[n]; ok {
		return box.kind
	}
	return T_deleted
}
func kind_first(n int) int { return (box_head[(n)]) }
func kind_next(n int) int  { return (bx[(n)].x_next_kind) }

// return exactly where a unit is.
// May point to another character, a structure, or a region.
func loc(n int) int       { return rp_loc_info(n).where }
func sub_first(n int) int { return (sub_head[(n)]) }
func sub_next(n int) int {
	if bx[n] != nil {
		return bx[n].x_next_sub
	}
	return 0
}
func subkind(n int) schar {
	if bx[n] != nil {
		return bx[n].skind
	}
	return 0
}
func valid_box(n int) bool {
	if box, ok := bx[n]; ok {
		return box.kind != T_deleted
	}
	return false
}

// guild stuff
func is_guild(n int) int {
	if subkind(n) == sub_guild && rp_subloc(n).guild != 0 {
		return rp_subloc(n).guild
	}
	return 0
}
func guild_member(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).guild
	}
	return 0
}

// religion stuff
func god_name(n int) string { return rp_relig_skill(n).name }
func holy_plant(n int) int {
	if is_priest(n) != 0 {
		return rp_relig_skill(is_priest(n)).plant
	}
	return 0
}
func is_follower(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).religion.priest
	}
	return 0
}
func is_priest(n int) int { return skill_school(has_subskill((n), sub_religion)) }
func is_temple(n int) int {
	if subkind(n) == sub_temple && rp_subloc(n).guild != 0 {
		return rp_subloc(n).guild
	}
	return 0
}
func is_wizard(n int) int { return has_subskill((n), sub_magic) }
func skill_aura(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).piety
	}
	return 0
}
func skill_piety(n int) int {
	if rp_skill(n) != nil {
		return rp_skill(n).piety
	}
	return 0
}

// movement tests
// _moving indicates that the unit has initiated movement
// _gone indicates that the unit has actually left the locations, and should not be interacted with anymore.
//
// The distinction allows zero time commands to interact with the entity on the day movement is begun.
func char_moving(n int) int {
	if rp_char(n) == nil {
		return 0
	}
	return rp_char(n).moving
}
func char_gone(n int) int {
	if char_moving(n) == 0 {
		return 0
	} else if evening {
		return sysclock.days_since_epoch - char_moving(n) + 1
	}
	return sysclock.days_since_epoch - char_moving(n)
}
func ship_moving(n int) int {
	if rp_subloc(n) == nil {
		return 0
	}
	return rp_subloc(n).moving
}
func ship_gone(n int) int {
	if ship_moving(n) == 0 {
		return 0
	} else if evening {
		return sysclock.days_since_epoch - ship_moving(n) + 1
	}
	return sysclock.days_since_epoch - ship_moving(n)
}

func add_s(n int) string {
	if n == 0 {
		return ""
	}
	return "s"
}

func add_ds(n int) string { panic(`n, ((n) == 1 ? "" : "s"`) }
func alive(n int) bool    { return kind(n) == T_char }
func char_alone(n int) bool {
	// alone except for angels & ninjas
	return count_stack_any_real(n, true, true) == 1
}
func char_alone_stealth(n int) bool {
	// alone except for ninjas (permit angels)
	return count_stack_any_real(n, true, true) == 1
}
func char_really_alone(n int) bool {
	// alone counting everything
	return count_stack_any_real(n, false, false) == 1
}
func char_really_hidden(n int) bool {
	// hidden & alone (w/ angels & ninjas)
	return (char_hidden(n) != 0 && char_alone_stealth(n) && !is_prisoner(n))
}
func is_npc(n int) bool {
	// subkind(n) == sub_ni || loyal_kind(n) == LOY_npc)
	return subkind(n) != 0 || loyal_kind(n) == LOY_npc || loyal_kind(n) == LOY_summon
}
func is_prisoner(n int) bool {
	if rp_char(n) == nil {
		return false
	}
	return rp_char(n).prisoner != 0
}
func is_real_npc(n int) bool    { return player(n) < 1000 }
func magic_skill(n int) bool    { return subkind(skill_school(n)) == sub_magic }
func refugee(n int) bool        { return nation(n) == 0 && subkind(n) != sub_ni }
func religion_skill(n int) bool { return subkind(skill_school(n)) == sub_religion }
func wait(n int) int {
	if rp_command(n) != nil {
		return rp_command(n).wait
	}
	return 0
}
func will_pay(n int) int {
	if rp_char(n) != nil {
		return rp_char(n).pay
	}
	return 0
}
