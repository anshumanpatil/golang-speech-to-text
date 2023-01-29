#!/usr/bin/env python3

import time
# from kafka import KafkaProducer
import speech_recognition as sr
import json
import os
# import pocketsphinx as ps
from pocketsphinx import Decoder

hypothesis = []
bag_of_words = []
buzz_words = ['<s>', '</s>', '[SPEECH]', '[NOISE]', '<sil>']

MODELDIR = os.path.join(os.path.dirname(__file__), "cmusphinx-en-in-8khz-5.2") # General english model
HMM_PATH = MODELDIR + "/en_in.cd_cont_5000" 
LM_PATH = MODELDIR + "/en-us.lm.bin"
DIC_PATH = MODELDIR + "/en_in.dic"
# bootstrap_servers = ['localhost:9092']
# topicName = 'speech'
# producer = KafkaProducer(bootstrap_servers = bootstrap_servers)
# producer = KafkaProducer()
def create_jsonlines(original):
    if isinstance(original, str):
        original = json.loads(original)

    return '\n'.join([json.dumps(original[outer_key], sort_keys=True) 
                      for outer_key in sorted(original.keys(), key=lambda x: int(x))])


def pro_recognize_q(audio_data):
        config = Decoder.default_config()
        # config.set_string("sampling_rate", 16000)
        # config.set_string("no_search", True)
        # config.set_string("-hmm", HMM_PATH)  # set the path of the hidden Markov model (HMM) parameter files
        config.set_string("-lm", LM_PATH)
        config.set_string("-dict", DIC_PATH)
        config.set_string("-logfn", os.devnull)  # disable logging (logging causes unwanted output in terminal)
        decoder = Decoder(config)
        raw_data = audio_data.get_raw_data(convert_rate=16000)
        decoder.start_utt()
        decoder.process_raw(raw_data, False, False)
        decoder.end_utt()
        hypothesis = decoder.hyp()
        if hypothesis is not None: return hypothesis.hypstr


def recognize_audio(audio_file, args):
    try:
        config = Decoder.default_config()
        # config.set_string("sampling_rate", 16000)
        # config.set_string("no_search", True)
        # config.set_string("-hmm", HMM_PATH)  # set the path of the hidden Markov model (HMM) parameter files
        config.set_string("-lm", LM_PATH)
        config.set_string("-dict", DIC_PATH)
        config.set_string("-logfn", os.devnull)  # disable logging (logging causes unwanted output in terminal)
        decoder = Decoder(config)
        decoder.start_utt()
        stream = open(audio_file, 'rb')
        in_speech_bf = False
        while True:
            buf = stream.read(len(stream))
            if buf:
                decoder.process_raw(buf, False, False)  # full_utt = False
                if decoder.hyp() is not None:
                    hypothesis.append(decoder.hyp().hypstr)
                    [bag_of_words.append(seg.word) for seg in decoder.seg() if seg.word not in buzz_words]
                    decoder.end_utt()
                    decoder.start_utt()
                    
            else:
                break
    except Exception as ex:
        print ('Error occurred with %s \n%s' % (audio_file, ex))


def process_audio(audio):
    try:
        dests = pro_recognize_q(audio)
        # dest = r.recognize_google(audio, show_all=True)
        # y = json.dumps(dest)
        # res = bytes(str(y), 'utf-8')
        # producer.send(topicName, res)
        # print(y)
        print(dests)
    except Exception as e:
        print("err "+str(e))

r = sr.Recognizer()
m = sr.Microphone()

def ListenToVoice():
    with m as source:
        r.adjust_for_ambient_noise(source)
        print("Keep Quite!")
        print("Speak")
        # print('\a')
        audio = r.listen(m)
        process_audio(audio)


from pocketsphinx import AudioFile

config = {
    'verbose': False,
    'audio_file': 'foo.wav',
    'hmm': HMM_PATH,
    'lm': LM_PATH,
    'dic': DIC_PATH
}

audio = AudioFile(**config)

for phrase in audio: print(phrase)
# recognize_audio("foo.wav", 4096)
# print(bag_of_words)
# print(hypothesis)
# try:
#     while True:
#         ListenToVoice()
# except KeyboardInterrupt:
#     print('interrupted!')


