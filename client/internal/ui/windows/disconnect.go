package windows

import (
	"context"
	"strings"
	"vpn2.0/app/lib/cmd"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (ac *AppContainer) Disconnect(ctx context.Context, name string, pass string, errCh chan error) fyne.CanvasObject {
	resp, _ := ac.clientManager.MakeLeaveRequest(ctx, name, pass)

	var modal *widget.PopUp
	var label *widget.Label

	if strings.TrimRight(resp, "\n") == cmd.SuccessResponse {
		label = widget.NewLabel("Вы успешно отключились от сети")
	} else {
		label = widget.NewLabel("Произошла ошибка. Попробуйте позже.")
	}

	modal = widget.NewModalPopUp(
		container.NewVBox(
			label,
			widget.NewButton("Закрыть", func() { modal.Hide() }),
		),
		ac.window.Canvas(),
	)
	modal.Show()

	return ac.MainWindow(ctx, errCh)
}
