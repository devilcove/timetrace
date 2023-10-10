package pages

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/assets"
)

type Route int

const (
	MainPage Route = iota
	LoginPage
)

func BuildMainPage(w fyne.Window) *fyne.Container {
	hello := widget.NewLabel("Hello World!")
	cookie, err := getCookie()
	if err != nil {
		return BuildLoginPage(w)
	}
	status, err := getStatus(cookie)
	if err != nil {
		return BuildLoginPage(w)
	}
	text := widget.NewTextGrid()
	text.SetText(status)
	c := container.NewVBox(
		hello,
		text,
		widget.NewButton("Hi", func() {
			hello.SetText("welcome)")
		}),
	)
	return c
}

func GetMainWindow(app fyne.App, title string) fyne.Window {
	w := app.NewWindow(title)
	w.SetMaster()
	if desktop, ok := app.(desktop.App); ok {
		tray := buildSystemTray(w)
		desktop.SetSystemTrayMenu(tray)
		icon := fyne.NewStaticResource("small", assets.SmallLogo)
		desktop.SetSystemTrayIcon(icon)
	}
	//w.SetCloseIntercept(func() {
	//w.Hide()
	//})
	buildMenu(w)
	buildWindow(w)
	return w
}

func buildWindow(w fyne.Window) error {
	w.Resize(fyne.Size{Width: 1024, Height: 768})
	Navigate(w, LoginPage)
	return nil
}

func buildSystemTray(w fyne.Window) *fyne.Menu {
	tray := fyne.NewMenu("Hello",
		fyne.NewMenuItem("open window", func() {
			log.Println("Tapped show")
			w.Show()
		}),
	)
	return tray
}

func buildMenu(w fyne.Window) error {
	//File
	fileMenuItem := fyne.NewMenuItem("File", func() {
	})
	logout := fyne.NewMenuItem("Logout", func() {
		w.SetContent(BuildLoginPage(w))
	})
	fileMenu := fyne.NewMenu("File")
	fileMenu.Items = make([]*fyne.MenuItem, 0)
	fileMenu.Items = append(fileMenu.Items, fileMenuItem)
	fileMenu.Items = append(fileMenu.Items, logout)
	fileMenu.Items = append(fileMenu.Items, fyne.NewMenuItem("Quit", func() {
		w.Close()
	}))
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
	menu.Items = append(menu.Items, aboutMenu)
	w.SetMainMenu(&menu)
	return nil
}

func Navigate(w fyne.Window, page Route) {
	switch page {
	case MainPage:
		w.SetContent(BuildMainPage(w))
	case LoginPage:
		w.SetContent(BuildLoginPage(w))
	}
}

func login(u, p string) (http.Cookie, error) {
	client := http.Client{}
	postData := struct {
		Username string
		Password string
	}{
		Username: u,
		Password: p,
	}
	body, err := json.Marshal(postData)
	if err != nil {
		return http.Cookie{}, err
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", bytes.NewBuffer(body))
	if err != nil {
		return http.Cookie{}, err
	}
	response, err := client.Do(req)
	if err != nil {
		return http.Cookie{}, err
	}
	defer response.Body.Close()
	ok := response.StatusCode >= 200 && response.StatusCode < 300
	if !ok {
		return http.Cookie{}, fmt.Errorf("status %s", response.Status)
	}
	for _, c := range response.Cookies() {
		if c.Name == "time" {
			return *c, nil
		}
	}
	return http.Cookie{}, fmt.Errorf("no cookie in response: status %s", response.Status)
}

func getCookie() (http.Cookie, error) {
	cookie := http.Cookie{}
	file, err := os.ReadFile("/tmp/cookie.timetrace")
	if err != nil {
		return cookie, err
	}
	if err := json.Unmarshal(file, &cookie); err != nil {
		return cookie, err
	}
	return cookie, nil
}

func getStatus(cookie http.Cookie) (string, error) {
	slog.Info("fetching current status")
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/projects/status", nil)
	if err != nil {
		slog.Error("http request", "error", err)
		return err.Error(), err
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		slog.Error("response", "error", err)
		return err.Error(), err
	}
	defer response.Body.Close()
	slog.Info("status response", "code", response.Status, "data", response.Body)
	if response.StatusCode != http.StatusOK {
		slog.Error("status code", "status", response.Status, "code", response.StatusCode)
		return response.Status, errors.New("response error")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err.Error(), err
	}
	return string(body), nil
}
