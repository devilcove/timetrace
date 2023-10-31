package pages

import (
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/devilcove/timetrace/assets"
	"github.com/devilcove/timetrace/standalone/database"
	"github.com/devilcove/timetrace/standalone/models"
	"github.com/google/uuid"
)

type record struct {
	id      uuid.UUID
	start   time.Time
	end     time.Time
	project string
}

func BuildResultsPage(w fyne.Window, r []models.Record) *fyne.Container {
	records := make(map[string][]record)
	for _, data := range r {
		records[data.Project] = append(records[data.Project], record{start: data.Start, end: data.End, id: data.ID, project: data.Project})
	}
	buildMenu(w)
	logo := canvas.NewImageFromResource(fyne.NewStaticResource("logo", assets.SmallLogo))
	logo.FillMode = canvas.ImageFillOriginal
	c := container.NewVBox(logo)
	for project, dates := range records {
		var duration time.Duration
		for _, date := range dates {
			duration = duration + date.end.Sub(date.start)
		}
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
			start := widget.NewEntry()
			start.Text = dates[id].start.Format(time.DateTime)
			end := widget.NewEntry()
			end.Text = dates[id].end.Format(time.DateTime)
			items := []*widget.FormItem{
				widget.NewFormItem("Start Time:", start),
				widget.NewFormItem("  End Time:", end),
			}

			d := dialog.NewForm("Edit Record", "Submit", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				if err := editRecord(dates[id], start.Text, end.Text); err != nil {
					dialog.ShowError(err, w)
				}
			}, w)
			d.Resize(fyne.Size{Width: 400})
			d.Show()
			c.Refresh()

		}
		padded := container.NewVScroll(list)
		padded.SetMinSize(fyne.Size{Height: 120})
		c.Add(widget.NewLabel(project + ": " + models.FmtDuration(duration)))
		c.Add(padded)
	}
	return c
}

func editRecord(record record, newStart, newEnd string) error {
	fmt.Println("editing record ", record, newStart, newEnd)
	start, err := time.Parse(time.DateTime, newStart)
	if err != nil {
		return errors.New("invalid start time")
	}
	fmt.Println("new start time", start)
	end, err := time.Parse(time.DateTime, newEnd)
	if err != nil {
		return errors.New("invalid end time")
	}
	fmt.Println("new end time", end)
	if end.Sub(start) < 0 {
		return errors.New("end time is earlier than start ttime")
	}
	newRecord := models.Record{
		ID:      record.id,
		Start:   start,
		End:     end,
		Project: record.project,
	}
	if err := database.SaveRecord(&newRecord); err != nil {
		return err
	}
	return nil
}
