import json
import codecs
import numpy as np
import jieba.posseg as seg

dic={}
cnt0={}
with codecs.open("./tensorflow/data/embedding/sgns.wiki.bigram-char","r","utf-8") as f:
    str0=f.readline()
    n=int(str0.split(' ')[0])
    cnt=0
    for i in range(0,30000):
        str0=f.readline()
        lis=[]
        lis0=str0.split(' ')
        if len(lis0)<301:
            continue
        if not ('\u4e00'<=lis0[0]<='\u9fff'):
            # print("jumped")
            continue
        ls1=seg.cut(lis0[0])
        for w in ls1:
            if w.flag in cnt0.keys():
                cnt0[w.flag]+=1
            else:
                cnt0[w.flag]=1
        for j in range(1,301):
            lis.append(int(133+79*float(lis0[j])))
        dic[lis0[0]]=lis
        print("\r{}".format(cnt),end='',flush=True)
        cnt=cnt+1
print(cnt0)
# with codecs.open("./embVectors.json","w","utf-8") as f:
    # f.write(json.dumps(dic))