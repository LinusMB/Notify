package parsing

import (
	"fmt"
	"strconv"
)

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
	s := newState(input)

	parseFloat := func(token string) (float64, error) {
		if token == "" {
			return 0, nil
		}
		return strconv.ParseFloat(token, 64)
	}

	var width string
	width, s, err = lexUntil(s, 'x')
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse width of input %s: %w",
			input,
			err,
		)
	}
	dim.Width, err = parseFloat(width)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse width of input %s: %w",
			input,
			err,
		)
	}

	var height string
	height, s, err = lexUntil(s, '+')
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse height of input %s: %w",
			input,
			err,
		)
	}
	dim.Height, err = parseFloat(height)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse height of input %s: %w",
			input,
			err,
		)
	}

	var xPos string
	xPos, s, err = lexUntil(s, '+')
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse x position of input %s: %w",
			input,
			err,
		)
	}
	dim.X, err = parseFloat(xPos)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse x position of input %s: %w",
			input,
			err,
		)
	}

	var yPos string
	yPos = s.remaining()
	dim.Y, err = parseFloat(yPos)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse y position of input %s: %w",
			input,
			err,
		)
	}

	return &dim, nil
}
