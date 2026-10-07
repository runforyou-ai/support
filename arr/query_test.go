package arr

import "testing"

func isEven(n int) bool { return n%2 == 0 }

func TestFirst(t *testing.T) {
	tests := []struct {
		name   string
		s      []int
		pred   func(int) bool
		want   int
		wantOK bool
	}{
		{"nil slice", nil, isEven, 0, false},
		{"nil pred", []int{3, 4}, nil, 3, true},
		{"nil pred empty", []int{}, nil, 0, false},
		{"match", []int{1, 2, 4}, isEven, 2, true},
		{"no match", []int{1, 3}, isEven, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := First(tt.s, tt.pred)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("First() = %v, %v; want %v, %v", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestLast(t *testing.T) {
	tests := []struct {
		name   string
		s      []int
		pred   func(int) bool
		want   int
		wantOK bool
	}{
		{"nil slice", nil, isEven, 0, false},
		{"nil pred", []int{3, 4, 5}, nil, 5, true},
		{"nil pred empty", nil, nil, 0, false},
		{"match", []int{1, 2, 4, 5}, isEven, 4, true},
		{"no match", []int{1, 3}, isEven, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Last(tt.s, tt.pred)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("Last() = %v, %v; want %v, %v", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestEvery(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want bool
	}{
		{"nil", nil, true},
		{"all", []int{2, 4}, true},
		{"some", []int{2, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Every(tt.s, isEven); got != tt.want {
				t.Fatalf("Every() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want int
	}{
		{"nil", nil, 0},
		{"none", []int{1, 3}, 0},
		{"some", []int{1, 2, 3, 4}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Count(tt.s, isEven); got != tt.want {
				t.Fatalf("Count() = %v; want %v", got, tt.want)
			}
		})
	}
}
