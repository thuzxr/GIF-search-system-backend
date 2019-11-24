import jieba.posseg as seg
import codecs

cnt0={}
with codecs.open("./ind_keyword.ind","r","utf-8") as f:
    str0=f.read()
    lis0=str0.split("#")
    # print(lis0)
    for i in lis0:
        ls1=seg.cut(i)
        for w in ls1:
            if w.flag in cnt0.keys():
                cnt0[w.flag]+=1
            else:
                cnt0[w.flag]=1

cnt1={}
with codecs.open("./tensorflow/data/embedding/sgns.wiki.bigram-char","r","utf-8") as f:
    str0=f.readline()
    n=int(str0.split(' ')[0])
    cnt00=0
    for i in range(0,n):
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
            cnt00=cnt00+1
            if w.flag in cnt1.keys():
                cnt1[w.flag]+=1
            else:
                cnt1[w.flag]=1
            break
        print("\r{}".format(cnt00),end='',flush=True)

print(cnt0)
print(cnt1)
lim=10

cnt01=0
for k in cnt0.keys():
    if k not in cnt1.keys():
        continue
    if cnt0[k]>lim:
        cnt01+=cnt1[k]
print(cnt01)
