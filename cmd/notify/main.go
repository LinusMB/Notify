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
	"strings"
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
	fontFaceRegular font.Face
	fontFaceBold    font.Face
	fontSize        float64
	winWidth        float64
	winHeight       float64
	winX            float64
	winY            float64
	autoDimension   bool
	bgColor         color.Color
	fgColor         color.Color
	outputString    string
	duration        time.Duration
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
	fontFamily := flag.String(
		"f",
		"",
		`font family. The font path to regular and bold font of type font family is searched using fc-match (make sure that fontconfig is installed).
If -f and -fp are unspecified the default font (Inconsolata) is used.
If both -f and -fp are specified, -fp is preferred.
Example: -f "Inconsolata"`,
	)
	fontPaths := flag.String(
		"fp",
		"",
		`font path to regular and bold font as "<regular-font-path>,<bold-font-path>.
The font is loaded from the font file at <*-font-path>.
If -f and -fp are unspecified the default font (Inconsolata) is used.
If both -f and -fp are specified, -fp is preferred.
Example: -fp "/usr/share/fonts/TTF/Inconsolata-Regular.ttf,/usr/share/fonts/TTF/Inconsolata-Bold.ttf"`,
	)
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
		var (
			fs  *fonts.FontSet
			err error
		)
		if *fontPaths != "" {
			const (
				REGULAR = iota
				BOLD    = iota
			)
			ts := strings.SplitN(*fontPaths, ",", 2)
			fs, err = fonts.LoadOpentypeFontSetFromPaths(
				ts[REGULAR],
				ts[BOLD],
				*fontSize,
			)
		} else if *fontFamily != "" {
			fs, err = fonts.LoadOpentypeFontSetFromFamily(*fontFamily, *fontSize)
		} else {
			fs, err = fonts.LoadOpentypeFontSetDefault(*fontSize)
		}
		failIf(err, "load font")
		config.fontFaceRegular = fs.Regular
		config.fontFaceBold = fs.Bold
	}
	{
		c, err := parsing.ParseColor(*backgroundColor)
		failIf(err, "parse background color")
		config.bgColor = c
	}
	{
		c, err := parsing.ParseColor(*foregroundColor)
		failIf(err, "parse foreground color")
		config.fgColor = c
	}

	config.fontSize = *fontSize
	config.duration = *duration
	config.outputString = *outputString
}

// TODO: register mouse clicks
// TODO: screen number
// TODO: image support
// TODO: display ... for text that is cut off
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
	c color.Color,
	imd *imdraw.IMDraw,
) {
	imd.Color = c
	vs := r.Vertices()
	imd.Push(vs[0], vs[1], vs[2], vs[3])
	imd.Polygon(0)
}

type textContent struct {
	title *text.Text
	body  *text.Text
}

func newTextContent(c color.Color, regular, bold *text.Atlas) *textContent {
	tc := textContent{
		title: text.New(pixel.ZV, bold),
		body:  text.New(pixel.ZV, regular),
	}
	tc.title.Color = c
	tc.body.Color = c
	return &tc
}

func (tc *textContent) writeTitle(title string) {
	fmt.Fprintln(tc.title, title)
}

func (tc *textContent) writeBody(body string) {
	tc.body.Dot = tc.title.Dot
	fmt.Fprint(tc.body, body)
}

func (tc *textContent) bounds() pixel.Rect {
	return tc.title.Bounds().Union(tc.body.Bounds())
}

func (tc *textContent) draw(t pixel.Target) {
	tc.title.Draw(t, pixel.IM)
	tc.body.Draw(t, pixel.IM)
}

func run() {
	const (
		minWinHeight = 40
		padding      = 10
		borderWidth  = 4
	)

	textContent := newTextContent(
		config.fgColor,
		text.NewAtlas(config.fontFaceRegular, text.ASCII),
		text.NewAtlas(config.fontFaceBold, text.ASCII),
	)

	var notification *parsing.Notification
	{
		bytes, err := io.ReadAll(os.Stdin)
		failIf(err, "read from stdin")
		input := string(bytes)
		notification = parsing.ParseNotification(input)
	}

	if notification.Title != "" {
		textContent.writeTitle(notification.Title)
	}
	textContent.writeBody(notification.Body)

	imd := imdraw.New(nil)

	var textBox, paddingBox, borderBox pixel.Rect
	{
		textBox = textContent.bounds()
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

	for !win.Closed() {
		imd.Draw(win)
		textContent.draw(win)
		win.Update()
	}
	// time.Sleep(config.duration)
	// fmt.Fprint(os.Stdout, config.outputString)
}

func main() {
	pixelgl.Run(run)
}
