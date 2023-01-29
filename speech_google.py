#!/usr/bin/env python3

import time
# from kafka import KafkaProducer
import speech_recognition as sr
import json

# bootstrap_servers = ['localhost:9092']
# topicName = 'speech'
# producer = KafkaProducer(bootstrap_servers = bootstrap_servers)
# producer = KafkaProducer()
def create_jsonlines(original):
    if isinstance(original, str):
        original = json.loads(original)

    return '\n'.join([json.dumps(original[outer_key], sort_keys=True) 
                      for outer_key in sorted(original.keys(), key=lambda x: int(x))])

def process_audio(audio):
    try:
        dest = r.recognize_google(audio, show_all=True)
        # y = json.dumps(dest)
        # res = bytes(str(y), 'utf-8')
        # producer.send(topicName, res)
        print(dest)
    except Exception as e:
        print("err "+str(e))

r = sr.Recognizer()
m = sr.Microphone(0)

def ListenToVoice():
    with m as source:
        r.adjust_for_ambient_noise(source)
        print("Speak")
        print('\a')
        audio = r.listen(m)
        process_audio(audio)

try:
    while True:
        ListenToVoice()
except KeyboardInterrupt:
    print('interrupted!')
