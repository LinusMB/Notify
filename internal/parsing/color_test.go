package parsing

import (
	"image/color"
	"testing"
)

func TestParseColor_ValidInput(t *testing.T) {
	tests := []struct {
		input string
		want  color.RGBA
	}{
		{"#ffffff", color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"#fff", color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
		{"#000000", color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}},
		{"#000", color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}},
		{"#0e0e0e", color.RGBA{R: 0x0e, G: 0x0e, B: 0x0e, A: 0xff}},
		{"#0f0", color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}},
	}
	for _, tt := range tests {
		testname := tt.input
		t.Run(testname, func(t *testing.T) {
			got, err := ParseColor(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseColor_InvalidInput(t *testing.T) {
	testInputs := []string{
		"",
		"#",
		"#ff",
		"#fffff",
		"#zzzzz",
		"1234",
		"1234567",
		"123456789",
		"#zzz",
	}
	for _, ti := range testInputs {
		testname := ti
		t.Run(testname, func(t *testing.T) {
			_, err := ParseColor(ti)
			if err == nil {
				t.Error("want error for invalid input")
			}
		})
	}
}
