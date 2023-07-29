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

// RadioButton represents a selectable option in
// a group of mutually exclusive choices.
type RadioButton struct {
	Group *widget.Enum // Group is the group this radio button belongs to.
	Key   string       // Key is this radio button's key in the group.

	Background         op.CallOp // Background is used to render the background of this radio button.
	BackgroundDisabled op.CallOp // BackgroundDisabled is used instead of Background in disabled mode.

	Outline         op.CallOp // Outline is used to render the outline of this radio button.
	OutlineDisabled op.CallOp // OutlineDisabled is used instead of Outline in disabled mode.
	OutlineWidth    unit.Dp   // OutlineWidth is the width of the outline.

	Tint op.CallOp // Tint is used to render the tint of the radio button when it's selected.

	Inset unit.Dp // Inset is the padding of the knob.

	Knob         op.CallOp // Knob is used to render the knob of the radio button when it's selected.
	KnobDisabled op.CallOp // KnobDisabled is used instead of Knob in disabled mode.
	KnobSize     unit.Dp   // KnobSize is the diameter if the knob.
}

// Layout lays the radio button out to the context.
func (b *RadioButton) Layout(gtx layout.Context) layout.Dimensions {
	var (
		disabled = gtx.Queue == nil
		active   = b.Group.Value == b.Key
	)
	return b.Group.Layout(
		gtx,
		b.Key,
		func(gtx layout.Context) layout.Dimensions {
			semantic.RadioButton.Add(gtx.Ops)
			return layout.Stack{Alignment: layout.Center}.Layout(
				gtx,
				layout.Expanded(
					func(gtx layout.Context) layout.Dimensions {
						var (
							size  = gtx.Constraints.Min
							shape = clip.Ellipse{Max: size}
							path  = shape.Path(gtx.Ops)
						)
						func() {
							defer shape.Push(gtx.Ops).Pop()
							if disabled {
								b.BackgroundDisabled.Add(gtx.Ops)
							} else {
								b.Background.Add(gtx.Ops)
								if active {
									b.Tint.Add(gtx.Ops)
								}
							}
						}()
						func() {
							var stroke = clip.Stroke{
								Width: float32(gtx.Dp(b.OutlineWidth)),
								Path:  path,
							}
							defer stroke.Op().Push(gtx.Ops).Pop()
							if disabled {
								b.OutlineDisabled.Add(gtx.Ops)
							} else {
								b.Outline.Add(gtx.Ops)
							}
						}()
						return layout.Dimensions{Size: size}
					},
				),
				layout.Stacked(
					func(gtx layout.Context) layout.Dimensions {
						var inset = layout.UniformInset(b.Inset)
						return inset.Layout(
							gtx,
							func(gtx layout.Context) layout.Dimensions {
								var (
									diameter = gtx.Dp(b.KnobSize)
									size     = image.Pt(diameter, diameter)
									shape    = clip.Ellipse{Max: size}
								)
								func() {
									defer shape.Push(gtx.Ops).Pop()
									if active {
										if disabled {
											b.KnobDisabled.Add(gtx.Ops)
										} else {
											b.Knob.Add(gtx.Ops)
										}
									}
								}()
								return layout.Dimensions{Size: size}
							},
						)
					},
				),
			)
		},
	)
}
