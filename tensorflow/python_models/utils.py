import python_models.constant as constant
import numpy as np
from tqdm import tqdm
import json
import jieba

class Loader():
    def __init__(self):
        self.emb_path = constant.embedding_file
    def load_emb(self):
        path = self.emb_path
        emb = []
        word2idx = {}
        with open(path,'r') as f:
            f.readline()
            print('Loading embedding')
            for line in f:
                splt = line.split(' ')[:-1]
                assert len(splt)==constant.embedding_dim+1
                word = splt[0]
                vec = [float(d) for d in splt[1:]]
                emb.append(vec)
                word2idx[word]=len(word2idx)+2
            print('Loading embedding finished')
        return word2idx,np.array(emb)


def get_word(tokens,word2idx,pad_len):
    words = []
    for token in tokens:
        if token in word2idx:
            words.append(word2idx[token])
        else:
            words.append(1)
    words += [0]*(pad_len-len(words))
    return words

def get_reduce_mean(path,word2idx,model,sess):
    maps = {}
    with open(path,'r') as f:
        gifs = json.load(f)['gifs']
    for gif in gifs:
        keywords = gif['keyword'].split(' ')
        words = []
        for keyword in keywords:
            words.append(get_word(list(jieba.cut(keyword)),word2idx,20))
        words = np.array(words,np.int32)
        emb = sess.run(model.emb1,feed_dict={model.input1:words})
        maps[gif['name']]=emb.tolist()
    with open('emb.json','w') as f:
        json.dump(maps,f)

if __name__=="__main__":
    from python_models.models import CBOW
    loader = Loader()
    word2idx,wordmat = loader.load_emb()
    model = CBOW(wordmat=wordmat)
    get_reduce_mean("/Users/saberrrrrrrr/go/src/backend/info.json",word2idx,model,model.sess)


