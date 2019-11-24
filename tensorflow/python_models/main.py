import tensorflow as tf
import sys
# sys.path.append("C:\\Users\\Mr Handsome\\Downloads\\Material_alt\\gif-dio\\backend\\tensorflow")
from utils import *
import models
import json
import os

flags = tf.flags
flags.DEFINE_string("embedding_file","./data/embedding/sgns.wiki.bigram-char","for normal embedding")
flags.DEFINE_string("word2idx_file","./data/word2idx.json","word2idx file, for normal embedding")
flags.DEFINE_integer("embedding_dim",300,"embedding dims, for normal embedding")
flags.DEFINE_string("BERT_DIR","./chinese_L-12_H-768_A-12","dir for bert")
flags.DEFINE_integer("bert_max_len",128,"max length for bert")
flags.DEFINE_string("mode","bert","mode of encoder, select normal or bert")
flags.DEFINE_string("gpu","0","gpu to run on")

def main(_):
    config = flags.FLAGS
    os.environ['CUDA_VISIBLE_DEVICES'] = config.gpu
    try:
        assert config.mode in ['bert','normal']
    except:
        raise ValueError("mode could only be bert or normal")
    model = models.BERT(config)
    recommend_bert("./info_old.json",config,model)
# /Users/saberrrrrrrr/go/src/backend/info_old.json
if __name__=="__main__":
    tf.app.run()