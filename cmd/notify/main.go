package main

import (
	"Notify/internal/fonts"
	"Notify/internal/parsing"
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/exp/constraints"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

type Configuration struct {
	fontFace      font.Face
	fontSize      float64
	winWidth      float64
	winHeight     float64
	winX          float64
	winY          float64
	autoDimension bool
	bgColor       color.Color
	fgColor       color.Color
	outputString  string
	duration      time.Duration
}

var (
	config  Configuration
	appName = "notify"
)

func failIf(err error, msg string) {
	if err != nil {
		log.Fatalf("error %s: %v", msg, err)
	}
}

func init() {
	help := func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(
			os.Stderr,
			"\nDisplays text read from stdin in a pop-up notification window\n",
		)
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Usage = help

	dimension := flag.String(
		"g",
		"0x0+20+20",
		`window dimensions as "<width>x<height>+<x>+<y>". 
If width = 0 or height = 0, window dimensions are set to fit the text`)
	fontName := flag.String(
		"f",
		"Inconsolata",
		"font pattern as passed to fc-match")
	fontSize := flag.Float64(
		"s",
		30,
		"font size")
	outputString := flag.String(
		"e",
		"echo 'done'",
		"string that is printed to stdout after notification closes")
	duration := flag.Duration(
		"d",
		3*time.Second,
		"duration after which notifcation closes")
	backgroundColor := flag.String(
		"B",
		"#000",
		"background color in hex format #rrggbb or #rgb")
	foregroundColor := flag.String(
		"F",
		"#fff",
		"foreground color in hex format #rrggbb or #rgb")

	flag.Parse()

	{
		dim, err := parsing.ParseDimension(*dimension)
		failIf(err, "parse dimension")

		if dim.Width == 0 || dim.Height == 0 {
			config.autoDimension = true
		} else {
			config.winWidth = dim.Width
			config.winHeight = dim.Height
		}
		config.winX = dim.X
		config.winY = dim.Y
	}

	{
		face, err := fonts.LoadTTFontFromPattern(*fontName, *fontSize)
		failIf(err, "load font")
		config.fontFace = face
	}
	{
		clr, err := parsing.ParseColor(*backgroundColor)
		failIf(err, "parse background color")
		config.bgColor = clr
	}
	{
		clr, err := parsing.ParseColor(*foregroundColor)
		failIf(err, "parse foreground color")
		config.fgColor = clr
	}

	config.fontSize = *fontSize
	config.duration = *duration
	config.outputString = *outputString
}

// TODO: register mouse clicks
// TODO: screen number
// TODO: image support
// TODO: display ... for text that is cut off
// TODO: accept input [title] body
// TODO: accept neg win position

func setupWindow(
	bounds pixel.Rect,
	position pixel.Vec,
) (*pixelgl.Window, error) {
	cfg := pixelgl.WindowConfig{
		Title:       appName,
		Bounds:      bounds,
		Position:    position,
		VSync:       true,
		Undecorated: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}
	win.SetSmooth(true)
	return win, nil
}

func max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func drawRectangle(
	r pixel.Rect,
	clr color.Color,
	imd *imdraw.IMDraw,
) {
	imd.Color = clr
	vs := r.Vertices()
	imd.Push(vs[0], vs[1], vs[2], vs[3])
	imd.Polygon(0)
}

func run() {
	const (
		minWinHeight = 40
		padding      = 10
		borderWidth  = 4
	)
	var displayText *text.Text
	{
		atlas := text.NewAtlas(config.fontFace, text.ASCII)
		displayText = text.New(pixel.ZV, atlas)
	}

	var text string
	{
		bytes, err := io.ReadAll(os.Stdin)
		failIf(err, "read from stdin")
		text = string(bytes)
	}

	imd := imdraw.New(nil)

	var textBox, paddingBox, borderBox pixel.Rect
	{
		textBox = displayText.BoundsOf(text)
		textBox = textBox.Resized(
			textBox.Center(),
			pixel.V(textBox.W(), max(textBox.H(), minWinHeight)),
		)

		if config.autoDimension {
			paddingBox = textBox.Resized(
				textBox.Center(),
				pixel.V(textBox.W()+2*padding, textBox.H()+2*padding),
			)
		} else {
			paddingBox = textBox.Resized(
				textBox.Center(),
				pixel.V(config.winWidth-2*borderWidth, config.winHeight-2*borderWidth),
			)
		}
		borderBox = paddingBox.Resized(
			paddingBox.Center(),
			pixel.V(paddingBox.W()+2*borderWidth, paddingBox.H()+2*borderWidth),
		)
	}

	borderColor := colornames.White

	drawRectangle(borderBox, borderColor, imd)

	drawRectangle(paddingBox, config.bgColor, imd)

	win, err := setupWindow(borderBox, pixel.V(config.winX, config.winY))
	failIf(err, "create window")

	fmt.Fprint(displayText, text)

	for !win.Closed() {
		imd.Draw(win)
		displayText.DrawColorMask(
			win,
			pixel.IM,
			config.fgColor,
		)
		win.Update()
	}
	// time.Sleep(config.duration)
	// fmt.Fprint(os.Stdout, config.outputString)
}

func main() {
	pixelgl.Run(run)
}
