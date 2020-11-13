/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/13 3:35 PM
 */

package main

import (
	"bufio"
	"fmt"
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/equation"
	"github.com/dengsgo/math-engine/source"
	"os"
	"strings"
	"time"

	"github.com/dengsgo/math-engine/engine"
)

func main() {
	loop()
}

func readStdin() (string, error) {
	f := bufio.NewReader(os.Stdin)
	s, err := f.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", err
	}
	if s == "exit" || s == "quit" || s == "q" {
		fmt.Println("bye")
		os.Exit(1)
	}
	return s, nil
}

// input loop
func loop() {
	engine.RegFunction("COMPARE", 2, engine.Compare)
	engine.RegFunction("COUNTONE", 0, engine.CountOne)
	for {

		var es []equation.Equation
		for {
			fmt.Print("input equation, if finished, enter 'end' to complete/> ")
			s, err := readStdin()
			if err != nil {
				break
			}
			if strings.ToLower(s) == "end" {
				break
			}

			e, err := equation.New(s)
			if err != nil {
				break
			}
			es = append(es, *e)
		}

		fmt.Println("====es====: ", es)
		start := time.Now()
		result, plog := equation.ExecEquation(es)

		// demo 测试，尝试去客户端请求解密后的结果
		var r int64
		switch result.Factor {
		case common.TypePaillier:
			r, _ = source.UploadResult(result.Cipher.Data)
		case common.TypeConst:
			fmt.Println("result: ", result.Number)
			r = result.Number
		}

		fmt.Println("result: ", r)
		fmt.Println("plog: ", plog)
		cost := time.Since(start)
		fmt.Println("time: " + cost.String())
	}
}
