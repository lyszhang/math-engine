/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/10 2:05 PM
 */

package entry

type LogStr struct {
	log string
}

var logStr LogStr

func Reset() {
	logStr.log = ""
}

func Append(s string) {
	logStr.log = logStr.log + s
}

func String() string {
	return logStr.log
}
