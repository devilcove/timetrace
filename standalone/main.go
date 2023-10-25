package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2/app"
	"github.com/devilcove/timetrace/standalone/database"
	"github.com/devilcove/timetrace/standalone/models"
	"github.com/devilcove/timetrace/standalone/pages"
	"github.com/joho/godotenv"
)

//func init() {
//pages.GetStatus()
//}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("read env", "error", err)
	}
	verbosity, _ := strconv.Atoi(os.Getenv("VERBOSITY"))
	setLogging(verbosity)
	database.InitializeDatabase()
	p := database.GetActiveProject()
	if p != nil {
		models.TrackingActive(*p)
	}
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

func setLogging(verbosity int) *slog.Logger {
	level := &slog.LevelVar{}
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source, ok := a.Value.Any().(*slog.Source)
			if ok {
				source.File = filepath.Base(source.File)
				source.Function = filepath.Base(source.Function)
			}
		}
		return a
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace, Level: level}))
	slog.SetDefault(logger)
	switch verbosity {
	case 4:
		level.Set(slog.LevelDebug)
	case 3:
		level.Set(slog.LevelInfo)
	case 2:
		level.Set(slog.LevelWarn)
	default:
		level.Set(slog.LevelError)
	}
	return logger
}
