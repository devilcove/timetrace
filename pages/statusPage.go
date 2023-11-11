package pages

import (
	"log/slog"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/assets"
	"github.com/devilcove/timetrace/database"
	"github.com/devilcove/timetrace/models"
)

var currentPage string

// SetCurrentPage saves the pages currently displayed
func SetCurrentPage(page string) {
	currentPage = page
}

// GetCurrentPage returns the page currently displayed
func GetCurrentPage() string {
	return currentPage
}

// StatusPage builds status page for display
func StatusPage(w fyne.Window) *fyne.Container {
	buildMenu(w)
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("logo", assets.SmallLogo))
	logo.FillMode = canvas.ImageFillOriginal
	status, err := GetStatus()
	if err != nil {
		slog.Error("get status", "error", err)
		os.Exit(1)
	}
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Current Project", Widget: widget.NewLabel(status.Current)},
			{Text: "Time This Session", Widget: widget.NewLabel(status.Elapsed)},
			{Text: "Time Today", Widget: widget.NewLabel(status.CurrentTotal)},
		},
	}
	stopButton := widget.NewButton("Stop    ", func() {
		stop()
		SetCurrentPage("status")
		w.SetContent(StatusPage(w))
	})
	stop := container.NewCenter(stopButton)
	todayTotals := widget.NewLabel("Total Time Today")
	todayTotals.Alignment = fyne.TextAlignCenter
	dailyForm := &widget.Form{}
	for _, duration := range status.Durations {
		dailyForm.Append(duration.Project, widget.NewLabel(duration.Elapsed))

	}
	dailyForm.Append("Total Today", widget.NewLabel(status.DailyTotal))
	c := container.NewCenter(container.NewVBox(logo, form, stop, todayTotals, dailyForm))
	return c
}

// GetMainWindow sets up the main window
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
	w.Resize(fyne.Size{Width: 512, Height: 256})
	return w
}

func buildSystemTray(w fyne.Window) *fyne.Menu {
	tray := fyne.NewMenu("Hello",
		fyne.NewMenuItem("Display", func() {
			slog.Info("Tapped show")
			w.Show()
		}),
	)
	return tray
}

// GetStatus retrieves the current day tracking events
func GetStatus() (models.StatusResponse, error) {
	durations := make(map[string]time.Duration)
	status := models.Status{}
	response := models.StatusResponse{}
	records, err := database.GetTodaysRecords()
	if err != nil {
		return response, err
	}
	status.Current = models.Tracked()
	for _, record := range records {
		if record.End.IsZero() {
			record.End = time.Now()
			status.Elapsed = record.Duration()
		}
		durations[record.Project] = durations[record.Project] + record.End.Sub(record.Start)
		status.DailyTotal = status.DailyTotal + record.Duration()
		if record.Project == status.Current {
			status.Total = status.Total + record.Duration()
		}
	}
	response.Current = status.Current
	response.Elapsed = models.FmtDuration(status.Elapsed)
	response.CurrentTotal = models.FmtDuration(status.Total)
	response.DailyTotal = models.FmtDuration(status.DailyTotal)
	for k := range durations {
		value := models.FmtDuration(durations[k])
		duration := models.Duration{
			Project: k,
			Elapsed: value,
		}
		response.Durations = append(response.Durations, duration)
	}
	return response, nil
}

func stop() error {
	records, err := database.GetAllRecords()
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.End.IsZero() {
			record.End = time.Now()
			if err := database.SaveRecord(&record); err != nil {
				return err
			}
		}
	}
	return nil
}
