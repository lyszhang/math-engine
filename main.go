package main

import (
	"bufio"
	"fmt"
	"github.com/dengsgo/math-engine/calculate"
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
	engine.RegFunction("compare", 2, engine.Compare)
	for {
		fmt.Print("input formulation/> ")
		s, err := readStdin()
		if err != nil {
			break
		}

		//TODO: 外部请求路径的输入
		//fmt.Print("input parameter query url/> ")

		start := time.Now()
		res := calculate.Exec(s)
		fmt.Println("result: ", res)
		cost := time.Since(start)
		fmt.Println("time: " + cost.String())
	}
}
