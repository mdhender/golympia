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

package main

import "github.com/mdhender/golympia/pkg/io"

type command struct{}
type harvest struct{}
type item_ent struct{}
type loc_control_ent struct{}
type make_ struct{}
type skill_ent struct{}
type sparse []int
type trade struct{}

func MM(item int) int                                                        { panic("!implemented") }
func add_aura(who int, aura int)                                             { panic("!implemented") }
func add_new_players()                                                       { panic("!implemented") }
func add_piety(who int, amount int, use_limit int)                           { panic("!implemented") }
func add_scrolls(where int)                                                  { panic("!implemented") }
func add_skill_experience(who int, skill int)                                { panic("!implemented") }
func add_unformed_sup(pl int)                                                { panic("!implemented") }
func alert_palantir_scry(who int, target int)                                { panic("!implemented") }
func alert_scry_attempt(who int, target int, s string)                       { panic("!implemented") }
func alert_scry_generic(who int, where int)                                  { panic("!implemented") }
func all_stack(who int, l *ilist)                                            { panic("!implemented") }
func alone_here(who int) int                                                 { panic("!implemented") }
func artifact_identify(header string, c *command) int                        { panic("!implemented") }
func auto_hades()                                                            { panic("!implemented") }
func auto_savage(who int)                                                    { panic("!implemented") }
func auto_undead(who int)                                                    { panic("!implemented") }
func autocharge(who int, total int) int                                      { panic("!implemented") }
func beast_wild(who int) int                                                 { panic("!implemented") }
func being_taught(who int, a1 int, item *int, bonus *int) int                { panic("!implemented") }
func building_owner(ship int) int                                            { panic("!implemented") }
func calc_beast_limit(who int, boolean int) int                              { panic("!implemented") }
func calc_entrance_fee(control *loc_control_ent, c *command, ruler int) int  { panic("!implemented") }
func calc_man_limit(who int, boolean int) int                                { panic("!implemented") }
func calculate_nation_nps()                                                  { panic("!implemented") }
func call_init_routines()                                                    { panic("!implemented") }
func can_join_guild(who int, school int) int                                 { panic("!implemented") }
func cast_check_char_here(who int, target int) int                           { panic("!implemented") }
func cast_where(who int) int                                                 { panic("!implemented") }
func change_box_subkind(where int, kind int)                                 { panic("!implemented") }
func char_here(who int, a1 int) int                                          { panic("!implemented") }
func char_reclaim(who int)                                                   { panic("!implemented") }
func char_rep_location(who int) string                                       { panic("!implemented") }
func char_rep_sup(i int, n int)                                              { panic("!implemented") }
func character_report()                                                      { panic("!implemented") }
func charge_entrance_fees(who int, owner int, amount int) int                { panic("!implemented") }
func check_all_auto_attacks()                                                { panic("!implemented") }
func check_captain_loses_sailors(qty int, target int, who int)               { panic("!implemented") }
func check_db()                                                              { panic("!implemented") }
func check_ocean_chars()                                                     { panic("!implemented") }
func check_skill(who int, skill int) int                                     { panic("!implemented") }
func check_win_conditions()                                                  { panic("!implemented") }
func clear_all_att(who int)                                                  { panic("!implemented") }
func clear_all_trades(where int)                                             { panic("!implemented") }
func clear_contacts(who int)                                                 { panic("!implemented") }
func clear_guard_flag(loser int)                                             { panic("!implemented") }
func clear_second_waits()                                                    { panic("!implemented") }
func clear_wait_parse(c *command)                                            { panic("!implemented") }
func cloak_lord(who int) int                                                 { panic("!implemented") }
func close_logfile()                                                         { panic("!implemented") }
func cmd_shift(c *command)                                                   { panic("!implemented") }
func combat_artifact_bonus(unit int, n int, item *int) int                   { panic("!implemented") }
func combat_comp(q1, q2 interface{}) int                                     { panic("!implemented") }
func command_done(c *command)                                                { panic("!implemented") }
func compute_dist()                                                          { panic("!implemented") }
func connect_locations(loc1, dir1, loc2, dir2 int)                           { panic("!implemented") }
func contacted(target int, who int) int                                      { panic("!implemented") }
func controlled_humans_here(province int) int                                { panic("!implemented") }
func controls_loc(destination int) int                                       { panic("!implemented") }
func count_any(unit int, b int, c int) int                                   { panic("!implemented") }
func count_any_real(unit int, b int, c int) int                              { panic("!implemented") }
func count_fighters(who int, item int) int                                   { panic("!implemented") }
func count_stack_any(who int) int                                            { panic("!implemented") }
func count_stack_any_real(who int, ignore_ninjas int, ignore_angels int) int { panic("!implemented") }
func create_cloudlands()                                                     { panic("!implemented") }
func create_faery()                                                          { panic("!implemented") }
func create_hades()                                                          { panic("!implemented") }
func create_monster_stack(item int, qty int, where int) int                  { panic("!implemented") }
func create_new_beasts(where int, sk int) int                                { panic("!implemented") }
func create_peasant_mob(where int) int                                       { panic("!implemented") }
func create_random_artifact(monster int) int                                 { panic("!implemented") }
func crosses_ocean(target int, who int) int                                  { panic("!implemented") }
func d_add_forts(c *command) int                                             { panic("!implemented") }
func d_add_iron_shoring(c *command) int                                      { panic("!implemented") }
func d_add_keels(c *command) int                                             { panic("!implemented") }
func d_add_ports(c *command) int                                             { panic("!implemented") }
func d_add_ram(c *command) int                                               { panic("!implemented") }
func d_add_sails(c *command) int                                             { panic("!implemented") }
func d_add_wooden_shoring(c *command) int                                    { panic("!implemented") }
func d_adv_med(c *command) int                                               { panic("!implemented") }
func d_archery(c *command) int                                               { panic("!implemented") }
func d_arrange_mugging(c *command) int                                       { panic("!implemented") }
func d_assassinate(c *command) int                                           { panic("!implemented") }
func d_aura_blast(c *command) int                                            { panic("!implemented") }
func d_banish_corpses(c *command) int                                        { panic("!implemented") }
func d_banish_undead(c *command) int                                         { panic("!implemented") }
func d_bar_loc(c *command) int                                               { panic("!implemented") }
func d_bind_storm(c *command) int                                            { panic("!implemented") }
func d_bird_spy(c *command) int                                              { panic("!implemented") }
func d_bless_fort(c *command) int                                            { panic("!implemented") }
func d_breed(c *command) int                                                 { panic("!implemented") }
func d_brew_death(c *command) int                                            { panic("!implemented") }
func d_brew_fiery(c *command) int                                            { panic("!implemented") }
func d_brew_heal(c *command) int                                             { panic("!implemented") }
func d_brew_slave(c *command) int                                            { panic("!implemented") }
func d_brew_weightlessness(c *command) int                                   { panic("!implemented") }
func d_bribe(c *command) int                                                 { panic("!implemented") }
func d_build_wagons(c *command) int                                          { panic("!implemented") }
func d_calm_peasants(c *command) int                                         { panic("!implemented") }
func d_capture_beasts(c *command) int                                        { panic("!implemented") }
func d_cloak_creat(c *command) int                                           { panic("!implemented") }
func d_cloak_reg(c *command) int                                             { panic("!implemented") }
func d_conceal_arts(c *command) int                                          { panic("!implemented") }
func d_conceal_location(c *command) int                                      { panic("!implemented") }
func d_conceal_mine(c *command) int                                          { panic("!implemented") }
func d_conceal_nation(c *command) int                                        { panic("!implemented") }
func d_create_dirt_golem(c *command) int                                     { panic("!implemented") }
func d_create_flesh_golem(c *command) int                                    { panic("!implemented") }
func d_create_holy_symbol(c *command) int                                    { panic("!implemented") }
func d_create_iron_golem(c *command) int                                     { panic("!implemented") }
func d_create_mist(c *command) int                                           { panic("!implemented") }
func d_create_mithril(c *command) int                                        { panic("!implemented") }
func d_create_ninja(c *command) int                                          { panic("!implemented") }
func d_curse_noncreat(c *command) int                                        { panic("!implemented") }
func d_death_fog(c *command) int                                             { panic("!implemented") }
func d_decrease_demand(c *command) int                                       { panic("!implemented") }
func d_decrease_supply(c *command) int                                       { panic("!implemented") }
func d_dedicate_temple(c *command) int                                       { panic("!implemented") }
func d_dedicate_tower(c *command) int                                        { panic("!implemented") }
func d_defense(c *command) int                                               { panic("!implemented") }
func d_destroy_art(c *command) int                                           { panic("!implemented") }
func d_detect_abil(c *command) int                                           { panic("!implemented") }
func d_detect_arts(c *command) int                                           { panic("!implemented") }
func d_detect_beasts(c *command) int                                         { panic("!implemented") }
func d_detect_gates(c *command) int                                          { panic("!implemented") }
func d_detect_scry(c *command) int                                           { panic("!implemented") }
func d_dispel_abil(c *command) int                                           { panic("!implemented") }
func d_dispel_region(c *command) int                                         { panic("!implemented") }
func d_dissipate(c *command) int                                             { panic("!implemented") }
func d_draw_crowds(c *command) int                                           { panic("!implemented") }
func d_eat_dead(c *command) int                                              { panic("!implemented") }
func d_edge_of_kireus(c *command) int                                        { panic("!implemented") }
func d_enchant_guard(c *command) int                                         { panic("!implemented") }
func d_fierce_wind(c *command) int                                           { panic("!implemented") }
func d_find_all_hidden_features(c *command) int                              { panic("!implemented") }
func d_find_food(c *command) int                                             { panic("!implemented") }
func d_find_hidden_features(c *command) int                                  { panic("!implemented") }
func d_find_rich(c *command) int                                             { panic("!implemented") }
func d_forge_art_x(c *command) int                                           { panic("!implemented") }
func d_forge_aura(c *command) int                                            { panic("!implemented") }
func d_forge_palantir(c *command) int                                        { panic("!implemented") }
func d_fortify_castle(c *command) int                                        { panic("!implemented") }
func d_gather_holy_plant(c *command) int                                     { panic("!implemented") }
func d_generic_trap(c *command) int                                          { panic("!implemented") }
func d_grow_pop(c *command) int                                              { panic("!implemented") }
func d_guard_loyalty(c *command) int                                         { panic("!implemented") }
func d_heal(c *command) int                                                  { panic("!implemented") }
func d_hide(c *command) int                                                  { panic("!implemented") }
func d_hide_item(c *command) int                                             { panic("!implemented") }
func d_hide_money(c *command) int                                            { panic("!implemented") }
func d_hinder_med(c *command) int                                            { panic("!implemented") }
func d_hinder_med_b(c *command) int                                          { panic("!implemented") }
func d_improve_charisma(c *command) int                                      { panic("!implemented") }
func d_improve_fort(c *command) int                                          { panic("!implemented") }
func d_improve_logging(c *command) int                                       { panic("!implemented") }
func d_improve_mining(c *command) int                                        { panic("!implemented") }
func d_improve_opium(c *command) int                                         { panic("!implemented") }
func d_improve_quarrying(c *command) int                                     { panic("!implemented") }
func d_improve_smithing(c *command) int                                      { panic("!implemented") }
func d_improve_taxes(c *command) int                                         { panic("!implemented") }
func d_incite(c *command) int                                                { panic("!implemented") }
func d_increase_demand(c *command) int                                       { panic("!implemented") }
func d_increase_supply(c *command) int                                       { panic("!implemented") }
func d_instill_fanaticism(c *command) int                                    { panic("!implemented") }
func d_keep_savage(c *command) int                                           { panic("!implemented") }
func d_keep_undead(c *command) int                                           { panic("!implemented") }
func d_last_rites(c *command) int                                            { panic("!implemented") }
func d_lead_to_gold(c *command) int                                          { panic("!implemented") }
func d_lightning(c *command) int                                             { panic("!implemented") }
func d_locate_char(c *command) int                                           { panic("!implemented") }
func d_mage_menial(c *command) int                                           { panic("!implemented") }
func d_meditate(c *command) int                                              { panic("!implemented") }
func d_mesmerize_crowd(c *command) int                                       { panic("!implemented") }
func d_mine_gate_crystal(c *command) int                                     { panic("!implemented") }
func d_mine_gold(c *command) int                                             { panic("!implemented") }
func d_mine_iron(c *command) int                                             { panic("!implemented") }
func d_mine_mithril(c *command) int                                          { panic("!implemented") }
func d_moat_castle(c *command) int                                           { panic("!implemented") }
func d_mutate_art(c *command) int                                            { panic("!implemented") }
func d_notify_jump(c *command) int                                           { panic("!implemented") }
func d_notify_unseal(c *command) int                                         { panic("!implemented") }
func d_npc_breed(c *command) int                                             { panic("!implemented") }
func d_obscure_art(c *command) int                                           { panic("!implemented") }
func d_obscure_forest_trail(c *command) int                                  { panic("!implemented") }
func d_obscure_mountain_trail(c *command) int                                { panic("!implemented") }
func d_persuade_oath(c *command) int                                         { panic("!implemented") }
func d_petty_thief(c *command) int                                           { panic("!implemented") }
func d_prep_ritual(c *command) int                                           { panic("!implemented") }
func d_proj_cast(c *command) int                                             { panic("!implemented") }
func d_protect_mine(c *command) int                                          { panic("!implemented") }
func d_quick_cast(c *command) int                                            { panic("!implemented") }
func d_raise(c *command) int                                                 { panic("!implemented") }
func d_rally(c *command) int                                                 { panic("!implemented") }
func d_raze(c *command) int                                                  { panic("!implemented") }
func d_recruit_elves(c *command) int                                         { panic("!implemented") }
func d_rem_art_cloak(c *command) int                                         { panic("!implemented") }
func d_rem_seal(c *command) int                                              { panic("!implemented") }
func d_remove_bless(c *command) int                                          { panic("!implemented") }
func d_remove_forts(c *command) int                                          { panic("!implemented") }
func d_remove_keels(c *command) int                                          { panic("!implemented") }
func d_remove_ports(c *command) int                                          { panic("!implemented") }
func d_remove_ram(c *command) int                                            { panic("!implemented") }
func d_remove_sails(c *command) int                                          { panic("!implemented") }
func d_renew_storm(c *command) int                                           { panic("!implemented") }
func d_resurrect(c *command) int                                             { panic("!implemented") }
func d_reveal_arts(c *command) int                                           { panic("!implemented") }
func d_reveal_key(c *command) int                                            { panic("!implemented") }
func d_reveal_mage(c *command) int                                           { panic("!implemented") }
func d_reveal_vision(c *command) int                                         { panic("!implemented") }
func d_sail(c *command) int                                                  { panic("!implemented") }
func d_save_proj(c *command) int                                             { panic("!implemented") }
func d_save_quick(c *command) int                                            { panic("!implemented") }
func d_scry_region(c *command) int                                           { panic("!implemented") }
func d_seal_gate(c *command) int                                             { panic("!implemented") }
func d_seize_storm(c *command) int                                           { panic("!implemented") }
func d_show_art_creat(c *command) int                                        { panic("!implemented") }
func d_show_art_reg(c *command) int                                          { panic("!implemented") }
func d_shroud_abil(c *command) int                                           { panic("!implemented") }
func d_shroud_region(c *command) int                                         { panic("!implemented") }
func d_smuggle_goods(c *command) int                                         { panic("!implemented") }
func d_smuggle_men(c *command) int                                           { panic("!implemented") }
func d_sneak(c *command) int                                                 { panic("!implemented") }
func d_spy_inv(c *command) int                                               { panic("!implemented") }
func d_spy_lord(c *command) int                                              { panic("!implemented") }
func d_spy_skills(c *command) int                                            { panic("!implemented") }
func d_strengthen_castle(c *command) int                                     { panic("!implemented") }
func d_study(c *command) int                                                 { panic("!implemented") }
func d_summon_fog(c *command) int                                            { panic("!implemented") }
func d_summon_rain(c *command) int                                           { panic("!implemented") }
func d_summon_wind(c *command) int                                           { panic("!implemented") }
func d_swordplay(c *command) int                                             { panic("!implemented") }
func d_tap_health(c *command) int                                            { panic("!implemented") }
func d_teach(c *command) int                                                 { panic("!implemented") }
func d_teleport_item(c *command) int                                         { panic("!implemented") }
func d_torture(c *command) int                                               { panic("!implemented") }
func d_trance(c *command) int                                                { panic("!implemented") }
func d_unbar_loc(c *command) int                                             { panic("!implemented") }
func d_undead_lord(c *command) int                                           { panic("!implemented") }
func d_unobscure_art(c *command) int                                         { panic("!implemented") }
func d_unseal_gate(c *command) int                                           { panic("!implemented") }
func d_urchin_spy(c *command) int                                            { panic("!implemented") }
func d_use_item(c *command) int                                              { panic("!implemented") }
func d_view_aura(c *command) int                                             { panic("!implemented") }
func d_vision_reg(c *command) int                                            { panic("!implemented") }
func d_weaken_fort(c *command) int                                           { panic("!implemented") }
func d_widen_entrance(c *command) int                                        { panic("!implemented") }
func d_write_spell(c *command) int                                           { panic("!implemented") }
func daily_events()                                                          { panic("!implemented") }
func delete_all_effects(what int, type_ int, sk int)                         { panic("!implemented") }
func delete_city_trade(where int, item int)                                  { panic("!implemented") }
func deliver_lore(pl int, a int)                                             { panic("!implemented") }
func deserted_s(n int) string                                                { panic("!implemented") }
func determine_output_order()                                                { panic("!implemented") }
func dissipate_storm(i int, boolean int)                                     { panic("!implemented") }
func do_command(c *command)                                                  { panic("!implemented") }
func do_db_url(c *command) int                                               { panic("!implemented") }
func do_npc_orders(new_ int, j int, k int)                                   { panic("!implemented") }
func do_production(i int, boolean int)                                       { panic("!implemented") }
func do_rules_url(c *command) int                                            { panic("!implemented") }
func do_times()                                                              { panic("!implemented") }
func do_wild_hunt(where int)                                                 { panic("!implemented") }
func eat_loop(mail_now int)                                                  { panic("!implemented") }
func effective_workers(who int) int                                          { panic("!implemented") }
func experience_use_speedup(c *command)                                      { panic("!implemented") }
func extra_item_info(who int, item int, qty int) string                      { panic("!implemented") }
func extract_stacked_unit(who int)                                           { panic("!implemented") }
func fetch_inside_name() string                                              { panic("!implemented") }
func find_command(cmd string) int                                            { panic("!implemented") }
func find_nation(name string) int                                            { panic("!implemented") }
func find_trade(i int, n int, item int) *trade                               { panic("!implemented") }
func find_use_entry(i int) int                                               { panic("!implemented") }
func finish_command(c *command) int                                          { panic("!implemented") }
func first_character(where int) int                                          { panic("!implemented") }
func float_cloudlands()                                                      { panic("!implemented") }
func flush_unit_orders(old_pl int, who int)                                  { panic("!implemented") }
func fort_default_defense(subloc int) int                                    { panic("!implemented") }
func garrison_here(where int) int                                            { panic("!implemented") }
func garrison_spot_check(garr int, target int) int                           { panic("!implemented") }
func garrison_summary(pl int)                                                { panic("!implemented") }
func gen_include_section()                                                   { panic("!implemented") }
func gen_include_sup(eat_pl int)                                             { panic("!implemented") }
func generate_one_treasure(who int)                                          { panic("!implemented") }
func generate_treasure(n int, i int)                                         { panic("!implemented") }
func get_rid_of_collapsed_mine(i int)                                        { panic("!implemented") }
func glob_init()                                                             { panic("!implemented") }
func gm_report(gm_player int)                                                { panic("!implemented") }
func gm_show_all_skills(skill_player int, boolean int)                       { panic("!implemented") }
func has_city(where int) int                                                 { panic("!implemented") }
func has_holy_plant(who int) int                                             { panic("!implemented") }
func has_holy_symbol(who int) int                                            { panic("!implemented") }
func here_pos(winner int) int                                                { panic("!implemented") }
func hinder_med_chance(who int) int                                          { panic("!implemented") }
func hinder_med_omen(target int, who int)                                    { panic("!implemented") }
func how_many(who1 int, who2 int, item int, qty int, have_left int) int      { panic("!implemented") }
func i_generic_harvest(c *command, t *harvest) int                           { panic("!implemented") }
func i_generic_make(c *command, t *make_) int                                { panic("!implemented") }
func i_petty_thief(c *command) int                                           { panic("!implemented") }
func i_repair(c *command) int                                                { panic("!implemented") }
func i_sail(c *command) int                                                  { panic("!implemented") }
func immediate_commands()                                                    { panic("!implemented") }
func init_collect_list()                                                     { panic("!implemented") }
func init_load_sup(who int)                                                  { panic("!implemented") }
func init_locs_touched()                                                     { panic("!implemented") }
func init_lower()                                                            { panic("!implemented") }
func init_ocean_chars()                                                      { panic("!implemented") }
func init_random()                                                           { panic("!implemented") }
func init_savage_attacks()                                                   { panic("!implemented") }
func init_spaces()                                                           { panic("!implemented") }
func init_weather_views()                                                    { panic("!implemented") }
func random() int                                                            { panic("!implemented") }
func int_comp(q1, q2 interface{}) int                                        { panic("!implemented") }
func inv_item_comp(q1, q2 *item_ent) int                                     { panic("!implemented") }
func is_defend(who int, target int) int                                      { panic("!implemented") }
func is_hostile(ruler int, who int) int                                      { panic("!implemented") }
func is_port_city_(where int) int                                            { panic("!implemented") }
func item_fly_cap(item int) int                                              { panic("!implemented") }
func item_land_cap(item int) int                                             { panic("!implemented") }
func item_ride_cap(item int) int                                             { panic("!implemented") }
func item_weight(item int) int                                               { panic("!implemented") }
func join_guild(who int, school int) int                                     { panic("!implemented") }
func list_order_templates()                                                  { panic("!implemented") }
func list_partial_skills(who int, who2 int, s string)                        { panic("!implemented") }
func list_pending_trades(who int, num int)                                   { panic("!implemented") }
func list_skill_sup(who int, e *skill_ent, prefix string)                    { panic("!implemented") }
func load_cmap_players() int                                                 { panic("!implemented") }
func load_db()                                                               { panic("!implemented") }
func loc_contains_hidden(target int) int                                     { panic("!implemented") }
func loc_depth(i int) int                                                    { panic("!implemented") }
func loc_target(target int) int                                              { panic("!implemented") }
func location_production()                                                   { panic("!implemented") }
func lock_tag()                                                              { panic("!implemented") }
func log_output(code int, format string, args ...interface{})                { panic("!implemented") }
func los_province_distance(who int, i int) int                               { panic("!implemented") }
func mail_reports()                                                          { panic("!implemented") }
func make_tower_guild(new_ int, skill int)                                   { panic("!implemented") }
func market_here(where int) int                                              { panic("!implemented") }
func market_report(who int, where int)                                       { panic("!implemented") }
func match_all_trades()                                                      { panic("!implemented") }
func mine_depth(where int) int                                               { panic("!implemented") }
func move_stack(who int, where int)                                          { panic("!implemented") }
func my_free(ptr interface{})                                                { panic("!implemented") }
func my_malloc(size int)                                                     { panic("!implemented") }
func my_realloc(ptr interface{}, size int)                                   { panic("!implemented") }
func nation_s(n int) string                                                  { panic("!implemented") }
func natural_weather()                                                       { panic("!implemented") }
func near_rocky_coast(loc int) int                                           { panic("!implemented") }
func nearby_grave(province int) int                                          { panic("!implemented") }
func new_player_top(mail_now int)                                            { panic("!implemented") }
func new_potion(who int) int                                                 { panic("!implemented") }
func new_suffuse_ring(where int) int                                         { panic("!implemented") }
func new_trade(who int, kind int, item int) *trade                           { panic("!implemented") }
func next_np_turn(pl int) int                                                { panic("!implemented") }
func notify_loc_shroud(where int)                                            { panic("!implemented") }
func np_req_s(skill int) string                                              { panic("!implemented") }
func npc_move(who int)                                                       { panic("!implemented") }
func nps_invested(who int) int                                               { panic("!implemented") }
func oly_parse(c *command, cmd string) int                                   { panic("!implemented") }
func oly_parse_cmd(_ *command, _ string) int                                 { panic("!implemented") }
func open_faery_hill(subloc int)                                             { panic("!implemented") }
func open_logfile()                                                          { panic("!implemented") }
func open_logfile_nondestruct()                                              { panic("!implemented") }
func open_times()                                                            { panic("!implemented") }
func opium_market_delta(where int)                                           { panic("!implemented") }
func orders_template(i int, n int)                                           { panic("!implemented") }
func out(who int, format string, args ...interface{})                        { panic("!implemented") }
func parse_arg(who int, s string) int                                        { panic("!implemented") }
func parse_wait_args(c *command) string                                      { panic("!implemented") }
func peaceful_enter(who int, where int, destination int) int                 { panic("!implemented") }
func pick_starting_city(nation int, n int) int                               { panic("!implemented") }
func ping_garrisons()                                                        { panic("!implemented") }
func place_here(where int, who int) int                                      { panic("!implemented") }
func player_banner()                                                         { panic("!implemented") }
func player_controls_loc(where int) int                                      { panic("!implemented") }
func player_ent_info()                                                       { panic("!implemented") }
func player_report()                                                         { panic("!implemented") }
func player_report_sup(who int)                                              { panic("!implemented") }
func post_month()                                                            { panic("!implemented") }
func post_production()                                                       { panic("!implemented") }
func pre_month()                                                             { panic("!implemented") }
func print_admit(pl int)                                                     { panic("!implemented") }
func print_att(who int, num int)                                             { panic("!implemented") }
func print_dot(ch int)                                                       { panic("!implemented") }
func print_hiring_status(pl int)                                             { panic("!implemented") }
func print_unformed(pl int)                                                  { panic("!implemented") }
func prisoner_movement_escape_check(who int)                                 { panic("!implemented") }
func process_orders()                                                        { panic("!implemented") }
func province_admin(i int) int                                               { panic("!implemented") }
func province_gate_here(where int) int                                       { panic("!implemented") }
func province_subloc(i int, subloc int) int                                  { panic("!implemented") }
func put_back_cookie(who int)                                                { panic("!implemented") }
func queue(who int, format string, args ...interface{})                      { panic("!implemented") }
func queue_npc_orders()                                                      { panic("!implemented") }
func random_beast(n int) int                                                 { panic("!implemented") }
func random_trade_good() int                                                 { panic("!implemented") }
func read_spool(mail_now int) int                                            { panic("!implemented") }
func regular_combat(_ int, _ int, _ int, _ int) int                          { panic("!implemented") }
func remove_comment(line string)                                             { panic("!implemented") }
func remove_ctrl_chars(line string)                                          { panic("!implemented") }
func rename_act_join_files()                                                 { panic("!implemented") }
func report_account_out(pl int, who int) int                                 { panic("!implemented") }
func reseed_monster_sublocs()                                                { panic("!implemented") }
func reset_cast_where(who int) int                                           { panic("!implemented") }
func restore_stack_actions(who int)                                          { panic("!implemented") }
func savage_hates(where int) int                                             { panic("!implemented") }
func save_box(fp *io.FILE, a int)                                            { panic("!implemented") }
func save_db()                                                               { panic("!implemented") }
func save_logdir()                                                           { panic("!implemented") }
func save_player_orders(pl int)                                              { panic("!implemented") }
func scan_char_item_lore()                                                   { panic("!implemented") }
func scan_char_skill_lore()                                                  { panic("!implemented") }
func seed_city(where int)                                                    { panic("!implemented") }
func seed_city_near_lists()                                                  { panic("!implemented") }
func seed_city_trade(where int)                                              { panic("!implemented") }
func seed_common_tradegoods()                                                { panic("!implemented") }
func seed_cookies()                                                          { panic("!implemented") }
func seed_initial_locations()                                                { panic("!implemented") }
func seed_monster_sublocs()                                                  { panic("!implemented") }
func seed_orcs()                                                             { panic("!implemented") }
func seed_population()                                                       { panic("!implemented") }
func seed_rare_tradegoods()                                                  { panic("!implemented") }
func seed_taxes()                                                            { panic("!implemented") }
func select_attacker(who int, target int) int                                { panic("!implemented") }
func send_rep(who int, turn int) int                                         { panic("!implemented") }
func set_att(who int, a int, attitude int)                                   { panic("!implemented") }
func set_bit(kr *sparse, i int)                                              { panic("!implemented") }
func set_html_pass(pl int)                                                   { panic("!implemented") }
func setup_html_all()                                                        { panic("!implemented") }
func setup_html_dir(pl int)                                                  { panic("!implemented") }
func ship_cap(ship int) int                                                  { panic("!implemented") }
func ship_storm_check(ship int)                                              { panic("!implemented") }
func ship_summary(pl int)                                                    { panic("!implemented") }
func show_carry_capacity(who int, who2 int)                                  { panic("!implemented") }
func show_item_skills(who int, who2 int)                                     { panic("!implemented") }
func show_item_where(who int, target int)                                    { panic("!implemented") }
func show_loc(who int, where int)                                            { panic("!implemented") }
func show_loc_posts(pl int, i int, boolean int)                              { panic("!implemented") }
func show_lore_sheets()                                                      { panic("!implemented") }
func show_unclaimed(i int, n int)                                            { panic("!implemented") }
func sort_for_output(l ilist)                                                { panic("!implemented") }
func sout(format string, args ...interface{}) string                         { panic("!implemented") }
func srandom(seed uint32)                                                    { panic("!implemented") }
func stage(s string)                                                         { panic("!implemented") }
func starting_noble_points(i int) int                                        { panic("!implemented") }
func storm_report(pl int)                                                    { panic("!implemented") }
func sub_item(where int, item int, something int) int                        { panic("!implemented") }
func summary_report()                                                        { panic("!implemented") }
func survive_fatal(pris int) int                                             { panic("!implemented") }
func tagout(who int, format string, args ...interface{})                     { panic("!implemented") }
func tags_off()                                                              { panic("!implemented") }
func tags_on()                                                               { panic("!implemented") }
func test_bit(kr sparse, i int) int                                          { panic("!implemented") }
func test_random() int                                                       { panic("!implemented") }
func text_list_free(l []string)                                              { panic("!implemented") }
func touch_loc_after_move(who int, where int)                                { panic("!implemented") }
func trade_suffuse_ring(where int)                                           { panic("!implemented") }
func update_city_trade(where int, action int, item int, qty int, base_price int, n int) {
	panic("!implemented")
}
func update_faery()                                              { panic("!implemented") }
func update_markets()                                            { panic("!implemented") }
func update_weather_view_locs(who int, where int)                { panic("!implemented") }
func v_add_forts(c *command) int                                 { panic("!implemented") }
func v_add_iron_shoring(c *command) int                          { panic("!implemented") }
func v_add_keels(c *command) int                                 { panic("!implemented") }
func v_add_ports(c *command) int                                 { panic("!implemented") }
func v_add_ram(c *command) int                                   { panic("!implemented") }
func v_add_sails(c *command) int                                 { panic("!implemented") }
func v_add_wooden_shoring(c *command) int                        { panic("!implemented") }
func v_adv_med(c *command) int                                   { panic("!implemented") }
func v_appear_common(c *command) int                             { panic("!implemented") }
func v_archery(c *command) int                                   { panic("!implemented") }
func v_arrange_mugging(c *command) int                           { panic("!implemented") }
func v_art_crown(c *command) int                                 { panic("!implemented") }
func v_art_destroy(c *command) int                               { panic("!implemented") }
func v_art_orb(c *command) int                                   { panic("!implemented") }
func v_art_teleport(c *command) int                              { panic("!implemented") }
func v_assassinate(c *command) int                               { panic("!implemented") }
func v_attack_tactics(c *command) int                            { panic("!implemented") }
func v_aura_blast(c *command) int                                { panic("!implemented") }
func v_aura_reflect(c *command) int                              { panic("!implemented") }
func v_banish_corpses(c *command) int                            { panic("!implemented") }
func v_banish_undead(c *command) int                             { panic("!implemented") }
func v_bar_loc(c *command) int                                   { panic("!implemented") }
func v_bind_storm(c *command) int                                { panic("!implemented") }
func v_bird_spy(c *command) int                                  { panic("!implemented") }
func v_bless_follower(c *command) int                            { panic("!implemented") }
func v_bless_fort(c *command) int                                { panic("!implemented") }
func v_breed(c *command) int                                     { panic("!implemented") }
func v_brew(c *command) int                                      { panic("!implemented") }
func v_bribe(c *command) int                                     { panic("!implemented") }
func v_build_wagons(c *command) int                              { panic("!implemented") }
func v_capture_beasts(c *command) int                            { panic("!implemented") }
func v_catch(c *command) int                                     { panic("!implemented") }
func v_cloak_creat(c *command) int                               { panic("!implemented") }
func v_cloak_reg(c *command) int                                 { panic("!implemented") }
func v_collect(c *command) int                                   { panic("!implemented") }
func v_conceal_arts(c *command) int                              { panic("!implemented") }
func v_conceal_location(c *command) int                          { panic("!implemented") }
func v_conceal_mine(c *command) int                              { panic("!implemented") }
func v_conceal_nation(c *command) int                            { panic("!implemented") }
func v_create_dirt_golem(c *command) int                         { panic("!implemented") }
func v_create_flesh_golem(c *command) int                        { panic("!implemented") }
func v_create_holy_symbol(c *command) int                        { panic("!implemented") }
func v_create_iron_golem(c *command) int                         { panic("!implemented") }
func v_create_mithril(c *command) int                            { panic("!implemented") }
func v_create_ninja(c *command) int                              { panic("!implemented") }
func v_curse_noncreat(c *command) int                            { panic("!implemented") }
func v_death_fog(c *command) int                                 { panic("!implemented") }
func v_decrease_demand(c *command) int                           { panic("!implemented") }
func v_decrease_supply(c *command) int                           { panic("!implemented") }
func v_dedicate_temple(c *command) int                           { panic("!implemented") }
func v_dedicate_tower(c *command) int                            { panic("!implemented") }
func v_deep_identify(c *command) int                             { panic("!implemented") }
func v_defense(c *command) int                                   { panic("!implemented") }
func v_defense_tactics(c *command) int                           { panic("!implemented") }
func v_destroy_art(c *command) int                               { panic("!implemented") }
func v_detect_abil(c *command) int                               { panic("!implemented") }
func v_detect_arts(c *command) int                               { panic("!implemented") }
func v_detect_beasts(c *command) int                             { panic("!implemented") }
func v_detect_gates(c *command) int                              { panic("!implemented") }
func v_detect_scry(c *command) int                               { panic("!implemented") }
func v_direct_storm(c *command) int                              { panic("!implemented") }
func v_dispel_abil(c *command) int                               { panic("!implemented") }
func v_dispel_region(c *command) int                             { panic("!implemented") }
func v_dissipate(c *command) int                                 { panic("!implemented") }
func v_draw_crowds(c *command) int                               { panic("!implemented") }
func v_eat_dead(c *command) int                                  { panic("!implemented") }
func v_edge_of_kireus(c *command) int                            { panic("!implemented") }
func v_enchant_guard(c *command) int                             { panic("!implemented") }
func v_fierce_wind(c *command) int                               { panic("!implemented") }
func v_fight_to_death(c *command) int                            { panic("!implemented") }
func v_find_food(c *command) int                                 { panic("!implemented") }
func v_find_forest_trail(c *command) int                         { panic("!implemented") }
func v_find_hidden_features(c *command) int                      { panic("!implemented") }
func v_find_mountain_trail(c *command) int                       { panic("!implemented") }
func v_find_rich(c *command) int                                 { panic("!implemented") }
func v_fish(c *command) int                                      { panic("!implemented") }
func v_forced_march(c *command) int                              { panic("!implemented") }
func v_forge_art_x(c *command) int                               { panic("!implemented") }
func v_forge_aura(c *command) int                                { panic("!implemented") }
func v_forge_palantir(c *command) int                            { panic("!implemented") }
func v_format(c *command) int                                    { panic("!implemented") }
func v_fortify_castle(c *command) int                            { panic("!implemented") }
func v_gather_holy_plant(c *command) int                         { panic("!implemented") }
func v_generic_trap(c *command) int                              { panic("!implemented") }
func v_grow_pop(c *command) int                                  { panic("!implemented") }
func v_heal(c *command) int                                      { panic("!implemented") }
func v_hide(c *command) int                                      { panic("!implemented") }
func v_hide_item(c *command) int                                 { panic("!implemented") }
func v_hide_money(c *command) int                                { panic("!implemented") }
func v_hinder_med(c *command) int                                { panic("!implemented") }
func v_hinder_med_b(c *command) int                              { panic("!implemented") }
func v_identify(c *command) int                                  { panic("!implemented") }
func v_implicit(c *command) int                                  { panic("!implemented") }
func v_improve_charisma(c *command) int                          { panic("!implemented") }
func v_improve_fort(c *command) int                              { panic("!implemented") }
func v_improve_logging(c *command) int                           { panic("!implemented") }
func v_improve_mining(c *command) int                            { panic("!implemented") }
func v_improve_opium(c *command) int                             { panic("!implemented") }
func v_improve_quarrying(c *command) int                         { panic("!implemented") }
func v_improve_smithing(c *command) int                          { panic("!implemented") }
func v_incite(c *command) int                                    { panic("!implemented") }
func v_increase_demand(c *command) int                           { panic("!implemented") }
func v_increase_supply(c *command) int                           { panic("!implemented") }
func v_jump_gate(c *command) int                                 { panic("!implemented") }
func v_keep_savage(c *command) int                               { panic("!implemented") }
func v_keep_undead(c *command) int                               { panic("!implemented") }
func v_last_rites(c *command) int                                { panic("!implemented") }
func v_lead_to_gold(c *command) int                              { panic("!implemented") }
func v_lightning(c *command) int                                 { panic("!implemented") }
func v_locate_char(c *command) int                               { panic("!implemented") }
func v_mage_menial(c *command) int                               { panic("!implemented") }
func v_meditate(c *command) int                                  { panic("!implemented") }
func v_mine_gate_crystal(c *command) int                         { panic("!implemented") }
func v_mine_gold(c *command) int                                 { panic("!implemented") }
func v_mine_iron(c *command) int                                 { panic("!implemented") }
func v_mine_mithril(c *command) int                              { panic("!implemented") }
func v_moat_castle(c *command) int                               { panic("!implemented") }
func v_move_attack(c *command) int                               { panic("!implemented") }
func v_mutate_art(c *command) int                                { panic("!implemented") }
func v_notab(c *command) int                                     { panic("!implemented") }
func v_notify_jump(c *command) int                               { panic("!implemented") }
func v_notify_unseal(c *command) int                             { panic("!implemented") }
func v_obscure_art(c *command) int                               { panic("!implemented") }
func v_obscure_forest_trail(c *command) int                      { panic("!implemented") }
func v_obscure_mountain_trail(c *command) int                    { panic("!implemented") }
func v_personal_fight_to_death(c *command) int                   { panic("!implemented") }
func v_persuade_oath(c *command) int                             { panic("!implemented") }
func v_petty_thief(c *command) int                               { panic("!implemented") }
func v_power_jewel(c *command) int                               { panic("!implemented") }
func v_prac_control(c *command) int                              { panic("!implemented") }
func v_prac_protect(c *command) int                              { panic("!implemented") }
func v_practice(c *command) int                                  { panic("!implemented") }
func v_prep_ritual(c *command) int                               { panic("!implemented") }
func v_proj_cast(c *command) int                                 { panic("!implemented") }
func v_proselytise(c *command) int                               { panic("!implemented") }
func v_protect_mine(c *command) int                              { panic("!implemented") }
func v_quarry(c *command) int                                    { panic("!implemented") }
func v_quick_cast(c *command) int                                { panic("!implemented") }
func v_raise(c *command) int                                     { panic("!implemented") }
func v_raise_corpses(c *command) int                             { panic("!implemented") }
func v_rally(c *command) int                                     { panic("!implemented") }
func v_raze(c *command) int                                      { panic("!implemented") }
func v_recruit_elves(c *command) int                             { panic("!implemented") }
func v_rem_art_cloak(c *command) int                             { panic("!implemented") }
func v_rem_seal(c *command) int                                  { panic("!implemented") }
func v_remove_bless(c *command) int                              { panic("!implemented") }
func v_remove_forts(c *command) int                              { panic("!implemented") }
func v_remove_keels(c *command) int                              { panic("!implemented") }
func v_remove_ports(c *command) int                              { panic("!implemented") }
func v_remove_ram(c *command) int                                { panic("!implemented") }
func v_remove_sails(c *command) int                              { panic("!implemented") }
func v_renew_storm(c *command) int                               { panic("!implemented") }
func v_resurrect(c *command) int                                 { panic("!implemented") }
func v_reveal_arts(c *command) int                               { panic("!implemented") }
func v_reveal_key(c *command) int                                { panic("!implemented") }
func v_reveal_mage(c *command) int                               { panic("!implemented") }
func v_reveal_vision(c *command) int                             { panic("!implemented") }
func v_reverse_jump(c *command) int                              { panic("!implemented") }
func v_sail(c *command) int                                      { panic("!implemented") }
func v_save_proj(c *command) int                                 { panic("!implemented") }
func v_save_quick(c *command) int                                { panic("!implemented") }
func v_scry_region(c *command) int                               { panic("!implemented") }
func v_seal_gate(c *command) int                                 { panic("!implemented") }
func v_seize_storm(c *command) int                               { panic("!implemented") }
func v_shipbuild(c *command) int                                 { panic("!implemented") }
func v_show_art_creat(c *command) int                            { panic("!implemented") }
func v_show_art_reg(c *command) int                              { panic("!implemented") }
func v_shroud_abil(c *command) int                               { panic("!implemented") }
func v_shroud_region(c *command) int                             { panic("!implemented") }
func v_smuggle_goods(c *command) int                             { panic("!implemented") }
func v_smuggle_men(c *command) int                               { panic("!implemented") }
func v_sneak(c *command) int                                     { panic("!implemented") }
func v_spy_inv(c *command) int                                   { panic("!implemented") }
func v_spy_lord(c *command) int                                  { panic("!implemented") }
func v_spy_skills(c *command) int                                { panic("!implemented") }
func v_strengthen_castle(c *command) int                         { panic("!implemented") }
func v_study(c *command) int                                     { panic("!implemented") }
func v_summon_aid(c *command) int                                { panic("!implemented") }
func v_summon_fog(c *command) int                                { panic("!implemented") }
func v_summon_rain(c *command) int                               { panic("!implemented") }
func v_summon_savage(c *command) int                             { panic("!implemented") }
func v_summon_wind(c *command) int                               { panic("!implemented") }
func v_swordplay(c *command) int                                 { panic("!implemented") }
func v_tap_health(c *command) int                                { panic("!implemented") }
func v_teach(c *command) int                                     { panic("!implemented") }
func v_teleport(c *command) int                                  { panic("!implemented") }
func v_teleport_item(c *command) int                             { panic("!implemented") }
func v_torture(c *command) int                                   { panic("!implemented") }
func v_trance(c *command) int                                    { panic("!implemented") }
func v_unbar_loc(c *command) int                                 { panic("!implemented") }
func v_undead_lord(c *command) int                               { panic("!implemented") }
func v_unobscure_art(c *command) int                             { panic("!implemented") }
func v_unseal_gate(c *command) int                               { panic("!implemented") }
func v_urchin_spy(c *command) int                                { panic("!implemented") }
func v_use_beasts(c *command) int                                { panic("!implemented") }
func v_use_cs(c *command) int                                    { panic("!implemented") }
func v_use_death(c *command) int                                 { panic("!implemented") }
func v_use_drum(c *command) int                                  { panic("!implemented") }
func v_use_faery_artifact(c *command) int                        { panic("!implemented") }
func v_use_fiery(c *command) int                                 { panic("!implemented") }
func v_use_heal(c *command) int                                  { panic("!implemented") }
func v_use_item(c *command) int                                  { panic("!implemented") }
func v_use_orb(c *command) int                                   { panic("!implemented") }
func v_use_proj_cast(c *command) int                             { panic("!implemented") }
func v_use_quick_cast(c *command) int                            { panic("!implemented") }
func v_use_slave(c *command) int                                 { panic("!implemented") }
func v_use_train_riding(c *command) int                          { panic("!implemented") }
func v_use_train_war(c *command) int                             { panic("!implemented") }
func v_use_weightlessness(c *command) int                        { panic("!implemented") }
func v_view_aura(c *command) int                                 { panic("!implemented") }
func v_vision_reg(c *command) int                                { panic("!implemented") }
func v_weaken_fort(c *command) int                               { panic("!implemented") }
func v_widen_entrance(c *command) int                            { panic("!implemented") }
func v_wood(c *command) int                                      { panic("!implemented") }
func v_write_spell(c *command) int                               { panic("!implemented") }
func v_yew(c *command) int                                       { panic("!implemented") }
func vector_players()                                            { panic("!implemented") }
func weekly_prisoner_escape_check()                              { panic("!implemented") }
func will_accept(target int, item int, who int, qty int) int     { panic("!implemented") }
func will_admit(pl int, who int, ruler int) int                  { panic("!implemented") }
func wiout(who int, ind int, format string, args ...interface{}) { panic("!implemented") }
func wout(who int, format string, args ...interface{})           { panic("!implemented") }
func write_email()                                               { panic("!implemented") }
func write_factions()                                            { panic("!implemented") }
func write_forwards()                                            { panic("!implemented") }
func write_nations_lists()                                       { panic("!implemented") }
func write_player(pl int)                                        { panic("!implemented") }
func write_player_list()                                         { panic("!implemented") }
func write_totimes()                                             { panic("!implemented") }
