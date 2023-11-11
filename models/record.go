package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Record contains the data for an individual tracking event
type Record struct {
	ID      uuid.UUID
	Project string
	Start   time.Time
	End     time.Time
}

// Durations is a map of project durations
type Durations map[string]string

// Duration holds information about time tracked for a project
type Duration struct {
	Project string
	Elapsed string
}

// Duration returns the duration of an individual tracking event
func (r *Record) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

// FmtDuration formats a duration into hours,minutes and decimal hours for display
func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	dm := int(m) / 6
	if m%6 != 0 {
		dm += 1
	}
	dh := h
	if dm == 10 {
		dh += 1
		dm = 0
	}
	return fmt.Sprintf("%02d:%02d (%2d.%d Hours)", h, m, dh, dm)
}

// StatusResponse contains a days tracking information for display
type StatusResponse struct {
	Current      string
	Elapsed      string
	CurrentTotal string
	DailyTotal   string
	Durations    []Duration
}

// Status contains a days tracking information
type Status struct {
	Current    string
	Elapsed    time.Duration
	Total      time.Duration
	DailyTotal time.Duration
}
