package main

import (
	ifont "Notify/internal/font"
	"Notify/internal/parsing"
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/exp/constraints"
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
	borderWidth     float64
	borderColor     color.Color
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
		"x+20+20",
		`window dimensions as "<width>x<height>+<x>+<y>".
If <width> or <height> are unspecified, the window size is set to fit the text.
If <x> or <y> are unspecified their value is set to 0 respectively.`)
	fontFamily := flag.String(
		"f",
		"",
		`font family. The font path to regular and bold font of type font family is searched using fc-match (make sure that fontconfig is installed).
Note that only otf and ttf fonts are supported.
If -f and -fp are unspecified the default font (Inconsolata) is used.
If both -f and -fp are specified, -fp is preferred.
Example: -f "Inconsolata"`,
	)
	fontPaths := flag.String(
		"fp",
		"",
		`font path to regular and bold font as "<regular-font-path>,<bold-font-path>.
The font is loaded from the font file at <*-font-path>.
Note that only otf and ttf fonts are supported.
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
		6*time.Second,
		`duration after which the notifcation window closes. 
If -d 0 is given, the notfication window will not close.`)
	borderWidth := flag.Float64(
		"bw",
		2,
		"border width")
	borderColor := flag.String(
		"bc",
		"#fff",
		"border color in hex format #rrggbb or #rgb")
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

		config.winWidth = dim.Width
		config.winHeight = dim.Height
		config.winX = dim.X
		config.winY = dim.Y
	}

	{
		var (
			fs  *ifont.FontSet
			err error
		)
		if *fontPaths != "" {
			const (
				REGULAR = iota
				BOLD    = iota
			)
			ts := strings.SplitN(*fontPaths, ",", 2)
			fs, err = ifont.LoadOpentypeFontSetFromPaths(
				ts[REGULAR],
				ts[BOLD],
				*fontSize,
			)
		} else if *fontFamily != "" {
			fs, err = ifont.LoadOpentypeFontSetFromFamily(*fontFamily, *fontSize)
		} else {
			fs, err = ifont.LoadOpentypeFontSetDefault(*fontSize)
		}
		failIf(err, "load font")
		config.fontFaceRegular = fs.Regular
		config.fontFaceBold = fs.Bold
	}
	{
		c, err := parsing.ParseColor(*borderColor)
		failIf(err, "parse border color")
		config.borderColor = c
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

	config.borderWidth = *borderWidth
	config.fontSize = *fontSize
	config.duration = *duration
	config.outputString = *outputString
}

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
	const minWinHeight = 40

	padding := math.Round(config.fontSize / 2)

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

		if config.winWidth == 0 || config.winHeight == 0 {
			paddingBox = textBox.Resized(
				textBox.Center(),
				pixel.V(textBox.W()+2*padding, textBox.H()+2*padding),
			)
		} else {
			paddingBox = textBox.Resized(
				textBox.Center(),
				pixel.V(config.winWidth-2*config.borderWidth, config.winHeight-2*config.borderWidth),
			)
		}
		borderBox = paddingBox.Resized(
			paddingBox.Center(),
			pixel.V(
				paddingBox.W()+2*config.borderWidth,
				paddingBox.H()+2*config.borderWidth,
			),
		)
	}

	drawRectangle(borderBox, config.borderColor, imd)

	drawRectangle(paddingBox, config.bgColor, imd)

	win, err := setupWindow(borderBox, pixel.V(config.winX, config.winY))
	failIf(err, "create window")

	imd.Draw(win)
	textContent.draw(win)

	var (
		exitCode int
		closeWin <-chan time.Time
	)

	if config.duration == 0 {
		closeWin = time.After(time.Duration(math.MaxInt64))
	} else {
		closeWin = time.After(config.duration)
	}

Loop:
	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			exitCode = 0
			break Loop
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			exitCode = 1
			break Loop
		}
		select {
		case <-closeWin:
			break Loop
		default:
			win.Update()
		}
	}

	fmt.Fprint(os.Stdout, config.outputString)
	os.Exit(exitCode)
}

// TODO: accept neg win position
// TODO: configure ci

func main() {
	pixelgl.Run(run)
}
