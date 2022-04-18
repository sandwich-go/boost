package xos

import (
	"os"
)

// SymlinkAction represents what to do on symlink.
type SymlinkAction int

const (
	// Deep creates hard-copy of contents.
	Deep SymlinkAction = iota
	// Shallow creates new symlink to the dest of symlink.
	Shallow
	// Skip does nothing with symlink.
	Skip
)

//go:generate optiongen --option_return_previous=false --option_prefix=WithCopy
func CopyOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// OnSymlink can specify what to do on symlink
		"OnSymlink": func(src string) SymlinkAction { return Shallow },
		// Skip can specify which files should be skipped
		"Skip": func(src string) (bool, error) { return false, nil },
		// AddPermission to every entities,
		"AddPermission": os.FileMode(0),
		// Sync file after copy.
		"Sync": false,
	}
}
