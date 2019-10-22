import tensorflow as tf
from python_models.utils import *
import python_models.constant as constant
import python_models.models as models
import json

def main(_):
    loader = Loader()
    word2idx,wordmat = loader.load_emb()
    model = models.CBOW(wordmat)
    model.dump()
    with open(constant.word2idx_file,'w') as f:
        json.dump(word2idx,f)

if __name__=="__main__":
    tf.app.run()