
import numpy as np
from tqdm import tqdm
import codecs
import json
import jieba
import bert.tokenization as tokenization
import models

def load_emb(path):
    emb = []
    word2idx = {}
    with codecs.open(path,'r',"utf-8") as f:
        f.readline()
        print('Loading embedding')
        for line in f:
            splt = line.split(' ')[:-1]
            assert len(splt)==config.embedding_dim+1
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


def convert_single_example( max_seq_length,
                           tokenizer,text_a):
  tokens_a = tokenizer.tokenize(text_a)
#   print(text_a)
#   print(tokens_a)
#   print('\n')
#   input()
  tokens_b = None
  tokens = []
  segment_ids = []
  tokens.append("[CLS]")
  segment_ids.append(0)
  for token in tokens_a:
    tokens.append(token)
    segment_ids.append(0)
  tokens.append("[SEP]")
  segment_ids.append(0)
  if tokens_b:
    for token in tokens_b:
      tokens.append(token)
      segment_ids.append(1)
    tokens.append("[SEP]")
    segment_ids.append(1)
  input_ids = tokenizer.convert_tokens_to_ids(tokens)# 将中文转换成ids

  input_mask = [1] * len(input_ids)
  ss = len(input_ids)
  while len(input_ids) < max_seq_length:
    input_ids.append(0)
    input_mask.append(0)
    segment_ids.append(0)
  return input_ids,input_mask,segment_ids,ss

def recommend_bert(path,config,model):
    with open(path,'r') as f:
        gifs = json.load(f)['gifs']
    tokenizer = tokenization.FullTokenizer(config.BERT_DIR+"/vocab.txt")
    input_ids,input_masks,segment_ids,numtags,bert_lens = [],[],[],[],[]
    maxtags = 0
    for i,gif in tqdm(enumerate(gifs)):
        keywords = gif['keyword'].split(' ')
        maxtags = max(maxtags,len(keywords))
    print(maxtags)
    max_d = 0
    for i, gif in tqdm(enumerate(gifs)):
        input_id,input_mask,segment_id,bert_len = [],[],[],[]
        keywords = gif['keyword'].split(' ')
        numtag = len(keywords)
        for keyword in keywords:
            a,b,c,d = convert_single_example(config.bert_max_len,tokenizer,keyword)
            max_d = max(max_d,d)
            input_id.append(a)
            input_mask.append(b)
            segment_id.append(c)
            bert_len.append(d)
        input_id+= [[0]*config.bert_max_len]*(maxtags-len(input_id))
        input_mask+= [[0]*config.bert_max_len]*(maxtags-len(input_mask))
        segment_id+= [[0]*config.bert_max_len]*(maxtags-len(segment_id))
        bert_len+= [0]*(maxtags-len(bert_len))
        input_ids.append(np.array(input_id))
        input_masks.append(np.array(input_mask))
        segment_ids.append(np.array(segment_id))
        bert_lens.append(np.array(bert_len))
        numtags.append(numtag)
    
    sess = model.sess
    i=-1
    for gif in tqdm(gifs):
        i+=1
        mean_sim = []
        for j,gif_judge in enumerate(gifs):
            mean_sim.append(sess.run(model.mean_sim,
            feed_dict={
                model.bert_input_ids: input_ids[i],
                model.bert_input_mask: input_masks[i],
                model.bert_segment_ids: segment_ids[i],
                model.bert_len: bert_lens[i],
                model.numtags1: np.array([numtags[i]]),

                model.bert_input_ids2: input_ids[j],
                model.bert_input_mask2: input_masks[j],
                model.bert_segment_ids2: segment_ids[j],
                model.bert_len2: bert_lens[j],
                model.numtags2: np.array([numtags[j]])
                }))
        sims = [(j,sim) for j,sim in enumerate(mean_sim) if j!=i ]
        sims = sorted(sims,key= lambda x:x[1],reverse=True)
        gifs[i]['recommend'] = " ".join([str(s[0]) for s in sims][:10])


    with open('info_bert.json','w') as f:
        json.dump({'gifs':gifs},f)



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

# if __name__=="__main__":
    # from python_models.models import CBOW
    # loader = Loader()
    # word2idx,wordmat = loader.load_emb()
    # model = CBOW(wordmat=wordmat)
    # recommend("/Users/saberrrrrrrr/go/src/backend/info.json",word2idx,model,model.sess)


