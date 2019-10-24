import json
import codecs
import numpy as np

dic={}
with codecs.open("./tensorflow/data/embedding/sgns.wiki.bigram-char","r","utf-8") as f:
    str0=f.readline()
    n=int(str0.split(' ')[0])
    for i in range(0,n):
        str0=f.readline()
        lis=[]
        lis0=str0.split(' ')
        if len(lis0)<301:
            continue
        if not ('\u4e00'<=lis0[0]<='\u9fff'):
            continue
        for j in range(1,301):
            lis.append(int(133+79*float(lis0[j])))
        dic[lis0[0]]=lis
        print("\r{}".format(i),end='',flush=True)
with codecs.open("./embVectors.json","w","utf-8") as f:
    f.write(json.dumps(dic))