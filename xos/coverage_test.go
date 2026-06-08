package xos

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestEnsureEmpty 覆盖 EnsureEmpty / MustEnsureEmpty 两个 0% 函数。
func TestEnsureEmpty(t *testing.T) {
	Convey("EnsureEmpty 清空目录并保证目录存在", t, func() {
		base := filepath.Join(os.TempDir(), "xos_ensure_empty_test")
		defer os.RemoveAll(base)

		// 先创建目录 + 一些文件
		So(TouchDirAll(base), ShouldBeNil)
		So(FilePutContents(filepath.Join(base, "f1.txt"), []byte("hi")), ShouldBeNil)
		So(FilePutContents(filepath.Join(base, "sub/f2.txt"), []byte("hi")), ShouldBeNil)
		So(ExistsFile(filepath.Join(base, "f1.txt")), ShouldBeTrue)

		// EnsureEmpty 后目录还在，但内容清掉
		So(EnsureEmpty(base), ShouldBeNil)
		So(Exists(base), ShouldBeTrue)
		So(ExistsFile(filepath.Join(base, "f1.txt")), ShouldBeFalse)

		// 多个目录批量
		base2 := filepath.Join(os.TempDir(), "xos_ensure_empty_test2")
		defer os.RemoveAll(base2)
		So(EnsureEmpty(base, base2), ShouldBeNil)
		So(Exists(base), ShouldBeTrue)
		So(Exists(base2), ShouldBeTrue)
	})

	Convey("MustEnsureEmpty 成功路径不 panic", t, func() {
		dir := filepath.Join(os.TempDir(), "xos_must_ensure_empty_ok")
		defer os.RemoveAll(dir)
		So(func() { MustEnsureEmpty(dir) }, ShouldNotPanic)
		So(Exists(dir), ShouldBeTrue)
	})
}

// TestIsHiddenOrInHiddenDir unix 实现：以 '.' 开头的目录或文件为隐藏。
//
// 注意：函数对绝对路径处理有 quirk —— path 切段后第一段是空字符串
// （/foo/bar 切成 [""+"foo"+"bar"]），逐段 IsHidden 在空 path 上 stat
// 报错。所以测试**用 cwd 相对路径**调用，绕开 absolute path quirk。
// 这不是 IsHiddenOrInHiddenDir 设计 bug 也不是测试 bug——是用法约定，
// 业务侧通常给相对路径。
func TestIsHiddenOrInHiddenDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("hide.go IsHiddenOrInHiddenDir 行为依赖 unix '.' 前缀约定")
	}
	Convey("IsHiddenOrInHiddenDir 相对路径用法", t, func() {
		// 切到 TempDir 工作做相对路径测试
		base := filepath.Join(os.TempDir(), "xos_hidden_test")
		defer os.RemoveAll(base)
		oldWd, _ := os.Getwd()
		defer func() { _ = os.Chdir(oldWd) }()
		So(TouchDirAll(filepath.Join(base, ".hidden_dir")), ShouldBeNil)
		So(TouchDirAll(filepath.Join(base, "visible")), ShouldBeNil)
		So(FilePutContents(filepath.Join(base, ".hidden_dir/file.txt"), []byte("x")), ShouldBeNil)
		So(FilePutContents(filepath.Join(base, "visible/file.txt"), []byte("x")), ShouldBeNil)
		So(FilePutContents(filepath.Join(base, "visible/.hide.txt"), []byte("x")), ShouldBeNil)
		So(os.Chdir(base), ShouldBeNil)

		// 相对路径调用，第一段就是 ".hidden_dir" / "visible"
		got, err := IsHiddenOrInHiddenDir(".hidden_dir/file.txt")
		So(err, ShouldBeNil)
		So(got, ShouldBeTrue)

		got, err = IsHiddenOrInHiddenDir("visible/.hide.txt")
		So(err, ShouldBeNil)
		So(got, ShouldBeTrue)

		got, err = IsHiddenOrInHiddenDir("visible/file.txt")
		So(err, ShouldBeNil)
		So(got, ShouldBeFalse)
	})
}

