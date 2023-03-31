package parsing

import (
	"errors"
	"fmt"
	"image/color"
)

func ParseColor(input string) (color.RGBA, error) {
	var (
		clr color.RGBA
		err error
	)

	switch len(input) {
	case 4:
		_, err = fmt.Sscanf(input, "#%1x%1x%1x", &clr.R, &clr.G, &clr.B)
		clr.R |= clr.R << 4
		clr.G |= clr.G << 4
		clr.B |= clr.B << 4
		clr.A = 0xff
	case 7:
		_, err = fmt.Sscanf(input, "#%02x%02x%02x", &clr.R, &clr.G, &clr.B)
		clr.A = 0xff
	case 9:
		_, err = fmt.Sscanf(
			input,
			"#%02x%02x%02x%02x",
			&clr.R,
			&clr.G,
			&clr.B,
			&clr.A,
		)
	default:
		err = errors.New("unexpected input length")
	}
	if err != nil {
		return clr, fmt.Errorf("could not parse color %s: %w", input, err)
	}
	return clr, nil
}
