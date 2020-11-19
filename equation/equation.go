/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/13 2:40 PM
 */

package equation

import (
	"errors"
	"fmt"
	"github.com/dengsgo/math-engine/calculate"
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/entry"
	"github.com/patrickmn/go-cache"
	"strings"
)

type Equation struct {
	exp       string
	resultTag string
}

// TODO: 更优雅的解析方式及错误检测
func New(equ string) (*Equation, error) {
	ss := strings.Split(equ, "=")
	if len(ss) != 2 {
		return nil, errors.New("may have wrong format equation")
	}
	return &Equation{
		exp:       ss[1],
		resultTag: ss[0],
	}, nil
}

func ExecEquation(es []Equation) (result *common.ArithmeticFactor, processLog string) {
	defer common.Cache.Flush()

	var log string
	for index, e := range es {
		r := calculate.Exec(e.exp)
		log = fmt.Sprintf("%s\n方程内容:%s\n%s", log, e.exp, entry.String())
		if index == (len(es) - 1) {
			return r, log
		}
		common.Cache.Set(e.resultTag, *r, cache.DefaultExpiration)
	}
	return
}
