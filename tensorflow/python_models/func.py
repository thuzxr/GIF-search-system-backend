import tensorflow as tf

def cosine(seq1, seq2):
    norm1 = tf.norm(seq1 + 1e-5, axis=2, keepdims=True)
    norm2 = tf.norm(seq2 + 1e-5, axis=2, keepdims=True)
    sim = tf.matmul(seq1 / norm1, tf.transpose(seq2 / norm2, [0, 2, 1]),name='cossim')
    return sim

def cosine_bert(seq1,seq2):
    norm1 = tf.norm(seq1 + 1e-5, axis=1, keepdims=True)
    norm2 = tf.norm(seq2 + 1e-5, axis=1, keepdims=True)
    sim = tf.matmul(seq1 / norm1, tf.transpose(seq2 / norm2, [1,0]),name='cossim')
    return sim
