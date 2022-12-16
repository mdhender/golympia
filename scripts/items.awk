BEGIN {
    id = 0;
    printf("[");
}

function qsafe(s) {
    gsub("\"", "", s)
    return s
}

# item section
/^[0-9]+ item / {
    if (section != "") {
        printf("\n  }");
    }
    if (id != 0) {
        printf("\n},");
    }
    id = $1;
    kind = $2;
    subkind = $3;
    section = "";
    printf("\n{ \"id\": %d", id);
    printf("\n, \"kind\": \"%s\"", kind);
    if (subkind != "0") {
        printf("\n, \"subkind\": \"%s\"", subkind);
    }
    next;
}

# name; string
/^na / {
    if (section != "") {
        printf("\n  }");
    }
    section = "";
    msg = $2;
    for (i = 3; i <= NF; i++) {
        msg = msg " " $i
    }
    printf("\n, \"name\": \"%s\"", qsafe(msg));
    next;
}

# character
/^CH/ {
    if (section != "") {
        printf("\n  }");
    }
    section = "character";
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

# character magic
/^CM/ {
    if (section != "") {
        printf("\n  }");
    }
    section = "character-magic";
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

# commands
/^CO/ {
    if (section != "") {
        printf("\n  }");
    }
    section = "commands";
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

# item magic
/^IM/ {
    if (section != "") {
        printf("\n  }");
    }
    section = "item-magic";
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

# item
/^IT/ {
    if (section != "") {
        printf("\n  }");
    }
    section = "item";
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

# misc
/^MI/ {
    if (section != "") {
        printf("\n  }");
    }
    section = "misc";
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

# item
/^PL/ {
    if (section != "") {
        printf("\n  }");
    }
    section = "player";
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

/^ ab / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "attack-bonus", $2);
    next;
}

/^ an / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if ($2 != "0") {
        printf("\n,   \"%s\": %s", "animal", "true");
    }
    next;
}

# ar can be auraculum or arguments
/^ ar / {
    if (NF == 2) {
        printf("\n,   \"%s\": %s", "auraculum", $2);
    } else {
        msg = $2;
        for (i = 3; i <= NF; i++) {
            msg = msg " " $i
        }
        printf("\n,   \"%s\": \"%s\"", "args", qsafe(msg));
    }
    next;
}

/^ at / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "attack", $2);
    next;
}

/^ au / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "aura", $2);
    next;
}

# ba may be barrier or aura bonus
/^ ba / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "item-magic") {
        printf("\n,   \"%s\": %d", "aura-bonus", $2);
    } else if (section == "locations") {
        printf("\n,   \"%s\": %d", "barrier", $2);
    } else {
        printf("\n,   \"%s\": %d", "ba", $2);
    }
    next;
}

/^ bh / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "behind", $2);
    next;
}

# bp can be break point or base price
/^ bp / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "character") {
        printf("\n,   \"%s\": %d", "break-point", $2);
    } else if (section == "item" || section == "items") {
        printf("\n,   \"%s\": %d", "base-price", $2);
    } else {
        printf("\n,   \"%s\": %d", "bp", $2);
    }
    next;
}

/^ ca / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if ($2 != "0") {
        printf("\n,   \"%s\": %s", "capturable", "true");
    }
    next;
}

/^ cc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "cloak-creator", $2);
    next;
}

/^ cl / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "curse-loyalty", $2);
    next;
}

/^ cr / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "cloak-region", $2);
    next;
}

# ct can be contact or creator
/^ ct .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", "contacts", list);
    next;
}

/^ ct / {
    if (NF == 2) {
        if (section == "character") {
            printf("\n,   \"%s\": [%d]", "contacts", $2);
        } else if (section == "item-magic") {
            printf("\n,   \"%s\": %d", "creator", $2);
        } else {
            printf("\n,   \"%s\": %d", "ct", $2);
        }
    } else {
        list = $2;
        for (i = 3; i < NF; i++) {
            list = list ", " $i
        }
        printf("\n,   \"%s\": [%s]", "contacts", list);
    }
    next;
}

# db can be defense bonus or db path
/^ db / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "player") {
        printf("\n,   \"%s\": \"%s\"", "db-path", $2);
    } else {
        printf("\n,   \"%s\": %d", "defense-bonus", $2);
    }
    next;
}

# de can be defense or days executing
/^ de / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "commands") {
        printf("\n,   \"%s\": %d", "days-executing", $2);
    } else {
        printf("\n,   \"%s\": %d", "de", $2);
    }
    next;
}

/^ df / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "defense", $2);
    next;
}

/^ di / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "npc-dir", $2);
    next;
}

# ds can be display name or distance from sea
/^ ds / {
    if (NF == 2) {
        printf("\n,   \"%s\": %s", "distance-from-sea", $2);
    } else {
        msg = $2;
        for (i = 3; i <= NF; i++) {
            msg = msg " " $i
        }
        printf("\n,   \"%s\": \"%s\"", "display-name", qsafe(msg));
    }
    next;
}

/^ dt / {
    if (NF != 4) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": {\"day\": %d, \"turn\": %d, \"days-since-epoch\": %d}", "oly-time", $2, $3, $4);
    next;
}

/^ fc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "fly-cap", $2);
    next;
}

/^ he / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "health", $2);
    next;
}

