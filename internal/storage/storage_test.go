package storage

import (
	"reflect"
	"testing"
)

func Test_Batcher(t *testing.T) {
	type args struct {
		items []string
		size  int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "13 into 5",
			args: args{
				items: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"},
				size:  5,
			},
			want: [][]string{{"a", "b", "c", "d", "e"}, {"f", "g", "h", "i", "j"}, {"k", "l", "m"}},
		},
		{
			name: "15 into 5",
			args: args{
				items: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"},
				size:  5,
			},
			want: [][]string{{"a", "b", "c", "d", "e"}, {"f", "g", "h", "i", "j"}, {"k", "l", "m", "n", "o"}},
		},
		{
			name: "1 into 5",
			args: args{
				items: []string{"a"},
				size:  5,
			},
			want: [][]string{{"a"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Batcher(tt.args.items, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("batcher() = %v, want %v", got, tt.want)
			}
		})
	}
}
