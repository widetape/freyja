package freyja

import (
	"image"
	"image/color"
	"math"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Shadow draws a component and a shadow below it.
//
// Note that non-opaque components may produce unexpected look.
type Shadow struct {
	Color  color.NRGBA // Color is the basic of the shadow.
	Layers int         // Layers determines how many shadow layer to draw.
	Spread unit.Dp     // Spread determines how far the shadow goes from the content.
	X      unit.Dp     // X is the horizontal offset of this shadow.
	Y      unit.Dp     // Y is the vertical offset of this shadow.
	Slope  float64
}

// Layout lays out the content and a shadow below it.
func (s *Shadow) Layout(gtx layout.Context, shape clip.PathSpec, content layout.Widget) layout.Dimensions {
	contentRecord := op.Record(gtx.Ops)
	contentDimensions := content(gtx)
	contentOp := contentRecord.Stop()
	func() {
		var (
			spread = gtx.Dp(s.Spread)
			offset = op.Offset(
				image.Pt(
					gtx.Dp(s.X),
					gtx.Dp(s.Y),
				),
			)
		)
		defer offset.Push(gtx.Ops).Pop()
		for i := 0; i < s.Layers; i++ {
			var (
				distance = float32(i+1) / float32(s.Layers)
				width    = float32(spread) * float32(1-math.Exp(-1*math.Pow(float64(distance), s.Slope)))
				alpha    = float32(s.Color.A) * distance
				color    = color.NRGBA{
					R: s.Color.R,
					G: s.Color.G,
					B: s.Color.B,
					A: uint8(alpha * (1 - distance)),
				}
			)
			shadowLayer(gtx.Ops, width, shape, color)
		}
	}()
	contentOp.Add(gtx.Ops)
	return contentDimensions
}

// shadowLayer draws a shadow layer into the ops.
func shadowLayer(ops *op.Ops, width float32, shape clip.PathSpec, color color.NRGBA) {
	var clip = clip.Stroke{
		Path:  shape,
		Width: width,
	}
	defer clip.Op().Push(ops).Pop()
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
