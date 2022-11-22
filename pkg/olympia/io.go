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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	ext_boxnum          = 0
	line                []byte
	monster_subloc_init = false
	population_init     = false
	seed                [3]int
)

/*
 *  advance lets us do a one-line lookahead in the scanning routines
 */
func advance() {
	line = readlin()
}

/*
 *  linehash macro
 *
 *    data is stored in the form:
 *        xy data
 *    were `xy` is a key followed by one character of whitespace
 *
 *    linehash macro returns `xy` crunched into an int
 *
 *  Example:
 *
 *    s = "na Name field";
 *    c = linehash(s);
 *    assert(c == `na`);
 */
func linehash(t []byte) string {
	// (strlen(t) < 2 ? 0 : (((t)[0])<<8 | ((t)[1])))
	if len(t) < 2 {
		return ""
	}
	return string(t[:2])
}

func t_string(t []byte) []byte {
	//(strlen(t) >= 4 ? &(t)[3] : "")
	if len(t) < 4 {
		return nil
	}
	return t[3:]
}

/*
 *  io.c -- load and save the entity database
 */

/*
 * Mon May 10 07:19:21 1999 -- Scott Turner
 * Go through an ilist and substitute the new skill numbers.
 *
 */
func convert_skill(skill int) int {
	if skill == 120 {
		return 1000
	}
	if skill == 9502 {
		return 1001
	}
	if skill == 9503 {
		return 1002
	}
	if skill == 9554 {
		return 1095
	}
	if skill == 9582 {
		return 1004
	}
	if skill == 9614 {
		return 1005
	}
	if skill == 9616 {
		return 1006
	}
	if skill == 9618 {
		return 1096
	}
	if skill == 9620 {
		return 1097
	}
	if skill == 9622 {
		return 1098
	}
	if skill == 9617 {
		return 1090
	}
	if skill == 9615 {
		return 1091
	}
	if skill == 9621 {
		return 1092
	}
	if skill == 9619 {
		return 1093
	}
	if skill == 9510 {
		return 1094
	}

	if skill == 121 {
		return 1100
	}
	if skill == 9501 {
		return 1193
	}
	if skill == 9505 {
		return 1102
	}
	if skill == 9541 {
		return 1194
	}
	if skill == 9570 {
		return 1104
	}
	if skill == 9580 {
		return 1105
	}
	if skill == 9581 {
		return 1106
	}
	if skill == 9594 {
		return 1107
	}
	if skill == 9598 {
		return 1108
	}
	if skill == 9507 {
		return 1195
	}
	if skill == 9595 {
		return 1131
	}
	if skill == 9596 {
		return 1132
	}
	if skill == 9612 {
		return 1190
	}
	if skill == 9613 {
		return 1191
	}
	if skill == 9599 {
		return 1192
	}

	if skill == 122 {
		return 1200
	}
	if skill == 9509 {
		return 1201
	}
	if skill == 9519 {
		return 1202
	}
	if skill == 9520 {
		return 1203
	}
	if skill == 9521 {
		return 1204
	}
	if skill == 9591 {
		return 1205
	}
	if skill == 9522 {
		return 1230
	}
	if skill == 9523 {
		return 1231
	}
	if skill == 9562 {
		return 1232
	}
	if skill == 9574 {
		return 1233
	}
	if skill == 9593 {
		return 1290
	}
	if skill == 9140 {
		return 1291
	}
	if skill == 9149 {
		return 1292
	}

	if skill == 124 {
		return 1300
	}
	if skill == 9515 {
		return 1301
	}
	if skill == 9537 {
		return 1302
	}
	if skill == 9538 {
		return 1303
	}
	if skill == 9592 {
		return 1304
	}
	if skill == 9536 {
		return 1330
	}
	if skill == 9539 {
		return 1331
	}
	if skill == 9585 {
		return 1332
	}
	if skill == 9586 {
		return 1333
	}
	if skill == 9610 {
		return 1390
	}

	if skill == 125 {
		return 1400
	}
	if skill == 9542 {
		return 1401
	}
	if skill == 9566 {
		return 1402
	}
	if skill == 9611 {
		return 1490
	}
	if skill == 9107 {
		return 1491
	}
	if skill == 9114 {
		return 1492
	}
	if skill == 9127 {
		return 1493
	}

	if skill == 126 {
		return 1500
	}
	if skill == 9551 {
		return 1501
	}
	if skill == 9552 {
		return 1502
	}
	if skill == 9584 {
		return 1503
	}
	if skill == 9549 {
		return 1530
	}
	if skill == 9550 {
		return 1531
	}
	if skill == 9573 {
		return 1532
	}
	if skill == 9590 {
		return 1533
	}
	if skill == 9553 {
		return 1590
	}
	if skill == 9623 {
		return 1591
	}

	if skill == 128 {
		return 1600
	}
	if skill == 9540 {
		return 1601
	}
	if skill == 9568 {
		return 1602
	}
	if skill == 9569 {
		return 1603
	}
	if skill == 9583 {
		return 1604
	}
	if skill == 9588 {
		return 1605
	}
	if skill == 9587 {
		return 1630
	}
	if skill == 9589 {
		return 1631
	}

	if skill == 129 {
		return 1700
	}
	if skill == 9563 {
		return 1701
	}
	if skill == 9564 {
		return 1702
	}
	if skill == 9128 {
		return 1703
	}
	if skill == 9565 {
		return 1730
	}
	if skill == 9567 {
		return 1731
	}
	if skill == 9129 {
		return 1790
	}
	if skill == 9130 {
		return 1791
	}

	if skill == 130 {
		return 1800
	}
	if skill == 9603 {
		return 1801
	}
	if skill == 9604 {
		return 1802
	}
	if skill == 9605 {
		return 1803
	}
	if skill == 9606 {
		return 1804
	}
	if skill == 9607 {
		return 1805
	}
	if skill == 9600 {
		return 1830
	}
	if skill == 9601 {
		return 1831
	}
	if skill == 9602 {
		return 1832
	}
	if skill == 9608 {
		return 1890
	}
	if skill == 9609 {
		return 1891
	}

	if skill == 131 {
		return 1900
	}
	if skill == 9145 {
		return 1901
	}
	if skill == 9579 {
		return 1902
	}
	if skill == 9150 {
		return 1903
	}
	if skill == 9517 {
		return 1930
	}
	if skill == 9529 {
		return 1931
	}
	if skill == 9530 {
		return 1932
	}
	if skill == 9146 {
		return 1990
	}

	if skill == 151 {
		return 2000
	}
	if skill == 9302 {
		return 2001
	}
	if skill == 9303 {
		return 2002
	}
	if skill == 9304 {
		return 2003
	}
	if skill == 9305 {
		return 2004
	}
	if skill == 9306 {
		return 2005
	}
	if skill == 9307 {
		return 2006
	}
	if skill == 9308 {
		return 2007
	}
	if skill == 9309 {
		return 2008
	}
	if skill == 9447 {
		return 2009
	}
	if skill == 9124 {
		return 2030
	}
	if skill == 9191 {
		return 2031
	}
	if skill == 9193 {
		return 2032
	}
	if skill == 9194 {
		return 2033
	}
	if skill == 9195 {
		return 2034
	}
	if skill == 9196 {
		return 2035
	}
	if skill == 9187 {
		return 2036
	}
	if skill == 9188 {
		return 2037
	}
	if skill == 9189 {
		return 2038
	}
	if skill == 9190 {
		return 2039
	}
	if skill == 9310 {
		return 2040
	}
	if skill == 9148 {
		return 2041
	}
	if skill == 9155 {
		return 2042
	}

	if skill == 152 {
		return 2100
	}
	if skill == 9312 {
		return 2101
	}
	if skill == 9313 {
		return 2102
	}
	if skill == 9314 {
		return 2103
	}
	if skill == 9315 {
		return 2104
	}
	if skill == 9316 {
		return 2105
	}
	if skill == 9317 {
		return 2106
	}
	if skill == 9318 {
		return 2107
	}
	if skill == 9319 {
		return 2108
	}
	if skill == 9440 {
		return 2109
	}
	if skill == 9311 {
		return 2130
	}
	if skill == 9400 {
		return 2131
	}
	if skill == 9401 {
		return 2132
	}
	if skill == 9402 {
		return 2133
	}
	if skill == 9403 {
		return 2134
	}
	if skill == 9404 {
		return 2135
	}
	if skill == 9405 {
		return 2136
	}
	if skill == 9406 {
		return 2137
	}
	if skill == 9433 {
		return 2138
	}
	if skill == 9156 {
		return 2139
	}

	if skill == 153 {
		return 2200
	}
	if skill == 9322 {
		return 2201
	}
	if skill == 9323 {
		return 2202
	}
	if skill == 9324 {
		return 2203
	}
	if skill == 9325 {
		return 2204
	}
	if skill == 9326 {
		return 2205
	}
	if skill == 9327 {
		return 2206
	}
	if skill == 9328 {
		return 2207
	}
	if skill == 9329 {
		return 2208
	}
	if skill == 9441 {
		return 2209
	}
	if skill == 9504 {
		return 2231
	}
	if skill == 9506 {
		return 2232
	}
	if skill == 9508 {
		return 2233
	}
	if skill == 9320 {
		return 2234
	}
	if skill == 9321 {
		return 2235
	}
	if skill == 9434 {
		return 2236
	}
	if skill == 9157 {
		return 2237
	}

	if skill == 154 {
		return 2300
	}
	if skill == 9332 {
		return 2301
	}
	if skill == 9333 {
		return 2302
	}
	if skill == 9334 {
		return 2303
	}
	if skill == 9335 {
		return 2304
	}
	if skill == 9336 {
		return 2305
	}
	if skill == 9337 {
		return 2306
	}
	if skill == 9338 {
		return 2307
	}
	if skill == 9339 {
		return 2308
	}
	if skill == 9442 {
		return 2309
	}
	if skill == 9407 {
		return 2330
	}
	if skill == 9408 {
		return 2331
	}
	if skill == 9409 {
		return 2332
	}
	if skill == 9410 {
		return 2333
	}
	if skill == 9411 {
		return 2334
	}
	if skill == 9412 {
		return 2335
	}
	if skill == 9413 {
		return 2336
	}
	if skill == 9435 {
		return 2337
	}
	if skill == 9158 {
		return 2338
	}

	if skill == 155 {
		return 2400
	}
	if skill == 9342 {
		return 2401
	}
	if skill == 9343 {
		return 2402
	}
	if skill == 9344 {
		return 2403
	}
	if skill == 9345 {
		return 2404
	}
	if skill == 9346 {
		return 2405
	}
	if skill == 9347 {
		return 2406
	}
	if skill == 9348 {
		return 2407
	}
	if skill == 9349 {
		return 2408
	}
	if skill == 9443 {
		return 2409
	}
	if skill == 9419 {
		return 2430
	}
	if skill == 9420 {
		return 2431
	}
	if skill == 9421 {
		return 2432
	}
	if skill == 9422 {
		return 2433
	}
	if skill == 9341 {
		return 2434
	}
	if skill == 9436 {
		return 2435
	}
	if skill == 9448 {
		return 2436
	}
	if skill == 9159 {
		return 2437
	}

	if skill == 156 {
		return 2500
	}
	if skill == 9352 {
		return 2501
	}
	if skill == 9353 {
		return 2502
	}
	if skill == 9354 {
		return 2503
	}
	if skill == 9355 {
		return 2504
	}
	if skill == 9356 {
		return 2505
	}
	if skill == 9357 {
		return 2506
	}
	if skill == 9358 {
		return 2507
	}
	if skill == 9359 {
		return 2508
	}
	if skill == 9444 {
		return 2509
	}
	if skill == 9414 {
		return 2530
	}
	if skill == 9415 {
		return 2531
	}
	if skill == 9416 {
		return 2532
	}
	if skill == 9417 {
		return 2533
	}
	if skill == 9418 {
		return 2534
	}
	if skill == 9437 {
		return 2535
	}
	if skill == 9162 {
		return 2536
	}

	if skill == 157 {
		return 2600
	}
	if skill == 9362 {
		return 2601
	}
	if skill == 9363 {
		return 2602
	}
	if skill == 9364 {
		return 2603
	}
	if skill == 9365 {
		return 2604
	}
	if skill == 9366 {
		return 2605
	}
	if skill == 9367 {
		return 2606
	}
	if skill == 9368 {
		return 2607
	}
	if skill == 9369 {
		return 2608
	}
	if skill == 9445 {
		return 2609
	}
	if skill == 9423 {
		return 2630
	}
	if skill == 9424 {
		return 2631
	}
	if skill == 9425 {
		return 2632
	}
	if skill == 9426 {
		return 2633
	}
	if skill == 9427 {
		return 2634
	}
	if skill == 9428 {
		return 2635
	}
	if skill == 9438 {
		return 2636
	}
	if skill == 9163 {
		return 2637
	}

	if skill == 158 {
		return 2700
	}
	if skill == 9372 {
		return 2701
	}
	if skill == 9373 {
		return 2702
	}
	if skill == 9374 {
		return 2703
	}
	if skill == 9375 {
		return 2704
	}
	if skill == 9376 {
		return 2705
	}
	if skill == 9377 {
		return 2706
	}
	if skill == 9378 {
		return 2707
	}
	if skill == 9379 {
		return 2708
	}
	if skill == 9446 {
		return 2709
	}
	if skill == 9429 {
		return 2730
	}
	if skill == 9430 {
		return 2731
	}
	if skill == 9431 {
		return 2732
	}
	if skill == 9432 {
		return 2733
	}
	if skill == 9439 {
		return 2734
	}
	if skill == 9166 {
		return 2735
	}

	if skill == 160 {
		return 2800
	}
	if skill == 9101 {
		return 2801
	}
	if skill == 9103 {
		return 2802
	}
	if skill == 9104 {
		return 2803
	}
	if skill == 9105 {
		return 2804
	}
	if skill == 9106 {
		return 2830
	}
	if skill == 9123 {
		return 2831
	}
	if skill == 9126 {
		return 2832
	}
	if skill == 9147 {
		return 2833
	}
	if skill == 9169 {
		return 2834
	}
	if skill == 9170 {
		return 2835
	}
	if skill == 9173 {
		return 2836
	}
	if skill == 9174 {
		return 2837
	}
	if skill == 9175 {
		return 2838
	}
	if skill == 9178 {
		return 2839
	}
	if skill == 9135 {
		return 2840
	}
	if skill == 9624 {
		return 2841
	}

	if skill == 162 {
		return 2900
	}
	if skill == 9141 {
		return 2901
	}
	if skill == 9142 {
		return 2930
	}
	if skill == 9143 {
		return 2931
	}
	if skill == 9151 {
		return 2932
	}
	if skill == 9160 {
		return 2933
	}
	if skill == 9161 {
		return 2934
	}
	if skill == 9167 {
		return 2935
	}
	if skill == 9171 {
		return 2936
	}
	if skill == 9172 {
		return 2937
	}
	if skill == 9177 {
		return 2938
	}
	if skill == 9197 {
		return 2939
	}
	if skill == 9136 {
		return 2940
	}

	if skill == 163 {
		return 3000
	}
	if skill == 9115 {
		return 3001
	}
	if skill == 9116 {
		return 3002
	}
	if skill == 9112 {
		return 3030
	}
	if skill == 9117 {
		return 3031
	}
	if skill == 9118 {
		return 3032
	}
	if skill == 9119 {
		return 3033
	}
	if skill == 9120 {
		return 3034
	}
	if skill == 9121 {
		return 3035
	}
	if skill == 9122 {
		return 3036
	}
	if skill == 9152 {
		return 3037
	}
	if skill == 9165 {
		return 3038
	}
	if skill == 9137 {
		return 3039
	}

	if skill == 164 {
		return 3100
	}
	if skill == 9108 {
		return 3101
	}
	if skill == 9109 {
		return 3102
	}
	if skill == 9164 {
		return 3103
	}
	if skill == 9168 {
		return 3104
	}
	if skill == 9102 {
		return 3130
	}
	if skill == 9110 {
		return 3131
	}
	if skill == 9111 {
		return 3132
	}
	if skill == 9113 {
		return 3133
	}
	if skill == 9131 {
		return 3134
	}
	if skill == 9132 {
		return 3135
	}
	if skill == 9133 {
		return 3136
	}
	if skill == 9144 {
		return 3137
	}
	if skill == 9153 {
		return 3138
	}
	if skill == 9176 {
		return 3139
	}

	if skill == 165 {
		return 3200
	}
	if skill == 9179 {
		return 3201
	}
	if skill == 9180 {
		return 3202
	}
	if skill == 9185 {
		return 3203
	}
	if skill == 9154 {
		return 3230
	}
	if skill == 9181 {
		return 3231
	}
	if skill == 9182 {
		return 3232
	}
	if skill == 9183 {
		return 3233
	}
	if skill == 9184 {
		return 3234
	}
	if skill == 9186 {
		return 3235
	}
	if skill == 9138 {
		return 3236
	}
	if skill == 9625 {
		return 3237
	}
	if skill == 9125 {
		return 3238
	}

	if skill == 170 {
		return 3300
	}
	if skill == 9201 {
		return 3330
	}
	if skill == 9202 {
		return 3331
	}
	if skill == 9139 {
		return 3332
	}
	if skill == 9626 {
		return 3333
	}
	return skill
}

