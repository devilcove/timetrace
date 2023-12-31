package pages

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetraced/models"
)

var currentUser models.User

func BuildLoginPage(w fyne.Window) *fyne.Container {
	buildMenu(w)
	// logo
	//logo := canvas.NewImageFromResource(fyne.NewStaticResource("logo", assets.Logo))
	//logo.SetMinSize(fyne.Size{Width: 150, Height: 80})
	//logoBox := container.NewCenter()
	//logoBox.Add(logo)

	// username input
	usernameTextbox := widget.NewEntry()
	usernameTextbox.SetPlaceHolder("Enter username")
	usernameBox := container.NewVBox(
		widget.NewLabel("Username"),
		usernameTextbox,
	)

	// password input
	passwordTextbox := widget.NewPasswordEntry()
	passwordTextbox.SetPlaceHolder("Enter password")
	passwordBox := container.NewVBox(
		widget.NewLabel("Password"),
		passwordTextbox,
	)

	// connect btn
	connectBtn := widget.NewButton("Connect", func() {
		err := login(usernameTextbox.Text, passwordTextbox.Text)
		if err != nil {
			slog.Error("failed to authenticate", "error", err)
			dialog.ShowError(fmt.Errorf("%w. failed to authenticate", err), w)
			return
		}
		w.SetContent(BuildMainPage(w))
	})

	// build layout
	vBox := container.NewVBox()
	//vBox.Add(logoBox)
	vBox.Add(usernameBox)
	vBox.Add(passwordBox)
	// TODO: add some top margin
	vBox.Add(connectBtn)
	w.SetContent(vBox)

	return vBox
}
