package pages

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/assets"
	"github.com/devilcove/timetrace/standalone/database"
	"github.com/devilcove/timetrace/standalone/models"
)

type Route int

const (
	MainPage Route = iota
	LoginPage
)

type DisplayStatus struct {
	Current      string
	SessionTime  string
	CurrentTotal string
	Totals       []struct {
		Project string
		Total   string
	}
}

func BuildMainPage(w fyne.Window) *fyne.Container {
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("logo", assets.SmallLogo))
	logo.FillMode = canvas.ImageFillOriginal
	status, err := GetStatus()
	if err != nil {
		slog.Error("get status", "error", err)
		os.Exit(1)
	}
	buildMenu(w)
	text := widget.NewTextGrid()
	text.SetText(fmt.Sprintf("Current Project:\t%s\nTime This Session:\t%s\nTime Today:\t\t\t%s\n", status.Current, status.Elapsed, status.CurrentTotal))
	stopButton := widget.NewButton("Stop    ", func() {
		stop()
		w.SetContent(BuildMainPage(w))
	})
	todayTotals := widget.NewLabel("Total Time Today")
	todayTotals.Alignment = fyne.TextAlignCenter
	var durations string
	for _, duration := range status.Durations {
		durations = durations + "\n" + duration.Project + "\t\t"
		durations = durations + duration.Elapsed
	}
	text2 := widget.NewTextGrid()
	text2.SetText(durations)
	text3 := widget.NewTextGrid()
	text3.SetText(fmt.Sprintf("\nTotal\t\t\t%s", status.DailyTotal))
	c := container.NewVBox()
	session := container.NewCenter()
	session.Add(text)
	stop := container.NewCenter()
	stop.Add(stopButton)
	summary := container.NewCenter()
	summary.Add(text2)
	dailyTotal := container.NewCenter()
	dailyTotal.Add(text3)
	c.Add(logo)
	c.Add(session)
	c.Add(stop)
	c.Add(todayTotals)
	c.Add(summary)
	c.Add(dailyTotal)
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
	w.Resize(fyne.Size{Width: 512, Height: 240})
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
