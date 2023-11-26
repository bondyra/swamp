package common

import (
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

func lst(args ...[]string) [][]string {
	return args
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
