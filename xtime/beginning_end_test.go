package xtime

import (
	"testing"
	"time"
)

var (
	format          = "2006-01-02 15:04:05.999999999"
	locationCaracas *time.Location
	locationBerlin  *time.Location
	timeCaracas     time.Time
)

func init() {
	var err error
	if locationCaracas, err = time.LoadLocation("America/Caracas"); err != nil {
		panic(err)
	}

	if locationBerlin, err = time.LoadLocation("Europe/Berlin"); err != nil {
		panic(err)
	}

	timeCaracas = time.Date(2016, 1, 1, 12, 10, 0, 0, locationCaracas)
}

func assertT(t *testing.T) func(time.Time, string, string) {
	return func(actual time.Time, expected string, msg string) {
		actualStr := actual.Format(format)
		if actualStr != expected {
			t.Errorf("Failed %s: actual: %v, expected: %v", msg, actualStr, expected)
		}
	}
}

func TestBeginningOf(t *testing.T) {
	assert := assertT(t)

	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)
	assert(BeginningOfMinute(n), "2013-11-18 17:51:00", "BeginningOfMinute")

	WeekStartDay = time.Monday
	assert(BeginningOfWeek(n), "2013-11-18 00:00:00", "BeginningOfWeek, FirstDayMonday")

	WeekStartDay = time.Tuesday
	assert(BeginningOfWeek(n), "2013-11-12 00:00:00", "BeginningOfWeek, FirstDayTuesday")

	WeekStartDay = time.Wednesday
	assert(BeginningOfWeek(n), "2013-11-13 00:00:00", "BeginningOfWeek, FirstDayWednesday")

	WeekStartDay = time.Thursday
	assert(BeginningOfWeek(n), "2013-11-14 00:00:00", "BeginningOfWeek, FirstDayThursday")

	WeekStartDay = time.Friday
	assert(BeginningOfWeek(n), "2013-11-15 00:00:00", "BeginningOfWeek, FirstDayFriday")

	WeekStartDay = time.Saturday
	assert(BeginningOfWeek(n), "2013-11-16 00:00:00", "BeginningOfWeek, FirstDaySaturday")

	WeekStartDay = time.Sunday
	assert(BeginningOfWeek(n), "2013-11-17 00:00:00", "BeginningOfWeek, FirstDaySunday")

	assert(BeginningOfHour(n), "2013-11-18 17:00:00", "BeginningOfHour")

	assert(BeginningOfHour(timeCaracas), "2016-01-01 12:00:00", "BeginningOfHour Caracas")

	assert(BeginningOfDay(n), "2013-11-18 00:00:00", "BeginningOfDay")

	location, err := time.LoadLocation("Japan")
	if err != nil {
		t.Fatalf("Error loading location: %v", err)
	}
	beginningOfDay := time.Date(2015, 05, 01, 0, 0, 0, 0, location)
	assert(BeginningOfDay(beginningOfDay), "2015-05-01 00:00:00", "BeginningOfDay")

	// DST
	dstBeginningOfDay := time.Date(2017, 10, 29, 10, 0, 0, 0, locationBerlin)
	assert(BeginningOfDay(dstBeginningOfDay), "2017-10-29 00:00:00", "BeginningOfDay DST")

	assert(BeginningOfWeek(n), "2013-11-17 00:00:00", "BeginningOfWeek")

	dstBegginingOfWeek := time.Date(2017, 10, 30, 12, 0, 0, 0, locationBerlin)
	assert(BeginningOfWeek(dstBegginingOfWeek), "2017-10-29 00:00:00", "BeginningOfWeek")

	dstBegginingOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(BeginningOfWeek(dstBegginingOfWeek), "2017-10-29 00:00:00", "BeginningOfWeek")

	WeekStartDay = time.Monday
	assert(BeginningOfWeek(n), "2013-11-18 00:00:00", "BeginningOfWeek, FirstDayMonday")
	dstBegginingOfWeek = time.Date(2017, 10, 24, 12, 0, 0, 0, locationBerlin)
	assert(BeginningOfWeek(dstBegginingOfWeek), "2017-10-23 00:00:00", "BeginningOfWeek, FirstDayMonday")

	dstBegginingOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(BeginningOfWeek(dstBegginingOfWeek), "2017-10-23 00:00:00", "BeginningOfWeek, FirstDayMonday")

	WeekStartDay = time.Sunday

	assert(BeginningOfMonth(n), "2013-11-01 00:00:00", "BeginningOfMonth")

	// DST
	dstBeginningOfMonth := time.Date(2017, 10, 31, 0, 0, 0, 0, locationBerlin)
	assert(BeginningOfMonth(dstBeginningOfMonth), "2017-10-01 00:00:00", "BeginningOfMonth DST")

	assert(BeginningOfYear(timeCaracas), "2016-01-01 00:00:00", "BeginningOfYear Caracas")
}

