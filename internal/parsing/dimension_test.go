package parsing

import "testing"

func TestParseDimension(t *testing.T) {
	tests := []struct {
		input string
		want  Dimension
	}{
		{"40x40+40+40", Dimension{Width: 40, Height: 40, X: 40, Y: 40}},
		{"0x0+0+0", Dimension{Width: 0, Height: 0, X: 0, Y: 0}},
		{"x+20+20", Dimension{Width: 0, Height: 0, X: 20, Y: 20}},
		{"x++", Dimension{Width: 0, Height: 0, X: 0, Y: 0}},
		{"42x42+-20+-20", Dimension{Width: 42, Height: 42, X: -20, Y: -20}},
	}
	for _, tt := range tests {
		testname := tt.input
		t.Run(testname, func(t *testing.T) {
			got, err := ParseDimension(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if *got != tt.want {
				t.Errorf("got %v, want %v", *got, tt.want)
			}
		})
	}
}
