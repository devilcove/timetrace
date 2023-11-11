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

// ReportPage builds a resort page for display
func ReportPage(w fyne.Window) *fyne.Container {
	buildMenu(w)
	// start time
	startTime := time.Now()
	startStr := binding.NewString()
	startStr.Set(time.Now().Format("2006-01-02"))
	startEntry := widget.NewEntryWithData(startStr)
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
	startText := widget.NewLabel("Start Time:")
	start := container.NewPadded(container.New(&datePicker{}, startText, startEntry, startButton))
	//endTime
	endTime := time.Now()
	endStr := binding.NewString()
	endStr.Set(time.Now().Format("2006-01-02"))
	endEntry := widget.NewEntryWithData(endStr)
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
	endText := widget.NewLabel("End Time:  ")
	end := container.NewPadded(container.New(&datePicker{}, endText, endEntry, endButton))

	//projects
	projectText := widget.NewLabel("Select Projects")
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
	projectBox := container.NewPadded(container.NewVBox(projectText, all, projectsCheckGroup))

	//Submit
	button := container.NewPadded(widget.NewButton("Submit", func() {

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
		w.SetContent(ResultsPage(w, reports))

	}))

	//layout
	vBox := &fyne.Container{}
	vBox = container.NewVBox(start, end, projectBox, button)
	w.SetContent(vBox)
	return vBox
}
