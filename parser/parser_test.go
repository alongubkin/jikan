package parser

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEveryMinute(t *testing.T) {
	schedule, err := parse("every minute")
	assert.NoError(t, err)
	assert.Equal(t, 1, schedule.Interval)
	assert.Equal(t, Minutes, schedule.TimeUnit)
	assert.Nil(t, schedule.Months)
}

func TestEvery5Seconds(t *testing.T) {
	schedule, err := parse("every 5 seconds")
	assert.NoError(t, err)
	assert.Equal(t, 5, schedule.Interval)
	assert.Equal(t, Seconds, schedule.TimeUnit)
	assert.Nil(t, schedule.Months)
}

func TestEvery2000Years(t *testing.T) {
	schedule, err := parse("every 2000 years")
	assert.NoError(t, err)
	assert.Equal(t, 2000, schedule.Interval)
	assert.Equal(t, Years, schedule.TimeUnit)
	assert.Nil(t, schedule.Months)
}

func TestEvery10HoursOfJanuary(t *testing.T) {
	schedule, err := parse("every 10 hours of january")
	assert.NoError(t, err)
	assert.Equal(t, 10, schedule.Interval)
	assert.Equal(t, Hours, schedule.TimeUnit)
	assert.EqualValues(t, []time.Month{1}, schedule.Months)
}

func TestEvery30HoursOfFebAprilMay(t *testing.T) {
	schedule, err := parse("every 30 hours of feb, april, may")
	assert.NoError(t, err)
	assert.Equal(t, 30, schedule.Interval)
	assert.Equal(t, Hours, schedule.TimeUnit)
	assert.EqualValues(t, []time.Month{2, 4, 5}, schedule.Months)
}

func TestEvery5SecondsOfMonth(t *testing.T) {
	schedule, err := parse("every 5 seconds of month")
	assert.NoError(t, err)
	assert.Equal(t, 5, schedule.Interval)
	assert.Equal(t, Seconds, schedule.TimeUnit)
	assert.Nil(t, schedule.Months)
}

func TestEveryMonthOfFebruary(t *testing.T) {
	_, err := parse("every month of february")
	assert.Error(t, err)
}

func TestEvery5YearOfFebJanMarch(t *testing.T) {
	_, err := parse("every 5 years of feb, jan, march")
	assert.Error(t, err)
}

func parse(s string) (*IntervalsSchedule, error) {
	parser := NewParser(strings.NewReader(s))
	schedule, err := parser.ParseIntervalsSchedule()
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
