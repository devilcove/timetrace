package pages

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/assets"
	"github.com/devilcove/timetrace/standalone/models"
)

type record struct {
	start time.Time
	end   time.Time
}

func BuildResultsPage(w fyne.Window, r []models.Record) *fyne.Container {
	records := make(map[string][]record)
	for _, data := range r {
		records[data.Project] = append(records[data.Project], record{start: data.Start, end: data.End})
	}
	buildMenu(w)
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("logo", assets.SmallLogo))
	logo.FillMode = canvas.ImageFillOriginal
	labelBox := container.NewVBox()
	for k, v := range records {
		labelBox.Add(widget.NewLabel(k))
		dateBox := container.NewVBox()
		for _, date := range v {
			startStop := fmt.Sprint(date.start.Format("2006-01-02:03:04"), " ", date.end.Format("2006-01-02:03:04"))
			dateBox.Add(widget.NewButton(startStop, func() {}))
		}
		labelBox.Add(dateBox)

	}
	c := container.NewVBox(logo, labelBox)
	return c
}
