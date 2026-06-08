package xos

import (
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBinary(t *testing.T) {
	Convey("search binary file", t, func() {
		file := "binary.go"
		So(SearchBinary(file), ShouldEqual, file)
		So(SearchBinaryPath(file), ShouldBeEmpty)
	})
}

// TestSearchBinaryPath_FoundInPath unix 上 /bin/sh 一定在 PATH 中，
// SearchBinaryPath 应找到。
func TestSearchBinaryPath_FoundInPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/sh which doesn't exist on windows")
	}
	Convey("SearchBinaryPath 在 PATH 中找到 sh", t, func() {
		// PATH 通常包含 /bin
		got := SearchBinaryPath("sh")
		So(got, ShouldNotBeEmpty)
		So(ExistsFile(got), ShouldBeTrue)
	})
}

// TestSearchBinary_DelegatesToPath 当 cwd 没文件时 fallback 到 PATH 搜索。
func TestSearchBinary_DelegatesToPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/sh which doesn't exist on windows")
	}
	Convey("SearchBinary 转 SearchBinaryPath（cwd 没'sh'文件）", t, func() {
		got := SearchBinary("sh")
		So(got, ShouldNotBeEmpty)
	})
}

// TestSearchBinaryPath_NotFound 不存在的二进制返空。
func TestSearchBinaryPath_NotFound(t *testing.T) {
	Convey("SearchBinaryPath 不存在的 binary 返空", t, func() {
		So(SearchBinaryPath("definitely-no-such-binary-x9z8"), ShouldBeEmpty)
	})
}

// TestSearchBinaryPath_EmptyPATH 空 PATH 走 'array len == 0' 早返。
func TestSearchBinaryPath_EmptyPATH(t *testing.T) {
	Convey("空 PATH env 时返空（array len==0 早返）", t, func() {
		t.Setenv("PATH", "")
		if runtime.GOOS == "windows" {
			t.Setenv("Path", "")
		}
		So(SearchBinaryPath("anything"), ShouldBeEmpty)
	})
}
