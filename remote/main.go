package main

import (
	"log/slog"
	"time"

	"fyne.io/fyne/v2/app"
	"github.com/devilcove/timetrace/remote/pages"
)

func init() {
	pages.GetStatus()
}

func main() {
	a := app.New()
	w := pages.GetMainWindow(a, "TimeTrace")
	//w := a.NewWindow("Timetrace")
	w.SetContent(pages.BuildMainPage(w))
	go func() {
		for range time.Tick(time.Minute) {
			slog.Info("refreshing page")
			w.SetContent(pages.BuildMainPage(w))
		}
	}()

	w.ShowAndRun()
}
