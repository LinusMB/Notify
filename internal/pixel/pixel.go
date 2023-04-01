package pixel

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

func SetupWindow(
	title string,
	winWidth, winHeight, winX, winY float64,
) (*pixelgl.Window, error) {
	winBox := createBox(winWidth, winHeight, pixel.ZV)
	cfg := pixelgl.WindowConfig{
		Title:       title,
		Bounds:      winBox,
		Position:    pixel.V(winX, winY),
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

func fillBox(
	imd *imdraw.IMDraw,
	r pixel.Rect,
	c color.Color,
) {
	imd.Color = c
	vs := r.Vertices()
	imd.Push(vs[0], vs[1], vs[2], vs[3])
	imd.Polygon(0)
}

func createBox(width, height float64, center pixel.Vec) pixel.Rect {
	return centerBox(pixel.R(0, 0, width, height), center)
}

func centerBox(box pixel.Rect, center pixel.Vec) pixel.Rect {
	return box.Moved(box.Center().To(center))
}

type NotificationWindow struct {
	imd *imdraw.IMDraw
}

func (nw *NotificationWindow) Draw(t pixel.Target) {
	nw.imd.Draw(t)
}

func SetupNotificationWindow(
	winWidth, winHeight, borderWidth float64,
	winColor, borderColor color.Color,
) *NotificationWindow {
	imd := imdraw.New(nil)
	contentBox := createBox(
		winWidth-2*borderWidth,
		winHeight-2*borderWidth,
		pixel.ZV,
	)
	borderBox := createBox(winWidth, winHeight, pixel.ZV)
	fillBox(imd, borderBox, borderColor)
	fillBox(imd, contentBox, winColor)
	nw := NotificationWindow{
		imd: imd,
	}
	return &nw
}

type NotificationText struct {
	title *text.Text
	body  *text.Text
}

func newNotificationText(
	c color.Color,
	regular, bold *text.Atlas,
) *NotificationText {
	nt := NotificationText{
		title: text.New(pixel.ZV, bold),
		body:  text.New(pixel.ZV, regular),
	}
	nt.title.Color = c
	nt.body.Color = c
	return &nt
}

func (nt *NotificationText) writeTitle(title string) {
	fmt.Fprintln(nt.title, title)
}

func (nt *NotificationText) writeBody(body string) {
	nt.body.Dot = nt.title.Dot
	fmt.Fprint(nt.body, body)
}

func (nt *NotificationText) bounds() pixel.Rect {
	return nt.title.Bounds().Union(nt.body.Bounds())
}

func (nt *NotificationText) W() float64 {
	box := nt.title.Bounds().Union(nt.body.Bounds())
	return box.W()
}

func (nt *NotificationText) H() float64 {
	box := nt.title.Bounds().Union(nt.body.Bounds())
	return box.H()
}

func (nt *NotificationText) Draw(t pixel.Target) {
	textBox := nt.bounds()
	v := textBox.Center().To(pixel.ZV)
	mat := pixel.IM.Moved(v)
	nt.title.Draw(t, mat)
	nt.body.Draw(t, mat)
}

func SetupNotificationText(
	fontRegular, fontBold font.Face,
	textColor color.Color,
	title, body string,
) *NotificationText {
	nt := newNotificationText(
		textColor,
		text.NewAtlas(fontRegular, text.ASCII),
		text.NewAtlas(fontBold, text.ASCII),
	)
	if title != "" {
		nt.writeTitle(title)
	}
	nt.writeBody(body)
	return nt
}
