/**
 * @Author: lyszhang
 * @Email: ericlyszhang@gmail.com
 * @Date: 2020/11/13 3:08 PM
 */

package common

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache *cache.Cache

func init() {
	// 默认过期时间为5min，每10min清理一次过期缓存
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}
