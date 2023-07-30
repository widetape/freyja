package freyja

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Slider struct {
	Origin widget.Float

	Background         color.NRGBA
	BackgroundDisabled color.NRGBA
	BackgroundWidth    unit.Dp

	Knob         color.NRGBA
	KnobDisabled color.NRGBA
	KnobSize     unit.Dp
	KnobShadow   Shadow

	Tint color.NRGBA
}

func (s *Slider) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Left: s.KnobSize / 2, Right: s.KnobSize / 2}.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {
			var (
				disabled = gtx.Queue == nil
				knobSize = gtx.Dp(s.KnobSize)
				size     = image.Pt(gtx.Constraints.Max.X, knobSize)
			)
			func() {
				var (
					width  = gtx.Dp(s.BackgroundWidth)
					radius = width / 2
					shape  = clip.UniformRRect(
						image.Rectangle{
							Max: image.Pt(
								size.X,
								width,
							),
						},
						radius,
					)
					color color.NRGBA
				)
				if disabled {
					color = s.BackgroundDisabled
				} else {
					color = s.Background
				}
				defer op.Offset(image.Pt(0, (knobSize/2)-(width/2))).Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, color, shape.Op(gtx.Ops))
			}()
			func() {
				var fgtx = gtx
				fgtx.Constraints.Min = image.Pt(size.X, size.Y)
				s.Origin.Layout(
					fgtx,
					layout.Horizontal,
					0, 1, false,
					knobSize,
				)
			}()
			func() {
				var (
					shape = clip.Ellipse{
						Max: image.Pt(knobSize, knobSize),
					}
					color color.NRGBA
				)
				if disabled {
					color = s.KnobDisabled
				} else {
					color = s.Knob
				}
				defer op.Offset(image.Pt(int(s.Origin.Pos()-float32(knobSize/2)), 0)).Push(gtx.Ops).Pop()
				s.KnobShadow.Layout(
					gtx,
					shape.Path(gtx.Ops),
					func(gtx layout.Context) layout.Dimensions {
						paint.FillShape(gtx.Ops, color, shape.Op(gtx.Ops))
						return layout.Dimensions{}
					},
				)
			}()
			return layout.Dimensions{
				Size: size,
			}
		},
	)
}
