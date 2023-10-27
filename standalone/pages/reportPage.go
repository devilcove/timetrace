package pages

import (
	"fmt"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/standalone/database"
	"github.com/devilcove/timetrace/standalone/models"
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

	//Submit
	button := widget.NewButton("Submit", func() {

		data := models.ReportRequest{
			Start:    startTime,
			End:      endTime,
			Projects: projectsCheckGroup.Selected,
		}
		reports, err := database.GetReportRecords(data)
		if err != nil {
			slog.Error("get records", "error", err)
		}
		SetCurrentPage("results")
		w.SetContent(BuildResultsPage(w, reports))

	})

	//layout
	vBox := &fyne.Container{}
	vBox = container.NewVBox(startBox, endBox, projectBox, button)
	vBox.Resize(fyne.Size{Width: 800})
	w.SetContent(vBox)
	return vBox
}
