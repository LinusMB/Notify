package parsing

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}

func ParseColor(hex string) (color.RGBA, error) {
	var (
		clr color.RGBA
		err error
	)

	switch len(hex) {
	case 4:
		_, err = fmt.Sscanf(hex, "#%1x%1x%1x", &clr.R, &clr.G, &clr.B)
		clr.R |= clr.R << 4
		clr.G |= clr.G << 4
		clr.B |= clr.B << 4
		clr.A = 0xff
	case 7:
		_, err = fmt.Sscanf(hex, "#%02x%02x%02x", &clr.R, &clr.G, &clr.B)
		clr.A = 0xff
	case 9:
		_, err = fmt.Sscanf(
			hex,
			"#%02x%02x%02x%02x",
			&clr.R,
			&clr.G,
			&clr.B,
			&clr.A,
		)
	default:
		err = errors.New("unexpected format")
	}
	if err != nil {
		return clr, fmt.Errorf("could not parse color %s: %w", hex, err)
	}
	return clr, nil
}

type Dimension struct {
	Width  float64
	Height float64
	X      float64
	Y      float64
}

func ParseDimension(dimension string) (*Dimension, error) {
	var (
		dim Dimension
		err error
	)
	re, err := regexp.Compile("[+x]")
	failIf(err, "compile regex")
	ts := re.Split(dimension, -1)
	const (
		WIDTH = iota
		HEIGHT
		X
		Y
	)
	dim.Width, err = strconv.ParseFloat(ts[WIDTH], 64)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse width as float in %s: %w",
			dimension,
			err,
		)
	}
	dim.Height, err = strconv.ParseFloat(ts[HEIGHT], 64)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse height as float in %s: %w",
			dimension,
			err,
		)
	}
	dim.X, err = strconv.ParseFloat(ts[X], 64)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse x as float in %s: %w",
			dimension,
			err,
		)
	}
	dim.Y, err = strconv.ParseFloat(ts[Y], 64)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse y as float in %s: %w",
			dimension,
			err,
		)
	}
	return &dim, err
}

type Notification struct {
	Title string
	Body  string
}

func ParseNotification(input string) *Notification {
	re, err := regexp.Compile(`(?:\[(.*)\])?(.*)`)
	failIf(err, "compile regex")
	const (
		TITLE = iota + 1
		BODY
	)
	ms := re.FindSubmatch([]byte(input))

	title := strings.TrimSpace(string(ms[TITLE]))
	body := strings.TrimSpace(string(ms[BODY]))
	n := Notification{
		Title: title,
		Body:  body,
	}

	return &n
}
