## Cache

存储在服务器端的Cache。Cache分别存储name,title

存储目录下的文件名对应base64编码的keyWord

第一次使用需要初始化目录，创建cache_name及cache_title

```go
cache.OfflineCacheInit()
//默认目录创建位置为utils.CACHE_DIR
```

提供的主要方法如下

```go
//返回keyword到对应的搜索缓存的映射
var m map[string][]utils.Gifs
m=cache.OfflineCacheReload()
//添加特定关键字及对应的搜索结果
var searchResult []utils.Gifs
var keyword string
cache.OfflineCacheAppend(keyword,searchResult)
//查询特定keyword的name list
var ans []string
ans=cache.OfflineCacheQuery(keyword)
//删除对应keyword对应的cache
cache.OfflineCacheDelete(keyword)
//清空cache
cache.OfflineCacheClear()
//TBD:客户端缓存的显示
```