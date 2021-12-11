package windows

import (
	"context"
	"strings"
	"vpn2.0/app/lib/cmd"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (ac *AppContainer) ConnectWindow(ctx context.Context, errCh chan error) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Соединение с защищенной сетью", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	name := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Название сети", name),
		widget.NewFormItem("Пароль", pass),
	)

	btn := widget.NewButton("Подтвердить", func() {
		ac.window.SetContent(ac.Connect(ctx, name.Text, pass.Text, errCh))
	})

	return container.NewVBox(
		label,
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		btn,
		layout.NewSpacer(),
	)
}

func (ac *AppContainer) Connect(ctx context.Context, name string, pass string, errCh chan error) fyne.CanvasObject {
	resp, _ := ac.clientManager.MakeConnectRequest(ctx, name, pass, errCh)

	var modal *widget.PopUp
	var label *widget.Label

	answer := strings.Split(resp, " ")

	if answer[0] == cmd.SuccessResponse {
		label = widget.NewLabel("Соединение с сетью " + name + " было успешно установлено")
		modal = widget.NewModalPopUp(
			container.NewVBox(
				label,
				widget.NewButton("Закрыть", func() { modal.Hide() }),
			),
			ac.window.Canvas(),
		)
		modal.Show()
		return ac.IPWindow(ctx, answer[1], name, pass, errCh)

	}

	label = widget.NewLabel("Произошла ошибка. Попробуйте позже.")
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

func (ac *AppContainer) IPWindow(ctx context.Context, addr string, name string, pass string, errCh chan error) fyne.CanvasObject {
	label1 := widget.NewLabel("Ваш адрес в сети: " + addr)
	label2 := widget.NewLabel("Введите ip адрес другого пользователя данной сети.")
	ip := widget.NewEntry()

	form := widget.NewForm(
		widget.NewFormItem("IP пользователя", ip),
	)
	btn := widget.NewButton("Отправить", func() {
		ac.window.SetContent(ac.ChatWindow(ctx, addr, ip.Text, name, pass, errCh))
	})

	return container.NewVBox(
		label1,
		label2,
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		btn,
		layout.NewSpacer(),
	)
}
