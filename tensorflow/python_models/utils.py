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

def recommend(path,word2idx,model,sess):
    maxtags = 0
    with open(path,'r') as f:
        gifs = json.load(f)['gifs']

    words = []
    for i,gif in tqdm(enumerate(gifs)):
        keywords = gif['keyword'].split(' ')
        maxtags = max(maxtags,len(keywords))
    for i, gif in tqdm(enumerate(gifs)):
        keywords = gif['keyword'].split(' ')
        word = []
        for keyword in keywords:
            word.append(get_word(list(jieba.cut(keyword)),word2idx,20))
        word += [[0]*20]*(maxtags-len(word))
        word = np.array(word,np.int32)
        words.append(word)

    arrs = np.array(words)
    for i,gif in tqdm(enumerate(gifs)):
        mean_sim = list(sess.run(model.mean_sim,feed_dict={model.input1:words[i],model.input2:arrs}))
        sims = [(j,sim) for j,sim in enumerate(mean_sim) if j!=i ]
        sims = sorted(sims,key= lambda x:x[1],reverse=True)
        gifs[i]['recommend'] = " ".join([str(s[0]) for s in sims][:10])


    with open('info2.json','w') as f:
        json.dump({'gifs':gifs},f)

if __name__=="__main__":
    from python_models.models import CBOW
    loader = Loader()
    word2idx,wordmat = loader.load_emb()
    model = CBOW(wordmat=wordmat)
    recommend("/Users/saberrrrrrrr/go/src/backend/info.json",word2idx,model,model.sess)


