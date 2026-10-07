package arr

import "testing"

type item struct {
	qty   int
	price float64
}

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want int
	}{
		{"nil", nil, 0},
		{"values", []int{1, 2, 3}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.s); got != tt.want {
				t.Fatalf("Sum() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestSumBy(t *testing.T) {
	tests := []struct {
		name string
		s    []item
		want float64
	}{
		{"nil", nil, 0},
		{"values", []item{{2, 1.5}, {1, 4}}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SumBy(tt.s, func(i item) float64 { return float64(i.qty) * i.price })
			if got != tt.want {
				t.Fatalf("SumBy() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestAvg(t *testing.T) {
	tests := []struct {
		name string
		s    []uint8
		want float64
	}{
		{"nil", nil, 0},
		{"no overflow", []uint8{200, 250}, 225},
		{"fraction", []uint8{1, 2}, 1.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Avg(tt.s); got != tt.want {
				t.Fatalf("Avg() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestAvgBy(t *testing.T) {
	tests := []struct {
		name string
		s    []item
		want float64
	}{
		{"nil", nil, 0},
		{"values", []item{{1, 0}, {4, 0}}, 2.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AvgBy(tt.s, func(i item) int { return i.qty }); got != tt.want {
				t.Fatalf("AvgBy() = %v; want %v", got, tt.want)
			}
		})
	}
}
