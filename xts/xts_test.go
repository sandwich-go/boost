package xts

import (
	"sync"
	"testing"
	"time"
)

// resetGlobal 把全局状态复位到默认，避免用例间相互污染。
// xts 是进程级单例，所有用例串行执行（不使用 t.Parallel）并在结束时复位。
func resetGlobal() {
	mu.Lock()
	provider = DefaultTimeProvider{}
	zoneOffset = 0
	loc = time.UTC
	nowFunc = nil
	mu.Unlock()
}

// makeTS 依据给定 UTC 偏移构造 Timestamp。
func makeTS(offsetHours int, y int, m time.Month, d, h, min, sec int) Timestamp {
	l := time.FixedZone("test", offsetHours*3600)
	return Timestamp(time.Date(y, m, d, h, min, sec, 0, l).Unix())
}

// setupFixedNow 初始化全局时钟为给定时区偏移，并把“当前时间”冻结为 now。
func setupFixedNow(offsetHours int64, now time.Time) {
	Initialize(WithTimeZoneOffset(offsetHours), WithNowFunc(func() time.Time { return now }))
}

// scaledProvider 是仅用于测试的可缩放时间提供者。
type scaledProvider struct {
	virtualBase time.Time
	actualBase  time.Time
	scale       float64
}

func (p scaledProvider) VirtualBase() time.Time { return p.virtualBase }
func (p scaledProvider) ActualBase() time.Time  { return p.actualBase }
func (p scaledProvider) Scale() float64         { return p.scale }

// ---- 常量 ----

func TestConstants(t *testing.T) {
	if Minute != 60 {
		t.Fatalf("Minute = %d, want 60", Minute)
	}
	if Hour != 3600 {
		t.Fatalf("Hour = %d, want 3600", Hour)
	}
	if Day != 86400 {
		t.Fatalf("Day = %d, want 86400", Day)
	}
	if Week != 7*86400 {
		t.Fatalf("Week = %d, want %d", Week, 7*86400)
	}
	if int(SundayInMondayAsWeekBeginning) != int(time.Saturday)+1 {
		t.Fatalf("SundayInMondayAsWeekBeginning = %d, want %d", int(SundayInMondayAsWeekBeginning), int(time.Saturday)+1)
	}
}

// ---- Initialize / Option ----

func TestInitializeDefault(t *testing.T) {
	defer resetGlobal()
	Initialize()
	if zoneOffset != 0 {
		t.Fatalf("default zoneOffset = %d, want 0", zoneOffset)
	}
	if _, ok := provider.(DefaultTimeProvider); !ok {
		t.Fatalf("default provider type = %T, want DefaultTimeProvider", provider)
	}
	if _, off := time.Now().In(Location()).Zone(); off != 0 {
		t.Fatalf("default Location offset = %d, want 0", off)
	}
}

func TestInitializeWithZoneOffset(t *testing.T) {
	defer resetGlobal()

	for _, off := range []int64{8, 0, -5, 14, -12} {
		Initialize(WithTimeZoneOffset(off))
		if zoneOffset != off {
			t.Fatalf("zoneOffset = %d, want %d", zoneOffset, off)
		}
		if _, got := time.Now().In(Location()).Zone(); int64(got) != off*3600 {
			t.Fatalf("Location offset = %d, want %d", got, off*3600)
		}
	}
}

func TestInitializeWithTimeProvider(t *testing.T) {
	defer resetGlobal()

	virtual := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	Initialize(WithTimeProvider(scaledProvider{
		virtualBase: virtual,
		actualBase:  time.Now().UTC(),
		scale:       1,
	}))
	// 未注入 nowFunc，Now 应基于 provider 的虚拟基准（2020 年）
	if y := Now().TimeUTC().Year(); y != 2020 {
		t.Fatalf("Now().Year() = %d, want 2020 (from provider virtual base)", y)
	}
}

func TestInitializePartialUpdatePreservesOthers(t *testing.T) {
	defer resetGlobal()

	p := scaledProvider{virtualBase: time.Now().UTC(), actualBase: time.Now().UTC(), scale: 1}
	Initialize(WithTimeProvider(p), WithTimeZoneOffset(8))

	// 仅更新 offset，provider 应保留
	Initialize(WithTimeZoneOffset(3))
	if zoneOffset != 3 {
		t.Fatalf("zoneOffset = %d, want 3", zoneOffset)
	}
	if _, ok := provider.(scaledProvider); !ok {
		t.Fatalf("provider not preserved across partial Initialize, got %T", provider)
	}

	// 仅更新 provider，offset 应保留
	Initialize(WithTimeProvider(DefaultTimeProvider{}))
	if zoneOffset != 3 {
		t.Fatalf("zoneOffset changed unexpectedly = %d, want 3", zoneOffset)
	}
	if _, ok := provider.(DefaultTimeProvider); !ok {
		t.Fatalf("provider = %T, want DefaultTimeProvider", provider)
	}
}

