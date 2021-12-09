package windows

import (
	"os"
	"os/exec"
	"strings"

	"vpn2.0/app/client/internal/client"
	"vpn2.0/app/client/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"vpn2.0/app/client/internal/config"
)

func ChatWindow(w fyne.Window, addr string, ip string, name string, pass string, errCh chan error) fyne.CanvasObject {
	label := widget.NewLabel("Ваш адрес в сети: " + addr)
	chat := widget.NewMultiLineEntry()
	//chat.Disable()
	input := widget.NewEntry()
	btnSend := widget.NewButton("Отправить", func() {
		SendMessage(input.Text, ip, chat)
	})
	btn := widget.NewButton("Отключиться", func() {
		w.SetContent(Disconnect(w, name, pass, errCh))
	})

	go processChatResponses(chat)

	return container.New(
		&ui.ChatLayout{},
		label,
		chat,
		input,
		btnSend,
		btn,
	)
}

func SendMessage(msg string, addr string, chat *widget.Entry) {
	cmd := exec.Command("sh", "send.sh", msg, addr, config.PORT)
	cmd.Start()
	chat.SetText(chat.Text + "you: " + msg + "\n")
}

func processChatResponses(chat *widget.Entry) {
	file, _ := os.Open("text.txt")

	var bufPool = make([]byte, 1500)
	for {
		n, err := file.Read(bufPool)

		if n < 1 || err != nil {
			continue
		}
		chat.SetText(chat.Text + "them: " + string(bufPool))
	}

}

func Disconnect(w fyne.Window, name string, pass string, errCh chan error) fyne.CanvasObject {
	manager, ctx := client.SetUpClient()

	resp, _ := manager.MakeLeaveRequest(ctx, name, pass)

	var modal *widget.PopUp
	var label *widget.Label

	if strings.TrimRight(resp, "\n") == "network deleted successfully" {
		label = widget.NewLabel("Вы успешно отключились от сети")
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
