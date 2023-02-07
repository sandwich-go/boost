package xos

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sandwich-go/boost/xpanic"
)

// FileCopyToDir the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func FileCopyToDir(src, dstDir string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()

	out, err0 := os.Create(filepath.Join(dstDir, filepath.Base(src)))
	if err0 != nil {
		return err0
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

// Ext 返回后缀，例如 'xxx.go' => '.go'
func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// FileGetContents 获取文件内容
func FileGetContents(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// MustFilePutContents 写入文件，如果发生错误则panic
func MustFilePutContents(filename string, content []byte) {
	dirName := filepath.Dir(filename)
	xpanic.WhenErrorAsFmtFirst(os.MkdirAll(dirName, os.ModePerm), "got error:%w while MkdirAll with:%s", dirName)
	xpanic.WhenErrorAsFmtFirst(ioutil.WriteFile(filename, content, 0644), "got error:%w while WriteFile with:%s", filename)
}

// FilePutContents 写入文件
func FilePutContents(filename string, content []byte) error {
	dirName := filepath.Dir(filename)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, content, 0644)
}

// MustGetFileWriter 获取写文件句柄
// filePath 指定文件
// prepend 是否保留源数据，如果保留，则源数据会被追加到文件尾
func MustGetFileWriter(filePath string, prepend bool) (writer io.Writer, deferFunc func()) {
	var prependData []byte
	if prepend {
		if ExistsFile(filePath) {
			fileContent, err := FileGetContents(filePath)
			xpanic.WhenErrorAsFmtFirst(err, "got error:%w while FileGetContents:%s", filePath)
			prependData = fileContent
		}
	}
	dirParent := filepath.Dir(filePath)
	xpanic.WhenErrorAsFmtFirst(os.MkdirAll(dirParent, os.ModePerm), "got error:%w while MkdirAll:%s", dirParent)
	output, err := os.Create(filePath)
	xpanic.WhenErrorAsFmtFirst(err, "got error:%w while Create:%s", filePath)
	return output, func() {
		_, _ = output.Write(prependData)
		_ = output.Close()
	}
}
