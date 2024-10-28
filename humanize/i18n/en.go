package i18n

import "strconv"

var LocaleEnglish = Locale{
	Days: func(i int64) string {
		if i == 1 {
			return "1 day"
		}
		return strconv.FormatInt(i, 10) + " days"
	},

	Hours: func(i int64) string {
		if i == 1 {
			return "1 hour"
		}
		return strconv.FormatInt(i, 10) + " hours"
	},

	Minutes: func(i int64) string {
		if i == 1 {
			return "1 minute"
		}
		return strconv.FormatInt(i, 10) + " minutes"
	},

	Seconds: func(i int64) string {
		if i == 1 {
			return "1 second"
		}
		return strconv.FormatInt(i, 10) + " seconds"
	},
}
