package prmrtr

import (
	"path"
	"strings"
)

const (
	skipParams = iota
	checkParams
)

type pathEntry struct {
	base          string
	parameterized bool
	paramNames    map[string]int
	chunks        []chunk
}

// newPathEntry creates a new pathEntry from a base path string
//
// it splits the path on "/" and checks to see if the path contains any parameters
// this is done by checking each "chunk" of the path for a ":" prefix
// so /users/:id/ will chunk to ["user", ":id"] where ":id" is prefixed with ":"
// meaning the route is parameterized.
func newPathEntry(base string, params int) *pathEntry {
	base = path.Clean(base)
	p := &pathEntry{base: base, paramNames: make(map[string]int)}
	p.chunks = make([]chunk, 0, strings.Count(base, "/"))

	var start int
	for i := 0; i < len(base); i++ {
		if base[i] != '/' || i == start {
			continue
		}
		s := base[start+1 : i]
		start = i
		c := chunk{value: s}
		if params == checkParams && s[0] == ':' {
			p.parameterized = true
			p.paramNames[s] = len(p.chunks)
			c.parameter = true
		}
		p.chunks = append(p.chunks, c)
	}
	s := base[start+1:]
	c := chunk{value: s}
	if params == checkParams && s[0] == ':' {
		p.parameterized = true
		p.paramNames[s] = len(p.chunks)
		c.parameter = true
	}
	p.chunks = append(p.chunks, c)

	return p
}

func (p *pathEntry) match(entry *pathEntry) bool {
	if !p.parameterized {
		return p.base == entry.base
	}
	if len(p.chunks) != len(entry.chunks) {
		return false
	}
	for i, c := range p.chunks {
		if c.parameter {
			if entry.chunks[i].value == "" {
				return false
			}
			continue
		}
		if c.value != entry.chunks[i].value {
			return false
		}
	}
	return true
}

type chunk struct {
	value     string
	parameter bool
}
