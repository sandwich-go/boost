package i18n

import (
	"golang.org/x/text/language"
)

type Locale struct {
	Days    func(int64) string
	Hours   func(int64) string
	Minutes func(int64) string
	Seconds func(int64) string
}

var Locales = map[language.Tag]Locale{
	language.Chinese: LocaleChinese,
	language.English: LocaleEnglish,
}
