package pages

import "fyne.io/fyne/v2"

type datePicker struct{}

func (d *datePicker) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		w = w + childSize.Width
		h = max(h, childSize.Height)
	}
	return fyne.NewSize(w, h)
}

func (d *datePicker) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		o.Move(pos)
		if i == 1 {
			o.Resize(fyne.Size{Width: o.MinSize().Width + 60, Height: o.MinSize().Height})
			pos = pos.Add(fyne.NewPos(o.MinSize().Width+60, 0))
		} else {
			o.Resize(o.MinSize())
			pos = pos.Add(fyne.NewPos(o.MinSize().Width, 0))
		}
	}
}
