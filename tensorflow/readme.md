# 说明

基于CBOW的gif图推荐系统。

# 使用

#### 请先进行环境配置：
1. python 3.6, tensorflow 1.10~1.13均可
2. 按照tf官网编译tf的C库
3. go get github.com/tensorflow/tensorflow/tensorflow/go

#### 需要的go包：
gse——用于分词
tensorflow——用于搭建模型

#### 具体使用：
1. 运行main.py文件，进行存图和计算word2idx。注意事先下载embedding文件，推荐https://github.com/Embedding/Chinese-Word-Vectors，存放位置和具体参数见constant.py。
2. cbow.go提供接口：
   1. Init函数：用来读图和word2idx，并对gif图进行预处理，先运行这个。注意要提供model的路径和word2idx的路径（见constant.py的dump路径）。
   2. Recommend函数：用来进行推荐，输入一个“喜欢”的gif图，和一个gif图集合，返回推荐的gif图。（前十个）