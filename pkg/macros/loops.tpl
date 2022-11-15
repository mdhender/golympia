{{- define "loop_all_here" -}}
{{- /*      loop_all_here(WHERE, I) */ -}}
{   var ll_l []int
    all_here({{ .WHERE }}, &ll_l)
    for ll_i := 0; ll_i < len(ll_l); ll_i++ {
        {{ .I }} = ll_l[ll_i]
{{- end -}}
{{- define "next_all_here" -}}
    }
    // ilist_reclaim(&ll_l)
}
{{- end -}}


{{- define "loop_castle" -}}
    {{- /*      loop_castle(i) */ -}}
    {{- template "loop_subkind" args "SK" "sub_castle" "I" .I -}}
{{- end -}}
{{- define "next_castle" -}}
    {{- template "next_subkind" -}}
{{- end -}}

{{- define "loop_city" -}}
    {{- /*      loop_city(i) */ -}}
    {{- template "loop_subkind" args "SK" "sub_city" "I" .I -}}
{{- end -}}
{{- define "next_city" -}}
    {{- template "next_subkind" -}}
{{- end -}}


{{- define "loop_kind" -}}
{{- /*      loop_kind(KIND, I) */ -}}
{
    ll_next := kind_first({{.KIND}})
    ll_i := ll_next
    for ll_next > 0 {
        ll_next = kind_next(ll_i)
        {{ .I }} = ll_i;
{{- end -}}
{{- define "next_kind" -}}
    }
}
{{- end -}}


{{- define "loop_loc" -}}
    {{- /*      loop_loc(i) */ -}}
    {{- template "loop_kind" args "KIND" "T_loc" "I" .I -}}
{{- end -}}
{{- define "next_loc" -}}
    {{- template "next_kind" -}}
{{- end -}}


{{- define "loop_nation" -}}
    {{- /*      loop_nation(i) */ -}}
    {{- template "loop_kind" args "KIND" "T_nation" "I" .I -}}
{{- end -}}
{{- define "next_nation" -}}
    {{- template "next_kind" -}}
{{- end -}}


{{- define "loop_subkind" -}}
{{- /*      loop_subkind(SK, I) */ -}}
{
    ll_next := sub_first({{.SK}})
    ll_i := ll_next
    for ll_next > 0 {
        ll_next = sub_next(ll_i)
        {{ .I }} = ll_i;
{{- end -}}
{{- define "next_subkind" -}}
    }
}
{{- end -}}

{{- define "loop_units" -}}
{{- /*      loop_units(PL, I) */ -}}
{
    var ll_l []int
    if rp_player({{ .PL }}) != nil {
        ll_l = append(ll_l, rp_player({{ .PL }}).units...)
    }
    for ll_i := 0; ll_i < len(ll_l); ll_i++ {
        {{ .I }} = ll_l[ll_i]
{{- end -}}
{{- define "next_unit" -}}
    }
}
{{- end -}}
