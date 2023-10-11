package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetraced/models"
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
	projects := getProjects()
	projectsMenu := fyne.NewMenu("Projects")
	projectsMenu.Items = make([]*fyne.MenuItem, len(projects))
	//data := binding.BindStringList(&[]string{})
	//list := widget.NewListWithData(data,
	//	func() fyne.CanvasObject{
	//		return widget.NewLabel("template")
	//	},
	//	func(i binding.DataItem, o fyne.CanvasObject){
	//		o.(*widget.Label).Bind(i.(binding.String))
	//	})
	for i, project := range projects {
		project := project
		projectsMenu.Items[i] = fyne.NewMenuItem(project.Name, func() {
			slog.Info("selected project", "project", project.Name)
			start(project.Name)
			w.SetContent(BuildMainPage(w))
		})
	}
	projectsMenu.Items = append(projectsMenu.Items, fyne.NewMenuItem("Add Project", func() {
		project := widget.NewEntry()
		project.Resize(fyne.Size{Width: 800})
		items := []*widget.FormItem{
			widget.NewFormItem("Project to Add", project),
		}
		d := dialog.NewForm("Add Project", "Add  ", "Cancel", items, func(b bool) {

			//})
			//dialog.ShowForm("Add Project", "Add  ", "Cancel", items, func(b bool) {
			if !b {
				return
			}
			if err := addProject(project.Text); err != nil {
				dialog.ShowError(err, w)
			}
		}, w)
		d.Resize(fyne.Size{Width: 400})
		d.Show()
	}))

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
