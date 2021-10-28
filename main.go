package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var authDB = map[string]string{"saxara": "1234"}

func LoginForm(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Login", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	login := widget.NewEntry()
	password := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Login", login),
		widget.NewFormItem("Password", password),
	)

	btn := widget.NewButton("Submit", func() {
		// fmt.Printf("%s %s\n", login.Text, password.Text)
		if _, exist := authDB[login.Text]; exist {
			w.SetContent(MainPage(w))
		} else {
			fmt.Printf("Error")
		}
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

func RegistrationForm(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Sign Up", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	nickname := widget.NewEntry()
	pass := widget.NewPasswordEntry()
	repeated_pass := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Nickname", nickname),
		widget.NewFormItem("Password", pass),
		widget.NewFormItem("Repeat password", repeated_pass),
	)

	btn := widget.NewButton("Submit", func() {
		if _, exist := authDB[nickname.Text]; !exist && pass.Text == repeated_pass.Text {
			authDB[nickname.Text] = pass.Text
			w.SetContent(MainPage(w))
		} else {
			fmt.Printf("Error")
		}
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

func MainPage(w fyne.Window) fyne.CanvasObject {
	label := widget.NewLabel("Main page")
	return container.NewVBox(
		label,
	)
}

func main() {
	a := app.New()
	w := a.NewWindow("VPN 2.0")
	w.Resize(fyne.NewSize(400, 250))

	label := widget.NewLabelWithStyle("Welcome to VPN 2.0!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	container := container.NewVBox(
		label,
		layout.NewSpacer(),
		widget.NewButton("Login", func() {
			w.SetContent(LoginForm(w))
		}),
		widget.NewButton("Register", func() {
			w.SetContent(RegistrationForm(w))
		}),
		layout.NewSpacer(),
	)
	w.SetContent(container)

	w.ShowAndRun()

}
