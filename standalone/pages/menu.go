package pages

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/standalone/database"
	"github.com/devilcove/timetrace/standalone/models"
	"github.com/google/uuid"
)

func buildMenu(w fyne.Window) error {
	//File
	status := fyne.NewMenuItem("Status", func() {
		w.SetContent(BuildMainPage(w))
	})
	quit := fyne.NewMenuItem("Quit", func() {
		w.Close()
	})
	fileMenu := fyne.NewMenu("File", status, quit)
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
		w.SetContent(BuildMainPage(w))
	}))

	//Reports
	reportButton := fyne.NewMenuItem("report", func() {
		w.SetContent(BuildReportPage(w))
	})
	reportsMenu := fyne.NewMenu("Reports", reportButton)

	// About
	helpMenuItem := fyne.NewMenuItem("Help", func() {
		dialog.ShowInformation("About", "v0.1.0", w)
	})
	aboutMenu := fyne.NewMenu("About", helpMenuItem)
	//MENU
	menu := fyne.NewMainMenu(fileMenu, projectsMenu, reportsMenu, aboutMenu)
	w.SetMainMenu(menu)
	return nil
}

func getProjects() (projects []models.Project) {
	projects, err := database.GetAllProjects()
	if err != nil {
		slog.Error("retrieve projects", "error", err)
	}
	return projects
}

func addProject(p string) error {
	if regexp.MustCompile(`\s+`).MatchString(p) {
		return errors.New("invalid project name")
	}
	if _, err := database.GetProject(p); err != nil && err != database.ErrNoResults {
		return fmt.Errorf("project exists %w", err)
	}
	project := models.Project{
		Name:    p,
		Active:  true,
		Updated: time.Now(),
	}
	if err := database.SaveProject(&project); err != nil {
		return fmt.Errorf("add new project %w", err)
	}
	return nil
}

func start(p string) error {
	project, err := database.GetProject(p)
	if err != nil {
		return fmt.Errorf("tracking start for project %s, %w", p, err)
	}
	if !project.Active {
		return errors.New("project is not active")
	}
	if models.IsTrackingActive() {
		if err := stop(); err != nil {
			slog.Error("stop tracking project", "error", err)
			return err
		}
	}
	record := models.Record{
		ID:      uuid.New(),
		Project: p,
		Start:   time.Now(),
	}
	if err := database.SaveRecord(&record); err != nil {
		slog.Error("save record", "error", err)
		return err
	}
	models.TrackingActive(project)
	slog.Info("tracking started", "project", p)
	return nil
}
