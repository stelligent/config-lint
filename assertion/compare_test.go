package assertion

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDaysOldForToday(t *testing.T) {
	now := time.Now().Format("2006-01-02T15:04:05Z")
	assert.Equal(t, 0, daysOld(now), "Expecting daysOld to return 0")
}

func TestDaysOldFor90DaysAgo(t *testing.T) {
	then := time.Now().Add(-time.Duration(90) * time.Hour * 24).Format("2006-01-02T15:04:05Z")
	assert.Equal(t, 90, daysOld(then), "Expecting daysOld to return 90")
}
