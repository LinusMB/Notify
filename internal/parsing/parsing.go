package parsing

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}

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

type Dimension struct {
	Width  float64
	Height float64
	X      float64
	Y      float64
}

func ParseDimension(input string) (*Dimension, error) {
	var (
		dim Dimension
		err error
	)

	re, err := regexp.Compile(`(\d*)x(\d*)\+(-?\d*)\+(-?\d*)`)
	failIf(err, "compile regex")
	if !re.Match([]byte(input)) {
		return nil, fmt.Errorf("unexpected format: %s", input)
	}
	ms := re.FindSubmatch([]byte(input))[1:]
	var ts [4]string

	for i := range ms {
		t := string(ms[i])
		if t == "" {
			ts[i] = "0"
			continue
		}
		ts[i] = t
	}

	const (
		WIDTH = iota
		HEIGHT
		X
		Y
	)

	dim.Width, err = strconv.ParseFloat(ts[WIDTH], 64)
	failIf(err, "parse dimension width")

	dim.Height, err = strconv.ParseFloat(ts[HEIGHT], 64)
	failIf(err, "parse dimension height")

	dim.X, err = strconv.ParseFloat(ts[X], 64)
	failIf(err, "parse dimension x")

	dim.Y, err = strconv.ParseFloat(ts[Y], 64)
	failIf(err, "parse dimension y")

	return &dim, err
}

type Notification struct {
	Title string
	Body  string
}

func ParseNotification(input string) *Notification {
	var notif Notification

	// assume opening has already been consumed, parse until closing unless opening has been encountered before
	parseBalanced := func(input string, opening, closing rune) (string, string) {
		i := 0
		balance := 0
		for i < len(input) {
			c, w := utf8.DecodeRuneInString(input[i:])
			switch c {
			case closing:
				if balance == 0 {
					return strings.TrimSpace(input[:i]), input[i+w:]
				}
				balance--
			case opening:
				balance++
			}
			i += w
		}
		return strings.TrimSpace(input), ""
	}

	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "[") {
		input, _ = strings.CutPrefix(input, "[")
		notif.Title, input = parseBalanced(input, '[', ']')
	}
	notif.Body = strings.TrimSpace(input)
	return &notif
}
