package ui

import "fyne.io/fyne/v2"

var WIDTH = 800
var HEIGHT = 500

type ChatLayout struct {
}

func (c *ChatLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}

func (c *ChatLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(containerSize.Width/2-100, 0)
	size := fyne.NewSize(containerSize.Width, 50)

	label := objects[0]
	label.Resize(size)
	label.Move(pos)

	pos = fyne.NewPos(0, size.Height)

	chatBox := objects[1]
	size = fyne.NewSize(containerSize.Width, 300)
	chatBox.Resize(size)
	chatBox.Move(pos)

	pos = pos.Add(fyne.NewPos(0, size.Height+20))

	input := objects[2]
	size = fyne.NewSize(containerSize.Width, 50)
	input.Resize(size)
	input.Move(pos)

	pos = pos.Add(fyne.NewPos(0, size.Height+20))

	btn := objects[3]
	btn.Resize(size)
	btn.Move(pos)
}
