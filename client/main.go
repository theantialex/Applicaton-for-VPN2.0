package main

import (
	"os/exec"

	"vpn2.0/app/client/internal/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"vpn2.0/app/client/internal/ui"
	"vpn2.0/app/client/internal/ui/windows"
)

func RunScript() {
	cmd := exec.Command("/bin/sh", "read.sh")
	cmd.Start()
}

func main() {
	RunScript()

	a := app.New()
	w := a.NewWindow("VPN 2.0")
	w.Resize(fyne.NewSize(float32(ui.WIDTH), float32(ui.HEIGHT)))

	errCh := make(chan error)
	go processError(a, w, errCh)

	manager, ctx := client.SetUpClient()

	clientApp := windows.NewAppContainer(w, manager)

	container := clientApp.MainWindow(ctx, errCh)
	w.SetContent(container)

	w.ShowAndRun()
}

func processError(a fyne.App, w fyne.Window, errCh chan error) {
	var caughtErr error

	var modal *widget.PopUp
	var label *widget.Label

	for caughtErr == nil {
		select {
		case errCheck := <-errCh:
			label = widget.NewLabel("Произошла ошибка: " + errCheck.Error())
			modal = widget.NewModalPopUp(
				container.NewVBox(
					label,
					widget.NewButton("Закрыть", func() { a.Quit() }),
				),
				w.Canvas(),
			)
			modal.Show()

			caughtErr = errCheck
		default:
		}
	}
	close(errCh)
}
