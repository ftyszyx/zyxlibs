package string

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"errors"
//"github.com/ftyszyx/libs/beego"
	"github.com/axgle/mahonia"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)


// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

func CheckPhone(phone string) bool {
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func GbkToUtf8(src []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func UTF82GBK(src string) (string, error) {
	enc := mahonia.NewEncoder("gbk")
	//converts a  string from UTF-8 to gbk encoding.
	return enc.ConvertString(src), nil
}

func GetStrArr(src []interface{}) []string {
	var strarr []string
	for _, item := range src {
		strarr = append(strarr, item.(string))
	}
	return strarr
}

func FormatTableTime(stamp string) (datestring string,err error) {
	if stamp == "" || stamp == "0" {
		return "",nil
	}
	return FormatTime(stamp, "Y-m-d H:i:s")
}

func FormatTime(stamp string, format string) (datestring string,err error) {
	timeint, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		return
	}
	timeinfo := time.Unix(timeint, 0)
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	datestring= timeinfo.Format(format)
	return
}

var (
	idnum_coefficient []int32 = []int32{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	idnum_code        []byte  = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
)

func CheckIdNum(idCardNo string) bool {
	if len(idCardNo) != 18 {
		return false
	}

	idByte := []byte(strings.ToUpper(idCardNo))

	sum := int32(0)
	for i := 0; i < 17; i++ {
		sum += int32(byte(idByte[i])-byte('0')) * idnum_coefficient[i]
	}
	return idnum_code[sum%11] == idByte[17]
}

func addZero(timestring string) string {
	if len([]rune(timestring)) == 1 {
		return "0" + timestring
	}
	return timestring
}
func ParseTime(timestr string) (error, time.Time) {
	re, err := regexp.Compile(`\s+`)
	if err != nil {
		return err, time.Time{}
	}
	areaarr := re.Split(timestr, 2)
	if len(areaarr) != 2 {
		return errors.New("fomart err"), time.Time{}
	}
	datearr := strings.Split(areaarr[0], "/")
	if len(datearr) != 3 {
		return errors.New("date fomart err"), time.Time{}
	}

	datearr[1] = addZero(datearr[1])
	datearr[2] = addZero(datearr[2])

	timearr := strings.Split(areaarr[1], ":")
	if len(timearr) == 2 {
		timearr = strings.Split(areaarr[1]+":00", ":")
	}
	if len(timearr) != 3 {
		return errors.New("time fomart err"), time.Time{}
	}
	timearr[0] = addZero(timearr[0])
	timearr[1] = addZero(timearr[1])
	timearr[2] = addZero(timearr[2])
	newstr := strings.Join(datearr, "/") + " " + strings.Join(timearr, ":")
	//starttime, err := beego.DateParse(newstr, "Y/m/d H:i:s")
	replacer := strings.NewReplacer(datePatterns...)
	format := replacer.Replace("Y/m/d H:i:s")
	starttime, err :=time.ParseInLocation(format, newstr, time.Local)
	if err != nil {
		return err, time.Time{}
	}

	return nil, starttime
}

func SliceIndex(itemlist []string, callback func(item string) bool) int {
	for i := 0; i < len(itemlist); i++ {
		if callback(itemlist[i]) {
			return i
		}
	}
	return -1
}

func GetScheme(r *http.Request) string {
	switch {
	case r.URL.Scheme == "https":
		return "https"
	case r.TLS != nil:
		return "https"
	case strings.HasPrefix(r.Proto, "HTTPS"):
		return "https"
	case r.Header.Get("X-Forwarded-Proto") == "https":
		return "https"
	default:
		return "http"
	}
}

func GetHost(r *http.Request) string {
	switch {
	case len(r.Host) != 0:
		return r.Host
	case len(r.URL.Host) != 0:
		return r.URL.Host
	case len(r.Header.Get("X-Forwarded-For")) != 0:
		return r.Header.Get("X-Forwarded-For")
	case len(r.Header.Get("X-Host")) != 0:
		return r.Header.Get("X-Host")
	case len(r.Header.Get("XFF")) != 0:
		return r.Header.Get("XFF")
	case len(r.Header.Get("X-Real-IP")) != 0:
		return r.Header.Get("X-Real-IP")
	default:
		return "localhost:8080"
	}
}

func GetURL(r *http.Request) string {
	return GetScheme(r) + "://" + GetHost(r)
}
