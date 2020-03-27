package string

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
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

func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetStrMD5(base64.URLEncoding.EncodeToString(b))
}


//每个字节可以存两个16进制的数字
func KDNGetByteMD5(src []byte) string {
	cipherStr := md5.Sum(src)
	var codestr []byte
	for _, item := range cipherStr {
		codestr = append(codestr, item)
	}
	ignsstr := strings.ToLower(base64.StdEncoding.EncodeToString(codestr))
	return ignsstr
}
