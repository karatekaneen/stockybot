package stockybot

import (
	"testing"

	"github.com/matryer/is"
)

func TestLastN(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		arr     []int
		wantErr bool
		want    []int
	}{
		{
			name: "simple",
			size: 3,
			arr:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			want: []int{8, 9, 0},
		},
		{
			name: "simple",
			size: 5,
			arr:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			want: []int{6, 7, 8, 9, 0},
		},
		{
			name: "simple",
			size: 0,
			arr:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			want: []int{},
		},
		{
			name:    "out of bounds",
			size:    11,
			wantErr: true,
			arr:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			got, err := LastN(tt.arr, tt.size)
			if tt.wantErr {
				is.True(err != nil)
				return
			}

			is.NoErr(err)
			is.Equal(got, tt.want)
		})
	}
}
