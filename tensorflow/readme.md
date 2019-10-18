# TensorFlow的使用

当前实现CBOW算法，使用python训练与存图，使用go进行读图，数据预处理与图运算。
由于当前上传模块还未做完，当前模块未并入总程序，而是作为一个独立的模块运行。

# 使用

请先进行环境配置：
1. python 3.6, tensorflow 1.13
2. 按照tf官网编译tf的C库
3. go get github.com/tensorflow/tensorflow/tensorflow/go

具体使用：
1. 运行.py文件，进行存图。
2. 运行go文件，进行读图和测例运算。