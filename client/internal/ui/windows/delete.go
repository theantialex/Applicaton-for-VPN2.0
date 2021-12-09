package windows

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"vpn2.0/app/client/internal/client"
)

func DeleteWindow(w fyne.Window, errCh chan error) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Удаление защищенной сети", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	name := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Название сети", name),
		widget.NewFormItem("Пароль", pass),
	)

	btnSubmit := widget.NewButton("Подтвердить", func() {
		w.SetContent(Delete(w, name.Text, pass.Text, errCh))
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

func Delete(w fyne.Window, name string, pass string, errCh chan error) fyne.CanvasObject {
	manager, ctx := client.SetUpClient()

	resp, _ := manager.MakeDeleteRequest(ctx, name, pass)

	var modal *widget.PopUp
	var label *widget.Label

	if strings.TrimRight(resp, "\n") == "network deleted successfully" {
		label = widget.NewLabel("Сеть " + name + " успешно удалена")
	} else {
		label = widget.NewLabel("Произошла ошибка. Попробуйте позже.")
	}

	modal = widget.NewModalPopUp(
		container.NewVBox(
			label,
			widget.NewButton("Закрыть", func() { modal.Hide() }),
		),
		w.Canvas(),
	)
	modal.Show()

	return MainWindow(w, errCh)
}
