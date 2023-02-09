package xcmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	flagPrefix    = "xcmd_"
	parsedOptions = make(map[string]string)
	argumentRegex = regexp.MustCompile(`^\-{1,2}([\w\?\.\-]+)(=){0,1}(.*)$`)
	formal        = make(map[string]string)
	invalidChars  = ".-/\\"
)

// SetFlagPrefix 设置 flag 的前缀
func SetFlagPrefix(prefix string) { flagPrefix = prefix }

// GetFlagPrefix 获取 flag 的前缀
func GetFlagPrefix() string { return flagPrefix }

func addFlagIfNotExists(fs *flag.FlagSet, name string, value string) {
	if fs.Lookup(name) == nil {
		usage := fmt.Sprintf("%sbase_command flag:%s default:%s", flagPrefix, name, value)
		_ = fs.String(name, value, usage)
	}
}

// DeclareInto 将 command 预留的 flag 定义添加到指定的 FlagSet 中
func DeclareInto(fs *flag.FlagSet) {
	for name, val := range formal {
		addFlagIfNotExists(fs, name, val)
	}
}

// AddFlag 添加框架内默认的 flag 数据
// name 必须有 flagPrefix 前缀，不能包含.-/\字符
func AddFlag(name, defaultVal string) error {
	if !strings.HasPrefix(name, flagPrefix) {
		return fmt.Errorf("command flag must has prefix: %s", flagPrefix)
	}
	if strings.ContainsAny(name, invalidChars) {
		return fmt.Errorf("command flag should only use _ as world separator, format must be:%s<package name>_<variable name>", flagPrefix)
	}
	_, alreadyThere := formal[name]
	if alreadyThere {
		return fmt.Errorf("flag redefined: %s", name)
	}
	formal[name] = defaultVal
	addFlagIfNotExists(flag.CommandLine, name, defaultVal)
	return nil
}

// MustAddFlag 添加框架内默认的 flag 数据，若失败，则panic
// name 必须有 flagPrefix 前缀，不能包含特殊的字符
func MustAddFlag(name, defaultVal string) {
	err := AddFlag(name, defaultVal)
	if err != nil {
		panic(err)
	}
}

// Init 根据给定的参数初始化预留的参数，如果args为空，则会使用os.Args初始化
func Init(args ...string) {
	if len(args) == 0 {
		if len(parsedOptions) != 0 { // 已经初始化过数据
			return
		}
		args = os.Args
	}

	if len(parsedOptions) == 0 {
		parsedOptions = make(map[string]string)
	}
	for i := 0; i < len(args); {
		match := argumentRegex.FindStringSubmatch(args[i])
		if len(match) > 2 {
			if match[2] == "=" {
				// -xcmd_debug=1
				// array  4 [-xcmd_debug=1 xcmd_debug = 1]
				parsedOptions[match[1]] = match[3]
			} else if i < len(args)-1 {
				if len(args[i+1]) > 0 && args[i+1][0] == '-' {
					parsedOptions[match[1]] = match[3]
				} else {
					parsedOptions[match[1]] = args[i+1]
					i += 2
					continue
				}
			} else {
				parsedOptions[match[1]] = match[3]
			}
		}
		i++
	}
}

// GetOptWithEnv 返回命令定义的行参数或Env中的参数
// 1. 命令行参数，小写，单词之间以_分割,$FlagPrefix_<package name>_<variable name>
// 2. 环境变量参数，小写(历史遗留),且单词之间以_分割,$FlagPrefix_<package name>_<variable name>
func GetOptWithEnv(key string, def ...string) string {
	Init()
	if v, ok := parsedOptions[key]; ok {
		return v
	}
	if r, ok := os.LookupEnv(key); ok {
		return r
	}
	if len(def) > 0 {
		return def[0]
	}
	if v, ok := formal[key]; ok {
		return v
	}
	return ""
}

// ContainsOpt checks whether option named `name` exist in the arguments.
func ContainsOpt(name string) bool {
	Init()
	_, ok := parsedOptions[name]
	return ok
}
