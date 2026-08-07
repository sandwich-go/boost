package md5

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sandwich-go/boost/xhash/nhash/jenkins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMd5(t *testing.T) {
	Convey("md5", t, func() {
		s, err := File("test")
		So(err, ShouldBeNil)

		f, err0 := os.Open("test")
		So(err0, ShouldBeNil)
		defer func() {
			_ = f.Close()
		}()
		s1, err1 := Buffer(f)
		So(err1, ShouldBeNil)
		So(s, ShouldEqual, s1)

		s2, err2 := Buffer(bytes.NewReader([]byte("aaaaaaaa")))
		So(err2, ShouldBeNil)
		t.Log(s2)

		hint, _ := jenkins.HashString("aaaaaaaa", 0, 0)
		fmt.Println(hint)
	})
}

func TestFileNormalizedLineEndings(t *testing.T) {
	dir := t.TempDir()
	writeFile := func(name, content string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	hashFile := func(path string) string {
		t.Helper()
		hash, err := File(path)
		if err != nil {
			t.Fatal(err)
		}
		return hash
	}

	withoutNewline := hashFile(writeFile("without-newline", "helloworld"))
	withLF := hashFile(writeFile("with-lf", "hello\nworld\n"))
	withCRLF := hashFile(writeFile("with-crlf", "hello\r\nworld\r\n"))
	withCR := hashFile(writeFile("with-cr", "hello\rworld\r"))
	withSpace := hashFile(writeFile("with-space", "hello world"))

	if withLF == withoutNewline {
		t.Fatalf("adding a line break should affect hash: both hashes are %q", withLF)
	}
	if withCRLF != withLF {
		t.Fatalf("CRLF and LF should produce the same hash: got %q, want %q", withCRLF, withLF)
	}
	if withCR != withLF {
		t.Fatalf("CR and LF should produce the same hash: got %q, want %q", withCR, withLF)
	}
	if withSpace == withoutNewline {
		t.Fatalf("space should affect hash: both hashes are %q", withSpace)
	}

	boundaryPrefix := strings.Repeat("a", 32*1024-1)
	boundaryCRLF := hashFile(writeFile("boundary-crlf", boundaryPrefix+"\r\nb"))
	boundaryLF := hashFile(writeFile("boundary-lf", boundaryPrefix+"\nb"))
	if boundaryCRLF != boundaryLF {
		t.Fatalf("CRLF split across reads should match LF: got %q, want %q", boundaryCRLF, boundaryLF)
	}
}

func TestFileMissingFile(t *testing.T) {
	if _, err := File(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("File should return an error for a missing file")
	}
}
