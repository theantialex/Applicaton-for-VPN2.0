package windows

import (
	"context"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (ac *AppContainer) MainWindow(ctx context.Context, errCh chan error) fyne.CanvasObject {
	if ac.window == nil {
		errCh <- errors.New("nil window")
	}

	label := widget.NewLabelWithStyle("Добро пожаловать в VPN 2.0!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	container := container.NewVBox(
		label,
		layout.NewSpacer(),
		widget.NewButton("Создать сеть", func() {
			ac.window.SetContent(ac.CreateWindow(ctx, errCh))
		}),
		widget.NewButton("Соединиться с сетью", func() {
			ac.window.SetContent(ac.ConnectWindow(ctx, errCh))
		}),
		widget.NewButton("Удалить сеть", func() {
			ac.window.SetContent(ac.DeleteWindow(ctx, errCh))
		}),
		layout.NewSpacer(),
	)
	return container
}
