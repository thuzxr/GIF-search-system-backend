# Overview
## 0.0.2
链接mysql服务器存储和简单模糊搜索

## 0.0.1
基础框架，链接oss

# Mysql
已在腾讯云配置好了mysql服务器（```49.233.71.202```）,从服务器登陆mysql可以使用```sudo mysql -uroot -p```，密码```QWEasd123_```. DATABASE使用```GIF_INFO```
增加可远程连接的新用户```GRANT ALL PRIVILEGES ON *.* TO YOURNAME@'%' IDENTIFIED BY "YOURPASSWORD"```.

# 简单使用
本地运行```main.go```,然后```127.0.0.1:8000/search?key=搜索内容```.