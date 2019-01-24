package string

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func GetFileMd5(path string) string {
	file, inerr := os.Open(path)
	defer file.Close()
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		return fmt.Sprintf("%x", md5h.Sum(nil))
	}
	return ""
}

func GetStrMD5(src string) string {
	return GetByteMD5([]byte(src))
}

func GetByteMD5(src []byte) string {
	cipherStr := md5.Sum(src)
	md5str1 := fmt.Sprintf("%x", cipherStr) //将[]byte转成16进制

	return md5str1
}
