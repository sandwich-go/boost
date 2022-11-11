package xcmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var FlagPrefix = "xcmd_"

func SetFlagPrefix(prefix string) {
	FlagPrefix = prefix
}

var (
	parsedOptions = make(map[string]string)
	argumentRegex = regexp.MustCompile(`^\-{1,2}([\w\?\.\-]+)(=){0,1}(.*)$`)
	formal        = make(map[string]string)
)

func cleanParsedOptions() { parsedOptions = make(map[string]string) }

// Empty strings.
var emptyStringMap = map[string]struct{}{
	"":      {},
	"0":     {},
	"no":    {},
	"off":   {},
	"false": {},
}

// IsFalse 判断command解析获取的数据是否为false
func IsFalse(v string) bool {
	_, ok := emptyStringMap[strings.ToLower(string(v))]
	return ok
}

// IsTrue 判断command解析获取的数据是否为true
func IsTrue(v string) bool { return !IsFalse(v) }

const defaultSliceSeparator = ","

// Slice 将 defaultSliceSeparator 分割的字符获取slice
func Slice(v string) []string {
	r := strings.Split(v, defaultSliceSeparator)
	var ss []string
	for _, s := range r {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}

// DeclareInto 将command预留的flag定义添加到指定的FlagSet中
func DeclareInto(fs *flag.FlagSet) {
	for name, val := range formal {
		// 防止使用默认的flag set报错: flag provided but not defined
		if fs.Lookup(name) == nil {
			fs.String(name, val, fmt.Sprintf("%sbase_command flag:%s default:%s", FlagPrefix, name, val))
		}
	}
}

// AddFlag 添加框架内默认的Flag数据，默认的Flag数据可以通过调用AutoFlag添加到指定的FlagSet以防止出现flag provided but not defined
// 默认会在Init调用的时候对flag.CommandLine调用AutoFlag
func AddFlag(name, defaultVal string) {
	if !strings.HasPrefix(name, FlagPrefix) {
		panic("command flag must has prefix: " + FlagPrefix)
	}
	if strings.ContainsAny(name, ".-/\\") {
		panic(fmt.Sprintf("command flag should only use _ as world separator, format must be:%s<package name>_<variable name>", FlagPrefix))
	}
	_, alreadyThere := formal[name]
	if alreadyThere {
		panic(fmt.Sprintf("flag redefined: %s", name))
	}
	formal[name] = defaultVal
	// 添加到默认flag
	if flag.CommandLine.Lookup(name) == nil {
		flag.CommandLine.String(name, defaultVal, fmt.Sprintf("%sbase_command flag:%s default:%s", FlagPrefix, name, defaultVal))
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
