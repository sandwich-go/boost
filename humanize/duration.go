package humanize

import (
	"math"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/sandwich-go/boost/humanize/i18n"
)

// Duration 打印人类易读时间字符串
func Duration(duration time.Duration, lang language.Tag) string {
	loc, ok := i18n.Locales[lang]
	if !ok {
		loc = i18n.LocaleEnglish
	}

	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		fun    func(int64) string
		amount int64
	}{
		{loc.Days, days},
		{loc.Hours, hours},
		{loc.Minutes, minutes},
		{loc.Seconds, seconds},
	}

	var parts []string
	for _, chunk := range chunks {
		if chunk.amount == 0 {
			continue
		}
		parts = append(parts, chunk.fun(chunk.amount))
	}

	return strings.Join(parts, " ")
}
