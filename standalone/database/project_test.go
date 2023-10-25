package database

import (
	"testing"

	"github.com/devilcove/timetrace/standalone/models"
	"github.com/stretchr/testify/assert"
)

func TestSaveProject(t *testing.T) {
	p := models.Project{
		Name:   "testProject",
		Active: true,
	}
	err := SaveProject(&p)
	assert.Nil(t, err)
}