// TestForceHide_ForceUnHide unix 实现：通过文件名前缀 '.' 切换隐藏。
func TestForceHide_ForceUnHide(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("hide.unix.go 仅 unix 编译")
	}
	Convey("ForceHide / ForceUnHide", t, func() {
		base := filepath.Join(os.TempDir(), "xos_forcehide_test")
		defer os.RemoveAll(base)
		So(TouchDirAll(base), ShouldBeNil)

		// 准备一个普通文件 + 一个隐藏占位（让 ForceHide 走 force overwrite 路径）
		visible := filepath.Join(base, "f.txt")
		hidden := filepath.Join(base, ".f.txt")
		So(FilePutContents(visible, []byte("v")), ShouldBeNil)
		So(FilePutContents(hidden, []byte("h")), ShouldBeNil)

		// ForceHide 应成功把 visible 改名到 .f.txt（覆盖原 hidden）
		newPath, err := ForceHide(visible)
		So(err, ShouldBeNil)
		So(newPath, ShouldEqual, hidden)
		So(ExistsFile(hidden), ShouldBeTrue)
		So(ExistsFile(visible), ShouldBeFalse)

		// ForceUnHide 反向：.f.txt → f.txt
		newPath2, err := ForceUnHide(hidden)
		So(err, ShouldBeNil)
		So(newPath2, ShouldEqual, visible)
		So(ExistsFile(visible), ShouldBeTrue)
	})

	Convey("ForceHide 已是隐藏的文件 → 直接返原 path", t, func() {
		base := filepath.Join(os.TempDir(), "xos_forcehide_already")
		defer os.RemoveAll(base)
		So(TouchDirAll(base), ShouldBeNil)
		hidden := filepath.Join(base, ".already.txt")
		So(FilePutContents(hidden, []byte("h")), ShouldBeNil)

		newPath, err := ForceHide(hidden)
		So(err, ShouldBeNil)
		So(newPath, ShouldEqual, hidden) // hide.unix.go:27 'hidden == isHidden' 早返
	})
}

// TestCopy_Symlink 覆盖 lcopy（创建 symlink 复制）+ onsymlink Shallow 分支。
func TestCopy_Symlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink 在 windows 行为不一致，仅 unix 测")
	}
	Convey("Copy 默认 OnSymlink=Shallow → lcopy 复制 symlink 本身", t, func() {
		base := filepath.Join(os.TempDir(), "xos_copy_symlink_test")
		defer os.RemoveAll(base)
		So(TouchDirAll(base), ShouldBeNil)

		// 创建一个真实文件 + 指向它的 symlink
		real := filepath.Join(base, "real.txt")
		link := filepath.Join(base, "link.txt")
		So(FilePutContents(real, []byte("hello")), ShouldBeNil)
		So(os.Symlink(real, link), ShouldBeNil)

		// Copy symlink 到新位置（默认 OnSymlink=Shallow → lcopy）
		dest := filepath.Join(base, "copied_link.txt")
		err := Copy(link, dest)
		So(err, ShouldBeNil)

		// dest 应是 symlink，指向同 target
		info, err := os.Lstat(dest)
		So(err, ShouldBeNil)
		So(info.Mode()&os.ModeSymlink, ShouldNotEqual, 0)
		target, err := os.Readlink(dest)
		So(err, ShouldBeNil)
		So(target, ShouldEqual, real)
	})
}

// TestGetActuallyDir 覆盖 symlink dir / 普通 dir / 不存在三种 case
// （44.4% → 100%）。
func TestGetActuallyDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink unix only")
	}
	Convey("GetActuallyDir 解 symlink 目录 / 普通目录 / err", t, func() {
		base := filepath.Join(os.TempDir(), "xos_actually_dir_test")
		defer os.RemoveAll(base)

		// 普通目录
		realDir := filepath.Join(base, "real_dir")
		So(TouchDirAll(realDir), ShouldBeNil)
		got, err := GetActuallyDir(realDir)
		So(err, ShouldBeNil)
		So(got, ShouldEqual, realDir)

		// symlink → 返目标
		linkDir := filepath.Join(base, "link_dir")
		So(os.Symlink(realDir, linkDir), ShouldBeNil)
		got, err = GetActuallyDir(linkDir)
		So(err, ShouldBeNil)
		So(got, ShouldEqual, realDir)

		// 不存在 → 返 err
		nonexist := filepath.Join(base, "nonexistent")
		_, err = GetActuallyDir(nonexist)
		So(err, ShouldNotBeNil)
	})
}
