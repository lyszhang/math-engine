/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/13 3:19 PM
 */

package common

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	_, flag := Cache.Get("a")
	assert.False(t, flag, errors.New("should be false"))

	Cache.Set("a", "value of a", cache.DefaultExpiration)

	value, _ := Cache.Get("a")
	assert.Equal(t, "value of a", value.(string), errors.New("content missed"))

	Cache.Flush()
	assert.False(t, flag, errors.New("should be false"))
}
