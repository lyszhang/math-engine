package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dengsgo/math-engine/engine"
)

func main() {
	loop()
}

func readStdin() (string, error){
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
	//engine.RegFunction("double", 1, func(expr ...engine.ExprAST) float64 {
	//	return engine.ExprASTResult(expr[0]) * 2
	//})
	for {
		fmt.Print("input formulation/> ")
		s, err := readStdin()
		if err != nil {
			break
		}

		//TODO: 外部请求路径的输入
		//fmt.Print("input parameter query url/> ")

		start := time.Now()
		exec(s)
		cost := time.Since(start)
		fmt.Println("time: " + cost.String())
	}
}

// call engine
func exec(exp string) {
	// input text -> []token
	toks, err := engine.Parse(exp)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	// []token -> AST Tree
	ast := engine.NewAST(toks, exp)
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	// AST builder
	ar := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	fmt.Printf("ExprAST: %+v\n", ar)
	// catch runtime errors
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR: ", e)
		}
	}()
	// AST traversal -> result
	r := engine.ExprASTResult(ar)
	fmt.Printf("%s = %v\n", exp, r)

	engine.UploadResult(r.Cipher.Data)
}