/*
 *  Returns the entity name in parenthesis, if available, to make the
 *  data files easier to read.
 */
func if_name(num int) string { /* to make the data files easier to read */
	if pretty_data_files == FALSE {
		return ""
	} else if !valid_box(num) {
		return ""
	} else if s := name(num); s != "" {
		return fmt.Sprintf(" (%s)", s)
	}
	return ""
}

func box_scan(t []byte) int {
	if n := atoi_b(t); valid_box(n) {
		return n
	}
	//#if 0
	//    /* temp fix */
	//    if (convert_skill(n) != n)
	//      return convert_skill(n);
	//#endif
	fprintf(os.Stderr, "box_scan(%d): bad reference: %s\n", ext_boxnum, line)
	return 0
}

func box_print(fp *os.File, header []byte, n int) {
	/* assert(!n || valid_box(n)); */
	if valid_box(n) {
		fprintf(fp, "%s%d%s\n", header, n, if_name(n))
	}
}

/*
 *  boxlist0_scan, boxlist0_print:
 *  same as boxlist_xxx, but allows zero
 */
func boxlist0_scan(s []byte, box_num int, l []int) []int {
	for len(s) != 0 {
		if iswhite(s[0]) {
			s = s[1:]
		} else if s[0] == '\\' { // continuation line follows
			s = readlin_ew()
		} else if isdigit(s[0]) {
			n := atoi_b(s)
			if n == 0 || valid_box(n) {
				l = append(l, n)
				//#if 0
				//          } else if (convert_skill(n) != n) { /* temp fix */
				//              l = append(l,  convert_skill(n));
				//#endif
			} else {
				fprintf(os.Stderr, "boxlist_scan(%d): bad box reference: %d\n", box_num, n)
			}
			for len(s) != 0 && isdigit(s[0]) {
				s = s[1:]
			}
		} else {
			break
		}
	}
	return l
}

func boxlist0_print(fp *os.File, header []byte, l []int) {
	count := 0
	for i := 0; i < len(l); i++ {
		if l[i] == 0 || valid_box(l[i]) {
			count++
			if count == 1 {
				fputb(header, fp)
			} else if count%11 == 10 { //continuation line
				fputs("\\\n\t", fp)
			}
			count++
			fprintf(fp, "%d ", l[i])
		}
	}

	if count != 0 {
		fprintf(fp, "\n")
	}
}

func boxlist_scan(s []byte, box_num int, l []int) []int {
	for len(s) != 0 {
		if iswhite(s[0]) {
			s = s[1:]
		} else if s[0] == '\\' { /* continuation line follows */
			s = readlin_ew()
		} else if isdigit(s[0]) || s[0] == '-' {
			n := atoi_b(s)
			if valid_box(n) || n == MONSTER_ATT {
				l = append(l, n)
				//#if 0
				//            } else if (convert_skill(n) != n) { /* temp fix */
				//              l = append(l,  convert_skill(n));
				//#endif
			} else {
				fprintf(os.Stderr, "boxlist_scan(%d): bad box reference: %d\n", box_num, n)
			}
			for len(s) != 0 && (isdigit(s[0]) || s[0] == '-') {
				s = s[1:]
			}
		} else {
			break
		}
	}
	return l
}

func boxlist_print(fp *os.File, header []byte, l []int) {
	if l == nil {
		return
	}

	count := 0
	for i := 0; i < len(l); i++ {
		if valid_box(l[i]) || l[i] == MONSTER_ATT {
			count++
			if count == 1 {
				fputb(header, fp)
			} else if count%11 == 10 { // commented out continuation line
				fputs("\\\n\t", fp)
			}
			fprintf(fp, "%d ", l[i])
		}
	}
	if count != 0 {
		fprintf(fp, "\n")
	}
}

func admit_print_sup(fp *os.File, p *admit) {
	if !valid_box(p.targ) {
		return
	} else if p.sense == 0 && len(p.l) == 0 {
		return
	}
	fprintf(fp, " am %d %d ", p.targ, p.sense)
	count := 2
	for i := 0; i < len(p.l); i++ {
		if valid_box(p.l[i]) {
			count++
			if count%11 == 10 { // continuation line
				fputs("\\\n\t", fp)
			}
			fprintf(fp, "%d ", p.l[i])
		}
	}
	if count != 0 {
		fprintf(fp, "\n")
	}
}

func admit_print(fp *os.File, p *entity_player) {
	for i := 0; i < len(p.admits); i++ {
		admit_print_sup(fp, p.admits[i])
	}
}

func admit_scan(s []byte, box_num int, pp *entity_player) {
	p := &admit{}

	count := 0
	for len(s) != 0 {
		if iswhite(s[0]) {
			s = s[1:]
		} else if s[0] == '\\' { /* continuation line follows */
			s = readlin_ew()
		} else if isdigit(s[0]) {
			count++
			n := atoi_b(s)
			switch count {
			case 1:
				p.targ = n
				break
			case 2:
				p.sense = n
				break
			default:
				/* Temp fix for nations */
				if n <= 1002 && n >= 1000 {
					n -= 3
				}
				if valid_box(n) {
					p.l = append(p.l, n)
				} else {
					fprintf(os.Stderr, "admit_scan(%d): bad box reference: %d\n", box_num, n)
				}
			}
			for len(s) != 0 && isdigit(s[0]) {
				s = s[1:]
			}
		} else {
			break
		}
	}

	if !valid_box(p.targ) {
		fprintf(os.Stderr, "admit_scan(%d): bad targ %d\n", box_num, p.targ)
		return
	}

	pp.admits = append(pp.admits, p)
}

func ilist_print(fp *os.File, header []byte, l []int) {
	if len(l) > 0 {
		fputb(header, fp)
		for i := 0; i < len(l); i++ {
			if i%11 == 10 { // continuation line
				fprintf(fp, "\\\n\t")
			}
			fprintf(fp, "%d ", l[i])
		}
		fprintf(fp, "\n")
	}
}

func ilist_scan(s []byte, l []int) {
	for len(s) != 0 {
		if iswhite(s[0]) {
			s = s[1:]
		} else if s[0] == '\\' { /* continuation line follows */
			s = readlin_ew()
		} else if isdigit(s[0]) {
			l = append(l, atoi_b(s))
			for len(s) != 0 && isdigit(s[0]) {
				s = s[1:]
			}
		} else {
			break
		}
	}
}

func known_print(fp *os.File, header []byte, kn sparse) {
	count, first := 0, true
	for _, i := range known_sparse_loop(kn) {
		if !valid_box(i) {
			continue
		}
		if first {
			fputb(header, fp)
			first = false
		} else if count%11 == 10 { // mdhender: commented out continuation line
			fprintf(fp, "\\\n\t")
		}
		count++
		fprintf(fp, "%d ", i)
	}
	if !first {
		fprintf(fp, "\n")
	}
}

func known_scan(s []byte, kn sparse, box_num int) sparse {
	for len(s) != 0 {
		if iswhite(s[0]) {
			s = s[1:]
		} else if s[0] == '\\' { /* continuation line follows */
			s = readlin_ew()
		} else if isdigit(s[0]) {
			n := atoi_b(s)
			if valid_box(n) {
				kn = set_bit(kn, n)
				//#if 0
				//          } else if (convert_skill(n) != n) { /* temp fix */
				//                set_bit(kn, convert_skill(n));
				//#endif
			} else {
				fprintf(os.Stderr, "known_scan(%d): bad box reference: %d\n", box_num, n)
			}

			for len(s) != 0 && isdigit(s[0]) {
				s = s[1:]
			}
		} else {
			break
		}
	}
	return kn
}

func skill_list_print(fp *os.File, header []byte, l skill_ent_l) {
	count := 0
	for i := 0; i < len(l); i++ {
		if valid_box(l[i].skill) {
			count++
			if count == 1 {
				fputb(header, fp)
			} else if count > 1 {
				fputs(" \\\n\t", fp)
			}
			fprintf(fp, "%d %d %d %d 0", l[i].skill, l[i].know, l[i].days_studied, l[i].experience)
		}
	}
	if count != 0 {
		fputs("\n", fp)
	}
}

// todo: ignores errors with line (too many values, too few, invalid, etc)
func skill_list_scan(s []byte, l []*skill_ent, box_num int) []*skill_ent {
	s = bytes.TrimSpace(s)
	for len(s) != 0 {
		foundContinuation := false
		newt := &skill_ent{}
		// sscanf(s, "%d %d %d %d %d", &newt.skill, &newt.know, &newt.days_studied, &newt.experience, &dummy);
		for i, f := range bytes.Fields(s) {
			switch i {
			case 0:
				newt.skill = convert_skill(atoi_b(f))
			case 1:
				newt.know = atoi_b(f)
			case 2:
				newt.days_studied = atoi_b(f)
			case 3:
				newt.experience = atoi_b(f)
			case 4: // ignore value
			case 5:
				if bytes.Equal(f, []byte{'\\'}) { // continuation line
					foundContinuation = true
				}
			}
		}
		if valid_box(newt.skill) {
			l = append(l, newt)
		} else {
			fprintf(os.Stderr, "skill_list_scan(%d): bad skill %d\n", box_num, newt.skill)
		}
		if foundContinuation { /* another entry follows */
			s = bytes.TrimSpace(readlin_ew())
		} else {
			break
		}
	}
	return l
}

/*
 *  Effect Functions
 *  Mon Aug  5 12:58:00 1996 -- Scott Turner
 *
 */
func effect_list_print(fp *os.File, header []byte, l []*effect) {
	count := 0
	for i := 0; i < len(l); i++ {
		count++
		if count == 1 {
			fputb(header, fp)
		} else if count > 1 { // continuation line
			fputs(" \\\n\t", fp)
		}
		fprintf(fp, "%d %d %d %d", l[i].type_, l[i].subtype, l[i].days, l[i].data)
	}
	if count != 0 {
		fputs("\n", fp)
	}
}

func effect_list_scan(s []byte, l []*effect) []*effect {
	s = bytes.TrimSpace(s)
	for len(s) != 0 {
		foundContinuation := false
		newt := &effect{}
		//sscanf(s, "%d %d %d %d", &newt.type_, &newt.subtype, &newt.days, &newt.data);
		for i, f := range bytes.Fields(s) {
			switch i {
			case 0:
				newt.type_ = atoi_b(f)
			case 1:
				newt.subtype = atoi_b(f)
			case 2:
				newt.days = atoi_b(f)
			case 3:
				newt.data = atoi_b(f)
			case 4:
				if bytes.Equal(f, []byte{'\\'}) { // continuation line
					foundContinuation = true
				}
			}
		}
		l = append(l, newt)
		if foundContinuation { /* another entry follows */
			s = bytes.TrimSpace(readlin_ew())
		} else {
			break
		}
	}
	return l
}

