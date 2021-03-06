package pretty

import (
	"testing"
)

type difftest struct {
	a   interface{}
	b   interface{}
	exp []string
}

type S struct {
	A int
	S *S
	I interface{}
	C []int
	M map[string]int
}

var diffs = []difftest{
	{a: nil, b: nil},
	{a: S{A: 1}, b: S{A: 1}},

	{0, "", []string{`int != string`}},
	{0, 1, []string{`0 != 1`}},
	{S{}, new(S), []string{`pretty.S != *pretty.S`}},
	{"a", "b", []string{`"a" != "b"`}},
	{S{}, S{A: 1}, []string{`A: 0 != 1`}},
	{new(S), &S{A: 1}, []string{`A: 0 != 1`}},
	{S{S: new(S)}, S{S: &S{A: 1}}, []string{`S.A: 0 != 1`}},
	{S{}, S{I: 0}, []string{`I: nil != 0`}},
	{S{I: 1}, S{I: "x"}, []string{`I: int != string`}},
	{S{}, S{C: []int{1}}, []string{`C: []int[0] != []int[1]`}},
	{S{C: []int{}}, S{C: []int{1}}, []string{`C: []int[0] != []int[1]`}},
	{S{C: []int{1, 2, 3}}, S{C: []int{1, 2, 4}}, []string{`C[2]: 3 != 4`}},
	{S{}, S{A: 1, S: new(S)}, []string{`A: 0 != 1`, `S: nil != &{0 <nil> <nil> [] map[]}`}},
	{S{}, S{M:map[string]int{"a":2}}, []string{`M["a"]: (missing) != "<int Value>"`}},
	{S{M:map[string]int{"a":2}}, S{}, []string{`M["a"]: "<int Value>" != (missing)`}},
	{S{M:map[string]int{"a":1}}, S{M:map[string]int{"a":2}}, []string{`M["a"]: 1 != 2`}},
}

func TestDiff(t *testing.T) {
	for _, tt := range diffs {
		got := Diff(tt.a, tt.b)
		eq := len(got) == len(tt.exp)
		if eq {
			for i := range got {
				eq = eq && got[i] == tt.exp[i]
			}
		}
		if !eq {
			t.Errorf("diffing % #v", tt.a)
			t.Errorf("with    % #v", tt.b)
			diffdiff(t, got, tt.exp)
			continue
		}
	}
}

func diffdiff(t *testing.T, got, exp []string) {
	minus(t, "unexpected:", got, exp)
	minus(t, "missing:", exp, got)
}

func minus(t *testing.T, s string, a, b []string) {
	var i, j int
	for i = 0; i < len(a); i++ {
		for j = 0; j < len(b); j++ {
			if a[i] == b[j] {
				break
			}
		}
		if j == len(b) {
			t.Error(s, a[i])
		}
	}
}