func TestEndOf(t *testing.T) {
	assert := assertT(t)

	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)

	assert(EndOfMinute(n), "2013-11-18 17:51:59.999999999", "EndOfMinute")

	assert(EndOfHour(n), "2013-11-18 17:59:59.999999999", "EndOfHour")

	assert(EndOfHour(timeCaracas), "2016-01-01 12:59:59.999999999", "EndOfHour Caracas")

	assert(EndOfDay(n), "2013-11-18 23:59:59.999999999", "EndOfDay")

	dstEndOfDay := time.Date(2017, 10, 29, 1, 0, 0, 0, locationBerlin)
	assert(EndOfDay(dstEndOfDay), "2017-10-29 23:59:59.999999999", "EndOfDay DST")

	WeekStartDay = time.Tuesday
	assert(EndOfWeek(n), "2013-11-18 23:59:59.999999999", "EndOfWeek, FirstDayTuesday")

	WeekStartDay = time.Wednesday
	assert(EndOfWeek(n), "2013-11-19 23:59:59.999999999", "EndOfWeek, FirstDayWednesday")

	WeekStartDay = time.Thursday
	assert(EndOfWeek(n), "2013-11-20 23:59:59.999999999", "EndOfWeek, FirstDayThursday")

	WeekStartDay = time.Friday
	assert(EndOfWeek(n), "2013-11-21 23:59:59.999999999", "EndOfWeek, FirstDayFriday")

	WeekStartDay = time.Saturday
	assert(EndOfWeek(n), "2013-11-22 23:59:59.999999999", "EndOfWeek, FirstDaySaturday")

	WeekStartDay = time.Sunday
	assert(EndOfWeek(n), "2013-11-23 23:59:59.999999999", "EndOfWeek, FirstDaySunday")

	WeekStartDay = time.Monday
	assert(EndOfWeek(n), "2013-11-24 23:59:59.999999999", "EndOfWeek, FirstDayMonday")

	dstEndOfWeek := time.Date(2017, 10, 24, 12, 0, 0, 0, locationBerlin)
	assert(EndOfWeek(dstEndOfWeek), "2017-10-29 23:59:59.999999999", "EndOfWeek, FirstDayMonday")

	dstEndOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(EndOfWeek(dstEndOfWeek), "2017-10-29 23:59:59.999999999", "EndOfWeek, FirstDayMonday")

	WeekStartDay = time.Sunday
	assert(EndOfWeek(n), "2013-11-23 23:59:59.999999999", "EndOfWeek")

	dstEndOfWeek = time.Date(2017, 10, 29, 0, 0, 0, 0, locationBerlin)
	assert(EndOfWeek(dstEndOfWeek), "2017-11-04 23:59:59.999999999", "EndOfWeek")

	dstEndOfWeek = time.Date(2017, 10, 29, 12, 0, 0, 0, locationBerlin)
	assert(EndOfWeek(dstEndOfWeek), "2017-11-04 23:59:59.999999999", "EndOfWeek")

	assert(EndOfMonth(n), "2013-11-30 23:59:59.999999999", "EndOfMonth")

	assert(EndOfYear(n), "2013-12-31 23:59:59.999999999", "EndOfYear")

	n1 := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.UTC)
	assert(EndOfMonth(n1), "2013-02-28 23:59:59.999999999", "EndOfMonth for 2013/02")

	n2 := time.Date(1900, 02, 18, 17, 51, 49, 123456789, time.UTC)
	assert(EndOfMonth(n2), "1900-02-28 23:59:59.999999999", "EndOfMonth")
}

func TestMondayAndSunday(t *testing.T) {
	assert := assertT(t)

	n := time.Date(2013, 11, 19, 17, 51, 49, 123456789, time.UTC)
	n2 := time.Date(2013, 11, 24, 17, 51, 49, 123456789, time.UTC)
	nDst := time.Date(2017, 10, 29, 10, 0, 0, 0, locationBerlin)

	assert(Monday(n), "2013-11-18 00:00:00", "Monday")

	assert(Monday(n2), "2013-11-18 00:00:00", "Monday")

	assert(Monday(timeCaracas), "2015-12-28 00:00:00", "Monday Caracas")

	assert(Monday(nDst), "2017-10-23 00:00:00", "Monday DST")

	assert(Sunday(n), "2013-11-24 00:00:00", "Sunday")

	assert(Sunday(n2), "2013-11-24 00:00:00", "Sunday")

	assert(Sunday(timeCaracas), "2016-01-03 00:00:00", "Sunday Caracas")

	assert(Sunday(nDst), "2017-10-29 00:00:00", "Sunday DST")

	assert(BeginningOfWeek(n), "2013-11-17 00:00:00", "BeginningOfWeek, FirstDayMonday")
	WeekStartDay = time.Monday
	assert(BeginningOfWeek(n), "2013-11-18 00:00:00", "BeginningOfWeek, FirstDayMonday")
}

func TestSameDay(t *testing.T) {
	n := time.Date(2013, 11, 19, 17, 51, 49, 123456789, time.UTC)
	if !IsToday(BeginningOfDay(n).Unix(), n) {
		t.Errorf("Failed IsToday: actual: false, expected: true")
	}
	if !IsTomorrow(EndOfDay(n).Unix()+10, n) {
		t.Errorf("Failed IsTomorrow: actual: false, expected: true")
	}
	if !IsYesterday(BeginningOfDay(n).Unix()-10, n) {
		t.Errorf("Failed IsYesterday: actual: false, expected: true")
	}
	if !IsThisWeek(BeginningOfDay(n).Unix(), n) {
		t.Errorf("Failed IsThisWeek: actual: false, expected: true")
	}
	{
		v := SubDays(n.Unix(), n.Unix()+10)
		if v != 0 {
			t.Errorf("Failed SubDays: actual: %d, expected: 0", v)
		}
	}
	{
		v := SubDays(n.Unix(), n.Unix()+24*60*60)
		if v != 1 {
			t.Errorf("Failed SubDays: actual: %d, expected: 1", v)
		}
	}
	{
		v := DaysInMonth(n)
		if DaysInMonth(n) != 30 {
			t.Errorf("Failed DaysInMonth: actual: %d, expected: 30", v)
		}
	}
}
