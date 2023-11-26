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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Union(tt.args.inputs...); !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DuplicatedElements(tt.args.input); !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersect(tt.args.inputs...); !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.a, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal[string](tt.args.input)
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
