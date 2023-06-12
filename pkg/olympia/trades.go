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

type TradeList []*Trade
type Trade struct {
	Id        int  `json:"id,omitempty"`
	Kind      int  `json:"kind,omitempty"` // BUY or SELL
	Item      int  `json:"item,omitempty"`
	Qty       int  `json:"qty,omitempty"`
	Cloak     bool `json:"cloak,omitempty"` // don't reveal identity of trader
	Cost      int  `json:"cost,omitempty"`
	Counter   int  `json:"counter,omitempty"`    // Counter to age and lose untraded goods
	HaveLeft  int  `json:"have-left,omitempty"`  // amount remaining
	MonthProd int  `json:"month-prod,omitempty"` // month city produces item
	OldQty    int  `json:"old-qty,omitempty"`    // qty at beginning of month, for trade goods
	sort      int  // temp key for sorting -- not saved
	who       int  // redundant -- not saved
}

type trade_l []*trade

func (l trade_l) rem_value(e *trade) trade_l {
	var cp trade_l
	for i := len(l) - 1; i >= 0; i-- {
		if l[i] != e {
			cp = append(cp, l[i])
		}
	}
	return cp
}

func (l trade_l) ToTradeList() (tl TradeList) {
	for _, e := range l {
		if !valid_box(e.item) {
			continue
		}
		// Weed out completed or cleared BUY and SELL trades, but don't touch PRODUCE or CONSUME zero-qty trades.
		// Scott Turner:
		//   Why not?  This causes problems because loop_trade ignores zero qty trades (as it probably should).
		//             (e.kind == BUY || e.kind == SELL) &&
		if e.qty <= 0 {
			continue
		}
		tl = append(tl, &Trade{
			Kind:      e.kind,
			Item:      e.item,
			Qty:       e.qty,
			Cost:      e.cost,
			Cloak:     e.cloak != FALSE,
			HaveLeft:  e.have_left,
			MonthProd: e.month_prod,
			OldQty:    e.old_qty,
			Counter:   e.counter,
		})
	}
	return tl
}
