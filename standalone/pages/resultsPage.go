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
	"github.com/google/uuid"
)

type record struct {
	id    uuid.UUID
	start time.Time
	end   time.Time
}

func BuildResultsPage(w fyne.Window, r []models.Record) *fyne.Container {
	records := make(map[string][]record)
	for _, data := range r {
		records[data.Project] = append(records[data.Project], record{start: data.Start, end: data.End, id: data.ID})
	}
	buildMenu(w)
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("logo", assets.SmallLogo))
	logo.FillMode = canvas.ImageFillOriginal
	c := container.NewVBox(logo)
	for project, dates := range records {
		list := widget.NewList(
			func() int {
				return len(dates)
			},
			func() fyne.CanvasObject {
				template := widget.NewLabel("placeholder")
				//template.Resize(fyne.Size{Width: 100})
				return template
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(fmt.Sprint(dates[i].start.Format(time.DateTime), "--", dates[i].end.Format(time.DateTime)))
			},
		)
		list.OnSelected = func(id widget.ListItemID) {
			fmt.Println("record selected", id)
			fmt.Println(dates[id].id, dates[id].start, dates[id].end)
		}
		padded := container.NewVScroll(list)
		padded.SetMinSize(fyne.Size{Height: 120})
		c.Add(widget.NewLabel(project))
		c.Add(padded)
	}
	return c
}
