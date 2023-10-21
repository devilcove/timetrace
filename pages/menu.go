package pages

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	wid "fyne.io/x/fyne/widget"
	"github.com/devilcove/timetraced/models"
)

var loggedIn bool = true

func buildMenu(w fyne.Window) error {
	slog.Info("building menu", "loggedin", loggedIn)
	//File
	status := fyne.NewMenuItem("Status", func() {
		w.SetContent(BuildMainPage(w))
	})
	logout := fyne.NewMenuItem("Logout", func() {
		loggedIn = false
		//clear cookie file
		cookie, err := os.Create(os.TempDir() + "/cookie.timetrace")
		if err != nil {
			slog.Error("delete cookie store", "error", err)
		}
		cookie.Close()
		w.SetContent(BuildLoginPage(w))
	})
	quit := fyne.NewMenuItem("Quit", func() {
		w.Close()
	})
	fileMenu := fyne.NewMenu("File", status, logout, quit)
	//Projects
	projects := getProjects()
	projectsMenu := fyne.NewMenu("Projects")
	projectsMenu.Items = make([]*fyne.MenuItem, len(projects))
	for i, project := range projects {
		project := project
		projectsMenu.Items[i] = fyne.NewMenuItem(project.Name, func() {
			slog.Info("selected project", "project", project.Name)
			start(project.Name)
			w.SetContent(BuildMainPage(w))
		})
	}
	if currentUser.IsAdmin {
		projectsMenu.Items = append(projectsMenu.Items, fyne.NewMenuItem("Add Project", func() {
			project := widget.NewEntry()
			project.Resize(fyne.Size{Width: 800})
			items := []*widget.FormItem{
				widget.NewFormItem("Project to Add", project),
			}
			d := dialog.NewForm("Add Project", "Add  ", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				if err := addProject(project.Text); err != nil {
					dialog.ShowError(err, w)
				}
			}, w)
			d.Resize(fyne.Size{Width: 400})
			d.Show()
			buildMenu(w)
			w.Content().Refresh()
		}))
	}

	//Reports
	reportButton := fyne.NewMenuItem("report", func() {
		w.SetContent(BuildReportPage(w))
	})
	reportMenuItem := fyne.NewMenuItem("Generate Report", func() {
		var startTime, endTime time.Time
		start := wid.NewCalendar(time.Now(), func(t time.Time) {
			startTime = t
		})

		end := wid.NewCalendar(time.Now(), func(t time.Time) {
			endTime = t
		})
		items := []*widget.FormItem{
			widget.NewFormItem("StartDate", start),
			widget.NewFormItem("EndDate", end),
		}
		d := dialog.NewForm("Generate Report", "Submit", "Cancel", items, func(b bool) {
			if !b {
				return
			}
			getReport(startTime, endTime)
		}, w)
		//p := widget.NewModalPopUp(x, d)
		w.Resize(fyne.Size{Width: 400, Height: 800})
		d.Show()
		//w.SetContent(ReportPage(w))
		w.SetContent(BuildMainPage(w))
	})
	reportsMenu := fyne.NewMenu("Reports", reportButton, reportMenuItem)

	// Users
	usersMenu := &fyne.Menu{}
	if currentUser.IsAdmin {
		users := getUsers()
		usersMenu = fyne.NewMenu("Users")
		usersMenu.Items = make([]*fyne.MenuItem, len(users))
		for i, user := range users {
			user := user
			usersMenu.Items[i] = fyne.NewMenuItem("Edit "+user.Username, func() {
				password := widget.NewPasswordEntry()
				items := []*widget.FormItem{
					widget.NewFormItem("New Password", password),
				}
				d := dialog.NewForm("Edit User", "Submit", "Cancel", items, func(b bool) {
					if !b {
						return
					}
					user.Password = password.Text
					if err := editUser(user); err != nil {
						dialog.ShowError(err, w)
					}
				}, w)
				d.Resize(fyne.Size{Width: 400})
				d.Show()
				w.SetContent(BuildMainPage(w))
			})
		}
		usersMenu.Items = append(usersMenu.Items, fyne.NewMenuItem("Add User", func() {
			user := widget.NewEntry()
			password := widget.NewPasswordEntry()
			user.Resize(fyne.Size{Width: 800})
			items := []*widget.FormItem{
				widget.NewFormItem("User to Add", user),
				widget.NewFormItem("Password", password),
			}
			d := dialog.NewForm("Add User", "Add  ", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				err := addUser(models.User{Username: user.Text, Password: password.Text})
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("added user", user.Text, w)
				}

			}, w)
			d.Resize(fyne.Size{Width: 400})
			d.Show()
			w.SetContent(BuildMainPage(w))

		}))
		usersMenu.Items = append(usersMenu.Items, fyne.NewMenuItem("Delete User", func() {
			user := widget.NewEntry()
			user.Resize(fyne.Size{Width: 800})
			items := []*widget.FormItem{
				widget.NewFormItem("User to Delete", user),
			}
			d := dialog.NewForm("Delete User", "Delete", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				if err := deleteUser(user.Text); err != nil {
					dialog.ShowError(err, w)
				}
			}, w)
			d.Resize(fyne.Size{Width: 400})
			d.Show()
			w.SetContent(BuildMainPage(w))
		}))
	} else {
		usersMenu = fyne.NewMenu("Users", fyne.NewMenuItem("Edit Password", func() {
			user := currentUser
			password := widget.NewPasswordEntry()
			items := []*widget.FormItem{
				widget.NewFormItem("New Password", password),
			}
			d := dialog.NewForm("Edit User", "Submit", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				user.Password = password.Text
				if err := editUser(user); err != nil {
					dialog.ShowError(err, w)
				}
			}, w)
			d.Resize(fyne.Size{Width: 400})
			d.Show()
			w.SetContent(BuildMainPage(w))

		}))
	}

	// About
	helpMenuItem := fyne.NewMenuItem("Help", func() {
		dialog.ShowInformation("About", "v0.1.0", w)
	})
	aboutMenu := fyne.NewMenu("About", helpMenuItem)
	//MENU
	menu := fyne.NewMainMenu(fileMenu, projectsMenu, usersMenu, reportsMenu, aboutMenu)
	w.SetMainMenu(menu)
	return nil
}

