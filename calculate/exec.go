/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/10 11:23 AM
 */

package calculate

import (
	"fmt"
	"github.com/dengsgo/math-engine/common"
	"github.com/dengsgo/math-engine/engine"
	"github.com/dengsgo/math-engine/entry"
)

// call engine
func Exec(exp string) (result *common.ArithmeticFactor) {
	entry.Reset()

	// input text -> []token
	toks, err := engine.Parse(exp)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	// []token -> AST Tree
	ast := engine.NewAST(toks, exp)
	if ast.Err != nil {
		entry.Append("ERROR: " + ast.Err.Error())
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	// AST builder
	ar := ast.ParseExpression()
	if ast.Err != nil {
		entry.Append("ERROR: " + ast.Err.Error())
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	fmt.Printf("Raw ExprAST: %+v\n", ar)
	entry.Append(fmt.Sprintf("ExprAST: %+v\n", ar))

	// catch runtime errors
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR: ", e)
		}
	}()
	// AST traversal -> result
	result = engine.ExprASTResult(ar)
	entry.Append(fmt.Sprintf("最终结果: %s\n", result.String()))
	return
}
