{{- define "win" -}}
// has one or the other nation won?
func check_nation_win() bool {
    var k, i, flag, ruler, j int
    total := 0

    // this is possibly redundant.
    calculate_nation_nps()

    // minimum # of turns.
    if sysclock.turn < MIN_TURNS {
        return false
    }

    // loop_nation(i)
    {{ template "loop_nation" args "I" "i" }}
    {
        // ignore neutral nations
        if p_nation(k).neutral != 0 {
            continue
        }

        flag = 0
        // loop_city(i)
        {{ template "loop_city" args "I" "i" }}
        {
            ruler = player_controls_loc(i)
            if ruler != 0 && nation(ruler) != 0 && nation(ruler) != k {
                flag = 1
                break
            }
        }
        {{ template "next_city" }}
        if flag != 0 {
            continue
        }

        // loop_castle(i)
        {{ template "loop_castle" args "I" "i" }}
        {
            ruler = player_controls_loc(i)
            if ruler != 0 && nation(ruler) != 0 && nation(ruler) != k {
                flag = 1
                break
            }
        }
        {{ template "next_castle" }}
        if flag != 0 {
            continue
        }

        total = 0
        // loop_nation(j)
        {{ template "loop_nation" args "I" "j" }}
        {
            if p_nation(j).neutral != 0 {
                continue
            }
            if (k != j) {
                total += rp_nation(j).nps
            }
        }
        {{ template "next_nation" }}

        // you haven't met the win conditions this turn, so we should zero you out.
        if total * 2 >= rp_nation(k).nps {
            rp_nation(k).win = 0
            continue
        }

        // add another turn...
        rp_nation(k).win++

        // two nations cannot win simultaneously because of the last condition.
        if rp_nation(k).win == 2 {
            // we have a winner!
            // this really needs to go into the front of the Times, but how?
            return true
        }
    }
    {{ template "next_nation" }}

    return false
}
{{- end -}}
