import tensorflow as tf
from python_models.func import cosine
import numpy as np
import python_models.constant as constant

class CBOW():
    def __init__(self,wordmat):
        self.wordmat = wordmat
        self.make_graph()
        self.init()
    def make_graph(self,scope='CBOW'):
        with tf.variable_scope(scope,reuse = tf.AUTO_REUSE):
            with tf.variable_scope('Init'):
                wordmat = tf.concat([
                    tf.zeros([1, constant.embedding_dim], tf.float32, name='pad'),
                    tf.get_variable('unk', shape=[1, constant.embedding_dim], dtype=tf.float32,
                                    initializer=tf.contrib.layers.xavier_initializer())
                    , tf.get_variable('wordmat', initializer=tf.constant(self.wordmat, dtype=tf.float32))], axis=0)
            with tf.variable_scope('placeholder'):
                self.input1 = input1 = tf.placeholder(tf.int32,[None,20],'input1')
                self.input2 = input2 = tf.placeholder(tf.int32,[None,None,20],'input2')
                len1 = tf.reduce_sum(tf.cast(tf.cast(input1,tf.bool),tf.float32),axis=1,keep_dims=True)
                len2 = tf.reduce_sum(tf.cast(tf.cast(input2, tf.bool),tf.float32), axis=2,keep_dims=True)
                numtag1 = tf.reduce_sum(tf.cast(tf.cast(tf.squeeze(len1), tf.bool),tf.float32))
                tag_mask = tf.tile(tf.expand_dims(tf.cast(tf.cast(tf.squeeze(len2), tf.bool),tf.float32),axis=1),[1,tf.shape(self.input1)[0],1])
            with tf.variable_scope('Emb'):
                emb1 = tf.nn.embedding_lookup(wordmat,input1)
                self.emb1 = emb1 = tf.div(tf.reduce_sum(emb1,axis=1),len1+1e-8,name='emb1')
                emb1 = tf.tile(tf.expand_dims(emb1,axis=0),[tf.shape(self.input2)[0],1,1])
                emb2 = tf.nn.embedding_lookup(wordmat,input2)
                self.emb2 = emb2 = tf.div(tf.reduce_sum(emb2, axis=2),len2+1e-8,name='emb2')
            with tf.variable_scope('Sim'):
                self.sim = cosine(emb1,emb2)-2*(1-tag_mask)
                self.mean_sim = tf.div(tf.reduce_sum(tf.reduce_max(self.sim,axis=2),axis=1),numtag1,name='mean_sim')
    def init(self):
        self.sess = sess = tf.Session()
        sess.run(tf.global_variables_initializer())
    def dump(self):
        builder = tf.saved_model.builder.SavedModelBuilder("CBOW")
        builder.add_meta_graph_and_variables(self.sess, ["var"])
        builder.save()
