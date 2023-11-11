package models

import "time"

// Report contains data tracking information for a project
type Report struct {
	Project string
	Total   time.Duration
	Items   []ReportRecord
}

// ReportRecord contians start/end times for an individual tracking event
type ReportRecord struct {
	Start time.Time
	End   time.Time
}

// ReportRequest contains data to retrieve a report
type ReportRequest struct {
	Start    time.Time
	End      time.Time
	Projects []string
}
