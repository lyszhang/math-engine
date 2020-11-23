/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/13 3:40 PM
 */

package equation

import (
	"fmt"
	"github.com/dengsgo/math-engine/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	buf := "z=x+y"
	_, err := New(buf)
	assert.NoError(t, err)
}

func TestDivide(t *testing.T) {
	fmt.Println(float64(10) / float64(3))
}

func BenchmarkExecEquation(b *testing.B) {
	engine.RegFunction("COMPARE", 2, engine.Compare)
	es := []Equation{{"COMPARE(x, y)", "z"}}
	for i := 0; i < b.N; i++ {
		ExecEquation(es)
	}
}
