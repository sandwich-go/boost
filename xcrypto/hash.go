package xcrypto

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func HashFileMd5(filePath string) (string, error) {
	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	//Tell the program to call the following function when the current function returns
	defer func() { _ = file.Close() }()

	return HashBufferMd5(file)

}

func HashBufferMd5(src io.Reader) (string, error) {
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
