package buildinfo

import (
	"fmt"
	"runtime/debug"

	"github.com/sandwich-go/boost/version"
)

// PrintBuildInfo 打印BuildInfo
func PrintBuildInfo() {
	fmt.Println("Build Version:")
	fmt.Printf("\t%s", version.String())
	fmt.Println()
	if info, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("Main ModuleName:")
		printModule(&info.Main)
		fmt.Println("Dependencies:")
		for _, dep := range info.Deps {
			printModule(dep)
		}
	} else {
		fmt.Println("Built without Go modules")
	}
}

func printModule(m *debug.Module) {
	fmt.Printf("\t%s", m.Path)
	fmt.Printf("@%s", m.Version)
	if m.Sum != "" {
		fmt.Printf(" (sum: %s)", m.Sum)
	}
	if m.Replace != nil {
		fmt.Printf(" (replace: %s)", m.Replace.Path)
	}
	fmt.Println()
}
