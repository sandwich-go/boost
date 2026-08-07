package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// File 对文件进行 md5 hash，换行符差异不影响 Windows 和 Unix 平台下的结果。
func File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := md5.New()
	buf := make([]byte, 32*1024)
	normalized := make([]byte, 0, len(buf))
	pendingCR := false

	for {
		n, readErr := file.Read(buf)
		if n > 0 {
			normalized = normalized[:0]
			for _, b := range buf[:n] {
				if pendingCR {
					normalized = append(normalized, '\n')
					pendingCR = false
					if b == '\n' {
						continue
					}
				}
				if b == '\r' {
					pendingCR = true
				} else {
					normalized = append(normalized, b)
				}
			}
			_, _ = hash.Write(normalized)
		}

		if readErr == io.EOF {
			if pendingCR {
				_, _ = hash.Write([]byte{'\n'})
			}
			break
		}
		if readErr != nil {
			return "", readErr
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Buffer 对数据流进行 md5 hash
func Buffer(src io.Reader) (string, error) {
	var returnMD5String string
	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, src); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	return hex.EncodeToString(hashInBytes), nil
}
