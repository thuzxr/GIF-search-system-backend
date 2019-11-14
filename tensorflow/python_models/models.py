import tensorflow as tf
from func import cosine,cosine_bert
import numpy as np
from bert import modeling

class CBOW():
    def __init__(self,wordmat,config):
        self.config = config
        self.wordmat = wordmat
        self.model()
        self.session()
    def model(self,scope='CBOW'):
        config = self.config
        with tf.variable_scope(scope,reuse = tf.AUTO_REUSE):
            with tf.variable_scope('Init'):
                wordmat = tf.concat([
                    tf.zeros([1, config.embedding_dim], tf.float32, name='pad'),
                    tf.get_variable('unk', shape=[1, config.embedding_dim], dtype=tf.float32,
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
    def session(self):
        self.sess = sess = tf.Session()
        sess.run(tf.global_variables_initializer())
    def dump(self):
        builder = tf.saved_model.builder.SavedModelBuilder("CBOW")
        builder.add_meta_graph_and_variables(self.sess, ["var"])
        builder.save()

class BERT():
    def __init__(self,config):
        self.config = config
        self.model()
        self.session()
    def model(self):
        config = self.config
        bert_config = modeling.BertConfig.from_json_file(config.BERT_DIR+'/bert_config.json')

        self.bert_input_ids = tf.placeholder (shape=[None,config.bert_max_len],dtype=tf.int32,name="input_ids")
        self.bert_input_mask = tf.placeholder (shape=[None,config.bert_max_len],dtype=tf.int32,name="input_masks")
        self.bert_segment_ids = tf.placeholder (shape=[None,config.bert_max_len],dtype=tf.int32,name="segment_ids")
        self.bert_len = tf.placeholder(shape=[None],dtype=tf.int32,name='bert_len')
        bert_model = modeling.BertModel(
            config=bert_config,
            is_training= False,
            input_ids=self.bert_input_ids,
            input_mask=self.bert_input_mask,
            token_type_ids=self.bert_segment_ids,
            use_one_hot_embeddings=False
        )

        self.bert_input_ids2 = tf.placeholder(shape=[None,config.bert_max_len],dtype=tf.int32,name="input_ids2")
        self.bert_input_mask2 = tf.placeholder (shape=[None,config.bert_max_len],dtype=tf.int32,name="input_masks2")
        self.bert_segment_ids2 = tf.placeholder (shape=[None,config.bert_max_len],dtype=tf.int32,name="segment_ids2")
        self.bert_len2 = tf.placeholder(shape=[None],dtype=tf.int32,name='bert_len2')
        bert_model2 = modeling.BertModel(
            config=bert_config,
            is_training= False,
            input_ids=self.bert_input_ids2,
            input_mask=self.bert_input_mask2,
            token_type_ids=self.bert_segment_ids2,
            use_one_hot_embeddings=False
        )

        self.numtags1 = tf.placeholder(shape=[1],dtype=tf.int32,name="numtag1")
        self.numtags2 = tf.placeholder(shape=[1],dtype=tf.int32,name="numtag2")

        tag_mask1 = tf.expand_dims(tf.reshape(tf.sequence_mask(self.numtags1,tf.shape(self.bert_input_ids)[0],dtype=tf.float32),[-1]),axis=1)
        tag_mask2 = tf.expand_dims(tf.reshape(tf.sequence_mask(self.numtags2,tf.shape(self.bert_input_ids)[0],dtype=tf.float32),[-1]),axis=1)

        tvars = tf.trainable_variables()

        init_checkpoint = config.BERT_DIR+'/bert_model.ckpt'
        (assignment_map, initialized_variable_names) = modeling.get_assignment_map_from_checkpoint(tvars,init_checkpoint)
        tf.train.init_from_checkpoint(init_checkpoint, assignment_map)

        output = bert_model.get_pooled_output()*tag_mask1
        output2 = bert_model2.get_pooled_output()*tag_mask2

        sim = cosine_bert(output,output2)
        self.mean_sim = tf.div(tf.reduce_sum(tf.reduce_max(sim,axis=1),axis=0),tf.cast(self.numtags1[0],tf.float32),name='mean_sim')
    def session(self):
        self.sess = sess = tf.Session()
        sess.run(tf.global_variables_initializer())