/*
 *  Build Functions
 *  Mon Aug  5 12:58:00 1996 -- Scott Turner
 *
 */
func build_list_print(fp *os.File, header []byte, l []*entity_build) {
	count := 0
	for i := 0; i < len(l); i++ {
		count++
		if count == 1 {
			fputb(header, fp)
		} else if count > 1 { // continuation line
			fputs(" \\\n\t", fp)
		}

		fprintf(fp, "%d %d %d %d", l[i].type_, l[i].build_materials, l[i].effort_required, l[i].effort_given)
	}
	if count != 0 {
		fputs("\n", fp)
	}
}

func build_list_scan(s []byte, l []*entity_build) []*entity_build {
	s = bytes.TrimSpace(s)
	for len(s) != 0 {
		foundContinuation := false
		newt := &entity_build{}
		//sscanf(s, "%d %d %d %d", &newt.type_, &newt.build_materials, &newt.effort_required, &newt.effort_given);
		for i, f := range bytes.Fields(s) {
			switch i {
			case 0:
				newt.type_ = atoi_b(f)
			case 1:
				newt.build_materials = atoi_b(f)
			case 2:
				newt.effort_required = atoi_b(f)
			case 3:
				newt.effort_given = atoi_b(f)
			case 4:
				if bytes.Equal(f, []byte{'\\'}) { // continuation line
					foundContinuation = true
				}
			}
		}
		l = append(l, newt)
		if foundContinuation { /* another entry follows */
			s = bytes.TrimSpace(readlin_ew())
		} else {
			break
		}
	}
	return l
}

func item_list_print(fp *os.File, header []byte, l []*item_ent) int {
	count := 0
	for i := 0; i < len(l); i++ {
		if valid_box(l[i].item) && l[i].qty > 0 {
			count++
			if count == 1 {
				fputb(header, fp)
			} else if count > 1 { // continuation line
				fputs(" \\\n\t", fp)
			}
			fprintf(fp, "%d %d", l[i].item, l[i].qty)
		}
	}
	if count != 0 {
		fputs("\n", fp)
	}

	return count
}

func item_list_scan(s []byte, l []*item_ent, box_num int) []*item_ent {
	s = bytes.TrimSpace(s)
	for len(s) != 0 {
		foundContinuation := false
		newt := &item_ent{}
		//sscanf(s, "%d %d", &newt.item, &newt.qty);
		for i, f := range bytes.Fields(s) {
			switch i {
			case 0:
				newt.item = atoi_b(f)
			case 1:
				newt.qty = atoi_b(f)
			case 2:
				if bytes.Equal(f, []byte{'\\'}) { // continuation line
					foundContinuation = true
				}
			}
		}
		if valid_box(newt.item) {
			l = append(l, newt)
		} else {
			fprintf(os.Stderr, "item_list_scan(%d): bad item %d\n", box_num, newt.item)
			fprintf(os.Stderr, "  bad item deleted from list.\n")
		}
		if foundContinuation { /* another entry follows */
			s = bytes.TrimSpace(readlin_ew())
		} else {
			break
		}
	}
	return l
}

func trade_list_print(fp *os.File, header []byte, l []*trade) {
	count := 0
	for i := 0; i < len(l); i++ {
		if valid_box(l[i].item) {
			/*
			 *  Weed out completed or cleared BUY and SELL trades, but don't
			 *  touch PRODUCE or CONSUME zero-qty trades.
			 *
			 *  Tue Jun  6 13:29:32 2000 -- Scott Turner
			 *
			 *  Why not?  This causes problems because loop_trade ignores zero
			 *  qty trades (as it probably should).
			 *
			 *      (l[i].kind == BUY || l[i].kind == SELL) &&
			 */
			if l[i].qty <= 0 {
				continue
			}
			count++
			if count == 1 {
				fputb(header, fp)
			} else if count > 1 {
				fputs(" \\\n\t", fp)
			}
			fprintf(fp, "%d %d %d %d %d %d %d %d %d", l[i].kind, l[i].item, l[i].qty, l[i].cost, l[i].cloak, l[i].have_left, l[i].month_prod, l[i].old_qty, l[i].counter)
		}
	}
	if count != 0 {
		fputs("\n", fp)
	}
}

func trade_list_scan(s []byte, l []*trade, box_num int) []*trade {
	s = bytes.TrimSpace(s)
	for len(s) != 0 {
		foundContinuation := false
		newt := &trade{who: box_num}
		//sscanf(s, "%d %d %d %d %d %d %d %d %d", &newt.kind, &newt.item, &newt.qty, &newt.cost, &newt.cloak, &newt.have_left, &newt.month_prod, &newt.old_qty, &newt.counter);
		for i, f := range bytes.Fields(s) {
			switch i {
			case 0:
				newt.kind = atoi_b(f)
			case 1:
				newt.item = atoi_b(f)
			case 2:
				newt.qty = atoi_b(f)
			case 3:
				newt.cost = atoi_b(f)
			case 4:
				newt.cloak = atoi_b(f)
			case 5:
				newt.have_left = atoi_b(f)
			case 6:
				newt.month_prod = atoi_b(f)
			case 7:
				newt.old_qty = atoi_b(f)
			case 8:
				newt.counter = atoi_b(f)
			case 9:
				if bytes.Equal(f, []byte{'\\'}) { // continuation line
					foundContinuation = true
				}
			}
		}
		if valid_box(newt.item) {
			l = append(l, newt)
		} else {
			fprintf(os.Stderr, "trade_list_scan(%d): bad item %d\n", box_num, newt.item)
		}

		if foundContinuation { /* another entry follows */
			s = bytes.TrimSpace(readlin_ew())
		} else {
			break
		}
	}
	return l
}

func req_list_print(fp *os.File, header []byte, l []*req_ent) {
	count := 0

	for i := 0; i < len(l); i++ {
		if valid_box(l[i].item) {
			count++
			if count == 1 {
				fputb(header, fp)
			}

			if count > 1 {
				fputs(" \\\n\t", fp)
			}

			fprintf(fp, "%d %d %d",
				l[i].item,
				l[i].qty,
				l[i].consume)
		}
	}

	if count != 0 {
		fputs("\n", fp)
	}
}

func eq_list_scan(s []byte, l []*req_ent, box_num int) []*req_ent {
	s = bytes.TrimSpace(s)
	for len(s) != 0 {
		foundContinuation := false
		newt := &req_ent{}
		//sscanf(s, "%d %d %d", &newt.item, &newt.qty, &newt.consume);
		for i, f := range bytes.Fields(s) {
			switch i {
			case 0:
				newt.item = atoi_b(f)
			case 1:
				newt.qty = atoi_b(f)
			case 2:
				newt.consume = atoi_b(f)
			case 3:
				if bytes.Equal(f, []byte{'\\'}) { // continuation line
					foundContinuation = true
				}
			}
		}
		if valid_box(newt.item) {
			l = append(l, newt)
		} else {
			fprintf(os.Stderr, "req_list_scan(%d): bad item %d\n", box_num, newt.item)
		}
		if foundContinuation { /* another entry follows */
			s = bytes.TrimSpace(readlin_ew())
		} else {
			break
		}
	}
	return l
}

func olytime_scan(s []byte, p *olytime) {
	//sscanf(s, "%hd %hd %d", &p.turn, &p.day, &p.days_since_epoch);
	for i, f := range bytes.Fields(s) {
		switch i {
		case 0:
			p.turn = atoi_b(f)
		case 1:
			p.day = atoi_b(f)
		case 2:
			p.days_since_epoch = atoi_b(f)
		}
	}
}

func olytime_print(fp *os.File, header []byte, p *olytime) {
	if p.turn != 0 || p.day != 0 || p.days_since_epoch != 0 {
		fprintf(fp, "%s%d %d %d\n", header, p.turn, p.day, p.days_since_epoch)
	}
}

/*
 *  Routine to check if a structure is completely empty.
 *
 *  Since structures may have elements which are not saved by save_db(),
 *  this routine may return true when, in fact, no data from the structure
 *  will be saved.  However, the next turn run will clear this up, since
 *  the re-loaded empty structure will now pass zero_check.
 *
 *  Using zero_check is more reliable than element testing, since we might
 *  forget to add one to the check list.  Also, our concern is long-term
 *  buildup of unused empty structure, so keeping one around for a turn
 *  is not a big deal.
 */
func zero_check(p interface{}) bool {
	// https://freshman.tech/snippets/go/check-empty-struct/
	return p == nil || reflect.ValueOf(p).IsZero()
}

func print_loc_info(fp *os.File, p *loc_info) {
	if zero_check(*p) {
		return
	}
	fprintf(fp, "LI\n")
	box_print(fp, []byte(" wh "), p.where)
	boxlist_print(fp, []byte(" hl "), p.here_list)
}

func scan_loc_info(p *loc_info, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)
		switch c {
		case `wh`:
			p.where = box_scan(t)
			break

		case `hl`:
			p.here_list = boxlist_scan(t, box_num, p.here_list)
			break

		default:
			fprintf(os.Stderr, "scan_loc_info(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_magic(fp *os.File, p *char_magic) {
	if zero_check(*p) {
		return
	}

	fprintf(fp, "CM\n")

	if p.magician != FALSE {
		fprintf(fp, " im %d\n", p.magician)
	}

	if p.max_aura != FALSE {
		fprintf(fp, " ma %d\n", p.max_aura)
	}

	if p.cur_aura != FALSE {
		fprintf(fp, " ca %d\n", p.cur_aura)
	}

	if p.ability_shroud != FALSE {
		fprintf(fp, " as %d\n", p.ability_shroud)
	}

	if p.hinder_meditation != FALSE {
		fprintf(fp, " hm %d\n", p.hinder_meditation)
	}

	if p.quick_cast != FALSE {
		fprintf(fp, " qc %d\n", p.quick_cast)
	}

	if p.aura_reflect != FALSE {
		fprintf(fp, " rb %d\n", p.aura_reflect)
	}

	if p.hide_self != FALSE {
		fprintf(fp, " hs %d\n", p.hide_self)
	}

	if p.hide_mage != FALSE {
		fprintf(fp, " cm %d\n", p.hide_mage)
	}

	if p.knows_weather != FALSE {
		fprintf(fp, " kw %d\n", p.knows_weather)
	}

	if p.swear_on_release != FALSE {
		fprintf(fp, " sr %d\n", p.swear_on_release)
	}

	//#if 0
	//  /* Pledgee might have died! */
	//    if (valid_box(p.pledge)) box_print(fp, " pl ", p.pledge);
	//#endif

	if valid_box(p.project_cast) {
		box_print(fp, []byte(" pc "), p.project_cast)
	}
	box_print(fp, []byte(" ar "), p.auraculum)
	box_print(fp, []byte(" ot "), p.token) /* our token artifact */
	known_print(fp, []byte(" vi "), p.visions)
}

func scan_magic(p *char_magic, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)
		switch c {
		case `im`:
			p.magician = atoi_b(t)
			break
		case `ma`:
			p.max_aura = atoi_b(t)
			break
		case `ca`:
			p.cur_aura = atoi_b(t)
			break
		case `as`:
			p.ability_shroud = atoi_b(t)
			break
		case `hm`:
			p.hinder_meditation = atoi_b(t)
			break
		case `pc`:
			p.project_cast = box_scan(t)
			break
		case `qc`:
			p.quick_cast = atoi_b(t)
			break
		case `ot`:
			p.token = box_scan(t)
			break
		case `pl`:
			box_scan(t)
			break
		case `ar`:
			p.auraculum = box_scan(t)
			break
		case `rb`:
			p.aura_reflect = atoi_b(t)
			break
		case `hs`:
			p.hide_self = atoi_b(t)
			break
		case `cm`:
			p.hide_mage = atoi_b(t)
			break
		case `sr`:
			p.swear_on_release = atoi_b(t)
			break
		case `kw`:
			p.knows_weather = atoi_b(t)
			break

		case `pr`: /* Old "prayer" flag */
			break

		case `vi`:
			p.visions = known_scan(t, p.visions, box_num)
			break

		default:
			fprintf(os.Stderr, "scan_magic(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_artifact(fp *os.File, p *entity_artifact) {
	if zero_check(*p) {
		return
	}

	fprintf(fp, "AR\n")

	if p.type_ != FALSE {
		fprintf(fp, " ty %d\n", p.type_)
	}

	if p.param1 != FALSE {
		fprintf(fp, " p1 %d\n", p.param1)
	}

	if p.param2 != FALSE {
		fprintf(fp, " p2 %d\n", p.param2)
	}

	if p.uses != FALSE {
		fprintf(fp, " us %d\n", p.uses)
	}
}

func scan_artifact(p *entity_artifact, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `ty`:
			p.type_ = atoi_b(t)
			break
		case `p1`:
			p.param1 = atoi_b(t)
			break
		case `p2`:
			p.param2 = atoi_b(t)
			break
		case `us`:
			p.uses = atoi_b(t)
			break
		default:
			fprintf(os.Stderr, "scan_artifact(%d):  bad line: %s\n",
				box_num, line)
		}
		advance()
	}
}

func accept_print_sup(fp *os.File, p *accept_ent) {
	/*
	 *  Trim out obviously bad "accepts".
	 *
	 */
	if p.from_who != 0 && !valid_box(p.from_who) {
		return
	}
	fprintf(fp, " ac %d %d %d\n", p.item, p.from_who, p.qty)
}

