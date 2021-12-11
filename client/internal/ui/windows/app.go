package windows

import (
	"fyne.io/fyne/v2"

	"vpn2.0/app/client/internal/client"
)

type AppContainer struct {
	window        fyne.Window
	clientManager *client.Manager
}

func NewAppContainer(w fyne.Window, client *client.Manager) *AppContainer {
	return &AppContainer{
		window:        w,
		clientManager: client,
	}
}
