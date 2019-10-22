import python_models.constant as constant
import numpy as np

class Loader():
    def __init__(self):
        self.emb_path = constant.embedding_file
    def load_emb(self):
        path = self.emb_path
        emb = []
        word2idx = {}
        with open(path,'r') as f:
            for line in f:
                splt = line.split()
                word = splt[0]
                vec = [float(d) for d in splt[1:]]
                emb.append(vec)
                word2idx[word]=len(word2idx)+2
        return word2idx,np.array(emb)