func accept_print(fp *os.File, p *entity_char) {
	for i := 0; i < len(p.accept); i++ {
		accept_print_sup(fp, p.accept[i])
	}
}

func print_char(fp *os.File, p *entity_char) {

	fprintf(fp, "CH\n")

	box_print(fp, []byte(" ni "), p.unit_item)
	box_print(fp, []byte(" lo "), p.unit_lord)

	if p.health != FALSE {
		fprintf(fp, " he %d\n", p.health)
	}

	if p.sick != FALSE {
		fprintf(fp, " si %d\n", p.sick)
	}

	if p.loy_kind != FALSE {
		fprintf(fp, " lk %d\n", p.loy_kind)
	}
	if p.loy_rate != FALSE {
		fprintf(fp, " lr %d\n", p.loy_rate)
	}

	skill_list_print(fp, []byte(" sl\t"), p.skills)

	//#if 0
	//    effect_list_print(fp, " el\t", p.effects);
	//#endif

	if p.prisoner != FALSE {
		fprintf(fp, " pr %d\n", p.prisoner)
	}

	if p.moving != FALSE {
		fprintf(fp, " mo %d\n", p.moving)
	}

	if p.behind != FALSE {
		fprintf(fp, " bh %d\n", p.behind)
	}

	if p.guard != FALSE {
		fprintf(fp, " gu %d\n", p.guard)
	}

	if p.time_flying != FALSE {
		fprintf(fp, " tf %d\n", p.time_flying)
	}

	if p.break_point != FALSE {
		fprintf(fp, " bp %d\n", p.break_point)
	}

	if p.personal_break_point != FALSE {
		fprintf(fp, " pb %d\n", p.personal_break_point)
	}

	if p.rank != FALSE {
		fprintf(fp, " ra %d\n", p.rank)
	}

	if p.attack != FALSE {
		fprintf(fp, " at %d\n", p.attack)
	}

	if p.defense != FALSE {
		fprintf(fp, " df %d\n", p.defense)
	}

	if p.missile != FALSE {
		fprintf(fp, " mi %d\n", p.missile)
	}

	if p.npc_prog != FALSE {
		fprintf(fp, " po %d\n", p.npc_prog)
	}

	if p.guild != FALSE {
		fprintf(fp, " gl %d\n", p.guild)
	}

	if p.pay != FALSE {
		fprintf(fp, " pa %d\n", p.pay)
	}

	boxlist_print(fp, []byte(" ct "), p.contact)

	olytime_print(fp, []byte(" dt "), &p.death_time)

	/*
	 *  Religion stuff...
	 *
	 */
	if p.religion.priest != FALSE {
		fprintf(fp, " pi %d\n", p.religion.priest)
	}

	if p.religion.piety != FALSE {
		fprintf(fp, " pt %d\n", p.religion.piety)
	}

	boxlist_print(fp, []byte(" fl "), p.religion.followers)

	if len(p.accept) != 0 {
		accept_print(fp, p)
	}
}

func accept_scan(s []byte, pp *entity_char) {
	p := &accept_ent{}
	//sscanf(s, "%d %d %d", &p.item, &p.from_who, &p.qty);
	for i, f := range bytes.Fields(s) {
		switch i {
		case 0:
			p.item = atoi_b(f)
		case 1:
			p.from_who = atoi_b(f)
		case 2:
			p.qty = atoi_b(f)
		}
	}
	pp.accept = append(pp.accept, p)
}

func scan_char(p *entity_char, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `ni`:
			p.unit_item = box_scan(t)
			break
		case `lo`:
			p.unit_lord = box_scan(t)
			break
		case `he`:
			p.health = atoi_b(t)
			break
		case `si`:
			p.sick = atoi_b(t)
			break
		case `pr`:
			p.prisoner = atoi_b(t)
			break
		case `mo`:
			p.moving = atoi_b(t)
			break
		case `bh`:
			p.behind = atoi_b(t)
			break
		case `lk`:
			p.loy_kind = atoi_b(t)
			break
		case `lr`:
			p.loy_rate = atoi_b(t)
			break
		case `gu`:
			p.guard = atoi_b(t)
			break
		case `tf`:
			p.time_flying = atoi_b(t)
			break
		case `bp`:
			p.break_point = atoi_b(t)
			break
		case `pb`:
			p.personal_break_point = atoi_b(t)
			break
		case `ra`:
			p.rank = atoi_b(t)
			break
		case `at`:
			p.attack = atoi_b(t)
			break
		case `df`:
			p.defense = atoi_b(t)
			break
		case `mi`:
			p.missile = atoi_b(t)
			break
		case `po`:
			p.npc_prog = atoi_b(t)
			break
		case `gl`:
			p.guild = convert_skill(atoi_b(t))
			break
		case `pa`:
			p.pay = atoi_b(t)
			break
		case `pi`:
			p.religion.priest = atoi_b(t)
			break
		case `pt`:
			p.religion.piety = atoi_b(t)
			break

		case `ct`:
			p.contact = boxlist_scan(t, box_num, p.contact)
			break

		case `fl`:
			p.religion.followers = boxlist_scan(t, box_num, p.religion.followers)
			break

		case `sl`:
			p.skills = skill_list_scan(t, p.skills, box_num)
			break

			//#if 0
			//         case `el`:
			//            effect_list_scan(t, &p.effects);
			//            break;
			//#endif

		case `dt`:
			olytime_scan(t, &p.death_time)
			break

		case `ac`:
			accept_scan(t, p)
			break

		default:
			fprintf(os.Stderr, "scan_char(%d):  bad line: %s\n",
				box_num, line)
		}
		advance()
	}
}

/*
 *  Mine Info Functions
 *  Fri Jan 24 12:35:16 1997 -- Scott Turner
 *
 */
func mine_info_print(fp *os.File, header []byte, m *entity_mine) {
	fprintf(fp, string(header))
	for i := 0; i < MINE_MAX; i++ {
		if item_list_print(fp, []byte("    ml\t"), m.mc[i].items) == 0 {
			fputs("    ml\t\n", fp)
		}
		fprintf(fp, "    ms %d\n", m.shoring[i])
	}
	fputs("\n", fp)
}

func mine_info_scan(s []byte, l *entity_mine, bn int) *entity_mine {
	/* Skip the initial line. */
	s = bytes.TrimSpace(readlin_ew())

	m := &entity_mine{}

	for i := 0; i < MINE_MAX; i++ {
		if len(s) > 4 {
			m.mc[i].items = item_list_scan(s[2:], m.mc[i].items, bn)
		}
		s = bytes.TrimSpace(readlin_ew())
		// sscanf(s, "ms %d", &m.shoring[i]);
		for idx, f := range bytes.Fields(s) {
			switch idx {
			case 0:
				m.shoring[i] = atoi_b(f) // todo: change to append
			}
		}
		s = bytes.TrimSpace(readlin_ew())
	}

	return m
}

func print_loc(fp *os.File, p *entity_loc) {

	if zero_check(*p) {
		return
	}

	fprintf(fp, "LO\n")

	boxlist0_print(fp, []byte(" pd "), p.prov_dest)

	if p.hidden != FALSE {
		fprintf(fp, " hi %d\n", p.hidden)
	}

	if p.shroud != FALSE {
		fprintf(fp, " sh %d\n", p.shroud)
	}

	//#if 0
	//   if (p.barrier) {
	//        fprintf(fp, " ba %d\n", p.barrier);
	//   }
	//#endif

	if p.dist_from_sea != FALSE {
		fprintf(fp, " ds %d\n", p.dist_from_sea)
	}

	if p.dist_from_swamp != FALSE {
		fprintf(fp, " dw %d\n", p.dist_from_swamp)
	}

	if p.dist_from_gate != FALSE {
		fprintf(fp, " dg %d\n", p.dist_from_gate)
	}

	if p.sea_lane != FALSE {
		fprintf(fp, " sl %d\n", p.sea_lane)
	}

	if p.tax_rate != FALSE {
		fprintf(fp, " tr %d\n", p.tax_rate)
	}

	if p.mine_info != nil {
		mine_info_print(fp, []byte(" mi\n"), p.mine_info)
	}

	if p.control.weight != FALSE {
		fprintf(fp, " fw %d\n", p.control.weight)
	}

	if p.control.men != FALSE {
		fprintf(fp, " fm %d\n", p.control.men)
	}

	if p.control.nobles != FALSE {
		fprintf(fp, " fn %d\n", p.control.nobles)
	}

	if p.control.closed {
		fprintf(fp, " cd %d\n", p.control.closed)
	}

	box_print(fp, []byte(" ng "), p.near_grave)
	//#if 0
	//    effect_list_print(fp, " el\t", p.effects);
	//#endif
}

func scan_loc(p *entity_loc, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `hi`:
			p.hidden = atoi_b(t)
			break
		case `sh`:
			p.shroud = atoi_b(t)
			break
			//#if 0
			//  case `ba`:
			//    p.barrier = atoi(t);        break;
			//#endif
		case `ds`:
			p.dist_from_sea = atoi_b(t)
			break
		case `dw`:
			p.dist_from_swamp = atoi_b(t)
			break
		case `dg`:
			p.dist_from_gate = atoi_b(t)
			break
		case `ng`:
			p.near_grave = box_scan(t)
			break
		case `sl`:
			p.sea_lane = atoi_b(t)
			break
		case `tr`:
			p.tax_rate = atoi_b(t)
			break
		case `fw`:
			p.control.weight = atoi_b(t)
			break
		case `fm`:
			p.control.men = atoi_b(t)
			break
		case `fn`:
			p.control.nobles = atoi_b(t)
			break
		case `cd`:
			p.control.closed = atoi_b(t) != FALSE
			break

		case `pd`:
			p.prov_dest = boxlist0_scan(t, box_num, p.prov_dest)
			break

			//#if 0
			//          case `el`:
			//                effect_list_scan(t, &p.effects);
			//                break;
			//#endif

		case `mi`:
			p.mine_info = mine_info_scan(t, p.mine_info, box_num)
			break

		default:
			fprintf(os.Stderr, "scan_loc(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_ship(fp *os.File, p *entity_ship) {
	fprintf(fp, "SP\n")

	if p.hulls != FALSE {
		fprintf(fp, " hu %d\n", p.hulls)
	}

	if p.forts != FALSE {
		fprintf(fp, " fo %d\n", p.forts)
	}

	if p.sails != FALSE {
		fprintf(fp, " sa %d\n", p.sails)
	}

	if p.ports != FALSE {
		fprintf(fp, " po %d\n", p.ports)
	}

	if p.keels != FALSE {
		fprintf(fp, " ke %d\n", p.keels)
	}

	if p.galley_ram != FALSE {
		fprintf(fp, " gr %d\n", p.galley_ram)
	}
}

func scan_ship(p *entity_ship, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {

		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `hu`:
			p.hulls = atoi_b(t)
			break
		case `fo`:
			p.forts = atoi_b(t)
			break
		case `sa`:
			p.sails = atoi_b(t)
			break
		case `po`:
			p.ports = atoi_b(t)
			break
		case `ke`:
			p.keels = atoi_b(t)
			break
		case `gr`:
			p.galley_ram = atoi_b(t)
			break
		}
		advance()
	}
}

func print_subloc(fp *os.File, p *entity_subloc) {

	fprintf(fp, "SL\n")

	boxlist_print(fp, []byte(" te "), p.teaches)

	if p.hp != FALSE {
		fprintf(fp, " hp %d\n", p.hp)
	}

	if p.moat != FALSE {
		fprintf(fp, " mt %d\n", p.moat)
	}

	if p.damage != FALSE {
		fprintf(fp, " da %d\n", p.damage)
	}

	if p.defense != FALSE {
		fprintf(fp, " de %d\n", p.defense)
	}

	/*
	   if (p.capacity) {
	     fprintf(fp, " ca %d\n", p.capacity);
	   }

	   if (p.build_materials)
	       fprintf(fp, " bm %d\n", p.build_materials);

	   if (p.effort_required)
	       fprintf(fp, " er %d\n", p.effort_required);

	   if (p.effort_given)
	       fprintf(fp, " eg %d\n", p.effort_given);

	*/

	if p.moving != FALSE {
		fprintf(fp, " mo %d\n", p.moving)
	}

	/*
	   if (p.galley_ram) {
	     fprintf(fp, " gr %d\n", p.galley_ram);
	   }

	   if (p.shaft_depth)
	       fprintf(fp, " sd %d\n", p.shaft_depth);
	*/

	if p.safe {
		fprintf(fp, " sh %d\n", TRUE) // mdhender: changed from p.safe to TRUE
	}

	if p.major != FALSE {
		fprintf(fp, " mc %d\n", p.major)
	}

	if p.opium_econ != FALSE {
		fprintf(fp, " op %d\n", p.opium_econ)
	}

	if p.loot != FALSE {
		fprintf(fp, " lo %d\n", p.loot)
	}

	if p.prominence != FALSE {
		fprintf(fp, " cp %d\n", p.prominence)
	}

	boxlist_print(fp, []byte(" nc "), p.near_cities)

	boxlist_print(fp, []byte(" lt "), p.link_to)
	boxlist_print(fp, []byte(" lf "), p.link_from)
	/*
	   boxlist_print(fp, " bs ", p.bound_storms);
	*/

	/*
	   if (p.link_when)
	       fprintf(fp, " lw %d\n", p.link_when);

	   if (p.link_open)
	       fprintf(fp, " lp %d\n", p.link_open);
	*/

	if p.guild != FALSE {
		fprintf(fp, " gl %d\n", p.guild)
	}

	if p.control.weight != FALSE {
		fprintf(fp, " fw %d\n", p.control.weight)
	}

	if p.control.men != FALSE {
		fprintf(fp, " fm %d\n", p.control.men)
	}

	if p.control.nobles != FALSE {
		fprintf(fp, " fn %d\n", p.control.nobles)
	}

	if p.control.closed {
		fprintf(fp, " cd %d\n", TRUE) // mdhender: changed to TRUE
	}

	if p.tax_market != FALSE {
		fprintf(fp, " tm %d\n", p.tax_market)
	}

	build_list_print(fp, []byte(" bl\t"), p.builds)
	//#if 0
	//    effect_list_print(fp, " el\t", p.effects);
	//#endif

	if p.entrance_size != FALSE {
		fprintf(fp, " es %d\n", p.entrance_size)
	}
}

