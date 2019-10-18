import tensorflow as tf
from python_models.func import cosine
import numpy as np

class CBOW():
    def __init__(self,wordmat):
        self.wordmat = wordmat
        # self.word2idx = word2idx
        self.make_graph()
        self.init()
        self.dump()
    def make_graph(self,scope='CBOW'):
        with tf.variable_scope(scope,reuse = tf.AUTO_REUSE):
            with tf.variable_scope('Init'):
                wordmat = tf.get_variable('wordmat',initializer=tf.constant(self.wordmat,dtype=tf.float32))
            with tf.variable_scope('placeholder'):
                input1 = tf.placeholder(tf.int32,[None,100],'input1')
                input2 = tf.placeholder(tf.int32,[None,100],'input2')
            with tf.variable_scope('Emb'):
                emb1 = tf.nn.embedding_lookup(wordmat,input1)
                emb1 = tf.reduce_mean(emb1,axis=1)
                emb2 = tf.nn.embedding_lookup(wordmat,input2)
                emb2 = tf.reduce_mean(emb2, axis=1)
            with tf.variable_scope('Sim'):
                self.sim = cosine(emb1,emb2)
    def init(self):
        self.sess = sess = tf.Session()
        sess.run(tf.global_variables_initializer())
    def dump(self):
        builder = tf.saved_model.builder.SavedModelBuilder("CBOW")
        # GOLANG note that we must tag our model so that we can retrieve it at inference-time
        builder.add_meta_graph_and_variables(self.sess, ["var"])
        builder.save()

if __name__=="__main__":
    wordmat = np.random.rand(40000,100)
    cbow = CBOW(wordmat)
