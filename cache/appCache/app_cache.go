package appCache

import "SQLGuardian/cache"

/*
*

	@author: XingGao
	@date: 2023/8/13

*
*/
var appCache = cache.NewCache()

// Set 存储全局缓存
func Set(key string, value interface{}) {
	appCache.Set(key, value)
}

// Get 获取全局缓存
func Get(key string) (interface{}, bool) {
	return appCache.Get(key)
}
