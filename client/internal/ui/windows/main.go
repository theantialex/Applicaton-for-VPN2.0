package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MainWindow(w fyne.Window, errCh chan error) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Добро пожаловать в VPN 2.0!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	container := container.NewVBox(
		label,
		layout.NewSpacer(),
		widget.NewButton("Создать сеть", func() {
			w.SetContent(CreateWindow(w, errCh))
		}),
		widget.NewButton("Соединиться с сетью", func() {
			w.SetContent(ConnectWindow(w, errCh))
		}),
		widget.NewButton("Удалить сеть", func() {
			w.SetContent(DeleteWindow(w, errCh))
		}),
		layout.NewSpacer(),
	)
	return container
}
