package models

import (
	"time"
)

var (
	trackingActive bool
	trackedProject string
)

// Project is a project that can be tracked
type Project struct {
	Name    string
	Active  bool
	Updated time.Time
}

// StartRequest is a request to start tracking time for a project
type StartRequest struct {
	Project string
}

// IstracktingActive idicates whether tracking is currently active
func IsTrackingActive() bool {
	return trackingActive
}

// TrackingActive sets active project
func TrackingActive(p Project) {
	trackingActive = true
	trackedProject = p.Name
}

// TrackingInactive turns off project tracking
func TrackingInactive() {
	trackingActive = false
	trackedProject = ""
}

// Tracked returns the actively tracked project
func Tracked() string {
	return trackedProject
}
