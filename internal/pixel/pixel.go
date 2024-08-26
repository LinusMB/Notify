package pixel

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func SetupWindow(
	title string,
	winWidth, winHeight, winX, winY float64,
) (*pixelgl.Window, error) {
	winBox := createBox(winWidth, winHeight, pixel.ZV)
	monW, monH := pixelgl.PrimaryMonitor().Size()
	position := pixel.V(winX, winY)
	if winX < 0 {
		position.X = monW + winX - winWidth
	}
	if winY < 0 {
		position.Y = monH + winY - winHeight
	}
	cfg := pixelgl.WindowConfig{
		Title:       title,
		Bounds:      winBox,
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