func TestInitializeNilProviderPanics(t *testing.T) {
	defer resetGlobal()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Initialize(WithTimeProvider(nil)) should panic")
		}
	}()
	Initialize(WithTimeProvider(nil))
}

// ---- Now / From / WithNowFunc 注入 ----

func TestFrom(t *testing.T) {
	t0 := time.Date(2025, 6, 15, 12, 30, 45, 999000000, time.UTC)
	if got := From(t0); got.Unix() != t0.Unix() {
		t.Fatalf("From(t).Unix() = %d, want %d (sub-second truncated)", got.Unix(), t0.Unix())
	}
}

func TestNowFallbackToProvider(t *testing.T) {
	defer resetGlobal()
	Initialize()
	// 未注入 nowFunc，默认走 DefaultTimeProvider（系统时钟）
	if delta := time.Since(Now().Time()); delta < -2*time.Second || delta > 2*time.Second {
		t.Fatalf("Now() default drifted from system now by %v", delta)
	}
}

func TestWithNowFuncInjectionAndReset(t *testing.T) {
	defer resetGlobal()

	fixed := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	Initialize(WithTimeZoneOffset(0), WithNowFunc(func() time.Time { return fixed }))

	if got := Now(); got != From(fixed) {
		t.Fatalf("Now() after WithNowFunc = %d, want %d", got, From(fixed))
	}
	// 冻结后多次调用返回一致
	if Now() != Now() {
		t.Fatal("injected Now() not stable")
	}

	// WithNowFunc(nil) 恢复默认（基于 provider），Now 应贴近系统时间
	Initialize(WithNowFunc(nil))
	if delta := time.Since(Now().Time()); delta < -2*time.Second || delta > 2*time.Second {
		t.Fatalf("Now() after reset drifted by %v", delta)
	}
}

// ---- Time / TimeUTC / Unix ----

func TestTimeConversion(t *testing.T) {
	defer resetGlobal()
	Initialize(WithTimeZoneOffset(8))

	// 2025-04-10 08:00:00 UTC = 2025-04-10 16:00:00 UTC+8
	ts := Timestamp(time.Date(2025, 4, 10, 8, 0, 0, 0, time.UTC).Unix())

	local := ts.Time()
	if local.Hour() != 16 || local.Day() != 10 {
		t.Fatalf("Time() = %v, want hour 16 day 10 (UTC+8)", local)
	}
	if _, off := local.Zone(); off != 8*3600 {
		t.Fatalf("Time() zone offset = %d, want %d", off, 8*3600)
	}

	utc := ts.TimeUTC()
	if utc.Hour() != 8 || utc.Day() != 10 {
		t.Fatalf("TimeUTC() = %v, want hour 8 day 10", utc)
	}

	if ts.Unix() != int64(ts) {
		t.Fatalf("Unix() = %d, want %d", ts.Unix(), int64(ts))
	}
}

// ---- Add / AddS / Sub / SubS ----

func TestAddSub(t *testing.T) {
	base := Timestamp(1000000)

	if base.AddS(100) != Timestamp(1000100) {
		t.Fatalf("AddS(100) = %d, want 1000100", base.AddS(100))
	}
	if base.AddS(-100) != Timestamp(999900) {
		t.Fatalf("AddS(-100) = %d, want 999900", base.AddS(-100))
	}
	if base.Add(2*time.Hour) != base.AddS(7200) {
		t.Fatal("Add(2h) != AddS(7200)")
	}
	if base.Add(1500*time.Millisecond) != base.AddS(1) {
		t.Fatal("Add should truncate sub-second")
	}
	if a := Timestamp(1000100); a.SubS(base) != 100 {
		t.Fatalf("SubS = %d, want 100", a.SubS(base))
	}
	if a := Timestamp(1000100); a.Sub(base) != 100*time.Second {
		t.Fatalf("Sub = %v, want 100s", a.Sub(base))
	}
	if base.Sub(Timestamp(1000100)) != -100*time.Second {
		t.Fatalf("Sub negative = %v, want -100s", base.Sub(Timestamp(1000100)))
	}
}

