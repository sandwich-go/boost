package i18n

import "strconv"

var LocaleChinese = Locale{
	Days: func(i int64) string {
		if i == 1 {
			return "1 天"
		}
		return strconv.FormatInt(i, 10) + " 天"
	},

	Hours: func(i int64) string {
		if i == 1 {
			return "1 小时"
		}
		return strconv.FormatInt(i, 10) + " 小时"
	},

	Minutes: func(i int64) string {
		if i == 1 {
			return "1 分钟"
		}
		return strconv.FormatInt(i, 10) + " 分钟"
	},

	Seconds: func(i int64) string {
		if i == 1 {
			return "1 秒"
		}
		return strconv.FormatInt(i, 10) + " 秒"
	},
}
