package freyja

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type TextField struct {
	Origin widget.Editor

	LeadingContent  layout.Widget
	TrailingContent layout.Widget
	Spacing         unit.Dp

	Shadow Shadow

	Background         color.NRGBA
	BackgroundDisabled color.NRGBA

	BorderColor         color.NRGBA
	BorderColorDisabled color.NRGBA
	BorderWidth         unit.Dp
	BorderRadius        unit.Dp

	Inset layout.Inset

	OutlineColor color.NRGBA
	OutlineWidth unit.Dp

	Font   font.Font
	Shaper *text.Shaper

	SelectionColor color.NRGBA

	FontColor         color.NRGBA
	FontColorDisabled color.NRGBA
	FontSize          unit.Sp

	Hint              string
	HintColor         color.NRGBA
	HintColorDisabled color.NRGBA
}

func (t *TextField) Layout(gtx layout.Context) layout.Dimensions {
	var disabled = gtx.Queue == nil
	return layout.Stack{Alignment: layout.Center}.Layout(
		gtx,
		layout.Expanded(
			func(gtx layout.Context) layout.Dimensions {
				var (
					size  = gtx.Constraints.Min
					shape = clip.UniformRRect(image.Rectangle{Max: size}, gtx.Dp(t.BorderRadius))
				)
				if t.Origin.Focused() {
					var stroke = clip.Stroke{
						Path:  shape.Path(gtx.Ops),
						Width: float32(gtx.Dp(t.OutlineWidth * 2)),
					}
					paint.FillShape(
						gtx.Ops,
						t.OutlineColor,
						stroke.Op(),
					)
				} else {
					t.Shadow.Layout(
						gtx,
						shape.Path(gtx.Ops),
						func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{}
						},
					)
				}
				defer shape.Push(gtx.Ops).Pop()
				if disabled {
					paint.Fill(gtx.Ops, t.BackgroundDisabled)
				} else {
					paint.Fill(gtx.Ops, t.Background)
				}
				var border = clip.Stroke{
					Path:  shape.Path(gtx.Ops),
					Width: float32(gtx.Dp(t.BorderWidth * 2)),
				}
				paint.FillShape(gtx.Ops, t.BorderColor, border.Op())
				return layout.Dimensions{Size: size}
			},
		),
		layout.Stacked(
			func(gtx layout.Context) layout.Dimensions {
				return t.Inset.Layout(
					gtx,
					func(gtx layout.Context) layout.Dimensions {
						var (
							flex = layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
								Spacing:   layout.SpaceAround,
							}
							spacer = layout.Spacer{Width: t.Spacing}
						)
						return flex.Layout(
							gtx,
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									if t.LeadingContent != nil {
										return t.LeadingContent(gtx)
									}
									return layout.Dimensions{}
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									if t.LeadingContent != nil {
										return spacer.Layout(gtx)
									}
									return layout.Dimensions{}
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									gtx.Constraints.Min.Y = 0
									textColorRecord := op.Record(gtx.Ops)
									if disabled {
										paint.Fill(gtx.Ops, t.FontColorDisabled)
									} else {
										paint.Fill(gtx.Ops, t.FontColor)
									}
									textColor := textColorRecord.Stop()
									selectionColorRecord := op.Record(gtx.Ops)
									if disabled {
										paint.Fill(gtx.Ops, color.NRGBA{})
									} else {
										paint.Fill(gtx.Ops, t.SelectionColor)
									}
									selectionColor := selectionColorRecord.Stop()
									hintColorRecord := op.Record(gtx.Ops)
									if disabled {
										paint.Fill(gtx.Ops, t.HintColorDisabled)
									} else {
										paint.Fill(gtx.Ops, t.HintColor)
									}
									hintColor := hintColorRecord.Stop()
									if t.Origin.Len() == 0 {
										widget.Label{MaxLines: 1}.Layout(
											gtx,
											t.Shaper,
											t.Font,
											t.FontSize,
											t.Hint,
											hintColor,
										)
									}
									return t.Origin.Layout(
										gtx,
										t.Shaper,
										t.Font,
										t.FontSize,
										textColor,
										selectionColor,
									)
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									if t.TrailingContent != nil {
										return spacer.Layout(gtx)
									}
									return layout.Dimensions{}
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									if t.TrailingContent != nil {
										return t.TrailingContent(gtx)
									}
									return layout.Dimensions{}
								},
							),
						)
					},
				)
			},
		),
	)
}
