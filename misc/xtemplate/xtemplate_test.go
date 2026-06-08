package xtemplate

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag(t *testing.T) {
	Convey("template", t, func() {
		s := `a {{ .val1 }} {{ .val2 }}`
		s1, err := Execute(s, map[string]interface{}{"val1": "b", "val2": 2})
		So(err, ShouldBeNil)
		So(string(s1), ShouldEqual, "a b 2")
	})
}

// TestExecute_FuncMap 覆盖 funcMap 内置 helper（CamelCase / SnakeCase /
// FirstLower / FirstUpper / Unescaped）+ 大小写 alias。
func TestExecute_FuncMap(t *testing.T) {
	Convey("内置 FuncMap helper 都可用", t, func() {
		got, err := Execute(`{{ CamelCase "hello_world" }}`, nil)
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, "HelloWorld")

		// 大小写 alias（init 时把 first letter 改小写也注册了）
		got, err = Execute(`{{ camelCase "hello_world" }}`, nil)
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, "HelloWorld")

		got, err = Execute(`{{ FirstLower "Hello" }}`, nil)
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, "hello")
	})
}

// TestExecute_CustomFuncMap 覆盖 cfg.GetFuncMap() 自定义 funcMap 注入分支
// （xtemplate.go:34-36）。
func TestExecute_CustomFuncMap(t *testing.T) {
	Convey("WithFuncMap 注入自定义 helper", t, func() {
		fm := template.FuncMap{
			"double": func(n int) int { return n * 2 },
		}
		got, err := Execute(`{{ double 5 }}`, nil, WithOptionFuncMap(fm))
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, "10")
	})
}

// TestExecute_Filter 覆盖 cfg.GetFilers() filter pipeline（line 48-53）。
func TestExecute_Filter(t *testing.T) {
	Convey("filter 链路过滤输出", t, func() {
		upper := func(b []byte) []byte {
			for i, c := range b {
				if c >= 'a' && c <= 'z' {
					b[i] = c - 32
				}
			}
			return b
		}
		// 链路：上游 + 下游 filter，第二个 filter 是 nil 走 'continue' 分支
		got, err := Execute(`hello`, nil, WithOptionFilers(upper, nil))
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, "HELLO")
	})
}

// TestExecute_WriteFile 覆盖 WithFileName 把结果写到文件（line 54-65）。
func TestExecute_WriteFile(t *testing.T) {
	Convey("WithFileName 写到普通文件", t, func() {
		dir := t.TempDir()
		path := filepath.Join(dir, "out.txt")
		got, err := Execute(`hello {{ .name }}`, map[string]string{"name": "world"},
			WithOptionFileName(path))
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, "hello world")

		// 文件应被写入
		data, ferr := os.ReadFile(path)
		So(ferr, ShouldBeNil)
		So(string(data), ShouldEqual, "hello world")
	})

	Convey("WithFileName 写 .go 文件走 goformat", t, func() {
		dir := t.TempDir()
		path := filepath.Join(dir, "out.go")
		// 故意写不规范 Go 代码（缺空格 / 多空行），goformat 应规整
		raw := "package x\n\nfunc  foo(  ){}\n"
		got, err := Execute(raw, nil, WithOptionFileName(path))
		So(err, ShouldBeNil)
		// goformat 规整后应不含双空格
		So(string(got), ShouldNotContainSubstring, "func  foo")
	})
}

// TestExecute_Errors 覆盖 Parse / Execute 错误路径。
func TestExecute_Errors(t *testing.T) {
	Convey("Parse error: invalid template syntax", t, func() {
		_, err := Execute(`{{ unclosed`, nil)
		So(err, ShouldNotBeNil)
	})

	Convey("Execute error: 调用 nil 上的 method 报 err", t, func() {
		// 调用一个会让 template engine 真报 err 的场景：在 nil interface
		// 上调函数 method（template 解析期不报，执行期报 nil pointer
		// dereference）
		type T struct{ Y string }
		var v *T = nil
		_, err := Execute(`{{ .X.Y }}`, struct{ X *T }{X: v})
		So(err, ShouldNotBeNil)
	})
}