# hs can be hide self or require holy symbol
/^ hs / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if ($2 != "0") {
        if (section == "character-magic") {
            printf("\n,   \"%s\": %s", "hide-self", "true");
        } else if (section == "skills") {
            printf("\n,   \"%s\": %s", "require-holy-symbol", "true");
        } else {
            printf("\n,   \"%s\": %s", "hs", "true");
        }
    }
    next;
}

/^ im / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if ($2 != "0") {
        printf("\n,   \"%s\": %s", "is-magician", "true");
    }
    next;
}

/^ kw / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "knows-weather", $2);
    next;
}

/^ li / {
    if (NF < 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    msg = $2;
    for (i = 3; i <= NF; i++) {
        msg = msg " " $i
    }
    printf("\n,   \"%s\": \"%s\"", "cmd", qsafe(msg));
    next;
}

/^ lc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "land-cap", $2);
    next;
}

/^ lk / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "loyalty-kind", $2);
    next;
}

/^ lo / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "lore", $2);
    next;
}

/^ lr / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "loyalty-rate", $2);
    next;
}

# ma can be magician or max aura
/^ ma / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "player") {
        printf("\n,   \"%s\": %d", "magic", $2);
    } else if (section == "character-magic") {
        printf("\n,   \"%s\": %d", "max-aura", $2);
    } else {
        printf("\n,   \"%s\": %d", "ma", $2);
    }
    next;
}

/^ mb / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "missile-bonus", $2);
    next;
}

/^ mi / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "missile", $2);
    next;
}

/^ ms .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", "may-study", list);
    next;
}

/^ ms / {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", "may-study", list);
    next;
}

/^ mu / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if ($2 != "0") {
        printf("\n,   \"%s\": %s", "is-man-item", "true");
    }
    next;
}

/^ oc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "orb-use-count", $2);
    next;
}

/^ ol / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "old-lord", $2);
    next;
}

# pc may be practice cost or project cast
/^ pc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "skill" || section == "skills") {
        printf("\n,   \"%s\": %d", "practice-cost", $2);
    } else if (section == "item-magic") {
        printf("\n,   \"%s\": %d", "project-cast", $2);
    } else {
        printf("\n,   \"%s\": %d", "pc", $2);
    }
    next;
}

/^ pl / {
    if (NF < 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    msg = $2;
    for (i = 3; i <= NF; i++) {
        msg = msg " " $i
    }
    printf("\n,   \"%s\": \"%s\"", "plural-name", qsafe(msg));
    next;
}

# po may be poll, npc program, or ports
/^ po / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "character" || section == "characters") {
        printf("\n,   \"%s\": %d", "npc-program", $2);
    } else if (section == "command" || section == "commands") {
        printf("\n,   \"%s\": %d", "poll", $2);
    } else if (section == "ship" || section == "ships") {
        printf("\n,   \"%s\": %d", "ports", $2);
    } else {
        printf("\n,   \"%s\": %d", "po", $2);
    }
    next;
}

/^ pr / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if ($2 != "0") {
        printf("\n,   \"%s\": %s", "prominent", "true");
    }
    next;
}

/^ ra / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "rank", $2);
    next;
}

/^ rb / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "aura-reflect", $2);
    next;
}

/^ rc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "ride-cap", $2);
    next;
}

/^ rd / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "unknown-rd", $2);
    next;
}

/^ si / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "sick", $2);
    next;
}

# sl can be skills
/^ sl\t.*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^ sl\t/ {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

/^ sn / {
    if (NF < 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    msg = $2;
    for (i = 3; i <= NF; i++) {
        msg = msg " " $i
    }
    printf("\n,   \"%s\": \"%s\"", "save-name", qsafe(msg));
    next;
}

/^ sr / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if ($2 != "0") {
        printf("\n,   \"%s\": %s", "swear-on-release", "true");
    }
    next;
}

# st can be strength or status
/^ st / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "commands") {
        printf("\n,   \"%s\": %d", "status", $2);
    } else if (section == "characters") {
        printf("\n,   \"%s\": %d", "strength", $2);
    } else {
        printf("\n,   \"%s\": %d", "st", $2);
    }
    next;
}

/^ ti / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "token-ni", $2);
    next;
}

# tl can be trades, time to learn, or to location
/^tl\t.*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^tl\t/ {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

/^ tn / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "token-num", $2);
    next;
}

/^ ue / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "use-experience", $2);
    next;
}

/^ uk / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "use-key", $2);
    next;
}

/^ un / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "units", $2);
    next;
}

# us can be use skill or uses
/^ us / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    if (section == "commands") {
        printf("\n,   \"%s\": %d", "use-skill", $2);
    } else if (section == "skills") {
        printf("\n,   \"%s\": %d", "uses", $2);
    } else {
        printf("\n,   \"%s\": %d", "us", $2);
    }
    next;
}

/^ vi .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", "visions", list);
    next;
}

/^ vi / {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", "visions", list);
    next;
}

/^ vp / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "unknown-vp", $2);
    next;
}

/^ wa / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "wait", $2);
    next;
}

/^ wt / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", "weight", $2);
    next;
}

# continuation line for a list
/^\t.*\\/ {
    list = $1;
    for (i = 2; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n          , %s", list);
    next;
}

/^\t/ {
    list = $1;
    for (i = 2; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n          , %s]", list);
    next;
}

NF != 0 {
    printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
    exit;
}

END {
    if (section != "") {
        printf("}", cb);
    }
    if (id != 0) {
        printf("\n}", cb);
    }
    printf("\n]\n");
}
