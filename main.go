package main

import (
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	client "vpn2.0/app/client/cmd"
	"vpn2.0/app/client/config"
)

func CreateWindow(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Создание защищенной сети", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	name := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Название сети", name),
		widget.NewFormItem("Пароль", pass),
	)

	btnSubmit := widget.NewButton("Подтвердить", func() {
		w.SetContent(Create(w, name.Text, pass.Text))
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

func Create(w fyne.Window, name string, pass string) fyne.CanvasObject {
	manager, ctx := client.SetUpClient()

	resp, _ := manager.MakeCreateRequest(ctx, name, pass)

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
		w.Canvas(),
	)
	modal.Show()

	return MainWindow(w)
}

func ConnectWindow(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Соединение с защищенной сетью", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	name := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Название сети", name),
		widget.NewFormItem("Пароль", pass),
	)

	btn := widget.NewButton("Подтвердить", func() {
		w.SetContent(Connect(w, name.Text, pass.Text))
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

func Connect(w fyne.Window, name string, pass string) fyne.CanvasObject {
	manager, ctx := client.SetUpClient()

	resp, _ := manager.MakeConnectRequest(ctx, name, pass)

	var modal *widget.PopUp
	var label *widget.Label

	answer := strings.Split(resp, " ")

	if answer[0] == "success" {
		label = widget.NewLabel("Соединение с сетью " + name + " было успешно установлено")
		modal = widget.NewModalPopUp(
			container.NewVBox(
				label,
				widget.NewButton("Закрыть", func() { modal.Hide() }),
			),
			w.Canvas(),
		)
		modal.Show()
		return ChatWindow(w, answer[1])

	}

	label = widget.NewLabel("Произошла ошибка. Попробуйте позже.")
	modal = widget.NewModalPopUp(
		container.NewVBox(
			label,
			widget.NewButton("Закрыть", func() { modal.Hide() }),
		),
		w.Canvas(),
	)
	modal.Show()
	return MainWindow(w)
}

func ChatWindow(w fyne.Window, addr string) fyne.CanvasObject {
	label := widget.NewLabel("Ваш адрес в сети: " + addr)
	chat := widget.NewMultiLineEntry()
	chat.Disable()

	ip := widget.NewEntry()
	form := widget.NewForm(
		widget.NewFormItem("IP пользователя", ip),
	)

	input := widget.NewEntry()
	btn := widget.NewButton("Отправить", func() {
		SendMessage(input.Text, ip.Text, chat)
	})

	go processChatResponces(chat)

	return container.NewVBox(
		label,
		form,
		chat,
		input,
		btn,
	)
}

func SendMessage(msg string, addr string, chat *widget.Entry) {
	cmd := exec.Command("sh", "send.sh", msg, addr, config.PORT)
	cmd.Start()
	chat.SetText(chat.Text + "\n" + msg)
}

func processChatResponces(chat *widget.Entry) {
	file, _ := os.Open("text.txt")

	var bufPool = make([]byte, 1500)
	for {
		n, err := file.Read(bufPool)

		if n < 1 || err != nil {
			continue
		}
		chat.SetText(chat.Text + "\n" + string(bufPool))
	}

}

func MainWindow(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Добро пожаловать в VPN 2.0!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	container := container.NewVBox(
		label,
		layout.NewSpacer(),
		widget.NewButton("Создать сеть", func() {
			w.SetContent(CreateWindow(w))
		}),
		widget.NewButton("Соединиться с сетью", func() {
			w.SetContent(ConnectWindow(w))
		}),
		layout.NewSpacer(),
	)
	return container
}

func RunScript() {
	cmd := exec.Command("sh", "read.sh")
	cmd.Start()
}

func main() {

	RunScript()

	a := app.New()
	w := a.NewWindow("VPN 2.0")
	w.Resize(fyne.NewSize(800, 500))

	container := MainWindow(w)
	w.SetContent(container)

	w.ShowAndRun()

}
