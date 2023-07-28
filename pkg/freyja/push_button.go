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

type PushButton struct {
	Button widget.Clickable

	Background         op.CallOp
	BackgroundDisabled op.CallOp
	CornerRadius       unit.Dp

	Inset layout.Inset

	Shaper *text.Shaper
	Font   font.Font

	Label              string
	FontSize           unit.Sp
	Foreground         op.CallOp
	ForegroundDisabled op.CallOp

	HoverColor color.NRGBA
	ClickColor color.NRGBA
}

func (b *PushButton) Layout(gtx layout.Context) layout.Dimensions {
	min := gtx.Constraints.Min
	disabled := gtx.Queue == nil
	return b.Button.Layout(
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
							if b.Button.Pressed() {
								paint.Fill(gtx.Ops, b.ClickColor)
							} else {
								if b.Button.Hovered() {
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
