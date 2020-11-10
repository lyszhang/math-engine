/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/10 11:23 AM
 */

package calculate

import (
	"fmt"
	"github.com/dengsgo/math-engine/engine"
	"github.com/dengsgo/math-engine/source"
)

// call engine
func Exec(exp string) (result int64) {
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

	switch r.Factor {
	case engine.TypePaillier:
		result, _ = source.UploadResult(r.Cipher.Data)
	case engine.TypeConst:
		fmt.Println("result: ", r.Number)
		result = r.Number
	}
	return
}
