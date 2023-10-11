package pages

import (
	"log/slog"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var loggedIn bool = true

func buildMenu(w fyne.Window) error {
	slog.Info("building menu", "loggedin", loggedIn)
	//File
	fileMenuItem := fyne.NewMenuItem("File", func() {
	})
	logout := fyne.NewMenuItem("Logout", func() {
		loggedIn = false
		//clear cookie file
		cookie, err := os.Create(os.TempDir() + "/cookie.timetrace")
		if err != nil {
			slog.Info("delete cookie store", "error", err)
		}
		cookie.Close()
		w.SetContent(BuildLoginPage(w))
	})
	fileMenu := fyne.NewMenu("File")
	fileMenu.Items = make([]*fyne.MenuItem, 0)
	//if loggedIn {
	fileMenu.Items = append(fileMenu.Items, fileMenuItem)
	fileMenu.Items = append(fileMenu.Items, logout)
	//}
	fileMenu.Items = append(fileMenu.Items, fyne.NewMenuItem("Quit", func() {
		w.Close()
	}))
	//Projects
	projectsMenu := fyne.NewMenu("Projects")
	projectsMenu.Items = make([]*fyne.MenuItem, 0)
	//Reports
	reportsMenu := fyne.NewMenu("Reports")
	reportsMenu.Items = make([]*fyne.MenuItem, 0)
	// About
	helpMenuItem := fyne.NewMenuItem("Help", func() {
		dialog.ShowInformation("About", "v0.1.0", w)
	})
	aboutMenu := fyne.NewMenu("About")
	aboutMenu.Items = make([]*fyne.MenuItem, 0)
	aboutMenu.Items = append(aboutMenu.Items, helpMenuItem)
	//MENU
	menu := fyne.MainMenu{}
	menu.Items = make([]*fyne.Menu, 0)
	menu.Items = append(menu.Items, fileMenu)
	//if loggedIn {
	menu.Items = append(menu.Items, projectsMenu)
	menu.Items = append(menu.Items, reportsMenu)
	//}
	menu.Items = append(menu.Items, aboutMenu)
	w.SetMainMenu(&menu)
	return nil
}
