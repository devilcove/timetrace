package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/devilcove/timetrace/pages"
)

func main() {
	a := app.New()
	w := pages.GetMainWindow(a, "TimeTrace")
	//w := a.NewWindow("Timetrace")
	w.SetContent(pages.BuildMainPage(w))

	w.ShowAndRun()
}