// ---- BeginningOfDay / BeginningOfNextDay ----

func TestBeginningOfDay(t *testing.T) {
	defer resetGlobal()

	cases := []struct {
		name        string
		offset      int64
		ts          Timestamp
		wantBOD     Timestamp
		wantNextDay Timestamp
	}{
		{"UTC+8 midnight stays", 8, makeTS(8, 2025, 4, 10, 0, 0, 0), makeTS(8, 2025, 4, 10, 0, 0, 0), makeTS(8, 2025, 4, 11, 0, 0, 0)},
		{"UTC+8 afternoon snaps back", 8, makeTS(8, 2025, 4, 10, 15, 30, 0), makeTS(8, 2025, 4, 10, 0, 0, 0), makeTS(8, 2025, 4, 11, 0, 0, 0)},
		{"UTC+8 23:59:59 same day", 8, makeTS(8, 2025, 4, 10, 23, 59, 59), makeTS(8, 2025, 4, 10, 0, 0, 0), makeTS(8, 2025, 4, 11, 0, 0, 0)},
		{"UTC+0", 0, makeTS(0, 2025, 4, 10, 15, 0, 0), makeTS(0, 2025, 4, 10, 0, 0, 0), makeTS(0, 2025, 4, 11, 0, 0, 0)},
		{"UTC-5", -5, makeTS(-5, 2025, 4, 10, 3, 0, 0), makeTS(-5, 2025, 4, 10, 0, 0, 0), makeTS(-5, 2025, 4, 11, 0, 0, 0)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Initialize(WithTimeZoneOffset(c.offset))
			if got := c.ts.BeginningOfDay(); got != c.wantBOD {
				t.Fatalf("BeginningOfDay() = %d, want %d", got, c.wantBOD)
			}
			if got := c.ts.BeginningOfNextDay(); got != c.wantNextDay {
				t.Fatalf("BeginningOfNextDay() = %d, want %d", got, c.wantNextDay)
			}
			// 幂等
			bod := c.ts.BeginningOfDay()
			if bod.BeginningOfDay() != bod {
				t.Fatal("BeginningOfDay not idempotent")
			}
		})
	}
}

// ---- InSameDay / DiffDays / IsSameDay ----

func TestInSameDayAndDiffDays(t *testing.T) {
	defer resetGlobal()
	Initialize(WithTimeZoneOffset(8))

	a := makeTS(8, 2025, 4, 10, 1, 0, 0)
	b := makeTS(8, 2025, 4, 10, 23, 59, 59)
	c := makeTS(8, 2025, 4, 11, 0, 0, 0)

	if !a.InSameDay(b) {
		t.Fatal("InSameDay(same day) want true")
	}
	if b.InSameDay(c) {
		t.Fatal("InSameDay(across midnight) want false")
	}
	if !IsSameDay(a, b) {
		t.Fatal("IsSameDay(same day) want true")
	}
	if IsSameDay(a, c) {
		t.Fatal("IsSameDay(across midnight) want false")
	}

	if a.DiffDays(a) != 0 {
		t.Fatalf("DiffDays(self) = %d, want 0", a.DiffDays(a))
	}
	if c.DiffDays(a) != 1 {
		t.Fatalf("DiffDays(next) = %d, want 1", c.DiffDays(a))
	}
	if a.DiffDays(c) != -1 {
		t.Fatalf("DiffDays(prev) = %d, want -1", a.DiffDays(c))
	}
	if d, e := makeTS(8, 2025, 4, 13, 12, 0, 0), makeTS(8, 2025, 4, 10, 6, 0, 0); d.DiffDays(e) != 3 {
		t.Fatalf("DiffDays(3 apart) = %d, want 3", d.DiffDays(e))
	}
}

// ---- WithinToday / IsWithinToday ----

