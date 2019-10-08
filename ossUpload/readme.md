# OSS交互模块

## 将Gif图上传至OSS

```go
ossUpload.OssUpload(gif, path)
```

gif为待上传的Gif类（仅要求name非空即可）,path为待上传gif在本地的**存储文件夹**相对路径地址

例如，上传 /home/usr/1.gif , 则
```go
path="\\home\\usr"
gif.name="1"
```

## 从OSS获取对应GIF图的临时访问链接

```go
ossUpload.OssSignLink(gif, timespan)
```

gif为待获取链接的Gif类（仅要求name非空）,timespan为gif图临时链接的有效时长（单位：秒，**格式为int64类型**）