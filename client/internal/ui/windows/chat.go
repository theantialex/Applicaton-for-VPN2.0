package windows

import (
	"os"
	"os/exec"
	"vpn2.0/app/client/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"vpn2.0/app/client/internal/config"
)

func ChatWindow(w fyne.Window, addr string, ip string) fyne.CanvasObject {
	label := widget.NewLabel("Ваш адрес в сети: " + addr)
	chat := widget.NewMultiLineEntry()
	//chat.Disable()
	input := widget.NewEntry()
	btn := widget.NewButton("Отправить", func() {
		SendMessage(input.Text, ip, chat)
	})

	go processChatResponses(chat)

	return container.New(
		&ui.ChatLayout{},
		label,
		chat,
		input,
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
