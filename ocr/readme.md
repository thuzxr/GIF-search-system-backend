# Gif元信息读取及OCR交互模块

## Gif元信息读取

通过读取info.json内的对应信息，返回当前Gif的信息

后续开发中将替换为同Mysql数据库的交互

```go
var gifs []Gifs
gifs=ocr.JsonParse(localPath)
```

其中localPath为info.json存储路径

## OCR交互

通过向OCR服务器发送GET请求实现交互，请求详情同官方文档一致

```go
var tags []string
tags=ocr.Ocr(gif)
```