func TestWithinToday(t *testing.T) {
	defer resetGlobal()
	// 冻结当前为 2025-04-10 12:00 UTC+8
	setupFixedNow(8, time.Date(2025, 4, 10, 4, 0, 0, 0, time.UTC))

	today := Now()
	if !today.WithinToday() {
		t.Fatal("today WithinToday() want true")
	}
	if !IsWithinToday(today) {
		t.Fatal("today IsWithinToday() want true")
	}
	if today.AddS(-Day).WithinToday() {
		t.Fatal("yesterday WithinToday() want false")
	}
	if IsWithinToday(today.AddS(Day)) {
		t.Fatal("tomorrow IsWithinToday() want false")
	}
	// 同一天不同时刻
	if !makeTS(8, 2025, 4, 10, 0, 0, 0).WithinToday() {
		t.Fatal("same-day 00:00 WithinToday() want true")
	}
	if !makeTS(8, 2025, 4, 10, 23, 59, 59).WithinToday() {
		t.Fatal("same-day 23:59:59 WithinToday() want true")
	}
}

// ---- InSameWeek / InSameMonth ----

func TestInSameWeek(t *testing.T) {
	defer resetGlobal()
	Initialize(WithTimeZoneOffset(8))

	// 2025-04-07 (Mon) ~ 2025-04-13 (Sun) is ISO week 15
	mon := makeTS(8, 2025, 4, 7, 12, 0, 0)
	sun := makeTS(8, 2025, 4, 13, 12, 0, 0)
	if !mon.InSameWeek(sun) {
		t.Fatal("Mon-Sun same ISO week want true")
	}
	sunEnd := makeTS(8, 2025, 4, 13, 23, 59, 59)
	nextMon := makeTS(8, 2025, 4, 14, 0, 0, 0)
	if sunEnd.InSameWeek(nextMon) {
		t.Fatal("across ISO week want false")
	}
	// 跨年 ISO 周边界：2024-12-30(Mon) 与 2025-01-01(Wed) 同属 ISO 2025-W01
	a := makeTS(8, 2024, 12, 30, 12, 0, 0)
	b := makeTS(8, 2025, 1, 1, 12, 0, 0)
	if !a.InSameWeek(b) {
		t.Fatal("2024-12-30 and 2025-01-01 should be same ISO week")
	}
}

func TestInSameMonth(t *testing.T) {
	defer resetGlobal()
	Initialize(WithTimeZoneOffset(8))

	if a, b := makeTS(8, 2025, 4, 1, 0, 0, 0), makeTS(8, 2025, 4, 30, 23, 59, 59); !a.InSameMonth(b) {
		t.Fatal("same month want true")
	}
	if a, b := makeTS(8, 2025, 4, 30, 23, 59, 59), makeTS(8, 2025, 5, 1, 0, 0, 0); a.InSameMonth(b) {
		t.Fatal("diff month want false")
	}
	if a, b := makeTS(8, 2025, 4, 10, 0, 0, 0), makeTS(8, 2026, 4, 10, 0, 0, 0); a.InSameMonth(b) {
		t.Fatal("same month diff year want false")
	}
}

// ---- BeginningOfDayOfSameWeek ----

func TestBeginningOfDayOfSameWeek(t *testing.T) {
	defer resetGlobal()
	Initialize(WithTimeZoneOffset(8))

	// 2025-04-09 Wednesday UTC+8
	wed := makeTS(8, 2025, 4, 9, 15, 0, 0)

	// Monday as week beginning
	if got, exp := wed.BeginningOfDayOfSameWeek(time.Monday, true), makeTS(8, 2025, 4, 7, 0, 0, 0); got != exp {
		t.Fatalf("Monday of week = %d, want %d", got, exp)
	}
	if got, exp := wed.BeginningOfDayOfSameWeek(time.Sunday, true), makeTS(8, 2025, 4, 13, 0, 0, 0); got != exp {
		t.Fatalf("Sunday of Mon-week = %d, want %d", got, exp)
	}
	if got, exp := wed.BeginningOfDayOfSameWeek(time.Wednesday, true), makeTS(8, 2025, 4, 9, 0, 0, 0); got != exp {
		t.Fatalf("Wednesday itself = %d, want %d", got, exp)
	}
	// 从周日回看本周一（Mon-based 下周日属于本周末）
	sun := makeTS(8, 2025, 4, 13, 10, 0, 0)
	if got, exp := sun.BeginningOfDayOfSameWeek(time.Monday, true), makeTS(8, 2025, 4, 7, 0, 0, 0); got != exp {
		t.Fatalf("from Sunday get Monday = %d, want %d", got, exp)
	}

	// Sunday as week beginning
	if got, exp := wed.BeginningOfDayOfSameWeek(time.Sunday, false), makeTS(8, 2025, 4, 6, 0, 0, 0); got != exp {
		t.Fatalf("Sunday start = %d, want %d", got, exp)
	}
	if got, exp := wed.BeginningOfDayOfSameWeek(time.Saturday, false), makeTS(8, 2025, 4, 12, 0, 0, 0); got != exp {
		t.Fatalf("Saturday of Sun-week = %d, want %d", got, exp)
	}
}

