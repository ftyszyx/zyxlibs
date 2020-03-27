package string

import (
	"crypto/md5"
	"encoding/base64"
	"crypto/rand"
	"fmt"
	"time"
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
	unix32bits := uint32(time.Now().UTC().Unix())
	buff := make([]byte, 12)
	_, err := rand.Read(buff)
	if  err != nil {
		panic(err)
	}
	return  fmt.Sprintf("%d-%s",unix32bits,base64.URLEncoding.EncodeToString(buff))
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
