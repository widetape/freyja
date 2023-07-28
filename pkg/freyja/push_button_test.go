package freyja_test

import (
	"image"
	"image/color"
	"testing"
	"time"

	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/widetape/freyja/pkg/freyja"
)

func ExamplePushButton() {
	fonts := gofont.Collection()
	font := fonts[0].Font
	shaper := text.NewShaper(fonts)
	gtx := layout.Context{}

	button := freyja.PushButton{
		Label:    "Press Me!", // This text will be shown on the button.
		FontSize: unit.Sp(13), // Size of the text.
		Font:     font,        // Font of the text.
		Shaper:   shaper,      // Shaper to draw the text.

		CornerRadius: unit.Dp(12),
		Inset: layout.Inset{
			Top: unit.Dp(8), Bottom: unit.Dp(8),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},

		Background: func() op.CallOp {
			ops := op.Ops{}
			r := op.Record(&ops)
			// Fill button'a background with pure red.
			paint.Fill(&ops, color.NRGBA{R: 0xFF, A: 0xFF})
			return r.Stop()
		}(),
		BackgroundDisabled: func() op.CallOp {
			ops := op.Ops{}
			r := op.Record(&ops)
			// Fill button's background with light gray when it's disabled.
			paint.Fill(&ops, color.NRGBA{A: 0x10})
			return r.Stop()
		}(),

		Foreground: func() op.CallOp {
			ops := op.Ops{}
			r := op.Record(&ops)
			// Fill button's text color with black.
			paint.Fill(&ops, color.NRGBA{A: 0xFF})
			return r.Stop()
		}(),
		ForegroundDisabled: func() op.CallOp {
			ops := op.Ops{}
			r := op.Record(&ops)
			// Fill button's text color with light gray when it's disabled.
			//
			// Note that this won't be the same color as the disabled background,
			// but instead it will blend with it, because they are layed
			// on top of each other, so the text will still be readable.
			paint.Fill(&ops, color.NRGBA{A: 0x10})
			return r.Stop()
		}(),

		// These do not affect the button in disabled mode...
		// Both, HoverColor and ClickColor, do not affect the text,
		// but only affect the background.

		// Drawn on top of the button's background when the pointer
		// is hover on top of the button.
		//
		// Note that this does not appear when the button is being pressed.
		HoverColor: color.NRGBA{A: 0x08},

		// Drawn of top of the button's background when the pointer
		// is pressing the button.
		ClickColor: color.NRGBA{A: 0x0F},
	}

	// Layout the button to the context.
	button.Layout(gtx)
}

func BenchmarkPushButton_Layout(b *testing.B) {
	fonts := gofont.Collection()
	font := fonts[0].Font
	shaper := text.NewShaper(fonts)
	for i := 0; i < b.N; i++ {
		button := freyja.PushButton{
			Shaper:   shaper,
			Font:     font,
			FontSize: 13,
			Label:    "Some Button Text",
		}
		ops := op.Ops{}
		gtx := layout.Context{
			Constraints: layout.Exact(image.Pt(1000, 1000)),
			Metric: unit.Metric{
				PxPerDp: 1.0,
				PxPerSp: 1.0,
			},
			Queue: nil,
			Now:   time.Now(),
			Locale: system.Locale{
				Language:  "go",
				Direction: system.LTR,
			},
			Ops: &ops,
		}
		b.StartTimer()
		button.Layout(gtx)
		b.StopTimer()
	}
}
