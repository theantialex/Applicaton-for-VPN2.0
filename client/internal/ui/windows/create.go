package windows

import (
	"context"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (ac *AppContainer) CreateWindow(ctx context.Context, errCh chan error) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Создание защищенной сети", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	name := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Название сети", name),
		widget.NewFormItem("Пароль", pass),
	)

	btnSubmit := widget.NewButton("Подтвердить", func() {
		ac.window.SetContent(ac.Create(ctx, name.Text, pass.Text, errCh))
	})

	return container.NewVBox(
		label,
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		btnSubmit,
		layout.NewSpacer(),
	)
}

func (ac *AppContainer) Create(ctx context.Context, name string, pass string, errCh chan error) fyne.CanvasObject {
	resp, _ := ac.clientManager.MakeCreateRequest(ctx, name, pass)

	var modal *widget.PopUp
	var label *widget.Label

	if strings.TrimRight(resp, "\n") == "network created successfully" {
		label = widget.NewLabel("Сеть " + name + " успешно создана")
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
