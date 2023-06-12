BEGIN {
    id = 0;
    printf("[");
}

/^[0-9]+ loc / {
    if (section != "") {
        printf("\n  }");
    }
    if (id != 0) {
        printf("\n},");
    }
    id = $1;
    kind = $3;
    section = "";
    printf("\n{ \"id\": %d", id);
    printf("\n, \"kind\": \"%s\"", kind);
    next;
}

/^na / {
    if (section != "") {
        printf("\n  }");
    }
    section = "";
    msg = $2;
    for (i = 3; i <= NF; i++) {
        msg = msg " " $i
    }
    printf("\n, \"name\": \"%s\"", msg);
    next;
}

/^CH/ {
    if (section != "") {
        printf("\n  }");
    }
    section = $1;
    printf("\n, \"%s\": {", section);
    next;
}

/^LI/ {
    if (section != "") {
        printf("\n  }");
    }
    section = $1;
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

/^LO/ {
    if (section != "") {
        printf("\n  }");
    }
    section = $1;
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

/^MI/ {
    if (section != "") {
        printf("\n  }");
    }
    section = $1;
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

/^SL/ {
    if (section != "") {
        printf("\n  }");
    }
    section = $1;
    printf("\n, \"%s\": { \"pid\": %d", section, id);
    next;
}

/^ ba / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ bm / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ cl / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ cp / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ da / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ dg / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ de / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ ds / {
    msg = $2;
    for (i = 3; i <= NF; i++) {
        msg = msg " " $i
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ eg / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ er / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ hi / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ hl .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^ hl / {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

/^il\t.*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^il\t/ {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

/^ lc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ lf .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^ lf / {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

/^ lp / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ lt / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ lw / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ mc / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ md / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ nc .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^ nc / {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

/^ pd .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^ pd / {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

/^ ql / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ sd / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ sg / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ sh / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ sl / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ td / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^ te .*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s", $1, list);
    next;
}

/^ te / {
    list = $2;
    for (i = 3; i <= NF; i++) {
        list = list ", " $i
    }
    printf("\n,   \"%s\": [%s]", $1, list);
    next;
}

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

/^ wh / {
    if (NF != 2) {
        printf("\n\nerror: %d: unhandled line:\n%s\n\n", NR, $0);
        exit;
    }
    printf("\n,   \"%s\": %d", $1, $2);
    next;
}

/^\t.*\\/ {
    list = $2;
    for (i = 3; i < NF; i++) {
        list = list ", " $i
    }
    printf("\n          , %s", list);
    next;
}

/^\t/ {
    list = $2;
    for (i = 3; i <= NF; i++) {
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
    if (id != 0) {
        printf("\n}", cb);
    }
    printf("\n]\n");
}
