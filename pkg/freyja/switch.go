package freyja

import (
	"image"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
)

// Switch is a toggle button with knob.
type Switch struct {
	Origin widget.Bool // Origin is the bool if this switch.

	Background         op.CallOp // Background is used to draw the background for this switch.
	BackgroundDisabled op.CallOp // BackgroundDisabled is used instead of Background when the switch is disabled.

	Tint         op.CallOp // Tint is used to draw the tinted background overlay of the switch in "On" mode.
	TintDisabled op.CallOp // TintDisabled is used instead of Tint when the switch is disabled.

	Knob         op.CallOp // Knob is used to draw the knob of this switch.
	KnobDisabled op.CallOp // KnobDisabled is used instead of Knob when the switch is disabled.
	KnobSize     unit.Dp   // KnobSize is the diameter of the knob.

	Inset unit.Dp // Inset is the gab between the knob and the borders of switch.
	Shift unit.Dp // Shift is the distance that the knob shifts to the right.
}

// Layout lays Switch out to the context.
func (s *Switch) Layout(gtx layout.Context) layout.Dimensions {
	disabled := gtx.Queue == nil
	return s.Origin.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {
			semantic.Switch.Add(gtx.Ops)
			return layout.Stack{Alignment: layout.Center}.Layout(
				gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					var (
						radius = gtx.Dp(s.KnobSize)/2 + gtx.Dp(s.Inset)
						size   = gtx.Constraints.Min
						shape  = clip.UniformRRect(image.Rectangle{Max: size}, radius)
					)
					func() {
						defer shape.Push(gtx.Ops).Pop()
						if disabled {
							s.BackgroundDisabled.Add(gtx.Ops)
							if s.Origin.Value {
								s.TintDisabled.Add(gtx.Ops)
							}
						} else {
							s.Background.Add(gtx.Ops)
							if s.Origin.Value {
								s.Tint.Add(gtx.Ops)
							}
						}
					}()
					return layout.Dimensions{Size: size}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					inset := layout.Inset{
						Top: s.Inset, Bottom: s.Inset,
						Left: s.Inset, Right: s.Inset,
					}
					return inset.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							var (
								size      = gtx.Dp(s.KnobSize)
								shiftSize = gtx.Dp(s.Shift)
								shift     image.Point
								shape     = clip.Ellipse{Max: image.Pt(size, size)}
							)
							if s.Origin.Value {
								shift.X = shiftSize
							}
							defer op.Offset(shift).Push(gtx.Ops).Pop()
							func() {
								defer shape.Push(gtx.Ops).Pop()
								if disabled {
									s.KnobDisabled.Add(gtx.Ops)
								} else {
									s.Knob.Add(gtx.Ops)
								}
							}()
							return layout.Dimensions{Size: image.Pt(size+shiftSize, size)}
						},
					)
				}),
			)
		},
	)
}
