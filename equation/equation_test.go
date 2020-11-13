/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/13 3:40 PM
 */

package equation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	buf := "z=x+y"
	_, err := New(buf)
	assert.NoError(t, err)
}
