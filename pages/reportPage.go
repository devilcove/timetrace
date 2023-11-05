package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetraced/models"
)

var (
	globalAppTime = time.Now()
)

func BuildReportPage(w fyne.Window) *fyne.Container {
	buildMenu(w)
	// start time
	startTime := time.Now()
	startStr := binding.NewString()
	startStr.Set(time.Now().Format("2006-01-02"))
	start := widget.NewEntryWithData(startStr)
	//start.Size(fyne.Size{Width: 100})
	startCal := newCalendar(time.Now(), func(t time.Time) {
		fmt.Println("time selected", t)
		startStr.Set(t.Format("2006-01-02"))
		startTime = t
	})
	var startButton *widget.Button
	startButton = widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(startButton)
		position.Y += startButton.Size().Height
		startCal.showAtPos(w.Canvas(), position)
	})
	startBox := container.NewVBox(start, startButton)
	startBox.Resize(fyne.Size{Width: 100})
	//endTime
	endTime := time.Now()
	endStr := binding.NewString()
	endStr.Set(time.Now().Format("2006-01-02"))
	end := widget.NewEntryWithData(endStr)
	endCal := newCalendar(time.Now(), func(t time.Time) {
		fmt.Println("time selected", t)
		endStr.Set(t.Format("2006-01-02"))
		endTime = t
	})
	var endButton *widget.Button
	endButton = widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(endButton)
		position.Y += endButton.Size().Height
		endCal.showAtPos(w.Canvas(), position)
	})
	endBox := container.NewVBox(end, endButton)

	//projects
	projects := getProjects()
	projectOptions := []string{}
	for _, project := range projects {
		project := project
		projectOptions = append(projectOptions, project.Name)
	}
	projectsCheckGroup := widget.NewCheckGroup(projectOptions, func(s []string) {})
	all := widget.NewCheck("Select All Projects", func(b bool) {
		if b {
			projectsCheckGroup.SetSelected(projectsCheckGroup.Options)
		} else {
			projectsCheckGroup.SetSelected([]string{})
		}
	})
	projectBox := container.NewVBox(all, projectsCheckGroup)

	//users
	userOptions := []string{}
	for _, user := range getUsers() {
		user := user
		userOptions = append(userOptions, user.Username)
	}
	usersCheckBoxGroup := widget.NewCheckGroup(userOptions, func(s []string) {})
	usersCheckBoxGroup.SetSelected([]string{currentUser.Username})
	all = widget.NewCheck("Select All Users", func(b bool) {
		if b {
			usersCheckBoxGroup.SetSelected(userOptions)
		} else {
			usersCheckBoxGroup.SetSelected([]string{})
		}
	})
	userBox := container.NewVBox(all, usersCheckBoxGroup)
	//Submit
	button := widget.NewButton("Submit", func() {
		cookie, err := getCookie()
		if err != nil {
			loggedIn = false
			return
		}
		client := &http.Client{Timeout: time.Second * 10}
		data := models.ReportRequest{
			Start:    startTime,
			End:      endTime,
			Projects: projectsCheckGroup.Selected,
			Users:    usersCheckBoxGroup.Selected,
		}
		payload, err := json.Marshal(data)
		if err != nil {
			slog.Error("json error", "error", err)
		}
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/reports", bytes.NewBuffer(payload))
		if err != nil {
			return
		}
		req.AddCookie(&cookie)
		response, err := client.Do(req)
		if err != nil {
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return
		}
		reports := []models.Record{}
		if err := json.Unmarshal(body, &reports); err != nil {
			return
		}
		if err := saveCookie(response.Cookies()); err != nil {
			return
		}
		w.SetContent(BuildResultsPage(w, reports))

	})
	//layout
	vBox := &fyne.Container{}
	if currentUser.IsAdmin {
		vBox = container.NewVBox(startBox, endBox, projectBox, userBox, button)
	} else {
		vBox = container.NewVBox(startBox, endBox, projectBox, button)

	}
	vBox.Resize(fyne.Size{Width: 800})
	w.SetContent(vBox)
	return vBox
}
