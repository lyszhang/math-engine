package main

import (
	"bufio"
	"fmt"
	"github.com/dengsgo/math-engine/calculate"
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/entry"
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
	engine.RegFunction("RATIO", 4, engine.Ratio)
	engine.RegFunction("MAX", 0, engine.Max)
	engine.RegFunction("MIN", 0, engine.Min)
	for {
		fmt.Print("input formulation/> ")
		s, err := readStdin()
		if err != nil {
			break
		}
		if len(s) == 0 {
			continue
		}

		start := time.Now()
		result := calculate.Exec(s)
		plog := entry.String()

		// demo 测试，尝试去客户端请求解密后的结果
		var r float64
		switch result.Factor {
		case common.TypePaillier:
			r, _ = source.UploadResult(result.Cipher.Data)
		case common.TypeConst:
			r, _ = result.Value()
		}

		fmt.Println("result: ", r)
		fmt.Println("plog: ", plog)
		cost := time.Since(start)
		fmt.Println("time: " + cost.String())
	}
}
