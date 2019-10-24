import python_models.constant as constant
import numpy as np
from tqdm import tqdm
import codecs
class Loader():
    def __init__(self):
        self.emb_path = constant.embedding_file
    def load_emb(self):
        path = self.emb_path
        emb = []
        word2idx = {}
        with codecs.open(path,'r',"utf-8") as f:
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