package prmrtr

import (
	"testing"
)

func Benchmark_pathEntry_match(b *testing.B) {
	benchmarks := []struct {
		name string
		base string
		test string
	}{
		{
			name: "bench small",
			base: "/user/:id",
			test: "/user/12",
		},
		{
			name: "bench medium",
			base: "/user/:id/lists/:listID/item/:itemID",
			test: "/user/12/lists/100/item/22",
		},
		{
			name: "bench long",
			base: "/api/user/:id/lists/:listID/item/:itemID/with/:lots/more/:fields/:to/check",
			test: "/api/user/12/lists/100/item/22/with/alot/more/fields/for/check",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			pe := newPathEntry(bm.base, checkParams)
			test := newPathEntry(bm.test, skipParams)
			for i := 0; i < b.N; i++ {
				if !pe.match(test) {
					b.Error("should have matched", bm.test, bm.base)
				}
			}
		})
	}
}

func Benchmark_newPathEntry(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{
			name:  "bench small params",
			input: "/user/:id",
		},
		{
			name:  "bench medium params",
			input: "/user/:id/lists/:listID/item/:itemID",
		},
		{
			name:  "bench long params",
			input: "/api/user/:id/lists/:listID/item/:itemID/with/:lots/more/:fields/:to/check",
		},
		{
			name:  "bench small",
			input: "/user/12",
		},
		{
			name:  "bench medium",
			input: "/user/12/lists/100/item/22",
		},
		{
			name:  "bench long",
			input: "/api/user/12/lists/100/item/22/with/alot/more/fields/for/check",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				newPathEntry(bm.input, checkParams)
			}
		})
	}
}
