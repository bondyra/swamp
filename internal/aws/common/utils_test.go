package common

import (
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

func lst(args ...[]string) [][]string {
	return args
}

func TestUnion(t *testing.T) {
	type args struct {
		inputs [][]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test nil",
			args: args{lst(nil, nil)},
			want: []string{},
		},
		{
			name: "test empty",
			args: args{inputs: lst([]string{}, []string{}, []string{})},
			want: []string{},
		},
		{
			name: "test some plus empty",
			args: args{inputs: lst([]string{"a", "b"}, []string{}, []string{})},
			want: []string{"a", "b"},
		},
		{
			name: "test some plus some",
			args: args{inputs: lst([]string{"a", "b"}, []string{"c"}, []string{"d"})},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "test some plus some overlapping",
			args: args{inputs: lst([]string{"a", "b"}, []string{"c", "d"}, []string{"d"})},
			want: []string{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Union(tt.args.inputs...)
			slices.Sort(got)
			slices.Sort(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuplicatedElements(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test nil",
			args: args{nil},
			want: []string{},
		},
		{
			name: "test returns nothing when unique",
			args: args{input: []string{"a", "b"}},
			want: []string{},
		},
		{
			name: "test duplicates",
			args: args{input: []string{"a", "b", "c", "a", "b"}},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DuplicatedElements(tt.args.input)
			slices.Sort(got)
			slices.Sort(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DuplicatedElements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	type args struct {
		inputs [][]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test nil",
			args: args{lst(nil, nil)},
			want: []string{},
		},
		{
			name: "test empty",
			args: args{inputs: lst([]string{}, []string{}, []string{})},
			want: []string{},
		},
		{
			name: "test no intersection",
			args: args{inputs: lst([]string{"a", "b"}, []string{"c"}, []string{"d"})},
			want: []string{},
		},
		{
			name: "test no intersection when common element is not in every input",
			args: args{inputs: lst([]string{"a", "b"}, []string{"c", "d"}, []string{"d"})},
			want: []string{},
		},
		{
			name: "test actual intersection",
			args: args{inputs: lst([]string{"a", "c", "b", "d"}, []string{"c", "d"}, []string{"d", "c", "b", "e"})},
			want: []string{"c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Intersect(tt.args.inputs...)
			slices.Sort(got)
			slices.Sort(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	type args struct {
		minuend     []string
		subtrahends [][]string
	}
	type s []string
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test nil",
			args: args{s{}, lst(nil, nil)},
			want: s{},
		},
		{
			name: "test empty",
			args: args{s{}, lst(s{}, s{})},
			want: s{},
		},
		{
			name: "test empty subtrahends",
			args: args{s{"a"}, lst(s{}, s{})},
			want: s{"a"},
		},
		{
			name: "test one subtrahend",
			args: args{s{"a", "b"}, lst(s{"b"})},
			want: s{"a"},
		},
		{
			name: "test one subtrahend and one empty subtrahend",
			args: args{s{"a", "b"}, lst(s{"b"}, s{})},
			want: s{"a"},
		},
		{
			name: "test multiple subtrahends",
			args: args{s{"a", "b", "c"}, lst(s{"b"}, s{"c"})},
			want: s{"a"},
		},
		{
			name: "test multiple subtrahends same value",
			args: args{s{"a", "b", "c"}, lst(s{"b"}, s{"b"})},
			want: s{"a", "c"},
		},
		{
			name: "test multiple subtrahends more values",
			args: args{s{"a", "b", "c", "d"}, lst(s{"b", "e"}, s{"c", "f"})},
			want: s{"a", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Difference(tt.args.minuend, tt.args.subtrahends...)
			slices.Sort(tt.want)
			slices.Sort(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type something struct {
		a string
		b int
	}
	type args struct {
		a []something
		f func(something) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test nil",
			args: args{a: nil, f: func(s something) string { return s.a }},
			want: nil,
		},
		{
			name: "test empty",
			args: args{a: []something{}, f: func(s something) string { return s.a }},
			want: make([]string, 0),
		},
		{
			name: "test actual map",
			args: args{a: []something{{a: "a", b: 1}, {a: "b", b: 2}}, f: func(s something) string { return s.a }},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.args.a, tt.args.f)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type something struct {
		a string
		b int
	}
	type args struct {
		a []something
		f func(something) bool
	}
	tests := []struct {
		name string
		args args
		want []something
	}{
		{
			name: "test",
			args: args{a: []something{{"match", 1}, {"no-match", 2}, {"match", 3}}, f: func(s something) bool { return s.a == "match" }},
			want: []something{{"match", 1}, {"match", 3}},
		},
		{
			name: "test nil",
			args: args{a: nil, f: func(s something) bool { return s.a == "match" }},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.a, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type testStruct struct {
		A string `json:"a"`
		B int    `json:"b,omitempty"`
	}
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *testStruct
		wantErr bool
	}{
		{
			name:    "test nil",
			args:    args{nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test empty",
			args:    args{[]byte("")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test invalid",
			args:    args{[]byte("{invalid json")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test required",
			args:    args{[]byte("{\"a\":\"a\"}")},
			want:    &testStruct{A: "a"},
			wantErr: false,
		},
		{
			name:    "test full",
			args:    args{[]byte("{\"a\":\"a\", \"b\":1}")},
			want:    &testStruct{A: "a", B: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal[testStruct](tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
