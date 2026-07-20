package xts

import (
	"time"
)

// From 返回时间 t 对应的 Timestamp（秒级，亚秒部分被截断）。
func From(t time.Time) Timestamp { return Timestamp(t.Unix()) }

// Now 返回当前时间对应的 Timestamp。
// 若已通过 WithNowFunc 注入自定义时间函数（如 Tasklet 执行时间），则返回注入函数的结果。
func Now() Timestamp { return From(now()) }

func (ts Timestamp) time() time.Time { return time.Unix(int64(ts), 0) }

// Time 返回该 Timestamp 在全局时区下的本地时间。
func (ts Timestamp) Time() time.Time { return ts.time().In(location()) }

// TimeUTC 返回该 Timestamp 对应的 UTC 时间。
func (ts Timestamp) TimeUTC() time.Time { return ts.time().UTC() }

// Unix 返回原始的 Unix 秒数。
func (ts Timestamp) Unix() int64 { return int64(ts) }

// Add 增加时长，亚秒部分会被截断（Timestamp 精度为秒）。
func (ts Timestamp) Add(d time.Duration) Timestamp { return ts.AddS(Second(d / time.Second)) }

// AddS 增加 s 秒。
func (ts Timestamp) AddS(s Second) Timestamp { return ts + Timestamp(s) }

// Sub 返回两个时间戳相减的时间间隔。
func (ts Timestamp) Sub(other Timestamp) time.Duration {
	return time.Duration(ts.SubS(other)) * time.Second
}

// SubS 返回两个时间戳相减的秒数。
func (ts Timestamp) SubS(other Timestamp) Second { return Second(ts - other) }

// BeginningOfDay 返回该时间戳所在“当地日”的 0 点时间戳，当地日按全局固定时区偏移划分。
func (ts Timestamp) BeginningOfDay() Timestamp {
	return ts.AddS((-Second(ts) - zoneOffsetHours()*Hour) % Day)
}

// BeginningOfNextDay 返回下一“当地日”的 0 点时间戳。
func (ts Timestamp) BeginningOfNextDay() Timestamp { return ts.BeginningOfDay().AddS(Day) }

// InSameDay 判断两个时间戳是否落在同一“当地日”。
func (ts Timestamp) InSameDay(o Timestamp) bool {
	return ts.BeginningOfDay() == o.BeginningOfDay()
}

// DiffDays 返回两个时间戳相差的天数（按当地日起点相减）。
func (ts Timestamp) DiffDays(o Timestamp) int64 {
	return int64(ts.BeginningOfDay()-o.BeginningOfDay()) / Day
}

// WithinToday 判断该时间戳是否落在当前“当地日”。
func (ts Timestamp) WithinToday() bool { return ts.InSameDay(Now()) }

// InSameWeek 判断两个时间戳是否落在同一 ISO 周（基于全局时区的本地时间）。
func (ts Timestamp) InSameWeek(o Timestamp) bool {
	year1, week1 := ts.Time().ISOWeek()
	year2, week2 := o.Time().ISOWeek()
	return year1 == year2 && week1 == week2
}

// InSameMonth 判断两个时间戳是否落在同一自然月（基于全局时区的本地时间）。
func (ts Timestamp) InSameMonth(o Timestamp) bool {
	t1 := ts.Time()
	t2 := o.Time()
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// BeginningOfDayOfSameWeek 返回同一周内指定星期的 0 点时间戳。
// mondayAsWeekBeginning 为 true 时以周一为一周起始，周日被视为该周最后一天。
func (ts Timestamp) BeginningOfDayOfSameWeek(weekday time.Weekday, mondayAsWeekBeginning bool) Timestamp {
	t := ts.BeginningOfDay().Time()

	var dayDiff int
	currWeekDay := t.Weekday()

	if mondayAsWeekBeginning {
		if weekday == time.Sunday {
			weekday = SundayInMondayAsWeekBeginning
		}

		if currWeekDay == time.Sunday {
			currWeekDay = SundayInMondayAsWeekBeginning
		}
	}

	dayDiff = int(weekday - currWeekDay)
	return From(t.AddDate(0, 0, dayDiff))
}

// Week0 返回本周第一天的 0 点时刻。
func Week0(mondayAsWeekBeginning bool) Timestamp {
	weekday := time.Sunday
	if mondayAsWeekBeginning {
		weekday = time.Monday
	}
	return Now().BeginningOfDayOfSameWeek(weekday, mondayAsWeekBeginning)
}

// NextWeek0 返回下周第一天的 0 点时刻。
func NextWeek0(mondayAsWeekBeginning bool) Timestamp {
	return Week0(mondayAsWeekBeginning).AddS(Week)
}

// IsSameDay 判断两个时间戳是否落在同一“当地日”。
func IsSameDay(lhs, rhs Timestamp) bool { return lhs.InSameDay(rhs) }

// IsWithinToday 判断时间戳是否落在当前“当地日”。
func IsWithinToday(ts Timestamp) bool { return ts.WithinToday() }
