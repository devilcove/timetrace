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

var currentPage = "status"

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("read env", "error", err)
	}
	verbosity, err := strconv.Atoi(os.Getenv("VERBOSITY"))
	if err != nil {
		verbosity = 3
	}
	setLogging(verbosity)
	database.InitializeDatabase()
	p := database.GetActiveProject()
	if p != nil {
		models.TrackingActive(*p)
	}
	a := app.New()
	w := pages.GetMainWindow(a, "TimeTrace")
	//w := a.NewWindow("Timetrace")
	pages.SetCurrentPage("status")
	w.SetContent(pages.StatusPage(w))
	go func() {
		for range time.Tick(time.Minute) {
			slog.Info("refreshing page")
			if pages.GetCurrentPage() == "status" {
				w.SetContent(pages.StatusPage(w))
			}
		}
	}()
	go func() {
		if os.Getenv("sync") != "true" {
			slog.Info("sync not set")
			return
		}
		for range time.Tick(time.Minute * 5) {
			slog.Info("syncing db")
			sync()
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

func sync() {
	file, err := os.ReadFile("./time.db")
	if err != nil {
		slog.Error("sync database", "error", err)
	}
	if err := os.WriteFile("./time.db.backup", file, os.ModePerm); err != nil {
		slog.Error("save sync", "error", err)
	}
}