// ---- Week0 / NextWeek0 ----

func TestWeek0AndNextWeek0(t *testing.T) {
	defer resetGlobal()
	// 固定当前为 2025-04-09 Wed 15:00 UTC+8
	setupFixedNow(8, time.Date(2025, 4, 9, 7, 0, 0, 0, time.UTC))

	wMon := Week0(true)
	if wd := wMon.Time().Weekday(); wd != time.Monday {
		t.Fatalf("Week0(true) weekday = %v, want Monday", wd)
	}
	if h, m := wMon.Time().Hour(), wMon.Time().Minute(); h != 0 || m != 0 {
		t.Fatalf("Week0(true) not midnight: %02d:%02d", h, m)
	}
	if exp := makeTS(8, 2025, 4, 7, 0, 0, 0); wMon != exp {
		t.Fatalf("Week0(true) = %d, want %d", wMon, exp)
	}

	wSun := Week0(false)
	if wd := wSun.Time().Weekday(); wd != time.Sunday {
		t.Fatalf("Week0(false) weekday = %v, want Sunday", wd)
	}
	if exp := makeTS(8, 2025, 4, 6, 0, 0, 0); wSun != exp {
		t.Fatalf("Week0(false) = %d, want %d", wSun, exp)
	}

	if nw := NextWeek0(true); nw.SubS(wMon) != Week || nw.Time().Weekday() != time.Monday {
		t.Fatalf("NextWeek0(true): subS=%d wd=%v", nw.SubS(wMon), nw.Time().Weekday())
	}
	if nw := NextWeek0(false); nw.SubS(wSun) != Week || nw.Time().Weekday() != time.Sunday {
		t.Fatalf("NextWeek0(false): subS=%d wd=%v", nw.SubS(wSun), nw.Time().Weekday())
	}
}

// ---- 跨时区边界 ----

func TestCrossTimezone(t *testing.T) {
	defer resetGlobal()

	// 2025-04-10 23:30:00 UTC → 2025-04-11 07:30:00 UTC+8
	utcTs := Timestamp(time.Date(2025, 4, 10, 23, 30, 0, 0, time.UTC).Unix())

	Initialize(WithTimeZoneOffset(0))
	if d := utcTs.BeginningOfDay().TimeUTC().Day(); d != 10 {
		t.Fatalf("UTC+0 BeginningOfDay day = %d, want 10", d)
	}

	Initialize(WithTimeZoneOffset(8))
	if d := utcTs.BeginningOfDay().Time().Day(); d != 11 {
		t.Fatalf("UTC+8 BeginningOfDay day = %d, want 11", d)
	}
}

// ---- 缩放时间提供者 ----

func TestScaledProviderNow(t *testing.T) {
	defer resetGlobal()

	// scale=2：现实流逝会被放大 2 倍映射到虚拟时间；这里只校验 Now 基于 provider 而非系统时钟。
	virtual := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	Initialize(WithTimeProvider(scaledProvider{
		virtualBase: virtual,
		actualBase:  time.Now().UTC(),
		scale:       2,
	}))
	if y := Now().TimeUTC().Year(); y != 2030 {
		t.Fatalf("Now().Year() = %d, want 2030 (provider virtual base)", y)
	}
}

// ---- 并发安全 ----

func TestConcurrentReadWrite(t *testing.T) {
	defer resetGlobal()
	Initialize(WithTimeZoneOffset(8))

	fixed := time.Date(2025, 4, 10, 4, 0, 0, 0, time.UTC)

	var wg sync.WaitGroup
	// 写方：反复 Initialize（含 WithNowFunc）
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				Initialize(WithTimeZoneOffset(int64(i)), WithNowFunc(func() time.Time { return fixed }))
			}
		}(i)
	}
	// 读方：反复调用依赖全局 clock 的方法
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ts := makeTS(8, 2025, 4, 10, 12, 0, 0)
			for j := 0; j < 200; j++ {
				_ = Now()
				_ = ts.BeginningOfDay()
				_ = ts.Time()
				_ = Location()
				_ = Week0(true)
			}
		}()
	}
	wg.Wait()
}
