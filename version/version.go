package version

import (
	"fmt"
)

// Build information. Populated at build-time.
var (
	Version   = "unknown"
	Revision  = "unknown"
	Branch    = "unknown"
	BuildUser = "unknown"
	BuildDate = "unknown"
	UserData  = "unknown"
)

// Info provides the iterable version information.
var Info = map[string]string{
	"version":    Version,
	"revision":   Revision,
	"branch":     Branch,
	"build_user": BuildUser,
	"build_date": BuildDate,
	"user_data":  UserData,
}

// Valid version info is valid
func Valid() bool { return Version != "unknown" }

// String format version info
func String() string {
	if UserData == "unknown" || UserData == "" {
		return fmt.Sprintf("%s_%s_%s_%s_%s", Version, Revision, Branch, BuildUser, BuildDate)
	}
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s", Version, Revision, Branch, BuildUser, BuildDate, UserData)
}
