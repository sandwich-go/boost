package xos

import "runtime"

// getShell returns the shell command depending on current working operation system.
// It returns "cmd.exe" for windows, and "bash" or "sh" for others.
func GetShell() string {
	switch runtime.GOOS {
	case "windows":
		return SearchBinary("cmd.exe")
	default:
		// Check the default binary storage path.
		if FileExists("/bin/bash") {
			return "/bin/bash"
		}
		if FileExists("/bin/sh") {
			return "/bin/sh"
		}
		// Else search the env PATH.
		path := SearchBinary("bash")
		if path == "" {
			path = SearchBinary("sh")
		}
		return path
	}
}

// getShellOption returns the shell option depending on current working operating system.
// It returns "/c" for windows, and "-c" for others.
// -c string
//      If the -c option is present, then commands are read from string.
//      If there are arguments after the string, they are assigned to the positional
//      parameters, starting with $0.
func GetShellOption() string {
	switch runtime.GOOS {
	case "windows":
		return "/c"
	default:
		return "-c"
	}
}