func getProjects() (projects []models.Project) {
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return
	}
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/projects", nil)
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
	if err := json.Unmarshal(body, &projects); err != nil {
		return
	}
	return projects
}

func addProject(p string) error {
	project := models.Project{Name: p}
	payload, err := json.Marshal(project)
	if err != nil {
		return err
	}
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return err
	}
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/projects", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status error %s", response.Status)
	}
	fmt.Println("add projects", p)
	return nil
}

func start(p string) error {
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return err
	}
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/projects/"+p+"/start", nil)
	if err != nil {
		return err
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status error %s", response.Status)
	}
	fmt.Println("start recording projects", p)
	return nil
}

func getUsers() (users []models.User) {
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return
	}
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/users", nil)
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
	if err := json.Unmarshal(body, &users); err != nil {
		return
	}
	return
}

func editUser(u models.User) error {
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return errors.New("not logged in")
	}
	client := &http.Client{Timeout: time.Second * 10}
	payload, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshal %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/users", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("create request %w", err)
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request %s", response.Status)
	}
	return nil
}

func addUser(u models.User) error {
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return errors.New("not logged in")
	}
	client := &http.Client{Timeout: time.Second * 10}
	payload, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshal %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/users", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("create request %w", err)
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("bad request %s", response.Status)
	}
	slog.Info("added user", "user", u.Username)
	return nil
}

func deleteUser(u string) error {
	cookie, err := getCookie()
	if err != nil {
		loggedIn = false
		return errors.New("not logged in")
	}
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/users/"+u, nil)
	if err != nil {
		return fmt.Errorf("create request %w", err)
	}
	req.AddCookie(&cookie)
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("bad request %s", response.Status)
	}
	return nil
}

func getReport(s, e time.Time) {
	slog.Info("get report", "start", s, "end", e)
}
