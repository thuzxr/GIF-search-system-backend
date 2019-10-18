## Search

基本的离线搜索模块。

Index.go提供离线索引的生成及读取模块
```go
db:=DB_connect()
//生成json格式的index,对应IndexParse()
search.IndexInit(db)
//生成.ind格式的index,index文件内容为#号分割的name title keyword,对应NameIndex() TitleIndex() KeywordIndex()
search.FastIndexInit(db)

var GifList []utils.GIfs
//从json格式的索引中读取gif列表
GifList=search.IndexParse()

//从.ind格式的索引中读取name title keyword的列表
names:=search.NameIndex()
titles:=search.TitleIndex()
keywords:=search.KeywordIndex()

```

Search.go提供基于离线索引的简单搜索

```go
res:=SimpleSearch(keyword, names, titles, keywords)
//keyword为待查询的关键字,keywords为读取的keyword列表
```