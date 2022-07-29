package prmrtr

import (
	"net/http"
	"strconv"
)

const (
	varCtxKey = iota + 1
)

func Vars(r *http.Request) Var {
	vars := r.Context().Value(varCtxKey)
	if vars == nil {
		return newVar()
	}
	return vars.(Var)
}

type Var struct {
	paramNames map[string]int
	chunks     []chunk
}

func newVar() Var {
	return Var{paramNames: make(map[string]int)}
}

func newVarFromEntries(pattern *pathEntry, actual *pathEntry) Var {
	return Var{
		paramNames: pattern.paramNames,
		chunks:     actual.chunks,
	}
}

func (v Var) String(name string) (string, bool) {
	if name[0] != ':' {
		name = ":" + name
	}
	idx, ok := v.paramNames[name]
	if !ok {
		return "", false
	}
	return v.chunks[idx].value, true
}

func (v Var) Int(name string) (int, bool) {
	s, ok := v.String(name)
	if !ok {
		return 0, false
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return i, true
}
