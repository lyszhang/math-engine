/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/20 4:35 PM
 */

package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat64ToInterger(t *testing.T) {
	f := 0.832847234
	a, l := Float64ToInterger(f)
	assert.Equal(t, a, int64(832847))
	assert.Equal(t, l, int64(6))

	f = 0.8328
	a, l = Float64ToInterger(f)
	assert.Equal(t, a, int64(8328))
	assert.Equal(t, l, int64(4))

	f = float64(1)
	a, l = Float64ToInterger(f)
	assert.Equal(t, a, int64(1))
	assert.Equal(t, l, int64(0))
}

func TestFloat64ToIntergerOffset6(t *testing.T) {
	f := 0.832847234
	a, l := Float64ToIntergerOffsetDefault(f)
	assert.Equal(t, a, int64(832847))
	assert.Equal(t, l, int64(6))

	f = 0.8328
	a, l = Float64ToIntergerOffsetDefault(f)
	assert.Equal(t, a, int64(832800))
	assert.Equal(t, l, int64(6))

	f = float64(1)
	a, l = Float64ToIntergerOffsetDefault(f)
	assert.Equal(t, a, int64(1000000))
	assert.Equal(t, l, int64(6))
}
