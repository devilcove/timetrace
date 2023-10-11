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
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/assets"
	"github.com/devilcove/timetraced/models"
)

type Route int

const (
	MainPage Route = iota
	LoginPage
)

func BuildMainPage(w fyne.Window) *fyne.Container {
	buildMenu(w)
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("logo", assets.SmallLogo))
	logo.FillMode = canvas.ImageFillOriginal
	hello := widget.NewLabel("Hello World!")
	hello.Alignment = fyne.TextAlignCenter
	status, err := GetStatus()
	if err != nil {
		return BuildLoginPage(w)
	}
	text := widget.NewTextGrid()
	text.SetText(fmt.Sprintf("Current Project:\t%s\nTime This Session:\t%s\nTime Today:\t\t\t%s\n", status.Current, status.Elapsed, status.Total))
	stopButton := widget.NewButton("Stop    ", stop)
	c := container.NewVBox()
	c.Add(hello)
	c.Add(logo)
	h := container.NewCenter()
	h.Add(text)
	h2 := container.NewCenter()
	h2.Add(stopButton)
	c.Add(h)
	c.Add(h2)
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
	//buildMenu(w)
	//buildWindow(w)
	w.Resize(fyne.Size{Width: 1024, Height: 768})
	return w
}

//func buildWindow(w fyne.Window) error {
//	w.Resize(fyne.Size{Width: 1024, Height: 768})
//	//Navigate(w, LoginPage)
//	return nil
//}

func buildSystemTray(w fyne.Window) *fyne.Menu {
	tray := fyne.NewMenu("Hello",
		fyne.NewMenuItem("open window", func() {
			log.Println("Tapped show")
			w.Show()
		}),
	)
	return tray
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
	file, err := os.ReadFile(os.TempDir() + "/cookie.timetrace")
	if err != nil {
		return cookie, err
	}
	if err := json.Unmarshal(file, &cookie); err != nil {
		return cookie, err
	}
	return cookie, nil
}

func GetStatus() (models.StatusResponse, error) {
	data := models.StatusResponse{}
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return data, errors.New("cookie not set")
	}
	slog.Info("fetching current status")
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/projects/status", nil)
	if err != nil {
		loggedIn = false
		slog.Error("http request", "error", err)
		return data, err
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		loggedIn = false
		slog.Error("response", "error", err)
		return data, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		loggedIn = false
		slog.Error("status code", "status", response.Status, "code", response.StatusCode)
		return data, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		loggedIn = false
		return data, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		loggedIn = false
		return data, err
	}
	slog.Info("status response", "code", response.Status, "data", data)
	loggedIn = true
	return data, nil
}

func stop() {
	cookie, err := getCookie()
	if err != nil {
		slog.Error("cookie retrieval", "error", err)
		return
	}
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/projects/stop", nil)
	if err != nil {
		slog.Error("http request", "error", err)
		return
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		slog.Error("response", "error", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		slog.Error("status code", "status", response.Status, "code", response.StatusCode)
		return
	}
}
