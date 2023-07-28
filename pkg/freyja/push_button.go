package freyja

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// PushButton is a button with text.
type PushButton struct {
	Origin widget.Clickable // Origin is the clickable of this push button.

	Background         op.CallOp // Background is called to fill the background of this button.
	BackgroundDisabled op.CallOp // BackgroundDisabled is used instead of Background in disabled mode.
	CornerRadius       unit.Dp   // CornerRadius is the radius of smooth corners.

	Inset layout.Inset // Inset is used to margin the text from corners of this button.

	Shaper *text.Shaper // Shaper is used to layout the text.
	Font   font.Font    // Font is used for the text.

	Label              string    // Label is the text.
	FontSize           unit.Sp   // FontSize is the size of the text.
	Foreground         op.CallOp // Foreground is the material operation for the text.
	ForegroundDisabled op.CallOp // ForegroundDisabled is used instead of Foreground in disabled mode.

	HoverColor color.NRGBA // HoverColor is drawn over the push button when it's hovered.
	ClickColor color.NRGBA // ClickColor is drawn over the push button while it's being pressed.
}

// Layout lays PushButton out to the context.
func (b *PushButton) Layout(gtx layout.Context) layout.Dimensions {
	min := gtx.Constraints.Min
	disabled := gtx.Queue == nil
	return b.Origin.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {
			semantic.Button.Add(gtx.Ops)
			return layout.Stack{Alignment: layout.Center}.Layout(
				gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Min
					shape := clip.UniformRRect(image.Rectangle{Max: size}, gtx.Dp(b.CornerRadius))
					func() {
						defer shape.Push(gtx.Ops).Pop()
						if disabled {
							b.BackgroundDisabled.Add(gtx.Ops)
						} else {
							b.Background.Add(gtx.Ops)
							if b.Origin.Pressed() {
								paint.Fill(gtx.Ops, b.ClickColor)
							} else {
								if b.Origin.Hovered() {
									paint.Fill(gtx.Ops, b.HoverColor)
								}
							}
						}
					}()
					return layout.Dimensions{Size: size}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = min
					return layout.Center.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							return b.Inset.Layout(
								gtx,
								func(gtx layout.Context) layout.Dimensions {
									var color op.CallOp
									if disabled {
										color = b.ForegroundDisabled
									} else {
										color = b.Foreground
									}
									return widget.Label{MaxLines: 1}.Layout(
										gtx,
										b.Shaper,
										b.Font,
										b.FontSize,
										b.Label,
										color,
									)
								},
							)
						},
					)
				}),
			)
		},
	)
}
