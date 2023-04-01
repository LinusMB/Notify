package pixel

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

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