func scan_subloc(p *entity_subloc, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `tm`:
			p.tax_market = atoi_b(t)
			break
		case `hp`:
			p.hp = atoi_b(t)
			break
		case `mt`:
			p.moat = atoi_b(t)
			break
		case `da`:
			p.damage = atoi_b(t)
			break
		case `de`:
			p.defense = atoi_b(t)
			break
			/*
			   case `ca`:   p.capacity = atoi(t);            break;
			   case `er`:    p.effort_required = atoi(t);        break;
			   case `eg`:    p.effort_given = atoi(t);        break;
			   case `bm`:    p.build_materials = atoi(t);        break;
			*/
		case `ca`, `er`, `eg`, `bm`:
			break

		case `mo`:
			p.moving = atoi_b(t)
			break
		case `gr`:
			//p.galley_ram = atoi(t);
			break
			//case `sd`:
		//    p.shaft_depth = atoi(t);
		//    break;
		case `sh`:
			p.safe = atoi_b(t) != FALSE
			break
		case `mc`:
			p.major = atoi_b(t)
			break
		case `op`:
			p.opium_econ = atoi_b(t)
			break
		case `lo`:
			p.loot = atoi_b(t)
			break
		case `cp`:
			p.prominence = atoi_b(t)
			break
			/*
			   case `lw`:    p.link_when = atoi(t);            break;
			   case `lp`:    p.link_open = atoi(t);            break;
			*/
		case `gl`:
			p.guild = convert_skill(atoi_b(t))
			break
		case `es`:
			p.entrance_size = atoi_b(t)
			break

		case `lt`:
			p.link_to = boxlist_scan(t, box_num, p.link_to)
			break

		case `lf`:
			p.link_from = boxlist_scan(t, box_num, p.link_from)
			break

		case `te`:
			p.teaches = boxlist_scan(t, box_num, p.teaches)
			break

		case `nc`:
			p.near_cities = boxlist_scan(t, box_num, p.near_cities)
			break

		case `bs`:
			/*
			   boxlist_scan(t, box_num, &(p.bound_storms));
			*/
			break

		case `bl`:
			p.builds = build_list_scan(t, p.builds)
			break

			//#if 0
			//      case `el`:
			//          effect_list_scan(t, &p.effects);
			//          break;
			//#endif

		case `fw`:
			p.control.weight = atoi_b(t)
			break
		case `fm`:
			p.control.men = atoi_b(t)
			break
		case `fn`:
			p.control.nobles = atoi_b(t)
			break
		case `cd`:
			p.control.closed = atoi_b(t) != FALSE
			break

		default:
			fprintf(os.Stderr, "scan_subloc(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}

}

func print_item(fp *os.File, p *entity_item) {

	if zero_check(*p) {
		return
	}

	fprintf(fp, "IT\n")

	if p.trade_good != FALSE {
		fprintf(fp, " tg %d\n", p.trade_good)
	}

	if len(p.plural_name) != 0 {
		fprintf(fp, " pl %s\n", p.plural_name)
	}

	if p.weight != FALSE {
		fprintf(fp, " wt %d\n", p.weight)
	}

	if p.land_cap != FALSE {
		fprintf(fp, " lc %d\n", p.land_cap)
	}

	if p.ride_cap != FALSE {
		fprintf(fp, " rc %d\n", p.ride_cap)
	}

	if p.fly_cap != FALSE {
		fprintf(fp, " fc %d\n", p.fly_cap)
	}

	if p.is_man_item != FALSE {
		fprintf(fp, " mu %d\n", p.is_man_item)
	}

	if p.prominent != FALSE {
		fprintf(fp, " pr %d\n", p.prominent)
	}

	if p.animal != FALSE {
		fprintf(fp, " an %d\n", p.animal)
	}

	if p.attack != FALSE {
		fprintf(fp, " at %d\n", p.attack)
	}

	if p.defense != FALSE {
		fprintf(fp, " df %d\n", p.defense)
	}

	if p.missile != FALSE {
		fprintf(fp, " mi %d\n", p.missile)
	}

	if p.base_price != FALSE {
		fprintf(fp, " bp %d\n", p.base_price)
	}

	if p.capturable != FALSE {
		fprintf(fp, " ca %d\n", p.capturable)
	}

	if p.ungiveable != FALSE {
		fprintf(fp, " ug %d\n", p.ungiveable)
	}

	if p.wild != FALSE {
		fprintf(fp, " wi %d\n", p.wild)
	}

	if p.maintenance != FALSE {
		fprintf(fp, " mt %d\n", p.maintenance)
	}

	if p.npc_split != FALSE {
		fprintf(fp, " sp %d\n", p.npc_split)
	}

	if p.animal_part != FALSE {
		fprintf(fp, " ap %d\n", p.animal_part)
	}

	box_print(fp, []byte(" un "), p.who_has)
}

func scan_item(p *entity_item, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `pl`:
			p.plural_name = string(t)
			break
		case `tg`:
			p.trade_good = atoi_b(t)
			break
		case `wt`:
			p.weight = atoi_b(t)
			break
		case `lc`:
			p.land_cap = atoi_b(t)
			break
		case `rc`:
			p.ride_cap = atoi_b(t)
			break
		case `fc`:
			p.fly_cap = atoi_b(t)
			break
		case `mu`:
			p.is_man_item = atoi_b(t)
			break
		case `pr`:
			p.prominent = atoi_b(t)
			break
		case `an`:
			p.animal = atoi_b(t)
			break
		case `un`:
			p.who_has = box_scan(t)
			break
		case `at`:
			p.attack = atoi_b(t)
			break
		case `df`:
			p.defense = atoi_b(t)
			break
		case `mi`:
			p.missile = atoi_b(t)
			break
		case `bp`:
			p.base_price = atoi_b(t)
			break
		case `ca`:
			p.capturable = atoi_b(t)
			break
		case `ug`:
			p.ungiveable = atoi_b(t)
			break
		case `wi`:
			p.wild = atoi_b(t)
			break
		case `mt`:
			p.maintenance = atoi_b(t)
			break
		case `sp`:
			p.npc_split = atoi_b(t)
			break
		case `ap`:
			p.animal_part = atoi_b(t)
			break

		default:
			fprintf(os.Stderr, "scan_item(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_item_magic(fp *os.File, p *item_magic) {

	if zero_check(*p) {
		return
	}

	fprintf(fp, "IM\n")

	if p.religion != FALSE {
		fprintf(fp, " rl %d\n", p.religion)
	}

	if p.aura != FALSE {
		fprintf(fp, " au %d\n", p.aura)
	}

	if p.curse_loyalty != FALSE {
		fprintf(fp, " cl %d\n", p.curse_loyalty)
	}

	if p.cloak_region != FALSE {
		fprintf(fp, " cr %d\n", p.cloak_region)
	}

	if p.cloak_creator != FALSE {
		fprintf(fp, " cc %d\n", p.cloak_creator)
	}

	if p.use_key != FALSE {
		fprintf(fp, " uk %d\n", p.use_key)
	}

	if p.quick_cast != FALSE {
		fprintf(fp, " qc %d\n", p.quick_cast)
	}

	if p.attack_bonus != FALSE {
		fprintf(fp, " ab %d\n", p.attack_bonus)
	}

	if p.defense_bonus != FALSE {
		fprintf(fp, " db %d\n", p.defense_bonus)
	}

	if p.missile_bonus != FALSE {
		fprintf(fp, " mb %d\n", p.missile_bonus)
	}

	if p.aura_bonus != FALSE {
		fprintf(fp, " ba %d\n", p.aura_bonus)
	}

	if p.token_num != FALSE {
		fprintf(fp, " tn %d\n", p.token_num)
	}

	if p.orb_use_count != FALSE {
		fprintf(fp, " oc %d\n", p.orb_use_count)
	}

	box_print(fp, []byte(" ti "), p.token_ni)

	//#if 0
	//    box_print(fp, " rc ", p.region_created);
	//#endif

	if valid_box(p.project_cast) {
		box_print(fp, []byte(" pc "), p.project_cast)
	}
	box_print(fp, []byte(" ct "), p.creator)
	if p.lore != FALSE {
		fprintf(fp, " lo %d\n", p.lore)
	}

	boxlist_print(fp, []byte(" mu "), p.may_use)
	boxlist_print(fp, []byte(" ms "), p.may_study)
}

func scan_item_magic(p *item_magic, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `au`:
			p.aura = atoi_b(t)
			break
		case `cl`:
			p.curse_loyalty = atoi_b(t)
			break
		case `cr`:
			p.cloak_region = atoi_b(t)
			break
		case `cc`:
			p.cloak_creator = atoi_b(t)
			break
		case `uk`:
			p.use_key = atoi_b(t)
			break
		case `rc`:
			/* p.region_created = box_scan(t); */
			break
		case `pc`:
			p.project_cast = box_scan(t)
			break
		case `ct`:
			p.creator = box_scan(t)
			break
		case `lo`:
			p.lore = box_scan(t)
			break
		case `qc`:
			p.quick_cast = atoi_b(t)
			break
		case `ab`:
			p.attack_bonus = atoi_b(t)
			break
		case `db`:
			p.defense_bonus = atoi_b(t)
			break
		case `mb`:
			p.missile_bonus = atoi_b(t)
			break
		case `ba`:
			p.aura_bonus = atoi_b(t)
			break
		case `tn`:
			p.token_num = atoi_b(t)
			break
		case `ti`:
			p.token_ni = atoi_b(t)
			break
		case `oc`:
			p.orb_use_count = atoi_b(t)
			break
		case `rl`:
			p.religion = atoi_b(t)
			break

		case `mu`:
			p.may_use = boxlist_scan(t, box_num, (p.may_use))
			break

		case `ms`:
			p.may_study = boxlist_scan(t, box_num, (p.may_study))
			break

		default:
			fprintf(os.Stderr, "scan_item_magic(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_player(fp *os.File, p *entity_player) {

	fprintf(fp, "PL\n")

	if len(p.full_name) != 0 {
		fprintf(fp, " fn %s\n", p.full_name)
	}

	if len(p.email) != 0 {
		fprintf(fp, " em %s\n", p.email)
	}

	if len(p.vis_email) != 0 {
		fprintf(fp, " ve %s\n", p.vis_email)
	}

	if len(p.password) != 0 {
		fprintf(fp, " pw %s\n", p.password)
	}

	if p.noble_points != FALSE {
		fprintf(fp, " np %d\n", p.noble_points)
	}

	if p.first_turn != FALSE {
		fprintf(fp, " ft %d\n", p.first_turn)
	}

	if p.format != FALSE {
		fprintf(fp, " fo %d\n", p.format)
	}

	if len(p.rules_path) != 0 {
		fprintf(fp, " rp %s\n", p.rules_path)
	}

	if len(p.db_path) != 0 {
		fprintf(fp, " db%s\n", p.db_path)
	}

	if p.notab {
		fprintf(fp, " nt %d\n", TRUE) // mdhender: changed to TRUE
	}

	if p.first_tower != FALSE {
		fprintf(fp, " tf %d\n", p.first_tower)
	}

	if p.split_lines != FALSE {
		fprintf(fp, " sl %d\n", p.split_lines)
	}

	if p.split_bytes != FALSE {
		fprintf(fp, " sb %d\n", p.split_bytes)
	}

	if p.sent_orders != FALSE {
		fprintf(fp, " so %d\n", p.sent_orders)
	}

	if p.dont_remind != FALSE {
		fprintf(fp, " dr %d\n", p.dont_remind)
	}

	if p.compuserve {
		fprintf(fp, " ci %d\n", TRUE) // mdhender: changed to TRUE
	}

	if p.broken_mailer != FALSE {
		fprintf(fp, " bm %d\n", p.broken_mailer)
	}

	if p.last_order_turn != FALSE {
		fprintf(fp, " lt %d\n", p.last_order_turn)
	}

	if p.nation != FALSE {
		fprintf(fp, " na %d\n", p.nation)
	}

	if p.magic != FALSE {
		fprintf(fp, " ma %d\n", p.magic)
	}

	if p.jump_start != FALSE {
		fprintf(fp, " js %d\n", p.jump_start)
	}

	known_print(fp, []byte(" kn "), p.known)
	boxlist_print(fp, []byte(" un "), p.units)
	boxlist_print(fp, []byte(" uf "), p.unformed)
	admit_print(fp, p)
}

func scan_player(p *entity_player, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `fn`:
			p.full_name = string(t)
			break
		case `em`:
			p.email = string(t)
			break
		case `ve`:
			p.vis_email = string(t)
			break
		case `pw`:
			p.password = string(t)
			break
		case `np`:
			p.noble_points = atoi_b(t)
			break
		case `ft`:
			p.first_turn = atoi_b(t)
			break
		case `fo`:
			p.format = atoi_b(t)
			break
		case `rp`:
			p.rules_path = string(t)
			break
		case `fp`:
		case `db`:
			p.db_path = string(t)
			break
		case `nt`:
			p.notab = atoi_b(t) != FALSE
			break
		case `tf`:
			p.first_tower = atoi_b(t)
			break
		case `so`:
			p.sent_orders = atoi_b(t)
			break
		case `lt`:
			p.last_order_turn = atoi_b(t)
			break
		case `sl`:
			p.split_lines = atoi_b(t)
			break
		case `sb`:
			p.split_bytes = atoi_b(t)
			break
		case `ci`:
			p.compuserve = atoi_b(t) != FALSE
			break
		case `ti`: /* ignore it... */
			break
		case `bm`:
			p.broken_mailer = atoi_b(t)
			break
		case `dr`:
			p.dont_remind = atoi_b(t)
			break
		case `na`:
			{
				p.nation = atoi_b(t)
				/* temp fix */ // todo: fix temp fix
				if p.nation <= 1002 && p.nation >= 1000 {
					p.nation -= 3
				}
				break
			}
		case `ma`:
			p.magic = atoi_b(t)
			break
		case `js`:
			p.jump_start = atoi_b(t)
			break

		case `kn`:
			p.known = known_scan(t, p.known, box_num)
			break

		case `un`:
			p.units = boxlist_scan(t, box_num, (p.units))
			break

		case `uf`:
			p.unformed = boxlist_scan(t, box_num, (p.unformed))
			break

		case `am`:
			admit_scan(t, box_num, p)
			break

		default:
			fprintf(os.Stderr, "scan_player(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_religion(fp *os.File, p *entity_religion_skill) {
	fprintf(fp, " RL\n")
	fprintf(fp, " na %s\n", p.name)
	fprintf(fp, " st %d\n", p.strength)
	fprintf(fp, " wk %d\n", p.weakness)
	fprintf(fp, " pl %d\n", p.plant)
	fprintf(fp, " an %d\n", p.animal)
	fprintf(fp, " tr %d\n", p.terrain)
	fprintf(fp, " hp %d\n", p.high_priest)
	fprintf(fp, " b0 %d\n", p.bishops[0])
	fprintf(fp, " b1 %d\n", p.bishops[1])
	fprintf(fp, " ER\n")
}

func scan_religion(box_num int) *entity_religion_skill {
	newt := &entity_religion_skill{}

	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		if c == `ER` {
			break
		}
		t := t_string(line)

		switch c {
		case `na`:
			newt.name = string(t)
			break
		case `st`:
			newt.strength = atoi_b(t)
			break
		case `wk`:
			newt.weakness = atoi_b(t)
			break
		case `pl`:
			newt.plant = atoi_b(t)
			break
		case `an`:
			newt.animal = atoi_b(t)
			break
		case `tr`:
			newt.terrain = atoi_b(t)
			break
		case `hp`:
			newt.high_priest = atoi_b(t)
			break
		case `b0`:
			newt.bishops[0] = atoi_b(t)
			break
		case `b1`:
			newt.bishops[1] = atoi_b(t)
			break
		}
		advance()
	}
	return newt
}

func print_skill(fp *os.File, p *entity_skill) {
	fprintf(fp, "SK\n")
	fprintf(fp, " tl %d\n", p.time_to_learn)
	fprintf(fp, " tu %d\n", p.time_to_use)
	/*
	 *  Flags
	 *  Mon Oct 21 12:50:46 1996 -- Scott Turner
	 *
	 */
	for i := 0; i < MAX_FLAGS; i++ {
		if (p.flags & (1 << i)) != 0 {
			switch 1 << i {
			case IS_POLLED:
				fprintf(fp, " pl 1\n")
				break
			case REQ_HOLY_SYMBOL:
				fprintf(fp, " hs 1\n")
				break
			case REQ_HOLY_PLANT:
				fprintf(fp, " hp 1\n")
				break
			case COMBAT_SKILL:
				fprintf(fp, " cs 1\n")
				break
			}
		}
	}

	box_print(fp, []byte(" rs "), p.required_skill)
	boxlist_print(fp, []byte(" of "), p.offered)
	boxlist_print(fp, []byte(" re "), p.research)
	boxlist_print(fp, []byte(" gl "), p.guild)
	req_list_print(fp, []byte(" rq\t"), p.req)
	box_print(fp, []byte(" pr "), p.produced)

	if p.practice_cost != FALSE {
		fprintf(fp, " pc %d\n", p.practice_cost)
	}
	if p.practice_time != FALSE {
		fprintf(fp, " tp %d\n", p.practice_time)
	}
	if p.practice_prog != FALSE {
		fprintf(fp, " pp %d\n", p.practice_prog)
	}
	if p.np_req != FALSE {
		fprintf(fp, " np %d\n", p.np_req)
	}
	if p.no_exp != FALSE {
		fprintf(fp, " ne %d\n", p.no_exp)
	}
	if p.piety != FALSE {
		fprintf(fp, " pt %d\n", p.piety)
	}
	if p.religion_skill != nil {
		print_religion(fp, p.religion_skill)
	}
}

func scan_skill(p *entity_skill, box_num int) {
	advance()
	p.flags = 0
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `tl`:
			p.time_to_learn = atoi_b(t)
			break
		case `tu`:
			p.time_to_use = atoi_b(t)
			break
		case `ne`:
			p.no_exp = atoi_b(t)
			break
		case `np`:
			p.np_req = atoi_b(t)
			break
		case `rs`:
			p.required_skill = box_scan(t)
			break
		case `pr`:
			p.produced = box_scan(t)
			break
		case `pt`:
			p.piety = atoi_b(t)
			break
		case `pc`:
			p.practice_cost = atoi_b(t)
			break
		case `tp`:
			p.practice_time = atoi_b(t)
			break
		case `pp`:
			p.practice_prog = atoi_b(t)
			break

		case `of`:
			p.offered = boxlist_scan(t, box_num, (p.offered))
			break

		case `re`:
			p.research = boxlist_scan(t, box_num, (p.research))
			break

		case `gl`:
			p.guild = boxlist_scan(t, box_num, (p.guild))
			break

		case `rq`:
			p.req = req_list_scan(t, p.req, box_num)
			break

		case `RL`:
			p.religion_skill = scan_religion(box_num)
			break

		case `pl`:
			p.flags |= IS_POLLED
			break

		case `hs`:
			p.flags |= REQ_HOLY_SYMBOL
			break

		case `hp`:
			p.flags |= REQ_HOLY_PLANT
			break

		case `cs`:
			p.flags |= COMBAT_SKILL
			break

		default:
			fprintf(os.Stderr, "scan_skill(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func req_list_scan(s []byte, l []*req_ent, box_num int) []*req_ent { panic("!implemented") }

func print_command(fp *os.File, p *command) {
	if p.cmd == 0 {
		return
	}

	fprintf(fp, "CO\n")
	fprintf(fp, " li %s\n", p.line)

	fprintf(fp, " ar %d %d %d %d %d %d %d %d %d\n", p.a, p.b, p.c, p.d, p.e, p.f, p.g, p.h, p.i)

	/*
	   fprintf(fp, " sv %d %d %d %d %d %d %d\n",
	       p.v.direction,
	       p.v.destination,
	       p.v.road,
	       p.v.dest_hidden,
	       p.v.distance,
	       p.v.orig,
	       p.v.orig_hidden);
	*/

	if p.state != FALSE {
		fprintf(fp, " cs %d\n", p.state)
	}

	if p.wait != FALSE {
		fprintf(fp, " wa %d\n", p.wait)
	}

	if p.status != FALSE {
		fprintf(fp, " st %d\n", p.status)
	}

	if p.use_skill != FALSE {
		box_print(fp, []byte(" us "), p.use_skill)
	}

	if p.use_exp != FALSE {
		fprintf(fp, " ue %d\n", p.use_exp)
	}

	if p.days_executing != FALSE {
		fprintf(fp, " de %d\n", p.days_executing)
	}

	if p.poll != FALSE {
		fprintf(fp, " po %d\n", p.poll)
	}

	if p.pri != FALSE {
		fprintf(fp, " pr %d\n", p.pri)
	}

	if p.inhibit_finish {
		fprintf(fp, " if %d\n", TRUE) // mdhender: changed to TRUE
	}

	/*
	 *  Thu Oct 24 15:05:21 1996 -- Scott Turner
	 *
	 *  This must now be saved, since it can be > 1
	 *
	 */
	if p.second_wait != FALSE {
		fprintf(fp, " sw %d\n", p.second_wait)
	}
}

func scan_command(p *command, box_num int) {
	p.who = box_num

	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `li`:
			if !oly_parse_cmd(p, t) {
				fprintf(os.Stderr, "scan_command(%d): bad cmd %s\n", box_num, t)
			}
			break

		case `cs`:
			p.state = atoi_b(t)
			break
		case `wa`:
			p.wait = atoi_b(t)
			break
		case `st`:
			p.status = atoi_b(t)
			break
		case `de`:
			p.days_executing = atoi_b(t)
			break
		case `po`:
			p.poll = atoi_b(t)
			break
		case `pr`:
			p.pri = atoi_b(t)
			break
		case `if`:
			p.inhibit_finish = atoi_b(t) != FALSE
			break
		case `us`:
			p.use_skill = box_scan(t)
			break
		case `ue`:
			p.use_exp = atoi_b(t)
			break
		case `sw`:
			p.second_wait = atoi_b(t)
			break

		case `ar`:
			//sscanf(t, "%d %d %d %d %d %d %d %d %d", &p.a, &p.b, &p.c, &p.d, &p.e, &p.f, &p.g, &p.h, &p.i);
			for i, f := range bytes.Fields(t) {
				switch i {
				case 0:
					p.a = convert_skill(atoi_b(f)) // todo: temp fix
				case 1:
					p.b = convert_skill(atoi_b(f)) // todo: temp fix
				case 2:
					p.c = atoi_b(f)
				case 3:
					p.d = atoi_b(f)
				case 4:
					p.e = atoi_b(f)
				case 5:
					p.f = atoi_b(f)
				case 6:
					p.g = atoi_b(f)
				case 7:
					p.h = atoi_b(f)
				case 8:
					p.i = atoi_b(f)
				}
			}
			break

			//case `sv`:
			//    sscanf(t, "%d %d %d %d %d %d %d", &p.v.direction, &p.v.destination, &p.v.road, &p.v.dest_hidden, &p.v.distance, &p.v.orig, &p.v.orig_hidden);

		default:
			fprintf(os.Stderr, "scan_command(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_gate(fp *os.File, p *entity_gate) {

	fprintf(fp, "GA\n")

	box_print(fp, []byte(" tl "), p.to_loc)

	if p.notify_jumps != FALSE {
		box_print(fp, []byte(" nj "), p.notify_jumps)
	}

	if p.notify_unseal != FALSE {
		box_print(fp, []byte(" nu "), p.notify_unseal)
	}

	if p.seal_key != FALSE {
		fprintf(fp, " sk %d\n", p.seal_key)
	}

	if p.road_hidden != FALSE {
		fprintf(fp, " rh %d\n", p.road_hidden)
	}
}

func scan_gate(p *entity_gate, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `tl`:
			p.to_loc = box_scan(t)
			break
		case `nj`:
			p.notify_jumps = box_scan(t)
			break
		case `nu`:
			p.notify_unseal = box_scan(t)
			break
		case `sk`:
			p.seal_key = atoi_b(t)
			break
		case `rh`:
			p.road_hidden = atoi_b(t)
			break

		default:
			fprintf(os.Stderr, "scan_gate(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

func print_misc(fp *os.File, p *entity_misc) {

	if zero_check(*p) {
		return
	}

	fprintf(fp, "MI\n")

	if valid_box(p.summoned_by) {
		box_print(fp, []byte(" sb "), p.summoned_by)
	}

	if p.npc_dir != FALSE {
		fprintf(fp, " di %d\n", p.npc_dir)
	}

	if p.npc_created != FALSE {
		fprintf(fp, " mc %d\n", p.npc_created)
	}

	if p.mine_delay != FALSE {
		fprintf(fp, " md %d\n", p.mine_delay)
	}

	if p.storm_str != FALSE {
		fprintf(fp, " ss %d\n", p.storm_str)
	}

	if p.cmd_allow != FALSE {
		fprintf(fp, " ca %c\n", p.cmd_allow)
	}

	box_print(fp, []byte(" gc "), p.garr_castle)
	box_print(fp, []byte(" mh "), p.npc_home)
	box_print(fp, []byte(" co "), p.npc_cookie)
	box_print(fp, []byte(" ov "), p.only_vulnerable)
	if p.old_lord != FALSE && valid_box(p.old_lord) {
		box_print(fp, []byte(" ol "), p.old_lord)
	}
	box_print(fp, []byte(" bs "), p.bind_storm)

	if len(p.save_name) != 0 {
		fprintf(fp, " sn %s\n", p.save_name)
	}

	if len(p.display) != 0 {
		fprintf(fp, " ds %s\n", p.display)
	}

	known_print(fp, []byte(" nm "), p.npc_memory)
}

func scan_misc(p *entity_misc, box_num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `di`:
			p.npc_dir = atoi_b(t)
			break
		case `mc`:
			p.npc_created = atoi_b(t)
			break
		case `md`:
			p.mine_delay = atoi_b(t)
			break
		case `ss`:
			p.storm_str = atoi_b(t)
			break
		case `mh`:
			p.npc_home = box_scan(t)
			break
		case `gc`:
			p.garr_castle = box_scan(t)
			break
		case `sb`:
			p.summoned_by = box_scan(t)
			break
		case `co`:
			p.npc_cookie = box_scan(t)
			break
		case `ov`:
			p.only_vulnerable = box_scan(t)
			break
		case `bs`:
			p.bind_storm = box_scan(t)
			break
		case `ol`:
			p.old_lord = box_scan(t)
			break
		case `sn`:
			p.save_name = string(t)
			break
		case `ds`:
			p.display = string(t)
			break
		case `ca`:
			p.cmd_allow = t[0]
			break

		case `nm`:
			p.npc_memory = known_scan(t, p.npc_memory, box_num)
			break

		default:
			fprintf(os.Stderr, "scan_misc(%d):  bad line: %s\n", box_num, line)
		}
		advance()
	}
}

/*
 *  Thu Jan 21 07:15:54 1999 -- Scott Turner
 *
 *  Full-fledged type version.
 *
 */
func scan_nation2(n *entity_nation, num int) {
	advance()
	for line != nil && len(line) != 0 && iswhite(line[0]) {
		if line[0] == '#' { // todo: should this check be above?
			continue
		}

		line = line[1:]
		c := linehash(line)
		t := t_string(line)

		switch c {
		case `na`:
			n.name = string(t)
			break
		case `ci`:
			n.citizen = string(t)
			break
		case `cp`:
			n.capital = box_scan(t)
			break
		case `wi`:
			n.win = atoi_b(t)
			break
		case `pl`:
			n.player_limit = atoi_b(t)
			break
		case `js`:
			n.jump_start = atoi_b(t)
			break
		case `ne`:
			n.neutral = atoi_b(t) != FALSE
			break
		case `ps`:
			n.proscribed_skills = boxlist_scan(t, -1, (n.proscribed_skills))
			break

		default:
			fprintf(os.Stderr, "scan_nation:  bad line: %s\n", line)
		}
		advance()
	}
}

func load_box(n int) {
	/*
	 *  The fast scan of libdir/master should have allocated all of the boxes
	 *  If one was later added manually, it won't have been pre-allocated by
	 *  the master fast scan.  Remove master file and let io.c do the slow
	 *  scan to allocate all the boxes.  Save the database to recreate master.
	 */

	if !valid_box(n) {
		fprintf(os.Stderr, "Unforseen box %d found in load phase.\n", n)
		fprintf(os.Stderr, "Remove %s/master and retry.\n", libdir)
		panic("invalid input")
	}

	//p := bx[n];
	ext_boxnum = n

	advance()
	for len(line) != 0 {
		if line[0] == '#' { // todo: should this check be above?
			advance()
			continue
		}

		c := linehash(line)
		t := t_string(line)

		switch c {
		case `na`:
			set_name(n, string(t))
			advance()
			break
		case `il`:
			bx[n].items = item_list_scan(t, bx[n].items, n)
			advance()
			break
		case `tl`:
			bx[n].trades = trade_list_scan(t, bx[n].trades, n)
			advance()
			break
		case `an`:
			p_disp(n).neutral = boxlist_scan(t, n, (p_disp(n).neutral))
			advance()
			break
		case `ad`:
			p_disp(n).defend = boxlist_scan(t, n, (p_disp(n).defend))
			advance()
			break
		case `ah`:
			p_disp(n).hostile = boxlist_scan(t, n, (p_disp(n).hostile))
			advance()
			break
		case `el`:
			bx[n].effects = effect_list_scan(t, bx[n].effects)
			advance()
			break
		case `CH`:
			scan_char(p_char(n), n)
			break
		case `CM`:
			scan_magic(p_magic(n), n)
			break
		case `LI`:
			scan_loc_info(p_loc_info(n), n)
			break
		case `LO`:
			scan_loc(p_loc(n), n)
			break
		case `SL`:
			scan_subloc(p_subloc(n), n)
			break
		case `IT`:
			scan_item(p_item(n), n)
			break
		case `PL`:
			scan_player(p_player(n), n)
			break
		case `SK`:
			scan_skill(p_skill(n), n)
			break
		case `GA`:
			scan_gate(p_gate(n), n)
			break
		case `MI`:
			scan_misc(p_misc(n), n)
			break
		case `IM`:
			scan_item_magic(p_item_magic(n), n)
			break
		case `AR`:
			scan_artifact(p_item_artifact(n), n)
			break
		case `CO`:
			scan_command(p_command(n), n)
			break
		case `SP`:
			scan_ship(p_ship(n), n)
			break
		case `NA`:
			scan_nation2(p_nation(n), n)
			break
		default:
			fprintf(os.Stderr, "load_box(%d):  bad line: %s\n", n, line)
			advance()
			for line != nil && len(line) != 0 && iswhite(line[0]) {
				advance()
			}
			continue
		}
	}

	if line != nil {
		assert(len(line) == 0)
		advance() /* advance over blank separating line */
	}
}

func print_nation(fp *os.File, n *entity_nation) {
	fprintf(fp, "NA\n")

	if n.name != "" {
		fprintf(fp, " na %s\n", n.name)
	}
	if n.citizen != "" {
		fprintf(fp, " ci %s\n", n.citizen)
	}
	if n.capital != 0 {
		box_print(fp, []byte(" cp "), n.capital)
	}
	if n.win != 0 {
		fprintf(fp, " wi %d\n", n.win)
	}
	if n.player_limit != 0 {
		fprintf(fp, " pl %d\n", n.player_limit)
	}
	if n.jump_start != 0 {
		fprintf(fp, " js %d\n", n.jump_start)
	}
	if n.neutral {
		fprintf(fp, " ne %d\n", TRUE) // mdhender: changed to TRUE
	}
	if n.proscribed_skills != nil {
		boxlist_print(fp, []byte(" ps "), n.proscribed_skills)
	}
	fprintf(fp, "\n")
}

func save_box(fp *os.File, n int) {
	if kind(n) == T_deleted {
		return
	}

	assert(valid_box(n))

	p := bx[n]

	if bx[n].skind != FALSE {
		fprintf(fp, "%d %s %s\n", n, kind_s[bx[n].kind], subkind_s[bx[n].skind])
	} else {
		fprintf(fp, "%d %s 0\n", n, kind_s[bx[n].kind])
	}

	if len(p.name) != 0 {
		fprintf(fp, "na %s\n", p.name)
	}

	item_list_print(fp, []byte("il\t"), p.items)
	trade_list_print(fp, []byte("tl\t"), p.trades)
	effect_list_print(fp, []byte("el\t"), p.effects)

	if pd := rp_disp(n); pd != nil {
		boxlist_print(fp, []byte("an "), pd.neutral)
		boxlist_print(fp, []byte("ad "), pd.defend)
		boxlist_print(fp, []byte("ah "), pd.hostile)
	}

	if vp := rp_loc_info(n); vp != nil {
		print_loc_info(fp, vp)
	}
	if vp := rp_char(n); vp != nil {
		print_char(fp, vp)
	}
	if vp := rp_magic(n); vp != nil {
		print_magic(fp, vp)
	}
	if vp := rp_loc(n); vp != nil {
		print_loc(fp, vp)
	}
	if vp := rp_subloc(n); vp != nil {
		print_subloc(fp, vp)
	}
	if vp := rp_item(n); vp != nil {
		print_item(fp, vp)
	}
	if vp := rp_item_magic(n); vp != nil {
		print_item_magic(fp, vp)
	}
	if vp := rp_item_artifact(n); vp != nil {
		print_artifact(fp, vp)
	}
	if vp := rp_player(n); vp != nil {
		print_player(fp, vp)
	}
	if vp := rp_skill(n); vp != nil {
		print_skill(fp, vp)
	}
	if vp := rp_gate(n); vp != nil {
		print_gate(fp, vp)
	}
	if vp := rp_misc(n); vp != nil {
		print_misc(fp, vp)
	}
	if vp := rp_command(n); vp != nil {
		print_command(fp, vp)
	}
	if vp := rp_ship(n); vp != nil {
		print_ship(fp, vp)
	}
	/*
	 *  Thu Jan 21 07:13:13 1999 -- Scott Turner
	 *
	 *  Don't actually need to save this info, since it should not change.
	 *
	 *  Fri Mar 30 17:35:08 2001 -- Scott Turner
	 *
	 *  It does now!
	 *
	 */
	if vp := rp_nation(n); vp != nil {
		print_nation(fp, vp)
	}

	fprintf(fp, "\n")

	bx[n].temp = 1 /* mark for write_leftovers() */
}

func open_write_fp(fnam string) (*os.File, error) {
	fnam = filepath.Join(libdir, fnam)
	fp, err := os.Create(fnam)
	if err != nil {
		fprintf(os.Stderr, "open_write_fp: can't open %s: %v\n", fnam, err)
		return nil, err
	}
	return fp, nil
}

func write_kind(box_kind int, fnam string) {
	fp, err := open_write_fp(fnam)
	if err != nil {
		log.Printf("write_kind: %+v\n", err)
		return
	}
	for _, i := range loop_kind(box_kind) {
		save_box(fp, i)
	}

	_ = fp.Close()
}

func write_player(pl int) {
	fp, err := open_write_fp(sout("fact/%d", pl))
	if err != nil {
		log.Printf("write_player: %+v\n", err)
		return
	}

	save_box(fp, pl)

	for _, who := range loop_units(pl) {
		assert(kind(who) == T_char || kind(who) == T_deleted)
		save_box(fp, who)
	}

	_ = fp.Close()
}

func write_chars() {
	for _, pl := range loop_player() {
		write_player(pl)
	}
}

func write_leftovers() {
	fp, err := open_write_fp("misc")
	if err != nil {
		log.Printf("write_leftovers: %+v\n", err)
		return
	}

	for i := 0; i < MAX_BOXES; i++ {
		if bx[i] != nil && kind(i) != T_nation && bx[i].temp == 0 {
			if kind(i) != T_deleted {
				save_box(fp, i)
			}
		}
	}

	_ = fp.Close()
}

func read_boxes(fnam string) {
	fnam = filepath.Join(libdir, fnam)
	if !readfile(fnam) {
		return
	}

	advance()

	for line != nil {
		if len(line) == 0 || line[0] == '#' {
			advance()
			continue /* skip blank and comment lines */
		}
		box_num := atoi_b(line)
		if box_num > 0 {
			load_box(box_num)
		} else {
			fprintf(os.Stderr, "read_boxes(%s): unexpected line %s\n", fnam, line)
			advance()
		}
	}
}

func read_chars() error {
	dirFact := filepath.Join(libdir, "fact")
	files, err := os.ReadDir(dirFact)
	if err != nil {
		log.Printf("read_chars: can't open %q: %v\n", dirFact, err)
		return err
	}

	for _, f := range files {
		if isdigit(f.Name()[0]) && !strings.HasSuffix(f.Name(), "~") {
			read_boxes(filepath.Join("fact", f.Name()))
		}
	}
	return nil
}

func fast_scan() bool {
	path := filepath.Join(libdir, "master")
	if !readfile(path) {
		return false
	}

	for s := readlin(); s != nil; s = readlin() {
		if len(s) == 0 {
			continue
		}
		num := atoi_b(s)

		p := s
		for len(p) != 0 && isdigit(p[0]) {
			p = p[1:]
		}
		for len(p) != 0 && iswhite(p[0]) {
			p = p[1:]
		}

		var q []byte
		for n, i := 0, bytes.IndexByte(q, '.'); i != -1; n, i = n+1, bytes.IndexByte(q, '.') {
			if n == 0 {
				p = p[:i]
			}
			q = q[i:]
		}

		kind, sk := atoi_b(p), atoi_b(q)

		alloc_box(num, kind, sk)
	}

	return true
}

func scan_boxes(fnam string) {
	fnam = filepath.Join(libdir, fnam)
	if !readfile(fnam) {
		return
	}

	for s := readlin(); s != nil; s = readlin() {
		if len(s) == 0 || !isdigit(s[0]) {
			continue
		}

		/*
		   if (*s == '#')
		       continue;
		*/

		/*
		 *  Parse something of the form: box-number kind subkind
		 *  example:  10 item artifact
		 */

		box_num := atoi_b(s)

		/* skip over space to kind */
		for len(s) != 0 && s[0] != ' ' {
			s = s[1:]
		}
		if len(s) != 0 && s[0] == ' ' {
			s = s[1:]
		}

		/* skip to subkind */
		var t []byte
		var i int
		for i = 0; i < len(s) && s[i] != ' '; i++ {
			//
		}
		if i < len(s) && s[i] == ' ' {
			s, t = s[:i], s[i+i:]
		}

		kind := lookup_sb(kind_s, s)

		if kind < 0 {
			fprintf(os.Stderr, "read_boxes(%d): bad kind: %s\n", box_num, s)
			kind = 0
		}

		// todo: this can pass nil to lookup()
		var sk int
		if len(t) != 0 && t[0] == '0' {
			sk = 0
		} else {
			sk = lookup_sb(subkind_s, t)
		}
		if sk < 0 {
			fprintf(os.Stderr, "read_boxes(%d): bad subkind: %s\n", box_num, t)
			sk = 0
		}

		alloc_box(box_num, kind, sk)

		for s = readlin(); s != nil; s = readlin() {
			if len(s) == 0 || s[0] == '#' { /* skip blank lines and comments */
				continue
			}

			for len(s) != 0 && iswhite(s[0]) {
				s = s[1:]
			}

			if len(s) == 0 { /* blank line, end of entry */
				break
			}
		}
	}
}

func scan_chars() error {
	dirFact := filepath.Join(libdir, "fact")
	files, err := os.ReadDir(dirFact)
	if err != nil {
		log.Printf("scan_chars: can't open %q: %v\n", dirFact, err)
		return err
	}

	for _, f := range files {
		if isdigit(f.Name()[0]) && !strings.HasSuffix(f.Name(), "~") {
			scan_boxes(filepath.Join("fact", f.Name()))
		}
	}

	return nil
}

/*
 *  Scan through all of the entity data files, calling alloc_box
 *  for each entity once its number and kind are known.
 *
 *  We do this so it is possible to perform type and sanity checking when
 *  the contents of the boxes are read in the second pass (read_boxes).
 */
func scan_all_boxes() {

	stage("fast_scan failed, scan_all_boxes()")

	scan_boxes("loc")
	scan_boxes("item")
	scan_boxes("skill")
	scan_boxes("gate")
	scan_boxes("road")
	scan_boxes("ship")
	scan_boxes("unform")
	scan_boxes("misc")
	scan_boxes("nation")

	if err := scan_chars(); err != nil {
		log.Printf("scan_all_boxes: %+v\n", err)
	}
}

func read_all_boxes() {

	read_boxes("loc")
	read_boxes("item")
	read_boxes("skill")
	read_boxes("gate")
	read_boxes("road")
	read_boxes("ship")
	read_boxes("unform")
	read_boxes("misc")
	read_boxes("nation")

	if err := read_chars(); err != nil {
		log.Printf("read_all_boxes: %+v\n", err)
	}
}

func write_all_boxes() {
	for i := 0; i < MAX_BOXES; i++ {
		if bx[i] != nil {
			bx[i].temp = 0
		}
	}

	dirFact := filepath.Join(libdir, "fact")
	if err := rmdir(dirFact); err != nil {
		log.Fatalf("write_all_boxes: %+v\n", err)
	} else if err := mkdir(dirFact); err != nil {
		log.Fatalf("write_all_boxes: %+v\n", err)
	}

	write_kind(T_loc, "loc")
	write_kind(T_item, "item")
	write_kind(T_skill, "skill")
	write_kind(T_gate, "gate")
	write_kind(T_road, "road")
	write_kind(T_ship, "ship")
	write_kind(T_unform, "unform")
	write_kind(T_nation, "nation")

	write_chars()
	write_leftovers()
}

func write_master() error {
	fnam := filepath.Join(libdir, "master")
	fp, err := os.Create(fnam)
	if err != nil {
		fprintf(os.Stderr, "can't write %s: %v\n", fnam, err)
		return err
	}

	for i := 0; i < MAX_BOXES; i++ {
		if bx[i] != nil {
			bx[i].temp = 0
		}
	}

	for i := 0; i < MAX_BOXES; i++ {
		if kind(i) != T_deleted {
			s := name(i)

			switch kind(i) {
			case 0:
				break

			case T_loc, T_item, T_skill, T_gate, T_ship, T_unform:
				fprintf(fp, "%d\t%d.%d\t%s\t\t%s\n", i, bx[i].kind, bx[i].skind, kind_s[bx[i].kind], s)
				bx[i].temp = 1
				break

			case T_player:
				fprintf(fp, "%d\t%d.%d\tfact/%d\t%s\n", i, bx[i].kind, bx[i].skind, i, s)
				bx[i].temp = 1
				break

			case T_char:
				fprintf(fp, "%d\t%d.%d\tfact/%d\t%s\n", i, bx[i].kind, bx[i].skind, player(i), s)
				bx[i].temp = 1
				break
			}
		}
	}

	for i := 0; i < MAX_BOXES; i++ {
		if kind(i) != T_deleted && bx[i].temp == 0 {
			s := name(i)
			fprintf(fp, "%d\t%d.%d\tmisc\t%s\n", i, bx[i].kind, bx[i].skind, s)
		}
	}

	_ = fp.Close()

	return nil
}

//#if 0
//static void
//scan_nation(struct entity_nation *n)
//{
//  char *t;
//  int c;
//
//  advance();
//  for line != nil && len(line) != 0 && iswhite(line[0]) {
//    if (*line == '#') continue;
//
//    line = line[1:]
//    c := linehash(line)
//    t := t_string(line)
//
//    switch (c) {
//    case `na`:    {
//      n.name = string(t)
//      break;
//    }
//    case `ci`:    n.citizen = string(t) break;
//    case `cp`:    n.capital = box_scan(t); break;
//    case `wi`:    n.win = atoi(t); break;
//    case `pl`:  n.player_limit = atoi(t); break;
//    case `js`:  n.jump_start = atoi(t); break;
//    case `ps`:
//      boxlist_scan(t, -1, &(n.proscribed_skills));
//      break;
//
//    case 0:
//    default:
//      fprintf(os.Stderr,  "scan_nation:  bad line: %s\n", line);
//    }
//    advance();
//  }
//}
//
///*
// *  Tue Apr  8 11:55:15 1997 -- Scott Turner
// *
// *  The nations file.
// *
// *  Fri Jan 15 14:20:53 1999 -- Scott Turner
// *
// */
//static void
//scan_nations()
//{
//  char *s;
//  char *fname;
//  int i;
//
//  fname = sout("%s/nations", libdir);
//  if (!readfile(fname)) {
//    fprintf(os.Stderr,  "scan_nations: can't read %s: ", fname);
//    perror("");
//    return;
//  }
//
//  s = readlin();
//
//  if (s == nil || (num_nations = atoi(s)) == 0) {
//    fprintf(os.Stderr, "No nations number in nations file?\n");
//    return;
//  }
//
//  if (num_nations >= MAX_NATIONS) {
//    fprintf(os.Stderr, "Too many nations (%d) in nations file.\n",num_nations);
//    return;
//  }
//
//  nations[0].name = "";
//  nations[0].citizen = "";
//
//  for(i=1;i<=num_nations;i++) {
//    scan_nation(&nations[i]);
//  }
//}
//
///*
// *  Tue Apr  8 11:55:15 1997 -- Scott Turner
// *
// *  The nations file.
// *
// */
//static void
//save_nations()
//{
//  FILE *fp;
//  int i;
//
//  fp = open_write_fp("nation");
//  if (fp == nil) return;
//
//  fprintf(fp,"%d\n",num_nations);
//  for(i=1;i<=num_nations;i++)
//    save_nation(fp, &nations[i]);
//
//  _ = fp.Close()
//
//}
//#endif

func load_system() error {
	fname := filepath.Join(libdir, "system.json")
	data, err := os.ReadFile(fname)
	if err != nil {
		log.Printf("load_system: system: %v\n", err)
		return err
	}
	var js xlat.SYSTEM
	if err := json.Unmarshal(data, &js); err != nil {
		log.Printf("load_system: system: %v\n", err)
		return err
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

	return nil
}

func save_system() error {
	log.Printf("save_system: please update for boolean values\n")

	fname := filepath.Join(libdir, "system")
	fp, err := os.Create(fname)
	if err != nil {
		fprintf(os.Stderr, "load_system: can't write %s: %v\n", fname, err)
		return err
	}

	olytime_print(fp, []byte("sysclock: "), &sysclock)
	fprintf(fp, "game_num=%d\n", game_number)
	fprintf(fp, "indep_player=%d\n", indep_player)
	fprintf(fp, "gm_player=%d\n", gm_player)
	fprintf(fp, "skill_player=%d\n", skill_player)
	fprintf(fp, "from_host=%s\n", from_host)
	fprintf(fp, "reply_host=%s\n", reply_host)
	fprintf(fp, "seed0=%d\n", seed[0])
	fprintf(fp, "seed1=%d\n", seed[1])
	fprintf(fp, "seed2=%d\n", seed[2])
	fprintf(fp, "post=%d\n", post_has_been_run)
	fprintf(fp, "init=%d\n", seed_has_been_run)
	fprintf(fp, "ds=%d\n", dist_sea_compute)
	fprintf(fp, "nc=%d\n", near_city_init)
	fprintf(fp, "mi=%d\n", cookie_init)
	fprintf(fp, "fr=%d\n", faery_region)
	fprintf(fp, "fp=%d\n", faery_player)
	fprintf(fp, "hr=%d\n", hades_region)
	fprintf(fp, "hp=%d\n", hades_pit)
	fprintf(fp, "hl=%d\n", hades_player)
	fprintf(fp, "np=%d\n", npc_pl)
	fprintf(fp, "ms=%d\n", monster_subloc_init)
	fprintf(fp, "pi=%d\n", population_init)
	fprintf(fp, "cr=%d\n", cloud_region)
	fprintf(fp, "cp=%d\n", combat_pl)
	fprintf(fp, "xsize=%d\n", xsize)
	fprintf(fp, "ysize=%d\n", ysize)
	fprintf(fp, "turn_limit=%d\n", options.turn_limit)
	fprintf(fp, "autodrop=%d\n", options.auto_drop)
	fprintf(fp, "free=%d\n", options.free)
	fprintf(fp, "turn_charge=%s\n", options.turn_charge)
	fprintf(fp, "mp_antipathy=%d\n", options.mp_antipathy)
	fprintf(fp, "survive_np=%d\n", options.survive_np)
	fprintf(fp, "death_nps=%d\n", options.death_nps)
	fprintf(fp, "guild_teaching=%d\n", options.guild_teaching)
	fprintf(fp, "accounting_dir=%s\n", options.accounting_dir)
	fprintf(fp, "accounting_prog=%s\n", options.accounting_prog)
	fprintf(fp, "html_path=%s\n", options.html_path)
	fprintf(fp, "html_passwords=%s\n", options.html_passwords)
	fprintf(fp, "times_pay=%d\n", options.times_pay)
	fprintf(fp, "open_ended=%d\n", options.open_ended)
	fprintf(fp, "cpp=%s\n", options.cpp)
	fprintf(fp, "full_markets=%d\n", options.full_markets)
	fprintf(fp, "output_tags=%d\n", options.output_tags)
	fprintf(fp, "num_books=%d\n", options.num_books)
	fprintf(fp, "market_age=%d\n", options.market_age)
	fprintf(fp, "min_piety=%d\n", options.min_piety)
	fprintf(fp, "piety_limit=%d\n", options.piety_limit)
	fprintf(fp, "head_priest_piety_limit=%d\n", options.head_priest_piety_limit)
	fprintf(fp, "top_piety=%d\n", options.top_piety)
	fprintf(fp, "middle_piety=%d\n", options.middle_piety)
	fprintf(fp, "bottom_piety=%d\n", options.bottom_piety)
	fprintf(fp, "claim_give=%d\n", options.claim_give)
	fprintf(fp, "check_balance=%d\n", options.check_balance)
	fprintf(fp, "free_np_limit=%d\n", options.free_np_limit)

	_ = fp.Close()

	return nil
}

func delete_deadchars() {
	for _, i := range loop_kind(T_deadchar) {
		/*
		 *  Loc should have been zeroed already in kill_char
		 */

		if loc(i) != 0 {
			set_where(i, 0)
		}
		delete_box(i)
	}
}

func load_db() error {

	stage("load_db()")

	/*
	 *  Assertions to verify the sanity of the linehash macro
	 *  Switch the byte ordering if this fails
	 */

	assert(linehash(nil) == ``)
	assert(linehash([]byte("")) == ``)
	assert(linehash([]byte("ab f")) == `ab`)
	assert(linehash([]byte("na")) == `na`)
	assert(linehash([]byte("ab ")) == `ab`)

	if err := load_system(); err != nil {
		log.Fatalf("load_db: %+v\n", err)
	}

	if !fast_scan() { /* pass 1: call alloc_box for each entity */
		scan_all_boxes()
	}

	read_all_boxes() /* pass 2: read the entity attributes */
	/*
	 *  At this point we should be able to set the MAX_MM
	 *  to be that of a dragon.
	 *
	 *  Mon Jun 14 07:27:07 1999 -- Scott Turner
	 *
	 *  Using nazgul, since dragon is an outlier.
	 *
	 */
	//MAX_MM = MM(item_nazgul);
	if !(MAX_MM < MM(item_nazgul)) {
		panic("assert(MM(item_nazgul) <= MAX_MM)")
	}

	if err := load_orders(); err != nil {
		log.Fatalf("load_db: %+v\n")
	}

	//#if 0
	//    scan_nations();         /* Here so the boxes are valid... */
	//#endif

	check_db()            /* check database integrity */
	determine_map_edges() /* initialization for map routines */

	if post_has_been_run == FALSE {
		stage("INIT: post_production()")
		post_production()
	}

	if seed_has_been_run == FALSE {
		stage("INIT: seed_initial_locations()")
		seed_initial_locations()
		seed_orcs()
	}

	if faery_region == 0 {
		create_faery()
	}

	if hades_region == 0 {
		create_hades()
	}

	if cloud_region == 0 {
		create_cloudlands()
	}

	if dist_sea_compute == FALSE {
		compute_dist()
		dist_sea_compute = TRUE
	}

	if near_city_init == FALSE {
		seed_city_near_lists()
		near_city_init = TRUE
	}

	if cookie_init == FALSE {
		seed_cookies()
		cookie_init = TRUE
	}

	if !monster_subloc_init {
		seed_monster_sublocs(true)
		monster_subloc_init = true
	}

	if !population_init {
		seed_population()
		population_init = true
	}

	if combat_pl == 0 || bx[combat_pl] == nil {
		combat_pl = 210

		alloc_box(combat_pl, T_player, sub_pl_npc)
		set_name(combat_pl, "Combat log")
		p_player(combat_pl).password = DEFAULT_PASSWORD
		fprintf(os.Stderr, "\tcreated combat player %d\n", combat_pl)
	}

	init_ocean_chars()
	delete_deadchars()

	return nil
}

func cleanup_posts() {
	for _, i := range loop_post() {
		set_where(i, 0)
		delete_box(i)
	}
}

func save_logdir() error {
	system(sout("rm -rf %s/save/%d", libdir, sysclock.turn))
	if err := mkdir(filepath.Join(libdir, "save")); err != nil {
		return err
	}

	s := filepath.Join(libdir, "log")
	t := filepath.Join(libdir, "save", fmt.Sprintf("%d", sysclock.turn))

	if err := rename(s, t); err != nil {
		fprintf(os.Stderr, "couldn't rename %s to %s: %v\n", s, t, err)
		return err
	}

	s = filepath.Join(libdir, "players.html", libdir)
	t = filepath.Join(libdir, "save", fmt.Sprintf("%d", sysclock.turn), "players.html")

	if err := rename(s, t); err != nil {
		fprintf(os.Stderr, "couldn't rename %s to %s: %v\n", s, t, err)
		return err
	}

	return mkdir(filepath.Join(libdir, "log"))
}

func save_db() error {
	stage("save_db()")

	cleanup_posts()
	if err := save_system(); err != nil {
		log.Printf("save_db: %v\n", err)
		return err
	}
	write_all_boxes()
	if err := write_master(); err != nil {
		log.Printf("save_db: %v\n", err)
		return err
	}
	save_orders()
	rename_act_join_files()

	return nil
}
