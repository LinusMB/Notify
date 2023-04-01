package pixel

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

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
