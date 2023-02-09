package xtime

import (
	"time"
)

// WeekStartDay set week start day, default is sunday
var WeekStartDay = time.Sunday

// ToLocal to local time
func ToLocal(tt time.Time, offset int) time.Time {
	return tt.In(time.FixedZone("", offset))
}

// BeginningOfMinute beginning of minute
func BeginningOfMinute(tt time.Time) time.Time {
	return tt.Truncate(time.Minute)
}

// BeginningOfHour beginning of hour
func BeginningOfHour(tt time.Time) time.Time {
	return BeginningOfHourSpec(tt, tt.Hour())
}

// BeginningOfHourSpec beginning of specified hour
func BeginningOfHourSpec(tt time.Time, hourSpec int) time.Time {
	y, m, d := tt.Date()
	return time.Date(y, m, d, hourSpec, 0, 0, 0, tt.Location())
}

// BeginningOfDay beginning of day
func BeginningOfDay(tt time.Time) time.Time {
	y, m, d := tt.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, tt.Location())
}

// BeginningOfWeek beginning of week
func BeginningOfWeek(tt time.Time) time.Time {
	t := BeginningOfDay(tt)
	weekday := int(t.Weekday())

	if WeekStartDay != time.Sunday {
		weekStartDayInt := int(WeekStartDay)
		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
	}
	return t.AddDate(0, 0, -weekday)
}

// BeginningOfMonth beginning of month
func BeginningOfMonth(tt time.Time) time.Time {
	y, m, _ := tt.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, tt.Location())
}

// BeginningOfYear beginning of year
func BeginningOfYear(tt time.Time) time.Time {
	y, _, _ := tt.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, tt.Location())
}

// EndOfMinute end of minute
func EndOfMinute(tt time.Time) time.Time {
	return BeginningOfMinute(tt).Add(time.Minute - time.Nanosecond)
}

// EndOfHour end of hour
func EndOfHour(tt time.Time) time.Time {
	return BeginningOfHour(tt).Add(time.Hour - time.Nanosecond)
}

// EndOfDay end of day
func EndOfDay(tt time.Time) time.Time {
	y, m, d := tt.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), tt.Location())
}

// EndOfWeek end of week
func EndOfWeek(tt time.Time) time.Time {
	return BeginningOfWeek(tt).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// EndOfMonth end of month
func EndOfMonth(tt time.Time) time.Time {
	return BeginningOfMonth(tt).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfYear end of year
func EndOfYear(tt time.Time) time.Time {
	return BeginningOfYear(tt).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// Monday the monday of tt
func Monday(tt time.Time) time.Time {
	t := BeginningOfDay(tt)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -weekday+1)
}

// Sunday the sunday of tt
func Sunday(tt time.Time) time.Time {
	t := BeginningOfDay(tt)
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	}
	return t.AddDate(0, 0, 7-weekday)
}

// DaysInMonth days of current month
func DaysInMonth(tt time.Time) int {
	switch tt.Month() {
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if IsLeapYear(tt.Year()) {
			return 29
		}
		return 28
	default:
		return 31
	}
}

// LeapYear year is leap year
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsToday tt is today
func IsToday(ts int64, tt time.Time) bool {
	y1, m1, d1 := time.Unix(ts, 0).In(tt.Location()).Date()
	y2, m2, d2 := tt.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsTomorrow tt is tomorrow
func IsTomorrow(ts int64, tt time.Time) bool {
	return IsToday(ts, tt.AddDate(0, 0, 1))
}

// IsYesterday tt is yesterday
func IsYesterday(ts int64, tt time.Time) bool {
	return IsToday(ts, tt.AddDate(0, 0, -1))
}

// IsThisWeek tt is current week
func IsThisWeek(ts int64, tt time.Time) (isThisWeek bool) {
	y1, w1 := tt.ISOWeek()
	y2, w2 := time.Unix(ts, 0).In(tt.Location()).ISOWeek()
	return y1 == y2 && w1 == w2
}

// SubDays days of between start and end
func SubDays(start, end int64) int64 {
	if start == end {
		return 0
	}
	return int64(time.Unix(end, 0).Sub(time.Unix(start, 0)).Hours() / 24)
}
